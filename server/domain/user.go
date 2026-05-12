package domain

import "time"

// User - Model đại diện cho bảng users trong database
// GORM sẽ tự động tạo bảng dựa trên struct này
type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Email           string    `json:"email" gorm:"unique;not null"`
	Password        string    `json:"-" gorm:"not null"`
	Phone           string    `json:"phone" gorm:"default:''"`
	AvatarURL       string    `json:"avatar_url" gorm:"default:''"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"not null;default:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
