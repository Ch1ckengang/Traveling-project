import { createContext, useContext, useEffect, useMemo, useReducer } from 'react';

const AUTH_STORAGE_KEY = 'accessToken';
const REFRESH_STORAGE_KEY = 'refreshToken';
const LOGOUT_EVENT_NAME = 'auth:logout';

const initialState = {
  user: null,
  accessToken: null,
  refreshToken: null,
  isAuthenticated: false,
  isLoading: false
};

const AuthContext = createContext(null);

function authReducer(state, action) {
  switch (action.type) {
    case 'LOGIN': {
      const nextAccessToken = action.payload?.accessToken || null;
      const nextRefreshToken = action.payload?.refreshToken || null;

      if (nextAccessToken) {
        localStorage.setItem(AUTH_STORAGE_KEY, nextAccessToken);
      } else {
        localStorage.removeItem(AUTH_STORAGE_KEY);
      }

      if (nextRefreshToken) {
        localStorage.setItem(REFRESH_STORAGE_KEY, nextRefreshToken);
      }

      return {
        ...state,
        user: action.payload?.user || null,
        accessToken: nextAccessToken,
        refreshToken: nextRefreshToken,
        isAuthenticated: Boolean(nextAccessToken),
        isLoading: false
      };
    }

    case 'LOGOUT': {
      localStorage.removeItem(AUTH_STORAGE_KEY);
      localStorage.removeItem(REFRESH_STORAGE_KEY);

      return {
        ...initialState
      };
    }

    case 'UPDATE_USER': {
      return {
        ...state,
        user: action.payload || null
      };
    }

    case 'SET_LOADING': {
      return {
        ...state,
        isLoading: Boolean(action.payload)
      };
    }

    default:
      return state;
  }
}

export function AuthProvider({ children }) {
  const [state, dispatch] = useReducer(authReducer, {
    ...initialState,
    accessToken: localStorage.getItem(AUTH_STORAGE_KEY),
    refreshToken: localStorage.getItem(REFRESH_STORAGE_KEY),
    isAuthenticated: Boolean(localStorage.getItem(AUTH_STORAGE_KEY))
  });

  useEffect(() => {
    const handleGlobalLogout = () => {
      dispatch({ type: 'LOGOUT' });
    };

    window.addEventListener(LOGOUT_EVENT_NAME, handleGlobalLogout);

    return () => {
      window.removeEventListener(LOGOUT_EVENT_NAME, handleGlobalLogout);
    };
  }, []);

  const value = useMemo(
    () => ({
      state,
      dispatch
    }),
    [state]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within AuthProvider');
  }

  return context;
}
