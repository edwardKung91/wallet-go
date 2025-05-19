package wallet

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateWallet(t *testing.T) {
	mock := &mockService{
		MockCreateWallet: func(userID uuid.UUID) (*wallet, error) {
			return &wallet{ID: uuid.New(), UserID: userID}, nil
		},
	}
	h := NewHandler(mock)

	t.Run("valid request", func(t *testing.T) {
		userID := uuid.New().String()
		body := []byte(`{"user_id":"` + userID + `"}`)
		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(body))
		res := httptest.NewRecorder()

		h.CreateWallet(res, req)
		if res.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d", res.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer([]byte("{invalid}")))
		res := httptest.NewRecorder()

		h.CreateWallet(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", res.Code)
		}
	})

	t.Run("invalid UUID", func(t *testing.T) {
		body := []byte(`{"user_id":"not-a-uuid"}`)
		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(body))
		res := httptest.NewRecorder()

		h.CreateWallet(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", res.Code)
		}
	})
}

// TestDeposit tests the Deposit handler using wallet_id in URL
func TestDeposit(t *testing.T) {
	mock := &mockService{
		MockDeposit: func(walletID uuid.UUID, amount int64) (uuid.UUID, error) {
			return uuid.New(), nil
		},
	}
	h := NewHandler(mock)

	t.Run("valid request", func(t *testing.T) {
		id := uuid.New()
		body := []byte(`{"amount": 100}`)
		req := httptest.NewRequest(http.MethodPost, "/wallet/"+id.String()+"/deposit", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"wallet_id": id.String()})
		res := httptest.NewRecorder()

		h.Deposit(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", res.Code)
		}
	})

	t.Run("invalid wallet_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wallet/invalid-uuid/deposit", nil)
		req = mux.SetURLVars(req, map[string]string{"wallet_id": "invalid-uuid"})
		res := httptest.NewRecorder()

		h.Deposit(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", res.Code)
		}
	})
}

// TestWithdraw tests the Withdraw handler using wallet_id in URL
func TestWithdraw(t *testing.T) {
	mock := &mockService{
		MockWithdraw: func(walletID uuid.UUID, amount int64) (uuid.UUID, error) {
			return uuid.New(), nil
		},
	}
	h := NewHandler(mock)

	t.Run("valid request", func(t *testing.T) {
		id := uuid.New()
		body := []byte(`{"amount": 50}`)
		req := httptest.NewRequest(http.MethodPost, "/wallet/"+id.String()+"/withdraw", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"wallet_id": id.String()})
		res := httptest.NewRecorder()

		h.Withdraw(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", res.Code)
		}
	})

	t.Run("invalid wallet_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wallet/invalid-uuid/withdraw", nil)
		req = mux.SetURLVars(req, map[string]string{"wallet_id": "invalid-uuid"})
		res := httptest.NewRecorder()

		h.Withdraw(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", res.Code)
		}
	})
}

// TestGetBalance with wallet_id in URL query param
func TestGetBalance(t *testing.T) {
	mock := &mockService{
		MockGetBalance: func(walletID uuid.UUID) (int64, error) {
			return 500, nil
		},
	}
	h := NewHandler(mock)

	t.Run("valid wallet_id", func(t *testing.T) {
		id := uuid.New().String()
		req := httptest.NewRequest(http.MethodGet, "/wallet/"+id+"/balance", nil)
		req = mux.SetURLVars(req, map[string]string{"wallet_id": id})
		log.Print(req)
		res := httptest.NewRecorder()

		h.GetBalance(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", res.Code)
		}
	})

	t.Run("invalid wallet_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wallet/invalid-uuid/balance", nil)
		req = mux.SetURLVars(req, map[string]string{"wallet_id": "invalid-uuid"})
		res := httptest.NewRecorder()

		h.GetBalance(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", res.Code)
		}
	})
}

// TestGetTransactions with wallet_id in URL query param
func TestGetTransactions(t *testing.T) {
	mock := &mockService{
		MockGetTransactions: func(walletID uuid.UUID) ([]transaction, error) {
			return []transaction{
				{Amount: 100, Type: "deposit"},
			}, nil
		},
	}
	h := NewHandler(mock)

	t.Run("valid wallet_id", func(t *testing.T) {
		id := uuid.New().String()
		req := httptest.NewRequest(http.MethodGet, "/wallet/"+id+"/transactions", nil)
		req = mux.SetURLVars(req, map[string]string{"wallet_id": id})
		res := httptest.NewRecorder()

		h.GetTransactions(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", res.Code)
		}
	})

	t.Run("invalid wallet_id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wallet/invalid-uuid/transactions", nil)
		req = mux.SetURLVars(req, map[string]string{"wallet_id": "invalid-uuid"})
		res := httptest.NewRecorder()

		h.GetTransactions(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", res.Code)
		}
	})
}

// TestTransfer still takes wallet IDs from request body
func TestTransfer(t *testing.T) {
	mock := &mockService{
		MockTransfer: func(from uuid.UUID, to uuid.UUID, amt int64) (uuid.UUID, error) {
			return uuid.New(), nil
		},
	}
	h := NewHandler(mock)

	t.Run("valid transfer", func(t *testing.T) {
		fromID := uuid.New().String()
		toID := uuid.New().String()
		body := []byte(`{"from_id":"` + fromID + `", "to_id":"` + toID + `", "amount":100}`)
		req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(body))
		res := httptest.NewRecorder()

		h.Transfer(res, req)
		if res.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", res.Code)
		}
	})

	t.Run("invalid from_id", func(t *testing.T) {
		body := []byte(`{"from_id":"invalid", "to_id":"` + uuid.New().String() + `", "amount":100}`)
		req := httptest.NewRequest(http.MethodPost, "/transfer", bytes.NewBuffer(body))
		res := httptest.NewRecorder()

		h.Transfer(res, req)
		if res.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", res.Code)
		}
	})
}
