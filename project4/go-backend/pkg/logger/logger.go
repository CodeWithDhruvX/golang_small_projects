package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	// Logger instance
	Logger *logrus.Logger
)

// InitLogger initializes the global logger
func InitLogger() {
	Logger = logrus.New()

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set formatter based on environment
	if os.Getenv("ENVIRONMENT") == "production" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set log level
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		level = "info"
	}

	switch level {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "info":
		Logger.SetLevel(logrus.InfoLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}
