import { Navigate, Outlet } from 'react-router-dom';

export function ProtectedRoute({ isAuthenticated }) {
  if (!isAuthenticated) {
    return <Navigate to='/auth/login' replace />;
  }

  return <Outlet />;
}

export function RoleRoute({ role, allowedRoles }) {
  if (!allowedRoles.includes(role)) {
    return <Navigate to='/' replace />;
  }

  return <Outlet />;
}
