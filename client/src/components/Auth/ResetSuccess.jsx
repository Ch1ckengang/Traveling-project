import { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import AuthLayout from './AuthLayout';
import { resetPassword } from '../../services/authService';
import '../../styles/Auth.css';

const ResetSuccess = () => {
  const location = useLocation();
  const navigate = useNavigate();
  const fromForgotPassword = Boolean(location.state?.fromForgotPassword);
  const serverMessage = location.state?.message || '';
  const initialEmail = location.state?.email || '';
  const [formData, setFormData] = useState({
    email: initialEmail,
    otpCode: '',
    newPassword: '',
    confirmPassword: ''
  });
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const description = serverMessage || (fromForgotPassword
    ? 'Mã OTP đã được gửi đến email của bạn. Nhập mã để đặt mật khẩu mới.'
    : 'Tài khoản của bạn đã được bảo mật với mật khẩu mới. Hãy đăng nhập để tiếp tục hành trình.');

  const handleChange = (event) => {
    setFormData((current) => ({
      ...current,
      [event.target.name]: event.target.value
    }));
    setError('');
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');

    if (!formData.email.trim() || !formData.otpCode.trim()) {
      setError('Vui lòng nhập email và mã OTP.');
      return;
    }

    if (formData.newPassword.length < 8) {
      setError('Mật khẩu mới phải có ít nhất 8 ký tự.');
      return;
    }

    if (formData.newPassword !== formData.confirmPassword) {
      setError('Mật khẩu xác nhận không khớp.');
      return;
    }

    setLoading(true);
    try {
      await resetPassword(formData.email.trim(), formData.otpCode.trim(), formData.newPassword);
      navigate('/reset-password', {
        state: {
          message: 'Đặt lại mật khẩu thành công. Hãy đăng nhập bằng mật khẩu mới.'
        }
      });
    } catch (requestError) {
      setError(requestError.response?.data?.message || 'Không thể đặt lại mật khẩu. Vui lòng kiểm tra mã OTP.');
    } finally {
      setLoading(false);
    }
  };

  if (fromForgotPassword) {
    return (
      <AuthLayout title="Đặt lại mật khẩu" subtitle={description}>
        {error && <div className="error-message">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="reset-email">Email</label>
            <input
              id="reset-email"
              name="email"
              type="email"
              value={formData.email}
              onChange={handleChange}
              placeholder="name@example.com"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="reset-otp">Mã OTP</label>
            <input
              id="reset-otp"
              name="otpCode"
              type="text"
              value={formData.otpCode}
              onChange={handleChange}
              placeholder="Nhập mã 6 số"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="reset-new-password">Mật khẩu mới</label>
            <input
              id="reset-new-password"
              name="newPassword"
              type="password"
              value={formData.newPassword}
              onChange={handleChange}
              autoComplete="new-password"
              required
            />
          </div>

          <div className="form-group">
            <label htmlFor="reset-confirm-password">Xác nhận mật khẩu</label>
            <input
              id="reset-confirm-password"
              name="confirmPassword"
              type="password"
              value={formData.confirmPassword}
              onChange={handleChange}
              autoComplete="new-password"
              required
            />
          </div>

          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Đang đặt lại...' : 'Đặt lại mật khẩu'}
          </button>
        </form>
      </AuthLayout>
    );
  }

  return (
    <AuthLayout title="Mật khẩu đã được đặt lại!" subtitle={description}>
      <div className="reset-success-icon" aria-hidden="true">
        <svg width="30" height="30" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path d="M20 7L9 18L4 13" stroke="#2F7A2E" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </div>

      <Link to="/login" className="btn-primary btn-block">
        Đăng nhập ngay
      </Link>
    </AuthLayout>
  );
};

export default ResetSuccess;
