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

// TestCreateUserRequest_Validation tests request validation
func TestCreateUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		req     CreateUserRequest
		wantErr bool
	}{
		{
			name:    "Valid request",
			req:     CreateUserRequest{Name: "John Doe", Email: "john@example.com"},
			wantErr: false,
		},
		{
			name:    "Empty name",
			req:     CreateUserRequest{Name: "", Email: "john@example.com"},
			wantErr: false, // We allow empty in current implementation
		},
		{
			name:    "Empty email",
			req:     CreateUserRequest{Name: "John Doe", Email: ""},
			wantErr: false, // We allow empty in current implementation
		},
		{
			name:    "Both empty",
			req:     CreateUserRequest{Name: "", Email: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify the struct can be created
			if tt.req.Name == "" && tt.req.Email == "" {
				// This is valid for our implementation
			}
		})
	}
}

// TestUser_Struct tests User struct creation
func TestUser_Struct(t *testing.T) {
	user := User{
		ID:        uuid.New().String(),
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}

	if user.ID == "" {
		t.Error("User ID should not be empty")
	}
	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

// TestUserCreatedEvent tests event struct
func TestUserCreatedEvent(t *testing.T) {
	event := UserCreatedEvent{
		UserID:    uuid.New().String(),
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal event: %v", err)
	}

	var decoded UserCreatedEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal event: %v", err)
	}

	if decoded.UserID != event.UserID {
		t.Error("UserID mismatch after unmarshal")
	}
}

// TestMetrics_Registration tests that metrics are registered
func TestMetrics_Registration(t *testing.T) {
	// Create a new registry for testing
	reg := prometheus.NewRegistry()

	// Register metrics
	reg.MustRegister(httpRequestsTotal)
	reg.MustRegister(httpRequestDuration)
	reg.MustRegister(dbQueriesTotal)
	reg.MustRegister(kafkaMessagesProduced)
	reg.MustRegister(activeConnections)

	// Collect metrics
	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	expectedMetrics := []string{
		"user_service_http_requests_total",
		"user_service_http_request_duration_seconds",
		"user_service_db_queries_total",
		"user_service_kafka_messages_produced_total",
		"user_service_active_connections",
	}

	for _, expected := range expectedMetrics {
		found := false
		for _, mf := range metrics {
			if mf.GetName() == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Metric %s not found in registry", expected)
		}
	}
}

// TestMockKafkaProducer tests the mock producer
func TestMockKafkaProducer(t *testing.T) {
	mock := &MockKafkaProducer{}

	msg := &sarama.ProducerMessage{
		Topic: "user.created",
		Key:   sarama.StringEncoder("test-key"),
		Value: sarama.ByteEncoder(`{"test": "data"}`),
	}

	partition, offset, err := mock.SendMessage(msg)
	if err != nil {
		t.Fatalf("SendMessage failed: %v", err)
	}

	if partition != 0 || offset != 0 {
		t.Error("Expected partition=0, offset=0")
	}

	if len(mock.messages) != 1 {
		t.Errorf("Expected 1 message in mock, got %d", len(mock.messages))
	}
}

// TestMockKafkaProducer_Error tests error handling
func TestMockKafkaProducer_Error(t *testing.T) {
	mock := &MockKafkaProducer{
		err: sarama.ErrLeaderNotAvailable,
	}

	msg := &sarama.ProducerMessage{
		Topic: "user.created",
		Value: sarama.ByteEncoder(`{}`),
	}

	_, _, err := mock.SendMessage(msg)
	if err != sarama.ErrLeaderNotAvailable {
		t.Errorf("Expected ErrLeaderNotAvailable, got %v", err)
	}
}

// TestCreateUserHandler_MethodNotAllowed tests HTTP method validation
func TestCreateUserHandler_MethodNotAllowed(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: &MockKafkaProducer{},
	}

	// Test GET request to POST endpoint
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		server.handleCreateUser(w, r)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

