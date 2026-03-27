package repositories

import (
	"travel-backend/database"
	"travel-backend/models"
)

// CreateBooking - Tạo bản ghi đặt tour mới
func CreateBooking(booking *models.Booking) error {
	return database.DB.Create(booking).Error
}
