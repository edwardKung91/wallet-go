package wallet

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
