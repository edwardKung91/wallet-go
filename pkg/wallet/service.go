package wallet

import (
	"database/sql"           // SQL DB operations
	"github.com/google/uuid" // UUID generation and parsing
	"log"
	"time" // For timestamps
)

// service struct holds the DB reference and encapsulates business logic.
type service struct {
	db *sql.DB
}

// NewService initializes a new service instance with the given DB connection.
func NewService(db *sql.DB) *service {
	return &service{db: db}
}

// WalletExists Checks if the wallet to be updated exists
func (s *service) WalletExists(walletID uuid.UUID) (bool, error) {
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1)`, walletID).Scan(&exists)
	if err != nil {
		log.Printf("DB Select error: %v", err)
		return false, err
	}
	return exists, nil
}

// CreateWallet inserts a new wallet with zero balance for a user.
func (s *service) CreateWallet(userID uuid.UUID) (*wallet, error) {
	id := uuid.New() // Generate a new wallet UUID
	_, err := s.db.Exec(`INSERT INTO wallets (id, user_id, balance) VALUES ($1, $2, $3)`, id, userID, 0)
	if err != nil {
		log.Printf("DB Insertion error: %v", err)
		return nil, err
	}
	return &wallet{ID: id, UserID: userID, Balance: 0}, nil
}

// Deposit adds money to a specific wallet and logs the transaction.
func (s *service) Deposit(walletID uuid.UUID, amount int64) (uuid.UUID, error) {
	if amount <= 0 {
		return uuid.Nil, ErrInvalidAmount
	}

	exists, err := s.WalletExists(walletID)
	if err != nil {
		return uuid.Nil, err
	}
	if !exists {
		return uuid.Nil, ErrWalletNotFound
	}

	// Begin transaction to ensure atomicity
	txn, err := s.db.Begin()
	if err != nil {
		log.Printf("DB Begin error: %v", err)
		return uuid.Nil, err
	}
	defer txn.Rollback()

	// Update wallet balance
	_, err = txn.Exec(`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, amount, walletID)
	if err != nil {
		log.Printf("DB Update error: %v", err)
		return uuid.Nil, err
	}

	txnId := uuid.New()

	// Log transaction as "deposit"
	_, err = txn.Exec(`INSERT INTO transactions (id, from_wallet, to_wallet, amount, type, created_at)
                      VALUES ($1, NULL, $2, $3, $4, $5)`,
		txnId, walletID, amount, TxnTypeDeposit, time.Now())
	if err != nil {
		log.Printf("DB Insert error: %v", err)
		return uuid.Nil, err
	}

	if err := txn.Commit(); err != nil {
		log.Printf("DB Commit error: %v", err)
		return uuid.Nil, err
	}

	return txnId, nil
}

// Withdraw subtracts money from a wallet if there's enough balance.
func (s *service) Withdraw(walletID uuid.UUID, amount int64) (uuid.UUID, error) {
	if amount <= 0 {
		return uuid.Nil, ErrInvalidAmount
	}

	exists, err := s.WalletExists(walletID)
	if err != nil {
		return uuid.Nil, err
	}
	if !exists {
		return uuid.Nil, ErrWalletNotFound
	}

	txn, err := s.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	defer txn.Rollback()

	// Check current balance
	var balance int64
	err = txn.QueryRow(`SELECT balance FROM wallets WHERE id = $1`, walletID).Scan(&balance)
	if err != nil {
		log.Printf("DB Select error: %v", err)
		return uuid.Nil, err
	}

	if balance < amount {
		return uuid.Nil, ErrInsufficientFunds
	}

	// Deduct from wallet
	_, err = txn.Exec(`UPDATE wallets SET balance = balance - $1 WHERE id = $2`, amount, walletID)
	if err != nil {
		log.Printf("DB UPDATE error: %v", err)
		return uuid.Nil, err
	}

	txnId := uuid.New()

	// Log transaction as "withdrawal"
	_, err = txn.Exec(`INSERT INTO transactions (id, from_wallet, to_wallet, amount, type, created_at)
                      VALUES ($1, $2, NULL, $3, $4, $5)`,
		txnId, walletID, amount, TxnTypeWithdrawal, time.Now())
	if err != nil {
		log.Printf("DB Insert error: %v", err)
		return uuid.Nil, err
	}

	if err := txn.Commit(); err != nil {
		return uuid.Nil, err
	}

	return txnId, nil
}

// Transfer moves funds from one wallet to another in a single atomic transaction.
func (s *service) Transfer(fromID, toID uuid.UUID, amount int64) (uuid.UUID, error) {
	if amount <= 0 {
		return uuid.Nil, ErrInvalidAmount
	}

	if fromID == toID {
		return uuid.Nil, ErrSameWalletTransfer
	}

	frmExists, err := s.WalletExists(fromID)
	if err != nil {
		return uuid.Nil, err
	}
	if !frmExists {
		return uuid.Nil, ErrSourceInvalid
	}

	toExists, err := s.WalletExists(toID)
	if err != nil {
		return uuid.Nil, err
	}
	if !toExists {
		return uuid.Nil, ErrDestinationInvalid
	}

	txn, err := s.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	defer txn.Rollback()

	var balance int64
	err = txn.QueryRow(`SELECT balance FROM wallets WHERE id = $1`, fromID).Scan(&balance)
	if err != nil {
		return uuid.Nil, err
	}

	if balance < amount {
		return uuid.Nil, ErrInsufficientFunds
	}

	// Subtract from sender
	_, err = txn.Exec(`UPDATE wallets SET balance = balance - $1 WHERE id = $2`, amount, fromID)
	if err != nil {
		return uuid.Nil, err
	}

	// Add to receiver
	_, err = txn.Exec(`UPDATE wallets SET balance = balance + $1 WHERE id = $2`, amount, toID)
	if err != nil {
		return uuid.Nil, err
	}

	txnId := uuid.New()

	// Log the transaction as "transfer"
	_, err = txn.Exec(`INSERT INTO transactions (id, from_wallet, to_wallet, amount, type, created_at)
                      VALUES ($1, $2, $3, $4, $5, $6)`,
		txnId, fromID, toID, amount, TxnTypeTransfer, time.Now())
	if err != nil {
		return uuid.Nil, err
	}

	if err := txn.Commit(); err != nil {
		return uuid.Nil, err
	}

	return txnId, nil
}

// GetBalance returns the current balance of a wallet.
func (s *service) GetBalance(walletID uuid.UUID) (int64, error) {
	var balance int64

	exists, err := s.WalletExists(walletID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, ErrWalletNotFound
	}

	err = s.db.QueryRow(`SELECT balance FROM wallets WHERE id = $1`, walletID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// GetTransactions fetches all transactions where the wallet was either sender or receiver.
func (s *service) GetTransactions(walletID uuid.UUID) ([]transaction, error) {

	exists, err := s.WalletExists(walletID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrWalletNotFound
	}

	rows, err := s.db.Query(`
        SELECT id, from_wallet, to_wallet, amount, type, created_at
        FROM transactions
        WHERE from_wallet = $1 OR to_wallet = $1
        ORDER BY created_at DESC`, walletID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txns []transaction
	for rows.Next() {
		var txn transaction
		err := rows.Scan(&txn.ID, &txn.FromWallet, &txn.ToWallet, &txn.Amount, &txn.Type, &txn.CreatedAt)
		if err != nil {
			return nil, err
		}
		txns = append(txns, txn)
	}

	return txns, nil
}

// Compile-time check to ensure service implements Service interface
var _ Service = (*service)(nil)
