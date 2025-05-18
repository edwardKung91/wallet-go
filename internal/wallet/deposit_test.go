package wallet

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
