package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

// MockKafkaProducer implements sarama.SyncProducer for testing
type MockKafkaProducer struct {
	messages []*sarama.ProducerMessage
	err      error
}

func (m *MockKafkaProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	if m.err != nil {
		return 0, 0, m.err
	}
	m.messages = append(m.messages, msg)
	return 0, 0, nil
}

func (m *MockKafkaProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	return nil
}

func (m *MockKafkaProducer) Close() error {
	return nil
}

// TestOrder_Struct tests Order struct creation
func TestOrder_Struct(t *testing.T) {
	order := Order{
		ID:          uuid.New().String(),
		UserID:      uuid.New().String(),
		ProductName: "Test Product",
		Price:       999,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	if order.ID == "" {
		t.Error("Order ID should not be empty")
	}
	if order.ProductName != "Test Product" {
		t.Errorf("Expected product name 'Test Product', got '%s'", order.ProductName)
	}
	if order.Price != 999 {
		t.Errorf("Expected price 999, got %d", order.Price)
	}
	if order.Status != "pending" {
		t.Errorf("Expected status 'pending', got '%s'", order.Status)
	}
}

// TestCreateOrderRequest_Validation tests request validation
func TestCreateOrderRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateOrderRequest
		wantErr bool
	}{
		{
			name: "Valid request",
			req: CreateOrderRequest{
				UserID:      uuid.New().String(),
				ProductName: "Laptop",
				Price:       1000,
			},
			wantErr: false,
		},
		{
			name: "Zero price",
			req: CreateOrderRequest{
				UserID:      uuid.New().String(),
				ProductName: "Free Item",
				Price:       0,
			},
			wantErr: false,
		},
		{
			name: "Negative price",
			req: CreateOrderRequest{
				UserID:      uuid.New().String(),
				ProductName: "Discount",
				Price:       -100,
			},
			wantErr: false,
		},
		{
			name: "Empty product name",
			req: CreateOrderRequest{
				UserID:      uuid.New().String(),
				ProductName: "",
				Price:       100,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the struct can be created
			if tt.req.UserID == "" {
				t.Skip("Skipping empty user ID test")
			}
		})
	}
}

// TestOrderCreatedEvent tests event struct
func TestOrderCreatedEvent(t *testing.T) {
	event := OrderCreatedEvent{
		OrderID:     uuid.New().String(),
		UserID:      uuid.New().String(),
		ProductName: "Gaming Laptop",
		Price:       1299,
		CreatedAt:   time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	var decoded OrderCreatedEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if decoded.OrderID != event.OrderID {
		t.Error("OrderID mismatch after unmarshal")
	}
	if decoded.Price != event.Price {
		t.Errorf("Price mismatch: expected %d, got %d", event.Price, decoded.Price)
	}
}

// TestUserCreatedEvent tests user event struct
func TestUserCreatedEvent(t *testing.T) {
	event := UserCreatedEvent{
		UserID:    uuid.New().String(),
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal user event: %v", err)
	}

	var decoded UserCreatedEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal user event: %v", err)
	}

	if decoded.Name != event.Name {
		t.Errorf("Name mismatch: expected '%s', got '%s'", event.Name, decoded.Name)
	}
}

// TestUserCache tests the user cache functionality
func TestUserCache(t *testing.T) {
	server := &Server{
		userCache: make(map[string]UserCreatedEvent),
	}

	// Add users to cache
	user1 := UserCreatedEvent{
		UserID:    "user-1",
		Name:      "User One",
		Email:     "user1@example.com",
		CreatedAt: time.Now(),
	}
	server.userCache["user-1"] = user1

	user2 := UserCreatedEvent{
		UserID:    "user-2",
		Name:      "User Two",
		Email:     "user2@example.com",
		CreatedAt: time.Now(),
	}
	server.userCache["user-2"] = user2

	// Verify cache size
	if len(server.userCache) != 2 {
		t.Errorf("Expected 2 users in cache, got %d", len(server.userCache))
	}

	// Verify retrieval
	if cached, ok := server.userCache["user-1"]; !ok {
		t.Error("User-1 not found in cache")
	} else if cached.Name != "User One" {
		t.Errorf("Expected 'User One', got '%s'", cached.Name)
	}
}

