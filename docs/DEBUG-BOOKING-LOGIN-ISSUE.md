# 🐛 DEBUG: Booking Form Redirect to Login Issue

## Vấn đề
Sau khi điền form đặt tour và nhấn "Xác nhận đặt tour", hệ thống redirect về trang login.

## Nguyên nhân có thể

### 1. **Token không tồn tại trong localStorage**
- User đã logout nhưng UI vẫn hiển thị như đã login
- Token bị xóa bởi code khác
- Browser clear cache/cookies

### 2. **Token đã hết hạn**
- Access token hết hạn (thường 15-30 phút)
- Refresh token cũng hết hạn
- Backend trả về 401 Unauthorized

### 3. **Token format không đúng**
- Token không được lưu đúng key (`accessToken` vs `auth_tokens`)
- Token bị corrupt
- Token không phải JWT hợp lệ

### 4. **Backend middleware reject token**
- Token signature không hợp lệ
- Token claims thiếu `user_id`
- Token đã bị revoke

---

## Cách Debug

### Bước 1: Kiểm tra localStorage

Mở **Browser DevTools** (F12) → **Console**, chạy:

```javascript
// Kiểm tra tất cả tokens
console.log('accessToken:', localStorage.getItem('accessToken'));
console.log('refreshToken:', localStorage.getItem('refreshToken'));
console.log('auth_tokens:', localStorage.getItem('auth_tokens'));
console.log('user:', localStorage.getItem('user'));
```

**Kết quả mong đợi:**
- `accessToken`: Có giá trị (JWT string dài)
- `refreshToken`: Có giá trị (JWT string dài)
- `auth_tokens`: Có giá trị (JSON object)
- `user`: Có giá trị (JSON object với id, name, email)

**Nếu thiếu:**
- → Token sync issue
- → Cần login lại

### Bước 2: Sử dụng Debug Tool

Mở file: `client/debug-token.html` trong browser:

```bash
# Từ thư mục client
open debug-token.html
# Hoặc
firefox debug-token.html
# Hoặc
google-chrome debug-token.html
```

**Các chức năng:**
1. **Refresh**: Xem tokens hiện tại
2. **Clear All Tokens**: Xóa hết tokens (test logout)
3. **Test Booking API**: Test trực tiếp API booking

### Bước 3: Xem Console Logs

Sau khi thêm logging vào code, khi submit booking form, bạn sẽ thấy:

```
🔍 Debug - Before booking:
  accessToken: EXISTS
  auth_tokens: EXISTS

🔐 Axios Request Interceptor:
  URL: /bookings
  Method: post
  AccessToken: eyJhbGciOiJIUzI1NiIs...
  Authorization header: Bearer eyJhbGciOiJIUzI1...

📤 Sending booking request: {...}
```

**Nếu thấy:**
```
❌ Axios Response Error:
  Status: 401
  URL: /bookings
  Error: {message: "Invalid access token"}
```

→ Token không hợp lệ hoặc đã hết hạn

**Nếu thấy:**
```
🔄 Attempting token refresh...
❌ Token refresh failed: ...
❌ No refresh token found, forcing logout
```

→ Refresh token cũng hết hạn, cần login lại

### Bước 4: Kiểm tra Backend Logs

Xem terminal đang chạy backend:

```bash
cd server
go run cmd/server/main.go
```

**Tìm log:**
```
[GIN] 2026/04/30 - 15:24:23 | 401 | 2.5ms | ::1 | POST "/v1/api/bookings"
```

**Nếu thấy 401:**
- Token không được gửi lên
- Token không hợp lệ
- Token đã hết hạn

### Bước 5: Test với cURL

Test trực tiếp API bằng cURL:

```bash
# Lấy token từ localStorage (copy từ browser console)
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

# Test booking API
curl -X POST http://localhost:8080/v1/api/bookings \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "tour_id": 1,
    "full_name": "Test User",
    "phone": "0123456789",
    "email": "test@example.com",
    "adult_count": 1,
    "child_count": 0,
    "infant_count": 0,
    "travel_date": "2026-05-15",
    "note": "Test"
  }'
```

**Kết quả mong đợi:**
```json
{
  "success": true,
  "message": "Đặt tour thành công...",
  "booking": {...}
}
```

**Nếu nhận 401:**
```json
{
  "success": false,
  "message": "Invalid access token",
  "code": "AUTH_TOKEN_INVALID"
}
```

---

## Giải pháp

### Giải pháp 1: Login lại

**Đơn giản nhất:**
1. Logout
2. Login lại
3. Thử đặt tour

### Giải pháp 2: Fix Token Sync

Nếu vấn đề là token sync, thêm vào `AuthContext.jsx`:

```javascript
// Sau khi login thành công
const login = (userData, tokenData) => {
  setUser(userData);
  setTokens(tokenData);
  
  // Lưu tất cả formats
  localStorage.setItem('user', JSON.stringify(userData));
  localStorage.setItem('auth_tokens', JSON.stringify(tokenData));
  localStorage.setItem('accessToken', tokenData.access_token);
  localStorage.setItem('refreshToken', tokenData.refresh_token);
};
```

### Giải pháp 3: Tăng Token Expiry

Nếu token hết hạn quá nhanh, sửa backend:

```go
// server/internal/auth/jwt.go
func generateAccessToken(userID uint) (string, error) {
    claims := &AccessTokenClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Tăng từ 15 phút lên 24 giờ
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    // ...
}
```

### Giải pháp 4: Disable Auth (Testing Only)

**CHỈ ĐỂ TEST**, tạm thời bỏ auth requirement:

```go
// server/cmd/server/main.go
// Comment out auth middleware
// protected := v1.Group("/api")
// protected.Use(auth.AuthRequired())
// {
//     protected.POST("/bookings", booking.CreateBookingHandler)
// }

// Thay bằng:
v1.POST("/api/bookings", booking.CreateBookingHandler)
```

**⚠️ CẢNH BÁO:** Chỉ dùng để test, KHÔNG deploy lên production!

---

## Checklist Debug

- [ ] Kiểm tra localStorage có tokens không
- [ ] Xem console logs khi submit form
- [ ] Kiểm tra backend logs
- [ ] Test với cURL
- [ ] Thử login lại
- [ ] Kiểm tra token expiry time
- [ ] Verify token sync code

---

## Kết quả mong đợi

Sau khi fix, flow đúng sẽ là:

1. User điền form booking
2. Click "Xác nhận đặt tour"
3. Frontend gửi POST `/bookings` với `Authorization: Bearer <token>`
4. Backend verify token → OK
5. Backend tạo booking → Success
6. Frontend hiển thị "Đặt tour thành công"
7. Redirect về `/account/bookings`

---

## Liên hệ

Nếu vẫn gặp vấn đề, cung cấp:
1. Screenshot console logs
2. Backend terminal logs
3. localStorage content
4. Browser và version

---

**Last Updated:** 2026-04-30
