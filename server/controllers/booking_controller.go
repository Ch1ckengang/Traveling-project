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
			errors.Is(err, services.ErrInvalidAdultCount) ||
			errors.Is(err, services.ErrInvalidChildCount) ||
			errors.Is(err, services.ErrInvalidInfantCount) ||
			errors.Is(err, services.ErrInvalidQuantity) ||
			errors.Is(err, services.ErrInsufficientSlots) ||
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
		Message: "Đặt tour thành công. Vui lòng thanh toán trong vòng 24 giờ để giữ chỗ.",
		Booking: booking,
	})
}

// GetUserBookings - Xử lý GET /api/users/:id/bookings
func GetUserBookings(c *gin.Context) {
	userID := c.Param("id")

	bookings, err := services.GetBookingsByUserID(userID)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, services.ErrInvalidUserID) {
			status = http.StatusBadRequest
		}

		c.JSON(status, models.BookingListResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.BookingListResponse{
		Success:  true,
		Bookings: bookings,
	})
}

// CancelBooking - Xử lý PUT /api/users/:id/bookings/:bookingId/cancel
func CancelBooking(c *gin.Context) {
	userID := c.Param("id")
	bookingID := c.Param("bookingId")

	booking, err := services.CancelBookingByUserID(userID, bookingID)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, services.ErrInvalidUserID) || errors.Is(err, services.ErrInvalidBookingID) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, services.ErrBookingNotFound) {
			status = http.StatusNotFound
		}

		if errors.Is(err, services.ErrBookingCannotCancel) {
			status = http.StatusConflict
		}

		c.JSON(status, models.BookingResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.BookingResponse{
		Success: true,
		Message: "Hủy tour thành công.",
		Booking: booking,
	})
}
