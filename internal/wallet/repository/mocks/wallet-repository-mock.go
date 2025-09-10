package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"wallet-service/internal/wallet/model"
	"wallet-service/internal/wallet/repository"
)

type WalletRepositoryMock struct {
	mock.Mock
}

var _ repository.WalletRepository = (*WalletRepositoryMock)(nil)

func (m *WalletRepositoryMock) GetWalletById(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	args := m.Called(ctx, id)
	if w, ok := args.Get(0).(*model.Wallet); ok {
		return w, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *WalletRepositoryMock) GetWalletByIdForUpdate(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	args := m.Called(ctx, id)
	if w, ok := args.Get(0).(*model.Wallet); ok {
		return w, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *WalletRepositoryMock) UpdateWalletTx(ctx context.Context, id uuid.UUID, newBalance float64) error {
	args := m.Called(ctx, id, newBalance)
	return args.Error(0)
}

func (m *WalletRepositoryMock) SaveOperationTx(ctx context.Context, op *model.Operation) error {
	args := m.Called(ctx, op)
	return args.Error(0)
}

func (m *WalletRepositoryMock) WithTx(ctx context.Context, fn func(txRepo repository.WalletRepository) error) error {
	args := m.Called(ctx, fn)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}
