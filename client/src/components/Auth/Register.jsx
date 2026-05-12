import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import axios from 'axios';
import AuthLayout from './AuthLayout';
import '../../styles/Auth.css';

const API_BASE_URL = 'http://localhost:8080/v1/api';

const getPasswordStrength = (value) => {
  const checks = [
    value.length >= 8,
    /[A-Z]/.test(value),
    /\d/.test(value),
    /[^A-Za-z0-9]/.test(value)
  ];

  const score = checks.filter(Boolean).length;
  if (score <= 1) {
    return { score, label: 'Yếu', tone: 'weak' };
  }

  if (score <= 3) {
    return { score, label: 'Trung bình', tone: 'medium' };
  }

  return { score, label: 'Mạnh', tone: 'strong' };
};

const Register = () => {
  const [formData, setFormData] = useState({
    name: '',
    email: '',
    password: '',
    confirmPassword: ''
  });
  const [acceptedTerms, setAcceptedTerms] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const navigate = useNavigate();
  const passwordStrength = getPasswordStrength(formData.password);

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    const normalizedName = formData.name.trim();
    const normalizedEmail = formData.email.trim().toLowerCase();
    const normalizedPassword = formData.password;

    if (!normalizedName || normalizedName.length < 2) {
      setError('Họ và tên phải có ít nhất 2 ký tự.');
      return;
    }

    if (!normalizedEmail) {
      setError('Vui lòng nhập email hợp lệ.');
      return;
    }

    // Kiểm tra mật khẩu khớp
    if (normalizedPassword !== formData.confirmPassword) {
      setError('Mật khẩu xác nhận không khớp');
      return;
    }

    // Đồng bộ rule với backend: tối thiểu 8 ký tự
    if (normalizedPassword.length < 8) {
      setError('Mật khẩu phải có ít nhất 8 ký tự');
      return;
    }

    if (!acceptedTerms) {
      setError('Bạn cần đồng ý Điều khoản sử dụng để tiếp tục.');
      return;
    }

    setLoading(true);

    try {
      const response = await axios.post(`${API_BASE_URL}/register`, {
        name: normalizedName,
        email: normalizedEmail,
        password: normalizedPassword
      });

      if (response.data.success) {
        sessionStorage.setItem('pending_verification_email', normalizedEmail);
        navigate('/otp-verification', {
          state: {
            email: normalizedEmail,
            message: response.data.message || 'Đăng ký thành công. Vui lòng nhập mã OTP đã gửi về email để xác thực.'
          }
        });
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Đăng ký thất bại. Vui lòng thử lại.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout title="Tạo tài khoản">
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="register-name">Họ và tên</label>
          <input
            id="register-name"
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            placeholder="Nhập họ và tên"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="register-email">Email</label>
          <input
            id="register-email"
            type="email"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="name@example.com"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="register-password">Mật khẩu</label>
          <input
            id="register-password"
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            placeholder="Tối thiểu 8 ký tự"
            required
          />

          <div className="password-meter" role="status" aria-live="polite">
            {[1, 2, 3, 4].map((index) => (
              <span
                key={`meter-${index}`}
                className={index <= passwordStrength.score ? passwordStrength.tone : ''}
              />
            ))}
          </div>
          <p className={`password-meter-label ${passwordStrength.tone}`}>Độ mạnh: {passwordStrength.label}</p>
        </div>

        <div className="form-group">
          <label htmlFor="register-confirm-password">Xác nhận mật khẩu</label>
          <input
            id="register-confirm-password"
            type="password"
            name="confirmPassword"
            value={formData.confirmPassword}
            onChange={handleChange}
            placeholder="Nhập lại mật khẩu"
            required
          />
        </div>

        <label className="checkbox-wrap checkbox-wrap-block">
          <input
            type="checkbox"
            checked={acceptedTerms}
            onChange={(event) => setAcceptedTerms(event.target.checked)}
          />
          Tôi đồng ý với <a href="#" className="inline-link">Điều khoản sử dụng</a>
        </label>

        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? 'Đang tạo tài khoản...' : 'Tạo tài khoản'}
        </button>
      </form>

      <div className="auth-footer-link">
        <p>Đã có tài khoản? <Link to="/login">Đăng nhập ngay</Link></p>
      </div>
    </AuthLayout>
  );
};

export default Register;
