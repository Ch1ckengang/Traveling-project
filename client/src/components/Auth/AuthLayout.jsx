const AuthLayout = ({ title, subtitle, children }) => {
  return (
    <div className="auth-shell">
      <div className="auth-grid">
        <aside className="auth-side-panel" aria-hidden="true">
          <p className="auth-side-label">Traveling System</p>
          <h1>Bao trọn hành trình, an tâm từng bước</h1>
          <p>
            Hệ thống đặt tour và quản lý tài khoản tập trung, giúp bạn xác thực nhanh,
            bảo mật ổn định và theo dõi chuyến đi liền mạch.
          </p>
        </aside>

        <section className="auth-card" aria-label={title}>
          <header className="auth-card-header">
            <h2>{title}</h2>
            {subtitle && <p>{subtitle}</p>}
          </header>
          {children}
        </section>
      </div>
    </div>
  );
};

export default AuthLayout;
