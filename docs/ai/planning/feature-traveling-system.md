---
phase: planning
title: Project Planning & Task Breakdown
feature: traveling-system
description: Implementation plan for the online tour booking system
---

# Project Planning & Task Breakdown — Traveling System

## Milestones

- [x] **M0 — Foundation**: Go Gin + GORM + MySQL backend, React + Vite frontend, basic Auth (Login/Register)
- [ ] **M1 — Security & Advanced Auth**: JWT, bcrypt, role-based middleware
- [ ] **M2 — Tour Management**: CRUD tours, tour schedules, destinations, services
- [ ] **M3 — Booking Flow**: Book a tour, view bookings
- [ ] **M4 — Invoice & Dashboard**: Invoices, staff dashboard
- [ ] **M5 — Polish & Testing**: UI/UX improvements, testing, bug fixes

---

## Task Breakdown

### Phase 1 — Security & Advanced Auth (M1)

#### Backend
- [ ] **1.1** Add `golang-jwt/jwt` and `golang.org/x/crypto/bcrypt` dependencies to `go.mod`
- [ ] **1.2** Update `models.go`: add `Role` field (`customer`/`staff`/`admin`) to `ThanhVien`
- [ ] **1.3** Migrate `Password` to bcrypt hash (one-time migration script)
- [ ] **1.4** Create `middleware/auth.go`: verify JWT token from `Authorization: Bearer <token>` header
- [ ] **1.5** Create `middleware/role.go`: enforce role-based access control on protected routes
- [ ] **1.6** Update `POST /api/login` — return JWT token instead of just the user object
- [ ] **1.7** Refactor `main.go` — extract handlers into the `handlers/` directory

#### Frontend
- [ ] **1.8** Update `AuthContext.jsx` — store JWT token in `localStorage`, add `logout` function
- [ ] **1.9** Create `axiosInstance.js` — automatically attach `Authorization` header to all requests
- [ ] **1.10** Add Protected Route — redirect to `/login` if not authenticated

---

### Phase 2 — Tour Management (M2)

#### Backend — Database Models
- [ ] **2.1** Update `models.go` — add all structs: `KhachHang`, `NhanVien`, full `Tour`, `LichTour`, `DiaDiem`, `TourDiaDiem`, `DichVu`, `DichVuDiaDiem`
- [ ] **2.2** Update `database.go` — AutoMigrate all new models
- [ ] **2.3** Update seed data — add sample tours with full field coverage

#### Backend — Tour APIs
- [ ] **2.4** Create `handlers/tour.go`:
  - `GET /api/tours` — list tours (with filter by location, price)
  - `GET /api/tours/:id` — tour detail + schedules
  - `POST /api/tours` — create tour (role: staff/admin)
  - `PUT /api/tours/:id` — update tour (role: staff/admin)
  - `DELETE /api/tours/:id` — delete tour (role: admin)
- [ ] **2.5** Create `handlers/lichtour.go`:
  - `GET /api/tours/:id/lich` — departure schedules for a tour
  - `POST /api/lichtour` — create new schedule (role: staff)
  - `DELETE /api/lichtour/:id` — delete schedule (role: staff)
- [ ] **2.6** Create `handlers/diadiem.go`:
  - `GET /api/diadiem` — list destinations
  - `POST /api/diadiem` — add destination
  - `POST /api/tour-diadiem` — assign destination to tour
- [ ] **2.7** Create `handlers/dichvu.go`:
  - `GET /api/dichvu` — list services
  - `POST /api/dichvu` — add new service

#### Frontend — Tour UI
- [ ] **2.8** Create `components/Tour/TourCard.jsx` — card showing name, destination, price, duration
- [ ] **2.9** Create `components/Tour/TourList.jsx` — tour grid with skeleton loading
- [ ] **2.10** Create `components/Tour/TourDetail.jsx` — tour detail page: description, destinations, schedules, services
- [ ] **2.11** Update `SearchBar.jsx` — call API to filter tours by destination and price range
- [ ] **2.12** Create `pages/TourManagePage.jsx` — CRUD tour page for staff (table + form)

