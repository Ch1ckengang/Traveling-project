# Apr 17Đặc Tả Dự Án: Hệ Thống Đặt Tour Du Lịch Trực Tuyến (B2C)

## 1. TỔNG QUAN DỰ ÁN

### Mô tả

Hệ thống web application cho phép khách hàng tìm kiếm, đặt và thanh toán tour du lịch hoàn toàn trực tuyến. Vận hành theo mô hình B2C — hãng tour trực tiếp phục vụ khách hàng cuối mà không qua trung gian.

### Mục tiêu

- Số hóa toàn bộ quy trình đặt tour từ tìm kiếm đến thanh toán
- Giảm tải công việc thủ công cho nhân viên
- Tăng trải nghiệm khách hàng qua giao diện trực quan
- Quản lý tập trung dữ liệu tour, booking, doanh thu

### Stack Kỹ Thuật

| Layer | Technology |
| --- | --- |
| Frontend | ReactJS (SPA) |
| Backend | Golang + Gin framework |
| Database | PostgreSQL |
| Cache | Redis |
| Auth | JWT (Access + Refresh Token) |
| File Storage | S3-compatible (MinIO/AWS S3) |
| Architecture | Monolithic |

## 2. ACTORS & PHÂN QUYỀN

```text
┌─────────────────────────────────────────────┐
│                  HỆ THỐNG                   │
│                                             │
│  GUEST ──────────────────────────────────  │
│  (Chưa đăng nhập)                          │
│    • Xem danh sách tour                    │
│    • Xem chi tiết tour                     │
│    • Tìm kiếm / lọc tour                   │
│    • Xem review                            │
│                                             │
│  CUSTOMER ───────────────────────────────  │
│  (Khách hàng đã đăng ký)                   │
│    • Tất cả quyền của Guest                │
│    • Đặt tour & thanh toán online          │
│    • Quản lý booking cá nhân              │
│    • Viết review sau khi đi               │
│    • Quản lý hồ sơ cá nhân               │
│    • Dùng mã giảm giá                     │
│                                             │
│  STAFF ──────────────────────────────────  │
│  (Nhân viên hãng tour)                     │
│    • Xử lý booking (xác nhận/hủy)         │
│    • Hỗ trợ thanh toán tại quầy           │
│    • Quản lý tour & lịch tour             │
│    • Xem báo cáo cơ bản                   │
│                                             │
│  ADMIN ──────────────────────────────────  │
│  (Quản trị viên)                           │
│    • Toàn quyền hệ thống                  │
│    • Quản lý tài khoản người dùng         │
│    • Cấu hình hệ thống                    │
│    • Xem tất cả báo cáo & thống kê        │
│    • Quản lý mã giảm giá                  │
└─────────────────────────────────────────────┘
```

## 3. DANH SÁCH MODULE

### MODULE 1 — Authentication & User Management

Mô tả: Xử lý toàn bộ vòng đời tài khoản người dùng

Chức năng:

- Đăng ký tài khoản (email + password)
- Đăng nhập / Đăng xuất
- Refresh JWT token
- Quên mật khẩu (gửi email reset)
- Đổi mật khẩu
- Xem & cập nhật hồ sơ cá nhân
- Upload ảnh đại diện

API Endpoints (Golang Gin):

```text
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/logout
POST   /api/v1/auth/refresh-token
POST   /api/v1/auth/forgot-password
POST   /api/v1/auth/reset-password
GET    /api/v1/users/me
PUT    /api/v1/users/me
PUT    /api/v1/users/me/password
PUT    /api/v1/users/me/avatar
```

### MODULE 2 — Tour Management

Mô tả: Quản lý danh mục tour và thông tin chi tiết

Chức năng:

Phía khách hàng (public):

- Xem danh sách tour (phân trang)
- Tìm kiếm tour theo từ khóa
- Lọc tour theo: điểm đến, giá, thời gian, loại hình, phương tiện
- Sắp xếp theo: giá tăng/giảm, mới nhất, phổ biến nhất
- Xem chi tiết tour (mô tả, lịch trình, dịch vụ, ảnh)
- Xem các lịch tour còn chỗ
- Xem review của tour

