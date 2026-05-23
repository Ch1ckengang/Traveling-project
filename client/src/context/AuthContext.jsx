import { createContext, useState, useContext, useCallback, useEffect, useRef } from 'react';
import axios from 'axios';


const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/v1/api';

const extractTokensFromAuthResponse = (payload) => {
  if (!payload || typeof payload !== 'object') {
    return null;
  }

  let tokens = null;
  
  // Try to extract tokens from nested structure
  if (payload.data && typeof payload.data === 'object' && payload.data.tokens) {
    tokens = payload.data.tokens;
  } else if (payload.tokens) {
    tokens = payload.tokens;
  }
  
  if (!tokens) {
    return null;
  }
  
  // Normalize token keys: Backend uses snake_case (access_token, refresh_token)
  return {
    access_token: tokens.AccessToken || tokens.access_token || tokens.accessToken,
    refresh_token: tokens.RefreshToken || tokens.refresh_token || tokens.refreshToken,
    token_type: tokens.TokenType || tokens.token_type || tokens.tokenType || 'Bearer',
    expires_in: tokens.ExpiresIn || tokens.expires_in || tokens.expiresIn
  };
};

// Tạo Context cho quản lý authentication toàn cục
const AuthContext = createContext(null);

/**
 * AuthProvider - Component bao bọc toàn bộ app để cung cấp auth state
 * Quản lý: user info, trạng thái đăng nhập, login/logout functions
 * Lưu trữ: sử dụng localStorage để duy trì session khi refresh trang
 */
