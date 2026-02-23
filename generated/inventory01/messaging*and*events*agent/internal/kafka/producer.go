// file: internal/kafka/producer.go
package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"

	"inventory-service/internal/outbox"
)

// ProducerConfig holds configuration for the Kafka producer
type ProducerConfig struct {
	Brokers           []string
	ClientID          string
	RequiredAcks      sarama.RequiredAcks
	MaxRetries        int
	RetryBackoff      time.Duration
	FlushFrequency    time.Duration
	FlushMessages     int
	EnableIdempotent  bool
	CompressionCodec  sarama.CompressionCodec
}

// DefaultProducerConfig returns sensible defaults
func DefaultProducerConfig(brokers []string) ProducerConfig {
	return ProducerConfig{
		Brokers:          brokers,
		ClientID:         "inventory-service",
		RequiredAcks:     sarama.WaitForAll,
		MaxRetries:       5,
		RetryBackoff:     100 * time.Millisecond,
		FlushFrequency:   500 * time.Millisecond,
		FlushMessages:    100,
		EnableIdempotent: true,
		CompressionCodec: sarama.CompressionSnappy,
	}
}

// Producer wraps Sarama async producer with outbox pattern support
type Producer struct {
	producer     sarama.AsyncProducer
	config       ProducerConfig
	logger       zerolog.Logger
	outboxRepo   *outbox.Repository
	
	wg           sync.WaitGroup
	shutdownCh   chan struct{}
	shutdownOnce sync.Once
}

// NewProducer creates a new Kafka producer
func NewProducer(cfg ProducerConfig, outboxRepo *outbox.Repository, logger zerolog.Logger) (*Producer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.ClientID = cfg.ClientID
	saramaConfig.Producer.RequiredAcks = cfg.RequiredAcks
	saramaConfig.Producer.Retry.Max = cfg.MaxRetries
	saramaConfig.Producer.Retry.Backoff = cfg.RetryBackoff
	saramaConfig.Producer.Flush.Frequency = cfg.FlushFrequency
	saramaConfig.Producer.Flush.Messages = cfg.FlushMessages
	saramaConfig.Producer.Idempotent = cfg.EnableIdempotent
	saramaConfig.Producer.Compression = cfg.CompressionCodec
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Net.MaxOpenRequests = 1 // Required for idempotent producer

	producer, err := sarama.NewAsyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	p := &Producer{
		producer:   producer,
		config:     cfg,
		logger:     logger,
		outboxRepo: outboxRepo,
		shutdownCh: make(chan struct{}),
	}

	// Start success/error handlers
	p.wg.Add(2)
	go p.handleSuccesses()
	go p.handleErrors()

	return p, nil
}

// handleSuccesses processes successful message deliveries
func (p *Producer) handleSuccesses() {
	defer p.wg.Done()

	for {
		select {
		case msg, ok := <-p.producer.Successes():
			if !ok {
				return
			}
			
			// Extract outbox ID from metadata
			if outboxID, ok := msg.Metadata.(string); ok {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := p.outboxRepo.MarkAsPublished(ctx, outboxID); err != nil {
					p.logger.Error().
						Err(err).
						Str("outbox_id", outboxID).
						Msg("Failed to mark outbox entry as published")
				} else {
					p.logger.Debug().
						Str("outbox_id", outboxID).
						Str("topic", msg.Topic).
						Int32("partition", msg.Partition).
						Int64("offset", msg.Offset).
						Msg("Message published successfully")
				}
				cancel()
			}

		case <-p.shutdownCh:
			return
		}
	}
}

// handleErrors processes failed message deliveries
func (p *Producer) handleErrors() {
	defer p.wg.Done()

	for {
		select {
		case err, ok := <-p.producer.Errors():
			if !ok {
				return
			}

			// Extract outbox ID and retry count from metadata
			if outboxID, ok := err.Msg.Metadata.(string); ok {
				p.logger.Error().
					Err(err.Err).
					Str("outbox_id", outboxID).
					Str("topic", err.Msg.Topic).
					Msg("Failed to publish message")

				// The outbox processor will handle retries
			}

		case <-p.shutdownCh:
			return
		}
	}
}

// PublishFromOutbox publishes a message from an outbox entry
func (p *Producer) PublishFromOutbox(entry *outbox.OutboxEntry) error {
	// Parse headers
	var headers []sarama.RecordHeader
	if len(entry.Headers) > 0 {
		var headerMap map[string]string
		if err := json.Unmarshal(entry.Headers, &headerMap); err == nil {
			for k, v := range headerMap {
				headers = append(headers, sarama.RecordHeader{
					Key:   []byte(k),
					Value: []byte(v),
				})
			}
		}
	}

	// Add standard headers
	headers = append(headers,
		sarama.RecordHeader{Key: []byte("event_type"), Value: []byte(entry.EventType)},
		sarama.RecordHeader{Key: []byte("aggregate_type"), Value: []byte(entry.AggregateType)},
		sarama.RecordHeader{Key: []byte("aggregate_id"), Value: []byte(entry.AggregateID)},
		sarama.RecordHeader{Key: []byte("outbox_id"), Value: []byte(entry.ID)},
	)

	msg := &sarama.ProducerMessage{
		Topic:    entry.Topic,
		Key:      sarama.StringEncoder(entry.Key),
		Value:    sarama.ByteEncoder(entry.Payload),
		Headers:  headers,
		Metadata: entry.ID, // Store outbox ID for success/error handling
	}

	select {
	case p.producer.Input() <- msg:
		return nil
	case <-p.shutdownCh:
		return fmt.Errorf("producer is shutting down")
	}
}

// Close gracefully shuts down the producer
func (p *Producer) Close() error {
	p.shutdownOnce.Do(func() {
		close(p.shutdownCh)
	})

	// Close the producer
	if err := p.producer.Close(); err != nil {
		p.logger.Error().Err(err).Msg("Error closing producer")
		return err
	}

	// Wait for handlers to finish
	p.wg.Wait()
	p.logger.Info().Msg("Producer shut down gracefully")
	return nil
}