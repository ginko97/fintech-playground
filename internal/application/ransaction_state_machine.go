package application

import (
	"context"
	"errors"
	"time"

	"github.com/ginko97/fintech-playground/internal/domain"
	"github.com/ginko97/fintech-playground/internal/infrastructure"
)

type TransactionStateMachine struct {
	redis *infrastructure.RedisClient
}

func NewTransactionStateMachine(redis *infrastructure.RedisClient) *TransactionStateMachine {
	return &TransactionStateMachine{redis: redis}
}

// UpdateState atomically updates transaction state with Redis lock + validation
func (sm *TransactionStateMachine) UpdateState(ctx context.Context, txID string, newStatus domain.TransactionStatus, reason string) error {
	lockKey := "lock:tx:" + txID
	lock, err := sm.redis.GetClient().SetNX(ctx, lockKey, "locked", 10*time.Second).Result()
	if err != nil || !lock {
		return errors.New("transaction is being processed by another worker")
	}
	defer sm.redis.GetClient().Del(ctx, lockKey) // release lock

	// TODO: Later load current state from Redis or DB
	state := domain.NewTransactionState(txID)

	if err := state.Transition(newStatus, reason); err != nil {
		return err
	}

	// Save to Redis (for fast access) + update DB in real flow
	// For now we just validate transition
	return nil
}

// GetState example
func (sm *TransactionStateMachine) GetState(ctx context.Context, txID string) (*domain.TransactionState, error) {
	// In real implementation we would get from Redis cache + DB fallback
	return domain.NewTransactionState(txID), nil
}
