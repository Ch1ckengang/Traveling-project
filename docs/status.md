# 📊 TRẠNG THÁI DỰ ÁN TRAVELING

**Cập nhật lần cuối:** 2026-05-13  
**Tiến độ tổng thể:** ~60-65% (Phase 0-2 hoàn thành, Phase 3 đang thực hiện)

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

### Backend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `payment/vnpay_client.go` | HMAC-SHA512, URL generation, signature validation | ✅ |
| `payment/service.go` | InitiatePayment, ProcessReturn, ProcessWebhook, GetPaymentStatus | ✅ |
| `payment/handler.go` | 5 HTTP endpoints | ✅ |
| `payment/config.go` | LoadVNPayConfig convenience wrapper | ✅ |

### API Routes
```
POST   /v1/api/payments/initiate         (Auth)
GET    /v1/api/payments/status/:ref      (Auth)
GET    /v1/api/bookings/:id/payments     (Auth)
GET    /v1/api/payments/return           (Public - VNPay callback)
POST   /v1/api/payments/webhook          (Public - VNPay IPN)
```

### Frontend
| File | Mô tả | Trạng thái |
|------|--------|------------|
| `services/paymentService.js` | All payment API calls + helpers | ✅ |
| `pages/payment/PaymentResult.jsx` | VNPay redirect result page | ✅ |
| `pages/payment/PaymentResult.css` | Styled result page | ✅ |
| `App.jsx` | Route /payment/result | ✅ |
| `pages/customer/BookingDetail.jsx` | Payment button integrated | ✅ |

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

## 🔄 PHASE 4: COUPON MODULE — ĐANG THỰC HIỆN

- [ ] `domain/coupon.go` + backend CRUD
- [ ] Customer: validate coupon API
- [ ] Admin CRUD + usage tracking
- [ ] Frontend: CouponInput + AdminCouponsPage

## ⬜ PHASE 5: ADMIN DASHBOARD & REPORT — CHƯA BẮT ĐẦU

- [ ] Dashboard summary/revenue/top-tours APIs
- [ ] AdminDashboardPage + charts (Recharts)
- [ ] CSV export

## ⬜ PHASE 6: NOTIFICATION — CHƯA BẮT ĐẦU

- [ ] Email templates
- [ ] Notification model + API
- [ ] Frontend notification bell
- [ ] Cronjob departure reminders

## ⬜ PHASE 7: POLISH — CHƯA BẮT ĐẦU

- [ ] Avatar upload
- [ ] Tour gallery (multi-image)
- [ ] Tour schedules
- [ ] Booking invoice PDF

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
