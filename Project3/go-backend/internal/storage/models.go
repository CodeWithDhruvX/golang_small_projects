package storage

import (
	"time"

	"github.com/google/uuid"
)

// Document represents a file uploaded to the knowledge base
type Document struct {
	ID           uuid.UUID      `json:"id" db:"id"`
	Filename     string         `json:"filename" db:"filename"`
	ContentType  string         `json:"content_type" db:"content_type"`
	FileSize     int            `json:"file_size" db:"file_size"`
	UploadTime   time.Time      `json:"upload_time" db:"upload_time"`
	Processed    bool           `json:"processed" db:"processed"`
	Metadata     interface{}   `json:"metadata" db:"metadata"`
}

// DocumentChunk represents a chunk of a document for RAG
type DocumentChunk struct {
	ID         uuid.UUID      `json:"id" db:"id"`
	DocumentID uuid.UUID      `json:"document_id" db:"document_id"`
	ChunkIndex int            `json:"chunk_index" db:"chunk_index"`
	Content    string         `json:"content" db:"content"`
	Embedding  []float32      `json:"embedding" db:"embedding"`
	PageNumber *int           `json:"page_number" db:"page_number"`
	Metadata   interface{}   `json:"metadata" db:"metadata"`
	CreatedAt  time.Time      `json:"created_at" db:"created_at"`
}

// ChatSession represents a conversation session
type ChatSession struct {
	ID           uuid.UUID `json:"id" db:"id"`
	SessionName  *string   `json:"session_name" db:"session_name"`
	ModelName    string    `json:"model_name" db:"model_name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	LastActivity time.Time `json:"last_activity" db:"last_activity"`
}

// ChatMessage represents a message in a chat session
type ChatMessage struct {
	ID           uuid.UUID   `json:"id" db:"id"`
	SessionID    uuid.UUID   `json:"session_id" db:"session_id"`
	MessageType  string      `json:"message_type" db:"message_type"` // "user" or "assistant"
	Content      string      `json:"content" db:"content"`
	Citations    interface{}   `json:"citations" db:"citations"`
	ModelUsed    *string     `json:"model_used" db:"model_used"`
	TokensUsed   *int        `json:"tokens_used" db:"tokens_used"`
	ResponseTime *int        `json:"response_time_ms" db:"response_time_ms"`
	CreatedAt    time.Time   `json:"created_at" db:"created_at"`
}

// SearchResult represents a document chunk search result with similarity score
type SearchResult struct {
	DocumentChunk
	Similarity float64 `json:"similarity"`
	Document   Document `json:"document"`
}

// ChatRequest represents a request to chat with the AI
type ChatRequest struct {
	SessionID uuid.UUID `json:"session_id" binding:"required"`
	Message   string    `json:"message" binding:"required"`
	Model     string    `json:"model,omitempty"`
}

// ChatResponse represents a response from the AI
type ChatResponse struct {
	SessionID   uuid.UUID            `json:"session_id"`
	MessageID   uuid.UUID            `json:"message_id"`
	Message     string               `json:"message"`
	Citations   []DocumentCitation   `json:"citations"`
	ModelUsed   string               `json:"model_used"`
	TokensUsed  int                  `json:"tokens_used"`
	ResponseTime int                  `json:"response_time_ms"`
	Timestamp   time.Time            `json:"timestamp"`
}

// DocumentCitation represents a citation from a document
type DocumentCitation struct {
	DocumentID   uuid.UUID `json:"document_id"`
	Filename     string    `json:"filename"`
	PageNumber   *int      `json:"page_number"`
	ChunkIndex   int       `json:"chunk_index"`
	Content      string    `json:"content"`
	Similarity   float64   `json:"similarity"`
}

// UploadResponse represents the response after uploading a document
type UploadResponse struct {
	DocumentID uuid.UUID `json:"document_id"`
	Filename   string    `json:"filename"`
	Status     string    `json:"status"`
	Message    string    `json:"message"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
	Ollama    string    `json:"ollama"`
}

// ReadyResponse represents the readiness check response
type ReadyResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

// MetricsResponse represents application metrics
type MetricsResponse struct {
	TotalDocuments      int `json:"total_documents"`
	ProcessedDocuments  int `json:"processed_documents"`
	TotalChunks         int `json:"total_chunks"`
	TotalSessions       int `json:"total_sessions"`
	TotalMessages       int `json:"total_messages"`
	AverageResponseTime int `json:"average_response_time_ms"`
}
