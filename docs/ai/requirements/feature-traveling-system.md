---
phase: requirements
title: Requirements & Problem Understanding
feature: traveling-system
description: Phân tích yêu cầu toàn bộ hệ thống đặt tour du lịch trực tuyến
---

# Requirements & Problem Understanding — Traveling System

## Problem Statement

**Vấn đề đang giải quyết:**

- Khách hàng khó tìm kiếm, so sánh và đặt tour du lịch một cách nhanh chóng, minh bạch
- Doanh nghiệp du lịch thiếu công cụ quản lý tour, lịch khởi hành, đặt chỗ và xuất hóa đơn tập trung
- Quy trình xử lý đặt tour thủ công dẫn đến sai sót và tốn nhiều nhân lực

**Đối tượng bị ảnh hưởng:**
- Khách hàng cá nhân muốn đặt tour du lịch trong nước
- Nhân viên công ty du lịch (tư vấn, lập hóa đơn, quản lý tour)
- Quản trị viên hệ thống

**Tình trạng hiện tại:**
- Hệ thống đang có backend Go (Gin + GORM + MySQL) với các API cơ bản: đăng nhập, đăng ký, cập nhật thông tin, lấy danh sách tour
- Frontend React (Vite) đã có giao diện đăng nhập, đăng ký, hồ sơ cá nhân và tìm kiếm tour
- Chưa có tính năng đặt tour, quản lý lịch tour, hóa đơn

---

## Goals & Objectives

**Mục tiêu chính (In-scope):**
1. Quản lý tài khoản thành viên — đăng ký, đăng nhập, chỉnh sửa hồ sơ
2. Quản lý thông tin tour — thêm, sửa, xóa, tìm kiếm tour
3. Quản lý lịch tour — tạo và theo dõi lịch khởi hành cho từng tour
4. Đặt tour — khách hàng chọn tour, chọn lịch và gửi phiếu đặt
5. Quản lý hóa đơn — nhân viên xác nhận và xuất hóa đơn
6. Quản lý địa điểm và dịch vụ — gán điểm đến và dịch vụ cho từng tour

**Mục tiêu phụ:**
- Hiển thị thông tin tour hấp dẫn với ảnh, mô tả, giá
- Thống kê doanh thu và số lượng khách

**Ngoài phạm vi (Non-goals):**
- Tích hợp cổng thanh toán trực tuyến (VNPay, Momo) — giai đoạn sau
- App mobile — chỉ Web
- Đa ngôn ngữ — chỉ Tiếng Việt
- Đánh giá / review tour

---

## User Stories & Use Cases

### Khách Hàng (tblKhachHang)

| # | User Story | Mức độ ưu tiên |
|---|------------|----------------|
| KH-01 | Là khách hàng, tôi muốn **đăng ký tài khoản** để có thể đặt tour | 🔴 Cao |
| KH-02 | Là khách hàng, tôi muốn **đăng nhập** để truy cập hệ thống | 🔴 Cao |
| KH-03 | Là khách hàng, tôi muốn **xem danh sách tour** để tìm tour phù hợp | 🔴 Cao |
| KH-04 | Là khách hàng, tôi muốn **tìm kiếm tour theo địa điểm** để thu hẹp lựa chọn | 🟡 Trung bình |
| KH-05 | Là khách hàng, tôi muốn **xem chi tiết tour** (ảnh, mô tả, lịch, giá) để quyết định đặt | 🔴 Cao |
| KH-06 | Là khách hàng, tôi muốn **đặt tour** (chọn lịch, số người lớn, số trẻ em) để xác nhận chuyến đi | 🔴 Cao |
| KH-07 | Là khách hàng, tôi muốn **xem lịch sử đặt tour** để theo dõi các chuyến đi đã đặt | 🟡 Trung bình |
| KH-08 | Là khách hàng, tôi muốn **chỉnh sửa thông tin cá nhân** để cập nhật hồ sơ | 🟡 Trung bình |

