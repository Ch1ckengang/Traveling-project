package review

import (
	"fmt"
	"travel-backend/database"
	"travel-backend/domain"

	"gorm.io/gorm"
)

// GetTourReviews - Lấy danh sách reviews của tour (public)
func GetTourReviews(tourID uint, rating, offset, limit int) ([]domain.Review, int64, error) {
	return FindReviewsByTourID(tourID, rating, offset, limit)
}

// GetTourStats - Lấy thống kê rating của tour
func GetTourStats(tourID uint) (*domain.TourReviewStats, error) {
	return GetTourReviewStats(tourID)
}

// SubmitReview - Customer tạo review mới
func SubmitReview(userID uint, req domain.CreateReviewRequest) (*domain.Review, error) {
	// 1. Validate request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// 2. Kiểm tra booking tồn tại và thuộc về user
	var booking domain.Booking
	if err := database.DB.First(&booking, req.BookingID).Error; err != nil {
		return nil, fmt.Errorf("không tìm thấy booking")
	}
	if booking.UserID != userID {
		return nil, fmt.Errorf("bạn không có quyền review booking này")
	}

	// 3. Kiểm tra booking đã confirmed/completed (đã thanh toán hoặc đã xác nhận)
	if booking.Status != "confirmed" && booking.Status != "completed" {
		return nil, fmt.Errorf("chỉ có thể review sau khi booking đã được xác nhận")
	}

	// 4. Kiểm tra tour_id khớp với booking
	if booking.TourID != req.TourID {
		return nil, fmt.Errorf("tour_id không khớp với booking")
	}

	// 5. Kiểm tra đã review booking này chưa
	_, err := FindReviewByBookingID(req.BookingID, userID)
	if err == nil {
		return nil, fmt.Errorf("bạn đã đánh giá booking này rồi")
	}
	if err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("lỗi kiểm tra review: %w", err)
	}

	// 6. Tạo review
	review := &domain.Review{
		TourID:    req.TourID,
		UserID:    userID,
		BookingID: req.BookingID,
		Rating:    req.Rating,
		Title:     req.Title,
		Content:   req.Content,
		Status:    domain.ReviewStatusPublished,
	}

	if err := CreateReview(review); err != nil {
		return nil, fmt.Errorf("không thể tạo review: %w", err)
	}

	// 7. Cập nhật tour rating
	UpdateTourRating(req.TourID)

	// Reload với user info
	review, _ = FindReviewByID(review.ID)
	return review, nil
}

// EditReview - Customer sửa review (trong 7 ngày)
func EditReview(reviewID, userID uint, req domain.UpdateReviewRequest) (*domain.Review, error) {
	review, err := FindReviewByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy review")
	}

	// Kiểm tra ownership
	if review.UserID != userID {
		return nil, fmt.Errorf("bạn không có quyền sửa review này")
	}

	// Kiểm tra thời hạn 7 ngày
	if !review.CanEdit() {
		return nil, fmt.Errorf("đã quá thời hạn 7 ngày để sửa review")
	}

	// Cập nhật fields
	if req.Rating != nil {
		review.Rating = *req.Rating
	}
	if req.Title != nil {
		review.Title = *req.Title
	}
	if req.Content != nil {
		review.Content = *req.Content
	}

	if err := UpdateReview(review); err != nil {
		return nil, fmt.Errorf("không thể cập nhật review: %w", err)
	}

	// Cập nhật tour rating
	UpdateTourRating(review.TourID)

	return review, nil
}

// AdminPublishReview - Admin publish review
func AdminPublishReview(reviewID uint) (*domain.Review, error) {
	review, err := FindReviewByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy review")
	}

	review.Status = domain.ReviewStatusPublished
	if err := UpdateReview(review); err != nil {
		return nil, fmt.Errorf("không thể publish review: %w", err)
	}

	UpdateTourRating(review.TourID)
	return review, nil
}

// AdminHideReview - Admin ẩn review
func AdminHideReview(reviewID uint) (*domain.Review, error) {
	review, err := FindReviewByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy review")
	}

	review.Status = domain.ReviewStatusHidden
	if err := UpdateReview(review); err != nil {
		return nil, fmt.Errorf("không thể ẩn review: %w", err)
	}

	UpdateTourRating(review.TourID)
	return review, nil
}

// AdminReplyReview - Admin phản hồi review
func AdminReplyReview(reviewID, adminUserID uint, reply string) (*domain.Review, error) {
	review, err := FindReviewByID(reviewID)
	if err != nil {
		return nil, fmt.Errorf("không tìm thấy review")
	}

	review.AdminReply = &reply
	now := review.UpdatedAt // use current time
	review.RepliedAt = &now
	review.RepliedBy = &adminUserID

	if err := UpdateReview(review); err != nil {
		return nil, fmt.Errorf("không thể lưu phản hồi: %w", err)
	}

	return review, nil
}
