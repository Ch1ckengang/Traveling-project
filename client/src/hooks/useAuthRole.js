import { useMemo } from 'react';
import { useAuth } from '../context/AuthContext';

const useAuthRole = () => {
  const { user, isLoggedIn } = useAuth();

  return useMemo(() => ({
    isLoggedIn,
    role: (user?.role || 'customer').toString().toLowerCase()
  }), [isLoggedIn, user?.role]);
};

export default useAuthRole;
