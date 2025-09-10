package integration_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"wallet-service/internal/config"
	"wallet-service/internal/db"
	"wallet-service/internal/dto"
	"wallet-service/internal/wallet/model"
	"wallet-service/internal/wallet/repository"
	"wallet-service/internal/wallet/service"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	cfg := config.LoadConfig()
	testDB = db.NewPostgres(cfg.DBUrl)
	m.Run()
}

func createTestWallet(db *gorm.DB, initialBalance float64) *model.Wallet {
	wallet := &model.Wallet{
		ID:      uuid.New(),
		Balance: initialBalance,
	}
	if err := db.Create(wallet).Error; err != nil {
		panic(err)
	}
	return wallet
}

func TestConcurrentDeposits(t *testing.T) {
	ctx := context.Background()
	initialBalance := 1000.0
	depositAmount := 10.0
	totalRequests := 1000
	maxParallel := 50

	wallet := createTestWallet(testDB, initialBalance)
	repo := repository.NewWalletRepository(testDB)
	svc := service.NewWalletService(repo)

	var wg sync.WaitGroup
	errCh := make(chan error, totalRequests)
	sem := make(chan struct{}, maxParallel)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			req := dto.WalletOperationRequest{
				WalletID:      wallet.ID,
				OperationType: "DEPOSIT",
				Amount:        depositAmount,
			}

			_, err := svc.UpdateWalletBalance(ctx, req)
			errCh <- err
		}()
	}

	wg.Wait()
	close(errCh)

	failed := 0
	for err := range errCh {
		if err != nil {
			t.Errorf("failed to deposit: %v", err)
			failed++
		}
	}

	if failed > 0 {
		t.Logf("%d operations failed", failed)
	}

	var updatedWallet model.Wallet
	err := testDB.First(&updatedWallet, "id = ?", wallet.ID).Error
	assert.NoError(t, err)

	expectedBalance := initialBalance + float64(totalRequests)*depositAmount
	t.Logf("Expected: %v, Actual: %v", expectedBalance, updatedWallet.Balance)
	assert.Equal(t, expectedBalance, updatedWallet.Balance, fmt.Sprintf("expected %v, got %v", expectedBalance, updatedWallet.Balance))
}
