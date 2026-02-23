# API Design Agent Output

# Payments Microservice API Design

## 1. Complete API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/v1/payments` | Initiate a new payment |
| `POST` | `/api/v1/payments/{payment_id}/confirm` | Confirm a pending payment |
| `GET` | `/api/v1/payments/{payment_id}` | Retrieve payment details |
| `POST` | `/api/v1/payments/{payment_id}/refunds` | Issue a full or partial refund |
| `GET` | `/api/v1/payments/{payment_id}/refunds` | List refunds for a payment |
| `GET` | `/api/v1/refunds/{refund_id}` | Retrieve refund details |
| `GET` | `/api/v1/transactions` | Retrieve transaction history with filtering |
| `GET` | `/api/v1/transactions/{transaction_id}` | Retrieve single transaction |
| `POST` | `/api/v1/webhooks/stripe` | Handle Stripe webhook callbacks |
| `GET` | `/api/v1/payment-methods/{customer_id}` | List saved payment methods for customer |
| `DELETE` | `/api/v1/payment-methods/{payment_method_id}` | Remove a saved payment method |
| `GET` | `/health` | Health check endpoint |
| `GET` | `/ready` | Readiness probe |

---

## 2. Go Structs for Request/Response Payloads

```go
// file: internal/api/types/common.go
package types

import "time"

// Money represents a monetary amount with currency
// Currency follows ISO 4217 (e.g., USD, EUR, GBP)
// Amount is in the smallest currency unit (cents for USD)
type Money struct {
	Amount   int64  `json:"amount" validate:"required,gt=0"`
	Currency string `json:"currency" validate:"required,len=3,uppercase"`
}

// PaginationParams for list endpoints
type PaginationParams struct {
	Limit  int    `json:"limit" validate:"omitempty,min=1,max=100"`
	Offset int    `json:"offset" validate:"omitempty,min=0"`
	Cursor string `json:"cursor,omitempty"`
}

// PaginationMeta returned in list responses
type PaginationMeta struct {
	Total      int64  `json:"total"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}

// Timestamp wrapper for consistent JSON formatting
type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format(time.RFC3339) + `"`), nil
}
```

```go
// file: internal/api/types/payment.go
package types

import "time"

// PaymentStatus represents the lifecycle state of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusRequiresConfirmation PaymentStatus = "requires_confirmation"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusSucceeded PaymentStatus = "succeeded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCanceled  PaymentStatus = "canceled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusPartiallyRefunded PaymentStatus = "partially_refunded"
)

// InitiatePaymentRequest represents the request to create a new payment
// @Description Request payload for initiating a payment
type InitiatePaymentRequest struct {
	// IdempotencyKey ensures exactly-once payment processing
	// Must be unique per payment attempt, typically a UUID
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid4"`

	// OrderID links this payment to an order in the Order Service
	OrderID string `json:"order_id" validate:"required,min=1,max=64"`

	// CustomerID identifies the customer making the payment
	CustomerID string `json:"customer_id" validate:"required,min=1,max=64"`

	// Amount specifies the payment amount and currency
	Amount Money `json:"amount" validate:"required"`

	// PaymentMethodID references a tokenized payment method from Stripe
	// Never contains raw card data (PCI-DSS compliance)
	PaymentMethodID string `json:"payment_method_id" validate:"required,startswith=pm_"`

	// Description appears on customer's statement
	Description string `json:"description,omitempty" validate:"omitempty,max=500"`

	// Metadata for additional context (order details, etc.)
	Metadata map[string]string `json:"metadata,omitempty" validate:"omitempty,dive,keys,max=40,endkeys,max=500"`

	// CaptureMethod determines when funds are captured
	// "automatic" captures immediately, "manual" requires confirmation
	CaptureMethod string `json:"capture_method,omitempty" validate:"omitempty,oneof=automatic manual"`

	// ReturnURL for redirect-based payment methods
	ReturnURL string `json:"return_url,omitempty" validate:"omitempty,url"`
}

