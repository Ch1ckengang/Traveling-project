-- =============================================
-- Migration: Drop payments table and related indexes
-- Version: 001 (DOWN)
-- Description: Rollback payments table creation for VNPay integration
-- =============================================

-- Remove payment tracking columns from bookings table
ALTER TABLE bookings DROP COLUMN IF EXISTS payment_id;
ALTER TABLE bookings DROP COLUMN IF EXISTS payment_deadline;

-- Drop indexes for audit logs
DROP INDEX IF EXISTS idx_audit_timestamp;
DROP INDEX IF EXISTS idx_audit_event_type;
DROP INDEX IF EXISTS idx_audit_payment_id;

-- Drop payment audit logs table
DROP TABLE IF EXISTS payment_audit_logs;

-- Drop indexes for payments
DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_payments_vnpay_txn_id;
DROP INDEX IF EXISTS idx_payments_transaction_ref;
DROP INDEX IF EXISTS idx_payments_booking_id;

-- Drop payments table
DROP TABLE IF EXISTS payments;