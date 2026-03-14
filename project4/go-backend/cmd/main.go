package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"ai-recruiter-assistant/internal/auth"
	"ai-recruiter-assistant/internal/metrics"
	"ai-recruiter-assistant/internal/storage"
	"ai-recruiter-assistant/internal/web/middleware"
)

// @title AI Recruiter Assistant API
// @version 1.0
// @description AI-powered recruiting assistant with email classification and reply generation
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Initialize metrics
	prometheusMetrics := metrics.NewPrometheusMetrics()
	
	// Start metrics server
	prometheusMetrics.StartMetricsServer(":9091")

	// Initialize storage
	db, err := storage.NewDatabase("postgres://postgres:admin123@localhost:5432/ai_recruiter?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer db.Close()

	// Initialize auth service
	authService := auth.NewAuthService("your-super-secret-jwt-key-change-in-production", 3600) // 1 hour TTL

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.NewMetricsMiddleware(prometheusMetrics).Middleware())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"service": "ai-recruiter-assistant",
		})
	})

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", func(c *gin.Context) {
				// TODO: Implement registration
				c.JSON(200, gin.H{"message": "Registration endpoint"})
			})
			auth.POST("/login", func(c *gin.Context) {
				// TODO: Implement login
				c.JSON(200, gin.H{"message": "Login endpoint"})
			})
			auth.POST("/refresh", func(c *gin.Context) {
				// TODO: Implement token refresh
				c.JSON(200, gin.H{"message": "Token refresh endpoint"})
			})
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(authService.AuthMiddleware())
		{
			// Profile routes
			profile := protected.Group("/profile")
			{
				profile.GET("", func(c *gin.Context) {
					// TODO: Implement get profile
					c.JSON(200, gin.H{"message": "Get profile endpoint"})
				})
				profile.PUT("", func(c *gin.Context) {
					// TODO: Implement update profile
					c.JSON(200, gin.H{"message": "Update profile endpoint"})
				})
				profile.POST("/resume", func(c *gin.Context) {
					// TODO: Implement resume upload
					c.JSON(200, gin.H{"message": "Resume upload endpoint"})
				})
				profile.GET("/resume", func(c *gin.Context) {
					// TODO: Implement get resume
					c.JSON(200, gin.H{"message": "Get resume endpoint"})
				})
			}

			// Application routes
			applications := protected.Group("/applications")
			{
				applications.GET("", func(c *gin.Context) {
					// TODO: Implement get applications
					c.JSON(200, gin.H{"applications": []interface{}{}})
				})
				applications.POST("", func(c *gin.Context) {
					// TODO: Implement create application
					c.JSON(200, gin.H{"message": "Create application endpoint"})
				})
				applications.GET("/:id", func(c *gin.Context) {
					// TODO: Implement get application
					c.JSON(200, gin.H{"message": "Get application endpoint"})
				})
				applications.PUT("/:id", func(c *gin.Context) {
					// TODO: Implement update application
					c.JSON(200, gin.H{"message": "Update application endpoint"})
				})
				applications.DELETE("/:id", func(c *gin.Context) {
					// TODO: Implement delete application
					c.JSON(200, gin.H{"message": "Delete application endpoint"})
				})
			}

			// Email routes
			emails := protected.Group("/emails")
			{
				emails.GET("", func(c *gin.Context) {
					// TODO: Implement email listing
					c.JSON(200, gin.H{"emails": []interface{}{}})
				})
				emails.POST("/import", func(c *gin.Context) {
					// TODO: Implement email import
					c.JSON(200, gin.H{"message": "Email import started"})
				})
				emails.POST("/upload", func(c *gin.Context) {
					// TODO: Implement email upload
					c.JSON(200, gin.H{"message": "Email uploaded"})
				})
			}

			// AI routes
			ai := protected.Group("/ai")
			{
				ai.POST("/classify-email", func(c *gin.Context) {
					// TODO: Implement email classification
					c.JSON(200, gin.H{"message": "Email classification endpoint"})
				})
				ai.POST("/generate-reply", func(c *gin.Context) {
					// TODO: Implement reply generation
					c.JSON(200, gin.H{"message": "Reply generation endpoint"})
				})
				ai.POST("/extract-requirements", func(c *gin.Context) {
					// TODO: Implement requirement extraction
					c.JSON(200, gin.H{"message": "Requirement extraction endpoint"})
				})
				ai.POST("/search-context", func(c *gin.Context) {
					// TODO: Implement context search
					c.JSON(200, gin.H{"message": "Context search endpoint"})
				})
			}

			// Follow-up routes
			followup := protected.Group("/followup")
			{
				followup.POST("/generate", func(c *gin.Context) {
					// TODO: Implement follow-up generation
					c.JSON(200, gin.H{"message": "Follow-up generation endpoint"})
				})
				followup.GET("/pending", func(c *gin.Context) {
					// TODO: Implement get pending follow-ups
					c.JSON(200, gin.H{"message": "Get pending follow-ups endpoint"})
				})
				followup.POST("/schedule", func(c *gin.Context) {
					// TODO: Implement schedule follow-up
					c.JSON(200, gin.H{"message": "Schedule follow-up endpoint"})
				})
			}
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Starting AI Recruiter Assistant on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
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
