package domain

import (
	"fmt"
	"time"
)

// DiscountType constants
const (
	DiscountTypePercent = "percent"
	DiscountTypeFixed   = "fixed"
)

// Coupon - Model đại diện cho mã giảm giá
type Coupon struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Code          string    `json:"code" gorm:"uniqueIndex;size:50;not null"`
	DiscountType  string    `json:"discount_type" gorm:"not null;default:'percent'"` // percent or fixed
	DiscountValue int64     `json:"discount_value" gorm:"not null"`
	MinOrderValue int64     `json:"min_order_value" gorm:"not null;default:0"`
	MaxDiscount   int64     `json:"max_discount" gorm:"not null;default:0"` // 0 means no limit
	StartDate     time.Time `json:"start_date" gorm:"not null"`
	EndDate       time.Time `json:"end_date" gorm:"not null"`
	UsageLimit    int       `json:"usage_limit" gorm:"not null;default:0"` // 0 means unlimited
	UsedCount     int       `json:"used_count" gorm:"not null;default:0"`
	IsActive      bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// IsValid - Kiểm tra xem coupon có đang hợp lệ không (thời gian, số lượt, trạng thái)
func (c *Coupon) IsValid() error {
	if !c.IsActive {
		return fmt.Errorf("mã giảm giá này đã bị vô hiệu hóa")
	}

	now := time.Now()
	if now.Before(c.StartDate) {
		return fmt.Errorf("mã giảm giá chưa đến thời gian áp dụng")
	}
	if now.After(c.EndDate) {
		return fmt.Errorf("mã giảm giá đã hết hạn")
	}

	if c.UsageLimit > 0 && c.UsedCount >= c.UsageLimit {
		return fmt.Errorf("mã giảm giá đã hết lượt sử dụng")
	}

	return nil
}

// CalculateDiscount - Tính toán số tiền được giảm
func (c *Coupon) CalculateDiscount(orderTotal int64) (int64, error) {
	if err := c.IsValid(); err != nil {
		return 0, err
	}

	if orderTotal < c.MinOrderValue {
		return 0, fmt.Errorf("đơn hàng phải đạt tối thiểu %s để áp dụng mã này", FormatPriceVND(c.MinOrderValue))
	}

	var discountAmount int64
	if c.DiscountType == DiscountTypeFixed {
		discountAmount = c.DiscountValue
	} else if c.DiscountType == DiscountTypePercent {
		discountAmount = (orderTotal * c.DiscountValue) / 100
		if c.MaxDiscount > 0 && discountAmount > c.MaxDiscount {
			discountAmount = c.MaxDiscount
		}
	}

	// Không giảm quá giá trị đơn hàng
	if discountAmount > orderTotal {
		discountAmount = orderTotal
	}

	return discountAmount, nil
}

// CreateCouponRequest - Request để tạo mới mã giảm giá
type CreateCouponRequest struct {
	Code          string    `json:"code" binding:"required,min=3,max=50"`
	DiscountType  string    `json:"discount_type" binding:"required,oneof=percent fixed"`
	DiscountValue int64     `json:"discount_value" binding:"required,min=1"`
	MinOrderValue int64     `json:"min_order_value" binding:"min=0"`
	MaxDiscount   int64     `json:"max_discount" binding:"min=0"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date" binding:"required"`
	UsageLimit    int       `json:"usage_limit" binding:"min=0"`
	IsActive      *bool     `json:"is_active" binding:"required"`
}

// ValidateCouponRequest - Request từ khách hàng để kiểm tra mã
type ValidateCouponRequest struct {
	Code       string `json:"code" binding:"required"`
	OrderTotal int64  `json:"order_total" binding:"required,min=0"`
}