Phía Staff/Admin:

- Tạo mới tour
- Cập nhật thông tin tour
- Upload ảnh tour
- Ẩn/Hiện tour
- Xóa tour (soft delete)
- Thêm / sửa / xóa lịch tour
- Quản lý số lượng chỗ theo lịch

API Endpoints:

```text
GET    /api/v1/tours                    (public - list + filter)
GET    /api/v1/tours/:slug              (public - detail)
GET    /api/v1/tours/:id/schedules      (public - available schedules)
GET    /api/v1/tours/:id/reviews        (public)

POST   /api/v1/admin/tours              (STAFF+)
PUT    /api/v1/admin/tours/:id          (STAFF+)
DELETE /api/v1/admin/tours/:id          (ADMIN)
POST   /api/v1/admin/tours/:id/images   (STAFF+)

POST   /api/v1/admin/schedules          (STAFF+)
PUT    /api/v1/admin/schedules/:id      (STAFF+)
DELETE /api/v1/admin/schedules/:id      (ADMIN)
```

### MODULE 3 — Booking Management

Mô tả: Xử lý toàn bộ luồng đặt tour từ tạo booking đến hoàn tất

Chức năng:

Phía khách hàng:

- Tạo booking (chọn lịch tour, số lượng khách, thông tin hành khách)
- Xem danh sách booking của bản thân
- Xem chi tiết booking
- Hủy booking (theo policy)
- Tải hóa đơn PDF

Phía Staff/Admin:

- Xem tất cả booking
- Xác nhận booking
- Hủy booking
- Tìm kiếm booking theo mã, SĐT, tên khách
- Xem danh sách khách theo lịch tour

Booking Flow:

```text
Khách chọn tour & lịch
        ↓
   Điền thông tin
   hành khách
        ↓
  Tạo Booking
  (status: PENDING)
        ↓
   Chuyển sang
   thanh toán
        ↓
  Thanh toán thành công
        ↓
  Booking CONFIRMED
  + Gửi email xác nhận
        ↓
     Tour kết thúc
        ↓
  Booking COMPLETED
  → Mở khóa Review
```

API Endpoints:

```text
POST   /api/v1/bookings                     (CUSTOMER)
GET    /api/v1/bookings                     (CUSTOMER - own bookings)
GET    /api/v1/bookings/:code               (CUSTOMER - own)
POST   /api/v1/bookings/:code/cancel        (CUSTOMER)
GET    /api/v1/bookings/:code/invoice       (CUSTOMER - download PDF)

GET    /api/v1/admin/bookings               (STAFF+)
GET    /api/v1/admin/bookings/:code         (STAFF+)
PUT    /api/v1/admin/bookings/:code/confirm (STAFF+)
PUT    /api/v1/admin/bookings/:code/cancel  (STAFF+)
GET    /api/v1/admin/schedules/:id/passengers (STAFF+)
```

### MODULE 4 — Payment

Mô tả: Xử lý thanh toán online và tại quầy

Chức năng:

Thanh toán online:

- Khởi tạo giao dịch thanh toán
- Tích hợp VNPay (redirect flow)
- Xử lý callback/webhook từ cổng thanh toán
- Xác nhận thanh toán thành công / thất bại

Thanh toán tại quầy:

- Staff tra cứu booking theo SĐT / mã booking
- Ghi nhận thanh toán tiền mặt / chuyển khoản
- In hóa đơn

Hoàn tiền:

- Tạo yêu cầu hoàn tiền
- Admin duyệt hoàn tiền
- Ghi nhận hoàn tiền

Payment Flow:

```text
Booking PENDING
      ↓
POST /payments/initiate
      ↓
  Redirect sang
  VNPay / MoMo
      ↓
Khách thanh toán
      ↓
Gateway callback
POST /payments/callback
      ↓
Verify signature
      ↓
Cập nhật Payment SUCCESS
+ Booking CONFIRMED
+ Gửi email
```

