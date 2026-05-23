import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/v1/api';

const getAuthHeader = () => {
  const token = localStorage.getItem('accessToken');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

export const getDashboardSummary = async () => {
  const response = await axios.get(`${API_BASE_URL}/admin/dashboard/summary`, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const getRevenueChart = async () => {
  const response = await axios.get(`${API_BASE_URL}/admin/dashboard/revenue-chart`, {
    headers: getAuthHeader(),
  });
  return response.data;
};

export const getTopTours = async () => {
  const response = await axios.get(`${API_BASE_URL}/admin/dashboard/top-tours`, {
    headers: getAuthHeader(),
  });
  return response.data;
};
