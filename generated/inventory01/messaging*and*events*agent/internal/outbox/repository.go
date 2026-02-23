// file: internal/outbox/repository.go
package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Repository handles outbox persistence operations
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new outbox repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// InsertInTx inserts an outbox entry within an existing transaction
func (r *Repository) InsertInTx(ctx context.Context, tx *sqlx.Tx, entry *OutboxEntry) error {
	query := `
		INSERT INTO event_outbox (
			id, aggregate_type, aggregate_id, event_type, topic, key,
			payload, headers, status, retry_count, max_retries, created_at, scheduled_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)
	`

	if entry.ID == "" {
		entry.ID = uuid.New().String()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	if entry.ScheduledAt.IsZero() {
		entry.ScheduledAt = entry.CreatedAt
	}
	if entry.Status == "" {
		entry.Status = StatusPending
	}
	if entry.MaxRetries == 0 {
		entry.MaxRetries = 5
	}

	_, err := tx.ExecContext(ctx, query,
		entry.ID,
		entry.AggregateType,
		entry.AggregateID,
		entry.EventType,
		entry.Topic,
		entry.Key,
		entry.Payload,
		entry.Headers,
		entry.Status,
		entry.RetryCount,
		entry.MaxRetries,
		entry.CreatedAt,
		entry.ScheduledAt,
	)

	return err
}

// FetchPendingBatch fetches a batch of pending outbox entries for processing
func (r *Repository) FetchPendingBatch(ctx context.Context, batchSize int) ([]OutboxEntry, error) {
	query := `
		SELECT id, aggregate_type, aggregate_id, event_type, topic, key,
			   payload, headers, status, retry_count, max_retries, last_error,
			   created_at, processed_at, scheduled_at
		FROM event_outbox
		WHERE status = $1
		  AND scheduled_at <= $2
		  AND retry_count < max_retries
		ORDER BY created_at ASC
		LIMIT $3
		FOR UPDATE SKIP LOCKED
	`

	var entries []OutboxEntry
	err := r.db.SelectContext(ctx, &entries, query, StatusPending, time.Now().UTC(), batchSize)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending entries: %w", err)
	}

	return entries, nil
}

// MarkAsPublished marks an outbox entry as successfully published
func (r *Repository) MarkAsPublished(ctx context.Context, id string) error {
	query := `
		UPDATE event_outbox
		SET status = $1, processed_at = $2
		WHERE id = $3
	`

	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query, StatusPublished, now, id)
	return err
}

// MarkAsFailed marks an outbox entry as failed with retry scheduling
func (r *Repository) MarkAsFailed(ctx context.Context, id string, errMsg string, retryCount int) error {
	// Exponential backoff: 2^retryCount seconds, max 5 minutes
	backoffSeconds := 1 << retryCount
	if backoffSeconds > 300 {
		backoffSeconds = 300
	}
	nextRetry := time.Now().UTC().Add(time.Duration(backoffSeconds) * time.Second)

	query := `
		UPDATE event_outbox
		SET retry_count = $1, last_error = $2, scheduled_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, retryCount, errMsg, nextRetry, id)
	return err
}

// MarkAsDeadLetter marks an entry as permanently failed
func (r *Repository) MarkAsDeadLetter(ctx context.Context, id string, errMsg string) error {
	query := `
		UPDATE event_outbox
		SET status = $1, last_error = $2, processed_at = $3
		WHERE id = $4
	`

	now := time.Now().UTC()
	_, err := r.db.ExecContext(ctx, query, StatusFailed, errMsg, now, id)
	return err
}

// CleanupOldEntries removes successfully processed entries older than retention period
func (r *Repository) CleanupOldEntries(ctx context.Context, retentionDays int) (int64, error) {
	query := `
		DELETE FROM event_outbox
		WHERE status = $1
		  AND processed_at < $2
	`

	cutoff := time.Now().UTC().AddDate(0, 0, -retentionDays)
	result, err := r.db.ExecContext(ctx, query, StatusPublished, cutoff)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// CreateOutboxEntry is a helper to create an outbox entry from an event
func CreateOutboxEntry(aggregateType, aggregateID, eventType, topic string, payload interface{}, headers map[string]string) (*OutboxEntry, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var headersBytes []byte
	if headers != nil {
		headersBytes, err = json.Marshal(headers)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal headers: %w", err)
		}
	}

	return &OutboxEntry{
		ID:            uuid.New().String(),
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Topic:         topic,
		Key:           aggregateID, // Use aggregate ID as partition key
		Payload:       payloadBytes,
		Headers:       headersBytes,
		Status:        StatusPending,
		RetryCount:    0,
		MaxRetries:    5,
		CreatedAt:     time.Now().UTC(),
		ScheduledAt:   time.Now().UTC(),
	}, nil
}