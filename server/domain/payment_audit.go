package domain

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// Payment audit event types
const (
	EventPaymentInitiated   = "payment_initiated"
	EventPaymentCompleted   = "payment_completed"
	EventPaymentFailed      = "payment_failed"
	EventWebhookReceived    = "webhook_received"
	EventRefundProcessed    = "refund_processed"
	EventSecurityViolation  = "security_violation"
	EventStatusChanged      = "status_changed"
	EventRetryAttempted     = "retry_attempted"
)

// PaymentAuditLog - Model đại diện cho bảng payment_audit_logs trong database
// Lưu trữ tất cả các sự kiện liên quan đến thanh toán để audit và debug
type PaymentAuditLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PaymentID *uint     `json:"payment_id" gorm:"index"`
	BookingID *uint     `json:"booking_id" gorm:"index"`
	EventType string    `json:"event_type" gorm:"not null;size:50;index"`
	EventData json.RawMessage `json:"event_data" gorm:"type:jsonb"`
	UserID    *uint     `json:"user_id" gorm:"index"`
	IPAddress *net.IP   `json:"ip_address" gorm:"type:inet"`
	UserAgent *string   `json:"user_agent" gorm:"type:text"`
	Timestamp time.Time `json:"timestamp" gorm:"not null;default:CURRENT_TIMESTAMP;index"`

	// Relationships
	Payment *Payment `json:"payment,omitempty" gorm:"foreignKey:PaymentID"`
	Booking *Booking `json:"booking,omitempty" gorm:"foreignKey:BookingID"`
	User    *User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// PaymentAuditData - Struct chứa dữ liệu audit cho các sự kiện thanh toán
type PaymentAuditData struct {
	TransactionReference *string                `json:"transaction_reference,omitempty"`
	Amount               *int64                 `json:"amount,omitempty"`
	Currency             *string                `json:"currency,omitempty"`
	Status               *string                `json:"status,omitempty"`
	PreviousStatus       *string                `json:"previous_status,omitempty"`
	VNPayResponseCode    *string                `json:"vnpay_response_code,omitempty"`
	VNPayMessage         *string                `json:"vnpay_message,omitempty"`
	ErrorMessage         *string                `json:"error_message,omitempty"`
	RequestData          map[string]interface{} `json:"request_data,omitempty"`
	ResponseData         map[string]interface{} `json:"response_data,omitempty"`
	ProcessingTime       *int64                 `json:"processing_time_ms,omitempty"`
	RetryCount           *int                   `json:"retry_count,omitempty"`
	SecurityDetails      map[string]interface{} `json:"security_details,omitempty"`
}

// Validation methods

// ValidateEventType - Kiểm tra loại sự kiện hợp lệ
func (p *PaymentAuditLog) ValidateEventType() error {
	validEventTypes := []string{
		EventPaymentInitiated,
		EventPaymentCompleted,
		EventPaymentFailed,
		EventWebhookReceived,
		EventRefundProcessed,
		EventSecurityViolation,
		EventStatusChanged,
		EventRetryAttempted,
	}

	for _, validType := range validEventTypes {
		if p.EventType == validType {
			return nil
		}
	}

	return fmt.Errorf("invalid event type: %s", p.EventType)
}

// Helper methods

// SetEventData - Đặt dữ liệu sự kiện từ struct
func (p *PaymentAuditLog) SetEventData(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}
	
	p.EventData = jsonData
	return nil
}

// GetEventData - Lấy dữ liệu sự kiện và unmarshal vào struct
func (p *PaymentAuditLog) GetEventData(target interface{}) error {
	if p.EventData == nil {
		return nil
	}
	
	return json.Unmarshal(p.EventData, target)
}

// IsSecurityEvent - Kiểm tra có phải sự kiện bảo mật không
func (p *PaymentAuditLog) IsSecurityEvent() bool {
	return p.EventType == EventSecurityViolation
}

// IsPaymentEvent - Kiểm tra có phải sự kiện thanh toán không
func (p *PaymentAuditLog) IsPaymentEvent() bool {
	paymentEvents := []string{
		EventPaymentInitiated,
		EventPaymentCompleted,
		EventPaymentFailed,
		EventStatusChanged,
		EventRetryAttempted,
	}
	
	for _, event := range paymentEvents {
		if p.EventType == event {
			return true
		}
	}
	
	return false
}

// GetEventDisplayName - Lấy tên hiển thị của sự kiện
func (p *PaymentAuditLog) GetEventDisplayName() string {
	eventNames := map[string]string{
		EventPaymentInitiated:   "Khởi tạo thanh toán",
		EventPaymentCompleted:   "Thanh toán thành công",
		EventPaymentFailed:      "Thanh toán thất bại",
		EventWebhookReceived:    "Nhận webhook",
		EventRefundProcessed:    "Xử lý hoàn tiền",
		EventSecurityViolation:  "Vi phạm bảo mật",
		EventStatusChanged:      "Thay đổi trạng thái",
		EventRetryAttempted:     "Thử lại thanh toán",
	}
	
	if displayName, exists := eventNames[p.EventType]; exists {
		return displayName
	}
	
	return "Sự kiện không xác định"
}

// Factory functions

