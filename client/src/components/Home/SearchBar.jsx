import '../../styles/SearchBar.css';

const SearchBar = () => {
  return (
    <div className="search-bar-wrapper">
      <div className="search-inputs">
        <div className="input-group">
          <label>Tỉnh/ thành phố</label>
          <input type="text" placeholder="Nhập địa điểm..." />
        </div>
        
        <div className="input-group">
          <label>Thời gian</label>
          <input type="date" />
        </div>
        
        <div className="input-group">
          <label>Giá tiền</label>
          <select>
            <option>Tất cả mức giá</option>
            <option>Dưới 2 triệu</option>
            <option>2 - 5 triệu</option>
            <option>Trên 5 triệu</option>
          </select>
        </div>
      </div>
      
      <div className="search-btn-container">
        <button className="btn-search">Tìm kiếm</button>
      </div>
    </div>
  );
};

export default SearchBar;