import { useState } from 'react';
import axiosInstance from '../../utils/axiosInstance';
import './AdminPage.css';

const AdminReportsPage = () => {
  const [downloading, setDownloading] = useState(false);

  const downloadCSV = async () => {
    try {
      setDownloading(true);
      const res = await axiosInstance.get('/admin/bookings', { params: { limit: 1000 } });
      const bookings = res.data?.data || [];
      
      if (bookings.length === 0) {
        alert('Không có dữ liệu để xuất');
        return;
      }

      // Convert JSON to CSV
      const headers = ['Mã Đặt Chỗ', 'Ngày Đặt', 'Tên Khách', 'Email', 'Điện Thoại', 'Tên Tour', 'Ngày Khởi Hành', 'Số Người', 'Tổng Tiền', 'Trạng Thái', 'Thanh Toán'];
      const csvRows = [];
      csvRows.push(headers.join(','));

      for (const b of bookings) {
        const row = [
          b.booking_code,
          new Date(b.created_at).toLocaleDateString('vi-VN'),
          `"${b.full_name}"`, // Quote strings that might contain commas
          b.email,
          `"${b.phone}"`,
          `"${b.tour?.name || b.tour_name || 'N/A'}"`,
          new Date(b.travel_date).toLocaleDateString('vi-VN'),
          (b.adult_count + b.child_count + b.infant_count),
          b.total_amount,
          b.status,
          b.payment_status
        ];
        csvRows.push(row.join(','));
      }

      const csvString = csvRows.join('\n');
      
      // Create a Blob and download it
      const blob = new Blob([new Uint8Array([0xEF, 0xBB, 0xBF]), csvString], { type: 'text/csv;charset=utf-8;' }); // Added BOM for Excel UTF-8 support
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', `bookings_report_${new Date().getTime()}.csv`);
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } catch (err) {
      alert('Lỗi xuất file CSV: ' + (err.response?.data?.message || err.message));
    } finally {
      setDownloading(false);
    }
  };

  return (
    <div className="admin-page">
      <div className="admin-header">
        <h1>Báo Cáo & Xuất Dữ Liệu</h1>
      </div>

      <div className="dashboard-section" style={{ maxWidth: '600px' }}>
        <h3>Xuất Dữ Liệu Đặt Chỗ</h3>
        <p style={{ color: '#64748b', marginBottom: '1.5rem', marginTop: '0.5rem' }}>
          Tải xuống toàn bộ danh sách đặt chỗ dưới định dạng CSV (có thể mở bằng Excel).
        </p>
        
        <button 
          onClick={downloadCSV} 
          disabled={downloading}
          className="btn-add" 
          style={{ width: '100%', padding: '1rem', fontSize: '1rem' }}
        >
          {downloading ? 'Đang chuẩn bị dữ liệu...' : '📥 Tải Xuống Báo Cáo CSV'}
        </button>
      </div>
    </div>
  );
};

export default AdminReportsPage;
