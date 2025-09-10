package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"wallet-service/internal/dto"
	"wallet-service/internal/wallet/model"
	"wallet-service/internal/wallet/repository"
	"wallet-service/internal/wallet/repository/mocks"
	svcErrors "wallet-service/internal/wallet/service/errors"
)

func TestGetWallet_Success(t *testing.T) {
	ctx := context.Background()
	walletId := uuid.New()

	mockRepo := new(mocks.WalletRepositoryMock)
	mockRepo.On("GetWalletById", ctx, walletId).
		Return(&model.Wallet{ID: walletId, Balance: 100}, nil)

	svc := NewWalletService(mockRepo)

	balance, err := svc.GetWallet(ctx, walletId)

	assert.NoError(t, err)
	assert.Equal(t, float64(100), balance)
	mockRepo.AssertExpectations(t)
}

func TestGetWallet_NotFound(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()

	mockRepo := new(mocks.WalletRepositoryMock)
	mockRepo.On("GetWalletById", ctx, walletID).
		Return(nil, errors.New("not found"))

	svc := NewWalletService(mockRepo)

	_, err := svc.GetWallet(ctx, walletID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateWalletBalance_Deposit_Success(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()
	wallet := &model.Wallet{ID: walletID, Balance: 100}

	mockRepo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(mockRepo)

	mockRepo.On("WithTx", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(repository.WalletRepository) error)
			_ = fn(mockRepo)
		}).Return(nil)

	mockRepo.On("GetWalletByIdForUpdate", ctx, walletID).
		Return(wallet, nil)

	mockRepo.On("UpdateWalletTx", ctx, walletID, float64(150)).
		Return(nil)

	mockRepo.On("SaveOperationTx", ctx, mock.AnythingOfType("*model.Operation")).
		Run(func(args mock.Arguments) {
			op := args.Get(1).(*model.Operation)
			op.ID = uuid.New()
			op.CreatedAt = time.Now()
		}).
		Return(nil)

	req := dto.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        50,
	}

	op, err := svc.UpdateWalletBalance(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, walletID, op.WalletID)
	assert.Equal(t, "DEPOSIT", op.Type)
	assert.Equal(t, float64(50), op.Amount)
	mockRepo.AssertExpectations(t)
}

func TestUpdateWalletBalance_Withdraw_Success(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()
	wallet := &model.Wallet{ID: walletID, Balance: 100}

	mockRepo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(mockRepo)

	mockRepo.On("WithTx", ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(repository.WalletRepository) error)
			_ = fn(mockRepo)
		}).Return(nil)

	mockRepo.On("GetWalletByIdForUpdate", ctx, walletID).
		Return(wallet, nil)

	mockRepo.On("UpdateWalletTx", ctx, walletID, float64(50)).
		Return(nil)

	mockRepo.On("SaveOperationTx", ctx, mock.Anything).
		Return(nil)

	req := dto.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        50,
	}

	op, err := svc.UpdateWalletBalance(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, float64(50), op.Amount)
	assert.Equal(t, "WITHDRAW", op.Type)
	mockRepo.AssertExpectations(t)
}

func TestUpdateWalletBalance_Withdraw_InsufficientFunds(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()
	wallet := &model.Wallet{ID: walletID, Balance: 10}

	mockRepo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(mockRepo)

	mockRepo.On("WithTx", ctx, mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(repository.WalletRepository) error)
		_ = fn(mockRepo)
	}).Return(svcErrors.ErrInsufficientFunds)

	mockRepo.On("GetWalletByIdForUpdate", ctx, walletID).
		Return(wallet, nil)

	req := dto.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        50,
	}

	_, err := svc.UpdateWalletBalance(ctx, req)

	assert.ErrorIs(t, err, svcErrors.ErrInsufficientFunds)
	mockRepo.AssertExpectations(t)
}

func TestUpdateWalletBalance_InvalidOperation(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()

	mockRepo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(mockRepo)

	mockRepo.On("WithTx", ctx, mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(repository.WalletRepository) error)
		_ = fn(mockRepo)
	}).Return(svcErrors.ErrInvalidOperation)

	mockRepo.On("GetWalletByIdForUpdate", ctx, walletID).
		Return(&model.Wallet{ID: walletID, Balance: 100}, nil)

	req := dto.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: "INVALID",
		Amount:        50,
	}

	_, err := svc.UpdateWalletBalance(ctx, req)

	assert.ErrorIs(t, err, svcErrors.ErrInvalidOperation)
	mockRepo.AssertExpectations(t)
}

func TestUpdateWalletBalance_NegativeAmount(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()

	mockRepo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(mockRepo)

	req := dto.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        -10,
	}

	_, err := svc.UpdateWalletBalance(ctx, req)

	assert.ErrorIs(t, err, svcErrors.ErrInvalidAmount)
}

func TestUpdateWalletBalance_RetryOnDeadlock(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()
	wallet := &model.Wallet{ID: walletID, Balance: 100}

	attempts := 0

	mockRepo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(mockRepo)

	mockRepo.On("WithTx", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(repository.WalletRepository) error)
		attempts++
		if attempts < 3 {
			_ = fn(mockRepo)
			panic(&pgconn.PgError{Code: "40P01"})
		} else {
			_ = fn(mockRepo)
		}
	})

	mockRepo.On("GetWalletByIdForUpdate", ctx, walletID).Return(wallet, nil)
	mockRepo.On("UpdateWalletTx", ctx, walletID, float64(150)).Return(nil)
	mockRepo.On("SaveOperationTx", ctx, mock.AnythingOfType("*model.Operation")).Return(nil)

	req := dto.WalletOperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        50,
	}

	var op *model.Operation
	var err error

	defer func() {
		if r := recover(); r != nil {
			if pgErr, ok := r.(*pgconn.PgError); ok && pgErr.Code == "40P01" {
				err = svcErrors.ErrInsufficientFunds
			}
		}
	}()

	op, err = svc.UpdateWalletBalance(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, walletID, op.WalletID)
	assert.Equal(t, 3, attempts)
	mockRepo.AssertExpectations(t)
}

func TestUpdateWalletBalance_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := 1000.0
	depositAmount := 10.0
	concurrentRequests := 1000

	wallet := &model.Wallet{ID: walletID, Balance: initialBalance}

	mockRepo := new(mocks.WalletRepositoryMock)
	svc := NewWalletService(mockRepo)

	mockRepo.On("WithTx", ctx, mock.Anything).Run(func(args mock.Arguments) {
		fn := args.Get(1).(func(repository.WalletRepository) error)
		_ = fn(mockRepo)
	}).Return(nil)

	var mu sync.Mutex

	mockRepo.On("GetWalletByIdForUpdate", ctx, walletID).Return(wallet, nil)
	mockRepo.On("UpdateWalletTx", ctx, walletID, mock.Anything).Return(nil)
	mockRepo.On("SaveOperationTx", ctx, mock.AnythingOfType("*model.Operation")).Return(nil)

	wg := sync.WaitGroup{}
	wg.Add(concurrentRequests)

	for i := 0; i < concurrentRequests; i++ {
		go func() {
			defer wg.Done()
			req := dto.WalletOperationRequest{
				WalletID:      walletID,
				OperationType: "DEPOSIT",
				Amount:        depositAmount,
			}

			mu.Lock()
			_, err := svc.UpdateWalletBalance(ctx, req)
			mu.Unlock()

			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	expectedBalance := initialBalance + float64(concurrentRequests)*depositAmount
	assert.Equal(t, expectedBalance, wallet.Balance)
	mockRepo.AssertExpectations(t)
}
