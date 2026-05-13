package coupon

import (
	"travel-backend/database"
	"travel-backend/domain"

	"gorm.io/gorm"
)

// CreateCoupon - Tạo mới mã giảm giá
func CreateCoupon(coupon *domain.Coupon) error {
	return database.DB.Create(coupon).Error
}

// FindCouponByCode - Lấy mã giảm giá theo mã nhập (Code)
func FindCouponByCode(code string) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := database.DB.Where("code = ?", code).First(&coupon).Error
	return &coupon, err
}

// FindCouponByID - Lấy mã giảm giá theo ID
func FindCouponByID(id uint) (*domain.Coupon, error) {
	var coupon domain.Coupon
	err := database.DB.First(&coupon, id).Error
	return &coupon, err
}

// UpdateCoupon - Cập nhật mã giảm giá
func UpdateCoupon(coupon *domain.Coupon) error {
	return database.DB.Save(coupon).Error
}

// DeleteCoupon - Xóa mã giảm giá (chỉ xóa khi chưa ai dùng)
func DeleteCoupon(id uint) error {
	return database.DB.Delete(&domain.Coupon{}, id).Error
}

// IncrementUsedCount - Tăng số lần sử dụng của mã giảm giá
func IncrementUsedCount(id uint) error {
	return database.DB.Model(&domain.Coupon{}).Where("id = ?", id).UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

// FindAllCoupons - Lấy danh sách mã giảm giá cho Admin (có phân trang)
func FindAllCoupons(search string, offset, limit int) ([]domain.Coupon, int64, error) {
	query := database.DB.Model(&domain.Coupon{})

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("code LIKE ?", like)
	}

	var total int64
	query.Count(&total)

	var coupons []domain.Coupon
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&coupons).Error

	return coupons, total, err
}
