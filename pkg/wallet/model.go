package wallet

import (
	"github.com/google/uuid" // UUID type for unique IDs
	"time"                   // To handle timestamps
)

// wallet struct represents a user's wallet with a unique ID, the owner's user ID, and current balance.
type wallet struct {
	ID      uuid.UUID `json:"id"`      // Unique wallet ID
	UserID  uuid.UUID `json:"user_id"` // Owner's user ID
	Balance int64     `json:"balance"` // Wallet balance (in smallest currency unit, e.g. cents)
}

// transaction struct represents a record of money movement involving wallets.
type transaction struct {
	ID         uuid.UUID  `json:"id"`          // Unique transaction ID
	FromWallet *uuid.UUID `json:"from_wallet"` // Wallet sending money (nullable for deposits)
	ToWallet   *uuid.UUID `json:"to_wallet"`   // Wallet receiving money (nullable for withdrawals)
	Amount     int64      `json:"amount"`      // transaction amount
	Type       string     `json:"type"`        // Type of transaction: deposit, withdrawal, transfer
	CreatedAt  time.Time  `json:"created_at"`  // Timestamp of the transaction
}
