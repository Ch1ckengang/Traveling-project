import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/v1/api/admin';

const getAuthHeader = () => {
  const token = localStorage.getItem('accessToken');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

// ===== TOUR ADMIN =====

export const adminGetTours = async (params = {}) => {
  const response = await axios.get(`${API_BASE_URL}/tours`, {
    headers: getAuthHeader(),
    params,
  });
  return response.data;
};

export const adminCreateTour = async (tourData) => {
  const response = await axios.post(`${API_BASE_URL}/tours`, tourData, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const adminUpdateTour = async (id, tourData) => {
  const response = await axios.put(`${API_BASE_URL}/tours/${id}`, tourData, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const adminDeleteTour = async (id) => {
  const response = await axios.delete(`${API_BASE_URL}/tours/${id}`, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const adminToggleTour = async (id) => {
  const response = await axios.put(`${API_BASE_URL}/tours/${id}/toggle`, {}, {
    headers: getAuthHeader(),
  });
  return response.data;
};

// ===== BOOKING ADMIN =====

export const adminGetBookings = async (params = {}) => {
  const response = await axios.get(`${API_BASE_URL}/bookings`, {
    headers: getAuthHeader(),
    params,
  });
  return response.data;
};

export const adminGetBookingStats = async () => {
  const response = await axios.get(`${API_BASE_URL}/bookings/stats`, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const adminGetBookingByCode = async (code) => {
  const response = await axios.get(`${API_BASE_URL}/bookings/${code}`, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const adminConfirmBooking = async (code) => {
  const response = await axios.put(`${API_BASE_URL}/bookings/${code}/confirm`, {}, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const adminCancelBooking = async (code, reason = '') => {
  const response = await axios.put(`${API_BASE_URL}/bookings/${code}/cancel`, { reason }, {
    headers: getAuthHeader(),
  });
  return response.data;
};

// ===== USER ADMIN =====

export const adminGetUsers = async (params = {}) => {
  const response = await axios.get(`${API_BASE_URL}/users`, {
    headers: getAuthHeader(),
    params,
  });
  return response.data;
};

export const adminUpdateUserStatus = async (id, isActive) => {
  const response = await axios.put(`${API_BASE_URL}/users/${id}/status`, { is_active: isActive }, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const adminUpdateUserRole = async (id, role) => {
  const response = await axios.put(`${API_BASE_URL}/users/${id}/role`, { role }, {
    headers: getAuthHeader(),
  });
  return response.data;
};
