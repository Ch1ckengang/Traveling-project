# Quản Lý Thông Tin Cá Nhân - Hoạt Động Nghiệp Vụ

## 1. Đăng Ký Tài Khoản (Register)

### Mô tả nghiệp vụ:
Cho phép người dùng mới tạo tài khoản trên hệ thống để có thể sử dụng các tính năng của ứng dụng du lịch.

### Yêu cầu nghiệp vụ:
- Thông tin bắt buộc:
  - Username: Tên đăng nhập duy nhất trong hệ thống
  - Email: Địa chỉ email hợp lệ và duy nhất
  - Password: Mật khẩu đủ mạnh (tối thiểu 6 ký tự)
  - Full Name: Họ và tên đầy đủ
  - Phone Number: Số điện thoại liên lạc

### Quy trình xử lý:
1. **Nhập thông tin**: Người dùng điền form đăng ký với các thông tin cần thiết
2. **Validation phía client**:
   - Kiểm tra các trường bắt buộc không được để trống
   - Kiểm tra định dạng email hợp lệ
   - Kiểm tra độ dài mật khẩu (≥ 6 ký tự)
   - Kiểm tra username không chứa ký tự đặc biệt
   - Kiểm tra số điện thoại đúng định dạng
3. **Gửi request đăng ký**: Gửi thông tin đến server qua API endpoint `/api/register`
4. **Validation phía server**:
   - Kiểm tra username đã tồn tại chưa
   - Kiểm tra email đã được sử dụng chưa
   - Validate lại tất cả các trường dữ liệu
5. **Mã hóa mật khẩu**: Hash password bằng bcrypt trước khi lưu vào database
6. **Lưu thông tin**: Tạo bản ghi mới trong bảng users
7. **Phản hồi**:
   - Thành công: Trả về thông báo đăng ký thành công
   - Thất bại: Trả về lỗi cụ thể (username/email đã tồn tại, dữ liệu không hợp lệ)

### Luồng xử lý lỗi:
- Username đã tồn tại → Hiển thị "Tên đăng nhập đã được sử dụng"
- Email đã tồn tại → Hiển thị "Email đã được đăng ký"
- Dữ liệu không hợp lệ → Hiển thị lỗi validation tương ứng
- Lỗi server → Hiển thị "Đăng ký thất bại, vui lòng thử lại"

---

## 2. Đăng Nhập (Login)

### Mô tả nghiệp vụ:
Cho phép người dùng đã có tài khoản xác thực và truy cập vào hệ thống.

### Yêu cầu nghiệp vụ:
- Username hoặc Email
- Password
- Hệ thống phải xác thực thông tin đăng nhập
- Tạo và quản lý phiên đăng nhập (session/token)

### Quy trình xử lý:
1. **Nhập thông tin đăng nhập**: Người dùng nhập username/email và password
2. **Validation phía client**:
   - Kiểm tra các trường không được để trống
   - Kiểm tra định dạng cơ bản
3. **Gửi request đăng nhập**: Gửi thông tin đến server qua API endpoint `/api/login`
4. **Xác thực phía server**:
   - Tìm kiếm user theo username hoặc email
   - So sánh password đã hash với password trong database
5. **Tạo phiên đăng nhập**:
   - Nếu xác thực thành công: Tạo JWT token hoặc session
   - Token chứa thông tin user (user_id, username, role)
6. **Phản hồi**:
   - Thành công: Trả về token và thông tin user cơ bản
   - Thất bại: Trả về lỗi "Thông tin đăng nhập không chính xác"
7. **Lưu trạng thái**: Client lưu token vào localStorage/cookie
8. **Chuyển hướng**: Redirect đến trang chủ hoặc trang trước đó

### Luồng xử lý lỗi:
- Tài khoản không tồn tại → "Tài khoản không tồn tại"
- Mật khẩu sai → "Mật khẩu không chính xác"
- Tài khoản bị khóa → "Tài khoản đã bị khóa"
- Lỗi server → "Đăng nhập thất bại, vui lòng thử lại"

### Bảo mật:
- Giới hạn số lần đăng nhập thất bại (rate limiting)
- Mã hóa password khi truyền (HTTPS)
- Token có thời gian hết hạn
- Không hiển thị thông tin chi tiết lỗi (tránh lộ thông tin hệ thống)

---

## 3. Chỉnh Sửa Thông Tin Cá Nhân (Edit Profile)

### Mô tả nghiệp vụ:
Cho phép người dùng đã đăng nhập cập nhật thông tin cá nhân của mình.

### Yêu cầu nghiệp vụ:
- Người dùng phải đã đăng nhập
- Chỉ có thể chỉnh sửa thông tin của chính mình
- Một số thông tin có thể chỉnh sửa:
  - Full Name (Họ tên)
  - Email (với xác thực)
  - Phone Number (Số điện thoại)
  - Avatar/Profile Picture (Ảnh đại diện)
  - Bio/Description (Giới thiệu)
  - Address (Địa chỉ)
- Username thường không cho phép thay đổi

### Quy trình xử lý:
1. **Hiển thị thông tin hiện tại**:
   - Load thông tin user từ database
   - Hiển thị trong form với các giá trị đã điền sẵn
