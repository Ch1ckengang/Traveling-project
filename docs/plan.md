# 📋 KẾ HOẠCH HOÀN THÀNH DỰ ÁN TRAVELING

**Ngày tạo:** 2026-05-13  
**Mục tiêu:** Hoàn thành 100% tất cả 8 modules theo đặc tả  
**Tiến độ hiện tại:** ~30-35%

---

## 🔴 PHASE 0: NỀN TẢNG HẠ TẦNG (Ưu tiên CAO NHẤT)

> Phải làm trước vì TẤT CẢ modules khác đều phụ thuộc vào đây.

### 0.1 — Hệ thống phân quyền RBAC
**Vấn đề:** User model không có trường `role`. Không phân biệt Guest/Customer/Staff/Admin.

**Tasks:**
- [ ] Thêm trường `Role string` vào `domain/user.go` (giá trị: `customer`, `staff`, `admin`)
- [ ] Migration: thêm cột `role` với default `customer`
- [ ] Tạo `internal/shared/role_middleware.go` — middleware kiểm tra role
- [ ] Cập nhật seed data: tạo user admin mẫu
- [ ] Cập nhật `AuthContext.jsx` frontend để lưu role

### 0.2 — Chuẩn hóa Tour Price (string → int64)
**Vấn đề:** `Price` đang lưu dạng string `"2.000.000đ"` → không filter/sort/tính toán được.

**Tasks:**
- [ ] Đổi `Price string` → `PriceAmount int64` (đơn vị VND)
- [ ] Thêm trường `Slug`, `Rating`, `ReviewCount`, `IsActive`
- [ ] Cập nhật seed data, repository, service, handler
- [ ] Cập nhật frontend format hiển thị giá

### 0.3 — Chuẩn hóa API Response
- [ ] Tạo `internal/shared/response.go` — helper response dùng chung
- [ ] Tạo `internal/shared/pagination.go` — pagination helper
- [ ] Áp dụng cho tất cả handlers

---

## 🔴 PHASE 1: HOÀN THIỆN PAYMENT MODULE (Ưu tiên 1)

> Đã có foundation (domain, config, repository), cần service + handler + frontend.

### 1.1 — VNPay Client
- [ ] Tạo `internal/payment/vnpay_client.go`
- [ ] HMAC-SHA512 signature generation + verification
- [ ] Payment URL generation
- [ ] Unit tests

### 1.2 — Payment Service
- [ ] Tạo `internal/payment/service.go`
- [ ] `InitiatePayment`, `ProcessWebhook`, `HandleReturn`, `GetPaymentStatus`, `RetryPayment`
- [ ] Cập nhật booking status khi payment success
- [ ] Gửi email xác nhận thanh toán

### 1.3 — Payment Handler + Routes
- [ ] Tạo `internal/payment/handler.go`
- [ ] `POST /v1/api/payments/initiate` | `GET /payments/return` | `POST /payments/webhook` | `GET /payments/status/:ref`
- [ ] Đăng ký routes trong `main.go`

### 1.4 — Payment Frontend
- [ ] Tạo `services/paymentService.js`
- [ ] Tạo `PaymentButton`, `PaymentStatus` components
- [ ] Tạo `PaymentReturn` page + route
- [ ] Tích hợp vào BookingDetail

---

## 🔴 PHASE 2: ADMIN TOUR & BOOKING (Ưu tiên 2)

### 2.1 — Admin Tour CRUD
- [ ] `POST /v1/api/admin/tours` | `PUT /:id` | `DELETE /:id` | `PUT /:id/toggle`
- [ ] Role middleware: Staff+ required

### 2.2 — Admin Booking Management
- [ ] `GET /v1/api/admin/bookings` | `GET /:code` | `PUT /:code/confirm` | `PUT /:code/cancel`
- [ ] Search theo mã/SĐT/tên

### 2.3 — Admin Frontend Pages
- [ ] Build `AdminToursPage` — CRUD table + form
- [ ] Build `AdminBookingsPage` — List + actions
- [ ] Build `AdminUsersPage` — List + lock/unlock

---

## 🟡 PHASE 3: REVIEW & RATING MODULE (Ưu tiên 3)

### 3.1 — Review Backend
- [ ] Tạo `domain/review.go`, `internal/review/{repository,service,handler}.go`
- [ ] Public: `GET /tours/:id/reviews`
- [ ] Customer: `POST /reviews` (chỉ COMPLETED booking) | `PUT /reviews/:id` (7 ngày)
- [ ] Admin: `GET /admin/reviews` | `PUT /:id/publish` | `PUT /:id/hide` | `POST /:id/reply`
- [ ] Cập nhật Tour rating khi có review mới

### 3.2 — Review Frontend
- [ ] Build `WriteReviewPage`, `ReviewList` component
- [ ] Build `AdminReviewsPage`
- [ ] Tích hợp vào TourDetail

---

## 🟡 PHASE 4: COUPON MODULE (Ưu tiên 4)

### 4.1 — Coupon Backend
- [ ] Tạo `domain/coupon.go`, `internal/coupon/{repository,service,handler}.go`
- [ ] Customer: `POST /coupons/validate`
- [ ] Admin CRUD: `POST/GET/PUT/DELETE /admin/coupons` + `GET /:id/usage`

