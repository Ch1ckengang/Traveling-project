package auth

import (
	"net/http"
	"strconv"
	"strings"
	"travel-backend/database"
	"travel-backend/domain"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// ===== ADMIN USER HANDLERS =====

// AdminGetUsersHandler - GET /v1/api/admin/users
// Lấy danh sách users với phân trang và search
func AdminGetUsersHandler(c *gin.Context) {
	pagination := shared.GetPaginationParams(c)
	search := c.DefaultQuery("search", "")
	role := c.DefaultQuery("role", "")

	query := database.DB.Model(&domain.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}
	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(email) LIKE ? OR LOWER(phone) LIKE ?", like, like, like)
	}

	var total int64
	query.Count(&total)

	var users []domain.User
	err := query.
		Order("created_at DESC").
		Offset(pagination.Offset()).
		Limit(pagination.Limit).
		Find(&users).Error

	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy danh sách users", "ADMIN_USER_FETCH_FAILED")
		return
	}

	meta := shared.BuildPaginationMeta(pagination, int(total))
	shared.RespondSuccessWithMeta(c, http.StatusOK, "", users, meta)
}

// AdminUpdateUserStatusHandler - PUT /v1/api/admin/users/:id/status
// Kích hoạt / khóa tài khoản user
func AdminUpdateUserStatusHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID user không hợp lệ", "ADMIN_USER_INVALID_ID")
		return
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "ADMIN_USER_INVALID_PAYLOAD")
		return
	}

	result := database.DB.Model(&domain.User{}).Where("id = ?", id).Update("is_active", req.IsActive)
	if result.Error != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể cập nhật trạng thái", "ADMIN_USER_UPDATE_FAILED")
		return
	}
	if result.RowsAffected == 0 {
		shared.RespondError(c, http.StatusNotFound, "Không tìm thấy user", "ADMIN_USER_NOT_FOUND")
		return
	}

	status := "kích hoạt"
	if !req.IsActive {
		status = "khóa"
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã "+status+" tài khoản thành công", gin.H{
		"user_id":   id,
		"is_active": req.IsActive,
	})
}

// AdminUpdateUserRoleHandler - PUT /v1/api/admin/users/:id/role
// Thay đổi role user (chỉ Admin)
func AdminUpdateUserRoleHandler(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id == 0 {
		shared.RespondError(c, http.StatusBadRequest, "ID user không hợp lệ", "ADMIN_USER_INVALID_ID")
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		shared.RespondError(c, http.StatusBadRequest, "Dữ liệu không hợp lệ", "ADMIN_USER_INVALID_PAYLOAD")
		return
	}

	if !domain.HasValidRole(req.Role) {
		shared.RespondError(c, http.StatusBadRequest, "Role không hợp lệ (customer/staff/admin)", "ADMIN_USER_INVALID_ROLE")
		return
	}

	result := database.DB.Model(&domain.User{}).Where("id = ?", id).Update("role", req.Role)
	if result.Error != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể thay đổi role", "ADMIN_USER_UPDATE_FAILED")
		return
	}
	if result.RowsAffected == 0 {
		shared.RespondError(c, http.StatusNotFound, "Không tìm thấy user", "ADMIN_USER_NOT_FOUND")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Đã thay đổi role thành công", gin.H{
		"user_id": id,
		"role":    req.Role,
	})
}
