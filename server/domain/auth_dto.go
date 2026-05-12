package domain

// LoginRequest - Dữ liệu đăng nhập từ client
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest - Dữ liệu đăng ký tài khoản từ client
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone"`
	Password string `json:"password" binding:"required"`
}

// OTPSendRequest - Dữ liệu yêu cầu gửi mã OTP
type OTPSendRequest struct {
	Email string `json:"email" binding:"required"`
}

// OTPVerifyRequest - Dữ liệu xác thực OTP
type OTPVerifyRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// ForgotPasswordRequest - Dữ liệu gửi yêu cầu quên mật khẩu
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

// RefreshTokenRequest - Dữ liệu làm mới phiên đăng nhập
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"` // Changed from refresh_token to refreshToken
}

// UpdateUserRequest - Dữ liệu cập nhật thông tin cá nhân từ client
// Tất cả các trường đều optional
type UpdateUserRequest struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
	Password  string `json:"password,omitempty"`
}

// ResetPasswordRequest - Dữ liệu đặt lại mật khẩu (dùng OTP token)
type ResetPasswordRequest struct {
	Email       string `json:"email" binding:"required"`
	OTPCode     string `json:"otp_code" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangePasswordRequest - Dữ liệu đổi mật khẩu (yêu cầu mật khẩu hiện tại)
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

// GetMeResponse - Response cho GET /users/me
type GetMeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// AuthResponse - Response cho các API authentication
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// TokenPair - Cặp token cho phiên đăng nhập
type TokenPair struct {
	TokenType        string `json:"token_type"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}
