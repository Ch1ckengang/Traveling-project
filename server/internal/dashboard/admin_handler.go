package dashboard

import (
	"net/http"
	"travel-backend/internal/shared"

	"github.com/gin-gonic/gin"
)

// AdminGetDashboardSummaryHandler - GET /api/v1/admin/dashboard/summary
func AdminGetDashboardSummaryHandler(c *gin.Context) {
	summary, err := GetDashboardSummary()
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy dữ liệu thống kê", "DASHBOARD_SUMMARY_ERROR")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Thành công", summary)
}

// AdminGetRevenueChartHandler - GET /api/v1/admin/dashboard/revenue-chart
func AdminGetRevenueChartHandler(c *gin.Context) {
	data, err := GetRevenueChartData()
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy dữ liệu biểu đồ doanh thu", "DASHBOARD_REVENUE_ERROR")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Thành công", data)
}

// AdminGetTopToursHandler - GET /api/v1/admin/dashboard/top-tours
func AdminGetTopToursHandler(c *gin.Context) {
	data, err := GetTopTours()
	if err != nil {
		shared.RespondError(c, http.StatusInternalServerError, "Không thể lấy dữ liệu top tours", "DASHBOARD_TOP_TOURS_ERROR")
		return
	}

	shared.RespondSuccess(c, http.StatusOK, "Thành công", data)
}
