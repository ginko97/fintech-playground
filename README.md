                          Fintech Transaction Service
────────────────────────────────────────────────────────────
Client
  │
  ▼
HTTP Handler (Gin)
  │
  ▼
Transaction Service (Application Layer)
  │
  ├─► Idempotency Check (SELECT first)
  │
  ▼
Transaction Repository (Interface)
  │
  ▼
PostgreSQL Ledger
  (idempotency_key UNIQUE + version column)

Sequence Flow

Client ──POST /transactions──► Handler ──Create(req)──► Service
                                            │
                                            ▼
                                      FindByIdempotencyKey()
                                            │
                                            ▼
                                      Repository ──SELECT──► PostgreSQL
                                            │
                                      (if exists) ──return existing
                                            │
                                      else ──INSERT (ON CONFLICT)──► DB
                                            │
                                      ◄──────────────────────────────
                                            │
                                      Return Transaction ID
                                            │
Client ◄──────────────────── 201 Created ─────────────────────── Handler