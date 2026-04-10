import { useState } from 'react';
import { useNavigate, Link, useLocation } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import axios from 'axios';
import AuthLayout from './AuthLayout';
import '../../styles/Auth.css';

const API_BASE_URL = 'http://localhost:8080/v1/api';

const Login = () => {
  const location = useLocation();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [rememberMe, setRememberMe] = useState(true);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const successMessage = location.state?.message || '';
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    const normalizedEmail = email.trim().toLowerCase();
    const normalizedPassword = password.trim();

    if (!normalizedEmail || !normalizedPassword) {
      setError('Vui lòng nhập đầy đủ email và mật khẩu.');
      return;
    }

    setLoading(true);

    try {
      const response = await axios.post(`${API_BASE_URL}/login`, {
        email: normalizedEmail,
        password: normalizedPassword
      });

      if (response.data.success) {
        login(response.data.user);
        if (!rememberMe) {
          localStorage.removeItem('user');
        }
        navigate('/');
      }
    } catch (err) {
      setError(err.response?.data?.message || 'Đăng nhập thất bại. Vui lòng thử lại.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout title="Chào mừng trở lại">
      {successMessage && <div className="success-message">{successMessage}</div>}
      {error && <div className="error-message">{error}</div>}

      <div className="social-row">
        <button type="button" className="social-btn" aria-label="Đăng nhập với Google">
          <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
            <path fill="#EA4335" d="M12 10.2v3.9h5.5c-.2 1.3-1.5 3.9-5.5 3.9-3.3 0-6-2.8-6-6.2s2.7-6.2 6-6.2c1.9 0 3.1.8 3.8 1.5l2.6-2.6C16.7 2.9 14.6 2 12 2 6.9 2 2.8 6.5 2.8 12s4.1 10 9.2 10c5.3 0 8.8-3.8 8.8-9.1 0-.6-.1-1.1-.2-1.6H12z" />
          </svg>
          Google
        </button>
        <button type="button" className="social-btn" aria-label="Đăng nhập với Facebook">
          <svg width="18" height="18" viewBox="0 0 24 24" aria-hidden="true">
            <path fill="#1877F2" d="M24 12.1C24 5.4 18.6 0 12 0S0 5.4 0 12.1C0 18.1 4.4 23 10.1 24v-8.5H7.1v-3.5h3V9.4c0-3 1.8-4.8 4.6-4.8 1.3 0 2.7.2 2.7.2v3h-1.5c-1.5 0-2 1-2 2v2.3h3.4l-.5 3.5h-2.9V24C19.6 23 24 18.1 24 12.1z" />
          </svg>
          Facebook
        </button>
      </div>

      <div className="auth-divider"><span>hoặc</span></div>

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="login-email">Email</label>
          <input
            id="login-email"
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="name@example.com"
            required
          />
        </div>

        <div className="form-group">
          <label htmlFor="login-password">Mật khẩu</label>
          <div className="input-with-action">
            <input
              id="login-password"
              type={showPassword ? 'text' : 'password'}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Nhập mật khẩu"
              required
            />
            <button
              type="button"
              className="icon-btn"
              onClick={() => setShowPassword((prev) => !prev)}
              aria-label={showPassword ? 'Ẩn mật khẩu' : 'Hiện mật khẩu'}
            >
              {showPassword ? (
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path d="M3 3L21 21" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
                  <path d="M10.6 10.6A2 2 0 0013.4 13.4" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
                  <path d="M9 5.2A10.8 10.8 0 0112 4.8c5.5 0 9.2 4.7 10 6.1-.5.8-1.9 2.9-4.1 4.5" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
                  <path d="M6.3 7.3C4.3 8.8 2.9 10.9 2 12.3c.8 1.4 4.5 6.1 10 6.1 1.1 0 2.1-.2 3.1-.5" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round" />
                </svg>
              ) : (
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
                  <path d="M2 12.2c.8-1.4 4.5-6.2 10-6.2s9.2 4.8 10 6.2c-.8 1.4-4.5 6.2-10 6.2s-9.2-4.8-10-6.2z" stroke="currentColor" strokeWidth="1.8" />
                  <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="1.8" />
                </svg>
              )}
            </button>
          </div>
        </div>

        <div className="auth-form-row">
          <label className="checkbox-wrap">
            <input
              type="checkbox"
              checked={rememberMe}
              onChange={(event) => setRememberMe(event.target.checked)}
            />
            Ghi nhớ đăng nhập
          </label>

          <Link className="inline-link" to="/forgot-password">
            Quên mật khẩu?
          </Link>
        </div>

        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? 'Đang đăng nhập...' : 'Đăng nhập'}
        </button>
      </form>

      <Link to="/register" className="btn-outline btn-block auth-secondary-action">
        Chưa có tài khoản? Đăng ký
      </Link>
    </AuthLayout>
  );
};

export default Login;
