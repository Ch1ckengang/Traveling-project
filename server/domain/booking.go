package domain

import "time"

// Booking - Model đại diện cho bảng bookings trong database
type Booking struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	UserID          uint       `json:"user_id" gorm:"not null;default:0;index"`
	TourID          uint       `json:"tour_id" gorm:"not null;index"`
	Tour            Tour       `json:"tour" gorm:"foreignKey:TourID"`
	FullName        string     `json:"full_name" gorm:"not null;default:''"`
	Phone           string     `json:"phone" gorm:"not null;default:''"`
	Email           string     `json:"email" gorm:"not null;default:''"`
	AdultCount      int        `json:"adult_count" gorm:"not null;default:1"`
	ChildCount      int        `json:"child_count" gorm:"not null;default:0"`
	InfantCount     int        `json:"infant_count" gorm:"not null;default:0"`
	Quantity        int        `json:"quantity" gorm:"not null"`
	TravelDate      string     `json:"travel_date" gorm:"not null"`
	ScheduleID      *uint      `json:"schedule_id" gorm:"index"`
	TotalAmount     int64      `json:"total_amount" gorm:"not null;default:0"`
	CouponCode      string     `json:"coupon_code" gorm:"size:50"`
	DiscountAmount  int64      `json:"discount_amount" gorm:"default:0"`
	BookingCode     string     `json:"booking_code" gorm:"uniqueIndex;default:''"`
	PaymentStatus   string     `json:"payment_status" gorm:"not null;default:'unpaid'"`
	PaymentID       *uint      `json:"payment_id" gorm:"index"`
	PaymentDeadline *time.Time `json:"payment_deadline"`
	Note            string     `json:"note"`
	Status          string     `json:"status" gorm:"not null;default:'booked'"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relationships
	Payment      *Payment      `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
	TourSchedule *TourSchedule `json:"schedule,omitempty" gorm:"foreignKey:ScheduleID"`
}
