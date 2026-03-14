package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ai-recruiter-assistant/internal/auth"
	"ai-recruiter-assistant/internal/storage"
	"github.com/google/uuid"
	"time"
)

// EmailHandler handles email-related HTTP requests
type EmailHandler struct {
	storage     storage.StorageInterface
	authService *auth.AuthService
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(storage storage.StorageInterface, authService *auth.AuthService) *EmailHandler {
	return &EmailHandler{
		storage:     storage,
		authService: authService,
	}
}

// Email represents an email object
type Email struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	SenderEmail  string `json:"sender_email"`
	SenderName   string `json:"sender_name,omitempty"`
	IsRecruiter  bool   `json:"is_recruiter"`
	Processed    bool   `json:"processed"`
	GmailID      string `json:"gmail_id,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// GetEmails retrieves all emails for the authenticated user
// @Summary Get user emails
// @Description Retrieve all emails for the authenticated user with optional date range filtering
// @Tags emails
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Param start_date query string false "Start date (YYYY-MM-DD format)"
// @Param end_date query string false "End date (YYYY-MM-DD format)"
// @Success 200 {array} Email
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /emails [get]
func (h *EmailHandler) GetEmails(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "25")) // Increased from 10 to 25
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 25 // Increased from 10 to 25
	}

	// Parse date range parameters with default to one week
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	
	var startDate, endDate time.Time
	var err error
	
	// Default to one week ago if no start date provided
	if startDateStr == "" {
		startDate = time.Now().AddDate(0, 0, -7)
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
	}
	
	// Default to today if no end date provided
	if endDateStr == "" {
		endDate = time.Now()
	} else {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
	}

	// Validate date range
	if startDate.After(endDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date cannot be after end date"})
		return
	}

	// Get emails from storage with date range filtering
	storageEmails, err := h.storage.GetUserEmailsByDateRange(userID, page, limit, startDate, endDate)
	if err != nil {
		logrus.Errorf("Failed to get emails: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve emails"})
		return
	}

	// Convert storage emails to handler emails
	emails := make([]Email, len(storageEmails))
	for i, storageEmail := range storageEmails {
		emails[i] = Email{
			ID:           storageEmail.ID,
			UserID:       storageEmail.UserID,
			Subject:      storageEmail.Subject,
			Body:         storageEmail.Body,
			SenderEmail:  storageEmail.SenderEmail,
			SenderName:   storageEmail.SenderName,
			IsRecruiter:  storageEmail.IsRecruiter,
			Processed:    storageEmail.Processed,
			GmailID:      storageEmail.GmailID,
			CreatedAt:    storageEmail.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    storageEmail.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"emails":     emails,
		"page":       page,
		"limit":      limit,
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
	})
}

// GetEmail retrieves a specific email by ID
// @Summary Get email by ID
// @Description Retrieve a specific email by its ID
// @Tags emails
// @Produce json
// @Security BearerAuth
// @Param id path string true "Email ID"
// @Success 200 {object} Email
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /emails/{id} [get]
func (h *EmailHandler) GetEmail(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	emailID := c.Param("id")
	if emailID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email ID is required"})
		return
	}

	storageEmail, err := h.storage.GetEmailByID(emailID, userID)
	if err != nil {
		if err.Error() == "email not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
			return
		}
		logrus.Errorf("Failed to get email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email"})
		return
	}

	// Convert storage email to handler email
	email := Email{
		ID:           storageEmail.ID,
		UserID:       storageEmail.UserID,
		Subject:      storageEmail.Subject,
		Body:         storageEmail.Body,
		SenderEmail:  storageEmail.SenderEmail,
		SenderName:   storageEmail.SenderName,
		IsRecruiter:  storageEmail.IsRecruiter,
		Processed:    storageEmail.Processed,
		GmailID:      storageEmail.GmailID,
		CreatedAt:    storageEmail.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    storageEmail.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusOK, email)
}

// ImportEmails imports emails from IMAP server
// @Summary Import emails from IMAP
// @Description Import emails from configured IMAP server
// @Tags emails
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /emails/import [post]
func (h *EmailHandler) ImportEmails(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Mock email data for demonstration
	mockEmails := []storage.Email{
		{
			ID:          uuid.New().String(),
			UserID:      userID,
			Subject:     "Senior Software Engineer Role - Google",
			Body:        "Hi,\n\nI'm a recruiter at Google. We're looking for a Senior Software Engineer with Go and Python experience. Your profile looks great! Are you interested?\n\nBest,\nGoogle Recruiting",
			SenderEmail: "hr@google.com",
			SenderName:  "Google Recruitment",
			IsRecruiter: true,
			Processed:   false,
			GmailID:     "mock_google_001",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New().String(),
			UserID:      userID,
			Subject:     "Opportunity at Meta",
			Body:        "Hello Dhruv,\n\nWe saw your GitHub and were impressed by your Go projects. Meta is hiring for infrastructure team. Let's chat!\n\nCheers,\nMeta HR",
			SenderEmail: "recruiter@meta.com",
			SenderName:  "Meta Recruitment",
			IsRecruiter: true,
			Processed:   false,
			GmailID:     "mock_meta_002",
			CreatedAt:   time.Now().Add(-1 * time.Hour),
			UpdatedAt:   time.Now().Add(-1 * time.Hour),
		},
	}

	for _, email := range mockEmails {
		_, err := h.storage.CreateEmail(&email)
		if err != nil {
			logrus.Errorf("Failed to create mock email: %v", err)
		}
	}

	logrus.Infof("Mock emails imported for user: %s", userID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Sample emails imported successfully for demonstration",
		"status":  "completed",
		"count":   len(mockEmails),
	})
}

// UploadEmail handles manual email upload
// @Summary Upload email manually
// @Description Upload an email manually (JSON format)
// @Tags emails
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param email body object{subject=string,body=string,sender_email=string,sender_name=string} true "Email data"
// @Success 201 {object} Email
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /emails/upload [post]
func (h *EmailHandler) UploadEmail(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var emailData struct {
		Subject     string `json:"subject" binding:"required"`
		Body        string `json:"body" binding:"required"`
		SenderEmail string `json:"sender_email" binding:"required,email"`
		SenderName  string `json:"sender_name,omitempty"`
	}

	if err := c.ShouldBindJSON(&emailData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create email
	email := &storage.Email{
		UserID:      userID,
		Subject:     emailData.Subject,
		Body:        emailData.Body,
		SenderEmail: emailData.SenderEmail,
		SenderName:  emailData.SenderName,
		IsRecruiter: false, // Will be determined by AI classification
		Processed:   false,
	}

	createdEmail, err := h.storage.CreateEmail(email)
	if err != nil {
		logrus.Errorf("Failed to create email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create email"})
		return
	}

	// Convert storage email to handler email
	responseEmail := Email{
		ID:           createdEmail.ID,
		UserID:       createdEmail.UserID,
		Subject:      createdEmail.Subject,
		Body:         createdEmail.Body,
		SenderEmail:  createdEmail.SenderEmail,
		SenderName:   createdEmail.SenderName,
		IsRecruiter:  createdEmail.IsRecruiter,
		Processed:    createdEmail.Processed,
		GmailID:      createdEmail.GmailID,
		CreatedAt:    createdEmail.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    createdEmail.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	c.JSON(http.StatusCreated, responseEmail)
}

// GetEmailStats retrieves email statistics for the dashboard
// @Summary Get email statistics
// @Description Retrieve email statistics for the dashboard
// @Tags emails
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /emails/stats [get]
func (h *EmailHandler) GetEmailStats(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get all emails to calculate statistics
	emails, err := h.storage.GetUserEmails(userID, 1, 1000) // Get first 1000 emails
	if err != nil {
		logrus.Errorf("Failed to get emails for stats: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve email statistics"})
		return
	}

	totalEmails := len(emails)
	recruiterEmails := 0
	processedEmails := 0

	for _, email := range emails {
		if email.IsRecruiter {
			recruiterEmails++
		}
		if email.Processed {
			processedEmails++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_emails":     totalEmails,
		"recruiter_emails": recruiterEmails,
		"processed_emails": processedEmails,
		"unprocessed_emails": totalEmails - processedEmails,
	})
}
