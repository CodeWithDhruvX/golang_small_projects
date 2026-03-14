package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ai-recruiter-assistant/internal/auth"
	"ai-recruiter-assistant/internal/storage"
)

// ProfileHandler handles profile-related HTTP requests
type ProfileHandler struct {
	storage     storage.StorageInterface
	authService *auth.AuthService
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(storage storage.StorageInterface, authService *auth.AuthService) *ProfileHandler {
	return &ProfileHandler{
		storage:     storage,
		authService: authService,
	}
}

// GetProfile retrieves the user's profile
// @Summary Get user profile
// @Description Retrieve the authenticated user's profile
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} auth.User
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile [get]
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.storage.GetUserByID(userID)
	if err != nil {
		logrus.Errorf("Failed to get user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateProfile updates the user's profile
// @Summary Update user profile
// @Description Update the authenticated user's profile
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body auth.User true "Profile data"
// @Success 200 {object} auth.User
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user auth.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure user can only update their own profile
	user.ID = userID

	err := h.storage.UpdateUser(&user)
	if err != nil {
		logrus.Errorf("Failed to update user profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, &user)
}

// UploadResume handles resume file upload
// @Summary Upload resume
// @Description Upload a resume file (PDF) and process it for RAG
// @Tags profile
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "Resume file (PDF)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile/resume [post]
func (h *ProfileHandler) UploadResume(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Check file type (should be PDF)
	if file.Header.Get("Content-Type") != "application/pdf" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF files are allowed"})
		return
	}

	// Check file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large (max 10MB)"})
		return
	}

	// Read file content
	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer fileContent.Close()

	// Read all bytes
	content := make([]byte, file.Size)
	_, err = fileContent.Read(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file content"})
		return
	}

	// TODO: Process the resume using resume processor
	// 1. Extract text from PDF using unipdf
	// 2. Generate embeddings using Ollama
	// 3. Store in documents table for RAG
	
	logrus.Infof("Resume uploaded for user: %s, file: %s, size: %d bytes", userID, file.Filename, file.Size)
	
	c.JSON(http.StatusOK, gin.H{
		"message":    "Resume uploaded and processed successfully",
		"filename":   file.Filename,
		"size":       file.Size,
		"status":     "processed",
		"chunks_processed": 0, // TODO: Return actual chunk count
	})
}

// GetResume retrieves the user's resume
// @Summary Get resume
// @Description Retrieve the user's uploaded resume
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /profile/resume [get]
func (h *ProfileHandler) GetResume(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// TODO: Retrieve resume information from storage
	// This would return the resume file or metadata
	
	logrus.Infof("Resume retrieval requested for user: %s", userID)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Resume retrieval not implemented yet",
		"user_id": userID,
	})
}
