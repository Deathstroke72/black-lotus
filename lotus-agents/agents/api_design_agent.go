package agents

import (
	"context"
	"fmt"

	"github.com/Deathstroke72/black-lotus/lotus-agents/config"

	"github.com/anthropics/anthropic-sdk-go"
)

const apiDesignResponsibilities = `- Design the HTTP interface layer (Clean Architecture: interfaces/http)
- Define DTOs for every HTTP request and response — NOT domain entities
- Write thin HTTP handlers that decode a DTO, call a use case port interface, and return a DTO
- Define the router that wires URL patterns to handlers and attaches middleware
- Reference use case interfaces (e.g. usecase.CreatePaymentUseCase) by name — do not implement them
- Follow REST best practices: correct HTTP verbs, status codes, and error response envelope
- Zero domain logic in handlers: validate input → call use case → map result to DTO`

const apiDesignOutputFormat = `Produce these files (every code block MUST start with // file: <path>):

  internal/interfaces/http/dto/<entity>_dto.go
      → Request/Response structs with json tags; no domain types embedded

  internal/interfaces/http/handler/<entity>_handler.go
      → Struct with use case port interface constructor; methods: decode → call use case → encode

  internal/interfaces/http/router/router.go
      → Registers all routes, attaches middleware chain

Format: ` + "```go\n// file: internal/interfaces/http/<subdir>/<filename>.go\n<code>\n```" + `

Do not generate domain entities, repository implementations, or use case logic.`

// APIDesignAgent designs REST API contracts for any microservice
type APIDesignAgent struct {
	*BaseAgent
}

func NewAPIDesignAgent(cfg *config.Config, svc *config.ServiceDefinition) *APIDesignAgent {
	return &APIDesignAgent{
		BaseAgent: NewBaseAgentForService(cfg, "API Design Agent", svc, apiDesignResponsibilities, apiDesignOutputFormat),
	}
}

func (a *APIDesignAgent) Description() string {
	return "Designs RESTful API contracts, route definitions, and request/response schemas"
}

func (a *APIDesignAgent) Run(ctx context.Context, svc *config.ServiceDefinition, agentContext map[string]string) (*AgentResult, error) {
	prompt := fmt.Sprintf(`Design the REST API for the following microservice:

%s

Please produce:

1. A complete list of API endpoints (HTTP method, path, description) for all operations listed above
1. Go structs for all request and response payloads with JSON tags
1. Router setup code (chi or net/http)
1. Standardized error response format
1. OpenAPI-style godoc comments for each endpoint
1. Any domain-specific validation rules or constraints`, svc.Prompt())

	if ctx, ok := agentContext["project_context"]; ok {
		prompt += "\n\nAdditional Context:\n" + ctx
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
		if art.Filename == "" && art.Language == "go" {
			artifacts[i].Filename = fmt.Sprintf("internal/interfaces/http/handler/api_%d.go", i+1)
		}
	}

	return &AgentResult{
		AgentName: a.Name(),
		Output:    output,
		Artifacts: artifacts,
	}, nil
}
