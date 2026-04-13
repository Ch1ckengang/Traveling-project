package auth

import (
	"errors"
	"net/http"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// UserController - Tầng xử lý HTTP cho User
// Giống như "người phục vụ" — nhận request từ client, gọi service, trả response
// KHÔNG chứa business logic hay câu SQL
// Chỉ làm 3 việc: đọc request → gọi service → trả JSON

// LoginHandler - Xử lý POST /api/login
func LoginHandler(c *gin.Context) {
	var req domain.LoginRequest

	// Bước 1: Đọc và validate dữ liệu từ request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.AuthResponse{
			Success: false,
			Message: shared.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	user, err := Login(req.Email, req.Password)
	if err != nil {
		c.JSON(authErrorStatus(err), domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Bước 3: Trả kết quả cho client
	c.JSON(http.StatusOK, domain.AuthResponse{
		Success: true,
		Message: "Đăng nhập thành công",
		User:    user,
	})
}

// RegisterHandler - Xử lý POST /api/register
func RegisterHandler(c *gin.Context) {
	var req domain.RegisterRequest

	// Bước 1: Đọc và validate dữ liệu từ request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.AuthResponse{
			Success: false,
			Message: shared.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	newUser, err := Register(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(authErrorStatus(err), domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Bước 3: Trả kết quả cho client
	c.JSON(http.StatusOK, domain.AuthResponse{
		Success: true,
		Message: "Đăng ký thành công. Vui lòng kiểm tra email đăng ký để lấy mã OTP xác thực.",
		User:    newUser,
	})
}

// SendOTPHandler - Xử lý POST /api/otp/send
func SendOTPHandler(c *gin.Context) {
	var req domain.OTPSendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.AuthResponse{
			Success: false,
			Message: shared.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	if err := SendOTPForEmail(req.Email); err != nil {
		c.JSON(authErrorStatus(err), domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.AuthResponse{
		Success: true,
		Message: "Mã xác thực đã được gửi",
	})
}

// VerifyOTPHandler - Xử lý POST /api/otp/verify
func VerifyOTPHandler(c *gin.Context) {
	var req domain.OTPVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.AuthResponse{
			Success: false,
			Message: shared.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	if err := VerifyOTPForEmail(req.Email, req.Code); err != nil {
		c.JSON(authErrorStatus(err), domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.AuthResponse{
		Success: true,
		Message: "Xác thực OTP thành công",
	})
}

// ForgotPasswordHandler - Xử lý POST /api/password/forgot
func ForgotPasswordHandler(c *gin.Context) {
	var req domain.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.AuthResponse{
			Success: false,
			Message: shared.ErrInvalidAuthPayload.Error(),
		})
		return
	}

	if err := RequestPasswordReset(req.Email); err != nil {
		c.JSON(authErrorStatus(err), domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.AuthResponse{
		Success: true,
		Message: shared.ErrPasswordResetQueued.Error(),
	})
}

func authErrorStatus(err error) int {
	switch {
	case errors.Is(err, shared.ErrInvalidAuthPayload),
		errors.Is(err, shared.ErrInvalidName),
		errors.Is(err, shared.ErrInvalidRegisterEmail),
		errors.Is(err, shared.ErrRegisterEmailMustGmail),
		errors.Is(err, shared.ErrWeakPassword):
		return http.StatusBadRequest
	case errors.Is(err, shared.ErrInvalidCredentials):
		return http.StatusUnauthorized
	case errors.Is(err, shared.ErrEmailNotVerified):
		return http.StatusForbidden
	case errors.Is(err, shared.ErrInvalidOTPEmail),
		errors.Is(err, shared.ErrOTPEmailNotRegistered),
		errors.Is(err, shared.ErrInvalidOTPCode):
		return http.StatusBadRequest
	case errors.Is(err, shared.ErrEmailAlreadyRegistered):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// UpdateUserHandler - Xử lý PUT /api/users/:id
func UpdateUserHandler(c *gin.Context) {
	// Bước 1: Đọc ID từ URL và body từ request
	userID := c.Param("id")
	var req domain.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.AuthResponse{
			Success: false,
			Message: "Dữ liệu không hợp lệ",
		})
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	updatedUser, err := UpdateUser(userID, req)
	if err != nil {
		// Xác định HTTP status phù hợp dựa vào thông báo lỗi
		status := http.StatusInternalServerError
		if err.Error() == "Không tìm thấy người dùng" {
			status = http.StatusNotFound
		} else if err.Error() == "Email đã được sử dụng bởi tài khoản khác" {
			status = http.StatusConflict
		}
		c.JSON(status, domain.AuthResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Bước 3: Trả kết quả cho client
	c.JSON(http.StatusOK, domain.AuthResponse{
		Success: true,
		Message: "Cập nhật thông tin thành công",
		User:    updatedUser,
	})
}
