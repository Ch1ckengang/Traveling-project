package shared

import (
	"github.com/gin-gonic/gin"
)

// APIError - Cấu trúc lỗi chuẩn cho API responses
type APIError struct {
	Code    string      `json:"code"`
	Details interface{} `json:"details,omitempty"`
}

// APIResponse - Cấu trúc response chuẩn cho toàn bộ API
type APIResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message,omitempty"`
	Data    interface{}            `json:"data"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
	Error   *APIError              `json:"error"`
}

// RespondSuccess - Trả response thành công với format chuẩn
func RespondSuccess(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    BuildMeta(c),
		Error:   nil,
	})
}

// RespondSuccessWithMeta - Trả response thành công kèm metadata (pagination, etc.)
func RespondSuccessWithMeta(c *gin.Context, status int, message string, data interface{}, meta map[string]interface{}) {
	baseMeta := BuildMeta(c)
	for k, v := range meta {
		baseMeta[k] = v
	}
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    baseMeta,
		Error:   nil,
	})
}

// RespondError - Trả response lỗi với format chuẩn
func RespondError(c *gin.Context, status int, message, code string) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Meta:    BuildMeta(c),
		Error:   &APIError{Code: code},
	})
}

// RespondErrorWithDetails - Trả response lỗi kèm chi tiết
func RespondErrorWithDetails(c *gin.Context, status int, message, code string, details interface{}) {
	c.JSON(status, APIResponse{
		Success: false,
		Message: message,
		Data:    nil,
		Meta:    BuildMeta(c),
		Error:   &APIError{Code: code, Details: details},
	})
}

// BuildMeta - Tạo metadata cho response (request_id, etc.)
func BuildMeta(c *gin.Context) map[string]interface{} {
	meta := make(map[string]interface{})
	requestID := c.GetHeader("X-Request-ID")
	if requestID != "" {
		meta["request_id"] = requestID
	}
	return meta
}
