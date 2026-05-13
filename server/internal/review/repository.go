package review

import (
	"travel-backend/database"
	"travel-backend/domain"
)

// FindReviewsByTourID - Lấy reviews theo tour với phân trang
func FindReviewsByTourID(tourID uint, rating int, offset, limit int) ([]domain.Review, int64, error) {
	query := database.DB.Model(&domain.Review{}).
		Where("tour_id = ? AND status = ?", tourID, domain.ReviewStatusPublished).
		Preload("User")

	if rating >= 1 && rating <= 5 {
		query = query.Where("rating = ?", rating)
	}

	var total int64
	query.Count(&total)

	var reviews []domain.Review
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&reviews).Error

	return reviews, total, err
}

// FindReviewByID - Lấy review theo ID
func FindReviewByID(id uint) (*domain.Review, error) {
	var review domain.Review
	err := database.DB.Preload("User").Preload("Tour").First(&review, id).Error
	return &review, err
}

// FindReviewByBookingID - Kiểm tra user đã review booking này chưa
func FindReviewByBookingID(bookingID, userID uint) (*domain.Review, error) {
	var review domain.Review
	err := database.DB.
		Where("booking_id = ? AND user_id = ?", bookingID, userID).
		First(&review).Error
	return &review, err
}

// CreateReview - Tạo review mới
func CreateReview(review *domain.Review) error {
	return database.DB.Create(review).Error
}

// UpdateReview - Cập nhật review
func UpdateReview(review *domain.Review) error {
	return database.DB.Save(review).Error
}

// FindAllReviews - Lấy tất cả reviews cho admin (phân trang)
func FindAllReviews(status, search string, offset, limit int) ([]domain.Review, int64, error) {
	query := database.DB.Model(&domain.Review{}).
		Preload("User").
		Preload("Tour")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		like := "%" + search + "%"
		query = query.Where("LOWER(content) LIKE ? OR LOWER(title) LIKE ?", like, like)
	}

	var total int64
	query.Count(&total)

	var reviews []domain.Review
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&reviews).Error

	return reviews, total, err
}

// GetTourReviewStats - Lấy thống kê rating cho tour
func GetTourReviewStats(tourID uint) (*domain.TourReviewStats, error) {
	stats := &domain.TourReviewStats{TourID: tourID}

	// Tổng reviews + average
	row := database.DB.Model(&domain.Review{}).
		Where("tour_id = ? AND status = ?", tourID, domain.ReviewStatusPublished).
		Select("COUNT(*) as total, COALESCE(AVG(rating), 0) as avg_rating").
		Row()

	var avgRating float64
	if err := row.Scan(&stats.TotalReviews, &avgRating); err != nil {
		return stats, err
	}
	stats.AverageRating = avgRating

	// Count per rating
	type RatingCount struct {
		Rating int
		Count  int
	}
	var counts []RatingCount
	database.DB.Model(&domain.Review{}).
		Where("tour_id = ? AND status = ?", tourID, domain.ReviewStatusPublished).
		Select("rating, COUNT(*) as count").
		Group("rating").
		Find(&counts)

	for _, rc := range counts {
		switch rc.Rating {
		case 5:
			stats.Rating5 = rc.Count
		case 4:
			stats.Rating4 = rc.Count
		case 3:
			stats.Rating3 = rc.Count
		case 2:
			stats.Rating2 = rc.Count
		case 1:
			stats.Rating1 = rc.Count
		}
	}

	return stats, nil
}

// UpdateTourRating - Cập nhật rating trung bình cho tour
func UpdateTourRating(tourID uint) error {
	stats, err := GetTourReviewStats(tourID)
	if err != nil {
		return err
	}

	return database.DB.Model(&domain.Tour{}).
		Where("id = ?", tourID).
		Updates(map[string]interface{}{
			"rating":       stats.AverageRating,
			"review_count": stats.TotalReviews,
		}).Error
}
