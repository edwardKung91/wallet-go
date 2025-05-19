# Go Wallet Service

A simple wallet microservice written in Go.

## This Code was developed in windows environment

## Recommended Versions
- Go -> go version go1.24.3 windows/amd64
- Postgres -> 17.5

## Features

- Create wallet
- Deposit / Withdraw funds
- Transfer funds
- View balance and transaction history

## Tech Stack

- Golang
- PostgreSQL
- Gorilla Mux

## Setup

1. Create PostgreSQL database named `wallet`
2. Run schema from `schema.sql`
3. Set your DB credentials in `.env`
4. Start the server:

```bash
go run ./cmd/server
```

## API Endpoints

| Method | Endpoint              | Description           |
|--------|-----------------------|-----------------------|
| POST   | /wallet               | Create wallet         |
| POST   | /wallet/deposit       | Deposit funds         |
| POST   | /wallet/withdraw      | Withdraw funds        |
| POST   | /wallet/transfer      | Transfer funds        |
| GET    | /wallet/balance       | Get wallet balance    |
| GET    | /wallet/transactions  | Get transaction history|

## API Endpoint Usage
### 1. Create Wallet
    POST /wallet

Request Body:
```
{
"user_id": "uuid-of-user"
}
```

Example:
```
curl -X POST http://localhost:8080/wallet \
-H "Content-Type: application/json" \
-d '{"user_id":"123e4567-e89b-12d3-a456-426614174000"}'
Response:
```
Response:
```
{
"id": "wallet-uuid",
"user_id": "123e4567-e89b-12d3-a456-426614174000",
"balance": 0
}
```

### 2. Deposit Money
    POST /wallet/deposit

Request Body:
```
{
    "wallet_id": "UUID-of-wallet",
    "amount": 5000
}
```

Example:
```
curl -X POST http://localhost:8080/wallet/deposit \
-H "Content-Type: application/json" \
-d '{"wallet_id":"wallet-uuid", "amount":5000}'
```

Response:
```
{
    "status": "success",
    "transaction_id": "UUID-of-transaction"
}
```

### 3. Withdraw Money
    POST /wallet/withdraw

Request Body:
```
{
    "wallet_id": "UUID-of-wallet",
    "amount": 2000
}
```

Example:
```
curl -X POST http://localhost:8080/wallet/withdraw \
-H "Content-Type: application/json" \
-d '{"wallet_id":"wallet-uuid", "amount":2000}'
```

Response:
```
{
    "status": "success",
    "transaction_id": "UUID-of-transaction"
}
```

### 4. Transfer Money
    POST /wallet/transfer

Request Body:
```
{
    "from_id": "UUID-of-sender-wallet",
    "to_id": "UUID-of-recipient-wallet",
    "amount": 1000
}
```

Example:
```
curl -X POST http://localhost:8080/wallet/transfer \
-H "Content-Type: application/json" \
-d '{"from_id":"wallet1-uuid", "to_id":"wallet2-uuid", "amount":1000}'
```

Response:
```
{
    "status": "success",
    "transaction_id": "UUID-of-transaction"
}
```

### 5. Get Wallet Balance
    GET /wallet/balance?wallet_id=UUID-of-wallet

Example:
```
curl http://localhost:8080/wallet/balance?wallet_id=wallet-uuid
```

Response:
```
{
    "balance": 3000
}
```

### 6. Get Transaction History
    GET /wallet/transactions?wallet_id=UUID-of-wallet

Example:
```
curl http://localhost:8080/wallet/transactions?wallet_id=wallet-uuid
```

Response:
```
[
    {
        "id": "txn-uuid",
        "from_wallet": "wallet1-uuid",
        "to_wallet": "wallet2-uuid",
        "amount": 1000,
        "type": "transfer",
        "created_at": "2025-05-17T12:34:56Z"
    },
    {
        "id": "txn-uuid",
        "from_wallet": null,
        "to_wallet": "wallet2-uuid",
        "amount": 5000,
        "type": "deposit",
        "created_at": "2025-05-16T10:00:00Z"
    }
]
```