### 4.2 — Coupon Frontend
- [ ] `CouponInput` component trong booking form
- [ ] Build `AdminCouponsPage`

---

## 🟡 PHASE 5: ADMIN DASHBOARD & REPORT (Ưu tiên 5)

### 5.1 — Dashboard Backend
- [ ] `GET /admin/dashboard/summary` | `/revenue` | `/top-tours` | `/booking-stats`
- [ ] `GET /admin/users` | `PUT /admin/users/:id/status`
- [ ] `GET /admin/reports/export` — CSV export

### 5.2 — Dashboard Frontend
- [ ] Build `AdminDashboardPage` — stats cards + charts (Recharts)
- [ ] Revenue chart, top tours, booking stats

---

## 🟢 PHASE 6: NOTIFICATION MODULE (Ưu tiên 6)

- [ ] Email templates: booking confirmed, payment success/fail, booking cancelled, departure reminder
- [ ] Tạo `domain/notification.go`, `internal/notification/service.go`
- [ ] API: `GET /notifications` | `PUT /notifications/:id/read`
- [ ] Frontend: notification bell + dropdown
- [ ] Cronjob: departure reminder (3 ngày, 1 ngày trước)

---

## 🟢 PHASE 7: POLISH (Ưu tiên 7)

- [ ] Upload avatar API + frontend integration
- [ ] Tour gallery (multi-image) + `tour_images` table
- [ ] Tour schedule system + `tour_schedules` table
- [ ] Booking invoice PDF download
- [ ] Frontend validation improvements

---

## 🔵 PHASE 8: INFRASTRUCTURE & DEPLOYMENT

### 8.1 — Redis Cache
- [ ] Thêm Redis vào docker-compose
- [ ] Migrate rate limiter sang Redis
- [ ] Cache tour listings + user sessions

### 8.2 — File Storage (MinIO)
- [ ] Thêm MinIO vào docker-compose
- [ ] Tạo S3 client: upload tour images, avatars, review images

### 8.3 — Docker Production
- [ ] Dockerfile cho Go backend + React frontend
- [ ] `docker-compose.prod.yml` (all services)
- [ ] Nginx reverse proxy + SSL

### 8.4 — CI/CD & Testing
- [ ] GitHub Actions: lint + test + build
- [ ] Health check endpoint, graceful shutdown
- [ ] Unit tests cho Auth/Tour/Booking/Payment services
- [ ] Integration tests cho booking→payment flow

---

## 🏗️ SƠ ĐỒ INFRASTRUCTURE ĐỀ XUẤT

```
┌──────────────────────────────────────────────────┐
│                 NGINX (Reverse Proxy)            │
│                SSL Termination                   │
│             Port 80/443 → Internal               │
├────────────────────┬─────────────────────────────┤
│  React Frontend    │     Go Backend (Gin)        │
│  (nginx serve)     │     Port 8080               │
│  Port 3000         │                             │
│                    │  ┌───────────────────────┐   │
│                    │  │ Auth │ Tour │ Booking │   │
│                    │  │ Payment │ Review     │   │
│                    │  │ Coupon │ Notification│   │
│                    │  │ Admin Dashboard      │   │
│                    │  └───────────────────────┘   │
├────────────────────┴─────────────────────────────┤
│                  Data Layer                       │
│ ┌────────────┐ ┌──────────┐ ┌──────────────────┐│
│ │ PostgreSQL │ │  Redis   │ │  MinIO (S3)      ││
│ │ Port 5433  │ │ Port 6379│ │  Port 9000       ││
│ │ Main DB    │ │ Cache    │ │  File Storage    ││
│ └────────────┘ └──────────┘ └──────────────────┘│
└──────────────────────────────────────────────────┘
```

### Biến môi trường cần bổ sung
```env
# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# MinIO / S3
S3_ENDPOINT=localhost:9000
S3_ACCESS_KEY=minioadmin
S3_SECRET_KEY=minioadmin
S3_BUCKET=traveling

# Email (production)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=app-password

# App
APP_ENV=development
APP_PORT=8080
FRONTEND_URL=http://localhost:5173
```

---

## 📊 TỔNG KẾT

| Phase | Mô tả | Ước lượng |
|-------|--------|-----------|
| 0 | Nền tảng (RBAC, Price, Response) | 1 ngày |
| 1 | Payment Module | 2 ngày |
| 2 | Admin Tour + Booking | 1.5 ngày |
| 3 | Review Module | 1.5 ngày |
| 4 | Coupon Module | 1 ngày |
| 5 | Admin Dashboard | 1.5 ngày |
| 6 | Notification | 1 ngày |
| 7 | Polish | 1 ngày |
| 8 | Infrastructure | 2 ngày |
| **Tổng** | **~50 files mới, ~40 files sửa** | **~13 ngày** |

### Thứ tự thực hiện:
```
Phase 0 → Phase 1 → Phase 2 → Phase 3/4 (song song) → Phase 5 → Phase 6 → Phase 7 → Phase 8
```

> [!IMPORTANT]
> Phase 0 là **BLOCKING** — phải hoàn thành trước vì RBAC ảnh hưởng tất cả admin APIs, Tour Price ảnh hưởng Payment calculations.
