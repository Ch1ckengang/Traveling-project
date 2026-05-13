import { useEffect, useState, useCallback } from 'react';
import { adminGetBookings, adminGetBookingStats, adminConfirmBooking, adminCancelBooking } from '../../services/adminService';
import './AdminPage.css';

const AdminBookingsPage = () => {
  const [bookings, setBookings] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [meta, setMeta] = useState({});
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');

  const fetchData = useCallback(async () => {
    setLoading(true);
    try {
      const [bookingsRes, statsRes] = await Promise.all([
        adminGetBookings({ page, limit: 10, status: statusFilter, search }),
        adminGetBookingStats(),
      ]);
      setBookings(bookingsRes.data || []);
      setMeta(bookingsRes.meta || {});
      if (statsRes.success) setStats(statsRes.data);
    } catch {
      setBookings([]);
    }
    setLoading(false);
  }, [page, search, statusFilter]);

  useEffect(() => { fetchData(); }, [fetchData]);

  const handleConfirm = async (code) => {
    if (!window.confirm(`Xác nhận booking ${code}?`)) return;
    try {
      await adminConfirmBooking(code);
      fetchData();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi');
    }
  };

  const handleCancel = async (code) => {
    const reason = prompt('Lý do hủy booking:');
    if (reason === null) return;
    try {
      await adminCancelBooking(code, reason);
      fetchData();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi');
    }
  };

  const formatDate = (d) => {
    if (!d) return '—';
    return new Date(d).toLocaleDateString('vi-VN');
  };

  const formatPrice = (amount) => {
    if (!amount) return '—';
    return amount.toLocaleString('vi-VN') + 'đ';
  };

  const statusLabels = {
    pending: 'Chờ xử lý',
    booked: 'Đã đặt',
    confirmed: 'Đã xác nhận',
    cancelled: 'Đã hủy',
    completed: 'Hoàn thành',
  };

  return (
    <div className="admin-page">
      <div className="admin-header">
        <h1>Quản lý Booking</h1>
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-number">{stats.total_bookings}</div>
            <div className="stat-label">Tổng booking</div>
          </div>
          <div className="stat-card stat-revenue">
            <div className="stat-number">{formatPrice(stats.total_revenue)}</div>
            <div className="stat-label">Doanh thu</div>
          </div>
          {stats.stats?.map(s => (
            <div className="stat-card" key={s.status}>
              <div className="stat-number">{s.count}</div>
              <div className="stat-label">{statusLabels[s.status] || s.status}</div>
            </div>
          ))}
        </div>
      )}

      <div className="admin-filters">
        <input
          type="text"
          placeholder="Tìm mã booking, tên, SĐT..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="filter-input"
        />
        <select value={statusFilter} onChange={(e) => { setStatusFilter(e.target.value); setPage(1); }} className="filter-select">
          <option value="">Tất cả trạng thái</option>
          <option value="pending">Chờ xử lý</option>
          <option value="booked">Đã đặt</option>
          <option value="confirmed">Đã xác nhận</option>
          <option value="cancelled">Đã hủy</option>
        </select>
      </div>

      {loading ? (
        <div className="admin-loading">Đang tải...</div>
      ) : (
        <>
          <div className="admin-table-wrapper">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Mã booking</th>
                  <th>Khách hàng</th>
                  <th>Tour</th>
                  <th>Ngày đi</th>
                  <th>Tổng tiền</th>
                  <th>Trạng thái</th>
                  <th>Thanh toán</th>
                  <th>Thao tác</th>
                </tr>
              </thead>
              <tbody>
                {bookings.length === 0 ? (
                  <tr><td colSpan="8" className="empty-row">Không có booking nào</td></tr>
                ) : bookings.map(b => (
                  <tr key={b.id}>
                    <td className="code-cell">{b.booking_code}</td>
                    <td>
                      <div>{b.full_name}</div>
                      <div className="text-muted">{b.phone}</div>
                    </td>
                    <td className="tour-name">{b.tour?.name || '—'}</td>
                    <td>{formatDate(b.travel_date)}</td>
                    <td>{formatPrice(b.total_amount)}</td>
                    <td>
                      <span className={`badge badge-${b.status}`}>
                        {statusLabels[b.status] || b.status}
                      </span>
                    </td>
                    <td>
                      <span className={`badge badge-pay-${b.payment_status}`}>
                        {b.payment_status}
                      </span>
                    </td>
                    <td className="actions">
                      {(b.status === 'pending' || b.status === 'booked') && (
                        <button className="btn-sm btn-confirm" onClick={() => handleConfirm(b.booking_code)}>
                          Xác nhận
                        </button>
                      )}
                      {b.status !== 'cancelled' && (
                        <button className="btn-sm btn-cancel-action" onClick={() => handleCancel(b.booking_code)}>
                          Hủy
                        </button>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {meta.total_pages > 1 && (
            <div className="admin-pagination">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>← Trước</button>
              <span>Trang {page} / {meta.total_pages} ({meta.total} bookings)</span>
              <button disabled={page >= meta.total_pages} onClick={() => setPage(p => p + 1)}>Sau →</button>
            </div>
          )}
        </>
      )}
    </div>
  );
};

export default AdminBookingsPage;
