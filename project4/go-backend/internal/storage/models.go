
package storage

import "time"

// Email represents an email in the system
type Email struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	SenderEmail  string    `json:"sender_email"`
	SenderName   string    `json:"sender_name,omitempty"`
	IsRecruiter  bool      `json:"is_recruiter"`
	Processed    bool      `json:"processed"`
	GmailID      string    `json:"gmail_id,omitempty"`
	Embedding    []float64 `json:"embedding,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Application represents a job application
type Application struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Company        string    `json:"company"`
	Role           string    `json:"role"`
	RecruiterEmail string    `json:"recruiter_email"`
	RecruiterName  string    `json:"recruiter_name,omitempty"`
	Status         string    `json:"status"` // Applied, Interview Scheduled, Offer, Rejected
	EmailID        string    `json:"email_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Document represents a document chunk for RAG
type Document struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	Embedding []float64 `json:"embedding,omitempty"`
	Source    string    `json:"source"` // resume, profile, email, etc.
	Metadata  string    `json:"metadata,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// AIReply represents an AI-generated reply
type AIReply struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	EmailID       string    `json:"email_id"`
	ReplyContent  string    `json:"reply_content"`
	ModelUsed     string    `json:"model_used"`
	TokensUsed    int       `json:"tokens_used"`
	ResponseTime  int       `json:"response_time_ms"`
	IsSent        bool      `json:"is_sent"`
	CreatedAt     time.Time `json:"created_at"`
}

// EmailProcessingLog represents a log entry for email processing
type EmailProcessingLog struct {
	ID             string    `json:"id"`
	EmailID        string    `json:"email_id,omitempty"`
	ProcessingStep string    `json:"processing_step"`
	Status         string    `json:"status"` // success, error, warning
	Message        string    `json:"message,omitempty"`
	Metadata       string    `json:"metadata,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// GmailIntegration represents Gmail OAuth2 integration for a user
type GmailIntegration struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenExpiry  time.Time `json:"token_expiry"`
	Email        string    `json:"email"`
	IsActive     bool      `json:"is_active"`
	LastSyncAt   time.Time `json:"last_sync_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
