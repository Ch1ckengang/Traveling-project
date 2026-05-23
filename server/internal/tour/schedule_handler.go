package tour

import (
	"net/http"
	"strconv"
	"time"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// AdminGetTourSchedules - GET /v1/api/admin/tours/:id/schedules
func AdminGetTourSchedules(c *gin.Context) {
	tourID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || tourID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID tour không hợp lệ", "INVALID_ID")
		return
	}

	var schedules []domain.TourSchedule
	if err := database.DB.Where("tour_id = ?", tourID).Order("departure_date asc").Find(&schedules).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi lấy lịch trình", "FETCH_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "", schedules)
}

// AdminCreateTourSchedule - POST /v1/api/admin/tours/:id/schedules
func AdminCreateTourSchedule(c *gin.Context) {
	tourID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || tourID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID tour không hợp lệ", "INVALID_ID")
		return
	}

	var req struct {
		DepartureDate time.Time `json:"departure_date" binding:"required"`
		ReturnDate    time.Time `json:"return_date"`
		TotalSlots    int       `json:"total_slots" binding:"required,min=1"`
		PriceModifier int64     `json:"price_modifier"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "INVALID_PAYLOAD")
		return
	}

	schedule := domain.TourSchedule{
		TourID:         uint(tourID),
		DepartureDate:  req.DepartureDate,
		ReturnDate:     req.ReturnDate,
		TotalSlots:     req.TotalSlots,
		RemainingSlots: req.TotalSlots,
		PriceModifier:  req.PriceModifier,
		Status:         "active",
	}

	if err := database.DB.Create(&schedule).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi tạo lịch trình", "CREATE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusCreated, "Tạo lịch trình thành công", schedule)
}

// AdminDeleteTourSchedule - DELETE /v1/api/admin/tours/:id/schedules/:scheduleId
func AdminDeleteTourSchedule(c *gin.Context) {
	scheduleID, err := strconv.ParseUint(c.Param("scheduleId"), 10, 32)
	if err != nil || scheduleID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID lịch trình không hợp lệ", "INVALID_ID")
		return
	}

	if err := database.DB.Where("id = ?", scheduleID).Delete(&domain.TourSchedule{}).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi xóa lịch trình", "DELETE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Xóa lịch trình thành công", nil)
}

// GetTourSchedules - GET /v1/api/tours/:id/schedules (Public)
func GetTourSchedules(c *gin.Context) {
	tourID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || tourID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID tour không hợp lệ", "INVALID_ID")
		return
	}

	var schedules []domain.TourSchedule
	// Only fetch active schedules with remaining slots from today onwards
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	
	if err := database.DB.
		Where("tour_id = ? AND status = ? AND departure_date >= ? AND remaining_slots > 0", tourID, "active", today).
		Order("departure_date asc").
		Find(&schedules).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi lấy lịch trình", "FETCH_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "", schedules)
}
