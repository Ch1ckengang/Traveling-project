package payment

import (
	"errors"
	"fmt"
	"time"
	"travel-backend/database"
	"travel-backend/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// paymentRepositoryImpl is the concrete implementation of PaymentRepository interface
// It uses GORM for database operations and follows the repository pattern
type paymentRepositoryImpl struct {
	db *gorm.DB
}

// NewPaymentRepository creates a new instance of PaymentRepository
// Uses the global database connection from the database package
func NewPaymentRepository() PaymentRepository {
	return &paymentRepositoryImpl{
		db: database.DB,
	}
}

// NewPaymentRepositoryWithDB creates a new instance of PaymentRepository with a custom DB connection
// Useful for testing with a different database connection
func NewPaymentRepositoryWithDB(db *gorm.DB) PaymentRepository {
	return &paymentRepositoryImpl{
		db: db,
	}
}

// CreatePayment creates a new payment record in the database
// Validates booking_id exists and is valid before creating the payment
// Sets timestamps automatically and supports atomic transactions
func (r *paymentRepositoryImpl) CreatePayment(payment *domain.Payment) error {
	if payment == nil {
		return errors.New("payment cannot be nil")
	}

	// Validate booking exists
	var booking domain.Booking
	if err := r.db.First(&booking, payment.BookingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("booking with id %d not found", payment.BookingID)
		}
		return fmt.Errorf("failed to validate booking: %w", err)
	}

	// Validate payment amount
	if err := payment.ValidateAmount(); err != nil {
		return fmt.Errorf("invalid payment amount: %w", err)
	}

	// Validate transaction reference format
	if err := payment.ValidateTransactionReference(); err != nil {
		return fmt.Errorf("invalid transaction reference: %w", err)
	}

	// Set timestamps if not already set
	now := time.Now()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = now
	}

	// Create payment record
	if err := r.db.Create(payment).Error; err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

// GetPaymentByID retrieves a payment by its primary key ID
func (r *paymentRepositoryImpl) GetPaymentByID(paymentID uint) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.First(&payment, paymentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by id: %w", err)
	}
	return &payment, nil
}

// GetPaymentByTransactionReference retrieves a payment by its unique transaction reference
func (r *paymentRepositoryImpl) GetPaymentByTransactionReference(txnRef string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("transaction_reference = ?", txnRef).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by transaction reference: %w", err)
	}
	return &payment, nil
}

// GetPaymentByVNPayTransactionID retrieves a payment by VNPay's transaction ID
func (r *paymentRepositoryImpl) GetPaymentByVNPayTransactionID(vnpayTxnID string) (*domain.Payment, error) {
	var payment domain.Payment
	err := r.db.Where("vnpay_transaction_id = ?", vnpayTxnID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get payment by vnpay transaction id: %w", err)
	}
	return &payment, nil
}

// GetPaymentsByBookingID retrieves all payments associated with a booking
func (r *paymentRepositoryImpl) GetPaymentsByBookingID(bookingID uint) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.db.Where("booking_id = ?", bookingID).
		Order("created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by booking id: %w", err)
	}
	return payments, nil
}

// GetPaymentsByUserID retrieves all payments for a specific user through their bookings
func (r *paymentRepositoryImpl) GetPaymentsByUserID(userID uint) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.db.Joins("JOIN bookings ON bookings.id = payments.booking_id").
		Where("bookings.user_id = ?", userID).
		Order("payments.created_at DESC").
		Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by user id: %w", err)
	}
	return payments, nil
}

// UpdatePaymentStatus updates the payment status with validation
func (r *paymentRepositoryImpl) UpdatePaymentStatus(paymentID uint, newStatus string) error {
	// Get current payment
	payment, err := r.GetPaymentByID(paymentID)
	if err != nil {
		return err
	}
	if payment == nil {
		return fmt.Errorf("payment with id %d not found", paymentID)
	}

	// Validate status transition
	if err := payment.ValidateStatusTransition(newStatus); err != nil {
		return fmt.Errorf("invalid status transition: %w", err)
	}

	// Update status
	err = r.db.Model(&domain.Payment{}).
		Where("id = ?", paymentID).
		Updates(map[string]interface{}{
			"status":     newStatus,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	return nil
}

// UpdatePayment updates a payment record with all fields
func (r *paymentRepositoryImpl) UpdatePayment(payment *domain.Payment) error {
	if payment == nil {
		return errors.New("payment cannot be nil")
	}

	payment.UpdatedAt = time.Now()

	err := r.db.Save(payment).Error
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	return nil
}

// UpdatePaymentWithTransaction updates a payment within a database transaction
func (r *paymentRepositoryImpl) UpdatePaymentWithTransaction(tx *gorm.DB, payment *domain.Payment) error {
	if payment == nil {
		return errors.New("payment cannot be nil")
	}

	payment.UpdatedAt = time.Now()

	err := tx.Save(payment).Error
	if err != nil {
		return fmt.Errorf("failed to update payment in transaction: %w", err)
	}

	return nil
}

// SoftDeletePayment marks a payment as deleted (sets deleted_at timestamp)
func (r *paymentRepositoryImpl) SoftDeletePayment(paymentID uint) error {
	err := r.db.Delete(&domain.Payment{}, paymentID).Error
	if err != nil {
		return fmt.Errorf("failed to soft delete payment: %w", err)
	}
	return nil
}

// GetPaymentsByStatus retrieves all payments with a specific status
func (r *paymentRepositoryImpl) GetPaymentsByStatus(status string, limit, offset int) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.db.Where("status = ?", status).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by status: %w", err)
	}
	return payments, nil
}

