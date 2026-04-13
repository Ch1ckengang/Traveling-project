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

// UpdateUserRequest - Dữ liệu cập nhật thông tin cá nhân từ client
// Tất cả các trường đều optional
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// AuthResponse - Response cho các API authentication
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
}
