package agents

import (
	"context"
	"fmt"

	"github.com/Deathstroke72/black-lotus/lotus-agents/config"

	"github.com/anthropics/anthropic-sdk-go"
)

const apiDesignResponsibilities = `- Design clean, RESTful API contracts tailored to this service’s domain

- Define route structures and OpenAPI-style documentation
- Specify request/response schemas with proper validation rules
- Handle domain-specific edge cases and error scenarios
- Follow REST best practices and correct HTTP semantics`

const apiDesignOutputFormat = `When generating code, always include:

- Route definitions using Go's net/http or chi router
- Request/Response structs with JSON tags and validation
- Proper HTTP status codes and error response formats
- Comments explaining design decisions

Format code blocks as:
` + "`go\n// file: <filename>\n<code>\n`"

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
			artifacts[i].Filename = fmt.Sprintf("api_%d.go", i+1)
		}
	}

	return &AgentResult{
		AgentName: a.Name(),
		Output:    output,
		Artifacts: artifacts,
	}, nil
}
