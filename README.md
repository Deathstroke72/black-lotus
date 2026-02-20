# black-lotus

# Inventory Microservice Agent Pipeline

A multi-agent system built in Go using the Claude API that autonomously designs and generates a production-ready inventory microservice for an e-commerce platform.

## Architecture

```
Orchestrator
├── API Design Agent        → REST endpoints, schemas, router setup
├── Backend & DB Agent      → Service logic, PostgreSQL schema, repositories
├── Messaging & Events Agent → Kafka producers/consumers, domain events
└── Testing & Security Agent → Tests, JWT auth, RBAC, rate limiting
```

Each agent receives the output of previous agents as context, ensuring the generated code is coherent and consistent across layers.

## Prerequisites

- Go 1.21+
- An Anthropic API key

## Setup

```bash
# Install dependencies
go mod tidy

# Set your API key
export ANTHROPIC_API_KEY=your_key_here

# Run the full pipeline
go run main.go

# Run with a custom task
go run main.go "Build an inventory service for a fashion retailer with size/color variants" ./my-output
```

## Output Structure

After running, the `generated/` directory will contain:

```
generated/
├── README.md                          # Pipeline summary
├── api_design_agent/
│   ├── output.md                      # Full agent output
│   └── *.go                           # Router, handlers, schemas
├── backend_and_database_agent/
│   ├── output.md
│   ├── *.go                           # Service layer, repositories
│   └── *.sql                          # Migrations
├── messaging_and_events_agent/
│   ├── output.md
│   └── *.go                           # Kafka producers, consumers, events
└── testing_and_security_agent/
    ├── output.md
    └── *.go                           # Tests, auth middleware, rate limiter
```

## Agents

### 1. API Design Agent

Designs the REST API contract including all endpoints, request/response structs, HTTP status codes, and router configuration using Go’s `net/http` or `chi`.

### 2. Backend & Database Agent

Implements the service layer and data access layer. Generates PostgreSQL schemas with proper indexing, repository implementations using `pgx`, and concurrency-safe stock operations using `SELECT FOR UPDATE`.

### 3. Messaging & Events Agent

Implements event-driven communication via Kafka. Produces domain events (StockReserved, StockDepleted, etc.) and consumes events from dependent services (orders, warehouse). Implements the transactional outbox pattern for reliability.

### 4. Testing & Security Agent

Writes table-driven unit tests, integration tests using `testcontainers-go`, concurrency tests, and implements JWT middleware with role-based access control and rate limiting.

## Adding New Agents

1. Create a new file in `agents/` implementing the `Agent` interface
1. Add it to the pipeline in `orchestrator/pipeline.go`
1. Add a context key so downstream agents can use its output

```go
// agents/my_new_agent.go
type MyNewAgent struct { *BaseAgent }

func (a *MyNewAgent) Run(ctx context.Context, task string, context map[string]string) (*AgentResult, error) {
    // Call Claude with a specialized system prompt
}
```