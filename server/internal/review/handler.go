package review

import (
	"net/http"
	"strconv"
	"travel-backend/domain"
	authmodule "travel-backend/internal/auth"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// GetTourReviewsHandler - GET /v1/api/tours/:id/reviews (public)
// Lấy danh sách reviews của tour với phân trang và filter rating
func GetTourReviewsHandler(c *gin.Context) {
	tourID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || tourID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID tour không hợp lệ", "REVIEW_INVALID_TOUR_ID")
		return
	}

	pagination := shared.GetPaginationParams(c)
	ratingFilter, _ := strconv.Atoi(c.DefaultQuery("rating", "0"))

	reviews, total, err := GetTourReviews(uint(tourID), ratingFilter, pagination.Offset(), pagination.Limit)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy reviews", "REVIEW_FETCH_FAILED")
		return
	}

	// Lấy current user ID (nếu đã login) để check canEdit
	var currentUserID uint
	if uid, ok := authmodule.GetAuthenticatedUserID(c); ok {
		currentUserID = uid
	}

	summaries := make([]*domain.ReviewSummary, len(reviews))
	for i, r := range reviews {
		summaries[i] = r.ToSummary(currentUserID)
	}

	// Lấy stats
	stats, _ := GetTourStats(uint(tourID))

	meta := shared.BuildPaginationMeta(pagination, int(total))
	meta["stats"] = stats

	shared.RespondSuccessWithMeta(c, http.StatusOK, "", summaries, meta)
}

// CreateReviewHandler - POST /v1/api/reviews (auth required)
// Customer tạo review cho tour đã book
func CreateReviewHandler(c *gin.Context) {
	userID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", "AUTH_TOKEN_INVALID")
		return
	}

	var req domain.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "REVIEW_INVALID_PAYLOAD")
		return
	}

	review, err := SubmitReview(userID, req)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "REVIEW_SUBMIT_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusCreated, "Đánh giá thành công!", gin.H{
		"review": review.ToSummary(userID),
	})
}

// UpdateReviewHandler - PUT /v1/api/reviews/:id (auth required, owner only, 7 days)
// Customer sửa review trong vòng 7 ngày
func UpdateReviewHandler(c *gin.Context) {
	userID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Yêu cầu đăng nhập", "AUTH_TOKEN_INVALID")
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || reviewID == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID review không hợp lệ", "REVIEW_INVALID_ID")
		return
	}

	var req domain.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "REVIEW_INVALID_PAYLOAD")
		return
	}

	review, err := EditReview(uint(reviewID), userID, req)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "REVIEW_UPDATE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Cập nhật đánh giá thành công", gin.H{
		"review": review.ToSummary(userID),
	})
}
