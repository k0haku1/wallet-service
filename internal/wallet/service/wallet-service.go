package service

import "wallet-service/internal/wallet/repository"

type WalletService struct {
	repo *repository.WalletRepository
}

func NewWalletService(repo *repository.WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}
