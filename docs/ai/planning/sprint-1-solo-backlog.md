# Sprint 1 Solo Backlog - Core + Authentication

## Sprint Goal
Build a stable authentication foundation with consistent API response format, role-aware access control baseline, and complete account lifecycle for Customer users.

## Sprint Duration
- 10 working days (2 weeks)
- Working style: solo, vertical slices, end-to-end every few days

## Execution Snapshot (Updated: 2026-04-18)

### Current Backend Auth Baseline
- Implemented endpoints: register, login, send OTP, verify OTP, forgot password, update user profile.
- Not implemented yet: refresh token, logout, reset password, change password, profile read endpoint.
- Auth handlers are now standardized with one envelope: `success`, `message`, `data`, `meta`, `error`.
- Auth errors are mapped to deterministic error codes via centralized `authErrorInfo` mapping.

### Sprint 1 Ticket Status Board

| Ticket | Status | Notes |
|---|---|---|
| S1-D1-01 | Done | Auth handlers/services/repositories audited; current scope and gaps identified. |
| S1-D1-02 | Done | Auth handlers now return unified envelope: `success`, `message`, `data`, `meta`, `error`. |
| S1-D1-03 | Done | Auth error code mapping is implemented in `authErrorInfo` and used in all auth handlers. |
| S1-D1-04 | Done | `server/.env.example` extended with JWT/OTP/reset/email auth config fields. |

### Day 1 Output: Unified API Response Contract (Draft)

Use this response envelope consistently for all auth endpoints:

```json
{
	"success": true,
	"message": "Human readable summary",
	"data": {},
	"meta": {
		"request_id": "optional-request-id"
	},
	"error": null
}
```

Error response:

```json
{
	"success": false,
	"message": "Validation failed",
	"data": null,
	"meta": {
		"request_id": "optional-request-id"
	},
	"error": {
		"code": "AUTH_INVALID_CREDENTIALS",
		"details": null
	}
}
```

### Day 1 Output: Auth Error Code Catalog (Draft)

| Error Code | HTTP Status | Current Trigger |
|---|---|---|
| AUTH_INVALID_PAYLOAD | 400 | Invalid/missing auth request body |
| AUTH_INVALID_NAME | 400 | Empty name in register flow |
| AUTH_INVALID_EMAIL | 400 | Invalid email format |
| AUTH_WEAK_PASSWORD | 400 | Password shorter than policy |
| AUTH_EMAIL_ALREADY_REGISTERED | 409 | Duplicate email on register |
| AUTH_INVALID_CREDENTIALS | 401 | Wrong email/password |
| AUTH_EMAIL_NOT_VERIFIED | 403 | Login blocked for unverified account |
| AUTH_OTP_INVALID_EMAIL | 400 | Invalid email in OTP flow |
| AUTH_OTP_EMAIL_NOT_REGISTERED | 400 | OTP requested for unknown email |
| AUTH_OTP_INVALID_OR_EXPIRED | 400 | Wrong or expired OTP |
| AUTH_EMAIL_CHECK_FAILED | 500 | Repository error while checking email |
| AUTH_INTERNAL_ERROR | 500 | Fallback for unmapped auth errors |

### Day 1 Exit Criteria Check
- Every auth handler has a target response shape defined: Done.
- No ambiguous error messages without code: Done (for auth handlers).

## Definition of Done (Sprint)
- Auth APIs are functional: register, login, refresh token, logout, forgot password, reset password, change password, profile update.
- API responses follow one consistent success/error envelope.
- Basic role-based guard is in place for protected endpoints.
- Frontend auth flows are integrated and tested against backend.
- Core happy path + critical failures are covered by tests.
- A short runbook exists for local setup and verification.

## Prioritized User Stories
1. As a guest, I can register with email/password and verify my account.
2. As a user, I can log in and receive access/refresh tokens.
3. As a user, I can refresh access token securely.
4. As a user, I can log out and invalidate my refresh token.
5. As a user, I can request password reset and set a new password.
6. As a user, I can change my password while logged in.
7. As a user, I can view/update my profile.
8. As a developer, I get consistent API responses and predictable error codes.

## Day-by-Day Plan

### Day 1 - Baseline and contract hardening

#### Tickets
- S1-D1-01: Audit current auth endpoints and data model.
- S1-D1-02: Define unified API response contract (success, message, data, meta, error).
- S1-D1-03: Create error code catalog for auth domain.
- S1-D1-04: Add/update env config template for auth secrets and token TTLs.

#### Deliverables
- Response contract documented.
- Error code mapping table documented.
- Updated .env.example fields for JWT and email/reset settings.

#### Acceptance Criteria
- Every auth handler has a target response shape defined.
- No ambiguous error messages without code.

---

### Day 2 - JWT foundation (access + refresh)

#### Tickets
- S1-D2-01: Implement token pair generation (access/refresh).
- S1-D2-02: Add refresh token persistence strategy (DB table or secure store).
- S1-D2-03: Implement refresh endpoint with token rotation.
- S1-D2-04: Add middleware to validate access token and extract user claims.

