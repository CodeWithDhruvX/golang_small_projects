package kafka

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type ConsumerGroupHandler struct {
	Ready   chan bool
	Process func(message *sarama.ConsumerMessage) error
	logger  *zap.Logger
	metrics *ConsumerMetrics
}

type ConsumerMetrics struct {
	MessagesConsumed int64
	ProcessingErrors int64
	LastConsumedTime time.Time
}

func NewConsumerGroupHandler(logger *zap.Logger, processMsg func(*sarama.ConsumerMessage) error) *ConsumerGroupHandler {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	
	return &ConsumerGroupHandler{
		Ready:   make(chan bool),
		Process: processMsg,
		logger:  logger,
		metrics: &ConsumerMetrics{},
	}
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	h.logger.Info("Consumer group setup completed")
	close(h.Ready)
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	h.logger.Info("Consumer group cleanup completed")
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	h.logger.Info("Starting to consume from claim",
		zap.String("topic", claim.Topic()),
		zap.Int32("partition", claim.Partition()),
		zap.Int64("initial_offset", claim.InitialOffset()))

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				h.logger.Warn("Message channel was closed", 
					zap.String("topic", claim.Topic()),
					zap.Int32("partition", claim.Partition()))
				return nil
			}

			start := time.Now()
			h.logger.Debug("Processing message",
				zap.String("topic", message.Topic),
				zap.Int32("partition", message.Partition),
				zap.Int64("offset", message.Offset),
				zap.Int("message_size", len(message.Value)))

			err := h.Process(message)
			duration := time.Since(start)
			
			if err != nil {
				atomic.AddInt64(&h.metrics.ProcessingErrors, 1)
				h.logger.Error("Error processing message",
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.Error(err),
					zap.Duration("processing_duration", duration))
			} else {
				atomic.AddInt64(&h.metrics.MessagesConsumed, 1)
				h.metrics.LastConsumedTime = time.Now()
				session.MarkMessage(message, "")
				
				h.logger.Info("Message processed successfully",
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset),
					zap.Duration("processing_duration", duration),
					zap.Int("message_size", len(message.Value)))
			}
		case <-session.Context().Done():
			h.logger.Info("Consumer context cancelled, stopping consumption")
			return nil
		}
	}
}

func (h *ConsumerGroupHandler) GetMetrics() ConsumerMetrics {
	return ConsumerMetrics{
		MessagesConsumed: atomic.LoadInt64(&h.metrics.MessagesConsumed),
		ProcessingErrors: atomic.LoadInt64(&h.metrics.ProcessingErrors),
		LastConsumedTime: h.metrics.LastConsumedTime,
	}
}

func StartConsumer(ctx context.Context, brokers []string, groupID string, topics []string, processMsg func(*sarama.ConsumerMessage) error) error {
	return StartConsumerWithLogger(ctx, brokers, groupID, topics, processMsg, nil)
}

func StartConsumerWithLogger(ctx context.Context, brokers []string, groupID string, topics []string, processMsg func(*sarama.ConsumerMessage) error, logger *zap.Logger) error {
	if logger == nil {
		logger, _ = zap.NewProduction()
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Session.Timeout = 10 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	
	// Enable Sarama logging
	sarama.Logger = log.New(log.Writer(), "[Sarama] ", log.LstdFlags)

	logger.Info("Initializing Kafka consumer",
		zap.Strings("brokers", brokers),
		zap.String("group_id", groupID),
		zap.Strings("topics", topics))

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		logger.Error("Failed to create consumer group client",
			zap.Error(err),
			zap.Strings("brokers", brokers),
			zap.String("group_id", groupID))
		return fmt.Errorf("error creating consumer group client: %w", err)
	}
	defer client.Close()

	logger.Info("Consumer group client created successfully")

	handler := NewConsumerGroupHandler(logger, processMsg)

	// Wait for the consumer to be ready
	select {
	case <-handler.Ready:
		logger.Info("Consumer is ready to start consuming messages")
	case <-time.After(10 * time.Second):
		logger.Warn("Consumer setup took longer than expected")
	case <-ctx.Done():
		logger.Info("Context cancelled before consumer was ready")
		return ctx.Err()
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("Context cancelled, stopping consumer")
			return nil
		default:
			logger.Debug("Starting to consume messages")
			err := client.Consume(ctx, topics, handler)
			if err != nil {
				logger.Error("Error from consumer", zap.Error(err))
				return fmt.Errorf("error from consumer: %w", err)
			}
			
			if ctx.Err() != nil {
				return nil
			}
			
			logger.Debug("Consumer session completed, restarting")
			handler.Ready = make(chan bool)
		}
	}
}
