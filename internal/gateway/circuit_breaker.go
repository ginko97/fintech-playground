package gateway

import (
	"context"
	"errors"
	"sync"
	"time"
)

type CircuitBreaker struct {
	failures    int
	threshold   int
	lastFailure time.Time
	state       string // CLOSED, OPEN, HALF_OPEN
	mu          sync.Mutex
}

func NewCircuitBreaker(threshold int) *CircuitBreaker {
	return &CircuitBreaker{
		threshold: threshold,
		state:     "CLOSED",
	}
}

func (cb *CircuitBreaker) Execute(ctx context.Context, fn func() error) error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.state == "OPEN" {
		if time.Since(cb.lastFailure) > 30*time.Second {
			cb.state = "HALF_OPEN"
		} else {
			return errors.New("circuit breaker is OPEN")
		}
	}

	err := fn()
	if err != nil {
		cb.failures++
		cb.lastFailure = time.Now()
		if cb.failures >= cb.threshold {
			cb.state = "OPEN"
		}
		return err
	}

	// Success
	cb.failures = 0
	cb.state = "CLOSED"
	return nil
}
