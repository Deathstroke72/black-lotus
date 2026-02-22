package agents

import (
"context"
"fmt"


"github.com/anthropics/anthropic-sdk-go"
"lotus-agents/config"
)

const testingResponsibilities = `- Write comprehensive unit tests using Go’s testing package and testify

- Design table-driven tests covering edge cases for all domain operations
- Write integration tests using testcontainers-go for real dependencies
- Implement JWT authentication middleware with role-based access control
- Add rate limiting, request ID generation, and audit logging middleware
- Identify and test security vulnerabilities specific to this service’s domain`

const testingOutputFormat = `When generating code, always include:

- Table-driven unit tests with mock repositories (using interfaces)
- Integration tests with testcontainers (PostgreSQL, Kafka as needed)
- Concurrency tests for any operations that modify shared state
- JWT middleware (RS256), RBAC roles appropriate to this service
- Rate limiter middleware (token bucket per IP/API key)
- Audit logging middleware for all mutating operations
- A Makefile with test targets and coverage reporting

Format code blocks as:
` + “`go\n// file: <filename>\n<code>\n`”

// TestingSecurityAgent writes tests and implements security for any microservice
type TestingSecurityAgent struct {
*BaseAgent
}

func NewTestingSecurityAgent(cfg *config.Config, svc *config.ServiceDefinition) *TestingSecurityAgent {
return &TestingSecurityAgent{
BaseAgent: NewBaseAgentForService(cfg, “Testing & Security Agent”, svc, testingResponsibilities, testingOutputFormat),
}
}

func (a *TestingSecurityAgent) Description() string {
return “Writes unit/integration tests and implements JWT auth, RBAC, rate limiting, and security middleware”
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
   
   if api, ok := agentContext[“api_design”]; ok {
   prompt += “\n\nAPI Design (write tests and middleware for these endpoints):\n” + api
   }
   if backend, ok := agentContext[“backend_db”]; ok {
   prompt += “\n\nService/Repo Layer (mock these interfaces in tests):\n” + backend
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
   artifacts[i].Filename = fmt.Sprintf(“test_%d.go”, i+1)
   }
   }
   
   return &AgentResult{
   AgentName: a.Name(),
   Output:    output,
   Artifacts: artifacts,
   }, nil
   }
