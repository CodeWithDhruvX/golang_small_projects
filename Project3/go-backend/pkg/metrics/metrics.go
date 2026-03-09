package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Collector holds all application metrics
type Collector struct {
	// HTTP metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpResponseSize     *prometheus.HistogramVec

	// Database metrics
	dbConnectionsActive  prometheus.Gauge
	dbQueryDuration      *prometheus.HistogramVec
	dbQueryTotal         *prometheus.CounterVec

	// AI/ML metrics
	aiRequestsTotal     *prometheus.CounterVec
	aiRequestDuration   *prometheus.HistogramVec
	aiTokensTotal       *prometheus.CounterVec

	// Document metrics
	documentsTotal      prometheus.Gauge
	documentsProcessed  prometheus.Gauge
	documentChunksTotal prometheus.Gauge

	// Chat metrics
	chatSessionsTotal   prometheus.Gauge
	chatMessagesTotal   prometheus.Gauge
	chatResponseTime    prometheus.Histogram
}

// NewCollector creates a new metrics collector
func NewCollector() *Collector {
	return &Collector{
		// HTTP metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		httpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "HTTP response size in bytes",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			[]string{"method", "endpoint"},
		),

		// Database metrics
		dbConnectionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Number of active database connections",
			},
		),
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"query_type"},
		),
		dbQueryTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"query_type", "status"},
		),

		// AI/ML metrics
		aiRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_requests_total",
				Help: "Total number of AI requests",
			},
			[]string{"model", "request_type", "status"},
		),
		aiRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ai_request_duration_seconds",
				Help:    "AI request duration in seconds",
				Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
			},
			[]string{"model", "request_type"},
		),
		aiTokensTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_tokens_total",
				Help: "Total number of AI tokens processed",
			},
			[]string{"model", "token_type"},
		),

		// Document metrics
		documentsTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "documents_total",
				Help: "Total number of documents",
			},
		),
		documentsProcessed: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "documents_processed_total",
				Help: "Total number of processed documents",
			},
		),
		documentChunksTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "document_chunks_total",
				Help: "Total number of document chunks",
			},
		),

		// Chat metrics
		chatSessionsTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "chat_sessions_total",
				Help: "Total number of chat sessions",
			},
		),
		chatMessagesTotal: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "chat_messages_total",
				Help: "Total number of chat messages",
			},
		),
		chatResponseTime: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "chat_response_time_seconds",
				Help:    "Chat response time in seconds",
				Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
			},
		),
	}
}

// HTTP metrics methods

func (c *Collector) RecordHTTPRequest(method, endpoint, status string, duration time.Duration, size int) {
	c.httpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	c.httpRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	c.httpResponseSize.WithLabelValues(method, endpoint).Observe(float64(size))
}

// Database metrics methods

func (c *Collector) SetDBConnections(count float64) {
	c.dbConnectionsActive.Set(count)
}

func (c *Collector) RecordDBQuery(queryType, status string, duration time.Duration) {
	c.dbQueryDuration.WithLabelValues(queryType).Observe(duration.Seconds())
	c.dbQueryTotal.WithLabelValues(queryType, status).Inc()
}

// AI metrics methods

func (c *Collector) RecordAIRequest(model, requestType, status string, duration time.Duration) {
	c.aiRequestsTotal.WithLabelValues(model, requestType, status).Inc()
	c.aiRequestDuration.WithLabelValues(model, requestType).Observe(duration.Seconds())
}

func (c *Collector) RecordAITokens(model, tokenType string, count float64) {
	c.aiTokensTotal.WithLabelValues(model, tokenType).Add(count)
}

// Document metrics methods

func (c *Collector) SetDocumentsTotal(count float64) {
	c.documentsTotal.Set(count)
}

func (c *Collector) SetDocumentsProcessed(count float64) {
	c.documentsProcessed.Set(count)
}

func (c *Collector) SetDocumentChunksTotal(count float64) {
	c.documentChunksTotal.Set(count)
}

// Chat metrics methods

func (c *Collector) SetChatSessionsTotal(count float64) {
	c.chatSessionsTotal.Set(count)
}

func (c *Collector) SetChatMessagesTotal(count float64) {
	c.chatMessagesTotal.Set(count)
}

func (c *Collector) RecordChatResponseTime(duration time.Duration) {
	c.chatResponseTime.Observe(duration.Seconds())
}
