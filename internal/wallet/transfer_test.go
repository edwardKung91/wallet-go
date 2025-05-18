package wallet

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
