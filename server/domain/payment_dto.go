package domain

import (
	"fmt"
	"time"
)

// Request DTOs

// InitiatePaymentRequest - Request để khởi tạo thanh toán
type InitiatePaymentRequest struct {
	BookingID uint `json:"booking_id" binding:"required,min=1"`
}

// RefundPaymentRequest - Request để hoàn tiền
type RefundPaymentRequest struct {
	PaymentID    uint   `json:"payment_id" binding:"required,min=1"`
	RefundAmount int64  `json:"refund_amount" binding:"required,min=1"`
	Reason       string `json:"reason" binding:"required,max=500"`
}

// Response DTOs

// InitiatePaymentResponse - Response cho việc khởi tạo thanh toán
type InitiatePaymentResponse struct {
	Success              bool      `json:"success"`
	PaymentURL           string    `json:"payment_url"`
	TransactionReference string    `json:"transaction_reference"`
	ExpiresAt            time.Time `json:"expires_at"`
	Message              string    `json:"message,omitempty"`
}

// PaymentStatusResponse - Response cho việc kiểm tra trạng thái thanh toán
type PaymentStatusResponse struct {
	Success bool            `json:"success"`
	Payment *PaymentSummary `json:"payment"`
	Message string          `json:"message,omitempty"`
}

// RefundPaymentResponse - Response cho việc hoàn tiền
type RefundPaymentResponse struct {
	Success       bool   `json:"success"`
	RefundID      string `json:"refund_id,omitempty"`
	RefundAmount  int64  `json:"refund_amount"`
	RefundStatus  string `json:"refund_status"`
	EstimatedTime string `json:"estimated_time,omitempty"`
	Message       string `json:"message,omitempty"`
}

// VNPay specific DTOs

// VNPayReturnParams - Parameters từ VNPay return URL
type VNPayReturnParams struct {
	VnpAmount          string `form:"vnp_Amount" binding:"required"`
	VnpBankCode        string `form:"vnp_BankCode"`
	VnpBankTranNo      string `form:"vnp_BankTranNo"`
	VnpCardType        string `form:"vnp_CardType"`
	VnpOrderInfo       string `form:"vnp_OrderInfo"`
	VnpPayDate         string `form:"vnp_PayDate" binding:"required"`
	VnpResponseCode    string `form:"vnp_ResponseCode" binding:"required"`
	VnpTmnCode         string `form:"vnp_TmnCode" binding:"required"`
	VnpTransactionNo   string `form:"vnp_TransactionNo"`
	VnpTransactionStatus string `form:"vnp_TransactionStatus" binding:"required"`
	VnpTxnRef          string `form:"vnp_TxnRef" binding:"required"`
	VnpSecureHash      string `form:"vnp_SecureHash" binding:"required"`
}

// VNPayWebhookParams - Parameters từ VNPay IPN webhook
type VNPayWebhookParams struct {
	VnpAmount          string `json:"vnp_Amount" binding:"required"`
	VnpBankCode        string `json:"vnp_BankCode"`
	VnpBankTranNo      string `json:"vnp_BankTranNo"`
	VnpCardType        string `json:"vnp_CardType"`
	VnpOrderInfo       string `json:"vnp_OrderInfo"`
	VnpPayDate         string `json:"vnp_PayDate" binding:"required"`
	VnpResponseCode    string `json:"vnp_ResponseCode" binding:"required"`
	VnpTmnCode         string `json:"vnp_TmnCode" binding:"required"`
	VnpTransactionNo   string `json:"vnp_TransactionNo"`
	VnpTransactionStatus string `json:"vnp_TransactionStatus" binding:"required"`
	VnpTxnRef          string `json:"vnp_TxnRef" binding:"required"`
	VnpSecureHash      string `json:"vnp_SecureHash" binding:"required"`
}

// VNPayPaymentRequest - Request để tạo VNPay payment URL
type VNPayPaymentRequest struct {
	TransactionReference string
	Amount               int64
	OrderInfo            string
	ClientIP             string
	ReturnURL            string
	ExpiresAt            time.Time
}

// Payment history and reporting DTOs

// PaymentHistoryRequest - Request để lấy lịch sử thanh toán
type PaymentHistoryRequest struct {
	UserID    *uint     `form:"user_id"`
	BookingID *uint     `form:"booking_id"`
	Status    *string   `form:"status"`
	FromDate  *time.Time `form:"from_date"`
	ToDate    *time.Time `form:"to_date"`
	Page      int       `form:"page,default=1" binding:"min=1"`
	Limit     int       `form:"limit,default=20" binding:"min=1,max=100"`
}

// PaymentHistoryResponse - Response cho lịch sử thanh toán
type PaymentHistoryResponse struct {
	Success    bool               `json:"success"`
	Payments   []*PaymentSummary  `json:"payments"`
	Pagination *PaginationInfo    `json:"pagination"`
	Message    string             `json:"message,omitempty"`
}

// PaginationInfo - Thông tin phân trang
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// Audit log DTOs

// PaymentAuditRequest - Request để lấy audit logs
type PaymentAuditRequest struct {
	PaymentID *uint     `form:"payment_id"`
	BookingID *uint     `form:"booking_id"`
	EventType *string   `form:"event_type"`
	UserID    *uint     `form:"user_id"`
	FromDate  *time.Time `form:"from_date"`
	ToDate    *time.Time `form:"to_date"`
	Page      int       `form:"page,default=1" binding:"min=1"`
	Limit     int       `form:"limit,default=50" binding:"min=1,max=200"`
}

// PaymentAuditResponse - Response cho audit logs
type PaymentAuditResponse struct {
	Success    bool                   `json:"success"`
	AuditLogs  []*PaymentAuditSummary `json:"audit_logs"`
	Pagination *PaginationInfo        `json:"pagination"`
	Message    string                 `json:"message,omitempty"`
}

