# Tổng Quan Hệ Thống Traveling (Bản Tóm Tắt Tiếng Việt)

Cập nhật: 2026-04-13

Tài liệu này là bản tổng hợp nhanh bằng tiếng Việt. Tài liệu chi tiết theo từng module được viết bằng tiếng Anh trong `doc/modules/`.

## 1. Kiến trúc tổng thể
Hệ thống backend đang theo hướng monolithic theo module:
- `auth`
- `tour`
- `booking`

Công nghệ chính:
- Gin (HTTP API)
- GORM (ORM)
- PostgreSQL (database)

Cấu trúc hiện tại:
- `server/cmd/server/main.go`: khởi tạo app + route + migrate + seed
- `server/database/postgres.go`: kết nối DB
- `server/domain/*`: entity + DTO dùng chung
- `server/internal/auth/*`: đăng ký, OTP, đăng nhập, quên mật khẩu
- `server/internal/tour/*`: danh sách tour + bộ lọc + sắp xếp
- `server/internal/booking/*`: đặt tour, lịch sử, hủy đặt tour
- `server/internal/shared/errors.go`: lỗi dùng chung

## 2. Luồng auth hiện tại
1. User đăng ký bằng email `@gmail.com`.
2. Backend tạo tài khoản với `is_email_verified = false`.
3. Backend tạo OTP 6 số (hết hạn 3 phút) và gửi cho email đã đăng ký.
4. User xác thực OTP thành công -> `is_email_verified = true`.
5. User chỉ được đăng nhập khi đã xác thực email.

Ghi chú:
- Hiện tại OTP đang ở dev mode (lưu in-memory).
- Để production, cần đưa OTP sang Redis/DB và tích hợp email service thực.

## 3. Luồng tour + booking
- Tour:
  - API công khai để lấy danh sách tour.
  - Hỗ trợ bộ lọc theo city, duration, price và sort.
- Booking:
  - Tạo booking với validate đầy đủ (tour, số lượng, ngày đi, liên hệ).
  - Trừ chỗ còn lại của tour theo transaction.
  - Có API xem lịch sử booking và hủy booking.
  - Khi hủy booking sẽ cộng lại số chỗ.

## 4. API chính
Base prefix: `/v1`

- Auth:
  - `POST /api/register`
  - `POST /api/login`
  - `POST /api/otp/send`
  - `POST /api/otp/verify`
  - `POST /api/password/forgot`
  - `PUT /api/users/:id`
- Tour:
  - `GET /api/tours`
  - `GET /api/tours/domestic`
  - `GET /api/tours/international`
- Booking:
  - `POST /api/bookings`
  - `GET /api/users/:id/bookings`
  - `PUT /api/users/:id/bookings/:bookingId/cancel`

## 5. Tài liệu chi tiết
- `doc/modules/system-architecture.md`
- `doc/modules/auth-module.md`
- `doc/modules/tour-module.md`
- `doc/modules/booking-module.md`
