# Traveling Backend - MySQL Setup

## Hướng dẫn thiết lập Database

### 1. Cài đặt MySQL
Đảm bảo MySQL đã được cài đặt trên máy của bạn.

### 2. Tạo Database
Mở MySQL command line hoặc phpMyAdmin và chạy:

```sql
CREATE DATABASE traveling_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. Cấu hình file .env
Chỉnh sửa file `.env` với thông tin MySQL của bạn:

```env
DB_USER=root
DB_PASSWORD=your_mysql_password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=traveling_db
```

### 4. Chạy Server
```bash
go run main.go
```

Server sẽ tự động:
- Kết nối đến MySQL
- Tạo các bảng (users, tours)
- Seed dữ liệu mẫu

### 5. Kiểm tra
- API Tours: http://localhost:8080/api/tours
- API Login: POST http://localhost:8080/api/login
- API Register: POST http://localhost:8080/api/register

### Tài khoản mẫu
- Email: `test@example.com`
- Password: `123456`
