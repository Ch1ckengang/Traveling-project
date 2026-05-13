import { useEffect, useState, useCallback } from 'react';
import { adminGetReviews, adminPublishReview, adminHideReview, adminReplyReview } from '../../services/reviewService';
import './AdminPage.css';

const AdminReviewsPage = () => {
  const [reviews, setReviews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [meta, setMeta] = useState({});
  const [search, setSearch] = useState('');
  const [statusFilter, setStatusFilter] = useState('');
  const [replyModal, setReplyModal] = useState(null);
  const [replyText, setReplyText] = useState('');

  const fetchReviews = useCallback(async () => {
    setLoading(true);
    try {
      const res = await adminGetReviews({ page, limit: 10, status: statusFilter, search });
      setReviews(res.data || []);
      setMeta(res.meta || {});
    } catch {
      setReviews([]);
    }
    setLoading(false);
  }, [page, search, statusFilter]);

  useEffect(() => { fetchReviews(); }, [fetchReviews]);

  const handlePublish = async (id) => {
    try { await adminPublishReview(id); fetchReviews(); }
    catch (err) { alert(err.response?.data?.message || 'Lỗi'); }
  };

  const handleHide = async (id) => {
    if (!window.confirm('Ẩn review này?')) return;
    try { await adminHideReview(id); fetchReviews(); }
    catch (err) { alert(err.response?.data?.message || 'Lỗi'); }
  };

  const handleReply = async () => {
    if (!replyModal || !replyText.trim()) return;
    try {
      await adminReplyReview(replyModal, replyText);
      setReplyModal(null);
      setReplyText('');
      fetchReviews();
    } catch (err) {
      alert(err.response?.data?.message || 'Lỗi');
    }
  };

  const renderStars = (rating) => '★'.repeat(rating) + '☆'.repeat(5 - rating);

  const statusLabels = {
    published: 'Hiển thị',
    hidden: 'Ẩn',
    pending: 'Chờ duyệt',
  };

  const formatDate = (d) => d ? new Date(d).toLocaleDateString('vi-VN') : '—';

  return (
    <div className="admin-page">
      <div className="admin-header">
        <h1>Quản lý Đánh giá</h1>
      </div>

      <div className="admin-filters">
        <input
          type="text"
          placeholder="Tìm nội dung, tiêu đề..."
          value={search}
          onChange={(e) => { setSearch(e.target.value); setPage(1); }}
          className="filter-input"
        />
        <select value={statusFilter} onChange={(e) => { setStatusFilter(e.target.value); setPage(1); }} className="filter-select">
          <option value="">Tất cả trạng thái</option>
          <option value="published">Hiển thị</option>
          <option value="hidden">Ẩn</option>
          <option value="pending">Chờ duyệt</option>
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
                  <th>Tour</th>
                  <th>Khách hàng</th>
                  <th>Rating</th>
                  <th>Nội dung</th>
                  <th>Trạng thái</th>
                  <th>Ngày</th>
                  <th>Thao tác</th>
                </tr>
              </thead>
              <tbody>
                {reviews.length === 0 ? (
                  <tr><td colSpan="8" className="empty-row">Không có review nào</td></tr>
                ) : reviews.map(r => (
                  <tr key={r.id}>
                    <td>{r.id}</td>
                    <td className="tour-name">{r.tour_name || '—'}</td>
                    <td>
                      <div>{r.user_name}</div>
                      <div className="text-muted">{r.user_email}</div>
                    </td>
                    <td style={{ color: '#f59e0b' }}>{renderStars(r.rating)}</td>
                    <td style={{ maxWidth: 250, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }}>
                      {r.title ? <strong>{r.title}: </strong> : ''}{r.content}
                    </td>
                    <td>
                      <span className={`badge badge-${r.status === 'published' ? 'confirmed' : r.status === 'hidden' ? 'cancelled' : 'pending'}`}>
                        {statusLabels[r.status] || r.status}
                      </span>
                    </td>
                    <td>{formatDate(r.created_at)}</td>
                    <td className="actions">
                      {r.status !== 'published' && (
                        <button className="btn-sm btn-confirm" onClick={() => handlePublish(r.id)}>Publish</button>
                      )}
                      {r.status !== 'hidden' && (
                        <button className="btn-sm btn-cancel-action" onClick={() => handleHide(r.id)}>Ẩn</button>
                      )}
                      <button className="btn-sm btn-edit" onClick={() => { setReplyModal(r.id); setReplyText(r.admin_reply || ''); }}>
                        {r.admin_reply ? 'Sửa PH' : 'Phản hồi'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {meta.total_pages > 1 && (
            <div className="admin-pagination">
              <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>← Trước</button>
              <span>Trang {page} / {meta.total_pages} ({meta.total} reviews)</span>
              <button disabled={page >= meta.total_pages} onClick={() => setPage(p => p + 1)}>Sau →</button>
            </div>
          )}
        </>
      )}

      {/* Reply Modal */}
      {replyModal && (
        <div className="modal-overlay" onClick={() => setReplyModal(null)}>
          <div className="modal-content" onClick={e => e.stopPropagation()} style={{ maxWidth: 500 }}>
            <h2>Phản hồi đánh giá</h2>
            <textarea
              rows="4"
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              placeholder="Nhập phản hồi..."
              style={{ width: '100%', padding: '0.75rem', borderRadius: 8, border: '1px solid #d1d5db', resize: 'vertical', fontFamily: 'inherit' }}
            />
            <div className="modal-actions">
              <button className="btn-cancel" onClick={() => setReplyModal(null)}>Hủy</button>
              <button className="btn-save" onClick={handleReply} disabled={!replyText.trim()}>Gửi phản hồi</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminReviewsPage;
