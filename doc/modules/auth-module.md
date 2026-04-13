# Auth Module

Updated: 2026-04-13
Module path: `server/internal/auth`

## 1. Purpose
Handles account registration, email OTP verification, login, forgot-password request, and profile update.

## 2. Key Files
- `server/internal/auth/handler.go`
- `server/internal/auth/service.go`
- `server/internal/auth/repository.go`
- `server/domain/auth_dto.go`
- `server/domain/user.go`
- `server/internal/shared/errors.go`

## 3. API Surface
- `POST /v1/api/register`
- `POST /v1/api/login`
- `POST /v1/api/otp/send`
- `POST /v1/api/otp/verify`
- `POST /v1/api/password/forgot`
- `PUT /v1/api/users/:id`

## 4. Current Business Rules
### Registration
- Name must not be empty.
- Email must be valid format.
- Register email must end with `@gmail.com`.
- Password length must be at least 8.
- Duplicate email is rejected.
- Password is hashed with bcrypt.
- User is created with `is_email_verified = false`.
- OTP is generated and sent right after successful registration.

### OTP Send
- Email must be valid format.
- Email must already exist in users table.
- OTP length is 6 digits.
- OTP expiry is 3 minutes.
- OTP record is stored in memory map keyed by email.

### OTP Verify
- Email and code are validated.
- OTP must exist, match exactly, and not expire.
- OTP is removed after successful verification.
- User `is_email_verified` is set to `true`.

### Login
- Email/password input is required.
- Invalid email/password returns generic auth error.
- Login is blocked when `is_email_verified = false`.

### Forgot Password
- Returns neutral response behavior.
- Current implementation is dev-friendly and does not send real reset email.

## 5. Error Mapping (Handler)
- `400`: invalid payload, invalid email format, invalid OTP, email not registered for OTP.
- `401`: invalid credentials.
- `403`: email not verified.
- `409`: duplicate register email.
- `500`: internal processing errors.

## 6. Security Notes
- Passwords are hashed with bcrypt.
- Login error avoids user enumeration by returning generic credential error.
- OTP should be moved from memory to Redis/DB for production reliability.
