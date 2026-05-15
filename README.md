# Fintech Playground - Transaction Service

**Immutable ledger with idempotency protection.**

## Architecture

Client
   |
   v
HTTP Handler
   |
   v
Transaction Service  ----> Idempotency Guard
   |
   v
Transaction Repository
   |
   v
PostgreSQL Ledger (immutable)



Client          Handler         Service          Repository         PostgreSQL
  |                 |               |                 |                 |
  |--- POST /tx --->|               |                 |                 |
  |                 |--- Create --->|                 |                 |
  |                 |               |--- FindByKey -->|                 |
  |                 |               |                 |--- SELECT ----->|
  |                 |               |<-- Exists ------|                 |
  |                 |<-- Return tx--|                 |                 |
  |<-- 201 Created--|               |                 |                 |