package agents

import (
	"context"
	"fmt"

	"github.com/Deathstroke72/black-lotus/lotus-agents/config"

	"github.com/anthropics/anthropic-sdk-go"
)

const backendResponsibilities = `- Implement the domain layer, application layer, and PostgreSQL infrastructure (Clean Architecture)
- domain/entity: pure Go structs with business invariants; zero external imports (no net/http, no database/sql)
- domain/repository: interfaces declaring data access contracts (not implementations)
- domain/service: domain services for business logic spanning multiple entities
- application/usecase: one file per use case; orchestrates domain via repository interfaces
- application/port: input port interfaces (what HTTP handlers call)
- infrastructure/postgres/repository: pgx concrete implementations of domain repository interfaces
- infrastructure/postgres/migration: SQL up/down migration files
- Dependency Rule: domain and application layers must never import net/http, database/sql, or any Kafka SDK`

const backendOutputFormat = `Produce these files (every code block MUST start with // file: or -- file:):

  internal/domain/entity/<entity>.go
      → Pure struct with business invariants; zero external imports

  internal/domain/repository/<entity>_repository.go
      → Interface only; methods receive/return domain entities

  internal/domain/service/<service>.go
      → Domain logic spanning multiple entities; zero external imports

  internal/application/port/input.go
      → Use case input port interfaces (what handlers call)

  internal/application/usecase/<operation>_usecase.go
      → One file per use case; depends on domain repository interfaces

  internal/infrastructure/postgres/repository/<entity>_repository.go
      → Implements domain repository interface using pgx

  internal/infrastructure/postgres/migration/<NNN>_<description>_up.sql
  internal/infrastructure/postgres/migration/<NNN>_<description>_down.sql

Format Go: ` + "```go\n// file: internal/<layer>/<subdir>/<filename>.go\n<code>\n```" + `
Format SQL: ` + "```sql\n-- file: internal/infrastructure/postgres/migration/<filename>.sql\n<sql>\n```" + `

The domain and application layers must contain zero external package imports.`

// BackendDBAgent implements service logic and database layer for any microservice
type BackendDBAgent struct {
	*BaseAgent
}

func NewBackendDBAgent(cfg *config.Config, svc *config.ServiceDefinition) *BackendDBAgent {
	return &BackendDBAgent{
		BaseAgent: NewBaseAgentForService(cfg, "Backend & Database Agent", svc, backendResponsibilities, backendOutputFormat),
	}
}

func (a *BackendDBAgent) Description() string {
	return "Implements business logic, service layer, and database schema/repositories"
}

func (a *BackendDBAgent) Run(ctx context.Context, svc *config.ServiceDefinition, agentContext map[string]string) (*AgentResult, error) {
	prompt := fmt.Sprintf(`Implement the backend service layer and database code for the following microservice:

%s

Please produce:

1. PostgreSQL schema for all entities listed above (tables, indexes, constraints)
1. Repository interfaces and implementations using pgx
1. Service layer structs with all business operations implemented
1. Database migration files (up + down)
1. Any concurrency or consistency mechanisms needed for the operations above
1. Dependency injection wiring (how repos plug into services)`, svc.Prompt())

	if apiDesign, ok := agentContext["api_design"]; ok {
		prompt += "\n\nAPI Design (implement these contracts):\n" + apiDesign
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
				artifacts[i].Filename = fmt.Sprintf("internal/application/usecase/usecase_%d.go", i+1)
			case "sql":
				artifacts[i].Filename = fmt.Sprintf("internal/infrastructure/postgres/migration/migration_%d.sql", i+1)
			}
		}
	}

	return &AgentResult{
		AgentName: a.Name(),
		Output:    output,
		Artifacts: artifacts,
	}, nil
}
