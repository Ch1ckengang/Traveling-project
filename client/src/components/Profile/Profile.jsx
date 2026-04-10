import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { useTheme } from '../../context/ThemeContext';
import axios from 'axios';
import '../../styles/Profile.css';

const API_BASE_URL = 'http://localhost:8080/v1/api';

const formatCurrency = (amount) => {
  if (!Number.isFinite(amount)) {
    return '0đ';
  }

  return new Intl.NumberFormat('vi-VN').format(amount) + 'đ';
};

const formatDateTime = (dateTime) => {
  if (!dateTime) {
    return 'Chưa có dữ liệu';
  }

  const date = new Date(dateTime);
  if (Number.isNaN(date.getTime())) {
    return 'Chưa có dữ liệu';
  }

  return new Intl.DateTimeFormat('vi-VN', {
    dateStyle: 'medium',
    timeStyle: 'short'
  }).format(date);
};

/**
 * Profile - Component trang quản lý thông tin cá nhân
 * Chức năng:
 * - Hiển thị thông tin user (tên, email, avatar)
 * - Chỉnh sửa thông tin cá nhân
 * - Đổi mật khẩu (tùy chọn)
 * - Kiểm tra email trùng lặp khi cập nhật
 * - Validation form đầy đủ
 */
const Profile = () => {
  const { user, login } = useAuth();
  const { theme, isDarkMode, toggleTheme } = useTheme();
  const navigate = useNavigate();

  const [activeSection, setActiveSection] = useState('profile');
  
  // State quản lý form data
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  const [originalData, setOriginalData] = useState({
    name: '',
    email: ''
  });
  
  const [isEditing, setIsEditing] = useState(false); // Trạng thái chế độ xem/chỉnh sửa
  const [error, setError] = useState(''); // Thông báo lỗi
  const [success, setSuccess] = useState(''); // Thông báo thành công
  const [loading, setLoading] = useState(false); // Trạng thái đang xử lý

  const [bookings, setBookings] = useState([]);
  const [bookingLoading, setBookingLoading] = useState(false);
  const [bookingError, setBookingError] = useState('');
  const [bookingSuccess, setBookingSuccess] = useState('');
  const [bookingReloadToken, setBookingReloadToken] = useState(0);
  const [cancelingBookingId, setCancelingBookingId] = useState(null);

  // Load thông tin user vào form khi component mount hoặc user thay đổi
  useEffect(() => {
    // Nếu chưa đăng nhập -> chuyển về trang login
    if (!user) {
      navigate('/login');
      return;
    }
    
    // Set giá trị mặc định cho form từ thông tin user hiện tại
    const nextData = {
      name: user.name || '',
      email: user.email || '',
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    };

    setFormData(nextData);
    setOriginalData({
      name: nextData.name,
      email: nextData.email
    });
    setIsEditing(false);
  }, [user, navigate]);

  useEffect(() => {
    if (activeSection !== 'bookings' || !user?.id) {
      return;
    }

    const fetchBookings = async () => {
      setBookingLoading(true);
      setBookingError('');
      setBookingSuccess('');

      try {
        const response = await axios.get(`${API_BASE_URL}/users/${user.id}/bookings`);
        setBookings(response.data?.bookings || []);
      } catch (err) {
        setBookingError(err.response?.data?.message || 'Không thể tải danh sách tour đã đặt.');
      } finally {
        setBookingLoading(false);
      }
    };

    fetchBookings();
  }, [activeSection, user?.id, bookingReloadToken]);

  /**
   * handleChange - Xử lý khi user nhập vào input
   * Cập nhật formData và clear error/success messages
   */
  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
    setError('');
    setSuccess('');
  };

  const canCancelBooking = (booking) => {
    const status = (booking?.status || '').toString().trim().toLowerCase();
    return status !== 'cancelled' && status !== 'canceled';
  };

  const handleCancelBooking = async (booking) => {
    if (!user?.id || !booking?.id || !canCancelBooking(booking)) {
      return;
    }

    const confirmed = window.confirm(`Bạn có chắc muốn hủy tour "${booking.tour?.name || 'đã đặt'}" không?`);
    if (!confirmed) {
      return;
    }

    setBookingError('');
    setBookingSuccess('');
    setCancelingBookingId(booking.id);

    try {
      const response = await axios.put(`${API_BASE_URL}/users/${user.id}/bookings/${booking.id}/cancel`);
      if (response.data?.success) {
        setBookingSuccess(response.data?.message || 'Hủy tour thành công.');
        setBookingReloadToken((prev) => prev + 1);
      }
    } catch (err) {
      setBookingError(err.response?.data?.message || 'Không thể hủy tour đã đặt.');
    } finally {
      setCancelingBookingId(null);
    }
  };

  /**
   * handleSubmit - Xử lý khi user submit form cập nhật thông tin
   * Luồng:
   * 1. Validate form (tên, email không rỗng, password match, length)
   * 2. Gửi PUT request đến /api/users/:id
   * 3. Backend kiểm tra email trùng lặp
   * 4. Nếu success -> cập nhật context & localStorage, hiển thị thông báo
   * 5. Nếu lỗi -> hiển thị thông báo lỗi (đặc biệt: email trùng)
   */
  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setSuccess('');

    if (!user?.id) {
      setError('Không xác định được tài khoản. Vui lòng đăng nhập lại.');
      return;
    }

    const cleanedName = formData.name.trim();
    const cleanedEmail = formData.email.trim();

    // VALIDATION
    if (cleanedName === '') {
      setError('Tên không được để trống');
      return;
    }

    if (cleanedEmail === '') {
      setError('Email không được để trống');
      return;
    }

    const hasProfileChange = cleanedName !== originalData.name || cleanedEmail !== originalData.email;
    const hasPasswordChange = formData.newPassword.trim() !== '';

    if (!hasProfileChange && !hasPasswordChange) {
      setError('Bạn chưa thay đổi thông tin nào để lưu.');
      return;
    }

    // Nếu muốn đổi mật khẩu, validate mật khẩu mới
    if (formData.newPassword) {
      if (formData.newPassword.length < 6) {
        setError('Mật khẩu mới phải có ít nhất 6 ký tự');
        return;
      }
      
      if (formData.newPassword !== formData.confirmPassword) {
        setError('Mật khẩu xác nhận không khớp');
        return;
      }
    }

    setLoading(true);

    try {
      // Chuẩn bị data để gửi
      const updateData = {
        name: cleanedName,
        email: cleanedEmail,
      };

      // Thêm password nếu user muốn đổi
      if (formData.newPassword) {
        updateData.password = formData.newPassword;
      }

      // GỬI REQUEST ĐẾN BACKEND
      const response = await axios.put(
        `http://localhost:8080/v1/api/users/${user.id}`,
        updateData
      );

      if (response.data.success) {
        // Cập nhật context với thông tin mới (tự động update header)
        login(response.data.user);

        const updatedName = response.data.user?.name || cleanedName;
        const updatedEmail = response.data.user?.email || cleanedEmail;

        setOriginalData({
          name: updatedName,
          email: updatedEmail
        });
        
        setSuccess('Cập nhật thông tin thành công!');
        setIsEditing(false); // Chuyển về chế độ xem
        
        // Reset password fields
        setFormData({
          ...formData,
          name: updatedName,
          email: updatedEmail,
          currentPassword: '',
          newPassword: '',
          confirmPassword: ''
        });
      }
    } catch (err) {
      // XỬ LÝ LỖI
      if (err.response?.status === 409) {
        // Status 409 = Conflict -> Email đã tồn tại
        setError('Email đã được sử dụng bởi tài khoản khác');
      } else {
        setError(err.response?.data?.message || 'Cập nhật thất bại. Vui lòng thử lại.');
      }
    } finally {
      setLoading(false);
    }
  };

  /**
   * handleCancel - Xử lý khi user hủy chỉnh sửa
   * Reset form về giá trị ban đầu và chuyển về chế độ xem
   */
  const handleCancel = () => {
    setFormData({
      name: originalData.name,
      email: originalData.email,
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    });
    setIsEditing(false);
    setError('');
    setSuccess('');
  };

  if (!user) {
    return null;
  }

  const renderBookedTours = () => {
    if (bookingLoading) {
      return <p className="bookings-placeholder">Đang tải danh sách tour đã đặt...</p>;
    }

    if (bookingError) {
      return <div className="error-message">{bookingError}</div>;
    }

    if (bookings.length === 0) {
      return <p className="bookings-placeholder">Bạn chưa có tour nào đã đặt.</p>;
    }

    return (
      <div className="bookings-list">
        {bookings.map((booking) => (
          <article className="booking-card" key={booking.id}>
            <div className="booking-card-top">
              <h4>{booking.tour?.name || 'Tour không xác định'}</h4>
              <span className="booking-status">{booking.status || 'booked'}</span>
            </div>
            <p className="booking-code">Mã đặt chỗ: {booking.booking_code || '---'}</p>

            <div className="booking-grid">
              <p><strong>Ngày đi:</strong> {booking.travel_date || '---'}</p>
              <p><strong>Địa điểm:</strong> {booking.tour?.location || '---'}</p>
              <p><strong>Số khách:</strong> {booking.quantity || 0}</p>
              <p><strong>Tổng tiền:</strong> {formatCurrency(booking.total_amount || 0)}</p>
              <p><strong>Thanh toán:</strong> {booking.payment_status || 'unpaid'}</p>
              <p><strong>Đặt lúc:</strong> {formatDateTime(booking.created_at)}</p>
            </div>

            <div className="booking-actions-row">
              <button
                type="button"
                className="btn-cancel-booking"
                onClick={() => handleCancelBooking(booking)}
                disabled={!canCancelBooking(booking) || cancelingBookingId === booking.id}
              >
                {cancelingBookingId === booking.id ? 'Đang hủy...' : 'Hủy tour'}
              </button>
            </div>
          </article>
        ))}
      </div>
    );
  };

  return (
    <div className="profile-container">
      <div className="profile-box">
        <aside className="profile-sidebar">
          <div className="profile-header">
            <div className="profile-avatar-large">
              {user.name?.charAt(0).toUpperCase()}
            </div>
            <h2>{user.name}</h2>
            <p className="profile-email">{user.email}</p>
          </div>

          <nav className="profile-menu" aria-label="Menu hồ sơ khách hàng">
            <button
              type="button"
              className={`profile-menu-item ${activeSection === 'profile' ? 'is-active' : ''}`}
              onClick={() => {
                setActiveSection('profile');
                setError('');
                setSuccess('');
              }}
            >
              Thông tin cá nhân
            </button>
            <button
              type="button"
              className={`profile-menu-item ${activeSection === 'bookings' ? 'is-active' : ''}`}
              onClick={() => {
                setActiveSection('bookings');
                setError('');
                setSuccess('');
              }}
            >
              Tour đã đặt
            </button>
          </nav>
        </aside>

        <section className="profile-content">
          {activeSection === 'profile' && (
            <>
              {error && <div className="error-message">{error}</div>}
              {success && <div className="success-message">{success}</div>}

              <form onSubmit={handleSubmit} className="profile-form">
                <div className="form-section theme-section">
                  <h3>Giao diện</h3>
                  <div className="theme-mode-row">
                    <div>
                      <p className="theme-mode-title">Chế độ sáng/tối</p>
                      <p className="theme-mode-subtitle">
                        Chuyển đổi giữa giao diện sáng và tối. Cài đặt được lưu cho lần truy cập sau.
                      </p>
                    </div>

                    <button
                      type="button"
                      className="theme-toggle-btn"
                      onClick={toggleTheme}
                      aria-label="Chuyển chế độ sáng tối"
                      aria-pressed={isDarkMode}
                    >
                      <span className={`theme-toggle-pill ${isDarkMode ? 'is-dark' : 'is-light'}`}>
                        <span className="theme-toggle-thumb" />
                      </span>
                      <span className="theme-toggle-label">{theme === 'dark' ? 'Đang dùng: Tối' : 'Đang dùng: Sáng'}</span>
                    </button>
                  </div>
                </div>

                <div className="form-section">
                  <h3>Thông tin cá nhân</h3>

                  <div className="form-group">
                    <label>Họ và tên</label>
                    <input
                      type="text"
                      name="name"
                      value={formData.name}
                      onChange={handleChange}
                      disabled={!isEditing}
                      required
                    />
                  </div>

                  <div className="form-group">
                    <label>Email</label>
                    <input
                      type="email"
                      name="email"
                      value={formData.email}
                      onChange={handleChange}
                      disabled={!isEditing}
                      required
                    />
                  </div>
                </div>

                {isEditing && (
                  <div className="form-section">
                    <h3>Đổi mật khẩu (tùy chọn)</h3>

                    <div className="form-group">
                      <label>Mật khẩu mới</label>
                      <input
                        type="password"
                        name="newPassword"
                        value={formData.newPassword}
                        onChange={handleChange}
                        placeholder="Nhập mật khẩu mới (nếu muốn đổi)"
                      />
                    </div>

                    <div className="form-group">
                      <label>Xác nhận mật khẩu mới</label>
                      <input
                        type="password"
                        name="confirmPassword"
                        value={formData.confirmPassword}
                        onChange={handleChange}
                        placeholder="Nhập lại mật khẩu mới"
                      />
                    </div>
                  </div>
                )}

                <div className="profile-actions">
                  {!isEditing ? (
                    <button
                      type="button"
                      onClick={() => setIsEditing(true)}
                      className="btn-edit"
                    >
                      Chỉnh sửa thông tin
                    </button>
                  ) : (
                    <>
                      <button
                        type="submit"
                        className="btn-save"
                        disabled={loading}
                      >
                        {loading ? 'Đang lưu...' : 'Lưu thay đổi'}
                      </button>
                      <button
                        type="button"
                        onClick={handleCancel}
                        className="btn-cancel"
                        disabled={loading}
                      >
                        Hủy
                      </button>
                    </>
                  )}
                </div>
              </form>
            </>
          )}

          {activeSection === 'bookings' && (
            <div className="bookings-section">
              {bookingSuccess && <div className="success-message">{bookingSuccess}</div>}
              <div className="section-head">
                <h3>Tour đã đặt</h3>
                <button
                  type="button"
                  className="btn-refresh-bookings"
                  onClick={() => setBookingReloadToken((prev) => prev + 1)}
                >
                  Làm mới
                </button>
              </div>
              {renderBookedTours()}
            </div>
          )}
        </section>
      </div>
    </div>
  );
};

export default Profile;
