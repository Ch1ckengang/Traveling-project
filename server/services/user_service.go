package services

import (
	"errors"
	"net/mail"
	"strings"
	"travel-backend/models"
	"travel-backend/repositories"

	"golang.org/x/crypto/bcrypt"
)

// UserService - Tầng chứa business logic cho User
// Giống như "đầu bếp" — biết phải làm gì với nguyên liệu
// Gọi repository để lấy/lưu dữ liệu, áp dụng các quy tắc nghiệp vụ
// KHÔNG biết gì về HTTP request hay SQL

// Login - Xử lý logic đăng nhập
// Business rule: email và password phải khớp với dữ liệu trong DB
func Login(email, password string) (*models.User, error) {
	email = normalizeEmail(email)
	if err := validateLoginInput(email, password); err != nil {
		return nil, err
	}

	user, err := repositories.FindUserByEmail(email)
	if err != nil {
		// Không tiết lộ email có tồn tại hay không
		if repositories.IsNotFoundError(err) {
			return nil, ErrInvalidCredentials
		}
		return nil, ErrLoginProcessFailed
	}

	if err := verifyPassword(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

// Register - Xử lý logic đăng ký tài khoản mới
// Business rules:
//  1. Email không được đã tồn tại trong hệ thống
//  2. Nếu email chưa tồn tại → tạo user mới
func Register(name, email, password string) (*models.User, error) {
	name = strings.TrimSpace(name)
	email = normalizeEmail(email)
	if err := validateRegisterInput(name, email, password); err != nil {
		return nil, err
	}

	// Rule 1: Kiểm tra email đã tồn tại chưa
	_, err := repositories.FindUserByEmail(email)
	if err == nil {
		// err == nil nghĩa là tìm thấy user → email đã được đăng ký
		return nil, ErrEmailAlreadyRegistered
	}
	if !repositories.IsNotFoundError(err) {
		return nil, ErrEmailCheckFailed
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, ErrPasswordHashFailed
	}

	// Email chưa tồn tại → tạo user mới
	newUser := &models.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	if err := repositories.CreateUser(newUser); err != nil {
		return nil, ErrCreateAccountFailed
	}

	return newUser, nil
}

// UpdateUser - Xử lý logic cập nhật thông tin cá nhân
// Business rules:
//  1. User phải tồn tại
//  2. Nếu đổi email → email mới không được trùng với user khác
//  3. Chỉ cập nhật các trường được gửi lên (không bỏ trống)
func UpdateUser(userID string, req models.UpdateUserRequest) (*models.User, error) {
	// Rule 1: Tìm user hiện tại
	user, err := repositories.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("Không tìm thấy người dùng")
	}

	// Rule 2: Nếu muốn đổi email → kiểm tra không được trùng với user khác
	if req.Email != "" && req.Email != user.Email {
		if repositories.EmailExistsForOtherUser(req.Email, userID) {
			return nil, errors.New("Email đã được sử dụng bởi tài khoản khác")
		}
		user.Email = req.Email
	}

	// Rule 3: Chỉ cập nhật các trường được gửi
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Password != "" {
		hashedPassword, err := hashPassword(req.Password)
		if err != nil {
			return nil, errors.New("Không thể bảo mật mật khẩu")
		}
		user.Password = hashedPassword
	}

	// Lưu thay đổi qua repository
	if err := repositories.SaveUser(user); err != nil {
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
		return ErrInvalidAuthPayload
	}
	return nil
}

func validateRegisterInput(name, email, password string) error {
	if name == "" {
		return ErrInvalidName
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidRegisterEmail
	}

	if len(password) < 8 {
		return ErrWeakPassword
	}

	return nil
}
