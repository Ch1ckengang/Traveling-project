package payment

import (
	"fmt"
	"log"
	"net"
	"time"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/notification"
	"travel-backend/internal/shared"
)

// PaymentService - Tầng xử lý business logic cho Payment
type PaymentService struct {
	repo       PaymentRepository
	vnpayClient *VNPayClient
	config     *VNPayConfig
}

// NewPaymentService - Tạo payment service mới
func NewPaymentService(config *VNPayConfig) *PaymentService {
	return &PaymentService{
		repo:        NewPaymentRepository(),
		vnpayClient: NewVNPayClient(config),
		config:      config,
	}
}

// InitiatePayment - Khởi tạo thanh toán cho booking
func (s *PaymentService) InitiatePayment(bookingID, userID uint, clientIP string) (*domain.InitiatePaymentResponse, error) {
	// 1. Lấy thông tin booking
	var booking domain.Booking
	if err := database.DB.Preload("Tour").First(&booking, bookingID).Error; err != nil {
		return nil, fmt.Errorf("không tìm thấy booking: %w", err)
	}

	// 2. Kiểm tra booking thuộc về user
	if booking.UserID != userID {
		return nil, fmt.Errorf("bạn không có quyền thanh toán booking này")
	}

	// 3. Kiểm tra booking status
	if booking.Status == "cancelled" {
		return nil, fmt.Errorf("booking đã bị hủy, không thể thanh toán")
	}
	if booking.PaymentStatus == "paid" {
		return nil, fmt.Errorf("booking đã được thanh toán")
	}

	// 4. Tính toán số tiền
	amount := booking.TotalAmount
	if amount <= 0 {
		// Fallback: tính từ tour price
		amount = booking.Tour.PriceAmount * int64(booking.Quantity)
	}
	if amount <= 0 {
		return nil, fmt.Errorf("không thể xác định số tiền thanh toán")
	}

	// 5. Tạo payment record
	payment := domain.NewPayment(bookingID, amount)
	if err := s.repo.CreatePayment(payment); err != nil {
		return nil, fmt.Errorf("không thể tạo phiên thanh toán: %w", err)
	}

	// 6. Tạo VNPay payment URL
	orderInfo := fmt.Sprintf("Thanh toan booking %s", booking.BookingCode)
	vnpayReq := &domain.VNPayPaymentRequest{
		TransactionReference: payment.TransactionReference,
		Amount:               amount,
		OrderInfo:            orderInfo,
		ClientIP:             clientIP,
		ReturnURL:            s.config.ReturnURL,
		ExpiresAt:            *payment.SessionExpiresAt,
	}

	paymentURL, err := s.vnpayClient.GeneratePaymentURL(vnpayReq)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo URL thanh toán: %w", err)
	}

	// 7. Cập nhật booking với payment info
	database.DB.Model(&booking).Updates(map[string]interface{}{
		"payment_id":       payment.ID,
		"payment_status":   "processing",
		"payment_deadline": payment.SessionExpiresAt,
	})

	// 8. Ghi audit log
	ipAddr := net.ParseIP(clientIP)
	auditLog := domain.NewPaymentInitiatedLog(payment, userID, ipAddr, "")
	s.repo.CreateAuditLog(auditLog)

	log.Printf("[PAYMENT][INITIATE] booking_id=%d payment_id=%d ref=%s amount=%d",
		bookingID, payment.ID, payment.TransactionReference, amount)

	return &domain.InitiatePaymentResponse{
		Success:              true,
		PaymentURL:           paymentURL,
		TransactionReference: payment.TransactionReference,
		ExpiresAt:            *payment.SessionExpiresAt,
		Message:              "Đã tạo phiên thanh toán thành công",
	}, nil
}

