# Traveling Backend

Tai lieu nay mo ta backend theo huong nghiep vu va luong xu ly cho tung chuc nang.

## 1. Tong quan kien truc

Backend dang dung:
- `Gin` cho HTTP API
- `GORM` cho truy cap DB
- `SQLite` (file `travel.db`) cho luu tru

Kien truc theo 3 tang:
- `controllers/`: nhan request, goi service, tra response
- `services/`: xu ly business logic va validate nghiep vu
- `repositories/`: tuong tac truc tiep voi database

Luong tong quat:
1. Client goi API.
2. Controller bind JSON/query va goi service.
3. Service validate + xu ly nghiep vu + goi repository.
4. Repository doc/ghi DB.
5. Controller map loi -> HTTP status va tra JSON.

## 2. Khoi dong backend

Tu thu muc `server/`:

```bash
go run main.go
```

Server chay tai `http://localhost:8080`.

Khi khoi dong:
1. `database.Connect()` ket noi SQLite (`travel.db`).
2. `AutoMigrate` tao/cap nhat bang `users`, `tours`, `bookings`.
3. `seedData()` seed user/tour mau neu bang trong.
4. Dang ky route duoi prefix `/v1`.

## 3. Danh sach API

Base path: `/v1`

- `GET /api/tours`
- `GET /api/tours/domestic`
- `POST /api/bookings`
- `POST /api/register`
- `POST /api/login`
- `PUT /api/users/:id`

## 4. Nghiep vu theo tung chuc nang

### 4.1 Dang ky tai khoan (`POST /api/register`)

File lien quan:
- `controllers/user_controller.go`
- `services/user_service.go`
- `repositories/user_repository.go`

Luong xu ly:
1. Controller bind body vao `RegisterRequest`.
2. Service chuan hoa input:
   - `name`: trim space
   - `email`: trim + lowercase
3. Service validate:
   - `name` khong duoc rong
   - `email` dung format
   - `password` >= 8 ky tu
4. Service kiem tra email da ton tai chua.
5. Neu hop le, service hash password bang `bcrypt.GenerateFromPassword`.
6. Tao user moi qua repository.
7. Tra ve `AuthResponse` thanh cong (khong tra password vi field password duoc an trong JSON).

Loi nghiep vu thuong gap:
- Email da dang ky -> `409 Conflict`
- Input khong hop le -> `400 Bad Request`
- Loi he thong -> `500 Internal Server Error`

### 4.2 Dang nhap (`POST /api/login`)

File lien quan:
- `controllers/user_controller.go`
- `services/user_service.go`
- `repositories/user_repository.go`

Luong xu ly:
1. Controller bind body vao `LoginRequest`.
2. Service chuan hoa email (`trim + lowercase`) va validate input khong rong.
3. Service tim user theo email.
4. Service so sanh password bang `bcrypt.CompareHashAndPassword`.
5. Neu dung, tra thong tin user.
6. Neu sai email hoac password, tra 1 thong diep chung de tranh lo user co ton tai hay khong.

Loi nghiep vu thuong gap:
- Sai thong tin dang nhap -> `401 Unauthorized`
- Input khong hop le -> `400 Bad Request`
- Loi he thong -> `500 Internal Server Error`

### 4.3 Cap nhat thong tin user (`PUT /api/users/:id`)

File lien quan:
- `controllers/user_controller.go`
- `services/user_service.go`
- `repositories/user_repository.go`

Luong xu ly:
1. Controller doc `:id` va bind `UpdateUserRequest`.
2. Service tim user hien tai theo ID.
3. Neu doi email, kiem tra email moi khong trung user khac.
4. Chi cap nhat cac truong duoc gui len (`name/email/password`).
5. Neu co doi password, service hash password bang bcrypt truoc khi luu.
6. Luu user qua repository.

### 4.4 Lay danh sach tour (`GET /api/tours`, `GET /api/tours/domestic`)

File lien quan:
- `controllers/tour_controller.go`
- `services/tour_service.go`
- `repositories/tour_repository.go`

Luong xu ly:
1. Controller dam bao co du lieu tour co ban (`CreateToursIfEmpty`).
2. Doc query filter (`category`, `q/city`, `duration`, `price`, `sort`).
3. Service lay danh sach theo category.
4. Service loc theo thanh pho, thoi luong, gia.
5. Service sap xep theo query (`price`, `duration`, `name`, `latest`).
6. Tra ve danh sach tours.

### 4.5 Dat tour (`POST /api/bookings`)

File lien quan:
- `controllers/booking_controller.go`
- `services/booking_service.go`
- `repositories/booking_repository.go`

Luong xu ly:
1. Controller bind body vao `CreateBookingRequest`.
2. Service normalize du lieu (`full_name`, `phone`, `email`, `travel_date`, `note`).
3. Service validate:
   - `tour_id` hop le
   - `full_name` khong rong
   - `phone` dung regex
   - `email` dung format
   - `quantity > 0`
   - `travel_date` dung format `YYYY-MM-DD` va khong o qua khu
4. Service kiem tra tour co ton tai.
5. Tao booking voi status mac dinh `booked`.
6. Luu DB qua repository va tra `201 Created`.

## 5. Rule bao mat chinh

1. Khong luu plain text password trong DB.
2. Tat ca password phai duoc hash bang bcrypt truoc khi luu.
3. Dang nhap chi so sanh bang hash (`CompareHashAndPassword`).
4. Truong `password` cua `User` khong duoc tra ve JSON response (`json:"-"`).

## 6. Cau truc thu muc backend

```text
server/
  main.go
  controllers/
  services/
  repositories/
  models/
  database/
  travel.db
```

## 7. Tai lieu lien quan

- `server/DATABASE_SETUP.md`: huong dan cau hinh database
- `server/init_database.sql`: script SQL tham khao
