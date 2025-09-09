package repository

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"wallet-service/internal/wallet/model"
)

type WalletRepository interface {
	GetWalletById(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
	UpdateWalletTx(ctx context.Context, id uuid.UUID, newBalance float64) error
	SaveOperationTx(ctx context.Context, op *model.Operation) error
	GetWalletByIdForUpdate(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
	WithTx(ctx context.Context, fn func(txRepo WalletRepository) error) error
}

type walletRepository struct {
	db *gorm.DB
}

func (w *walletRepository) GetWalletById(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	var wallet model.Wallet

	if err := w.db.First(&wallet, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (w *walletRepository) GetWalletByIdForUpdate(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	var wallet model.Wallet
	if err := w.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&wallet, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (w *walletRepository) UpdateWalletTx(ctx context.Context, id uuid.UUID, newBalance float64) error {
	return w.db.Model(&model.Wallet{}).Where("id = ?", id).Update("balance", newBalance).Error
}

func (w *walletRepository) SaveOperationTx(ctx context.Context, op *model.Operation) error {
	return w.db.Create(op).Error
}

func (w *walletRepository) WithTx(ctx context.Context, fn func(txRepo WalletRepository) error) error {
	return w.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txRepo := &walletRepository{db: tx}
		return fn(txRepo)
	})
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}