// Error DTOs

// PaymentErrorResponse - Response cho lỗi thanh toán
type PaymentErrorResponse struct {
	Success   bool   `json:"success"`
	Error     string `json:"error"`
	ErrorCode string `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
}

// Validation methods

// Validate - Validate InitiatePaymentRequest
func (r *InitiatePaymentRequest) Validate() error {
	if r.BookingID == 0 {
		return fmt.Errorf("booking_id is required and must be greater than 0")
	}
	return nil
}

// Validate - Validate RefundPaymentRequest
func (r *RefundPaymentRequest) Validate() error {
	if r.PaymentID == 0 {
		return fmt.Errorf("payment_id is required and must be greater than 0")
	}
	if r.RefundAmount <= 0 {
		return fmt.Errorf("refund_amount must be greater than 0")
	}
	if r.Reason == "" {
		return fmt.Errorf("reason is required")
	}
	if len(r.Reason) > 500 {
		return fmt.Errorf("reason must not exceed 500 characters")
	}
	return nil
}

// Validate - Validate PaymentHistoryRequest
func (r *PaymentHistoryRequest) Validate() error {
	if r.Page < 1 {
		r.Page = 1
	}
	if r.Limit < 1 {
		r.Limit = 20
	}
	if r.Limit > 100 {
		r.Limit = 100
	}
	
	if r.FromDate != nil && r.ToDate != nil && r.FromDate.After(*r.ToDate) {
		return fmt.Errorf("from_date must be before to_date")
	}
	
	return nil
}

// Helper methods

// ToMap - Convert VNPayReturnParams to map for signature validation
func (p *VNPayReturnParams) ToMap() map[string]string {
	return map[string]string{
		"vnp_Amount":           p.VnpAmount,
		"vnp_BankCode":         p.VnpBankCode,
		"vnp_BankTranNo":       p.VnpBankTranNo,
		"vnp_CardType":         p.VnpCardType,
		"vnp_OrderInfo":        p.VnpOrderInfo,
		"vnp_PayDate":          p.VnpPayDate,
		"vnp_ResponseCode":     p.VnpResponseCode,
		"vnp_TmnCode":          p.VnpTmnCode,
		"vnp_TransactionNo":    p.VnpTransactionNo,
		"vnp_TransactionStatus": p.VnpTransactionStatus,
		"vnp_TxnRef":           p.VnpTxnRef,
	}
}

// ToMap - Convert VNPayWebhookParams to map for signature validation
func (p *VNPayWebhookParams) ToMap() map[string]string {
	return map[string]string{
		"vnp_Amount":           p.VnpAmount,
		"vnp_BankCode":         p.VnpBankCode,
		"vnp_BankTranNo":       p.VnpBankTranNo,
		"vnp_CardType":         p.VnpCardType,
		"vnp_OrderInfo":        p.VnpOrderInfo,
		"vnp_PayDate":          p.VnpPayDate,
		"vnp_ResponseCode":     p.VnpResponseCode,
		"vnp_TmnCode":          p.VnpTmnCode,
		"vnp_TransactionNo":    p.VnpTransactionNo,
		"vnp_TransactionStatus": p.VnpTransactionStatus,
		"vnp_TxnRef":           p.VnpTxnRef,
	}
}

// IsSuccessful - Check if VNPay return indicates successful payment
func (p *VNPayReturnParams) IsSuccessful() bool {
	return p.VnpResponseCode == "00" && p.VnpTransactionStatus == "00"
}

// IsSuccessful - Check if VNPay webhook indicates successful payment
func (p *VNPayWebhookParams) IsSuccessful() bool {
	return p.VnpResponseCode == "00" && p.VnpTransactionStatus == "00"
}

// GetVNPayErrorMessage - Get error message from VNPay response code
func GetVNPayErrorMessage(responseCode string) string {
	errorMessages := map[string]string{
		"00": "Giao dịch thành công",
		"07": "Trừ tiền thành công. Giao dịch bị nghi ngờ (liên quan tới lừa đảo, giao dịch bất thường).",
		"09": "Giao dịch không thành công do: Thẻ/Tài khoản của khách hàng chưa đăng ký dịch vụ InternetBanking tại ngân hàng.",
		"10": "Giao dịch không thành công do: Khách hàng xác thực thông tin thẻ/tài khoản không đúng quá 3 lần",
		"11": "Giao dịch không thành công do: Đã hết hạn chờ thanh toán. Xin quý khách vui lòng thực hiện lại giao dịch.",
		"12": "Giao dịch không thành công do: Thẻ/Tài khoản của khách hàng bị khóa.",
		"13": "Giao dịch không thành công do Quý khách nhập sai mật khẩu xác thực giao dịch (OTP). Xin quý khách vui lòng thực hiện lại giao dịch.",
		"24": "Giao dịch không thành công do: Khách hàng hủy giao dịch",
		"51": "Giao dịch không thành công do: Tài khoản của quý khách không đủ số dư để thực hiện giao dịch.",
		"65": "Giao dịch không thành công do: Tài khoản của Quý khách đã vượt quá hạn mức giao dịch trong ngày.",
		"75": "Ngân hàng thanh toán đang bảo trì.",
		"79": "Giao dịch không thành công do: KH nhập sai mật khẩu thanh toán quá số lần quy định. Xin quý khách vui lòng thực hiện lại giao dịch",
		"99": "Các lỗi khác (lỗi còn lại, không có trong danh sách mã lỗi đã liệt kê)",
	}
	
	if message, exists := errorMessages[responseCode]; exists {
		return message
	}
	
	return "Lỗi không xác định"
}