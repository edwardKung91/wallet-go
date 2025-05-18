package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"go-wallet/internal/wallet"
)

func Setup(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	h := wallet.NewHandler(db)

	r.HandleFunc("/wallet", h.CreateWallet).Methods("POST")
	r.HandleFunc("/wallet/deposit", h.Deposit).Methods("POST")
	r.HandleFunc("/wallet/withdraw", h.Withdraw).Methods("POST")
	r.HandleFunc("/wallet/transfer", h.Transfer).Methods("POST")
	r.HandleFunc("/wallet/balance", h.GetBalance).Methods("GET")
	r.HandleFunc("/wallet/transactions", h.GetTransactions).Methods("GET")

	return r
}
