# black-lotus
# Microservice Agent Pipeline

A generic multi-agent system in Go using the Claude API that can design and generate any microservice from a simple `ServiceDefinition` struct. All agents automatically adapt their system prompts, code generation, and tests to whatever service you define.

## Architecture

```
ServiceDefinition  ←── you define this
       ↓
  Orchestrator (Pipeline)
       ↓
  ┌────────────────────────────────────────────┐
  │  API Design Agent     → routes, schemas    │
  │        ↓ (output passed as context)        │
  │  Backend & DB Agent   → service + schema   │
  │        ↓                                   │
  │  Messaging Agent      → Kafka events       │
  │        ↓                                   │
  │  Testing & Security   → tests + auth       │
  └────────────────────────────────────────────┘
       ↓
  generated/<service-name>/
```

## Quick Start

```bash
export ANTHROPIC_API_KEY=your_key_here
go mod tidy
go run main.go
```

## Defining a Microservice

Edit `main.go` and swap in your own `ServiceDefinition`:

```go
svc := &config.ServiceDefinition{
    Name:        "orders",
    Description: "Manages the full order lifecycle for an e-commerce platform",
    Language:    "Go",
    Entities:    []string{"Order", "OrderItem", "ShippingAddress"},
    Operations: []string{
        "Place an order",
        "Cancel an order",
        "Update order status",
        "Track shipment",
    },
    Integrations: []string{
        "Inventory Service (Kafka)",
        "Payment Service (REST)",
        "PostgreSQL (primary store)",
    },
    ExtraRequirements: []string{
        "Idempotency on order placement",
        "Full order history audit trail",
    },
}
```

Three example services are provided out of the box in `config/service_definition.go`:

- `config.InventoryService()` — stock management across warehouses
- `config.PaymentsService()` — payment processing and refunds
- `config.NotificationsService()` — email, SMS, and push notifications

## Output Structure

```
generated/
└── <service-name>/
    ├── README.md
    ├── api_design_agent/
    │   ├── output.md        # Full agent output
    │   └── *.go             # Router, handlers, schemas
    ├── backend_and_database_agent/
    │   ├── output.md
    │   ├── *.go             # Service layer, repositories
    │   └── *.sql            # Migrations
    ├── messaging_and_events_agent/
    │   ├── output.md
    │   └── *.go             # Kafka producers, consumers, events
    └── testing_and_security_agent/
        ├── output.md
        └── *.go             # Tests, JWT middleware, rate limiter
```

## Adding a New Agent

1. Create `agents/my_agent.go` implementing the `Agent` interface:

```go
type MyAgent struct{ *BaseAgent }

func NewMyAgent(cfg *config.Config, svc *config.ServiceDefinition) *MyAgent {
    return &MyAgent{
        BaseAgent: NewBaseAgentForService(cfg, "My Agent", svc, responsibilities, outputFormat),
    }
}

func (a *MyAgent) Run(ctx context.Context, svc *config.ServiceDefinition, agentCtx map[string]string) (*AgentResult, error) {
    // Build prompt using svc.Prompt() + agentCtx from prior agents
    // Call a.Chat(ctx, messages)
    // Return AgentResult with parsed artifacts
}
```

1. Register it in `orchestrator/pipeline.go`:

```go
func(svc *config.ServiceDefinition) agents.Agent { return agents.NewMyAgent(cfg, svc) },
```

1. Add a context key so downstream agents can access its output.```