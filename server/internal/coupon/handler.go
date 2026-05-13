package coupon

import (
	"net/http"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// ValidateCouponHandler - POST /v1/api/coupons/validate (Public/Customer)
func ValidateCouponHandler(c *gin.Context) {
	var req domain.ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "COUPON_INVALID_PAYLOAD")
		return
	}

	coupon, discountAmount, err := ValidateCoupon(req.Code, req.OrderTotal)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, err.Error(), "COUPON_INVALID")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Áp dụng mã giảm giá thành công", gin.H{
		"coupon_code":     coupon.Code,
		"discount_type":   coupon.DiscountType,
		"discount_value":  coupon.DiscountValue,
		"discount_amount": discountAmount,
		"final_total":     req.OrderTotal - discountAmount,
	})
}
