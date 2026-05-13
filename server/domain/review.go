package domain

import (
	"fmt"
	"time"
)

// Review status constants
const (
	ReviewStatusPending   = "pending"
	ReviewStatusPublished = "published"
	ReviewStatusHidden    = "hidden"
)

// Review - Model đại diện cho bảng reviews trong database
type Review struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	TourID     uint      `json:"tour_id" gorm:"not null;index"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	BookingID  uint      `json:"booking_id" gorm:"not null;index"`
	Rating     int       `json:"rating" gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Title      string    `json:"title" gorm:"size:200"`
	Content    string    `json:"content" gorm:"type:text;not null"`
	Status     string    `json:"status" gorm:"not null;default:'published';size:20;index"`
	AdminReply *string   `json:"admin_reply" gorm:"type:text"`
	RepliedAt  *time.Time `json:"replied_at"`
	RepliedBy  *uint     `json:"replied_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	Tour    *Tour    `json:"tour,omitempty" gorm:"foreignKey:TourID"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Booking *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
}

// CanEdit - Kiểm tra review có thể sửa không (trong 7 ngày)
func (r *Review) CanEdit() bool {
	return time.Since(r.CreatedAt) <= 7*24*time.Hour
}

// IsPublished - Kiểm tra review đã publish chưa
func (r *Review) IsPublished() bool {
	return r.Status == ReviewStatusPublished
}

// CreateReviewRequest - Request tạo review mới
type CreateReviewRequest struct {
	TourID    uint   `json:"tour_id" binding:"required,min=1"`
	BookingID uint   `json:"booking_id" binding:"required,min=1"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Title     string `json:"title" binding:"max=200"`
	Content   string `json:"content" binding:"required,min=10,max=2000"`
}

// Validate - Validate CreateReviewRequest
func (r *CreateReviewRequest) Validate() error {
	if r.TourID == 0 {
		return fmt.Errorf("tour_id là bắt buộc")
	}
	if r.BookingID == 0 {
		return fmt.Errorf("booking_id là bắt buộc")
	}
	if r.Rating < 1 || r.Rating > 5 {
		return fmt.Errorf("rating phải từ 1 đến 5")
	}
	if len(r.Content) < 10 {
		return fmt.Errorf("nội dung đánh giá tối thiểu 10 ký tự")
	}
	if len(r.Content) > 2000 {
		return fmt.Errorf("nội dung đánh giá tối đa 2000 ký tự")
	}
	return nil
}

// UpdateReviewRequest - Request cập nhật review
type UpdateReviewRequest struct {
	Rating  *int    `json:"rating" binding:"omitempty,min=1,max=5"`
	Title   *string `json:"title" binding:"omitempty,max=200"`
	Content *string `json:"content" binding:"omitempty,min=10,max=2000"`
}

// AdminReplyRequest - Request admin phản hồi review
type AdminReplyRequest struct {
	Reply string `json:"reply" binding:"required,min=1,max=1000"`
}

// ReviewSummary - Tóm tắt review cho public API
type ReviewSummary struct {
	ID         uint       `json:"id"`
	Rating     int        `json:"rating"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	UserName   string     `json:"user_name"`
	UserAvatar string     `json:"user_avatar"`
	AdminReply *string    `json:"admin_reply,omitempty"`
	RepliedAt  *time.Time `json:"replied_at,omitempty"`
	CanEdit    bool       `json:"can_edit"`
	CreatedAt  time.Time  `json:"created_at"`
}

// ToSummary - Chuyển Review thành ReviewSummary
func (r *Review) ToSummary(currentUserID uint) *ReviewSummary {
	summary := &ReviewSummary{
		ID:         r.ID,
		Rating:     r.Rating,
		Title:      r.Title,
		Content:    r.Content,
		AdminReply: r.AdminReply,
		RepliedAt:  r.RepliedAt,
		CanEdit:    r.UserID == currentUserID && r.CanEdit(),
		CreatedAt:  r.CreatedAt,
	}

	if r.User != nil {
		summary.UserName = r.User.Name
		summary.UserAvatar = r.User.AvatarURL
	}

	return summary
}

// TourReviewStats - Thống kê rating cho tour
type TourReviewStats struct {
	TourID       uint    `json:"tour_id"`
	AverageRating float64 `json:"average_rating"`
	TotalReviews int     `json:"total_reviews"`
	Rating5      int     `json:"rating_5"`
	Rating4      int     `json:"rating_4"`
	Rating3      int     `json:"rating_3"`
	Rating2      int     `json:"rating_2"`
	Rating1      int     `json:"rating_1"`
}
