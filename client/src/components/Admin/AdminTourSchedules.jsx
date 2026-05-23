import { useState, useEffect } from 'react';
import { adminGetTourSchedules, adminCreateTourSchedule, adminDeleteTourSchedule } from '../../services/adminService';

const AdminTourSchedules = ({ tour, onClose }) => {
  const [schedules, setSchedules] = useState([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  
  const [form, setForm] = useState({
    departure_date: '',
    return_date: '',
    total_slots: 30,
    price_modifier: 0
  });

  useEffect(() => {
    if (tour?.id) {
      fetchSchedules();
    }
  }, [tour?.id]);

  const fetchSchedules = async () => {
    setLoading(true);
    try {
      const res = await adminGetTourSchedules(tour.id);
      setSchedules(res.data || []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!form.departure_date || !form.total_slots) {
      alert('Vui lòng điền ngày khởi hành và số chỗ');
      return;
    }

    setSaving(true);
    try {
      const data = {
        departure_date: new Date(form.departure_date).toISOString(),
        return_date: form.return_date ? new Date(form.return_date).toISOString() : new Date(form.departure_date).toISOString(),
        total_slots: parseInt(form.total_slots),
        price_modifier: parseInt(form.price_modifier) || 0
      };
      
      await adminCreateTourSchedule(tour.id, data);
      
      // Reset form
      setForm({
        departure_date: '',
        return_date: '',
        total_slots: 30,
        price_modifier: 0
      });
      
      fetchSchedules();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi tạo lịch trình');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async (scheduleId) => {
    if (!window.confirm('Bạn có chắc muốn xóa lịch trình này?')) return;
    try {
      await adminDeleteTourSchedule(tour.id, scheduleId);
      fetchSchedules();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi xóa lịch trình');
    }
  };

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-content" style={{ maxWidth: '800px' }} onClick={(e) => e.stopPropagation()}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
          <h2>Quản lý lịch trình: {tour?.name}</h2>
          <button onClick={onClose} style={{ background: 'none', border: 'none', fontSize: '24px', cursor: 'pointer' }}>×</button>
        </div>

        <div style={{ display: 'flex', gap: '20px' }}>
          {/* Form */}
          <div style={{ flex: 1 }}>
            <h3 style={{ marginBottom: '15px' }}>Thêm lịch khởi hành</h3>
            <form onSubmit={handleSubmit}>
              <div className="form-group">
                <label>Ngày khởi hành *</label>
                <input 
                  type="date" 
                  required 
                  value={form.departure_date} 
                  onChange={(e) => setForm({...form, departure_date: e.target.value})} 
                />
              </div>
              <div className="form-group">
                <label>Ngày về (Tùy chọn)</label>
                <input 
                  type="date" 
                  value={form.return_date} 
                  onChange={(e) => setForm({...form, return_date: e.target.value})} 
                />
              </div>
              <div className="form-group">
                <label>Tổng số chỗ *</label>
                <input 
                  type="number" 
                  required 
                  min="1"
                  value={form.total_slots} 
                  onChange={(e) => setForm({...form, total_slots: e.target.value})} 
                />
              </div>
              <div className="form-group">
                <label>Phụ thu/Giảm giá (VND)</label>
                <input 
                  type="number" 
                  value={form.price_modifier} 
                  onChange={(e) => setForm({...form, price_modifier: e.target.value})} 
                  placeholder="Ví dụ: 500000 để phụ thu lễ"
                />
              </div>
              <button type="submit" className="btn-save" style={{ width: '100%', marginTop: '10px' }} disabled={saving}>
                {saving ? 'Đang thêm...' : 'Thêm lịch trình'}
              </button>
            </form>
          </div>

          {/* List */}
          <div style={{ flex: 2 }}>
            <h3 style={{ marginBottom: '15px' }}>Danh sách lịch trình</h3>
            {loading ? (
              <p>Đang tải...</p>
            ) : schedules.length === 0 ? (
              <p>Chưa có lịch trình nào.</p>
            ) : (
              <table className="admin-table">
                <thead>
                  <tr>
                    <th>Khởi hành</th>
                    <th>Còn lại / Tổng</th>
                    <th>Giá điều chỉnh</th>
                    <th>Thao tác</th>
                  </tr>
                </thead>
                <tbody>
                  {schedules.map(s => (
                    <tr key={s.id}>
                      <td>{new Date(s.departure_date).toLocaleDateString('vi-VN')}</td>
                      <td>{s.remaining_slots} / {s.total_slots}</td>
                      <td style={{ color: s.price_modifier > 0 ? 'red' : s.price_modifier < 0 ? 'green' : 'inherit' }}>
                        {s.price_modifier !== 0 ? s.price_modifier.toLocaleString('vi-VN') + 'đ' : '-'}
                      </td>
                      <td>
                        <button className="btn-sm btn-delete" onClick={() => handleDelete(s.id)}>Xóa</button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminTourSchedules;
