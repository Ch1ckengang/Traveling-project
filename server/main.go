package main

import (
	"log"
	"net/http"
	"travel-backend/database"
	"travel-backend/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// LoginRequest - Cấu trúc dữ liệu nhận từ client khi đăng nhập
// Bao gồm email và password, cả 2 đều bắt buộc (required)
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest - Cấu trúc dữ liệu nhận từ client khi đăng ký tài khoản mới
// Bao gồm tên, email và password, tất cả đều bắt buộc
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest - Cấu trúc dữ liệu nhận từ client khi cập nhật thông tin cá nhân
// Các trường đều optional, chỉ cập nhật những trường được gửi lên
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // omitempty: không bắt buộc
}

// AuthResponse - Cấu trúc dữ liệu trả về cho các API liên quan đến authentication
// Success: trạng thái thành công/thất bại
// Message: thông báo lỗi hoặc thành công
// User: thông tin user (chỉ trả về khi thành công)
type AuthResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	User    *models.User `json:"user,omitempty"`
}

// main - Hàm chính khởi tạo và chạy server
// 1. Load biến môi trường từ file .env
// 2. Kết nối database MySQL
// 3. Tự động tạo/cập nhật bảng (migration)
// 4. Seed dữ liệu mẫu nếu database trống
// 5. Khởi tạo Gin router và định nghĩa các API endpoints
// 6. Chạy server trên port 8080
func main() {
	// Load file .env để đọc thông tin database (DB_USER, DB_PASSWORD, etc.)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Không tìm thấy file .env, sử dụng biến môi trường hệ thống")
	}

	// Kết nối đến MySQL database
	database.Connect()

	// Auto migrate: tự động tạo/cập nhật cấu trúc bảng users và tours
	database.DB.AutoMigrate(&models.User{}, &models.Tour{})

	// Seed dữ liệu mẫu vào database nếu chưa có data
	seedData()

	// Khởi tạo Gin router (framework web)
	r := gin.Default()

	// Cấu hình CORS để cho phép React (frontend) gọi API từ domain khác
	r.Use(cors.Default())

	// ===== API ENDPOINTS =====
	
	// GET /api/tours - Lấy danh sách tất cả các tour
	// Response: Array của Tour objects
	r.GET("/api/tours", func(c *gin.Context) {
		var tours []models.Tour
		// Lấy tất cả tours từ database
		database.DB.Find(&tours)
		// Trả về JSON với status 200
		c.JSON(200, tours)
	})

	// POST /api/login - Xử lý đăng nhập
	// Request body: {email: string, password: string}
	// Response: {success: bool, message: string, user: User}
	// Luồng: Nhận data -> Validate -> Tìm user -> Kiểm tra password -> Trả kết quả
	r.POST("/api/login", func(c *gin.Context) {
		var loginReq LoginRequest

		// Bind JSON từ request body vào struct LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, AuthResponse{
				Success: false,
				Message: "Dữ liệu không hợp lệ",
			})
			return
		}

		// Tìm user trong database với email và password khớp
		var user models.User
		result := database.DB.Where("email = ? AND password = ?", loginReq.Email, loginReq.Password).First(&user)

		// Nếu không tìm thấy -> email hoặc password sai
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, AuthResponse{
				Success: false,
				Message: "Email hoặc mật khẩu không đúng",
			})
			return
		}

		// Đăng nhập thành công, trả về user (password tự động bị ẩn bởi json:"-" tag)
		c.JSON(http.StatusOK, AuthResponse{
			Success: true,
			Message: "Đăng nhập thành công",
			User:    &user,
		})
	})

	// POST /api/register - Đăng ký tài khoản mới
	// Request body: {name: string, email: string, password: string}
	// Response: {success: bool, message: string, user: User}
	// Luồng: Nhận data -> Validate -> Kiểm tra email trùng -> Tạo user -> Lưu DB
	r.POST("/api/register", func(c *gin.Context) {
		var registerReq RegisterRequest

		// Bind và validate JSON request
		if err := c.ShouldBindJSON(&registerReq); err != nil {
			c.JSON(http.StatusBadRequest, AuthResponse{
				Success: false,
				Message: "Dữ liệu không hợp lệ",
			})
			return
		}

		// Kiểm tra email đã tồn tại trong database chưa
		var existingUser models.User
		if err := database.DB.Where("email = ?", registerReq.Email).First(&existingUser).Error; err == nil {
			// Nếu tìm thấy (err == nil) -> email đã tồn tại
			c.JSON(http.StatusConflict, AuthResponse{
				Success: false,
				Message: "Email đã được đăng ký",
			})
			return
		}

		// Tạo đối tượng user mới
		newUser := models.User{
			Name:     registerReq.Name,
			Email:    registerReq.Email,
			Password: registerReq.Password,
		}

		// Lưu user mới vào database
		if err := database.DB.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, AuthResponse{
				Success: false,
				Message: "Không thể tạo tài khoản",
			})
			return
		}

		// Đăng ký thành công, trả về thông tin user
		c.JSON(http.StatusOK, AuthResponse{
			Success: true,
			Message: "Đăng ký thành công",
			User:    &newUser,
		})
	})

	// PUT /api/users/:id - Cập nhật thông tin cá nhân
	// URL params: id (user ID)
	// Request body: {name?: string, email?: string, password?: string}
	// Response: {success: bool, message: string, user: User}
	// Luồng: Lấy ID -> Tìm user -> Validate email trùng -> Update fields -> Lưu DB
	// Chức năng quan trọng: Kiểm tra email không bị trùng với user khác
	r.PUT("/api/users/:id", func(c *gin.Context) {
		// Lấy user ID từ URL parameter
		userID := c.Param("id")
		var updateReq UpdateUserRequest

		// Bind JSON request
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, AuthResponse{
				Success: false,
				Message: "Dữ liệu không hợp lệ",
			})
			return
		}

		// Tìm user hiện tại trong database theo ID
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, AuthResponse{
				Success: false,
				Message: "Không tìm thấy người dùng",
			})
			return
		}

		// KIỂM TRA TRÙNG LẶP EMAIL
		// Nếu user muốn đổi email (email mới khác email cũ)
		if updateReq.Email != "" && updateReq.Email != user.Email {
			var existingUser models.User
			// Tìm xem có user nào khác (id khác) đang dùng email này không
			if err := database.DB.Where("email = ? AND id != ?", updateReq.Email, userID).First(&existingUser).Error; err == nil {
				// Nếu tìm thấy -> email đã được dùng bởi user khác
				c.JSON(http.StatusConflict, AuthResponse{
					Success: false,
					Message: "Email đã được sử dụng bởi tài khoản khác",
				})
				return
			}
			// Email không trùng, cập nhật email mới
			user.Email = updateReq.Email
		}

		// Cập nhật tên nếu được gửi lên
		if updateReq.Name != "" {
			user.Name = updateReq.Name
		}

		// Cập nhật mật khẩu nếu được gửi lên
		if updateReq.Password != "" {
			user.Password = updateReq.Password
		}

		// Lưu các thay đổi vào database
		if err := database.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, AuthResponse{
				Success: false,
				Message: "Không thể cập nhật thông tin",
			})
			return
		}

		// Cập nhật thành công
		c.JSON(http.StatusOK, AuthResponse{
			Success: true,
			Message: "Cập nhật thông tin thành công",
			User:    &user,
		})
	})

	// Khởi động server HTTP trên port 8080
	// Server sẽ lắng nghe các request tại localhost:8080
	r.Run(":8080")
}

