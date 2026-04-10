# Frontend Figma Design Spec - Module 1 (Project-Aligned)

## Objective
Create a Figma wireframe and high-fidelity design package for **Module 1 - Customer** of the Traveling Management System.

This version is adjusted to the **current codebase** so UI/UX design is feasible for implementation with minimal backend changes.

## Current Implementation Baseline (from repository)
- Frontend stack: React + Vite.
- Existing routes: home tour browsing/search, booking flow, auth pages, profile.
- Existing backend support:
  - Tour listing and filtering.
  - Booking creation.
  - Booking data includes adult/child/infant counts, booking code, payment status.
- Booking flow is currently implemented as multi-step screens in app logic.

## Scope
- Product: Traveling Tour Management System.
- Module: Module 1 - Customer.
- Primary device: Desktop 1440x900.
- Optional: Tablet 768 width adaptation.
- Style: Corporate, clean, reliable.

## Fixed Design System
Use the design tokens exactly as provided in your prompt:
- Color palette.
- Typography scale.
- 8px spacing grid.
- Component specs (button/input/table/card/badge/modal/header/sidebar).

No additional custom colors except documented exceptions.

## Figma File Structure
Create these pages in Figma:
1. `Design System`
2. `Wireframes`
3. `Module 1 - QLTTCN`
4. `Module 1 - Dat Tour`
5. `Module 1 - Thanh Toan`

## Screen Inventory (11 Screens)

### A. QLTTCN (3)
1. `1.1 GDKhachHang - Profile Overview`
2. `1.2 GDSuaTTCaNhan - Edit Profile`
3. `1.3 ModalXacNhanLuu - Save Confirmation`

### B. Dat Tour (4)
4. `2.1 TimKiemTourPage - Search`
5. `2.2 ChiTietTourPage - Tour Detail`
6. `2.3 PhieuDatTourPage - Booking Draft`
7. `2.4 DatTourThanhCongPage - Booking Success`

### C. Thanh Toan (4)
8. `3.1 DanhSachChoThanhToanPage`
9. `3.2 ChonPhuongThucThanhToanPage`
10. `3.3 DangXuLyThanhToanPage`
11. `3.4 KetQuaThanhToanPage (Success + Failed variants)`

## Project-Aligned Data Notes
When designing fields and labels, align to current backend naming:
- Tour:
  - `name`, `location`, `duration`, `price`, `departure_date`, `remaining_slots`, `itinerary`, `services`
- Booking:
  - `booking_code`, `payment_status`, `adult_count`, `child_count`, `infant_count`, `quantity`, `total_amount`

Recommended display mapping:
- `remaining_slots` -> `So cho con lai`.
- `payment_status: unpaid` -> `Chua thanh toan` badge style.

## Step Flow Mapping (important)
Current frontend behavior should be represented as separate screens:
1. Search tour.
2. View selected tour detail.
3. Create booking draft (customer info + guest counts + summary).
4. Confirm booking.
5. Show success with booking code and payment instruction.

## Wireframe Rules
For each of the 11 screens:
- Build low-fidelity grayscale first.
- Include empty state.
- Include loading state.
- Include error state.
- Keep layout in desktop frame `1440x900`.

## High-Fidelity Rules
After wireframes are approved:
- Apply design system tokens consistently.
- Use Auto Layout on all major containers.
- Add annotations:
  - Component name.
  - Size.
  - Token references (color/text/spacing).
  - Interaction state notes.

## Mandatory Components in Library
Create reusable components with variants:
- Button: Primary, Secondary, Danger, Disabled.
- Input: Default, Focus, Error, Disabled, Readonly.
- Table: Header, odd row, even row, hover row.
- Badge: unpaid, paid, canceled, warning.
- Card: default and elevated.
- Modal: confirmation and alert.
- Stepper: minus/value/plus.
- Payment Option Card: default/selected.
- Timeline item.
- Sidebar menu item: default/active/hover.

## Handoff Checklist for Developer
Before handoff, ensure each screen has:
1. Frame name uses module prefix and screen id.
2. Notes for hover/focus/error/disabled states.
3. Spacing values on key gaps/paddings.
4. Redlines for critical dimensions.
5. Export setup for assets/icons if needed.

## Suggested Starting Screen
`2.1 TimKiemTourPage`

Reason:
- It anchors the booking journey.
- It is already present in current frontend flow.
- It defines patterns reused by 2.2, 2.3, and 2.4.

## Suggested Work Sequence
1. Build `Design System` page.
2. Create all 11 grayscale wireframes.
3. Convert booking flow screens first (2.1 -> 2.4).
4. Convert profile flow (1.1 -> 1.3).
5. Convert payment flow (3.1 -> 3.4).
6. Annotate and publish library.

## Implementation Risk Notes (for planning)
- Payment gateway integration is not fully implemented yet; screens 3.x should be designed as UI-ready templates.
- Sidebar-heavy profile/payment pages may need route-shell update in frontend architecture.
- Tour image is currently placeholder in UI; keep placeholder component in design to match current state.
