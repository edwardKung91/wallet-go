package wallet

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func newTestService(t *testing.T) (*service, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	svc := &service{db: db}
	return svc, mock, func() { db.Close() }
}

/*
*

	CREATE WALLET Test Cases

*
*/
func TestCreateWallet_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	userID := uuid.New()

	mock.ExpectExec(`INSERT INTO wallets \(id, user_id, balance\)`).
		WithArgs(sqlmock.AnyArg(), userID, int64(0)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	wallet, err := svc.CreateWallet(userID)

	assert.NoError(t, err)
	assert.NotEqual(t, nil, wallet)
}

func TestCreateWallet_InsertFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	userID := uuid.New()

	mock.ExpectExec(`INSERT INTO wallets \(id, user_id, balance\)`).
		WithArgs(sqlmock.AnyArg(), userID, int64(0)).
		WillReturnError(assert.AnError)

	wallet, err := svc.CreateWallet(userID)
	assert.Error(t, err)
	assert.Nil(t, wallet)
}

/*
*

	DEPOSIT Test Cases

*
*/

func TestDeposit_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	amount := int64(200)

	// Expect wallet exists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Begin transaction
	mock.ExpectBegin()

	// Expect update wallet balance
	mock.ExpectExec(`UPDATE wallets SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(amount, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect insert transaction
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(sqlmock.AnyArg(), walletID, amount, TxnTypeDeposit, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect commit
	mock.ExpectCommit()

	txnID, err := svc.Deposit(walletID, amount)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, txnID)
}

func TestDeposit_InvalidAmount(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	txnID, err := svc.Deposit(uuid.New(), 0)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAmount, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestDeposit_WalletNotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	amount := int64(100)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	txnID, err := svc.Deposit(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, ErrWalletNotFound, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestDeposit_BeginTxFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	amount := int64(100)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin().WillReturnError(errors.New("db error"))

	txnID, err := svc.Deposit(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestDeposit_InsertTxnFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	amount := int64(100)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE wallets SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(amount, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(sqlmock.AnyArg(), nil, walletID, amount, TxnTypeDeposit, sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	mock.ExpectRollback()

	txnID, err := svc.Deposit(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, txnID)
}

/*
*

	WITHDRAW Test Cases

*
*/

func TestWithdraw_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	amount := int64(100)
	initialBalance := int64(200)

	// Expect wallet existence check
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Begin transaction
	mock.ExpectBegin()

	// Expect SELECT balance
	mock.ExpectQuery(`SELECT balance FROM wallets WHERE id = \$1`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(initialBalance))

	// Expect UPDATE balance
	mock.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE id = \$2`).
		WithArgs(amount, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect INSERT transaction
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(sqlmock.AnyArg(), walletID, amount, TxnTypeWithdrawal, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Expect Commit
	mock.ExpectCommit()

	id, err := svc.Withdraw(walletID, amount)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, id)
}

