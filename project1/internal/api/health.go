package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type HealthChecker interface {
	GetHealthStatus() map[string]interface{}
}

type HealthResponse struct {
	Status    string                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Services  map[string]interface{} `json:"services"`
}

type HealthHandler struct {
	logger   *zap.Logger
	checkers map[string]HealthChecker
}

func NewHealthHandler(logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		logger:   logger,
		checkers: make(map[string]HealthChecker),
	}
}

func (h *HealthHandler) RegisterHealthCheck(name string, checker HealthChecker) {
	h.checkers[name] = checker
	h.logger.Info("Registered health check", zap.String("service", name))
}

func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := &HealthResponse{
		Timestamp: time.Now(),
		Services:  make(map[string]interface{}),
	}

	overallStatus := "healthy"

	for name, checker := range h.checkers {
		status := checker.GetHealthStatus()
		response.Services[name] = status
		
		if statusStr, ok := status["status"].(string); ok && statusStr != "healthy" {
			overallStatus = "unhealthy"
		}
	}

	response.Status = overallStatus

	w.Header().Set("Content-Type", "application/json")
	if overallStatus == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	h.logger.Debug("Health check completed", 
		zap.String("overall_status", overallStatus),
		zap.Int("services_checked", len(h.checkers)))

	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) HandleReady(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "ready",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *HealthHandler) HandleLive(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
