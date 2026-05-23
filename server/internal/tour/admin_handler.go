package tour

import (
	"net/http"
	"strconv"
	"log"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// ===== ADMIN TOUR HANDLERS =====
// Yêu cầu StaffRequired() middleware đã chạy trước

// AdminGetToursHandler - GET /v1/api/admin/tours
// Lấy tất cả tours (bao gồm cả inactive) với phân trang
func AdminGetToursHandler(c *gin.Context) {
	pagination := shared.GetPaginationParams(c)
	category := c.DefaultQuery("type", "")
	search := c.DefaultQuery("search", "")

	query := database.DB.Model(&domain.Tour{})

	if category != "" {
		query = query.Where("type = ?", category)
	}
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(location) LIKE ?", like, like)
	}

	var total int64
	query.Count(&total)

	var tours []domain.Tour
	err := query.
		Order("created_at DESC").
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		Find(&tours).Error

	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy danh sách tour", "ADMIN_TOUR_FETCH_FAILED")
		return
	}

	meta := shared.BuildPaginationMeta(pagination, int(total))
	shared.RespondSuccessWithMeta(c, http.StatusOK, "", tours, meta)
}

// AdminCreateTourHandler - POST /v1/api/admin/tours
// Tạo tour mới
func AdminCreateTourHandler(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Type          string `json:"type" binding:"required"`
		PriceAmount   int64  `json:"price_amount" binding:"required,min=0"`
		Description   string `json:"description"`
		Location      string `json:"location"`
		Country       string `json:"country"`
		Duration      string   `json:"duration"`
		DepartureDate string   `json:"departure_date"`
		Itinerary     string   `json:"itinerary"`
		Services      string   `json:"services"`
		ImageURL      string   `json:"image_url"`
		Images        []string `json:"images"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "ADMIN_TOUR_INVALID_PAYLOAD")
		return
	}

	tour := domain.Tour{
		Name:          req.Name,
		Slug:          domain.GenerateSlug(req.Name),
		Type:          req.Type,
		PriceAmount:   req.PriceAmount,
		Price:         domain.FormatPriceVND(req.PriceAmount),
		Description:   req.Description,
		Location:      req.Location,
		Country:       req.Country,
		Duration:      req.Duration,
		DepartureDate: req.DepartureDate,
		Itinerary:     req.Itinerary,
		Services:      req.Services,
		ImageURL:      req.ImageURL,
		IsActive:      true,
		RemainingSlots: 30,
	}

	if err := database.DB.Create(&tour).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể tạo tour", "ADMIN_TOUR_CREATE_FAILED")
		return
	}

	// Insert images
	if len(req.Images) > 0 {
		var tourImages []domain.TourImage
		for _, imgURL := range req.Images {
			tourImages = append(tourImages, domain.TourImage{
				TourID:   tour.ID,
				ImageURL: imgURL,
			})
		}
		if err := database.DB.Create(&tourImages).Error; err != nil {
			// Just log the error, don't fail the creation
			log.Printf("Warning: failed to insert tour images: %v", err)
		}
	}
	
	// Reload tour to include images
	database.DB.Preload("Images").First(&tour, tour.ID)

	shared.RespondSuccess(c, http.StatusCreated, "Tạo tour thành công", gin.H{"tour": tour})
}

// AdminUpdateTourHandler - PUT /v1/api/admin/tours/:id
// Cập nhật thông tin tour
func AdminUpdateTourHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID tour không hợp lệ", "ADMIN_TOUR_INVALID_ID")
		return
	}

	var tour domain.Tour
	if err := database.DB.First(&tour, id).Error; err != nil {
		shared.RespondError(c, http.StatusNotFound, "Không tìm thấy tour", "ADMIN_TOUR_NOT_FOUND")
		return
	}

	var req struct {
		Name          *string `json:"name"`
		Type          *string `json:"type"`
		PriceAmount   *int64  `json:"price_amount"`
		Description   *string `json:"description"`
		Location      *string `json:"location"`
		Country       *string `json:"country"`
		Duration      *string `json:"duration"`
		DepartureDate *string `json:"departure_date"`
		Itinerary     *string   `json:"itinerary"`
		Services      *string   `json:"services"`
		ImageURL      *string   `json:"image_url"`
		Images        *[]string `json:"images"`
		RemainingSlots *int     `json:"remaining_slots"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "ADMIN_TOUR_INVALID_PAYLOAD")
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
		updates["slug"] = domain.GenerateSlug(*req.Name)
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.PriceAmount != nil {
		updates["price_amount"] = *req.PriceAmount
		updates["price"] = domain.FormatPriceVND(*req.PriceAmount)
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Location != nil {
		updates["location"] = *req.Location
	}
	if req.Country != nil {
		updates["country"] = *req.Country
	}
	if req.Duration != nil {
		updates["duration"] = *req.Duration
	}
	if req.DepartureDate != nil {
		updates["departure_date"] = *req.DepartureDate
	}
	if req.Itinerary != nil {
		updates["itinerary"] = *req.Itinerary
	}
	if req.Services != nil {
		updates["services"] = *req.Services
	}
	if req.ImageURL != nil {
		updates["image_url"] = *req.ImageURL
	}
	if req.RemainingSlots != nil {
		updates["remaining_slots"] = *req.RemainingSlots
	}

	if len(updates) == 0 {
		shared.RespondError(c, http.StatusBadRequest, "Không có dữ liệu cập nhật", "ADMIN_TOUR_NO_UPDATES")
		return
	}

	if err := database.DB.Model(&tour).Updates(updates).Error; err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể cập nhật tour", "ADMIN_TOUR_UPDATE_FAILED")
		return
	}

	// Xử lý update mảng ảnh
	if req.Images != nil {
		// Xóa các ảnh cũ
		database.DB.Where("tour_id = ?", tour.ID).Delete(&domain.TourImage{})
		
		// Thêm ảnh mới
		if len(*req.Images) > 0 {
			var tourImages []domain.TourImage
			for _, imgURL := range *req.Images {
				tourImages = append(tourImages, domain.TourImage{
					TourID:   tour.ID,
					ImageURL: imgURL,
				})
			}
			database.DB.Create(&tourImages)
		}
	}

	// Reload tour
	database.DB.Preload("Images").First(&tour, id)
	shared.RespondSuccess(c, http.StatusOK, "Cập nhật tour thành công", gin.H{"tour": tour})
}

