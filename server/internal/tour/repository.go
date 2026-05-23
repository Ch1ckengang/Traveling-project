package tour

import (
	"strings"
	"travel-backend/database"
	"travel-backend/domain"

	"gorm.io/gorm"
)

// TourRepository - Tầng duy nhất được phép nói chuyện với database cho Tour
// Chỉ biết cách lấy/lưu dữ liệu Tour — không biết gì về HTTP hay business logic

// FindAllTours - Lấy tất cả tours từ database
func FindAllTours() ([]domain.Tour, error) {
	var tours []domain.Tour
	result := database.DB.Find(&tours)
	if result.Error != nil {
		return nil, result.Error
	}
	normalizeTourSlots(tours)
	return tours, nil
}

// FindDomesticTours - Lấy danh sách tour du lịch trong nước (Việt Nam)
func FindDomesticTours() ([]domain.Tour, error) {
	var tours []domain.Tour
	result := database.DB.Where("type = ? OR country = ? OR country = ? OR country = ''", "domestic", "Việt Nam", "Vietnam").Find(&tours)
	if result.Error != nil {
		return nil, result.Error
	}
	normalizeTourSlots(tours)
	return tours, nil
}

// FindToursByCategory - Lấy danh sách tour theo nhóm nghiệp vụ trên menu
func FindToursByCategory(category string) ([]domain.Tour, error) {
	var tours []domain.Tour
	query := database.DB

	switch strings.ToLower(category) {
	case "domestic":
		query = query.Where(
			"LOWER(TRIM(type)) = ? OR LOWER(REPLACE(TRIM(country), ' ', '')) IN (?, ?, ?) OR TRIM(country) = ''",
			"domestic",
			"vietnam",
			"vietnam",
			"việtnam",
		)
	case "international":
		// Chỉ lấy tour quốc tế và loại trừ rõ ràng các tour Việt Nam.
		query = query.Where(
			"LOWER(TRIM(type)) = ? AND LOWER(REPLACE(TRIM(country), ' ', '')) NOT IN (?, ?, ?) AND TRIM(country) <> ''",
			"international",
			"vietnam",
			"vietnam",
			"việtnam",
		)
	case "service":
		query = query.Where("type = ?", "service")
	}

	if err := query.Find(&tours).Error; err != nil {
		return nil, err
	}

	normalizeTourSlots(tours)

	return tours, nil
}

// FindTourByID - Tìm tour theo id (uint)
func FindTourByID(id uint) (*domain.Tour, error) {
	var tour domain.Tour
	if err := database.DB.Preload("Images").First(&tour, id).Error; err != nil {
		return nil, err
	}

	if tour.RemainingSlots <= 0 {
		tour.RemainingSlots = 30
	}

	return &tour, nil
}

// FindTourByIDString - Tìm tour theo id (string từ URL param)
func FindTourByIDString(id string) (*domain.Tour, error) {
	var tour domain.Tour
	if err := database.DB.Preload("Images").First(&tour, id).Error; err != nil {
		return nil, err
	}

	if tour.RemainingSlots <= 0 {
		tour.RemainingSlots = 30
	}

	return &tour, nil
}

// SearchTours - Tìm kiếm tour theo keyword (name, location, description)
func SearchTours(keyword string) ([]domain.Tour, error) {
	var tours []domain.Tour
	like := "%" + strings.ToLower(keyword) + "%"
	err := database.DB.
		Where("LOWER(name) LIKE ? OR LOWER(location) LIKE ? OR LOWER(description) LIKE ? OR LOWER(country) LIKE ?",
			like, like, like, like).
		Find(&tours).Error
	if err != nil {
		return nil, err
	}
	normalizeTourSlots(tours)
	return tours, nil
}

