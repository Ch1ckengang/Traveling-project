-- =============================================
-- Migration: Create payments table and related indexes
-- Version: 001
-- Description: Add payments table for VNPay integration
-- =============================================

-- Create payments table
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id),
    vnpay_transaction_id VARCHAR(100) UNIQUE,
    transaction_reference VARCHAR(50) UNIQUE NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'VND',
    payment_method VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    vnpay_response_code VARCHAR(10),
    vnpay_message TEXT,
    session_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create indexes for performance
CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_transaction_ref ON payments(transaction_reference);
CREATE INDEX idx_payments_vnpay_txn_id ON payments(vnpay_transaction_id);
CREATE INDEX idx_payments_status ON payments(status);

-- Create payment audit logs table
CREATE TABLE payment_audit_logs (
    id SERIAL PRIMARY KEY,
    payment_id INTEGER REFERENCES payments(id),
    booking_id INTEGER REFERENCES bookings(id),
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB,
    user_id INTEGER REFERENCES users(id),
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for audit logs
CREATE INDEX idx_audit_payment_id ON payment_audit_logs(payment_id);
CREATE INDEX idx_audit_event_type ON payment_audit_logs(event_type);
CREATE INDEX idx_audit_timestamp ON payment_audit_logs(timestamp);

-- Add payment tracking columns to bookings table
ALTER TABLE bookings 
ADD COLUMN payment_id INTEGER REFERENCES payments(id),
ADD COLUMN payment_deadline TIMESTAMP;

-- Update payment_status column type if needed (it's already VARCHAR(50))
-- ALTER TABLE bookings ALTER COLUMN payment_status TYPE VARCHAR(20);

-- Grant permissions for the payment tables
GRANT SELECT, INSERT, UPDATE, DELETE ON payments TO postgres;
GRANT SELECT, INSERT, UPDATE, DELETE ON payment_audit_logs TO postgres;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE payments_id_seq TO postgres;
GRANT USAGE, SELECT, UPDATE ON SEQUENCE payment_audit_logs_id_seq TO postgres;