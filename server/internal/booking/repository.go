package booking

import (
	"errors"
	"strings"
	"time"
	"travel-backend/database"
	"travel-backend/domain"

	"gorm.io/gorm"
)

var ErrBookingAlreadyCancelled = errors.New("booking already cancelled")

// CreateBooking - Tạo bản ghi đặt tour mới
func createBookingRecord(booking *domain.Booking) error {
	return database.DB.Create(booking).Error
}

// FindBookingsByUserID - Lấy danh sách tour đã đặt theo user
func FindBookingsByUserID(userID uint) ([]domain.Booking, error) {
	var bookings []domain.Booking
	err := database.DB.
		Preload("Tour").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&bookings).Error

	if err != nil {
		return nil, err
	}

	return bookings, nil
}

// FindBookingByID - Lấy thông tin chi tiết 1 booking theo ID
func FindBookingByID(bookingID uint) (*domain.Booking, error) {
	var booking domain.Booking
	err := database.DB.
		Preload("Tour").
		First(&booking, bookingID).Error

	if err != nil {
		return nil, err
	}

	return &booking, nil
}

// FindBookingByCode - Tìm booking theo mã booking_code
func FindBookingByCode(bookingCode string) (*domain.Booking, error) {
	var booking domain.Booking
	err := database.DB.
		Preload("Tour").
		Where("booking_code = ?", bookingCode).
		First(&booking).Error

	if err != nil {
		return nil, err
	}

	return &booking, nil
}


// CancelBookingByUserID - Hủy booking của đúng người dùng và hoàn lại số chỗ cho tour.
func cancelBookingByUserID(userID, bookingID uint) (*domain.Booking, error) {
	var booking domain.Booking

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ? AND user_id = ?", bookingID, userID).First(&booking).Error; err != nil {
			return err
		}

		status := strings.ToLower(strings.TrimSpace(booking.Status))
		if status == "cancelled" || status == "canceled" {
			return ErrBookingAlreadyCancelled
		}

		booking.Status = "cancelled"
		booking.PaymentStatus = "cancelled"
		if err := tx.Save(&booking).Error; err != nil {
			return err
		}

		var tour domain.Tour
		if err := tx.First(&tour, booking.TourID).Error; err != nil {
			return err
		}

		remaining := tour.RemainingSlots
		if remaining < 0 {
			remaining = 0
		}

		newSlots := remaining + booking.Quantity
		if err := tx.Model(&tour).Update("remaining_slots", newSlots).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if err := database.DB.Preload("Tour").First(&booking, booking.ID).Error; err != nil {
		return nil, err
	}

	return &booking, nil
}
// updateBookingPaymentDeadline - Cập nhật deadline thanh toán cho booking
func updateBookingPaymentDeadline(bookingID uint, deadline *time.Time) error {
	return database.DB.Model(&domain.Booking{}).
		Where("id = ?", bookingID).
		Update("payment_deadline", deadline).Error
}