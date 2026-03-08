import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import '../../styles/Header.css';

const Header = () => {
  const { user, isLoggedIn, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  }; 

  return (
    <header className="header-container">
      <div className="header-content">
        {/* Menu chính */}
        <nav className="nav-links">
          <Link to="/" className="nav-item">Trang chủ</Link>
          <Link to="/vietnam" className="nav-item">Du lịch Việt Nam</Link>
          <Link to="/quocte" className="nav-item">Du lịch Quốc tế</Link>
          <Link to="/dichvu" className="nav-item">Dịch vụ</Link>
        </nav>

        {/* Khu vực Đăng nhập / User */}
        <div className="auth-area">
          {isLoggedIn ? (
            <div className="user-profile">
              <Link to="/profile" className="avatar-link">
                <div className="avatar">{user?.name?.charAt(0).toUpperCase()}</div>
              </Link>
              <span>Xin chào, {user?.name}</span>
              <button onClick={handleLogout} className="btn-logout">Đăng xuất</button>
            </div>
          ) : (
            <div className="guest-actions">
              <Link to="/login">
                <button className="btn-login">ĐĂNG NHẬP</button>
              </Link>
              <Link to="/register">
                <button className="btn-register">ĐĂNG KÝ</button>
              </Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;