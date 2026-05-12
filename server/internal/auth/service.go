package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/mail"
	"strings"
	"time"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"golang.org/x/crypto/bcrypt"
)

// Removed in-memory OTP storage - now using database

// UserService - Tầng chứa business logic cho User
// Giống như "đầu bếp" — biết phải làm gì với nguyên liệu
// Gọi repository để lấy/lưu dữ liệu, áp dụng các quy tắc nghiệp vụ
// KHÔNG biết gì về HTTP request hay SQL

// Login - Xử lý logic đăng nhập
// Business rule: email và password phải khớp với dữ liệu trong DB
func Login(email, password string) (*domain.User, error) {
	email = normalizeEmail(email)
	if err := validateLoginInput(email, password); err != nil {
		return nil, err
	}

	user, err := FindUserByEmail(email)
	if err != nil {
		// Không tiết lộ email có tồn tại hay không
		if IsNotFoundError(err) {
			return nil, shared.ErrInvalidCredentials
		}
		return nil, shared.ErrLoginProcessFailed
	}

	if err := verifyPassword(user.Password, password); err != nil {
		return nil, shared.ErrInvalidCredentials
	}

	if !user.IsEmailVerified {
		return nil, shared.ErrEmailNotVerified
	}

	return user, nil
}

// Register - Xử lý logic đăng ký tài khoản mới
// Business rules:
//  1. Email không được đã tồn tại trong hệ thống
//  2. Nếu email chưa tồn tại → tạo user mới
func Register(name, email, phone, password string) (*domain.User, error) {
	name = strings.TrimSpace(name)
	email = normalizeEmail(email)
	phone = strings.TrimSpace(phone)
	if err := validateRegisterInput(name, email, password); err != nil {
		return nil, err
	}

	// Rule 1: Kiểm tra email đã tồn tại chưa
	_, err := FindUserByEmail(email)
	if err == nil {
		// err == nil nghĩa là tìm thấy user → email đã được đăng ký
		return nil, shared.ErrEmailAlreadyRegistered
	}
	if !IsNotFoundError(err) {
		return nil, shared.ErrEmailCheckFailed
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, shared.ErrPasswordHashFailed
	}

	// Email chưa tồn tại → tạo user mới
	newUser := &domain.User{
		Name:            name,
		Email:           email,
		Phone:           phone,
		Password:        hashedPassword,
		IsEmailVerified: false,
	}

	if err := CreateUser(newUser); err != nil {
		return nil, shared.ErrCreateAccountFailed
	}

	if err := SendOTPForEmail(email); err != nil {
		return nil, err
	}

	return newUser, nil
}

// SendOTPForEmail - Tạo mã OTP 6 chữ số cho email hợp lệ và lưu vào database
func SendOTPForEmail(email string) error {
	email = normalizeEmail(email)
	if _, err := mail.ParseAddress(email); err != nil {
		return shared.ErrInvalidOTPEmail
	}

	if _, err := FindUserByEmail(email); err != nil {
		if IsNotFoundError(err) {
			return shared.ErrOTPEmailNotRegistered
		}
		return shared.ErrEmailCheckFailed
	}

	code, err := generateOTPCode(6)
	if err != nil {
		return shared.ErrOTPProcessFailed
	}

	// Xóa các OTP cũ của email này trước khi tạo mới
	_ = DeleteOTPByEmail(email)

	// Lưu OTP vào database
	otp := &domain.OTP{
		Email:     email,
		Code:      code,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}

	if err := SaveOTP(otp); err != nil {
		return shared.ErrOTPProcessFailed
	}

	// Gửi email chứa OTP
	if err := shared.SendOTPEmail(email, code); err != nil {
		log.Printf("⚠️ Failed to send OTP email to %s: %v", email, err)
		// Không return error vì OTP đã lưu thành công
		// User vẫn có thể thấy OTP trong log (dev mode)
	}

	return nil
}

// VerifyOTPForEmail - Kiểm tra mã OTP còn hạn và khớp với email từ database
func VerifyOTPForEmail(email, code string) error {
	email = normalizeEmail(email)
	code = strings.TrimSpace(code)

	if _, err := mail.ParseAddress(email); err != nil {
		return shared.ErrInvalidOTPEmail
	}

	if len(code) != 6 {
		return shared.ErrInvalidOTPCode
	}

	// Tìm OTP từ database
	otp, err := FindOTPByEmail(email)
	if err != nil {
		return shared.ErrInvalidOTPCode
	}

	// Kiểm tra OTP đã hết hạn chưa
	if otp.IsExpired() {
		_ = DeleteOTPByEmail(email)
		return shared.ErrInvalidOTPCode
	}

	// Kiểm tra code có khớp không
	if otp.Code != code {
		return shared.ErrInvalidOTPCode
	}

	// Xóa OTP sau khi verify thành công
	_ = DeleteOTPByEmail(email)

	// Cập nhật user là đã verify email
	user, err := FindUserByEmail(email)
	if err != nil {
		if IsNotFoundError(err) {
			return shared.ErrOTPEmailNotRegistered
		}
		return shared.ErrEmailCheckFailed
	}

	if !user.IsEmailVerified {
		user.IsEmailVerified = true
		if err := SaveUser(user); err != nil {
			return shared.ErrCreateAccountFailed
		}
	}

	return nil
}

