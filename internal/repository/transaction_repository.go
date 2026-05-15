package repository

import (
	"context"
	"errors"

	"github.com/ginko97/fintech-playground/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrNotFound = errors.New("transaction not found")
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	FindByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error)
	FindByID(ctx context.Context, id string) (*domain.Transaction, error)
}

type transactionRepository struct {
	db *pgxpool.Pool
}

func NewTransactionRepository(db *pgxpool.Pool) TransactionRepository {
	return &transactionRepository{db: db}
}

// Create inserts with idempotency protection (DB constraint handles race)
func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	const query = `
		INSERT INTO transactions (
			id, idempotency_key, request_id, source_wallet_id, 
			dest_wallet_id, amount, currency, status, description,
			created_at, updated_at, version
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
		)
		ON CONFLICT (idempotency_key) DO NOTHING;  -- safety net
	`

	_, err := r.db.Exec(ctx, query,
		tx.ID,
		tx.IdempotencyKey,
		tx.RequestID,
		tx.SourceWalletID,
		tx.DestWalletID,
		tx.Amount,
		tx.Currency,
		tx.Status,
		tx.Description,
		tx.CreatedAt,
		tx.UpdatedAt,
		1, // initial version
	)

	if err != nil {
		return err
	}
	return nil
}

// FindByIdempotencyKey for idempotency guard
func (r *transactionRepository) FindByIdempotencyKey(ctx context.Context, key string) (*domain.Transaction, error) {
	const query = `
		SELECT id, idempotency_key, request_id, source_wallet_id, dest_wallet_id,
		       amount, currency, status, description, created_at, updated_at, version
		FROM transactions
		WHERE idempotency_key = $1
	`

	var t domain.Transaction
	err := r.db.QueryRow(ctx, query, key).Scan(
		&t.ID,
		&t.IdempotencyKey,
		&t.RequestID,
		&t.SourceWalletID,
		&t.DestWalletID,
		&t.Amount,
		&t.Currency,
		&t.Status,
		&t.Description,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

// FindByID example read
func (r *transactionRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	const query = `
		SELECT id, idempotency_key, request_id, source_wallet_id, dest_wallet_id,
		       amount, currency, status, description, created_at, updated_at, version
		FROM transactions
		WHERE id = $1
	`

	var t domain.Transaction
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.IdempotencyKey, &t.RequestID, &t.SourceWalletID,
		&t.DestWalletID, &t.Amount, &t.Currency, &t.Status,
		&t.Description, &t.CreatedAt, &t.UpdatedAt, &t.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}
