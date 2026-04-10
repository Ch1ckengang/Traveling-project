-- =============================================
-- PostgreSQL init script for Traveling App
-- =============================================
-- This file is intended to run INSIDE database travel_db
-- (works in pgAdmin Query Tool without psql meta-commands).
--
-- If travel_db does not exist yet, run this once in database postgres:
-- CREATE DATABASE travel_db;

-- =============================================
-- users table
-- =============================================
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- =============================================
-- tours table
-- =============================================
CREATE TABLE IF NOT EXISTS tours (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'domestic',
    price VARCHAR(100) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    country VARCHAR(255) NOT NULL DEFAULT 'Viet Nam',
    duration VARCHAR(100),
    departure_date VARCHAR(50) NOT NULL DEFAULT '',
    remaining_slots INT NOT NULL DEFAULT 30,
    itinerary TEXT,
    services TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tours_location ON tours(location);
CREATE INDEX IF NOT EXISTS idx_tours_type ON tours(type);
CREATE INDEX IF NOT EXISTS idx_tours_country ON tours(country);

-- =============================================
-- bookings table
-- =============================================
CREATE TABLE IF NOT EXISTS bookings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL DEFAULT 0,
    tour_id BIGINT NOT NULL,
    full_name VARCHAR(255) NOT NULL DEFAULT '',
    phone VARCHAR(50) NOT NULL DEFAULT '',
    email VARCHAR(255) NOT NULL DEFAULT '',
    adult_count INT NOT NULL DEFAULT 1,
    child_count INT NOT NULL DEFAULT 0,
    infant_count INT NOT NULL DEFAULT 0,
    quantity INT NOT NULL,
    travel_date VARCHAR(50) NOT NULL,
    total_amount BIGINT NOT NULL DEFAULT 0,
    booking_code VARCHAR(120) NOT NULL DEFAULT '',
    payment_status VARCHAR(50) NOT NULL DEFAULT 'unpaid',
    note TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'booked',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_bookings_booking_code UNIQUE (booking_code)
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_tour_id ON bookings(tour_id);

-- Ensure backend connection user can access schema/table/sequence.
GRANT CONNECT ON DATABASE travel_db TO postgres;
GRANT USAGE ON SCHEMA public TO postgres;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO postgres;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO postgres;
ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO postgres;

-- Clean existing duplicated tours from previous runs (keep the oldest row by name).
WITH duplicated AS (
    SELECT id,
           ROW_NUMBER() OVER (PARTITION BY name ORDER BY id) AS rn
    FROM tours
)
DELETE FROM tours
WHERE id IN (
    SELECT id
    FROM duplicated
    WHERE rn > 1
);

-- =============================================
-- Seed sample users (password is bcrypt hash for: 123456)
-- =============================================
INSERT INTO users (name, email, password)
VALUES
('Nguyen Van A', 'test@example.com', '$2a$10$mI21ZHy1/IJz9M0zkhjsgOAliwChmmA5RNU2HaCghFUvf669xQVnS'),
('Tran Thi B', 'user@example.com', '$2a$10$mI21ZHy1/IJz9M0zkhjsgOAliwChmmA5RNU2HaCghFUvf669xQVnS')
ON CONFLICT (email) DO NOTHING;

INSERT INTO tours (name, type, price, description, location, country, duration)
SELECT src.name, src.type, src.price, src.description, src.location, src.country, src.duration
FROM (
    VALUES
    ('Tour Da Nang - Hoi An', 'domestic', '2.000.000d', 'Kham pha ve dep cua Da Nang va pho co Hoi An', 'Da Nang', 'Viet Nam', '3 ngay 2 dem'),
    ('Tour Ha Noi - Sa Pa', 'domestic', '3.500.000d', 'Chinh phuc dinh Fansipan va kham pha Sapa', 'Ha Noi', 'Viet Nam', '4 ngay 3 dem'),
    ('Tour Phu Quoc', 'domestic', '5.000.000d', 'Nghi duong tai dao ngoc Phu Quoc', 'Phu Quoc', 'Viet Nam', '5 ngay 4 dem'),
    ('Tour Nha Trang - Da Lat', 'domestic', '4.800.000d', 'Ket hop du lich bien va nghi duong cao nguyen', 'Nha Trang', 'Viet Nam', '4 ngay 3 dem'),
    ('Tour Bangkok - Pattaya', 'international', '8.200.000d', 'Tham quan chua vang va thanh pho bien Pattaya', 'Bangkok', 'Thai Lan', '5 ngay 4 dem'),
    ('Tour Seoul Mua Hoa', 'international', '12.500.000d', 'Kham pha Seoul va van hoa Han Quoc hien dai', 'Seoul', 'Han Quoc', '6 ngay 5 dem'),
    ('Tour Tokyo - Nui Phu Si', 'international', '15.900.000d', 'Trai nghiem Tokyo va bieu tuong Nhat Ban', 'Tokyo', 'Nhat Ban', '6 ngay 5 dem'),
    ('Tour Paris - Lyon', 'international', '18.500.000d', 'Hanh trinh Phap voi kien truc va am thuc chau Au', 'Paris', 'Phap', '7 ngay 6 dem'),
    ('Tour Singapore - Sentosa', 'international', '9.600.000d', 'Kham pha dao quoc su tu va Sentosa', 'Singapore', 'Singapore', '4 ngay 3 dem'),
    ('Tour Bali - Ubud', 'international', '10.800.000d', 'Nghi duong bien dao Bali va tham quan Ubud', 'Bali', 'Indonesia', '5 ngay 4 dem'),
    ('Tour Sydney - Melbourne', 'international', '21.900.000d', 'Kham pha hai thanh pho noi tieng cua Uc', 'Sydney', 'Uc', '7 ngay 6 dem'),
    ('Tour Dubai - Abu Dhabi', 'international', '19.500.000d', 'Trai nghiem Trung Dong xa hoa va doc dao', 'Dubai', 'UAE', '6 ngay 5 dem')
) AS src(name, type, price, description, location, country, duration)
WHERE NOT EXISTS (
    SELECT 1
    FROM tours t
    WHERE t.name = src.name
);

-- =============================================
-- Quick summary
-- =============================================
SELECT 'Database created successfully!' AS status,
       (SELECT COUNT(*) FROM users) AS total_users,
       (SELECT COUNT(*) FROM tours) AS total_tours,
       (SELECT COUNT(*) FROM bookings) AS total_bookings;
