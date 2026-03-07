package api

import (
	"encoding/json"
	"net/http"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type Metrics struct {
	RequestsTotal    int64     `json:"requests_total"`
	RequestsSuccess  int64     `json:"requests_success"`
	RequestsError    int64     `json:"requests_error"`
	LastRequestTime  time.Time `json:"last_request_time"`
	AverageLatency   int64     `json:"average_latency_ms"`
	TotalLatency     int64     `json:"total_latency_ms"`
}

type MetricsHandler struct {
	logger  *zap.Logger
	metrics Metrics
}

func NewMetricsHandler(logger *zap.Logger) *MetricsHandler {
	return &MetricsHandler{
		logger: logger,
	}
}

func (m *MetricsHandler) RecordRequest(success bool, latency time.Duration) {
	atomic.AddInt64(&m.metrics.RequestsTotal, 1)
	atomic.AddInt64(&m.metrics.TotalLatency, latency.Milliseconds())
	
	if success {
		atomic.AddInt64(&m.metrics.RequestsSuccess, 1)
	} else {
		atomic.AddInt64(&m.metrics.RequestsError, 1)
	}
	
	m.metrics.LastRequestTime = time.Now()
	
	// Update average latency
	total := atomic.LoadInt64(&m.metrics.RequestsTotal)
	if total > 0 {
		totalLatency := atomic.LoadInt64(&m.metrics.TotalLatency)
		atomic.StoreInt64(&m.metrics.AverageLatency, totalLatency/total)
	}
}

func (m *MetricsHandler) GetMetrics() Metrics {
	return Metrics{
		RequestsTotal:   atomic.LoadInt64(&m.metrics.RequestsTotal),
		RequestsSuccess: atomic.LoadInt64(&m.metrics.RequestsSuccess),
		RequestsError:   atomic.LoadInt64(&m.metrics.RequestsError),
		LastRequestTime: m.metrics.LastRequestTime,
		AverageLatency:  atomic.LoadInt64(&m.metrics.AverageLatency),
	}
}

func (m *MetricsHandler) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := m.GetMetrics()
	
	response := map[string]interface{}{
		"http": metrics,
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	m.logger.Debug("Metrics requested", zap.Int64("total_requests", metrics.RequestsTotal))
	json.NewEncoder(w).Encode(response)
}

type KafkaMetrics struct {
	ProducerMessagesProduced int64     `json:"producer_messages_produced"`
	ProducerErrors           int64     `json:"producer_errors"`
	ProducerLastProducedTime time.Time `json:"producer_last_produced_time"`
	ConsumerMessagesConsumed int64     `json:"consumer_messages_consumed"`
	ConsumerProcessingErrors int64     `json:"consumer_processing_errors"`
	ConsumerLastConsumedTime time.Time `json:"consumer_last_consumed_time"`
}

func (m *MetricsHandler) HandleKafkaMetrics(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"kafka": KafkaMetrics{
			// These will be populated by the actual Kafka components
			ProducerMessagesProduced: 0,
			ProducerErrors:           0,
			ProducerLastProducedTime: time.Time{},
			ConsumerMessagesConsumed: 0,
			ConsumerProcessingErrors: 0,
			ConsumerLastConsumedTime: time.Time{},
		},
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
