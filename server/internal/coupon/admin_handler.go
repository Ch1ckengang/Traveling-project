package coupon

import (
	"net/http"
	"strconv"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// AdminGetCouponsHandler - GET /v1/api/admin/coupons (Staff+)
func AdminGetCouponsHandler(c *gin.Context) {
	pagination := shared.GetPaginationParams(c)
	search := c.DefaultQuery("search", "")

	coupons, total, err := FindAllCoupons(search, pagination.Offset(), pagination.Limit)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy danh sách mã giảm giá", "ADMIN_COUPON_FETCH_FAILED")
		return
	}

	meta := shared.BuildPaginationMeta(pagination, int(total))
	shared.RespondSuccessWithMeta(c, http.StatusOK, "", coupons, meta)
}

// AdminCreateCouponHandler - POST /v1/api/admin/coupons (Staff+)
func AdminCreateCouponHandler(c *gin.Context) {
	var req domain.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "ADMIN_COUPON_INVALID_PAYLOAD")
		return
	}

	coupon, err := AdminCreateCoupon(req)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "ADMIN_COUPON_CREATE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusCreated, "Tạo mã giảm giá thành công", coupon)
}

// AdminUpdateCouponHandler - PUT /v1/api/admin/coupons/:id (Staff+)
func AdminUpdateCouponHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID mã giảm giá không hợp lệ", "ADMIN_COUPON_INVALID_ID")
		return
	}

	var req domain.CreateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "ADMIN_COUPON_INVALID_PAYLOAD")
		return
	}

	coupon, err := AdminUpdateCoupon(uint(id), req)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "ADMIN_COUPON_UPDATE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Cập nhật mã giảm giá thành công", coupon)
}

// AdminDeleteCouponHandler - DELETE /v1/api/admin/coupons/:id (Staff+)
func AdminDeleteCouponHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID mã giảm giá không hợp lệ", "ADMIN_COUPON_INVALID_ID")
		return
	}

	if err := AdminDeleteCoupon(uint(id)); err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "ADMIN_COUPON_DELETE_FAILED")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Xóa mã giảm giá thành công", nil)
}
