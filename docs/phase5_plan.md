# Implementation Plan: Phase 5 (Admin Dashboard & Reports)

## Mục tiêu
Hoàn thành tính năng thống kê và báo cáo cho trang Admin Dashboard. Module này giúp quản trị viên theo dõi hiệu suất hệ thống như tổng doanh thu, số lượng tour, số lượng người dùng và thống kê đặt chỗ.

## Kế hoạch triển khai

### 1. Backend: Dashboard APIs (`server/internal/dashboard/`)
Tạo module `dashboard` chuyên xử lý các thống kê:
- **`GET /api/v1/admin/dashboard/summary`**: Trả về các chỉ số tổng quan (Total Users, Total Tours, Total Bookings, Total Revenue).
- **`GET /api/v1/admin/dashboard/revenue-chart`**: Trả về dữ liệu doanh thu theo ngày/tháng để vẽ biểu đồ.
- **`GET /api/v1/admin/dashboard/top-tours`**: Trả về danh sách top 5 tour có lượng đặt nhiều nhất.
- Đăng ký các route này trong `main.go`.

### 2. Frontend: Cập nhật Dashboard Component (`client/src/pages/admin/Dashboard.jsx`)
- Sử dụng thư viện `recharts` để hiển thị biểu đồ trực quan.
- Tích hợp 4 cards thống kê tổng quan (Người dùng, Đặt chỗ, Doanh thu, Tour).
- Hiển thị bảng "Các Tour Nổi Bật Nhất".
- Gọi API từ `services/dashboardService.js`.

### 3. Frontend: Cập nhật Reports Component (`client/src/pages/admin/Reports.jsx`)
- Xây dựng trang báo cáo chi tiết.
- Cung cấp chức năng tải xuống danh sách Bookings dưới định dạng `.csv`.