// RequestPasswordReset - Gửi email reset password dùng OTP system
func RequestPasswordReset(email string) error {
	email = normalizeEmail(email)
	if _, err := mail.ParseAddress(email); err != nil {
		return shared.ErrInvalidRegisterEmail
	}

	_, err := FindUserByEmail(email)
	if err != nil {
		if IsNotFoundError(err) {
			// Không tiết lộ email có tồn tại hay không
			return nil
		}
		return shared.ErrEmailCheckFailed
	}

	// Tái sử dụng OTP system để gửi reset code
	if err := SendOTPForEmail(email); err != nil {
		log.Printf("⚠️ Failed to send password reset OTP to %s: %v", email, err)
	}

	return nil
}

// GetUserByID - Lấy thông tin user theo uint ID (dùng cho GET /users/me)
func GetUserByID(userID uint) (*domain.User, error) {
	user, err := FindUserByIDUint(userID)
	if err != nil {
		if IsNotFoundError(err) {
			return nil, errors.New("Không tìm thấy người dùng")
		}
		return nil, errors.New("Không thể lấy thông tin người dùng")
	}
	return user, nil
}

// ChangePassword - Đổi mật khẩu (yêu cầu mật khẩu hiện tại)
func ChangePassword(userID uint, currentPassword, newPassword string) error {
	user, err := FindUserByIDUint(userID)
	if err != nil {
		return errors.New("Không tìm thấy người dùng")
	}

	if err := verifyPassword(user.Password, currentPassword); err != nil {
		return shared.ErrInvalidCredentials
	}

	if len(newPassword) < 8 {
		return shared.ErrWeakPassword
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return shared.ErrPasswordHashFailed
	}

	user.Password = hashedPassword
	if err := SaveUser(user); err != nil {
		return errors.New("Không thể cập nhật mật khẩu")
	}

	return nil
}

// ResetPassword - Đặt lại mật khẩu bằng OTP token
func ResetPassword(email, otpCode, newPassword string) error {
	email = normalizeEmail(email)

	// Verify OTP trước
	if err := VerifyOTPForEmail(email, otpCode); err != nil {
		return err
	}

	if len(newPassword) < 8 {
		return shared.ErrWeakPassword
	}

	user, err := FindUserByEmail(email)
	if err != nil {
		return errors.New("Không tìm thấy người dùng")
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return shared.ErrPasswordHashFailed
	}

	user.Password = hashedPassword
	if err := SaveUser(user); err != nil {
		return errors.New("Không thể cập nhật mật khẩu")
	}

	return nil
}

// UpdateUser - Xử lý logic cập nhật thông tin cá nhân
// Business rules:
//  1. User phải tồn tại
//  2. Nếu đổi email → email mới không được trùng với user khác
//  3. Chỉ cập nhật các trường được gửi lên (không bỏ trống)
func UpdateUser(userID string, req domain.UpdateUserRequest) (*domain.User, error) {
	// Rule 1: Tìm user hiện tại
	user, err := FindUserByID(userID)
	if err != nil {
		return nil, errors.New("Không tìm thấy người dùng")
	}

	// Rule 2: Nếu muốn đổi email → kiểm tra không được trùng với user khác
	if req.Email != "" && req.Email != user.Email {
		if EmailExistsForOtherUser(req.Email, userID) {
			return nil, errors.New("Email đã được sử dụng bởi tài khoản khác")
		}
		user.Email = req.Email
	}

	// Rule 3: Chỉ cập nhật các trường được gửi
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Phone != "" {
		user.Phone = strings.TrimSpace(req.Phone)
	}
	if req.AvatarURL != "" {
		user.AvatarURL = strings.TrimSpace(req.AvatarURL)
	}
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return nil, errors.New("Không thể bảo mật mật khẩu")
		}
		user.Password = hashedPassword
	}

	// Lưu thay đổi qua repository
	if err := SaveUser(user); err != nil {
		return nil, errors.New("Không thể cập nhật thông tin")
	}

	return user, nil
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func verifyPassword(hashedPassword, rawPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func validateLoginInput(email, password string) error {
	if email == "" || strings.TrimSpace(password) == "" {
		return shared.ErrInvalidAuthPayload
	}
	return nil
}

func validateRegisterInput(name, email, password string) error {
	if name == "" {
		return shared.ErrInvalidName
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return shared.ErrInvalidRegisterEmail
	}

	// Removed @gmail.com restriction - accept all valid email domains

	if len(password) < 8 {
		return shared.ErrWeakPassword
	}

	return nil
}

func generateOTPCode(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid otp length")
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		result[i] = byte('0' + n.Int64())
	}

	return string(result), nil
}
