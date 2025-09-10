package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"gorm.io/gorm"
	"time"
	"wallet-service/internal/dto"
	"wallet-service/internal/wallet/model"
	"wallet-service/internal/wallet/repository"
	svcErrors "wallet-service/internal/wallet/service/errors"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{
		repo: repo,
	}
}
func (s *WalletService) GetWallet(ctx context.Context, id uuid.UUID) (float64, error) {
	wallet, err := s.repo.GetWalletById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, svcErrors.ErrWalletNotFound
		}
		return 0, err
	}
	return wallet.Balance, nil
}

func (s *WalletService) UpdateWalletBalance(ctx context.Context, req dto.WalletOperationRequest) (*model.Operation, error) {
	if req.Amount < 0 {
		return nil, svcErrors.ErrInvalidAmount
	}
	var op *model.Operation
	var err error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {

		err = s.repo.WithTx(ctx, func(txRepo repository.WalletRepository) error {
			wallet, err := txRepo.GetWalletByIdForUpdate(ctx, req.WalletID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return svcErrors.ErrWalletNotFound
				}
				return err
			}

			switch req.OperationType {
			case "DEPOSIT":
				wallet.Balance += req.Amount
			case "WITHDRAW":
				if wallet.Balance < req.Amount {
					return svcErrors.ErrInsufficientFunds
				}
				wallet.Balance -= req.Amount
			default:
				return svcErrors.ErrInvalidOperation
			}

			if err := txRepo.UpdateWalletTx(ctx, wallet.ID, wallet.Balance); err != nil {
				return err
			}

			op = &model.Operation{
				ID:       uuid.New(),
				WalletID: wallet.ID,
				Type:     req.OperationType,
				Amount:   req.Amount,
			}

			return txRepo.SaveOperationTx(ctx, op)
		})

		if err == nil || !isDeadlock(err) || isSerializationFailure(err) {
			break
		}
		time.Sleep(time.Millisecond * time.Duration(50*(i+1)))
	}

	return op, err
}

func isDeadlock(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40P01"
	}
	return false
}

func isSerializationFailure(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "40001"
	}
	return false
}
