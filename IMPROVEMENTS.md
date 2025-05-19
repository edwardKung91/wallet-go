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

### 3. Make the transaction schema more extensible for new transaction types
- **Why**: There could be other transaction types like authorisations or even for other types of accounts like lending which have a drawdown
- **How**: We can just keep the type column as kind of like a metadata that can be of any value instead of checking against specific set of types

### 4. Better mapping of internal errors to HTTP error codes
- **Why**: This ensures that all errors produced by the system are mapped to an appropriate client facing error.
- **How**: Improve the repository of errors we currently have and create better mapping of those errors to corresponding http codes

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

### 4. Hiding Secrets like passwords and DB info
- **Why**: No files on the repository should contain passwords and even potentially usernames for security
- **How**: Hashicorp vault can be used to keep such secrets and loaded from there by the service

### 5. Improve masking of internal errors
- **Why**: Internal errors should only be visible in internal logging to ensure that no information on our infrastructure is exposed
- **How**: With better mapping of errors and better catching of these in the handler we can respond to the client with more generic error messages

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

