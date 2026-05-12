package auth

import (
	"travel-backend/database"
	"travel-backend/domain"
)

// SaveOTP - Lưu OTP vào database
func SaveOTP(otp *domain.OTP) error {
	return database.DB.Create(otp).Error
}

// FindOTPByEmail - Tìm OTP mới nhất cho email
func FindOTPByEmail(email string) (*domain.OTP, error) {
	var otp domain.OTP
	err := database.DB.
		Where("email = ?", email).
		Order("created_at DESC").
		First(&otp).Error
	return &otp, err
}

// DeleteOTPByEmail - Xóa tất cả OTP của email (sau khi verify thành công)
func DeleteOTPByEmail(email string) error {
	return database.DB.Where("email = ?", email).Delete(&domain.OTP{}).Error
}

// CleanupExpiredOTPs - Xóa các OTP đã hết hạn (có thể chạy định kỳ)
func CleanupExpiredOTPs() error {
	return database.DB.Where("expires_at < NOW()").Delete(&domain.OTP{}).Error
}
