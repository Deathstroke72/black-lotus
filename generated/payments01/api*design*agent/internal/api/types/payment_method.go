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