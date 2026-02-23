// file: internal/api/types/transaction.go
package types

import "time"

// TransactionType categorizes the transaction
type TransactionType string

const (
	TransactionTypePayment TransactionType = "payment"
	TransactionTypeRefund  TransactionType = "refund"
	TransactionTypeChargeback TransactionType = "chargeback"
	TransactionTypeFee     TransactionType = "fee"
)

// TransactionResponse represents a transaction record
type TransactionResponse struct {
	ID            string            `json:"id"`
	PaymentID     string            `json:"payment_id"`
	RefundID      string            `json:"refund_id,omitempty"`
	Type          TransactionType   `json:"type"`
	Amount        Money             `json:"amount"`
	CustomerID    string            `json:"customer_id"`
	OrderID       string            `json:"order_id"`
	Status        string            `json:"status"`
	Description   string            `json:"description,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}

// ListTransactionsRequest for filtering transactions
type ListTransactionsRequest struct {
	// CustomerID filters by customer
	CustomerID string `json:"customer_id,omitempty" validate:"omitempty,max=64"`

	// OrderID filters by order
	OrderID string `json:"order_id,omitempty" validate:"omitempty,max=64"`

	// Type filters by transaction type
	Type TransactionType `json:"type,omitempty" validate:"omitempty,oneof=payment refund chargeback fee"`

	// StartDate filters transactions from this date
	StartDate *time.Time `json:"start_date,omitempty"`

	// EndDate filters transactions until this date
	EndDate *time.Time `json:"end_date,omitempty"`

	// Pagination
	Limit  int    `json:"limit,omitempty" validate:"omitempty,min=1,max=100"`
	Cursor string `json:"cursor,omitempty"`
}

// ListTransactionsResponse for paginated transaction list
type ListTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Pagination   PaginationMeta        `json:"pagination"`
}