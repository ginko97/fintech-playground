package application

import (
	"context"
	"errors"
	"time"

	"github.com/ginko97/fintech-playground/internal/domain"
	"github.com/ginko97/fintech-playground/internal/fraud"
	"github.com/ginko97/fintech-playground/internal/gateway"
	"github.com/ginko97/fintech-playground/internal/repository"
	"github.com/ginko97/fintech-playground/internal/wallet"
	"github.com/google/uuid"
)

// CreateTransactionRequest - DTO from handler
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
	fraudEngine     *fraud.FraudEngine
	walletService   *wallet.WalletService
}

func NewTransactionService(
	repo repository.TransactionRepository,
	sm *TransactionStateMachine,
	wp *WorkerPool,
	fraudEngine *fraud.FraudEngine,
	walletService *wallet.WalletService,
) *TransactionService {
	return &TransactionService{
		repo:            repo,
		stateMachine:    sm,
		workerPool:      wp,
		externalGateway: gateway.NewAdvancedGateway(),
		circuitBreaker:  gateway.NewCircuitBreaker(5),
		fraudEngine:     fraudEngine,
		walletService:   walletService,
	}
}

// Create handles idempotency + business rules + fraud check
func (s *TransactionService) Create(ctx context.Context, req CreateTransactionRequest) (*domain.Transaction, error) {
	// 1. Idempotency check
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

	// 3. Fraud Check
	if err := s.fraudEngine.Check(ctx, tx); err != nil {
		tx.Status = domain.StatusFailed
		s.repo.Create(ctx, tx)
		return tx, err
	}

	// 4. Save transaction
	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, err
	}

	// 5. Submit to background worker
	s.ProcessAsync(tx)

	return tx, nil
}

// ProcessAsync submits to worker pool
func (s *TransactionService) ProcessAsync(tx *domain.Transaction) {
	s.workerPool.Submit(tx)
}

// GetByID example
func (s *TransactionService) GetByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.repo.FindByID(ctx, id)
}
