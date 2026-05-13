package shared

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams - Tham số phân trang lấy từ query string
type PaginationParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// Offset - Tính offset cho SQL query
func (p PaginationParams) Offset() int {
	return (p.Page - 1) * p.Limit
}

// PaginationMeta - Metadata phân trang trả về cho frontend
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// GetPaginationParams - Lấy tham số phân trang từ query string
// Mặc định: page=1, limit=10. Giới hạn: limit tối đa 100
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return PaginationParams{Page: page, Limit: limit}
}

// BuildPaginationMeta - Tạo metadata phân trang từ total count
func BuildPaginationMeta(params PaginationParams, total int) map[string]interface{} {
	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	return map[string]interface{}{
		"page":        params.Page,
		"limit":       params.Limit,
		"total":       total,
		"total_pages": totalPages,
	}
}
