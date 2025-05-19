package wallet

import (
	"database/sql"  // Provides SQL DB interaction functions
	"encoding/json" // Used to parse and return JSON
	"github.com/gorilla/mux"
	"net/http" // Used for HTTP request/response handling
	"strings"

	"github.com/google/uuid" // Used to generate and parse UUIDs
)

// handler struct holds a reference to the wallet service.
type handler struct {
	service *service
}

// NewHandler creates a new handler instance with an initialized wallet service.
func NewHandler(db *sql.DB) *handler {
	return &handler{service: NewService(db)}
}

type TransactionResponse struct {
	Status        string     `json:"status"`                   // e.g. "success" or "error"
	TransactionID *uuid.UUID `json:"transaction_id,omitempty"` // optional
	Error         string     `json:"error,omitempty"`          // optional
}

// writeJSON is a helper to write a JSON response with the correct headers.
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// CreateWallet handles the API request to create a wallet for a user.
func (h *handler) CreateWallet(w http.ResponseWriter, r *http.Request) {
	var body struct {
		UserID string `json:"user_id"` // The user ID that the wallet is for
	}

	// Decode JSON request body into `body`
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid JSON body",
		})
		return
	}

	// Validate UUID format
	userID, err := uuid.Parse(strings.TrimSpace(body.UserID))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid user_id format (must be UUID)",
		})
		return
	}

	// Call the service to create a wallet
	wallet, err := h.service.CreateWallet(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the created wallet data
	writeJSON(w, http.StatusCreated, wallet)
}

// Deposit handles the API request to deposit funds into a wallet.
func (h *handler) Deposit(w http.ResponseWriter, r *http.Request) {
	// Extract wallet_id from URL path
	vars := mux.Vars(r)
	walletIDStr := vars["wallet_id"]

	// Validate UUID
	walletID, err := uuid.Parse(strings.TrimSpace(walletIDStr))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid wallet_id format (must be UUID)",
		})
		return
	}

	var body struct {
		Amount int64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid JSON body",
		})
		return
	}

	// Call the service to perform the deposit
	txnId, err := h.service.Deposit(walletID, body.Amount)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status:        "error",
			TransactionID: nil,
			Error:         err.Error(),
		})
		return
	}

	// Respond with HTTP 200 and the txn id of successful txn
	writeJSON(w, http.StatusOK, TransactionResponse{
		Status:        "success",
		Error:         "",
		TransactionID: &txnId,
	})
}

// Withdraw handles the API request to withdraw funds from a wallet.
func (h *handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	// Extract wallet_id from URL path
	vars := mux.Vars(r)
	walletIDStr := vars["wallet_id"]

	// Validate UUID
	walletID, err := uuid.Parse(strings.TrimSpace(walletIDStr))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid wallet_id format (must be UUID)",
		})
		return
	}

	var body struct {
		Amount int64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid JSON body",
		})
		return
	}

	// Call the service to perform the withdrawal
	txnId, err := h.service.Withdraw(walletID, body.Amount)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status:        "error",
			TransactionID: nil,
			Error:         err.Error(),
		})
		return
	}

	// Respond with HTTP 200 and the txn id of successful txn
	writeJSON(w, http.StatusOK, TransactionResponse{
		Status:        "success",
		Error:         "",
		TransactionID: &txnId,
	})
}

// Transfer handles transferring funds from one wallet to another.
func (h *handler) Transfer(w http.ResponseWriter, r *http.Request) {
	var body struct {
		FromID string `json:"from_id"`
		ToID   string `json:"to_id"`
		Amount int64  `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid JSON body",
		})
		return
	}

	// Validate UUID format
	frmWalletID, err := uuid.Parse(strings.TrimSpace(body.FromID))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid Source wallet format (must be UUID)",
		})
		return
	}

	// Validate UUID format
	toWalletID, err := uuid.Parse(strings.TrimSpace(body.ToID))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid Destination wallet format (must be UUID)",
		})
		return
	}

	// Call the service to perform the transfer
	txnId, err := h.service.Transfer(frmWalletID, toWalletID, body.Amount)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status:        "error",
			TransactionID: nil,
			Error:         err.Error(),
		})
		return
	}

	// Respond with HTTP 200 and the txn id of successful txn
	writeJSON(w, http.StatusOK, TransactionResponse{
		Status:        "success",
		Error:         "",
		TransactionID: &txnId,
	})
}

// GetBalance handles retrieving the wallet balance.
func (h *handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	// Extract wallet_id from URL path
	vars := mux.Vars(r)
	walletIDStr := vars["wallet_id"]

	// Validate UUID
	walletID, err := uuid.Parse(strings.TrimSpace(walletIDStr))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid wallet_id format (must be UUID)",
		})
		return
	}

	// Get balance from the service
	balance, err := h.service.GetBalance(walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return the balance as JSON
	writeJSON(w, http.StatusOK, map[string]int64{"balance": balance})
}

// GetTransactions returns the transaction history for a wallet.
func (h *handler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// Extract wallet_id from URL path
	vars := mux.Vars(r)
	walletIDStr := vars["wallet_id"]

	// Validate UUID
	walletID, err := uuid.Parse(strings.TrimSpace(walletIDStr))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, TransactionResponse{
			Status: "error",
			Error:  "Invalid wallet_id format (must be UUID)",
		})
		return
	}

	// Retrieve transactions
	txns, err := h.service.GetTransactions(walletID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the transactions in JSON format
	writeJSON(w, http.StatusOK, txns)
}
