package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// LogLevel represents different log levels
type LogLevel string

const (
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelFatal LogLevel = "FATAL"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Data      map[string]interface{} `json:"data,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
}

// Logger provides structured logging functionality
type Logger struct {
	level  LogLevel
	output *log.Logger
	file   *os.File
}

// NewLogger creates a new logger instance
func NewLogger(level LogLevel, logToFile bool) (*Logger, error) {
	logger := &Logger{
		level:  level,
		output: log.New(os.Stdout, "", 0),
	}

	if logToFile {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.file = logFile
		logger.output = log.New(logFile, "", 0)
	}

	return logger, nil
}

// Close closes the logger and any open files
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// log is the internal logging method
func (l *Logger) log(level LogLevel, message string, data map[string]interface{}, requestID string) {
	if !l.shouldLog(level) {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Data:      data,
		RequestID: requestID,
	}

	jsonEntry, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple logging if JSON marshaling fails
		l.output.Printf("[%s] %s - %s", level, message, err.Error())
		return
	}

	l.output.Println(string(jsonEntry))
}

// shouldLog checks if the given level should be logged
func (l *Logger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		LogLevelDebug: 0,
		LogLevelInfo:  1,
		LogLevelWarn:  2,
		LogLevelError: 3,
		LogLevelFatal: 4,
	}

	return levels[level] >= levels[l.level]
}

// Debug logs a debug message
func (l *Logger) Debug(message string, data map[string]interface{}) {
	l.log(LogLevelDebug, message, data, "")
}

// Info logs an info message
func (l *Logger) Info(message string, data map[string]interface{}) {
	l.log(LogLevelInfo, message, data, "")
}

// Warn logs a warning message
func (l *Logger) Warn(message string, data map[string]interface{}) {
	l.log(LogLevelWarn, message, data, "")
}

// Error logs an error message
func (l *Logger) Error(message string, data map[string]interface{}) {
	l.log(LogLevelError, message, data, "")
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(message string, data map[string]interface{}) {
	l.log(LogLevelFatal, message, data, "")
	l.Close()
	os.Exit(1)
}

// DebugWithRequestID logs a debug message with request ID
func (l *Logger) DebugWithRequestID(message string, data map[string]interface{}, requestID string) {
	l.log(LogLevelDebug, message, data, requestID)
}

// InfoWithRequestID logs an info message with request ID
func (l *Logger) InfoWithRequestID(message string, data map[string]interface{}, requestID string) {
	l.log(LogLevelInfo, message, data, requestID)
}

// WarnWithRequestID logs a warning message with request ID
func (l *Logger) WarnWithRequestID(message string, data map[string]interface{}, requestID string) {
	l.log(LogLevelWarn, message, data, requestID)
}

// ErrorWithRequestID logs an error message with request ID
func (l *Logger) ErrorWithRequestID(message string, data map[string]interface{}, requestID string) {
	l.log(LogLevelError, message, data, requestID)
}

// LogRequest logs an HTTP request
func (l *Logger) LogRequest(method, path, userAgent, remoteAddr string, requestID string) {
	l.InfoWithRequestID("HTTP Request", map[string]interface{}{
		"method":      method,
		"path":        path,
		"user_agent":  userAgent,
		"remote_addr": remoteAddr,
	}, requestID)
}

// LogResponse logs an HTTP response
func (l *Logger) LogResponse(statusCode int, duration time.Duration, requestID string) {
	level := LogLevelInfo
	if statusCode >= 400 {
		level = LogLevelWarn
	}
	if statusCode >= 500 {
		level = LogLevelError
	}

	l.log(level, "HTTP Response", map[string]interface{}{
		"status_code": statusCode,
		"duration_ms": duration.Milliseconds(),
	}, requestID)
}

// LogSecurityEvent logs security-related events
func (l *Logger) LogSecurityEvent(event string, data map[string]interface{}, requestID string) {
	l.ErrorWithRequestID(fmt.Sprintf("Security Event: %s", event), data, requestID)
}

// LogAPIError logs API errors
func (l *Logger) LogAPIError(apiErr APIError, requestID string) {
	l.ErrorWithRequestID("API Error", map[string]interface{}{
		"error_code": apiErr.Code,
		"message":    apiErr.Message,
		"details":    apiErr.Details,
	}, requestID)
}

// LogPanic logs a panic event
func (l *Logger) LogPanic(panicValue interface{}, requestID string) {
	l.ErrorWithRequestID("Panic occurred", map[string]interface{}{
		"panic": panicValue,
	}, requestID)
}
