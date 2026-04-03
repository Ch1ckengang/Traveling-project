package controllers

import (
	"errors"
	"net/http"
	"travel-backend/models"
	"travel-backend/services"

	"github.com/gin-gonic/gin"
)

// CreateBooking - Xử lý POST /api/bookings
func CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.BookingResponse{
			Success: false,
			Message: "Dữ liệu đặt tour không hợp lệ",
		})
		return
	}

	booking, err := services.CreateBooking(req)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, services.ErrInvalidBookingPayload) ||
			errors.Is(err, services.ErrInvalidFullName) ||
			errors.Is(err, services.ErrInvalidPhone) ||
			errors.Is(err, services.ErrInvalidEmail) ||
			errors.Is(err, services.ErrInvalidQuantity) ||
			errors.Is(err, services.ErrInvalidTravelDate) ||
			errors.Is(err, services.ErrTravelDateInPast) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, services.ErrTourNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, models.BookingResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, models.BookingResponse{
		Success: true,
		Message: "Đặt tour thành công",
		Booking: booking,
	})
}
