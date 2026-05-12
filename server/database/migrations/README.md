# Database Migrations

This directory contains database migration files for the Travel Booking System.

## Migration Naming Convention

Migration files follow the pattern: `{version}_{description}.{direction}.sql`

- `version`: Sequential number (001, 002, etc.)
- `description`: Brief description of the migration (snake_case)
- `direction`: Either `up` (apply migration) or `down` (rollback migration)

## Current Migrations

### 001_create_payments_table
- **Purpose**: Add payments table and audit logging for VNPay integration
- **Tables Created**:
  - `payments`: Store payment transaction data
  - `payment_audit_logs`: Store audit trail for payment events
- **Tables Modified**:
  - `bookings`: Add payment_id and payment_deadline columns
- **Indexes Created**:
  - Performance indexes on payments table
  - Audit log indexes for efficient querying

## Running Migrations

Currently, migrations need to be run manually. To apply the payments migration:

```sql
-- Connect to your travel_db database and run:
\i server/database/migrations/001_create_payments_table.up.sql
```

To rollback:

```sql
-- Connect to your travel_db database and run:
\i server/database/migrations/001_create_payments_table.down.sql
```

## Future Improvements

Consider implementing an automated migration system using tools like:
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [goose](https://github.com/pressly/goose)
- Custom migration runner in Go

## Migration Guidelines

1. Always create both `up` and `down` migration files
2. Test migrations on a copy of production data
3. Include proper indexes for performance
4. Grant appropriate permissions
5. Use `IF EXISTS` and `IF NOT EXISTS` where appropriate
6. Document breaking changes in migration comments