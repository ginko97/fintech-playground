-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS transactions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_id   TEXT NOT NULL UNIQUE,
    amount        BIGINT NOT NULL CHECK (amount > 0),
    currency      TEXT NOT NULL DEFAULT 'IDR',
    type          TEXT NOT NULL,
    status        TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    source_wallet TEXT,
    dest_wallet   TEXT,
    metadata      JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    idempotency_key TEXT UNIQUE NOT NULL,
    request_id      TEXT NOT NULL,
    ledger_balance  BIGINT DEFAULT 0,
    version       INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_tx_external_id ON transactions(external_id);
CREATE INDEX IF NOT EXISTS idx_tx_status ON transactions(status);
CREATE INDEX IF NOT EXISTS idx_tx_created ON transactions(created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS transactions;