// TestCreateOrderHandler_MethodNotAllowed tests HTTP method validation
func TestCreateOrderHandler_MethodNotAllowed(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: &MockKafkaProducer{},
		userCache:     make(map[string]UserCreatedEvent),
	}

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	rr := httptest.NewRecorder()

	server.handleCreateOrder(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

// TestCreateOrderHandler_InvalidBody tests invalid JSON handling
func TestCreateOrderHandler_InvalidBody(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: &MockKafkaProducer{},
		userCache:     make(map[string]UserCreatedEvent),
	}

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.handleCreateOrder(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", rr.Code)
	}
}

// TestGetOrderHandler_MethodNotAllowed tests method validation for GET
func TestGetOrderHandler_MethodNotAllowed(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: &MockKafkaProducer{},
		userCache:     make(map[string]UserCreatedEvent),
	}

	req := httptest.NewRequest(http.MethodPost, "/orders/123", nil)
	rr := httptest.NewRecorder()

	server.handleGetOrder(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

// TestGetOrderHandler_MissingID tests missing order ID
func TestGetOrderHandler_MissingID(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: &MockKafkaProducer{},
		userCache:     make(map[string]UserCreatedEvent),
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/", nil)
	rr := httptest.NewRecorder()

	server.handleGetOrder(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing ID, got %d", rr.Code)
	}
}

// TestHealthHandler_Response tests health check response
func TestHealthHandler_Response(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: nil,
		userCache:     make(map[string]UserCreatedEvent),
	}

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	server.handleHealth(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	if response["service"] != "order-service" {
		t.Errorf("Expected service 'order-service', got '%v'", response["service"])
	}

	if _, ok := response["cached_users"]; !ok {
		t.Error("Expected 'cached_users' field in health response")
	}
}

// TestStatsHandler_Response tests stats endpoint
func TestStatsHandler_Response(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: nil,
		userCache:     make(map[string]UserCreatedEvent),
	}

	req := httptest.NewRequest(http.MethodGet, "/internal/stats", nil)
	rr := httptest.NewRecorder()

	server.handleStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse stats response: %v", err)
	}

	if response["service"] != "order-service" {
		t.Errorf("Expected service 'order-service', got '%v'", response["service"])
	}

	if _, ok := response["cachedUsers"]; !ok {
		t.Error("Expected 'cachedUsers' field in stats response")
	}
}

