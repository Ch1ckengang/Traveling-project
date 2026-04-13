# Tour Module

Updated: 2026-04-13
Module path: `server/internal/tour`

## 1. Purpose
Provides public tour listing APIs with category filtering, search filtering, and sorting.

## 2. Key Files
- `server/internal/tour/handler.go`
- `server/internal/tour/service.go`
- `server/internal/tour/repository.go`
- `server/domain/tour.go`

## 3. API Surface
- `GET /v1/api/tours`
- `GET /v1/api/tours/domestic`
- `GET /v1/api/tours/international`

## 4. Filter Model
`TourFilter` fields:
- `Category`
- `City`
- `Duration`
- `Price`
- `Sort`

Input comes from query params in handlers.

## 5. Current Business Logic
- Ensures seed tours exist before listing (`CreateToursIfEmpty`).
- Supports category-level filtering in repository (`domestic`, `international`, `service`).
- Supports service-level filters:
  - city/name contains match
  - duration buckets: `short`, `medium`, `long`
  - price buckets: `low`, `mid`, `high`
- Supports sorting:
  - `price_asc`, `price_desc`
  - `duration_asc`, `duration_desc`
  - `name_asc`, `name_desc`
  - `latest`

## 6. Data and Slot Rules
- `remaining_slots` defaults to 30 when missing/invalid.
- `DecreaseTourRemainingSlots` runs in DB transaction.
- International route excludes Vietnam by normalized country rules.

## 7. Operational Notes
- Seed function updates existing tours by name for consistency.
- Tour list endpoint is used by homepage and booking entry flow.
