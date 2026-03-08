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
    price VARCHAR(100) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    duration VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_location (location)
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
INSERT INTO tours (name, price, description, location, duration) VALUES
('Tour Đà Nẵng - Hội An', '2.000.000đ', 'Khám phá vẻ đẹp của Đà Nẵng và phố cổ Hội An', 'Đà Nẵng', '3 ngày 2 đêm'),
('Tour Hà Nội - Sa Pa', '3.500.000đ', 'Chinh phục đỉnh Fansipan và khám phá Sapa', 'Hà Nội', '4 ngày 3 đêm'),
('Tour Phú Quốc', '5.000.000đ', 'Nghỉ dưỡng tại đảo ngọc Phú Quốc', 'Phú Quốc', '5 ngày 4 đêm'),
('Tour Nha Trang', '3.000.000đ', 'Tắm biển và khám phá vịnh Nha Trang', 'Nha Trang', '3 ngày 2 đêm'),
('Tour Đà Lạt', '2.500.000đ', 'Thành phố ngàn hoa với khí hậu mát mẻ', 'Đà Lạt', '3 ngày 2 đêm');

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