// GetExpiredPayments retrieves all payments where session_expires_at has passed
func (r *paymentRepositoryImpl) GetExpiredPayments() ([]domain.Payment, error) {
	var payments []domain.Payment
	now := time.Now()
	err := r.db.Where("session_expires_at < ? AND status IN (?)", now, []string{domain.PaymentStatusPending, domain.PaymentStatusProcessing}).
		Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get expired payments: %w", err)
	}
	return payments, nil
}

// GetPaymentsByDateRange retrieves payments created within a date range
func (r *paymentRepositoryImpl) GetPaymentsByDateRange(startDate, endDate string, limit, offset int) ([]domain.Payment, error) {
	var payments []domain.Payment
	err := r.db.Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&payments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get payments by date range: %w", err)
	}
	return payments, nil
}

// CountPaymentsByStatus returns the total count of payments with a specific status
func (r *paymentRepositoryImpl) CountPaymentsByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Payment{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count payments by status: %w", err)
	}
	return count, nil
}

// CreateAuditLog creates a new payment audit log entry
func (r *paymentRepositoryImpl) CreateAuditLog(log *domain.PaymentAuditLog) error {
	if log == nil {
		return errors.New("audit log cannot be nil")
	}

	// Validate event type
	if err := log.ValidateEventType(); err != nil {
		return fmt.Errorf("invalid event type: %w", err)
	}

	// Set timestamp if not already set
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	err := r.db.Create(log).Error
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return nil
}

// GetAuditLogsByPaymentID retrieves all audit logs for a specific payment
func (r *paymentRepositoryImpl) GetAuditLogsByPaymentID(paymentID uint) ([]domain.PaymentAuditLog, error) {
	var logs []domain.PaymentAuditLog
	err := r.db.Where("payment_id = ?", paymentID).
		Order("timestamp DESC").
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by payment id: %w", err)
	}
	return logs, nil
}

// GetAuditLogsByBookingID retrieves all audit logs for a specific booking
func (r *paymentRepositoryImpl) GetAuditLogsByBookingID(bookingID uint) ([]domain.PaymentAuditLog, error) {
	var logs []domain.PaymentAuditLog
	err := r.db.Where("booking_id = ?", bookingID).
		Order("timestamp DESC").
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by booking id: %w", err)
	}
	return logs, nil
}

// GetAuditLogsByEventType retrieves audit logs filtered by event type
func (r *paymentRepositoryImpl) GetAuditLogsByEventType(eventType string, limit, offset int) ([]domain.PaymentAuditLog, error) {
	var logs []domain.PaymentAuditLog
	err := r.db.Where("event_type = ?", eventType).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by event type: %w", err)
	}
	return logs, nil
}

// GetAuditLogsByDateRange retrieves audit logs within a date range
func (r *paymentRepositoryImpl) GetAuditLogsByDateRange(startDate, endDate string, limit, offset int) ([]domain.PaymentAuditLog, error) {
	var logs []domain.PaymentAuditLog
	err := r.db.Where("timestamp >= ? AND timestamp <= ?", startDate, endDate).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs by date range: %w", err)
	}
	return logs, nil
}

// GetSecurityViolationLogs retrieves all security violation audit logs
func (r *paymentRepositoryImpl) GetSecurityViolationLogs(limit, offset int) ([]domain.PaymentAuditLog, error) {
	var logs []domain.PaymentAuditLog
	err := r.db.Where("event_type = ?", domain.EventSecurityViolation).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		Find(&logs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get security violation logs: %w", err)
	}
	return logs, nil
}

// BeginTransaction starts a new database transaction
func (r *paymentRepositoryImpl) BeginTransaction() (*gorm.DB, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}
	return tx, nil
}

// WithTransaction executes a function within a database transaction
func (r *paymentRepositoryImpl) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// LockPaymentForUpdate locks a payment record for update (SELECT FOR UPDATE)
func (r *paymentRepositoryImpl) LockPaymentForUpdate(tx *gorm.DB, paymentID uint) (*domain.Payment, error) {
	var payment domain.Payment
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&payment, paymentID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment with id %d not found", paymentID)
		}
		return nil, fmt.Errorf("failed to lock payment for update: %w", err)
	}
	return &payment, nil
}

// GetPaymentWithLock retrieves and locks a payment by transaction reference
func (r *paymentRepositoryImpl) GetPaymentWithLock(tx *gorm.DB, txnRef string) (*domain.Payment, error) {
	var payment domain.Payment
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("transaction_reference = ?", txnRef).
		First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment with transaction reference %s not found", txnRef)
		}
		return nil, fmt.Errorf("failed to get payment with lock: %w", err)
	}
	return &payment, nil
}
