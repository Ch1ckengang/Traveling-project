---
phase: design
title: System Design & Architecture
feature: traveling-system
description: Thiết kế kiến trúc hệ thống đặt tour du lịch trực tuyến
---

# System Design & Architecture — Traveling System

## Architecture Overview

```mermaid
graph TD
    Browser["🌐 Browser\n(React + Vite)"]
    API["⚙️ API Server\n(Go + Gin)"]
    DB[("🗄️ MySQL Database")]
    Auth["🔐 Auth Middleware\n(JWT)"]

    Browser -->|"HTTP/REST JSON"| Auth
    Auth --> API
    API --> DB

    subgraph Frontend
        Login["Login / Register"]
        Profile["Profile"]
        TourList["Tour List / Search"]
        TourDetail["Tour Detail"]
        BookingForm["Booking Form"]
        InvoiceView["Invoice View"]
    end

    subgraph Backend Routes
        AuthRoute["/api/login\n/api/register"]
        UserRoute["/api/users/:id"]
        TourRoute["/api/tours"]
        LichTourRoute["/api/lichtour"]
        PDTourRoute["/api/pdtour"]
        HoaDonRoute["/api/hoadon"]
        DiaDiemRoute["/api/diadiem"]
        DichVuRoute["/api/dichvu"]
    end

    Browser --> Frontend
    Frontend --> Backend Routes
```

**Công nghệ sử dụng:**

| Tầng | Công nghệ | Lý do |
|------|-----------|-------|
| Frontend | React 18 + Vite | SPA nhanh, hot reload tốt |
| State Management | React Context API | Đủ dùng cho quy mô hiện tại |
| Backend | Go 1.21 + Gin | Hiệu năng cao, cú pháp rõ ràng |
| ORM | GORM v2 | Auto-migrate, query đơn giản |
| Database | MySQL 8.0 | Quan hệ chặt chẽ, phù hợp mô hình nghiệp vụ |
| Authentication | JWT (golang-jwt) | Stateless, dễ mở rộng |
| Password | bcrypt | Bảo mật chuẩn công nghiệp |

---

## Data Models

### Sơ đồ quan hệ đầy đủ

```
tblThanhVien ────┬──── tblKhachHang ──── tblPDTour ────┬──── tblLichTour ──── tblTour ────┬──── tblTourDiaDiem ──── tblDiaDiem
                 │                           │          │                                  │
                 ├──── tblNhanVien           └──── tblHoaDon                              └──── tblTourDiaDiem ──── tblDichvuDiaDiem ──── tblDichvu
                 │
                 └──── tblHoaDon (nhân viên lập)
```

### Chi tiết Model Go (GORM)

```go
// tblThanhVien
type ThanhVien struct {
    ID       uint   `gorm:"primaryKey"`
    Username string `gorm:"size:25;not null"`
    Password string `gorm:"size:255;not null"`   // bcrypt hash
    Ngaysinh string `gorm:"type:date;not null"`
    Email    string `gorm:"size:25;not null;unique"`
}

// tblKhachHang
type KhachHang struct {
    ID            uint      `gorm:"primaryKey"`
    MaKH          string    `gorm:"size:25;not null"`
    ThanhVienID   uint      `gorm:"not null"`
    ThanhVien     ThanhVien `gorm:"foreignKey:ThanhVienID"`
}

// tblNhanVien
type NhanVien struct {
    ID           uint      `gorm:"primaryKey"`
    MaNV         string    `gorm:"size:25;not null"`
    Chucvu       string    `gorm:"size:25;not null"`
    ThanhVienID  uint      `gorm:"not null"`
    ThanhVien    ThanhVien `gorm:"foreignKey:ThanhVienID"`
}

// tblTour
type Tour struct {
    ID         uint   `gorm:"primaryKey"`
    MaTour     string `gorm:"size:25;not null"`
    TenTour    string `gorm:"size:25;not null"`
    ThoiGian   string `gorm:"size:25;not null"`
    PhuongTien string `gorm:"size:25;not null"`
    SLKhachMax int    `gorm:"not null"`
    Mota       string `gorm:"size:255;not null"`
    ChiPhi     int    `gorm:"not null"`
}

// tblLichTour
type LichTour struct {
    ID      uint   `gorm:"primaryKey"`
    NgayVe  string `gorm:"type:date;not null"`
    TourID  uint   `gorm:"not null"`
    Tour    Tour   `gorm:"foreignKey:TourID"`
}

// tblPDTour
type PDTour struct {
    ID           uint      `gorm:"primaryKey"`
    SoKhachNL    int       `gorm:"not null"`
    SoKhachTreEm int       `gorm:"not null"`
    KhachHangID  uint      `gorm:"not null"`
    LichTourID   uint      `gorm:"not null"`
    KhachHang    KhachHang `gorm:"foreignKey:KhachHangID"`
    LichTour     LichTour  `gorm:"foreignKey:LichTourID"`
}

// tblHoaDon
type HoaDon struct {
    ID          uint      `gorm:"primaryKey"`
    MaHD        string    `gorm:"size:25;not null"`
    PDTourID    uint      `gorm:"not null"`
    ThanhVienID uint      `gorm:"not null"`
    PDTour      PDTour    `gorm:"foreignKey:PDTourID"`
    ThanhVien   ThanhVien `gorm:"foreignKey:ThanhVienID"`
}
```

---

## API Design

### Authentication

| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| POST | `/api/login` | Đăng nhập → trả JWT token | ❌ |
| POST | `/api/register` | Đăng ký tài khoản mới | ❌ |

### Quản lý người dùng

| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| GET | `/api/users/:id` | Lấy thông tin cá nhân | ✅ JWT |
| PUT | `/api/users/:id` | Cập nhật thông tin cá nhân | ✅ JWT |

### Tour

| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| GET | `/api/tours` | Danh sách tất cả tour | ❌ |
| GET | `/api/tours/:id` | Chi tiết một tour | ❌ |
| POST | `/api/tours` | Tạo tour mới | ✅ NV |
| PUT | `/api/tours/:id` | Cập nhật tour | ✅ NV |
| DELETE | `/api/tours/:id` | Xóa tour | ✅ NV |

### Lịch Tour

| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| GET | `/api/tours/:id/lich` | Lấy danh sách lịch của tour | ❌ |
| POST | `/api/lichtour` | Tạo lịch tour mới | ✅ NV |
| DELETE | `/api/lichtour/:id` | Xóa lịch tour | ✅ NV |

### Phiếu Đặt Tour

| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| GET | `/api/pdtour` | Danh sách phiếu đặt (nhân viên) | ✅ NV |
| GET | `/api/pdtour/my` | Phiếu đặt của khách hàng hiện tại | ✅ KH |
| POST | `/api/pdtour` | Tạo phiếu đặt mới | ✅ KH |

### Hóa Đơn

| Method | Endpoint | Mô tả | Auth |
|--------|----------|-------|------|
| GET | `/api/hoadon` | Danh sách hóa đơn | ✅ NV |
| POST | `/api/hoadon` | Tạo hóa đơn từ phiếu đặt | ✅ NV |

### Request / Response mẫu

**POST /api/login**
```json
// Request
{ "email": "user@example.com", "password": "123456" }

// Response 200
{ "success": true, "token": "eyJ...", "user": { "id": 1, "name": "Nguyễn Văn A" } }

// Response 401
{ "success": false, "message": "Email hoặc mật khẩu không đúng" }
```

**POST /api/pdtour**
```json
// Request
{
  "lich_tour_id": 3,
  "so_khach_nl": 2,
  "so_khach_tre_em": 1
}

// Response 201
{
  "success": true,
  "phieu_dat": { "id": 10, "ma_pd": "PD202503001", ... }
}
```

---

## Component Breakdown

### Frontend Components

```
src/
├── components/
│   ├── Auth/
│   │   ├── Login.jsx         ✅ Có sẵn
│   │   └── Register.jsx      ✅ Có sẵn
│   ├── Home/
│   │   └── SearchBar.jsx     ✅ Có sẵn
│   ├── Layout/
│   │   └── Header.jsx        ✅ Có sẵn
│   ├── Profile/
│   │   └── Profile.jsx       ✅ Có sẵn
│   ├── Tour/
│   │   ├── TourList.jsx      🆕 Cần tạo
│   │   ├── TourCard.jsx      🆕 Cần tạo
│   │   └── TourDetail.jsx    🆕 Cần tạo
│   ├── Booking/
│   │   └── BookingForm.jsx   🆕 Cần tạo
│   └── Invoice/
│       └── InvoiceList.jsx   🆕 Cần tạo
├── context/
│   ├── AuthContext.jsx        ✅ Có sẵn
│   └── TourContext.jsx        🆕 Cần tạo
└── pages/
    ├── HomePage.jsx           🆕 Cần tạo
    ├── TourPage.jsx           🆕 Cần tạo
    └── DashboardPage.jsx      🆕 Cần tạo (nhân viên)
```

### Backend Modules

```
server/
├── main.go              ✅ Có sẵn (cần refactor)
├── database/
│   └── database.go      ✅ Có sẵn
├── models/
│   └── models.go        ✅ Có sẵn (cần bổ sung models)
├── middleware/
│   └── auth.go          🆕 Cần tạo (JWT middleware)
├── handlers/
│   ├── auth.go          🆕 Cần tạo
│   ├── tour.go          🆕 Cần tạo
│   ├── booking.go       🆕 Cần tạo
│   └── invoice.go       🆕 Cần tạo
└── routes/
    └── routes.go        🆕 Cần tạo
```

---

## Design Decisions

| Quyết định | Lựa chọn | Lý do |
|-----------|---------|-------|
| Auth strategy | JWT (stateless) | Dễ scale, không cần session store |
| Password hashing | bcrypt (cost=10) | Chuẩn bảo mật, chống brute-force |
| API style | RESTful JSON | Đơn giản, dễ test |
| DB migrations | GORM AutoMigrate | Nhanh cho giai đoạn dev |
| Error format | `{ success, message, data }` | Nhất quán, dễ xử lý ở client |
| Role phân quyền | Field `role` trong ThanhVien (`customer`/`staff`/`admin`) | Đơn giản, đủ dùng |

---

## Non-Functional Requirements

**Hiệu năng:**
- API response < 500ms cho các query thông thường
- Hỗ trợ tối thiểu 50 concurrent users

**Bảo mật:**
- Password lưu dưới dạng bcrypt hash, không bao giờ trả về plaintext
- JWT token hết hạn sau 24 giờ
- Input validation và sanitize ở cả client và server
- CORS chỉ cho phép origin của frontend

**Độ tin cậy:**
- Database transaction cho các thao tác tạo hóa đơn
- Xử lý lỗi rõ ràng với HTTP status code chuẩn
