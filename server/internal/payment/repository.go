package payment

import (
	"travel-backend/domain"

	"gorm.io/gorm"
)

// PaymentRepository defines all database operations for payment management
// This interface provides methods for creating, retrieving, and updating payment records,
// as well as managing payment audit logs with support for transactions and concurrent processing.
type PaymentRepository interface {
	// Payment CRUD Operations

	// CreatePayment creates a new payment record in the database
	// Returns error if booking_id is invalid or database operation fails
	CreatePayment(payment *domain.Payment) error

	// GetPaymentByID retrieves a payment by its primary key ID
	// Returns nil if payment not found
	GetPaymentByID(paymentID uint) (*domain.Payment, error)

	// GetPaymentByTransactionReference retrieves a payment by its unique transaction reference
	// This is used for webhook processing and status checks
	// Returns nil if payment not found
	GetPaymentByTransactionReference(txnRef string) (*domain.Payment, error)

	// GetPaymentByVNPayTransactionID retrieves a payment by VNPay's transaction ID
	// Used for reconciliation and duplicate detection
	// Returns nil if payment not found
	GetPaymentByVNPayTransactionID(vnpayTxnID string) (*domain.Payment, error)

	// GetPaymentsByBookingID retrieves all payments associated with a booking
	// Returns empty slice if no payments found
	// Ordered by created_at DESC to show most recent first
	GetPaymentsByBookingID(bookingID uint) ([]domain.Payment, error)

	// GetPaymentsByUserID retrieves all payments for a specific user through their bookings
	// Returns empty slice if no payments found
	// Ordered by created_at DESC
	GetPaymentsByUserID(userID uint) ([]domain.Payment, error)

	// UpdatePaymentStatus updates the payment status with validation
	// Performs status transition validation before updating
	// Updates updated_at timestamp automatically
	// Returns error if status transition is invalid
	UpdatePaymentStatus(paymentID uint, newStatus string) error

	// UpdatePayment updates a payment record with all fields
	// Used for updating VNPay response data, payment method, etc.
	// Updates updated_at timestamp automatically
	UpdatePayment(payment *domain.Payment) error

	// UpdatePaymentWithTransaction updates a payment within a database transaction
	// Used for atomic operations that need to update both payment and booking
	// Caller is responsible for committing or rolling back the transaction
	UpdatePaymentWithTransaction(tx *gorm.DB, payment *domain.Payment) error

	// SoftDeletePayment marks a payment as deleted (sets deleted_at timestamp)
	// Maintains audit trail by not physically deleting the record
	SoftDeletePayment(paymentID uint) error

	// Payment Query Operations

	// GetPaymentsByStatus retrieves all payments with a specific status
	// Supports pagination with limit and offset
	// Returns empty slice if no payments found
	GetPaymentsByStatus(status string, limit, offset int) ([]domain.Payment, error)

	// GetExpiredPayments retrieves all payments where session_expires_at has passed
	// and status is still pending or processing
	// Used for cleanup jobs to mark expired payments
	GetExpiredPayments() ([]domain.Payment, error)

	// GetPaymentsByDateRange retrieves payments created within a date range
	// Used for reporting and analytics
	// Supports pagination with limit and offset
	GetPaymentsByDateRange(startDate, endDate string, limit, offset int) ([]domain.Payment, error)

	// CountPaymentsByStatus returns the total count of payments with a specific status
	// Used for dashboard statistics
	CountPaymentsByStatus(status string) (int64, error)

	// Audit Log Operations

	// CreateAuditLog creates a new payment audit log entry
	// All payment operations should generate corresponding audit logs
	CreateAuditLog(log *domain.PaymentAuditLog) error

	// GetAuditLogsByPaymentID retrieves all audit logs for a specific payment
	// Ordered by timestamp DESC to show most recent events first
	GetAuditLogsByPaymentID(paymentID uint) ([]domain.PaymentAuditLog, error)

	// GetAuditLogsByBookingID retrieves all audit logs for a specific booking
	// Includes logs from all payment attempts for that booking
	// Ordered by timestamp DESC
	GetAuditLogsByBookingID(bookingID uint) ([]domain.PaymentAuditLog, error)

	// GetAuditLogsByEventType retrieves audit logs filtered by event type
	// Supports pagination with limit and offset
	// Used for security monitoring and event analysis
	GetAuditLogsByEventType(eventType string, limit, offset int) ([]domain.PaymentAuditLog, error)

	// GetAuditLogsByDateRange retrieves audit logs within a date range
	// Used for compliance reporting and audit trail exports
	// Supports pagination with limit and offset
	GetAuditLogsByDateRange(startDate, endDate string, limit, offset int) ([]domain.PaymentAuditLog, error)

	// GetSecurityViolationLogs retrieves all security violation audit logs
	// Used for security monitoring and incident response
	// Ordered by timestamp DESC
	GetSecurityViolationLogs(limit, offset int) ([]domain.PaymentAuditLog, error)

	// Transaction Support

	// BeginTransaction starts a new database transaction
	// Used for operations that need to update multiple tables atomically
	// Caller must call Commit() or Rollback() on the returned transaction
	BeginTransaction() (*gorm.DB, error)

	// WithTransaction executes a function within a database transaction
	// Automatically commits on success or rolls back on error
	// This is the preferred method for transactional operations
	WithTransaction(fn func(tx *gorm.DB) error) error

	// Concurrency Support

	// LockPaymentForUpdate locks a payment record for update (SELECT FOR UPDATE)
	// Used to prevent concurrent modifications and race conditions
	// Must be called within a transaction
	// Returns error if payment not found or already locked
	LockPaymentForUpdate(tx *gorm.DB, paymentID uint) (*domain.Payment, error)

	// GetPaymentWithLock retrieves and locks a payment by transaction reference
	// Combines GetPaymentByTransactionReference with row-level locking
	// Must be called within a transaction
	GetPaymentWithLock(tx *gorm.DB, txnRef string) (*domain.Payment, error)
}