API Endpoints:

```text
POST   /api/v1/payments/initiate         (CUSTOMER)
POST   /api/v1/payments/callback/vnpay   (Webhook - public)
GET    /api/v1/payments/:code            (CUSTOMER - own)

POST   /api/v1/admin/payments/offline    (STAFF+ - thanh toán quầy)
POST   /api/v1/admin/payments/:id/refund (ADMIN)
GET    /api/v1/admin/payments            (STAFF+)
```

### MODULE 5 — Review & Rating

Mô tả: Hệ thống đánh giá tour sau khi sử dụng dịch vụ

Chức năng:

- Viết review (chỉ khách đã hoàn thành booking)
- Upload ảnh kèm review (tối đa 5 ảnh)
- Chỉnh sửa review trong 7 ngày
- Xem danh sách review của tour
- Phân trang review
- Lọc review theo số sao
- Admin duyệt / ẩn review
- Admin phản hồi review

API Endpoints:

```text
GET    /api/v1/tours/:id/reviews         (public)
POST   /api/v1/reviews                   (CUSTOMER - verified booking)
PUT    /api/v1/reviews/:id               (CUSTOMER - own, within 7 days)

GET    /api/v1/admin/reviews             (STAFF+)
PUT    /api/v1/admin/reviews/:id/publish (STAFF+)
PUT    /api/v1/admin/reviews/:id/hide    (STAFF+)
POST   /api/v1/admin/reviews/:id/reply   (STAFF+)
```

### MODULE 6 — Promotion & Coupon

Mô tả: Quản lý mã giảm giá và chương trình khuyến mãi

Chức năng:

- Tạo mã giảm giá (% hoặc số tiền cố định)
- Giới hạn lượt dùng, thời gian hiệu lực, giá trị đơn tối thiểu
- Khách nhập mã khi đặt tour → preview số tiền giảm
- Validate mã trước khi thanh toán
- Thống kê sử dụng mã

API Endpoints:

```text
POST   /api/v1/coupons/validate          (CUSTOMER - check coupon)

POST   /api/v1/admin/coupons             (ADMIN)
GET    /api/v1/admin/coupons             (ADMIN)
PUT    /api/v1/admin/coupons/:id         (ADMIN)
DELETE /api/v1/admin/coupons/:id         (ADMIN)
GET    /api/v1/admin/coupons/:id/usage   (ADMIN)
```

### MODULE 7 — Notification

Mô tả: Gửi thông báo tự động đến khách hàng

Chức năng:

- Email xác nhận đặt tour thành công
- Email xác nhận thanh toán
- Email nhắc nhở trước ngày khởi hành (3 ngày, 1 ngày)
- Email thông báo hủy tour
- Email reset mật khẩu
- In-app notification (realtime qua WebSocket hoặc polling)

Trigger Events:

```text
booking.confirmed     → Email xác nhận booking
payment.success       → Email xác nhận thanh toán + Invoice PDF
payment.failed        → Email thông báo thất bại
booking.cancelled     → Email thông báo hủy
departure.reminder    → Cronjob chạy hàng ngày
password.reset        → Email link reset
```

### MODULE 8 — Admin Dashboard & Report

Mô tả: Giao diện quản trị và báo cáo thống kê

Chức năng:

- Dashboard tổng quan (doanh thu, booking, khách hàng mới)
- Biểu đồ doanh thu theo ngày / tháng / năm
- Báo cáo tour phổ biến nhất
- Báo cáo tỷ lệ hủy booking
- Quản lý tài khoản người dùng (kích hoạt / khóa)
- Xuất báo cáo Excel/CSV

API Endpoints:

```text
GET    /api/v1/admin/dashboard/summary
GET    /api/v1/admin/dashboard/revenue?period=monthly
GET    /api/v1/admin/dashboard/top-tours
GET    /api/v1/admin/dashboard/booking-stats

GET    /api/v1/admin/users
PUT    /api/v1/admin/users/:id/status
GET    /api/v1/admin/reports/export?type=booking&from=...&to=...
```

