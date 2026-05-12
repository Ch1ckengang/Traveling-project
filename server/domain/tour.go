package domain

import "time"

// Tour - Model đại diện cho bảng tours trong database
// Chứa thông tin về các tour du lịch
type Tour struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name" gorm:"not null"`
	Type           string    `json:"type" gorm:"not null;default:'domestic'"`
	Price          string    `json:"price" gorm:"not null"`
	Description    string    `json:"description"`
	Location       string    `json:"location"`
	Country        string    `json:"country" gorm:"not null;default:'Việt Nam'"`
	Duration       string    `json:"duration"`
	DepartureDate  string    `json:"departure_date" gorm:"default:''"`
	RemainingSlots int       `json:"remaining_slots" gorm:"not null;default:30"`
	Itinerary      string    `json:"itinerary"`
	Services       string    `json:"services"`
	ImageURL       string    `json:"image_url" gorm:"default:''"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
