# Booking Module

Updated: 2026-04-13
Module path: `server/internal/booking`

## 1. Purpose
Handles booking creation, user booking history retrieval, and booking cancellation.

## 2. Key Files
- `server/internal/booking/handler.go`
- `server/internal/booking/service.go`
- `server/internal/booking/repository.go`
- `server/domain/booking.go`
- `server/domain/booking_dto.go`

## 3. API Surface
- `POST /v1/api/bookings`
- `GET /v1/api/users/:id/bookings`
- `PUT /v1/api/users/:id/bookings/:bookingId/cancel`

## 4. Create Booking Rules
Request requirements:
- `tour_id` is required.
- `full_name`, `phone`, `email`, `travel_date` are required.
- At least 1 adult.
- Child/infant counts cannot be negative.
- Total quantity must be greater than 0.
- Travel date must be valid `YYYY-MM-DD` and not in the past.

Business steps:
1. Normalize request values.
2. Validate payload.
3. Read tour by ID.
4. Check remaining slots.
5. Calculate total amount (children priced at 75% of base price).
6. Create booking record.
7. Decrease tour remaining slots in transaction.

## 5. Booking Retrieval
- `GetBookingsByUserID` validates numeric user ID.
- Returns latest bookings first (`created_at DESC`).
- Preloads `Tour` relation.

## 6. Cancel Booking Rules
- Validates user ID and booking ID.
- Booking must belong to the user.
- Already cancelled booking is rejected.
- Cancellation updates:
  - booking status -> `cancelled`
  - payment status -> `cancelled`
  - restores reserved slots back to tour
- All updates are executed in one DB transaction.

## 7. Error Behavior
- `400`: invalid input or invalid IDs.
- `404`: tour/booking not found.
- `409`: booking already cancelled.
- `500`: persistence/transaction failure.
