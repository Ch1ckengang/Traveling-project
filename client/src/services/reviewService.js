import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/v1/api';

const getAuthHeader = () => {
  const token = localStorage.getItem('accessToken');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

// ===== PUBLIC =====

/** Lấy reviews của tour (public) */
export const getTourReviews = async (tourId, params = {}) => {
  const response = await axios.get(`${API_BASE_URL}/tours/${tourId}/reviews`, { params });
  return response.data;
};

// ===== CUSTOMER =====

/** Tạo review mới */
export const createReview = async (data) => {
  const response = await axios.post(`${API_BASE_URL}/reviews`, data, {
    headers: getAuthHeader(),
  });
  return response.data;
};

/** Sửa review (trong 7 ngày) */
export const updateReview = async (reviewId, data) => {
  const response = await axios.put(`${API_BASE_URL}/reviews/${reviewId}`, data, {
    headers: getAuthHeader(),
  });
  return response.data;
};

// ===== ADMIN =====

/** Admin lấy tất cả reviews */
export const adminGetReviews = async (params = {}) => {
  const response = await axios.get(`${API_BASE_URL}/admin/reviews`, {
    headers: getAuthHeader(),
    params,
  });
  return response.data;
};

/** Admin publish review */
export const adminPublishReview = async (reviewId) => {
  const response = await axios.put(`${API_BASE_URL}/admin/reviews/${reviewId}/publish`, {}, {
    headers: getAuthHeader(),
  });
  return response.data;
};

/** Admin ẩn review */
export const adminHideReview = async (reviewId) => {
  const response = await axios.put(`${API_BASE_URL}/admin/reviews/${reviewId}/hide`, {}, {
    headers: getAuthHeader(),
  });
  return response.data;
};

/** Admin reply review */
export const adminReplyReview = async (reviewId, reply) => {
  const response = await axios.post(`${API_BASE_URL}/admin/reviews/${reviewId}/reply`, { reply }, {
    headers: getAuthHeader(),
  });
  return response.data;
};
