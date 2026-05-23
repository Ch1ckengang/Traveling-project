import { useState, useEffect } from 'react';
import { createBooking } from '../../services/bookingService';
import { validateCoupon } from '../../services/couponService';
import { getTourSchedules } from '../../services/tourService';
import DateInput from '../ui/DateInput';

const BookingModal = ({ tour, isOpen, onClose, onSuccess }) => {
  const [formData, setFormData] = useState({
    fullName: '',
    phone: '',
    email: '',
    adultCount: 1,
    childCount: 0,
    infantCount: 0,
    travelDate: '',
    scheduleId: '',
    notes: '',
    couponCode: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  
  // Coupon state
  const [validatingCoupon, setValidatingCoupon] = useState(false);
  const [couponError, setCouponError] = useState(null);
  const [appliedCoupon, setAppliedCoupon] = useState(null);

  // Schedules state
  const [schedules, setSchedules] = useState([]);
  const [loadingSchedules, setLoadingSchedules] = useState(false);

  useEffect(() => {
    if (isOpen && tour?.id) {
      fetchSchedules();
    }
  }, [isOpen, tour?.id]);

  const fetchSchedules = async () => {
    try {
      setLoadingSchedules(true);
      const res = await getTourSchedules(tour.id);
      if (res.success && res.data) {
        setSchedules(res.data);
      }
    } catch (err) {
      console.error('Lỗi lấy lịch trình:', err);
    } finally {
      setLoadingSchedules(false);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    
    // Reset coupon if changing people count since it affects total
    if (['adultCount', 'childCount', 'infantCount'].includes(name)) {
      setAppliedCoupon(null);
      setCouponError(null);
    }
    
    if (name === 'couponCode') {
      setCouponError(null);
      setFormData(prev => ({
        ...prev,
        [name]: value.toUpperCase()
      }));
      return;
    }

    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Helper to calculate raw total amount
  const calculateTotal = () => {
    const priceRaw = tour.price_amount || parseInt(tour.price.replace(/\D/g, '')) || 0;
    
    // Add schedule price modifier if selected
    let scheduleModifier = 0;
    if (formData.scheduleId) {
      const selectedSchedule = schedules.find(s => s.id === parseInt(formData.scheduleId));
      if (selectedSchedule) {
        scheduleModifier = selectedSchedule.price_modifier || 0;
      }
    }
    
    const baseAdultPrice = priceRaw + scheduleModifier;

    const adults = parseInt(formData.adultCount || 0);
    const children = parseInt(formData.childCount || 0);
    const infants = parseInt(formData.infantCount || 0);
    
    const adultTotal = baseAdultPrice * adults;
    const childTotal = (baseAdultPrice * 75 / 100) * children;
    const infantTotal = (baseAdultPrice * 25 / 100) * infants;
    
    return adultTotal + childTotal + infantTotal;
  };

  const handleValidateCoupon = async () => {
    if (!formData.couponCode.trim()) return;
    
    setValidatingCoupon(true);
    setCouponError(null);
    try {
      const totalAmount = calculateTotal();
      const res = await validateCoupon(formData.couponCode, totalAmount);
      if (res.success) {
        setAppliedCoupon(res.data);
      }
    } catch (err) {
      setAppliedCoupon(null);
      setCouponError(err.response?.data?.message || 'Mã giảm giá không hợp lệ');
    }
    setValidatingCoupon(false);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    // Validation
    if (!formData.fullName || !formData.phone || !formData.email) {
      setError('Vui lòng điền đầy đủ thông tin bắt buộc');
      return;
    }

    if (schedules.length > 0 && !formData.scheduleId) {
      setError('Vui lòng chọn ngày khởi hành');
      return;
    }

    if (schedules.length === 0 && !formData.travelDate) {
      setError('Vui lòng nhập ngày khởi hành');
      return;
    }

    const totalPeople = parseInt(formData.adultCount) + parseInt(formData.childCount) + parseInt(formData.infantCount);
    
    // Check against schedule slots if selected, else tour slots
    let availableSlots = tour.remaining_slots;
    if (formData.scheduleId) {
      const selectedSchedule = schedules.find(s => s.id === parseInt(formData.scheduleId));
      if (selectedSchedule) {
        availableSlots = selectedSchedule.remaining_slots;
      }
    }

    if (totalPeople > availableSlots) {
      setError(`Chỉ còn ${availableSlots} chỗ trống`);
      return;
    }

    // Debug: Check tokens before booking
    const accessToken = localStorage.getItem('accessToken');
    const authTokens = localStorage.getItem('auth_tokens');
    console.log('🔍 Debug - Before booking:');
    console.log('  accessToken:', accessToken ? 'EXISTS' : 'MISSING');
    console.log('  auth_tokens:', authTokens ? 'EXISTS' : 'MISSING');

    if (!accessToken && !authTokens) {
      setError('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.');
      return;
    }

    try {
      setLoading(true);
      setError(null);

      const bookingData = {
        tour_id: tour.id,
        full_name: formData.fullName,
        phone: formData.phone,
        email: formData.email,
        adult_count: parseInt(formData.adultCount),
        child_count: parseInt(formData.childCount),
        infant_count: parseInt(formData.infantCount),
        travel_date: schedules.length > 0 ? '' : formData.travelDate,
        schedule_id: formData.scheduleId ? parseInt(formData.scheduleId) : undefined,
        coupon_code: appliedCoupon ? appliedCoupon.coupon_code : '',
        note: formData.notes
      };

      console.log('📤 Sending booking request:', bookingData);
      const response = await createBooking(bookingData);
      console.log('📥 Booking response:', response);
      
      if (response.success) {
        onSuccess(response.booking);
        onClose();
      } else {
        setError(response.message || 'Đặt tour thất bại');
      }
    } catch (err) {
      console.error('❌ Booking error:', err);
      console.error('  Status:', err.response?.status);
      console.error('  Data:', err.response?.data);
      
      // Check if it's an auth error
      if (err.response?.status === 401) {
        setError('Phiên đăng nhập đã hết hạn. Vui lòng đăng nhập lại.');
      } else {
        setError(err.response?.data?.message || 'Có lỗi xảy ra khi đặt tour');
      }
    } finally {
      setLoading(false);
    }
  };

  if (!isOpen) return null;

  const totalPeople = parseInt(formData.adultCount || 0) + parseInt(formData.childCount || 0) + parseInt(formData.infantCount || 0);
  const baseTotal = calculateTotal();
  const finalTotal = appliedCoupon ? appliedCoupon.final_total : baseTotal;
  const formatPrice = (val) => val.toLocaleString('vi-VN') + 'đ';

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="sticky top-0 bg-white border-b px-6 py-4 flex justify-between items-center">
          <h2 className="text-2xl font-bold">Đặt Tour</h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl"
          >
            ×
          </button>
        </div>

        {/* Tour Info */}
        <div className="px-6 py-4 bg-gray-50 border-b">
          <h3 className="font-bold text-lg mb-2">{tour.name}</h3>
          <div className="flex gap-4 text-sm text-gray-600">
            <span>📍 {tour.location}</span>
            <span>⏱️ {tour.duration}</span>
            <span className="font-bold text-blue-600">{tour.price}</span>
          </div>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} className="px-6 py-4">
          {error && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded">
              {error}
            </div>
          )}

          {/* Contact Info */}
          <div className="mb-6">
            <h4 className="font-bold mb-3">Thông tin liên hệ</h4>
            
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">
                Họ và tên <span className="text-red-500">*</span>
              </label>
              <input
                type="text"
                name="fullName"
                value={formData.fullName}
                onChange={handleChange}
                className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Nguyễn Văn A"
                required
              />
            </div>

            <div className="grid grid-cols-2 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium mb-1">
                  Số điện thoại <span className="text-red-500">*</span>
                </label>
                <input
                  type="tel"
                  name="phone"
                  value={formData.phone}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="0912345678"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">
                  Email <span className="text-red-500">*</span>
                </label>
                <input
                  type="email"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="email@example.com"
                  required
                />
              </div>
            </div>
          </div>

          {/* Booking Details */}
          <div className="mb-6">
            <h4 className="font-bold mb-3">Chi tiết đặt tour</h4>
            
            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">
                Ngày khởi hành <span className="text-red-500">*</span>
              </label>
              {loadingSchedules ? (
                <div className="text-sm text-gray-500">Đang tải lịch trình...</div>
              ) : schedules.length > 0 ? (
                <select
                  name="scheduleId"
                  value={formData.scheduleId}
                  onChange={handleChange}
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                >
                  <option value="">-- Chọn ngày khởi hành --</option>
                  {schedules.map(schedule => (
                    <option key={schedule.id} value={schedule.id}>
                      {new Date(schedule.departure_date).toLocaleDateString('vi-VN')} 
                      {schedule.price_modifier > 0 ? ` (+${schedule.price_modifier.toLocaleString('vi-VN')}đ)` : ''} 
                      {schedule.price_modifier < 0 ? ` (${schedule.price_modifier.toLocaleString('vi-VN')}đ)` : ''} 
                      {' '}- Còn {schedule.remaining_slots} chỗ
                    </option>
                  ))}
                </select>
              ) : (
                <>
                  <DateInput
                    name="travelDate"
                    value={formData.travelDate}
                    onChange={handleChange}
                    min={new Date().toISOString().split('T')[0]}
                    required
                  />
                  <p className="text-xs text-gray-500 mt-1">
                    Nhập theo định dạng: ngày/tháng/năm (ví dụ: 15/05/2026)
                  </p>
                </>
              )}
            </div>

            <div className="grid grid-cols-3 gap-4 mb-4">
              <div>
                <label className="block text-sm font-medium mb-1">
                  Người lớn <span className="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  name="adultCount"
                  value={formData.adultCount}
                  onChange={handleChange}
                  min="1"
                  max={tour.remaining_slots}
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">
                  Trẻ em (2-12 tuổi)
                </label>
                <input
                  type="number"
                  name="childCount"
                  value={formData.childCount}
                  onChange={handleChange}
                  min="0"
                  max={tour.remaining_slots}
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-1">
                  Em bé (&lt;2 tuổi)
                </label>
                <input
                  type="number"
                  name="infantCount"
                  value={formData.infantCount}
                  onChange={handleChange}
                  min="0"
                  max={tour.remaining_slots}
                  className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
            </div>

            <div className="mb-4">
              <label className="block text-sm font-medium mb-1">
                Ghi chú
              </label>
              <textarea
                name="notes"
                value={formData.notes}
                onChange={handleChange}
                rows="3"
                className="w-full px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                placeholder="Yêu cầu đặc biệt, câu hỏi..."
              />
            </div>
          </div>

          {/* Coupon */}
          <div className="mb-6">
            <h4 className="font-bold mb-3">Mã giảm giá</h4>
            <div className="flex gap-2">
              <input
                type="text"
                name="couponCode"
                value={formData.couponCode}
                onChange={handleChange}
                placeholder="Nhập mã (nếu có)"
                className="flex-1 px-3 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                type="button"
                onClick={handleValidateCoupon}
                disabled={validatingCoupon || !formData.couponCode}
                className="px-4 py-2 bg-gray-100 border rounded hover:bg-gray-200 disabled:opacity-50 font-medium"
              >
                {validatingCoupon ? 'Đang ktra...' : 'Áp dụng'}
              </button>
            </div>
            {couponError && <p className="text-red-500 text-sm mt-1">{couponError}</p>}
            {appliedCoupon && (
              <p className="text-green-600 text-sm mt-1 font-medium">
                ✅ Đã áp dụng mã {appliedCoupon.coupon_code} (Giảm {formatPrice(appliedCoupon.discount_amount)})
              </p>
            )}
          </div>

          {/* Summary */}
          <div className="bg-gray-50 p-4 rounded mb-4">
            <div className="flex justify-between mb-2">
              <span>Tổng số người:</span>
              <span className="font-bold">{totalPeople} người</span>
            </div>
            <div className="flex justify-between mb-2">
              <span>Còn lại:</span>
              <span className={totalPeople > tour.remaining_slots ? 'text-red-500 font-bold' : ''}>
                {tour.remaining_slots} chỗ
              </span>
            </div>
            
            {appliedCoupon && (
              <div className="flex justify-between mb-2 text-gray-600">
                <span>Tổng ban đầu:</span>
                <span className="line-through">{formatPrice(baseTotal)}</span>
              </div>
            )}
            {appliedCoupon && (
              <div className="flex justify-between mb-2 text-green-600 font-medium">
                <span>Được giảm:</span>
                <span>-{formatPrice(appliedCoupon.discount_amount)}</span>
              </div>
            )}

            <div className="flex justify-between text-lg font-bold text-blue-600 pt-2 border-t mt-2">
              <span>Tổng thanh toán:</span>
              <span>{formatPrice(finalTotal)}</span>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-3">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-3 border rounded hover:bg-gray-50"
              disabled={loading}
            >
              Hủy
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-3 bg-blue-600 text-white rounded hover:bg-blue-700 disabled:bg-gray-400"
              disabled={loading || totalPeople > tour.remaining_slots}
            >
              {loading ? 'Đang xử lý...' : 'Xác nhận đặt tour'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default BookingModal;
