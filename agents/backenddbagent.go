package agents

import (
“context”
“fmt”

```
"github.com/anthropics/anthropic-sdk-go"
"inventory-agents/config"
```

)

const backendSystemPrompt = `You are an expert Backend Development and Database Agent specializing in Go microservices.

Your responsibilities:

- Implement service layer business logic for inventory operations
- Design PostgreSQL schemas with proper indexing and constraints
- Write repository pattern code for data access
- Handle concurrency: optimistic locking, row-level locks, atomic stock updates
- Implement the Saga pattern for distributed transactions with order services
- Use sqlx or pgx for database interactions

Key concerns for inventory:

- Preventing overselling (stock going negative) with proper locking
- Efficient bulk stock updates
- Audit trails for stock changes
- Multi-warehouse support

Format code blocks as:
`+ "```go\n// file: <filename>\n<code>\n```" +`

For SQL, use:
` + “`sql\n-- file: <filename>\n<sql>\n`”

// BackendDBAgent implements service logic and database layer
type BackendDBAgent struct {
*BaseAgent
}

func NewBackendDBAgent(cfg *config.Config) *BackendDBAgent {
return &BackendDBAgent{
BaseAgent: NewBaseAgent(cfg, “Backend & Database Agent”, backendSystemPrompt),
}
}

func (a *BackendDBAgent) Description() string {
return “Implements business logic, service layer, and database schema/repository for the inventory microservice”
}

func (a *BackendDBAgent) Run(ctx context.Context, task string, agentContext map[string]string) (*AgentResult, error) {
prompt := fmt.Sprintf(`Implement the backend service layer and database code for an inventory microservice.

Task: %s

Please produce:

1. PostgreSQL schema (tables, indexes, constraints) for:
- products, inventory_items, warehouses, stock_movements, reservations
1. Repository interfaces and implementations using pgx/sqlx
1. Service layer with business logic for:
- Stock reservation and release
- Atomic stock decrements (prevent negative stock)
- Stock replenishment
- Multi-warehouse aggregation
1. Database migrations (up/down)
1. Concurrency handling with SELECT FOR UPDATE or optimistic locking`, task)
   
   if apiDesign, ok := agentContext[“api_design”]; ok {
   prompt += “\n\nAPI Design from API Design Agent (implement these contracts):\n” + apiDesign
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
