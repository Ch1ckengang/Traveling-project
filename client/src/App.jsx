import { BrowserRouter as Router, Navigate, Route, Routes, useLocation } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import { AuthProvider } from './context/AuthContext';
import { useAuth } from './context/AuthContext';
import Header from './components/Layout/Header';
import SearchBar from './components/Home/SearchBar';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import Profile from './components/Profile/Profile';
import './App.css';

const API_BASE_URL = 'http://localhost:8080/v1/api';
const ITEMS_PER_VIEW = 3;

const CATEGORY_BY_PATH = {
  '/': 'all',
  '/vietnam': 'domestic',
  '/quocte': 'international',
  '/dichvu': 'service'
};

const CATEGORY_TITLES = {
  all: 'Danh sách tour nổi bật',
  domestic: 'Du lịch Việt Nam',
  international: 'Du lịch quốc tế',
  service: 'Dịch vụ du lịch'
};

const defaultFilters = {
  city: '',
  duration: 'all',
  price: 'all',
  sort: 'default'
};

function HomePage() {
  const location = useLocation();
  const { user } = useAuth();

  const [activeCategory, setActiveCategory] = useState('all');
  const [filters, setFilters] = useState(defaultFilters);
  const [appliedFilters, setAppliedFilters] = useState(defaultFilters);
  const [tours, setTours] = useState([]);
  const [loadingTours, setLoadingTours] = useState(false);
  const [loadError, setLoadError] = useState('');
  const [carouselIndex, setCarouselIndex] = useState(0);

  const [selectedTour, setSelectedTour] = useState(null);
  const [fullName, setFullName] = useState('');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [quantity, setQuantity] = useState(1);
  const [travelDate, setTravelDate] = useState('');
  const [note, setNote] = useState('');
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingMessage, setBookingMessage] = useState('');
  const [bookingError, setBookingError] = useState('');

  useEffect(() => {
    setActiveCategory(CATEGORY_BY_PATH[location.pathname] || 'all');
    setCarouselIndex(0);
  }, [location.pathname]);

  useEffect(() => {
    const nextWeek = new Date();
    nextWeek.setDate(nextWeek.getDate() + 7);
    setTravelDate(nextWeek.toISOString().slice(0, 10));
  }, []);

  useEffect(() => {
    if (!user) {
      return;
    }

    setFullName(user?.name || '');
    setEmail(user?.email || '');
  }, [user]);

  useEffect(() => {
    const fetchTours = async () => {
      setLoadingTours(true);
      setLoadError('');
      setBookingMessage('');
      setBookingError('');

      const params = {
        category: activeCategory,
        city: appliedFilters.city.trim(),
        duration: appliedFilters.duration,
        price: appliedFilters.price,
        sort: appliedFilters.sort
      };

      if (!params.city) {
        delete params.city;
      }

      if (params.duration === 'all') {
        delete params.duration;
      }

      if (params.price === 'all') {
        delete params.price;
      }

      if (params.category === 'all') {
        delete params.category;
      }

      if (params.sort === 'default') {
        delete params.sort;
      }

      try {
        const response = await axios.get(`${API_BASE_URL}/tours`, { params });
        const nextTours = response.data || [];

        setTours(nextTours);
        setCarouselIndex(0);
      } catch (error) {
        setLoadError(error.response?.data?.message || 'Không thể tải danh sách tour.');
      } finally {
        setLoadingTours(false);
      }
    };

    fetchTours();
  }, [activeCategory, appliedFilters]);

  useEffect(() => {
    if (!selectedTour) {
      return;
    }

    if (!tours.some((tour) => tour.id === selectedTour.id)) {
      setSelectedTour(null);
    }
  }, [selectedTour, tours]);

  const handleSearch = () => {
    setAppliedFilters({ ...filters });
  };

  const visibleTours = useMemo(() => {
    if (tours.length <= ITEMS_PER_VIEW) {
      return tours;
    }

    return Array.from({ length: ITEMS_PER_VIEW }, (_, offset) => {
      const index = (carouselIndex + offset) % tours.length;
      return tours[index];
    });
  }, [carouselIndex, tours]);

  const canSlide = tours.length > ITEMS_PER_VIEW;

  const handlePrev = () => {
    if (!canSlide) {
      return;
    }

    setCarouselIndex((prev) => (prev - 1 + tours.length) % tours.length);
  };

  const handleNext = () => {
    if (!canSlide) {
      return;
    }

    setCarouselIndex((prev) => (prev + 1) % tours.length);
  };

  const handleSelectTour = (tour) => {
    setSelectedTour(tour);
    setBookingMessage('');
    setBookingError('');
  };

  const handleCreateBooking = async () => {
    setBookingMessage('');
    setBookingError('');

    if (!selectedTour) {
      setBookingError('Vui lòng chọn tour trước khi đặt.');
      return;
    }

    if (!fullName.trim() || !phone.trim() || !email.trim()) {
      setBookingError('Vui lòng điền đầy đủ họ tên, số điện thoại và email.');
      return;
    }

    if (!travelDate) {
      setBookingError('Vui lòng chọn ngày đi.');
      return;
    }

    if (quantity <= 0) {
      setBookingError('Số lượng khách phải lớn hơn 0.');
      return;
    }

    const userId = user?.id || 0;

    setBookingLoading(true);

    try {
      const response = await axios.post(`${API_BASE_URL}/bookings`, {
        user_id: userId,
        tour_id: selectedTour.id,
        full_name: fullName,
        phone,
        email,
        quantity,
        travel_date: travelDate,
        note
      });

      if (response.data.success) {
        setBookingMessage(`Đặt tour thành công. Mã đơn: #${response.data.booking.id}`);
        setQuantity(1);
        setNote('');
      }
    } catch (error) {
      setBookingError(error.response?.data?.message || 'Không thể đặt tour, vui lòng thử lại.');
    } finally {
      setBookingLoading(false);
    }
  };

  return (
    <main className="home-page">
      <SearchBar filters={filters} onFiltersChange={setFilters} onSearch={handleSearch} />

      <section className="tour-section">
        <div className="tour-section-header">
          <h2>{CATEGORY_TITLES[activeCategory] || CATEGORY_TITLES.all}</h2>
          <p>{tours.length} kết quả</p>
        </div>

        {loadingTours && <p className="state-text">Đang tải dữ liệu...</p>}
        {!loadingTours && loadError && <p className="state-error">{loadError}</p>}

        {!loadingTours && !loadError && tours.length === 0 && (
          <p className="state-text">Không tìm thấy tour phù hợp với bộ lọc hiện tại.</p>
        )}

        {!loadingTours && !loadError && tours.length > 0 && (
          <div className="carousel-shell">
            <button type="button" className="carousel-control" onClick={handlePrev} disabled={!canSlide} aria-label="Tour trước">
              &lt;
            </button>

            <div className="tour-grid">
              {visibleTours.map((tour) => (
                <article className="tour-card" key={tour.id}>
                  <div className="tour-card-thumb" />
                  <h3>{tour.name}</h3>
                  <p className="tour-meta">{tour.location} | {tour.duration}</p>
                  <p className="tour-price">{tour.price}</p>
                  <p className="tour-description">{tour.description || 'Trọn gói bao gồm tham quan và hỗ trợ lịch trình.'}</p>
                  <button type="button" className="tour-action" onClick={() => handleSelectTour(tour)}>
                    Chọn tour
                  </button>
                </article>
              ))}
            </div>

            <button type="button" className="carousel-control" onClick={handleNext} disabled={!canSlide} aria-label="Tour tiếp theo">
              &gt;
            </button>
          </div>
        )}
      </section>

      {selectedTour && (
        <section className="booking-panel">
          <h3>Đặt tour: {selectedTour.name}</h3>

          <div className="booking-grid">
            <label>
              Họ và tên
              <input value={fullName} onChange={(event) => setFullName(event.target.value)} />
            </label>

            <label>
              Số điện thoại
              <input value={phone} onChange={(event) => setPhone(event.target.value)} />
            </label>

            <label>
              Email
              <input type="email" value={email} onChange={(event) => setEmail(event.target.value)} />
            </label>

            <label>
              Số lượng
              <input type="number" min="1" value={quantity} onChange={(event) => setQuantity(Number(event.target.value))} />
            </label>

            <label>
              Ngày đi
              <input type="date" value={travelDate} onChange={(event) => setTravelDate(event.target.value)} />
            </label>
          </div>

          <label className="booking-note">
            Ghi chú
            <textarea rows="3" value={note} onChange={(event) => setNote(event.target.value)} />
          </label>

          <div className="booking-actions">
            <button type="button" onClick={handleCreateBooking} disabled={bookingLoading}>
              {bookingLoading ? 'Đang xử lý...' : 'Xác nhận đặt tour'}
            </button>
            <button type="button" className="ghost" onClick={() => setSelectedTour(null)}>
              Đóng
            </button>
          </div>

          {bookingMessage && <p className="state-success">{bookingMessage}</p>}
          {bookingError && <p className="state-error">{bookingError}</p>}
        </section>
      )}
    </main>
  );
}

function AppShell() {
  const location = useLocation();
  const hideHeader = location.pathname === '/login' || location.pathname === '/register';

  return (
    <div className="app-shell">
      {!hideHeader && <Header />}

      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/vietnam" element={<HomePage />} />
        <Route path="/quocte" element={<HomePage />} />
        <Route path="/dichvu" element={<HomePage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </div>
  );
}

function App() {
  return (
    <AuthProvider>
      <Router>
        <AppShell />
      </Router>
    </AuthProvider>
  );
}

export default App;
