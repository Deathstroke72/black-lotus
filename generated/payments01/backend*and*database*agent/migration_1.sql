-- file: migrations/000001_initial_schema.up.sql

-- Enable required extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Custom types for payment statuses
CREATE TYPE payment_status AS ENUM (
    'pending',
    'requires_confirmation',
    'processing',
    'succeeded',
    'failed',
    'canceled',
    'refunded',
    'partially_refunded'
);

CREATE TYPE refund_status AS ENUM (
    'pending',
    'processing',
    'succeeded',
    'failed',
    'canceled'
);

CREATE TYPE transaction_type AS ENUM (
    'payment',
    'refund',
    'chargeback',
    'fee',
    'adjustment'
);

CREATE TYPE transaction_status AS ENUM (
    'pending',
    'completed',
    'failed',
    'reversed'
);

CREATE TYPE payment_method_type AS ENUM (
    'card',
    'bank_account',
    'digital_wallet'
);

CREATE TYPE card_brand AS ENUM (
    'visa',
    'mastercard',
    'amex',
    'discover',
    'unknown'
);

-- Payment Methods table (PCI-DSS compliant - no raw card data)
CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    type payment_method_type NOT NULL,
    
    -- Stripe references (we never store actual card data)
    stripe_payment_method_id VARCHAR(255) UNIQUE,
    stripe_customer_id VARCHAR(255),
    
    -- Card metadata (safe to store per PCI-DSS)
    card_brand card_brand,
    card_last_four VARCHAR(4),
    card_exp_month SMALLINT CHECK (card_exp_month >= 1 AND card_exp_month <= 12),
    card_exp_year SMALLINT CHECK (card_exp_year >= 2020),
    card_fingerprint VARCHAR(255), -- For duplicate detection
    
    -- Bank account metadata
    bank_name VARCHAR(255),
    bank_last_four VARCHAR(4),
    
    -- Digital wallet metadata
    wallet_type VARCHAR(50),
    
    billing_name VARCHAR(255),
    billing_email VARCHAR(255),
    billing_address_line1 VARCHAR(255),
    billing_address_line2 VARCHAR(255),
    billing_city VARCHAR(100),
    billing_state VARCHAR(100),
    billing_postal_code VARCHAR(20),
    billing_country VARCHAR(2),
    
    is_default BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    
    -- Indexes defined below
    CONSTRAINT valid_card_data CHECK (
        type != 'card' OR (
            card_brand IS NOT NULL AND 
            card_last_four IS NOT NULL AND 
            card_exp_month IS NOT NULL AND 
            card_exp_year IS NOT NULL
        )
    )
);

-- Payments table
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Idempotency
    idempotency_key VARCHAR(255) NOT NULL,
    
    -- Business references
    customer_id UUID NOT NULL,
    order_id UUID NOT NULL,
    payment_method_id UUID REFERENCES payment_methods(id),
    
    -- Stripe references
    stripe_payment_intent_id VARCHAR(255) UNIQUE,
    stripe_customer_id VARCHAR(255),
    
    -- Amount details
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    amount_refunded BIGINT NOT NULL DEFAULT 0 CHECK (amount_refunded >= 0),
    
    -- Status tracking
    status payment_status NOT NULL DEFAULT 'pending',
    failure_code VARCHAR(100),
    failure_message TEXT,
    
    -- Metadata
    description TEXT,
    statement_descriptor VARCHAR(22),
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,
    
    -- Version for optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT unique_idempotency_key UNIQUE (idempotency_key),
    CONSTRAINT valid_refund_amount CHECK (amount_refunded <= amount)
);

-- Refunds table
CREATE TABLE refunds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- Idempotency
    idempotency_key VARCHAR(255) NOT NULL,
    
    -- References
    payment_id UUID NOT NULL REFERENCES payments(id),
    stripe_refund_id VARCHAR(255) UNIQUE,
    
    -- Amount
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency VARCHAR(3) NOT NULL,
    
    -- Status
    status refund_status NOT NULL DEFAULT 'pending',
    failure_code VARCHAR(100),
    failure_message TEXT,
    
    -- Reason
    reason VARCHAR(50),
    description TEXT,
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    
    -- Version for optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT unique_refund_idempotency UNIQUE (idempotency_key)
);

