package web

import (
	"fmt"
	"strings"
	"time"

	"private-knowledge-base-go/internal/auth"
	"private-knowledge-base-go/internal/ingestion"
	"private-knowledge-base-go/internal/rag"
	"private-knowledge-base-go/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	router *gin.Engine,
	authService *auth.Service,
	ingestionService *ingestion.Service,
	ragService *rag.Service,
	db *storage.PostgresDB,
	logger *logrus.Logger,
) {
	// Create handlers
	authMiddleware := auth.NewMiddleware(authService, logger)
	chatHandler := NewChatHandler(ragService, db, logger)
	documentsHandler := NewDocumentsHandler(ingestionService, db, logger)
	healthHandler := NewHealthHandler(db, logger)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check routes (no auth required)
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	router.GET("/live", healthHandler.Live)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handleLogin(authService))
			auth.POST("/refresh", authMiddleware.RequireAuth(), handleRefresh(authService))
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// User routes
			user := protected.Group("/user")
			{
				user.GET("/profile", handleUserProfile(authService))
			}

			// Chat routes
			chat := protected.Group("/chat")
			{
				chat.POST("/sessions", chatHandler.CreateSession)
				chat.GET("/sessions", chatHandler.GetSessions)
				chat.GET("/sessions/:id", chatHandler.GetSession)
				chat.DELETE("/sessions/:id", chatHandler.DeleteSession)
				chat.GET("/sessions/:id/messages", chatHandler.GetMessages)
				chat.POST("/", chatHandler.Chat)
				chat.POST("/stream", chatHandler.StreamChat)
				// Handle CORS preflight explicitly for streaming endpoint
				chat.OPTIONS("/stream", func(c *gin.Context) {
					c.Header("Access-Control-Allow-Origin", "*")
					c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Cache-Control, Accept")
					c.Header("Access-Control-Allow-Methods", "POST, OPTIONS")
					c.Status(204)
				})
			}

			// Document routes
			documents := protected.Group("/documents")
			{
				documents.POST("/upload", documentsHandler.UploadDocument)
				documents.GET("/", documentsHandler.ListDocuments)
				documents.GET("", documentsHandler.ListDocuments) // Handle /documents without trailing slash
				documents.GET("/:id", documentsHandler.GetDocument)
				documents.DELETE("/:id", documentsHandler.DeleteDocument)
				documents.POST("/:id/reindex", documentsHandler.ReindexDocument)
				documents.GET("/:id/chunks", documentsHandler.GetDocumentChunks)
			}

			// Metrics routes
			metrics := protected.Group("/metrics")
			{
				metrics.GET("/", healthHandler.Metrics)
			}
		}

		// Public routes (optional auth)
		public := v1.Group("")
		public.Use(authMiddleware.OptionalAuth())
		{
			// Add any public routes here
		}
	}
}

// Authentication handlers

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
}

type RefreshResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserProfileResponse struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
	JoinedAt time.Time `json:"joined_at"`
}

// handleUserProfile returns user profile information
// @Summary Get user profile
// @Description Returns the current user's profile information
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} UserProfileResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/user/profile [get]
func handleUserProfile(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user information from context (set by auth middleware)
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, ErrorResponse{
				Error: "User not authenticated",
			})
			return
		}

		usernameInterface, exists := c.Get("username")
		if !exists {
			c.JSON(401, ErrorResponse{
				Error: "User not authenticated",
			})
			return
		}

		// Type assertions - handle UUID for user_id
		var userID string
		switch v := userIDInterface.(type) {
		case string:
			userID = v
		case uuid.UUID:
			userID = v.String()
		default:
			c.JSON(500, ErrorResponse{
				Error: "Invalid user ID type",
			})
			return
		}

		username, ok := usernameInterface.(string)
		if !ok {
			c.JSON(500, ErrorResponse{
				Error: "Invalid username type",
			})
			return
		}

		// Return user profile (mock data for now)
		c.JSON(200, UserProfileResponse{
			UserID:   userID,
			Username: username,
			Email:    fmt.Sprintf("%s@example.com", username),
			FullName: fmt.Sprintf("%s User", username),
			JoinedAt: time.Now().AddDate(0, -1, 0), // Joined 1 month ago
		})
	}
}

// handleLogin handles user login
// @Summary User login
// @Description Authenticates a user and returns a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func handleLogin(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, ErrorResponse{
				Error:   "Invalid request format",
				Details: err.Error(),
			})
			return
		}

		// TODO: Implement actual authentication logic
		// For now, accept any username/password and generate a token
		// In production, you would validate against a user database
		
		// Generate mock user ID
	userIDStr := "550e8400-e29b-41d4-a716-446655440000"
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(500, ErrorResponse{
			Error:   "Failed to parse user ID",
			Details: err.Error(),
		})
		return
	}
		
		token, err := authService.GenerateToken(userID, req.Username)
		if err != nil {
			c.JSON(500, ErrorResponse{
				Error:   "Failed to generate token",
				Details: err.Error(),
			})
			return
		}

		c.JSON(200, LoginResponse{
			Token:     token,
			ExpiresAt: time.Now().Add(24 * time.Hour),
			UserID:    userID.String(),
			Username:  req.Username,
		})
	}
}

// handleRefresh handles token refresh
// @Summary Refresh token
// @Description Refreshes an existing JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} RefreshResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func handleRefresh(authService *auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, ErrorResponse{
				Error: "Authorization header required",
			})
			return
		}

		// Extract token from "Bearer <token>" format
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(401, ErrorResponse{
				Error: "Invalid authorization header format",
			})
			return
		}

		tokenString := tokenParts[1]

		// Refresh token
		newToken, err := authService.RefreshToken(tokenString)
		if err != nil {
			c.JSON(401, ErrorResponse{
				Error:   "Invalid or expired token",
				Details: err.Error(),
			})
			return
		}

		c.JSON(200, RefreshResponse{
			Token:     newToken,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
	}
}
