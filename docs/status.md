# 📊 TRẠNG THÁI DỰ ÁN TRAVELING

**Cập nhật lần cuối:** 2026-05-17
**Tiến độ tổng thể:** ~75-80% (Payment hardened through Option A-C, core modules completed)

---

## ✅ PHASE 0: NỀN TẢNG HẠ TẦNG — HOÀN THÀNH

### 0.1 — RBAC (Role-Based Access Control)
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `domain/user.go` | Thêm Role (customer/staff/admin), IsActive, helpers | ✅ |
| `shared/role_middleware.go` | StaffRequired(), AdminRequired() middleware | ✅ |
| `auth/middleware.go` | Set role vào context từ JWT claims | ✅ |
| `auth/token.go` | Thêm Role vào JWT claims (sign + parse) | ✅ |
| `cmd/server/main.go` | Seed admin + staff users | ✅ |
| `client/AuthContext.jsx` | Expose isAdmin, isStaff, userRole | ✅ |

### 0.2 — Tour Price Normalization
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `domain/tour.go` | PriceAmount int64, Slug, Rating, ReviewCount, IsActive | ✅ |
| `tour/service.go` | Filter/sort dùng PriceAmount, thêm sort rating/popular | ✅ |
| `tour/handler.go` | Dùng shared.RespondSuccess/Error | ✅ |
| `cmd/server/main.go` | Seed 14 tours với PriceAmount + Slug | ✅ |

### 0.3 — API Response Standardization
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `shared/response.go` | RespondSuccess, RespondError, RespondSuccessWithMeta | ✅ |
| `shared/pagination.go` | GetPaginationParams, BuildPaginationMeta | ✅ |

---

## ✅ PHASE 1: PAYMENT MODULE — HOÀN THÀNH

### Completion Scope
- [x] Option A: Sandbox-ready payment flow
  - Chuẩn hóa amount theo VND, chỉ nhân x100 tại VNPay gateway boundary
  - ReturnURL chỉ verify checksum và redirect UI; IPN là nguồn cập nhật trạng thái chính
  - Hỗ trợ VNPay IPN qua GET và POST query/form params
  - Dùng env cho API/frontend URLs thay vì hardcode localhost
  - Thêm focused tests cho amount, VNPay amount encoding, config, redirect URL
- [x] Option B: Production-ready VNPay foundation
  - IPN xử lý idempotent với transaction + row lock
  - Verify `vnp_Amount`, `vnp_ResponseCode`, `vnp_TransactionStatus` trước khi cập nhật paid/failed
  - Atomic update payment + booking + audit log
  - Production config yêu cầu HTTPS cho ReturnURL, IPNURL, FrontendURL
- [x] Option C: Multi-method foundation
  - Thêm `PaymentGateway` interface để tách orchestration khỏi VNPay implementation
  - VNPay client implement provider boundary (`ProviderName`, URL generation, signature validation, response parsing)

### Backend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `domain/payment.go` | Payment amount stored in VND, validation, summaries | ✅ |
| `payment/gateway.go` | Provider abstraction for VNPay/future gateways | ✅ |
| `payment/vnpay_client.go` | HMAC-SHA512, URL generation, x100 amount boundary, signature validation | ✅ |
| `payment/service.go` | InitiatePayment, display-only ReturnURL, authoritative IPN, status checks | ✅ |
| `payment/handler.go` | Return/IPN param collection, frontend redirect URL builder | ✅ |
| `payment/config.go` | VNPay + frontend env config, production HTTPS validation | ✅ |
| `payment/*_test.go`, `domain/payment_amount_test.go` | Focused amount/config/gateway tests | ✅ |

### API Routes
```
POST   /v1/api/payments/initiate         (Auth)
GET    /v1/api/payments/status/:ref      (Auth)
GET    /v1/api/bookings/:id/payments     (Auth)
GET    /v1/api/payments/return           (Public - VNPay callback)
GET    /v1/api/payments/webhook          (Public - VNPay IPN)
POST   /v1/api/payments/webhook          (Public - VNPay IPN)
```

### Frontend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `services/paymentService.js` | Payment API calls + helpers using `VITE_API_URL` | ✅ |
| `pages/payment/PaymentResult.jsx` | VNPay redirect result page | ✅ |
| `pages/payment/PaymentResult.css` | Styled result page | ✅ |
| `App.jsx` | Route /payment/result | ✅ |
| `pages/customer/BookingDetail.jsx` | Payment button integrated | ✅ |
| `pages/customer/Bookings.jsx` | Payment action routes to booking detail payment flow | ✅ |

