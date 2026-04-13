package main

import (
	"log"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/auth"
	"travel-backend/internal/booking"
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

	// Auto migrate: tự động tạo/cập nhật cấu trúc bảng
	database.DB.AutoMigrate(&domain.User{}, &domain.Tour{}, &domain.Booking{})

	// Seed dữ liệu mẫu vào database nếu chưa có data
	seedData()

	// Khởi tạo Gin router
	r := gin.Default()

	// Cấu hình CORS để cho phép React (frontend) gọi API từ domain khác
	r.Use(cors.Default())

	// ===== ĐĂNG KÝ ROUTES =====
	// Nhóm tất cả API dưới prefix /v1
	v1 := r.Group("/v1")
	{
		// Tour routes
		v1.GET("/api/tours", tour.GetToursHandler)
		v1.GET("/api/tours/domestic", tour.GetDomesticToursHandler)
		v1.GET("/api/tours/international", tour.GetInternationalToursHandler)
		v1.POST("/api/bookings", booking.CreateBookingHandler)
		v1.GET("/api/users/:id/bookings", booking.GetUserBookingsHandler)
		v1.PUT("/api/users/:id/bookings/:bookingId/cancel", booking.CancelBookingHandler)

		// User routes
		v1.POST("/api/login", auth.LoginHandler)
		v1.POST("/api/register", auth.RegisterHandler)
		v1.POST("/api/otp/send", auth.SendOTPHandler)
		v1.POST("/api/otp/verify", auth.VerifyOTPHandler)
		v1.POST("/api/password/forgot", auth.ForgotPasswordHandler)
		v1.PUT("/api/users/:id", auth.UpdateUserHandler)
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
			{Name: "Tour Đà Nẵng - Hội An", Type: "domestic", Price: "2.000.000đ", Location: "Đà Nẵng", Country: "Việt Nam", Duration: "3 ngày 2 đêm", Description: "Khám phá phố cổ Hội An và biển Mỹ Khê."},
			{Name: "Tour Hà Nội - Sa Pa", Type: "domestic", Price: "3.500.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Trải nghiệm khí hậu vùng cao và bản làng Tây Bắc."},
			{Name: "Tour Phú Quốc", Type: "domestic", Price: "5.000.000đ", Location: "Phú Quốc", Country: "Việt Nam", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo và thưởng thức hải sản địa phương."},
			{Name: "Tour Nha Trang - Đà Lạt", Type: "domestic", Price: "4.800.000đ", Location: "Nha Trang", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Kết hợp nghỉ biển Nha Trang và khí hậu mát mẻ Đà Lạt."},
			{Name: "Tour Bangkok - Pattaya", Type: "international", Price: "8.200.000đ", Location: "Bangkok", Country: "Thái Lan", Duration: "5 ngày 4 đêm", Description: "Lộ trình quốc tế phù hợp gia đình và nhóm bạn."},
			{Name: "Tour Seoul Mùa Hoa", Type: "international", Price: "12.500.000đ", Location: "Seoul", Country: "Hàn Quốc", Duration: "6 ngày 5 đêm", Description: "Tham quan cung điện, phố mua sắm và ẩm thực Hàn."},
			{Name: "Tour Tokyo - Núi Phú Sĩ", Type: "international", Price: "15.900.000đ", Location: "Tokyo", Country: "Nhật Bản", Duration: "6 ngày 5 đêm", Description: "Khám phá Tokyo hiện đại và trải nghiệm văn hóa Nhật Bản."},
			{Name: "Tour Paris - Lyon", Type: "international", Price: "18.500.000đ", Location: "Paris", Country: "Pháp", Duration: "7 ngày 6 đêm", Description: "Hành trình châu Âu với điểm nhấn ẩm thực và kiến trúc cổ điển."},
			{Name: "Tour Singapore - Sentosa", Type: "international", Price: "9.600.000đ", Location: "Singapore", Country: "Singapore", Duration: "4 ngày 3 đêm", Description: "Khám phá đảo quốc sư tử với lịch trình hiện đại và thân thiện gia đình."},
			{Name: "Tour Bali - Ubud", Type: "international", Price: "10.800.000đ", Location: "Bali", Country: "Indonesia", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo, check-in ruộng bậc thang và đền cổ Bali."},
			{Name: "Tour Sydney - Melbourne", Type: "international", Price: "21.900.000đ", Location: "Sydney", Country: "Úc", Duration: "7 ngày 6 đêm", Description: "Hành trình nước Úc qua hai thành phố biểu tượng."},
			{Name: "Tour Dubai - Abu Dhabi", Type: "international", Price: "19.500.000đ", Location: "Dubai", Country: "UAE", Duration: "6 ngày 5 đêm", Description: "Trải nghiệm thành phố xa hoa và văn hóa Trung Đông đặc sắc."},
			{Name: "Combo Visa + Vé Máy Bay", Type: "service", Price: "1.800.000đ", Location: "Hồ Chí Minh", Country: "Việt Nam", Duration: "2 ngày", Description: "Dịch vụ làm visa nhanh và hỗ trợ đặt vé trọn gói."},
			{Name: "Đưa Đón Sân Bay Cao Cấp", Type: "service", Price: "900.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "Trong ngày", Description: "Đưa đón đúng giờ với xe riêng và tài xế kinh nghiệm."},
		}
		database.DB.Create(&tours)
		log.Println("✅ Đã seed dữ liệu Tour mẫu")
	}
}
