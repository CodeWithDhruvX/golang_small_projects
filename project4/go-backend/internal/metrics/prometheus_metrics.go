package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

// PrometheusMetrics handles application metrics
type PrometheusMetrics struct {
	// API metrics
	apiRequestsTotal *prometheus.CounterVec
	apiRequestDuration *prometheus.HistogramVec
	
	// AI metrics
	aiRequestsTotal *prometheus.CounterVec
	aiResponseDuration *prometheus.HistogramVec
	aiTokenUsage *prometheus.CounterVec
	
	// Email metrics
	emailsProcessed *prometheus.CounterVec
	emailProcessingDuration *prometheus.HistogramVec
	recruiterEmailsDetected *prometheus.CounterVec
	
	// Vector search metrics
	vectorSearchDuration *prometheus.HistogramVec
	vectorSearchResults *prometheus.HistogramVec
	
	// Application metrics
	applicationsCreated *prometheus.CounterVec
	duplicateApplicationsDetected *prometheus.CounterVec
	
	// System metrics
	redisConnections *prometheus.GaugeVec
	databaseConnections *prometheus.GaugeVec
	cacheHitRate *prometheus.GaugeVec
}

// NewPrometheusMetrics creates a new Prometheus metrics instance
func NewPrometheusMetrics() *PrometheusMetrics {
	metrics := &PrometheusMetrics{
		// API metrics
		apiRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_recruiter_api_requests_total",
				Help: "Total number of API requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		apiRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ai_recruiter_api_request_duration_seconds",
				Help:    "Duration of API requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "endpoint"},
		),
		
		// AI metrics
		aiRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_recruiter_ai_requests_total",
				Help: "Total number of AI requests",
			},
			[]string{"model", "operation"},
		),
		aiResponseDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ai_recruiter_ai_response_duration_seconds",
				Help:    "Duration of AI responses in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"model", "operation"},
		),
		aiTokenUsage: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_recruiter_ai_tokens_total",
				Help: "Total number of AI tokens used",
			},
			[]string{"model", "operation"},
		),
		
		// Email metrics
		emailsProcessed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_recruiter_emails_processed_total",
				Help: "Total number of emails processed",
			},
			[]string{"source", "status"},
		),
		emailProcessingDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ai_recruiter_email_processing_duration_seconds",
				Help:    "Duration of email processing in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"source"},
		),
		recruiterEmailsDetected: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_recruiter_recruiter_emails_detected_total",
				Help: "Total number of recruiter emails detected",
			},
			[]string{"confidence_range"},
		),
		
		// Vector search metrics
		vectorSearchDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ai_recruiter_vector_search_duration_seconds",
				Help:    "Duration of vector searches in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"operation"},
		),
		vectorSearchResults: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ai_recruiter_vector_search_results_count",
				Help:    "Number of results from vector searches",
				Buckets: []float64{0, 1, 2, 5, 10, 20, 50, 100},
			},
			[]string{"operation"},
		),
		
		// Application metrics
		applicationsCreated: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_recruiter_applications_created_total",
				Help: "Total number of applications created",
			},
			[]string{"company"},
		),
		duplicateApplicationsDetected: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ai_recruiter_duplicate_applications_detected_total",
				Help: "Total number of duplicate applications detected",
			},
			[]string{"company"},
		),
		
		// System metrics
		redisConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "ai_recruiter_redis_connections",
				Help: "Number of Redis connections",
			},
			[]string{"state"},
		),
		databaseConnections: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "ai_recruiter_database_connections",
				Help: "Number of database connections",
			},
			[]string{"state"},
		),
		cacheHitRate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "ai_recruiter_cache_hit_rate",
				Help: "Cache hit rate as percentage",
			},
			[]string{"cache_type"},
		),
	}
	
	// Register all metrics with Prometheus
	prometheus.MustRegister(
		metrics.apiRequestsTotal,
		metrics.apiRequestDuration,
		metrics.aiRequestsTotal,
		metrics.aiResponseDuration,
		metrics.aiTokenUsage,
		metrics.emailsProcessed,
		metrics.emailProcessingDuration,
		metrics.recruiterEmailsDetected,
		metrics.vectorSearchDuration,
		metrics.vectorSearchResults,
		metrics.applicationsCreated,
		metrics.duplicateApplicationsDetected,
		metrics.redisConnections,
		metrics.databaseConnections,
		metrics.cacheHitRate,
	)
	
	return metrics
}

