package review

import (
	"net/http"
	"strconv"
	"travel-backend/domain"
	authmodule "travel-backend/internal/auth"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// ===== ADMIN REVIEW HANDLERS =====

// AdminGetReviewsHandler - GET /v1/api/admin/reviews (Staff+)
// Lấy tất cả reviews với phân trang, filter status, search
func AdminGetReviewsHandler(c *gin.Context) {
	pagination := shared.GetPaginationParams(c)
	status := c.DefaultQuery("status", "")
	search := c.DefaultQuery("search", "")

	reviews, total, err := FindAllReviews(status, search, pagination.Offset(), pagination.Limit)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy reviews", "ADMIN_REVIEW_FETCH_FAILED")
		return
	}

	// Convert to admin view (include full info)
	results := make([]gin.H, len(reviews))
	for i, r := range reviews {
		item := gin.H{
			"id":          r.ID,
			"tour_id":     r.TourID,
			"user_id":     r.UserID,
			"booking_id":  r.BookingID,
			"rating":      r.Rating,
			"title":       r.Title,
			"content":     r.Content,
			"status":      r.Status,
			"admin_reply": r.AdminReply,
			"replied_at":  r.RepliedAt,
			"created_at":  r.CreatedAt,
			"updated_at":  r.UpdatedAt,
		}
		if r.Tour != nil {
			item["tour_name"] = r.Tour.Name
		}
		if r.User != nil {
			item["user_name"] = r.User.Name
			item["user_email"] = r.User.Email
		}
		results[i] = item
	}

	meta := shared.BuildPaginationMeta(pagination, int(total))
	shared.RespondSuccessWithMeta(c, http.StatusOK, "", results, meta)
}

// AdminPublishReviewHandler - PUT /v1/api/admin/reviews/:id/publish (Staff+)
func AdminPublishReviewHandler(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || reviewID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID review không hợp lệ", "ADMIN_REVIEW_INVALID_ID")
		return
	}

	review, err := AdminPublishReview(uint(reviewID))
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "ADMIN_REVIEW_PUBLISH_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã publish review", gin.H{
		"review_id": review.ID,
		"status":    review.Status,
	})
}

// AdminHideReviewHandler - PUT /v1/api/admin/reviews/:id/hide (Staff+)
func AdminHideReviewHandler(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || reviewID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID review không hợp lệ", "ADMIN_REVIEW_INVALID_ID")
		return
	}

	review, err := AdminHideReview(uint(reviewID))
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "ADMIN_REVIEW_HIDE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã ẩn review", gin.H{
		"review_id": review.ID,
		"status":    review.Status,
	})
}

// AdminReplyReviewHandler - POST /v1/api/admin/reviews/:id/reply (Staff+)
func AdminReplyReviewHandler(c *gin.Context) {
	adminUserID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", "AUTH_TOKEN_INVALID")
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || reviewID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID review không hợp lệ", "ADMIN_REVIEW_INVALID_ID")
		return
	}

	var req domain.AdminReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "ADMIN_REVIEW_INVALID_PAYLOAD")
		return
	}

	review, err := AdminReplyReview(uint(reviewID), adminUserID, req.Reply)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "ADMIN_REVIEW_REPLY_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã phản hồi review", gin.H{
		"review_id":   review.ID,
		"admin_reply": review.AdminReply,
		"replied_at":  review.RepliedAt,
	})
}
