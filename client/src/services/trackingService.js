import axiosInstance from '../utils/axiosInstance';
import { v4 as uuidv4 } from 'uuid';

// Hàm lấy session id (lưu trong localStorage để track user chưa login)
export const getSessionId = () => {
  let sessionId = localStorage.getItem('tracking_session_id');
  if (!sessionId) {
    sessionId = uuidv4();
    localStorage.setItem('tracking_session_id', sessionId);
  }
  return sessionId;
};

// Gửi tracking log lên server
export const logActivity = async (actionType, data = {}) => {
  try {
    const payload = {
      session_id: getSessionId(),
      action_type: actionType,
      ...data
    };
    // Sử dụng axiosInstance.post (không await để không block UI)
    axiosInstance.post('/tracking/log', payload).catch(err => {
      console.warn('Failed to log activity:', err.message);
    });
  } catch (err) {
    console.warn('Tracking error:', err);
  }
};
