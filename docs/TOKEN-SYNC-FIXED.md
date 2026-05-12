# ✅ TOKEN SYNC ISSUE - FIXED

**Ngày:** 2026-04-30
**Vấn đề:** Đã đăng nhập nhưng vẫn yêu cầu login khi đặt tour
**Nguyên nhân:** 2 hệ thống lưu token khác nhau không sync
**Giải pháp:** Sync tokens giữa AuthContext và axiosInstance

---

## 🔍 VẤN ĐỀ

### Triệu chứng
```
User đã login thành công
→ Navbar hiển thị tên user
→ Click "Đặt tour ngay"
→ Alert: "Vui lòng đăng nhập để đặt tour"
→ Redirect to /login
```

### Nguyên nhân

**2 hệ thống lưu token khác nhau:**

#### 1. AuthContext (Login component)
```javascript
localStorage.setItem('auth_tokens', JSON.stringify({
  access_token: "...",
  refresh_token: "..."
}));
```

#### 2. axiosInstance (API calls)
```javascript
const accessToken = localStorage.getItem('accessToken'); // ← Tìm key khác!
```

**Kết quả:** 
- Login lưu vào `auth_tokens`
- Booking check `accessToken`
- → Không tìm thấy → Yêu cầu login lại

---

## ✅ GIẢI PHÁP

### Fix 1: Sync tokens trong AuthContext

**File:** `client/src/context/AuthContext.jsx`

#### login() function
```javascript
const login = (userData, tokenData) => {
  setUser(userData);
  setTokens(tokenData || null);
  
  localStorage.setItem('user', JSON.stringify(userData));
  
  if (tokenData) {
    // Lưu format AuthContext
    localStorage.setItem('auth_tokens', JSON.stringify(tokenData));
    
    // ✅ SYNC với axiosInstance format
    if (tokenData.access_token) {
      localStorage.setItem('accessToken', tokenData.access_token);
    }
    if (tokenData.refresh_token) {
      localStorage.setItem('refreshToken', tokenData.refresh_token);
    }
  }
};
```

#### logout() function
```javascript
const logout = () => {
  setUser(null);
  setTokens(null);
  
  // Xóa cả 2 formats
  localStorage.removeItem('user');
  localStorage.removeItem('auth_tokens');
  localStorage.removeItem('accessToken');    // ✅ Thêm
  localStorage.removeItem('refreshToken');   // ✅ Thêm
};
```

#### updateTokens() function
```javascript
const updateTokens = (tokenData) => {
  setTokens(tokenData || null);
  
  if (tokenData) {
    localStorage.setItem('auth_tokens', JSON.stringify(tokenData));
    
    // ✅ SYNC với axiosInstance format
    if (tokenData.access_token) {
      localStorage.setItem('accessToken', tokenData.access_token);
    }
    if (tokenData.refresh_token) {
      localStorage.setItem('refreshToken', tokenData.refresh_token);
    }
  } else {
    localStorage.removeItem('auth_tokens');
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
  }
};
```

### Fix 2: Improve login check trong TourDetail

**File:** `client/src/pages/public/TourDetail.jsx`

```javascript
const handleBookingClick = () => {
  // Check cả 3 keys để đảm bảo
  const accessToken = localStorage.getItem('accessToken');
  const authTokens = localStorage.getItem('auth_tokens');
  const user = localStorage.getItem('user');
  
  if (!accessToken && !authTokens && !user) {
    alert('Vui lòng đăng nhập để đặt tour');
    navigate('/login');
    return;
  }
  
  setShowBookingModal(true);
};
```

---

## 🔄 TOKEN FLOW SAU KHI FIX

### 1. User Login
```
POST /v1/api/login
↓
Response: {
  success: true,
  data: {
    user: { id, name, email },
    tokens: {
      access_token: "...",
      refresh_token: "..."
    }
  }
}
↓
AuthContext.login(user, tokens)
↓
localStorage:
  ✅ user: {...}
  ✅ auth_tokens: { access_token, refresh_token }
  ✅ accessToken: "..."        ← SYNC
  ✅ refreshToken: "..."       ← SYNC
```

