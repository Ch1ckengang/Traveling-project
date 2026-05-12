package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Payment status constants
const (
	PaymentStatusPending    = "pending"
	PaymentStatusProcessing = "processing"
	PaymentStatusPaid       = "paid"
	PaymentStatusFailed     = "failed"
	PaymentStatusRefunded   = "refunded"
	PaymentStatusExpired    = "expired"
)

// Payment - Model đại diện cho bảng payments trong database
// Chứa thông tin về các giao dịch thanh toán qua VNPay
type Payment struct {
	ID                   uint      `json:"id" gorm:"primaryKey"`
	BookingID            uint      `json:"booking_id" gorm:"not null;index"`
	VNPayTransactionID   *string   `json:"vnpay_transaction_id" gorm:"unique;size:100"`
	TransactionReference string    `json:"transaction_reference" gorm:"uniqueIndex;not null;size:50"`
	Amount               int64     `json:"amount" gorm:"not null"` // Amount in VND cents (multiply by 100)
	Currency             string    `json:"currency" gorm:"not null;default:'VND';size:3"`
	PaymentMethod        *string   `json:"payment_method" gorm:"size:50"`
	Status               string    `json:"status" gorm:"not null;default:'pending';size:20;index"`
	VNPayResponseCode    *string   `json:"vnpay_response_code" gorm:"size:10"`
	VNPayMessage         *string   `json:"vnpay_message" gorm:"type:text"`
	SessionExpiresAt     *time.Time `json:"session_expires_at"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	DeletedAt            *time.Time `json:"deleted_at" gorm:"index"`

	// Relationships
	Booking *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
}

// Validation methods for business rules

// ValidateAmount - Kiểm tra số tiền thanh toán hợp lệ
func (p *Payment) ValidateAmount() error {
	if p.Amount <= 0 {
		return errors.New("payment amount must be positive")
	}
	
	// VNPay minimum amount is 5,000 VND (500,000 cents)
	if p.Amount < 500000 {
		return errors.New("payment amount must be at least 5,000 VND")
	}
	
	// VNPay maximum amount is 500,000,000 VND (50,000,000,000 cents)
	if p.Amount > 50000000000 {
		return errors.New("payment amount exceeds maximum limit of 500,000,000 VND")
	}
	
	return nil
}

// ValidateStatusTransition - Kiểm tra chuyển đổi trạng thái thanh toán hợp lệ
func (p *Payment) ValidateStatusTransition(newStatus string) error {
	if !isValidPaymentStatus(newStatus) {
		return fmt.Errorf("invalid payment status: %s", newStatus)
	}
	
	currentStatus := p.Status
	
	// Define valid status transitions
	validTransitions := map[string][]string{
		PaymentStatusPending: {
			PaymentStatusProcessing,
			PaymentStatusExpired,
			PaymentStatusFailed,
		},
		PaymentStatusProcessing: {
			PaymentStatusPaid,
			PaymentStatusFailed,
			PaymentStatusExpired,
		},
		PaymentStatusPaid: {
			PaymentStatusRefunded,
		},
		PaymentStatusFailed: {
			PaymentStatusPending, // Allow retry
			PaymentStatusProcessing,
		},
		PaymentStatusExpired: {
			PaymentStatusPending, // Allow retry
		},
		PaymentStatusRefunded: {
			// Terminal state - no transitions allowed
		},
	}
	
	allowedStatuses, exists := validTransitions[currentStatus]
	if !exists {
		return fmt.Errorf("no valid transitions defined for status: %s", currentStatus)
	}
	
	for _, allowedStatus := range allowedStatuses {
		if newStatus == allowedStatus {
			return nil
		}
	}
	
	return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
}

// ValidateTransactionReference - Kiểm tra định dạng mã tham chiếu giao dịch
func (p *Payment) ValidateTransactionReference() error {
	if p.TransactionReference == "" {
		return errors.New("transaction reference is required")
	}
	
	// Transaction reference format: PAY + YYYYMMDD + sequential number (e.g., PAY20240430001)
	pattern := `^PAY\d{8}\d{3,6}$`
	matched, err := regexp.MatchString(pattern, p.TransactionReference)
	if err != nil {
		return fmt.Errorf("error validating transaction reference format: %w", err)
	}
	
	if !matched {
		return errors.New("transaction reference must follow format: PAY + YYYYMMDD + sequential number")
	}
	
	return nil
}

// Helper methods for common operations

// IsExpired - Kiểm tra phiên thanh toán đã hết hạn chưa
func (p *Payment) IsExpired() bool {
	if p.SessionExpiresAt == nil {
		return false
	}
	return time.Now().After(*p.SessionExpiresAt)
}

// CanRetry - Kiểm tra có thể thử lại thanh toán không
func (p *Payment) CanRetry() bool {
	// Can retry if payment failed or expired, and not older than 24 hours
	if p.Status != PaymentStatusFailed && p.Status != PaymentStatusExpired {
		return false
	}
	
	// Check if payment is not older than 24 hours
	retryDeadline := p.CreatedAt.Add(24 * time.Hour)
	return time.Now().Before(retryDeadline)
}

// IsSuccessful - Kiểm tra thanh toán thành công
func (p *Payment) IsSuccessful() bool {
	return p.Status == PaymentStatusPaid
}

// IsPending - Kiểm tra thanh toán đang chờ xử lý
func (p *Payment) IsPending() bool {
	return p.Status == PaymentStatusPending || p.Status == PaymentStatusProcessing
}

// IsTerminal - Kiểm tra trạng thái thanh toán đã kết thúc
func (p *Payment) IsTerminal() bool {
	return p.Status == PaymentStatusPaid || p.Status == PaymentStatusRefunded
}

// GetAmountInVND - Lấy số tiền theo đơn vị VND (chia cho 100)
func (p *Payment) GetAmountInVND() float64 {
	return float64(p.Amount) / 100.0
}

// SetAmountFromVND - Đặt số tiền từ đơn vị VND (nhân với 100)
func (p *Payment) SetAmountFromVND(amountVND float64) error {
	if amountVND <= 0 {
		return errors.New("amount must be positive")
	}
	
	// Convert to cents and round to avoid floating point precision issues
	p.Amount = int64(amountVND * 100 + 0.5)
	
	return p.ValidateAmount()
}

// GetStatusDisplayName - Lấy tên hiển thị của trạng thái thanh toán
func (p *Payment) GetStatusDisplayName() string {
	statusNames := map[string]string{
		PaymentStatusPending:    "Chờ thanh toán",
		PaymentStatusProcessing: "Đang xử lý",
		PaymentStatusPaid:       "Đã thanh toán",
		PaymentStatusFailed:     "Thanh toán thất bại",
		PaymentStatusRefunded:   "Đã hoàn tiền",
		PaymentStatusExpired:    "Hết hạn",
	}
	
	if displayName, exists := statusNames[p.Status]; exists {
		return displayName
	}
	
	return "Không xác định"
}

// UpdateStatus - Cập nhật trạng thái thanh toán với validation
func (p *Payment) UpdateStatus(newStatus string) error {
	if err := p.ValidateStatusTransition(newStatus); err != nil {
		return err
	}
	
	p.Status = newStatus
	p.UpdatedAt = time.Now()
	
	return nil
}

// SetVNPayResponse - Đặt thông tin phản hồi từ VNPay
func (p *Payment) SetVNPayResponse(transactionID, responseCode, message string) {
	if transactionID != "" {
		p.VNPayTransactionID = &transactionID
	}
	if responseCode != "" {
		p.VNPayResponseCode = &responseCode
	}
	if message != "" {
		p.VNPayMessage = &message
	}
	p.UpdatedAt = time.Now()
}

// Helper functions

// isValidPaymentStatus - Kiểm tra trạng thái thanh toán hợp lệ
func isValidPaymentStatus(status string) bool {
	validStatuses := []string{
		PaymentStatusPending,
		PaymentStatusProcessing,
		PaymentStatusPaid,
		PaymentStatusFailed,
		PaymentStatusRefunded,
		PaymentStatusExpired,
	}
	
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	
	return false
}

// GenerateTransactionReference - Tạo mã tham chiếu giao dịch duy nhất
func GenerateTransactionReference() string {
	now := time.Now()
	dateStr := now.Format("20060102")
	
	// Generate a sequential number based on timestamp (microseconds)
	seqNum := now.UnixMicro() % 1000000 // Get last 6 digits of microseconds
	
	return fmt.Sprintf("PAY%s%06d", dateStr, seqNum)
}

// NewPayment - Tạo payment mới với các giá trị mặc định
func NewPayment(bookingID uint, amount int64) *Payment {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute) // 15 minutes expiration
	
	return &Payment{
		BookingID:            bookingID,
		TransactionReference: GenerateTransactionReference(),
		Amount:               amount,
		Currency:             "VND",
		Status:               PaymentStatusPending,
		SessionExpiresAt:     &expiresAt,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// PaymentSummary - Struct để hiển thị thông tin tóm tắt thanh toán
type PaymentSummary struct {
	ID                   uint      `json:"id"`
	TransactionReference string    `json:"transaction_reference"`
	Amount               float64   `json:"amount"` // Amount in VND
	Currency             string    `json:"currency"`
	Status               string    `json:"status"`
	StatusDisplayName    string    `json:"status_display_name"`
	PaymentMethod        *string   `json:"payment_method"`
	IsExpired            bool      `json:"is_expired"`
	CanRetry             bool      `json:"can_retry"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// ToSummary - Chuyển đổi Payment thành PaymentSummary
func (p *Payment) ToSummary() *PaymentSummary {
	return &PaymentSummary{
		ID:                   p.ID,
		TransactionReference: p.TransactionReference,
		Amount:               p.GetAmountInVND(),
		Currency:             p.Currency,
		Status:               p.Status,
		StatusDisplayName:    p.GetStatusDisplayName(),
		PaymentMethod:        p.PaymentMethod,
		IsExpired:            p.IsExpired(),
		CanRetry:             p.CanRetry(),
		CreatedAt:            p.CreatedAt,
		UpdatedAt:            p.UpdatedAt,
	}
}