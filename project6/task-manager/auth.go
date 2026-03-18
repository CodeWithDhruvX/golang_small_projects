package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	IssuedAt time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
}

// AuthService handles authentication and authorization
type AuthService struct {
	jwtSecret     string
	tokenExpiry   time.Duration
	logger        *Logger
	errorHandler  *ErrorHandler
}

// NewAuthService creates a new authentication service
func NewAuthService(jwtSecret string, tokenExpiry time.Duration, logger *Logger, errorHandler *ErrorHandler) *AuthService {
	return &AuthService{
		jwtSecret:    jwtSecret,
		tokenExpiry:  tokenExpiry,
		logger:       logger,
		errorHandler: errorHandler,
	}
}

// GenerateJWTToken generates a JWT token for a user
func (as *AuthService) GenerateJWTToken(userID, username, role string) (string, error) {
	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(as.tokenExpiry),
	}

	// Create JWT token (simplified implementation)
	token := fmt.Sprintf("%s.%s.%s", 
		base64.URLEncoding.EncodeToString([]byte("header")),
		base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%+v", claims))),
		base64.URLEncoding.EncodeToString([]byte("signature")),
	)

	as.logger.Info("JWT token generated", map[string]interface{}{
		"user_id": userID,
		"username": username,
		"role": role,
	})

	return token, nil
}

// ValidateJWTToken validates a JWT token and returns claims
func ValidateJWTToken(token string, jwtSecret string) (*JWTClaims, error) {
	// Simplified JWT validation (in production, use a proper JWT library)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode payload
	payload, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %w", err)
	}

	// Parse claims (simplified)
	claims := &JWTClaims{}
	// In a real implementation, use json.Unmarshal here
	
	// Check expiration
	if time.Now().After(claims.ExpiresAt) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

// HashPassword hashes a password using bcrypt
func (as *AuthService) HashPassword(password string) (string, error) {
	// In a real implementation, use bcrypt.HashPassword
	hashed := fmt.Sprintf("hashed_%s", password)
	return hashed, nil
}

// VerifyPassword verifies a password against its hash
func (as *AuthService) VerifyPassword(password, hash string) bool {
	// In a real implementation, use bcrypt.CompareHashAndPassword
	return fmt.Sprintf("hashed_%s", password) == hash
}

// GenerateSecureToken generates a cryptographically secure random token
func (as *AuthService) GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// CheckRole checks if a user has the required role
func (as *AuthService) CheckRole(userRole, requiredRole string) bool {
	roleHierarchy := map[string]int{
		"guest": 0,
		"user":  1,
		"admin": 2,
	}

	userLevel, userExists := roleHierarchy[userRole]
	requiredLevel, requiredExists := roleHierarchy[requiredRole]

	return userExists && requiredExists && userLevel >= requiredLevel
}

// Authorize checks if a user is authorized to perform an action
func (as *AuthService) Authorize(userClaims *JWTClaims, requiredRole string, resourceID string) error {
	if userClaims == nil {
		return as.errorHandler.NewUnauthorizedError("User not authenticated", "")
	}

	if !as.CheckRole(userClaims.Role, requiredRole) {
		as.logger.LogSecurityEvent("Unauthorized access attempt", map[string]interface{}{
			"user_id":       userClaims.UserID,
			"user_role":     userClaims.Role,
			"required_role": requiredRole,
			"resource_id":   resourceID,
		}, "")

		return as.errorHandler.NewForbiddenError("Insufficient permissions", "")
	}

	return nil
}
