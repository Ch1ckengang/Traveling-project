package main

import (
	"log"
	"time"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/auth"
	"travel-backend/internal/booking"
	"travel-backend/internal/coupon"
	"travel-backend/internal/dashboard"
	"travel-backend/internal/notification"
	"travel-backend/internal/payment"
	"travel-backend/internal/review"
	"travel-backend/internal/shared"
	"travel-backend/internal/tour"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// main - Khởi tạo và chạy server
// Sau khi refactor, main.go chỉ còn 3 nhiệm vụ:
//  1. Kết nối database
//  2. Đăng ký các routes (gắn URL với controller)
//  3. Chạy server
//
// Tất cả logic đã được chuyển sang đúng tầng của nó
func main() {
	// Kết nối đến PostgreSQL database
	database.Connect()

	// Khởi tạo email service
	shared.InitEmailService()

	// Auto migrate: tự động tạo/cập nhật cấu trúc bảng
	if err := database.DB.AutoMigrate(&domain.User{}, &domain.Tour{}, &domain.Booking{}, &domain.Payment{}, &domain.PaymentAuditLog{}, &domain.OTP{}, &domain.Review{}, &domain.Coupon{}, &domain.Notification{}); err != nil {
		log.Printf("⚠️ AutoMigrate error: %v", err)
	} else {
		log.Println("✅ AutoMigrate completed successfully")
	}

	// Seed dữ liệu mẫu vào database nếu chưa có data
	seedData()

	// Khởi tạo Gin router
	r := gin.Default()

	// Cấu hình CORS để cho phép React (frontend) gọi API từ domain khác
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rate limiting cho auth endpoints (10 requests per minute)
	authRateLimiter := shared.NewRateLimiter(10, 1*time.Minute)

	// Khởi tạo VNPay payment handler
	vnpayConfig := payment.LoadVNPayConfig()
	paymentHandler := payment.NewPaymentHandler(vnpayConfig)

	// ===== ĐĂNG KÝ ROUTES =====
	// Nhóm tất cả API dưới prefix /v1
	v1 := r.Group("/v1")
	{
		// Tour routes (public)
		v1.GET("/api/tours", tour.GetToursHandler)
		v1.GET("/api/tours/domestic", tour.GetDomesticToursHandler)
		v1.GET("/api/tours/international", tour.GetInternationalToursHandler)
		v1.GET("/api/tours/search", tour.SearchToursHandler)
		v1.GET("/api/tours/:id", tour.GetTourByIDHandler)

		// Review routes (public)
		v1.GET("/api/tours/:tourId/reviews", review.GetTourReviewsHandler)

		// User routes (with rate limiting)
		authRoutes := v1.Group("/api")
		authRoutes.Use(shared.RateLimitMiddleware(authRateLimiter))
		{
			authRoutes.POST("/login", auth.LoginHandler)
			authRoutes.POST("/register", auth.RegisterHandler)
			authRoutes.POST("/otp/send", auth.SendOTPHandler)
			authRoutes.POST("/otp/verify", auth.VerifyOTPHandler)
			authRoutes.POST("/password/forgot", auth.ForgotPasswordHandler)
			authRoutes.POST("/password/reset", auth.ResetPasswordHandler)
		}

		v1.POST("/api/token/refresh", auth.RefreshTokenHandler)

		// Protected routes (yêu cầu access token)
		protected := v1.Group("/api")
		protected.Use(auth.AuthRequired())
		{
			protected.POST("/bookings", booking.CreateBookingHandler)
			protected.GET("/bookings/code/:code", booking.GetBookingByCodeHandler)
			protected.GET("/bookings/:id", booking.GetBookingByIDHandler)
			protected.GET("/users/:id/bookings", booking.GetUserBookingsHandler)
			protected.PUT("/users/:id/bookings/:bookingId/cancel", booking.CancelBookingHandler)
			protected.PUT("/users/:id", auth.UpdateUserHandler)
			protected.GET("/users/me", auth.GetMeHandler)
			protected.PUT("/users/me/password", auth.ChangePasswordHandler)
			protected.POST("/logout", auth.LogoutHandler)

			// Payment routes (authenticated)
			protected.POST("/payments/initiate", paymentHandler.InitiatePaymentHandler)
			protected.GET("/payments/status/:ref", paymentHandler.GetPaymentStatusHandler)
			protected.GET("/bookings/:id/payments", paymentHandler.GetBookingPaymentsHandler)

			// Review routes (authenticated)
			protected.POST("/reviews", review.CreateReviewHandler)
			protected.PUT("/reviews/:id", review.UpdateReviewHandler)

			// Coupon routes
			protected.POST("/coupons/validate", coupon.ValidateCouponHandler)

			// Notification routes
			protected.GET("/notifications", notification.GetNotificationsHandler)
			protected.PUT("/notifications/:id/read", notification.MarkAsReadHandler)
			protected.PUT("/notifications/read-all", notification.MarkAllAsReadHandler)
		}

		// Payment routes (public - VNPay callbacks)
		v1.GET("/api/payments/return", paymentHandler.PaymentReturnHandler)
		v1.POST("/api/payments/webhook", paymentHandler.PaymentWebhookHandler)

		// Admin routes (yêu cầu Staff hoặc Admin role)
		admin := v1.Group("/api/admin")
		admin.Use(auth.AuthRequired(), shared.StaffRequired())
		{
			// Admin Tour CRUD
			admin.GET("/tours", tour.AdminGetToursHandler)
			admin.POST("/tours", tour.AdminCreateTourHandler)
			admin.PUT("/tours/:id", tour.AdminUpdateTourHandler)
			admin.DELETE("/tours/:id", tour.AdminDeleteTourHandler)
			admin.PUT("/tours/:id/toggle", tour.AdminToggleTourHandler)

			// Admin Dashboard
			admin.GET("/dashboard/summary", dashboard.AdminGetDashboardSummaryHandler)
			admin.GET("/dashboard/revenue-chart", dashboard.AdminGetRevenueChartHandler)
			admin.GET("/dashboard/top-tours", dashboard.AdminGetTopToursHandler)

			// Admin Booking Management
			admin.GET("/bookings", booking.AdminGetBookingsHandler)
			admin.GET("/bookings/stats", booking.AdminGetBookingStatsHandler)
			admin.GET("/bookings/:code", booking.AdminGetBookingByCodeHandler)
			admin.PUT("/bookings/:code/confirm", booking.AdminConfirmBookingHandler)
			admin.PUT("/bookings/:code/cancel", booking.AdminCancelBookingHandler)

			// Admin User Management
			admin.GET("/users", auth.AdminGetUsersHandler)
			admin.PUT("/users/:id/status", auth.AdminUpdateUserStatusHandler)
			admin.PUT("/users/:id/role", auth.AdminUpdateUserRoleHandler)

			// Admin Review Management
			admin.GET("/reviews", review.AdminGetReviewsHandler)
			admin.PUT("/reviews/:id/publish", review.AdminPublishReviewHandler)
			admin.PUT("/reviews/:id/hide", review.AdminHideReviewHandler)
			admin.POST("/reviews/:id/reply", review.AdminReplyReviewHandler)

			// Admin Coupon Management
			admin.GET("/coupons", coupon.AdminGetCouponsHandler)
			admin.POST("/coupons", coupon.AdminCreateCouponHandler)
			admin.PUT("/coupons/:id", coupon.AdminUpdateCouponHandler)
			admin.DELETE("/coupons/:id", coupon.AdminDeleteCouponHandler)
		}
	}

	// Khởi động server HTTP trên port 8080
	r.Run(":8080")
}

// seedData - Thêm dữ liệu mẫu vào database khi lần đầu khởi động
func seedData() {
	seedUsers()
	seedTours()
}

func seedUsers() {
	var userCount int64
	database.DB.Model(&domain.User{}).Count(&userCount)
	if userCount > 0 {
		// Đảm bảo user cũ có role (migration cho data cũ)
		database.DB.Model(&domain.User{}).Where("role = '' OR role IS NULL").Update("role", domain.RoleCustomer)
		return
	}

	testPasswordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		log.Println("⚠️ Không thể tạo hash mật khẩu mẫu, bỏ qua seed user")
		return
	}

	users := []domain.User{
		{Name: "Admin", Email: "admin@traveling.com", Password: string(testPasswordHash), Role: domain.RoleAdmin, IsEmailVerified: true, IsActive: true},
		{Name: "Nhân viên A", Email: "staff@traveling.com", Password: string(testPasswordHash), Role: domain.RoleStaff, IsEmailVerified: true, IsActive: true},
		{Name: "Nguyễn Văn A", Email: "test@example.com", Password: string(testPasswordHash), Role: domain.RoleCustomer, IsEmailVerified: true, IsActive: true},
		{Name: "Trần Thị B", Email: "user@example.com", Password: string(testPasswordHash), Role: domain.RoleCustomer, IsEmailVerified: true, IsActive: true},
	}
	database.DB.Create(&users)
	log.Println("✅ Đã seed dữ liệu User mẫu (bao gồm admin + staff)")
}

