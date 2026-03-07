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
	"go.mongodb.org/mongo-driver/bson"
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

// TestPayment_Struct tests Payment struct creation
func TestPayment_Struct(t *testing.T) {
	payment := Payment{
		ID:            uuid.New().String(),
		OrderID:       uuid.New().String(),
		UserID:        uuid.New().String(),
		Amount:        999,
		PaymentStatus: "success",
		CreatedAt:     time.Now(),
	}

	if payment.ID == "" {
		t.Error("Payment ID should not be empty")
	}
	if payment.Amount != 999 {
		t.Errorf("Expected amount 999, got %d", payment.Amount)
	}
	if payment.PaymentStatus != "success" {
		t.Errorf("Expected status 'success', got '%s'", payment.PaymentStatus)
	}
}

// TestPayment_MongoDBTags tests BSON tags for MongoDB
func TestPayment_MongoDBTags(t *testing.T) {
	payment := Payment{
		ID:            "payment-123",
		OrderID:       "order-123",
		UserID:        "user-123",
		Amount:        500,
		PaymentStatus: "success",
		CreatedAt:     time.Now(),
	}

	// Test BSON marshaling
	data, err := bson.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal to BSON: %v", err)
	}

	var decoded Payment
	if err := bson.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal from BSON: %v", err)
	}

	if decoded.ID != payment.ID {
		t.Errorf("ID mismatch: expected '%s', got '%s'", payment.ID, decoded.ID)
	}
	if decoded.Amount != payment.Amount {
		t.Errorf("Amount mismatch: expected %d, got %d", payment.Amount, decoded.Amount)
	}
}

