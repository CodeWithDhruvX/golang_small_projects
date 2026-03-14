package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

// LogoutHandler handles user logout by invalidating the token
func LogoutHandler(redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No authorization header provided"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bearer token required"})
			return
		}

		// Add token to blacklist in Redis (optional - for additional security)
		if redisClient != nil {
			ctx := c.Request.Context()
			err := redisClient.Set(ctx, "blacklist:"+tokenString, "true", 0).Err()
			if err != nil {
				logrus.Warnf("Failed to blacklist token: %v", err)
			}
		}

		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	}
}

// ErrorHandler is a custom error handler middleware
func ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			logrus.Errorf("Panic recovered: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
				"message": "Something went wrong. Please try again later.",
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}

// RequestLogger logs all incoming requests
func RequestLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"user_agent": param.Request.UserAgent(),
		}).Info("HTTP Request")
		return ""
	})
}

// RateLimiter implements a simple rate limiter middleware
func RateLimiter() gin.HandlerFunc {
	// This is a placeholder implementation
	// In a production environment, you'd want to use a more sophisticated rate limiting solution
	return func(c *gin.Context) {
		// For now, just pass through
		c.Next()
	}
}
