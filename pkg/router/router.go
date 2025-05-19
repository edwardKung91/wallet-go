package router

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"wallet-go/pkg/wallet"
)

func Setup(db *sql.DB) http.Handler {
	r := mux.NewRouter()
	s := wallet.NewService(db)
	h := wallet.NewHandler(s)

	r.HandleFunc("/wallet", h.CreateWallet).Methods("POST")
	r.HandleFunc("/wallet/{wallet_id}/deposit", h.Deposit).Methods("POST")
	r.HandleFunc("/wallet/{wallet_id}/withdraw", h.Withdraw).Methods("POST")
	r.HandleFunc("/wallet/transfer", h.Transfer).Methods("POST")
	r.HandleFunc("/wallet/{wallet_id}/balance", h.GetBalance).Methods("GET")
	r.HandleFunc("/wallet/{wallet_id}/transactions", h.GetTransactions).Methods("GET")

	return r
}
