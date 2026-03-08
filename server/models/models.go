package models

import (
	"time"
)

// User - Model đại diện cho bảng users trong database
// GORM sẽ tự động tạo bảng dựa trên struct này
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                      // Primary key, auto increment
	Name      string    `json:"name" gorm:"not null"`                      // Tên user, không được null
	Email     string    `json:"email" gorm:"unique;not null"`              // Email, phải unique và không null
	Password  string    `json:"-" gorm:"not null"`                         // Password, json:"-" không trả về frontend (bảo mật)
	CreatedAt time.Time `json:"created_at"`                                 // Thời gian tạo, tự động set bởi GORM
	UpdatedAt time.Time `json:"updated_at"`                                 // Thời gian cập nhật cuối, tự động update bởi GORM
}

// Tour - Model đại diện cho bảng tours trong database
// Chứa thông tin về các tour du lịch
type Tour struct {
	ID          uint      `json:"id" gorm:"primaryKey"`      // Primary key, auto increment
	Name        string    `json:"name" gorm:"not null"`      // Tên tour
	Price       string    `json:"price" gorm:"not null"`     // Giá tour (string để lưu format "2.000.000đ")
	Description string    `json:"description"`               // Mô tả tour
	Location    string    `json:"location"`                  // Địa điểm tour
	Duration    string    `json:"duration"`                  // Thời gian tour (ví dụ: "3 ngày 2 đêm")
	CreatedAt   time.Time `json:"created_at"`                // Thời gian tạo
	UpdatedAt   time.Time `json:"updated_at"`                // Thời gian cập nhật
}
