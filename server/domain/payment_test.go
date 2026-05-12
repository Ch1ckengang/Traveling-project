package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPayment_ValidateAmount(t *testing.T) {
	tests := []struct {
		name    string
		payment *Payment
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid amount",
			payment: &Payment{
				Amount: 1000000, // 10,000 VND
			},
			wantErr: false,
		},
		{
			name: "minimum valid amount",
			payment: &Payment{
				Amount: 500000, // 5,000 VND (minimum)
			},
			wantErr: false,
		},
		{
			name: "maximum valid amount",
			payment: &Payment{
				Amount: 50000000000, // 500,000,000 VND (maximum)
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			payment: &Payment{
				Amount: 0,
			},
			wantErr: true,
			errMsg:  "payment amount must be positive",
		},
		{
			name: "negative amount",
			payment: &Payment{
				Amount: -100000,
			},
			wantErr: true,
			errMsg:  "payment amount must be positive",
		},
		{
			name: "amount below minimum",
			payment: &Payment{
				Amount: 499999, // Below 5,000 VND
			},
			wantErr: true,
			errMsg:  "payment amount must be at least 5,000 VND",
		},
		{
			name: "amount above maximum",
			payment: &Payment{
				Amount: 50000000001, // Above 500,000,000 VND
			},
			wantErr: true,
			errMsg:  "payment amount exceeds maximum limit of 500,000,000 VND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payment.ValidateAmount()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPayment_ValidateStatusTransition(t *testing.T) {
	tests := []struct {
		name           string
		currentStatus  string
		newStatus      string
		wantErr        bool
		errMsg         string
	}{
		// Valid transitions from pending
		{
			name:          "pending to processing",
			currentStatus: PaymentStatusPending,
			newStatus:     PaymentStatusProcessing,
			wantErr:       false,
		},
		{
			name:          "pending to expired",
			currentStatus: PaymentStatusPending,
			newStatus:     PaymentStatusExpired,
			wantErr:       false,
		},
		{
			name:          "pending to failed",
			currentStatus: PaymentStatusPending,
			newStatus:     PaymentStatusFailed,
			wantErr:       false,
		},
		// Valid transitions from processing
		{
			name:          "processing to paid",
			currentStatus: PaymentStatusProcessing,
			newStatus:     PaymentStatusPaid,
			wantErr:       false,
		},
		{
			name:          "processing to failed",
			currentStatus: PaymentStatusProcessing,
			newStatus:     PaymentStatusFailed,
			wantErr:       false,
		},
		{
			name:          "processing to expired",
			currentStatus: PaymentStatusProcessing,
			newStatus:     PaymentStatusExpired,
			wantErr:       false,
		},
		// Valid transitions from paid
		{
			name:          "paid to refunded",
			currentStatus: PaymentStatusPaid,
			newStatus:     PaymentStatusRefunded,
			wantErr:       false,
		},
		// Valid transitions from failed
		{
			name:          "failed to pending (retry)",
			currentStatus: PaymentStatusFailed,
			newStatus:     PaymentStatusPending,
			wantErr:       false,
		},
		{
			name:          "failed to processing (retry)",
			currentStatus: PaymentStatusFailed,
			newStatus:     PaymentStatusProcessing,
			wantErr:       false,
		},
		// Valid transitions from expired
		{
			name:          "expired to pending (retry)",
			currentStatus: PaymentStatusExpired,
			newStatus:     PaymentStatusPending,
			wantErr:       false,
		},
		// Invalid transitions
		{
			name:          "pending to paid (skip processing)",
			currentStatus: PaymentStatusPending,
			newStatus:     PaymentStatusPaid,
			wantErr:       true,
			errMsg:        "invalid status transition from pending to paid",
		},
		{
			name:          "paid to pending (invalid reverse)",
			currentStatus: PaymentStatusPaid,
			newStatus:     PaymentStatusPending,
			wantErr:       true,
			errMsg:        "invalid status transition from paid to pending",
		},
		{
			name:          "refunded to any status (terminal)",
			currentStatus: PaymentStatusRefunded,
			newStatus:     PaymentStatusPending,
			wantErr:       true,
			errMsg:        "invalid status transition from refunded to pending",
		},
		{
			name:          "invalid new status",
			currentStatus: PaymentStatusPending,
			newStatus:     "invalid_status",
			wantErr:       true,
			errMsg:        "invalid payment status: invalid_status",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := &Payment{Status: tt.currentStatus}
			err := payment.ValidateStatusTransition(tt.newStatus)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPayment_ValidateTransactionReference(t *testing.T) {
	tests := []struct {
		name    string
		payment *Payment
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid transaction reference",
			payment: &Payment{
				TransactionReference: "PAY20240430001",
			},
			wantErr: false,
		},
		{
			name: "valid transaction reference with longer sequence",
			payment: &Payment{
				TransactionReference: "PAY20240430123456",
			},
			wantErr: false,
		},
		{
			name: "empty transaction reference",
			payment: &Payment{
				TransactionReference: "",
			},
			wantErr: true,
			errMsg:  "transaction reference is required",
		},
		{
			name: "invalid format - missing PAY prefix",
			payment: &Payment{
				TransactionReference: "20240430001",
			},
			wantErr: true,
			errMsg:  "transaction reference must follow format: PAY + YYYYMMDD + sequential number",
		},
		{
			name: "invalid format - wrong date format",
			payment: &Payment{
				TransactionReference: "PAY2024043001",
			},
			wantErr: true,
			errMsg:  "transaction reference must follow format: PAY + YYYYMMDD + sequential number",
		},
		{
			name: "invalid format - missing sequence number",
			payment: &Payment{
				TransactionReference: "PAY20240430",
			},
			wantErr: true,
			errMsg:  "transaction reference must follow format: PAY + YYYYMMDD + sequential number",
		},
		{
			name: "invalid format - sequence too short",
			payment: &Payment{
				TransactionReference: "PAY2024043001",
			},
			wantErr: true,
			errMsg:  "transaction reference must follow format: PAY + YYYYMMDD + sequential number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payment.ValidateTransactionReference()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPayment_IsExpired(t *testing.T) {
	tests := []struct {
		name     string
		payment  *Payment
		expected bool
	}{
		{
			name: "no expiration time",
			payment: &Payment{
				SessionExpiresAt: nil,
			},
			expected: false,
		},
		{
			name: "expired session",
			payment: &Payment{
				SessionExpiresAt: func() *time.Time {
					past := time.Now().Add(-1 * time.Hour)
					return &past
				}(),
			},
			expected: true,
		},
		{
			name: "not expired session",
			payment: &Payment{
				SessionExpiresAt: func() *time.Time {
					future := time.Now().Add(1 * time.Hour)
					return &future
				}(),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_CanRetry(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		payment  *Payment
		expected bool
	}{
		{
			name: "can retry - failed payment within 24 hours",
			payment: &Payment{
				Status:    PaymentStatusFailed,
				CreatedAt: now.Add(-12 * time.Hour),
			},
			expected: true,
		},
		{
			name: "can retry - expired payment within 24 hours",
			payment: &Payment{
				Status:    PaymentStatusExpired,
				CreatedAt: now.Add(-12 * time.Hour),
			},
			expected: true,
		},
		{
			name: "cannot retry - failed payment older than 24 hours",
			payment: &Payment{
				Status:    PaymentStatusFailed,
				CreatedAt: now.Add(-25 * time.Hour),
			},
			expected: false,
		},
		{
			name: "cannot retry - expired payment older than 24 hours",
			payment: &Payment{
				Status:    PaymentStatusExpired,
				CreatedAt: now.Add(-25 * time.Hour),
			},
			expected: false,
		},
		{
			name: "cannot retry - paid payment",
			payment: &Payment{
				Status:    PaymentStatusPaid,
				CreatedAt: now.Add(-12 * time.Hour),
			},
			expected: false,
		},
		{
			name: "cannot retry - pending payment",
			payment: &Payment{
				Status:    PaymentStatusPending,
				CreatedAt: now.Add(-12 * time.Hour),
			},
			expected: false,
		},
		{
			name: "cannot retry - processing payment",
			payment: &Payment{
				Status:    PaymentStatusProcessing,
				CreatedAt: now.Add(-12 * time.Hour),
			},
			expected: false,
		},
		{
			name: "cannot retry - refunded payment",
			payment: &Payment{
				Status:    PaymentStatusRefunded,
				CreatedAt: now.Add(-12 * time.Hour),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.CanRetry()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_IsSuccessful(t *testing.T) {
	tests := []struct {
		name     string
		payment  *Payment
		expected bool
	}{
		{
			name: "successful payment",
			payment: &Payment{
				Status: PaymentStatusPaid,
			},
			expected: true,
		},
		{
			name: "pending payment",
			payment: &Payment{
				Status: PaymentStatusPending,
			},
			expected: false,
		},
		{
			name: "failed payment",
			payment: &Payment{
				Status: PaymentStatusFailed,
			},
			expected: false,
		},
		{
			name: "refunded payment",
			payment: &Payment{
				Status: PaymentStatusRefunded,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.IsSuccessful()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_IsPending(t *testing.T) {
	tests := []struct {
		name     string
		payment  *Payment
		expected bool
	}{
		{
			name: "pending payment",
			payment: &Payment{
				Status: PaymentStatusPending,
			},
			expected: true,
		},
		{
			name: "processing payment",
			payment: &Payment{
				Status: PaymentStatusProcessing,
			},
			expected: true,
		},
		{
			name: "paid payment",
			payment: &Payment{
				Status: PaymentStatusPaid,
			},
			expected: false,
		},
		{
			name: "failed payment",
			payment: &Payment{
				Status: PaymentStatusFailed,
			},
			expected: false,
		},
		{
			name: "expired payment",
			payment: &Payment{
				Status: PaymentStatusExpired,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.IsPending()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_IsTerminal(t *testing.T) {
	tests := []struct {
		name     string
		payment  *Payment
		expected bool
	}{
		{
			name: "paid payment (terminal)",
			payment: &Payment{
				Status: PaymentStatusPaid,
			},
			expected: true,
		},
		{
			name: "refunded payment (terminal)",
			payment: &Payment{
				Status: PaymentStatusRefunded,
			},
			expected: true,
		},
		{
			name: "pending payment (not terminal)",
			payment: &Payment{
				Status: PaymentStatusPending,
			},
			expected: false,
		},
		{
			name: "processing payment (not terminal)",
			payment: &Payment{
				Status: PaymentStatusProcessing,
			},
			expected: false,
		},
		{
			name: "failed payment (not terminal)",
			payment: &Payment{
				Status: PaymentStatusFailed,
			},
			expected: false,
		},
		{
			name: "expired payment (not terminal)",
			payment: &Payment{
				Status: PaymentStatusExpired,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.IsTerminal()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_GetAmountInVND(t *testing.T) {
	tests := []struct {
		name     string
		payment  *Payment
		expected float64
	}{
		{
			name: "convert cents to VND",
			payment: &Payment{
				Amount: 2500000, // 25,000 VND in cents
			},
			expected: 25000.0,
		},
		{
			name: "zero amount",
			payment: &Payment{
				Amount: 0,
			},
			expected: 0.0,
		},
		{
			name: "large amount",
			payment: &Payment{
				Amount: 50000000000, // 500,000,000 VND in cents
			},
			expected: 500000000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.GetAmountInVND()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_SetAmountFromVND(t *testing.T) {
	tests := []struct {
		name      string
		payment   *Payment
		amountVND float64
		wantErr   bool
		expected  int64
	}{
		{
			name:      "valid VND amount",
			payment:   &Payment{},
			amountVND: 25000.0,
			wantErr:   false,
			expected:  2500000, // 25,000 VND in cents
		},
		{
			name:      "VND amount with decimals",
			payment:   &Payment{},
			amountVND: 25000.50,
			wantErr:   false,
			expected:  2500050, // 25,000.50 VND in cents
		},
		{
			name:      "zero VND amount",
			payment:   &Payment{},
			amountVND: 0.0,
			wantErr:   true,
		},
		{
			name:      "negative VND amount",
			payment:   &Payment{},
			amountVND: -1000.0,
			wantErr:   true,
		},
		{
			name:      "amount below minimum after conversion",
			payment:   &Payment{},
			amountVND: 4999.0, // Below 5,000 VND minimum
			wantErr:   true,
		},
		{
			name:      "amount above maximum after conversion",
			payment:   &Payment{},
			amountVND: 500000001.0, // Above 500,000,000 VND maximum
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payment.SetAmountFromVND(tt.amountVND)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, tt.payment.Amount)
			}
		})
	}
}

func TestPayment_GetStatusDisplayName(t *testing.T) {
	tests := []struct {
		name     string
		payment  *Payment
		expected string
	}{
		{
			name: "pending status",
			payment: &Payment{
				Status: PaymentStatusPending,
			},
			expected: "Chờ thanh toán",
		},
		{
			name: "processing status",
			payment: &Payment{
				Status: PaymentStatusProcessing,
			},
			expected: "Đang xử lý",
		},
		{
			name: "paid status",
			payment: &Payment{
				Status: PaymentStatusPaid,
			},
			expected: "Đã thanh toán",
		},
		{
			name: "failed status",
			payment: &Payment{
				Status: PaymentStatusFailed,
			},
			expected: "Thanh toán thất bại",
		},
		{
			name: "refunded status",
			payment: &Payment{
				Status: PaymentStatusRefunded,
			},
			expected: "Đã hoàn tiền",
		},
		{
			name: "expired status",
			payment: &Payment{
				Status: PaymentStatusExpired,
			},
			expected: "Hết hạn",
		},
		{
			name: "unknown status",
			payment: &Payment{
				Status: "unknown_status",
			},
			expected: "Không xác định",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.payment.GetStatusDisplayName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_UpdateStatus(t *testing.T) {
	tests := []struct {
		name           string
		payment        *Payment
		newStatus      string
		wantErr        bool
		expectedStatus string
	}{
		{
			name: "valid status update",
			payment: &Payment{
				Status:    PaymentStatusPending,
				UpdatedAt: time.Now().Add(-1 * time.Hour),
			},
			newStatus:      PaymentStatusProcessing,
			wantErr:        false,
			expectedStatus: PaymentStatusProcessing,
		},
		{
			name: "invalid status update",
			payment: &Payment{
				Status: PaymentStatusPaid,
			},
			newStatus: PaymentStatusPending,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := tt.payment.UpdatedAt
			err := tt.payment.UpdateStatus(tt.newStatus)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, tt.payment.Status)
				assert.True(t, tt.payment.UpdatedAt.After(oldUpdatedAt))
			}
		})
	}
}

func TestPayment_SetVNPayResponse(t *testing.T) {
	payment := &Payment{
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}
	oldUpdatedAt := payment.UpdatedAt

	transactionID := "VNP123456789"
	responseCode := "00"
	message := "Success"

	payment.SetVNPayResponse(transactionID, responseCode, message)

	assert.Equal(t, &transactionID, payment.VNPayTransactionID)
	assert.Equal(t, &responseCode, payment.VNPayResponseCode)
	assert.Equal(t, &message, payment.VNPayMessage)
	assert.True(t, payment.UpdatedAt.After(oldUpdatedAt))
}

func TestPayment_SetVNPayResponse_EmptyValues(t *testing.T) {
	payment := &Payment{}

	payment.SetVNPayResponse("", "", "")

	assert.Nil(t, payment.VNPayTransactionID)
	assert.Nil(t, payment.VNPayResponseCode)
	assert.Nil(t, payment.VNPayMessage)
}

func TestIsValidPaymentStatusHelper(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{PaymentStatusPending, true},
		{PaymentStatusProcessing, true},
		{PaymentStatusPaid, true},
		{PaymentStatusFailed, true},
		{PaymentStatusRefunded, true},
		{PaymentStatusExpired, true},
		{"invalid_status", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := isValidPaymentStatus(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGenerateTransactionReference(t *testing.T) {
	// Generate multiple references to test uniqueness
	refs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		ref := GenerateTransactionReference()
		
		// Check format
		assert.Regexp(t, `^PAY\d{8}\d{6}$`, ref)
		
		// Check uniqueness
		assert.False(t, refs[ref], "Generated duplicate transaction reference: %s", ref)
		refs[ref] = true
		
		// Check date part
		today := time.Now().Format("20060102")
		assert.Contains(t, ref, today)
	}
}

func TestNewPayment(t *testing.T) {
	bookingID := uint(123)
	amount := int64(2500000)

	payment := NewPayment(bookingID, amount)

	assert.Equal(t, bookingID, payment.BookingID)
	assert.Equal(t, amount, payment.Amount)
	assert.Equal(t, "VND", payment.Currency)
	assert.Equal(t, PaymentStatusPending, payment.Status)
	assert.NotEmpty(t, payment.TransactionReference)
	assert.Regexp(t, `^PAY\d{8}\d{6}$`, payment.TransactionReference)
	assert.NotNil(t, payment.SessionExpiresAt)
	assert.True(t, payment.SessionExpiresAt.After(time.Now()))
	assert.True(t, payment.SessionExpiresAt.Before(time.Now().Add(16*time.Minute)))
	assert.WithinDuration(t, time.Now(), payment.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), payment.UpdatedAt, time.Second)
}

func TestPayment_ToSummary(t *testing.T) {
	now := time.Now()
	expiresAt := now.Add(15 * time.Minute)
	paymentMethod := "ATM"

	payment := &Payment{
		ID:                   123,
		TransactionReference: "PAY20240430001",
		Amount:               2500000,
		Currency:             "VND",
		Status:               PaymentStatusPending,
		PaymentMethod:        &paymentMethod,
		SessionExpiresAt:     &expiresAt,
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	summary := payment.ToSummary()

	assert.Equal(t, payment.ID, summary.ID)
	assert.Equal(t, payment.TransactionReference, summary.TransactionReference)
	assert.Equal(t, 25000.0, summary.Amount) // Converted to VND
	assert.Equal(t, payment.Currency, summary.Currency)
	assert.Equal(t, payment.Status, summary.Status)
	assert.Equal(t, "Chờ thanh toán", summary.StatusDisplayName)
	assert.Equal(t, payment.PaymentMethod, summary.PaymentMethod)
	assert.Equal(t, payment.IsExpired(), summary.IsExpired)
	assert.Equal(t, payment.CanRetry(), summary.CanRetry)
	assert.Equal(t, payment.CreatedAt, summary.CreatedAt)
	assert.Equal(t, payment.UpdatedAt, summary.UpdatedAt)
}

// Edge case tests for business logic

func TestPayment_ValidateAmount_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		amount  int64
		wantErr bool
	}{
		{
			name:    "exactly minimum amount",
			amount:  500000,
			wantErr: false,
		},
		{
			name:    "one cent below minimum",
			amount:  499999,
			wantErr: true,
		},
		{
			name:    "exactly maximum amount",
			amount:  50000000000,
			wantErr: false,
		},
		{
			name:    "one cent above maximum",
			amount:  50000000001,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := &Payment{Amount: tt.amount}
			err := payment.ValidateAmount()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPayment_CanRetry_EdgeCases(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name      string
		status    string
		createdAt time.Time
		expected  bool
	}{
		{
			name:      "failed payment exactly 24 hours ago",
			status:    PaymentStatusFailed,
			createdAt: now.Add(-24 * time.Hour),
			expected:  false, // Should be false as it's exactly at the boundary
		},
		{
			name:      "failed payment 23 hours 59 minutes ago",
			status:    PaymentStatusFailed,
			createdAt: now.Add(-23*time.Hour - 59*time.Minute),
			expected:  true,
		},
		{
			name:      "failed payment 24 hours 1 minute ago",
			status:    PaymentStatusFailed,
			createdAt: now.Add(-24*time.Hour - 1*time.Minute),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := &Payment{
				Status:    tt.status,
				CreatedAt: tt.createdAt,
			}
			result := payment.CanRetry()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPayment_IsExpired_EdgeCases(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name      string
		expiresAt *time.Time
		expected  bool
	}{
		{
			name:      "expires exactly now",
			expiresAt: &now,
			expected:  true, // Should be expired if exactly at current time
		},
		{
			name: "expires 1 second from now",
			expiresAt: func() *time.Time {
				future := now.Add(1 * time.Second)
				return &future
			}(),
			expected: false,
		},
		{
			name: "expired 1 second ago",
			expiresAt: func() *time.Time {
				past := now.Add(-1 * time.Second)
				return &past
			}(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payment := &Payment{SessionExpiresAt: tt.expiresAt}
			result := payment.IsExpired()
			assert.Equal(t, tt.expected, result)
		})
	}
}