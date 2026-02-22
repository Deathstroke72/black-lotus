package config

import (
	"os"
	"strconv"
)

const (
	defaultModel     = "claude-opus-4-5"
	defaultMaxTokens = 8192
)

// Config holds runtime configuration loaded from environment variables.
type Config struct {
	// AnthropicAPIKey is read from ANTHROPIC_API_KEY.
	AnthropicAPIKey string

	// Model is the Claude model to use, read from CLAUDE_MODEL.
	// Defaults to "claude-opus-4-5".
	Model string

	// MaxTokens is the maximum number of tokens per response, read from CLAUDE_MAX_TOKENS.
	// Defaults to 8192.
	MaxTokens int
}

// Load reads configuration from environment variables and returns a populated Config.
func Load() *Config {
	model := os.Getenv("CLAUDE_MODEL")
	if model == "" {
		model = defaultModel
	}

	maxTokens := defaultMaxTokens
	if v := os.Getenv("CLAUDE_MAX_TOKENS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			maxTokens = n
		}
	}

	return &Config{
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		Model:           model,
		MaxTokens:       maxTokens,
	}
}
