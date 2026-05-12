package shared

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter - Simple in-memory rate limiter
// For production, consider using Redis-based rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

// NewRateLimiter - Tạo rate limiter mới
// limit: số request tối đa trong window
// window: khoảng thời gian (ví dụ: 1 phút)
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Cleanup expired entries every minute
	go rl.cleanup()

	return rl
}

// Allow - Kiểm tra xem IP có được phép request không
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	windowStart := now.Add(-rl.window)

	// Lấy danh sách requests của IP này
	requests := rl.requests[ip]

	// Lọc bỏ các requests ngoài window
	validRequests := []time.Time{}
	for _, reqTime := range requests {
		if reqTime.After(windowStart) {
			validRequests = append(validRequests, reqTime)
		}
	}

	// Kiểm tra có vượt limit không
	if len(validRequests) >= rl.limit {
		rl.requests[ip] = validRequests
		return false
	}

	// Thêm request mới
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests

	return true
}

// cleanup - Xóa các entries cũ để tránh memory leak
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		windowStart := now.Add(-rl.window)

		for ip, requests := range rl.requests {
			validRequests := []time.Time{}
			for _, reqTime := range requests {
				if reqTime.After(windowStart) {
					validRequests = append(validRequests, reqTime)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware - Gin middleware cho rate limiting
func RateLimitMiddleware(limiter *RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "Quá nhiều yêu cầu. Vui lòng thử lại sau.",
				"error": gin.H{
					"code": "RATE_LIMIT_EXCEEDED",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
