package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"wallet-service/internal/config"
	"wallet-service/internal/db"
	"wallet-service/internal/wallet/handler"
	"wallet-service/internal/wallet/handler/middleware"
	"wallet-service/internal/wallet/repository"
	"wallet-service/internal/wallet/service"
)

func main() {
	cfg := config.LoadConfig()
	gormDb := db.NewPostgres(cfg.DBUrl)

	var (
		walletRepo    repository.WalletRepository = repository.NewWalletRepository(gormDb)
		walletService *service.WalletService      = service.NewWalletService(walletRepo)
		walletHandler *handler.WalletHandler      = handler.NewWalletHandler(walletService)
	)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	api := app.Group("/api/v1")
	api.Get("/wallets/:wallet_uuid", walletHandler.GetWalletBalance)
	api.Post("/wallet", walletHandler.UpdateWalletBalance)

	log.Printf("Starting server on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
