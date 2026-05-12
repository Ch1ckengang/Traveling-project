package domain

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentAuditLog_ValidateEventType(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		wantErr   bool
	}{
		{
			name:      "valid payment initiated event",
			eventType: EventPaymentInitiated,
			wantErr:   false,
		},
		{
			name:      "valid payment completed event",
			eventType: EventPaymentCompleted,
			wantErr:   false,
		},
		{
			name:      "valid payment failed event",
			eventType: EventPaymentFailed,
			wantErr:   false,
		},
		{
			name:      "valid webhook received event",
			eventType: EventWebhookReceived,
			wantErr:   false,
		},
		{
			name:      "valid refund processed event",
			eventType: EventRefundProcessed,
			wantErr:   false,
		},
		{
			name:      "valid security violation event",
			eventType: EventSecurityViolation,
			wantErr:   false,
		},
		{
			name:      "valid status changed event",
			eventType: EventStatusChanged,
			wantErr:   false,
		},
		{
			name:      "valid retry attempted event",
			eventType: EventRetryAttempted,
			wantErr:   false,
		},
		{
			name:      "invalid event type",
			eventType: "invalid_event",
			wantErr:   true,
		},
		{
			name:      "empty event type",
			eventType: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &PaymentAuditLog{
				EventType: tt.eventType,
			}

			err := log.ValidateEventType()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPaymentAuditLog_SetAndGetEventData(t *testing.T) {
	log := &PaymentAuditLog{}

	// Test data
	originalData := &PaymentAuditData{
		TransactionReference: stringPtr("PAY20240430001"),
		Amount:               int64Ptr(2500000),
		Currency:             stringPtr("VND"),
		Status:               stringPtr("completed"),
		VNPayResponseCode:    stringPtr("00"),
		VNPayMessage:         stringPtr("Success"),
	}

	// Set event data
	err := log.SetEventData(originalData)
	require.NoError(t, err)
	assert.NotNil(t, log.EventData)

	// Get event data
	var retrievedData PaymentAuditData
	err = log.GetEventData(&retrievedData)
	require.NoError(t, err)

	// Verify data integrity
	assert.Equal(t, *originalData.TransactionReference, *retrievedData.TransactionReference)
	assert.Equal(t, *originalData.Amount, *retrievedData.Amount)
	assert.Equal(t, *originalData.Currency, *retrievedData.Currency)
	assert.Equal(t, *originalData.Status, *retrievedData.Status)
	assert.Equal(t, *originalData.VNPayResponseCode, *retrievedData.VNPayResponseCode)
	assert.Equal(t, *originalData.VNPayMessage, *retrievedData.VNPayMessage)
}

func TestPaymentAuditLog_SetEventData_InvalidData(t *testing.T) {
	log := &PaymentAuditLog{}

	// Test with data that cannot be marshaled (circular reference)
	type circularRef struct {
		Self *circularRef `json:"self"`
	}
	circular := &circularRef{}
	circular.Self = circular

	err := log.SetEventData(circular)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to marshal event data")
}

func TestPaymentAuditLog_GetEventData_EmptyData(t *testing.T) {
	log := &PaymentAuditLog{
		EventData: nil,
	}

	var data PaymentAuditData
	err := log.GetEventData(&data)
	assert.NoError(t, err)
}

func TestPaymentAuditLog_IsSecurityEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		expected  bool
	}{
		{
			name:      "security violation event",
			eventType: EventSecurityViolation,
			expected:  true,
		},
		{
			name:      "payment initiated event",
			eventType: EventPaymentInitiated,
			expected:  false,
		},
		{
			name:      "webhook received event",
			eventType: EventWebhookReceived,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &PaymentAuditLog{
				EventType: tt.eventType,
			}

			result := log.IsSecurityEvent()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaymentAuditLog_IsPaymentEvent(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		expected  bool
	}{
		{
			name:      "payment initiated event",
			eventType: EventPaymentInitiated,
			expected:  true,
		},
		{
			name:      "payment completed event",
			eventType: EventPaymentCompleted,
			expected:  true,
		},
		{
			name:      "payment failed event",
			eventType: EventPaymentFailed,
			expected:  true,
		},
		{
			name:      "status changed event",
			eventType: EventStatusChanged,
			expected:  true,
		},
		{
			name:      "retry attempted event",
			eventType: EventRetryAttempted,
			expected:  true,
		},
		{
			name:      "webhook received event",
			eventType: EventWebhookReceived,
			expected:  false,
		},
		{
			name:      "refund processed event",
			eventType: EventRefundProcessed,
			expected:  false,
		},
		{
			name:      "security violation event",
			eventType: EventSecurityViolation,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &PaymentAuditLog{
				EventType: tt.eventType,
			}

			result := log.IsPaymentEvent()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPaymentAuditLog_GetEventDisplayName(t *testing.T) {
	tests := []struct {
		name      string
		eventType string
		expected  string
	}{
		{
			name:      "payment initiated event",
			eventType: EventPaymentInitiated,
			expected:  "Khởi tạo thanh toán",
		},
		{
			name:      "payment completed event",
			eventType: EventPaymentCompleted,
			expected:  "Thanh toán thành công",
		},
		{
			name:      "payment failed event",
			eventType: EventPaymentFailed,
			expected:  "Thanh toán thất bại",
		},
		{
			name:      "webhook received event",
			eventType: EventWebhookReceived,
			expected:  "Nhận webhook",
		},
		{
			name:      "refund processed event",
			eventType: EventRefundProcessed,
			expected:  "Xử lý hoàn tiền",
		},
		{
			name:      "security violation event",
			eventType: EventSecurityViolation,
			expected:  "Vi phạm bảo mật",
		},
		{
			name:      "status changed event",
			eventType: EventStatusChanged,
			expected:  "Thay đổi trạng thái",
		},
		{
			name:      "retry attempted event",
			eventType: EventRetryAttempted,
			expected:  "Thử lại thanh toán",
		},
		{
			name:      "unknown event type",
			eventType: "unknown_event",
			expected:  "Sự kiện không xác định",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &PaymentAuditLog{
				EventType: tt.eventType,
			}

			result := log.GetEventDisplayName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewPaymentAuditLog(t *testing.T) {
	eventType := EventPaymentInitiated
	paymentID := uint(123)
	bookingID := uint(456)
	userID := uint(789)

	log := NewPaymentAuditLog(eventType, &paymentID, &bookingID, &userID)

	assert.Equal(t, eventType, log.EventType)
	assert.Equal(t, &paymentID, log.PaymentID)
	assert.Equal(t, &bookingID, log.BookingID)
	assert.Equal(t, &userID, log.UserID)
	assert.WithinDuration(t, time.Now(), log.Timestamp, time.Second)
}

func TestNewPaymentInitiatedLog(t *testing.T) {
	payment := &Payment{
		ID:                   123,
		BookingID:            456,
		TransactionReference: "PAY20240430001",
		Amount:               2500000,
		Currency:             "VND",
		Status:               "pending",
	}
	userID := uint(789)
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Mozilla/5.0 Test Browser"

	log := NewPaymentInitiatedLog(payment, userID, ipAddress, userAgent)

	assert.Equal(t, EventPaymentInitiated, log.EventType)
	assert.Equal(t, &payment.ID, log.PaymentID)
	assert.Equal(t, &payment.BookingID, log.BookingID)
	assert.Equal(t, &userID, log.UserID)
	assert.Equal(t, &ipAddress, log.IPAddress)
	assert.Equal(t, &userAgent, log.UserAgent)

	// Verify event data
	var eventData PaymentAuditData
	err := log.GetEventData(&eventData)
	require.NoError(t, err)
	assert.Equal(t, payment.TransactionReference, *eventData.TransactionReference)
	assert.Equal(t, payment.Amount, *eventData.Amount)
	assert.Equal(t, payment.Currency, *eventData.Currency)
	assert.Equal(t, payment.Status, *eventData.Status)
}

func TestNewPaymentCompletedLog(t *testing.T) {
	payment := &Payment{
		ID:                   123,
		BookingID:            456,
		TransactionReference: "PAY20240430001",
		Amount:               2500000,
		Status:               "completed",
	}
	vnpayResponseCode := "00"
	vnpayMessage := "Success"

	log := NewPaymentCompletedLog(payment, vnpayResponseCode, vnpayMessage)

	assert.Equal(t, EventPaymentCompleted, log.EventType)
	assert.Equal(t, &payment.ID, log.PaymentID)
	assert.Equal(t, &payment.BookingID, log.BookingID)

	// Verify event data
	var eventData PaymentAuditData
	err := log.GetEventData(&eventData)
	require.NoError(t, err)
	assert.Equal(t, payment.TransactionReference, *eventData.TransactionReference)
	assert.Equal(t, payment.Amount, *eventData.Amount)
	assert.Equal(t, payment.Status, *eventData.Status)
	assert.Equal(t, vnpayResponseCode, *eventData.VNPayResponseCode)
	assert.Equal(t, vnpayMessage, *eventData.VNPayMessage)
}

func TestNewPaymentFailedLog(t *testing.T) {
	payment := &Payment{
		ID:                   123,
		BookingID:            456,
		TransactionReference: "PAY20240430001",
		Amount:               2500000,
		Status:               "failed",
	}
	errorMessage := "Payment processing failed"
	vnpayResponseCode := "24"

	log := NewPaymentFailedLog(payment, errorMessage, vnpayResponseCode)

	assert.Equal(t, EventPaymentFailed, log.EventType)
	assert.Equal(t, &payment.ID, log.PaymentID)
	assert.Equal(t, &payment.BookingID, log.BookingID)

	// Verify event data
	var eventData PaymentAuditData
	err := log.GetEventData(&eventData)
	require.NoError(t, err)
	assert.Equal(t, payment.TransactionReference, *eventData.TransactionReference)
	assert.Equal(t, payment.Amount, *eventData.Amount)
	assert.Equal(t, payment.Status, *eventData.Status)
	assert.Equal(t, errorMessage, *eventData.ErrorMessage)
	assert.Equal(t, vnpayResponseCode, *eventData.VNPayResponseCode)
}

func TestNewWebhookReceivedLog(t *testing.T) {
	paymentID := uint(123)
	bookingID := uint(456)
	requestData := map[string]interface{}{
		"vnp_TxnRef":      "PAY20240430001",
		"vnp_ResponseCode": "00",
		"vnp_Amount":      "250000000",
	}

	log := NewWebhookReceivedLog(&paymentID, &bookingID, requestData)

	assert.Equal(t, EventWebhookReceived, log.EventType)
	assert.Equal(t, &paymentID, log.PaymentID)
	assert.Equal(t, &bookingID, log.BookingID)

	// Verify event data
	var eventData PaymentAuditData
	err := log.GetEventData(&eventData)
	require.NoError(t, err)
	assert.Equal(t, requestData, eventData.RequestData)
}

func TestNewSecurityViolationLog(t *testing.T) {
	paymentID := uint(123)
	bookingID := uint(456)
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Malicious Bot"
	securityDetails := map[string]interface{}{
		"violation_type": "invalid_signature",
		"attempted_signature": "invalid_hash",
		"expected_signature": "valid_hash",
	}

	log := NewSecurityViolationLog(&paymentID, &bookingID, ipAddress, userAgent, securityDetails)

	assert.Equal(t, EventSecurityViolation, log.EventType)
	assert.Equal(t, &paymentID, log.PaymentID)
	assert.Equal(t, &bookingID, log.BookingID)
	assert.Equal(t, &ipAddress, log.IPAddress)
	assert.Equal(t, &userAgent, log.UserAgent)

	// Verify event data
	var eventData PaymentAuditData
	err := log.GetEventData(&eventData)
	require.NoError(t, err)
	assert.Equal(t, securityDetails, eventData.SecurityDetails)
}

func TestNewStatusChangedLog(t *testing.T) {
	payment := &Payment{
		ID:                   123,
		BookingID:            456,
		TransactionReference: "PAY20240430001",
		Status:               "completed",
	}
	previousStatus := "pending"

	log := NewStatusChangedLog(payment, previousStatus)

	assert.Equal(t, EventStatusChanged, log.EventType)
	assert.Equal(t, &payment.ID, log.PaymentID)
	assert.Equal(t, &payment.BookingID, log.BookingID)

	// Verify event data
	var eventData PaymentAuditData
	err := log.GetEventData(&eventData)
	require.NoError(t, err)
	assert.Equal(t, payment.TransactionReference, *eventData.TransactionReference)
	assert.Equal(t, payment.Status, *eventData.Status)
	assert.Equal(t, previousStatus, *eventData.PreviousStatus)
}

func TestPaymentAuditLog_ToSummary(t *testing.T) {
	ipAddress := net.ParseIP("192.168.1.1")
	userAgent := "Test Browser"
	timestamp := time.Now()

	log := &PaymentAuditLog{
		ID:        123,
		EventType: EventPaymentInitiated,
		PaymentID: uintPtr(456),
		BookingID: uintPtr(789),
		UserID:    uintPtr(101112),
		IPAddress: &ipAddress,
		UserAgent: &userAgent,
		Timestamp: timestamp,
	}

	summary := log.ToSummary()

	assert.Equal(t, log.ID, summary.ID)
	assert.Equal(t, log.EventType, summary.EventType)
	assert.Equal(t, "Khởi tạo thanh toán", summary.EventDisplayName)
	assert.Equal(t, log.PaymentID, summary.PaymentID)
	assert.Equal(t, log.BookingID, summary.BookingID)
	assert.Equal(t, log.UserID, summary.UserID)
	assert.Equal(t, ipAddress.String(), *summary.IPAddress)
	assert.False(t, summary.IsSecurityEvent)
	assert.True(t, summary.IsPaymentEvent)
	assert.Equal(t, timestamp, summary.Timestamp)
}

func TestPaymentAuditLog_ToSummary_NilIPAddress(t *testing.T) {
	log := &PaymentAuditLog{
		ID:        123,
		EventType: EventPaymentInitiated,
		IPAddress: nil,
	}

	summary := log.ToSummary()

	assert.Nil(t, summary.IPAddress)
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func int64Ptr(i int64) *int64 {
	return &i
}

func uintPtr(u uint) *uint {
	return &u
}