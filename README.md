## Architecture

```mermaid
flowchart LR
    Client["Client"] 
    Handler["HTTP Handler (Gin)"]
    Service["Transaction Service"]
    Guard["Idempotency Guard + State Machine"]
    Repo["Repository"]
    DB[("PostgreSQL Ledger")]
    Redis[("Redis")]

    Client -->|"POST /transactions"| Handler
    Handler -->|"Create(req)"| Service
    Service -->|"check"| Guard
    Service -->|"Find/Create"| Repo
    Repo -->|"INSERT / SELECT"| DB
    Service -.->|"Async Processing"| Redis
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
    Note over Service,Redis: Worker Pool processes async
```

## Key Production Features
Idempotent Transaction API
Finite State Machine with safe transitions
Redis distributed locking
Worker Pool for background processing
Webhook handler with HMAC signature validation
Tokenization service (PCI DSS scope reduction)
Rate limiting using Redis
Clean Architecture (Domain / Application / Infrastructure)

How to Run
