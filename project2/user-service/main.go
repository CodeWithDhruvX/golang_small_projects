package main

import (
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

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
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
			Name: "user_service_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_service_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	dbQueriesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_service_db_queries_total",
			Help: "Total number of database queries",
		},
		[]string{"operation", "status"},
	)

	kafkaMessagesProduced = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "user_service_kafka_messages_produced_total",
			Help: "Total number of Kafka messages produced",
		},
		[]string{"topic", "status"},
	)

	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "user_service_active_connections",
			Help: "Number of active connections",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(dbQueriesTotal)
	prometheus.MustRegister(kafkaMessagesProduced)
	prometheus.MustRegister(activeConnections)
}

type Server struct {
	db            *sql.DB
	kafkaProducer sarama.SyncProducer
}

func main() {
	log.Println("Starting User Service...")

	// Get environment variables
	dbURL := os.Getenv("USER_DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/userdb?sslmode=disable"
	}

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		kafkaBrokers = "localhost:9092"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
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
	log.Println("Connected to Kafka")

	server := &Server{
		db:            db,
		kafkaProducer: kafkaProducer,
	}

	// Setup HTTP handlers
	http.HandleFunc("/users", server.handleCreateUser)
	http.HandleFunc("/users/", server.handleGetUser)
	http.HandleFunc("/health", server.handleHealth)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/internal/stats", server.handleStats)

	log.Printf("User Service running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL UNIQUE,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := db.Exec(query)
	return err
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	activeConnections.Inc()
	defer activeConnections.Dec()

	if r.Method != http.MethodPost {
		httpRequestsTotal.WithLabelValues("POST", "/users", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpRequestsTotal.WithLabelValues("POST", "/users", "400").Inc()
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user := User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	// Insert into database
	dbStart := time.Now()
	_, err := s.db.Exec(
		"INSERT INTO users (id, name, email, created_at) VALUES ($1, $2, $3, $4)",
		user.ID, user.Name, user.Email, user.CreatedAt,
	)
	dbQueriesTotal.WithLabelValues("INSERT", fmt.Sprintf("%t", err == nil)).Inc()

	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		httpRequestsTotal.WithLabelValues("POST", "/users", "500").Inc()
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	log.Printf("DB query duration: %v", time.Since(dbStart))

	// Publish Kafka event
	event := UserCreatedEvent{
		UserID:    user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("Failed to marshal event: %v", err)
	} else {
		msg := &sarama.ProducerMessage{
			Topic: "user.created",
			Key:   sarama.StringEncoder(user.ID),
			Value: sarama.ByteEncoder(eventBytes),
		}

		_, _, err = s.kafkaProducer.SendMessage(msg)
		kafkaMessagesProduced.WithLabelValues("user.created", fmt.Sprintf("%t", err == nil)).Inc()

		if err != nil {
			log.Printf("Failed to produce Kafka message: %v", err)
		} else {
			log.Printf("Published user.created event for user %s", user.ID)
		}
	}

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("POST", "/users").Observe(duration)
	httpRequestsTotal.WithLabelValues("POST", "/users", "201").Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)

	log.Printf("Created user: %s (%s) in %v", user.Name, user.ID, time.Since(start))
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	activeConnections.Inc()
	defer activeConnections.Dec()

	if r.Method != http.MethodGet {
		httpRequestsTotal.WithLabelValues("GET", "/users/{id}", "405").Inc()
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Path[len("/users/"):]
	if userID == "" {
		httpRequestsTotal.WithLabelValues("GET", "/users/{id}", "400").Inc()
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var user User
	dbStart := time.Now()
	err := s.db.QueryRow(
		"SELECT id, name, email, created_at FROM users WHERE id = $1",
		userID,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	dbQueriesTotal.WithLabelValues("SELECT", fmt.Sprintf("%t", err == nil)).Inc()
	log.Printf("DB query duration: %v", time.Since(dbStart))

	if err == sql.ErrNoRows {
		httpRequestsTotal.WithLabelValues("GET", "/users/{id}", "404").Inc()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("Failed to query user: %v", err)
		httpRequestsTotal.WithLabelValues("GET", "/users/{id}", "500").Inc()
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	httpRequestDuration.WithLabelValues("GET", "/users/{id}").Observe(duration)
	httpRequestsTotal.WithLabelValues("GET", "/users/{id}", "200").Inc()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)

	log.Printf("Retrieved user: %s", user.ID)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"service":   "user-service",
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"service":           "user-service",
		"timestamp":         time.Now().Format(time.RFC3339),
		"activeConnections": getGaugeValue(activeConnections),
	}

	// Get user count
	var userCount int
	s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	stats["totalUsers"] = userCount

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func getGaugeValue(g prometheus.Gauge) float64 {
	// Return current value by collecting metrics
	return 0 // Simplified - prometheus client doesn't expose direct value reading
}
