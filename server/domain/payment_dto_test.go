package domain

import (
	"testing"
	"time"
)

func TestInitiatePaymentRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request InitiatePaymentRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: InitiatePaymentRequest{
				BookingID: 123,
			},
			wantErr: false,
		},
		{
			name: "zero booking ID",
			request: InitiatePaymentRequest{
				BookingID: 0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitiatePaymentRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRefundPaymentRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request RefundPaymentRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: RefundPaymentRequest{
				PaymentID:    123,
				RefundAmount: 100000,
				Reason:       "Customer cancellation",
			},
			wantErr: false,
		},
		{
			name: "zero payment ID",
			request: RefundPaymentRequest{
				PaymentID:    0,
				RefundAmount: 100000,
				Reason:       "Customer cancellation",
			},
			wantErr: true,
		},
		{
			name: "zero refund amount",
			request: RefundPaymentRequest{
				PaymentID:    123,
				RefundAmount: 0,
				Reason:       "Customer cancellation",
			},
			wantErr: true,
		},
		{
			name: "negative refund amount",
			request: RefundPaymentRequest{
				PaymentID:    123,
				RefundAmount: -100000,
				Reason:       "Customer cancellation",
			},
			wantErr: true,
		},
		{
			name: "empty reason",
			request: RefundPaymentRequest{
				PaymentID:    123,
				RefundAmount: 100000,
				Reason:       "",
			},
			wantErr: true,
		},
		{
			name: "reason too long",
			request: RefundPaymentRequest{
				PaymentID:    123,
				RefundAmount: 100000,
				Reason:       string(make([]byte, 501)), // 501 characters
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("RefundPaymentRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaymentHistoryRequest_Validate(t *testing.T) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	yesterday := now.Add(-24 * time.Hour)

	tests := []struct {
		name    string
		request PaymentHistoryRequest
		wantErr bool
	}{
		{
			name: "valid request with defaults",
			request: PaymentHistoryRequest{
				Page:  0, // Should be set to 1
				Limit: 0, // Should be set to 20
			},
			wantErr: false,
		},
		{
			name: "valid request with custom values",
			request: PaymentHistoryRequest{
				Page:  2,
				Limit: 50,
			},
			wantErr: false,
		},
		{
			name: "limit too high",
			request: PaymentHistoryRequest{
				Page:  1,
				Limit: 200, // Should be capped to 100
			},
			wantErr: false,
		},
		{
			name: "valid date range",
			request: PaymentHistoryRequest{
				Page:     1,
				Limit:    20,
				FromDate: &yesterday,
				ToDate:   &tomorrow,
			},
			wantErr: false,
		},
		{
			name: "invalid date range - from after to",
			request: PaymentHistoryRequest{
				Page:     1,
				Limit:    20,
				FromDate: &tomorrow,
				ToDate:   &yesterday,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PaymentHistoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check that defaults are applied
			if tt.request.Page == 0 {
				if tt.request.Page != 1 {
					t.Errorf("Expected Page to be set to 1, got %d", tt.request.Page)
				}
			}
			if tt.request.Limit == 0 {
				if tt.request.Limit != 20 {
					t.Errorf("Expected Limit to be set to 20, got %d", tt.request.Limit)
				}
			}
			if tt.request.Limit > 100 {
				if tt.request.Limit != 100 {
					t.Errorf("Expected Limit to be capped to 100, got %d", tt.request.Limit)
				}
			}
		})
	}
}

func TestVNPayReturnParams_ToMap(t *testing.T) {
	params := VNPayReturnParams{
		VnpAmount:            "250000000",
		VnpBankCode:          "NCB",
		VnpBankTranNo:        "VNP01234567",
		VnpCardType:          "ATM",
		VnpOrderInfo:         "Thanh toan tour ABC",
		VnpPayDate:           "20240430153000",
		VnpResponseCode:      "00",
		VnpTmnCode:           "MERCHANT123",
		VnpTransactionNo:     "13123456",
		VnpTransactionStatus: "00",
		VnpTxnRef:            "PAY20240430001",
		VnpSecureHash:        "abcd1234",
	}

	paramMap := params.ToMap()

	// Check that all fields are included except VnpSecureHash
	expectedFields := []string{
		"vnp_Amount", "vnp_BankCode", "vnp_BankTranNo", "vnp_CardType",
		"vnp_OrderInfo", "vnp_PayDate", "vnp_ResponseCode", "vnp_TmnCode",
		"vnp_TransactionNo", "vnp_TransactionStatus", "vnp_TxnRef",
	}

	for _, field := range expectedFields {
		if _, exists := paramMap[field]; !exists {
			t.Errorf("Expected field %s to be in map", field)
		}
	}

	// Check that VnpSecureHash is not included
	if _, exists := paramMap["vnp_SecureHash"]; exists {
		t.Error("VnpSecureHash should not be included in ToMap() result")
	}

	// Check specific values
	if paramMap["vnp_Amount"] != params.VnpAmount {
		t.Errorf("Expected vnp_Amount to be %s, got %s", params.VnpAmount, paramMap["vnp_Amount"])
	}
	if paramMap["vnp_TxnRef"] != params.VnpTxnRef {
		t.Errorf("Expected vnp_TxnRef to be %s, got %s", params.VnpTxnRef, paramMap["vnp_TxnRef"])
	}
}

func TestVNPayWebhookParams_ToMap(t *testing.T) {
	params := VNPayWebhookParams{
		VnpAmount:            "250000000",
		VnpBankCode:          "NCB",
		VnpBankTranNo:        "VNP01234567",
		VnpCardType:          "ATM",
		VnpOrderInfo:         "Thanh toan tour ABC",
		VnpPayDate:           "20240430153000",
		VnpResponseCode:      "00",
		VnpTmnCode:           "MERCHANT123",
		VnpTransactionNo:     "13123456",
		VnpTransactionStatus: "00",
		VnpTxnRef:            "PAY20240430001",
		VnpSecureHash:        "abcd1234",
	}

	paramMap := params.ToMap()

	// Check that all fields are included except VnpSecureHash
	expectedFields := []string{
		"vnp_Amount", "vnp_BankCode", "vnp_BankTranNo", "vnp_CardType",
		"vnp_OrderInfo", "vnp_PayDate", "vnp_ResponseCode", "vnp_TmnCode",
		"vnp_TransactionNo", "vnp_TransactionStatus", "vnp_TxnRef",
	}

	for _, field := range expectedFields {
		if _, exists := paramMap[field]; !exists {
			t.Errorf("Expected field %s to be in map", field)
		}
	}

	// Check that VnpSecureHash is not included
	if _, exists := paramMap["vnp_SecureHash"]; exists {
		t.Error("VnpSecureHash should not be included in ToMap() result")
	}
}

func TestVNPayReturnParams_IsSuccessful(t *testing.T) {
	tests := []struct {
		name   string
		params VNPayReturnParams
		want   bool
	}{
		{
			name: "successful payment",
			params: VNPayReturnParams{
				VnpResponseCode:      "00",
				VnpTransactionStatus: "00",
			},
			want: true,
		},
		{
			name: "failed payment - response code not 00",
			params: VNPayReturnParams{
				VnpResponseCode:      "24",
				VnpTransactionStatus: "00",
			},
			want: false,
		},
		{
			name: "failed payment - transaction status not 00",
			params: VNPayReturnParams{
				VnpResponseCode:      "00",
				VnpTransactionStatus: "01",
			},
			want: false,
		},
		{
			name: "failed payment - both not 00",
			params: VNPayReturnParams{
				VnpResponseCode:      "24",
				VnpTransactionStatus: "01",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsSuccessful(); got != tt.want {
				t.Errorf("VNPayReturnParams.IsSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVNPayWebhookParams_IsSuccessful(t *testing.T) {
	tests := []struct {
		name   string
		params VNPayWebhookParams
		want   bool
	}{
		{
			name: "successful payment",
			params: VNPayWebhookParams{
				VnpResponseCode:      "00",
				VnpTransactionStatus: "00",
			},
			want: true,
		},
		{
			name: "failed payment - response code not 00",
			params: VNPayWebhookParams{
				VnpResponseCode:      "24",
				VnpTransactionStatus: "00",
			},
			want: false,
		},
		{
			name: "failed payment - transaction status not 00",
			params: VNPayWebhookParams{
				VnpResponseCode:      "00",
				VnpTransactionStatus: "01",
			},
			want: false,
		},
		{
			name: "failed payment - both not 00",
			params: VNPayWebhookParams{
				VnpResponseCode:      "24",
				VnpTransactionStatus: "01",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.params.IsSuccessful(); got != tt.want {
				t.Errorf("VNPayWebhookParams.IsSuccessful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetVNPayErrorMessage(t *testing.T) {
	tests := []struct {
		name         string
		responseCode string
		want         string
	}{
		{
			name:         "success code",
			responseCode: "00",
			want:         "Giao dịch thành công",
		},
		{
			name:         "customer cancelled",
			responseCode: "24",
			want:         "Giao dịch không thành công do: Khách hàng hủy giao dịch",
		},
		{
			name:         "insufficient funds",
			responseCode: "51",
			want:         "Giao dịch không thành công do: Tài khoản của quý khách không đủ số dư để thực hiện giao dịch.",
		},
		{
			name:         "unknown error code",
			responseCode: "999",
			want:         "Lỗi không xác định",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVNPayErrorMessage(tt.responseCode); got != tt.want {
				t.Errorf("GetVNPayErrorMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationInfo(t *testing.T) {
	pagination := &PaginationInfo{
		Page:       2,
		Limit:      20,
		Total:      100,
		TotalPages: 5,
		HasNext:    true,
		HasPrev:    true,
	}

	// Test that all fields are properly set
	if pagination.Page != 2 {
		t.Errorf("Expected Page to be 2, got %d", pagination.Page)
	}
	if pagination.Limit != 20 {
		t.Errorf("Expected Limit to be 20, got %d", pagination.Limit)
	}
	if pagination.Total != 100 {
		t.Errorf("Expected Total to be 100, got %d", pagination.Total)
	}
	if pagination.TotalPages != 5 {
		t.Errorf("Expected TotalPages to be 5, got %d", pagination.TotalPages)
	}
	if !pagination.HasNext {
		t.Error("Expected HasNext to be true")
	}
	if !pagination.HasPrev {
		t.Error("Expected HasPrev to be true")
	}
}