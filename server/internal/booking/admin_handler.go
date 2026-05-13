package booking

import (
	"net/http"
	"strconv"
	"strings"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// ===== ADMIN BOOKING HANDLERS =====
// Yêu cầu StaffRequired() middleware đã chạy trước

// AdminGetBookingsHandler - GET /v1/api/admin/bookings
// Lấy tất cả bookings với phân trang, filter, search
func AdminGetBookingsHandler(c *gin.Context) {
	pagination := shared.GetPaginationParams(c)
	status := c.DefaultQuery("status", "")
	search := c.DefaultQuery("search", "")
	paymentStatus := c.DefaultQuery("payment_status", "")

	query := database.DB.Model(&domain.Booking{}).Preload("Tour")

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if paymentStatus != "" {
		query = query.Where("payment_status = ?", paymentStatus)
	}
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(booking_code) LIKE ? OR LOWER(full_name) LIKE ? OR LOWER(phone) LIKE ? OR LOWER(email) LIKE ?",
			like, like, like, like,
		)
	}

	var total int64
	query.Count(&total)

	var bookings []domain.Booking
	err := query.
		Order("created_at DESC").
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		Find(&bookings).Error

	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy danh sách booking", "ADMIN_BOOKING_FETCH_FAILED")
		return
	}

	// Convert to BookingWithPayment
	results := make([]*domain.BookingWithPayment, len(bookings))
	for i, b := range bookings {
		results[i] = b.ToBookingWithPayment()
	}

	meta := shared.BuildPaginationMeta(pagination, int(total))
	shared.RespondSuccessWithMeta(c, http.StatusOK, "", results, meta)
}

// AdminGetBookingByCodeHandler - GET /v1/api/admin/bookings/:code
// Lấy chi tiết booking theo mã code (admin không cần check ownership)
func AdminGetBookingByCodeHandler(c *gin.Context) {
	code := c.Param("code")

	var booking domain.Booking
	err := database.DB.Preload("Tour").Where("booking_code = ?", code).First(&booking).Error
	if err != nil {
		// Thử tìm theo ID nếu không phải booking code
		id, parseErr := strconv.ParseUint(code, 10, 32)
		if parseErr == nil {
			err = database.DB.Preload("Tour").First(&booking, id).Error
		}
		if err != nil {
			shared.RespondError(c, http.StatusNotFound, "Không tìm thấy booking", "ADMIN_BOOKING_NOT_FOUND")
			return
		}
	}

	shared.RespondSuccess(c, http.StatusOK, "", gin.H{
		"booking": booking.ToBookingWithPayment(),
	})
}

// AdminConfirmBookingHandler - PUT /v1/api/admin/bookings/:code/confirm
// Xác nhận booking (chuyển status → confirmed)
func AdminConfirmBookingHandler(c *gin.Context) {
	code := c.Param("code")

	var booking domain.Booking
	if err := database.DB.Where("booking_code = ?", code).First(&booking).Error; err != nil {
		shared.RespondError(c, http.StatusNotFound, "Không tìm thấy booking", "ADMIN_BOOKING_NOT_FOUND")
		return
	}

	status := strings.ToLower(strings.TrimSpace(booking.Status))
	if status != "booked" && status != "pending" {
		shared.RespondError(c, http.StatusConflict,
			"Chỉ có thể xác nhận booking ở trạng thái 'booked' hoặc 'pending'",
			"ADMIN_BOOKING_INVALID_STATUS")
		return
	}

	updates := map[string]interface{}{
		"status": "confirmed",
	}

	if err := database.DB.Model(&booking).Updates(updates).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể xác nhận booking", "ADMIN_BOOKING_CONFIRM_FAILED")
		return
	}

	database.DB.Preload("Tour").First(&booking, booking.ID)

	shared.RespondSuccess(c, http.StatusOK, "Xác nhận booking thành công", gin.H{
		"booking": booking.ToBookingWithPayment(),
	})
}

// AdminCancelBookingHandler - PUT /v1/api/admin/bookings/:code/cancel
// Admin hủy booking (hoàn lại slots cho tour)
func AdminCancelBookingHandler(c *gin.Context) {
	code := c.Param("code")

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	var booking domain.Booking
	if err := database.DB.Where("booking_code = ?", code).First(&booking).Error; err != nil {
		shared.RespondError(c, http.StatusNotFound, "Không tìm thấy booking", "ADMIN_BOOKING_NOT_FOUND")
		return
	}

	status := strings.ToLower(strings.TrimSpace(booking.Status))
	if status == "cancelled" || status == "canceled" {
		shared.RespondError(c, http.StatusConflict, "Booking đã bị hủy trước đó", "ADMIN_BOOKING_ALREADY_CANCELLED")
		return
	}

	// Cập nhật status
	updates := map[string]interface{}{
		"status":         "cancelled",
		"payment_status": "cancelled",
	}
	if req.Reason != "" {
		updates["note"] = booking.Note + " [Admin hủy: " + req.Reason + "]"
	}

	if err := database.DB.Model(&booking).Updates(updates).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể hủy booking", "ADMIN_BOOKING_CANCEL_FAILED")
		return
	}

	// Hoàn lại slots cho tour
	database.DB.Model(&domain.Tour{}).
		Where("id = ?", booking.TourID).
		UpdateColumn("remaining_slots", database.DB.Raw("remaining_slots + ?", booking.Quantity))

	database.DB.Preload("Tour").First(&booking, booking.ID)

	shared.RespondSuccess(c, http.StatusOK, "Hủy booking thành công", gin.H{
		"booking": booking.ToBookingWithPayment(),
	})
}

// AdminGetBookingStatsHandler - GET /v1/api/admin/bookings/stats
// Thống kê nhanh số lượng booking theo status
func AdminGetBookingStatsHandler(c *gin.Context) {
	type StatusCount struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}

	var stats []StatusCount
	database.DB.Model(&domain.Booking{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&stats)

	var totalBookings int64
	database.DB.Model(&domain.Booking{}).Count(&totalBookings)

	var totalRevenue int64
	database.DB.Model(&domain.Booking{}).
		Where("payment_status = ?", "paid").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue)

	shared.RespondSuccess(c, http.StatusOK, "", gin.H{
		"stats":          stats,
		"total_bookings": totalBookings,
		"total_revenue":  totalRevenue,
	})
}
