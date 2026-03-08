import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import axios from 'axios';
import '../../styles/Profile.css';

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
  const navigate = useNavigate();
  
  // State quản lý form data
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  
  const [isEditing, setIsEditing] = useState(false); // Trạng thái chế độ xem/chỉnh sửa
  const [error, setError] = useState(''); // Thông báo lỗi
  const [success, setSuccess] = useState(''); // Thông báo thành công
  const [loading, setLoading] = useState(false); // Trạng thái đang xử lý

  // Load thông tin user vào form khi component mount hoặc user thay đổi
  useEffect(() => {
    // Nếu chưa đăng nhập -> chuyển về trang login
    if (!user) {
      navigate('/login');
      return;
    }
    
    // Set giá trị mặc định cho form từ thông tin user hiện tại
    setFormData({
      name: user.name || '',
      email: user.email || '',
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    });
  }, [user, navigate]);

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

    // VALIDATION
    if (formData.name.trim() === '') {
      setError('Tên không được để trống');
      return;
    }

    if (formData.email.trim() === '') {
      setError('Email không được để trống');
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
        name: formData.name,
        email: formData.email,
      };

      // Thêm password nếu user muốn đổi
      if (formData.newPassword) {
        updateData.password = formData.newPassword;
      }

      // GỬI REQUEST ĐẾN BACKEND
      const response = await axios.put(
        `http://localhost:8080/api/users/${user.id}`,
        updateData
      );

      if (response.data.success) {
        // Cập nhật context với thông tin mới (tự động update header)
        login(response.data.user);
        
        setSuccess('Cập nhật thông tin thành công!');
        setIsEditing(false); // Chuyển về chế độ xem
        
        // Reset password fields
        setFormData({
          ...formData,
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
      name: user.name || '',
      email: user.email || '',
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

  return (
    <div className="profile-container">
      <div className="profile-box">
        <div className="profile-header">
          <div className="profile-avatar-large">
            {user.name?.charAt(0).toUpperCase()}
          </div>
          <h2>{user.name}</h2>
          <p className="profile-email">{user.email}</p>
        </div>

        {error && <div className="error-message">{error}</div>}
        {success && <div className="success-message">{success}</div>}

        <form onSubmit={handleSubmit} className="profile-form">
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
      </div>
    </div>
  );
};

export default Profile;
