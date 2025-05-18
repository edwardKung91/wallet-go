package wallet

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
