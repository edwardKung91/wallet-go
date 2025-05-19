package wallet

import "errors"

var (
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrInsufficientFunds  = errors.New("insufficient balance")
	ErrInvalidAmount      = errors.New("amount must be greater than zero")
	ErrSameWalletTransfer = errors.New("cannot transfer to the same wallet")
	ErrSourceInvalid      = errors.New("sender wallet does not exist")
	ErrDestinationInvalid = errors.New("recipient wallet does not exist")
)
