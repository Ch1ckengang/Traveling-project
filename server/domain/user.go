package domain

import "time"

// Role constants
const (
	RoleCustomer = "customer"
	RoleStaff    = "staff"
	RoleAdmin    = "admin"
)

// User - Model đại diện cho bảng users trong database
// GORM sẽ tự động tạo bảng dựa trên struct này
type User struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	Name            string    `json:"name" gorm:"not null"`
	Email           string    `json:"email" gorm:"unique;not null"`
	Password        string    `json:"-" gorm:"not null"`
	Phone           string    `json:"phone" gorm:"default:''"`
	AvatarURL       string    `json:"avatar_url" gorm:"default:''"`
	Role            string    `json:"role" gorm:"not null;default:'customer';size:20;index"`
	IsEmailVerified bool      `json:"is_email_verified" gorm:"not null;default:false"`
	IsActive        bool      `json:"is_active" gorm:"not null;default:true"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// IsAdmin - Kiểm tra user có phải admin không
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsStaff - Kiểm tra user có phải staff không
func (u *User) IsStaff() bool {
	return u.Role == RoleStaff
}

// IsStaffOrAbove - Kiểm tra user có quyền Staff hoặc Admin
func (u *User) IsStaffOrAbove() bool {
	return u.Role == RoleStaff || u.Role == RoleAdmin
}

// HasValidRole - Kiểm tra role hợp lệ
func HasValidRole(role string) bool {
	return role == RoleCustomer || role == RoleStaff || role == RoleAdmin
}
