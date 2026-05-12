package payment

import (
	"testing"
	"travel-backend/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err, "Failed to open test database")

	// Auto-migrate the schema
	err = db.AutoMigrate(&domain.Payment{}, &domain.PaymentAuditLog{}, &domain.Booking{}, &domain.User{}, &domain.Tour{})
	require.NoError(t, err, "Failed to migrate test database")

	return db
}

// createTestBooking creates a test booking in the database
func createTestBooking(t *testing.T, db *gorm.DB) *domain.Booking {
	// Create a test user first
	user := &domain.User{
		Name:            "Test User",
		Email:           "test@example.com",
		Password:        "hashedpassword",
		IsEmailVerified: true,
	}
	err := db.Create(user).Error
	require.NoError(t, err, "Failed to create test user")

	// Create a test tour
	tour := &domain.Tour{
		Name:           "Test Tour",
		Type:           "domestic",
		Price:          "1000000",
		Description:    "Test Description",
		Duration:       "3",
		RemainingSlots: 10,
		Location:       "Test Location",
		Country:        "Việt Nam",
	}
	err = db.Create(tour).Error
	require.NoError(t, err, "Failed to create test tour")

	// Create a test booking
	booking := &domain.Booking{
		UserID:        user.ID,
		TourID:        tour.ID,
		Quantity:      2,
		TotalAmount:   2000000,
		Status:        "booked",
		PaymentStatus: "unpaid",
	}
	err = db.Create(booking).Error
	require.NoError(t, err, "Failed to create test booking")

	return booking
}

func TestCreatePayment_Success(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking
	booking := createTestBooking(t, db)

	// Create payment
	payment := domain.NewPayment(booking.ID, 2000000)

	err := repo.CreatePayment(payment)
	assert.NoError(t, err, "CreatePayment should succeed")
	assert.NotZero(t, payment.ID, "Payment ID should be set after creation")
	assert.Equal(t, booking.ID, payment.BookingID, "Booking ID should match")
	assert.Equal(t, int64(2000000), payment.Amount, "Amount should match")
	assert.Equal(t, domain.PaymentStatusPending, payment.Status, "Status should be pending")
	assert.NotEmpty(t, payment.TransactionReference, "Transaction reference should be set")
}

func TestCreatePayment_InvalidBookingID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create payment with non-existent booking ID
	payment := domain.NewPayment(999, 2000000)

	err := repo.CreatePayment(payment)
	assert.Error(t, err, "CreatePayment should fail with invalid booking ID")
	assert.Contains(t, err.Error(), "booking with id 999 not found", "Error should mention booking not found")
}

func TestCreatePayment_NilPayment(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	err := repo.CreatePayment(nil)
	assert.Error(t, err, "CreatePayment should fail with nil payment")
	assert.Contains(t, err.Error(), "payment cannot be nil", "Error should mention nil payment")
}

func TestCreatePayment_InvalidAmount(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking
	booking := createTestBooking(t, db)

	// Create payment with invalid amount (too small)
	payment := domain.NewPayment(booking.ID, 100) // Less than minimum 5,000 VND (500,000 cents)

	err := repo.CreatePayment(payment)
	assert.Error(t, err, "CreatePayment should fail with invalid amount")
	assert.Contains(t, err.Error(), "invalid payment amount", "Error should mention invalid amount")
}

func TestCreatePayment_SetsTimestamps(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking
	booking := createTestBooking(t, db)

	// Create payment without timestamps
	payment := &domain.Payment{
		BookingID:            booking.ID,
		TransactionReference: domain.GenerateTransactionReference(),
		Amount:               2000000,
		Currency:             "VND",
		Status:               domain.PaymentStatusPending,
	}

	err := repo.CreatePayment(payment)
	assert.NoError(t, err, "CreatePayment should succeed")
	assert.False(t, payment.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, payment.UpdatedAt.IsZero(), "UpdatedAt should be set")
}

func TestCreatePayment_TransactionSupport(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking
	booking := createTestBooking(t, db)

	// Test transaction support by creating payment within a transaction
	err := repo.WithTransaction(func(tx *gorm.DB) error {
		payment := domain.NewPayment(booking.ID, 2000000)
		
		// Create payment within transaction
		if err := tx.Create(payment).Error; err != nil {
			return err
		}

		// Verify payment was created
		var count int64
		tx.Model(&domain.Payment{}).Where("booking_id = ?", booking.ID).Count(&count)
		assert.Equal(t, int64(1), count, "Payment should exist in transaction")

		return nil
	})

	assert.NoError(t, err, "Transaction should succeed")

	// Verify payment exists after transaction commit
	var count int64
	db.Model(&domain.Payment{}).Where("booking_id = ?", booking.ID).Count(&count)
	assert.Equal(t, int64(1), count, "Payment should exist after transaction commit")
}

