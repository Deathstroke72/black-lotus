// file: internal/infrastructure/kafka/producer/config.go
package producer

import (
	"time"
)

// Config holds Kafka producer configuration
type Config struct {
	Brokers           []string
	ClientID          string
	Acks              string // "all", "1", "0"
	MaxRetries        int
	RetryBackoff      time.Duration
	MaxBackoff        time.Duration
	BatchSize         int
	LingerMs          int
	CompressionType   string // "none", "gzip", "snappy", "lz4", "zstd"
	IdempotentEnabled bool
}

// DefaultConfig returns a production-ready default configuration
func DefaultConfig() Config {
	return Config{
		Brokers:           []string{"localhost:9092"},
		ClientID:          "inventory-service",
		Acks:              "all",
		MaxRetries:        5,
		RetryBackoff:      100 * time.Millisecond,
		MaxBackoff:        10 * time.Second,
		BatchSize:         16384,
		LingerMs:          5,
		CompressionType:   "snappy",
		IdempotentEnabled: true,
	}
}

// OutboxConfig holds configuration for outbox processing
type OutboxConfig struct {
	PollInterval    time.Duration
	BatchSize       int
	MaxRetries      int
	RetentionPeriod time.Duration // How long to keep published entries
}

// DefaultOutboxConfig returns default outbox configuration
func DefaultOutboxConfig() OutboxConfig {
	return OutboxConfig{
		PollInterval:    100 * time.Millisecond,
		BatchSize:       100,
		MaxRetries:      5,
		RetentionPeriod: 7 * 24 * time.Hour,
	}
}