package services

import (
	"errors"
	"travel-backend/models"
	"travel-backend/repositories"
)

// UserService - Tầng chứa business logic cho User
// Giống như "đầu bếp" — biết phải làm gì với nguyên liệu
// Gọi repository để lấy/lưu dữ liệu, áp dụng các quy tắc nghiệp vụ
// KHÔNG biết gì về HTTP request hay SQL

// Login - Xử lý logic đăng nhập
// Business rule: email và password phải khớp với dữ liệu trong DB
func Login(email, password string) (*models.User, error) {
	user, err := repositories.FindUserByEmailAndPassword(email, password)
	if err != nil {
		// Repository trả về lỗi → email hoặc password sai
		return nil, errors.New("Email hoặc mật khẩu không đúng")
	}
	return user, nil
}

// Register - Xử lý logic đăng ký tài khoản mới
// Business rules:
//  1. Email không được đã tồn tại trong hệ thống
//  2. Nếu email chưa tồn tại → tạo user mới
func Register(name, email, password string) (*models.User, error) {
	// Rule 1: Kiểm tra email đã tồn tại chưa
	_, err := repositories.FindUserByEmail(email)
	if err == nil {
		// err == nil nghĩa là tìm thấy user → email đã được đăng ký
		return nil, errors.New("Email đã được đăng ký")
	}

	// Email chưa tồn tại → tạo user mới
	newUser := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := repositories.CreateUser(newUser); err != nil {
		return nil, errors.New("Không thể tạo tài khoản")
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
		user.Password = req.Password
	}

	// Lưu thay đổi qua repository
	if err := repositories.SaveUser(user); err != nil {
		return nil, errors.New("Không thể cập nhật thông tin")
	}

	return user, nil
}
