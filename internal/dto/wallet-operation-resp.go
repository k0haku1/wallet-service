package dto

import "github.com/google/uuid"

type WalletOperationResponse struct {
	OperationID   uuid.UUID `json:"operationId"`
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        float64   `json:"amount"`
	CreatedAt     string    `json:"createdAt"`
}
