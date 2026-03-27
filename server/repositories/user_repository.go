package repositories

import (
	"travel-backend/database"
	"travel-backend/models"
)

// UserRepository - Tầng duy nhất được phép nói chuyện với database cho User
// Giống như "kho nguyên liệu" trong bếp — chỉ biết cách lấy/lưu dữ liệu
// Không biết gì về HTTP request hay business rules

// FindByEmailAndPassword - Tìm user theo email và password
// Dùng cho chức năng đăng nhập
func FindUserByEmailAndPassword(email, password string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ? AND password = ?", email, password).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindUserByEmail - Tìm user theo email
// Dùng để kiểm tra email đã tồn tại khi đăng ký
func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindUserByID - Tìm user theo ID
// Dùng để lấy thông tin user trước khi cập nhật
func FindUserByID(id string) (*models.User, error) {
	var user models.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// EmailExistsForOtherUser - Kiểm tra email đã được dùng bởi user khác chưa
// Dùng khi cập nhật email: không được trùng với user khác (trừ chính mình)
func EmailExistsForOtherUser(email, currentUserID string) bool {
	var existingUser models.User
	result := database.DB.Where("email = ? AND id != ?", email, currentUserID).First(&existingUser)
	return result.Error == nil // nil error = tìm thấy = email đã tồn tại
}

// CreateUser - Tạo user mới trong database
func CreateUser(user *models.User) error {
	return database.DB.Create(user).Error
}

// SaveUser - Lưu (cập nhật) thông tin user vào database
func SaveUser(user *models.User) error {
	return database.DB.Save(user).Error
}
