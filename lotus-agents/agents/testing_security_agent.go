package agents

import (
	"context"
	"fmt"

	"github.com/Deathstroke72/black-lotus/lotus-agents/config"

	"github.com/anthropics/anthropic-sdk-go"
)

const testingResponsibilities = `- Write tests at every CA layer (unit tests mock the layer beneath; integration tests use testcontainers)
- internal/interfaces/http/middleware: JWT RS256, RBAC per-endpoint, rate limiter (token bucket), request ID, audit logging
- cmd/server/main.go: dependency wiring — pgx pool → repos → use cases → handlers → router → HTTP server
- Place test files next to the code they test (e.g. internal/domain/entity/payment_test.go)
- Use table-driven tests; mock repository interfaces with hand-written or mockery-generated mocks
- Integration tests use testcontainers-go (PostgreSQL, Kafka as needed)
- Makefile: test, test-integration, test-race, coverage, lint, build, run targets`

const testingOutputFormat = `Produce these files (every code block MUST start with // file: <path>):

  internal/interfaces/http/middleware/jwt_middleware.go
  internal/interfaces/http/middleware/rbac_middleware.go
  internal/interfaces/http/middleware/rate_limit_middleware.go
  internal/interfaces/http/middleware/request_id_middleware.go
  internal/interfaces/http/middleware/audit_middleware.go

  cmd/server/main.go
      → Wire all layers: pgx pool → postgres repos → use cases → handlers → middleware → HTTP server

  internal/domain/entity/<entity>_test.go          (unit tests for domain entities)
  internal/application/usecase/<usecase>_test.go   (unit tests with mock repos)
  internal/interfaces/http/handler/<handler>_test.go (handler tests with mock use cases)
  internal/infrastructure/postgres/repository/<repo>_integration_test.go (testcontainers)

  Makefile

Format Go: ` + "```go\n// file: <path>/<filename>.go\n<code>\n```" + `
Format Makefile: ` + "```makefile\n// file: Makefile\n<content>\n```"

// TestingSecurityAgent writes tests and implements security for any microservice
type TestingSecurityAgent struct {
	*BaseAgent
}

func NewTestingSecurityAgent(cfg *config.Config, svc *config.ServiceDefinition) *TestingSecurityAgent {
	return &TestingSecurityAgent{
		BaseAgent: NewBaseAgentForService(cfg, "Testing & Security Agent", svc, testingResponsibilities, testingOutputFormat),
	}
}

func (a *TestingSecurityAgent) Description() string {
	return "Writes unit/integration tests and implements JWT auth, RBAC, rate limiting, and security middleware"
}

func (a *TestingSecurityAgent) Run(ctx context.Context, svc *config.ServiceDefinition, agentContext map[string]string) (*AgentResult, error) {
	prompt := fmt.Sprintf(`Write tests and implement security for the following microservice:

%s

Please produce:

1. Unit tests for the service layer — one test file per major operation
- Table-driven tests with success and failure cases
- Mock repositories generated from interfaces
- Concurrency tests for any operations that mutate shared state
1. Integration tests using testcontainers-go
1. Security middleware stack:
- JWT validation (RS256) with roles appropriate to this service
- Role-based access control per endpoint
- Rate limiter (token bucket, configurable per role)
- Request ID + audit logging middleware for all mutations
1. Security-focused test cases:
- Unauthorized access attempts
- Input validation / injection attempts
- Any domain-specific security concerns
1. Makefile with: test, test-integration, coverage, lint targets`, svc.Prompt())

	if api, ok := agentContext["api_design"]; ok {
		prompt += "\n\nAPI Design (write tests and middleware for these endpoints):\n" + api
	}
	if backend, ok := agentContext["backend_db"]; ok {
		prompt += "\n\nService/Repo Layer (mock these interfaces in tests):\n" + backend
	}

	messages := []anthropic.MessageParam{
		anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
	}

	output, err := a.Chat(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("[%s] failed: %w", a.Name(), err)
	}

	artifacts := ParseArtifacts(output)
	for i, art := range artifacts {
		if art.Filename == "" {
			switch art.Language {
			case "go":
				artifacts[i].Filename = fmt.Sprintf("internal/interfaces/http/middleware/middleware_%d.go", i+1)
			case "makefile":
				artifacts[i].Filename = "Makefile"
			}
		}
	}

	return &AgentResult{
		AgentName: a.Name(),
		Output:    output,
		Artifacts: artifacts,
	}, nil
}
