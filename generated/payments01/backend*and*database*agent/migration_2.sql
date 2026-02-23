-- file: migrations/000001_initial_schema.down.sql

-- Drop triggers first
DROP TRIGGER IF EXISTS increment_refunds_version ON refunds;
DROP TRIGGER IF EXISTS increment_payments_version ON payments;
DROP TRIGGER IF EXISTS update_refunds_updated_at ON refunds;
DROP TRIGGER IF EXISTS update_payments_updated_at ON payments;
DROP TRIGGER IF EXISTS update_payment_methods_updated_at ON payment_methods;

-- Drop functions
DROP FUNCTION IF EXISTS increment_version();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in dependency order
DROP TABLE IF EXISTS webhook_events;
DROP TABLE IF EXISTS outbox_events;
DROP TABLE IF EXISTS idempotency_keys;
DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS refunds;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS payment_methods;

-- Drop custom types
DROP TYPE IF EXISTS payment_method_type;
DROP TYPE IF EXISTS card_brand;
DROP TYPE IF EXISTS transaction_status;
DROP TYPE IF EXISTS transaction_type;
DROP TYPE IF EXISTS refund_status;
DROP TYPE IF EXISTS payment_status;