package repository

import (
	"context"
	"errors"

	"github.com/ginko97/fintech-playground/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrWalletNotFound = errors.New("wallet not found")

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(ctx context.Context, wallet *domain.Wallet) error {
	const query = `
		INSERT INTO wallets (id, user_id, balance, currency, status, created_at, updated_at, version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(ctx, query,
		wallet.ID, wallet.UserID, wallet.Balance, wallet.Currency,
		wallet.Status, wallet.CreatedAt, wallet.UpdatedAt, wallet.Version,
	)
	return err
}

func (r *WalletRepository) FindByID(ctx context.Context, id string) (*domain.Wallet, error) {
	const query = `
		SELECT id, user_id, balance, currency, status, created_at, updated_at, version
		FROM wallets WHERE id = $1
	`

	var w domain.Wallet
	err := r.db.QueryRow(ctx, query, id).Scan(
		&w.ID, &w.UserID, &w.Balance, &w.Currency,
		&w.Status, &w.CreatedAt, &w.UpdatedAt, &w.Version,
	)

	if err != nil {
		return nil, ErrWalletNotFound
	}
	return &w, nil
}

// UpdateBalance with optimistic locking
func (r *WalletRepository) UpdateBalance(ctx context.Context, walletID string, amount int64, expectedVersion int) error {
	const query = `
		UPDATE wallets 
		SET balance = balance + $1, 
		    updated_at = NOW(), 
		    version = version + 1
		WHERE id = $2 AND version = $3
	`

	result, err := r.db.Exec(ctx, query, amount, walletID, expectedVersion)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("balance update failed - version mismatch or wallet not found")
	}

	return nil
}
