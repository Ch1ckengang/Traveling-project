package auth

import (
	"net/http"
	"strconv"
	"strings"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

const contextUserIDKey = "auth_user_id"

// AuthRequired - middleware yêu cầu Bearer access token hợp lệ.
// Sau khi verify token, set user_id và role vào context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if rawHeader == "" {
			respondError(c, http.StatusUnauthorized, shared.ErrMissingBearerToken.Error(), "AUTH_TOKEN_MISSING")
			c.Abort()
			return
		}

		parts := strings.SplitN(rawHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			respondError(c, http.StatusUnauthorized, shared.ErrMissingBearerToken.Error(), "AUTH_TOKEN_MISSING")
			c.Abort()
			return
		}

		claims, err := parseAccessToken(parts[1])
		if err != nil {
			respondError(c, http.StatusUnauthorized, shared.ErrInvalidAccessToken.Error(), "AUTH_TOKEN_INVALID")
			c.Abort()
			return
		}

		c.Set(contextUserIDKey, claims.UserID)
		c.Set("auth_user_role", claims.Role)
		c.Next()
	}
}

// GetAuthenticatedUserID - lấy user id từ access token đã verify.
func GetAuthenticatedUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get(contextUserIDKey)
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	if !ok || userID == 0 {
		return 0, false
	}

	return userID, true
}

// GetAuthenticatedUserRole - lấy role từ access token đã verify.
func GetAuthenticatedUserRole(c *gin.Context) string {
	value, exists := c.Get("auth_user_role")
	if !exists {
		return ""
	}

	role, ok := value.(string)
	if !ok {
		return ""
	}

	return role
}

// EnsurePathUserMatchesToken - đảm bảo /users/:id là chính chủ.
func EnsurePathUserMatchesToken(c *gin.Context, paramName string) bool {
	authUserID, ok := GetAuthenticatedUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, shared.ErrInvalidAccessToken.Error(), "AUTH_TOKEN_INVALID")
		return false
	}

	rawID := strings.TrimSpace(c.Param(paramName))
	parsedID, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil || parsedID == 0 {
		respondError(c, http.StatusBadRequest, shared.ErrInvalidUserID.Error(), "AUTH_USER_ID_INVALID")
		return false
	}

	if uint(parsedID) != authUserID {
		// Admin có thể truy cập tài nguyên của user khác
		role := GetAuthenticatedUserRole(c)
		if role == "admin" {
			return true
		}
		respondError(c, http.StatusForbidden, shared.ErrForbiddenResource.Error(), "AUTH_FORBIDDEN_RESOURCE")
		return false
	}

	return true
}
