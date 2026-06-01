package fraud

import (
	"context"
	"errors"

	"github.com/ginko97/fintech-playground/internal/domain"
)

type FraudEngine struct {
	rules []FraudRule
}

type FraudRule interface {
	Check(ctx context.Context, tx *domain.Transaction) error
}

func NewFraudEngine() *FraudEngine {
	return &FraudEngine{
		rules: []FraudRule{
			&AmountLimitRule{},
			&VelocityRule{},
		},
	}
}

func (e *FraudEngine) Check(ctx context.Context, tx *domain.Transaction) error {
	for _, rule := range e.rules {
		if err := rule.Check(ctx, tx); err != nil {
			return err
		}
	}
	return nil
}

// Example Rules

type AmountLimitRule struct{}

func (r *AmountLimitRule) Check(_ context.Context, tx *domain.Transaction) error {
	if tx.Amount > 50000000 { // 50 million IDR limit
		return errors.New("transaction amount exceeds fraud limit")
	}
	return nil
}

type VelocityRule struct{}

func (r *VelocityRule) Check(_ context.Context, tx *domain.Transaction) error {
	// In real system we would check Redis for recent transactions
	// For now, simple simulation
	if tx.Amount > 10000000 && tx.Currency == "IDR" {
		return errors.New("velocity check failed - suspicious high amount")
	}
	return nil
}
