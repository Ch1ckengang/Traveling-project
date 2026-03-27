import '../../styles/SearchBar.css';

const SearchBar = ({ filters, onFiltersChange, onSearch }) => {
  const handleChange = (field, value) => {
    onFiltersChange({
      ...filters,
      [field]: value
    });
  };

  return (
    <div className="search-bar-wrapper">
      <div className="search-inputs">
        <div className="input-group">
          <label>Tên thành phố</label>
          <input
            type="text"
            placeholder="Nhập địa điểm..."
            value={filters.city}
            onChange={(event) => handleChange('city', event.target.value)}
          />
        </div>

        <div className="input-group">
          <label>Thời gian</label>
          <select value={filters.duration} onChange={(event) => handleChange('duration', event.target.value)}>
            <option value="all">Tất cả</option>
            <option value="short">2 - 3 ngày</option>
            <option value="medium">4 - 5 ngày</option>
            <option value="long">Từ 6 ngày</option>
          </select>
        </div>

        <div className="input-group">
          <label>Giá tiền</label>
          <select value={filters.price} onChange={(event) => handleChange('price', event.target.value)}>
            <option value="all">Tất cả mức giá</option>
            <option value="low">Dưới 3 triệu</option>
            <option value="mid">3 - 7 triệu</option>
            <option value="high">Trên 7 triệu</option>
          </select>
        </div>

        <div className="input-group">
          <label>Sắp xếp</label>
          <select value={filters.sort} onChange={(event) => handleChange('sort', event.target.value)}>
            <option value="default">Mặc định</option>
            <option value="price_asc">Giá tăng dần</option>
            <option value="price_desc">Giá giảm dần</option>
            <option value="duration_asc">Thời lượng ngắn trước</option>
            <option value="duration_desc">Thời lượng dài trước</option>
            <option value="name_asc">Tên A-Z</option>
            <option value="name_desc">Tên Z-A</option>
            <option value="latest">Mới cập nhật</option>
          </select>
        </div>
      </div>

      <div className="search-btn-container">
        <button className="btn-search" type="button" onClick={onSearch}>Tìm kiếm</button>
      </div>
    </div>
  );
};

export default SearchBar;