func seedTours() {
	var tourCount int64
	database.DB.Model(&domain.Tour{}).Count(&tourCount)
	if tourCount > 0 {
		// Migration: cập nhật PriceAmount cho tours cũ nếu chưa có
		database.DB.Model(&domain.Tour{}).Where("price_amount = 0 OR price_amount IS NULL").Updates(map[string]interface{}{
			"is_active": true,
		})
		return
	}

	tours := []domain.Tour{
		{Name: "Tour Đà Nẵng - Hội An", Slug: "tour-da-nang-hoi-an", Type: "domestic", PriceAmount: 2000000, Price: "2.000.000đ", Location: "Đà Nẵng", Country: "Việt Nam", Duration: "3 ngày 2 đêm", Description: "Khám phá phố cổ Hội An và biển Mỹ Khê.", Itinerary: "Ngày 1: Đà Nẵng - Sơn Trà. Ngày 2: Hội An. Ngày 3: Mỹ Khê - mua sắm.", Services: "Xe đưa đón, khách sạn, ăn sáng, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1559592413-7cec4d0cae2b?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Hà Nội - Sa Pa", Slug: "tour-ha-noi-sa-pa", Type: "domestic", PriceAmount: 3500000, Price: "3.500.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Trải nghiệm khí hậu vùng cao và bản làng Tây Bắc.", Itinerary: "Hà Nội - Sa Pa - Fansipan - bản Cát Cát - Hà Nội.", Services: "Xe giường nằm, khách sạn, vé tham quan cơ bản.", ImageURL: "https://images.unsplash.com/photo-1528127269322-539801943592?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Phú Quốc", Slug: "tour-phu-quoc", Type: "domestic", PriceAmount: 5000000, Price: "5.000.000đ", Location: "Phú Quốc", Country: "Việt Nam", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo và thưởng thức hải sản địa phương.", Itinerary: "Đông đảo - Nam đảo - cáp treo Hòn Thơm - chợ đêm.", Services: "Resort, ăn sáng, xe đưa đón sân bay.", ImageURL: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Nha Trang - Đà Lạt", Slug: "tour-nha-trang-da-lat", Type: "domestic", PriceAmount: 4800000, Price: "4.800.000đ", Location: "Nha Trang", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Kết hợp nghỉ biển Nha Trang và khí hậu mát mẻ Đà Lạt.", Itinerary: "Nha Trang - VinWonders - Đà Lạt - Langbiang - chợ đêm.", Services: "Khách sạn, xe tour, vé tham quan theo chương trình.", ImageURL: "https://images.unsplash.com/photo-1540541338287-41700207dee6?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Bangkok - Pattaya", Slug: "tour-bangkok-pattaya", Type: "international", PriceAmount: 8200000, Price: "8.200.000đ", Location: "Bangkok", Country: "Thái Lan", Duration: "5 ngày 4 đêm", Description: "Lộ trình quốc tế phù hợp gia đình và nhóm bạn.", Itinerary: "Bangkok - chùa Phật Vàng - Pattaya - đảo Coral - mua sắm.", Services: "Vé máy bay, khách sạn, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1508009603885-50cf7c579365?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Seoul Mùa Hoa", Slug: "tour-seoul-mua-hoa", Type: "international", PriceAmount: 12500000, Price: "12.500.000đ", Location: "Seoul", Country: "Hàn Quốc", Duration: "6 ngày 5 đêm", Description: "Tham quan cung điện, phố mua sắm và ẩm thực Hàn.", Itinerary: "Gyeongbokgung - Namsan - Myeongdong - đảo Nami.", Services: "Vé máy bay, visa hỗ trợ, khách sạn, ăn theo chương trình.", ImageURL: "https://images.unsplash.com/photo-1538485399081-7191377e8241?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Tokyo - Núi Phú Sĩ", Slug: "tour-tokyo-nui-phu-si", Type: "international", PriceAmount: 15900000, Price: "15.900.000đ", Location: "Tokyo", Country: "Nhật Bản", Duration: "6 ngày 5 đêm", Description: "Khám phá Tokyo hiện đại và trải nghiệm văn hóa Nhật Bản.", Itinerary: "Tokyo - Asakusa - Shibuya - Phú Sĩ - Gotemba.", Services: "Vé máy bay, khách sạn, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Paris - Lyon", Slug: "tour-paris-lyon", Type: "international", PriceAmount: 18500000, Price: "18.500.000đ", Location: "Paris", Country: "Pháp", Duration: "7 ngày 6 đêm", Description: "Hành trình châu Âu với điểm nhấn ẩm thực và kiến trúc cổ điển.", Itinerary: "Paris - Louvre - Eiffel - Lyon - phố cổ.", Services: "Khách sạn, xe tour, vé tham quan chính.", ImageURL: "https://images.unsplash.com/photo-1502602898657-3e91760cbb34?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Singapore - Sentosa", Slug: "tour-singapore-sentosa", Type: "international", PriceAmount: 9600000, Price: "9.600.000đ", Location: "Singapore", Country: "Singapore", Duration: "4 ngày 3 đêm", Description: "Khám phá đảo quốc sư tử với lịch trình hiện đại và thân thiện gia đình.", Itinerary: "Marina Bay - Gardens by the Bay - Sentosa - Orchard.", Services: "Vé máy bay, khách sạn, xe đưa đón.", ImageURL: "https://images.unsplash.com/photo-1525625293386-3f8f99389edd?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Bali - Ubud", Slug: "tour-bali-ubud", Type: "international", PriceAmount: 10800000, Price: "10.800.000đ", Location: "Bali", Country: "Indonesia", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo, check-in ruộng bậc thang và đền cổ Bali.", Itinerary: "Kuta - Ubud - Tegallalang - Tanah Lot.", Services: "Resort, ăn sáng, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1537996194471-e657df975ab4?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Sydney - Melbourne", Slug: "tour-sydney-melbourne", Type: "international", PriceAmount: 21900000, Price: "21.900.000đ", Location: "Sydney", Country: "Úc", Duration: "7 ngày 6 đêm", Description: "Hành trình nước Úc qua hai thành phố biểu tượng.", Itinerary: "Sydney Opera House - Blue Mountains - Melbourne - Great Ocean Road.", Services: "Vé máy bay, khách sạn, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1506973035872-a4ec16b8e8d9?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Tour Dubai - Abu Dhabi", Slug: "tour-dubai-abu-dhabi", Type: "international", PriceAmount: 19500000, Price: "19.500.000đ", Location: "Dubai", Country: "UAE", Duration: "6 ngày 5 đêm", Description: "Trải nghiệm thành phố xa hoa và văn hóa Trung Đông đặc sắc.", Itinerary: "Dubai Mall - Burj Khalifa - sa mạc Safari - Abu Dhabi.", Services: "Vé máy bay, khách sạn, xe tour, visa hỗ trợ.", ImageURL: "https://images.unsplash.com/photo-1512453979798-5ea266f8880c?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Combo Visa + Vé Máy Bay", Slug: "combo-visa-ve-may-bay", Type: "service", PriceAmount: 1800000, Price: "1.800.000đ", Location: "Hồ Chí Minh", Country: "Việt Nam", Duration: "2 ngày", Description: "Dịch vụ làm visa nhanh và hỗ trợ đặt vé trọn gói.", Itinerary: "Tư vấn hồ sơ - nộp hồ sơ - theo dõi kết quả - bàn giao.", Services: "Tư vấn visa, hỗ trợ form, đặt vé theo yêu cầu.", ImageURL: "https://images.unsplash.com/photo-1436491865332-7a61a109cc05?auto=format&fit=crop&w=1200&q=80", IsActive: true},
		{Name: "Đưa Đón Sân Bay Cao Cấp", Slug: "dua-don-san-bay-cao-cap", Type: "service", PriceAmount: 900000, Price: "900.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "Trong ngày", Description: "Đưa đón đúng giờ với xe riêng và tài xế kinh nghiệm.", Itinerary: "Xác nhận lịch - đón tại điểm hẹn - hỗ trợ hành lý - trả khách.", Services: "Xe riêng, tài xế, hỗ trợ hành lý.", ImageURL: "https://images.unsplash.com/photo-1449965408869-eaa3f722e40d?auto=format&fit=crop&w=1200&q=80", IsActive: true},
	}
	database.DB.Create(&tours)
	log.Println("✅ Đã seed dữ liệu Tour mẫu (với PriceAmount)")
}
