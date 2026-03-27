---
phase: implementation
title: Implementation Guide
feature: international-tour-database-refresh
description: Hướng dẫn triển khai chuẩn hóa dữ liệu tour quốc tế
---

# Implementation Guide

## Development Setup
**How do we get started?**

- Backend: `cd server && go run .`
- Frontend: `cd client && npm run dev`
- Build check: `go build ./...` và `npm run build`

## Code Structure
**How is the code organized?**

- `server/models/models.go`: Tour schema
- `server/repositories/tour_repository.go`: seed/backfill tour chuẩn
- `server/services/tour_service.go`: filter/sort
- `server/controllers/tour_controller.go`: query mapping
- `server/init_database.sql`: SQL setup cho MySQL

## Implementation Notes
**Key technical details to remember:**

### Core Features
- Feature 1: Seed dữ liệu tour quốc tế chuẩn theo danh sách tên tour.
- Feature 2: Backfill idempotent theo `name` để tránh trùng dữ liệu.
- Feature 3: Đồng bộ schema SQL (`type`, `country`) với runtime model.

### Patterns & Best Practices
- Không xóa dữ liệu cũ hàng loạt trong seed/backfill.
- Update chọn lọc các trường chuẩn hóa cho bản ghi đã tồn tại.
- Giữ danh sách tour seed tập trung một nơi để dễ bảo trì.

## Integration Points
**How do pieces connect?**

- Startup gọi `seedData()` và luồng API tour gọi `CreateToursIfEmpty()`.
- Frontend đọc dữ liệu qua `GET /v1/api/tours` với `category=international`.

## Error Handling
**How do we handle failures?**

- Trả lỗi khi migrate/query/update fail ở repository.
- Controller trả HTTP 500 với message rõ ràng cho lỗi seed/lấy tour.

## Performance Considerations
**How do we keep it fast?**

- Dataset nhỏ, backfill theo vòng lặp là đủ.
- Có thể tối ưu sang upsert bulk nếu số lượng tour tăng mạnh.

## Security Notes
**What security measures are in place?**

- Seed không chứa dữ liệu nhạy cảm.
- API danh sách tour read-only, không ghi từ phía client public.
