---
phase: implementation
title: Implementation Guide
feature: traveling-system
description: Hướng dẫn kỹ thuật triển khai hệ thống đặt tour du lịch
---

# Implementation Guide — Traveling System

## Development Setup

**Yêu cầu môi trường:**
- Go 1.21+
- Node.js 18+
- MySQL 8.0 (chạy local hoặc Docker)

**Chạy Backend:**
```bash
cd server
cp .env.example .env      # Điền DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET
go mod tidy
go run main.go            # Server chạy tại http://localhost:8080
```

**Chạy Frontend:**
```bash
cd client
npm install
npm run dev               # App chạy tại http://localhost:5173
```

**Biến môi trường `.env` (server):**
```env
DB_USER=root
DB_PASSWORD=yourpassword
DB_HOST=localhost
DB_PORT=3306
DB_NAME=travel_db
JWT_SECRET=your_super_secret_key_here
JWT_EXPIRE_HOURS=24
```

---

## Code Structure

```
server/
├── main.go                  # Entry point, khởi tạo router
├── go.mod / go.sum
├── .env
├── database/
│   └── database.go          # Kết nối MySQL, expose DB instance
├── models/
│   └── models.go            # Tất cả GORM struct models
├── middleware/
│   ├── auth.go              # JWT verification middleware
│   └── role.go              # Role-based access control
├── handlers/
│   ├── auth.go              # Login, Register handlers
│   ├── user.go              # GetProfile, UpdateProfile handlers
│   ├── tour.go              # CRUD tour handlers
│   ├── lichtour.go          # CRUD lịch tour handlers
│   ├── booking.go           # Tạo & xem phiếu đặt
│   ├── invoice.go           # Tạo & xem hóa đơn
│   ├── diadiem.go           # CRUD địa điểm
│   └── dichvu.go            # CRUD dịch vụ
└── routes/
    └── routes.go            # Đăng ký tất cả routes

client/src/
├── main.jsx
├── App.jsx                  # Router setup
├── api/
│   └── axiosInstance.js     # Axios với interceptor JWT
├── context/
│   ├── AuthContext.jsx      # Trạng thái đăng nhập toàn cục
│   └── TourContext.jsx      # Cache danh sách tour
├── components/
│   ├── Auth/                # Login, Register
│   ├── Layout/              # Header, Footer
│   ├── Profile/             # Xem & sửa hồ sơ
│   ├── Tour/                # TourList, TourCard, TourDetail
│   ├── Booking/             # BookingForm, BookingHistory
│   └── Invoice/             # InvoiceList, InvoiceDetail
└── pages/
    ├── HomePage.jsx
    ├── TourPage.jsx
    └── DashboardPage.jsx
```

---

## Implementation Notes

### Core Features

#### JWT Authentication (middleware/auth.go)
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")
        if tokenStr == "" || !strings.HasPrefix(tokenStr, "Bearer ") {
            c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "Chưa đăng nhập"})
            return
        }
        tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
        // Parse & validate JWT, set userID và role vào context
        claims, err := parseJWT(tokenStr)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "Token không hợp lệ"})
            return
        }
        c.Set("userID", claims.UserID)
        c.Set("role", claims.Role)
        c.Next()
    }
}
```

#### bcrypt Password (handlers/auth.go)
```go
// Đăng ký — hash password
hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
newUser.Password = string(hashedPwd)

// Đăng nhập — verify
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
if err != nil {
    // password sai
}
```

#### Đặt tour — kiểm tra số chỗ còn lại (handlers/booking.go)
```go
// Đếm số khách đã đặt trong lịch tour này
var totalBooked int64
db.Model(&PDTour{}).
    Where("lich_tour_id = ?", req.LichTourID).
    Select("COALESCE(SUM(so_khach_nl + so_khach_tre_em), 0)").
    Scan(&totalBooked)

newGuests := req.SoKhachNL + req.SoKhachTreEm
if int(totalBooked)+newGuests > lichTour.Tour.SLKhachMax {
    c.JSON(400, gin.H{"success": false, "message": "Tour đã hết chỗ"})
    return
}
```

#### Axio Instance với JWT (client/src/api/axiosInstance.js)
```js
import axios from 'axios'

const api = axios.create({ baseURL: 'http://localhost:8080/api' })

api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

api.interceptors.response.use(
  res => res,
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(err)
  }
)

export default api
```

### Patterns & Best Practices

- **Handler pattern**: Mỗi handler chỉ làm 1 việc — bind request, validate, gọi DB, trả response
- **Error response nhất quán**: Luôn trả `{ success: bool, message: string, data?: ... }`
- **GORM Preload**: Dùng `db.Preload("LichTour.Tour").Find(&pdtours)` thay vì nhiều query
- **React Context**: Chỉ dùng cho state toàn cục (auth, tour list) — component state cho form data
- **CSS Modules**: Mỗi component có file `.css` riêng, không dùng global class conflict

---

## Integration Points

**Frontend → Backend:**
- Base URL: `http://localhost:8080/api`
- Format: JSON
- Auth: `Authorization: Bearer <token>` header

**Backend → Database:**
- Connection string: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True`
- GORM tự động handle connection pool

---

## Error Handling

**Backend — HTTP Status codes chuẩn:**

| Tình huống | Status Code |
|-----------|-------------|
| Thành công | 200 / 201 |
| Dữ liệu không hợp lệ | 400 Bad Request |
| Chưa đăng nhập | 401 Unauthorized |
| Không có quyền | 403 Forbidden |
| Không tìm thấy | 404 Not Found |
| Trùng lặp (email/mã) | 409 Conflict |
| Lỗi server | 500 Internal Server Error |

**Frontend — Error handling pattern:**
```jsx
const [error, setError] = useState(null)
try {
  const res = await api.post('/pdtour', data)
  if (res.data.success) { /* xử lý thành công */ }
} catch (err) {
  setError(err.response?.data?.message || 'Có lỗi xảy ra')
}
```

---

## Performance Considerations

- **DB Index**: Thêm index trên `email` (tblThanhVien), `tour_id` (tblLichTour), `lich_tour_id` (tblPDTour)
- **GORM Select**: Chỉ select các cột cần thiết thay vì `SELECT *` khi trả danh sách
- **Frontend**: Dùng `React.memo` cho `TourCard` tránh re-render không cần thiết
- **Pagination**: Thêm `?page=1&limit=10` cho API `/api/tours` và `/api/pdtour`

---

## Security Notes

- Password **không bao giờ** được trả về trong response (dùng `json:"-"` tag trong GORM struct)
- JWT secret phải đủ dài (>= 32 ký tự) và lưu trong `.env`, không commit vào git
- Input validation trước khi lưu DB: kiểm tra độ dài, ký tự đặc biệt
- CORS chỉ whitelist `http://localhost:5173` (dev) và domain production thực tế
- File `.env` phải có trong `.gitignore`
