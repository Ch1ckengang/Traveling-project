import { useEffect, useState } from 'react';
import { getDashboardSummary, getRevenueChart, getTopTours } from '../../services/dashboardService';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import './AdminPage.css';

const AdminDashboardPage = () => {
  const [summary, setSummary] = useState({
    total_users: 0,
    total_tours: 0,
    total_bookings: 0,
    total_revenue: 0,
  });
  const [chartData, setChartData] = useState([]);
  const [topTours, setTopTours] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [summaryRes, chartRes, topToursRes] = await Promise.all([
          getDashboardSummary(),
          getRevenueChart(),
          getTopTours()
        ]);
        
        if (summaryRes.success) setSummary(summaryRes.data);
        if (chartRes.success) {
          // Format date for chart
          const formattedData = (chartRes.data || []).map(item => ({
            ...item,
            shortDate: new Date(item.date).toLocaleDateString('vi-VN', { day: '2-digit', month: '2-digit' })
          }));
          setChartData(formattedData);
        }
        if (topToursRes.success) setTopTours(topToursRes.data || []);
      } catch (err) {
        console.error('Error fetching dashboard data:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  const formatCurrency = (val) => val.toLocaleString('vi-VN') + 'đ';

  if (loading) return <div className="admin-loading">Đang tải dữ liệu tổng quan...</div>;

  return (
    <div className="admin-page dashboard-page">
      <div className="admin-header">
        <h1>Dashboard Thống Kê</h1>
      </div>

      <div className="summary-cards">
        <div className="summary-card">
          <div className="card-icon" style={{ background: '#e0f2fe', color: '#0284c7' }}>👥</div>
          <div className="card-content">
            <p>Tổng Người Dùng</p>
            <h3>{summary.total_users.toLocaleString()}</h3>
          </div>
        </div>
        
        <div className="summary-card">
          <div className="card-icon" style={{ background: '#fef3c7', color: '#d97706' }}>🗺️</div>
          <div className="card-content">
            <p>Tổng Tour</p>
            <h3>{summary.total_tours.toLocaleString()}</h3>
          </div>
        </div>
        
        <div className="summary-card">
          <div className="card-icon" style={{ background: '#dcfce7', color: '#16a34a' }}>📋</div>
          <div className="card-content">
            <p>Số Bookings (Thành công)</p>
            <h3>{summary.total_bookings.toLocaleString()}</h3>
          </div>
        </div>
        
        <div className="summary-card">
          <div className="card-icon" style={{ background: '#fce7f3', color: '#db2777' }}>💰</div>
          <div className="card-content">
            <p>Tổng Doanh Thu</p>
            <h3>{formatCurrency(summary.total_revenue)}</h3>
          </div>
        </div>
      </div>

      <div className="dashboard-grid">
        <div className="dashboard-section chart-section">
          <h3>Doanh thu 30 ngày qua</h3>
          <div className="chart-container" style={{ height: '350px', marginTop: '1rem' }}>
            {chartData.length > 0 ? (
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartData} margin={{ top: 5, right: 20, left: 20, bottom: 5 }}>
                  <CartesianGrid strokeDasharray="3 3" vertical={false} stroke="#eee" />
                  <XAxis dataKey="shortDate" axisLine={false} tickLine={false} tick={{fill: '#64748b'}} />
                  <YAxis 
                    tickFormatter={(val) => val >= 1000000 ? `${(val / 1000000).toFixed(1)}M` : val} 
                    axisLine={false} 
                    tickLine={false} 
                    tick={{fill: '#64748b'}} 
                  />
                  <Tooltip 
                    formatter={(value) => [formatCurrency(value), 'Doanh thu']}
                    labelStyle={{ color: '#334155', fontWeight: 'bold' }}
                    contentStyle={{ borderRadius: '8px', border: 'none', boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1)' }}
                  />
                  <Line type="monotone" dataKey="revenue" stroke="#3b82f6" strokeWidth={3} dot={{ r: 4, fill: '#3b82f6', strokeWidth: 2, stroke: '#fff' }} activeDot={{ r: 6 }} />
                </LineChart>
              </ResponsiveContainer>
            ) : (
              <div className="empty-chart">Chưa có dữ liệu doanh thu trong 30 ngày qua</div>
            )}
          </div>
        </div>

        <div className="dashboard-section top-tours-section">
          <h3>Top Tour Yêu Thích</h3>
          <div className="top-tours-list" style={{ marginTop: '1rem' }}>
            {topTours.length > 0 ? (
              topTours.map((tour, index) => (
                <div key={tour.tour_id} className="top-tour-item" style={{ display: 'flex', alignItems: 'center', padding: '1rem', borderBottom: '1px solid #f1f5f9' }}>
                  <div className="rank" style={{ width: '30px', height: '30px', borderRadius: '50%', backgroundColor: index < 3 ? '#fef08a' : '#f1f5f9', display: 'flex', alignItems: 'center', justifyContent: 'center', fontWeight: 'bold', marginRight: '1rem' }}>
                    {index + 1}
                  </div>
                  <div className="tour-info" style={{ flex: 1 }}>
                    <p style={{ fontWeight: 600, margin: 0, color: '#1e293b' }}>{tour.tour_name}</p>
                    <p style={{ margin: 0, fontSize: '0.85rem', color: '#64748b' }}>{tour.booking_count} lượt đặt</p>
                  </div>
                  <div className="tour-revenue" style={{ fontWeight: 600, color: '#059669' }}>
                    {formatCurrency(tour.revenue)}
                  </div>
                </div>
              ))
            ) : (
              <div className="empty-list">Chưa có dữ liệu tour</div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminDashboardPage;
