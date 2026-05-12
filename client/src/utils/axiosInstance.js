import axios from 'axios';

const ACCESS_TOKEN_KEY = 'accessToken';
const REFRESH_TOKEN_KEY = 'refreshToken';
const LOGOUT_EVENT_NAME = 'auth:logout';

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/v1/api'
});

const refreshClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080/v1/api'
});

let isRefreshing = false;
let requestQueue = [];

const processQueue = (error, token = null) => {
  requestQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
      return;
    }

    resolve(token);
  });

  requestQueue = [];
};

const forceLogout = () => {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  window.dispatchEvent(new Event(LOGOUT_EVENT_NAME));
  window.location.href = '/login';
};

axiosInstance.interceptors.request.use(
  (config) => {
    const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
    
    // Debug logging
    console.log('🔐 Axios Request Interceptor:');
    console.log('  URL:', config.url);
    console.log('  Method:', config.method);
    console.log('  AccessToken:', accessToken ? `${accessToken.substring(0, 20)}...` : 'MISSING');
    
    if (accessToken) {
      config.headers = config.headers || {};
      config.headers.Authorization = `Bearer ${accessToken}`;
      console.log('  Authorization header:', config.headers.Authorization.substring(0, 30) + '...');
    } else {
      console.warn('  ⚠️ No access token found!');
    }

    return config;
  },
  (error) => Promise.reject(error)
);

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;
    const status = error.response?.status;

    console.log('❌ Axios Response Error:');
    console.log('  Status:', status);
    console.log('  URL:', originalRequest?.url);
    console.log('  Error:', error.response?.data);

    if (status !== 401 || !originalRequest || originalRequest._retry) {
      return Promise.reject(error);
    }

    console.log('🔄 Attempting token refresh...');
    originalRequest._retry = true;

    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    if (!refreshToken) {
      console.log('❌ No refresh token found, forcing logout');
      forceLogout();
      return Promise.reject(error);
    }

    if (isRefreshing) {
      console.log('⏳ Already refreshing, adding to queue...');
      return new Promise((resolve, reject) => {
        requestQueue.push({ resolve, reject });
      })
        .then((nextAccessToken) => {
          originalRequest.headers = originalRequest.headers || {};
          originalRequest.headers.Authorization = `Bearer ${nextAccessToken}`;
          return axiosInstance(originalRequest);
        })
        .catch((queueError) => Promise.reject(queueError));
    }

    isRefreshing = true;

    try {
      console.log('📤 Sending refresh token request...');
      // Backend expects "refreshToken" (PascalCase) in request body
      const refreshResponse = await refreshClient.post('/token/refresh', {
        refreshToken: refreshToken  // ← Backend expects this key
      });

      // Backend returns PascalCase: AccessToken, RefreshToken
      const nextAccessToken =
        refreshResponse.data?.data?.AccessToken ||
        refreshResponse.data?.AccessToken ||
        refreshResponse.data?.data?.accessToken ||
        refreshResponse.data?.accessToken ||
        null;

      const nextRefreshToken =
        refreshResponse.data?.data?.RefreshToken ||
        refreshResponse.data?.RefreshToken ||
        refreshResponse.data?.data?.refreshToken ||
        refreshResponse.data?.refreshToken ||
        refreshToken;

      if (!nextAccessToken) {
        throw new Error('Refresh token response missing access token');
      }

      console.log('✅ Token refreshed successfully');
      localStorage.setItem(ACCESS_TOKEN_KEY, nextAccessToken);
      localStorage.setItem(REFRESH_TOKEN_KEY, nextRefreshToken);

      processQueue(null, nextAccessToken);

      originalRequest.headers = originalRequest.headers || {};
      originalRequest.headers.Authorization = `Bearer ${nextAccessToken}`;

      return axiosInstance(originalRequest);
    } catch (refreshError) {
      console.log('❌ Token refresh failed:', refreshError);
      processQueue(refreshError, null);
      forceLogout();
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  }
);

export default axiosInstance;
