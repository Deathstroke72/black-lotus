package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Deathstroke72/black-lotus/lotus-agents/config"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// AgentResult holds the output of an agent’s work
type AgentResult struct {
	AgentName string
	Output    string
	Artifacts []Artifact
	Error     error
}

// Artifact represents a file or piece of code produced by an agent
type Artifact struct {
	Filename string
	Content  string
	Language string
}

// Agent defines the interface every specialized agent must implement
type Agent interface {
	Name() string
	Description() string
	Run(ctx context.Context, svc *config.ServiceDefinition, agentContext map[string]string) (*AgentResult, error)
}

// BaseAgent provides shared Claude API functionality
type BaseAgent struct {
	client       anthropic.Client
	cfg          *config.Config
	agentName    string
	systemPrompt string
}

// buildSystemPrompt creates a dynamic system prompt incorporating the service definition
func buildSystemPrompt(role, svcName, language, responsibilities, outputFormat string) string {
	return fmt.Sprintf(`You are an expert %s specializing in %s microservices written in %s.

Service you are building: %s

Your responsibilities:
%s

%s`, role, svcName, language, svcName, responsibilities, outputFormat)
}

func NewBaseAgent(cfg *config.Config, name, systemPrompt string) *BaseAgent {
	client := anthropic.NewClient(option.WithAPIKey(cfg.AnthropicAPIKey))
	return &BaseAgent{
		client:       client,
		cfg:          cfg,
		agentName:    name,
		systemPrompt: systemPrompt,
	}
}

func NewBaseAgentForService(cfg *config.Config, name string, svc *config.ServiceDefinition, responsibilities, outputFormat string) *BaseAgent {
	role := name
	prompt := buildSystemPrompt(role, svc.Name, svc.Language, responsibilities, outputFormat)
	return NewBaseAgent(cfg, name, prompt)
}

func (b *BaseAgent) Name() string { return b.agentName }

// WithSystemPrompt returns a copy of the base agent with an updated system prompt
func (b *BaseAgent) WithSystemPrompt(prompt string) *BaseAgent {
	clone := *b
	clone.systemPrompt = prompt
	return &clone
}

// Chat sends a message to Claude and returns the response text
func (b *BaseAgent) Chat(ctx context.Context, messages []anthropic.MessageParam) (string, error) {
	resp, err := b.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(b.cfg.Model),
		MaxTokens: int64(b.cfg.MaxTokens),
		System: []anthropic.TextBlockParam{
			{Text: b.systemPrompt},
		},
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("claude API error: %w", err)
	}

	var sb strings.Builder
	for _, block := range resp.Content {
		if block.Type == "text" {
			sb.WriteString(block.Text)
		}
	}
	return sb.String(), nil

}

// ParseArtifacts extracts code blocks from markdown-style output
func ParseArtifacts(output string) []Artifact {
	var artifacts []Artifact
	lines := strings.Split(output, "\n")
	var inBlock bool
	var lang, filename string
	var blockLines []string

	for _, line := range lines {
		if strings.HasPrefix(line, "```") && !inBlock {
			inBlock = true
			lang = strings.TrimPrefix(line, "```")
			filename = ""
			blockLines = nil
		} else if line == "```" && inBlock {
			inBlock = false
			artifacts = append(artifacts, Artifact{
				Filename: filename,
				Language: lang,
				Content:  strings.Join(blockLines, "\n"),
			})
		} else if inBlock {
			// Detect filename hints like // file: main.go
			if strings.HasPrefix(line, "// file:") || strings.HasPrefix(line, "# file:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					filename = strings.TrimSpace(parts[1])
				}
			}
			blockLines = append(blockLines, line)
		}
	}
	return artifacts
}

// ToJSON is a helper to pretty-print structs for context passing
func ToJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
