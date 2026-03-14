package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	gmailv1 "google.golang.org/api/gmail/v1"

	"ai-recruiter-assistant/internal/auth"
	gmailsvc "ai-recruiter-assistant/internal/gmail"
	"ai-recruiter-assistant/internal/storage"
)

// GmailHandler handles Gmail OAuth2 operations
type GmailHandler struct {
	storage      storage.StorageInterface
	gmailService *gmailsvc.GmailService
	authService  *auth.AuthService
}

// NewGmailHandler creates a new Gmail handler
func NewGmailHandler(storage storage.StorageInterface, gmailService *gmailsvc.GmailService, authService *auth.AuthService) *GmailHandler {
	return &GmailHandler{
		storage:      storage,
		gmailService: gmailService,
		authService:  authService,
	}
}

// AuthRequest initiates Gmail OAuth2 authentication
// @Summary Initiate Gmail OAuth2 authentication
// @Description Starts the Gmail OAuth2 flow by returning an authorization URL
// @Tags gmail
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/gmail [get]
func (h *GmailHandler) AuthRequest(c *gin.Context) {
	// Generate state parameter for security
	state, err := gmailsvc.GenerateState()
	if err != nil {
		logrus.Errorf("Failed to generate state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	// Store state in session/cache for later validation
	// For now, we'll pass it through the OAuth flow and validate in callback
	
	authURL := h.gmailService.GetAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
		"message":  "Visit the authorization URL to connect your Gmail account",
	})
}

// AuthCallback handles the OAuth2 callback from Google
// @Summary Handle Gmail OAuth2 callback
// @Description Processes the OAuth2 callback and creates/updates user account
// @Tags gmail
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/gmail/callback [get]
func (h *GmailHandler) AuthCallback(c *gin.Context) {
	code := c.Query("code")
	_ = c.Query("state") // Store state for validation in production

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code is required"})
		return
	}

	// Exchange code for tokens
	token, err := h.gmailService.ExchangeCodeForToken(c.Request.Context(), code)
	if err != nil {
		logrus.Errorf("Failed to exchange code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange authorization code"})
		return
	}

	// Create a temporary Gmail client to get user profile
	tempClient := h.gmailService.GetConfig().Client(c.Request.Context(), token)
	tempService, err := gmailv1.NewService(c.Request.Context(), option.WithHTTPClient(tempClient))
	if err != nil {
		logrus.Errorf("Failed to create temporary Gmail service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Gmail service"})
		return
	}

	// Get user profile from Google
	profile, err := tempService.Users.GetProfile("me").Do()
	if err != nil {
		logrus.Errorf("Failed to get user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user profile"})
		return
	}

	// Check if user exists by email
	existingUser, err := h.storage.GetUserByEmail(profile.EmailAddress)
	if err != nil && err.Error() != "user not found" {
		logrus.Errorf("Failed to check existing user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process user"})
		return
	}

	var userID string
	if existingUser != nil {
		// User exists, use existing ID
		userID = existingUser.ID
	} else {
		// Create new user
		newUser := &auth.User{
			ID:            uuid.New().String(),
			Email:         profile.EmailAddress,
			Name:          profile.EmailAddress, // Use email as name for now
			PasswordHash: "", // OAuth2 users don't have passwords
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		
		err = h.storage.CreateUser(newUser)
		if err != nil {
			logrus.Errorf("Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user account"})
			return
		}
		userID = newUser.ID
	}

	// Store Gmail token
	err = h.gmailService.StoreToken(c.Request.Context(), userID, profile.EmailAddress, token)
	if err != nil {
		logrus.Errorf("Failed to store token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store Gmail integration"})
		return
	}

	// Generate JWT token for the user
	jwtToken, err := h.authService.GenerateToken(userID, profile.EmailAddress)
	if err != nil {
		logrus.Errorf("Failed to generate JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate authentication token"})
		return
	}

	// Redirect to frontend with the token
	frontendURL := "http://localhost:4200/auth/gmail/callback"
	redirectURL := fmt.Sprintf("%s?token=%s&id=%s&email=%s", frontendURL, jwtToken, userID, profile.EmailAddress)
	
	c.Redirect(http.StatusSeeOther, redirectURL)
}

