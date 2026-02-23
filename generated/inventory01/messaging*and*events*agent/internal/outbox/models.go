// file: internal/outbox/models.go
package outbox

import (
	"time"
)

// OutboxStatus represents the status of an outbox entry
type OutboxStatus string

const (
	StatusPending   OutboxStatus = "pending"
	StatusPublished OutboxStatus = "published"
	StatusFailed    OutboxStatus = "failed"
)

// OutboxEntry represents a message in the outbox table
type OutboxEntry struct {
	ID            string       `db:"id"`
	AggregateType string       `db:"aggregate_type"`
	AggregateID   string       `db:"aggregate_id"`
	EventType     string       `db:"event_type"`
	Topic         string       `db:"topic"`
	Key           string       `db:"key"`
	Payload       []byte       `db:"payload"`
	Headers       []byte       `db:"headers"`
	Status        OutboxStatus `db:"status"`
	RetryCount    int          `db:"retry_count"`
	MaxRetries    int          `db:"max_retries"`
	LastError     *string      `db:"last_error"`
	CreatedAt     time.Time    `db:"created_at"`
	ProcessedAt   *time.Time   `db:"processed_at"`
	ScheduledAt   time.Time    `db:"scheduled_at"`
}