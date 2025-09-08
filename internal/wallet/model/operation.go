package model

import (
	"github.com/google/uuid"
	"time"
)

type Operation struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	WalletID  uuid.UUID `gorm:"type:uuid;not null"`
	Type      string    `gorm:"type:varchar(10);not null"`
	Amount    float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
