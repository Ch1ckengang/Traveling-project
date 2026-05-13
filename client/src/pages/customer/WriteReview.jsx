import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { createReview } from '../../services/reviewService';
import './WriteReview.css';

const WriteReviewPage = () => {
  const navigate = useNavigate();
  const { user } = useAuth();
  const [searchParams] = useSearchParams();
  const tourId = searchParams.get('tourId');
  const bookingId = searchParams.get('bookingId');
  const tourName = searchParams.get('tourName') || 'Tour';

  const [rating, setRating] = useState(0);
  const [hoverRating, setHoverRating] = useState(0);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  if (!user) {
    return (
      <div className="write-review-page">
        <div className="review-form-card">
          <p>Vui lòng đăng nhập để viết đánh giá.</p>
          <button onClick={() => navigate('/login')}>Đăng nhập</button>
        </div>
      </div>
    );
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (rating === 0) {
      setError('Vui lòng chọn số sao đánh giá');
      return;
    }
    if (content.length < 10) {
      setError('Nội dung đánh giá tối thiểu 10 ký tự');
      return;
    }

    setSubmitting(true);
    try {
      const res = await createReview({
        tour_id: parseInt(tourId),
        booking_id: parseInt(bookingId),
        rating,
        title,
        content,
      });

      if (res.success) {
        alert('Đánh giá thành công! Cảm ơn bạn.');
        navigate('/account/bookings');
      } else {
        setError(res.message || 'Không thể gửi đánh giá');
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Có lỗi xảy ra');
    }
    setSubmitting(false);
  };

  return (
    <div className="write-review-page">
      <div className="review-form-card">
        <h1>Đánh giá tour</h1>
        <p className="tour-name-label">{tourName}</p>

        {error && <div className="review-error">{error}</div>}

        <form onSubmit={handleSubmit}>
          {/* Star Rating */}
          <div className="rating-selector">
            <label>Đánh giá của bạn *</label>
            <div className="stars-row">
              {[1, 2, 3, 4, 5].map(star => (
                <button
                  key={star}
                  type="button"
                  className={`star-btn ${star <= (hoverRating || rating) ? 'active' : ''}`}
                  onMouseEnter={() => setHoverRating(star)}
                  onMouseLeave={() => setHoverRating(0)}
                  onClick={() => setRating(star)}
                >
                  ★
                </button>
              ))}
              <span className="rating-text">
                {rating === 1 && 'Rất tệ'}
                {rating === 2 && 'Tệ'}
                {rating === 3 && 'Bình thường'}
                {rating === 4 && 'Tốt'}
                {rating === 5 && 'Tuyệt vời'}
              </span>
            </div>
          </div>

          {/* Title */}
          <div className="form-field">
            <label>Tiêu đề (tùy chọn)</label>
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="VD: Trải nghiệm tuyệt vời!"
              maxLength={200}
            />
          </div>

          {/* Content */}
          <div className="form-field">
            <label>Nội dung đánh giá *</label>
            <textarea
              rows="5"
              value={content}
              onChange={(e) => setContent(e.target.value)}
              placeholder="Chia sẻ trải nghiệm của bạn về tour này..."
              maxLength={2000}
            />
            <span className="char-count">{content.length}/2000</span>
          </div>

          <div className="form-actions">
            <button type="button" className="btn-back" onClick={() => navigate(-1)}>
              Quay lại
            </button>
            <button type="submit" className="btn-submit" disabled={submitting || rating === 0}>
              {submitting ? 'Đang gửi...' : 'Gửi đánh giá'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default WriteReviewPage;
