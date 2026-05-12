package domain

import "time"

// OTP - Model đại diện cho bảng otp_codes trong database
// Lưu trữ mã OTP với thời gian hết hạn
type OTP struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"index;not null"`
	Code      string    `json:"code" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// IsExpired - Kiểm tra OTP đã hết hạn chưa
func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}
