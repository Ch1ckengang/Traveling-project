package domain

import "time"

// NotificationType constants
const (
	NotifTypeBooking = "booking"
	NotifTypeSystem  = "system"
	NotifTypeReview  = "review"
	NotifTypePayment = "payment"
)

// Notification - Model lưu trữ thông báo cho người dùng
type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Title     string    `json:"title" gorm:"not null"`
	Message   string    `json:"message" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null;default:'system'"`
	IsRead    bool      `json:"is_read" gorm:"not null;default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
