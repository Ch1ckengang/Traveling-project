import { useEffect, useMemo, useRef, useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import axios from 'axios';
import AuthLayout from './AuthLayout';
import '../../styles/Auth.css';

const OTP_LENGTH = 6;
const OTP_EXPIRE_SECONDS = 180;
const API_BASE_URL = 'http://localhost:8080/v1/api';

const maskEmail = (email) => {
  if (!email || !email.includes('@')) {
    return 'n***@email.com';
  }

  const [name, domain] = email.split('@');
  if (!name) {
    return `n***@${domain}`;
  }

  return `${name[0]}***@${domain}`;
};

const formatCountdown = (seconds) => {
  const minutes = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${String(minutes).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
};

const OtpVerification = () => {
  const location = useLocation();
  const navigate = useNavigate();

  const initialEmail = location.state?.email || sessionStorage.getItem('pending_verification_email') || '';
  const infoMessage = location.state?.message || '';
  const [digits, setDigits] = useState(Array(OTP_LENGTH).fill(''));
  const [statusMessage, setStatusMessage] = useState(infoMessage);
  const [error, setError] = useState('');
  const [verifying, setVerifying] = useState(false);
  const [resending, setResending] = useState(false);
  const [countdown, setCountdown] = useState(OTP_EXPIRE_SECONDS);
  const inputRefs = useRef([]);

  const obfuscatedEmail = useMemo(() => maskEmail(initialEmail), [initialEmail]);

  useEffect(() => {
    inputRefs.current[0]?.focus();
  }, []);

  useEffect(() => {
    if (countdown <= 0) {
      return undefined;
    }

    const timer = setInterval(() => {
      setCountdown((prev) => prev - 1);
    }, 1000);

    return () => clearInterval(timer);
  }, [countdown]);

  const handleChange = (index, value) => {
    if (!/^\d?$/.test(value)) {
      return;
    }

    const nextDigits = [...digits];
    nextDigits[index] = value;
    setDigits(nextDigits);

    if (value && index < OTP_LENGTH - 1) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handleKeyDown = (index, event) => {
    if (event.key === 'Backspace' && !digits[index] && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
  };

  const handlePaste = (event) => {
    const pasted = event.clipboardData.getData('text').trim();
    if (!/^\d{6}$/.test(pasted)) {
      return;
    }

    const nextDigits = pasted.split('').slice(0, OTP_LENGTH);
    setDigits(nextDigits);
    inputRefs.current[OTP_LENGTH - 1]?.focus();
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setStatusMessage('');

    const otpCode = digits.join('');
    if (otpCode.length < OTP_LENGTH) {
      setError('Vui lòng nhập đủ 6 chữ số.');
      return;
    }

    if (!initialEmail) {
      setError('Không tìm thấy email để xác thực. Vui lòng đăng ký lại.');
      return;
    }

    setVerifying(true);

    try {
      const response = await axios.post(`${API_BASE_URL}/otp/verify`, {
        email: initialEmail,
        code: otpCode
      });

      if (response.data.success) {
        sessionStorage.removeItem('pending_verification_email');
        navigate('/login', {
          state: {
            message: response.data.message || 'Xác thực thành công. Bạn có thể đăng nhập ngay.'
          }
        });
      }
    } catch (requestError) {
      setError(requestError.response?.data?.message || 'Không thể xác thực OTP. Vui lòng thử lại.');
    } finally {
      setVerifying(false);
    }
  };

  const handleResend = async () => {
    if (!initialEmail) {
      setError('Không tìm thấy email để gửi lại mã.');
      return;
    }

    setResending(true);

    try {
      const response = await axios.post(`${API_BASE_URL}/otp/send`, {
        email: initialEmail
      });

      setStatusMessage(response.data.message || 'Mã mới đã được gửi.');
      setDigits(Array(OTP_LENGTH).fill(''));
      setCountdown(OTP_EXPIRE_SECONDS);
      setError('');
      inputRefs.current[0]?.focus();
    } catch (requestError) {
      setError(requestError.response?.data?.message || 'Không thể gửi lại mã OTP.');
    } finally {
      setResending(false);
    }
  };

  return (
    <AuthLayout title="Nhập mã xác thực" subtitle={`Mã 6 chữ số đã gửi đến ${obfuscatedEmail}`}>
      <div className="auth-step-indicator" aria-hidden="true">
        <span />
        <span className="active" />
        <span />
      </div>

      {error && <div className="error-message">{error}</div>}
      {statusMessage && <div className="success-message">{statusMessage}</div>}

      <form onSubmit={handleSubmit}>
        <div className="otp-row" onPaste={handlePaste}>
          {digits.map((digit, index) => (
            <input
              key={`otp-${index + 1}`}
              className="otp-input"
              inputMode="numeric"
              maxLength="1"
              value={digit}
              onChange={(event) => handleChange(index, event.target.value)}
              onKeyDown={(event) => handleKeyDown(index, event)}
              ref={(element) => {
                inputRefs.current[index] = element;
              }}
              aria-label={`Mã OTP số ${index + 1}`}
              required
            />
          ))}
        </div>

        <p className="otp-countdown">Hết hạn sau: <strong>{formatCountdown(Math.max(countdown, 0))}</strong></p>

        <button type="submit" className="btn-primary" disabled={verifying}>
          {verifying ? 'Đang xác thực...' : 'Xác nhận'}
        </button>
      </form>

      <p className="auth-inline-help">
        Không nhận được mã?{' '}
        <button type="button" className="auth-link-btn" onClick={handleResend} disabled={resending}>
          {resending ? 'Đang gửi...' : 'Gửi lại'}
        </button>
      </p>

      <div className="auth-footer-link">
        <Link to="/login">Quay lại đăng nhập</Link>
      </div>
    </AuthLayout>
  );
};

export default OtpVerification;
