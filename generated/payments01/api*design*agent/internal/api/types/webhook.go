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