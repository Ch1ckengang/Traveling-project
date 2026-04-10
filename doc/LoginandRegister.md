# User Authentication and Profile Management - Business Flow

## 1. Account Registration

### Business description
Allows new users to create an account and access travel platform features.

### Required input
- Full Name
- Email (valid and unique)
- Password
- Confirm Password
- Acceptance of Terms

### Processing flow
1. User submits registration form on frontend.
2. Client validates required fields, password match, and terms checkbox.
3. Frontend calls `POST /v1/api/register`.
4. Backend validates input and checks duplicate email.
5. Backend hashes password with bcrypt and creates user record.
6. Frontend redirects directly to login screen with success message.
7. User logs in immediately after successful registration.

### Error handling
- Duplicate email -> conflict response.
- Invalid email/password format -> validation error.
- Server issue -> generic registration failure message.

## 2. Login

### Business description
Authenticates an existing user and grants access to protected user features.

### Required input
- Email
- Password
- Optional: Remember login

### Processing flow
1. User submits login form.
2. Frontend calls `POST /v1/api/login`.
3. Backend normalizes email and verifies bcrypt password hash.
4. On success, frontend stores user in `AuthContext`.
5. If "Remember me" is enabled, user data remains in `localStorage`.

### Error handling and security
- Invalid credentials return generic message.
- Backend does not expose whether email exists.

## 3. OTP Verification

### Business description
Temporarily disabled for registration flow while Gmail/email verification is turned off.

### Processing flow
1. OTP endpoints remain available in backend for future use.
2. Frontend currently does not route users to OTP screen after registration.
3. OTP screen is not part of active auth navigation.

### OTP rules
- 6 digits
- 3-minute expiration
- OTP is deleted after successful verify or expiration

## 4. Forgot Password

### Business description
Handles reset request entry point in a safe, non-enumerable manner.

### Processing flow
1. User enters email on forgot-password screen.
2. Frontend calls `POST /v1/api/password/forgot`.
3. Backend validates request and returns neutral message.
4. Frontend redirects to reset-success screen.

### Security behavior
- Response must not reveal whether the email exists.

## 5. Profile Update (Current API)

### Business description
Allows logged-in users to update their profile information.

### API
- `PUT /v1/api/users/:id`

### Processing flow
1. Frontend sends update payload (name/email/password if changed).
2. Backend verifies target user and uniqueness constraints.
3. Password updates are hashed with bcrypt.
4. Backend stores and returns updated user info.

## 6. Frontend Routes (Current)

- `/login`
- `/register`
- `/forgot-password`
- `/reset-success`
- `/profile`

## 7. Backend Auth APIs (Current)

- `POST /v1/api/login`
- `POST /v1/api/register`
- `POST /v1/api/otp/send`
- `POST /v1/api/otp/verify`
- `POST /v1/api/password/forgot`
- `PUT /v1/api/users/:id`

## 8. Notes for Future Implementation

- Frontend password rule is aligned to backend minimum length (8).
- OTP/Gmail verification is temporarily disabled in registration flow.
- OTP is currently in-memory and should be moved to Redis/DB with TTL for production.
- Forgot-password currently returns safe message only; full email token reset flow is pending.
- Social login buttons are UI-only and not connected to OAuth providers yet.
