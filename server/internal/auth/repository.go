package auth

import (
	"errors"
	"travel-backend/database"
	"travel-backend/domain"

	"gorm.io/gorm"
)

// UserRepository - Tầng duy nhất được phép nói chuyện với database cho User
// Giống như "kho nguyên liệu" trong bếp — chỉ biết cách lấy/lưu dữ liệu
// Không biết gì về HTTP request hay business rules

// FindUserByEmail - Tìm user theo email
// Dùng để kiểm tra email đã tồn tại khi đăng ký và lấy hash password khi đăng nhập
func FindUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// IsNotFoundError - Kiểm tra lỗi không tìm thấy bản ghi
func IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// FindUserByID - Tìm user theo ID (string, dùng cho UpdateUser)
// Dùng để lấy thông tin user trước khi cập nhật
func FindUserByID(id string) (*domain.User, error) {
	var user domain.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// FindUserByIDUint - Tìm user theo ID (uint, dùng cho middleware/JWT claims)
func FindUserByIDUint(id uint) (*domain.User, error) {
	var user domain.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}


// EmailExistsForOtherUser - Kiểm tra email đã được dùng bởi user khác chưa
// Dùng khi cập nhật email: không được trùng với user khác (trừ chính mình)
func EmailExistsForOtherUser(email, currentUserID string) bool {
	var existingUser domain.User
	result := database.DB.Where("email = ? AND id != ?", email, currentUserID).First(&existingUser)
	return result.Error == nil // nil error = tìm thấy = email đã tồn tại
}

// CreateUser - Tạo user mới trong database
func CreateUser(user *domain.User) error {
	return database.DB.Create(user).Error
}

// SaveUser - Lưu (cập nhật) thông tin user vào database
func SaveUser(user *domain.User) error {
	return database.DB.Save(user).Error
}
