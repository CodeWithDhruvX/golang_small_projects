package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Order represents an order in the system
type Order struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProductName string    `json:"product_name"`
	Price       int       `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateOrderRequest represents a request to create an order
type CreateOrderRequest struct {
	UserID      string `json:"user_id"`
	ProductName string `json:"product_name"`
	Price       int    `json:"price"`
}

// OrderCreatedEvent represents the order.created Kafka event
type OrderCreatedEvent struct {
	OrderID     string    `json:"order_id"`
	UserID      string    `json:"user_id"`
	ProductName string    `json:"product_name"`
	Price       int       `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

// UserCreatedEvent represents the user.created Kafka event
type UserCreatedEvent struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Metrics
var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_service_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	dbQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_service_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	kafkaMessagesProduced = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_service_kafka_messages_produced_total",
			Help: "Total number of Kafka messages produced",
		},
		[]string{"topic", "status"},
	)

	kafkaMessagesConsumed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_service_kafka_messages_consumed_total",
			Help: "Total number of Kafka messages consumed",
		},
		[]string{"topic", "status"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "order_service_active_connections",
			Help: "Number of active connections",
		},
	)

	usersCached = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "order_service_users_cached",
			Help: "Number of users in local cache from Kafka events",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(dbQueriesTotal)
	prometheus.MustRegister(kafkaMessagesProduced)
	prometheus.MustRegister(kafkaMessagesConsumed)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(usersCached)
}

type Server struct {
	db            *sql.DB
	kafkaProducer sarama.SyncProducer
	userCache     map[string]UserCreatedEvent
}

func main() {
	log.Println("Starting Order Service...")

	// Get environment variables
	dbURL := os.Getenv("ORDER_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/orderdb?sslmode=disable"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Connected to PostgreSQL")

	// Create table if not exists
	if err := createTable(db); err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

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
		db:            db,
		kafkaProducer: kafkaProducer,
		userCache:     make(map[string]UserCreatedEvent),
	}

	// Start Kafka consumer in background
	go server.startKafkaConsumer(kafkaBrokers)

	// Setup HTTP handlers
	http.HandleFunc("/orders", server.handleCreateOrder)
	http.HandleFunc("/orders/", server.handleGetOrder)
	http.HandleFunc("/health", server.handleHealth)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/internal/stats", server.handleStats)
	http.HandleFunc("/internal/cache", server.handleCacheView)

	log.Printf("Order Service running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id UUID PRIMARY KEY,
		user_id UUID NOT NULL,
		product_name VARCHAR(255) NOT NULL,
		price INTEGER NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'pending',
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func (s *Server) startKafkaConsumer(brokers string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup([]string{brokers}, "order-service-group", config)
	if err != nil {
		log.Printf("Failed to create Kafka consumer group: %v", err)
		return
	}
	defer consumer.Close()

	log.Println("Kafka consumer started, listening for user.created events...")

	handler := &ConsumerGroupHandler{server: s}
	ctx := context.Background()
	for {
		err := consumer.Consume(ctx, []string{"user.created"}, handler)
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

		var event UserCreatedEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Failed to unmarshal user.created event: %v", err)
			kafkaMessagesConsumed.WithLabelValues("user.created", "error").Inc()
			continue
		}

		// Store in cache
		h.server.userCache[event.UserID] = event
		usersCached.Set(float64(len(h.server.userCache)))
		kafkaMessagesConsumed.WithLabelValues("user.created", "success").Inc()

		log.Printf("Cached user from Kafka event: %s (%s)", event.Name, event.UserID)
		session.MarkMessage(message, "")
	}
	return nil
}

func (s *Server) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	activeConnections.Inc()
	defer activeConnections.Dec()

	if r.Method != http.MethodPost {
		httpRequestsTotal.WithLabelValues("POST", "/orders", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpRequestsTotal.WithLabelValues("POST", "/orders", "400").Inc()
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	order := Order{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		ProductName: req.ProductName,
		Price:       req.Price,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	// Insert into database
	dbStart := time.Now()
	_, err := s.db.Exec(
		"INSERT INTO orders (id, user_id, product_name, price, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)",
		order.ID, order.UserID, order.ProductName, order.Price, order.Status, order.CreatedAt,
	)
	dbQueriesTotal.WithLabelValues("INSERT", fmt.Sprintf("%t", err == nil)).Inc()

	if err != nil {
		log.Printf("Failed to insert order: %v", err)
		httpRequestsTotal.WithLabelValues("POST", "/orders", "500").Inc()
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}
	log.Printf("DB query duration: %v", time.Since(dbStart))

	// Publish Kafka event
	event := OrderCreatedEvent{
		OrderID:     order.ID,
		UserID:      order.UserID,
		ProductName: order.ProductName,
		Price:       order.Price,
		CreatedAt:   order.CreatedAt,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
	} else {
		msg := &sarama.ProducerMessage{
			Topic: "order.created",
			Key:   sarama.StringEncoder(order.ID),
			Value: sarama.ByteEncoder(eventBytes),
		}

		_, _, err = s.kafkaProducer.SendMessage(msg)
		kafkaMessagesProduced.WithLabelValues("order.created", fmt.Sprintf("%t", err == nil)).Inc()

		if err != nil {
			log.Printf("Failed to produce Kafka message: %v", err)
		} else {
			log.Printf("Published order.created event for order %s", order.ID)
		}
	}

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("POST", "/orders").Observe(duration)
	httpRequestsTotal.WithLabelValues("POST", "/orders", "201").Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)

	log.Printf("Created order: %s for user %s in %v", order.ID, order.UserID, time.Since(start))
}

func (s *Server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	activeConnections.Inc()
	defer activeConnections.Dec()

	if r.Method != http.MethodGet {
		httpRequestsTotal.WithLabelValues("GET", "/orders/{id}", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orderID := r.URL.Path[len("/orders/"):]
	if orderID == "" {
		httpRequestsTotal.WithLabelValues("GET", "/orders/{id}", "400").Inc()
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	var order Order
	dbStart := time.Now()
	err := s.db.QueryRow(
		"SELECT id, user_id, product_name, price, status, created_at FROM orders WHERE id = $1",
		orderID,
	).Scan(&order.ID, &order.UserID, &order.ProductName, &order.Price, &order.Status, &order.CreatedAt)
	dbQueriesTotal.WithLabelValues("SELECT", fmt.Sprintf("%t", err == nil)).Inc()
	log.Printf("DB query duration: %v", time.Since(dbStart))

	if err == sql.ErrNoRows {
		httpRequestsTotal.WithLabelValues("GET", "/orders/{id}", "404").Inc()
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to query order: %v", err)
		httpRequestsTotal.WithLabelValues("GET", "/orders/{id}", "500").Inc()
		http.Error(w, "Failed to get order", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("GET", "/orders/{id}").Observe(duration)
	httpRequestsTotal.WithLabelValues("GET", "/orders/{id}", "200").Inc()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)

	log.Printf("Retrieved order: %s", order.ID)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"service":   "order-service",
		"status":    "healthy",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	// Check database
	dbStatus := "connected"
	if err := s.db.Ping(); err != nil {
		dbStatus = "disconnected"
		health["status"] = "unhealthy"
	}
	health["database"] = dbStatus
	health["cached_users"] = len(s.userCache)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"service":           "order-service",
		"timestamp":         time.Now().Format(time.RFC3339),
		"activeConnections": 0,
	}

	// Get order count
	var orderCount int
	s.db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&orderCount)
	stats["totalOrders"] = orderCount
	stats["cachedUsers"] = len(s.userCache)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleCacheView(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.userCache)
}
