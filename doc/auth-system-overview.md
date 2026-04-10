# Authentication System Overview (Frontend + Backend)

Updated: 2026-04-05
Scope: Full summary of implemented features for the authentication UI redesign and OTP/forgot-password flow.

## 1) Completed Features

### Frontend (React)
- Redesigned authentication UI in flat design with blue primary theme `#185FA5`, minimal and professional style.
- Implemented 5 auth screens:
  - Login
  - Register
  - OTP Verification
  - Forgot Password
  - Reset Success
- Added responsive 2-column layout, card max width 380px, 0.5px border, 12px radius.
- Added blue focus ring for inputs and updated primary/outline button styles.
- Added dark mode support through CSS variables.
- Added auth routes and hid header on auth pages.
- Integrated OTP and forgot-password API calls.

### Backend (Go + Gin + GORM)
- Added auth endpoints:
  - `POST /v1/api/otp/send`
  - `POST /v1/api/otp/verify`
  - `POST /v1/api/password/forgot`
- Added request models for OTP and forgot-password APIs.
- Added in-memory OTP logic:
  - Generate 6-digit OTP
  - Expiration after 3 minutes
  - Remove OTP after successful verification or expiry
- Added safe forgot-password behavior:
  - Do not disclose whether an email exists
  - Return a neutral response message
- Dev-mode note: OTP and reset requests are logged in backend terminal for local testing.

## 2) Updated Files

### Frontend
- `client/src/components/Auth/AuthLayout.jsx`
- `client/src/components/Auth/Login.jsx`
- `client/src/components/Auth/Register.jsx`
- `client/src/components/Auth/OtpVerification.jsx`
- `client/src/components/Auth/ForgotPassword.jsx`
- `client/src/components/Auth/ResetSuccess.jsx`
- `client/src/styles/Auth.css`
- `client/src/App.jsx`
- `client/src/context/AuthContext.jsx` (state optimization + lint fix)
- `client/src/context/ThemeContext.jsx` (lint fix)

### Backend
- `server/main.go`
- `server/controllers/user_controller.go`
- `server/services/user_service.go`
- `server/services/errors.go`
- `server/models/models.go`

## 3) Current Routes and APIs

### Frontend auth routes
- `/login`
- `/register`
- `/otp-verification`
- `/forgot-password`
- `/reset-success`

### Backend auth APIs
- `POST /v1/api/login`
- `POST /v1/api/register`
- `POST /v1/api/otp/send`
- `POST /v1/api/otp/verify`
- `POST /v1/api/password/forgot`
- `PUT /v1/api/users/:id`

## 4) End-to-End Flows (Frontend -> Backend)

### Flow A: Register + OTP
1. User opens `/register`, enters form data, and submits.
2. Frontend calls `POST /v1/api/register`.
3. Backend validates input and creates account (bcrypt password hash).
4. Frontend calls `POST /v1/api/otp/send` with the registered email.
5. Backend generates 6-digit OTP, stores it in-memory for 3 minutes, and logs it in dev mode.
6. Frontend redirects to `/otp-verification` and shows status message.
7. User enters OTP, frontend calls `POST /v1/api/otp/verify`.
8. Backend validates email + code + expiry.
9. On success, frontend redirects to `/login` with success message.

### Flow B: Login
1. User opens `/login` and submits email/password.
2. Frontend calls `POST /v1/api/login`.
3. Backend authenticates credentials.
4. On success, frontend stores user in `AuthContext`.
5. If "Remember me" is enabled, user info remains in `localStorage`.
6. If disabled, `localStorage` is cleared after login.

### Flow C: Forgot Password
1. User opens `/forgot-password` and submits email.
2. Frontend calls `POST /v1/api/password/forgot`.
3. Backend validates email and returns neutral response behavior.
4. Frontend redirects to `/reset-success` and displays server message.

## 5) Business Rules

### 5.1 Frontend rules

#### Login
- Social buttons (Google/Facebook) are UI-only for now (no real OAuth integration yet).
- Supports password show/hide.
- Supports "Remember me".
- Includes "Forgot password" entry.

#### Register
- Client-side validation:
  - Password and confirm password must match.
  - Password minimum length is currently 6 on client.
  - Terms checkbox is required.
- Includes 4-bar password strength indicator (Weak/Medium/Strong).
- On successful register, OTP verification is required.

Important: Backend requires minimum 8 characters. Frontend currently validates 6, so 6-7 character passwords can pass client checks but fail on backend.

#### OTP Verification
- 6 separate input boxes, auto-focus forward, backspace focus backward.
- Supports paste of 6-digit code.
- UI countdown in `MM:SS` format (3 minutes).
- "Resend" button calls OTP send endpoint again.

#### Forgot Password + Reset Success
- Simple email form calls forgot-password API.
- Success screen shows confirmation message and login CTA.

### 5.2 Backend rules

#### Login (`services.Login`)
- Normalize email (trim + lowercase).
- Validate login payload.
- Return generic auth failure message to reduce account enumeration risk.

#### Register (`services.Register`)
- Validate:
  - Non-empty name
  - Valid email format
  - Password minimum 8 characters
- Check duplicate email.
- Hash password with bcrypt.
- Create user in database.

#### OTP send/verify
- `SendOTPForEmail`:
  - Validate email format.
  - Generate random 6-digit code.
  - Store in-memory map `otpStore[email] = {code, expiresAt}`.
- `VerifyOTPForEmail`:
  - Validate email and code length.
  - Check OTP existence, expiration, and exact code match.
  - Delete OTP after successful verification.

#### Forgot password (`RequestPasswordReset`)
- Validate email format.
- Check user existence only for internal flow.
- Never expose whether email exists.
- Current implementation does not send real email yet (dev logging only).

## 6) Applied UI/CSS Standards
- Primary color: `#185FA5`, hover: `#0C447C`.
- Border width for card/input/button: `0.5px`.
- Input padding: `8px 10px`.
- OTP box size: `44x48` (mobile: `40x44`).
- Auth module font: `system-ui, sans-serif`.
- Dark mode via `:root` and `:root[data-theme='dark']` variables.

## 7) Validation and Build Status
- Frontend:
  - `npm run lint`: pass
  - `npm run build`: pass
- Backend:
  - `go build ./...`: pass

## 8) Current Limitations
- OTP is in-memory and resets when server restarts.
- No persistent OTP storage in Redis/DB yet.
- No real email provider integration for OTP/reset.
- No complete password reset form/token link flow yet.
- Social login is UI-only; OAuth is not implemented.

## 9) Recommended Next Steps
1. Align frontend password rule with backend minimum (8+).
2. Integrate real email provider for OTP and reset messages.
3. Store OTP with TTL in Redis/DB instead of in-memory map.
4. Implement full reset-password flow:
   - Generate reset token
   - Send reset link email
   - Add new-password form
   - Add password update confirmation endpoint
5. Add rate limiting for OTP send/verify and forgot-password endpoints.
