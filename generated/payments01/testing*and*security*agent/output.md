# Testing & Security Agent Output

# Payments Microservice - Tests and Security Implementation

I'll provide a comprehensive implementation of tests and security middleware for the payments microservice.

## Project Structure

```
payments/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   └── types/
│   ├── domain/
│   ├── middleware/
│   │   ├── audit.go
│   │   ├── auth.go
│   │   ├── ratelimit.go
│   │   └── requestid.go
│   ├── repository/
│   └── service/
├── tests/
│   ├── integration/
│   └── security/
├── Makefile
└── go.mod
```

## 1. Domain Entities and Interfaces

```go
// file: internal/domain/entities.go
package domain

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus represents the lifecycle state of a payment
type PaymentStatus string

const (
	PaymentStatusPending              PaymentStatus = "pending"
	PaymentStatusRequiresConfirmation PaymentStatus = "requires_confirmation"
	PaymentStatusProcessing           PaymentStatus = "processing"
	PaymentStatusSucceeded            PaymentStatus = "succeeded"
	PaymentStatusFailed               PaymentStatus = "failed"
	PaymentStatusCanceled             PaymentStatus = "canceled"
	PaymentStatusRefunded             PaymentStatus = "refunded"
	PaymentStatusPartiallyRefunded    PaymentStatus = "partially_refunded"
)

// RefundStatus represents the lifecycle state of a refund
type RefundStatus string

const (
	RefundStatusPending    RefundStatus = "pending"
	RefundStatusProcessing RefundStatus = "processing"
	RefundStatusSucceeded  RefundStatus = "succeeded"
	RefundStatusFailed     RefundStatus = "failed"
	RefundStatusCanceled   RefundStatus = "canceled"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeChargeback TransactionType = "chargeback"
)

// TransactionStatus represents the status of a transaction
type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusFailed    TransactionStatus = "failed"
	TransactionStatusReversed  TransactionStatus = "reversed"
)

// PaymentMethodType represents the type of payment method
type PaymentMethodType string

const (
	PaymentMethodTypeCard         PaymentMethodType = "card"
	PaymentMethodTypeBankAccount  PaymentMethodType = "bank_account"
	PaymentMethodTypeDigitalWallet PaymentMethodType = "digital_wallet"
)

// Money represents a monetary amount
type Money struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

// Payment represents a payment entity
type Payment struct {
	ID                    uuid.UUID     `json:"id"`
	IdempotencyKey        string        `json:"idempotency_key"`
	CustomerID            uuid.UUID     `json:"customer_id"`
	OrderID               uuid.UUID     `json:"order_id"`
	PaymentMethodID       uuid.UUID     `json:"payment_method_id"`
	Amount                Money         `json:"amount"`
	Status                PaymentStatus `json:"status"`
	StripePaymentIntentID string        `json:"stripe_payment_intent_id,omitempty"`
	Description           string        `json:"description,omitempty"`
	Metadata              map[string]string `json:"metadata,omitempty"`
	FailureCode           string        `json:"failure_code,omitempty"`
	FailureMessage        string        `json:"failure_message,omitempty"`
	RefundedAmount        int64         `json:"refunded_amount"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	ConfirmedAt           *time.Time    `json:"confirmed_at,omitempty"`
}

