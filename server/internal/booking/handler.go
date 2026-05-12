package booking

import (
	"errors"
	"net/http"
	"time"
	"travel-backend/domain"
	authmodule "travel-backend/internal/auth"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// CreateBookingHandler - Xử lý POST /api/bookings
func CreateBookingHandler(c *gin.Context) {
	authUserID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.BookingResponse{
			Success: false,
			Message: shared.ErrInvalidAccessToken.Error(),
		})
		return
	}

	var req domain.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, domain.BookingResponse{
			Success: false,
			Message: "Dữ liệu đặt tour không hợp lệ",
		})
		return
	}

	req.UserID = authUserID

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

	// Convert to BookingWithPayment and set payment deadline (24 hours from now)
	bookingWithPayment := booking.ToBookingWithPayment()
	
	// Set payment deadline to 24 hours from now
	paymentDeadline := time.Now().Add(24 * time.Hour)
	booking.PaymentDeadline = &paymentDeadline
	bookingWithPayment.PaymentDeadline = &paymentDeadline
	
	// Update the booking in database with payment deadline
	if err := updateBookingPaymentDeadline(booking.ID, &paymentDeadline); err != nil {
		// Log error but don't fail the request
		// The booking was created successfully, just the deadline wasn't set
	}

	c.JSON(http.StatusCreated, domain.BookingResponse{
		Success: true,
		Message: "Đặt tour thành công. Vui lòng thanh toán trong vòng 24 giờ để giữ chỗ.",
		Booking: bookingWithPayment,
	})
}

// GetUserBookingsHandler - Xử lý GET /api/users/:id/bookings
func GetUserBookingsHandler(c *gin.Context) {
	if !authmodule.EnsurePathUserMatchesToken(c, "id") {
		return
	}

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

	// Convert bookings to BookingWithPayment
	bookingsWithPayment := make([]*domain.BookingWithPayment, len(bookings))
	for i, booking := range bookings {
		bookingsWithPayment[i] = booking.ToBookingWithPayment()
		// TODO: Load payment info for each booking if needed
		// This can be optimized later with batch loading
	}

	c.JSON(http.StatusOK, domain.BookingListResponse{
		Success:  true,
		Bookings: bookingsWithPayment,
	})
}

// GetBookingByIDHandler - Xử lý GET /api/bookings/:id
func GetBookingByIDHandler(c *gin.Context) {
	authUserID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.BookingResponse{
			Success: false,
			Message: shared.ErrInvalidAccessToken.Error(),
		})
		return
	}

	bookingID := c.Param("id")

	booking, err := GetBookingByID(bookingID)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, shared.ErrInvalidBookingID) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, shared.ErrBookingNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, domain.BookingResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Check ownership: user chỉ được xem booking của mình
	if booking.UserID != authUserID {
		c.JSON(http.StatusForbidden, domain.BookingResponse{
			Success: false,
			Message: shared.ErrForbiddenResource.Error(),
		})
		return
	}

	// Convert to BookingWithPayment
	bookingWithPayment := booking.ToBookingWithPayment()
	// TODO: Load payment info and history for this booking
	// This can be implemented when payment service is available

	c.JSON(http.StatusOK, domain.BookingResponse{
		Success: true,
		Booking: bookingWithPayment,
	})
}

// CancelBookingHandler - Xử lý PUT /api/users/:id/bookings/:bookingId/cancel
func CancelBookingHandler(c *gin.Context) {
	if !authmodule.EnsurePathUserMatchesToken(c, "id") {
		return
	}

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

	// Convert to BookingWithPayment
	bookingWithPayment := booking.ToBookingWithPayment()

	c.JSON(http.StatusOK, domain.BookingResponse{
		Success: true,
		Message: "Hủy tour thành công.",
		Booking: bookingWithPayment,
	})
}

// GetBookingByCodeHandler - Xử lý GET /api/bookings/code/:code
func GetBookingByCodeHandler(c *gin.Context) {
	authUserID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, domain.BookingResponse{
			Success: false,
			Message: shared.ErrInvalidAccessToken.Error(),
		})
		return
	}

	bookingCode := c.Param("code")

	booking, err := GetBookingByCode(bookingCode)
	if err != nil {
		status := http.StatusInternalServerError

		if errors.Is(err, shared.ErrInvalidBookingID) {
			status = http.StatusBadRequest
		}

		if errors.Is(err, shared.ErrBookingNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, domain.BookingResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Check ownership
	if booking.UserID != authUserID {
		c.JSON(http.StatusForbidden, domain.BookingResponse{
			Success: false,
			Message: shared.ErrForbiddenResource.Error(),
		})
		return
	}

	bookingWithPayment := booking.ToBookingWithPayment()

	c.JSON(http.StatusOK, domain.BookingResponse{
		Success: true,
		Booking: bookingWithPayment,
	})
}

