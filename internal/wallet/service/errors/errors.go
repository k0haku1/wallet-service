package errors

import "errors"

var (
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInvalidOperation  = errors.New("invalid operation type")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrInvalidAmount     = errors.New("amount must be positive")
)
