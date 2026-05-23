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

- Idempotent Transaction API
- Finite State Machine with safe state transitions
- Redis distributed locking
- Worker pool for background processing
- Webhook handler with HMAC signature validation
- Tokenization service for PCI DSS scope reduction
- Redis-based rate limiting
- Clean Architecture:
  - Domain
  - Application
  - Infrastructure

---

## How to Run

### Start Dependencies

```bash
docker-compose up -d
```

### Run API

```bash
go run cmd/api/main.go
```

---

## Test API

```bash
curl -X POST http://localhost:8080/api/v1/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "idempotency_key": "test-001",
    "source_wallet_id": "w001",
    "dest_wallet_id": "w002",
    "amount": 100000,
    "currency": "IDR"
  }'
```