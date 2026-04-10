import { Link, useLocation } from 'react-router-dom';
import AuthLayout from './AuthLayout';
import '../../styles/Auth.css';

const ResetSuccess = () => {
  const location = useLocation();
  const fromForgotPassword = Boolean(location.state?.fromForgotPassword);
  const serverMessage = location.state?.message || '';

  const description = serverMessage || (fromForgotPassword
    ? 'Link đặt lại đã được gửi. Sau khi đổi mật khẩu thành công, tài khoản của bạn sẽ được bảo mật tốt hơn.'
    : 'Tài khoản của bạn đã được bảo mật với mật khẩu mới. Hãy đăng nhập để tiếp tục hành trình.');

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
