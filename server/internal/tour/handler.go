package tour

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TourController - Tầng xử lý HTTP cho Tour
// Nhận request → gọi service → trả JSON

// createTours - Tạo tour mẫu khi hệ thống chưa có dữ liệu tour
func createTours() error {
	return CreateToursIfEmpty()
}

func buildTourFilter(c *gin.Context, category string) TourFilter {
	return TourFilter{
		Category: category,
		City:     c.DefaultQuery("city", c.Query("q")),
		Duration: c.DefaultQuery("duration", "all"),
		Price:    c.DefaultQuery("price", "all"),
		Sort:     c.DefaultQuery("sort", "default"),
	}
}

// GetToursHandler - Xử lý GET /api/tours
func GetToursHandler(c *gin.Context) {
	// Đảm bảo luôn có dữ liệu tour cơ bản cho chức năng đặt tour.
	if err := createTours(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể khởi tạo dữ liệu tour",
		})
		return
	}

	filter := buildTourFilter(c, c.DefaultQuery("category", "all"))

	// Gọi helper để lấy danh sách tour
	tours, err := GetToursByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách tour",
		})
		return
	}

	// Trả về danh sách tour dưới dạng JSON
	c.JSON(http.StatusOK, tours)
}

// GetDomesticToursHandler - Xử lý GET /api/tours/domestic
func GetDomesticToursHandler(c *gin.Context) {
	if err := createTours(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể khởi tạo dữ liệu tour",
		})
		return
	}

	filter := buildTourFilter(c, "domestic")

	tours, err := GetToursByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách tour du lịch Việt Nam",
		})
		return
	}

	c.JSON(http.StatusOK, tours)
}

// GetInternationalToursHandler - Xử lý GET /api/tours/international
func GetInternationalToursHandler(c *gin.Context) {
	if err := createTours(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể khởi tạo dữ liệu tour",
		})
		return
	}

	filter := buildTourFilter(c, "international")

	tours, err := GetToursByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách tour du lịch quốc tế",
		})
		return
	}

	c.JSON(http.StatusOK, tours)
}
