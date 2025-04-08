package repository

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wallet-app/internal/domain/entity"
	"wallet-app/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInsufficientFunds = errors.New("low balance")
	ErrInvalidOperation  = errors.New("unknown error")
	ErrTooManyRetries    = errors.New("too many retries")

	pgRetryableErrors = map[string]bool{
		"40001": true, // serialization fail
		"40P01": true, // deadlock
	}
)

type walletRepo struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) storage.WalletRepository {
	return &walletRepo{db: db}
}

func isRetryableError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgRetryableErrors[pgErr.Code]
}

func withTxRetry(
	ctx context.Context,
	db *pgxpool.Pool,
	maxRetries int,
	txFunc func(pgx.Tx) error,
) error {
	var lastError error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * 10 * time.Millisecond
			time.Sleep(backoff)
		}

		tx, err := db.BeginTx(ctx, pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		})
		if err != nil {
			lastError = err
			continue
		}

		if err := txFunc(tx); err != nil {
			_ = tx.Rollback(ctx)
			if isRetryableError(err) {
				lastError = err
				continue
			}

			return err
		}

		if err := tx.Commit(ctx); err != nil {
			if isRetryableError(err) {
				lastError = err
				continue
			}

			return fmt.Errorf("transaction error: %w", err)
		}

		return nil
	}

	if lastError != nil {

		return fmt.Errorf("%w", lastError)
	}
	return ErrTooManyRetries
}

func (r *walletRepo) CreateIfNotExists(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO wallets (id, balance)
		VALUES ($1, 0)
		ON CONFLICT (id) DO NOTHING
	`, id)
	if err != nil {
		return fmt.Errorf("create wallet error: %w", err)
	}
	return nil
}

func (r *walletRepo) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	var balance int64
	err := r.db.QueryRow(ctx, `
		SELECT balance FROM wallets WHERE id = $1
	`, id).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("error getting balance: %w", err)
	}
	return balance, nil
}

func (r *walletRepo) UpdateBalance(ctx context.Context, id uuid.UUID, amount int64, opType entity.OperationType) error {
	return withTxRetry(ctx, r.db, 5, func(tx pgx.Tx) error {
		switch opType {
		case entity.DEPOSIT:
			_, err := tx.Exec(ctx, `
				UPDATE wallets
				SET balance = balance + $1
				WHERE id = $2
			`, amount, id)
			if err != nil {
				return fmt.Errorf("deposit error: %w", err)
			}

		case entity.WITHDRAW:
			res, err := tx.Exec(ctx, `
				UPDATE wallets
				SET balance = balance - $1
				WHERE id = $2 AND balance >= $1
			`, amount, id)
			if err != nil {
				return fmt.Errorf("withdraw error: %w", err)
			}
			if res.RowsAffected() == 0 {
				return ErrInsufficientFunds
			}

		default:
			return ErrInvalidOperation
		}

		return nil
	})
}
