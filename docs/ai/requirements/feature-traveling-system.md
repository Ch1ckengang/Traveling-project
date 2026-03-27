---
phase: requirements
title: Requirements & Problem Understanding
feature: traveling-system
description: Full requirements analysis for the online tour booking system
---

# Requirements & Problem Understanding — Traveling System

## Problem Statement

**Problem being solved:**

- Customers struggle to search, compare, and book tours quickly and transparently
- Travel companies lack a centralized tool for managing tours, schedules, bookings, and invoices
- Manual booking processes lead to errors and high labor costs

**Who is affected:**
- Individual customers looking to book domestic travel tours
- Travel company staff (consultants, invoicing, tour management)
- System administrators

**Current state:**
- Backend in Go (Gin + GORM + MySQL) with basic APIs: login, register, update profile, list tours
- React (Vite) frontend has login, register, profile, and tour search screens
- No tour booking, tour schedule management, or invoice features yet

---

## Goals & Objectives

**Primary goals (In-scope):**
1. Member account management — registration, login, profile editing
2. Tour information management — create, update, delete, search tours
3. Tour schedule management — create and track departure schedules per tour
4. Tour booking — customers select a tour, choose a schedule, and submit a booking
5. Invoice management — staff confirm and issue invoices
6. Destination & service management — assign destinations and services to tours

**Secondary goals:**
- Display rich tour information with images, descriptions, and pricing
- Revenue and guest count statistics

**Non-goals (Out of scope):**
- Online payment gateway integration (VNPay, Momo) — future phase
- Mobile app — Web only
- Multi-language support — Vietnamese only
- Tour ratings / reviews

---

## User Stories & Use Cases

### Customer (tblKhachHang)

| # | User Story | Priority |
|---|------------|----------|
| C-01 | As a customer, I want to **register an account** so that I can book tours | 🔴 High |
| C-02 | As a customer, I want to **log in** so that I can access the system | 🔴 High |
| C-03 | As a customer, I want to **browse the tour list** so that I can find a suitable tour | 🔴 High |
| C-04 | As a customer, I want to **search tours by destination** so that I can narrow my options | 🟡 Medium |
| C-05 | As a customer, I want to **view tour details** (images, description, schedule, price) to decide whether to book | 🔴 High |
| C-06 | As a customer, I want to **book a tour** (select schedule, number of adults and children) to confirm my trip | 🔴 High |
| C-07 | As a customer, I want to **view my booking history** to track my booked trips | 🟡 Medium |
| C-08 | As a customer, I want to **edit my personal profile** to keep my information up to date | 🟡 Medium |

### Staff (tblNhanVien)

| # | User Story | Priority |
|---|------------|----------|
| S-01 | As a staff member, I want to **manage tour information** (CRUD) to maintain the tour catalog | 🔴 High |
| S-02 | As a staff member, I want to **create tour schedules** (departure/return dates) for customers to choose from | 🔴 High |
| S-03 | As a staff member, I want to **view all booking orders** to process them | 🔴 High |
| S-04 | As a staff member, I want to **create invoices** for confirmed bookings | 🔴 High |
| S-05 | As a staff member, I want to **manage destinations** to assign them to tours | 🟡 Medium |
| S-06 | As a staff member, I want to **manage services** (hotel, dining) per tour stop | 🟢 Low |

### Administrator

| # | User Story | Priority |
|---|------------|----------|
| A-01 | As an admin, I want to **manage member accounts** to control access rights | 🟡 Medium |
| A-02 | As an admin, I want to **view summary statistics** (revenue, tours, guests) | 🟢 Low |

---

## Success Criteria

| Criterion | Measurement |
|-----------|-------------|
| Customer successfully books a tour | Booking is created and saved in DB, visible to staff |
| Staff successfully creates an invoice | Invoice correctly linked to the booking and member |
| Tour search works correctly | Returns accurate results in < 1 second |
| Account authentication is secure | Passwords are hashed and never exposed in responses |
| UI displays correctly | Responsive on both desktop and mobile |
| No data loss | DB transaction rolls back on error |

---

## Constraints & Assumptions

**Technical constraints:**
- Backend: Go 1.21+ with Gin framework and GORM ORM
- Database: MySQL 8.0 with utf8mb4 charset
- Frontend: React 18 + Vite, no TypeScript
- No JWT authentication middleware yet — needs to be added
- Passwords currently stored as plaintext — must migrate to bcrypt

**Business constraints:**
- Each tour schedule has a maximum guest capacity (`SLKhachMax` in `tblTour`)
- Invoices can only be created by staff (not by customers)
- A booking must include at least 1 adult guest

**Assumptions:**
- The system runs internally; no CDN or load balancer required
- Tour images are stored locally or via external URLs
- Timezone: UTC+7 (Vietnam)

---

## Questions & Open Items

| # | Question | Status |
|---|----------|--------|
| Q-01 | Passwords are currently plaintext — when to migrate to bcrypt? | ⏳ Pending |
| Q-02 | JWT token or session-based authentication? | ⏳ Pending |
| Q-03 | Booking status flow: pending → confirmed → cancelled? | ⏳ Pending |
| Q-04 | Does tour pricing differentiate between adults and children? | ⏳ Pending |
| Q-05 | Tour image upload: store on server or use cloud storage? | ⏳ Pending |
| Q-06 | Do staff and customers share the same login endpoint? | ⏳ Pending |
