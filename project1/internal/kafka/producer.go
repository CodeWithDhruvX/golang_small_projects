package kafka

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Producer struct {
	syncProducer sarama.SyncProducer
	logger       *zap.Logger
	brokers      []string
}

type ProducerMetrics struct {
	MessagesProduced int64
	Errors           int64
	LastProducedTime time.Time
}

func NewProducer(brokers []string) (*Producer, error) {
	return NewProducerWithLogger(brokers, nil)
}

func NewProducerWithLogger(brokers []string, logger *zap.Logger) (*Producer, error) {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Retry.Backoff = 100 * time.Millisecond
	config.Producer.Flush.Frequency = 100 * time.Millisecond
	
	// Enable Sarama logging
	sarama.Logger = log.New(log.Writer(), "[Sarama] ", log.LstdFlags)

	logger.Info("Initializing Kafka producer", 
		zap.Strings("brokers", brokers),
		zap.Int("retry_max", config.Producer.Retry.Max),
		zap.Duration("retry_backoff", config.Producer.Retry.Backoff))

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Error("Failed to create Kafka producer", 
			zap.Error(err),
			zap.Strings("brokers", brokers))
		return nil, fmt.Errorf("failed to completely start producer: %w", err)
	}

	logger.Info("Kafka producer initialized successfully")

	return &Producer{
		syncProducer: producer,
		logger:       logger,
		brokers:      brokers,
	}, nil
}

func (p *Producer) PublishEvent(topic string, event interface{}) error {
	start := time.Now()
	
	p.logger.Debug("Starting to publish event", 
		zap.String("topic", topic),
		zap.String("event_type", fmt.Sprintf("%T", event)))

	eventBytes, err := json.Marshal(event)
	if err != nil {
		p.logger.Error("Failed to serialize event", 
			zap.String("topic", topic),
			zap.Error(err))
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	p.logger.Debug("Event serialized successfully", 
		zap.String("topic", topic),
		zap.Int("event_size_bytes", len(eventBytes)))

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(eventBytes),
	}

	partition, offset, err := p.syncProducer.SendMessage(msg)
	duration := time.Since(start)
	
	if err != nil {
		p.logger.Error("Failed to send message to Kafka", 
			zap.String("topic", topic),
			zap.Error(err),
			zap.Duration("attempt_duration", duration))
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	p.logger.Info("Message published successfully", 
		zap.String("topic", topic),
		zap.Int32("partition", partition),
		zap.Int64("offset", offset),
		zap.Duration("publish_duration", duration),
		zap.Int("message_size", len(eventBytes)))

	return nil
}

func (p *Producer) GetHealthStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":  "healthy",
		"brokers": p.brokers,
		"type":    "producer",
	}
}

func (p *Producer) Close() error {
	p.logger.Info("Closing Kafka producer")
	return p.syncProducer.Close()
}
