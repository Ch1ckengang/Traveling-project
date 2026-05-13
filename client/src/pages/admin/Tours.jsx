import { useEffect, useState, useCallback } from 'react';
import { adminGetTours, adminCreateTour, adminUpdateTour, adminDeleteTour, adminToggleTour } from '../../services/adminService';
import './AdminPage.css';

const AdminToursPage = () => {
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [meta, setMeta] = useState({});
  const [search, setSearch] = useState('');
  const [typeFilter, setTypeFilter] = useState('');
  const [showModal, setShowModal] = useState(false);
  const [editTour, setEditTour] = useState(null);
  const [saving, setSaving] = useState(false);

  const [form, setForm] = useState({
    name: '', type: 'domestic', price_amount: '', description: '',
    location: '', country: 'Việt Nam', duration: '', departure_date: '',
    itinerary: '', services: '', image_url: '',
  });

  const fetchTours = useCallback(async () => {
    setLoading(true);
    try {
      const res = await adminGetTours({ page, limit: 10, type: typeFilter, search });
      setTours(res.data || []);
      setMeta(res.meta || {});
    } catch {
      setTours([]);
    }
    setLoading(false);
  }, [page, search, typeFilter]);

  useEffect(() => { fetchTours(); }, [fetchTours]);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setSaving(true);
    try {
      const data = { ...form, price_amount: parseInt(form.price_amount) || 0 };
      if (editTour) {
        await adminUpdateTour(editTour.id, data);
      } else {
        await adminCreateTour(data);
      }
      setShowModal(false);
      setEditTour(null);
      resetForm();
      fetchTours();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi khi lưu tour');
    }
    setSaving(false);
  };

  const handleEdit = (tour) => {
    setEditTour(tour);
    setForm({
      name: tour.name, type: tour.type, price_amount: tour.price_amount?.toString() || '',
      description: tour.description || '', location: tour.location || '',
      country: tour.country || 'Việt Nam', duration: tour.duration || '',
      departure_date: tour.departure_date || '', itinerary: tour.itinerary || '',
      services: tour.services || '', image_url: tour.image_url || '',
    });
    setShowModal(true);
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Bạn có chắc muốn ẩn tour này?')) return;
    try {
      await adminDeleteTour(id);
      fetchTours();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi');
    }
  };

  const handleToggle = async (id) => {
    try {
      await adminToggleTour(id);
      fetchTours();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi');
    }
  };

  const resetForm = () => {
    setForm({
      name: '', type: 'domestic', price_amount: '', description: '',
      location: '', country: 'Việt Nam', duration: '', departure_date: '',
      itinerary: '', services: '', image_url: '',
    });
  };

  const formatPrice = (amount) => {
    if (!amount) return '—';
    return amount.toLocaleString('vi-VN') + 'đ';
  };

  return (
    <div className="admin-page">
      <div className="admin-header">
        <h1>Quản lý Tour</h1>
        <button className="btn-add" onClick={() => { resetForm(); setEditTour(null); setShowModal(true); }}>
          + Tạo tour mới
        </button>
      </div>

      <div className="admin-filters">
        <input
          type="text"
          placeholder="Tìm kiếm theo tên, địa điểm..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="filter-input"
        />
        <select value={typeFilter} onChange={(e) => { setTypeFilter(e.target.value); setPage(1); }} className="filter-select">
          <option value="">Tất cả loại</option>
          <option value="domestic">Trong nước</option>
          <option value="international">Quốc tế</option>
          <option value="service">Dịch vụ</option>
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
                  <th>ID</th>
                  <th>Tên tour</th>
                  <th>Loại</th>
                  <th>Giá</th>
                  <th>Địa điểm</th>
                  <th>Trạng thái</th>
                  <th>Slots</th>
                  <th>Thao tác</th>
                </tr>
              </thead>
              <tbody>
                {tours.length === 0 ? (
                  <tr><td colSpan="8" className="empty-row">Không có tour nào</td></tr>
                ) : tours.map(tour => (
                  <tr key={tour.id}>
                    <td>{tour.id}</td>
                    <td className="tour-name">{tour.name}</td>
                    <td><span className={`badge badge-${tour.type}`}>{tour.type}</span></td>
                    <td>{formatPrice(tour.price_amount)}</td>
                    <td>{tour.location}</td>
                    <td>
                      <span className={`status-dot ${tour.is_active ? 'active' : 'inactive'}`}>
                        {tour.is_active ? 'Hoạt động' : 'Ẩn'}
                      </span>
                    </td>
                    <td>{tour.remaining_slots}</td>
                    <td className="actions">
                      <button className="btn-sm btn-edit" onClick={() => handleEdit(tour)}>Sửa</button>
                      <button className="btn-sm btn-toggle" onClick={() => handleToggle(tour.id)}>
                        {tour.is_active ? 'Ẩn' : 'Hiện'}
                      </button>
                      {!tour.is_active && (
                        <button className="btn-sm btn-delete" onClick={() => handleDelete(tour.id)}>Xóa</button>
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
              <span>Trang {page} / {meta.total_pages} ({meta.total} tours)</span>
              <button disabled={page >= meta.total_pages} onClick={() => setPage(p => p + 1)}>Sau →</button>
            </div>
          )}
        </>
      )}

      {/* Modal Create/Edit */}
      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>{editTour ? 'Cập nhật tour' : 'Tạo tour mới'}</h2>
            <form onSubmit={handleSubmit}>
              <div className="form-grid">
                <div className="form-group">
                  <label>Tên tour *</label>
                  <input required value={form.name} onChange={(e) => setForm({...form, name: e.target.value})} />
                </div>
                <div className="form-group">
                  <label>Loại *</label>
                  <select value={form.type} onChange={(e) => setForm({...form, type: e.target.value})}>
                    <option value="domestic">Trong nước</option>
                    <option value="international">Quốc tế</option>
                    <option value="service">Dịch vụ</option>
                  </select>
                </div>
                <div className="form-group">
                  <label>Giá (VND) *</label>
                  <input type="number" required value={form.price_amount} onChange={(e) => setForm({...form, price_amount: e.target.value})} />
                </div>
                <div className="form-group">
                  <label>Địa điểm</label>
                  <input value={form.location} onChange={(e) => setForm({...form, location: e.target.value})} />
                </div>
                <div className="form-group">
                  <label>Quốc gia</label>
                  <input value={form.country} onChange={(e) => setForm({...form, country: e.target.value})} />
                </div>
                <div className="form-group">
                  <label>Thời gian</label>
                  <input value={form.duration} onChange={(e) => setForm({...form, duration: e.target.value})} placeholder="VD: 3 ngày 2 đêm" />
                </div>
                <div className="form-group full-width">
                  <label>Mô tả</label>
                  <textarea rows="3" value={form.description} onChange={(e) => setForm({...form, description: e.target.value})} />
                </div>
                <div className="form-group full-width">
                  <label>Lịch trình</label>
                  <textarea rows="3" value={form.itinerary} onChange={(e) => setForm({...form, itinerary: e.target.value})} />
                </div>
                <div className="form-group full-width">
                  <label>Dịch vụ bao gồm</label>
                  <textarea rows="2" value={form.services} onChange={(e) => setForm({...form, services: e.target.value})} />
                </div>
                <div className="form-group full-width">
                  <label>URL hình ảnh</label>
                  <input value={form.image_url} onChange={(e) => setForm({...form, image_url: e.target.value})} placeholder="https://..." />
                </div>
              </div>
              <div className="modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setShowModal(false)}>Hủy</button>
                <button type="submit" className="btn-save" disabled={saving}>
                  {saving ? 'Đang lưu...' : (editTour ? 'Cập nhật' : 'Tạo tour')}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminToursPage;
