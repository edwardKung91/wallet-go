package wallet

import (
	"github.com/google/uuid" // UUID type for unique IDs
	"time"                   // To handle timestamps
)

// Wallet struct represents a user's wallet with a unique ID, the owner's user ID, and current balance.
type Wallet struct {
	ID      uuid.UUID `json:"id"`      // Unique wallet ID
	UserID  uuid.UUID `json:"user_id"` // Owner's user ID
	Balance int64     `json:"balance"` // Wallet balance (in smallest currency unit, e.g. cents)
}

// Transaction struct represents a record of money movement involving wallets.
type Transaction struct {
	ID         uuid.UUID  `json:"id"`          // Unique transaction ID
	FromWallet *uuid.UUID `json:"from_wallet"` // Wallet sending money (nullable for deposits)
	ToWallet   *uuid.UUID `json:"to_wallet"`   // Wallet receiving money (nullable for withdrawals)
	Amount     int64      `json:"amount"`      // Transaction amount
	Type       string     `json:"type"`        // Type of transaction: deposit, withdrawal, transfer
	CreatedAt  time.Time  `json:"created_at"`  // Timestamp of the transaction
}
