package booking

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
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
	"gorm.io/gorm/clause"
)

var phoneRegexp = regexp.MustCompile(`^[0-9+\-\s]{8,20}$`)

// CreateBooking - Xử lý nghiệp vụ đặt tour
// Sử dụng database transaction với row-level lock để prevent race condition (overbooking).
// Tất cả operations (check slots, create booking, decrement slots, increment coupon usage)
// đều nằm trong cùng một transaction.
func CreateBooking(req domain.CreateBookingRequest) (*domain.Booking, error) {
	normalizedReq := normalizeBookingRequest(req)

	if err := validateBookingRequest(normalizedReq); err != nil {
		return nil, err
	}

	// Validate tour exists (read-only, outside transaction)
	tourInfo, err := tourmodule.FindTourByID(normalizedReq.TourID)
	if err != nil {
		return nil, shared.ErrTourNotFound
	}

	basePrice := tourInfo.PriceAmount
	if basePrice <= 0 && tourInfo.Price != "" {
		// Fallback: parse Price string (e.g. "3.500.000đ" → 3500000)
		basePrice = domain.ParsePriceVND(tourInfo.Price)
	}

	totalAmount := calculateBookingTotal(basePrice, normalizedReq.AdultCount, normalizedReq.ChildCount)

	// If ScheduleID is provided, validate schedule belongs to tour
	if normalizedReq.ScheduleID != nil {
		schedule, err := tourmodule.FindTourScheduleByID(*normalizedReq.ScheduleID)
		if err != nil || schedule.TourID != normalizedReq.TourID {
			return nil, fmt.Errorf("lịch trình không hợp lệ")
		}

		// Update TravelDate and price modifier
		normalizedReq.TravelDate = schedule.DepartureDate.Format("02/01/2006")
		totalAmount += schedule.PriceModifier * int64(normalizedReq.Quantity)
	}

	var discountAmount int64 = 0
	var couponData *domain.Coupon

	// Validate coupon BEFORE transaction (read-only check)
	if normalizedReq.CouponCode != "" {
		var discount int64
		couponData, discount, err = coupon.ValidateCoupon(normalizedReq.CouponCode, totalAmount)
		if err != nil || couponData == nil {
			return nil, fmt.Errorf("mã giảm giá không hợp lệ: %v", err)
		}
		discountAmount = discount
		totalAmount = totalAmount - discountAmount
	}

	bookingCode := generateBookingCode()

	var booking *domain.Booking

	// === BEGIN TRANSACTION ===
	// Tất cả operations phải thành công hoặc rollback toàn bộ
	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		if normalizedReq.ScheduleID != nil {
			// Lock schedule row và check slots
			var schedule domain.TourSchedule
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&schedule, *normalizedReq.ScheduleID).Error; err != nil {
				return fmt.Errorf("lịch trình không tồn tại: %v", err)
			}
			if normalizedReq.Quantity > schedule.RemainingSlots {
				return shared.ErrInsufficientSlots
			}

			// Decrease schedule slots
			if err := tx.Model(&schedule).
				Update("remaining_slots", gorm.Expr("remaining_slots - ?", normalizedReq.Quantity)).Error; err != nil {
				return fmt.Errorf("không thể cập nhật chỗ trống: %v", err)
			}
		} else {
			// Lock tour row và check slots
			var tour domain.Tour
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				First(&tour, normalizedReq.TourID).Error; err != nil {
				return fmt.Errorf("tour không tồn tại: %v", err)
			}

			remaining := tour.RemainingSlots
			if remaining <= 0 {
				remaining = 30 // Normalize like FindTourByID does
			}
			if normalizedReq.Quantity > remaining {
				return shared.ErrInsufficientSlots
			}

			// Decrease tour slots
			if err := tx.Model(&tour).
				Update("remaining_slots", gorm.Expr("remaining_slots - ?", normalizedReq.Quantity)).Error; err != nil {
				return fmt.Errorf("không thể cập nhật chỗ trống: %v", err)
			}
		}

		// Create booking record
		booking = &domain.Booking{
			UserID:         normalizedReq.UserID,
			TourID:         normalizedReq.TourID,
			FullName:       normalizedReq.FullName,
			Phone:          normalizedReq.Phone,
			Email:          normalizedReq.Email,
			AdultCount:     normalizedReq.AdultCount,
			ChildCount:     normalizedReq.ChildCount,
			InfantCount:    normalizedReq.InfantCount,
			Quantity:       normalizedReq.Quantity,
			TravelDate:     normalizedReq.TravelDate,
			ScheduleID:     normalizedReq.ScheduleID,
			TotalAmount:    totalAmount,
			CouponCode:     normalizedReq.CouponCode,
			DiscountAmount: discountAmount,
			BookingCode:    bookingCode,
			PaymentStatus:  "unpaid",
			Note:           normalizedReq.Note,
			Status:         "pending",
		}

		if err := tx.Create(booking).Error; err != nil {
			return fmt.Errorf("không thể tạo booking: %v", err)
		}

		// Increment coupon usage INSIDE transaction (atomic with booking)
		if couponData != nil {
			if err := tx.Model(&domain.Coupon{}).
				Where("id = ?", couponData.ID).
				UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
				return fmt.Errorf("không thể cập nhật coupon usage: %v", err)
			}
		}

		return nil
	})
	// === END TRANSACTION ===

	if txErr != nil {
		if errors.Is(txErr, shared.ErrInsufficientSlots) {
			return nil, shared.ErrInsufficientSlots
		}
		return nil, shared.ErrCreateBookingFailed
	}

	// Gửi email xác nhận đặt tour (bất đồng bộ, không block request)
	go func() {
		if err := shared.SendBookingConfirmationEmail(booking.Email, booking.BookingCode, tourInfo.Name, booking.TravelDate); err != nil {
			log.Printf("[BOOKING][EMAIL] Failed to send confirmation email for %s: %v", booking.BookingCode, err)
		}

		// Gửi thông báo trong app
		msg := fmt.Sprintf("Bạn đã đặt thành công tour %s. Mã đặt chỗ: %s.", tourInfo.Name, booking.BookingCode)
		if err := notification.SendNotification(booking.UserID, "Đặt tour thành công", msg, domain.NotifTypeBooking); err != nil {
			log.Printf("[BOOKING][NOTIFICATION] Failed to send notification for %s: %v", booking.BookingCode, err)
		}
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

// generateBookingCode - Tạo mã booking unique sử dụng crypto/rand
// Format: TOUR-YYYYMMDD-XXXXXX (X = hex random)
func generateBookingCode() string {
	b := make([]byte, 4) // 4 bytes = 8 hex chars
	if _, err := rand.Read(b); err != nil {
		// Fallback nếu crypto/rand fail
		return fmt.Sprintf("TOUR-%s-%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("TOUR-%s-%s", time.Now().Format("20060102"), hex.EncodeToString(b))
}
