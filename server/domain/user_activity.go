package domain

import (
	"time"
)

// UserActivity - Lưu trữ lịch sử hành vi của người dùng
type UserActivity struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        *uint     `json:"user_id" gorm:"index"`             // Có thể null nếu user chưa đăng nhập
	SessionID     string    `json:"session_id" gorm:"index"`          // Cookie/Session ID lưu ở client
	ActionType    string    `json:"action_type" gorm:"index;size:50"` // view_tour, search, view_category
	TourID        *uint     `json:"tour_id" gorm:"index"`             // ID tour nếu action = view_tour
	SearchKeyword string    `json:"search_keyword"`                   // Từ khóa nếu action = search
	Category      string    `json:"category"`                         // Tên category nếu view_category
	Timestamp     time.Time `json:"timestamp" gorm:"index"`           // Thời điểm xảy ra hành động
}
