---
phase: planning
title: Project Planning & Task Breakdown
feature: traveling-system
description: Kế hoạch triển khai hệ thống đặt tour du lịch trực tuyến
---

# Project Planning & Task Breakdown — Traveling System

## Milestones

- [x] **M0 — Foundation**: Backend Gin + GORM + MySQL, Frontend React + Vite, Auth cơ bản (Login/Register)
- [ ] **M1 — Security & Auth nâng cao**: JWT, bcrypt, middleware phân quyền
- [ ] **M2 — Tour Management**: CRUD tour, lịch tour, địa điểm, dịch vụ
- [ ] **M3 — Booking Flow**: Đặt tour, xem phiếu đặt
- [ ] **M4 — Invoice & Dashboard**: Hóa đơn, dashboard nhân viên
- [ ] **M5 — Polish & Testing**: Tối ưu UI/UX, kiểm thử, sửa lỗi

---

## Task Breakdown

### Phase 1 — Security & Auth nâng cao (M1)

#### Backend
- [ ] **1.1** Thêm dependency `golang-jwt/jwt` và `golang.org/x/crypto/bcrypt` vào `go.mod`
- [ ] **1.2** Cập nhật `models.go`: thêm field `Role` (`customer`/`staff`/`admin`) vào `ThanhVien`
- [ ] **1.3** Migrate `Password` sang bcrypt hash (viết script migration một lần)
- [ ] **1.4** Tạo `middleware/auth.go`: xác thực JWT token từ header `Authorization: Bearer <token>`
- [ ] **1.5** Tạo `middleware/role.go`: kiểm tra role (chỉ nhân viên/admin mới truy cập được một số route)
- [ ] **1.6** Cập nhật `POST /api/login` — trả về JWT token thay vì chỉ trả user
- [ ] **1.7** Refactor `main.go` — tách handlers ra thư mục `handlers/`

#### Frontend
- [ ] **1.8** Cập nhật `AuthContext.jsx` — lưu JWT token vào `localStorage`, thêm hàm `logout`
- [ ] **1.9** Tạo `axiosInstance.js` — tự động đính kèm `Authorization` header cho mọi request
- [ ] **1.10** Thêm Protected Route — redirect về `/login` nếu chưa đăng nhập

---

### Phase 2 — Tour Management (M2)

#### Backend — Database Models
- [ ] **2.1** Cập nhật `models.go` — thêm đầy đủ structs: `KhachHang`, `NhanVien`, `Tour` (đầy đủ fields), `LichTour`, `DiaDiem`, `TourDiaDiem`, `DichVu`, `DichVuDiaDiem`
- [ ] **2.2** Cập nhật `database.go` — AutoMigrate tất cả models mới
- [ ] **2.3** Cập nhật seed data — thêm dữ liệu mẫu cho tour đầy đủ fields

#### Backend — Tour APIs
- [ ] **2.4** Tạo `handlers/tour.go`:
  - `GET /api/tours` — danh sách tour (có filter theo location, giá)
  - `GET /api/tours/:id` — chi tiết tour + lịch tour
  - `POST /api/tours` — tạo tour mới (role: staff/admin)
  - `PUT /api/tours/:id` — cập nhật tour (role: staff/admin)
  - `DELETE /api/tours/:id` — xóa tour (role: admin)
- [ ] **2.5** Tạo `handlers/lichtour.go`:
  - `GET /api/tours/:id/lich` — lịch khởi hành của tour
  - `POST /api/lichtour` — tạo lịch mới (role: staff)
  - `DELETE /api/lichtour/:id` — xóa lịch (role: staff)
- [ ] **2.6** Tạo `handlers/diadiem.go`:
  - `GET /api/diadiem` — danh sách địa điểm
  - `POST /api/diadiem` — thêm địa điểm
  - `POST /api/tour-diadiem` — gán địa điểm vào tour
- [ ] **2.7** Tạo `handlers/dichvu.go`:
  - `GET /api/dichvu` — danh sách dịch vụ
  - `POST /api/dichvu` — thêm dịch vụ mới

#### Frontend — Tour UI
- [ ] **2.8** Tạo `components/Tour/TourCard.jsx` — card hiển thị tên, địa điểm, giá, thời gian
- [ ] **2.9** Tạo `components/Tour/TourList.jsx` — grid danh sách tour + skeleton loading
- [ ] **2.10** Tạo `components/Tour/TourDetail.jsx` — trang chi tiết tour: mô tả, địa điểm, lịch, dịch vụ
- [ ] **2.11** Cập nhật `SearchBar.jsx` — gọi API filter tour theo địa điểm và khoảng giá
- [ ] **2.12** Tạo `pages/TourManagePage.jsx` — trang CRUD tour cho nhân viên (bảng + form)