## 4. KIẾN TRÚC MONOLITHIC — GOLANG GIN + REACTJS

### Cấu trúc thư mục Backend (Golang)

```text
backend/
├── cmd/
│   └── server/
│       └── main.go              # Entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Load env, app config
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── rate_limit.go
│   ├── router/
│   │   └── router.go            # Đăng ký tất cả routes
│   ├── modules/
│   │   ├── auth/
│   │   │   ├── handler.go       # Gin handlers
│   │   │   ├── service.go       # Business logic
│   │   │   ├── repository.go    # DB queries
│   │   │   └── dto.go           # Request/Response structs
│   │   ├── tour/
│   │   │   ├── handler.go
│   │   │   ├── service.go
│   │   │   ├── repository.go
│   │   │   └── dto.go
│   │   ├── booking/
│   │   ├── payment/
│   │   ├── review/
│   │   ├── coupon/
│   │   └── admin/
│   ├── models/
│   │   ├── user.go              # GORM models
│   │   ├── tour.go
│   │   ├── booking.go
│   │   ├── payment.go
│   │   └── ...
│   ├── pkg/
│   │   ├── database/
│   │   │   └── postgres.go      # DB connection
│   │   ├── cache/
│   │   │   └── redis.go
│   │   ├── email/
│   │   │   └── sender.go
│   │   ├── storage/
│   │   │   └── s3.go
│   │   ├── payment/
│   │   │   └── vnpay.go
│   │   └── jwt/
│   │       └── jwt.go
│   └── utils/
│       ├── response.go          # Chuẩn hóa API response
│       ├── validator.go
│       └── pagination.go
├── migrations/                  # SQL migration files
├── .env
├── go.mod
└── go.sum
```

### Cấu trúc thư mục Frontend (ReactJS)

```text
frontend/
├── public/
├── src/
│   ├── assets/
│   ├── components/              # Reusable UI components
│   │   ├── common/
│   │   │   ├── Button/
│   │   │   ├── Input/
│   │   │   ├── Modal/
│   │   │   └── Pagination/
│   │   ├── tour/
│   │   │   ├── TourCard/
│   │   │   ├── TourFilter/
│   │   │   └── TourGallery/
│   │   └── booking/
│   │       ├── BookingForm/
│   │       └── BookingStatus/
│   ├── pages/                   # Route-level pages
│   │   ├── Home/
│   │   ├── TourList/
│   │   ├── TourDetail/
│   │   ├── Booking/
│   │   │   ├── BookingForm/
│   │   │   ├── BookingPayment/
│   │   │   └── BookingConfirm/
│   │   ├── Profile/
│   │   ├── MyBookings/
│   │   └── Admin/
│   │       ├── Dashboard/
│   │       ├── Tours/
│   │       ├── Bookings/
│   │       └── Users/
│   ├── services/                # API call functions
│   │   ├── authService.js
│   │   ├── tourService.js
│   │   ├── bookingService.js
│   │   └── paymentService.js
│   ├── store/                   # State management (Redux/Zustand)
│   │   ├── authSlice.js
│   │   ├── tourSlice.js
│   │   └── bookingSlice.js
│   ├── hooks/                   # Custom hooks
│   │   ├── useAuth.js
│   │   └── useTours.js
│   ├── utils/
│   │   ├── axiosInstance.js     # Axios config + interceptor
│   │   ├── formatters.js
│   │   └── constants.js
│   ├── routes/
│   │   └── AppRouter.jsx        # React Router config
│   └── App.jsx
├── .env
└── package.json
```

### API Response Standard

```json
{
  "success": true,
  "message": "Đặt tour thành công",
  "data": { ... },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  },
  "error": null
}
```

### Error Response Standard

```json
{
  "success": false,
  "message": "Validation failed",
  "data": null,
  "error": {
    "code": "VALIDATION_ERROR",
    "details": [
      { "field": "email", "message": "Email không hợp lệ" }
    ]
  }
}
```
