package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// Error implements the error interface
func (v ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", v.Field, v.Message)
}

// ValidationResult contains validation results
type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
}

// Validator provides validation functionality
type Validator struct {
	encryptionService *EncryptionService
}

// NewValidator creates a new validator instance
func NewValidator(encryptionService *EncryptionService) *Validator {
	return &Validator{
		encryptionService: encryptionService,
	}
}

// ValidateEntity validates any entity that implements the Entity interface
func (v *Validator) ValidateEntity(entity Entity) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	switch e := entity.(type) {
	case *Task:
		return v.validateTask(e)
	case *User:
		return v.validateUser(e)
	case *Project:
		return v.validateProject(e)
	default:
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "entity",
			Message: "unknown entity type",
		})
	}

	return result
}

// validateTask validates a task entity
func (v *Validator) validateTask(task *Task) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// Validate title
	if strings.TrimSpace(task.Title) == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "title",
			Message: "title is required",
		})
	} else if len(task.Title) > 100 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "title",
			Message: "title must be less than 100 characters",
			Value:   task.Title,
		})
	}

	// Validate description
	if len(task.Description) > 500 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "description",
			Message: "description must be less than 500 characters",
			Value:   task.Description,
		})
	}

	// Validate priority
	validPriorities := []string{"low", "medium", "high"}
	if task.Priority != "" && !contains(validPriorities, task.Priority) {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "priority",
			Message: "priority must be one of: low, medium, high",
			Value:   task.Priority,
		})
	}

	// Validate due date
	if task.DueDate != "" {
		if _, err := time.Parse("2006-01-02", task.DueDate); err != nil {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "due_date",
				Message: "due_date must be in YYYY-MM-DD format",
				Value:   task.DueDate,
			})
		}
	}

	// Validate assigned to (email format)
	if task.AssignedTo != "" {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(task.AssignedTo) {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "assigned_to",
				Message: "assigned_to must be a valid email address",
				Value:   task.AssignedTo,
			})
		}
	}

	return result
}

// validateUser validates a user entity
func (v *Validator) validateUser(user *User) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// Validate username
	if strings.TrimSpace(user.Username) == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "username",
			Message: "username is required",
		})
	} else if len(user.Username) < 3 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "username",
			Message: "username must be at least 3 characters",
			Value:   user.Username,
		})
	} else if len(user.Username) > 50 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "username",
			Message: "username must be less than 50 characters",
			Value:   user.Username,
		})
	}

	// Validate email
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user.Email) {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "email",
			Message: "email must be a valid email address",
			Value:   user.Email,
		})
	}

	// Validate password
	if len(user.Password) < 8 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "password",
			Message: "password must be at least 8 characters",
		})
	}

	// Validate role
	validRoles := []string{"admin", "user", "guest"}
	if !contains(validRoles, user.Role) {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "role",
			Message: "role must be one of: admin, user, guest",
			Value:   user.Role,
		})
	}

	return result
}

// validateProject validates a project entity
func (v *Validator) validateProject(project *Project) *ValidationResult {
	result := &ValidationResult{
		IsValid: true,
		Errors:  []ValidationError{},
	}

	// Validate name
	if strings.TrimSpace(project.Name) == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: "name is required",
		})
	} else if len(project.Name) > 100 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "name",
			Message: "name must be less than 100 characters",
			Value:   project.Name,
		})
	}

	// Validate description
	if len(project.Description) > 1000 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "description",
			Message: "description must be less than 1000 characters",
			Value:   project.Description,
		})
	}

	// Validate status
	validStatuses := []string{"planning", "active", "completed", "archived"}
	if project.Status != "" && !contains(validStatuses, project.Status) {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "status",
			Message: "status must be one of: planning, active, completed, archived",
			Value:   project.Status,
		})
	}

	// Validate budget
	if project.Budget < 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "budget",
			Message: "budget must be non-negative",
			Value:   fmt.Sprintf("%.2f", project.Budget),
		})
	}

	return result
}

// contains checks if a string slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
