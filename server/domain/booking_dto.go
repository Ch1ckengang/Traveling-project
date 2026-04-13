package domain

// CreateBookingRequest - Dữ liệu đặt tour từ client
type CreateBookingRequest struct {
	UserID      uint   `json:"user_id"`
	TourID      uint   `json:"tour_id" binding:"required"`
	FullName    string `json:"full_name" binding:"required"`
	Phone       string `json:"phone" binding:"required"`
	Email       string `json:"email" binding:"required"`
	AdultCount  int    `json:"adult_count"`
	ChildCount  int    `json:"child_count"`
	InfantCount int    `json:"infant_count"`
	Quantity    int    `json:"quantity"`
	TravelDate  string `json:"travel_date" binding:"required"`
	Note        string `json:"note,omitempty"`
}

// BookingResponse - Response cho API đặt tour
type BookingResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message,omitempty"`
	Booking *Booking `json:"booking,omitempty"`
}

// BookingListResponse - Response cho API lấy danh sách booking
type BookingListResponse struct {
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
	Bookings []Booking `json:"bookings,omitempty"`
}
