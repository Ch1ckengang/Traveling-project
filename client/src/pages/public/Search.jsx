import { useState, useEffect, useCallback } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { searchTours } from '../../services/tourService';
import '../../styles/Search.css';

const SearchPage = () => {
  const [searchParams, setSearchParams] = useSearchParams();
  const [keyword, setKeyword] = useState(searchParams.get('q') || '');
  const [inputValue, setInputValue] = useState(searchParams.get('q') || '');
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [hasSearched, setHasSearched] = useState(false);

  const doSearch = useCallback(async (kw) => {
    setLoading(true);
    setError('');
    setHasSearched(true);
    try {
      const result = await searchTours(kw);
      // API trả về { success, data: [...], total } hoặc array trực tiếp
      const list = result?.data || (Array.isArray(result) ? result : []);
      setTours(list);
    } catch {
      setError('Không thể tìm kiếm. Vui lòng thử lại.');
      setTours([]);
    } finally {
      setLoading(false);
    }
  }, []);

  // Tìm kiếm khi q param thay đổi
  useEffect(() => {
    const q = searchParams.get('q');
    if (q) {
      setKeyword(q);
      setInputValue(q);
      doSearch(q);
    }
  }, [searchParams, doSearch]);

  const handleSearch = (e) => {
    e.preventDefault();
    const kw = inputValue.trim();
    if (!kw) return;
    setKeyword(kw);
    setSearchParams({ q: kw });
  };

  const getTourTypeLabel = (type) => {
    switch (type) {
      case 'domestic': return '🇻🇳 Việt Nam';
      case 'international': return '🌍 Quốc tế';
      case 'service': return '✈️ Dịch vụ';
      default: return type;
    }
  };

  return (
    <div className="search-page">
      {/* Search Header */}
      <div className="search-hero">
        <div className="search-hero-content">
          <h1>Tìm kiếm Tour</h1>
          <p>Khám phá hàng chục tour hấp dẫn trong và ngoài nước</p>
          <form className="search-form" onSubmit={handleSearch}>
            <div className="search-input-wrap">
              <span className="search-icon">🔍</span>
              <input
                id="search-input"
                type="text"
                className="search-input"
                placeholder="Nhập địa điểm, tên tour..."
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                autoFocus
              />
              {inputValue && (
                <button
                  type="button"
                  className="search-clear-btn"
                  onClick={() => { setInputValue(''); }}
                  aria-label="Xóa"
                >✕</button>
              )}
            </div>
            <button type="submit" className="search-btn" id="search-submit-btn">
              Tìm kiếm
            </button>
          </form>
        </div>
      </div>

      {/* Results */}
      <div className="search-results-container">
        {keyword && hasSearched && (
          <div className="search-results-header">
            {loading ? (
              <p className="search-status">Đang tìm kiếm &quot;{keyword}&quot;...</p>
            ) : (
              <p className="search-status">
                {error ? '' : `Tìm thấy `}
                {!error && <strong>{tours.length} tour</strong>}
                {!error && ` cho từ khóa `}
                {!error && <strong>&quot;{keyword}&quot;</strong>}
              </p>
            )}
          </div>
        )}

        {/* Loading */}
        {loading && (
          <div className="search-loading">
            <div className="search-spinner"></div>
            <p>Đang tìm kiếm...</p>
          </div>
        )}

        {/* Error */}
        {error && !loading && (
          <div className="search-error">
            <span>⚠️</span>
            <p>{error}</p>
          </div>
        )}

        {/* Empty State */}
        {!loading && !error && hasSearched && tours.length === 0 && (
          <div className="search-empty">
            <div className="search-empty-icon">🔍</div>
            <h3>Không tìm thấy kết quả</h3>
            <p>Không có tour nào khớp với &quot;{keyword}&quot;. Hãy thử từ khóa khác.</p>
            <Link to="/tours" className="search-browse-link">
              Xem tất cả tour →
            </Link>
          </div>
        )}

        {/* Initial State (no search yet) */}
        {!loading && !error && !hasSearched && (
          <div className="search-empty">
            <div className="search-empty-icon">✈️</div>
            <h3>Nhập từ khóa để tìm tour</h3>
            <p>Bạn có thể tìm theo tên tour, địa điểm, hoặc quốc gia.</p>
          </div>
        )}

        {/* Results Grid */}
        {!loading && !error && tours.length > 0 && (
          <div className="search-grid">
            {tours.map((tour) => (
              <Link
                key={tour.id}
                to={`/tours/${tour.id}`}
                className="search-card"
              >
                <div className="search-card-image">
                  {tour.image_url ? (
                    <img src={tour.image_url} alt={tour.name} />
                  ) : (
                    <div className="search-card-placeholder">
                      <span>📷</span>
                      <span>{tour.location}</span>
                    </div>
                  )}
                  <span className="search-card-type">{getTourTypeLabel(tour.type)}</span>
                </div>
                <div className="search-card-body">
                  <h3 className="search-card-name">{tour.name}</h3>
                  <p className="search-card-location">📍 {tour.location}{tour.country && tour.country !== 'Việt Nam' ? `, ${tour.country}` : ''}</p>
                  <p className="search-card-desc">{tour.description}</p>
                  <div className="search-card-footer">
                    <span className="search-card-price">{tour.price}</span>
                    <span className="search-card-duration">⏱ {tour.duration}</span>
                  </div>
                  <div className="search-card-slots">
                    Còn {tour.remaining_slots} chỗ
                  </div>
                </div>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default SearchPage;

