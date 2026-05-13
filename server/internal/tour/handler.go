package tour

import (
	"net/http"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// TourController - Tầng xử lý HTTP cho Tour
// Nhận request → gọi service → trả JSON (sử dụng response chuẩn)

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
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy danh sách tour", "TOUR_FETCH_FAILED")
		return
	}

	shared.RespondSuccessWithMeta(c, http.StatusOK, "", tours, map[string]interface{}{
		"total": len(tours),
	})
}

// GetDomesticToursHandler - Xử lý GET /api/tours/domestic
func GetDomesticToursHandler(c *gin.Context) {
	filter := buildTourFilter(c, "domestic")

	tours, err := GetToursByFilter(filter)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy danh sách tour du lịch Việt Nam", "TOUR_FETCH_FAILED")
		return
	}

	shared.RespondSuccessWithMeta(c, http.StatusOK, "", tours, map[string]interface{}{
		"total": len(tours),
	})
}

// GetInternationalToursHandler - Xử lý GET /api/tours/international
func GetInternationalToursHandler(c *gin.Context) {
	filter := buildTourFilter(c, "international")

	tours, err := GetToursByFilter(filter)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy danh sách tour du lịch quốc tế", "TOUR_FETCH_FAILED")
		return
	}

	shared.RespondSuccessWithMeta(c, http.StatusOK, "", tours, map[string]interface{}{
		"total": len(tours),
	})
}

// GetTourByIDHandler - Xử lý GET /api/tours/:id
func GetTourByIDHandler(c *gin.Context) {
	id := c.Param("id")

	tour, err := GetTourByID(id)
	if err != nil {
		shared.RespondError(c, http.StatusNotFound, err.Error(), "TOUR_NOT_FOUND")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "", gin.H{"tour": tour})
}

// SearchToursHandler - Xử lý GET /api/tours/search?q=keyword
func SearchToursHandler(c *gin.Context) {
	keyword := c.Query("q")

	tours, err := SearchToursByKeyword(keyword)
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể tìm kiếm tour", "TOUR_SEARCH_FAILED")
		return
	}

	shared.RespondSuccessWithMeta(c, http.StatusOK, "", tours, map[string]interface{}{
		"total": len(tours),
	})
}
