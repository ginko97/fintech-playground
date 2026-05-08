package domain

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
	StatusRefund    TransactionStatus = "refunded"
)

type Transaction struct {
	ID           uuid.UUID         `json:"id"`
	ExternalID   string            `json:"external_id"`
	Amount       int64             `json:"amount"`
	Currency     string            `json:"currency"`
	Type         string            `json:"type"`
	Status       TransactionStatus `json:"status"`
	SourceWallet string            `json:"source_of_wallet,omitempty"`
	DestWallet   string            `json:"dest_wallet,omitempty"`
	MetaData     map[string]any    `json:"meta_data"`
	Created_at   time.Time         `json:"created_at"`
	Updated_at   time.Time         `json:"upated_at"`
	Version      int               `json:"version"`
}

func (t Transaction) IsFinal() bool {
	return t.Status == StatusCompleted || t.Status == StatusFailed || t.Status == StatusPending || t.Status == StatusRefund

}
