package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"wallet-service/internal/wallet/model"
)

type WalletRepository interface {
	GetWalletById(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
	UpdateWalletTx(ctx context.Context, tx *gorm.DB, id uuid.UUID, newBalance float64) error
	SaveOperationTx(ctx context.Context, tx *gorm.DB, op *model.Operation) error
}

type walletRepository struct {
	db *gorm.DB
}

func (w *walletRepository) GetWalletById(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w *walletRepository) UpdateWalletTx(ctx context.Context, tx *gorm.DB, id uuid.UUID, newBalance float64) error {
	//TODO implement me
	panic("implement me")
}

func (w *walletRepository) SaveOperationTx(ctx context.Context, tx *gorm.DB, op *model.Operation) error {
	//TODO implement me
	panic("implement me")
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}
