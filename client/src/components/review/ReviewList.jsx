import { useEffect, useState, useCallback } from 'react';
import { getTourReviews } from '../../services/reviewService';
import './ReviewList.css';

/**
 * ReviewList - Component hiển thị danh sách review của tour
 * Bao gồm rating stats + danh sách review + phân trang + filter
 */
const ReviewList = ({ tourId }) => {
  const [reviews, setReviews] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [ratingFilter, setRatingFilter] = useState(0);

  const fetchReviews = useCallback(async () => {
    setLoading(true);
    try {
      const res = await getTourReviews(tourId, { page, limit: 5, rating: ratingFilter || undefined });
      setReviews(res.data || []);
      if (res.meta) {
        setTotalPages(res.meta.total_pages || 1);
        if (res.meta.stats) setStats(res.meta.stats);
      }
    } catch {
      setReviews([]);
    }
    setLoading(false);
  }, [tourId, page, ratingFilter]);

  useEffect(() => { fetchReviews(); }, [fetchReviews]);

  const renderStars = (rating) => {
    return '★'.repeat(rating) + '☆'.repeat(5 - rating);
  };

  const timeAgo = (dateStr) => {
    const date = new Date(dateStr);
    const diff = Date.now() - date.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    if (days === 0) return 'Hôm nay';
    if (days === 1) return 'Hôm qua';
    if (days < 30) return `${days} ngày trước`;
    if (days < 365) return `${Math.floor(days / 30)} tháng trước`;
    return `${Math.floor(days / 365)} năm trước`;
  };

  return (
    <div className="review-list-section">
      <h2 className="review-section-title">Đánh giá từ khách hàng</h2>

      {/* Rating Overview */}
      {stats && stats.total_reviews > 0 && (
        <div className="rating-overview">
          <div className="rating-big">
            <span className="rating-number">{stats.average_rating?.toFixed(1)}</span>
            <span className="rating-stars">{renderStars(Math.round(stats.average_rating))}</span>
            <span className="rating-count">{stats.total_reviews} đánh giá</span>
          </div>
          <div className="rating-bars">
            {[5, 4, 3, 2, 1].map(star => {
              const count = stats[`rating_${star}`] || 0;
              const pct = stats.total_reviews > 0 ? (count / stats.total_reviews * 100) : 0;
              return (
                <button
                  key={star}
                  className={`rating-bar-row ${ratingFilter === star ? 'active' : ''}`}
                  onClick={() => { setRatingFilter(ratingFilter === star ? 0 : star); setPage(1); }}
                >
                  <span className="bar-label">{star} ★</span>
                  <div className="bar-track">
                    <div className="bar-fill" style={{ width: `${pct}%` }}></div>
                  </div>
                  <span className="bar-count">{count}</span>
                </button>
              );
            })}
          </div>
        </div>
      )}

      {/* Filter indicator */}
      {ratingFilter > 0 && (
        <div className="filter-indicator">
          Lọc: {ratingFilter} sao
          <button onClick={() => { setRatingFilter(0); setPage(1); }}>✕ Bỏ lọc</button>
        </div>
      )}

      {/* Reviews List */}
      {loading ? (
        <div className="review-loading">Đang tải đánh giá...</div>
      ) : reviews.length === 0 ? (
        <div className="review-empty">
          {ratingFilter > 0 ? 'Không có đánh giá nào với bộ lọc này.' : 'Chưa có đánh giá nào. Hãy là người đầu tiên!'}
        </div>
      ) : (
        <div className="review-items">
          {reviews.map(review => (
            <div key={review.id} className="review-card">
              <div className="review-header">
                <div className="review-user">
                  <div className="review-avatar">
                    {review.user_avatar
                      ? <img src={review.user_avatar} alt="" />
                      : <span>{(review.user_name || 'K')[0].toUpperCase()}</span>
                    }
                  </div>
                  <div>
                    <div className="review-name">{review.user_name || 'Khách hàng'}</div>
                    <div className="review-date">{timeAgo(review.created_at)}</div>
                  </div>
                </div>
                <div className="review-rating">
                  <span className="stars">{renderStars(review.rating)}</span>
                </div>
              </div>

              {review.title && <h4 className="review-title">{review.title}</h4>}
              <p className="review-content">{review.content}</p>

              {review.admin_reply && (
                <div className="admin-reply">
                  <strong>Phản hồi từ Traveling:</strong>
                  <p>{review.admin_reply}</p>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="review-pagination">
          <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>← Trước</button>
          <span>Trang {page} / {totalPages}</span>
          <button disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>Sau →</button>
        </div>
      )}
    </div>
  );
};

export default ReviewList;