func TestCreatePayment_ConcurrentCreation(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking
	booking := createTestBooking(t, db)

	// Create multiple payments concurrently
	done := make(chan bool, 3)
	errors := make(chan error, 3)
	paymentIDs := make(chan uint, 3)

	for i := 0; i < 3; i++ {
		go func() {
			payment := domain.NewPayment(booking.ID, 2000000)
			err := repo.CreatePayment(payment)
			if err != nil {
				errors <- err
			} else {
				paymentIDs <- payment.ID
			}
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 3; i++ {
		<-done
	}
	close(errors)
	close(paymentIDs)

	// Check for any errors
	errorCount := 0
	for err := range errors {
		if err != nil {
			t.Logf("Concurrent creation error: %v", err)
			errorCount++
		}
	}

	// Count successful payments
	successCount := 0
	for range paymentIDs {
		successCount++
	}

	// All payments should be created successfully
	assert.Equal(t, 3, successCount, "All 3 payments should be created successfully")
	assert.Equal(t, 0, errorCount, "No errors should occur during concurrent creation")

	// Verify in database
	var count int64
	db.Model(&domain.Payment{}).Where("booking_id = ?", booking.ID).Count(&count)
	assert.Equal(t, int64(3), count, "All 3 payments should exist in database")
}

func TestGetPaymentByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking and payment
	booking := createTestBooking(t, db)
	payment := domain.NewPayment(booking.ID, 2000000)
	err := repo.CreatePayment(payment)
	require.NoError(t, err)

	// Retrieve payment by ID
	retrieved, err := repo.GetPaymentByID(payment.ID)
	assert.NoError(t, err, "GetPaymentByID should succeed")
	assert.NotNil(t, retrieved, "Retrieved payment should not be nil")
	assert.Equal(t, payment.ID, retrieved.ID, "Payment IDs should match")
	assert.Equal(t, payment.TransactionReference, retrieved.TransactionReference, "Transaction references should match")
}

func TestGetPaymentByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Try to retrieve non-existent payment
	retrieved, err := repo.GetPaymentByID(999)
	assert.NoError(t, err, "GetPaymentByID should not return error for not found")
	assert.Nil(t, retrieved, "Retrieved payment should be nil when not found")
}

func TestGetPaymentByTransactionReference(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking and payment
	booking := createTestBooking(t, db)
	payment := domain.NewPayment(booking.ID, 2000000)
	err := repo.CreatePayment(payment)
	require.NoError(t, err)

	// Retrieve payment by transaction reference
	retrieved, err := repo.GetPaymentByTransactionReference(payment.TransactionReference)
	assert.NoError(t, err, "GetPaymentByTransactionReference should succeed")
	assert.NotNil(t, retrieved, "Retrieved payment should not be nil")
	assert.Equal(t, payment.ID, retrieved.ID, "Payment IDs should match")
	assert.Equal(t, payment.TransactionReference, retrieved.TransactionReference, "Transaction references should match")
}

func TestCreateAuditLog(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create test booking and payment
	booking := createTestBooking(t, db)
	payment := domain.NewPayment(booking.ID, 2000000)
	err := repo.CreatePayment(payment)
	require.NoError(t, err)

	// Create audit log
	auditLog := domain.NewPaymentAuditLog(domain.EventPaymentInitiated, &payment.ID, &booking.ID, nil)
	err = repo.CreateAuditLog(auditLog)
	assert.NoError(t, err, "CreateAuditLog should succeed")
	assert.NotZero(t, auditLog.ID, "Audit log ID should be set after creation")
}

func TestCreateAuditLog_InvalidEventType(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPaymentRepositoryWithDB(db)

	// Create audit log with invalid event type
	auditLog := &domain.PaymentAuditLog{
		EventType: "invalid_event_type",
	}

	err := repo.CreateAuditLog(auditLog)
	assert.Error(t, err, "CreateAuditLog should fail with invalid event type")
	assert.Contains(t, err.Error(), "invalid event type", "Error should mention invalid event type")
}