func TestWithdraw_InvalidAmount(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	id, err := svc.Withdraw(uuid.New(), 0)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAmount, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestWithdraw_WalletNotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	id, err := svc.Withdraw(walletID, 100)
	assert.Error(t, err)
	assert.Equal(t, ErrWalletNotFound, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestWithdraw_InsufficientBalance(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	amount := int64(500)
	balance := int64(100)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT balance FROM wallets WHERE id = \$1`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(balance))

	mock.ExpectRollback()

	id, err := svc.Withdraw(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientFunds, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestWithdraw_DBError(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin().WillReturnError(errors.New("db down"))

	id, err := svc.Withdraw(walletID, 100)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestWithdraw_InsertTxnFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	amount := int64(100)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin()

	mock.ExpectExec(`UPDATE wallets SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(amount, walletID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(sqlmock.AnyArg(), nil, walletID, amount, TxnTypeWithdrawal, sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	mock.ExpectRollback()

	txnID, err := svc.Withdraw(walletID, amount)
	assert.Error(t, err)
	assert.Equal(t, uuid.Nil, txnID)
}

/*
*

	TRANSFER Test Cases

*
*/

func TestTransfer_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	fromID := uuid.New()
	toID := uuid.New()
	amount := int64(100)

	// WalletExists checks
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(toID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin()

	// Balance check
	mock.ExpectQuery(`SELECT balance FROM wallets WHERE id = \$1`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(1000))

	// Subtract from sender
	mock.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE id = \$2`).
		WithArgs(amount, fromID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Add to receiver
	mock.ExpectExec(`UPDATE wallets SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(amount, toID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Insert transaction
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(sqlmock.AnyArg(), fromID, toID, amount, TxnTypeTransfer, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	txnID, err := svc.Transfer(fromID, toID, amount)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, txnID)
}

func TestTransfer_InsufficientFunds(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	fromID := uuid.New()
	toID := uuid.New()
	amount := int64(1000)

	// WalletExists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(toID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin()

	// Balance is too low
	mock.ExpectQuery(`SELECT balance FROM wallets WHERE id = \$1`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(200)) // < amount

	mock.ExpectRollback()

	txnID, err := svc.Transfer(fromID, toID, amount)
	assert.Error(t, err)
	assert.Equal(t, ErrInsufficientFunds, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestTransfer_SameWallet(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	txnID, err := svc.Transfer(walletID, walletID, 500)
	assert.Error(t, err)
	assert.Equal(t, ErrSameWalletTransfer, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestTransfer_InvalidAmount(t *testing.T) {
	svc, _, cleanup := newTestService(t)
	defer cleanup()

	txnID, err := svc.Transfer(uuid.New(), uuid.New(), -50)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAmount, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestTransfer_SourceWalletNotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	fromID := uuid.New()
	toID := uuid.New()
	amount := int64(100)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	txnID, err := svc.Transfer(fromID, toID, amount)
	assert.Error(t, err)
	assert.Equal(t, ErrSourceInvalid, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestTransfer_DestinationWalletNotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	fromID := uuid.New()
	toID := uuid.New()
	amount := int64(100)

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(toID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	txnID, err := svc.Transfer(fromID, toID, amount)
	assert.Error(t, err)
	assert.Equal(t, ErrDestinationInvalid, err)
	assert.Equal(t, uuid.Nil, txnID)
}

func TestTransfer_InsertFailsAndRollback(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	fromID := uuid.New()
	toID := uuid.New()
	amount := int64(500)

	// WalletExists checks
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(toID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectBegin()

	// Balance check returns enough balance
	mock.ExpectQuery(`SELECT balance FROM wallets WHERE id = \$1`).
		WithArgs(fromID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(1000))

	// Subtract from sender
	mock.ExpectExec(`UPDATE wallets SET balance = balance - \$1 WHERE id = \$2`).
		WithArgs(amount, fromID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Add to receiver
	mock.ExpectExec(`UPDATE wallets SET balance = balance \+ \$1 WHERE id = \$2`).
		WithArgs(amount, toID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Simulate INSERT failure
	mock.ExpectExec(`INSERT INTO transactions`).
		WithArgs(sqlmock.AnyArg(), fromID, toID, amount, TxnTypeTransfer, sqlmock.AnyArg()).
		WillReturnError(errors.New("insert failed"))

	mock.ExpectRollback()

	// Call the actual service
	txnID, err := svc.Transfer(fromID, toID, amount)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insert failed")
	assert.Equal(t, uuid.Nil, txnID)
}

/*
*

	GET BALANCE Test Cases

*
*/

func TestGetBalance_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()
	expectedBalance := int64(1000)

	// Expect WalletExists to return true
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Expect balance query
	mock.ExpectQuery(`SELECT balance FROM wallets WHERE id = \$1`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"balance"}).AddRow(expectedBalance))

	// Run test
	balance, err := svc.GetBalance(walletID)

	assert.NoError(t, err)
	assert.Equal(t, expectedBalance, balance)
}

func TestGetBalance_WalletNotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	// Expect WalletExists to return false
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	// Run test
	balance, err := svc.GetBalance(walletID)

	assert.Error(t, err)
	assert.Equal(t, ErrWalletNotFound, err)
	assert.Equal(t, int64(0), balance)
}

func TestGetBalance_WalletExistsQueryFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	// Simulate query error
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnError(assert.AnError)

	// Run test
	balance, err := svc.GetBalance(walletID)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Equal(t, int64(0), balance)
}

func TestGetBalance_SelectFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(`SELECT balance FROM wallets WHERE id = \$1`).
		WithArgs(walletID).
		WillReturnError(assert.AnError)

	// Run test
	balance, err := svc.GetBalance(walletID)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	assert.Equal(t, int64(0), balance)
}

/*
*

	GET TRANSACTION Test Cases

*
*/

func TestGetTransactions_Success(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	// Expect wallet exists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Prepare mock transactions
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "from_wallet", "to_wallet", "amount", "type", "created_at"}).
		AddRow(uuid.New(), walletID, uuid.New(), int64(100), TxnTypeTransfer, now).
		AddRow(uuid.New(), nil, walletID, int64(200), TxnTypeDeposit, now.Add(-time.Minute))

	mock.ExpectQuery(`SELECT id, from_wallet, to_wallet, amount, type, created_at FROM transactions`).
		WithArgs(walletID).
		WillReturnRows(rows)

	// Execute
	txns, err := svc.GetTransactions(walletID)

	assert.NoError(t, err)
	assert.Len(t, txns, 2)
	assert.Equal(t, int64(100), txns[0].Amount)
	assert.Equal(t, TxnTypeTransfer, txns[0].Type)
	assert.Equal(t, TxnTypeDeposit, txns[1].Type)
}

func TestGetTransactions_WalletNotFound(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

	txns, err := svc.GetTransactions(walletID)

	assert.Error(t, err)
	assert.Equal(t, ErrWalletNotFound, err)
	assert.Nil(t, txns)
}

func TestGetTransactions_WalletExistsQueryFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnError(errors.New("query failed"))

	txns, err := svc.GetTransactions(walletID)

	assert.Error(t, err)
	assert.Nil(t, txns)
}

func TestGetTransactions_QueryFails(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	mock.ExpectQuery(`SELECT id, from_wallet, to_wallet, amount, type, created_at FROM transactions`).
		WithArgs(walletID).
		WillReturnError(errors.New("query failed"))

	txns, err := svc.GetTransactions(walletID)

	assert.Error(t, err)
	assert.Nil(t, txns)
}

func TestGetTransactions_ScanFails_MissingFromWallet(t *testing.T) {
	svc, mock, cleanup := newTestService(t)
	defer cleanup()

	walletID := uuid.New()

	// Simulate wallet exists
	mock.ExpectQuery(`SELECT EXISTS\(SELECT 1 FROM wallets WHERE id = \$1\)`).
		WithArgs(walletID).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	// Simulate corrupted row missing `from_wallet`
	badRows := sqlmock.NewRows([]string{"id" /* missing from_wallet */, "to_wallet", "amount", "type", "created_at"}).
		AddRow(uuid.New(), uuid.New(), int64(100), TxnTypeTransfer, time.Now())

	mock.ExpectQuery(`SELECT id, from_wallet, to_wallet, amount, type, created_at FROM transactions`).
		WithArgs(walletID).
		WillReturnRows(badRows)

	// Execute the service call
	txns, err := svc.GetTransactions(walletID)

	assert.Error(t, err)
	assert.Nil(t, txns)
}
