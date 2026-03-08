import { createContext, useState, useContext, useEffect } from 'react';

// Tạo Context cho quản lý authentication toàn cục
const AuthContext = createContext(null);

/**
 * AuthProvider - Component bao bọc toàn bộ app để cung cấp auth state
 * Quản lý: user info, trạng thái đăng nhập, login/logout functions
 * Lưu trữ: sử dụng localStorage để duy trì session khi refresh trang
 */
export const AuthProvider = ({ children }) => {
  // State lưu thông tin user (id, name, email)
  const [user, setUser] = useState(null);
  // State đánh dấu user đã đăng nhập hay chưa
  const [isLoggedIn, setIsLoggedIn] = useState(false);

  // Kiểm tra localStorage khi component mount (load trang)
  // Nếu có thông tin user đã lưu -> khôi phục trạng thái đăng nhập
  useEffect(() => {
    const savedUser = localStorage.getItem('user');
    if (savedUser) {
      const userData = JSON.parse(savedUser);
      setUser(userData);
      setIsLoggedIn(true);
    }
  }, []);

  /**
   * login - Hàm xử lý khi user đăng nhập thành công
   * @param {Object} userData - Thông tin user từ API {id, name, email}
   * Cập nhật state và lưu vào localStorage
   */
  const login = (userData) => {
    setUser(userData);
    setIsLoggedIn(true);
    // Lưu vào localStorage để giữ session khi refresh
    localStorage.setItem('user', JSON.stringify(userData));
  };

  /**
   * logout - Hàm xử lý khi user đăng xuất
   * Xoá state và localStorage
   */
  const logout = () => {
    setUser(null);
    setIsLoggedIn(false);
    localStorage.removeItem('user');
  };

  // Cung cấp state và functions cho các component con
  return (
    <AuthContext.Provider value={{ user, isLoggedIn, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
};

/**
 * useAuth - Custom hook để sử dụng auth context
 * Sử dụng: const { user, isLoggedIn, login, logout } = useAuth();
 * @returns {Object} - {user, isLoggedIn, login, logout}
 */
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth phải được sử dụng trong AuthProvider');
  }
  return context;
};
