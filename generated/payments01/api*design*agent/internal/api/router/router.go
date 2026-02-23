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