package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// New creates a new logger instance
func New(level string) *logrus.Logger {
	logger := logrus.New()

	// Set output format
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	// Set output destination
	logger.SetOutput(os.Stdout)

	// Set log level
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

// NewDevelopment creates a logger for development environment
func NewDevelopment() *logrus.Logger {
	logger := logrus.New()

	// Use text formatter for development
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(os.Stdout)

	return logger
}

// NewProduction creates a logger for production environment
func NewProduction() *logrus.Logger {
	logger := logrus.New()

	// Use JSON formatter for production
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
	})

	logger.SetLevel(logrus.InfoLevel)
	logger.SetOutput(os.Stdout)

	return logger
}