export const AuthProvider = ({ children }) => {
  // State lưu thông tin user (id, name, email)
  const [user, setUser] = useState(() => {
    const savedUser = localStorage.getItem('user');
    return savedUser ? JSON.parse(savedUser) : null;
  });

  const [tokens, setTokens] = useState(() => {
    const savedTokens = localStorage.getItem('auth_tokens');
    return savedTokens ? JSON.parse(savedTokens) : null;
  });

  // Guard against re-entrant logout calls
  const isLoggingOut = useRef(false);

  // Trạng thái đăng nhập được suy ra trực tiếp từ user VÀ tokens
  const isLoggedIn = Boolean(user && tokens?.access_token);

  // Role-based computed values
  const userRole = user?.role || 'customer';
  const isAdmin = userRole === 'admin';
  const isStaff = userRole === 'staff' || userRole === 'admin';

  /**
   * login - Hàm xử lý khi user đăng nhập thành công
   * @param {Object} userData - Thông tin user từ API {id, name, email}
   * Cập nhật state và lưu vào localStorage
   */
  const login = (userData, tokenData) => {
    isLoggingOut.current = false;
    setUser(userData);
    setTokens(tokenData || null);
    
    // Lưu vào localStorage để giữ session khi refresh
    localStorage.setItem('user', JSON.stringify(userData));
    
    if (tokenData) {
      localStorage.setItem('auth_tokens', JSON.stringify(tokenData));
      
      // Sync với axiosInstance format
      if (tokenData.access_token) {
        localStorage.setItem('accessToken', tokenData.access_token);
      }
      if (tokenData.refresh_token) {
        localStorage.setItem('refreshToken', tokenData.refresh_token);
      }
    } else {
      localStorage.removeItem('auth_tokens');
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
    }
  };

  /**
   * logout - Hàm xử lý khi user đăng xuất
   * Xoá state và localStorage
   */
  const logout = useCallback(() => {
    if (isLoggingOut.current) return; // Prevent re-entrant calls
    isLoggingOut.current = true;
    
    setUser(null);
    setTokens(null);
    localStorage.removeItem('user');
    localStorage.removeItem('auth_tokens');
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
  }, []);

  // Listen for auth:logout event from axiosInstance's forceLogout()
  // This syncs AuthContext state when axiosInstance detects token expiry
  useEffect(() => {
    const handleForceLogout = () => {
      logout();
    };

    window.addEventListener('auth:logout', handleForceLogout);
    return () => window.removeEventListener('auth:logout', handleForceLogout);
  }, [logout]);

  const updateTokens = (tokenData) => {
    setTokens(tokenData || null);
    if (tokenData) {
      localStorage.setItem('auth_tokens', JSON.stringify(tokenData));
      
      // Sync với axiosInstance format
      if (tokenData.access_token) {
        localStorage.setItem('accessToken', tokenData.access_token);
      }
      if (tokenData.refresh_token) {
        localStorage.setItem('refreshToken', tokenData.refresh_token);
      }
      return;
    }

    localStorage.removeItem('auth_tokens');
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
  };

  const getAccessToken = () => tokens?.access_token || '';

  const refreshAccessToken = async () => {
    const refreshToken = tokens?.refresh_token;
    if (!refreshToken) {
      logout();
      throw new Error('Không tìm thấy refresh token.');
    }

    try {
      // Backend expects "refreshToken" (camelCase) in request body
      const response = await axios.post(`${API_BASE_URL}/token/refresh`, {
        refreshToken: refreshToken
      });

      const nextTokens = extractTokensFromAuthResponse(response.data);
      if (!nextTokens?.access_token || !nextTokens?.refresh_token) {
        logout();
        throw new Error('Phản hồi refresh token không hợp lệ.');
      }

      updateTokens(nextTokens);
      return nextTokens.access_token;
    } catch (error) {
      logout();
      throw error;
    }
  };

  const requestWithAuth = async (requestFactory) => {
    if (isLoggingOut.current) {
      throw new Error('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.');
    }

    const firstAccessToken = getAccessToken();
    if (!firstAccessToken) {
      throw new Error('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.');
    }

    try {
      return await requestFactory(firstAccessToken);
    } catch (error) {
      const status = error?.response?.status;
      if (status !== 401) {
        throw error;
      }

      const nextAccessToken = await refreshAccessToken();
      return requestFactory(nextAccessToken);
    }
  };

  // fetchUser - Gọi GET /users/me để lấy thông tin user mới nhất từ server
  const fetchUser = useCallback(async () => {
    const accessToken = getAccessToken();
    if (!accessToken || isLoggingOut.current) return null;
    try {
      const response = await axios.get(`${API_BASE_URL}/users/me`, {
        headers: { Authorization: `Bearer ${accessToken}` }
      });
      const userData = response.data?.data?.user || response.data?.user;
      if (userData) {
        setUser(userData);
        localStorage.setItem('user', JSON.stringify(userData));
        return userData;
      }
    } catch (err) {
      // If 401, token is expired — logout to stop retry loops
      if (err?.response?.status === 401) {
        logout();
      }
    }
    return null;
  }, [tokens?.access_token]);

  useEffect(() => {
    if (tokens?.access_token && !isLoggingOut.current) {
      fetchUser();
    }
  }, [tokens?.access_token, fetchUser]);

  // updateUser - Cập nhật thông tin user trong context và localStorage
  const updateUser = useCallback((userData) => {
    if (!userData) return;
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData));
  }, []);


  // Cung cấp state và functions cho các component con
  return (
    <AuthContext.Provider value={{
      user,
      tokens,
      isLoggedIn,
      userRole,
      isAdmin,
      isStaff,
      login,
      logout,
      updateTokens,
      updateUser,
      fetchUser,
      getAccessToken,
      refreshAccessToken,
      requestWithAuth
    }}>
      {children}
    </AuthContext.Provider>
  );

};

/**
 * useAuth - Custom hook để sử dụng auth context
 * Sử dụng: const { user, tokens, isLoggedIn, login, logout } = useAuth();
 * @returns {Object} - {user, tokens, isLoggedIn, login, logout, updateTokens, getAccessToken, refreshAccessToken, requestWithAuth}
 */
// eslint-disable-next-line react-refresh/only-export-components
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth phải được sử dụng trong AuthProvider');
  }
  return context;
};

