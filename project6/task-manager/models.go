package main

import (
	"time"
)

// Entity interface defines common behavior for all entities (polymorphism)
type Entity interface {
	GetID() interface{}
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	Validate() error
	EncryptSensitiveData(key []byte) error
	DecryptSensitiveData(key []byte) error
}

// BaseModel provides common fields for all entities
type BaseModel struct {
	ID        interface{} `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func (b BaseModel) GetID() interface{} {
	return b.ID
}

func (b BaseModel) GetCreatedAt() time.Time {
	return b.CreatedAt
}

func (b BaseModel) GetUpdatedAt() time.Time {
	return b.UpdatedAt
}

// Task represents a task entity
type Task struct {
	BaseModel
	Title       string `json:"title" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
	Completed   bool   `json:"completed"`
	Priority    string `json:"priority" validate:"oneof=low medium high"`
	DueDate     string `json:"due_date,omitempty" validate:"omitempty,datetime=2006-01-02"`
	AssignedTo  string `json:"assigned_to,omitempty" validate:"omitempty,email"`
}

// User represents a user entity (demonstrating polymorphism)
type User struct {
	BaseModel
	Username    string `json:"username" validate:"required,min=3,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"-" validate:"required,min=8"` // Hidden in JSON
	Role        string `json:"role" validate:"oneof=admin user guest"`
	IsActive    bool   `json:"is_active"`
	LastLogin   string `json:"last_login,omitempty"`
}

// Project represents a project entity (demonstrating polymorphism)
type Project struct {
	BaseModel
	Name        string   `json:"name" validate:"required,min=1,max=100"`
	Description string   `json:"description,omitempty" validate:"max=1000"`
	Status      string   `json:"status" validate:"oneof=planning active completed archived"`
	TeamMembers []string `json:"team_members,omitempty"`
	Budget      float64  `json:"budget,omitempty"`
}

// Repository interface for data access operations (polymorphism)
type Repository[T Entity] interface {
	Create(entity T) (T, error)
	GetByID(id interface{}) (T, error)
	GetAll() ([]T, error)
	Update(entity T) (T, error)
	Delete(id interface{}) error
	Search(criteria map[string]interface{}) ([]T, error)
}

// Service interface for business logic (polymorphism)
type Service[T Entity] interface {
	Create(entity T) (T, error)
	GetByID(id interface{}) (T, error)
	GetAll() ([]T, error)
	Update(entity T) (T, error)
	Delete(id interface{}) error
	ValidateAndCreate(entity T) (T, error)
}
