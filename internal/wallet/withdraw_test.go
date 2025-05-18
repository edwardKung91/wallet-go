package wallet

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
