package coupon

import (
	"fmt"
	"strings"
	"travel-backend/domain"

	"gorm.io/gorm"
)

// ValidateCoupon - Khách hàng kiểm tra mã giảm giá
func ValidateCoupon(code string, orderTotal int64) (*domain.Coupon, int64, error) {
	code = strings.ToUpper(strings.TrimSpace(code))

	coupon, err := FindCouponByCode(code)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, fmt.Errorf("mã giảm giá không tồn tại")
		}
		return nil, 0, fmt.Errorf("lỗi kiểm tra mã giảm giá: %v", err)
	}

	discountAmount, err := coupon.CalculateDiscount(orderTotal)
	if err != nil {
		return nil, 0, err
	}

	return coupon, discountAmount, nil
}

// AdminCreateCoupon - Admin tạo mã giảm giá mới
func AdminCreateCoupon(req domain.CreateCouponRequest) (*domain.Coupon, error) {
	code := strings.ToUpper(strings.TrimSpace(req.Code))

	// Check exist
	existing, err := FindCouponByCode(code)
	if err == nil && existing.ID != 0 {
		return nil, fmt.Errorf("mã giảm giá này đã tồn tại")
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	if req.StartDate.After(req.EndDate) {
		return nil, fmt.Errorf("ngày bắt đầu không thể sau ngày kết thúc")
	}

	coupon := &domain.Coupon{
		Code:          code,
		DiscountType:  req.DiscountType,
		DiscountValue: req.DiscountValue,
		MinOrderValue: req.MinOrderValue,
		MaxDiscount:   req.MaxDiscount,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		UsageLimit:    req.UsageLimit,
		IsActive:      isActive,
	}

	if err := CreateCoupon(coupon); err != nil {
		return nil, fmt.Errorf("không thể tạo mã giảm giá: %v", err)
	}

	return coupon, nil
}

// AdminUpdateCoupon - Admin cập nhật mã giảm giá
func AdminUpdateCoupon(id uint, req domain.CreateCouponRequest) (*domain.Coupon, error) {
	coupon, err := FindCouponByID(id)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy mã giảm giá")
	}

	code := strings.ToUpper(strings.TrimSpace(req.Code))
	if code != coupon.Code {
		existing, err := FindCouponByCode(code)
		if err == nil && existing.ID != 0 {
			return nil, fmt.Errorf("mã giảm giá '%s' đã tồn tại", code)
		}
	}

	if req.StartDate.After(req.EndDate) {
		return nil, fmt.Errorf("ngày bắt đầu không thể sau ngày kết thúc")
	}

	coupon.Code = code
	coupon.DiscountType = req.DiscountType
	coupon.DiscountValue = req.DiscountValue
	coupon.MinOrderValue = req.MinOrderValue
	coupon.MaxDiscount = req.MaxDiscount
	coupon.StartDate = req.StartDate
	coupon.EndDate = req.EndDate
	coupon.UsageLimit = req.UsageLimit
	
	if req.IsActive != nil {
		coupon.IsActive = *req.IsActive
	}

	if err := UpdateCoupon(coupon); err != nil {
		return nil, fmt.Errorf("không thể cập nhật mã giảm giá: %v", err)
	}

	return coupon, nil
}

// AdminDeleteCoupon - Admin xóa mã giảm giá
func AdminDeleteCoupon(id uint) error {
	coupon, err := FindCouponByID(id)
	if err != nil {
		return fmt.Errorf("không tìm thấy mã giảm giá")
	}

	if coupon.UsedCount > 0 {
		return fmt.Errorf("không thể xóa mã giảm giá đã có người sử dụng. Vui lòng vô hiệu hóa thay vì xóa.")
	}

	return DeleteCoupon(id)
}
