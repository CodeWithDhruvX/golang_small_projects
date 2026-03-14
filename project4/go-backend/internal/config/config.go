package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Port         int
	Environment  string
	JWTSecret    string
	JWTExpiry    int // in hours

	// Database configuration
	DatabaseURL string

	// Redis configuration
	RedisURL string

	// Ollama configuration
	OllamaURL string

	// Email configuration
	IMAPServer string
	IMAPPort   int
	IMAPUser   string
	IMAPPass   string

	// Gmail OAuth2 configuration
	GmailClientID     string
	GmailClientSecret string
	GmailRedirectURL  string

	// File upload configuration
	MaxFileSize int64 // in bytes
	UploadPath  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:         getEnvInt("PORT", 8080),
		Environment:  getEnv("ENVIRONMENT", "development"),
		JWTSecret:    getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
		JWTExpiry:    getEnvInt("JWT_EXPIRY", 24),
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ai_recruiter?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		OllamaURL:    getEnv("OLLAMA_URL", "http://localhost:11434"),
		IMAPServer:   getEnv("IMAP_SERVER", "imap.gmail.com"),
		IMAPPort:     getEnvInt("IMAP_PORT", 993),
		IMAPUser:     getEnv("IMAP_USER", ""),
		IMAPPass:     getEnv("IMAP_PASS", ""),
		GmailClientID:     getEnv("GMAIL_CLIENT_ID", ""),
		GmailClientSecret: getEnv("GMAIL_CLIENT_SECRET", ""),
		GmailRedirectURL:  getEnv("GMAIL_REDIRECT_URL", "http://localhost:8080/api/v1/auth/gmail/callback"),
		MaxFileSize:  getEnvInt64("MAX_FILE_SIZE", 10*1024*1024), // 10MB
		UploadPath:   getEnv("UPLOAD_PATH", "./uploads"),
	}

	return cfg, nil
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt gets environment variable as integer with default value
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvInt64 gets environment variable as int64 with default value
func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
