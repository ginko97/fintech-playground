package gateway

import (
	"context"
	"time"

	"github.com/ginko97/fintech-playground/internal/domain"
)

type ExternalGateway interface {
	ProcessPayment(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)
}

type APIGateway struct{}

func NewApiGateway() *APIGateway {
	return &APIGateway{}
}

// ProcessPayment simulates external PSP call
func (g *APIGateway) ProcessPayment(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	// Simulate network latency
	select {
	case <-time.After(400 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Simulate success / failure
	if tx.Amount%13 == 0 {
		tx.Status = domain.StatusFailed
		return tx, nil
	}

	tx.Status = domain.StatusCompleted
	return tx, nil
}
