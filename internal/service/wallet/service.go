package service

import (
	"context"
	"wallet-app/internal/domain/entity"
	"wallet-app/internal/service"
	"wallet-app/internal/storage"

	"github.com/google/uuid"
)

type walletService struct {
	repo storage.WalletRepository
}

func NewWalletService(repo storage.WalletRepository) service.WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) ProcessTransaction(ctx context.Context, id uuid.UUID, op entity.OperationType, amount int64) error {
	if err := s.repo.CreateIfNotExists(ctx, id); err != nil {
		return err
	}
	return s.repo.UpdateBalance(ctx, id, amount, op)
}

func (s *walletService) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	return s.repo.GetBalance(ctx, id)
}
