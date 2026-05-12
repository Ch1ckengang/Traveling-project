package main

import (
	"log"
	"time"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/auth"
	"travel-backend/internal/booking"
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
	if err := database.DB.AutoMigrate(&domain.User{}, &domain.Tour{}, &domain.Booking{}, &domain.OTP{}); err != nil {
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

	// ===== ĐĂNG KÝ ROUTES =====
	// Nhóm tất cả API dưới prefix /v1
	v1 := r.Group("/v1")
	{
		// Tour routes
		v1.GET("/api/tours", tour.GetToursHandler)
		v1.GET("/api/tours/domestic", tour.GetDomesticToursHandler)
		v1.GET("/api/tours/international", tour.GetInternationalToursHandler)
		v1.GET("/api/tours/search", tour.SearchToursHandler)
		v1.GET("/api/tours/:id", tour.GetTourByIDHandler)

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
		}
	}

	// Khởi động server HTTP trên port 8080
	r.Run(":8080")
}

// seedData - Thêm dữ liệu mẫu vào database khi lần đầu khởi động
// Hàm này vẫn ở main.go vì nó là logic khởi tạo ứng dụng,
// không phải business logic của một feature cụ thể
func seedData() {
	// Seed Users mẫu nếu bảng còn trống
	var userCount int64
	database.DB.Model(&domain.User{}).Count(&userCount)
	if userCount == 0 {
		testPasswordHash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Println("⚠️ Không thể tạo hash mật khẩu mẫu, bỏ qua seed user")
			return
		}

		users := []domain.User{
			{Name: "Nguyễn Văn A", Email: "test@example.com", Password: string(testPasswordHash), IsEmailVerified: true},
			{Name: "Trần Thị B", Email: "user@example.com", Password: string(testPasswordHash), IsEmailVerified: true},
		}
		database.DB.Create(&users)
		log.Println("✅ Đã seed dữ liệu User mẫu")
	}

	// Seed Tours mẫu nếu bảng còn trống
	var tourCount int64
	database.DB.Model(&domain.Tour{}).Count(&tourCount)
	if tourCount == 0 {
		tours := []domain.Tour{
			{Name: "Tour Đà Nẵng - Hội An", Type: "domestic", Price: "2.000.000đ", Location: "Đà Nẵng", Country: "Việt Nam", Duration: "3 ngày 2 đêm", Description: "Khám phá phố cổ Hội An và biển Mỹ Khê.", Itinerary: "Ngày 1: Đà Nẵng - Sơn Trà. Ngày 2: Hội An. Ngày 3: Mỹ Khê - mua sắm.", Services: "Xe đưa đón, khách sạn, ăn sáng, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1559592413-7cec4d0cae2b?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Hà Nội - Sa Pa", Type: "domestic", Price: "3.500.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Trải nghiệm khí hậu vùng cao và bản làng Tây Bắc.", Itinerary: "Hà Nội - Sa Pa - Fansipan - bản Cát Cát - Hà Nội.", Services: "Xe giường nằm, khách sạn, vé tham quan cơ bản.", ImageURL: "https://images.unsplash.com/photo-1528127269322-539801943592?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Phú Quốc", Type: "domestic", Price: "5.000.000đ", Location: "Phú Quốc", Country: "Việt Nam", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo và thưởng thức hải sản địa phương.", Itinerary: "Đông đảo - Nam đảo - cáp treo Hòn Thơm - chợ đêm.", Services: "Resort, ăn sáng, xe đưa đón sân bay.", ImageURL: "https://images.unsplash.com/photo-1507525428034-b723cf961d3e?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Nha Trang - Đà Lạt", Type: "domestic", Price: "4.800.000đ", Location: "Nha Trang", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Kết hợp nghỉ biển Nha Trang và khí hậu mát mẻ Đà Lạt.", Itinerary: "Nha Trang - VinWonders - Đà Lạt - Langbiang - chợ đêm.", Services: "Khách sạn, xe tour, vé tham quan theo chương trình.", ImageURL: "https://images.unsplash.com/photo-1540541338287-41700207dee6?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Bangkok - Pattaya", Type: "international", Price: "8.200.000đ", Location: "Bangkok", Country: "Thái Lan", Duration: "5 ngày 4 đêm", Description: "Lộ trình quốc tế phù hợp gia đình và nhóm bạn.", Itinerary: "Bangkok - chùa Phật Vàng - Pattaya - đảo Coral - mua sắm.", Services: "Vé máy bay, khách sạn, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1508009603885-50cf7c579365?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Seoul Mùa Hoa", Type: "international", Price: "12.500.000đ", Location: "Seoul", Country: "Hàn Quốc", Duration: "6 ngày 5 đêm", Description: "Tham quan cung điện, phố mua sắm và ẩm thực Hàn.", Itinerary: "Gyeongbokgung - Namsan - Myeongdong - đảo Nami.", Services: "Vé máy bay, visa hỗ trợ, khách sạn, ăn theo chương trình.", ImageURL: "https://images.unsplash.com/photo-1538485399081-7191377e8241?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Tokyo - Núi Phú Sĩ", Type: "international", Price: "15.900.000đ", Location: "Tokyo", Country: "Nhật Bản", Duration: "6 ngày 5 đêm", Description: "Khám phá Tokyo hiện đại và trải nghiệm văn hóa Nhật Bản.", Itinerary: "Tokyo - Asakusa - Shibuya - Phú Sĩ - Gotemba.", Services: "Vé máy bay, khách sạn, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1540959733332-eab4deabeeaf?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Paris - Lyon", Type: "international", Price: "18.500.000đ", Location: "Paris", Country: "Pháp", Duration: "7 ngày 6 đêm", Description: "Hành trình châu Âu với điểm nhấn ẩm thực và kiến trúc cổ điển.", Itinerary: "Paris - Louvre - Eiffel - Lyon - phố cổ.", Services: "Khách sạn, xe tour, vé tham quan chính.", ImageURL: "https://images.unsplash.com/photo-1502602898657-3e91760cbb34?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Singapore - Sentosa", Type: "international", Price: "9.600.000đ", Location: "Singapore", Country: "Singapore", Duration: "4 ngày 3 đêm", Description: "Khám phá đảo quốc sư tử với lịch trình hiện đại và thân thiện gia đình.", Itinerary: "Marina Bay - Gardens by the Bay - Sentosa - Orchard.", Services: "Vé máy bay, khách sạn, xe đưa đón.", ImageURL: "https://images.unsplash.com/photo-1525625293386-3f8f99389edd?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Bali - Ubud", Type: "international", Price: "10.800.000đ", Location: "Bali", Country: "Indonesia", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo, check-in ruộng bậc thang và đền cổ Bali.", Itinerary: "Kuta - Ubud - Tegallalang - Tanah Lot.", Services: "Resort, ăn sáng, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1537996194471-e657df975ab4?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Sydney - Melbourne", Type: "international", Price: "21.900.000đ", Location: "Sydney", Country: "Úc", Duration: "7 ngày 6 đêm", Description: "Hành trình nước Úc qua hai thành phố biểu tượng.", Itinerary: "Sydney Opera House - Blue Mountains - Melbourne - Great Ocean Road.", Services: "Vé máy bay, khách sạn, xe tour, hướng dẫn viên.", ImageURL: "https://images.unsplash.com/photo-1506973035872-a4ec16b8e8d9?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Tour Dubai - Abu Dhabi", Type: "international", Price: "19.500.000đ", Location: "Dubai", Country: "UAE", Duration: "6 ngày 5 đêm", Description: "Trải nghiệm thành phố xa hoa và văn hóa Trung Đông đặc sắc.", Itinerary: "Dubai Mall - Burj Khalifa - sa mạc Safari - Abu Dhabi.", Services: "Vé máy bay, khách sạn, xe tour, visa hỗ trợ.", ImageURL: "https://images.unsplash.com/photo-1512453979798-5ea266f8880c?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Combo Visa + Vé Máy Bay", Type: "service", Price: "1.800.000đ", Location: "Hồ Chí Minh", Country: "Việt Nam", Duration: "2 ngày", Description: "Dịch vụ làm visa nhanh và hỗ trợ đặt vé trọn gói.", Itinerary: "Tư vấn hồ sơ - nộp hồ sơ - theo dõi kết quả - bàn giao.", Services: "Tư vấn visa, hỗ trợ form, đặt vé theo yêu cầu.", ImageURL: "https://images.unsplash.com/photo-1436491865332-7a61a109cc05?auto=format&fit=crop&w=1200&q=80"},
			{Name: "Đưa Đón Sân Bay Cao Cấp", Type: "service", Price: "900.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "Trong ngày", Description: "Đưa đón đúng giờ với xe riêng và tài xế kinh nghiệm.", Itinerary: "Xác nhận lịch - đón tại điểm hẹn - hỗ trợ hành lý - trả khách.", Services: "Xe riêng, tài xế, hỗ trợ hành lý.", ImageURL: "https://images.unsplash.com/photo-1449965408869-eaa3f722e40d?auto=format&fit=crop&w=1200&q=80"},
		}
		database.DB.Create(&tours)
		log.Println("✅ Đã seed dữ liệu Tour mẫu")
	}
}