---

### Phase 3 — Booking Flow (M3)

#### Backend
- [ ] **3.1** Create `handlers/booking.go`:
  - `POST /api/pdtour` — create booking (role: customer)
    - Check remaining seats for the selected schedule
    - Auto-create `KhachHang` record if user doesn't have one yet
  - `GET /api/pdtour/my` — current customer's bookings
  - `GET /api/pdtour` — all bookings (role: staff)
  - `PUT /api/pdtour/:id/cancel` — cancel a booking

#### Frontend
- [ ] **3.2** Create `components/Booking/BookingForm.jsx` — form to select schedule, enter guest count, confirm booking
- [ ] **3.3** Create `components/Booking/BookingHistory.jsx` — list of customer's bookings
- [ ] **3.4** Add "Book Tour" button to `TourDetail.jsx`
- [ ] **3.5** Add "Booking History" tab to `Profile.jsx`

---

### Phase 4 — Invoice & Dashboard (M4)

#### Backend
- [ ] **4.1** Create `handlers/invoice.go`:
  - `GET /api/hoadon` — list all invoices (role: staff)
  - `POST /api/hoadon` — create invoice from booking (role: staff)
    - Calculate total = (numAdults + numChildren * 0.7) * tourCost
  - `GET /api/hoadon/:id` — invoice detail

#### Frontend
- [ ] **4.2** Create `pages/DashboardPage.jsx` — staff dashboard with summary statistics
- [ ] **4.3** Create `components/Invoice/InvoiceList.jsx` — invoice list
- [ ] **4.4** Create `components/Invoice/InvoiceDetail.jsx` — invoice detail with print/PDF button

---

### Phase 5 — Polish & Testing (M5)

- [ ] **5.1** Add loading states and error boundaries throughout the app
- [ ] **5.2** Responsive CSS for mobile (< 768px)
- [ ] **5.3** Unit tests for backend handlers using Go testing package
- [ ] **5.4** End-to-end testing: register → book tour → create invoice flow
- [ ] **5.5** Write project setup and run guide (`README.md`)

---

## Dependencies

```
M1 (Auth) ═► M2 (Tour) ═► M3 (Booking) ═► M4 (Invoice)
                                               │
                                          M5 (Testing)
```

- **M1 must complete before M2**: Tour APIs require JWT middleware
- **M2 must complete before M3**: Booking requires Tour and LichTour data
- **M3 must complete before M4**: Invoice creation requires PDTour data
- **M5 can run in parallel with M3 and M4**

---

## Timeline & Estimates

| Milestone | Estimate | Notes |
|-----------|---------|-------|
| M1 — Security & Auth | 2–3 days | Critical refactor, proceed carefully |
| M2 — Tour Management | 3–4 days | Many CRUDs, sufficient seed data needed |
| M3 — Booking Flow | 2–3 days | Most complex business logic |
| M4 — Invoice & Dashboard | 2 days | Relatively straightforward |
| M5 — Polish & Testing | 2 days | Can run in parallel with M3/M4 |
| **Total** | **~12–14 days** | |

---

## Risks & Mitigation

| Risk | Level | Mitigation |
|------|-------|------------|
| Migrating passwords from plaintext to bcrypt may corrupt test data | 🟡 Medium | Run migration in a transaction, backup DB first |
| Invoice pricing logic (adults/children) not yet confirmed | 🟡 Medium | Confirm Q-04 before implementing M4 |
| Frontend out of sync with refactored APIs | 🟡 Medium | Agree on API contract first, use Postman for testing |
| GORM AutoMigrate not suitable for production schema changes | 🟢 Low | Use manual migration scripts for production deployments |

---

## Resources Needed

- **Backend Developer**: Go, GORM, JWT, bcrypt
- **Frontend Developer**: React, Axios, CSS modules
- **Database**: MySQL 8.0 installed with CREATE/DROP privileges
- **Tools**: Go 1.21+, Node.js 18+, Postman (API testing), MySQL Workbench (data inspection)
