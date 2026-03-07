package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dhruv/kafka-microservices/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EventPublisher interface {
	PublishEvent(topic string, event interface{}) error
}

type OrderHandler struct {
	logger   *zap.Logger
	producer EventPublisher
	topic    string
}

func NewOrderHandler(logger *zap.Logger, producer EventPublisher, topic string) *OrderHandler {
	return &OrderHandler{
		logger:   logger,
		producer: producer,
		topic:    topic,
	}
}

func (h *OrderHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Enrich order
	order.ID = uuid.New().String()
	order.Status = "PENDING"
	order.CreatedAt = time.Now().UTC().Format(time.RFC3339)

	// Publish to Kafka
	if err := h.producer.PublishEvent(h.topic, order); err != nil {
		h.logger.Error("Failed to publish order event", zap.Error(err))
		http.Error(w, "Failed to process order", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Order received and published", zap.String("order_id", order.ID))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
