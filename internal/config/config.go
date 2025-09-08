package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port  string
	DBUrl string
}

func LoadConfig() *Config {
	if err := godotenv.Load("config.env"); err != nil {
		log.Println("config.env not found, using system env")
	}

	return &Config{
		Port:  getEnv("APP_PORT", "8080"),
		DBUrl: getEnv("DB_URL", "postgres://wallet_user:password@postgres:5432/wallet_db?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
