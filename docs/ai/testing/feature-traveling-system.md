---
phase: testing
title: Testing Strategy
feature: traveling-system
description: Chiến lược kiểm thử hệ thống đặt tour du lịch
---

# Testing Strategy — Traveling System

## Test Coverage Goals

- **Unit test**: 100% handlers backend (auth, tour, booking, invoice)
- **Integration test**: Toàn bộ critical paths — đăng ký → đặt tour → tạo hóa đơn
- **End-to-end**: Các luồng chính của khách hàng và nhân viên
- **Manual test**: UI/UX, responsive, edge cases giao diện

---

## Unit Tests

### handlers/auth.go

- [ ] **Test đăng ký thành công** — email mới, password hợp lệ → HTTP 200, user được tạo
- [ ] **Test đăng ký email trùng** → HTTP 409, message "Email đã được đăng ký"
- [ ] **Test đăng ký thiếu trường bắt buộc** → HTTP 400
- [ ] **Test đăng nhập thành công** — email + password đúng → HTTP 200 + JWT token
- [ ] **Test đăng nhập sai password** → HTTP 401
- [ ] **Test đăng nhập email không tồn tại** → HTTP 401
- [ ] **Test password được hash** — kiểm tra DB lưu bcrypt hash, không phải plaintext

### handlers/tour.go

- [ ] **Test lấy danh sách tour** → HTTP 200, trả về array
- [ ] **Test lấy chi tiết tour tồn tại** → HTTP 200, đúng fields
- [ ] **Test lấy chi tiết tour không tồn tại** → HTTP 404
- [ ] **Test tạo tour — role staff** → HTTP 201
- [ ] **Test tạo tour — role customer** → HTTP 403
- [ ] **Test cập nhật tour** → HTTP 200, data được cập nhật
- [ ] **Test xóa tour — role admin** → HTTP 200
- [ ] **Test xóa tour — role staff** → HTTP 403

### handlers/booking.go

- [ ] **Test đặt tour thành công** — còn đủ chỗ → HTTP 201, phiếu được tạo
- [ ] **Test đặt tour hết chỗ** — tổng khách vượt `SLKhachMax` → HTTP 400 "Tour đã hết chỗ"
- [ ] **Test đặt tour không tồn tại** → HTTP 404
- [ ] **Test đặt tour khi chưa đăng nhập** → HTTP 401
- [ ] **Test xem phiếu đặt của khách hàng** — chỉ thấy phiếu của mình → HTTP 200
- [ ] **Test nhân viên xem tất cả phiếu đặt** → HTTP 200, trả toàn bộ

### handlers/invoice.go

- [ ] **Test tạo hóa đơn thành công** — role staff, PDTour tồn tại → HTTP 201
- [ ] **Test tạo hóa đơn — role customer** → HTTP 403
- [ ] **Test tạo hóa đơn cho PDTour đã có hóa đơn** → HTTP 409
- [ ] **Test tính tổng tiền hóa đơn** — (NL + TreEm * 0.7) * chiPhi

### middleware/auth.go

- [ ] **Test request không có token** → HTTP 401
- [ ] **Test token hết hạn** → HTTP 401 "Token không hợp lệ"
- [ ] **Test token hợp lệ** → cho phép đi tiếp, `userID` được set vào context
- [ ] **Test token sai chữ ký** → HTTP 401

---

## Integration Tests

- [ ] **Luồng đăng ký → đăng nhập → lấy profile** — liên tục 3 API calls
- [ ] **Luồng nhân viên tạo tour → tạo lịch tour → khách đặt tour** — end-to-end nghiệp vụ
- [ ] **Luồng khách đặt tour → nhân viên tạo hóa đơn** — dữ liệu liên kết đúng
- [ ] **Test transaction rollback** — nếu tạo hóa đơn lỗi giữa chừng, PDTour không bị thay đổi trạng thái
- [ ] **Test concurrent booking** — 2 khách cùng đặt tour còn 1 chỗ → chỉ 1 thành công

---

## End-to-End Tests

### Luồng Khách Hàng

- [ ] **E2E-KH-01**: Truy cập trang chủ → xem danh sách tour → click vào tour → xem chi tiết
- [ ] **E2E-KH-02**: Đăng ký tài khoản → đăng nhập → đặt tour → xem lịch sử đặt tour
- [ ] **E2E-KH-03**: Cập nhật thông tin cá nhân (tên, email) → kiểm tra hiển thị đúng
- [ ] **E2E-KH-04**: Tìm kiếm tour theo địa điểm → kết quả đúng