// PaymentResponse represents a payment resource
// @Description Payment resource returned by the API
type PaymentResponse struct {
	ID              string            `json:"id"`
	OrderID         string            `json:"order_id"`
	CustomerID      string            `json:"customer_id"`
	Amount          Money             `json:"amount"`
	Status          PaymentStatus     `json:"status"`
	Description     string            `json:"description,omitempty"`
	PaymentMethodID string            `json:"payment_method_id"`
	
	// StripePaymentIntentID for reference (not sensitive)
	StripePaymentIntentID string `json:"stripe_payment_intent_id,omitempty"`
	
	// RefundedAmount tracks total refunded
	RefundedAmount  Money             `json:"refunded_amount,omitempty"`
	
	// FailureCode and FailureMessage for failed payments
	FailureCode    string            `json:"failure_code,omitempty"`
	FailureMessage string            `json:"failure_message,omitempty"`
	
	Metadata       map[string]string `json:"metadata,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	ConfirmedAt    *time.Time        `json:"confirmed_at,omitempty"`
}

// ConfirmPaymentRequest for confirming a payment requiring manual capture
type ConfirmPaymentRequest struct {
	// IdempotencyKey for confirmation request
	IdempotencyKey string `json:"idempotency_key" validate:"required,uuid4"`

	// AmountToCapture allows partial capture (optional)
	// If omitted, captures full authorized amount
	AmountToCapture *int64 `json:"amount_to_capture,omitempty" validate:"omitempty,gt=0"`
}

// ConfirmPaymentResponse after successful confirmation
type ConfirmPaymentResponse struct {
	Payment PaymentResponse `json:"payment"`
}
```

```go
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
```

```go
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
```

```go
// file: internal/api/types/webhook.go
package types

// StripeWebhookPayload represents the incoming Stripe webhook
// Note: This is parsed from Stripe's webhook library, not directly
type StripeWebhookPayload struct {
	// Raw payload is verified using Stripe signature
	// Actual event type and data extracted after verification
}

// WebhookResponse returned to Stripe
type WebhookResponse struct {
	Received bool `json:"received"`
}
```

```go
// file: internal/api/types/payment_method.go
package types

import "time"

// PaymentMethodType categorizes the payment method
type PaymentMethodType string

const (
	PaymentMethodTypeCard       PaymentMethodType = "card"
	PaymentMethodTypeBankAccount PaymentMethodType = "bank_account"
	PaymentMethodTypeWallet     PaymentMethodType = "wallet"
)

// CardBrand represents the card network
type CardBrand string

const (
	CardBrandVisa       CardBrand = "visa"
	CardBrandMastercard CardBrand = "mastercard"
	CardBrandAmex       CardBrand = "amex"
	CardBrandDiscover   CardBrand = "discover"
	CardBrandUnknown    CardBrand = "unknown"
)

// PaymentMethodResponse represents a tokenized payment method
// Note: We NEVER store or return raw card numbers (PCI-DSS compliance)
type PaymentMethodResponse struct {
	ID         string            `json:"id"`
	CustomerID string            `json:"customer_id"`
	Type       PaymentMethodType `json:"type"`
	
	// Card details (only non-sensitive data)
	Card *CardDetails `json:"card,omitempty"`
	
	// IsDefault indicates if this is the customer's default payment method
	IsDefault  bool              `json:"is_default"`
	CreatedAt  time.Time         `json:"created_at"`
}

// CardDetails contains non-sensitive card information
// PCI-DSS Compliant: Only stores last 4 digits and expiry
type CardDetails struct {
	Brand       CardBrand `json:"brand"`
	Last4       string    `json:"last4"`
	ExpMonth    int       `json:"exp_month"`
	ExpYear     int       `json:"exp_year"`
	Fingerprint string    `json:"fingerprint"` // For duplicate detection
}

// ListPaymentMethodsResponse for customer's payment methods
type ListPaymentMethodsResponse struct {
	PaymentMethods []PaymentMethodResponse `json:"payment_methods"`
}
```

---

## 3. Router Setup Code

```go
// file: internal/api/router/router.go
package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	
	"payments-service/internal/api/handlers"
	custommiddleware "payments-service/internal/api/middleware"
)

// Config holds router configuration
type Config struct {
	Environment     string
	AllowedOrigins  []string
	StripeWebhookSecret string
}

// NewRouter creates and configures the chi router with all payment endpoints
func NewRouter(cfg Config, h *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(custommiddleware.StructuredLogger())
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	
	// CORS configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Idempotency-Key"},
		ExposedHeaders:   []string{"X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Health check endpoints (no auth required)
	r.Get("/health", h.HealthCheck)
	r.Get("/ready", h.ReadinessCheck)

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		// Apply authentication middleware to all API routes
		r.Use(custommiddleware.Authenticate)
		r.Use(custommiddleware.RateLimiter(100, time.Minute))

		// Payment endpoints
		r.Route("/payments", func(r chi.Router) {
			// POST /api/v1/payments - Initiate a new payment
			// @Summary Initiate a new payment
			// @Description Creates a new payment intent with the specified amount and payment method
			// @Tags payments
			// @Accept json
			// @Produce json
			// @Param X-Idempotency-Key header string true "Idempotency key for request deduplication"
			// @Param request body types.InitiatePaymentRequest true "Payment initiation request"
			// @Success 201 {object} types.PaymentResponse
			// @Failure 400 {object} types.ErrorResponse "Invalid request"
			// @Failure 409 {object} types.ErrorResponse "Duplicate idempotency key with different payload"
			// @Failure 422 {object} types.ErrorResponse "Payment processing failed"
			// @Router /payments [post]
			r.With(custommiddleware.IdempotencyCheck).Post("/", h.InitiatePayment)

			r.Route("/{paymentID}", func(r chi.Router) {
				r.Use(custommiddleware.PaymentCtx) // Load payment into context
				
				// GET /api/v1/payments/{paymentID} - Get payment details
				// @Summary Get payment details
				// @Description Retrieves the details of an existing payment
				// @Tags payments
				// @Produce json
				// @Param paymentID path string true "Payment ID"
				// @Success 200 {object} types.PaymentResponse
				// @Failure 404 {object} types.ErrorResponse "Payment not found"
				// @Router /payments/{paymentID} [get]
				r.Get("/", h.GetPayment)

				// POST /api/v1/payments/{paymentID}/confirm - Confirm a payment
				// @Summary Confirm a pending payment
				// @Description Confirms and captures a payment that requires manual confirmation
				// @Tags payments
				// @Accept json
				// @Produce json
				// @Param paymentID path string true "Payment ID"
				// @Param X-Idempotency-Key header string true "Idempotency key"
				// @Param request body types.ConfirmPaymentRequest true "Confirmation request"
				// @Success 200 {object} types.ConfirmPaymentResponse
				// @Failure 400 {object} types.ErrorResponse "Payment cannot be confirmed"
				// @Failure 404 {object} types.ErrorResponse "Payment not found"
				// @Router /payments/{paymentID}/confirm [post]
				r.With(custommiddleware.IdempotencyCheck).Post("/confirm", h.ConfirmPayment)

				// Refund endpoints nested under payment
				r.Route("/refunds", func(r chi.Router) {
					// POST /api/v1/payments/{paymentID}/refunds - Issue a refund
					// @Summary Issue a refund
					// @Description Creates a full or partial refund for a payment
					// @Tags refunds
					// @Accept json
					// @Produce json
					// @Param paymentID path string true "Payment ID"
					// @Param X-Idempotency-Key header string true "Idempotency key"
					// @Param request body types.CreateRefundRequest true "Refund request"
					// @Success 201 {object} types.RefundResponse
					// @Failure 400 {object} types.ErrorResponse "Invalid refund request"
					// @Failure 404 {object} types.ErrorResponse "Payment not found"
					// @Failure 422 {object} types.ErrorResponse "Refund exceeds available amount"
					// @Router /payments/{paymentID}/refunds [post]
					r.With(custommiddleware.IdempotencyCheck).Post("/", h.CreateRefund)

					// GET /api/v1/payments/{paymentID}/refunds - List refunds for payment
					// @Summary List refunds for a payment
					// @Description Retrieves all refunds associated with a payment
					// @Tags refunds
					// @Produce json
					// @Param paymentID path string true "Payment ID"
					// @Param limit query int false "Max results (default 20, max 100)"
					// @Param cursor query string false "Pagination cursor"
					// @Success 200 {object} types.ListRefundsResponse
					// @Failure 404 {object} types.ErrorResponse "Payment not found"
					// @Router /payments/{paymentID}/refunds [get]
					r.Get("/", h.ListRefunds)
				})
			})
		})

		// Standalone refund endpoint for direct access
		r.Route("/refunds", func(r chi.Router) {
			// GET /api/v1/refunds/{refundID} - Get refund details
			// @Summary Get refund details
			// @Description Retrieves the details of a specific refund
			// @Tags refunds
			// @Produce json
			// @Param refundID path string true "Refund ID"
			// @Success 200 {object} types.RefundResponse
			// @Failure 404 {object} types.ErrorResponse "Refund not found"
			// @Router /refunds/{refundID} [get]
			r.Get("/{refundID}", h.GetRefund)
		})

		// Transaction history endpoints
		r.Route("/transactions", func(r chi.Router) {
			// GET /api/v1/transactions - List transactions
			// @Summary List transactions
			// @Description Retrieves transaction history with optional filtering
			// @Tags transactions
			// @Produce json
			// @Param customer_id query string false "Filter by customer ID"
			// @Param order_id query string false "Filter by order ID"
			// @Param type query string false "Filter by type (payment, refund, chargeback)"
			// @Param start_date query string false "Start date (RFC3339)"
			// @Param end_date query string false "End date (RFC3339)"
			// @Param limit query int false "Max results (default 20, max 100)"
			// @Param cursor query string false "Pagination cursor"
			// @Success 200 {object} types.ListTransactionsResponse
			// @Router /transactions [get]
			r.Get("/", h.ListTransactions)

			// GET /api/v1/transactions/{transactionID} - Get transaction details
			// @Summary Get transaction details
			// @Description Retrieves a specific transaction record
			// @Tags transactions
			// @Produce json
			// @Param transactionID path string true "Transaction ID"
			// @Success 200 {object} types.TransactionResponse
			// @Failure 404 {object} types.ErrorResponse "Transaction not found"
			// @Router /transactions/{transactionID} [get]
			r.Get("/{transactionID}", h.GetTransaction)
		})

		// Payment methods endpoints
		r.Route("/payment-methods", func(r chi.Router) {
			// GET /api/v1/payment-methods/{customerID} - List customer's payment methods
			// @Summary List payment methods
			// @Description Retrieves saved payment methods for a customer
			// @Tags payment-methods
			// @Produce json
			// @Param customerID path string true "Customer ID"
			// @Success 200 {object} types.ListPaymentMethodsResponse
			// @Router /payment-methods/{customerID} [get]
			r.Get("/{customerID}", h.ListPaymentMethods)

			// DELETE /api/v1/payment-methods/{paymentMethodID} - Remove payment method
			// @Summary Delete a payment method
			// @Description Removes a saved payment method
			// @Tags payment-methods
			// @Param paymentMethodID path string true "Payment Method ID"
			// @Success 204 "Successfully deleted"
			// @Failure 404 {object} types.ErrorResponse "Payment method not found"
			// @Router /payment-methods/{paymentMethodID} [delete]
			r.Delete("/{paymentMethodID}", h.DeletePaymentMethod)
		})
	})

	// Webhook endpoints (separate from main API, different auth)
	r.Route("/api/v1/webhooks", func(r chi.Router) {
		// POST /api/v1/webhooks/stripe - Stripe webhook handler
		// @Summary Handle Stripe webhooks
		// @Description Receives and processes webhook events from Stripe
		// @Tags webhooks
		// @Accept json
		// @Produce json
		// @Param Stripe-Signature header string true "Stripe webhook signature"
		// @Success 200 {object} types.WebhookResponse
		// @Failure 400 {object} types.ErrorResponse "Invalid signature"
		// @Router /webhooks/stripe [post]
		r.With(custommiddleware.StripeWebhookVerifier(cfg.StripeWebhookSecret)).
			Post("/stripe", h.HandleStripeWebhook)
	})

	return r
}
```

---

## 4. Standardized Error Response Format

```go
// file: internal/api/types/errors.go
package types

import (
	"encoding/json"
	"net/http"
)

// ErrorCode represents machine-readable error codes
type ErrorCode string

const (
	// General errors
	ErrCodeInvalidRequest     ErrorCode = "INVALID_REQUEST"
	ErrCodeValidationFailed   ErrorCode = "VALIDATION_FAILED"
	ErrCodeResourceNotFound   ErrorCode = "RESOURCE_NOT_FOUND"
	ErrCodeInternalError      ErrorCode = "INTERNAL_ERROR"
	ErrCodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden          ErrorCode = "FORBIDDEN"
	ErrCodeRateLimitExceeded  ErrorCode = "RATE_LIMIT_EXCEEDED"
	
	// Idempotency errors
	ErrCodeIdempotencyKeyMissing  ErrorCode = "IDEMPOTENCY_KEY_MISSING"
	ErrCodeIdempotencyKeyConflict ErrorCode = "IDEMPOTENCY_KEY_CONFLICT"
	
	// Payment-specific errors
	ErrCodePaymentFailed          ErrorCode = "PAYMENT_FAILED"
	ErrCodePaymentDeclined        ErrorCode = "PAYMENT_DECLINED"
	ErrCodeInsufficientFunds      ErrorCode = "INSUFFICIENT_FUNDS"
	ErrCodeInvalidPaymentMethod   ErrorCode = "INVALID_PAYMENT_METHOD"
	ErrCodePaymentAlreadyConfirmed ErrorCode = "PAYMENT_ALREADY_CONFIRMED"
	ErrCodePaymentNotConfirmable  ErrorCode = "PAYMENT_NOT_CONFIRMABLE"
	ErrCodePaymentExpired         ErrorCode = "PAYMENT_EXPIRED"
	
	// Refund-specific errors
	ErrCodeRefundFailed           ErrorCode = "REFUND_FAILED"
	ErrCodeRefundExceedsPayment   ErrorCode = "REFUND_EXCEEDS_PAYMENT"
	ErrCodePaymentNotRefundable   ErrorCode = "PAYMENT_NOT_REFUNDABLE"
	ErrCodeRefundAlreadyProcessed ErrorCode = "REFUND_ALREADY_PROCESSED"
	
	// Gateway errors
	ErrCodeGatewayTimeout    ErrorCode = "GATEWAY_TIMEOUT"
	ErrCodeGatewayUnavailable ErrorCode = "GATEWAY_UNAVAILABLE"
)

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains the error information
type ErrorDetail struct {
	// Code is a machine-readable error code
	Code ErrorCode `json:"code"`
	
	// Message is a human-readable error description
	Message string `json:"message"`
	
	// Details provides additional context about the error
	Details map[string]interface{} `json:"details,omitempty"`
	
	// ValidationErrors lists field-level validation failures
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
	
	// RequestID for tracing/debugging
	RequestID string `json:"request_id,omitempty"`
	
	// RetryAfter indicates when to retry (for rate limiting)
	RetryAfter *int `json:"retry_after,omitempty"`
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Common error constructors

// NewErrorResponse creates a standard error response
func NewErrorResponse(code ErrorCode, message string, requestID string) *ErrorResponse {
	return &ErrorResponse{
		Error: ErrorDetail{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	}
}

// WithDetails adds details to the error response
func (e *ErrorResponse) WithDetails(details map[string]interface{}) *ErrorResponse {
	e.Error.Details = details
	return e
}

// WithValidationErrors adds validation errors to the response
func (e *ErrorResponse) WithValidationErrors(errors []ValidationError) *ErrorResponse {
	e.Error.ValidationErrors = errors
	return e
}

// WriteJSON writes the error response to the http.ResponseWriter
func (e *ErrorResponse) WriteJSON(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(e)
}

// HTTPStatusForErrorCode maps error codes to HTTP status codes
func HTTPStatusForErrorCode(code ErrorCode) int {
	switch code {
	case ErrCodeInvalidRequest, ErrCodeValidationFailed, 
	     ErrCodeIdempotencyKeyMissing:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeResourceNotFound:
		return http.StatusNotFound
	case ErrCodeIdempotencyKeyConflict:
		return http.StatusConflict
	case ErrCodePaymentFailed, ErrCodePaymentDeclined, ErrCodeInsufficientFunds,
	     ErrCodeInvalidPaymentMethod, ErrCodeRefundFailed, 
	     ErrCodeRefundExceedsPayment, ErrCodePaymentNotRefundable,
	     ErrCodePaymentNotConfirmable, ErrCodePaymentAlreadyConfirmed:
		return http.StatusUnprocessableEntity
	case ErrCodeRateLimitExceeded:
		return http.StatusTooManyRequests
	case ErrCodeGatewayTimeout:
		return http.StatusGatewayTimeout
	case ErrCodeGatewayUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
```

```go
// file: internal/api/handlers/errors.go
package handlers

import (
	"net/http"
	
	"payments-service/internal/api/types"
	"github.com/go-chi/chi/v5/middleware"
)

// respondError sends a standardized error response
func respondError(w http.ResponseWriter, r *http.Request, code types.ErrorCode, message string) {
	requestID := middleware.GetReqID(r.Context())
	errResp := types.NewErrorResponse(code, message, requestID)
	errResp.WriteJSON(w, types.HTTPStatusFor