// seedData - Hàm thêm dữ liệu mẫu vào database
// Chỉ chạy khi database còn trống (lần đầu khởi động)
// Seed: Users mẫu và Tours mẫu
func seedData() {
	// SEED USERS MẪU
	var userCount int64
	// Đếm số lượng user hiện có trong database
	database.DB.Model(&models.User{}).Count(&userCount)
	// Nếu chưa có user nào -> thêm user mẫu
	if userCount == 0 {
		users := []models.User{
			{Name: "Nguyễn Văn A", Email: "test@example.com", Password: "123456"},
			{Name: "Trần Thị B", Email: "user@example.com", Password: "123456"},
		}
		database.DB.Create(&users)
		log.Println("✅ Đã seed dữ liệu User mẫu")
	}

	// SEED TOURS MẪU
	var tourCount int64
	// Đếm số lượng tour hiện có
	database.DB.Model(&models.Tour{}).Count(&tourCount)
	// Nếu chưa có tour nào -> thêm tour mẫu
	if tourCount == 0 {
		tours := []models.Tour{
			{Name: "Tour Đà Nẵng - Hội An", Price: "2.000.000đ", Location: "Đà Nẵng", Duration: "3 ngày 2 đêm"},
			{Name: "Tour Hà Nội - Sa Pa", Price: "3.500.000đ", Location: "Hà Nội", Duration: "4 ngày 3 đêm"},
			{Name: "Tour Phú Quốc", Price: "5.000.000đ", Location: "Phú Quốc", Duration: "5 ngày 4 đêm"},
		}
		database.DB.Create(&tours)
		log.Println("✅ Đã seed dữ liệu Tour mẫu")
	}
}