// SyncEmails syncs emails from Gmail
// @Summary Sync emails from Gmail
// @Description Fetches recent emails from Gmail and stores them in the database
// @Tags gmail
// @Produce json
// @Security BearerAuth
// @Param max_results query int false "Maximum number of emails to fetch" default(50)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gmail/sync [post]
func (h *GmailHandler) SyncEmails(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse max_results parameter
	maxResults := int64(100) // Increased from 50 to get more emails from last week
	if maxResultsStr := c.Query("max_results"); maxResultsStr != "" {
		if parsed, err := strconv.Atoi(maxResultsStr); err == nil && parsed > 0 {
			maxResults = int64(parsed)
		}
	}

	// Fetch emails from Gmail
	messages, err := h.gmailService.FetchEmails(c.Request.Context(), userID, maxResults)
	if err != nil {
		logrus.Errorf("Failed to fetch emails: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch emails from Gmail"})
		return
	}

	// Process and store emails
	processedCount := 0
	for _, msg := range messages {
		subject, from, _, body, err := gmailsvc.ParseMessageContent(msg)
		if err != nil {
			logrus.Errorf("Failed to parse message %s: %v", msg.Id, err)
			continue
		}

		// Create email record
		email := &storage.Email{
			ID:          uuid.New().String(),
			UserID:      userID,
			Subject:     subject,
			Body:        body,
			SenderEmail: from,
			SenderName:  extractNameFromEmail(from),
			IsRecruiter: false, // Will be determined by AI classification
			Processed:   false,
			GmailID:     msg.Id,
			CreatedAt:   time.Unix(msg.InternalDate/1000, 0),
			UpdatedAt:   time.Now(),
		}

		_, err = h.storage.CreateEmail(email)
		if err != nil {
			// If it's a duplicate (pgx.ErrNoRows because of ON CONFLICT DO NOTHING), just skip it
			if err.Error() == "no rows in result set" {
				continue
			}
			logrus.Errorf("Failed to create email %s: %v", msg.Id, err)
			continue
		}

		processedCount++
	}

	// Update last sync time
	integration, err := h.storage.GetGmailIntegration(userID)
	if err == nil {
		integration.LastSyncAt = time.Now()
		integration.UpdatedAt = time.Now()
		h.storage.UpdateGmailIntegration(integration)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Email sync completed",
		"processed_count": processedCount,
		"total_fetched":   len(messages),
	})
}

// SendEmail sends an email via Gmail
// @Summary Send email via Gmail
// @Description Sends an email using the authenticated Gmail account
// @Tags gmail
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param email body object{to=string,subject=string,body=string} true "Email data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gmail/send [post]
func (h *GmailHandler) SendEmail(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var emailData struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&emailData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.gmailService.SendEmail(c.Request.Context(), userID, emailData.To, emailData.Subject, emailData.Body)
	if err != nil {
		logrus.Errorf("Failed to send email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Email sent successfully",
		"to":      emailData.To,
	})
}

// GetStatus returns Gmail integration status
// @Summary Get Gmail integration status
// @Description Returns the current status of Gmail integration
// @Tags gmail
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gmail/status [get]
func (h *GmailHandler) GetStatus(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	integration, err := h.storage.GetGmailIntegration(userID)
	if err != nil {
		if err.Error() == "gmail integration not found" {
			c.JSON(http.StatusOK, gin.H{
				"connected": false,
				"message":   "Gmail account not connected",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get Gmail integration status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"connected":    integration.IsActive,
		"email":        integration.Email,
		"last_sync_at": integration.LastSyncAt,
		"created_at":   integration.CreatedAt,
	})
}

// Disconnect removes Gmail integration
// @Summary Disconnect Gmail account
// @Description Removes the Gmail integration and deletes stored tokens
// @Tags gmail
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /gmail/disconnect [delete]
func (h *GmailHandler) Disconnect(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	err := h.storage.DeleteGmailIntegration(userID)
	if err != nil {
		logrus.Errorf("Failed to delete Gmail integration: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect Gmail account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gmail account disconnected successfully",
	})
}

// Helper function to extract name from email string
func extractNameFromEmail(emailStr string) string {
	// Simple extraction - in production, you might want more sophisticated parsing
	if len(emailStr) > 0 && emailStr[0] == '<' {
		// Format: <email@domain.com>
		return ""
	}
	
	// Format: Name <email@domain.com> or email@domain.com
	if idx := strings.Index(emailStr, "<"); idx != -1 {
		name := strings.TrimSpace(emailStr[:idx])
		if len(name) > 0 && name[len(name)-1] == ' ' {
			name = name[:len(name)-1]
		}
		return name
	}
	
	return ""
}
