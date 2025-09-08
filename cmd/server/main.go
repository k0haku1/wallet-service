package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"wallet-service/internal/config"
)

func main() {
	cfg := config.LoadConfig()
	//gormDb := db.NewPostgres(cfg.DBUrl)

	app := fiber.New()

	log.Printf("Starting server on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
