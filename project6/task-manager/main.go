package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Task struct {
	BaseModel
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date,omitempty"`
	AssignedTo  string `json:"assigned_to,omitempty"`
}

var (
	tasks   = make(map[int]Task)
	users   = make(map[string]User)
	projects = make(map[string]Project)
	nextTaskID    = 1
	tasksMu   sync.RWMutex
	usersMu   sync.RWMutex	
	projectsMu sync.RWMutex

	// Security components
	logger            *Logger
	errorHandler      *ErrorHandler
	validator         *Validator
	encryptionService *EncryptionService
	authService       *AuthService
	sslManager        *SSLManager
)

func main() {
	// Initialize security components
	initializeSecurityComponents()
	defer logger.Close()

	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// Add middleware in the correct order
	r.Use(RecoveryMiddleware(errorHandler))
	r.Use(RequestIDMiddleware())
	r.Use(LoggingMiddleware(logger))
	r.Use(SecurityHeadersMiddleware())
	r.Use(CORSMiddleware([]string{"*"}))
	r.Use(TimeoutMiddleware(30 * time.Second))

	// Rate limiting
	r.Use(RateLimitMiddleware(60, 10)) // 60 requests per minute, burst of 10

	// SSL/HTTPS redirect middleware
	if sslManager.config.RedirectToHTTPS {
		r.Use(sslManager.HTTPSToHTTPSRedirectMiddleware())
	}

	// Initialize sample data
	initializeSampleData()

	// Public routes (no authentication required)
	public := r.Group("/api/v1/public")
	{
		public.GET("/health", healthCheck)
		public.GET("/metrics", metrics)
		public.POST("/login", login)
		public.POST("/register", register)
	}

	// Protected routes (authentication required)
	protected := r.Group("/api/v1")
	protected.Use(JWTAuthenticationMiddleware("your-super-secret-jwt-key", errorHandler))
	protected.Use(ValidationMiddleware(validator, errorHandler))
	protected.Use(EncryptionMiddleware(encryptionService, errorHandler))
	{
		// Task routes
		protected.GET("/tasks", getTasks)
		protected.GET("/tasks/:id", getTask)
		protected.POST("/tasks", createTask)
		protected.PUT("/tasks/:id", updateTask)
		protected.DELETE("/tasks/:id", deleteTask)

		// User routes (admin only)
		admin := protected.Group("/users")
		admin.Use(RequireRoleMiddleware("admin", errorHandler))
		{
			admin.GET("", getUsers)
			admin.GET("/:id", getUser)
			admin.POST("", createUser)
			admin.PUT("/:id", updateUser)
			admin.DELETE("/:id", deleteUser)
		}

		// Project routes
		protected.GET("/projects", getProjects)
		protected.GET("/projects/:id", getProject)
		protected.POST("/projects", createProject)
		protected.PUT("/projects/:id", updateProject)
		protected.DELETE("/projects/:id", deleteProject)
	}

	// Setup SSL and start server
	if err := sslManager.SetupSSL(); err != nil {
		logger.Fatal("Failed to setup SSL", map[string]interface{}{"error": err.Error()})
	}

	if sslManager.config.UseHTTPS {
		// Start HTTPS server
		tlsConfig, err := sslManager.CreateTLSConfig()
		if err != nil {
			logger.Fatal("Failed to create TLS config", map[string]interface{}{"error": err.Error()})
		}

		httpServer := &http.Server{
			Addr:      fmt.Sprintf(":%d", sslManager.config.HTTPSPort),
			Handler:   r,
			TLSConfig: tlsConfig,
		}

		logger.Info("Starting HTTPS server", map[string]interface{}{
			"port": sslManager.config.HTTPSPort,
			"cert_file": sslManager.config.CertFile,
			"key_file": sslManager.config.KeyFile,
		})

		if err := httpServer.ListenAndServeTLS("", ""); err != nil {
			logger.Fatal("Failed to start HTTPS server", map[string]interface{}{"error": err.Error()})
		}
	} else {
		// Start HTTP server (for development)
		logger.Info("Starting HTTP server (development mode)", map[string]interface{}{
			"port": sslManager.config.Port,
		})

		if err := r.Run(fmt.Sprintf(":%d", sslManager.config.Port)); err != nil {
			logger.Fatal("Failed to start HTTP server", map[string]interface{}{"error": err.Error()})
		}
	}
}

func initializeSecurityComponents() {
	var err error

	// Initialize logger
	logger, err = NewLogger(LogLevelInfo, true)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize encryption service
	encryptionService = NewEncryptionService("your-encryption-key-32-chars-long")

	// Initialize error handler
	errorHandler = NewErrorHandler(encryptionService, logger)

	// Initialize validator
	validator = NewValidator(encryptionService)

	// Initialize auth service
	authService = NewAuthService("your-super-secret-jwt-key", 24*time.Hour, logger, errorHandler)

	// Initialize SSL manager
	sslConfig := GetDefaultSSLConfig()
	sslConfig.UseHTTPS = os.Getenv("USE_HTTPS") != "false" // Default to HTTPS
	sslManager = NewSSLManager(sslConfig, logger)

	logger.Info("Security components initialized", nil)
}

func initializeSampleData() {
	tasksMu.Lock()
	defer tasksMu.Unlock()

	usersMu.Lock()
	defer usersMu.Unlock()

	projectsMu.Lock()
	defer projectsMu.Unlock()

	// Initialize sample tasks
	tasks[1] = Task{
		BaseModel: BaseModel{
			ID:        1,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-12 * time.Hour),
		},
		Title:       "Learn Go Security Features",
		Description: "Study encryption, validation, and error handling",
		Completed:   false,
		Priority:    "high",
		DueDate:     "2026-03-25",
		AssignedTo:  "admin@example.com",
	}
	tasks[2] = Task{
		BaseModel: BaseModel{
			ID:        2,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-6 * time.Hour),
		},
		Title:       "Build Secure REST API",
		Description: "Implement security best practices",
		Completed:   true,
		Priority:    "medium",
		DueDate:     "2026-03-20",
		AssignedTo:  "user@example.com",
	}
	nextTaskID = 3

	// Initialize sample users
	users["admin"] = User{
		BaseModel: BaseModel{
			ID:        "admin",
			CreatedAt: time.Now().Add(-72 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  "hashed_admin_password",
		Role:      "admin",
		IsActive:  true,
		LastLogin: time.Now().Add(-2 * time.Hour).Format(time.RFC3339),
	}
	users["user"] = User{
		BaseModel: BaseModel{
			ID:        "user",
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-30 * time.Minute),
		},
		Username:  "user",
		Email:     "user@example.com",
		Password:  "hashed_user_password",
		Role:      "user",
		IsActive:  true,
		LastLogin: time.Now().Add(-4 * time.Hour).Format(time.RFC3339),
	}

	// Initialize sample projects
	projects["proj1"] = Project{
		BaseModel: BaseModel{
			ID:        "proj1",
			CreatedAt: time.Now().Add(-120 * time.Hour),
			UpdatedAt: time.Now().Add(-24 * time.Hour),
		},
		Name:        "Secure API Development",
		Description: "Building a secure REST API with Go",
		Status:      "active",
		TeamMembers: []string{"admin", "user"},
		Budget:      10000.0,
	}

	logger.Info("Sample data initialized", map[string]interface{}{
		"tasks":    len(tasks),
		"users":    len(users),
		"projects": len(projects),
	})
}