// ProcessReturn - Xử lý khi user được redirect về từ VNPay
func (s *PaymentService) ProcessReturn(params map[string]string, secureHash string) (*domain.PaymentStatusResponse, error) {
	txnRef := params["vnp_TxnRef"]

	// 1. Validate signature
	if !s.vnpayClient.ValidateSignature(params, secureHash) {
		log.Printf("[PAYMENT][RETURN] INVALID SIGNATURE for ref=%s", txnRef)
		return &domain.PaymentStatusResponse{
			Success: false,
			Message: "Chữ ký không hợp lệ",
		}, nil
	}

	// 2. Lấy payment
	payment, err := s.repo.GetPaymentByTransactionReference(txnRef)
	if err != nil || payment == nil {
		return &domain.PaymentStatusResponse{
			Success: false,
			Message: "Không tìm thấy giao dịch",
		}, nil
	}

	// 3. Xử lý kết quả
	responseCode := params["vnp_ResponseCode"]
	transactionNo := params["vnp_TransactionNo"]
	isSuccess, message := s.vnpayClient.ParseVNPayResponseCode(responseCode)

	if isSuccess {
		s.handlePaymentSuccess(payment, responseCode, transactionNo, message)
	} else {
		s.handlePaymentFailure(payment, responseCode, message)
	}

	return &domain.PaymentStatusResponse{
		Success: isSuccess,
		Payment: payment.ToSummary(),
		Message: message,
	}, nil
}

// ProcessWebhook - Xử lý IPN webhook từ VNPay (server-to-server)
func (s *PaymentService) ProcessWebhook(params map[string]string, secureHash string) (string, string) {
	txnRef := params["vnp_TxnRef"]

	// 1. Validate signature
	if !s.vnpayClient.ValidateSignature(params, secureHash) {
		log.Printf("[PAYMENT][WEBHOOK] INVALID SIGNATURE for ref=%s", txnRef)
		return "97", "Invalid Signature"
	}

	// 2. Lấy payment
	payment, err := s.repo.GetPaymentByTransactionReference(txnRef)
	if err != nil || payment == nil {
		log.Printf("[PAYMENT][WEBHOOK] Payment not found for ref=%s", txnRef)
		return "01", "Order Not Found"
	}

	// 3. Kiểm tra đã xử lý chưa
	if payment.Status == domain.PaymentStatusPaid || payment.Status == domain.PaymentStatusRefunded {
		return "02", "Order Already Confirmed"
	}

	// 4. Xử lý kết quả
	responseCode := params["vnp_ResponseCode"]
	transactionNo := params["vnp_TransactionNo"]
	isSuccess, message := s.vnpayClient.ParseVNPayResponseCode(responseCode)

	// Ghi audit log webhook
	webhookData := make(map[string]interface{})
	for k, v := range params {
		webhookData[k] = v
	}
	auditLog := domain.NewWebhookReceivedLog(&payment.ID, &payment.BookingID, webhookData)
	s.repo.CreateAuditLog(auditLog)

	if isSuccess {
		s.handlePaymentSuccess(payment, responseCode, transactionNo, message)
		return "00", "Confirm Success"
	}

	s.handlePaymentFailure(payment, responseCode, message)
	return "00", "Confirm Success"
}

// GetPaymentStatus - Kiểm tra trạng thái thanh toán
func (s *PaymentService) GetPaymentStatus(txnRef string, userID uint) (*domain.PaymentStatusResponse, error) {
	payment, err := s.repo.GetPaymentByTransactionReference(txnRef)
	if err != nil || payment == nil {
		return nil, fmt.Errorf("không tìm thấy giao dịch")
	}

	// Kiểm tra quyền: lấy booking để verify user ownership
	var booking domain.Booking
	if err := database.DB.First(&booking, payment.BookingID).Error; err != nil {
		return nil, fmt.Errorf("không tìm thấy booking liên quan")
	}
	if booking.UserID != userID {
		return nil, fmt.Errorf("bạn không có quyền xem giao dịch này")
	}

	return &domain.PaymentStatusResponse{
		Success: true,
		Payment: payment.ToSummary(),
		Message: payment.GetStatusDisplayName(),
	}, nil
}

// GetPaymentsByBooking - Lấy danh sách thanh toán theo booking
func (s *PaymentService) GetPaymentsByBooking(bookingID, userID uint) ([]domain.PaymentSummary, error) {
	// Verify ownership
	var booking domain.Booking
	if err := database.DB.First(&booking, bookingID).Error; err != nil {
		return nil, fmt.Errorf("không tìm thấy booking")
	}
	if booking.UserID != userID {
		return nil, fmt.Errorf("bạn không có quyền xem")
	}

	payments, err := s.repo.GetPaymentsByBookingID(bookingID)
	if err != nil {
		return nil, err
	}

	summaries := make([]domain.PaymentSummary, len(payments))
	for i, p := range payments {
		summaries[i] = *p.ToSummary()
	}
	return summaries, nil
}