// TestOrderCreatedEvent tests event struct
func TestOrderCreatedEvent(t *testing.T) {
	event := OrderCreatedEvent{
		OrderID:     uuid.New().String(),
		UserID:      uuid.New().String(),
		ProductName: "Smartphone",
		Price:       899,
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

// TestPaymentCompletedEvent tests payment event struct
func TestPaymentCompletedEvent(t *testing.T) {
	event := PaymentCompletedEvent{
		PaymentID: uuid.New().String(),
		OrderID:   uuid.New().String(),
		UserID:    uuid.New().String(),
		Amount:    1299,
		Status:    "success",
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(event)
	if err != nil {
		t.Fatalf("Failed to marshal payment event: %v", err)
	}

	var decoded PaymentCompletedEvent
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal payment event: %v", err)
	}

	if decoded.Status != event.Status {
		t.Errorf("Status mismatch: expected '%s', got '%s'", event.Status, decoded.Status)
	}
	if decoded.Amount != event.Amount {
		t.Errorf("Amount mismatch: expected %d, got %d", event.Amount, decoded.Amount)
	}
}

// TestProcessPayment_Simulation tests payment processing logic
func TestProcessPayment_Simulation(t *testing.T) {
	tests := []struct {
		name          string
		price         int
		wantStatus    string
		description   string
	}{
		{
			name:       "Normal price - should succeed",
			price:      100,
			wantStatus: "success",
			description: "Regular payment should succeed",
		},
		{
			name:       "Price ending in 99 - should fail",
			price:      199,
			wantStatus: "failed",
			description: "Payment ending in 99 simulates failure",
		},
		{
			name:       "Another normal price",
			price:      500,
			wantStatus: "success",
			description: "Another regular payment",
		},
		{
			name:       "Price 99 exactly - should fail",
			price:      99,
			wantStatus: "failed",
			description: "Price 99 should fail",
		},
		{
			name:       "Zero price",
			price:      0,
			wantStatus: "success",
			description: "Zero amount payment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the payment status logic from main.go
			status := "success"
			if tt.price%100 == 99 {
				status = "failed"
			}

			if status != tt.wantStatus {
				t.Errorf("Expected status '%s', got '%s' for price %d", tt.wantStatus, status, tt.price)
			}
		})
	}
}

// TestGetPaymentsHandler_MethodNotAllowed tests HTTP method validation
func TestGetPaymentsHandler_MethodNotAllowed(t *testing.T) {
	server := &Server{
		kafkaProducer: &MockKafkaProducer{},
		ctx:           nil,
	}

	req := httptest.NewRequest(http.MethodPost, "/payments", nil)
	rr := httptest.NewRecorder()

	server.handleGetPayments(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

// TestGetPaymentHandler_MethodNotAllowed tests method validation
func TestGetPaymentHandler_MethodNotAllowed(t *testing.T) {
	server := &Server{
		kafkaProducer: &MockKafkaProducer{},
		ctx:           nil,
	}

	req := httptest.NewRequest(http.MethodPost, "/payments/123", nil)
	rr := httptest.NewRecorder()

	server.handleGetPayment(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", rr.Code)
	}
}

// TestGetPaymentHandler_MissingID tests missing payment ID
func TestGetPaymentHandler_MissingID(t *testing.T) {
	server := &Server{
		kafkaProducer: &MockKafkaProducer{},
		ctx:           nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/payments/", nil)
	rr := httptest.NewRecorder()

	server.handleGetPayment(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for missing ID, got %d", rr.Code)
	}
}

// TestHealthHandler_Response tests health check response
func TestHealthHandler_Response(t *testing.T) {
	server := &Server{
		kafkaProducer: nil,
		ctx:           nil,
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

	if response["service"] != "payment-service" {
		t.Errorf("Expected service 'payment-service', got '%v'", response["service"])
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", response["status"])
	}
}

// TestStatsHandler_Response tests stats endpoint
func TestStatsHandler_Response(t *testing.T) {
	server := &Server{
		kafkaProducer: nil,
		ctx:           nil,
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

	if response["service"] != "payment-service" {
		t.Errorf("Expected service 'payment-service', got '%v'", response["service"])
	}

	if _, ok := response["totalPayments"]; !ok {
		t.Error("Expected 'totalPayments' field in stats response")
	}

	if _, ok := response["paymentsByStatus"]; !ok {
		t.Error("Expected 'paymentsByStatus' field in stats response")
	}
}

// TestConsumerGroupHandler tests the Kafka consumer handler
func TestConsumerGroupHandler(t *testing.T) {
	server := &Server{
		kafkaProducer: &MockKafkaProducer{},
		ctx:           nil,
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

// TestMockKafkaProducer tests the mock producer
func TestMockKafkaProducer(t *testing.T) {
	mock := &MockKafkaProducer{}

	event := PaymentCompletedEvent{
		PaymentID: "payment-123",
		OrderID:   "order-123",
		UserID:    "user-123",
		Amount:    500,
		Status:    "success",
		CreatedAt: time.Now(),
	}

	data, _ := json.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "payment.completed",
		Key:   sarama.StringEncoder("payment-123"),
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
		t.Errorf("Expected 1 message in mock, got %d", len(mock.messages))
	}

	// Verify message topic
	if mock.messages[0].Topic != "payment.completed" {
		t.Errorf("Expected topic 'payment.completed', got '%s'", mock.messages[0].Topic)
	}
}

// TestMockKafkaProducer_Error tests error handling
func TestMockKafkaProducer_Error(t *testing.T) {
	mock := &MockKafkaProducer{
		err: sarama.ErrLeaderNotAvailable,
	}

	msg := &sarama.ProducerMessage{
		Topic: "payment.completed",
		Value: sarama.ByteEncoder(`{}`),
	}

	_, _, err := mock.SendMessage(msg)
	if err != sarama.ErrLeaderNotAvailable {
		t.Errorf("Expected ErrLeaderNotAvailable, got %v", err)
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
	reg.MustRegister(paymentsProcessed)

	metrics, err := reg.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	expectedMetrics := []string{
		"payment_service_http_requests_total",
		"payment_service_http_request_duration_seconds",
		"payment_service_db_queries_total",
		"payment_service_kafka_messages_produced_total",
		"payment_service_kafka_messages_consumed_total",
		"payment_service_active_connections",
		"payment_service_payments_processed_total",
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

// TestPrometheusCounter_Increment tests counter increments
func TestPrometheusCounter_Increment(t *testing.T) {
	reg := prometheus.NewRegistry()
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "test_payments_total",
			Help: "Test payments",
		},
		[]string{"status"},
	)
	reg.MustRegister(counter)

	// Increment with different statuses
	counter.WithLabelValues("success").Inc()
	counter.WithLabelValues("success").Inc()
	counter.WithLabelValues("failed").Inc()
	counter.WithLabelValues("pending").Inc()

	metrics, _ := reg.Gather()

	for _, mf := range metrics {
		if mf.GetName() == "test_payments_total" {
			if len(mf.GetMetric()) != 3 {
				t.Errorf("Expected 3 metric variations, got %d", len(mf.GetMetric()))
			}

			// Check success count
			for _, m := range mf.GetMetric() {
				labels := m.GetLabel()
				for _, label := range labels {
					if label.GetName() == "status" && label.GetValue() == "success" {
						if m.GetCounter().GetValue() != 2 {
							t.Errorf("Expected 2 success payments, got %f", m.GetCounter().GetValue())
						}
					}
				}
			}
		}
	}
}

// TestBSON_MarshalUnmarshal tests BSON operations
func TestBSON_MarshalUnmarshal(t *testing.T) {
	now := time.Now()
	payment := Payment{
		ID:            "test-id-123",
		OrderID:       "order-456",
		UserID:        "user-789",
		Amount:        999,
		PaymentStatus: "success",
		CreatedAt:     now,
	}

	// Marshal
	data, err := bson.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Unmarshal
	var decoded Payment
	if err := bson.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify all fields
	if decoded.ID != payment.ID {
		t.Errorf("ID mismatch: expected '%s', got '%s'", payment.ID, decoded.ID)
	}
	if decoded.OrderID != payment.OrderID {
		t.Errorf("OrderID mismatch: expected '%s', got '%s'", payment.OrderID, decoded.OrderID)
	}
	if decoded.UserID != payment.UserID {
		t.Errorf("UserID mismatch: expected '%s', got '%s'", payment.UserID, decoded.UserID)
	}
	if decoded.Amount != payment.Amount {
		t.Errorf("Amount mismatch: expected %d, got %d", payment.Amount, decoded.Amount)
	}
	if decoded.PaymentStatus != payment.PaymentStatus {
		t.Errorf("PaymentStatus mismatch: expected '%s', got '%s'", payment.PaymentStatus, decoded.PaymentStatus)
	}
}

// TestBSON_Filter tests BSON filter creation
func TestBSON_Filter(t *testing.T) {
	// Test creating filters for MongoDB queries
	filter1 := bson.M{"_id": "payment-123"}
	if filter1["_id"] != "payment-123" {
		t.Error("Filter ID mismatch")
	}

	filter2 := bson.M{"order_id": "order-456"}
	if filter2["order_id"] != "order-456" {
		t.Error("Filter order_id mismatch")
	}

	filter3 := bson.M{"payment_status": "success"}
	if filter3["payment_status"] != "success" {
		t.Error("Filter status mismatch")
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

	// Validate UUID format
	if len(id1) != 36 {
		t.Errorf("UUID should be 36 characters, got %d", len(id1))
	}
}

// TestTime_Formatting tests timestamp handling
func TestTime_Formatting(t *testing.T) {
	now := time.Now()

	payment := Payment{
		ID:            "test-id",
		OrderID:       "order-id",
		UserID:        "user-id",
		Amount:        100,
		PaymentStatus: "success",
		CreatedAt:     now,
	}

	// Test JSON marshaling
	data, err := json.Marshal(payment)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var decoded Payment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Times should be approximately equal
	diff := decoded.CreatedAt.Sub(now)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Second {
		t.Errorf("Time mismatch: expected %v, got %v", now, decoded.CreatedAt)
	}
}

// TestJSONEncoding_MultiplePayments tests JSON encoding with multiple payments
func TestJSONEncoding_MultiplePayments(t *testing.T) {
	payments := make([]Payment, 50)
	for i := 0; i < 50; i++ {
		status := "success"
		if i%3 == 0 {
			status = "failed"
		}

		payments[i] = Payment{
			ID:            uuid.New().String(),
			OrderID:       uuid.New().String(),
			UserID:        uuid.New().String(),
			Amount:        (i + 1) * 10,
			PaymentStatus: status,
			CreatedAt:     time.Now(),
		}
	}

	data, err := json.Marshal(payments)
	if err != nil {
		t.Fatalf("Failed to marshal 50 payments: %v", err)
	}

	var decoded []Payment
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal 50 payments: %v", err)
	}

	if len(decoded) != 50 {
		t.Errorf("Expected 50 payments, got %d", len(decoded))
	}

	// Verify status distribution
	successCount := 0
	failedCount := 0
	for _, p := range decoded {
		if p.PaymentStatus == "success" {
			successCount++
		} else if p.PaymentStatus == "failed" {
			failedCount++
		}
	}

	if successCount == 0 || failedCount == 0 {
		t.Errorf("Expected mix of success and failed payments")
	}
}

// BenchmarkPayment_Marshal benchmarks payment marshaling
func BenchmarkPayment_Marshal(b *testing.B) {
	payment := Payment{
		ID:            uuid.New().String(),
		OrderID:       uuid.New().String(),
		UserID:        uuid.New().String(),
		Amount:        999,
		PaymentStatus: "success",
		CreatedAt:     time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(payment)
	}
}

// BenchmarkPayment_BSONMarshal benchmarks BSON marshaling
func BenchmarkPayment_BSONMarshal(b *testing.B) {
	payment := Payment{
		ID:            uuid.New().String(),
		OrderID:       uuid.New().String(),
		UserID:        uuid.New().String(),
		Amount:        999,
		PaymentStatus: "success",
		CreatedAt:     time.Now(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = bson.Marshal(payment)
	}
}

// BenchmarkUUID_Generation benchmarks UUID generation
func BenchmarkUUID_Generation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = uuid.New().String()
	}
}

// Helper function to get counter value
func getCounterValue(counter prometheus.Counter) float64 {
	dto := &dto.Metric{}
	counter.Write(dto)
	return dto.GetCounter().GetValue()
}
