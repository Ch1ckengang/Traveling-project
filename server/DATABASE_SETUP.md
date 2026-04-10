# Traveling Backend - PostgreSQL Setup

## Hướng dẫn thiết lập Database

### 1. Cài đặt PostgreSQL
Đảm bảo PostgreSQL đã được cài đặt trên máy của bạn.

### 2. Tạo Database
Mở `psql` và chạy:

```sql
CREATE DATABASE travel_db;
```

Hoặc chạy nhanh bằng Docker (khuyên dùng để theo dõi trên pgAdmin):

```bash
docker compose -f docker-compose.pgadmin.yml up -d
```

Mặc định sau khi chạy compose:
- PostgreSQL: `localhost:5432`
- pgAdmin: `http://localhost:5050`
- pgAdmin login: `admin@traveling.local` / `admin123`
- PostgreSQL login: `postgres` / `123456`

### 3. Cấu hình biến môi trường
Backend đọc cấu hình từ các biến sau (có default cho local):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=travel_db
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Ho_Chi_Minh
```

Bạn có thể sao chép từ file mẫu:

```bash
cp .env.example .env
```

Lưu ý:
- Nếu không đặt biến môi trường, backend sẽ dùng default.
- `DB_PASSWORD` nên được cấu hình rõ ràng khi chạy thật.
- Có thể dùng `DATABASE_URL` để override toàn bộ `DB_*`.

Ví dụ:

```env
DATABASE_URL=host=localhost user=postgres password=123456 dbname=travel_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh
```

### 4. Kết nối bằng pgAdmin
1. Mở `http://localhost:5050` và đăng nhập.
2. Chọn `Add New Server`.
3. Tab `General`: đặt tên ví dụ `travel-local`.
4. Tab `Connection`:
	- Host: `postgres` (nếu backend chạy trong cùng docker network) hoặc `localhost` (nếu backend chạy local)
	- Port: `5432`
	- Username: `postgres`
	- Password: `123456`
5. Save và mở database `travel_db` để xem bảng `users`, `tours`, `bookings`.

### 5. Chạy Server
```bash
go run main.go
```

Server sẽ tự động:
- Kết nối đến PostgreSQL
- Tạo/cập nhật các bảng (`users`, `tours`, `bookings`) bằng `AutoMigrate`
- Seed dữ liệu mẫu nếu bảng đang trống

### 6. Kiểm tra
- API Tours: `http://localhost:8080/v1/api/tours`
- API Login: `POST http://localhost:8080/v1/api/login`
- API Register: `POST http://localhost:8080/v1/api/register`

### 7. Import dữ liệu mẫu bằng SQL (tuỳ chọn)
Nếu cần seed bằng tay qua pgAdmin Query Tool, chạy file `init_database.sql` trong database `travel_db`.

### Tài khoản mẫu
- Email: `test@example.com`
- Password: `123456`
