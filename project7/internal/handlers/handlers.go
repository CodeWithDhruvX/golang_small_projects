package handlers

import (
	"net/http"
	"project7/internal/database"
	"project7/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	PostgresDB *database.PostgresDB
	MongoDB    *database.MongoDB
}

func NewHandler(postgresDB *database.PostgresDB, mongoDB *database.MongoDB) *Handler {
	return &Handler{
		PostgresDB: postgresDB,
		MongoDB:    mongoDB,
	}
}

// PostgreSQL Handlers

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.PostgresDB.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Log the user creation event
	logEntry := &models.LogEntry{
		Level:     "info",
		Message:   "User created",
		Service:   "user-service",
		UserID:    &user.ID,
		RequestID: c.GetHeader("X-Request-ID"),
		Timestamp: time.Now(),
	}
	h.MongoDB.CreateLogEntry(logEntry)

	c.JSON(http.StatusCreated, user)
}

func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := h.PostgresDB.GetUser(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.PostgresDB.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = uint(id)
	if err := h.PostgresDB.UpdateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.PostgresDB.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func (h *Handler) CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.PostgresDB.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *Handler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.PostgresDB.GetPost(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) GetAllPosts(c *gin.Context) {
	posts, err := h.PostgresDB.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Raw SQL Query Handlers

func (h *Handler) GetUsersByAgeRange(c *gin.Context) {
	minAge, err1 := strconv.Atoi(c.DefaultQuery("min_age", "0"))
	maxAge, err2 := strconv.Atoi(c.DefaultQuery("max_age", "100"))
	
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid age parameters"})
		return
	}

	users, err := h.PostgresDB.GetUsersByAgeRange(minAge, maxAge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *Handler) GetUserEmails(c *gin.Context) {
	emails, err := h.PostgresDB.GetUserEmails()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"emails": emails})
}

func (h *Handler) GetUsersWithPostCount(c *gin.Context) {
	results, err := h.PostgresDB.GetUsersWithPostCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user_stats": results})
}

func (h *Handler) SearchPosts(c *gin.Context) {
	searchTerm := c.Query("q")
	if searchTerm == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search term is required"})
		return
	}

	posts, err := h.PostgresDB.SearchPosts(searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func (h *Handler) GetCategoryStats(c *gin.Context) {
	stats, err := h.PostgresDB.GetCategoryStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"category_stats": stats})
}

func (h *Handler) GetActiveUsersWithFilters(c *gin.Context) {
	var ageMin, ageMax *int
	
	if minAgeStr := c.Query("min_age"); minAgeStr != "" {
		if minAge, err := strconv.Atoi(minAgeStr); err == nil {
			ageMin = &minAge
		}
	}
	
	if maxAgeStr := c.Query("max_age"); maxAgeStr != "" {
		if maxAge, err := strconv.Atoi(maxAgeStr); err == nil {
			ageMax = &maxAge
		}
	}
	
	emailDomain := c.Query("email_domain")

	users, err := h.PostgresDB.GetActiveUsersWithFilters(ageMin, ageMax, emailDomain)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// MongoDB Handlers

func (h *Handler) CreateUserProfile(c *gin.Context) {
	var profile models.UserProfile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	if err := h.MongoDB.CreateUserProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

func (h *Handler) GetUserProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	profile, err := h.MongoDB.GetUserProfile(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User profile not found"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *Handler) CreateLogEntry(c *gin.Context) {
	var entry models.LogEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry.Timestamp = time.Now()

	if err := h.MongoDB.CreateLogEntry(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *Handler) GetLogEntries(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	if limit > 1000 {
		limit = 1000
	}

	entries, err := h.MongoDB.GetLogEntries(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"entries": entries})
}

func (h *Handler) CreateAnalyticsEvent(c *gin.Context) {
	var event models.AnalyticsEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event.Timestamp = time.Now()
	event.IPAddress = c.ClientIP()
	event.UserAgent = c.GetHeader("User-Agent")

	if err := h.MongoDB.CreateAnalyticsEvent(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (h *Handler) GetAnalyticsEvents(c *gin.Context) {
	eventType := c.Query("event_type")
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "100"), 10, 64)
	if limit > 1000 {
		limit = 1000
	}

	events, err := h.MongoDB.GetAnalyticsEvents(eventType, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// Health check handler
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now(),
		"databases": gin.H{
			"postgres": "connected",
			"mongodb":  "connected",
		},
	})
}
