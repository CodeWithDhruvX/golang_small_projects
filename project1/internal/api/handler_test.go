package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dhruv/kafka-microservices/internal/models"
	"go.uber.org/zap/zaptest"
)

// MockPublisher is a mock implementation of EventPublisher
type MockPublisher struct {
	PublishFunc func(topic string, event interface{}) error
}

func (m *MockPublisher) PublishEvent(topic string, event interface{}) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(topic, event)
	}
	return nil
}

func TestHandleCreateOrder_Success(t *testing.T) {
	logger := zaptest.NewLogger(t)

	publisher := &MockPublisher{
		PublishFunc: func(topic string, event interface{}) error {
			if topic != "test-topic" {
				t.Errorf("expected topic test-topic, got %s", topic)
			}
			return nil
		},
	}

	handler := NewOrderHandler(logger, publisher, "test-topic")

	body := []byte(`{"item": "Laptop", "quantity": 1, "price": 1200.5}`)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandleCreateOrder(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, res.StatusCode)
	}

	var responseOrder models.Order
	if err := json.NewDecoder(res.Body).Decode(&responseOrder); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if responseOrder.ID == "" {
		t.Error("expected non-empty ID")
	}
	if responseOrder.Status != "PENDING" {
		t.Errorf("expected status PENDING, got %s", responseOrder.Status)
	}
}

func TestHandleCreateOrder_MethodNotAllowed(t *testing.T) {
	logger := zaptest.NewLogger(t)
	handler := NewOrderHandler(logger, &MockPublisher{}, "test-topic")

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	w := httptest.NewRecorder()

	handler.HandleCreateOrder(w, req)

	if w.Result().StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Result().StatusCode)
	}
}

func TestHandleCreateOrder_PublishFail(t *testing.T) {
	logger := zaptest.NewLogger(t)

	publisher := &MockPublisher{
		PublishFunc: func(topic string, event interface{}) error {
			return errors.New("kafka is down")
		},
	}

	handler := NewOrderHandler(logger, publisher, "test-topic")

	body := []byte(`{"item": "Laptop", "quantity": 1, "price": 1200.5}`)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.HandleCreateOrder(w, req)

	if w.Result().StatusCode != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d", http.StatusInternalServerError, w.Result().StatusCode)
	}
}

func TestHandleCreateOrder_BadRequest(t *testing.T) {
	logger := zaptest.NewLogger(t)
	handler := NewOrderHandler(logger, &MockPublisher{}, "test-topic")

	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer([]byte(`{invalid-json`)))
	w := httptest.NewRecorder()

	handler.HandleCreateOrder(w, req)

	if w.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Result().StatusCode)
	}
}