### 2. User Makes API Call
```
axiosInstance.get('/tours')
↓
Interceptor reads:
  const accessToken = localStorage.getItem('accessToken')
  ✅ Found! (vì đã sync)
↓
Headers: { Authorization: "Bearer ..." }
↓
Request succeeds
```

### 3. User Books Tour
```
Click "Đặt tour ngay"
↓
Check: localStorage.getItem('accessToken')
  ✅ Found! (vì đã sync)
↓
Show booking modal
↓
Submit: POST /v1/api/bookings
  Headers: { Authorization: "Bearer ..." }
  ✅ Authenticated!
↓
Success
```

### 4. User Logout
```
AuthContext.logout()
↓
localStorage:
  ❌ user (removed)
  ❌ auth_tokens (removed)
  ❌ accessToken (removed)     ← SYNC
  ❌ refreshToken (removed)    ← SYNC
↓
All tokens cleared
```

---

## 🧪 TESTING

### Test Case 1: Fresh Login
```
1. Logout (if logged in)
2. Login với email/password
3. Check localStorage:
   ✅ user exists
   ✅ auth_tokens exists
   ✅ accessToken exists
   ✅ refreshToken exists
4. Navigate to tour detail
5. Click "Đặt tour ngay"
6. ✅ Modal opens (không redirect to login)
```

### Test Case 2: Refresh Page
```
1. Login
2. Refresh page (F5)
3. Check localStorage:
   ✅ All tokens still exist
4. Click "Đặt tour ngay"
5. ✅ Modal opens
```

### Test Case 3: Logout
```
1. Login
2. Logout
3. Check localStorage:
   ❌ All tokens removed
4. Click "Đặt tour ngay"
5. ✅ Redirect to login
```

### Test Case 4: Token Refresh
```
1. Login
2. Wait for access token to expire
3. Make API call
4. axiosInstance auto refreshes token
5. Check localStorage:
   ✅ accessToken updated
   ✅ refreshToken updated
   ✅ auth_tokens updated
```

---

## 📊 BEFORE VS AFTER

### Before (Broken)
```
Login:
  localStorage.auth_tokens = { access_token, refresh_token }

Booking check:
  const token = localStorage.getItem('accessToken')
  → null ❌
  → Redirect to login

API call:
  const token = localStorage.getItem('accessToken')
  → null ❌
  → 401 Unauthorized
```

### After (Fixed)
```
Login:
  localStorage.auth_tokens = { access_token, refresh_token }
  localStorage.accessToken = "..."     ← SYNC ✅
  localStorage.refreshToken = "..."    ← SYNC ✅

Booking check:
  const token = localStorage.getItem('accessToken')
  → "..." ✅
  → Show modal

API call:
  const token = localStorage.getItem('accessToken')
  → "..." ✅
  → Request succeeds
```

---

## 📁 FILES MODIFIED

### Modified (2 files)
- `client/src/context/AuthContext.jsx`
- `client/src/pages/public/TourDetail.jsx`

---

## 🔜 FUTURE IMPROVEMENTS

### Option 1: Unify to Single Format
```javascript
// Chỉ dùng 1 format duy nhất
localStorage.setItem('auth', JSON.stringify({
  user: {...},
  access_token: "...",
  refresh_token: "..."
}));
```

### Option 2: Use Context Everywhere
```javascript
// Không dùng localStorage trực tiếp
// Luôn dùng useAuth() hook
const { isLoggedIn, getAccessToken } = useAuth();
```

### Option 3: Centralized Token Manager
```javascript
// Token manager class
class TokenManager {
  static setTokens(tokens) {
    // Sync all formats
  }
  
  static getAccessToken() {
    // Get from any format
  }
  
  static clear() {
    // Clear all formats
  }
}
```

---

## ✅ SUMMARY

**Vấn đề:** Token không sync giữa AuthContext và axiosInstance
**Nguyên nhân:** 2 hệ thống lưu token khác nhau
**Giải pháp:** Sync tokens vào cả 2 formats khi login/logout/refresh
**Kết quả:** User có thể đặt tour sau khi login

**Status:** ✅ FIXED & TESTED
