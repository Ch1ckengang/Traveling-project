# Next Session Handover - Authentication Scope

Updated: 2026-04-05
Purpose: Short reference document to quickly recover context and avoid requirement mismatch in future sessions.

## 1) Confirmed Requirements
- Authentication UI follows flat design with blue theme `#185FA5`, minimal and professional.
- Five screens are required:
  - Login
  - Register
  - OTP Verification
  - Forgot Password
  - Reset Success
- Auth layout: responsive 2 columns, max card width 380px, 0.5px border, 12px radius.
- Dark mode support through CSS variables.
- No external UI framework; implemented with existing React components and custom CSS.

## 2) Completed Scope
- Frontend includes all required auth routes and screen UI.
- Frontend is integrated with real APIs for:
  - Login
  - Register
  - Send OTP
  - Verify OTP
  - Forgot password
- Backend includes OTP and forgot-password endpoints.
- Frontend lint/build pass. Backend compile pass.

## 3) Current Auth APIs (Source of Truth)
- `POST /v1/api/login`
- `POST /v1/api/register`
- `POST /v1/api/otp/send`
- `POST /v1/api/otp/verify`
- `POST /v1/api/password/forgot`
- `PUT /v1/api/users/:id`

## 4) Business Rules That Must Be Preserved
- Login must not reveal whether an email exists.
- Register must hash password with bcrypt on backend.
- OTP rules:
  - 6 digits
  - 3-minute expiration
  - remove OTP after successful verification
- Forgot-password response must be neutral to avoid account enumeration.

## 5) Known Mismatch to Avoid Confusion
- Frontend register screen currently validates password length >= 6.
- Backend requires password length >= 8.
- Result: 6-7 character passwords may pass frontend but fail backend.
- Next priority: align frontend rule to 8.

## 6) Current Technical Limitations
- OTP is stored in-memory (lost after server restart).
- No real email delivery service for OTP/reset yet.
- Social login is UI-only (no OAuth flow).
- No full reset-password token-link + new-password form yet.

## 7) Quick Test Checklist
1. Run backend (`server`) and frontend (`client`).
2. Register a new account -> should redirect to OTP screen.
3. Read OTP from backend log (`[OTP] Email=... Code=...`).
4. Verify OTP -> should redirect back to login.
5. Submit forgot-password -> success screen should show neutral message.

## 8) Files to Read Before Continuing
- `doc/auth-system-overview.md` (detailed version)
- `doc/next-session-handover.md` (quick version)
- `client/src/components/Auth/*`
- `server/controllers/user_controller.go`
- `server/services/user_service.go`