### Required Env
```bash
# Backend
VNPAY_ENVIRONMENT=sandbox
VNPAY_MERCHANT_ID=your-vnpay-tmn-code
VNPAY_SECRET_KEY=your-vnpay-secret-key
VNPAY_RETURN_URL=http://localhost:8080/v1/api/payments/return
VNPAY_IPN_URL=http://localhost:8080/v1/api/payments/webhook
FRONTEND_BASE_URL=http://localhost:5173

# Frontend
VITE_API_URL=http://localhost:8080/v1/api
```

---

## ✅ PHASE 2: ADMIN DASHBOARD — HOÀN THÀNH

### Backend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `tour/admin_handler.go` | List, Create, Update, Delete(soft), Toggle | ✅ |
| `booking/admin_handler.go` | List, Detail, Confirm, Cancel, Stats | ✅ |
| `auth/admin_handler.go` | List users, Toggle status, Change role | ✅ |
| `booking/service.go` | Uses tour.PriceAmount instead of parsing | ✅ |

### API Routes
```
GET/POST/PUT/DELETE  /v1/api/admin/tours/*       (Staff+)
GET/PUT              /v1/api/admin/bookings/*     (Staff+)
GET    /v1/api/admin/bookings/stats              (Staff+)
GET/PUT              /v1/api/admin/users/*        (Staff+)
PUT    /v1/api/admin/users/:id/role              (Admin only - handled in logic)
```

### Frontend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `services/adminService.js` | Tour/Booking/User admin API calls | ✅ |
| `pages/admin/Tours.jsx` | Table + search + filter + CRUD modal | ✅ |
| `pages/admin/Bookings.jsx` | Table + stats cards + confirm/cancel | ✅ |
| `pages/admin/Users.jsx` | Table + role dropdown + lock/unlock | ✅ |
| `pages/admin/AdminPage.css` | Shared admin styles | ✅ |

---

## ✅ PHASE 3: REVIEW & RATING — HOÀN THÀNH

### Backend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `domain/review.go` | Review model + validation | ✅ |
| `review/repository.go` | CRUD + pagination queries | ✅ |
| `review/service.go` | Business logic + tour rating update | ✅ |
| `review/handler.go` | Public + Customer endpoints | ✅ |
| `review/admin_handler.go` | Admin: list, publish, hide, reply | ✅ |

### Frontend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `services/reviewService.js` | Review API calls | ✅ |
| `pages/customer/WriteReview.jsx` | Review form | ✅ |
| `components/review/ReviewList.jsx` | Review list for TourDetail | ✅ |
| `pages/admin/Reviews.jsx` | Admin review management | ✅ |

---

## ✅ PHASE 4: COUPON MODULE — HOÀN THÀNH

- [x] `domain/coupon.go` + backend CRUD
- [x] Customer: validate coupon API
- [x] Admin CRUD + usage tracking
- [x] Frontend: CouponInput + AdminCouponsPage

## ✅ PHASE 5: ADMIN DASHBOARD & REPORT — HOÀN THÀNH

### Backend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `dashboard/service.go` | Logic thống kê summary, revenue, top tours | ✅ |
| `dashboard/admin_handler.go` | Dashboard APIs | ✅ |

### Frontend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `services/dashboardService.js` | Gọi API thống kê | ✅ |
| `pages/admin/Dashboard.jsx` | Giao diện tổng quan & Biểu đồ (Recharts) | ✅ |
| `pages/admin/Reports.jsx` | Báo cáo chi tiết & CSV Export | ✅ |

## ✅ PHASE 6: NOTIFICATION — HOÀN THÀNH

### Backend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `domain/notification.go` | Notification model | ✅ |
| `notification/repository.go` | CRUD + queries | ✅ |
| `notification/service.go` | Send Notification logic | ✅ |
| `notification/handler.go` | Get & Mark as read API | ✅ |
| `booking/service.go` | Trigger Notification on Booking | ✅ |

### Frontend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `services/notificationService.js` | Gọi API thông báo | ✅ |
| `components/common/NotificationBell.jsx` | UI chuông & Dropdown | ✅ |
| `components/Layout/Header.jsx` | Gắn chuông vào header | ✅ |

## ✅ PHASE 7: POLISH — HOÀN THÀNH

- [x] Avatar upload
- [x] Tour gallery (multi-image)
- [x] Tour schedules
- [x] Booking invoice PDF

## ⬜ PHASE 8: INFRASTRUCTURE — CHƯA BẮT ĐẦU

- [ ] Redis cache
- [ ] MinIO file storage
- [ ] Docker production
- [ ] CI/CD + Testing

---

## 🔑 TÀI KHOẢN TEST

| Email | Password | Role |
|-------|----------|------|
| admin@traveling.com | 123456 | Admin |
| staff@traveling.com | 123456 | Staff |
| test@example.com | 123456 | Customer |
| user@example.com | 123456 | Customer |

## 🛠️ COMMANDS

```bash
# Start backend
cd server && go run cmd/server/main.go

# Start frontend
cd client && npm run dev

# Build check
cd server && go build ./...
```
