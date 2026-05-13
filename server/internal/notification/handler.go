package notification

import (
	"net/http"
	"strconv"
	authmodule "travel-backend/internal/auth"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// GetNotificationsHandler - GET /api/v1/notifications
func GetNotificationsHandler(c *gin.Context) {
	userID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Không có quyền truy cập", "UNAUTHORIZED")
		return
	}

	pagination := shared.GetPaginationParams(c)

	notifs, total, err := GetUserNotifications(userID, pagination.Offset(), pagination.Limit)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi lấy thông báo", "NOTIF_FETCH_ERROR")
		return
	}

	unreadCount, _ := GetUnreadCount(userID)

	meta := shared.BuildPaginationMeta(pagination, int(total))
	
	shared.RespondSuccessWithMeta(c, http.StatusOK, "Thành công", gin.H{
		"notifications": notifs,
		"unread_count":  unreadCount,
	}, meta)
}

// MarkAsReadHandler - PUT /api/v1/notifications/:id/read
func MarkAsReadHandler(c *gin.Context) {
	userID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Không có quyền truy cập", "UNAUTHORIZED")
		return
	}

	notifID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		shared.RespondError(c, http.StatusBadRequest, "ID thông báo không hợp lệ", "INVALID_NOTIF_ID")
		return
	}

	if err := MarkAsRead(uint(notifID), userID); err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi đánh dấu đã đọc", "NOTIF_READ_ERROR")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Đánh dấu đã đọc thành công", nil)
}

// MarkAllAsReadHandler - PUT /api/v1/notifications/read-all
func MarkAllAsReadHandler(c *gin.Context) {
	userID, ok := authmodule.GetAuthenticatedUserID(c)
	if !ok {
		shared.RespondError(c, http.StatusUnauthorized, "Không có quyền truy cập", "UNAUTHORIZED")
		return
	}

	if err := MarkAllAsRead(userID); err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Lỗi đánh dấu tất cả đã đọc", "NOTIF_READ_ALL_ERROR")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã đánh dấu tất cả là đã đọc", nil)
}
