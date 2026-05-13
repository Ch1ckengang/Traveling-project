import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { parsePaymentResult, getPaymentStatus } from '../../services/paymentService';
import './PaymentResult.css';

/**
 * PaymentResult - Trang hiển thị kết quả thanh toán
 * Hiển thị sau khi VNPay redirect user về
 */
const PaymentResult = () => {
  const navigate = useNavigate();
  const [result, setResult] = useState(null);
  const [payment, setPayment] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const loadResult = async () => {
      const parsed = parsePaymentResult();
      setResult(parsed);

      // Nếu có transaction ref, lấy chi tiết payment
      if (parsed.ref) {
        try {
          const statusRes = await getPaymentStatus(parsed.ref);
          if (statusRes.success && statusRes.data?.payment) {
            setPayment(statusRes.data.payment);
          }
        } catch {
          // Không block UI nếu không lấy được chi tiết
        }
      }
      setLoading(false);
    };

    loadResult();
  }, []);

  if (loading) {
    return (
      <div className="payment-result-container">
        <div className="payment-result-card">
          <div className="payment-loading">
            <div className="spinner"></div>
            <p>Đang xử lý kết quả thanh toán...</p>
          </div>
        </div>
      </div>
    );
  }

  const isSuccess = result?.status === 'success';

  return (
    <div className="payment-result-container">
      <div className="payment-result-card">
        <div className={`payment-icon ${isSuccess ? 'success' : 'failed'}`}>
          {isSuccess ? '✓' : '✕'}
        </div>

        <h1 className="payment-title">
          {isSuccess ? 'Thanh toán thành công!' : 'Thanh toán không thành công'}
        </h1>

        <p className="payment-message">
          {result?.message || (isSuccess
            ? 'Booking của bạn đã được xác nhận.'
            : 'Vui lòng thử lại hoặc liên hệ hỗ trợ.'
          )}
        </p>

        {payment && (
          <div className="payment-details">
            <div className="detail-row">
              <span className="detail-label">Mã giao dịch:</span>
              <span className="detail-value">{payment.transaction_reference}</span>
            </div>
            {payment.amount > 0 && (
              <div className="detail-row">
                <span className="detail-label">Số tiền:</span>
                <span className="detail-value">
                  {payment.amount.toLocaleString('vi-VN')}đ
                </span>
              </div>
            )}
            <div className="detail-row">
              <span className="detail-label">Trạng thái:</span>
              <span className={`detail-value status-${payment.status}`}>
                {payment.status_display_name}
              </span>
            </div>
          </div>
        )}

        <div className="payment-actions">
          <button
            className="btn-primary"
            onClick={() => navigate('/bookings')}
          >
            Xem danh sách booking
          </button>
          <button
            className="btn-secondary"
            onClick={() => navigate('/')}
          >
            Về trang chủ
          </button>
        </div>
      </div>
    </div>
  );
};

export default PaymentResult;