// AdminDeleteTourHandler - DELETE /v1/api/admin/tours/:id
// Soft delete (set is_active = false)
func AdminDeleteTourHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID tour không hợp lệ", "ADMIN_TOUR_INVALID_ID")
		return
	}

	result := database.DB.Model(&domain.Tour{}).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể xóa tour", "ADMIN_TOUR_DELETE_FAILED")
		return
	}
	if result.RowsAffected == 0 {
		shared.RespondError(c, http.StatusNotFound, "Không tìm thấy tour", "ADMIN_TOUR_NOT_FOUND")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã ẩn tour thành công", nil)
}

// AdminToggleTourHandler - PUT /v1/api/admin/tours/:id/toggle
// Ẩn/hiện tour
func AdminToggleTourHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID tour không hợp lệ", "ADMIN_TOUR_INVALID_ID")
		return
	}

	var tour domain.Tour
	if err := database.DB.First(&tour, id).Error; err != nil {
		shared.RespondError(c, http.StatusNotFound, "Không tìm thấy tour", "ADMIN_TOUR_NOT_FOUND")
		return
	}

	newStatus := !tour.IsActive
	database.DB.Model(&tour).Update("is_active", newStatus)

	status := "hiện"
	if !newStatus {
		status = "ẩn"
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã "+status+" tour thành công", gin.H{
		"tour_id":   tour.ID,
		"is_active": newStatus,
	})
}
