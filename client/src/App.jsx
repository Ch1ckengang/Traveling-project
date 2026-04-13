import { BrowserRouter as Router, Navigate, Route, Routes, useLocation } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import { AuthProvider } from './context/AuthContext';
import { useAuth } from './context/AuthContext';
import { ThemeProvider } from './context/ThemeContext';
import Header from './components/Layout/Header';
import SearchBar from './components/Home/SearchBar';
import Login from './components/Auth/Login';
import Register from './components/Auth/Register';
import OtpVerification from './components/Auth/OtpVerification';
import ForgotPassword from './components/Auth/ForgotPassword';
import ResetSuccess from './components/Auth/ResetSuccess';
import Profile from './components/Profile/Profile';
import './App.css';

const API_BASE_URL = 'http://localhost:8080/v1/api';
const ITEMS_PER_VIEW = 3;

const formatCurrency = (amount) => {
  if (!Number.isFinite(amount)) {
    return '0đ';
  }

  return new Intl.NumberFormat('vi-VN').format(amount) + 'đ';
};

const normalizeCountry = (value) => {
  return (value || '')
    .toString()
    .trim()
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/\s+/g, '');
};

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
  const [adultCount, setAdultCount] = useState(1);
  const [childCount, setChildCount] = useState(0);
  const [infantCount, setInfantCount] = useState(0);
  const [travelDate, setTravelDate] = useState('');
  const [note, setNote] = useState('');
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingMessage, setBookingMessage] = useState('');
  const [bookingError, setBookingError] = useState('');
  const [bookingStep, setBookingStep] = useState(1);
  const [confirmedBooking, setConfirmedBooking] = useState(null);

  useEffect(() => {
    const normalizedPath = location.pathname.replace(/\/+$/, '') || '/';
    setActiveCategory(CATEGORY_BY_PATH[normalizedPath] || 'all');
    setCarouselIndex(0);

    // Khi đổi mục trên menu, luôn quay về màn danh sách tour của mục đó.
    setBookingStep(1);
    setSelectedTour(null);
    setConfirmedBooking(null);
    setBookingMessage('');
    setBookingError('');
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
        const rawTours = response.data || [];
        const nextTours = activeCategory === 'international'
          ? rawTours.filter((tour) => {
            const type = (tour.type || '').toString().trim().toLowerCase();
            const country = normalizeCountry(tour.country);
            return type === 'international' && country !== '' && country !== 'vietnam';
          })
          : rawTours;

        setTours(nextTours);
        setCarouselIndex(0);
      } catch (error) {
        setTours([]);
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
    setConfirmedBooking(null);
    setBookingStep(2);
  };

  const handleStartOver = () => {
    setBookingStep(1);
    setSelectedTour(null);
    setConfirmedBooking(null);
    setBookingMessage('');
    setBookingError('');
  };

  const validateBookingInput = () => {
    if (!selectedTour) {
      return 'Vui lòng chọn tour trước khi đặt.';
    }

    if (!fullName.trim() || !phone.trim() || !email.trim()) {
      return 'Vui lòng điền đầy đủ họ tên, số điện thoại và email.';
    }

    if (!travelDate) {
      return 'Vui lòng chọn ngày đi.';
    }

    if (adultCount <= 0) {
      return 'Cần ít nhất 1 người lớn để đặt tour.';
    }

    if (childCount < 0 || infantCount < 0) {
      return 'Số lượng trẻ em/trẻ nhỏ không hợp lệ.';
    }

    const totalGuest = adultCount + childCount + infantCount;
    if (totalGuest <= 0) {
      return 'Tổng số khách phải lớn hơn 0.';
    }

    if (selectedTour.remaining_slots && totalGuest > selectedTour.remaining_slots) {
      return `Tour chỉ còn ${selectedTour.remaining_slots} chỗ, vui lòng giảm số khách.`;
    }

    return '';
  };

  const handleCreateBooking = async () => {
    setBookingMessage('');
    setBookingError('');

    const validationError = validateBookingInput();
    if (validationError) {
      setBookingError(validationError);
      return;
    }

    const totalGuest = adultCount + childCount + infantCount;

    const userId = user?.id || 0;

    setBookingLoading(true);

    try {
      const response = await axios.post(`${API_BASE_URL}/bookings`, {
        user_id: userId,
        tour_id: selectedTour.id,
        full_name: fullName,
        phone,
        email,
        adult_count: adultCount,
        child_count: childCount,
        infant_count: infantCount,
        quantity: totalGuest,
        travel_date: travelDate,
        note
      });

      if (response.data.success) {
        const booking = response.data.booking;
        const successMessage = `Đặt tour thành công. Mã: ${booking.booking_code || ('#' + booking.id)} | Tổng tiền: ${formatCurrency(booking.total_amount || 0)} | Trạng thái: ${booking.payment_status === 'unpaid' ? 'Chưa thanh toán' : booking.payment_status}`;
        setBookingMessage(successMessage);
        setConfirmedBooking(booking);
        setAdultCount(1);
        setChildCount(0);
        setInfantCount(0);
        setNote('');
        setBookingStep(4);
      }
    } catch (error) {
      setBookingError(error.response?.data?.message || 'Không thể đặt tour, vui lòng thử lại.');
    } finally {
      setBookingLoading(false);
    }
  };

  return (
    <main className="home-page">
      {bookingStep === 1 && (
        <>
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
                        Xem chi tiết tour
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
        </>
      )}

      {bookingStep === 2 && selectedTour && (
        <section className="booking-panel booking-screen">
          <h3>Bước 2 - Chi tiết tour và phiếu đặt tour</h3>
          <div className="booking-tour-detail">
            <p><strong>Tên tour:</strong> {selectedTour.name}</p>
            <p><strong>Khởi hành:</strong> {selectedTour.departure_date || travelDate}</p>
            <p><strong>Thời gian:</strong> {selectedTour.duration}</p>
            <p><strong>Chi phí:</strong> {selectedTour.price}</p>
            <p><strong>Số chỗ còn lại:</strong> {selectedTour.remaining_slots ?? 'Chưa cập nhật'}</p>
            <p><strong>Lịch trình:</strong> {selectedTour.itinerary || selectedTour.description || 'Đang cập nhật'}</p>
          </div>

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
              Người lớn
              <input type="number" min="1" value={adultCount} onChange={(event) => setAdultCount(Number(event.target.value))} />
            </label>

            <label>
              Trẻ em (5-11 tuổi)
              <input type="number" min="0" value={childCount} onChange={(event) => setChildCount(Number(event.target.value))} />
            </label>

            <label>
              Trẻ nhỏ (dưới 5 tuổi)
              <input type="number" min="0" value={infantCount} onChange={(event) => setInfantCount(Number(event.target.value))} />
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

          <p className="state-text">Tổng số khách: {adultCount + childCount + infantCount}</p>

          {bookingError && <p className="state-error">{bookingError}</p>}

          <div className="booking-actions">
            <button type="button" className="ghost" onClick={() => setBookingStep(1)}>
              Quay lại tìm kiếm
            </button>
            <button type="button" onClick={handleCreateBooking} disabled={bookingLoading}>
              {bookingLoading ? 'Đang xử lý...' : 'Xác nhận đặt tour'}
            </button>
          </div>
        </section>
      )}

      {bookingStep === 4 && confirmedBooking && (
        <section className="booking-panel booking-screen">
          <h3>Đặt tour thành công</h3>
          <p className="state-text">Mã đặt tour: <strong>{confirmedBooking.booking_code}</strong></p>
          <p className="state-text">{bookingMessage || 'Vui lòng thanh toán trong vòng 24 giờ để giữ chỗ.'}</p>
          <div className="booking-actions">
            <button type="button" className="ghost" onClick={handleStartOver}>
              Đặt tour khác
            </button>
          </div>
        </section>
      )}
    </main>
  );
}

function AppShell() {
  const location = useLocation();
  const authPaths = ['/login', '/register', '/forgot-password', '/reset-success'];
  const hideHeader = authPaths.includes(location.pathname);

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
        <Route path="/otp-verification" element={<OtpVerification />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-success" element={<ResetSuccess />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </div>
  );
}

function App() {
  return (
    <ThemeProvider>
      <AuthProvider>
        <Router>
          <AppShell />
        </Router>
      </AuthProvider>
    </ThemeProvider>
  );
}

export default App;