// TestCacheHandler_Response tests cache viewer endpoint
func TestCacheHandler_Response(t *testing.T) {
	server := &Server{
		userCache: map[string]UserCreatedEvent{
			"user-1": {
				UserID: "user-1",
				Name:   "Test User",
				Email:  "test@example.com",
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/internal/cache", nil)
	rr := httptest.NewRecorder()

	server.handleCacheView(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	var response map[string]UserCreatedEvent
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse cache response: %v", err)
	}

	if len(response) != 1 {
		t.Errorf("Expected 1 user in cache response, got %d", len(response))
	}
}

// TestConsumerGroupHandler tests the Kafka consumer handler
func TestConsumerGroupHandler(t *testing.T) {
	server := &Server{
		userCache: make(map[string]UserCreatedEvent),
	}

	handler := &ConsumerGroupHandler{server: server}

	// Test Setup
	if err := handler.Setup(nil); err != nil {
		t.Errorf("Setup should not return error, got %v", err)
	}

	// Test Cleanup
	if err := handler.Cleanup(nil); err != nil {
		t.Errorf("Cleanup should not return error, got %v", err)
	}
}

// TestMetrics_Registration tests that metrics are registered
func TestMetrics_Registration(t *testing.T) {
	reg := prometheus.NewRegistry()

	reg.MustRegister(httpRequestsTotal)
	reg.MustRegister(httpRequestDuration)
	reg.MustRegister(dbQueriesTotal)
	reg.MustRegister(kafkaMessagesProduced)
	reg.MustRegister(kafkaMessagesConsumed)
	reg.MustRegister(activeConnections)
	reg.MustRegister(usersCached)

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	expectedMetrics := []string{
		"order_service_http_requests_total",
		"order_service_http_request_duration_seconds",
		"order_service_db_queries_total",
		"order_service_kafka_messages_produced_total",
		"order_service_kafka_messages_consumed_total",
		"order_service_active_connections",
		"order_service_users_cached",
	}

	foundCount := 0
	for _, expected := range expectedMetrics {
		for _, mf := range metrics {
			if mf.GetName() == expected {
				foundCount++
				break
			}
		}
	}

	if foundCount != len(expectedMetrics) {
		t.Errorf("Expected %d metrics, found %d", len(expectedMetrics), foundCount)
	}
}

// TestPrometheusCounterVec tests counter vector operations
func TestPrometheusCounterVec(t *testing.T) {
	reg := prometheus.NewRegistry()
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_requests_total",
			Help: "Test requests",
		},
		[]string{"method", "status"},
	)
	reg.MustRegister(counter)

	// Increment with different labels
	counter.WithLabelValues("GET", "200").Inc()
	counter.WithLabelValues("GET", "200").Inc()
	counter.WithLabelValues("POST", "201").Inc()

	metrics, _ := reg.Gather()

	var foundMetrics int
	for _, mf := range metrics {
		if mf.GetName() == "test_requests_total" {
			foundMetrics = len(mf.GetMetric())
		}
	}

	if foundMetrics != 2 { // Should have 2 different label combinations
		t.Errorf("Expected 2 metric variations, got %d", foundMetrics)
	}
}

// TestPrometheusHistogram tests histogram operations
func TestPrometheusHistogram(t *testing.T) {
	reg := prometheus.NewRegistry()
	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "test_duration_seconds",
			Help:    "Test duration",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)
	reg.MustRegister(histogram)

	// Observe values
	histogram.WithLabelValues("/users").Observe(0.1)
	histogram.WithLabelValues("/users").Observe(0.2)
	histogram.WithLabelValues("/orders").Observe(0.3)

	metrics, _ := reg.Gather()

	for _, mf := range metrics {
		if mf.GetName() == "test_duration_seconds" {
			if len(mf.GetMetric()) != 2 {
				t.Errorf("Expected 2 histograms, got %d", len(mf.GetMetric()))
			}
		}
	}
}

// TestMockKafkaProducer_SendMessage tests mock producer
func TestMockKafkaProducer_SendMessage(t *testing.T) {
	mock := &MockKafkaProducer{}

	event := OrderCreatedEvent{
		OrderID:     "order-123",
		UserID:      "user-123",
		ProductName: "Test Product",
		Price:       100,
		CreatedAt:   time.Now(),
	}

	data, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "order.created",
		Key:   sarama.StringEncoder("order-123"),
		Value: sarama.ByteEncoder(data),
	}

	partition, offset, err := mock.SendMessage(msg)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if partition != 0 || offset != 0 {
		t.Error("Expected partition=0, offset=0")
	}

	if len(mock.messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(mock.messages))
	}
}

// BenchmarkCreateOrder_Baseline benchmarks order creation
func BenchmarkCreateOrder_Baseline(b *testing.B) {
	mock := &MockKafkaProducer{}
	server := &Server{
		db:            nil,
		kafkaProducer: mock,
		userCache:     make(map[string]UserCreatedEvent),
	}

	body := []byte(`{"user_id": "user-123", "product_name": "Test Product", "price": 100}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.handleCreateOrder(rr, req)
	}
}

// BenchmarkJSONMarshal_Order benchmarks order marshaling
func BenchmarkJSONMarshal_Order(b *testing.B) {
	order := Order{
		ID:          uuid.New().String(),
		UserID:      uuid.New().String(),
		ProductName: "Gaming Laptop",
		Price:       1299,
		Status:      "pending",
		CreatedAt:   time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(order)
	}
}

// BenchmarkUserCache_Access benchmarks cache access
func BenchmarkUserCache_Access(b *testing.B) {
	server := &Server{
		userCache: make(map[string]UserCreatedEvent),
	}

	// Populate cache
	for i := 0; i < 100; i++ {
		server.userCache[string(rune(i))] = UserCreatedEvent{
			UserID: string(rune(i)),
			Name:   "User",
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = server.userCache["50"]
	}
}

// Helper function to get counter value
func getCounterValue(counter prometheus.Counter) float64 {
	dto := &dto.Metric{}
	counter.Write(dto)
	return dto.GetCounter().GetValue()
}
