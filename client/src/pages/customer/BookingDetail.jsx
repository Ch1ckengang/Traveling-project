import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { getBookingById, cancelBooking } from '../../services/bookingService';
import { formatCurrency } from '../../utils/formatCurrency';
import Spinner from '../../components/ui/Spinner';

const BookingDetailPage = () => {
  const { bookingId } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [booking, setBooking] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [cancelling, setCancelling] = useState(false);

  useEffect(() => {
    if (!user?.id) {
      navigate('/login');
      return;
    }
    fetchBookingDetail();
  }, [user, bookingId]);

  const fetchBookingDetail = async () => {
    try {
      setLoading(true);
      setError(null);
      
      // Use new endpoint GET /bookings/:id
      const response = await getBookingById(bookingId);
      
      if (response.success) {
        setBooking(response.booking);
      } else {
        setError(response.message || 'Không thể tải thông tin đặt tour');
      }
    } catch (err) {
      console.error('Error fetching booking:', err);
      
      if (err.response?.status === 404) {
        setError('Không tìm thấy đặt tour');
      } else if (err.response?.status === 403) {
        setError('Bạn không có quyền xem đặt tour này');
      } else {
        setError('Có lỗi xảy ra khi tải thông tin đặt tour');
      }
    } finally {
      setLoading(false);
    }
  };

  const handleCancelBooking = async () => {
    if (!window.confirm('Bạn có chắc chắn muốn hủy tour này?')) {
      return;
    }

    try {
      setCancelling(true);
      const response = await cancelBooking(user.id, bookingId);
      
      if (response.success) {
        alert('Hủy tour thành công');
        navigate('/account/bookings');
      } else {
        alert(response.message || 'Không thể hủy tour');
      }
    } catch (err) {
      console.error('Error cancelling booking:', err);
      alert(err.response?.data?.message || 'Có lỗi xảy ra khi hủy tour');
    } finally {
      setCancelling(false);
    }
  };

  const getStatusBadge = (status) => {
    const statusConfig = {
      pending: { label: 'Chờ thanh toán', color: 'bg-yellow-100 text-yellow-800' },
      confirmed: { label: 'Đã xác nhận', color: 'bg-green-100 text-green-800' },
      cancelled: { label: 'Đã hủy', color: 'bg-red-100 text-red-800' },
      completed: { label: 'Hoàn thành', color: 'bg-blue-100 text-blue-800' }
    };

    const config = statusConfig[status] || { label: status, color: 'bg-gray-100 text-gray-800' };
    return (
      <span className={`px-4 py-2 rounded-full text-sm font-medium ${config.color}`}>
        {config.label}
      </span>
    );
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return 'N/A';
    
    const weekdays = ['Chủ Nhật', 'Thứ Hai', 'Thứ Ba', 'Thứ Tư', 'Thứ Năm', 'Thứ Sáu', 'Thứ Bảy'];
    const weekday = weekdays[date.getDay()];
    
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    
    return `${weekday}, ${day}/${month}/${year}`;
  };

  const formatDateTime = (dateString) => {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return 'N/A';
    
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    
    return `${day}/${month}/${year} ${hours}:${minutes}`;
  };

  const canCancelBooking = (booking) => {
    return booking.status === 'pending' || booking.status === 'confirmed';
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <Spinner />
      </div>
    );
  }

  if (error || !booking) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
          {error || 'Không tìm thấy đặt tour'}
        </div>
        <button
          onClick={() => navigate('/account/bookings')}
          className="text-blue-600 hover:text-blue-800"
        >
          ← Quay lại danh sách đặt tour
        </button>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Breadcrumb */}
      <div className="mb-6 text-sm text-gray-600">
        <button onClick={() => navigate('/account/bookings')} className="hover:text-blue-600">
          Đặt tour của tôi
        </button>
        <span className="mx-2">›</span>
        <span>{booking.booking_code}</span>
      </div>

      {/* Header */}
      <div className="bg-white rounded-lg shadow-sm border p-6 mb-6">
        <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4 mb-4">
          <div>
            <h1 className="text-3xl font-bold mb-2">{booking.tour?.name || 'Tour không xác định'}</h1>
            <p className="text-gray-600">Mã đặt tour: <span className="font-mono font-bold">{booking.booking_code}</span></p>
          </div>
          {getStatusBadge(booking.status)}
        </div>

        {booking.status === 'pending' && (
          <div className="bg-yellow-50 border border-yellow-200 text-yellow-800 px-4 py-3 rounded">
            ⚠️ Vui lòng thanh toán trong vòng 24 giờ để giữ chỗ
          </div>
        )}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Booking Information */}
          <div className="bg-white rounded-lg shadow-sm border p-6">
            <h2 className="text-xl font-bold mb-4">Thông tin đặt tour</h2>
            
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="text-sm text-gray-600">Họ và tên</label>
                  <p className="font-medium">{booking.full_name}</p>
                </div>
                <div>
                  <label className="text-sm text-gray-600">Số điện thoại</label>
                  <p className="font-medium">{booking.phone}</p>
                </div>
                <div>
                  <label className="text-sm text-gray-600">Email</label>
                  <p className="font-medium">{booking.email}</p>
                </div>
                <div>
                  <label className="text-sm text-gray-600">Ngày khởi hành</label>
                  <p className="font-medium">{formatDate(booking.travel_date)}</p>
                </div>
              </div>

              <div className="border-t pt-4">
                <label className="text-sm text-gray-600">Số lượng hành khách</label>
                <div className="mt-2 space-y-2">
                  <div className="flex justify-between">
                    <span>Người lớn</span>
                    <span className="font-medium">{booking.adult_count} người</span>
                  </div>
                  {booking.child_count > 0 && (
                    <div className="flex justify-between">
                      <span>Trẻ em (2-12 tuổi)</span>
                      <span className="font-medium">{booking.child_count} người</span>
                    </div>
                  )}
                  {booking.infant_count > 0 && (
                    <div className="flex justify-between">
                      <span>Em bé (&lt;2 tuổi)</span>
                      <span className="font-medium">{booking.infant_count} người</span>
                    </div>
                  )}
                  <div className="flex justify-between font-bold border-t pt-2">
                    <span>Tổng số người</span>
                    <span>{booking.adult_count + booking.child_count + booking.infant_count} người</span>
                  </div>
                </div>
              </div>

              {booking.note && (
                <div className="border-t pt-4">
                  <label className="text-sm text-gray-600">Ghi chú</label>
                  <p className="mt-1">{booking.note}</p>
                </div>
              )}
            </div>
          </div>

          {/* Timeline */}
          <div className="bg-white rounded-lg shadow-sm border p-6">
            <h2 className="text-xl font-bold mb-4">Lịch sử</h2>
            
            <div className="space-y-4">
              <div className="flex gap-4">
                <div className="flex flex-col items-center">
                  <div className="w-3 h-3 bg-blue-600 rounded-full"></div>
                  <div className="w-0.5 h-full bg-gray-300"></div>
                </div>
                <div className="pb-4">
                  <p className="font-medium">Đặt tour thành công</p>
                  <p className="text-sm text-gray-600">{formatDateTime(booking.created_at)}</p>
                </div>
              </div>

              {booking.status === 'cancelled' && booking.updated_at && (
                <div className="flex gap-4">
                  <div className="flex flex-col items-center">
                    <div className="w-3 h-3 bg-red-600 rounded-full"></div>
                  </div>
                  <div>
                    <p className="font-medium">Tour đã bị hủy</p>
                    <p className="text-sm text-gray-600">{formatDateTime(booking.updated_at)}</p>
                  </div>
                </div>
              )}

              {booking.status === 'confirmed' && (
                <div className="flex gap-4">
                  <div className="flex flex-col items-center">
                    <div className="w-3 h-3 bg-green-600 rounded-full"></div>
                  </div>
                  <div>
                    <p className="font-medium">Tour đã được xác nhận</p>
                    <p className="text-sm text-gray-600">{formatDateTime(booking.updated_at)}</p>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Sidebar */}
        <div className="lg:col-span-1">
          {/* Price Summary */}
          <div className="bg-white rounded-lg shadow-sm border p-6 mb-6 sticky top-4">
            <h2 className="text-xl font-bold mb-4">Tổng quan</h2>
            
            <div className="space-y-3 mb-6">
              <div className="flex justify-between">
                <span className="text-gray-600">Giá tour</span>
                <span className="font-medium">{formatCurrency(booking.total_amount)}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Số người</span>
                <span className="font-medium">{booking.adult_count + booking.child_count + booking.infant_count}</span>
              </div>
              <div className="border-t pt-3 flex justify-between text-lg font-bold">
                <span>Tổng tiền</span>
                <span className="text-blue-600">{formatCurrency(booking.total_amount)}</span>
              </div>
            </div>

            {/* Actions */}
            <div className="space-y-3">
              {booking.status === 'pending' && (
                <button
                  onClick={() => navigate(`/payment/${booking.id}`)}
                  className="w-full bg-green-600 text-white py-3 rounded-lg hover:bg-green-700 transition-colors font-medium"
                >
                  Thanh toán ngay
                </button>
              )}

              {canCancelBooking(booking) && (
                <button
                  onClick={handleCancelBooking}
                  disabled={cancelling}
                  className="w-full border border-red-600 text-red-600 py-3 rounded-lg hover:bg-red-50 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {cancelling ? 'Đang hủy...' : 'Hủy tour'}
                </button>
              )}

              <button
                onClick={() => navigate('/account/bookings')}
                className="w-full border border-gray-300 text-gray-700 py-3 rounded-lg hover:bg-gray-50 transition-colors font-medium"
              >
                Quay lại
              </button>
            </div>
          </div>

          {/* Help */}
          <div className="bg-blue-50 rounded-lg border border-blue-200 p-6">
            <h3 className="font-bold mb-2">Cần hỗ trợ?</h3>
            <p className="text-sm text-gray-700 mb-4">
              Liên hệ với chúng tôi nếu bạn có bất kỳ câu hỏi nào về đặt tour của mình.
            </p>
            <div className="space-y-2 text-sm">
              <div>📞 Hotline: 1900-xxxx</div>
              <div>📧 Email: support@travel.com</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BookingDetailPage;
