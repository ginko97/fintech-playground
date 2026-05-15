# Fintech Playground - Transaction Service

**Immutable ledger with idempotency protection.**

## Architecture
```mermaid
flowchart TD
    Client[Client] --> Handler[HTTP Handler]
    Handler --> Service[Transaction Service]
    Service --> Repository[Transaction Repository]
    Repository --> PostgreSQL[(PostgreSQL Ledger)]
    Service -.-> Idempotency[Idempotency Guard]

## Sequence Diagram - Create Transaction

```mermaid
sequenceDiagram
    participant C as Client
    participant H as Handler
    participant S as Service
    participant R as Repository
    participant DB as PostgreSQL

    C->>H: POST /api/v1/transactions
    H->>S: Create(req)
    S->>R: FindByIdempotencyKey(key)
    alt Transaction Exists
        R-->>S: Return existing transaction
    else New Transaction
        S->>R: Create(tx)
        R->>DB: INSERT ... ON CONFLICT DO NOTHING
    end
    S-->>H: Transaction response
    H-->>C: 201 Created