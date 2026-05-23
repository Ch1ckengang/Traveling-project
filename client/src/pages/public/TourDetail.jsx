import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axiosInstance from '../../utils/axiosInstance';
import BookingModal from '../../components/booking/BookingModal';
import ReviewList from '../../components/review/ReviewList';
import { logActivity } from '../../services/trackingService';

const TourDetailPage = () => {
  const { tourId } = useParams();
  const navigate = useNavigate();
  const [tour, setTour] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showBookingModal, setShowBookingModal] = useState(false);
  const [activeImage, setActiveImage] = useState(null);

  useEffect(() => {
    fetchTourDetail();
  }, [tourId]);

  const fetchTourDetail = async () => {
    try {
      setLoading(true);
      const response = await axiosInstance.get(`/tours/${tourId}`);
      // API trả về { success: true, data: { tour: {...} } } hoặc { success: true, data: {...} } hoặc tour trực tiếp
      const responseData = response.data?.data || response.data;
      const tourData = responseData.tour || responseData;
      
      if (tourData && tourData.id) {
        setTour(tourData);
        setActiveImage(tourData.image_url);
        setError(null);
        
        // Log view_tour activity sau 3 giây (đảm bảo người dùng thực sự xem)
        setTimeout(() => {
          logActivity('view_tour', { tour_id: tourData.id });
        }, 3000);
      } else {
        setError('Không tìm thấy tour');
      }
    } catch (err) {
      console.error('Error fetching tour:', err);
      if (err.response?.status === 404) {
        setError('Tour không tồn tại');
      } else {
        setError('Không thể tải thông tin tour');
      }
    } finally {
      setLoading(false);
    }
  };


  const handleBookingClick = () => {
    // Check if user is logged in - check both token systems
    const accessToken = localStorage.getItem('accessToken');
    const authTokens = localStorage.getItem('auth_tokens');
    const user = localStorage.getItem('user');
    
    if (!accessToken && !authTokens && !user) {
      alert('Vui lòng đăng nhập để đặt tour');
      navigate('/login');
      return;
    }
    
    setShowBookingModal(true);
  };

  const handleBookingSuccess = (booking) => {
    alert(`Đặt tour thành công! Mã đặt tour: ${booking.booking_code}`);
    // Refresh tour data to update remaining slots
    fetchTourDetail();
    // Navigate to bookings page
    setTimeout(() => {
      navigate('/account/bookings');
    }, 1500);
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">Đang tải...</div>
      </div>
    );
  }

  if (error || !tour) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-red-500">{error || 'Tour không tồn tại'}</div>
        <div className="text-center mt-4">
          <button
            onClick={() => navigate('/tours')}
            className="text-blue-600 hover:text-blue-800"
          >
            ← Quay lại danh sách tours
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Breadcrumb */}
      <div className="mb-6 text-sm text-gray-600">
        <button onClick={() => navigate('/tours')} className="hover:text-blue-600">
          Tours
        </button>
        <span className="mx-2">›</span>
        <span>{tour.name}</span>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Main Content */}
        <div className="lg:col-span-2">
          {/* Image Gallery */}
          <div className="mb-6">
            <div className="h-96 bg-gray-200 rounded-lg flex items-center justify-center overflow-hidden mb-2">
              {activeImage ? (
                <img src={activeImage} alt={tour.name} className="w-full h-full object-cover" />
              ) : (
                <span className="text-gray-400 text-2xl">📷 {tour.location}</span>
              )}
            </div>
            
            {/* Thumbnail list */}
            {tour.images && tour.images.length > 0 && (
              <div className="flex gap-2 overflow-x-auto py-2">
                {/* Main image as first thumbnail */}
                {tour.image_url && (
                  <div 
                    className={`w-20 h-16 flex-shrink-0 rounded cursor-pointer overflow-hidden border-2 ${activeImage === tour.image_url ? 'border-blue-500' : 'border-transparent'}`}
                    onClick={() => setActiveImage(tour.image_url)}
                  >
                    <img src={tour.image_url} alt="Main" className="w-full h-full object-cover" />
                  </div>
                )}
                {/* Additional images */}
                {tour.images.map((img) => (
                  <div 
                    key={img.id}
                    className={`w-20 h-16 flex-shrink-0 rounded cursor-pointer overflow-hidden border-2 ${activeImage === img.image_url ? 'border-blue-500' : 'border-transparent'}`}
                    onClick={() => setActiveImage(img.image_url)}
                  >
                    <img src={img.image_url} alt="Gallery" className="w-full h-full object-cover" />
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Tour Info */}
          <h1 className="text-3xl font-bold mb-4">{tour.name}</h1>

          <div className="flex items-center gap-4 mb-6">
            <span className="px-3 py-1 bg-blue-100 text-blue-600 rounded">
              {tour.type === 'domestic' ? 'Việt Nam' : tour.type === 'international' ? 'Quốc tế' : 'Dịch vụ'}
            </span>
            <span className="text-gray-600">📍 {tour.location}, {tour.country}</span>
            <span className="text-gray-600">⏱️ {tour.duration}</span>
          </div>

          {/* Description */}
          <div className="mb-8">
            <h2 className="text-xl font-bold mb-4">Mô tả</h2>
            <p className="text-gray-700 leading-relaxed">{tour.description}</p>
          </div>

          {/* Itinerary */}
          {tour.itinerary && (
            <div className="mb-8">
              <h2 className="text-xl font-bold mb-4">Lịch trình</h2>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="text-gray-700">{tour.itinerary}</p>
              </div>
            </div>
          )}

          {/* Services */}
          {tour.services && (
            <div className="mb-8">
              <h2 className="text-xl font-bold mb-4">Dịch vụ bao gồm</h2>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="text-gray-700">{tour.services}</p>
              </div>
            </div>
          )}

          {/* Reviews */}
          <ReviewList tourId={tour.id} />
        </div>

        {/* Sidebar - Booking Card */}
        <div className="lg:col-span-1">
          <div className="border rounded-lg p-6 sticky top-4">
            <div className="mb-6">
              <div className="text-3xl font-bold text-blue-600 mb-2">{tour.price}</div>
              <div className="text-sm text-gray-600">Giá mỗi người</div>
            </div>

            <div className="space-y-4 mb-6">
              <div className="flex justify-between">
                <span className="text-gray-600">Thời gian:</span>
                <span className="font-medium">{tour.duration}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Còn lại:</span>
                <span className="font-medium">{tour.remaining_slots} chỗ</span>
              </div>
            </div>

            <button
              onClick={handleBookingClick}
              className="w-full bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors font-medium"
            >
              Đặt tour ngay
            </button>

            <div className="mt-4 text-center text-sm text-gray-500">
              Miễn phí hủy trong 24h
            </div>
          </div>
        </div>
      </div>

      {/* Booking Modal */}
      <BookingModal
        tour={tour}
        isOpen={showBookingModal}
        onClose={() => setShowBookingModal(false)}
        onSuccess={handleBookingSuccess}
      />
    </div>
  );
};

export default TourDetailPage;