// CountTours - Đếm tổng số tour, optional filter theo category.
func CountTours(category string) (int64, error) {
	var count int64
	query := database.DB.Model(&domain.Tour{})

	switch strings.ToLower(strings.TrimSpace(category)) {
	case "domestic":
		query = query.Where("type = ?", "domestic")
	case "international":
		query = query.Where("type = ?", "international")
	case "service":
		query = query.Where("type = ?", "service")
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func normalizeTourSlots(tours []domain.Tour) {
	for i := range tours {
		if tours[i].RemainingSlots <= 0 {
			tours[i].RemainingSlots = 30
		}
	}
}

// DecreaseTourRemainingSlots - Trừ số chỗ còn lại sau khi tạo booking.
func DecreaseTourRemainingSlots(tourID uint, seats int) error {
	if seats <= 0 {
		return nil
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var tour domain.Tour
		if err := tx.First(&tour, tourID).Error; err != nil {
			return err
		}

		remaining := tour.RemainingSlots
		if remaining <= 0 {
			remaining = 30
		}

		if remaining < seats {
			return gorm.ErrInvalidData
		}

		tour.RemainingSlots = remaining - seats
		return tx.Model(&tour).Update("remaining_slots", tour.RemainingSlots).Error
	})
}

// FindTourScheduleByID
func FindTourScheduleByID(id uint) (*domain.TourSchedule, error) {
	var schedule domain.TourSchedule
	if err := database.DB.First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

// DecreaseTourScheduleRemainingSlots
func DecreaseTourScheduleRemainingSlots(scheduleID uint, seats int) error {
	if seats <= 0 {
		return nil
	}

	return database.DB.Transaction(func(tx *gorm.DB) error {
		var schedule domain.TourSchedule
		if err := tx.First(&schedule, scheduleID).Error; err != nil {
			return err
		}

		if schedule.RemainingSlots < seats {
			return gorm.ErrInvalidData
		}

		schedule.RemainingSlots -= seats
		return tx.Save(&schedule).Error
	})
}

// CreateToursIfEmpty - Tạo dữ liệu tour mẫu nếu bảng tours đang trống
func seedToursIfEmpty() error {
	if err := database.DB.AutoMigrate(&domain.Tour{}); err != nil {
		return err
	}

	tours := []domain.Tour{
		{Name: "Tour Đà Nẵng - Hội An", Type: "domestic", Price: "2.000.000đ", Location: "Đà Nẵng", Country: "Việt Nam", Duration: "3 ngày 2 đêm", Description: "Khám phá phố cổ Hội An và biển Mỹ Khê."},
		{Name: "Tour Hà Nội - Sa Pa", Type: "domestic", Price: "3.500.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Trải nghiệm khí hậu vùng cao và bản làng Tây Bắc."},
		{Name: "Tour Phú Quốc", Type: "domestic", Price: "5.000.000đ", Location: "Phú Quốc", Country: "Việt Nam", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo và thưởng thức hải sản địa phương."},
		{Name: "Tour Nha Trang - Đà Lạt", Type: "domestic", Price: "4.800.000đ", Location: "Nha Trang", Country: "Việt Nam", Duration: "4 ngày 3 đêm", Description: "Kết hợp nghỉ biển Nha Trang và khí hậu mát mẻ Đà Lạt."},
		{Name: "Tour Bangkok - Pattaya", Type: "international", Price: "8.200.000đ", Location: "Bangkok", Country: "Thái Lan", Duration: "5 ngày 4 đêm", Description: "Lộ trình quốc tế phù hợp gia đình và nhóm bạn."},
		{Name: "Tour Seoul Mùa Hoa", Type: "international", Price: "12.500.000đ", Location: "Seoul", Country: "Hàn Quốc", Duration: "6 ngày 5 đêm", Description: "Tham quan cung điện, phố mua sắm và ẩm thực Hàn."},
		{Name: "Tour Tokyo - Núi Phú Sĩ", Type: "international", Price: "15.900.000đ", Location: "Tokyo", Country: "Nhật Bản", Duration: "6 ngày 5 đêm", Description: "Khám phá Tokyo hiện đại và trải nghiệm văn hóa Nhật Bản."},
		{Name: "Tour Paris - Lyon", Type: "international", Price: "18.500.000đ", Location: "Paris", Country: "Pháp", Duration: "7 ngày 6 đêm", Description: "Hành trình châu Âu với điểm nhấn ẩm thực và kiến trúc cổ điển."},
		{Name: "Tour Singapore - Sentosa", Type: "international", Price: "9.600.000đ", Location: "Singapore", Country: "Singapore", Duration: "4 ngày 3 đêm", Description: "Khám phá đảo quốc sư tử với lịch trình hiện đại và thân thiện gia đình."},
		{Name: "Tour Bali - Ubud", Type: "international", Price: "10.800.000đ", Location: "Bali", Country: "Indonesia", Duration: "5 ngày 4 đêm", Description: "Nghỉ dưỡng biển đảo, check-in ruộng bậc thang và đền cổ Bali."},
		{Name: "Tour Sydney - Melbourne", Type: "international", Price: "21.900.000đ", Location: "Sydney", Country: "Úc", Duration: "7 ngày 6 đêm", Description: "Hành trình nước Úc qua hai thành phố biểu tượng."},
		{Name: "Tour Dubai - Abu Dhabi", Type: "international", Price: "19.500.000đ", Location: "Dubai", Country: "UAE", Duration: "6 ngày 5 đêm", Description: "Trải nghiệm thành phố xa hoa và văn hóa Trung Đông đặc sắc."},
		{Name: "Combo Visa + Vé Máy Bay", Type: "service", Price: "1.800.000đ", Location: "Hồ Chí Minh", Country: "Việt Nam", Duration: "2 ngày", Description: "Dịch vụ làm visa nhanh và hỗ trợ đặt vé trọn gói."},
		{Name: "Đưa Đón Sân Bay Cao Cấp", Type: "service", Price: "900.000đ", Location: "Hà Nội", Country: "Việt Nam", Duration: "Trong ngày", Description: "Đưa đón đúng giờ với xe riêng và tài xế kinh nghiệm."},
	}

	for _, tour := range tours {
		var exists int64
		if err := database.DB.Model(&domain.Tour{}).Where("name = ?", tour.Name).Count(&exists).Error; err != nil {
			return err
		}

		if exists > 0 {
			if err := database.DB.Model(&domain.Tour{}).
				Where("name = ?", tour.Name).
				Updates(map[string]interface{}{
					"type":        tour.Type,
					"country":     tour.Country,
					"location":    tour.Location,
					"duration":    tour.Duration,
					"description": tour.Description,
				}).Error; err != nil {
				return err
			}
			continue
		}

		if err := database.DB.Create(&tour).Error; err != nil {
			return err
		}
	}

	return nil
}
