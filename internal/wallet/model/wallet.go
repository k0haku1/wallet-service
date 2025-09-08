package model

import (
	"github.com/google/uuid"
	"time"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	Balance   float64   `gorm:"type:decimal(10,2);default:0"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
