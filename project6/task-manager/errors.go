package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// ErrorCode represents different types of errors
type ErrorCode string

const (
	ErrorCodeValidation    ErrorCode = "VALIDATION_ERROR"
	ErrorCodeNotFound      ErrorCode = "NOT_FOUND"
	ErrorCodeUnauthorized  ErrorCode = "UNAUTHORIZED"
	ErrorCodeForbidden     ErrorCode = "FORBIDDEN"
	ErrorCodeConflict      ErrorCode = "CONFLICT"
	ErrorCodeInternal      ErrorCode = "INTERNAL_ERROR"
	ErrorCodeDatabase      ErrorCode = "DATABASE_ERROR"
	ErrorCodeEncryption    ErrorCode = "ENCRYPTION_ERROR"
	ErrorCodeRateLimit     ErrorCode = "RATE_LIMIT_EXCEEDED"
	ErrorCodeTimeout       ErrorCode = "TIMEOUT"
	ErrorCodeBadRequest    ErrorCode = "BAD_REQUEST"
)

// APIError represents a structured API error
type APIError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	RequestID  string                 `json:"request_id,omitempty"`
	StackTrace string                 `json:"stack_trace,omitempty"`
}

// Error implements the error interface
func (e APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// WithDetails adds details to the error
func (e APIError) WithDetails(key string, value interface{}) APIError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithRequestID adds a request ID to the error
func (e APIError) WithRequestID(requestID string) APIError {
	e.RequestID = requestID
	return e
}

// WithStackTrace adds stack trace to the error
func (e APIError) WithStackTrace() APIError {
	e.StackTrace = getStackTrace()
	return e
}

// HTTPStatus returns the appropriate HTTP status code for the error
func (e APIError) HTTPStatus() int {
	switch e.Code {
	case ErrorCodeValidation, ErrorCodeBadRequest:
		return http.StatusBadRequest
	case ErrorCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrorCodeForbidden:
		return http.StatusForbidden
	case ErrorCodeNotFound:
		return http.StatusNotFound
	case ErrorCodeConflict:
		return http.StatusConflict
	case ErrorCodeRateLimit:
		return http.StatusTooManyRequests
	case ErrorCodeTimeout:
		return http.StatusRequestTimeout
	case ErrorCodeDatabase, ErrorCodeEncryption, ErrorCodeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// ErrorHandler handles errors consistently across the application
type ErrorHandler struct {
	encryptionService *EncryptionService
	logger            *Logger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(encryptionService *EncryptionService, logger *Logger) *ErrorHandler {
	return &ErrorHandler{
		encryptionService: encryptionService,
		logger:            logger,
	}
}

// HandleError processes an error and returns an APIError
func (eh *ErrorHandler) HandleError(err error, requestID string) APIError {
	if apiErr, ok := err.(APIError); ok {
		return apiErr.WithRequestID(requestID)
	}

	// Handle known error types
	switch {
	case strings.Contains(err.Error(), "validation"):
		return eh.NewAPIError(ErrorCodeValidation, "Validation failed", requestID).
			WithDetails("original_error", err.Error())
	case strings.Contains(err.Error(), "not found"):
		return eh.NewAPIError(ErrorCodeNotFound, "Resource not found", requestID).
			WithDetails("original_error", err.Error())
	case strings.Contains(err.Error(), "unauthorized"):
		return eh.NewAPIError(ErrorCodeUnauthorized, "Unauthorized access", requestID).
			WithDetails("original_error", err.Error())
	case strings.Contains(err.Error(), "forbidden"):
		return eh.NewAPIError(ErrorCodeForbidden, "Access forbidden", requestID).
			WithDetails("original_error", err.Error())
	case strings.Contains(err.Error(), "conflict"):
		return eh.NewAPIError(ErrorCodeConflict, "Resource conflict", requestID).
			WithDetails("original_error", err.Error())
	case strings.Contains(err.Error(), "timeout"):
		return eh.NewAPIError(ErrorCodeTimeout, "Operation timeout", requestID).
			WithDetails("original_error", err.Error())
	case strings.Contains(err.Error(), "database"):
		return eh.NewAPIError(ErrorCodeDatabase, "Database operation failed", requestID).
			WithDetails("original_error", err.Error())
	case strings.Contains(err.Error(), "encrypt") || strings.Contains(err.Error(), "decrypt"):
		return eh.NewAPIError(ErrorCodeEncryption, "Encryption/decryption failed", requestID).
			WithDetails("original_error", err.Error())
	default:
		return eh.NewAPIError(ErrorCodeInternal, "Internal server error", requestID).
			WithDetails("original_error", err.Error()).
			WithStackTrace()
	}
}

// NewAPIError creates a new API error
func (eh *ErrorHandler) NewAPIError(code ErrorCode, message string, requestID string) APIError {
	apiErr := APIError{
		Code:      code,
		Message:   message,
		Timestamp: time.Now(),
		RequestID: requestID,
	}

	// Log the error
	eh.logger.Error("API Error", map[string]interface{}{
		"code":       code,
		"message":    message,
		"request_id": requestID,
	})

	return apiErr
}

// NewValidationError creates a validation error
func (eh *ErrorHandler) NewValidationError(validationResult *ValidationResult, requestID string) APIError {
	return eh.NewAPIError(ErrorCodeValidation, "Validation failed", requestID).
		WithDetails("validation_errors", validationResult.Errors)
}

// NewNotFoundError creates a not found error
func (eh *ErrorHandler) NewNotFoundError(resourceType string, resourceID string, requestID string) APIError {
	return eh.NewAPIError(ErrorCodeNotFound, fmt.Sprintf("%s not found", resourceType), requestID).
		WithDetails("resource_type", resourceType).
		WithDetails("resource_id", resourceID)
}

// NewUnauthorizedError creates an unauthorized error
func (eh *ErrorHandler) NewUnauthorizedError(message string, requestID string) APIError {
	if message == "" {
		message = "Unauthorized access"
	}
	return eh.NewAPIError(ErrorCodeUnauthorized, message, requestID)
}

// NewForbiddenError creates a forbidden error
func (eh *ErrorHandler) NewForbiddenError(message string, requestID string) APIError {
	if message == "" {
		message = "Access forbidden"
	}
	return eh.NewAPIError(ErrorCodeForbidden, message, requestID)
}

// NewConflictError creates a conflict error
func (eh *ErrorHandler) NewConflictError(message string, requestID string) APIError {
	if message == "" {
		message = "Resource conflict"
	}
	return eh.NewAPIError(ErrorCodeConflict, message, requestID)
}

// NewInternalError creates an internal error
func (eh *ErrorHandler) NewInternalError(err error, requestID string) APIError {
	return eh.NewAPIError(ErrorCodeInternal, "Internal server error", requestID).
		WithDetails("original_error", err.Error()).
		WithStackTrace()
}

// RecoverFromPanic recovers from panics and converts them to errors
func (eh *ErrorHandler) RecoverFromPanic(requestID string) APIError {
	if r := recover(); r != nil {
		var err error
		switch v := r.(type) {
		case string:
			err = fmt.Errorf("panic: %s", v)
		case error:
			err = v
		default:
			err = fmt.Errorf("panic: %v", v)
		}

		eh.logger.Error("Panic recovered", map[string]interface{}{
			"panic":      r,
			"request_id": requestID,
			"stack_trace": getStackTrace(),
		})

		return eh.NewAPIError(ErrorCodeInternal, "Internal server error", requestID).
			WithDetails("panic", r).
			WithStackTrace()
	}
	return APIError{}
}

// getStackTrace captures the current stack trace
func getStackTrace() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}
