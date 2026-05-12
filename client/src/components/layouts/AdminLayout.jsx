import { useMemo, useState } from 'react';
import { Navigate, NavLink, Outlet, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';

const adminItems = [
  { label: 'Dashboard', icon: '📊', to: '/admin/dashboard', title: 'Dashboard' },
  { type: 'group', label: 'Quản lý' },
  { label: 'Tour', icon: '🗺️', to: '/admin/tours', title: 'Tour' },
  { label: 'Lịch tour', icon: '📅', to: '/admin/schedules', title: 'Lịch tour' },
  { label: 'Đặt chỗ', icon: '📋', to: '/admin/bookings', title: 'Đặt chỗ' },
  { label: 'Thanh toán', icon: '💳', to: '/admin/payments', title: 'Thanh toán' },
  { type: 'group', label: 'Khách hàng' },
  { label: 'Người dùng', icon: '👥', to: '/admin/users', title: 'Người dùng' },
  { label: 'Đánh giá', icon: '⭐', to: '/admin/reviews', title: 'Đánh giá' },
  { type: 'group', label: 'Marketing' },
  { label: 'Mã giảm giá', icon: '🎫', to: '/admin/coupons', title: 'Mã giảm giá' },
  { type: 'group', label: 'Hệ thống' },
  { label: 'Báo cáo', icon: '📈', to: '/admin/reports', title: 'Báo cáo' }
];

const navLinkClass = (collapsed) => ({ isActive }) =>
  `flex items-center ${collapsed ? 'justify-center' : 'gap-3'} rounded-button px-3 py-2 text-sm transition-colors ${
    isActive ? 'bg-primary-500 text-white' : 'text-slate-200 hover:bg-white/10'
  }`;

const AdminLayout = () => {
  const [collapsed, setCollapsed] = useState(false);
  const [mobileSheetOpen, setMobileSheetOpen] = useState(false);
  const { isLoggedIn, user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  const role = (user?.role || '').toUpperCase();
  const canAccess = isLoggedIn && (role === 'ADMIN' || role === 'STAFF');

  if (!canAccess) {
    return <Navigate to='/' replace />;
  }

  const title = useMemo(() => {
    const matchedItem = adminItems.find((item) => item.to && location.pathname.startsWith(item.to));
    return matchedItem?.title || 'Quản trị hệ thống';
  }, [location.pathname]);

  const handleLogout = () => {
    logout();
    navigate('/login', { replace: true });
  };

  const renderMenu = (isCollapsed) => (
    <nav className='space-y-1'>
      {adminItems.map((item) => {
        if (item.type === 'group') {
          return (
            <p
              key={item.label}
              className={`px-2 pt-3 pb-1 text-[11px] uppercase tracking-wide text-slate-400 ${isCollapsed ? 'text-center' : ''}`}
            >
              {isCollapsed ? '---' : `--- ${item.label} ---`}
            </p>
          );
        }

        return (
          <NavLink key={item.to} to={item.to} className={navLinkClass(isCollapsed)} title={item.label}>
            <span aria-hidden='true'>{item.icon}</span>
            {!isCollapsed && <span>{item.label}</span>}
          </NavLink>
        );
      })}
    </nav>
  );

  return (
    <div className='min-h-screen bg-slate-100 md:pb-0 pb-16'>
      <div className='flex min-h-screen'>
        <aside
          className={`hidden md:flex shrink-0 flex-col bg-slate-900 text-white border-r border-slate-800 transition-all duration-200 ${
            collapsed ? 'w-[60px]' : 'w-[240px]'
          }`}
        >
          <div className='h-14 flex items-center justify-between px-3 border-b border-slate-800'>
            {!collapsed && <p className='font-semibold text-sm'>Admin Panel</p>}
            <button
              type='button'
              onClick={() => setCollapsed((prev) => !prev)}
              className='h-8 w-8 rounded-md bg-slate-800 hover:bg-slate-700'
              aria-label='Toggle sidebar'
            >
              {collapsed ? '»' : '«'}
            </button>
          </div>

          <div className='flex-1 overflow-y-auto p-2'>{renderMenu(collapsed)}</div>
        </aside>

        <div className='min-w-0 flex-1 flex flex-col'>
          <header className='h-14 bg-white border-b border-slate-200 px-4 md:px-6 flex items-center justify-between sticky top-0 z-30'>
            <div className='flex items-center gap-2'>
              <button
                type='button'
                onClick={() => setMobileSheetOpen(true)}
                className='md:hidden h-9 px-3 rounded-md border border-slate-300 text-sm'
              >
                Menu
              </button>
              <h1 className='text-base md:text-lg font-semibold text-slate-900'>{title}</h1>
            </div>

            <div className='flex items-center gap-3'>
              <p className='text-sm text-slate-600 hidden sm:block'>{user?.name || 'Admin'}</p>
              <button
                type='button'
                onClick={handleLogout}
                className='rounded-button border border-rose-200 px-3 py-1.5 text-sm text-rose-600 hover:bg-rose-50'
              >
                Đăng xuất
              </button>
            </div>
          </header>

          <main className='flex-1 p-4 md:p-6'>
            <div className='rounded-card border border-slate-200 bg-white p-4 md:p-6 min-h-[calc(100vh-120px)]'>
              <Outlet />
            </div>
          </main>
        </div>
      </div>

      {mobileSheetOpen && (
        <div className='md:hidden fixed inset-0 z-50'>
          <button
            type='button'
            className='absolute inset-0 bg-black/40'
            onClick={() => setMobileSheetOpen(false)}
            aria-label='Close menu overlay'
          />

          <div className='absolute inset-x-0 bottom-0 rounded-t-2xl bg-slate-900 text-white p-4 max-h-[72vh] overflow-y-auto'>
            <div className='mb-3 flex items-center justify-between'>
              <p className='text-sm font-semibold'>Điều hướng quản trị</p>
              <button
                type='button'
                className='rounded-md bg-slate-800 px-2 py-1 text-xs'
                onClick={() => setMobileSheetOpen(false)}
              >
                Đóng
              </button>
            </div>

            {renderMenu(false)}
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminLayout;
