package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware(logger *Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetString("request_id")

		// Log request
		logger.LogRequest(
			c.Request.Method,
			c.Request.URL.Path,
			c.Request.UserAgent(),
			c.Request.RemoteAddr,
			requestID,
		)

		// Process request
		c.Next()

		// Log response
		duration := time.Since(start)
		logger.LogResponse(c.Writer.Status(), duration, requestID)
	}
}

// SecurityHeadersMiddleware adds security headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent content type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")
		
		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Strict transport security (only for HTTPS)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		
		// Content security policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Hide server information
		c.Header("Server", "SecureAPI")
		
		c.Next()
	}
}

// RateLimitMiddleware implements rate limiting
func RateLimitMiddleware(requestsPerMinute int, burstSize int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerMinute), burstSize)
	
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		
		if !limiter.Allow() {
			apiErr := errorHandler.NewAPIError(ErrorCodeRateLimit, "Rate limit exceeded", requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// ValidationMiddleware validates request payloads
func ValidationMiddleware(validator *Validator, errorHandler *ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		
		// Only validate POST, PUT, PATCH requests
		if !isValidationRequired(c.Request.Method) {
			c.Next()
			return
		}
		
		// Read and validate based on the expected entity type
		entityType := getEntityTypeFromPath(c.Request.URL.Path)
		var entity Entity
		
		switch entityType {
		case "task":
			entity = &Task{}
		case "user":
			entity = &User{}
		case "project":
			entity = &Project{}
		default:
			c.Next()
			return
		}
		
		if err := c.ShouldBindJSON(entity); err != nil {
			apiErr := errorHandler.NewAPIError(ErrorCodeBadRequest, "Invalid JSON format", requestID).
				WithDetails("validation_error", err.Error())
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}
		
		// Validate entity
		validationResult := validator.ValidateEntity(entity)
		if !validationResult.IsValid {
			apiErr := errorHandler.NewValidationError(validationResult, requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}
		
		// Store validated entity in context
		c.Set("validated_entity", entity)
		c.Next()
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}
		
		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		
		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	}
}

// EncryptionMiddleware encrypts/decrypts sensitive data
func EncryptionMiddleware(encryptionService *EncryptionService, errorHandler *ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		
		// Get validated entity from context
		entity, exists := c.Get("validated_entity")
		if !exists {
			c.Next()
			return
		}
		
		// Encrypt sensitive fields before processing
		if err := encryptSensitiveFields(entity, encryptionService); err != nil {
			apiErr := errorHandler.HandleError(err, requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}
		
		// Store encrypted entity
		c.Set("encrypted_entity", entity)
		c.Next()
		
		// Decrypt sensitive fields in response
		if c.Writer.Status() < 400 {
			if err := decryptSensitiveFields(entity, encryptionService); err != nil {
				// Log error but don't fail the request
				errorHandler.logger.ErrorWithRequestID("Failed to decrypt response data", map[string]interface{}{
					"error": err.Error(),
				}, requestID)
			}
		}
	}
}

// JWTAuthenticationMiddleware handles JWT authentication
func JWTAuthenticationMiddleware(jwtSecret string, errorHandler *ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		
		// Skip authentication for certain paths
		if isPublicPath(c.Request.URL.Path) {
			c.Next()
			return
		}
		
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiErr := errorHandler.NewUnauthorizedError("Authorization header required", requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}
		
		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			apiErr := errorHandler.NewUnauthorizedError("Invalid authorization header format", requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}
		
		token := tokenParts[1]
		
		// Validate JWT token
		claims, err := validateJWTToken(token, jwtSecret)
		if err != nil {
			apiErr := errorHandler.NewUnauthorizedError("Invalid token", requestID).
				WithDetails("token_error", err.Error())
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
			return
		}
		
		// Store user claims in context
		c.Set("user_claims", claims)
		c.Next()
	}
}

// RecoveryMiddleware recovers from panics and handles them gracefully
func RecoveryMiddleware(errorHandler *ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			requestID := c.GetString("request_id")
			apiErr := errorHandler.RecoverFromPanic(requestID)
			if apiErr.Code != "" {
				c.JSON(apiErr.HTTPStatus(), apiErr)
				c.Abort()
			}
		}()
		c.Next()
	}
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetString("request_id")
		
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()
		
		c.Request = c.Request.WithContext(ctx)
		
		finished := make(chan struct{})
		go func() {
			c.Next()
			close(finished)
		}()
		
		select {
		case <-finished:
			// Request completed normally
		case <-ctx.Done():
			// Request timed out
			apiErr := errorHandler.NewAPIError(ErrorCodeTimeout, "Request timeout", requestID)
			c.JSON(apiErr.HTTPStatus(), apiErr)
			c.Abort()
		}
	}
}

// Helper functions

func generateRequestID() string {
	return fmt.Sprintf("req_%d_%s", time.Now().UnixNano(), randomString(8))
}

func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func isValidationRequired(method string) bool {
	return method == "POST" || method == "PUT" || method == "PATCH"
}

func getEntityTypeFromPath(path string) string {
	if strings.Contains(path, "/tasks") {
		return "task"
	} else if strings.Contains(path, "/users") {
		return "user"
	} else if strings.Contains(path, "/projects") {
		return "project"
	}
	return ""
}

func isPublicPath(path string) bool {
	publicPaths := []string{"/health", "/metrics", "/docs", "/login", "/register"}
	for _, publicPath := range publicPaths {
		if strings.HasPrefix(path, publicPath) {
			return true
		}
	}
	return false
}

func encryptSensitiveFields(entity interface{}, encryptionService *EncryptionService) error {
	switch e := entity.(type) {
	case *User:
		if e.Password != "" {
			encrypted, err := encryptionService.EncryptField(e.Password)
			if err != nil {
				return fmt.Errorf("failed to encrypt password: %w", err)
			}
			e.Password = encrypted
		}
	case *Task:
		if e.AssignedTo != "" {
			encrypted, err := encryptionService.EncryptField(e.AssignedTo)
			if err != nil {
				return fmt.Errorf("failed to encrypt assigned_to field: %w", err)
			}
			e.AssignedTo = encrypted
		}
	}
	return nil
}

func decryptSensitiveFields(entity interface{}, encryptionService *EncryptionService) error {
	switch e := entity.(type) {
	case *User:
		if e.Password != "" {
			decrypted, err := encryptionService.DecryptField(e.Password)
			if err != nil {
				return fmt.Errorf("failed to decrypt password: %w", err)
			}
			e.Password = decrypted
		}
	case *Task:
		if e.AssignedTo != "" {
			decrypted, err := encryptionService.DecryptField(e.AssignedTo)
			if err != nil {
				return fmt.Errorf("failed to decrypt assigned_to field: %w", err)
			}
			e.AssignedTo = decrypted
		}
	}
	return nil
}
