package repositories

import (
	"strings"
	"travel-backend/database"
	"travel-backend/models"
)

// TourRepository - Tầng duy nhất được phép nói chuyện với database cho Tour
// Chỉ biết cách lấy/lưu dữ liệu Tour — không biết gì về HTTP hay business logic

// FindAllTours - Lấy tất cả tours từ database
func FindAllTours() ([]models.Tour, error) {
	var tours []models.Tour
	result := database.DB.Find(&tours)
	if result.Error != nil {
		return nil, result.Error
	}
	return tours, nil
}

// FindDomesticTours - Lấy danh sách tour du lịch trong nước (Việt Nam)
func FindDomesticTours() ([]models.Tour, error) {
	var tours []models.Tour
	result := database.DB.Where("type = ? OR country = ? OR country = ? OR country = ''", "domestic", "Việt Nam", "Vietnam").Find(&tours)
	if result.Error != nil {
		return nil, result.Error
	}
	return tours, nil
}

// FindToursByCategory - Lấy danh sách tour theo nhóm nghiệp vụ trên menu
func FindToursByCategory(category string) ([]models.Tour, error) {
	var tours []models.Tour
	query := database.DB

	switch strings.ToLower(category) {
	case "domestic":
		query = query.Where("type = ? OR country = ? OR country = ? OR country = ''", "domestic", "Việt Nam", "Vietnam")
	case "international":
		query = query.Where("type = ? OR (country <> ? AND country <> ? AND country <> '')", "international", "Việt Nam", "Vietnam")
	case "service":
		query = query.Where("type = ?", "service")
	}

	if err := query.Find(&tours).Error; err != nil {
		return nil, err
	}

	return tours, nil
}

// FindTourByID - Tìm tour theo id
func FindTourByID(id uint) (*models.Tour, error) {
	var tour models.Tour
	if err := database.DB.First(&tour, id).Error; err != nil {
		return nil, err
	}
	return &tour, nil
}

// CreateToursIfEmpty - Tạo dữ liệu tour mẫu nếu bảng tours đang trống
func CreateToursIfEmpty() error {
	if err := database.DB.AutoMigrate(&models.Tour{}); err != nil {
		return err
	}

	tours := []models.Tour{
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
		if err := database.DB.Model(&models.Tour{}).Where("name = ?", tour.Name).Count(&exists).Error; err != nil {
			return err
		}

		if exists > 0 {
			if err := database.DB.Model(&models.Tour{}).
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
