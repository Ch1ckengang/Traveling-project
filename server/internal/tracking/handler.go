package tracking

import (
	"net/http"
	"time"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/auth"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// LogActivityRequest - DTO nhận log từ client
type LogActivityRequest struct {
	SessionID     string `json:"session_id" binding:"required"`
	ActionType    string `json:"action_type" binding:"required"` // "view_tour", "search"
	TourID        *uint  `json:"tour_id"`
	SearchKeyword string `json:"search_keyword"`
	Category      string `json:"category"`
}

// LogUserActivity - POST /v1/api/tracking/log
func LogUserActivity(c *gin.Context) {
	var req LogActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "INVALID_PAYLOAD")
		return
	}

	activity := domain.UserActivity{
		SessionID:     req.SessionID,
		ActionType:    req.ActionType,
		TourID:        req.TourID,
		SearchKeyword: req.SearchKeyword,
		Category:      req.Category,
		Timestamp:     time.Now(),
	}

	// Gắn UserID nếu đang đăng nhập
	if userID, ok := auth.GetAuthenticatedUserID(c); ok {
		uid := userID
		activity.UserID = &uid
	}

	// Lưu bất đồng bộ để không chặn luồng API chính
	go func(act domain.UserActivity) {
		database.DB.Create(&act)
	}(activity)

	shared.RespondSuccess(c, http.StatusOK, "Log received", nil)
}