### Luồng Nhân Viên

- [ ] **E2E-NV-01**: Đăng nhập nhân viên → vào Dashboard → xem danh sách phiếu đặt
- [ ] **E2E-NV-02**: Tạo tour mới → tạo lịch tour → kiểm tra hiển thị trên trang khách
- [ ] **E2E-NV-03**: Xem phiếu đặt của khách → tạo hóa đơn → hóa đơn xuất hiện trong danh sách

---

## Test Data

**Seed data cho môi trường test:**

```sql
-- Thành viên mẫu
INSERT INTO tblThanhVien VALUES
(1, 'khach01', '<bcrypt_hash>', '1990-05-15', 'khach01@test.com', 'customer'),
(2, 'nhanvien01', '<bcrypt_hash>', '1985-03-20', 'nv01@test.com', 'staff'),
(3, 'admin01', '<bcrypt_hash>', '1980-01-01', 'admin@test.com', 'admin');

-- Tour mẫu
INSERT INTO tblTour VALUES
(1, 'T001', 'Tour Đà Nẵng', '3 ngày 2 đêm', 'Máy bay', 20, 'Mô tả tour', 2000000),
(2, 'T002', 'Tour Phú Quốc', '5 ngày 4 đêm', 'Máy bay', 15, 'Mô tả tour', 5000000);

-- Lịch tour mẫu
INSERT INTO tblLichTour VALUES
(1, '2026-04-10', 1),
(2, '2026-04-20', 1),
(3, '2026-05-01', 2);
```

**Mock objects cho unit test (Go):**

```go
func mockDB(t *testing.T) *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&ThanhVien{}, &Tour{}, &LichTour{}, &PDTour{}, &HoaDon{})
    return db
}
```

---

## Test Reporting & Coverage

**Chạy tests backend:**
```bash
cd server
go test ./... -v -cover
go test ./handlers/... -coverprofile=coverage.out
go tool cover -html=coverage.out    # Xem HTML report
```

**Mục tiêu coverage:**
- `handlers/` : ≥ 90%
- `middleware/` : 100%
- `models/` : không cần unit test (chỉ là struct)

---

## Manual Testing

**Checklist UI/UX:**
- [ ] Form đăng ký hiển thị lỗi validation rõ ràng
- [ ] Thông báo lỗi không lộ thông tin kỹ thuật (stack trace, SQL error)
- [ ] Loading spinner hiển thị khi gọi API
- [ ] Trang responsive trên mobile (< 768px) và desktop
- [ ] Nút "Đặt Tour" disabled khi tour đã hết chỗ
- [ ] Redirect về `/login` khi token hết hạn
- [ ] Form đặt tour validate số người > 0 trước khi submit

**Kiểm tra bảo mật thủ công:**
- [ ] Thử truy cập `/api/tours` POST khi không có token → phải nhận 401
- [ ] Thử dùng token của khách hàng để gọi API tạo tour → phải nhận 403
- [ ] Kiểm tra password trong response API → không được xuất hiện

---

## Performance Testing

- [ ] **Load test** `/api/tours` với 50 concurrent requests → response < 500ms
- [ ] **Load test** `POST /api/pdtour` với 10 concurrent requests đặt cùng 1 lịch tour còn 5 chỗ → đúng số phiếu được tạo, không quá số chỗ
- Công cụ: `wrk` hoặc `hey` (Go HTTP load testing)

```bash
# Ví dụ với hey
hey -n 100 -c 10 http://localhost:8080/api/tours
```

---

## Bug Tracking

**Mức độ ưu tiên bug:**

| Mức | Mô tả | SLA |
|-----|-------|-----|
| 🔴 Critical | Mất dữ liệu, lỗ hổng bảo mật, không đặt được tour | Sửa ngay |
| 🟡 High | Tính năng chính bị sai, UI hiển thị sai dữ liệu | Sửa trong ngày |
| 🟢 Low | Lỗi UI nhỏ, text sai, không ảnh hưởng chức năng | Sửa trong sprint |

**Quy trình báo bug:**
1. Mô tả steps to reproduce
2. Expected vs Actual behavior
3. Screenshot / logs
4. Tạo issue trên Git repository
