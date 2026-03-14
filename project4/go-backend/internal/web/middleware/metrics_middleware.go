package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"ai-recruiter-assistant/internal/metrics"
)

// MetricsMiddleware provides Prometheus metrics collection
type MetricsMiddleware struct {
	prometheusMetrics *metrics.PrometheusMetrics
}

// NewMetricsMiddleware creates a new metrics middleware
func NewMetricsMiddleware(prometheusMetrics *metrics.PrometheusMetrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		prometheusMetrics: prometheusMetrics,
	}
}

// Middleware returns the Gin middleware function
func (mm *MetricsMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Get status code
		status := c.Writer.Status()
		statusStr := getStatusCategory(status)
		
		// Record metrics
		mm.prometheusMetrics.RecordAPIRequest(
			c.Request.Method,
			c.FullPath(),
			statusStr,
			duration,
		)
		
		// Log request if it's an error
		if status >= 400 {
			logrus.WithFields(logrus.Fields{
				"method":    c.Request.Method,
				"path":      c.FullPath(),
				"status":    status,
				"duration":  duration,
				"client_ip": c.ClientIP(),
			}).Error("API request failed")
		}
	}
}

// getStatusCategory categorizes HTTP status codes
func getStatusCategory(status int) string {
	switch {
	case status >= 200 && status < 300:
		return "2xx"
	case status >= 300 && status < 400:
		return "3xx"
	case status >= 400 && status < 500:
		return "4xx"
	case status >= 500:
		return "5xx"
	default:
		return "unknown"
	}
}