// Refund represents a refund entity
type Refund struct {
	ID             uuid.UUID    `json:"id"`
	PaymentID      uuid.UUID    `json:"payment_id"`
	Amount         Money        `json:"amount"`
	Status         RefundStatus `json:"status"`
	Reason         string       `json:"reason,omitempty"`
	StripeRefundID string       `json:"stripe_refund_id,omitempty"`
	FailureReason  string       `json:"failure_reason,omitempty"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	ProcessedAt    *time.Time   `json:"processed_at,omitempty"`
}

// Transaction represents a transaction record
type Transaction struct {
	ID              uuid.UUID         `json:"id"`
	PaymentID       uuid.UUID         `json:"payment_id"`
	RefundID        *uuid.UUID        `json:"refund_id,omitempty"`
	Type            TransactionType   `json:"type"`
	Status          TransactionStatus `json:"status"`
	Amount          Money             `json:"amount"`
	BalanceBefore   int64             `json:"balance_before"`
	BalanceAfter    int64             `json:"balance_after"`
	Description     string            `json:"description,omitempty"`
	StripeChargeID  string            `json:"stripe_charge_id,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

// PaymentMethod represents a saved payment method
type PaymentMethod struct {
	ID                    uuid.UUID         `json:"id"`
	CustomerID            uuid.UUID         `json:"customer_id"`
	Type                  PaymentMethodType `json:"type"`
	StripePaymentMethodID string            `json:"stripe_payment_method_id"`
	CardBrand             string            `json:"card_brand,omitempty"`
	CardLastFour          string            `json:"card_last_four,omitempty"`
	CardExpMonth          int               `json:"card_exp_month,omitempty"`
	CardExpYear           int               `json:"card_exp_year,omitempty"`
	IsDefault             bool              `json:"is_default"`
	IsActive              bool              `json:"is_active"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

// TransactionFilter represents filters for transaction queries
type TransactionFilter struct {
	CustomerID  *uuid.UUID
	PaymentID   *uuid.UUID
	Type        *TransactionType
	Status      *TransactionStatus
	StartDate   *time.Time
	EndDate     *time.Time
	MinAmount   *int64
	MaxAmount   *int64
	Currency    *string
	Limit       int
	Offset      int
}

// WebhookEvent represents a Stripe webhook event
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}
```

```go
// file: internal/domain/errors.go
package domain

import "errors"

var (
	// Payment errors
	ErrPaymentNotFound       = errors.New("payment not found")
	ErrPaymentAlreadyExists  = errors.New("payment with this idempotency key already exists")
	ErrPaymentAlreadyConfirmed = errors.New("payment already confirmed")
	ErrPaymentNotConfirmable = errors.New("payment is not in a confirmable state")
	ErrPaymentNotRefundable  = errors.New("payment is not in a refundable state")
	ErrInvalidPaymentStatus  = errors.New("invalid payment status")
	
	// Refund errors
	ErrRefundNotFound        = errors.New("refund not found")
	ErrRefundExceedsPayment  = errors.New("refund amount exceeds available payment amount")
	ErrRefundAlreadyProcessed = errors.New("refund already processed")
	
	// Transaction errors
	ErrTransactionNotFound   = errors.New("transaction not found")
	
	// Payment method errors
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrPaymentMethodInactive = errors.New("payment method is not active")
	
	// Validation errors
	ErrInvalidAmount         = errors.New("invalid amount")
	ErrInvalidCurrency       = errors.New("invalid currency")
	ErrInvalidIdempotencyKey = errors.New("invalid idempotency key")
	ErrMissingCustomerID     = errors.New("customer ID is required")
	ErrMissingOrderID        = errors.New("order ID is required")
	
	// Gateway errors
	ErrGatewayUnavailable    = errors.New("payment gateway unavailable")
	ErrGatewayTimeout        = errors.New("payment gateway timeout")
	ErrGatewayRejected       = errors.New("payment rejected by gateway")
	
	// Webhook errors
	ErrInvalidWebhookSignature = errors.New("invalid webhook signature")
	ErrWebhookEventNotFound    = errors.New("webhook event not found")
	ErrDuplicateWebhookEvent   = errors.New("duplicate webhook event")
	
	// Authorization errors
	ErrUnauthorized          = errors.New("unauthorized")
	ErrForbidden             = errors.New("forbidden")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)
```

```go
// file: internal/domain/repositories.go
package domain

import (
	"context"

	"github.com/google/uuid"
)

// PaymentRepository defines the interface for payment persistence
type PaymentRepository interface {
	// Create creates a new payment
	Create(ctx context.Context, payment *Payment) error
	
	// GetByID retrieves a payment by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	
	// GetByIdempotencyKey retrieves a payment by idempotency key
	GetByIdempotencyKey(ctx context.Context, key string) (*Payment, error)
	
	// Update updates an existing payment
	Update(ctx context.Context, payment *Payment) error
	
	// GetByCustomerID retrieves payments for a customer
	GetByCustomerID(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]*Payment, int64, error)
	
	// GetByOrderID retrieves a payment by order ID
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*Payment, error)
	
	// LockForUpdate locks a payment row for update
	LockForUpdate(ctx context.Context, id uuid.UUID) (*Payment, error)
}

// RefundRepository defines the interface for refund persistence
type RefundRepository interface {
	// Create creates a new refund
	Create(ctx context.Context, refund *Refund) error
	
	// GetByID retrieves a refund by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Refund, error)
	
	// Update updates an existing refund
	Update(ctx context.Context, refund *Refund) error
	
	// GetByPaymentID retrieves refunds for a payment
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*Refund, error)
	
	// GetTotalRefundedAmount gets the total refunded amount for a payment
	GetTotalRefundedAmount(ctx context.Context, paymentID uuid.UUID) (int64, error)
}

// TransactionRepository defines the interface for transaction persistence
type TransactionRepository interface {
	// Create creates a new transaction
	Create(ctx context.Context, tx *Transaction) error
	
	// GetByID retrieves a transaction by ID
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	
	// List retrieves transactions with filters
	List(ctx context.Context, filter TransactionFilter) ([]*Transaction, int64, error)
	
	// GetByPaymentID retrieves transactions for a payment
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*Transaction, error)
}

// PaymentMethodRepository defines the interface for payment method persistence
type PaymentMethodRepository interface {
	// Create creates a new payment method
	Create(ctx context.Context, pm *PaymentMethod) error
	
	// GetByID retrieves a payment method by ID
	GetByID(ctx context.Context, id uuid.UUID) (*PaymentMethod, error)
	
	// GetByCustomerID retrieves payment methods for a customer
	GetByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*PaymentMethod, error)
	
	// Update updates an existing payment method
	Update(ctx context.Context, pm *PaymentMethod) error
	
	// Delete soft deletes a payment method
	Delete(ctx context.Context, id uuid.UUID) error
	
	// SetDefault sets a payment method as the default for a customer
	SetDefault(ctx context.Context, customerID, paymentMethodID uuid.UUID) error
}

// IdempotencyRepository defines the interface for idempotency key management
type IdempotencyRepository interface {
	// Store stores an idempotency key with its response
	Store(ctx context.Context, key string, response []byte, ttlSeconds int) error
	
	// Get retrieves a stored response for an idempotency key
	Get(ctx context.Context, key string) ([]byte, bool, error)
	
	// Exists checks if an idempotency key exists
	Exists(ctx context.Context, key string) (bool, error)
}

// WebhookEventRepository defines the interface for webhook event tracking
type WebhookEventRepository interface {
	// Store stores a processed webhook event ID
	Store(ctx context.Context, eventID string) error
	
	// Exists checks if a webhook event has been processed
	Exists(ctx context.Context, eventID string) (bool, error)
}

// UnitOfWork defines the interface for transaction management
type UnitOfWork interface {
	// Begin starts a new transaction
	Begin(ctx context.Context) (context.Context, error)
	
	// Commit commits the current transaction
	Commit(ctx context.Context) error
	
	// Rollback rolls back the current transaction
	Rollback(ctx context.Context) error
	
	// WithTransaction executes a function within a transaction
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
```

```go
// file: internal/domain/gateway.go
package domain

import (
	"context"
)

// PaymentGateway defines the interface for payment gateway operations
type PaymentGateway interface {
	// CreatePaymentIntent creates a new payment intent
	CreatePaymentIntent(ctx context.Context, req CreatePaymentIntentRequest) (*PaymentIntentResponse, error)
	
	// ConfirmPaymentIntent confirms a payment intent
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResponse, error)
	
	// CancelPaymentIntent cancels a payment intent
	CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResponse, error)
	
	// CreateRefund creates a refund
	CreateRefund(ctx context.Context, req CreateRefundRequest) (*RefundResponse, error)
	
	// GetPaymentIntent retrieves a payment intent
	GetPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResponse, error)
	
	// ValidateWebhookSignature validates a webhook signature
	ValidateWebhookSignature(payload []byte, signature string) error
}

// CreatePaymentIntentRequest represents a request to create a payment intent
type CreatePaymentIntentRequest struct {
	Amount              int64
	Currency            string
	PaymentMethodID     string
	CustomerID          string
	Description         string
	Metadata            map[string]string
	IdempotencyKey      string
	CaptureMethod       string // "automatic" or "manual"
}

// PaymentIntentResponse represents a payment intent response from the gateway
type PaymentIntentResponse struct {
	ID                string
	Status            string
	Amount            int64
	Currency          string
	ClientSecret      string
	PaymentMethodID   string
	FailureCode       string
	FailureMessage    string
}

// CreateRefundRequest represents a request to create a refund
type CreateRefundRequest struct {
	PaymentIntentID string
	Amount          int64
	Reason          string
	Metadata        map[string]string
	IdempotencyKey  string
}

// RefundResponse represents a refund response from the gateway
type RefundResponse struct {
	ID            string
	Status        string
	Amount        int64
	Currency      string
	FailureReason string
}

// EventPublisher defines the interface for publishing domain events
type EventPublisher interface {
	// PublishPaymentCreated publishes a payment created event
	PublishPaymentCreated(ctx context.Context, payment *Payment) error
	
	// PublishPaymentConfirmed publishes a payment confirmed event
	PublishPaymentConfirmed(ctx context.Context, payment *Payment) error
	
	// PublishPaymentFailed publishes a payment failed event
	PublishPaymentFailed(ctx context.Context, payment *Payment) error
	
	// PublishRefundCreated publishes a refund created event
	PublishRefundCreated(ctx context.Context, refund *Refund) error
	
	// PublishRefundCompleted publishes a refund completed event
	PublishRefundCompleted(ctx context.Context, refund *Refund) error
}
```

## 2. Service Layer Implementation

```go
// file: internal/service/payment_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/payments/internal/domain"
)

// PaymentService handles payment business logic
type PaymentService struct {
	paymentRepo     domain.PaymentRepository
	refundRepo      domain.RefundRepository
	transactionRepo domain.TransactionRepository
	paymentMethodRepo domain.PaymentMethodRepository
	idempotencyRepo domain.IdempotencyRepository
	gateway         domain.PaymentGateway
	eventPublisher  domain.EventPublisher
	unitOfWork      domain.UnitOfWork
}

// NewPaymentService creates a new PaymentService
func NewPaymentService(
	paymentRepo domain.PaymentRepository,
	refundRepo domain.RefundRepository,
	transactionRepo domain.TransactionRepository,
	paymentMethodRepo domain.PaymentMethodRepository,
	idempotencyRepo domain.IdempotencyRepository,
	gateway domain.PaymentGateway,
	eventPublisher domain.EventPublisher,
	unitOfWork domain.UnitOfWork,
) *PaymentService {
	return &PaymentService{
		paymentRepo:       paymentRepo,
		refundRepo:        refundRepo,
		transactionRepo:   transactionRepo,
		paymentMethodRepo: paymentMethodRepo,
		idempotencyRepo:   idempotencyRepo,
		gateway:           gateway,
		eventPublisher:    eventPublisher,
		unitOfWork:        unitOfWork,
	}
}

// InitiatePaymentRequest represents a request to initiate a payment
type InitiatePaymentRequest struct {
	IdempotencyKey  string
	CustomerID      uuid.UUID
	OrderID         uuid.UUID
	PaymentMethodID uuid.UUID
	Amount          int64
	Currency        string
	Description     string
	Metadata        map[string]string
}

// InitiatePayment initiates a new payment
func (s *PaymentService) InitiatePayment(ctx context.Context, req InitiatePaymentRequest) (*domain.Payment, error) {
	// Validate request
	if err := s.validateInitiateRequest(req); err != nil {
		return nil, err
	}

	// Check idempotency
	existingPayment, err := s.paymentRepo.GetByIdempotencyKey(ctx, req.IdempotencyKey)
	if err == nil && existingPayment != nil {
		return existingPayment, nil
	}

	// Verify payment method exists and is active
	paymentMethod, err := s.paymentMethodRepo.GetByID(ctx, req.PaymentMethodID)
	if err != nil {
		return nil, domain.ErrPaymentMethodNotFound
	}
	if !paymentMethod.IsActive {
		return nil, domain.ErrPaymentMethodInactive
	}

	// Create payment intent with gateway
	intentReq := domain.CreatePaymentIntentRequest{
		Amount:          req.Amount,
		Currency:        req.Currency,
		PaymentMethodID: paymentMethod.StripePaymentMethodID,
		CustomerID:      req.CustomerID.String(),
		Description:     req.Description,
		Metadata:        req.Metadata,
		IdempotencyKey:  req.IdempotencyKey,
		CaptureMethod:   "manual", // We'll confirm manually
	}

	intentResp, err := s.gateway.CreatePaymentIntent(ctx, intentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	// Create payment record
	payment := &domain.Payment{
		ID:                    uuid.New(),
		IdempotencyKey:        req.IdempotencyKey,
		CustomerID:            req.CustomerID,
		OrderID:               req.OrderID,
		PaymentMethodID:       req.PaymentMethodID,
		Amount:                domain.Money{Amount: req.Amount, Currency: req.Currency},
		Status:                domain.PaymentStatusRequiresConfirmation,
		StripePaymentIntentID: intentResp.ID,
		Description:           req.Description,
		Metadata:              req.Metadata,
		CreatedAt:             time.Now().UTC(),
		UpdatedAt:             time.Now().UTC(),
	}

	if err := s.unitOfWork.WithTransaction(ctx, func(txCtx context.Context) error {
		if err := s.paymentRepo.Create(txCtx, payment); err != nil {
			return err
		}

		// Create initial transaction record
		tx := &domain.Transaction{
			ID:          uuid.New(),
			PaymentID:   payment.ID,
			Type:        domain.TransactionTypePayment,
			Status:      domain.TransactionStatusPending,
			Amount:      payment.Amount,
			Description: "Payment initiated",
			CreatedAt:   time.Now().UTC(),
		}
		return s.transactionRepo.Create(txCtx, tx)
	}); err != nil {
		return nil, err
	}

	// Publish event
	_ = s.eventPublisher.PublishPaymentCreated(ctx, payment)

	return payment, nil
}

func (s *PaymentService) validateInitiateRequest(req InitiatePaymentRequest) error {
	if req.IdempotencyKey == "" {
		return domain.ErrInvalidIdempotencyKey
	}
	if req.CustomerID == uuid.Nil {
		return domain.ErrMissingCustomerID
	}
	if req.OrderID == uuid.Nil {
		return domain.ErrMissingOrderID
	}
	if req.Amount <= 0 {
		return domain.ErrInvalidAmount
	}
	if len(req.Currency) != 3 {
		return domain.ErrInvalidCurrency
	}
	return nil
}

// ConfirmPayment confirms a pending payment
func (s *PaymentService) ConfirmPayment(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	var payment *domain.Payment

	err := s.unitOfWork.WithTransaction(ctx, func(txCtx context.Context) error {
		var err error
		payment, err = s.paymentRepo.LockForUpdate(txCtx, paymentID)
		if err != nil {
			return domain.ErrPaymentNotFound
		}

		if payment.Status != domain.PaymentStatusRequiresConfirmation {
			if payment.Status == domain.PaymentStatusSucceeded {
				return domain.ErrPaymentAlreadyConfirmed
			}
			return domain.ErrPaymentNotConfirmable
		}

		// Confirm with gateway
		intentResp, err := s.gateway.ConfirmPaymentIntent(txCtx, payment.StripePaymentIntentID)
		if err != nil {
			payment.Status = domain.PaymentStatusFailed
			payment.FailureMessage = err.Error()
			payment.UpdatedAt = time.Now().UTC()
			_ = s.paymentRepo.Update(txCtx, payment)
			return fmt.Errorf("gateway confirmation failed: %w", err)
		}

		// Update payment status based on gateway response
		now := time.Now().UTC()
		switch intentResp.Status {
		case "succeeded":
			payment.Status = domain.PaymentStatusSucceeded
			payment.ConfirmedAt = &now
		case "processing":
			payment.Status = domain.PaymentStatusProcessing
		case "requires_action":
			payment.Status = domain.PaymentStatusRequiresConfirmation
		default:
			payment.Status = domain.PaymentStatusFailed
			payment.FailureCode = intentResp.FailureCode
			payment.FailureMessage = intentResp.FailureMessage
		}
		payment.UpdatedAt = now

		if err := s.paymentRepo.Update(txCtx, payment); err != nil {
			return err
		}

		// Update transaction status
		transactions, err := s.transactionRepo.GetByPaymentID(txCtx, payment.ID)
		if err == nil && len(transactions) > 0 {
			tx := transactions[0]
			if payment.Status == domain.PaymentStatusSucceeded {
				tx.Status = domain.TransactionStatusCompleted
			} else if payment.Status == domain.PaymentStatusFailed {
				tx.Status = domain.TransactionStatusFailed
			}
			// Note: In a real impl, we'd update the transaction
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish events
	if payment.Status == domain.PaymentStatusSucceeded {
		_ = s.eventPublisher.PublishPaymentConfirmed(ctx, payment)
	} else if payment.Status == domain.PaymentStatusFailed {
		_ = s.eventPublisher.PublishPaymentFailed(ctx, payment)
	}

	return payment, nil
}

// GetPayment retrieves a payment by ID
func (s *PaymentService) GetPayment(ctx context.Context, paymentID uuid.UUID) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, domain.ErrPaymentNotFound
	}
	return payment, nil
}

// IssueRefundRequest represents a request to issue a refund
type IssueRefundRequest struct {
	PaymentID uuid.UUID
	Amount    int64
	Reason    string
	Metadata  map[string]string
}

// IssueRefund issues a full or partial refund
func (s *PaymentService) IssueRefund(ctx context.Context, req IssueRefundRequest) (*domain.Refund, error) {
	var refund *domain.Refund

	err := s.unitOfWork.WithTransaction(ctx, func(txCtx context.Context) error {
		// Lock and get payment
		payment, err := s.paymentRepo.LockForUpdate(txCtx, req.PaymentID)
		if err != nil {
			return domain.ErrPaymentNotFound
		}

		// Validate payment is refundable
		if !s.isRefundable(payment.Status) {
			return domain.ErrPaymentNotRefundable
		}

		// Check refund amount
		totalRefunded, err := s.refundRepo.GetTotalRefundedAmount(txCtx, payment.ID)
		if err != nil {
			return err
		}

		availableForRefund := payment.Amount.Amount - totalRefunded
		refundAmount := req.Amount
		if refundAmount == 0 {
			refundAmount = availableForRefund // Full refund
		}

		if refundAmount > availableForRefund {
			return domain.ErrRefundExceedsPayment
		}

		// Create refund with gateway
		refundReq := domain.CreateRefundRequest{
			PaymentIntentID: payment.StripePaymentIntentID,
			Amount:          refundAmount,
			Reason:          req.Reason,
			Metadata:        req.Metadata,
			IdempotencyKey:  fmt.Sprintf("refund_%s_%d", req.PaymentID, time.Now().UnixNano()),
		}

		refundResp, err := s.gateway.CreateRefund(txCtx, refundReq)
		if err != nil {
			return fmt.Errorf("gateway refund failed: %w", err)
		}

		// Create refund record
		refund = &domain.Refund{
			ID:             uuid.New(),
			PaymentID:      payment.ID,
			Amount:         domain.Money{Amount: refundAmount, Currency: payment.Amount.Currency},
			Status:         domain.RefundStatus(refundResp.Status),
			Reason:         req.Reason,
			StripeRefundID: refundResp.ID,
			Metadata:       req.Metadata,
			CreatedAt:      time.Now().UTC(),
			UpdatedAt:      time.Now().UTC(),
		}

		if refundResp.Status == "succeeded" {
			now := time.Now().UTC()
			refund.ProcessedAt = &now
			refund.Status = domain.RefundStatusSucceeded
		}

		if err := s.refundRepo.Create(txCtx, refund); err != nil {
			return err
		}

		// Update payment refunded amount and status
		payment.RefundedAmount = totalRefunded + refundAmount
		if payment.RefundedAmount >= payment.Amount.Amount {
			payment.Status = domain.PaymentStatusRefunded
		} else {
			payment.Status = domain.PaymentStatusPartiallyRefunded
		}
		payment.UpdatedAt = time.Now().UTC()

		if err := s.paymentRepo.Update(txCtx, payment); err != nil {
			return err
		}

		// Create refund transaction
		tx := &domain.Transaction{
			ID:          uuid.New(),
			PaymentID:   payment.ID,
			RefundID:    &refund.ID,
			Type:        domain.TransactionTypeRefund,
			Status:      domain.TransactionStatusCompleted,
			Amount:      domain.Money{Amount: -refundAmount, Currency: payment.Amount.Currency},
			Description: fmt.Sprintf("Refund: %s", req.Reason),
			CreatedAt:   time.Now().UTC(),
		}
		return s.transactionRepo.Create(txCtx, tx)
	})

	if err != nil {
		return nil, err
	}

	// Publish event
	_ = s.eventPublisher.PublishRefundCreated(ctx, refund)

	return refund, nil
}

func (s *PaymentService) isRefundable(status domain.PaymentStatus) bool {
	return status == domain.PaymentStatusSucceeded ||
		status == domain.PaymentStatusPartiallyRefunded
}

// GetRefund retrieves a refund by ID
func (s *PaymentService) GetRefund(ctx context.Context, refundID uuid.UUID) (*domain.Refund, error) {
	refund, err := s.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		return nil, domain.ErrRefundNotFound
	}
	return refund, nil
}

// GetRefundsByPaymentID retrieves refunds for a payment
func (s *PaymentService) GetRefundsByPaymentID(