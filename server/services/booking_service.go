package services

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
	"travel-backend/models"
	"travel-backend/repositories"
)

var phoneRegexp = regexp.MustCompile(`^[0-9+\-\s]{8,20}$`)

// CreateBooking - Xử lý nghiệp vụ đặt tour
func CreateBooking(req models.CreateBookingRequest) (*models.Booking, error) {
	normalizedReq := normalizeBookingRequest(req)

	if err := validateBookingRequest(normalizedReq); err != nil {
		return nil, err
	}

	if _, err := repositories.FindTourByID(normalizedReq.TourID); err != nil {
		return nil, ErrTourNotFound
	}

	booking := &models.Booking{
		UserID:     normalizedReq.UserID,
		TourID:     normalizedReq.TourID,
		FullName:   normalizedReq.FullName,
		Phone:      normalizedReq.Phone,
		Email:      normalizedReq.Email,
		Quantity:   normalizedReq.Quantity,
		TravelDate: normalizedReq.TravelDate,
		Note:       normalizedReq.Note,
		Status:     "booked",
	}

	if err := repositories.CreateBooking(booking); err != nil {
		return nil, ErrCreateBookingFailed
	}

	return booking, nil
}

func normalizeBookingRequest(req models.CreateBookingRequest) models.CreateBookingRequest {
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.TravelDate = strings.TrimSpace(req.TravelDate)
	req.Note = strings.TrimSpace(req.Note)
	return req
}

func validateBookingRequest(req models.CreateBookingRequest) error {
	if req.TourID == 0 {
		return ErrInvalidBookingPayload
	}

	if req.FullName == "" {
		return ErrInvalidFullName
	}

	if !phoneRegexp.MatchString(req.Phone) {
		return ErrInvalidPhone
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return ErrInvalidEmail
	}

	if req.Quantity <= 0 {
		return ErrInvalidQuantity
	}

	travelDate, err := time.Parse("2006-01-02", req.TravelDate)
	if err != nil {
		return ErrInvalidTravelDate
	}

	startOfToday := time.Now().Truncate(24 * time.Hour)
	if travelDate.Before(startOfToday) {
		return ErrTravelDateInPast
	}

	return nil
}
