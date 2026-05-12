package domain

import (
	"fmt"
	"time"
)

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

// BookingResponse - Response cho API đặt tour
type BookingResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message,omitempty"`
	Booking *BookingWithPayment `json:"booking,omitempty"`
}

// BookingListResponse - Response cho API lấy danh sách booking
type BookingListResponse struct {
	Success  bool                  `json:"success"`
	Message  string                `json:"message,omitempty"`
	Bookings []*BookingWithPayment `json:"bookings,omitempty"`
}

// BookingWithPayment - Booking với thông tin thanh toán chi tiết
type BookingWithPayment struct {
	ID              uint       `json:"id"`
	UserID          uint       `json:"user_id"`
	TourID          uint       `json:"tour_id"`
	Tour            Tour       `json:"tour"`
	FullName        string     `json:"full_name"`
	Phone           string     `json:"phone"`
	Email           string     `json:"email"`
	AdultCount      int        `json:"adult_count"`
	ChildCount      int        `json:"child_count"`
	InfantCount     int        `json:"infant_count"`
	Quantity        int        `json:"quantity"`
	TravelDate      string     `json:"travel_date"`
	TotalAmount     int64      `json:"total_amount"`
	BookingCode     string     `json:"booking_code"`
	PaymentStatus   string     `json:"payment_status"`
	PaymentID       *uint      `json:"payment_id"`
	PaymentDeadline *time.Time `json:"payment_deadline"`
	Note            string     `json:"note"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Payment information
	PaymentInfo    *PaymentSummary   `json:"payment_info,omitempty"`
	PaymentHistory []*PaymentSummary `json:"payment_history,omitempty"`
}

// BookingSummary - Booking tóm tắt với thông tin thanh toán cơ bản
type BookingSummary struct {
	ID              uint       `json:"id"`
	BookingCode     string     `json:"booking_code"`
	TourName        string     `json:"tour_name"`
	TravelDate      string     `json:"travel_date"`
	TotalAmount     int64      `json:"total_amount"`
	PaymentStatus   string     `json:"payment_status"`
	PaymentDeadline *time.Time `json:"payment_deadline"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`

	// Payment summary
	CanPay               bool       `json:"can_pay"`
	PaymentURL           string     `json:"payment_url,omitempty"`
	TransactionReference string     `json:"transaction_reference,omitempty"`
	PaymentExpiresAt     *time.Time `json:"payment_expires_at,omitempty"`
}

// UpdateBookingPaymentRequest - Request để cập nhật thông tin thanh toán của booking
type UpdateBookingPaymentRequest struct {
	BookingID       uint       `json:"booking_id" binding:"required"`
	PaymentStatus   string     `json:"payment_status" binding:"required"`
	PaymentDeadline *time.Time `json:"payment_deadline"`
}

// BookingPaymentStatusResponse - Response cho việc kiểm tra trạng thái thanh toán của booking
type BookingPaymentStatusResponse struct {
	Success         bool            `json:"success"`
	BookingID       uint            `json:"booking_id"`
	BookingCode     string          `json:"booking_code"`
	PaymentStatus   string          `json:"payment_status"`
	PaymentDeadline *time.Time      `json:"payment_deadline"`
	CanPay          bool            `json:"can_pay"`
	PaymentInfo     *PaymentSummary `json:"payment_info,omitempty"`
	Message         string          `json:"message,omitempty"`
}

// Helper methods and conversions

// ToBookingWithPayment - Chuyển đổi Booking thành BookingWithPayment
func (b *Booking) ToBookingWithPayment() *BookingWithPayment {
	return &BookingWithPayment{
		ID:              b.ID,
		UserID:          b.UserID,
		TourID:          b.TourID,
		Tour:            b.Tour,
		FullName:        b.FullName,
		Phone:           b.Phone,
		Email:           b.Email,
		AdultCount:      b.AdultCount,
		ChildCount:      b.ChildCount,
		InfantCount:     b.InfantCount,
		Quantity:        b.Quantity,
		TravelDate:      b.TravelDate,
		TotalAmount:     b.TotalAmount,
		BookingCode:     b.BookingCode,
		PaymentStatus:   b.PaymentStatus,
		PaymentID:       b.PaymentID,
		PaymentDeadline: b.PaymentDeadline,
		Note:            b.Note,
		Status:          b.Status,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       b.UpdatedAt,
		PaymentInfo:     nil, // Will be populated separately
		PaymentHistory:  nil, // Will be populated separately
	}
}

// ToBookingSummary - Chuyển đổi Booking thành BookingSummary
func (b *Booking) ToBookingSummary() *BookingSummary {
	tourName := ""
	if b.Tour.ID != 0 {
		tourName = b.Tour.Name
	}

	return &BookingSummary{
		ID:                   b.ID,
		BookingCode:          b.BookingCode,
		TourName:             tourName,
		TravelDate:           b.TravelDate,
		TotalAmount:          b.TotalAmount,
		PaymentStatus:        b.PaymentStatus,
		PaymentDeadline:      b.PaymentDeadline,
		Status:               b.Status,
		CreatedAt:            b.CreatedAt,
		CanPay:               b.CanPay(),
		PaymentURL:           "",
		TransactionReference: "",
		PaymentExpiresAt:     nil,
	}
}

// CanPay - Kiểm tra xem booking có thể thanh toán hay không
func (b *Booking) CanPay() bool {
	// Chỉ có thể thanh toán khi:
	// 1. Trạng thái booking là "pending" hoặc "booked" (legacy)
	// 2. Trạng thái thanh toán là "unpaid" hoặc "failed"
	// 3. Chưa hết hạn thanh toán (nếu có)

	if b.Status != "pending" && b.Status != "booked" {
		return false
	}

	if b.PaymentStatus != "unpaid" && b.PaymentStatus != "failed" {
		return false
	}

	// Kiểm tra deadline thanh toán
	if b.PaymentDeadline != nil && time.Now().After(*b.PaymentDeadline) {
		return false
	}

	return true
}

// IsPaymentExpired - Kiểm tra xem thanh toán có hết hạn hay không
func (b *Booking) IsPaymentExpired() bool {
	if b.PaymentDeadline == nil {
		return false
	}
	return time.Now().After(*b.PaymentDeadline)
}

// GetPaymentTimeRemaining - Lấy thời gian còn lại để thanh toán
func (b *Booking) GetPaymentTimeRemaining() *time.Duration {
	if b.PaymentDeadline == nil {
		return nil
	}

	remaining := time.Until(*b.PaymentDeadline)
	if remaining < 0 {
		remaining = 0
	}

	return &remaining
}

// WithPaymentInfo - Thêm thông tin thanh toán vào BookingWithPayment
func (bwp *BookingWithPayment) WithPaymentInfo(payment *PaymentSummary) *BookingWithPayment {
	bwp.PaymentInfo = payment
	return bwp
}

// WithPaymentHistory - Thêm lịch sử thanh toán vào BookingWithPayment
func (bwp *BookingWithPayment) WithPaymentHistory(history []*PaymentSummary) *BookingWithPayment {
	bwp.PaymentHistory = history
	return bwp
}

// WithPaymentURL - Thêm URL thanh toán vào BookingSummary
func (bs *BookingSummary) WithPaymentURL(paymentURL, transactionRef string, expiresAt *time.Time) *BookingSummary {
	bs.PaymentURL = paymentURL
	bs.TransactionReference = transactionRef
	bs.PaymentExpiresAt = expiresAt
	return bs
}

// Validation methods

// Validate - Validate UpdateBookingPaymentRequest
func (r *UpdateBookingPaymentRequest) Validate() error {
	if r.BookingID == 0 {
		return fmt.Errorf("booking_id is required and must be greater than 0")
	}

	validStatuses := []string{"unpaid", "processing", "paid", "failed", "refunded", "expired"}
	isValidStatus := false
	for _, status := range validStatuses {
		if r.PaymentStatus == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return fmt.Errorf("payment_status must be one of: %v", validStatuses)
	}

	return nil
}

// IsValidPaymentStatus - Kiểm tra trạng thái thanh toán có hợp lệ hay không
func IsValidPaymentStatus(status string) bool {
	validStatuses := []string{"unpaid", "processing", "paid", "failed", "refunded", "expired"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// IsValidBookingStatus - Kiểm tra trạng thái booking có hợp lệ hay không
func IsValidBookingStatus(status string) bool {
	validStatuses := []string{"pending", "booked", "confirmed", "cancelled", "completed"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// GetPaymentStatusDisplayName - Lấy tên hiển thị cho trạng thái thanh toán
func GetPaymentStatusDisplayName(status string) string {
	statusNames := map[string]string{
		"unpaid":     "Chưa thanh toán",
		"processing": "Đang xử lý",
		"paid":       "Đã thanh toán",
		"failed":     "Thanh toán thất bại",
		"refunded":   "Đã hoàn tiền",
		"expired":    "Hết hạn thanh toán",
	}

	if displayName, exists := statusNames[status]; exists {
		return displayName
	}

	return "Không xác định"
}

// GetBookingStatusDisplayName - Lấy tên hiển thị cho trạng thái booking
func GetBookingStatusDisplayName(status string) string {
	statusNames := map[string]string{
		"booked":    "Đã đặt",
		"pending":   "Chờ thanh toán",
		"confirmed": "Đã xác nhận",
		"cancelled": "Đã hủy",
		"completed": "Hoàn thành",
	}

	if displayName, exists := statusNames[status]; exists {
		return displayName
	}

	return "Không xác định"
}