### Nhân Viên (tblNhanVien)

| # | User Story | Mức độ ưu tiên |
|---|------------|----------------|
| NV-01 | Là nhân viên, tôi muốn **quản lý thông tin tour** (CRUD) để duy trì catalog tour | 🔴 Cao |
| NV-02 | Là nhân viên, tôi muốn **tạo lịch tour** (ngày khởi hành, ngày về) để khách hàng lựa chọn | 🔴 Cao |
| NV-03 | Là nhân viên, tôi muốn **xem danh sách phiếu đặt** để xử lý đơn hàng | 🔴 Cao |
| NV-04 | Là nhân viên, tôi muốn **tạo hóa đơn** cho phiếu đặt được xác nhận | 🔴 Cao |
| NV-05 | Là nhân viên, tôi muốn **quản lý địa điểm** trong hệ thống để gán vào tour | 🟡 Trung bình |
| NV-06 | Là nhân viên, tôi muốn **quản lý dịch vụ** (khách sạn, ăn uống) theo từng điểm dừng | 🟢 Thấp |

### Quản Trị Viên

| # | User Story | Mức độ ưu tiên |
|---|------------|----------------|
| QTV-01 | Là quản trị viên, tôi muốn **quản lý tài khoản thành viên** để kiểm soát quyền truy cập | 🟡 Trung bình |
| QTV-02 | Là quản trị viên, tôi muốn **xem thống kê tổng hợp** (doanh thu, số tour, số khách) | 🟢 Thấp |

---

## Success Criteria

| Tiêu chí | Cách đo lường |
|----------|---------------|
| Khách hàng đặt tour thành công | Phiếu đặt được tạo và lưu vào DB, nhân viên nhìn thấy |
| Nhân viên tạo hóa đơn thành công | Hóa đơn liên kết đúng với phiếu đặt và thành viên |
| Tìm kiếm tour hoạt động | Trả về kết quả đúng trong < 1 giây |
| Xác thực tài khoản an toàn | Password được hash, không lộ ra response |
| Giao diện hiển thị đúng | Responsive trên desktop và mobile |
| Không mất dữ liệu | Transaction DB rollback khi lỗi |

---

## Constraints & Assumptions

**Ràng buộc kỹ thuật:**
- Backend: Go 1.21+ với Gin framework và GORM ORM
- Database: MySQL 8.0 với charset utf8mb4
- Frontend: React 18 + Vite, không dùng TypeScript
- Chưa có authentication middleware (JWT) — cần bổ sung
- Password hiện tại lưu plaintext — phải chuyển sang bcrypt

**Ràng buộc nghiệp vụ:**
- Một lịch tour có số lượng khách tối đa (`SLKhachMax` trong `tblTour`)
- Hóa đơn chỉ được tạo bởi nhân viên (không phải khách tự tạo)
- Phiếu đặt cần có ít nhất 1 khách người lớn

**Giả định:**
- Hệ thống chạy nội bộ, không yêu cầu CDN hay load balancer
- Ảnh tour được lưu local hoặc qua URL bên ngoài
- Múi giờ: UTC+7 (Việt Nam)

---

## Questions & Open Items

| # | Câu hỏi | Trạng thái |
|---|---------|------------|
| Q-01 | Password hiện tại lưu plaintext — khi nào migrate sang bcrypt? | ⏳ Cần xác nhận |
| Q-02 | JWT token hay Session-based authentication? | ⏳ Cần xác nhận |
| Q-03 | Trạng thái phiếu đặt tour: pending → confirmed → cancelled? | ⏳ Cần xác nhận |
| Q-04 | Giá tour có phân biệt người lớn / trẻ em không? | ⏳ Cần xác nhận |
| Q-05 | Upload ảnh tour: lưu server hay dùng cloud storage? | ⏳ Cần xác nhận |
| Q-06 | Nhân viên và khách hàng login cùng endpoint hay khác nhau? | ⏳ Cần xác nhận |
