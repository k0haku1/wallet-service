package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewPostgres(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to postgres database")
	return db
}
