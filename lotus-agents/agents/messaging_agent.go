package agents

import (
"context"
"fmt"

"github.com/anthropics/anthropic-sdk-go"
"lotus-agents/config"
)

const messagingResponsibilities = `- Design domain event schemas appropriate for this service

- Implement Kafka producers with the transactional outbox pattern
- Implement Kafka consumers with idempotency and dead letter queue handling
- Define which events this service publishes and which it consumes
- Ensure at-least-once delivery with retry and backoff logic
- Handle graceful shutdown of consumers`

const messagingOutputFormat = `When generating code, always include:

- Domain event structs with versioning, event_id, correlation_id, and timestamp
- Kafka producer with transactional outbox (events written to DB before publishing)
- Kafka consumer group with idempotency tracking
- Dead letter queue handling
- Graceful shutdown

Format code blocks as:
` + “`go\n// file: <filename>\n<code>\n`”

// MessagingAgent handles event-driven communication for any microservice
type MessagingAgent struct {
*BaseAgent
}

func NewMessagingAgent(cfg *config.Config, svc *config.ServiceDefinition) *MessagingAgent {
return &MessagingAgent{
BaseAgent: NewBaseAgentForService(cfg, “Messaging & Events Agent”, svc, messagingResponsibilities, messagingOutputFormat),
}
}

func (a *MessagingAgent) Description() string {
return “Designs and implements Kafka-based domain events, producers, consumers, and async communication”
}

func (a *MessagingAgent) Run(ctx context.Context, svc *config.ServiceDefinition, agentContext map[string]string) (*AgentResult, error) {
prompt := fmt.Sprintf(`Design and implement the messaging/eventing layer for the following microservice:

%s

Please produce:

1. Domain event structs this service will PUBLISH (derived from its operations and entities)
1. Events this service will CONSUME from its integrations
1. Kafka producer implementation with:
- Transactional outbox pattern
- Exponential backoff retry
- JSON serialization with schema versioning
1. Kafka consumer with:
- Consumer group setup
- Idempotency key tracking to prevent duplicate processing
- Dead letter queue for poison messages
1. Event handler functions for each consumed event type
1. Topic naming conventions and configuration recommendations
1. Graceful shutdown logic`, svc.Prompt())
   
   if backend, ok := agentContext[“backend_db”]; ok {
   prompt += “\n\nDatabase/Service Context (outbox table should align with this schema):\n” + backend
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
