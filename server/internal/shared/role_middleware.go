package shared

import (
	"net/http"
	"travel-backend/domain"

	"github.com/gin-gonic/gin"
)

const contextUserRoleKey = "auth_user_role"

// RoleRequired - Middleware kiểm tra user có role phù hợp
// Yêu cầu AuthRequired() đã chạy trước để có user_id trong context
func RoleRequired(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy user_id từ context (đã set bởi AuthRequired)
		userIDValue, exists := c.Get("auth_user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Yêu cầu đăng nhập",
				"error":   &APIError{Code: "AUTH_TOKEN_MISSING"},
			})
			c.Abort()
			return
		}

		userID, ok := userIDValue.(uint)
		if !ok || userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token không hợp lệ",
				"error":   &APIError{Code: "AUTH_TOKEN_INVALID"},
			})
			c.Abort()
			return
		}

		// Lấy role từ context (đã set bởi AuthRequired nếu có)
		roleValue, roleExists := c.Get(contextUserRoleKey)
		if !roleExists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Không có quyền truy cập",
				"error":   &APIError{Code: "AUTH_FORBIDDEN"},
			})
			c.Abort()
			return
		}

		userRole, ok := roleValue.(string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Không có quyền truy cập",
				"error":   &APIError{Code: "AUTH_FORBIDDEN"},
			})
			c.Abort()
			return
		}

		// Admin luôn được phép
		if userRole == domain.RoleAdmin {
			c.Next()
			return
		}

		// Kiểm tra role có nằm trong danh sách cho phép
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Bạn không có quyền thực hiện thao tác này",
			"error":   &APIError{Code: "AUTH_INSUFFICIENT_ROLE"},
		})
		c.Abort()
	}
}

// StaffRequired - Middleware yêu cầu role Staff hoặc Admin
func StaffRequired() gin.HandlerFunc {
	return RoleRequired(domain.RoleStaff, domain.RoleAdmin)
}

// AdminRequired - Middleware yêu cầu role Admin
func AdminRequired() gin.HandlerFunc {
	return RoleRequired(domain.RoleAdmin)
}
