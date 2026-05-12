import axiosInstance from '../utils/axiosInstance';

export const createBooking = async (bookingData) => {
  const response = await axiosInstance.post('/bookings', bookingData);
  return response.data;
};

export const getUserBookings = async (userId) => {
  const response = await axiosInstance.get(`/users/${userId}/bookings`);
  return response.data;
};

export const getBookingById = async (bookingId) => {
  const response = await axiosInstance.get(`/bookings/${bookingId}`);
  return response.data;
};

export const getBookingByCode = async (bookingCode) => {
  const response = await axiosInstance.get(`/bookings/code/${bookingCode}`);
  return response.data;
};

export const cancelBooking = async (userId, bookingId) => {
  const response = await axiosInstance.put(`/users/${userId}/bookings/${bookingId}/cancel`);
  return response.data;
};
