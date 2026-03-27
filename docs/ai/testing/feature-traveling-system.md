---
phase: testing
title: Testing Strategy
feature: traveling-system
description: Testing strategy for the online tour booking system
---

# Testing Strategy — Traveling System

## Test Coverage Goals

- **Unit tests**: 100% of backend handlers (auth, tour, booking, invoice)
- **Integration tests**: All critical paths — register → book tour → create invoice
- **End-to-end**: Key customer and staff user flows
- **Manual testing**: UI/UX, responsive design, interface edge cases

---

## Unit Tests

### handlers/auth.go

- [ ] **Register success** — new email, valid password → HTTP 200, user created
- [ ] **Register with duplicate email** → HTTP 409, message "Email already registered"
- [ ] **Register missing required fields** → HTTP 400
- [ ] **Login success** — correct email + password → HTTP 200 + JWT token
- [ ] **Login with wrong password** → HTTP 401
- [ ] **Login with non-existent email** → HTTP 401
- [ ] **Password is hashed** — verify DB stores bcrypt hash, not plaintext

### handlers/tour.go

- [ ] **Get tour list** → HTTP 200, returns array
- [ ] **Get existing tour detail** → HTTP 200, correct fields
- [ ] **Get non-existent tour** → HTTP 404
- [ ] **Create tour — role staff** → HTTP 201
- [ ] **Create tour — role customer** → HTTP 403
- [ ] **Update tour** → HTTP 200, data updated
- [ ] **Delete tour — role admin** → HTTP 200
- [ ] **Delete tour — role staff** → HTTP 403

### handlers/booking.go

- [ ] **Book tour successfully** — seats available → HTTP 201, booking created
- [ ] **Book tour when fully booked** — total guests exceed `SLKhachMax` → HTTP 400 "Tour is fully booked"
- [ ] **Book non-existent tour** → HTTP 404
- [ ] **Book tour while not authenticated** → HTTP 401
- [ ] **Customer views own bookings** — only sees their own → HTTP 200
- [ ] **Staff views all bookings** → HTTP 200, returns all records

### handlers/invoice.go

- [ ] **Create invoice successfully** — role staff, PDTour exists → HTTP 201
- [ ] **Create invoice — role customer** → HTTP 403
- [ ] **Create invoice for PDTour that already has one** → HTTP 409
- [ ] **Invoice total calculation** — (adults + children * 0.7) * tourCost

### middleware/auth.go

- [ ] **Request with no token** → HTTP 401
- [ ] **Expired token** → HTTP 401 "Invalid token"
- [ ] **Valid token** → request passes, `userID` set in context
- [ ] **Token with invalid signature** → HTTP 401

---

## Integration Tests

- [ ] **Register → login → get profile** — 3 sequential API calls
- [ ] **Staff creates tour → creates schedule → customer books tour** — full business flow
- [ ] **Customer books tour → staff creates invoice** — data correctly linked
- [ ] **Transaction rollback test** — if invoice creation fails mid-way, PDTour status must not change
- [ ] **Concurrent booking test** — 2 customers book the last seat simultaneously → only 1 succeeds

---

## End-to-End Tests

### Customer Flow

- [ ] **E2E-C-01**: Visit homepage → browse tour list → click a tour → view detail
- [ ] **E2E-C-02**: Register account → login → book tour → view booking history
- [ ] **E2E-C-03**: Update personal information (name, email) → verify displayed correctly
- [ ] **E2E-C-04**: Search tours by destination → correct results returned

### Staff Flow

- [ ] **E2E-S-01**: Staff login → go to Dashboard → view booking list
- [ ] **E2E-S-02**: Create new tour → create schedule → verify tour appears on customer page
- [ ] **E2E-S-03**: View customer booking → create invoice → invoice appears in list

---

## Test Data

**Seed data for test environment:**

```sql
-- Sample members
INSERT INTO tblThanhVien VALUES
(1, 'customer01', '<bcrypt_hash>', '1990-05-15', 'customer01@test.com', 'customer'),
(2, 'staff01',    '<bcrypt_hash>', '1985-03-20', 'staff01@test.com',    'staff'),
(3, 'admin01',    '<bcrypt_hash>', '1980-01-01', 'admin@test.com',      'admin');

-- Sample tours
INSERT INTO tblTour VALUES
(1, 'T001', 'Da Nang Tour',  '3 days 2 nights', 'Airplane', 20, 'Tour description', 2000000),
(2, 'T002', 'Phu Quoc Tour', '5 days 4 nights', 'Airplane', 15, 'Tour description', 5000000);

-- Sample schedules
INSERT INTO tblLichTour VALUES
(1, '2026-04-10', 1),
(2, '2026-04-20', 1),
(3, '2026-05-01', 2);
```

**Mock DB for unit tests (Go):**

```go
func mockDB(t *testing.T) *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&ThanhVien{}, &Tour{}, &LichTour{}, &PDTour{}, &HoaDon{})
    return db
}
```

---

## Test Reporting & Coverage

**Run backend tests:**
```bash
cd server
go test ./... -v -cover
go test ./handlers/... -coverprofile=coverage.out
go tool cover -html=coverage.out    # View HTML report
```

**Coverage targets:**
- `handlers/` : ≥ 90%
- `middleware/` : 100%
- `models/` : no unit tests needed (structs only)

---

## Manual Testing

**UI/UX Checklist:**
- [ ] Registration form shows clear validation errors
- [ ] Error messages do not expose technical details (stack trace, SQL errors)
- [ ] Loading spinner displays while awaiting API responses
- [ ] Pages are responsive on mobile (< 768px) and desktop
- [ ] "Book Tour" button is disabled when tour is fully booked
- [ ] Redirects to `/login` when JWT token expires
- [ ] Booking form validates guest count > 0 before submission

**Security Manual Checks:**
- [ ] Try `POST /api/tours` without a token → must receive 401
- [ ] Try using a customer token to call create tour API → must receive 403
- [ ] Inspect password field in API response → must never appear

---

## Performance Testing

- [ ] **Load test** `GET /api/tours` with 50 concurrent requests → response < 500ms
- [ ] **Load test** `POST /api/pdtour` with 10 concurrent requests on a schedule with 5 seats → correct number of bookings created, never exceeding capacity
- Tool: `wrk` or `hey` (Go HTTP load testing)

```bash
# Example with hey
hey -n 100 -c 10 http://localhost:8080/api/tours
```

---

## Bug Tracking

**Bug priority levels:**

| Level | Description | SLA |
|-------|-------------|-----|
| 🔴 Critical | Data loss, security vulnerability, cannot book tour | Fix immediately |
| 🟡 High | Core feature broken, UI shows wrong data | Fix within the day |
| 🟢 Low | Minor UI glitch, wrong text, no functional impact | Fix in next sprint |

**Bug reporting process:**
1. Describe steps to reproduce
2. Expected vs Actual behavior
3. Screenshot / logs
4. Create an issue in the Git repository
