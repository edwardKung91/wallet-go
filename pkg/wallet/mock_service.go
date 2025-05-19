package wallet

import (
	"github.com/google/uuid"
)

// MockService implements the Service interface for testing.
type mockService struct {
	MockCreateWallet    func(uuid.UUID) (*wallet, error)
	MockDeposit         func(uuid.UUID, int64) (uuid.UUID, error)
	MockWithdraw        func(uuid.UUID, int64) (uuid.UUID, error)
	MockTransfer        func(uuid.UUID, uuid.UUID, int64) (uuid.UUID, error)
	MockGetBalance      func(uuid.UUID) (int64, error)
	MockGetTransactions func(uuid.UUID) ([]transaction, error)
}

func (m *mockService) CreateWallet(userID uuid.UUID) (*wallet, error) {
	return m.MockCreateWallet(userID)
}
func (m *mockService) Deposit(walletID uuid.UUID, amount int64) (uuid.UUID, error) {
	return m.MockDeposit(walletID, amount)
}
func (m *mockService) Withdraw(walletID uuid.UUID, amount int64) (uuid.UUID, error) {
	return m.MockWithdraw(walletID, amount)
}
func (m *mockService) Transfer(from uuid.UUID, to uuid.UUID, amount int64) (uuid.UUID, error) {
	return m.MockTransfer(from, to, amount)
}
func (m *mockService) GetBalance(walletID uuid.UUID) (int64, error) {
	return m.MockGetBalance(walletID)
}
func (m *mockService) GetTransactions(walletID uuid.UUID) ([]transaction, error) {
	return m.MockGetTransactions(walletID)
}
