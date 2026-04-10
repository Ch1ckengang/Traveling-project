package controllers

import (
	"errors"
	"net/http"
	"travel-backend/models"
	"travel-backend/services"

	"github.com/gin-gonic/gin"
)

// UserController - Tầng xử lý HTTP cho User
// Giống như "người phục vụ" — nhận request từ client, gọi service, trả response
// KHÔNG chứa business logic hay câu SQL
// Chỉ làm 3 việc: đọc request → gọi service → trả JSON

// Login - Xử lý POST /api/login
func Login(c *gin.Context) {
	var req models.LoginRequest

	// Bước 1: Đọc và validate dữ liệu từ request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Message: services.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	user, err := services.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(authErrorStatus(err), models.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Bước 3: Trả kết quả cho client
	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Đăng nhập thành công",
		User:    user,
	})
}

// Register - Xử lý POST /api/register
func Register(c *gin.Context) {
	var req models.RegisterRequest

	// Bước 1: Đọc và validate dữ liệu từ request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Message: services.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	newUser, err := services.Register(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(authErrorStatus(err), models.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Bước 3: Trả kết quả cho client
	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Đăng ký thành công",
		User:    newUser,
	})
}

// SendOTP - Xử lý POST /api/otp/send
func SendOTP(c *gin.Context) {
	var req models.OTPSendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Message: services.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	if err := services.SendOTPForEmail(req.Email); err != nil {
		c.JSON(authErrorStatus(err), models.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Mã xác thực đã được gửi",
	})
}

// VerifyOTP - Xử lý POST /api/otp/verify
func VerifyOTP(c *gin.Context) {
	var req models.OTPVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Message: services.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	if err := services.VerifyOTPForEmail(req.Email, req.Code); err != nil {
		c.JSON(authErrorStatus(err), models.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Xác thực OTP thành công",
	})
}

// ForgotPassword - Xử lý POST /api/password/forgot
func ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Message: services.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	if err := services.RequestPasswordReset(req.Email); err != nil {
		c.JSON(authErrorStatus(err), models.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: services.ErrPasswordResetQueued.Error(),
	})
}

func authErrorStatus(err error) int {
	switch {
	case errors.Is(err, services.ErrInvalidAuthPayload),
		errors.Is(err, services.ErrInvalidName),
		errors.Is(err, services.ErrInvalidRegisterEmail),
		errors.Is(err, services.ErrWeakPassword):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, services.ErrInvalidOTPEmail),
		errors.Is(err, services.ErrInvalidOTPCode):
		return http.StatusBadRequest
	case errors.Is(err, services.ErrEmailAlreadyRegistered):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// UpdateUser - Xử lý PUT /api/users/:id
func UpdateUser(c *gin.Context) {
	// Bước 1: Đọc ID từ URL và body từ request
	userID := c.Param("id")
	var req models.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.AuthResponse{
			Success: false,
			Message: "Dữ liệu không hợp lệ",
		})
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	updatedUser, err := services.UpdateUser(userID, req)
	if err != nil {
		// Xác định HTTP status phù hợp dựa vào thông báo lỗi
		status := http.StatusInternalServerError
		if err.Error() == "Không tìm thấy người dùng" {
			status = http.StatusNotFound
		} else if err.Error() == "Email đã được sử dụng bởi tài khoản khác" {
			status = http.StatusConflict
		}
		c.JSON(status, models.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Bước 3: Trả kết quả cho client
	c.JSON(http.StatusOK, models.AuthResponse{
		Success: true,
		Message: "Cập nhật thông tin thành công",
		User:    updatedUser,
	})
}
