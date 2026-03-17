package web

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"ai-recruiter-assistant/internal/auth"
	"ai-recruiter-assistant/internal/config"
	"ai-recruiter-assistant/internal/gmail"
	"ai-recruiter-assistant/internal/storage"
	"ai-recruiter-assistant/internal/web/handlers"
	"ai-recruiter-assistant/internal/web/middleware"
	"ai-recruiter-assistant/internal/ai"
)

// SetupRoutes configures all API routes
func SetupRoutes(router *gin.RouterGroup, db *pgxpool.Pool, redisClient *redis.Client, cfg *config.Config) {
	// Initialize storage
	storageInstance := storage.NewStorage(db, redisClient)

	// Initialize auth service
	authService := auth.NewAuthService(cfg.JWTSecret, time.Duration(cfg.JWTExpiry)*time.Hour)

	// Initialize Gmail service
	gmailService := gmail.NewGmailService(cfg.GmailClientID, cfg.GmailClientSecret, cfg.GmailRedirectURL, storageInstance)
	if gmailService == nil {
		logrus.Error("Failed to initialize Gmail service - OAuth2 credentials missing")
		// Continue without Gmail routes
	} else {
		logrus.Info("Gmail service initialized successfully")
	}

	// Initialize AI service
	ollamaService := ai.NewOllamaService(cfg.OllamaURL)
	aiService := ai.NewAIService(ollamaService)
	logrus.Info("AI service initialized successfully")

	// Initialize handlers
	authHandler := auth.NewAuthHandler(authService, storageInstance)
	emailHandler := handlers.NewEmailHandler(storageInstance, authService)
	profileHandler := handlers.NewProfileHandler(storageInstance, authService)
	applicationHandler := handlers.NewApplicationHandler(storageInstance, authService)
	aiHandler := handlers.NewAIHandler(storageInstance, authService)
	
	var gmailHandler *handlers.GmailHandler
	if gmailService != nil {
		gmailHandler = handlers.NewGmailHandler(storageInstance, gmailService, authService, aiService)
	}

	// Public routes (no authentication required)
	public := router.Group("/auth")
	{
		public.POST("/register", authHandler.Register)
		public.POST("/login", authHandler.Login)
		public.POST("/logout", middleware.LogoutHandler(redisClient))
		
		// Add Gmail routes only if service is available
		if gmailHandler != nil {
			public.GET("/gmail", gmailHandler.AuthRequest) // OAuth2 initiation
			public.GET("/gmail/callback", gmailHandler.AuthCallback) // OAuth2 callback
		}
	}

	// Protected routes (authentication required)
	protected := router.Group("/")
	protected.Use(corsMiddleware())
	protected.Use(authService.AuthMiddleware())
	{
		// Auth routes
		auth := protected.Group("/auth")
		{
			auth.GET("/profile", authHandler.GetProfile)
			auth.PUT("/profile", authHandler.UpdateProfile)
		}

		// Gmail routes (only if service is available)
		if gmailHandler != nil {
			gmail := protected.Group("/gmail")
			{
				gmail.GET("/status", gmailHandler.GetStatus)
				gmail.POST("/sync", gmailHandler.SyncEmails)
				gmail.POST("/send", gmailHandler.SendEmail)
				gmail.POST("/generate-response", gmailHandler.GenerateResponse)
				gmail.DELETE("/disconnect", gmailHandler.Disconnect)
			}
		}

		// Email routes
		emails := protected.Group("/emails")
		{
			emails.GET("", emailHandler.GetEmails)
			emails.GET("/stats", emailHandler.GetEmailStats)
			emails.GET("/:id", emailHandler.GetEmail)
			emails.POST("/import", emailHandler.ImportEmails)
			emails.POST("/upload", emailHandler.UploadEmail)
		}

		// Profile routes
		profile := protected.Group("/profile")
		{
			profile.GET("", profileHandler.GetProfile)
			profile.PUT("", profileHandler.UpdateProfile)
			profile.POST("/resume", profileHandler.UploadResume)
			profile.GET("/resume", profileHandler.GetResume)
		}

		// Application routes
		applications := protected.Group("/applications")
		{
			applications.GET("", applicationHandler.GetApplications)
			applications.POST("", applicationHandler.CreateApplication)
			applications.GET("/:id", applicationHandler.GetApplication)
			applications.PUT("/:id", applicationHandler.UpdateApplication)
			applications.DELETE("/:id", applicationHandler.DeleteApplication)
		}

		// AI routes
		ai := protected.Group("/ai")
		{
			ai.POST("/classify-email", aiHandler.ClassifyEmail)
			ai.POST("/generate-reply", aiHandler.GenerateReply)
			ai.POST("/extract-requirements", aiHandler.ExtractRequirements)
			ai.POST("/search-context", aiHandler.SearchContext)
			ai.GET("/gpu-status", aiHandler.GetGPUStatus)
		}
	}

	// Optional auth routes (work with or without authentication)
	optional := router.Group("/")
	optional.Use(authService.OptionalAuthMiddleware())
	{
		// Add any routes that should work with or without authentication
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
