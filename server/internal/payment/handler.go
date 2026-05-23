package payment

import (
	"fmt"
	"net/http"
	"net/url"
	"travel-backend/domain"
	"travel-backend/internal/auth"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// PaymentHandler - Tầng xử lý HTTP cho Payment
type PaymentHandler struct {
	service *PaymentService
}

// NewPaymentHandler - Tạo payment handler mới
func NewPaymentHandler(config *VNPayConfig) *PaymentHandler {
	return &PaymentHandler{
		service: NewPaymentService(config),
	}
}

// InitiatePaymentHandler - POST /v1/api/payments/initiate
// Khởi tạo thanh toán cho 1 booking, trả về VNPay payment URL
func (h *PaymentHandler) InitiatePaymentHandler(c *gin.Context) {
	userID, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", "AUTH_TOKEN_INVALID")
		return
	}

	var req domain.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "PAYMENT_INVALID_PAYLOAD")
		return
	}

	if err := req.Validate(); err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "PAYMENT_VALIDATION_FAILED")
		return
	}

	clientIP := c.ClientIP()
	result, err := h.service.InitiatePayment(req.BookingID, userID, clientIP)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "PAYMENT_INITIATE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, result.Message, gin.H{
		"payment_url":           result.PaymentURL,
		"transaction_reference": result.TransactionReference,
		"expires_at":            result.ExpiresAt,
	})
}

// PaymentReturnHandler - GET /v1/api/payments/return
// VNPay redirect user về đây sau khi thanh toán
func (h *PaymentHandler) PaymentReturnHandler(c *gin.Context) {
	// Collect all vnp_ params from query string
	params := collectPaymentParams(c)

	secureHash := params["vnp_SecureHash"]
	delete(params, "vnp_SecureHash")
	delete(params, "vnp_SecureHashType")
	txnRef := params["vnp_TxnRef"]

	result, err := h.service.ProcessReturn(params, secureHash)
	if err != nil {
		// Redirect to frontend with error
		redirectURL := buildPaymentResultURL(h.service.config.FrontendURL, "error", txnRef, err.Error())
		c.Redirect(http.StatusFound, redirectURL)
		return
	}

	// Redirect to frontend with result
	status := "failed"
	if result.Success {
		status = "success"
	}

	redirectURL := buildPaymentResultURL(h.service.config.FrontendURL, status, txnRef, result.Message)
	c.Redirect(http.StatusFound, redirectURL)
}

// PaymentWebhookHandler - POST /v1/api/payments/webhook
// VNPay gọi endpoint này (IPN) để thông báo kết quả thanh toán
func (h *PaymentHandler) PaymentWebhookHandler(c *gin.Context) {
	// Collect all vnp_ params
	params := collectPaymentParams(c)

	secureHash := params["vnp_SecureHash"]
	delete(params, "vnp_SecureHash")
	delete(params, "vnp_SecureHashType")

	rspCode, message := h.service.ProcessWebhook(params, secureHash)

	// VNPay IPN yêu cầu trả về format cụ thể
	c.JSON(http.StatusOK, gin.H{
		"RspCode": rspCode,
		"Message": message,
	})
}

func collectPaymentParams(c *gin.Context) map[string]string {
	params := make(map[string]string)
	addValues := func(values url.Values) {
		for key, vals := range values {
			if len(vals) > 0 {
				params[key] = vals[0]
			}
		}
	}

	addValues(c.Request.URL.Query())
	if err := c.Request.ParseForm(); err == nil {
		addValues(c.Request.PostForm)
	}

	return params
}

func buildPaymentResultURL(frontendBase, status, ref, message string) string {
	resultURL, err := url.Parse(frontendBase)
	if err != nil {
		resultURL = &url.URL{Scheme: "http", Host: "localhost:5173"}
	}
	resultURL.Path = "/payment/result"
	query := resultURL.Query()
	query.Set("status", status)
	query.Set("ref", ref)
	query.Set("message", message)
	resultURL.RawQuery = query.Encode()
	return resultURL.String()
}

// GetPaymentStatusHandler - GET /v1/api/payments/status/:ref
// Kiểm tra trạng thái thanh toán theo transaction reference
func (h *PaymentHandler) GetPaymentStatusHandler(c *gin.Context) {
	userID, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", "AUTH_TOKEN_INVALID")
		return
	}

	txnRef := c.Param("ref")
	if txnRef == "" {
		shared.RespondError(c, http.StatusBadRequest, "Mã giao dịch không hợp lệ", "PAYMENT_INVALID_REF")
		return
	}

	result, err := h.service.GetPaymentStatus(txnRef, userID)
	if err != nil {
		shared.RespondError(c, http.StatusNotFound, err.Error(), "PAYMENT_NOT_FOUND")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, result.Message, gin.H{
		"payment": result.Payment,
	})
}

// GetBookingPaymentsHandler - GET /v1/api/bookings/:id/payments
// Lấy danh sách thanh toán theo booking
func (h *PaymentHandler) GetBookingPaymentsHandler(c *gin.Context) {
	userID, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", "AUTH_TOKEN_INVALID")
		return
	}

	var bookingID uint
	if _, err := fmt.Sscanf(c.Param("id"), "%d", &bookingID); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "ID booking không hợp lệ", "PAYMENT_INVALID_BOOKING_ID")
		return
	}

	payments, err := h.service.GetPaymentsByBooking(bookingID, userID)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "PAYMENT_FETCH_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "", gin.H{
		"payments": payments,
		"total":    len(payments),
	})
}
