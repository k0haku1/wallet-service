package dto

import "github.com/google/uuid"

type WalletOperationRequest struct {
	WalletID      uuid.UUID `json:"walletId" validate:"required"`
	OperationType string    `json:"operationType" validate:"required,oneof=DEPOSIT WITHDRAW"`
	Amount        float64   `json:"amount" validate:"required,gt=0"`
}
