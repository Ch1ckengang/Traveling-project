import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import axios from 'axios';
import AuthLayout from './AuthLayout';
import '../../styles/Auth.css';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/v1/api';

const ForgotPassword = () => {
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');

    if (!email.trim()) {
      setError('Vui lòng nhập email hợp lệ.');
      return;
    }

    setLoading(true);

    try {
      const response = await axios.post(`${API_BASE_URL}/password/forgot`, {
        email
      });

      navigate('/reset-success', {
        state: {
          fromForgotPassword: true,
          email,
          message: response.data.message
        }
      });
    } catch (requestError) {
      setError(requestError.response?.data?.message || 'Không thể gửi yêu cầu đặt lại.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthLayout title="Đặt lại mật khẩu" subtitle="Nhập email để nhận link đặt lại mật khẩu">
      {error && <div className="error-message">{error}</div>}

      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="forgot-email">Email</label>
          <input
            id="forgot-email"
            type="email"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
            placeholder="name@example.com"
            required
          />
        </div>

        <button type="submit" className="btn-primary" disabled={loading}>
          {loading ? 'Đang gửi...' : 'Gửi link đặt lại'}
        </button>
      </form>

      <Link to="/login" className="btn-outline btn-block auth-secondary-action">
        Quay lại đăng nhập
      </Link>
    </AuthLayout>
  );
};

export default ForgotPassword;
