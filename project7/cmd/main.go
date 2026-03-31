package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project7/internal/config"
	"project7/internal/database"
	"project7/internal/handlers"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize databases
	postgresDB, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}
	defer postgresDB.Close()

	mongoDB, err := database.NewMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize MongoDB: %v", err)
	}
	defer mongoDB.Close(context.Background())

	// Initialize handlers
	handler := handlers.NewHandler(postgresDB, mongoDB)

	// Setup Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check endpoint
	router.GET("/health", handler.HealthCheck)

	// API routes
	api := router.Group("/api/v1")
	{
		// PostgreSQL routes - Users
		users := api.Group("/users")
		{
			users.POST("", handler.CreateUser)
			users.GET("", handler.GetAllUsers)
			users.GET("/:id", handler.GetUser)
			users.PUT("/:id", handler.UpdateUser)
			users.DELETE("/:id", handler.DeleteUser)
		}

		// PostgreSQL routes - Posts
		posts := api.Group("/posts")
		{
			posts.POST("", handler.CreatePost)
			posts.GET("", handler.GetAllPosts)
			posts.GET("/:id", handler.GetPost)
			posts.GET("/search", handler.SearchPosts)
		}

		// Raw SQL Query endpoints
		sql := api.Group("/sql")
		{
			sql.GET("/users/age-range", handler.GetUsersByAgeRange)
			sql.GET("/users/emails", handler.GetUserEmails)
			sql.GET("/users/post-count", handler.GetUsersWithPostCount)
			sql.GET("/categories/stats", handler.GetCategoryStats)
			sql.GET("/users/filters", handler.GetActiveUsersWithFilters)
		}

		// MongoDB routes - User Profiles
		profiles := api.Group("/profiles")
		{
			profiles.POST("", handler.CreateUserProfile)
			profiles.GET("/:user_id", handler.GetUserProfile)
		}

		// MongoDB routes - Logs
		logs := api.Group("/logs")
		{
			logs.POST("", handler.CreateLogEntry)
			logs.GET("", handler.GetLogEntries)
		}

		// MongoDB routes - Analytics
		analytics := api.Group("/analytics")
		{
			analytics.POST("/events", handler.CreateAnalyticsEvent)
			analytics.GET("/events", handler.GetAnalyticsEvents)
		}
	}

	// Start server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Request-ID")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
