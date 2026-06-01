package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	Balance   int64     `json:"balance" db:"balance"`
	Currency  string    `json:"currency" db:"currency"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	Version   int       `json:"version" db:"version"`
}

func NewWallet(userID, currency string) *Wallet {
	return &Wallet{
		ID:        "wall_" + uuid.New().String()[:8], // simple unique ID
		UserID:    userID,
		Balance:   0,
		Currency:  currency,
		Status:    "active",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Version:   1,
	}
}

// AddBalance with optimistic locking
func (w *Wallet) AddBalance(amount int64, expectedVersion int) error {
	if w.Version != expectedVersion {
		return fmt.Errorf("version mismatch - wallet was updated by another transaction")
	}
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	w.Balance += amount
	w.UpdatedAt = time.Now().UTC()
	w.Version++
	return nil
}