2. **Chỉnh sửa thông tin**: Người dùng cập nhật các trường muốn thay đổi
3. **Validation phía client**:
   - Kiểm tra định dạng email hợp lệ (nếu thay đổi)
   - Kiểm tra số điện thoại đúng định dạng
   - Kiểm tra kích thước ảnh (nếu upload avatar)
4. **Upload file** (nếu có):
   - Upload avatar lên server/cloud storage
   - Validate: loại file, kích thước, định dạng
5. **Gửi request cập nhật**: Gửi dữ liệu đến API endpoint `/api/profile/update`
6. **Xác thực quyền truy cập**:
   - Verify JWT token
   - Kiểm tra user đang chỉnh sửa thông tin của chính mình
7. **Validation phía server**:
   - Kiểm tra email mới chưa được sử dụng bởi user khác
   - Kiểm tra số điện thoại chưa được sử dụng
   - Validate tất cả các trường dữ liệu
8. **Cập nhật database**:
   - Update bản ghi trong bảng users
   - Lưu URL của avatar (nếu có)
9. **Phản hồi**:
   - Thành công: Trả về thông tin user đã cập nhật
   - Thất bại: Trả về lỗi cụ thể
10. **Cập nhật UI**: Hiển thị thông tin mới và thông báo thành công

### Các trường hợp đặc biệt:
- **Thay đổi Email**:
  - Gửi email xác nhận đến địa chỉ mới
  - Email chỉ được cập nhật sau khi xác nhận
  - Gửi thông báo đến email cũ
  
- **Thay đổi Password**:
  - Yêu cầu nhập password cũ
  - Password mới phải đủ mạnh
  - Yêu cầu nhập lại password mới
  - Hash password mới trước khi lưu
  - Logout tất cả các session cũ (tùy chọn)

- **Upload Avatar**:
  - Giới hạn dung lượng (VD: 2MB)
  - Chỉ chấp nhận file ảnh (jpg, png, gif)
  - Resize ảnh về kích thước chuẩn
  - Xóa ảnh cũ khỏi storage

### Luồng xử lý lỗi:
- Không có quyền truy cập → "Unauthorized" - redirect to login
- Email đã được sử dụng → "Email đã tồn tại trong hệ thống"
- File ảnh quá lớn → "Kích thước ảnh vượt quá giới hạn cho phép"
- Định dạng file không hợp lệ → "Chỉ chấp nhận file ảnh jpg, png, gif"
- Lỗi server → "Cập nhật thất bại, vui lòng thử lại"

### Bảo mật:
- Verify JWT token cho mọi request
- Kiểm tra quyền sở hữu (user chỉ sửa được thông tin của mình)
- Sanitize input để tránh XSS
- Giới hạn kích thước và loại file upload
- Rate limiting để tránh spam

---

## 4. Xem Thông Tin Cá Nhân (View Profile)

### Mô tả nghiệp vụ:
Hiển thị thông tin chi tiết của người dùng.

### Quy trình xử lý:
1. **Request thông tin**: Gửi GET request đến `/api/profile/:userId`
2. **Xác thực**: Verify JWT token
3. **Lấy dữ liệu**: Query thông tin user từ database
4. **Phân quyền hiển thị**:
   - Thông tin công khai: Hiển thị cho mọi người
   - Thông tin riêng tư: Chỉ hiển thị cho chính user
5. **Phản hồi**: Trả về thông tin user
6. **Hiển thị UI**: Render thông tin trên giao diện

### Phân quyền thông tin:
- **Công khai**: Username, Full Name, Avatar, Bio
- **Riêng tư**: Email, Phone Number, Address (chỉ chính user thấy)

---

## 5. Đăng Xuất (Logout)

### Mô tả nghiệp vụ:
Kết thúc phiên đăng nhập của người dùng.

### Quy trình xử lý:
1. **Gửi request logout**: Gọi API endpoint `/api/logout`
2. **Xóa token phía server**: Invalidate token (nếu sử dụng blacklist)
3. **Xóa token phía client**: 
   - Xóa token khỏi localStorage/cookie
   - Clear auth context/state
4. **Chuyển hướng**: Redirect về trang login hoặc trang chủ
5. **Cập nhật UI**: Hiển thị trạng thái chưa đăng nhập

---

## 6. Sơ Đồ Luồng Dữ Liệu

User -> Client (React) -> API Server (Go) -> Database (PostgreSQL/MySQL)
                ↓                 ↓
         Validation        Authentication
         State Mgmt        Authorization
         Error Handle      Business Logic

## 7. API Endpoints

- `POST /api/register` - Đăng ký tài khoản mới
- `POST /api/login` - Đăng nhập
- `POST /api/logout` - Đăng xuất
- `GET /api/profile/:userId` - Xem thông tin cá nhân
- `PUT /api/profile/update` - Cập nhật thông tin cá nhân
- `PUT /api/profile/change-password` - Thay đổi mật khẩu
- `POST /api/profile/upload-avatar` - Upload ảnh đại diện

## 8. Mô Hình Dữ Liệu (User Model)

users {
  id: int (primary key)
  username: string (unique)
  email: string (unique)
  password: string (hashed)
  full_name: string
  phone_number: string
  avatar_url: string
  bio: text
  address: string
  created_at: timestamp
  updated_at: timestamp
  is_active: boolean
  role: string (user/admin)
}
