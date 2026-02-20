package agents

import (
“context”
“fmt”

```
"github.com/anthropics/anthropic-sdk-go"
"lotus-agents/config"
```

)

const backendResponsibilities = `- Implement service layer business logic for all domain operations

- Design the database schema with proper indexing and constraints
- Write repository pattern code for data access using pgx or sqlx
- Handle concurrency: locking strategies, atomic updates, race conditions
- Implement database migrations (up/down)
- Apply domain-appropriate patterns (e.g. Saga, outbox, event sourcing)`

const backendOutputFormat = `When generating code, always include:

- PostgreSQL schema (tables, indexes, constraints, foreign keys)
- Repository interfaces and concrete implementations
- Service structs with dependency injection
- Concurrency-safe operations where relevant
- Up/down migration SQL files

Format Go code blocks as:
`+ "```go\n// file: <filename>\n<code>\n```" +`

Format SQL blocks as:
` + “`sql\n-- file: <filename>\n<sql>\n`”

// BackendDBAgent implements service logic and database layer for any microservice
type BackendDBAgent struct {
*BaseAgent
}

func NewBackendDBAgent(cfg *config.Config, svc *config.ServiceDefinition) *BackendDBAgent {
return &BackendDBAgent{
BaseAgent: NewBaseAgentForService(cfg, “Backend & Database Agent”, svc, backendResponsibilities, backendOutputFormat),
}
}

func (a *BackendDBAgent) Description() string {
return “Implements business logic, service layer, and database schema/repositories”
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
   
   if apiDesign, ok := agentContext[“api_design”]; ok {
   prompt += “\n\nAPI Design (implement these contracts):\n” + apiDesign
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
   if art.Filename == “” {
   switch art.Language {
   case “go”:
   artifacts[i].Filename = fmt.Sprintf(“service_%d.go”, i+1)
   case “sql”:
   artifacts[i].Filename = fmt.Sprintf(“migration_%d.sql”, i+1)
   }
   }
   }
   
   return &AgentResult{
   AgentName: a.Name(),
   Output:    output,
   Artifacts: artifacts,
   }, nil
   }
