package dto

import "github.com/google/uuid"

type GetWalletResponse struct {
	WalletID uuid.UUID `json:"walletId"`
	Balance  float64   `json:"balance"`
}
