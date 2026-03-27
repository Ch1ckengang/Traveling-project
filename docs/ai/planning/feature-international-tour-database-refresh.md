---
phase: planning
title: Project Planning & Task Breakdown
feature: international-tour-database-refresh
description: Kế hoạch chuẩn hóa dữ liệu tour quốc tế trong CSDL
---

# Project Planning & Task Breakdown

## Milestones
**What are the major checkpoints?**

- [x] Milestone 1: Xác định hiện trạng schema và nguồn seed/backfill tour
- [x] Milestone 2: Cập nhật dữ liệu tour quốc tế trong runtime seed/backfill
- [x] Milestone 3: Đồng bộ script SQL khởi tạo với schema và dữ liệu mới
- [x] Milestone 4: Build và smoke test API lọc/sắp xếp theo tour quốc tế

## Task Breakdown
**What specific work needs to be done?**

### Phase 1: Foundation
- [x] Task 1.1: Rà soát model `Tour` và các cột bắt buộc (`type`, `country`)
- [x] Task 1.2: Rà soát `CreateToursIfEmpty` và dữ liệu quốc tế hiện có

### Phase 2: Core Features
- [x] Task 2.1: Mở rộng danh sách tour quốc tế chuẩn (nhiều tên tour rõ ràng)
- [x] Task 2.2: Đảm bảo backfill cập nhật record đã tồn tại theo tên
- [x] Task 2.3: Đồng bộ `main.go` seed với danh sách quốc tế chuẩn
- [x] Task 2.4: Cập nhật `init_database.sql` để có schema + dữ liệu quốc tế

### Phase 3: Integration & Polish
- [x] Task 3.1: Build backend/frontend
- [x] Task 3.2: Smoke test API `category=international`
- [x] Task 3.3: Ghi tài liệu feature requirements/design/planning/implementation/testing

## Dependencies
**What needs to happen in what order?**

- Cần hoàn tất chuẩn hóa dữ liệu ở repository trước khi chỉnh SQL script.
- Cần build sau khi chỉnh code để đảm bảo API không regress.

## Timeline & Estimates
**When will things be done?**

- Phase 1: ~0.5 giờ
- Phase 2: ~1 giờ
- Phase 3: ~0.5 giờ
- Tổng: ~2 giờ

## Risks & Mitigation
**What could go wrong?**

- Technical risks:
- Trùng tên tour gây cập nhật nhầm record.

- Mitigation strategies:
- Chỉ update các trường metadata khi match đúng `name` chuẩn.
- Tăng cường kiểm tra API theo category international sau khi backfill.

## Resources Needed
**What do we need to succeed?**

- Team members and roles:
- 1 backend dev + 1 reviewer QA.

- Tools and services:
- Go toolchain, Vite build, curl smoke test.

- Infrastructure:
- Local SQLite DB và script MySQL cho môi trường cài đặt.

- Documentation/knowledge:
- docs/ai templates và quy ước seed dữ liệu hiện tại.
