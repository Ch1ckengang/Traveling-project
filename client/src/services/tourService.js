import axiosInstance from '../utils/axiosInstance';

export const getTours = async (params = {}) => {
  const response = await axiosInstance.get('/tours', { params });
  return response.data;
};

export const getTourById = async (tourId) => {
  const response = await axiosInstance.get(`/tours/${tourId}`);
  return response.data;
};

export const getDomesticTours = async (params = {}) => {
  const response = await axiosInstance.get('/tours/domestic', { params });
  return response.data;
};

export const getInternationalTours = async (params = {}) => {
  const response = await axiosInstance.get('/tours/international', { params });
  return response.data;
};

export const searchTours = async (keyword) => {
  const response = await axiosInstance.get('/tours/search', { params: { q: keyword } });
  return response.data;
};

export const getTourSchedules = async (tourId) => {
  const response = await axiosInstance.get(`/tours/${tourId}/schedules`);
  return response.data;
};

