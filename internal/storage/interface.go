package storage

import (
	"context"
	"github.com/google/uuid"
	"wallet-app/internal/domain/entity"
)

type WalletRepository interface {
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
	UpdateBalance(ctx context.Context, id uuid.UUID, amount int64, opType entity.OperationType) error
	CreateIfNotExists(ctx context.Context, id uuid.UUID) error
}
