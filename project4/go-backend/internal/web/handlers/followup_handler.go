package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ai-recruiter-assistant/internal/ai"
	"ai-recruiter-assistant/internal/auth"
	"ai-recruiter-assistant/internal/followup"
	"ai-recruiter-assistant/internal/storage"
)

// FollowUpHandler handles follow-up email generation requests
type FollowUpHandler struct {
	storage     storage.StorageInterface
	authService *auth.AuthService
	ollama      *ai.OllamaService
	followUp    *followup.FollowUpService
}

// NewFollowUpHandler creates a new follow-up handler
func NewFollowUpHandler(storage storage.StorageInterface, authService *auth.AuthService) *FollowUpHandler {
	// Initialize Ollama service
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	
	ollama := ai.NewOllamaService(ollamaURL)
	followUpService := followup.NewFollowUpService(storage, ollama)
	
	return &FollowUpHandler{
		storage:     storage,
		authService: authService,
		ollama:      ollama,
		followUp:    followUpService,
	}
}

// GenerateFollowUpRequest represents the request for follow-up generation
type GenerateFollowUpRequest struct {
	ApplicationID string `json:"application_id" binding:"required"`
	Type          string `json:"type" binding:"required"`
}

// GenerateFollowUp generates a follow-up email
// @Summary Generate follow-up email
// @Description Generate a professional follow-up email for a job application
// @Tags followup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GenerateFollowUpRequest true "Follow-up generation request"
// @Success 200 {object} followup.FollowUpEmail
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /followup/generate [post]
func (h *FollowUpHandler) GenerateFollowUp(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req GenerateFollowUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get application to calculate days since last contact
	application, err := h.storage.GetApplicationByID(req.ApplicationID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Calculate days since last update
	daysSinceLast := int(time.Since(application.UpdatedAt).Hours() / 24)

	// Create follow-up request
	followUpReq := followup.FollowUpRequest{
		ApplicationID: req.ApplicationID,
		DaysSinceLast: daysSinceLast,
		Type:          req.Type,
	}

	// Generate follow-up email
	followUpEmail, err := h.followUp.GenerateFollowUpEmail(c.Request.Context(), userID, followUpReq)
	if err != nil {
		logrus.Errorf("Failed to generate follow-up email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate follow-up email"})
		return
	}

	c.JSON(http.StatusOK, followUpEmail)
}

// GetPendingFollowUps returns applications that need follow-up
// @Summary Get pending follow-ups
// @Description Get applications that need follow-up emails
// @Tags followup
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /followup/pending [get]
func (h *FollowUpHandler) GetPendingFollowUps(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	pendingFollowUps, err := h.followUp.GetPendingFollowUps(c.Request.Context(), userID)
	if err != nil {
		logrus.Errorf("Failed to get pending follow-ups: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pending follow-ups"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"pending_follow_ups": pendingFollowUps,
		"count":              len(pendingFollowUps),
	})
}

// ScheduleFollowUp schedules a follow-up reminder
// @Summary Schedule follow-up
// @Description Schedule a follow-up reminder for an application
// @Tags followup
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body map[string]interface{} true "Schedule follow-up request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /followup/schedule [post]
func (h *FollowUpHandler) ScheduleFollowUp(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req struct {
		ApplicationID string `json:"application_id" binding:"required"`
		FollowUpDate  string `json:"follow_up_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse follow-up date
	followUpDate, err := time.Parse("2006-01-02", req.FollowUpDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Verify application belongs to user
	_, err = h.storage.GetApplicationByID(req.ApplicationID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
		return
	}

	// Schedule follow-up
	err = h.followUp.ScheduleFollowUp(c.Request.Context(), userID, req.ApplicationID, followUpDate)
	if err != nil {
		logrus.Errorf("Failed to schedule follow-up: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to schedule follow-up"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Follow-up scheduled successfully",
		"application_id": req.ApplicationID,
		"follow_up_date": req.FollowUpDate,
	})
}
