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

let isForceLoggingOut = false;

const forceLogout = () => {
  if (isForceLoggingOut) return; // Prevent multiple concurrent logouts
  isForceLoggingOut = true;

  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem('auth_tokens');
  localStorage.removeItem('user');
  window.dispatchEvent(new Event(LOGOUT_EVENT_NAME));
  
  // Delay redirect to let React process the auth:logout event
  // This prevents the 401 loop by clearing state BEFORE redirect
  setTimeout(() => {
    isForceLoggingOut = false;
    window.location.href = '/login';
  }, 100);
};

axiosInstance.interceptors.request.use(
  (config) => {
    const accessToken = localStorage.getItem(ACCESS_TOKEN_KEY);
    
    if (accessToken) {
      config.headers = config.headers || {};
      config.headers.Authorization = `Bearer ${accessToken}`;
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

    if (status !== 401 || !originalRequest || originalRequest._retry) {
      return Promise.reject(error);
    }

    originalRequest._retry = true;

    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    if (!refreshToken) {
      forceLogout();
      return Promise.reject(error);
    }

    if (isRefreshing) {
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
      // Backend expects "refreshToken" (camelCase) in request body
      const refreshResponse = await refreshClient.post('/token/refresh', {
        refreshToken: refreshToken
      });

      // Backend response format: {success, data: {tokens: {access_token, refresh_token}}}
      const responseData = refreshResponse.data;
      const tokensObj =
        responseData?.data?.tokens ||   // Primary: {data: {tokens: {...}}}
        responseData?.tokens ||          // Fallback: {tokens: {...}}
        responseData?.data ||            // Fallback: {data: {access_token, ...}}
        null;

      const nextAccessToken = tokensObj?.access_token || tokensObj?.AccessToken || null;
      const nextRefreshToken = tokensObj?.refresh_token || tokensObj?.RefreshToken || refreshToken;

      if (!nextAccessToken) {
        throw new Error('Refresh token response missing access token');
      }

      localStorage.setItem(ACCESS_TOKEN_KEY, nextAccessToken);
      localStorage.setItem(REFRESH_TOKEN_KEY, nextRefreshToken);

      processQueue(null, nextAccessToken);

      originalRequest.headers = originalRequest.headers || {};
      originalRequest.headers.Authorization = `Bearer ${nextAccessToken}`;

      return axiosInstance(originalRequest);
    } catch (refreshError) {
      processQueue(refreshError, null);
      forceLogout();
      return Promise.reject(refreshError);
    } finally {
      isRefreshing = false;
    }
  }
);

export default axiosInstance;
