package domain

import "time"

type TransactionStatus string

const (
	StatusPending    TransactionStatus = "pending"
	StatusProcessing TransactionStatus = "processing"
	StatusCompleted  TransactionStatus = "completed"
	StatusFailed     TransactionStatus = "failed"
	StatusRefunded   TransactionStatus = "refunded"
)

// ValidTransitions defines allowed state changes
var ValidTransitions = map[TransactionStatus][]TransactionStatus{
	StatusPending:    {StatusProcessing, StatusFailed},
	StatusProcessing: {StatusCompleted, StatusFailed, StatusRefunded},
	StatusCompleted:  {},
	StatusFailed:     {StatusRefunded},
	StatusRefunded:   {},
}

func IsValidTransition(from, to TransactionStatus) bool {
	for _, allowed := range ValidTransitions[from] {
		if allowed == to {
			return true
		}
	}
	return false
}

type Transaction struct {
	ID             string            `json:"id" db:"id"`
	IdempotencyKey string            `json:"idempotency_key" db:"idempotency_key"`
	RequestID      string            `json:"request_id" db:"request_id"`
	SourceWalletID string            `json:"source_wallet_id" db:"source_wallet_id"`
	DestWalletID   string            `json:"dest_wallet_id" db:"dest_wallet_id"`
	Amount         int64             `json:"amount" db:"amount"`
	Currency       string            `json:"currency" db:"currency"`
	Status         TransactionStatus `json:"status" db:"status"`
	Description    string            `json:"description,omitempty" db:"description"`
	CreatedAt      time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at" db:"updated_at"`
	LedgerBalance  int64             `json:"ledger_balance,omitempty" db:"ledger_balance"`
	Version        int               `json:"version" db:"version"`
}

func (t Transaction) IsFinal() bool {
	return t.Status == StatusCompleted ||
		t.Status == StatusFailed ||
		t.Status == StatusRefunded
}