-- Transactions table (immutable audit log)
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    
    -- References
    payment_id UUID REFERENCES payments(id),
    refund_id UUID REFERENCES refunds(id),
    
    -- Transaction details
    type transaction_type NOT NULL,
    status transaction_status NOT NULL DEFAULT 'pending',
    
    -- Amount (can be negative for refunds)
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL,
    
    -- Net amounts after fees
    net_amount BIGINT,
    fee_amount BIGINT DEFAULT 0,
    
    -- External references
    stripe_transaction_id VARCHAR(255),
    stripe_balance_transaction_id VARCHAR(255),
    
    -- Context
    description TEXT,
    metadata JSONB DEFAULT '{}',
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    
    -- Transactions are immutable - no updated_at
    CONSTRAINT valid_transaction_ref CHECK (
        (type = 'payment' AND payment_id IS NOT NULL) OR
        (type = 'refund' AND refund_id IS NOT NULL) OR
        (type IN ('chargeback', 'fee', 'adjustment'))
    )
);

-- Idempotency keys table for request deduplication
CREATE TABLE idempotency_keys (
    key VARCHAR(255) PRIMARY KEY,
    request_path VARCHAR(255) NOT NULL,
    request_hash VARCHAR(64) NOT NULL, -- SHA-256 of request body
    response_code INTEGER,
    response_body JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW() + INTERVAL '24 hours',
    locked_until TIMESTAMPTZ
);

-- Outbox table for reliable event publishing
CREATE TABLE outbox_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    retry_count INTEGER NOT NULL DEFAULT 0,
    last_error TEXT,
    next_retry_at TIMESTAMPTZ
);

-- Webhook events table for idempotent webhook processing
CREATE TABLE webhook_events (
    id VARCHAR(255) PRIMARY KEY, -- Stripe event ID
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    processed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0
);

-- Indexes for payment_methods
CREATE INDEX idx_payment_methods_customer_id ON payment_methods(customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_methods_stripe_customer ON payment_methods(stripe_customer_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_payment_methods_fingerprint ON payment_methods(card_fingerprint) WHERE card_fingerprint IS NOT NULL;

-- Indexes for payments
CREATE INDEX idx_payments_customer_id ON payments(customer_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(status);
CREATE INDEX idx_payments_created_at ON payments(created_at DESC);
CREATE INDEX idx_payments_stripe_intent ON payments(stripe_payment_intent_id) WHERE stripe_payment_intent_id IS NOT NULL;
CREATE INDEX idx_payments_customer_status ON payments(customer_id, status);

-- Indexes for refunds
CREATE INDEX idx_refunds_payment_id ON refunds(payment_id);
CREATE INDEX idx_refunds_status ON refunds(status);
CREATE INDEX idx_refunds_created_at ON refunds(created_at DESC);
CREATE INDEX idx_refunds_stripe_id ON refunds(stripe_refund_id) WHERE stripe_refund_id IS NOT NULL;

-- Indexes for transactions
CREATE INDEX idx_transactions_payment_id ON transactions(payment_id) WHERE payment_id IS NOT NULL;
CREATE INDEX idx_transactions_refund_id ON transactions(refund_id) WHERE refund_id IS NOT NULL;
CREATE INDEX idx_transactions_type ON transactions(type);
CREATE INDEX idx_transactions_created_at ON transactions(created_at DESC);
CREATE INDEX idx_transactions_status ON transactions(status);

-- Composite index for transaction history queries
CREATE INDEX idx_transactions_history ON transactions(created_at DESC, type, status);

-- Indexes for idempotency_keys
CREATE INDEX idx_idempotency_expires ON idempotency_keys(expires_at);

-- Indexes for outbox_events
CREATE INDEX idx_outbox_unpublished ON outbox_events(created_at) WHERE published_at IS NULL;
CREATE INDEX idx_outbox_retry ON outbox_events(next_retry_at) WHERE published_at IS NULL AND next_retry_at IS NOT NULL;

-- Indexes for webhook_events
CREATE INDEX idx_webhook_unprocessed ON webhook_events(created_at) WHERE processed_at IS NULL;

-- Trigger function for updating updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply triggers
CREATE TRIGGER update_payment_methods_updated_at BEFORE UPDATE ON payment_methods
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_payments_updated_at BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_refunds_updated_at BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to increment version on update (optimistic locking)
CREATE OR REPLACE FUNCTION increment_version()
RETURNS TRIGGER AS $$
BEGIN
    NEW.version = OLD.version + 1;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER increment_payments_version BEFORE UPDATE ON payments
    FOR EACH ROW EXECUTE FUNCTION increment_version();

CREATE TRIGGER increment_refunds_version BEFORE UPDATE ON refunds
    FOR EACH ROW EXECUTE FUNCTION increment_version();