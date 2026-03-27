---
phase: implementation
title: Implementation Guide
feature: traveling-system
description: Technical implementation guide for the online tour booking system
---

# Implementation Guide — Traveling System

## Development Setup

**Prerequisites:**
- Go 1.21+
- Node.js 18+
- MySQL 8.0 (local or Docker)

**Run Backend:**
```bash
cd server
cp .env.example .env      # Fill in DB_USER, DB_PASSWORD, DB_NAME, JWT_SECRET
go mod tidy
go run main.go            # Server runs at http://localhost:8080
```

**Run Frontend:**
```bash
cd client
npm install
npm run dev               # App runs at http://localhost:5173
```

**Environment variables `.env` (server):**
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
├── main.go                  # Entry point, router initialization
├── go.mod / go.sum
├── .env
├── database/
│   └── database.go          # MySQL connection, exposes DB instance
├── models/
│   └── models.go            # All GORM struct definitions
├── middleware/
│   ├── auth.go              # JWT verification middleware
│   └── role.go              # Role-based access control
├── handlers/
│   ├── auth.go              # Login, Register handlers
│   ├── user.go              # GetProfile, UpdateProfile handlers
│   ├── tour.go              # CRUD tour handlers
│   ├── lichtour.go          # CRUD tour schedule handlers
│   ├── booking.go           # Create & view bookings
│   ├── invoice.go           # Create & view invoices
│   ├── diadiem.go           # CRUD destinations
│   └── dichvu.go            # CRUD services
└── routes/
    └── routes.go            # Register all routes

client/src/
├── main.jsx
├── App.jsx                  # Router setup
├── api/
│   └── axiosInstance.js     # Axios with JWT interceptor
├── context/
│   ├── AuthContext.jsx      # Global auth state
│   └── TourContext.jsx      # Tour list cache
├── components/
│   ├── Auth/                # Login, Register
│   ├── Layout/              # Header, Footer
│   ├── Profile/             # View & edit profile
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
            c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "Not authenticated"})
            return
        }
        tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
        // Parse & validate JWT, set userID and role into context
        claims, err := parseJWT(tokenStr)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"success": false, "message": "Invalid token"})
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
// Register — hash password
hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
newUser.Password = string(hashedPwd)

// Login — verify
err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
if err != nil {
    // wrong password
}
```

#### Booking — check remaining seats (handlers/booking.go)
```go
// Count total guests already booked for this schedule
var totalBooked int64
db.Model(&PDTour{}).
    Where("lich_tour_id = ?", req.LichTourID).
    Select("COALESCE(SUM(so_khach_nl + so_khach_tre_em), 0)").
    Scan(&totalBooked)

newGuests := req.SoKhachNL + req.SoKhachTreEm
if int(totalBooked)+newGuests > lichTour.Tour.SLKhachMax {
    c.JSON(400, gin.H{"success": false, "message": "Tour is fully booked"})
    return
}
```

#### Axios Instance with JWT (client/src/api/axiosInstance.js)
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

- **Handler pattern**: Each handler does one thing — bind request, validate, call DB, return response
- **Consistent error response**: Always return `{ success: bool, message: string, data?: ... }`
- **GORM Preload**: Use `db.Preload("LichTour.Tour").Find(&pdtours)` instead of multiple queries
- **React Context**: Only for global state (auth, tour list) — use component state for form data
- **CSS Modules**: Each component has its own `.css` file to avoid global class conflicts

---

## Integration Points

**Frontend → Backend:**
- Base URL: `http://localhost:8080/api`
- Format: JSON
- Auth: `Authorization: Bearer <token>` header

**Backend → Database:**
- Connection string: `user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True`
- GORM handles connection pooling automatically

---

## Error Handling

**Backend — Standard HTTP Status Codes:**

| Situation | Status Code |
|-----------|-------------|
| Success | 200 / 201 |
| Invalid input | 400 Bad Request |
| Not authenticated | 401 Unauthorized |
| Insufficient permissions | 403 Forbidden |
| Resource not found | 404 Not Found |
| Duplicate (email/code) | 409 Conflict |
| Server error | 500 Internal Server Error |

**Frontend — Error handling pattern:**
```jsx
const [error, setError] = useState(null)
try {
  const res = await api.post('/pdtour', data)
  if (res.data.success) { /* handle success */ }
} catch (err) {
  setError(err.response?.data?.message || 'An error occurred')
}
```

---

## Performance Considerations

- **DB Index**: Add indexes on `email` (tblThanhVien), `tour_id` (tblLichTour), `lich_tour_id` (tblPDTour)
- **GORM Select**: Only select required columns instead of `SELECT *` when returning lists
- **Frontend**: Use `React.memo` on `TourCard` to prevent unnecessary re-renders
- **Pagination**: Add `?page=1&limit=10` support to `/api/tours` and `/api/pdtour`

---

## Security Notes

- Passwords are **never** returned in responses (use `json:"-"` tag on GORM struct fields)
- JWT secret must be long (>= 32 characters) and stored in `.env`, never committed to git
- Validate and sanitize all inputs before saving to DB: check length, special characters
- CORS should only whitelist `http://localhost:5173` (dev) and the actual production domain
- `.env` file must be listed in `.gitignore`