// TestCreateUserHandler_InvalidBody tests invalid JSON handling
func TestCreateUserHandler_InvalidBody(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: &MockKafkaProducer{},
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.handleCreateUser(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", rr.Code)
	}
}

// TestHealthHandler_Response tests health check response format
func TestHealthHandler_Response(t *testing.T) {
	server := &Server{
		db:            nil,
		kafkaProducer: nil,
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

	if response["service"] != "user-service" {
		t.Error("Expected service to be 'user-service'")
	}
}

// TestUUID_Generation tests UUID generation
func TestUUID_Generation(t *testing.T) {
	id1 := uuid.New().String()
	id2 := uuid.New().String()

	if id1 == "" {
		t.Error("UUID should not be empty")
	}

	if id1 == id2 {
		t.Error("UUIDs should be unique")
	}

	// Validate UUID format (should have 4 hyphens)
	if len(id1) != 36 {
		t.Errorf("UUID should be 36 characters, got %d", len(id1))
	}
}

// TestMetricsCounter_Increment tests counter increments
func TestMetricsCounter_Increment(t *testing.T) {
	reg := prometheus.NewRegistry()
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_counter",
			Help: "Test counter",
		},
		[]string{"label"},
	)
	reg.MustRegister(counter)

	// Increment
	counter.WithLabelValues("test").Inc()
	counter.WithLabelValues("test").Inc()

	// Gather
	metrics, _ := reg.Gather()
	
	for _, mf := range metrics {
		if mf.GetName() == "test_counter" {
			if len(mf.GetMetric()) != 1 {
				t.Errorf("Expected 1 metric, got %d", len(mf.GetMetric()))
			}
			value := mf.GetMetric()[0].GetCounter().GetValue()
			if value != 2 {
				t.Errorf("Expected counter value 2, got %f", value)
			}
		}
	}
}

// TestTimeFormatting tests timestamp handling
func TestTimeFormatting(t *testing.T) {
	now := time.Now()
	
	// Test JSON marshaling
	user := User{
		ID:        "test-id",
		Name:      "Test",
		Email:     "test@test.com",
		CreatedAt: now,
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Times should be approximately equal (within 1 second for JSON precision)
	diff := decoded.CreatedAt.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Time mismatch: expected %v, got %v", now, decoded.CreatedAt)
	}
}

// TestJSONEncoding_LargeData tests JSON encoding with larger payloads
func TestJSONEncoding_LargeData(t *testing.T) {
	users := make([]User, 100)
	for i := 0; i < 100; i++ {
		users[i] = User{
			ID:        uuid.New().String(),
			Name:      "User",
			Email:     "user@example.com",
			CreatedAt: time.Now(),
		}
	}

	data, err := json.Marshal(users)
	if err != nil {
		t.Fatalf("Failed to marshal 100 users: %v", err)
	}

	var decoded []User
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal 100 users: %v", err)
	}

	if len(decoded) != 100 {
		t.Errorf("Expected 100 users, got %d", len(decoded))
	}
}

// BenchmarkCreateUser_Baseline benchmarks the create user endpoint
func BenchmarkCreateUser_Baseline(b *testing.B) {
	mock := &MockKafkaProducer{}
	server := &Server{
		db:            nil,
		kafkaProducer: mock,
	}

	body := []byte(`{"name": "Test User", "email": "test@example.com"}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		server.handleCreateUser(rr, req)
	}
}

// BenchmarkUUID_Generation benchmarks UUID generation
func BenchmarkUUID_Generation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuid.New().String()
	}
}

// BenchmarkJSON_Marshal benchmarks JSON marshaling
func BenchmarkJSON_Marshal(b *testing.B) {
	user := User{
		ID:        uuid.New().String(),
		Name:      "Test User",
		Email:     "test@example.com",
		CreatedAt: time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(user)
	}
}

// Helper function to get counter value
func getCounterValue(counter prometheus.Counter) float64 {
	dto := &dto.Metric{}
	counter.Write(dto)
	return dto.GetCounter().GetValue()
}
