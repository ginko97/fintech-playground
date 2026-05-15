-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS transactions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key   TEXT NOT NULL UNIQUE,
    request_id        TEXT NOT NULL,
    external_id       TEXT UNIQUE,
    source_wallet_id  TEXT NOT NULL,
    dest_wallet_id    TEXT NOT NULL,
    amount            BIGINT NOT NULL CHECK (amount > 0),
    currency          TEXT NOT NULL DEFAULT 'IDR',
    status            TEXT NOT NULL CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    description       TEXT,
    metadata          JSONB DEFAULT '{}',
    ledger_balance    BIGINT DEFAULT 0,
    version           INTEGER NOT NULL DEFAULT 1,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tx_idempotency_key ON transactions(idempotency_key);
CREATE INDEX IF NOT EXISTS idx_tx_source_wallet   ON transactions(source_wallet_id);
CREATE INDEX IF NOT EXISTS idx_tx_status          ON transactions(status);

-- +goose Down
DROP TABLE IF EXISTS transactions;