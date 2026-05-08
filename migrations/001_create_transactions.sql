-- 001_create_transactions.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS transactions (
    id            UUID PRIMARY KEY,
    external_id   TEXT NOT NULL UNIQUE,
    amount        BIGINT NOT NULL CHECK (amount > 0),
    currency      TEXT NOT NULL,
    type          TEXT NOT NULL,
    status        TEXT NOT NULL,
    source_wallet TEXT,
    dest_wallet   TEXT,
    metadata      JSONB,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version       INTEGER NOT NULL DEFAULT 1,

    CONSTRAINT chk_status CHECK (status IN ('pending', 'completed', 'failed', 'refunded'))
);

CREATE INDEX IF NOT EXISTS idx_transactions_external_id ON transactions(external_id);
CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at DESC);