import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/v1/api';

const getAuthHeader = () => {
  const token = localStorage.getItem('accessToken');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

// ===== PUBLIC/CUSTOMER =====

/** Validate coupon */
export const validateCoupon = async (code, orderTotal) => {
  const response = await axios.post(`${API_BASE_URL}/coupons/validate`, 
    { code, order_total: orderTotal },
    { headers: getAuthHeader() }
  );
  return response.data;
};

// ===== ADMIN =====

/** Lấy danh sách coupons */
export const adminGetCoupons = async (params = {}) => {
  const response = await axios.get(`${API_BASE_URL}/admin/coupons`, {
    headers: getAuthHeader(),
    params,
  });
  return response.data;
};

/** Tạo mới coupon */
export const adminCreateCoupon = async (data) => {
  const response = await axios.post(`${API_BASE_URL}/admin/coupons`, data, {
    headers: getAuthHeader(),
  });
  return response.data;
};

/** Cập nhật coupon */
export const adminUpdateCoupon = async (id, data) => {
  const response = await axios.put(`${API_BASE_URL}/admin/coupons/${id}`, data, {
    headers: getAuthHeader(),
  });
  return response.data;
};

/** Xóa coupon */
export const adminDeleteCoupon = async (id) => {
  const response = await axios.delete(`${API_BASE_URL}/admin/coupons/${id}`, {
    headers: getAuthHeader(),
  });
  return response.data;
};
