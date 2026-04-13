# Core Module Documentation (English)

This folder contains the main backend module documentation for the current monolithic architecture.

## Module Documents
- `system-architecture.md`: architecture, package layout, dependency rules, request flow.
- `auth-module.md`: registration, OTP verification, login, forgot-password behavior.
- `tour-module.md`: tour listing, filtering, sorting, and seed behavior.
- `booking-module.md`: booking creation, validation, cancellation, and slot management.

## Source of Truth
The docs in this folder reflect the current code structure under:
- `server/cmd/server/main.go`
- `server/internal/auth/*`
- `server/internal/tour/*`
- `server/internal/booking/*`
- `server/domain/*`
- `server/database/postgres.go`

## Vietnamese Summary
For a Vietnamese end-to-end summary, see:
- `doc/tong-quan-he-thong.md`
