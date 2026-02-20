package agents

import (
“context”
“fmt”

```
"github.com/anthropics/anthropic-sdk-go"
"inventory-agents/config"
```

)

const messagingSystemPrompt = `You are an expert Messaging and Event-Driven Architecture Agent for Go microservices.

Your responsibilities:

- Design event schemas for inventory domain events
- Implement Kafka producers and consumers in Go
- Handle event sourcing and outbox patterns for reliability
- Design dead letter queues and retry logic
- Ensure at-least-once delivery with idempotency keys
- Implement event-driven stock synchronization across services

Key inventory events to handle:

- StockReserved, StockReleased, StockDepleted
- OrderFulfilled, OrderCancelled (consumed from order service)
- LowStockAlert, StockReplenished
- InventoryAdjusted (manual corrections)

Use the confluent-kafka-go or segmentio/kafka-go library.

Format code blocks as:
` + “`go\n// file: <filename>\n<code>\n`”

// MessagingAgent handles event-driven communication
type MessagingAgent struct {
*BaseAgent
}

func NewMessagingAgent(cfg *config.Config) *MessagingAgent {
return &MessagingAgent{
BaseAgent: NewBaseAgent(cfg, “Messaging & Events Agent”, messagingSystemPrompt),
}
}

func (a *MessagingAgent) Description() string {
return “Designs and implements Kafka-based event streaming, domain events, and async communication patterns”
}

func (a *MessagingAgent) Run(ctx context.Context, task string, agentContext map[string]string) (*AgentResult, error) {
prompt := fmt.Sprintf(`Design and implement the messaging layer for an inventory microservice.

Task: %s

Please produce:

1. Domain event structs with versioning and metadata (event_id, timestamp, correlation_id)
1. Kafka producer with:
- Transactional outbox pattern (events stored in DB before publishing)
- Retry with exponential backoff
- Message serialization (JSON)
1. Kafka consumer with:
- Consumer group management
- Idempotency handling (track processed event IDs)
- Dead letter queue for failed events
1. Event handlers for: order.created, order.cancelled, inventory.adjust
1. Topic configuration recommendations
1. Graceful shutdown handling`, task)
   
   if backendContext, ok := agentContext[“backend_db”]; ok {
   prompt += “\n\nDatabase/Service Context:\n” + backendContext
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
   artifacts[i].Filename = fmt.Sprintf(“messaging_%d.go”, i+1)
   }
   }
   
   return &AgentResult{
   AgentName: a.Name(),
   Output:    output,
   Artifacts: artifacts,
   }, nil
   }