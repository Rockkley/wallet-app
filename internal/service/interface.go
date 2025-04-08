package service

import (
	"context"
	"github.com/google/uuid"
	"wallet-app/internal/domain/entity"
)

type WalletService interface {
	ProcessTransaction(ctx context.Context, id uuid.UUID, op entity.OperationType, amount int64) error
	GetBalance(ctx context.Context, id uuid.UUID) (int64, error)
}
