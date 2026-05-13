package notification

import (
	"travel-backend/database"
	"travel-backend/domain"
)

// CreateNotification - Lưu thông báo mới
func CreateNotification(notif *domain.Notification) error {
	return database.DB.Create(notif).Error
}

// FindUserNotifications - Lấy danh sách thông báo của user
func FindUserNotifications(userID uint, offset, limit int) ([]domain.Notification, int64, error) {
	var notifs []domain.Notification
	var total int64

	query := database.DB.Model(&domain.Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notifs).Error

	return notifs, total, err
}

// FindUnreadCount - Lấy số lượng thông báo chưa đọc
func FindUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&domain.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

// MarkAsRead - Đánh dấu 1 thông báo là đã đọc
func MarkAsRead(notifID uint, userID uint) error {
	return database.DB.Model(&domain.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("is_read", true).Error
}

// MarkAllAsRead - Đánh dấu tất cả thông báo của user là đã đọc
func MarkAllAsRead(userID uint) error {
	return database.DB.Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
