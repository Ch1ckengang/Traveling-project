-- =============================================
-- Script tạo Database cho Traveling App
-- =============================================

-- Xóa database nếu đã tồn tại và tạo mới
DROP DATABASE IF EXISTS travel_db;
CREATE DATABASE travel_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Sử dụng database vừa tạo
USE travel_db;

-- =============================================
-- Tạo bảng USERS
-- =============================================
CREATE TABLE users (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- Tạo bảng TOURS
-- =============================================
CREATE TABLE tours (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'domestic',
    price VARCHAR(100) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    country VARCHAR(255) NOT NULL DEFAULT 'Việt Nam',
    duration VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_location (location),
    INDEX idx_type (type),
    INDEX idx_country (country)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- =============================================
-- Insert dữ liệu mẫu vào bảng USERS
-- =============================================
INSERT INTO users (name, email, password) VALUES
('Nguyễn Văn A', 'test@example.com', '123456'),
('Trần Thị B', 'user@example.com', '123456'),
('Lê Văn C', 'admin@example.com', '123456');

-- =============================================
-- Insert dữ liệu mẫu vào bảng TOURS
-- =============================================
INSERT INTO tours (name, type, price, description, location, country, duration) VALUES
('Tour Đà Nẵng - Hội An', 'domestic', '2.000.000đ', 'Khám phá vẻ đẹp của Đà Nẵng và phố cổ Hội An', 'Đà Nẵng', 'Việt Nam', '3 ngày 2 đêm'),
('Tour Hà Nội - Sa Pa', 'domestic', '3.500.000đ', 'Chinh phục đỉnh Fansipan và khám phá Sapa', 'Hà Nội', 'Việt Nam', '4 ngày 3 đêm'),
('Tour Phú Quốc', 'domestic', '5.000.000đ', 'Nghỉ dưỡng tại đảo ngọc Phú Quốc', 'Phú Quốc', 'Việt Nam', '5 ngày 4 đêm'),
('Tour Nha Trang - Đà Lạt', 'domestic', '4.800.000đ', 'Kết hợp du lịch biển và nghỉ dưỡng cao nguyên', 'Nha Trang', 'Việt Nam', '4 ngày 3 đêm'),
('Tour Bangkok - Pattaya', 'international', '8.200.000đ', 'Tham quan chùa vàng và thành phố biển Pattaya', 'Bangkok', 'Thái Lan', '5 ngày 4 đêm'),
('Tour Seoul Mùa Hoa', 'international', '12.500.000đ', 'Khám phá Seoul và văn hóa Hàn Quốc hiện đại', 'Seoul', 'Hàn Quốc', '6 ngày 5 đêm'),
('Tour Tokyo - Núi Phú Sĩ', 'international', '15.900.000đ', 'Trải nghiệm Tokyo và biểu tượng Nhật Bản', 'Tokyo', 'Nhật Bản', '6 ngày 5 đêm'),
('Tour Paris - Lyon', 'international', '18.500.000đ', 'Hành trình Pháp với kiến trúc và ẩm thực châu Âu', 'Paris', 'Pháp', '7 ngày 6 đêm'),
('Tour Singapore - Sentosa', 'international', '9.600.000đ', 'Khám phá đảo quốc sư tử và Sentosa', 'Singapore', 'Singapore', '4 ngày 3 đêm'),
('Tour Bali - Ubud', 'international', '10.800.000đ', 'Nghỉ dưỡng biển đảo Bali và tham quan Ubud', 'Bali', 'Indonesia', '5 ngày 4 đêm'),
('Tour Sydney - Melbourne', 'international', '21.900.000đ', 'Khám phá hai thành phố nổi tiếng của Úc', 'Sydney', 'Úc', '7 ngày 6 đêm'),
('Tour Dubai - Abu Dhabi', 'international', '19.500.000đ', 'Trải nghiệm Trung Đông xa hoa và độc đáo', 'Dubai', 'UAE', '6 ngày 5 đêm');

-- =============================================
-- Kiểm tra dữ liệu đã insert
-- =============================================
SELECT 'Users Table:' as '';
SELECT * FROM users;

SELECT 'Tours Table:' as '';
SELECT * FROM tours;

-- =============================================
-- Thống kê
-- =============================================
SELECT 
    'Database created successfully!' as 'Status',
    (SELECT COUNT(*) FROM users) as 'Total Users',
    (SELECT COUNT(*) FROM tours) as 'Total Tours';
