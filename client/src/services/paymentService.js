import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/v1/api';

/**
 * Payment Service - Quản lý các API calls cho thanh toán
 */

// Lấy access token từ localStorage
const getAuthHeader = () => {
  const token = localStorage.getItem('accessToken');
  return token ? { Authorization: `Bearer ${token}` } : {};
};

/**
 * Khởi tạo thanh toán - tạo phiên thanh toán VNPay
 * @param {number} bookingId - ID của booking cần thanh toán
 * @returns {Object} { payment_url, transaction_reference, expires_at }
 */
export const initiatePayment = async (bookingId) => {
  const response = await axios.post(
    `${API_BASE_URL}/payments/initiate`,
    { booking_id: bookingId },
    { headers: getAuthHeader() }
  );
  return response.data;
};

/**
 * Kiểm tra trạng thái thanh toán
 * @param {string} transactionRef - Mã giao dịch
 * @returns {Object} payment status info
 */
export const getPaymentStatus = async (transactionRef) => {
  const response = await axios.get(
    `${API_BASE_URL}/payments/status/${transactionRef}`,
    { headers: getAuthHeader() }
  );
  return response.data;
};

/**
 * Lấy danh sách thanh toán theo booking
 * @param {number} bookingId - ID booking
 * @returns {Object} { payments, total }
 */
export const getBookingPayments = async (bookingId) => {
  const response = await axios.get(
    `${API_BASE_URL}/bookings/${bookingId}/payments`,
    { headers: getAuthHeader() }
  );
  return response.data;
};

/**
 * Redirect user đến VNPay để thanh toán
 * @param {number} bookingId - ID booking
 */
export const redirectToPayment = async (bookingId) => {
  try {
    const result = await initiatePayment(bookingId);
    if (result.success && result.data?.payment_url) {
      window.location.href = result.data.payment_url;
      return { success: true };
    }
    return { success: false, message: result.message || 'Không thể tạo phiên thanh toán' };
  } catch (error) {
    const message = error.response?.data?.message || 'Lỗi kết nối đến server';
    return { success: false, message };
  }
};

/**
 * Format giá tiền VND
 * @param {number} amount - Số tiền (đơn vị VND)
 * @returns {string} Chuỗi format "1.000.000đ"
 */
export const formatPriceVND = (amount) => {
  if (!amount || amount <= 0) return 'Liên hệ';
  return amount.toLocaleString('vi-VN') + 'đ';
};

/**
 * Parse kết quả payment return từ URL query params
 * @returns {Object} { status, ref, message }
 */
export const parsePaymentResult = () => {
  const params = new URLSearchParams(window.location.search);
  return {
    status: params.get('status') || 'unknown',
    ref: params.get('ref') || '',
    message: decodeURIComponent(params.get('message') || ''),
  };
};

export default {
  initiatePayment,
  getPaymentStatus,
  getBookingPayments,
  redirectToPayment,
  formatPriceVND,
  parsePaymentResult,
};
