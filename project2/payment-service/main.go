package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Payment represents a payment in the system
type Payment struct {
	ID            string    `json:"id" bson:"_id"`
	OrderID       string    `json:"order_id" bson:"order_id"`
	UserID        string    `json:"user_id" bson:"user_id"`
	Amount        int       `json:"amount" bson:"amount"`
	PaymentStatus string    `json:"payment_status" bson:"payment_status"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
}

// OrderCreatedEvent represents the order.created Kafka event
type OrderCreatedEvent struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	ProductName string    `json:"product_name"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// PaymentCompletedEvent represents the payment.completed Kafka event
type PaymentCompletedEvent struct {
	PaymentID string    `json:"payment_id"`
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	Amount    int       `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "payment_service_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	dbQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	kafkaMessagesProduced = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_kafka_messages_produced_total",
			Help: "Total number of Kafka messages produced",
		},
		[]string{"topic", "status"},
	)

	kafkaMessagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic", "status"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "payment_service_active_connections",
			Help: "Number of active connections",
		},
	)

	paymentsProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "payment_service_payments_processed_total",
			Help: "Total number of payments processed",
		},
		[]string{"status"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(dbQueriesTotal)
	prometheus.MustRegister(kafkaMessagesProduced)
	prometheus.MustRegister(kafkaMessagesConsumed)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(paymentsProcessed)
}

type Server struct {
	mongoClient   *mongo.Client
	database      *mongo.Database
	paymentsColl  *mongo.Collection
	kafkaProducer sarama.SyncProducer
	ctx           context.Context
}

func main() {
	log.Println("Starting Payment Service...")

	ctx := context.Background()

	// Get environment variables
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}
	log.Println("Connected to MongoDB")

	database := client.Database("paymentdb")
	paymentsColl := database.Collection("payments")

	// Create indexes
	createIndexes(ctx, paymentsColl)

	// Connect to Kafka
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll

	kafkaProducer, err := sarama.NewSyncProducer([]string{kafkaBrokers}, kafkaConfig)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()
	log.Println("Connected to Kafka producer")

	server := &Server{
		mongoClient:   client,
		database:      database,
		paymentsColl:  paymentsColl,
		kafkaProducer: kafkaProducer,
		ctx:           ctx,
	}

	// Start Kafka consumer in background
	go server.startKafkaConsumer(kafkaBrokers)

	// Setup HTTP handlers
	http.HandleFunc("/payments", server.handleGetPayments)
	http.HandleFunc("/payments/", server.handleGetPayment)
	http.HandleFunc("/health", server.handleHealth)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/internal/stats", server.handleStats)

	log.Printf("Payment Service running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createIndexes(ctx context.Context, coll *mongo.Collection) {
	// Create index on order_id
	orderIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "order_id", Value: 1}},
		Options: options.Index().SetUnique(false),
	}
	// Create index on user_id
	userIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "user_id", Value: 1}},
		Options: options.Index().SetUnique(false),
	}
	// Create index on payment_status
	statusIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "payment_status", Value: 1}},
	}

	_, err := coll.Indexes().CreateMany(ctx, []mongo.IndexModel{orderIndex, userIndex, statusIndex})
	if err != nil {
		log.Printf("Failed to create indexes: %v", err)
	} else {
		log.Println("Created MongoDB indexes")
	}
}

func (s *Server) startKafkaConsumer(brokers string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup([]string{brokers}, "payment-service-group", config)
	if err != nil {
		log.Printf("Failed to create Kafka consumer group: %v", err)
		return
	}
	defer consumer.Close()

	log.Println("Kafka consumer started, listening for order.created events...")

	handler := &ConsumerGroupHandler{server: s}
	ctx := context.Background()
	for {
		err := consumer.Consume(ctx, []string{"order.created"}, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}
	}
}

// ConsumerGroupHandler handles Kafka messages
type ConsumerGroupHandler struct {
	server *Server
}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s",
			string(message.Value), message.Timestamp, message.Topic)

		var event OrderCreatedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Failed to unmarshal order.created event: %v", err)
			kafkaMessagesConsumed.WithLabelValues("order.created", "error").Inc()
			continue
		}

		// Process payment
		h.server.processPayment(event)
		kafkaMessagesConsumed.WithLabelValues("order.created", "success").Inc()

		session.MarkMessage(message, "")
	}
	return nil
}

func (s *Server) processPayment(event OrderCreatedEvent) {
	log.Printf("Processing payment for order %s (amount: %d)", event.OrderID, event.Price)

	// Simulate payment processing
	time.Sleep(100 * time.Millisecond)

	// Determine payment status (95% success rate for simulation)
	status := "success"
	if event.Price%100 == 99 { // Fail payments where price ends in 99
		status = "failed"
	}

	payment := Payment{
		ID:            uuid.New().String(),
		OrderID:       event.OrderID,
		UserID:        event.UserID,
		Amount:        event.Price,
		PaymentStatus: status,
		CreatedAt:     time.Now(),
	}

	// Store in MongoDB
	dbStart := time.Now()
	_, err := s.paymentsColl.InsertOne(s.ctx, payment)
	dbQueriesTotal.WithLabelValues("INSERT", fmt.Sprintf("%t", err == nil)).Inc()

	if err != nil {
		log.Printf("Failed to insert payment: %v", err)
		paymentsProcessed.WithLabelValues("error").Inc()
		return
	}
	log.Printf("DB insert duration: %v", time.Since(dbStart))

	// Publish Kafka event
	completedEvent := PaymentCompletedEvent{
		PaymentID: payment.ID,
		OrderID:   payment.OrderID,
		UserID:    payment.UserID,
		Amount:    payment.Amount,
		Status:    payment.PaymentStatus,
		CreatedAt: payment.CreatedAt,
	}

	eventBytes, err := json.Marshal(completedEvent)
	if err != nil {
		log.Printf("Failed to marshal payment event: %v", err)
	} else {
		msg := &sarama.ProducerMessage{
			Topic: "payment.completed",
			Key:   sarama.StringEncoder(payment.ID),
			Value: sarama.ByteEncoder(eventBytes),
		}

		_, _, err = s.kafkaProducer.SendMessage(msg)
		kafkaMessagesProduced.WithLabelValues("payment.completed", fmt.Sprintf("%t", err == nil)).Inc()

		if err != nil {
			log.Printf("Failed to produce Kafka message: %v", err)
		} else {
			log.Printf("Published payment.completed event for payment %s (status: %s)", payment.ID, status)
		}
	}

	paymentsProcessed.WithLabelValues(status).Inc()
	log.Printf("Payment %s processed with status: %s", payment.ID, status)
}

func (s *Server) handleGetPayments(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	activeConnections.Inc()
	defer activeConnections.Dec()

	if r.Method != http.MethodGet {
		httpRequestsTotal.WithLabelValues("GET", "/payments", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dbStart := time.Now()
	cursor, err := s.paymentsColl.Find(s.ctx, bson.M{})
	dbQueriesTotal.WithLabelValues("FIND", fmt.Sprintf("%t", err == nil)).Inc()

	if err != nil {
		log.Printf("Failed to query payments: %v", err)
		httpRequestsTotal.WithLabelValues("GET", "/payments", "500").Inc()
		http.Error(w, "Failed to get payments", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(s.ctx)

	var payments []Payment
	if err := cursor.All(s.ctx, &payments); err != nil {
		log.Printf("Failed to decode payments: %v", err)
		httpRequestsTotal.WithLabelValues("GET", "/payments", "500").Inc()
		http.Error(w, "Failed to get payments", http.StatusInternalServerError)
		return
	}
	log.Printf("DB query duration: %v", time.Since(dbStart))

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("GET", "/payments").Observe(duration)
	httpRequestsTotal.WithLabelValues("GET", "/payments", "200").Inc()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)

	log.Printf("Retrieved %d payments", len(payments))
}

func (s *Server) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	activeConnections.Inc()
	defer activeConnections.Dec()

	if r.Method != http.MethodGet {
		httpRequestsTotal.WithLabelValues("GET", "/payments/{id}", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	paymentID := r.URL.Path[len("/payments/"):]
	if paymentID == "" {
		httpRequestsTotal.WithLabelValues("GET", "/payments/{id}", "400").Inc()
		http.Error(w, "Payment ID required", http.StatusBadRequest)
		return
	}

	var payment Payment
	dbStart := time.Now()
	err := s.paymentsColl.FindOne(s.ctx, bson.M{"_id": paymentID}).Decode(&payment)
	dbQueriesTotal.WithLabelValues("FIND_ONE", fmt.Sprintf("%t", err == nil)).Inc()
	log.Printf("DB query duration: %v", time.Since(dbStart))

	if err == mongo.ErrNoDocuments {
		httpRequestsTotal.WithLabelValues("GET", "/payments/{id}", "404").Inc()
		http.Error(w, "Payment not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to query payment: %v", err)
		httpRequestsTotal.WithLabelValues("GET", "/payments/{id}", "500").Inc()
		http.Error(w, "Failed to get payment", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("GET", "/payments/{id}").Observe(duration)
	httpRequestsTotal.WithLabelValues("GET", "/payments/{id}", "200").Inc()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payment)

	log.Printf("Retrieved payment: %s", payment.ID)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"service":   "payment-service",
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Check MongoDB
	dbStatus := "connected"
	if err := s.mongoClient.Ping(s.ctx, nil); err != nil {
		dbStatus = "disconnected"
		health["status"] = "unhealthy"
	}
	health["database"] = dbStatus

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"service":           "payment-service",
		"timestamp":         time.Now().Format(time.RFC3339),
		"activeConnections": 0,
	}

	// Get payment count
	count, err := s.paymentsColl.CountDocuments(s.ctx, bson.M{})
	if err != nil {
		log.Printf("Failed to count payments: %v", err)
	}
	stats["totalPayments"] = count

	// Get count by status
	statusCounts := make(map[string]int64)
	for _, status := range []string{"success", "failed", "pending"} {
		cnt, _ := s.paymentsColl.CountDocuments(s.ctx, bson.M{"payment_status": status})
		statusCounts[status] = cnt
	}
	stats["paymentsByStatus"] = statusCounts

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
