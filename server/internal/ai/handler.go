package ai

import (
	"net/http"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/auth"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	Message string `json:"message" binding:"required"`
}

// ChatHandler - POST /v1/api/ai/chat
func ChatHandler(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "INVALID_PAYLOAD")
		return
	}

	reply, err := ChatWithBot(c.Request.Context(), req.Message)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi khi gọi AI: "+err.Error(), "AI_ERROR")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Thành công", gin.H{
		"reply": reply,
	})
}

// RecommendHandler - GET /v1/api/ai/recommendations
func RecommendHandler(c *gin.Context) {
	userID, ok := auth.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", "AUTH_REQUIRED")
		return
	}

	tourIDs, err := AnalyzeBehavior(c.Request.Context(), userID)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể phân tích: "+err.Error(), "AI_ERROR")
		return
	}

	if len(tourIDs) == 0 {
		shared.RespondSuccess(c, http.StatusOK, "Không đủ dữ liệu", []domain.Tour{})
		return
	}

	var tours []domain.Tour
	if err := database.DB.Where("id IN ?", tourIDs).Find(&tours).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi lấy dữ liệu tour", "FETCH_ERROR")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Thành công", tours)
}
