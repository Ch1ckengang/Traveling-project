package booking

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
	"strings"
	"time"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/coupon"
	"travel-backend/internal/notification"
	"travel-backend/internal/shared"
	tourmodule "travel-backend/internal/tour"
	"unicode"

	"gorm.io/gorm"
)

var phoneRegexp = regexp.MustCompile(`^[0-9+\-\s]{8,20}$`)

// CreateBooking - Xử lý nghiệp vụ đặt tour
func CreateBooking(req domain.CreateBookingRequest) (*domain.Booking, error) {
	normalizedReq := normalizeBookingRequest(req)

	if err := validateBookingRequest(normalizedReq); err != nil {
		return nil, err
	}

	tour, err := tourmodule.FindTourByID(normalizedReq.TourID)
	if err != nil {
		return nil, shared.ErrTourNotFound
	}

	if normalizedReq.Quantity > tour.RemainingSlots {
		return nil, shared.ErrInsufficientSlots
	}

	totalAmount := calculateBookingTotal(tour.PriceAmount, normalizedReq.AdultCount, normalizedReq.ChildCount)
	var discountAmount int64 = 0

	// Handle coupon if provided
	if normalizedReq.CouponCode != "" {
		couponData, discount, err := coupon.ValidateCoupon(normalizedReq.CouponCode, totalAmount)
		if err == nil && couponData != nil {
			discountAmount = discount
			totalAmount = totalAmount - discountAmount
		} else {
			// If coupon is invalid, return error instead of proceeding
			return nil, fmt.Errorf("mã giảm giá không hợp lệ: %v", err)
		}
	}

	bookingCode := generateBookingCode()

	booking := &domain.Booking{
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
		TotalAmount:    totalAmount,
		CouponCode:     normalizedReq.CouponCode,
		DiscountAmount: discountAmount,
		BookingCode:    bookingCode,
		PaymentStatus:  "unpaid",
		Note:           normalizedReq.Note,
		Status:         "pending", // Changed from "pending_payment" to "pending"
	}

	if err := createBookingRecord(booking); err != nil {
		return nil, shared.ErrCreateBookingFailed
	}

	// Increment coupon usage if applied
	if booking.CouponCode != "" {
		couponData, _ := coupon.FindCouponByCode(booking.CouponCode)
		if couponData != nil {
			coupon.IncrementUsedCount(couponData.ID)
		}
	}

	if err := tourmodule.DecreaseTourRemainingSlots(normalizedReq.TourID, normalizedReq.Quantity); err != nil {
		if err == gorm.ErrInvalidData {
			return nil, shared.ErrInsufficientSlots
		}
		return nil, shared.ErrCreateBookingFailed
	}

	// Gửi email xác nhận đặt tour (bất đồng bộ, không block request)
	go func() {
		_ = shared.SendBookingConfirmationEmail(booking.Email, booking.BookingCode, tour.Name, booking.TravelDate)
		
		// Gửi thông báo trong app
		msg := fmt.Sprintf("Bạn đã đặt thành công tour %s. Mã đặt chỗ: %s.", tour.Name, booking.BookingCode)
		_ = notification.SendNotification(booking.UserID, "Đặt tour thành công", msg, domain.NotifTypeBooking)
	}()

	return booking, nil
}

// GetBookingsByUserID - Lấy lịch sử tour đã đặt của người dùng
func GetBookingsByUserID(userID string) ([]domain.Booking, error) {
	parsedUserID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil || parsedUserID == 0 {
		return nil, shared.ErrInvalidUserID
	}

	bookings, err := FindBookingsByUserID(uint(parsedUserID))
	if err != nil {
		return nil, shared.ErrFetchBookingsFailed
	}

	return bookings, nil
}

// GetBookingByID - Lấy chi tiết 1 booking theo ID
func GetBookingByID(bookingID string) (*domain.Booking, error) {
	parsedBookingID, err := strconv.ParseUint(bookingID, 10, 32)
	if err != nil || parsedBookingID == 0 {
		return nil, shared.ErrInvalidBookingID
	}

	booking, err := FindBookingByID(uint(parsedBookingID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrBookingNotFound
		}
		return nil, shared.ErrFetchBookingsFailed
	}

	return booking, nil
}

// GetBookingByCode - Lấy chi tiết booking theo mã booking code
func GetBookingByCode(bookingCode string) (*domain.Booking, error) {
	if bookingCode == "" {
		return nil, shared.ErrInvalidBookingID
	}

	booking, err := FindBookingByCode(bookingCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrBookingNotFound
		}
		return nil, shared.ErrFetchBookingsFailed
	}

	return booking, nil
}

// ConfirmBooking - Chuyển trạng thái booking từ pending → confirmed
func ConfirmBooking(bookingID string) (*domain.Booking, error) {
	parsedBookingID, err := strconv.ParseUint(bookingID, 10, 32)
	if err != nil || parsedBookingID == 0 {
		return nil, shared.ErrInvalidBookingID
	}

	booking, err := FindBookingByID(uint(parsedBookingID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrBookingNotFound
		}
		return nil, shared.ErrFetchBookingsFailed
	}

	if strings.ToLower(strings.TrimSpace(booking.Status)) != "pending" {
		return nil, shared.ErrInvalidBookingID
	}

	booking.Status = "confirmed"
	if err := database.DB.Save(booking).Error; err != nil {
		return nil, shared.ErrCreateBookingFailed
	}

	return booking, nil
}

// CancelBookingByUserID - Hủy tour đã đặt theo user và booking id.
func CancelBookingByUserID(userID, bookingID string) (*domain.Booking, error) {
	parsedUserID, err := strconv.ParseUint(userID, 10, 32)
	if err != nil || parsedUserID == 0 {
		return nil, shared.ErrInvalidUserID
	}

	parsedBookingID, err := strconv.ParseUint(bookingID, 10, 32)
	if err != nil || parsedBookingID == 0 {
		return nil, shared.ErrInvalidBookingID
	}

	booking, err := cancelBookingByUserID(uint(parsedUserID), uint(parsedBookingID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrBookingNotFound
		}

		if errors.Is(err, ErrBookingAlreadyCancelled) {
			return nil, shared.ErrBookingCannotCancel
		}

		return nil, shared.ErrCancelBookingFailed
	}

	return booking, nil
}

func normalizeBookingRequest(req domain.CreateBookingRequest) domain.CreateBookingRequest {
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

func validateBookingRequest(req domain.CreateBookingRequest) error {
	if req.TourID == 0 {
		return shared.ErrInvalidBookingPayload
	}

	if req.FullName == "" {
		return shared.ErrInvalidFullName
	}

	if !phoneRegexp.MatchString(req.Phone) {
		return shared.ErrInvalidPhone
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return shared.ErrInvalidEmail
	}

	if req.AdultCount < 1 {
		return shared.ErrInvalidAdultCount
	}

	if req.ChildCount < 0 {
		return shared.ErrInvalidChildCount
	}

	if req.InfantCount < 0 {
		return shared.ErrInvalidInfantCount
	}

	if req.Quantity <= 0 {
		return shared.ErrInvalidQuantity
	}

	travelDate, err := time.Parse("2006-01-02", req.TravelDate)
	if err != nil {
		return shared.ErrInvalidTravelDate
	}

	startOfToday := time.Now().Truncate(24 * time.Hour)
	if travelDate.Before(startOfToday) {
		return shared.ErrTravelDateInPast
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
