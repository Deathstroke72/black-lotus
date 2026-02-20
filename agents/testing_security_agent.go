package agents

import (
“context”
“fmt”

```
"github.com/anthropics/anthropic-sdk-go"
"inventory-agents/config"
```

)

const testingSecuritySystemPrompt = `You are an expert Testing and Security Agent for Go microservices.

Your responsibilities:

- Write comprehensive unit and integration tests using Go’s testing package and testify
- Design table-driven tests for business logic edge cases
- Implement security middleware: JWT validation, RBAC, rate limiting
- Write load/concurrency tests for stock operations
- Set up test containers for PostgreSQL and Kafka integration tests
- Identify and mitigate security vulnerabilities (SQL injection, race conditions, auth bypass)

Testing priorities for inventory:

- Concurrency tests: simultaneous stock reservations
- Boundary tests: stock going to zero, negative stock prevention
- Integration tests: full request lifecycle with real DB
- Event handler tests with mock Kafka

Security priorities:

- JWT middleware with role-based access (admin, warehouse, readonly)
- Rate limiting per API key
- Input sanitization and validation
- Audit logging for all stock mutations

Format code blocks as:
` + “`go\n// file: <filename>\n<code>\n`”

// TestingSecurityAgent writes tests and implements security
type TestingSecurityAgent struct {
*BaseAgent
}

func NewTestingSecurityAgent(cfg *config.Config) *TestingSecurityAgent {
return &TestingSecurityAgent{
BaseAgent: NewBaseAgent(cfg, “Testing & Security Agent”, testingSecuritySystemPrompt),
}
}

func (a *TestingSecurityAgent) Description() string {
return “Writes unit/integration tests, implements JWT auth, RBAC, rate limiting, and security middleware”
}

func (a *TestingSecurityAgent) Run(ctx context.Context, task string, agentContext map[string]string) (*AgentResult, error) {
prompt := fmt.Sprintf(`Write tests and implement security for an inventory microservice.

Task: %s

Please produce:

1. Unit tests for service layer (stock reservation, depletion, replenishment)
- Table-driven tests
- Mock repository using interfaces
- Concurrency test: 100 goroutines reserving last item simultaneously
1. Integration tests using testcontainers-go for PostgreSQL
1. Security middleware:
- JWT validation middleware (RS256)
- Role-based access control (admin, warehouse_manager, readonly)
- Rate limiter (token bucket per IP)
- Request ID and audit logging middleware
1. Security-focused test cases:
- Auth bypass attempts
- SQL injection via product IDs
- Race condition in stock updates
1. Makefile targets for running tests with coverage`, task)
   
   // Provide context from other agents
   contextSections := []string{}
   if api, ok := agentContext[“api_design”]; ok {
   contextSections = append(contextSections, “API Design:\n”+api)
   }
   if backend, ok := agentContext[“backend_db”]; ok {
   contextSections = append(contextSections, “Backend/DB Layer:\n”+backend)
   }
   for _, section := range contextSections {
   prompt += “\n\n” + section
   }
   
   messages := []anthropic.MessageParam{
   anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
   }
   
   output, err := a.Chat(ctx, messages)
   if err != nil {
   return nil, fmt.Errorf(”[%s] failed: %w”, a.Name(), err)
   }
   
   artifacts := ParseArtifacts(output)
   for i, art := range artifacts {
   if art.Filename == “” && art.Language == “go” {
   artifacts[i].Filename = fmt.Sprintf(“testing_security_%d.go”, i+1)
   }
   }
   
   return &AgentResult{
   AgentName: a.Name(),
   Output:    output,
   Artifacts: artifacts,
   }, nil
   }
