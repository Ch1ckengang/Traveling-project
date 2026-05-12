import { Navigate, NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom';
import Navbar from '../common/Navbar';
import { useAuth } from '../../context/AuthContext';

const linkBaseClass =
  'w-full flex items-center gap-3 rounded-button px-3 py-2 text-sm transition-colors border border-transparent';

const desktopNavLinkClass = ({ isActive }) =>
  `${linkBaseClass} ${isActive ? 'bg-primary-500 text-white' : 'text-slate-700 hover:bg-primary-50'}`;

const mobileNavLinkClass = ({ isActive }) =>
  `flex-1 text-center py-2 text-xs ${isActive ? 'text-primary-600 font-semibold' : 'text-slate-500'}`;

const getInitials = (name, email) => {
  const seed = (name || email || '').trim();
  if (!seed) {
    return 'U';
  }

  const parts = seed.split(/\s+/).filter(Boolean);
  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase();
  }

  return `${parts[0][0]}${parts[parts.length - 1][0]}`.toUpperCase();
};

const AccountLayout = () => {
  const { isLoggedIn, user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  if (!isLoggedIn) {
    const from = `${location.pathname}${location.search}${location.hash}`;
    return <Navigate to='/login' replace state={{ from }} />;
  }

  const initials = getInitials(user?.name, user?.email);

  const handleLogout = () => {
    logout();
    navigate('/login', { replace: true });
  };

  return (
    <div className='min-h-screen bg-slate-50 md:pb-0 pb-16'>
      <div className='sticky top-0 z-40'>
        <Navbar />
      </div>

      <div className='mx-auto w-full max-w-7xl px-4 py-6 md:px-6 lg:px-8'>
        <div className='flex gap-6'>
          <aside className='hidden md:flex w-[240px] shrink-0 flex-col rounded-card border border-slate-200 bg-white p-4 shadow-sm'>
            <div className='mb-6 flex items-center gap-3'>
              {user?.avatarUrl ? (
                <img
                  src={user.avatarUrl}
                  alt={user?.name || 'User avatar'}
                  className='h-14 w-14 rounded-full object-cover border border-slate-200'
                />
              ) : (
                <div className='h-14 w-14 rounded-full bg-primary-500 text-white grid place-items-center font-semibold'>
                  {initials}
                </div>
              )}

              <div className='min-w-0'>
                <p className='truncate text-sm font-semibold text-slate-900'>{user?.name || 'Người dùng'}</p>
                <p className='truncate text-xs text-slate-500'>{user?.email || 'user@email.com'}</p>
              </div>
            </div>

            <nav className='space-y-1'>
              <NavLink to='/account/profile' className={desktopNavLinkClass}>
                <span aria-hidden='true'>👤</span>
                Hồ sơ cá nhân
              </NavLink>

              <NavLink to='/account/bookings' className={desktopNavLinkClass}>
                <span aria-hidden='true'>📋</span>
                Đơn đặt tour
              </NavLink>

              <NavLink to='/account/password' className={desktopNavLinkClass}>
                <span aria-hidden='true'>🔒</span>
                Đổi mật khẩu
              </NavLink>

              <button
                type='button'
                onClick={handleLogout}
                className='w-full flex items-center gap-3 rounded-button px-3 py-2 text-sm text-rose-600 hover:bg-rose-50'
              >
                <span aria-hidden='true'>🚪</span>
                Đăng xuất
              </button>
            </nav>
          </aside>

          <main className='min-w-0 flex-1 rounded-card border border-slate-200 bg-white p-4 md:p-6 shadow-sm'>
            <Outlet />
          </main>
        </div>
      </div>

      <nav className='md:hidden fixed bottom-0 inset-x-0 z-50 border-t border-slate-200 bg-white/95 backdrop-blur px-2 py-1 flex items-center'>
        <NavLink to='/account/profile' className={mobileNavLinkClass}>
          <span className='block text-base' aria-hidden='true'>👤</span>
          Hồ sơ
        </NavLink>

        <NavLink to='/account/bookings' className={mobileNavLinkClass}>
          <span className='block text-base' aria-hidden='true'>📋</span>
          Đơn tour
        </NavLink>

        <NavLink to='/account/password' className={mobileNavLinkClass}>
          <span className='block text-base' aria-hidden='true'>🔒</span>
          Mật khẩu
        </NavLink>

        <button
          type='button'
          onClick={handleLogout}
          className='flex-1 text-center py-2 text-xs text-rose-600'
        >
          <span className='block text-base' aria-hidden='true'>🚪</span>
          Thoát
        </button>
      </nav>
    </div>
  );
};

export default AccountLayout;
