## Architecture

```mermaid
flowchart LR
    C[Client]
    H[HTTP Handler Gin]
    S[Transaction Service]
    G[Idempotency Guard]
    R[Repository]
    DB[(PostgreSQL Ledger)]

    C -->|POST /transactions| H
    H -->|Create(req)| S
    S -->|check| G
    S -->|FindByIdempotencyKey / Create| R
    R -->|INSERT / SELECT| DB
```

## Sequence Diagram

```mermaid
sequenceDiagram
    actor Client
    participant Handler
    participant Service
    participant Repository
    participant PostgreSQL

    Client->>Handler: POST /transactions
    Handler->>Service: Create(req)
    Service->>Repository: FindByIdempotencyKey()

    alt exists
        Repository-->>Service: return existing
    else new
        Service->>Repository: Create(tx)
        Repository->>PostgreSQL: INSERT (ON CONFLICT DO NOTHING)
    end

    Service-->>Handler: Transaction
    Handler-->>Client: 201 Created
```