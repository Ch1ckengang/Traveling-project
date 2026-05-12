import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import { getUserBookings, cancelBooking } from '../../services/bookingService';
import { formatCurrency } from '../../utils/formatCurrency';
import Spinner from '../../components/ui/Spinner';
import EmptyState from '../../components/ui/EmptyState';

const BookingsPage = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [bookings, setBookings] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [cancellingId, setCancellingId] = useState(null);

  useEffect(() => {
    if (!user?.id) {
      navigate('/login');
      return;
    }
    fetchBookings();
  }, [user]);

  const fetchBookings = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getUserBookings(user.id);
      
      if (response.success) {
        setBookings(response.bookings || []);
      } else {
        setError(response.message || 'Không thể tải danh sách đặt tour');
      }
    } catch (err) {
      console.error('Error fetching bookings:', err);
      setError('Có lỗi xảy ra khi tải danh sách đặt tour');
    } finally {
      setLoading(false);
    }
  };

  const handleCancelBooking = async (bookingId) => {
    if (!window.confirm('Bạn có chắc chắn muốn hủy tour này?')) {
      return;
    }

    try {
      setCancellingId(bookingId);
      const response = await cancelBooking(user.id, bookingId);
      
      if (response.success) {
        alert('Hủy tour thành công');
        fetchBookings(); // Refresh list
      } else {
        alert(response.message || 'Không thể hủy tour');
      }
    } catch (err) {
      console.error('Error cancelling booking:', err);
      alert(err.response?.data?.message || 'Có lỗi xảy ra khi hủy tour');
    } finally {
      setCancellingId(null);
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
      <span className={`px-3 py-1 rounded-full text-sm font-medium ${config.color}`}>
        {config.label}
      </span>
    );
  };

  const formatDate = (dateString) => {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    if (isNaN(date.getTime())) return 'N/A';
    
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = date.getFullYear();
    
    return `${day}/${month}/${year}`;
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

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
          {error}
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-2">Đặt tour của tôi</h1>
        <p className="text-gray-600">Quản lý các tour bạn đã đặt</p>
      </div>

      {/* Bookings List */}
      {bookings.length === 0 ? (
        <EmptyState
          icon="📅"
          title="Chưa có đặt tour nào"
          description="Bạn chưa đặt tour nào. Hãy khám phá các tour du lịch hấp dẫn!"
          actionLabel="Xem tours"
          onAction={() => navigate('/tours')}
        />
      ) : (
        <div className="space-y-4">
          {bookings.map((booking) => (
            <div
              key={booking.id}
              className="bg-white border rounded-lg p-6 hover:shadow-md transition-shadow"
            >
              <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
                {/* Left: Booking Info */}
                <div className="flex-1">
                  <div className="flex items-start justify-between mb-3">
                    <div>
                      <h3 className="text-xl font-bold mb-1">{booking.tour?.name || 'Tour không xác định'}</h3>
                      <p className="text-sm text-gray-500">Mã đặt tour: {booking.booking_code}</p>
                    </div>
                    {getStatusBadge(booking.status)}
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
                    <div>
                      <span className="text-gray-600">👤 Người đặt:</span>
                      <span className="ml-2 font-medium">{booking.full_name}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">📞 Điện thoại:</span>
                      <span className="ml-2 font-medium">{booking.phone}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">📧 Email:</span>
                      <span className="ml-2 font-medium">{booking.email}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">📅 Ngày khởi hành:</span>
                      <span className="ml-2 font-medium">{formatDate(booking.travel_date)}</span>
                    </div>
                    <div>
                      <span className="text-gray-600">👥 Số người:</span>
                      <span className="ml-2 font-medium">
                        {booking.adult_count} người lớn
                        {booking.child_count > 0 && `, ${booking.child_count} trẻ em`}
                        {booking.infant_count > 0 && `, ${booking.infant_count} em bé`}
                      </span>
                    </div>
                    <div>
                      <span className="text-gray-600">💰 Tổng tiền:</span>
                      <span className="ml-2 font-bold text-blue-600">{formatCurrency(booking.total_amount)}</span>
                    </div>
                  </div>

                  {booking.note && (
                    <div className="mt-3 text-sm">
                      <span className="text-gray-600">📝 Ghi chú:</span>
                      <span className="ml-2">{booking.note}</span>
                    </div>
                  )}

                  <div className="mt-3 text-xs text-gray-500">
                    Đặt lúc: {formatDate(booking.created_at)}
                  </div>
                </div>

                {/* Right: Actions */}
                <div className="flex lg:flex-col gap-2">
                  <button
                    onClick={() => navigate(`/account/bookings/${booking.id}`)}
                    className="flex-1 lg:flex-none px-4 py-2 border border-blue-600 text-blue-600 rounded hover:bg-blue-50 transition-colors"
                  >
                    Chi tiết
                  </button>
                  
                  {canCancelBooking(booking) && (
                    <button
                      onClick={() => handleCancelBooking(booking.id)}
                      disabled={cancellingId === booking.id}
                      className="flex-1 lg:flex-none px-4 py-2 border border-red-600 text-red-600 rounded hover:bg-red-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      {cancellingId === booking.id ? 'Đang hủy...' : 'Hủy tour'}
                    </button>
                  )}

                  {booking.status === 'pending' && (
                    <button
                      onClick={() => navigate(`/payment/${booking.id}`)}
                      className="flex-1 lg:flex-none px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
                    >
                      Thanh toán
                    </button>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default BookingsPage;
