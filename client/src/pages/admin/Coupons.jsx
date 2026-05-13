import { useEffect, useState, useCallback } from 'react';
import { adminGetCoupons, adminCreateCoupon, adminUpdateCoupon, adminDeleteCoupon } from '../../services/couponService';
import './AdminPage.css';

const AdminCouponsPage = () => {
  const [coupons, setCoupons] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [meta, setMeta] = useState({});
  const [search, setSearch] = useState('');
  
  const [showModal, setShowModal] = useState(false);
  const [editingId, setEditingId] = useState(null);
  const [formData, setFormData] = useState({
    code: '',
    discount_type: 'percent',
    discount_value: '',
    min_order_value: 0,
    max_discount: 0,
    start_date: '',
    end_date: '',
    usage_limit: 0,
    is_active: true
  });

  const fetchCoupons = useCallback(async () => {
    setLoading(true);
    try {
      const res = await adminGetCoupons({ page, limit: 10, search });
      setCoupons(res.data || []);
      setMeta(res.meta || {});
    } catch {
      setCoupons([]);
    }
    setLoading(false);
  }, [page, search]);

  useEffect(() => { fetchCoupons(); }, [fetchCoupons]);

  const handleOpenModal = (coupon = null) => {
    if (coupon) {
      setEditingId(coupon.id);
      setFormData({
        code: coupon.code,
        discount_type: coupon.discount_type,
        discount_value: coupon.discount_value,
        min_order_value: coupon.min_order_value,
        max_discount: coupon.max_discount,
        start_date: new Date(coupon.start_date).toISOString().split('T')[0],
        end_date: new Date(coupon.end_date).toISOString().split('T')[0],
        usage_limit: coupon.usage_limit,
        is_active: coupon.is_active
      });
    } else {
      setEditingId(null);
      setFormData({
        code: '',
        discount_type: 'percent',
        discount_value: '',
        min_order_value: 0,
        max_discount: 0,
        start_date: new Date().toISOString().split('T')[0],
        end_date: '',
        usage_limit: 0,
        is_active: true
      });
    }
    setShowModal(true);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const payload = {
        ...formData,
        code: formData.code.toUpperCase(),
        discount_value: Number(formData.discount_value),
        min_order_value: Number(formData.min_order_value),
        max_discount: Number(formData.max_discount),
        usage_limit: Number(formData.usage_limit),
        start_date: new Date(formData.start_date).toISOString(),
        end_date: new Date(formData.end_date).toISOString()
      };

      if (editingId) {
        await adminUpdateCoupon(editingId, payload);
      } else {
        await adminCreateCoupon(payload);
      }
      setShowModal(false);
      fetchCoupons();
    } catch (err) {
      alert(err.response?.data?.message || 'Có lỗi xảy ra');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Bạn có chắc muốn xóa mã giảm giá này?')) return;
    try {
      await adminDeleteCoupon(id);
      fetchCoupons();
    } catch (err) {
      alert(err.response?.data?.message || 'Có lỗi xảy ra');
    }
  };

  const formatPrice = (amount) => amount ? amount.toLocaleString('vi-VN') + 'đ' : '0đ';
  const formatDate = (d) => d ? new Date(d).toLocaleDateString('vi-VN') : '—';

  return (
    <div className="admin-page">
      <div className="admin-header">
        <h1>Quản lý Khuyến mãi (Coupons)</h1>
        <button className="btn-add" onClick={() => handleOpenModal()}>+ Thêm Mã Mới</button>
      </div>

      <div className="admin-filters">
        <input
          type="text"
          placeholder="Tìm mã giảm giá..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="filter-input"
        />
      </div>

      {loading ? (
        <div className="admin-loading">Đang tải...</div>
      ) : (
        <>
          <div className="admin-table-wrapper">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>Mã</th>
                  <th>Loại</th>
                  <th>Giá trị</th>
                  <th>Hạn sử dụng</th>
                  <th>Đã dùng / Giới hạn</th>
                  <th>Trạng thái</th>
                  <th>Thao tác</th>
                </tr>
              </thead>
              <tbody>
                {coupons.length === 0 ? (
                  <tr><td colSpan="7" className="empty-row">Không có mã giảm giá nào</td></tr>
                ) : coupons.map(c => (
                  <tr key={c.id}>
                    <td className="code-cell">{c.code}</td>
                    <td>
                      <span className={`badge ${c.discount_type === 'percent' ? 'badge-international' : 'badge-domestic'}`}>
                        {c.discount_type === 'percent' ? 'Phần trăm' : 'Cố định'}
                      </span>
                    </td>
                    <td style={{ color: '#059669', fontWeight: 600 }}>
                      {c.discount_type === 'percent' ? `${c.discount_value}%` : formatPrice(c.discount_value)}
                    </td>
                    <td>{formatDate(c.start_date)} - {formatDate(c.end_date)}</td>
                    <td>{c.used_count} / {c.usage_limit === 0 ? '∞' : c.usage_limit}</td>
                    <td>
                      <span className={`status-dot ${c.is_active ? 'active' : 'inactive'}`}>
                        {c.is_active ? 'Hoạt động' : 'Tạm khóa'}
                      </span>
                    </td>
                    <td className="actions">
                      <button className="btn-sm btn-edit" onClick={() => handleOpenModal(c)}>Sửa</button>
                      <button className="btn-sm btn-delete" onClick={() => handleDelete(c.id)}>Xóa</button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {meta.total_pages > 1 && (
            <div className="admin-pagination">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>← Trước</button>
              <span>Trang {page} / {meta.total_pages}</span>
              <button disabled={page >= meta.total_pages} onClick={() => setPage(p => p + 1)}>Sau →</button>
            </div>
          )}
        </>
      )}

      {/* Modal */}
      {showModal && (
        <div className="modal-overlay" onClick={() => setShowModal(false)}>
          <div className="modal-content" onClick={e => e.stopPropagation()}>
            <h2>{editingId ? 'Sửa Mã Giảm Giá' : 'Thêm Mã Mới'}</h2>
            <form onSubmit={handleSubmit} className="form-grid">
              
              <div className="form-group">
                <label>Mã nhập (Code)</label>
                <input required value={formData.code} onChange={e => setFormData({...formData, code: e.target.value.toUpperCase()})} placeholder="VD: SUMMER26" />
              </div>

              <div className="form-group">
                <label>Loại giảm giá</label>
                <select value={formData.discount_type} onChange={e => setFormData({...formData, discount_type: e.target.value})}>
                  <option value="percent">Theo phần trăm (%)</option>
                  <option value="fixed">Số tiền cố định (VNĐ)</option>
                </select>
              </div>

              <div className="form-group">
                <label>Giá trị giảm ({formData.discount_type === 'percent' ? '%' : 'VNĐ'})</label>
                <input required type="number" min="1" value={formData.discount_value} onChange={e => setFormData({...formData, discount_value: e.target.value})} />
              </div>

              <div className="form-group">
                <label>Giảm tối đa (VNĐ) {formData.discount_type !== 'percent' && '(Chỉ dùng cho %)'}</label>
                <input type="number" min="0" value={formData.max_discount} onChange={e => setFormData({...formData, max_discount: e.target.value})} disabled={formData.discount_type !== 'percent'} />
              </div>

              <div className="form-group">
                <label>Đơn tối thiểu (VNĐ)</label>
                <input type="number" min="0" value={formData.min_order_value} onChange={e => setFormData({...formData, min_order_value: e.target.value})} />
              </div>

              <div className="form-group">
                <label>Giới hạn lượt dùng (0 = ∞)</label>
                <input type="number" min="0" value={formData.usage_limit} onChange={e => setFormData({...formData, usage_limit: e.target.value})} />
              </div>

              <div className="form-group">
                <label>Ngày bắt đầu</label>
                <input required type="date" value={formData.start_date} onChange={e => setFormData({...formData, start_date: e.target.value})} />
              </div>

              <div className="form-group">
                <label>Ngày kết thúc</label>
                <input required type="date" value={formData.end_date} onChange={e => setFormData({...formData, end_date: e.target.value})} />
              </div>

              <div className="form-group full-width" style={{ flexDirection: 'row', alignItems: 'center', gap: '0.5rem' }}>
                <input type="checkbox" id="isActive" checked={formData.is_active} onChange={e => setFormData({...formData, is_active: e.target.checked})} style={{ width: 'auto' }} />
                <label htmlFor="isActive" style={{ margin: 0 }}>Đang hoạt động</label>
              </div>

              <div className="form-group full-width modal-actions">
                <button type="button" className="btn-cancel" onClick={() => setShowModal(false)}>Hủy</button>
                <button type="submit" className="btn-save">Lưu lại</button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminCouponsPage;
