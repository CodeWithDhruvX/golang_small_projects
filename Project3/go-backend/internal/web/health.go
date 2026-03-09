package web

import (
	"net/http"
	"time"

	"private-knowledge-base-go/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db     *storage.PostgresDB
	logger *logrus.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *storage.PostgresDB, logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// Health performs a basic health check
// @Summary Health check
// @Description Basic health check endpoint
// @Tags health
// @Produce json
// @Success 200 {object} storage.HealthResponse
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, storage.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC(),
		Version:   "1.0.0",
	})
}

// Ready performs a readiness check
// @Summary Readiness check
// @Description Checks if the application is ready to serve requests
// @Tags health
// @Produce json
// @Success 200 {object} storage.ReadyResponse
// @Failure 503 {object} ErrorResponse
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	services := make(map[string]string)
	
	// Check database connection
	if err := h.db.Ping(); err != nil {
		services["database"] = "unhealthy"
		c.JSON(http.StatusServiceUnavailable, storage.ReadyResponse{
			Status:    "not ready",
			Timestamp: time.Now().UTC(),
			Services:  services,
		})
		return
	}
	services["database"] = "healthy"
	
	// TODO: Check Ollama connection
	services["ollama"] = "unknown" // Would check Ollama health
	
	// TODO: Check Redis connection
	services["redis"] = "unknown" // Would check Redis health

	c.JSON(http.StatusOK, storage.ReadyResponse{
		Status:    "ready",
		Timestamp: time.Now().UTC(),
		Services:  services,
	})
}

// Metrics returns application metrics
// @Summary Application metrics
// @Description Returns various application metrics
// @Tags metrics
// @Produce json
// @Success 200 {object} storage.MetricsResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/metrics [get]
func (h *HealthHandler) Metrics(c *gin.Context) {
	metrics, err := h.db.GetMetrics(c.Request.Context())
	if err != nil {
		h.logger.Errorf("Failed to get metrics: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve metrics",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// Live performs a liveness check
// @Summary Liveness check
// @Description Checks if the application is alive
// @Tags health
// @Produce json
// @Success 200 {object} LivenessResponse
// @Router /live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, LivenessResponse{
		Status:    "alive",
		Timestamp: time.Now().UTC(),
		Uptime:    "0s", // TODO: Calculate actual uptime
	})
}

// LivenessResponse represents liveness check response
type LivenessResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime"`
}
