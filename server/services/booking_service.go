package services

import (
	"errors"
	"travel-backend/models"
	"travel-backend/repositories"
)

// CreateBooking - Xử lý nghiệp vụ đặt tour
func CreateBooking(req models.CreateBookingRequest) (*models.Booking, error) {
	if req.Quantity <= 0 {
		return nil, errors.New("Số lượng khách phải lớn hơn 0")
	}

	if _, err := repositories.FindTourByID(req.TourID); err != nil {
		return nil, errors.New("Tour không tồn tại")
	}

	booking := &models.Booking{
		UserID:     req.UserID,
		TourID:     req.TourID,
		FullName:   req.FullName,
		Phone:      req.Phone,
		Email:      req.Email,
		Quantity:   req.Quantity,
		TravelDate: req.TravelDate,
		Note:       req.Note,
		Status:     "booked",
	}

	if err := repositories.CreateBooking(booking); err != nil {
		return nil, errors.New("Không thể tạo đơn đặt tour")
	}

	return booking, nil
}
