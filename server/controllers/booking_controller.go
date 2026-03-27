package controllers

import (
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
		if err.Error() == "Số lượng khách phải lớn hơn 0" || err.Error() == "Tour không tồn tại" {
			status = http.StatusBadRequest
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
