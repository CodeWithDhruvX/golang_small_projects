package main

import (
	"net/http"
	"os"
	"time"

	"github.com/dhruv/kafka-microservices/internal/api"
	"github.com/dhruv/kafka-microservices/internal/kafka"
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
	port := "8080"
	if p := os.Getenv("SERVICE_PORT"); p != "" {
		port = p
	}

	// Initialize monitoring components
	healthHandler := api.NewHealthHandler(logger)
	metricsHandler := api.NewMetricsHandler(logger)

	// Initialize Kafka Producer (or mock if Kafka is not available)
	var producer api.EventPublisher
	var err error
	
	// Try to connect to Kafka first
	realProducer, err := kafka.NewProducerWithLogger(brokers, logger)
	if err != nil {
		logger.Warn("Failed to connect to Kafka, using mock producer", zap.Error(err))
		mockProducer, _ := kafka.NewMockProducer()
		producer = mockProducer
		logger.Info("Using mock producer - orders will be logged but not sent to Kafka")
	} else {
		producer = realProducer
		logger.Info("Connected to Kafka successfully")
		// Register producer for health checks
		healthHandler.RegisterHealthCheck("kafka-producer", realProducer)
	}
	
	if closer, ok := producer.(interface{ Close() error }); ok {
		defer closer.Close()
	}

	// Set up HTTP handlers
	orderHandler := api.NewOrderHandler(logger, producer, topic)
	
	// Add middleware for metrics
	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer func() {
			latency := time.Since(start)
			success := true
			// Check if response status indicates error (this is basic, you might want to enhance it)
			metricsHandler.RecordRequest(success, latency)
		}()
		orderHandler.HandleCreateOrder(w, r)
	})

	// Health check endpoints
	http.HandleFunc("/health", healthHandler.HandleHealth)
	http.HandleFunc("/ready", healthHandler.HandleReady)
	http.HandleFunc("/live", healthHandler.HandleLive)
	
	// Metrics endpoints
	http.HandleFunc("/metrics", metricsHandler.HandleMetrics)
	http.HandleFunc("/metrics/kafka", metricsHandler.HandleKafkaMetrics)

	logger.Info("Starting Order Service", 
		zap.String("port", port),
		zap.Strings("kafka_brokers", brokers),
		zap.String("kafka_topic", topic))
	
	logger.Info("Available endpoints:",
		zap.String("orders", "/orders"),
		zap.String("health", "/health"),
		zap.String("ready", "/ready"),
		zap.String("live", "/live"),
		zap.String("metrics", "/metrics"),
		zap.String("kafka_metrics", "/metrics/kafka"))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
