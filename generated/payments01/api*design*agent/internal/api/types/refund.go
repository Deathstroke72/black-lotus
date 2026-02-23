// file: internal/api/types/refund.go
package types

import "time"

// RefundStatus represents the lifecycle state of a refund
type RefundStatus string

const (
	RefundStatusPending   RefundStatus = "pending"
	RefundStatusSucceeded RefundStatus = "succeeded"
	RefundStatusFailed    RefundStatus = "failed"
	RefundStatusCanceled  RefundStatus = "canceled"
)

// RefundReason categorizes why the refund was issued
type RefundReason string

const (
	RefundReasonDuplicate           RefundReason = "duplicate"
	RefundReasonFraudulent          RefundReason = "fraudulent"
	RefundReasonRequestedByCustomer RefundReason = "requested_by_customer"
	RefundReasonOrderCanceled       RefundReason = "order_canceled"
	RefundReasonOther               RefundReason = "other"
)

// CreateRefundRequest for issuing a refund
// @Description Request payload for creating a refund
type CreateRefundRequest struct {
	// IdempotencyKey ensures exactly-once refund processing
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid4"`

	// Amount to refund; if omitted, refunds the full remaining amount
	Amount *int64 `json:"amount,omitempty" validate:"omitempty,gt=0"`

	// Reason for the refund (required for audit trail)
	Reason RefundReason `json:"reason" validate:"required,oneof=duplicate fraudulent requested_by_customer order_canceled other"`

	// ReasonDescription provides additional context
	ReasonDescription string `json:"reason_description,omitempty" validate:"omitempty,max=500"`

	// Metadata for additional context
	Metadata map[string]string `json:"metadata,omitempty" validate:"omitempty,dive,keys,max=40,endkeys,max=500"`
}

// RefundResponse represents a refund resource
type RefundResponse struct {
	ID                string            `json:"id"`
	PaymentID         string            `json:"payment_id"`
	Amount            Money             `json:"amount"`
	Status            RefundStatus      `json:"status"`
	Reason            RefundReason      `json:"reason"`
	ReasonDescription string            `json:"reason_description,omitempty"`
	StripeRefundID    string            `json:"stripe_refund_id,omitempty"`
	FailureReason     string            `json:"failure_reason,omitempty"`
	Metadata          map[string]string `json:"metadata,omitempty"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// ListRefundsResponse for paginated refund list
type ListRefundsResponse struct {
	Refunds    []RefundResponse `json:"refunds"`
	Pagination PaginationMeta   `json:"pagination"`
}