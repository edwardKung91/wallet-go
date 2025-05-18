# Wallet Service – Improvements and Future Enhancements

This document outlines areas for potential improvements in the wallet service, organized by category.

---
## Functional Improvements

### 1. Double entry book keeping
- **Why**: Currently withdraws and deposits seem like they come from thin air. (i.e no proper source or destination) This enables proper tracing of money movement between entities outside the company and the company. This also allows for proper reporting and statement generation
- **How**: Creating support for separate general ledger/financial ledger accounts which do not correspond to an individual 
    client but potentially an entity or even accounts internal to the bank (i.e income)

### 2. Actual association of wallets to user
- **Why**: A user/wallet owner may own multiple wallets or even multiple types of accounts from the company. Associating an account/wallet to a user enables this
- **How**: Create a separate user table and make the existing user_id column in the wallets table a foreign key

---

## Performance Improvements

### 1. Use Redis for Balance Caching
- **Why**: To reduce DB load on high-frequency balance checks and transfers.
- **How**: Cache wallet balances in Redis with TTL and invalidate on write operations (deposit/withdraw/transfer).

### 2. Batch Transaction History Queries
- **Why**: Pagination prevents performance bottlenecks with long histories.
- **How**: Implement `limit` and `offset` query parameters on `/wallet/transactions`.

### 3. Prepared Statements
- **Why**: Improves performance by avoiding repeated query planning.
- **How**: Use `sql.Prepare` or a query builder (e.g., `sqlc`, `gorm`).

---

## Scalability Enhancements

### 1. Horizontal Scaling of Application
- **Why**: To handle increased API traffic.
- **How**: Make app stateless; use Redis/PostgreSQL for shared state.

### 2. Sharding Wallet Tables (for huge scale)
- **Why**: Distribute data across DBs when user base grows.
- **How**: Use wallet ID hash mod N to select shard.

### 3. Async Transaction Processing
- **Why**: Decouple DB writes and reduce latency.
- **How**: Use a message queue (e.g., RabbitMQ or Kafka) and worker service.

---

## Security Improvements

### 1. Authentication & Authorization
- **Why**: Ensure only wallet owners can access/modify their data.
- **How**: Integrate JWT-based auth and user context propagation.

### 2. Input Validation & Rate Limiting
- **Why**: Prevent abuse and injection attacks.
- **How**: Validate input types and set rate limits using middleware.

### 3. Transaction Integrity (Double Spending Prevention)
- **Why**: To ensure no race conditions in wallet updates.
- **How**: Use DB row-level locks (`SELECT FOR UPDATE`) or distributed locks if needed. Optimistic locks such as using an update count can also be used

---

## Data Reliability & Resilience

### 1. Retry Logic with Idempotency Keys
- **Why**: Avoid duplicate transactions on network errors.
- **How**: Accept `idempotency_key` headers and store keys with outcomes.

### 2. Backups and Disaster Recovery
- **Why**: Prevent data loss.
- **How**: Automate periodic DB backups; test restore procedures regularly.

---

## Testing Enhancements

### 1. Integration Testing with Docker
- **Why**: Ensure realistic testing environments.
- **How**: Use `docker-compose` to spin up test DB/Redis containers.

### 2. Fuzz Testing
- **Why**: Find edge-case bugs in input handling.
- **How**: Use Go's native fuzzing support.

---

## Developer Experience

### 1. Use an ORM or SQL Generator
- **Why**: Simplify query building and reduce boilerplate.
- **How**: Consider `sqlc` (compile-time safe), `gorm`, or `ent`.

### 2. OpenAPI Documentation
- **Why**: Auto-generate client SDKs and API docs.
- **How**: Use Swagger annotations or `go-swagger`.

---