#### Deliverables
- Working login + refresh token flow in backend.
- Middleware usable by protected routes.

#### Acceptance Criteria
- Refresh returns new access token (and optionally rotated refresh token).
- Expired/invalid token returns standardized unauthorized error.

---

### Day 3 - Register and login hardening

#### Tickets
- S1-D3-01: Finalize register validation (email format, password policy).
- S1-D3-02: Finalize login validation and error handling.
- S1-D3-03: Add account status check (active/locked/unverified baseline).
- S1-D3-04: Add rate limit guard for login/register (basic anti-abuse).

#### Deliverables
- Register/login endpoints production-safe baseline.
- Validation errors standardized.

#### Acceptance Criteria
- Duplicate email returns deterministic conflict error.
- Wrong credentials do not leak sensitive detail.

---

### Day 4 - Logout and session controls

#### Tickets
- S1-D4-01: Implement logout endpoint to revoke refresh token.
- S1-D4-02: Implement "logout all sessions" option (if time allows).
- S1-D4-03: Add refresh token expiry and revoked checks.

#### Deliverables
- Session invalidation flow completed.

#### Acceptance Criteria
- Revoked refresh token cannot be reused.
- Logout endpoint is idempotent and safe.

---

### Day 5 - Forgot/reset password flow

#### Tickets
- S1-D5-01: Implement forgot password request endpoint.
- S1-D5-02: Implement reset token generation with short expiry.
- S1-D5-03: Implement reset password endpoint with token verification.
- S1-D5-04: Integrate email sender abstraction (real provider later, local stub now).

#### Deliverables
- End-to-end reset flow works locally.

#### Acceptance Criteria
- Reset token is one-time use.
- Expired/invalid reset token returns standardized error.

---

### Day 6 - Change password + profile APIs

#### Tickets
- S1-D6-01: Implement change password for authenticated user.
- S1-D6-02: Add profile read endpoint.
- S1-D6-03: Add profile update endpoint (name/phone/avatar placeholder).
- S1-D6-04: Add ownership and authorization checks.

#### Deliverables
- Complete account lifecycle APIs.

#### Acceptance Criteria
- Change password requires current password.
- Users cannot update another user profile.

---

### Day 7 - Frontend auth integration

#### Tickets
- S1-D7-01: Align frontend API client with backend response contract.
- S1-D7-02: Implement token storage and refresh strategy in frontend.
- S1-D7-03: Wire login/register/forgot/reset/change-password screens.
- S1-D7-04: Add route guards for protected pages.

#### Deliverables
- Functional auth UX linked to real backend APIs.

#### Acceptance Criteria
- Page refresh preserves session (until token expiry).
- Unauthorized user is redirected correctly.

---

### Day 8 - RBAC baseline and protected route checks

#### Tickets
- S1-D8-01: Add role claim in JWT and user model.
- S1-D8-02: Implement role guard middleware (customer/staff/admin baseline).
- S1-D8-03: Protect at least one non-auth endpoint as RBAC proof.

#### Deliverables
- Reusable authorization foundation for next sprints.

#### Acceptance Criteria
- Forbidden role receives standardized forbidden response.
- Role checks are centralized, not duplicated in handlers.

---

### Day 9 - Testing and reliability day

#### Tickets
- S1-D9-01: Unit tests for auth service logic (success + error paths).
- S1-D9-02: API integration tests for login/refresh/logout/reset.
- S1-D9-03: Regression checklist for frontend auth flows.

#### Deliverables
- Test suite for Sprint 1 critical flows.

#### Acceptance Criteria
- Core auth flow tests pass locally.
- No critical auth regression in manual checklist.

---

### Day 10 - Hardening, docs, release tag

#### Tickets
- S1-D10-01: Security pass (sensitive logs, token leakage, validation gaps).
- S1-D10-02: Refactor duplicated auth code and error mapping.
- S1-D10-03: Write runbook: local setup, env vars, test commands, curl examples.
- S1-D10-04: Prepare sprint demo script and release note.

#### Deliverables
- Sprint 1 releasable baseline.

#### Acceptance Criteria
- A new developer can run and verify auth flows from documentation only.
- Demo script can be executed in under 10 minutes.

## Suggested Ticket Status Columns
- Todo
- In Progress
- Blocked
- In Review (self-review)
- Done

## Risk Register (Sprint 1)
1. Token strategy changes mid-sprint.
- Mitigation: freeze JWT schema by Day 2 end.

2. Email provider delays.
- Mitigation: use local mail stub and interface abstraction first.

3. Scope creep from non-auth features.
- Mitigation: reject any ticket not tied to Sprint Goal.

## Daily Routine (Solo)
1. Start of day (15 minutes): select max 2 tickets.
2. Midday (10 minutes): update progress and blockers.
3. End of day (20 minutes): commit, short note, next-day plan.

## Sprint 1 Exit Checklist
- Register/login/refresh/logout verified.
- Forgot/reset/change password verified.
- Profile read/update verified.
- Response format unified.
- RBAC baseline active.
- Tests green for critical auth flows.
- Runbook updated.
