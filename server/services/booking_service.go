package services

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"travel-backend/models"
	"travel-backend/repositories"
	"unicode"

	"gorm.io/gorm"
)

var phoneRegexp = regexp.MustCompile(`^[0-9+\-\s]{8,20}$`)

// CreateBooking - Xử lý nghiệp vụ đặt tour
func CreateBooking(req models.CreateBookingRequest) (*models.Booking, error) {
	normalizedReq := normalizeBookingRequest(req)

	if err := validateBookingRequest(normalizedReq); err != nil {
		return nil, err
	}

	tour, err := repositories.FindTourByID(normalizedReq.TourID)
	if err != nil {
		return nil, ErrTourNotFound
	}

	if normalizedReq.Quantity > tour.RemainingSlots {
		return nil, ErrInsufficientSlots
	}

	totalAmount := calculateBookingTotal(extractPriceAmount(tour.Price), normalizedReq.AdultCount, normalizedReq.ChildCount)
	bookingCode := generateBookingCode()

	booking := &models.Booking{
		UserID:        normalizedReq.UserID,
		TourID:        normalizedReq.TourID,
		FullName:      normalizedReq.FullName,
		Phone:         normalizedReq.Phone,
		Email:         normalizedReq.Email,
		AdultCount:    normalizedReq.AdultCount,
		ChildCount:    normalizedReq.ChildCount,
		InfantCount:   normalizedReq.InfantCount,
		Quantity:      normalizedReq.Quantity,
		TravelDate:    normalizedReq.TravelDate,
		TotalAmount:   totalAmount,
		BookingCode:   bookingCode,
		PaymentStatus: "unpaid",
		Note:          normalizedReq.Note,
		Status:        "pending_payment",
	}

	if err := repositories.CreateBooking(booking); err != nil {
		return nil, ErrCreateBookingFailed
	}

	if err := repositories.DecreaseTourRemainingSlots(normalizedReq.TourID, normalizedReq.Quantity); err != nil {
		if err == gorm.ErrInvalidData {
			return nil, ErrInsufficientSlots
		}
		return nil, ErrCreateBookingFailed
	}

	return booking, nil
}

// GetBookingsByUserID - Lấy lịch sử tour đã đặt của người dùng
func GetBookingsByUserID(userID string) ([]models.Booking, error) {
	parsedUserID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil || parsedUserID == 0 {
		return nil, ErrInvalidUserID
	}

	bookings, err := repositories.FindBookingsByUserID(uint(parsedUserID))
	if err != nil {
		return nil, ErrFetchBookingsFailed
	}

	return bookings, nil
}

// CancelBookingByUserID - Hủy tour đã đặt theo user và booking id.
func CancelBookingByUserID(userID, bookingID string) (*models.Booking, error) {
	parsedUserID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil || parsedUserID == 0 {
		return nil, ErrInvalidUserID
	}

	parsedBookingID, err := strconv.ParseUint(bookingID, 10, 32)
	if err != nil || parsedBookingID == 0 {
		return nil, ErrInvalidBookingID
	}

	booking, err := repositories.CancelBookingByUserID(uint(parsedUserID), uint(parsedBookingID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBookingNotFound
		}

		if errors.Is(err, repositories.ErrBookingAlreadyCancelled) {
			return nil, ErrBookingCannotCancel
		}

		return nil, ErrCancelBookingFailed
	}

	return booking, nil
}

func normalizeBookingRequest(req models.CreateBookingRequest) models.CreateBookingRequest {
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.TravelDate = strings.TrimSpace(req.TravelDate)
	req.Note = strings.TrimSpace(req.Note)

	// Backward compatibility: legacy clients only send quantity.
	if req.AdultCount == 0 && req.ChildCount == 0 && req.Quantity > 0 {
		req.AdultCount = req.Quantity
	}

	req.Quantity = req.AdultCount + req.ChildCount + req.InfantCount
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

	if req.AdultCount < 1 {
		return ErrInvalidAdultCount
	}

	if req.ChildCount < 0 {
		return ErrInvalidChildCount
	}

	if req.InfantCount < 0 {
		return ErrInvalidInfantCount
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

func extractPriceAmount(priceText string) int64 {
	b := strings.Builder{}
	for _, r := range priceText {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}

	if b.Len() == 0 {
		return 0
	}

	amount, err := strconv.ParseInt(b.String(), 10, 64)
	if err != nil {
		return 0
	}

	return amount
}

func calculateBookingTotal(basePrice int64, adultCount, childCount int) int64 {
	adultTotal := basePrice * int64(adultCount)
	childTotal := (basePrice * 75 / 100) * int64(childCount)
	return adultTotal + childTotal
}

func generateBookingCode() string {
	return fmt.Sprintf("TOUR-%s-%04d", time.Now().Format("20060102"), time.Now().UnixNano()%10000)
}
