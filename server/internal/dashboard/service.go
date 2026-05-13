package dashboard

import (
	"time"
	"travel-backend/database"
	"travel-backend/domain"
)

// SummaryData contains high-level metrics for the dashboard
type SummaryData struct {
	TotalUsers     int64 `json:"total_users"`
	TotalTours     int64 `json:"total_tours"`
	TotalBookings  int64 `json:"total_bookings"`
	TotalRevenue   int64 `json:"total_revenue"`
}

// GetDashboardSummary retrieves overall counts and revenue
func GetDashboardSummary() (*SummaryData, error) {
	var summary SummaryData

	// Count Users
	if err := database.DB.Model(&domain.User{}).Count(&summary.TotalUsers).Error; err != nil {
		return nil, err
	}

	// Count Tours
	if err := database.DB.Model(&domain.Tour{}).Count(&summary.TotalTours).Error; err != nil {
		return nil, err
	}

	// Count Bookings (Only completed or confirmed for revenue relevance)
	if err := database.DB.Model(&domain.Booking{}).Where("status IN ?", []string{"completed", "confirmed"}).Count(&summary.TotalBookings).Error; err != nil {
		return nil, err
	}

	// Calculate Total Revenue from paid bookings
	type Result struct {
		Total int64
	}
	var res Result
	if err := database.DB.Model(&domain.Booking{}).Select("COALESCE(SUM(total_amount), 0) as total").Where("payment_status = ?", "paid").Scan(&res).Error; err != nil {
		return nil, err
	}
	summary.TotalRevenue = res.Total

	return &summary, nil
}

// RevenueData represents revenue for a specific period
type RevenueData struct {
	Date    string `json:"date"`
	Revenue int64  `json:"revenue"`
}

// GetRevenueChartData retrieves revenue grouped by date for the last 30 days
func GetRevenueChartData() ([]RevenueData, error) {
	var data []RevenueData

	// 30 days ago
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	// PostgreSQL specific query to group by date
	query := `
		SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as date, COALESCE(SUM(total_amount), 0) as revenue
		FROM bookings
		WHERE payment_status = 'paid' AND created_at >= ?
		GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
		ORDER BY date ASC
	`

	if err := database.DB.Raw(query, thirtyDaysAgo).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

// TopTourData represents a tour with its booking count
type TopTourData struct {
	TourID       uint   `json:"tour_id"`
	TourName     string `json:"tour_name"`
	BookingCount int64  `json:"booking_count"`
	Revenue      int64  `json:"revenue"`
}

// GetTopTours retrieves top 5 tours by booking count
func GetTopTours() ([]TopTourData, error) {
	var data []TopTourData

	query := `
		SELECT t.id as tour_id, t.name as tour_name, COUNT(b.id) as booking_count, COALESCE(SUM(b.total_amount), 0) as revenue
		FROM tours t
		LEFT JOIN bookings b ON t.id = b.tour_id AND b.payment_status = 'paid'
		GROUP BY t.id, t.name
		ORDER BY booking_count DESC
		LIMIT 5
	`

	if err := database.DB.Raw(query).Scan(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}
