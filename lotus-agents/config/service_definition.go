
package config

import “fmt”

// ServiceDefinition describes the microservice to be built.
// This is the single place where you define what you want built —
// all agents receive this and tailor their output accordingly.
type ServiceDefinition struct {
// Name is the short name of the microservice, e.g. “inventory”, “payments”, “notifications”
Name string


// Description is a plain-English summary of what the service does
Description string

// Language is the programming language to use, e.g. "Go", "Python", "Node.js"
Language string

// Entities are the core domain objects, e.g. ["Product", "StockItem", "Warehouse"]
Entities []string

// Operations are the key business operations, e.g. ["Reserve stock", "Process refund"]
Operations []string

// Integrations are external services this microservice talks to
// e.g. ["Order Service (Kafka)", "Payment Gateway (REST)", "Postgres"]
Integrations []string

// ExtraRequirements are any freeform additional requirements
ExtraRequirements []string


}

// Prompt builds a structured prompt string from the service definition,
// used as the base task injected into every agent.
func (s *ServiceDefinition) Prompt() string {
p := fmt.Sprintf(“Microservice Name: %s\n\n”, s.Name)
p += fmt.Sprintf(“Description:\n%s\n\n”, s.Description)
p += fmt.Sprintf(“Language: %s\n\n”, s.Language)


if len(s.Entities) > 0 {
	p += "Core Domain Entities:\n"
	for _, e := range s.Entities {
		p += fmt.Sprintf("  - %s\n", e)
	}
	p += "\n"
}

if len(s.Operations) > 0 {
	p += "Key Business Operations:\n"
	for _, o := range s.Operations {
		p += fmt.Sprintf("  - %s\n", o)
	}
	p += "\n"
}

if len(s.Integrations) > 0 {
	p += "External Integrations:\n"
	for _, i := range s.Integrations {
		p += fmt.Sprintf("  - %s\n", i)
	}
	p += "\n"
}

if len(s.ExtraRequirements) > 0 {
	p += "Additional Requirements:\n"
	for _, r := range s.ExtraRequirements {
		p += fmt.Sprintf("  - %s\n", r)
	}
	p += "\n"
}

return p


}

// –– Example service definitions you can use out of the box ––

// InventoryService returns a ServiceDefinition for an e-commerce inventory microservice
func InventoryService() *ServiceDefinition {
return &ServiceDefinition{
Name:        “inventory”,
Description: “Tracks product stock levels across multiple warehouses for an e-commerce platform. Handles reservations, replenishment, and low-stock alerting.”,
Language:    “Go”,
Entities:    []string{“Product”, “StockItem”, “Warehouse”, “StockMovement”, “Reservation”},
Operations: []string{
“Reserve stock for an order”,
“Release reserved stock on cancellation”,
“Decrement stock on fulfillment”,
“Replenish stock”,
“Aggregate stock across warehouses”,
“Trigger low-stock alerts”,
},
Integrations: []string{
“Order Service (Kafka events)”,
“PostgreSQL (primary store)”,
“Redis (stock level cache)”,
},
ExtraRequirements: []string{
“Prevent negative stock using atomic updates”,
“Full audit trail of all stock movements”,
“Support product variants (size, color)”,
},
}
}

// PaymentsService returns a ServiceDefinition for a payments microservice
func PaymentsService() *ServiceDefinition {
return &ServiceDefinition{
Name:        “payments”,
Description: “Handles payment processing, refunds, and transaction history for an e-commerce platform.”,
Language:    “Go”,
Entities:    []string{“Payment”, “Refund”, “Transaction”, “PaymentMethod”},
Operations: []string{
“Initiate a payment”,
“Confirm payment”,
“Issue a full or partial refund”,
“Retrieve transaction history”,
“Handle webhook callbacks from payment gateway”,
},
Integrations: []string{
“Stripe API (payment gateway)”,
“Order Service (Kafka events)”,
“PostgreSQL (transaction store)”,
},
ExtraRequirements: []string{
“Idempotency keys on all payment requests”,
“PCI-DSS compliant — never store raw card data”,
“Retry logic for transient gateway failures”,
},
}
}

// NotificationsService returns a ServiceDefinition for a notifications microservice
func NotificationsService() *ServiceDefinition {
return &ServiceDefinition{
Name:        “notifications”,
Description: “Sends email, SMS, and push notifications triggered by events across the platform.”,
Language:    “Go”,
Entities:    []string{“Notification”, “Template”, “Recipient”, “DeliveryLog”},
Operations: []string{
“Send email notification”,
“Send SMS notification”,
“Send push notification”,
“Render template with dynamic data”,
“Track delivery status”,
“Manage user notification preferences”,
},
Integrations: []string{
“SendGrid (email)”,
“Twilio (SMS)”,
“Firebase Cloud Messaging (push)”,
“Kafka (consume events from other services)”,
“PostgreSQL (delivery logs, preferences)”,
},
ExtraRequirements: []string{
“Respect user opt-out preferences”,
“Retry failed deliveries with backoff”,
“Deduplicate notifications using idempotency keys”,
},
}
}
