import { Link, NavLink, useNavigate } from 'react-router-dom';
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
        <Link to="/" className="brand-link" aria-label="Traveling home">
          <span className="brand-mark">T</span>
          <span className="brand-copy">
            <strong>Traveling</strong>
            <small>Explore with calm</small>
          </span>
        </Link>

        <nav className="nav-links">
          <NavLink to="/" className="nav-item">Trang chủ</NavLink>
          <NavLink to="/tours?category=domestic" className="nav-item">Việt Nam</NavLink>
          <NavLink to="/tours?category=international" className="nav-item">Quốc tế</NavLink>
          <NavLink to="/tours?category=service" className="nav-item">Dịch vụ</NavLink>
        </nav>

        <div className="auth-area">
          {isLoggedIn ? (
            <div className="user-profile">
              <Link to="/profile" className="avatar-link">
                <div className="avatar">{user?.name?.charAt(0).toUpperCase()}</div>
              </Link>
              <span className="user-text">{user?.name}</span>
              <button type="button" onClick={handleLogout} className="btn-logout">Đăng xuất</button>
            </div>
          ) : (
            <div className="guest-actions">
              <span className="avatar guest-avatar">U</span>
              <span className="user-text">Khách</span>
              <Link to="/login" className="auth-link">Đăng nhập</Link>
              <Link to="/register" className="auth-link">Đăng ký</Link>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header;