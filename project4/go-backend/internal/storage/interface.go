package storage

import "ai-recruiter-assistant/internal/auth"

// StorageInterface defines all storage operations
type StorageInterface interface {
	// User operations
	CreateUser(user *auth.User) error
	GetUserByEmail(email string) (*auth.User, error)
	GetUserByID(id string) (*auth.User, error)
	UpdateUser(user *auth.User) error

	// Email operations
	CreateEmail(email *Email) (*Email, error)
	GetEmailByID(id, userID string) (*Email, error)
	GetUserEmails(userID string, page, limit int) ([]Email, error)
	UpdateEmail(email *Email) error
	DeleteEmail(id, userID string) error

	// Application operations
	CreateApplication(application *Application) (*Application, error)
	GetApplicationByID(id, userID string) (*Application, error)
	GetUserApplications(userID string, page, limit int, status string) ([]Application, error)
	UpdateApplication(application *Application) (*Application, error)
	DeleteApplication(id, userID string) error
	CheckDuplicateApplication(userID, company, recruiterEmail string) (bool, error)

	// Document operations
	CreateDocument(document *Document) error
	GetDocumentsByUserID(userID string, source string) ([]Document, error)
	SearchDocuments(userID, query string, topK int) ([]Document, error)
	DeleteDocument(id, userID string) error

	// AI Reply operations
	CreateAIReply(reply *AIReply) error
	GetAIRepliesByUserID(userID string) ([]AIReply, error)
	GetAIRepliesByEmailID(emailID, userID string) ([]AIReply, error)
	UpdateAIReply(reply *AIReply) error

	// Vector operations
	StoreEmbedding(id string, embedding []float64, table string) error
	SearchSimilar(embedding []float64, userID string, topK int, table string) ([]Document, error)

	// Gmail integration operations
	CreateGmailIntegration(integration *GmailIntegration) error
	GetGmailIntegration(userID string) (*GmailIntegration, error)
	UpdateGmailIntegration(integration *GmailIntegration) error
	DeleteGmailIntegration(userID string) error
}
