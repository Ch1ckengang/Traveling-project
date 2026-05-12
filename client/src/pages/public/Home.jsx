import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { getTours } from '../../services/tourService';

const HomePage = () => {
  const [tours, setTours] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [searchKeyword, setSearchKeyword] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    fetchTours();
  }, []);

  const fetchTours = async () => {
    try {
      setLoading(true);
      const data = await getTours();
      setTours(data.slice(0, 6)); // Hiển thị 6 tours đầu tiên
      setError(null);
    } catch (err) {
      console.error('Error fetching tours:', err);
      setError('Không thể tải danh sách tour');
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = (e) => {
    e.preventDefault();
    const kw = searchKeyword.trim();
    if (!kw) return;
    navigate(`/search?q=${encodeURIComponent(kw)}`);
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

  return (
    <div className="container mx-auto px-4 py-8">
      {/* Hero Section */}
      <div className="text-center mb-12">
        <h1 className="text-4xl font-bold mb-4">Chào mừng đến với Traveling</h1>
        <p className="text-xl text-gray-600 mb-6">Khám phá những điểm đến tuyệt vời</p>

        {/* Search Bar */}
        <form onSubmit={handleSearch} className="flex justify-center gap-3 max-w-lg mx-auto">
          <input
            type="text"
            value={searchKeyword}
            onChange={(e) => setSearchKeyword(e.target.value)}
            placeholder="Tìm tour theo địa điểm..."
            className="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:border-blue-500"
            id="home-search-input"
          />
          <button
            type="submit"
            className="px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
            id="home-search-btn"
          >
            Tìm kiếm
          </button>
        </form>
      </div>


      {/* Tours Section */}
      <div className="mb-8">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">Tours Nổi Bật</h2>
          <Link to="/tours" className="text-blue-600 hover:text-blue-800">
            Xem tất cả →
          </Link>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {tours.map((tour) => (
            <Link
              key={tour.id}
              to={`/tours/${tour.id}`}
              className="border rounded-lg overflow-hidden hover:shadow-lg transition-shadow"
            >
              <div className="h-48 bg-gradient-to-br from-blue-100 to-cyan-50 flex items-center justify-center relative overflow-hidden">
                {tour.image_url ? (
                  <img src={tour.image_url} alt={tour.name} className="w-full h-full object-cover" />
                ) : (
                  <div className="flex flex-col items-center gap-1 text-blue-400">
                    <span className="text-3xl">🗺️</span>
                    <span className="text-sm font-medium">{tour.location}</span>
                  </div>
                )}
              </div>
              <div className="p-4">
                <h3 className="font-bold text-lg mb-2">{tour.name}</h3>
                <p className="text-gray-600 text-sm mb-2 line-clamp-2">
                  {tour.description}
                </p>
                <div className="flex justify-between items-center">
                  <span className="text-blue-600 font-bold">{tour.price}</span>
                  <span className="text-sm text-gray-500">{tour.duration}</span>
                </div>
                <div className="text-xs text-green-600 mt-1">Còn {tour.remaining_slots} chỗ</div>
              </div>
            </Link>
          ))}

        </div>
      </div>

      {/* Categories */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12">
        <Link
          to="/tours?category=domestic"
          className="p-6 border rounded-lg text-center hover:shadow-lg transition-shadow"
        >
          <div className="text-4xl mb-4">🇻🇳</div>
          <h3 className="font-bold text-xl mb-2">Du lịch Việt Nam</h3>
          <p className="text-gray-600">Khám phá đất nước</p>
        </Link>

        <Link
          to="/tours?category=international"
          className="p-6 border rounded-lg text-center hover:shadow-lg transition-shadow"
        >
          <div className="text-4xl mb-4">🌍</div>
          <h3 className="font-bold text-xl mb-2">Du lịch Quốc tế</h3>
          <p className="text-gray-600">Khám phá thế giới</p>
        </Link>

        <Link
          to="/tours?category=service"
          className="p-6 border rounded-lg text-center hover:shadow-lg transition-shadow"
        >
          <div className="text-4xl mb-4">✈️</div>
          <h3 className="font-bold text-xl mb-2">Dịch vụ</h3>
          <p className="text-gray-600">Visa, vé máy bay</p>
        </Link>
      </div>
    </div>
  );
};

export default HomePage;
