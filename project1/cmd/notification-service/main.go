package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/dhruv/kafka-microservices/internal/kafka"
	"github.com/dhruv/kafka-microservices/internal/models"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load configuration
	brokers := []string{"localhost:9092"}
	if b := os.Getenv("KAFKA_BROKERS"); b != "" {
		brokers = []string{b}
	}
	topic := "orders.created"
	if t := os.Getenv("KAFKA_TOPIC_ORDERS"); t != "" {
		topic = t
	}
	groupID := "notification-service-group"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Process message function
	processMessage := func(msg *sarama.ConsumerMessage) error {
		var order models.Order
		if err := json.Unmarshal(msg.Value, &order); err != nil {
			logger.Error("Failed to unmarshal order event", zap.Error(err))
			return err
		}

		logger.Info("Notification triggered for new order",
			zap.String("order_id", order.ID),
			zap.String("customer_id", order.CustomerID),
			zap.Float64("total_amount", order.TotalAmount),
			zap.Int("items_count", len(order.Items)),
		)
		return nil
	}

	logger.Info("Starting Notification Service Consumer")

	// Try to start real consumer or fallback to mock
	go func() {
		if err := kafka.StartConsumer(ctx, brokers, groupID, []string{topic}, processMessage); err != nil {
			logger.Warn("Failed to start Kafka consumer, using mock consumer", zap.Error(err))
			
			mockConsumer, _ := kafka.NewMockConsumer()
			logger.Info("Using mock consumer - will simulate order notifications")
			
			if err := mockConsumer.StartMockConsumer(ctx, topic, processMessage); err != nil {
				logger.Fatal("Error running mock consumer", zap.Error(err))
			}
		}
	}()

	// Wait for termination signal
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm
	logger.Info("Shutting down Notification Service")
}
