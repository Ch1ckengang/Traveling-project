package controllers

import (
	"net/http"
	"travel-backend/models"
	"travel-backend/services"

	"github.com/gin-gonic/gin"
)

// TourController - Tầng xử lý HTTP cho Tour
// Nhận request → gọi service → trả JSON

// createTours - Tạo tour mẫu khi hệ thống chưa có dữ liệu tour
func createTours() error {
	return services.CreateToursIfEmpty()
}

// getTours - Lấy danh sách tour phục vụ màn hình đặt tour
func getTours() ([]models.Tour, error) {
	return services.GetAllTours()
}

// getDomesticTours - Lấy danh sách tour du lịch trong nước
func getDomesticTours() ([]models.Tour, error) {
	return services.GetDomesticTours()
}

// GetTours - Xử lý GET /api/tours
func GetTours(c *gin.Context) {
	// Đảm bảo luôn có dữ liệu tour cơ bản cho chức năng đặt tour.
	if err := createTours(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể khởi tạo dữ liệu tour",
		})
		return
	}

	filter := services.TourFilter{
		Category: c.DefaultQuery("category", "all"),
		City:     c.DefaultQuery("city", c.Query("q")),
		Duration: c.DefaultQuery("duration", "all"),
		Price:    c.DefaultQuery("price", "all"),
		Sort:     c.DefaultQuery("sort", "default"),
	}

	// Gọi helper để lấy danh sách tour
	tours, err := services.GetToursByFilter(filter)
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

// GetDomesticTours - Xử lý GET /api/tours/domestic
func GetDomesticTours(c *gin.Context) {
	if err := createTours(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể khởi tạo dữ liệu tour",
		})
		return
	}

	filter := services.TourFilter{
		Category: "domestic",
		City:     c.DefaultQuery("city", c.Query("q")),
		Duration: c.DefaultQuery("duration", "all"),
		Price:    c.DefaultQuery("price", "all"),
		Sort:     c.DefaultQuery("sort", "default"),
	}

	tours, err := services.GetToursByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách tour du lịch Việt Nam",
		})
		return
	}

	c.JSON(http.StatusOK, tours)
}
