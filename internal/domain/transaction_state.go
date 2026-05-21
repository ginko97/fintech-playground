package domain

import (
	"fmt"
	"time"
)

type TransactionState struct {
	TransactionID string            `json:"transaction_id"`
	Status        TransactionStatus `json:"status"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Reason        string            `json:"reason,omitempty"`
	Version       int               `json:"version"`
}

func NewTransactionState(txID string) *TransactionState {
	return &TransactionState{
		TransactionID: txID,
		Status:        StatusPending,
		UpdatedAt:     time.Now().UTC(),
		Version:       1,
	}
}

func (s *TransactionState) Transition(newStatus TransactionStatus, reason string) error {
	if !IsValidTransition(s.Status, newStatus) {
		return fmt.Errorf("invalid state transition: %s → %s", s.Status, newStatus)
	}
	s.Status = newStatus
	s.Reason = reason
	s.UpdatedAt = time.Now().UTC()
	s.Version++
	return nil
}

func (s *TransactionState) IsFinal() bool {
	return s.Status == StatusCompleted || s.Status == StatusFailed || s.Status == StatusRefunded
}
