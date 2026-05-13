import axiosInstance from '../../utils/axiosInstance';

export const getNotifications = async (params = {}) => {
  const res = await axiosInstance.get('/notifications', { params });
  return res.data;
};

export const markAsRead = async (id) => {
  const res = await axiosInstance.put(`/notifications/${id}/read`);
  return res.data;
};

export const markAllAsRead = async () => {
  const res = await axiosInstance.put('/notifications/read-all');
  return res.data;
};
