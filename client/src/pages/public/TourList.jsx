import { useEffect, useState } from 'react';
import { Link, useSearchParams } from 'react-router-dom';
import { getTours, getDomesticTours, getInternationalTours } from '../../services/tourService';

const TourListPage = () => {
  const [searchParams] = useSearchParams();
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const category = searchParams.get('category') || 'all';

  useEffect(() => {
    fetchTours();
  }, [category]);

  const fetchTours = async () => {
    try {
      setLoading(true);
      let data;
      
      if (category === 'domestic') {
        data = await getDomesticTours();
      } else if (category === 'international') {
        data = await getInternationalTours();
      } else {
        data = await getTours();
      }
      
      if (data && data.success && Array.isArray(data.data)) {
        setTours(data.data);
      } else if (Array.isArray(data)) {
        setTours(data); // Fallback if API hasn't been standardized
      } else {
        setTours([]);
      }
      setError(null);
    } catch (err) {
      console.error('Error fetching tours:', err);
      setError('Không thể tải danh sách tour');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center">Đang tải...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="text-center text-red-500">{error}</div>
      </div>
    );
  }

  const getCategoryTitle = () => {
    switch (category) {
      case 'domestic':
        return 'Du lịch Việt Nam';
      case 'international':
        return 'Du lịch Quốc tế';
      case 'service':
        return 'Dịch vụ';
      default:
        return 'Tất cả Tours';
    }
  };

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold mb-4">{getCategoryTitle()}</h1>
        <p className="text-gray-600">Tìm thấy {tours.length} tours</p>
      </div>

      {/* Filter Tabs */}
      <div className="flex gap-4 mb-8 border-b">
        <Link
          to="/tours"
          className={`pb-2 px-4 ${
            category === 'all' ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-600'
          }`}
        >
          Tất cả
        </Link>
        <Link
          to="/tours?category=domestic"
          className={`pb-2 px-4 ${
            category === 'domestic' ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-600'
          }`}
        >
          Việt Nam
        </Link>
        <Link
          to="/tours?category=international"
          className={`pb-2 px-4 ${
            category === 'international' ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-600'
          }`}
        >
          Quốc tế
        </Link>
        <Link
          to="/tours?category=service"
          className={`pb-2 px-4 ${
            category === 'service' ? 'border-b-2 border-blue-600 text-blue-600' : 'text-gray-600'
          }`}
        >
          Dịch vụ
        </Link>
      </div>

      {/* Tours Grid */}
      {tours.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          Không tìm thấy tour nào
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {tours.map((tour) => (
            <Link
              key={tour.id}
              to={`/tours/${tour.id}`}
              className="border rounded-lg overflow-hidden hover:shadow-lg transition-shadow"
            >
              <div className="h-48 bg-gray-200 flex items-center justify-center">
                <span className="text-gray-400">📷 {tour.location}</span>
              </div>
              <div className="p-4">
                <div className="flex items-center gap-2 mb-2">
                  <span className="text-xs px-2 py-1 bg-blue-100 text-blue-600 rounded">
                    {tour.type === 'domestic' ? 'Việt Nam' : tour.type === 'international' ? 'Quốc tế' : 'Dịch vụ'}
                  </span>
                  <span className="text-xs text-gray-500">{tour.country}</span>
                </div>
                <h3 className="font-bold text-lg mb-2">{tour.name}</h3>
                <p className="text-gray-600 text-sm mb-3 line-clamp-2">
                  {tour.description}
                </p>
                <div className="flex justify-between items-center">
                  <span className="text-blue-600 font-bold">{tour.price}</span>
                  <span className="text-sm text-gray-500">{tour.duration}</span>
                </div>
                <div className="mt-2 text-sm text-gray-500">
                  Còn {tour.remaining_slots} chỗ
                </div>
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
};

export default TourListPage;
