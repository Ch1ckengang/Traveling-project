package payment

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/notification"
	"travel-backend/internal/shared"

	"gorm.io/gorm"
)

// PaymentService - Tầng xử lý business logic cho Payment
type PaymentService struct {
	repo    PaymentRepository
	gateway PaymentGateway
	config  *VNPayConfig
}

// NewPaymentService - Tạo payment service mới
func NewPaymentService(config *VNPayConfig) *PaymentService {
	return &PaymentService{
		repo:    NewPaymentRepository(),
		gateway: NewVNPayClient(config),
		config:  config,
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
		// Fallback: tính từ tour price_amount
		amount = booking.Tour.PriceAmount * int64(booking.Quantity)
	}
	if amount <= 0 && booking.Tour.Price != "" {
		// Fallback: parse Price string (e.g. "5.000.000đ" → 5000000)
		parsedPrice := domain.ParsePriceVND(booking.Tour.Price)
		if parsedPrice > 0 {
			amount = parsedPrice * int64(booking.Quantity)
			log.Printf("[PAYMENT][WARN] booking_id=%d used Price string fallback: %s → %d", bookingID, booking.Tour.Price, parsedPrice)
		}
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

	paymentURL, err := s.gateway.GeneratePaymentURL(vnpayReq)
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
	if !s.gateway.ValidateSignature(params, secureHash) {
		log.Printf("[PAYMENT][RETURN] INVALID SIGNATURE for ref=%s", txnRef)
		return nil, fmt.Errorf("chữ ký không hợp lệ")
	}

	// 2. Lấy payment
	payment, err := s.repo.GetPaymentByTransactionReference(txnRef)
	if err != nil || payment == nil {
		return &domain.PaymentStatusResponse{
			Success: false,
			Message: "Không tìm thấy giao dịch",
		}, nil
	}

	// ReturnURL is browser-facing. Do not mutate payment state here; VNPay IPN
	// is the authoritative server-to-server update path.
	responseCode := params["vnp_ResponseCode"]
	isSuccess, message := s.gateway.ParseResponseCode(responseCode)
	isSuccess = isSuccess && params["vnp_TransactionStatus"] == "00"

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
	if !s.gateway.ValidateSignature(params, secureHash) {
		log.Printf("[PAYMENT][WEBHOOK] INVALID SIGNATURE for ref=%s", txnRef)
		return "97", "Invalid Signature"
	}

	webhookData := make(map[string]interface{})
	for k, v := range params {
		webhookData[k] = v
	}

	var updatedPayment *domain.Payment
	var completed bool
	rspCode := "00"
	rspMessage := "Confirm Success"

	err := s.repo.WithTransaction(func(tx *gorm.DB) error {
		payment, err := s.repo.GetPaymentWithLock(tx, txnRef)
		if err != nil {
			log.Printf("[PAYMENT][WEBHOOK] Payment not found for ref=%s: %v", txnRef, err)
			rspCode = "01"
			rspMessage = "Order Not Found"
			return nil
		}

		auditLog := domain.NewWebhookReceivedLog(&payment.ID, &payment.BookingID, webhookData)
		if err := tx.Create(auditLog).Error; err != nil {
			return fmt.Errorf("failed to create webhook audit log: %w", err)
		}

		if payment.Status == domain.PaymentStatusPaid || payment.Status == domain.PaymentStatusRefunded {
			rspCode = "02"
			rspMessage = "Order Already Confirmed"
			updatedPayment = payment
			return nil
		}

		if !s.isValidWebhookAmount(payment, params["vnp_Amount"]) {
			log.Printf("[PAYMENT][WEBHOOK] Amount mismatch ref=%s expected=%d got=%s", txnRef, payment.Amount, params["vnp_Amount"])
			rspCode = "04"
			rspMessage = "Invalid Amount"
			return nil
		}

		responseCode := params["vnp_ResponseCode"]
		transactionStatus := params["vnp_TransactionStatus"]
		transactionNo := params["vnp_TransactionNo"]
		isSuccess, message := s.gateway.ParseResponseCode(responseCode)
		isSuccess = isSuccess && transactionStatus == "00"

		if isSuccess {
			if err := s.handlePaymentSuccessTx(tx, payment, responseCode, transactionNo, message); err != nil {
				return err
			}
			completed = true
		} else {
			if err := s.handlePaymentFailureTx(tx, payment, responseCode, message); err != nil {
				return err
			}
		}
		updatedPayment = payment
		return nil
	})
	if err != nil {
		log.Printf("[PAYMENT][WEBHOOK] Failed to process ref=%s: %v", txnRef, err)
		return "99", "Unknown Error"
	}

	if completed && updatedPayment != nil {
		s.sendPaymentConfirmationEmail(updatedPayment)
		s.sendPaymentSuccessNotification(updatedPayment)
	}

	return rspCode, rspMessage
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
func (s *PaymentService) handlePaymentSuccessTx(tx *gorm.DB, payment *domain.Payment, responseCode, transactionNo, message string) error {
	previousStatus := payment.Status

	payment.SetVNPayResponse(transactionNo, responseCode, message)
	payment.Status = domain.PaymentStatusPaid
	payment.UpdatedAt = time.Now()
	method := s.gateway.ProviderName()
	payment.PaymentMethod = &method

	if err := s.repo.UpdatePaymentWithTransaction(tx, payment); err != nil {
		return fmt.Errorf("failed to update payment %d: %w", payment.ID, err)
	}

	// Cập nhật booking status
	if err := tx.Model(&domain.Booking{}).Where("id = ?", payment.BookingID).Updates(map[string]interface{}{
		"payment_status": "paid",
		"status":         "confirmed",
	}).Error; err != nil {
		return fmt.Errorf("failed to update booking payment status: %w", err)
	}

	// Audit log
	auditLog := domain.NewPaymentCompletedLog(payment, responseCode, message)
	if err := tx.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to create completed audit log: %w", err)
	}

	log.Printf("[PAYMENT][SUCCESS] payment_id=%d ref=%s prev_status=%s new_status=paid",
		payment.ID, payment.TransactionReference, previousStatus)
	return nil
}

// handlePaymentFailure - Xử lý khi thanh toán thất bại
func (s *PaymentService) handlePaymentFailureTx(tx *gorm.DB, payment *domain.Payment, responseCode, message string) error {
	previousStatus := payment.Status

	payment.SetVNPayResponse("", responseCode, message)
	payment.Status = domain.PaymentStatusFailed
	payment.UpdatedAt = time.Now()

	if err := s.repo.UpdatePaymentWithTransaction(tx, payment); err != nil {
		return fmt.Errorf("failed to update payment %d: %w", payment.ID, err)
	}

	// Cập nhật booking
	if err := tx.Model(&domain.Booking{}).Where("id = ?", payment.BookingID).Updates(map[string]interface{}{
		"payment_status": "failed",
	}).Error; err != nil {
		return fmt.Errorf("failed to update booking payment status: %w", err)
	}

	// Audit log
	auditLog := domain.NewPaymentFailedLog(payment, message, responseCode)
	if err := tx.Create(auditLog).Error; err != nil {
		return fmt.Errorf("failed to create failed audit log: %w", err)
	}

	log.Printf("[PAYMENT][FAILED] payment_id=%d ref=%s code=%s prev=%s",
		payment.ID, payment.TransactionReference, responseCode, previousStatus)
	return nil
}

func (s *PaymentService) isValidWebhookAmount(payment *domain.Payment, rawAmount string) bool {
	vnpayAmount, err := strconv.ParseInt(rawAmount, 10, 64)
	if err != nil {
		return false
	}
	return vnpayAmount == payment.Amount*100
}

func (s *PaymentService) sendPaymentSuccessNotification(payment *domain.Payment) {
	go func(bookingID uint, amount int64) {
		var b domain.Booking
		if err := database.DB.Preload("Tour").First(&b, bookingID).Error; err == nil {
			msg := fmt.Sprintf("Thanh toán thành công cho tour %s. Số tiền: %s.", b.Tour.Name, domain.FormatPriceVND(amount))
			_ = notification.SendNotification(b.UserID, "Thanh toán thành công", msg, domain.NotifTypePayment)
		}
	}(payment.BookingID, payment.Amount)
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
