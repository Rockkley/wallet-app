package entity

import (
	"github.com/google/uuid"
)

type Wallet struct {
	ID      uuid.UUID `gorm:"type:uuid;primary_key"`
	Balance int64
}
