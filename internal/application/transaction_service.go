package application

import (
	"context"
	"errors"
	"time"

	"github.com/ginko97/fintech-playground/internal/domain"
	"github.com/ginko97/fintech-playground/internal/gateway"
	"github.com/ginko97/fintech-playground/internal/repository"
	"github.com/google/uuid"
)

// CreateTransactionRequest - DTO from handler (input model)
type CreateTransactionRequest struct {
	IdempotencyKey string `json:"idempotency_key" validate:"required"`
	SourceWalletID string `json:"source_wallet_id" validate:"required"`
	DestWalletID   string `json:"dest_wallet_id" validate:"required"`
	Amount         int64  `json:"amount" validate:"required,gt=0"`
	Currency       string `json:"currency" validate:"required"`
	Description    string `json:"description,omitempty"`
}

type TransactionService struct {
	repo            repository.TransactionRepository
	stateMachine    *TransactionStateMachine
	workerPool      *WorkerPool
	externalGateway *gateway.AdvancedGateway
	circuitBreaker  *gateway.CircuitBreaker
}

func NewTransactionService(
	repo repository.TransactionRepository,
	sm *TransactionStateMachine,
	wp *WorkerPool,
) *TransactionService {
	return &TransactionService{
		repo:            repo,
		stateMachine:    sm,
		workerPool:      wp,
		externalGateway: gateway.NewAdvancedGateway(),
		circuitBreaker:  gateway.NewCircuitBreaker(5),
	}
}

// Create handles idempotency + business rules
func (s *TransactionService) Create(ctx context.Context, req CreateTransactionRequest) (*domain.Transaction, error) {
	// 1. Idempotency check (safe retry)
	if tx, err := s.repo.FindByIdempotencyKey(ctx, req.IdempotencyKey); err == nil {
		return tx, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// 2. Build domain entity
	tx := &domain.Transaction{
		ID:             uuid.Must(uuid.NewV7()).String(),
		IdempotencyKey: req.IdempotencyKey,
		RequestID:      uuid.Must(uuid.NewV7()).String(),
		SourceWalletID: req.SourceWalletID,
		DestWalletID:   req.DestWalletID,
		Amount:         req.Amount,
		Currency:       req.Currency,
		Status:         domain.StatusPending,
		Description:    req.Description,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		Version:        1,
	}

	// 3. Persist
	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, err
	}

	// Call external PSP with circuit breaker
	err := s.circuitBreaker.Execute(ctx, func() error {
		_, err := s.externalGateway.ProcessPayment(ctx, tx)
		return err
	})

	if err != nil {
		tx.Status = domain.StatusFailed
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, err
	}

	s.ProcessAsync(tx) // send to worker pool

	return tx, nil
}

// GetByID example
func (s *TransactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.repo.FindByID(ctx, id)
}

// ProcessAsync submits transaction to worker pool for background processing
func (s *TransactionService) ProcessAsync(tx *domain.Transaction) {
	s.workerPool.Submit(tx) // we will inject workerPool next
}
