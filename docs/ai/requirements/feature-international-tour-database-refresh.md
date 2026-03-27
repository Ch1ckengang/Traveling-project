---
phase: requirements
title: Requirements & Problem Understanding
feature: international-tour-database-refresh
description: Làm lại dữ liệu tour, bổ sung và chuẩn hóa danh sách tour quốc tế
---

# Requirements & Problem Understanding

## Problem Statement
**What problem are we solving?**

- Dữ liệu tour hiện chưa đồng nhất giữa các nguồn seed (GORM seed/backfill và SQL khởi tạo), khiến lúc có/lúc thiếu tên tour quốc tế.
- Người dùng khi chọn tab Du lịch quốc tế có thể thấy ít tour hoặc dữ liệu không chuẩn (thiếu quốc gia/type).
- Team vận hành khó khởi tạo môi trường mới do script SQL chưa phản ánh mô hình dữ liệu hiện tại (`type`, `country`).

## Goals & Objectives
**What do we want to achieve?**

- Primary goals:
- Chuẩn hóa bộ dữ liệu tour để luôn có danh sách tour quốc tế rõ tên tour, quốc gia, thời lượng, mô tả.
- Đồng bộ dữ liệu giữa code seed/backfill và script SQL khởi tạo.
- Đảm bảo dữ liệu cũ được backfill không phá vỡ dữ liệu hiện có.

- Secondary goals:
- Cải thiện khả năng mở rộng để thêm tour quốc tế mới theo danh sách chuẩn.
- Tạo nền tảng ổn định cho lọc/sắp xếp tour theo nghiệp vụ.

- Non-goals (what's explicitly out of scope):
- Không triển khai thanh toán hoặc quản lý vé máy bay thực tế.
- Không thêm module quản trị CMS cho tour trong phạm vi feature này.

## User Stories & Use Cases
**How will users interact with the solution?**

- As a khách hàng, I want to xem được danh sách tour quốc tế đầy đủ để chọn chuyến đi phù hợp.
- As a QA/dev, I want to khởi tạo môi trường mới với dữ liệu chuẩn để test lọc/sắp xếp nhất quán.
- As a backend dev, I want to backfill dữ liệu an toàn để không làm mất dữ liệu tour hiện có.

- Key workflows and scenarios:
- Khởi động backend lần đầu: hệ thống seed đủ tour quốc tế theo danh sách chuẩn.
- Khởi động backend với DB cũ: hệ thống update/bổ sung tour quốc tế còn thiếu.
- Khởi tạo MySQL từ `init_database.sql`: có cột và dữ liệu tour quốc tế phù hợp.

- Edge cases to consider:
- Tour đã tồn tại cùng tên nhưng thiếu `type` hoặc `country`.
- Dữ liệu tên không dấu/có dấu khác nhau gây trùng logic.
- DB cũ thiếu cột mới so với mô hình hiện tại.

## Success Criteria
**How will we know when we're done?**

- API `GET /v1/api/tours?category=international` luôn trả về danh sách tour quốc tế có tên tour cụ thể.
- Backfill chạy idempotent: chạy nhiều lần không sinh bản ghi trùng theo tên chuẩn.
- Script `init_database.sql` tạo được bảng `tours` có `type`, `country` và dữ liệu quốc tế mẫu.
- Build backend/frontend pass.

## Constraints & Assumptions
**What limitations do we need to work within?**

- Technical constraints:
- Backend hiện dùng SQLite runtime qua GORM.
- Dự án vẫn duy trì script MySQL trong `init_database.sql` cho môi trường khác.

- Business constraints:
- Phân loại tour theo 3 nhóm chính: domestic/international/service.

- Time/budget constraints:
- Thực hiện trong phạm vi sửa dữ liệu và tài liệu, không thay đổi kiến trúc lớn.

- Assumptions we're making:
- Tên tour là khóa logic đủ ổn định để backfill/upsert ở mức seed.
- Danh sách tour quốc tế mẫu ban đầu do team kỹ thuật định nghĩa.

## Questions & Open Items
**What do we still need to clarify?**

- Bộ danh sách tour quốc tế chuẩn cuối cùng cần bao nhiêu tour?
- Có cần thêm trường tiền tệ chuẩn dạng số (`price_value`) thay cho string giá hiển thị không?
- Có cần quản lý đa ngôn ngữ cho tên/mô tả tour quốc tế trong giai đoạn tiếp theo không?
