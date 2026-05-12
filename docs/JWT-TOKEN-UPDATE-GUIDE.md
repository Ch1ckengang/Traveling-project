# 🔐 JWT Token Update Guide

## Vấn đề hiện tại

Sau khi cập nhật JWT secrets trong file `server/.env`, tất cả các token cũ đã trở nên không hợp lệ. Điều này gây ra lỗi 401 Unauthorized khi truy cập các API yêu cầu authentication.

## Nguyên nhân

JWT tokens được ký (signed) bằng secret key. Khi secret key thay đổi:
- Token cũ (signed với secret cũ) → Backend không thể verify → 401 Error
- Token mới (signed với secret mới) → Backend verify thành công → OK

## Giải pháp

### Bước 1: Xóa token cũ và đăng nhập lại

**Cách 1: Logout từ UI (Khuyến nghị)**
1. Mở trình duyệt
2. Click vào tên user ở góc phải navbar
3. Click "Đăng xuất"
4. Đăng nhập lại với tài khoản của bạn

**Cách 2: Xóa localStorage thủ công**
1. Mở trình duyệt
2. Nhấn F12 để mở Developer Tools
3. Vào tab "Console"
4. Chạy lệnh:
```javascript
localStorage.clear();
location.reload();
```
5. Đăng nhập lại

**Cách 3: Xóa cookies và cache**
1. Mở trình duyệt
2. Nhấn Ctrl+Shift+Delete (hoặc Cmd+Shift+Delete trên Mac)
3. Chọn "Cookies and other site data"
4. Chọn "Cached images and files"
5. Click "Clear data"
6. Reload trang và đăng nhập lại

### Bước 2: Verify token mới hoạt động

Sau khi đăng nhập lại:
1. Kiểm tra localStorage có token mới:
   - Mở Developer Tools (F12)
   - Vào tab "Application" → "Local Storage"
   - Kiểm tra có các key: `accessToken`, `refreshToken`, `auth_tokens`, `user`

2. Test các chức năng:
   - ✅ Xem danh sách tour
   - ✅ Xem chi tiết tour
   - ✅ Đặt tour (booking)
   - ✅ Xem danh sách booking
   - ✅ Xem chi tiết booking

## Cấu hình JWT hiện tại

```env
# Access Token: 24 giờ (1440 phút)
JWT_ACCESS_TTL_MINUTES=1440

# Refresh Token: 30 ngày (720 giờ)
JWT_REFRESH_TTL_HOURS=720

# JWT Secrets
JWT_ACCESS_SECRET=traveling-super-secret-access-key-2026
JWT_REFRESH_SECRET=traveling-super-secret-refresh-key-2026
```

### Ý nghĩa:
- **Access Token (24h)**: Token dùng cho mọi API request. Hết hạn sau 24 giờ.
- **Refresh Token (30 days)**: Token dùng để lấy access token mới. Hết hạn sau 30 ngày.
- **Auto Refresh**: Khi access token hết hạn, hệ thống tự động dùng refresh token để lấy access token mới (không cần login lại).

### Trải nghiệm người dùng:
- ✅ Đăng nhập 1 lần → Sử dụng liên tục 24 giờ không bị logout
- ✅ Sau 24 giờ → Hệ thống tự động refresh token (không cần login lại)
- ✅ Sau 30 ngày → Phải đăng nhập lại

## Debug Token Issues

### Kiểm tra token trong localStorage

```javascript
// Mở Console (F12) và chạy:
console.log('Access Token:', localStorage.getItem('accessToken'));
console.log('Refresh Token:', localStorage.getItem('refreshToken'));
console.log('Auth Tokens:', localStorage.getItem('auth_tokens'));
console.log('User:', localStorage.getItem('user'));
```

### Kiểm tra token expiry

```javascript
// Decode JWT token (không cần secret)
function parseJwt(token) {
  try {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(c => {
      return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    return JSON.parse(jsonPayload);
  } catch (e) {
    return null;
  }
}

// Kiểm tra expiry
const token = localStorage.getItem('accessToken');
if (token) {
  const payload = parseJwt(token);
  console.log('Token Payload:', payload);
  console.log('Expires At:', new Date(payload.exp * 1000));
  console.log('Is Expired:', Date.now() > payload.exp * 1000);
}
```

### Test API với token

```javascript
// Test API call
fetch('http://localhost:8080/v1/api/users/14/bookings', {
  headers: {
    'Authorization': 'Bearer ' + localStorage.getItem('accessToken')
  }
})
.then(res => res.json())
.then(data => console.log('API Response:', data))
.catch(err => console.error('API Error:', err));
```

## Lưu ý quan trọng

### ⚠️ Khi nào cần logout và login lại?

1. **Thay đổi JWT secrets** (như lần này)
2. **Thay đổi JWT token structure** (thêm/bớt claims)
3. **Backend restart với code thay đổi token logic**
4. **Lỗi 401 liên tục** sau khi đã thử refresh

### ✅ Khi nào KHÔNG cần logout?

1. **Backend restart** (không thay đổi secrets)
2. **Frontend rebuild**
3. **Refresh trang**
4. **Đóng/mở trình duyệt**

## Troubleshooting

### Vẫn bị 401 sau khi login lại?

1. **Kiểm tra backend đang chạy:**
```bash
curl http://localhost:8080/health
```

2. **Kiểm tra .env file được load:**
```bash
# Trong server directory
cat .env | grep JWT
```

3. **Restart backend:**
```bash
cd server
go run cmd/server/main.go
```

4. **Kiểm tra logs:**
- Xem terminal backend có log `[TOKEN][ISSUE]` khi login
- Xem có error message gì không

### Token bị expired ngay sau khi login?

Kiểm tra system time:
```bash
date
```

Nếu system time sai → JWT expiry sẽ sai → Token invalid ngay lập tức.

## Best Practices

### Development
- ✅ Dùng JWT secrets đơn giản, dễ nhớ
- ✅ Access token: 15-60 phút
- ✅ Refresh token: 7-30 ngày
- ✅ Log token issue/refresh events

### Production
- ✅ Dùng JWT secrets phức tạp, random
- ✅ Access token: 5-15 phút
- ✅ Refresh token: 7-14 ngày
- ✅ Rotate secrets định kỳ
- ✅ Implement token blacklist
- ✅ Monitor token usage

---

**Tóm tắt:** Sau khi thay đổi JWT secrets, logout và login lại để lấy token mới. Token cũ không thể dùng được nữa.