// NewPaymentAuditLog - Tạo audit log mới
func NewPaymentAuditLog(eventType string, paymentID, bookingID, userID *uint) *PaymentAuditLog {
	return &PaymentAuditLog{
		EventType: eventType,
		PaymentID: paymentID,
		BookingID: bookingID,
		UserID:    userID,
		Timestamp: time.Now(),
	}
}

// NewPaymentInitiatedLog - Tạo log cho sự kiện khởi tạo thanh toán
func NewPaymentInitiatedLog(payment *Payment, userID uint, ipAddress net.IP, userAgent string) *PaymentAuditLog {
	log := NewPaymentAuditLog(EventPaymentInitiated, &payment.ID, &payment.BookingID, &userID)
	log.IPAddress = &ipAddress
	log.UserAgent = &userAgent
	
	data := &PaymentAuditData{
		TransactionReference: &payment.TransactionReference,
		Amount:               &payment.Amount,
		Currency:             &payment.Currency,
		Status:               &payment.Status,
	}
	
	log.SetEventData(data)
	return log
}

// NewPaymentCompletedLog - Tạo log cho sự kiện thanh toán thành công
func NewPaymentCompletedLog(payment *Payment, vnpayResponseCode, vnpayMessage string) *PaymentAuditLog {
	log := NewPaymentAuditLog(EventPaymentCompleted, &payment.ID, &payment.BookingID, nil)
	
	data := &PaymentAuditData{
		TransactionReference: &payment.TransactionReference,
		Amount:               &payment.Amount,
		Status:               &payment.Status,
		VNPayResponseCode:    &vnpayResponseCode,
		VNPayMessage:         &vnpayMessage,
	}
	
	log.SetEventData(data)
	return log
}

// NewPaymentFailedLog - Tạo log cho sự kiện thanh toán thất bại
func NewPaymentFailedLog(payment *Payment, errorMessage, vnpayResponseCode string) *PaymentAuditLog {
	log := NewPaymentAuditLog(EventPaymentFailed, &payment.ID, &payment.BookingID, nil)
	
	data := &PaymentAuditData{
		TransactionReference: &payment.TransactionReference,
		Amount:               &payment.Amount,
		Status:               &payment.Status,
		ErrorMessage:         &errorMessage,
		VNPayResponseCode:    &vnpayResponseCode,
	}
	
	log.SetEventData(data)
	return log
}

// NewWebhookReceivedLog - Tạo log cho sự kiện nhận webhook
func NewWebhookReceivedLog(paymentID, bookingID *uint, requestData map[string]interface{}) *PaymentAuditLog {
	log := NewPaymentAuditLog(EventWebhookReceived, paymentID, bookingID, nil)
	
	data := &PaymentAuditData{
		RequestData: requestData,
	}
	
	log.SetEventData(data)
	return log
}

// NewSecurityViolationLog - Tạo log cho sự kiện vi phạm bảo mật
func NewSecurityViolationLog(paymentID, bookingID *uint, ipAddress net.IP, userAgent string, securityDetails map[string]interface{}) *PaymentAuditLog {
	log := NewPaymentAuditLog(EventSecurityViolation, paymentID, bookingID, nil)
	log.IPAddress = &ipAddress
	log.UserAgent = &userAgent
	
	data := &PaymentAuditData{
		SecurityDetails: securityDetails,
	}
	
	log.SetEventData(data)
	return log
}

// NewStatusChangedLog - Tạo log cho sự kiện thay đổi trạng thái
func NewStatusChangedLog(payment *Payment, previousStatus string) *PaymentAuditLog {
	log := NewPaymentAuditLog(EventStatusChanged, &payment.ID, &payment.BookingID, nil)
	
	data := &PaymentAuditData{
		TransactionReference: &payment.TransactionReference,
		Status:               &payment.Status,
		PreviousStatus:       &previousStatus,
	}
	
	log.SetEventData(data)
	return log
}

// PaymentAuditSummary - Struct để hiển thị thông tin tóm tắt audit log
type PaymentAuditSummary struct {
	ID                uint      `json:"id"`
	EventType         string    `json:"event_type"`
	EventDisplayName  string    `json:"event_display_name"`
	PaymentID         *uint     `json:"payment_id"`
	BookingID         *uint     `json:"booking_id"`
	UserID            *uint     `json:"user_id"`
	IPAddress         *string   `json:"ip_address"`
	IsSecurityEvent   bool      `json:"is_security_event"`
	IsPaymentEvent    bool      `json:"is_payment_event"`
	Timestamp         time.Time `json:"timestamp"`
}

// ToSummary - Chuyển đổi PaymentAuditLog thành PaymentAuditSummary
func (p *PaymentAuditLog) ToSummary() *PaymentAuditSummary {
	summary := &PaymentAuditSummary{
		ID:               p.ID,
		EventType:        p.EventType,
		EventDisplayName: p.GetEventDisplayName(),
		PaymentID:        p.PaymentID,
		BookingID:        p.BookingID,
		UserID:           p.UserID,
		IsSecurityEvent:  p.IsSecurityEvent(),
		IsPaymentEvent:   p.IsPaymentEvent(),
		Timestamp:        p.Timestamp,
	}
	
	if p.IPAddress != nil {
		ipStr := p.IPAddress.String()
		summary.IPAddress = &ipStr
	}
	
	return summary
}