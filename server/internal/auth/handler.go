package auth

import (
	"errors"
	"net/http"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type apiError struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

type apiResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    interface{}            `json:"data"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Error   *apiError              `json:"error"`
}

func respondSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, apiResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    buildMeta(c),
		Error:   nil,
	})
}

func respondError(c *gin.Context, status int, message, code string) {
	c.JSON(status, apiResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Meta:    buildMeta(c),
		Error: &apiError{
			Code: code,
		},
	})
}

func buildMeta(c *gin.Context) map[string]interface{} {
	requestID := c.GetHeader("X-Request-ID")
	if requestID == "" {
		return nil
	}

	return map[string]interface{}{
		"request_id": requestID,
	}
}

// UserController - Tầng xử lý HTTP cho User
// Giống như "người phục vụ" — nhận request từ client, gọi service, trả response
// KHÔNG chứa business logic hay câu SQL
// Chỉ làm 3 việc: đọc request → gọi service → trả JSON

// LoginHandler - Xử lý POST /api/login
func LoginHandler(c *gin.Context) {
	var req domain.LoginRequest

	// Bước 1: Đọc và validate dữ liệu từ request body
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, shared.ErrInvalidAuthPayload.Error(), "AUTH_INVALID_PAYLOAD")
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	user, err := Login(req.Email, req.Password)
	if err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	tokens, err := issueTokenPair(user)
	if err != nil {
		respondError(c, http.StatusInternalServerError, shared.ErrTokenIssueFailed.Error(), "AUTH_TOKEN_ISSUE_FAILED")
		return
	}

	// Bước 3: Trả kết quả cho client
	respondSuccess(c, http.StatusOK, "Đăng nhập thành công", gin.H{"user": user, "tokens": tokens})
}

// RefreshTokenHandler - Xử lý POST /api/token/refresh
func RefreshTokenHandler(c *gin.Context) {
	var req domain.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, shared.ErrInvalidAuthPayload.Error(), "AUTH_INVALID_PAYLOAD")
		return
	}

	tokens, err := refreshTokenPair(req.RefreshToken)
	if err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	respondSuccess(c, http.StatusOK, "Làm mới token thành công", gin.H{"tokens": tokens})
}

// RegisterHandler - Xử lý POST /api/register
func RegisterHandler(c *gin.Context) {
	var req domain.RegisterRequest

	// Bước 1: Đọc và validate dữ liệu từ request body
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, shared.ErrInvalidAuthPayload.Error(), "AUTH_INVALID_PAYLOAD")
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	newUser, err := Register(req.Name, req.Email, req.Phone, req.Password)
	if err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	// Bước 3: Trả kết quả cho client
	respondSuccess(c, http.StatusOK, "Đăng ký thành công. Vui lòng kiểm tra email đăng ký để lấy mã OTP xác thực.", gin.H{"user": newUser})
}

// SendOTPHandler - Xử lý POST /api/otp/send
func SendOTPHandler(c *gin.Context) {
	var req domain.OTPSendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, shared.ErrInvalidAuthPayload.Error(), "AUTH_INVALID_PAYLOAD")
		return
	}

	if err := SendOTPForEmail(req.Email); err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	respondSuccess(c, http.StatusOK, "Mã xác thực đã được gửi", nil)
}

// VerifyOTPHandler - Xử lý POST /api/otp/verify
func VerifyOTPHandler(c *gin.Context) {
	var req domain.OTPVerifyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, shared.ErrInvalidAuthPayload.Error(), "AUTH_INVALID_PAYLOAD")
		return
	}

	if err := VerifyOTPForEmail(req.Email, req.Code); err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	respondSuccess(c, http.StatusOK, "Xác thực OTP thành công", nil)
}

// ForgotPasswordHandler - Xử lý POST /api/password/forgot
func ForgotPasswordHandler(c *gin.Context) {
	var req domain.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, shared.ErrInvalidAuthPayload.Error(), "AUTH_INVALID_PAYLOAD")
		return
	}

	if err := RequestPasswordReset(req.Email); err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	respondSuccess(c, http.StatusOK, shared.ErrPasswordResetQueued.Error(), nil)
}

func authErrorInfo(err error) (int, string) {
	switch {
	case errors.Is(err, shared.ErrInvalidAuthPayload),
		errors.Is(err, shared.ErrInvalidName),
		errors.Is(err, shared.ErrInvalidRegisterEmail),
		errors.Is(err, shared.ErrRegisterEmailMustGmail),
		errors.Is(err, shared.ErrWeakPassword):
		if errors.Is(err, shared.ErrInvalidName) {
			return http.StatusBadRequest, "AUTH_INVALID_NAME"
		}
		if errors.Is(err, shared.ErrWeakPassword) {
			return http.StatusBadRequest, "AUTH_WEAK_PASSWORD"
		}
		if errors.Is(err, shared.ErrInvalidRegisterEmail) || errors.Is(err, shared.ErrRegisterEmailMustGmail) {
			return http.StatusBadRequest, "AUTH_INVALID_EMAIL"
		}
		return http.StatusBadRequest, "AUTH_INVALID_PAYLOAD"
	case errors.Is(err, shared.ErrInvalidCredentials):
		return http.StatusUnauthorized, "AUTH_INVALID_CREDENTIALS"
	case errors.Is(err, shared.ErrEmailNotVerified):
		return http.StatusForbidden, "AUTH_EMAIL_NOT_VERIFIED"
	case errors.Is(err, shared.ErrInvalidOTPEmail),
		errors.Is(err, shared.ErrOTPEmailNotRegistered),
		errors.Is(err, shared.ErrInvalidOTPCode):
		if errors.Is(err, shared.ErrInvalidOTPEmail) {
			return http.StatusBadRequest, "AUTH_OTP_INVALID_EMAIL"
		}
		if errors.Is(err, shared.ErrOTPEmailNotRegistered) {
			return http.StatusBadRequest, "AUTH_OTP_EMAIL_NOT_REGISTERED"
		}
		return http.StatusBadRequest, "AUTH_OTP_INVALID_OR_EXPIRED"
	case errors.Is(err, shared.ErrEmailAlreadyRegistered):
		return http.StatusConflict, "AUTH_EMAIL_ALREADY_REGISTERED"
	case errors.Is(err, shared.ErrInvalidRefreshToken):
		return http.StatusUnauthorized, "AUTH_REFRESH_TOKEN_INVALID"
	case errors.Is(err, shared.ErrEmailCheckFailed):
		return http.StatusInternalServerError, "AUTH_EMAIL_CHECK_FAILED"
	default:
		return http.StatusInternalServerError, "AUTH_INTERNAL_ERROR"
	}
}

// UpdateUserHandler - Xử lý PUT /api/users/:id
func UpdateUserHandler(c *gin.Context) {
	// Bước 1: Đọc ID từ URL và body từ request
	userID := c.Param("id")
	var req domain.UpdateUserRequest

	if !EnsurePathUserMatchesToken(c, "id") {
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "AUTH_INVALID_PAYLOAD")
		return
	}

	// Bước 2: Gọi service để xử lý business logic
	updatedUser, err := UpdateUser(userID, req)
	if err != nil {
		// Xác định HTTP status phù hợp dựa vào thông báo lỗi
		status := http.StatusInternalServerError
		code := "AUTH_INTERNAL_ERROR"
		if err.Error() == "Không tìm thấy người dùng" {
			status = http.StatusNotFound
			code = "AUTH_USER_NOT_FOUND"
		} else if err.Error() == "Email đã được sử dụng bởi tài khoản khác" {
			status = http.StatusConflict
			code = "AUTH_EMAIL_ALREADY_REGISTERED"
		}
		respondError(c, status, err.Error(), code)
		return
	}

	// Bước 3: Trả kết quả cho client
	respondSuccess(c, http.StatusOK, "Cập nhật thông tin thành công", gin.H{"user": updatedUser})
}

// GetMeHandler - Xử lý GET /api/users/me
func GetMeHandler(c *gin.Context) {
	authUserID, ok := GetAuthenticatedUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, shared.ErrInvalidAccessToken.Error(), "AUTH_TOKEN_INVALID")
		return
	}

	user, err := GetUserByID(authUserID)
	if err != nil {
		respondError(c, http.StatusNotFound, err.Error(), "AUTH_USER_NOT_FOUND")
		return
	}

	respondSuccess(c, http.StatusOK, "", gin.H{"user": user})
}

// ChangePasswordHandler - Xử lý PUT /api/users/me/password
func ChangePasswordHandler(c *gin.Context) {
	authUserID, ok := GetAuthenticatedUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, shared.ErrInvalidAccessToken.Error(), "AUTH_TOKEN_INVALID")
		return
	}

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "AUTH_INVALID_PAYLOAD")
		return
	}

	if err := ChangePassword(authUserID, req.CurrentPassword, req.NewPassword); err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	respondSuccess(c, http.StatusOK, "Đổi mật khẩu thành công", nil)
}

// ResetPasswordHandler - Xử lý POST /api/password/reset
func ResetPasswordHandler(c *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "AUTH_INVALID_PAYLOAD")
		return
	}

	if err := ResetPassword(req.Email, req.OTPCode, req.NewPassword); err != nil {
		status, code := authErrorInfo(err)
		respondError(c, status, err.Error(), code)
		return
	}

	respondSuccess(c, http.StatusOK, "Đặt lại mật khẩu thành công", nil)
}

// LogoutHandler - Xử lý POST /api/logout
// Backend stateless - token không lưu server, logout chỉ cần xác nhận client clear token
func LogoutHandler(c *gin.Context) {
	respondSuccess(c, http.StatusOK, "Đăng xuất thành công", nil)
}
