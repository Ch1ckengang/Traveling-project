package notification

import "travel-backend/domain"

// SendNotification - Tạo và gửi thông báo cho user
func SendNotification(userID uint, title, message, notifType string) error {
	notif := &domain.Notification{
		UserID:  userID,
		Title:   title,
		Message: message,
		Type:    notifType,
		IsRead:  false,
	}

	return CreateNotification(notif)
}

// GetUserNotifications - Lấy danh sách thông báo
func GetUserNotifications(userID uint, offset, limit int) ([]domain.Notification, int64, error) {
	return FindUserNotifications(userID, offset, limit)
}

// GetUnreadCount - Lấy số lượng chưa đọc
func GetUnreadCount(userID uint) (int64, error) {
	return FindUnreadCount(userID)
}
