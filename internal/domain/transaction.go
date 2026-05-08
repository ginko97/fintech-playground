package domain

import (
	"time"

	"github.com/gofrs/uuid/v5"
)

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
	StatusRefunded  TransactionStatus = "refunded"
)

type Transaction struct {
	ID           uuid.UUID         `json:"id" db:"id"`
	ExternalID   string            `json:"external_id" db:"external_id"`
	Amount       int64             `json:"amount" db:"amount"`
	Currency     string            `json:"currency" db:"currency"`
	Type         string            `json:"type" db:"type"`
	Status       TransactionStatus `json:"status" db:"status"`
	SourceWallet string            `json:"source_of_wallet,omitempty" db:"source_wallet"`
	DestWallet   string            `json:"dest_wallet,omitempty" db:"dest_wallet"`
	MetaData     map[string]any    `json:"meta_data" db:"metadata"`
	Created_at   time.Time         `json:"created_at" db:"created_at"`
	Updated_at   time.Time         `json:"upated_at" db:"updated_at"`
	Version      int               `json:"-" db:"version"`
}

func (t Transaction) IsFinal() bool {
	return t.Status == StatusCompleted ||
		t.Status == StatusFailed ||
		t.Status == StatusPending ||
		t.Status == StatusRefunded
}
