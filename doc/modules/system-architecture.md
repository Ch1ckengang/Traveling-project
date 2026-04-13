# System Architecture Module (Monolithic)

Updated: 2026-04-13

## 1. Architecture Style
The backend uses a modular monolithic architecture with clear package boundaries.

Stack:
- HTTP framework: Gin
- ORM: GORM
- Database: PostgreSQL

## 2. Current Backend Structure

```text
server/
├── cmd/server/main.go
├── database/postgres.go
├── domain/
│   ├── user.go
│   ├── tour.go
│   ├── booking.go
│   ├── auth_dto.go
│   └── booking_dto.go
├── internal/
│   ├── auth/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── tour/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   ├── booking/
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── repository.go
│   └── shared/errors.go
└── go.mod
```

## 3. Layer Responsibilities
- Handler: parse HTTP request, map errors to HTTP status, return JSON.
- Service: business rules, validation, orchestration across repositories.
- Repository: database reads/writes and transactions.
- Domain: entities and DTOs shared across modules.

## 4. Dependency Direction
Required one-way dependency:
- handler -> service -> repository -> database
- domain is referenced by handler/service/repository
- shared errors are referenced by handlers/services

Disallowed pattern:
- repository calling handler
- service importing another module handler

## 5. Runtime Flow
1. `main.go` connects to PostgreSQL and runs AutoMigrate.
2. `main.go` registers routes under `/v1`.
3. Incoming request goes to module handler.
4. Handler calls service.
5. Service applies business rules and uses repository.
6. Repository executes SQL through GORM.
7. Response is returned as module DTO/domain object.

## 6. Route Prefix and Main Endpoints
Base prefix: `/v1`

Main endpoints:
- `POST /api/register`
- `POST /api/login`
- `POST /api/otp/send`
- `POST /api/otp/verify`
- `POST /api/password/forgot`
- `GET /api/tours`
- `GET /api/tours/domestic`
- `GET /api/tours/international`
- `POST /api/bookings`
- `GET /api/users/:id/bookings`
- `PUT /api/users/:id/bookings/:bookingId/cancel`

## 7. Operational Notes
- OTP is currently in-memory (dev-oriented).
- Seed users are marked `is_email_verified = true` for local login convenience.
- New registered users must verify OTP before login.
- Register email domain is restricted to `@gmail.com`.