---

### Phase 3 — Booking Flow (M3)

#### Backend
- [ ] **3.1** Tạo `handlers/booking.go`:
  - `POST /api/pdtour` — tạo phiếu đặt (role: customer)
    - Kiểm tra số lượng chỗ còn lại của lịch tour
    - Tự động tạo `maKH` nếu user chưa có record `KhachHang`
  - `GET /api/pdtour/my` — phiếu đặt của khách hàng hiện tại
  - `GET /api/pdtour` — tất cả phiếu đặt (role: staff)
  - `PUT /api/pdtour/:id/cancel` — hủy phiếu đặt

#### Frontend
- [ ] **3.2** Tạo `components/Booking/BookingForm.jsx` — form chọn lịch, nhập số người, xác nhận đặt
- [ ] **3.3** Tạo `components/Booking/BookingHistory.jsx` — danh sách phiếu đặt của khách
- [ ] **3.4** Thêm nút "Đặt Tour" vào `TourDetail.jsx`
- [ ] **3.5** Thêm tab "Lịch sử đặt tour" vào `Profile.jsx`

---

### Phase 4 — Invoice & Dashboard (M4)

#### Backend
- [ ] **4.1** Tạo `handlers/invoice.go`:
  - `GET /api/hoadon` — danh sách hóa đơn (role: staff)
  - `POST /api/hoadon` — tạo hóa đơn từ phiếu đặt (role: staff)
    - Tính tổng tiền = (soKhachNL + soKhachTreEm * 0.7) * chiPhi
  - `GET /api/hoadon/:id` — chi tiết hóa đơn

#### Frontend
- [ ] **4.2** Tạo `pages/DashboardPage.jsx` — dashboard nhân viên: thống kê tổng quan
- [ ] **4.3** Tạo `components/Invoice/InvoiceList.jsx` — danh sách hóa đơn
- [ ] **4.4** Tạo `components/Invoice/InvoiceDetail.jsx` — chi tiết hóa đơn, nút in PDF

---

### Phase 5 — Polish & Testing (M5)

- [ ] **5.1** Thêm loading states và error boundaries toàn bộ app
- [ ] **5.2** Responsive CSS cho mobile (màn hình < 768px)
- [ ] **5.3** Unit test backend handlers với Go testing package
- [ ] **5.4** Kiểm thử end-to-end: luồng đăng ký → đặt tour → tạo hóa đơn
- [ ] **5.5** Viết hướng dẫn cài đặt và chạy dự án (`README.md`)

---

## Dependencies

```
M1 (Auth) ──► M2 (Tour) ──► M3 (Booking) ──► M4 (Invoice)
                                                    │
                                               M5 (Testing)
```

- **M1 phải hoàn thành trước M2**: Các API tour cần JWT middleware
- **M2 phải hoàn thành trước M3**: Cần có Tour và LichTour để đặt tour
- **M3 phải hoàn thành trước M4**: Cần có PDTour để tạo HoaDon
- **M5 có thể song song với M3, M4**

---

## Timeline & Estimates

| Milestone | Ước tính | Ghi chú |
|-----------|---------|---------|
| M1 — Security & Auth | 2–3 ngày | Refactor quan trọng, cẩn thận |
| M2 — Tour Management | 3–4 ngày | Nhiều CRUD, cần seed đủ data |
| M3 — Booking Flow | 2–3 ngày | Logic nghiệp vụ phức tạp nhất |
| M4 — Invoice & Dashboard | 2 ngày | Tương đối đơn giản |
| M5 — Polish & Testing | 2 ngày | Chạy song song với M3/M4 |
| **Tổng** | **~12–14 ngày** | |

---

## Risks & Mitigation

| Rủi ro | Mức độ | Biện pháp |
|--------|--------|-----------|
| Migration password từ plaintext sang bcrypt làm mất data test | 🟡 Trung bình | Chạy script migration trong transaction, backup DB trước |
| Logic tính giá hóa đơn (người lớn/trẻ em) chưa được xác nhận | 🟡 Trung bình | Xác nhận với Q-04 trước khi implement M4 |
| Frontend không đồng bộ với API mới sau refactor | 🟡 Trung bình | Thống nhất API contract trước, dùng Postman test |
| Overfit database schema với GORM AutoMigrate | 🟢 Thấp | Dùng migration scripts thủ công khi deploy production |

---

## Resources Needed

- **Backend Developer**: Go, GORM, JWT, bcrypt
- **Frontend Developer**: React, Axios, CSS modules
- **Database**: MySQL 8.0 đã cài đặt, user có quyền CREATE/DROP
- **Tools**: Go 1.21+, Node.js 18+, Postman (test API), MySQL Workbench (xem data)
