package tour

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TourController - Tầng xử lý HTTP cho Tour
// Nhận request → gọi service → trả JSON

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
	filter := buildTourFilter(c, c.DefaultQuery("category", "all"))

	tours, err := GetToursByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể lấy danh sách tour",
		})
		return
	}

	c.JSON(http.StatusOK, tours)
}

// GetDomesticToursHandler - Xử lý GET /api/tours/domestic
func GetDomesticToursHandler(c *gin.Context) {
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

// GetTourByIDHandler - Xử lý GET /api/tours/:id
func GetTourByIDHandler(c *gin.Context) {
	id := c.Param("id")

	tour, err := GetTourByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tour,
	})
}

// SearchToursHandler - Xử lý GET /api/tours/search?q=keyword
func SearchToursHandler(c *gin.Context) {
	keyword := c.Query("q")

	tours, err := SearchToursByKeyword(keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Không thể tìm kiếm tour",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tours,
		"total":   len(tours),
	})
}

