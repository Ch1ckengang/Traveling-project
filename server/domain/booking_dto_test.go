package domain

import (
	"testing"
	"time"
)

func TestBookingCanPay(t *testing.T) {
	tests := []struct {
		name     string
		booking  *Booking
		expected bool
	}{
		{
			name: "can pay - unpaid booking",
			booking: &Booking{
				Status:        "booked",
				PaymentStatus: "unpaid",
			},
			expected: true,
		},
		{
			name: "can pay - failed payment",
			booking: &Booking{
				Status:        "booked",
				PaymentStatus: "failed",
			},
			expected: true,
		},
		{
			name: "cannot pay - already paid",
			booking: &Booking{
				Status:        "booked",
				PaymentStatus: "paid",
			},
			expected: false,
		},
		{
			name: "cannot pay - cancelled booking",
			booking: &Booking{
				Status:        "cancelled",
				PaymentStatus: "unpaid",
			},
			expected: false,
		},
		{
			name: "cannot pay - expired deadline",
			booking: &Booking{
				Status:          "booked",
				PaymentStatus:   "unpaid",
				PaymentDeadline: &time.Time{}, // Zero time (past)
			},
			expected: false,
		},
		{
			name: "can pay - future deadline",
			booking: &Booking{
				Status:        "booked",
				PaymentStatus: "unpaid",
				PaymentDeadline: func() *time.Time {
					future := time.Now().Add(24 * time.Hour)
					return &future
				}(),
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.booking.CanPay()
			if result != tt.expected {
				t.Errorf("CanPay() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBookingIsPaymentExpired(t *testing.T) {
	tests := []struct {
		name     string
		booking  *Booking
		expected bool
	}{
		{
			name: "no deadline - not expired",
			booking: &Booking{
				PaymentDeadline: nil,
			},
			expected: false,
		},
		{
			name: "past deadline - expired",
			booking: &Booking{
				PaymentDeadline: func() *time.Time {
					past := time.Now().Add(-1 * time.Hour)
					return &past
				}(),
			},
			expected: true,
		},
		{
			name: "future deadline - not expired",
			booking: &Booking{
				PaymentDeadline: func() *time.Time {
					future := time.Now().Add(1 * time.Hour)
					return &future
				}(),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.booking.IsPaymentExpired()
			if result != tt.expected {
				t.Errorf("IsPaymentExpired() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestBookingToBookingWithPayment(t *testing.T) {
	booking := &Booking{
		ID:              1,
		UserID:          100,
		TourID:          200,
		FullName:        "John Doe",
		Phone:           "0123456789",
		Email:           "john@example.com",
		AdultCount:      2,
		ChildCount:      1,
		InfantCount:     0,
		Quantity:        3,
		TravelDate:      "2024-05-01",
		TotalAmount:     5000000,
		BookingCode:     "BOOK001",
		PaymentStatus:   "unpaid",
		PaymentID:       nil,
		PaymentDeadline: nil,
		Note:            "Test booking",
		Status:          "booked",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	result := booking.ToBookingWithPayment()

	if result.ID != booking.ID {
		t.Errorf("ID = %v, expected %v", result.ID, booking.ID)
	}
	if result.FullName != booking.FullName {
		t.Errorf("FullName = %v, expected %v", result.FullName, booking.FullName)
	}
	if result.PaymentStatus != booking.PaymentStatus {
		t.Errorf("PaymentStatus = %v, expected %v", result.PaymentStatus, booking.PaymentStatus)
	}
	if result.PaymentInfo != nil {
		t.Errorf("PaymentInfo should be nil initially")
	}
	if result.PaymentHistory != nil {
		t.Errorf("PaymentHistory should be nil initially")
	}
}

func TestBookingToBookingSummary(t *testing.T) {
	booking := &Booking{
		ID:              1,
		BookingCode:     "BOOK001",
		TravelDate:      "2024-05-01",
		TotalAmount:     5000000,
		PaymentStatus:   "unpaid",
		PaymentDeadline: nil,
		Status:          "booked",
		CreatedAt:       time.Now(),
		Tour: Tour{
			ID:   200,
			Name: "Ha Long Bay Tour",
		},
	}

	result := booking.ToBookingSummary()

	if result.ID != booking.ID {
		t.Errorf("ID = %v, expected %v", result.ID, booking.ID)
	}
	if result.BookingCode != booking.BookingCode {
		t.Errorf("BookingCode = %v, expected %v", result.BookingCode, booking.BookingCode)
	}
	if result.TourName != booking.Tour.Name {
		t.Errorf("TourName = %v, expected %v", result.TourName, booking.Tour.Name)
	}
	if result.CanPay != booking.CanPay() {
		t.Errorf("CanPay = %v, expected %v", result.CanPay, booking.CanPay())
	}
}

func TestUpdateBookingPaymentRequestValidate(t *testing.T) {
	tests := []struct {
		name      string
		request   *UpdateBookingPaymentRequest
		expectErr bool
	}{
		{
			name: "valid request",
			request: &UpdateBookingPaymentRequest{
				BookingID:     1,
				PaymentStatus: "paid",
			},
			expectErr: false,
		},
		{
			name: "invalid booking ID",
			request: &UpdateBookingPaymentRequest{
				BookingID:     0,
				PaymentStatus: "paid",
			},
			expectErr: true,
		},
		{
			name: "invalid payment status",
			request: &UpdateBookingPaymentRequest{
				BookingID:     1,
				PaymentStatus: "invalid_status",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.expectErr {
				t.Errorf("Validate() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestIsValidPaymentStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{"unpaid", true},
		{"processing", true},
		{"paid", true},
		{"failed", true},
		{"refunded", true},
		{"expired", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := IsValidPaymentStatus(tt.status)
			if result != tt.expected {
				t.Errorf("IsValidPaymentStatus(%s) = %v, expected %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestIsValidBookingStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{"booked", true},
		{"confirmed", true},
		{"cancelled", true},
		{"completed", true},
		{"invalid", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := IsValidBookingStatus(tt.status)
			if result != tt.expected {
				t.Errorf("IsValidBookingStatus(%s) = %v, expected %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetPaymentStatusDisplayName(t *testing.T) {
	tests := []struct {
		status   string
		expected string
	}{
		{"unpaid", "Chưa thanh toán"},
		{"processing", "Đang xử lý"},
		{"paid", "Đã thanh toán"},
		{"failed", "Thanh toán thất bại"},
		{"refunded", "Đã hoàn tiền"},
		{"expired", "Hết hạn thanh toán"},
		{"invalid", "Không xác định"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := GetPaymentStatusDisplayName(tt.status)
			if result != tt.expected {
				t.Errorf("GetPaymentStatusDisplayName(%s) = %v, expected %v", tt.status, result, tt.expected)
			}
		})
	}
}

func TestGetBookingStatusDisplayName(t *testing.T) {
	tests := []struct {
		status   string
		expected string
	}{
		{"booked", "Đã đặt"},
		{"confirmed", "Đã xác nhận"},
		{"cancelled", "Đã hủy"},
		{"completed", "Hoàn thành"},
		{"invalid", "Không xác định"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := GetBookingStatusDisplayName(tt.status)
			if result != tt.expected {
				t.Errorf("GetBookingStatusDisplayName(%s) = %v, expected %v", tt.status, result, tt.expected)
			}
		})
	}
}