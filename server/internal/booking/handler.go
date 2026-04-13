package booking

import (
	"errors"
	"net/http"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// CreateBookingHandler - Xử lý POST /api/bookings
func CreateBookingHandler(c *gin.Context) {
	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.BookingResponse{
			Success: false,
			Message: "Dữ liệu đặt tour không hợp lệ",
		})
		return
	}

	booking, err := CreateBooking(req)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, shared.ErrInvalidBookingPayload) ||
			errors.Is(err, shared.ErrInvalidFullName) ||
			errors.Is(err, shared.ErrInvalidPhone) ||
			errors.Is(err, shared.ErrInvalidEmail) ||
			errors.Is(err, shared.ErrInvalidAdultCount) ||
			errors.Is(err, shared.ErrInvalidChildCount) ||
			errors.Is(err, shared.ErrInvalidInfantCount) ||
			errors.Is(err, shared.ErrInvalidQuantity) ||
			errors.Is(err, shared.ErrInsufficientSlots) ||
			errors.Is(err, shared.ErrInvalidTravelDate) ||
			errors.Is(err, shared.ErrTravelDateInPast) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, shared.ErrTourNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, domain.BookingResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, domain.BookingResponse{
		Success: true,
		Message: "Đặt tour thành công. Vui lòng thanh toán trong vòng 24 giờ để giữ chỗ.",
		Booking: booking,
	})
}

// GetUserBookingsHandler - Xử lý GET /api/users/:id/bookings
func GetUserBookingsHandler(c *gin.Context) {
	userID := c.Param("id")

	bookings, err := GetBookingsByUserID(userID)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, shared.ErrInvalidUserID) {
			status = http.StatusBadRequest
		}

		c.JSON(status, domain.BookingListResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.BookingListResponse{
		Success:  true,
		Bookings: bookings,
	})
}

// CancelBookingHandler - Xử lý PUT /api/users/:id/bookings/:bookingId/cancel
func CancelBookingHandler(c *gin.Context) {
	userID := c.Param("id")
	bookingID := c.Param("bookingId")

	booking, err := CancelBookingByUserID(userID, bookingID)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, shared.ErrInvalidUserID) || errors.Is(err, shared.ErrInvalidBookingID) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, shared.ErrBookingNotFound) {
			status = http.StatusNotFound
		}

		if errors.Is(err, shared.ErrBookingCannotCancel) {
			status = http.StatusConflict
		}

		c.JSON(status, domain.BookingResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, domain.BookingResponse{
		Success: true,
		Message: "Hủy tour thành công.",
		Booking: booking,
	})
}
