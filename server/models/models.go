package models

import (
	"time"
)

// ===== REQUEST STRUCTS =====
// Mô tả hình dạng dữ liệu nhận vào từ client (HTTP request body)

// LoginRequest - Dữ liệu đăng nhập từ client
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest - Dữ liệu đăng ký tài khoản từ client
type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// OTPSendRequest - Dữ liệu yêu cầu gửi mã OTP
type OTPSendRequest struct {
	Email string `json:"email" binding:"required"`
}

// OTPVerifyRequest - Dữ liệu xác thực OTP
type OTPVerifyRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

// ForgotPasswordRequest - Dữ liệu gửi yêu cầu quên mật khẩu
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

// UpdateUserRequest - Dữ liệu cập nhật thông tin cá nhân từ client
// Tất cả các trường đều optional
type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// CreateBookingRequest - Dữ liệu đặt tour từ client
type CreateBookingRequest struct {
	UserID      uint   `json:"user_id"`
	TourID      uint   `json:"tour_id" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"required"`
	AdultCount  int    `json:"adult_count"`
	ChildCount  int    `json:"child_count"`
	InfantCount int    `json:"infant_count"`
	Quantity    int    `json:"quantity"`
	TravelDate  string `json:"travel_date" binding:"required"`
	Note        string `json:"note,omitempty"`
}

// ===== RESPONSE STRUCTS =====
// Mô tả hình dạng dữ liệu trả về cho client (HTTP response body)

// AuthResponse - Response cho các API authentication
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	User    *User  `json:"user,omitempty"`
}

// BookingResponse - Response cho API đặt tour
type BookingResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Booking *Booking `json:"booking,omitempty"`
}

// BookingListResponse - Response cho API lấy danh sách booking
type BookingListResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Bookings []Booking `json:"bookings,omitempty"`
}

// User - Model đại diện cho bảng users trong database
// GORM sẽ tự động tạo bảng dựa trên struct này
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`         // Primary key, auto increment
	Name      string    `json:"name" gorm:"not null"`         // Tên user, không được null
	Email     string    `json:"email" gorm:"unique;not null"` // Email, phải unique và không null
	Password  string    `json:"-" gorm:"not null"`            // Password, json:"-" không trả về frontend (bảo mật)
	CreatedAt time.Time `json:"created_at"`                   // Thời gian tạo, tự động set bởi GORM
	UpdatedAt time.Time `json:"updated_at"`                   // Thời gian cập nhật cuối, tự động update bởi GORM
}

// Tour - Model đại diện cho bảng tours trong database
// Chứa thông tin về các tour du lịch
type Tour struct {
	ID             uint      `json:"id" gorm:"primaryKey"` // Primary key, auto increment
	Name           string    `json:"name" gorm:"not null"` // Tên tour
	Type           string    `json:"type" gorm:"not null;default:'domestic'"`
	Price          string    `json:"price" gorm:"not null"` // Giá tour (string để lưu format "2.000.000đ")
	Description    string    `json:"description"`           // Mô tả tour
	Location       string    `json:"location"`              // Địa điểm tour
	Country        string    `json:"country" gorm:"not null;default:'Việt Nam'"`
	Duration       string    `json:"duration"` // Thời gian tour (ví dụ: "3 ngày 2 đêm")
	DepartureDate  string    `json:"departure_date" gorm:"default:''"`
	RemainingSlots int       `json:"remaining_slots" gorm:"not null;default:30"`
	Itinerary      string    `json:"itinerary"`
	Services       string    `json:"services"`
	CreatedAt      time.Time `json:"created_at"` // Thời gian tạo
	UpdatedAt      time.Time `json:"updated_at"` // Thời gian cập nhật
}

// Booking - Model đại diện cho bảng bookings trong database
type Booking struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	UserID        uint      `json:"user_id" gorm:"not null;default:0;index"`
	TourID        uint      `json:"tour_id" gorm:"not null;index"`
	Tour          Tour      `json:"tour" gorm:"foreignKey:TourID"`
	FullName      string    `json:"full_name" gorm:"not null;default:''"`
	Phone         string    `json:"phone" gorm:"not null;default:''"`
	Email         string    `json:"email" gorm:"not null;default:''"`
	AdultCount    int       `json:"adult_count" gorm:"not null;default:1"`
	ChildCount    int       `json:"child_count" gorm:"not null;default:0"`
	InfantCount   int       `json:"infant_count" gorm:"not null;default:0"`
	Quantity      int       `json:"quantity" gorm:"not null"`
	TravelDate    string    `json:"travel_date" gorm:"not null"`
	TotalAmount   int64     `json:"total_amount" gorm:"not null;default:0"`
	BookingCode   string    `json:"booking_code" gorm:"uniqueIndex;default:''"`
	PaymentStatus string    `json:"payment_status" gorm:"not null;default:'unpaid'"`
	Note          string    `json:"note"`
	Status        string    `json:"status" gorm:"not null;default:'booked'"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
