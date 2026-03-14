package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ai-recruiter-assistant/internal/auth"
	"ai-recruiter-assistant/internal/storage"
)

// ApplicationHandler handles application-related HTTP requests
type ApplicationHandler struct {
	storage     storage.StorageInterface
	authService *auth.AuthService
}

// NewApplicationHandler creates a new application handler
func NewApplicationHandler(storage storage.StorageInterface, authService *auth.AuthService) *ApplicationHandler {
	return &ApplicationHandler{
		storage:     storage,
		authService: authService,
	}
}

// Application represents a job application
type Application struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	Company        string `json:"company"`
	Role           string `json:"role"`
	RecruiterEmail string `json:"recruiter_email"`
	RecruiterName  string `json:"recruiter_name,omitempty"`
	Status         string `json:"status"`
	EmailID        string `json:"email_id,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// CreateApplicationRequest represents the request to create a new application
type CreateApplicationRequest struct {
	Company        string `json:"company" binding:"required"`
	Role           string `json:"role" binding:"required"`
	RecruiterEmail string `json:"recruiter_email" binding:"required,email"`
	RecruiterName  string `json:"recruiter_name,omitempty"`
	EmailID        string `json:"email_id,omitempty"`
}

// GetApplications retrieves all applications for the authenticated user
// @Summary Get user applications
// @Description Retrieve all job applications for the authenticated user
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications [get]
func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Parse status filter
	status := c.Query("status")

	// Get applications from storage
	applications, err := h.storage.GetUserApplications(userID, page, limit, status)
	if err != nil {
		logrus.Errorf("Failed to get applications: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": applications,
		"page":         page,
		"limit":        limit,
		"status":       status,
	})
}

// CreateApplication creates a new job application
// @Summary Create application
// @Description Create a new job application
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param application body CreateApplicationRequest true "Application data"
// @Success 201 {object} Application
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate application
	existing, err := h.storage.CheckDuplicateApplication(userID, req.Company, req.RecruiterEmail)
	if err == nil && existing {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Application already exists for this company and recruiter",
			"warning": "already applied",
		})
		return
	}

	// Create application
	application := &storage.Application{
		UserID:         userID,
		Company:        req.Company,
		Role:           req.Role,
		RecruiterEmail: req.RecruiterEmail,
		RecruiterName:  req.RecruiterName,
		Status:         "Applied",
		EmailID:        req.EmailID,
	}

	createdApplication, err := h.storage.CreateApplication(application)
	if err != nil {
		logrus.Errorf("Failed to create application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, createdApplication)
}

// GetApplication retrieves a specific application by ID
// @Summary Get application by ID
// @Description Retrieve a specific application by its ID
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Success 200 {object} Application
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id} [get]
func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	applicationID := c.Param("id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
		return
	}

	application, err := h.storage.GetApplicationByID(applicationID, userID)
	if err != nil {
		if err.Error() == "application not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		logrus.Errorf("Failed to get application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve application"})
		return
	}

	c.JSON(http.StatusOK, application)
}

// UpdateApplication updates an existing application
// @Summary Update application
// @Description Update an existing job application
// @Tags applications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Param application body Application true "Application data"
// @Success 200 {object} Application
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id} [put]
func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	applicationID := c.Param("id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
		return
	}

	var application storage.Application
	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure user can only update their own application
	application.ID = applicationID
	application.UserID = userID

	updatedApplication, err := h.storage.UpdateApplication(&application)
	if err != nil {
		if err.Error() == "application not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		logrus.Errorf("Failed to update application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update application"})
		return
	}

	c.JSON(http.StatusOK, updatedApplication)
}

// DeleteApplication deletes an application
// @Summary Delete application
// @Description Delete a job application
// @Tags applications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Application ID"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /applications/{id} [delete]
func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	applicationID := c.Param("id")
	if applicationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Application ID is required"})
		return
	}

	err := h.storage.DeleteApplication(applicationID, userID)
	if err != nil {
		if err.Error() == "application not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		logrus.Errorf("Failed to delete application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application deleted successfully"})
}
