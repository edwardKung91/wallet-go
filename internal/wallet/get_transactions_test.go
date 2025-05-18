package wallet

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
