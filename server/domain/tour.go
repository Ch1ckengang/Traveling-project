package domain

import (
	"fmt"
	"strings"
	"time"
)

// Tour - Model đại diện cho bảng tours trong database
// Chứa thông tin về các tour du lịch
type Tour struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	Slug           string    `json:"slug" gorm:"uniqueIndex;size:255"`
	Type           string    `json:"type" gorm:"not null;default:'domestic';index"`
	PriceAmount    int64     `json:"price_amount" gorm:"not null;default:0"`
	Price          string    `json:"price" gorm:"not null;default:''"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	Country        string    `json:"country" gorm:"not null;default:'Việt Nam'"`
	Duration       string    `json:"duration"`
	DepartureDate  string    `json:"departure_date" gorm:"default:''"`
	RemainingSlots int       `json:"remaining_slots" gorm:"not null;default:30"`
	Rating         float64   `json:"rating" gorm:"not null;default:0"`
	ReviewCount    int       `json:"review_count" gorm:"not null;default:0"`
	IsActive       bool      `json:"is_active" gorm:"not null;default:true"`
	Itinerary      string    `json:"itinerary"`
	Services       string    `json:"services"`
	ImageURL       string    `json:"image_url" gorm:"default:''"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// FormatPriceVND - Format giá thành chuỗi VND đẹp (ví dụ: "2.000.000đ")
func FormatPriceVND(amount int64) string {
	if amount <= 0 {
		return "Liên hệ"
	}

	str := fmt.Sprintf("%d", amount)
	n := len(str)
	if n <= 3 {
		return str + "đ"
	}

	var result strings.Builder
	for i, ch := range str {
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune('.')
		}
		result.WriteRune(ch)
	}
	result.WriteString("đ")
	return result.String()
}

// GenerateSlug - Tạo slug từ tên tour (đơn giản)
func GenerateSlug(name string) string {
	slug := strings.ToLower(strings.TrimSpace(name))
	slug = strings.ReplaceAll(slug, " ", "-")
	// Loại bỏ ký tự đặc biệt cơ bản
	replacer := strings.NewReplacer(
		"đ", "d", "Đ", "d",
		"à", "a", "á", "a", "ả", "a", "ã", "a", "ạ", "a",
		"ă", "a", "ắ", "a", "ằ", "a", "ẳ", "a", "ẵ", "a", "ặ", "a",
		"â", "a", "ấ", "a", "ầ", "a", "ẩ", "a", "ẫ", "a", "ậ", "a",
		"è", "e", "é", "e", "ẻ", "e", "ẽ", "e", "ẹ", "e",
		"ê", "e", "ế", "e", "ề", "e", "ể", "e", "ễ", "e", "ệ", "e",
		"ì", "i", "í", "i", "ỉ", "i", "ĩ", "i", "ị", "i",
		"ò", "o", "ó", "o", "ỏ", "o", "õ", "o", "ọ", "o",
		"ô", "o", "ố", "o", "ồ", "o", "ổ", "o", "ỗ", "o", "ộ", "o",
		"ơ", "o", "ớ", "o", "ờ", "o", "ở", "o", "ỡ", "o", "ợ", "o",
		"ù", "u", "ú", "u", "ủ", "u", "ũ", "u", "ụ", "u",
		"ư", "u", "ứ", "u", "ừ", "u", "ử", "u", "ữ", "u", "ự", "u",
		"ỳ", "y", "ý", "y", "ỷ", "y", "ỹ", "y", "ỵ", "y",
	)
	slug = replacer.Replace(slug)
	return slug
}
