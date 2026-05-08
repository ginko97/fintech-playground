CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- will use UUIDv7
    external_id TEXT NOT NULL UNIQUE,
    amount BIGINT NOT NULL CHECK (amount > 0),
    currency TEXT NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    source_wallet TEXT,
    dest_wallet TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1,
    
    CONSTRAINT uk_external_id UNIQUE (external_id)
);

CREATE INDEX idx_transactions_external_id ON transactions(external_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transactions_created_at ON transactions(created_at);