// handlePaymentSuccess - Xử lý khi thanh toán thành công
func (s *PaymentService) handlePaymentSuccess(payment *domain.Payment, responseCode, transactionNo, message string) {
	previousStatus := payment.Status

	payment.SetVNPayResponse(transactionNo, responseCode, message)
	payment.Status = domain.PaymentStatusPaid
	payment.UpdatedAt = time.Now()
	method := "vnpay"
	payment.PaymentMethod = &method

	if err := s.repo.UpdatePayment(payment); err != nil {
		log.Printf("[PAYMENT][ERROR] Failed to update payment %d: %v", payment.ID, err)
		return
	}

	// Cập nhật booking status
	database.DB.Model(&domain.Booking{}).Where("id = ?", payment.BookingID).Updates(map[string]interface{}{
		"payment_status": "paid",
		"status":         "confirmed",
	})

	// Audit log
	auditLog := domain.NewPaymentCompletedLog(payment, responseCode, message)
	s.repo.CreateAuditLog(auditLog)

	// Gửi email xác nhận
	s.sendPaymentConfirmationEmail(payment)

	// Gửi thông báo trong app
	go func(bookingID uint) {
		var b domain.Booking
		if err := database.DB.Preload("Tour").First(&b, bookingID).Error; err == nil {
			msg := fmt.Sprintf("Thanh toán thành công cho tour %s. Số tiền: %s.", b.Tour.Name, domain.FormatPriceVND(payment.Amount))
			_ = notification.SendNotification(b.UserID, "Thanh toán thành công", msg, domain.NotifTypePayment)
		}
	}(payment.BookingID)

	log.Printf("[PAYMENT][SUCCESS] payment_id=%d ref=%s prev_status=%s new_status=paid",
		payment.ID, payment.TransactionReference, previousStatus)
}

// handlePaymentFailure - Xử lý khi thanh toán thất bại
func (s *PaymentService) handlePaymentFailure(payment *domain.Payment, responseCode, message string) {
	previousStatus := payment.Status

	payment.SetVNPayResponse("", responseCode, message)
	payment.Status = domain.PaymentStatusFailed
	payment.UpdatedAt = time.Now()

	if err := s.repo.UpdatePayment(payment); err != nil {
		log.Printf("[PAYMENT][ERROR] Failed to update payment %d: %v", payment.ID, err)
		return
	}

	// Cập nhật booking
	database.DB.Model(&domain.Booking{}).Where("id = ?", payment.BookingID).Updates(map[string]interface{}{
		"payment_status": "failed",
	})

	// Audit log
	auditLog := domain.NewPaymentFailedLog(payment, message, responseCode)
	s.repo.CreateAuditLog(auditLog)

	log.Printf("[PAYMENT][FAILED] payment_id=%d ref=%s code=%s prev=%s",
		payment.ID, payment.TransactionReference, responseCode, previousStatus)
}

// sendPaymentConfirmationEmail - Gửi email xác nhận thanh toán
func (s *PaymentService) sendPaymentConfirmationEmail(payment *domain.Payment) {
	var booking domain.Booking
	if err := database.DB.Preload("Tour").First(&booking, payment.BookingID).Error; err != nil {
		log.Printf("[PAYMENT][EMAIL] Failed to load booking %d: %v", payment.BookingID, err)
		return
	}

	if booking.Email == "" {
		return
	}

	subject := fmt.Sprintf("Xác nhận thanh toán - Booking %s", booking.BookingCode)
	body := fmt.Sprintf(`
		<h2>Thanh toán thành công!</h2>
		<p>Xin chào %s,</p>
		<p>Chúng tôi xác nhận đã nhận được thanh toán cho booking của bạn:</p>
		<ul>
			<li><strong>Mã booking:</strong> %s</li>
			<li><strong>Tour:</strong> %s</li>
			<li><strong>Số tiền:</strong> %s</li>
			<li><strong>Mã giao dịch:</strong> %s</li>
		</ul>
		<p>Cảm ơn bạn đã sử dụng dịch vụ!</p>
	`, booking.FullName, booking.BookingCode, booking.Tour.Name,
		domain.FormatPriceVND(payment.Amount), payment.TransactionReference)

	if err := shared.SendEmail(booking.Email, subject, body); err != nil {
		log.Printf("[PAYMENT][EMAIL] Failed to send to %s: %v", booking.Email, err)
	}
}
