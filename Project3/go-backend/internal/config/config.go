package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port        int
	Environment string
	LogLevel    string
	AllowedOrigins []string

	// Database configuration
	DatabaseURL string

	// Authentication
	JWTSecret string

	// AI/ML configuration
	OllamaURL string

	// Vector embedding configuration
	EmbeddingModel string
	EmbeddingURL   string

	// File upload configuration
	MaxFileSize    int64
	UploadPath     string
	AllowedTypes   []string

	// Redis configuration
	RedisURL string

	// Monitoring
	MetricsEnabled bool
}

// Load configuration from environment variables and .env file
func Load() (*Config, error) {
	cfg := &Config{
		Port:           8080,
		Environment:    "development",
		LogLevel:       "info",
		AllowedOrigins: []string{
			"http://localhost:4200",
			"http://127.0.0.1:4200",
			"http://localhost:4300",
			"http://127.0.0.1:4300",
		},
		DatabaseURL:    "postgres://postgres:postgres@localhost:5432/knowledge_base?sslmode=disable",
		JWTSecret:      "your-secret-key-change-in-production",
		OllamaURL:      "http://localhost:11434",
		EmbeddingModel: "nomic-embed-text",
		EmbeddingURL:   "http://localhost:11434",
		MaxFileSize:    50 * 1024 * 1024, // 50MB
		UploadPath:     "./uploads",
		AllowedTypes:   []string{"pdf", "txt", "md", "go"},
		RedisURL:       "redis://localhost:6379",
		MetricsEnabled: true,
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		fmt.Sscanf(port, "%d", &cfg.Port)
	}

	if env := os.Getenv("ENVIRONMENT"); env != "" {
		cfg.Environment = env
	}

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	if allowed := os.Getenv("ALLOWED_ORIGINS"); allowed != "" {
		var list []string
		start := 0
		for i := 0; i <= len(allowed); i++ {
			if i == len(allowed) || allowed[i] == ',' {
				part := allowed[start:i]
				for len(part) > 0 && (part[0] == ' ' || part[0] == '\t') {
					part = part[1:]
				}
				for len(part) > 0 && (part[len(part)-1] == ' ' || part[len(part)-1] == '\t') {
					part = part[:len(part)-1]
				}
				if part != "" {
					list = append(list, part)
				}
				start = i + 1
			}
		}
		if len(list) > 0 {
			cfg.AllowedOrigins = list
		}
	}

	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		cfg.DatabaseURL = dbURL
	}

	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		cfg.JWTSecret = jwtSecret
	}

	if ollamaURL := os.Getenv("OLLAMA_URL"); ollamaURL != "" {
		cfg.OllamaURL = ollamaURL
	}

	if embeddingModel := os.Getenv("EMBEDDING_MODEL"); embeddingModel != "" {
		cfg.EmbeddingModel = embeddingModel
	}

	if embeddingURL := os.Getenv("EMBEDDING_URL"); embeddingURL != "" {
		cfg.EmbeddingURL = embeddingURL
	}

	if uploadPath := os.Getenv("UPLOAD_PATH"); uploadPath != "" {
		cfg.UploadPath = uploadPath
	}

	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		cfg.RedisURL = redisURL
	}

	if metricsEnabled := os.Getenv("METRICS_ENABLED"); metricsEnabled != "" {
		cfg.MetricsEnabled = metricsEnabled == "true"
	}

	// Validate required configuration
	if cfg.JWTSecret == "your-secret-key-change-in-production" && cfg.Environment == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production")
	}

	return cfg, nil
}