// RecordAPIRequest records an API request
func (pm *PrometheusMetrics) RecordAPIRequest(method, endpoint, status string, duration time.Duration) {
	pm.apiRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	pm.apiRequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// RecordAIRequest records an AI request
func (pm *PrometheusMetrics) RecordAIRequest(model, operation string, duration time.Duration, tokens int) {
	pm.aiRequestsTotal.WithLabelValues(model, operation).Inc()
	pm.aiResponseDuration.WithLabelValues(model, operation).Observe(duration.Seconds())
	pm.aiTokenUsage.WithLabelValues(model, operation).Add(float64(tokens))
}

// RecordEmailProcessing records email processing
func (pm *PrometheusMetrics) RecordEmailProcessing(source, status string, duration time.Duration) {
	pm.emailsProcessed.WithLabelValues(source, status).Inc()
	pm.emailProcessingDuration.WithLabelValues(source).Observe(duration.Seconds())
}

// RecordRecruiterEmailDetection records recruiter email detection
func (pm *PrometheusMetrics) RecordRecruiterEmailDetection(confidence float64) {
	confidenceRange := getConfidenceRange(confidence)
	pm.recruiterEmailsDetected.WithLabelValues(confidenceRange).Inc()
}

// RecordVectorSearch records vector search performance
func (pm *PrometheusMetrics) RecordVectorSearch(operation string, duration time.Duration, resultCount int) {
	pm.vectorSearchDuration.WithLabelValues(operation).Observe(duration.Seconds())
	pm.vectorSearchResults.WithLabelValues(operation).Observe(float64(resultCount))
}

// RecordApplicationCreated records application creation
func (pm *PrometheusMetrics) RecordApplicationCreated(company string) {
	pm.applicationsCreated.WithLabelValues(company).Inc()
}

// RecordDuplicateApplicationDetected records duplicate application detection
func (pm *PrometheusMetrics) RecordDuplicateApplicationDetected(company string) {
	pm.duplicateApplicationsDetected.WithLabelValues(company).Inc()
}

// UpdateRedisConnections updates Redis connection metrics
func (pm *PrometheusMetrics) UpdateRedisConnections(active, idle int) {
	pm.redisConnections.WithLabelValues("active").Set(float64(active))
	pm.redisConnections.WithLabelValues("idle").Set(float64(idle))
}

// UpdateDatabaseConnections updates database connection metrics
func (pm *PrometheusMetrics) UpdateDatabaseConnections(active, idle int) {
	pm.databaseConnections.WithLabelValues("active").Set(float64(active))
	pm.databaseConnections.WithLabelValues("idle").Set(float64(idle))
}

// UpdateCacheHitRate updates cache hit rate metrics
func (pm *PrometheusMetrics) UpdateCacheHitRate(cacheType string, hitRate float64) {
	pm.cacheHitRate.WithLabelValues(cacheType).Set(hitRate)
}

// getConfidenceRange categorizes confidence into ranges
func getConfidenceRange(confidence float64) string {
	if confidence >= 0.9 {
		return "high"
	} else if confidence >= 0.7 {
		return "medium"
	} else if confidence >= 0.5 {
		return "low"
	}
	return "very_low"
}

// MetricsHandler returns the Prometheus metrics HTTP handler
func (pm *PrometheusMetrics) MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// StartMetricsServer starts the metrics HTTP server
func (pm *PrometheusMetrics) StartMetricsServer(addr string) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", pm.MetricsHandler())
	
	logrus.Infof("Starting Prometheus metrics server on %s", addr)
	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			logrus.Errorf("Failed to start metrics server: %v", err)
		}
	}()
}
