package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"private-knowledge-base-go/internal/auth"
	"private-knowledge-base-go/internal/config"
	"private-knowledge-base-go/internal/ingestion"
	"private-knowledge-base-go/internal/rag"
	"private-knowledge-base-go/internal/storage"
	"private-knowledge-base-go/internal/web"
	"private-knowledge-base-go/pkg/logger"
	"private-knowledge-base-go/pkg/metrics"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title Private Knowledge Base API
// @version 1.0
// @description AI-powered knowledge base with RAG capabilities
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.LogLevel)

	// Initialize database
	db, err := storage.NewPostgresDB(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := db.RunMigrations(); err != nil {
		logger.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize services
	authService := auth.NewService(cfg.JWTSecret, logger)
	ragService := rag.NewService(db, cfg.OllamaURL, logger)
	ingestionService := ingestion.NewService(db, logger, ragService)
	_ = metrics.NewCollector() // TODO: Use metrics collector

	// Setup Gin router
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.New()
	
	router.Use(corsMiddleware(cfg.AllowedOrigins))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// Ensure preflight requests always succeed
	router.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNoContent)
	})

	// Setup routes
	web.SetupRoutes(router, authService, ingestionService, ragService, db, logger)

	// Metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting server on port %d", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func corsMiddleware(allowed []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		allowAll := len(allowed) == 0
		originAllowed := false
		if origin != "" {
			for _, o := range allowed {
				if o == "*" || o == origin {
					originAllowed = true
					break
				}
			}
		}

		if allowAll || originAllowed {
			if origin != "" && originAllowed && origin != "*" {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
			} else {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH, HEAD")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With, Cache-Control, Pragma")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
