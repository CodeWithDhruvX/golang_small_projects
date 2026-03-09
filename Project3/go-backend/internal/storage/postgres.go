package storage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// PostgresDB handles all database operations
type PostgresDB struct {
	pool   *pgxpool.Pool
	logger *logrus.Logger
}

// NewPostgresDB creates a new database connection pool
func NewPostgresDB(databaseURL string, logger *logrus.Logger) (*PostgresDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database")

	return &PostgresDB{
		pool:   pool,
		logger: logger,
	}, nil
}

// float32ArrayToVectorString converts []float32 to PostgreSQL vector format string
func (db *PostgresDB) float32ArrayToVectorString(embedding []float32) string {
	if embedding == nil {
		return "[]"
	}
	
	strValues := make([]string, len(embedding))
	for i, val := range embedding {
		strValues[i] = strconv.FormatFloat(float64(val), 'f', -1, 64)
	}
	
	return "[" + strings.Join(strValues, ",") + "]"
}

// vectorStringToFloat32Array converts PostgreSQL vector format string to []float32
func (db *PostgresDB) vectorStringToFloat32Array(vectorStr string) []float32 {
	if vectorStr == "" || vectorStr == "[]" {
		return nil
	}
	
	// Remove brackets and split by comma
	vectorStr = strings.Trim(vectorStr, "[]")
	if vectorStr == "" {
		return nil
	}
	
	strValues := strings.Split(vectorStr, ",")
	embedding := make([]float32, len(strValues))
	
	for i, strVal := range strValues {
		if val, err := strconv.ParseFloat(strings.TrimSpace(strVal), 32); err == nil {
			embedding[i] = float32(val)
		}
	}
	
	return embedding
}

// Close closes the database connection pool
func (db *PostgresDB) Close() error {
	db.pool.Close()
	return nil
}

// Ping tests the database connection
func (db *PostgresDB) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return db.pool.Ping(ctx)
}

// RunMigrations runs database migrations
func (db *PostgresDB) RunMigrations() error {
	db.logger.Info("Running database migrations...")
	
	// The migrations are handled by the init script in docker-config/postgres/init-pgvector.sql
	// This method can be used for future migrations
	
	db.logger.Info("Database migrations completed")
	return nil
}

// Document operations

// CreateDocument creates a new document record
func (db *PostgresDB) CreateDocument(ctx context.Context, doc *Document) error {
	query := `
		INSERT INTO documents (id, filename, content_type, file_size, processed, metadata)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := db.pool.Exec(ctx, query, doc.ID, doc.Filename, doc.ContentType, doc.FileSize, doc.Processed, doc.Metadata)
	if err != nil {
		return fmt.Errorf("failed to create document: %w", err)
	}
	
	return nil
}

// GetDocument retrieves a document by ID
func (db *PostgresDB) GetDocument(ctx context.Context, id uuid.UUID) (*Document, error) {
	query := `
		SELECT id, filename, content_type, file_size, upload_time, processed, metadata
		FROM documents
		WHERE id = $1
	`
	
	var doc Document
	err := db.pool.QueryRow(ctx, query, id).Scan(
		&doc.ID, &doc.Filename, &doc.ContentType, &doc.FileSize,
		&doc.UploadTime, &doc.Processed, &doc.Metadata,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("document not found")
		}
		return nil, fmt.Errorf("failed to get document: %w", err)
	}
	
	return &doc, nil
}

// ListDocuments retrieves all documents with pagination
func (db *PostgresDB) ListDocuments(ctx context.Context, limit, offset int) ([]Document, error) {
	query := `
		SELECT id, filename, content_type, file_size, upload_time, processed, metadata
		FROM documents
		ORDER BY upload_time DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := db.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %w", err)
	}
	defer rows.Close()
	
	var documents []Document
	for rows.Next() {
		var doc Document
		if err := rows.Scan(
			&doc.ID, &doc.Filename, &doc.ContentType, &doc.FileSize,
			&doc.UploadTime, &doc.Processed, &doc.Metadata,
		); err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		documents = append(documents, doc)
	}
	
	return documents, nil
}

// UpdateDocumentProcessed updates the processed status of a document
func (db *PostgresDB) UpdateDocumentProcessed(ctx context.Context, id uuid.UUID, processed bool) error {
	query := `UPDATE documents SET processed = $1 WHERE id = $2`
	
	_, err := db.pool.Exec(ctx, query, processed, id)
	if err != nil {
		return fmt.Errorf("failed to update document processed status: %w", err)
	}
	
	return nil
}

// DeleteDocument deletes a document and its chunks
func (db *PostgresDB) DeleteDocument(ctx context.Context, id uuid.UUID) error {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Delete chunks first (foreign key constraint)
	_, err = tx.Exec(ctx, "DELETE FROM document_chunks WHERE document_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete document chunks: %w", err)
	}
	
	// Delete document
	_, err = tx.Exec(ctx, "DELETE FROM documents WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

// DocumentChunk operations

// CreateDocumentChunk creates a new document chunk
func (db *PostgresDB) CreateDocumentChunk(ctx context.Context, chunk *DocumentChunk) error {
	query := `
		INSERT INTO document_chunks (id, document_id, chunk_index, content, embedding, page_number, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	// Convert embedding to PostgreSQL vector format
	var embedding interface{}
	if chunk.Embedding != nil {
		embedding = db.float32ArrayToVectorString(chunk.Embedding)
	}
	
	_, err := db.pool.Exec(ctx, query, 
		chunk.ID, chunk.DocumentID, chunk.ChunkIndex, chunk.Content, 
		embedding, chunk.PageNumber, chunk.Metadata)
	if err != nil {
		return fmt.Errorf("failed to create document chunk: %w", err)
	}
	
	return nil
}

// GetDocumentChunks retrieves chunks for a specific document
func (db *PostgresDB) GetDocumentChunks(ctx context.Context, documentID uuid.UUID) ([]DocumentChunk, error) {
	query := `
		SELECT id, document_id, chunk_index, content, page_number, metadata, created_at
		FROM document_chunks
		WHERE document_id = $1
		ORDER BY chunk_index ASC
	`
	rows, err := db.pool.Query(ctx, query, documentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document chunks: %w", err)
	}
	defer rows.Close()

	var chunks []DocumentChunk
	for rows.Next() {
		var chunk DocumentChunk
		var metadata interface{}
		if err := rows.Scan(
			&chunk.ID, &chunk.DocumentID, &chunk.ChunkIndex, &chunk.Content,
			&chunk.PageNumber, &metadata, &chunk.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan document chunk: %w", err)
		}
		chunk.Metadata = metadata
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

// SearchSimilarChunks searches for similar chunks using vector similarity
func (db *PostgresDB) SearchSimilarChunks(ctx context.Context, embedding []float32, limit int) ([]SearchResult, error) {
	query := `
		SELECT 
			dc.id, dc.document_id, dc.chunk_index, dc.content, dc.embedding, 
			dc.page_number, dc.metadata, dc.created_at,
			1 - (dc.embedding <=> $1) as similarity,
			d.id, d.filename, d.content_type, d.file_size, d.upload_time, d.processed, d.metadata
		FROM document_chunks dc
		JOIN documents d ON dc.document_id = d.id
		ORDER BY dc.embedding <=> $1
		LIMIT $2
	`
	
	// Convert embedding to PostgreSQL vector format
	embeddingStr := db.float32ArrayToVectorString(embedding)
	
	rows, err := db.pool.Query(ctx, query, embeddingStr, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search similar chunks: %w", err)
	}
	defer rows.Close()
	
	var results []SearchResult
	for rows.Next() {
		var result SearchResult
		var docMetadata, chunkMetadata []byte
		var similarity float64
		var embeddingStr string
		
		err := rows.Scan(
			&result.ID, &result.DocumentID, &result.ChunkIndex, &result.Content, 
			&embeddingStr, &result.PageNumber, &chunkMetadata, &result.CreatedAt,
			&similarity,
			&result.Document.ID, &result.Document.Filename, &result.Document.ContentType, 
			&result.Document.FileSize, &result.Document.UploadTime, &result.Document.Processed, 
			&docMetadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}
		
		// Convert embedding string back to []float32 if needed
		if embeddingStr != "" && embeddingStr != "[]" {
			result.Embedding = db.vectorStringToFloat32Array(embeddingStr)
		}
		
		result.Similarity = similarity
		result.Document.Metadata = docMetadata
		result.Metadata = chunkMetadata
		
		results = append(results, result)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating search results: %w", err)
	}
	
	return results, nil
}

// CreateChatSession creates a new chat session
func (db *PostgresDB) CreateChatSession(ctx context.Context, session *ChatSession) error {
	query := `
		INSERT INTO chat_sessions (id, session_name, model_name, created_at, last_activity)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := db.pool.Exec(ctx, query, 
		session.ID, session.SessionName, session.ModelName, 
		session.CreatedAt, session.LastActivity)
	if err != nil {
		return fmt.Errorf("failed to create chat session: %w", err)
	}
	
	return nil
}

// GetChatSession retrieves a chat session by ID
func (db *PostgresDB) GetChatSession(ctx context.Context, id uuid.UUID) (*ChatSession, error) {
	query := `
		SELECT id, session_name, model_name, created_at, last_activity
		FROM chat_sessions
		WHERE id = $1
	`
	
	var session ChatSession
	err := db.pool.QueryRow(ctx, query, id).Scan(
		&session.ID, &session.SessionName, &session.ModelName,
		&session.CreatedAt, &session.LastActivity,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("chat session not found")
		}
		return nil, fmt.Errorf("failed to get chat session: %w", err)
	}
	
	return &session, nil
}

// ListChatSessions retrieves all chat sessions
func (db *PostgresDB) ListChatSessions(ctx context.Context) ([]ChatSession, error) {
	query := `
		SELECT id, session_name, model_name, created_at, last_activity
		FROM chat_sessions
		ORDER BY last_activity DESC
	`
	
	rows, err := db.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list chat sessions: %w", err)
	}
	defer rows.Close()
	
	var sessions []ChatSession
	for rows.Next() {
		var session ChatSession
		if err := rows.Scan(
			&session.ID, &session.SessionName, &session.ModelName,
			&session.CreatedAt, &session.LastActivity,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chat session: %w", err)
		}
		sessions = append(sessions, session)
	}
	
	return sessions, nil
}

// ChatMessage operations

// CreateChatMessage creates a new chat message
func (db *PostgresDB) CreateChatMessage(ctx context.Context, message *ChatMessage) error {
	query := `
		INSERT INTO chat_messages (id, session_id, message_type, content, citations, model_used, tokens_used, response_time_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err := db.pool.Exec(ctx, query,
		message.ID, message.SessionID, message.MessageType, message.Content,
		message.Citations, message.ModelUsed, message.TokensUsed, message.ResponseTime, message.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create chat message: %w", err)
	}
	
	return nil
}

// GetChatMessages retrieves messages for a session
func (db *PostgresDB) GetChatMessages(ctx context.Context, sessionID uuid.UUID, limit int) ([]ChatMessage, error) {
	query := `
		SELECT id, session_id, message_type, content, citations, model_used, tokens_used, response_time_ms, created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
		LIMIT $2
	`
	
	rows, err := db.pool.Query(ctx, query, sessionID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}
	defer rows.Close()
	
	var messages []ChatMessage
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(
			&message.ID, &message.SessionID, &message.MessageType, &message.Content,
			&message.Citations, &message.ModelUsed, &message.TokensUsed, &message.ResponseTime, &message.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan chat message: %w", err)
		}
		messages = append(messages, message)
	}
	
	return messages, nil
}

// DeleteChatSession deletes a chat session and cascades messages
func (db *PostgresDB) DeleteChatSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM chat_sessions WHERE id = $1`
	_, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete chat session: %w", err)
	}
	return nil
}

// Metrics operations

// GetMetrics retrieves application metrics
func (db *PostgresDB) GetMetrics(ctx context.Context) (*MetricsResponse, error) {
	query := `
		SELECT 
			(SELECT COUNT(*) FROM documents) as total_documents,
			(SELECT COUNT(*) FROM documents WHERE processed = true) as processed_documents,
			(SELECT COUNT(*) FROM document_chunks) as total_chunks,
			(SELECT COUNT(*) FROM chat_sessions) as total_sessions,
			(SELECT COUNT(*) FROM chat_messages) as total_messages,
			COALESCE(AVG(response_time_ms), 0) as average_response_time
	`
	
	var metrics MetricsResponse
	err := db.pool.QueryRow(ctx, query).Scan(
		&metrics.TotalDocuments, &metrics.ProcessedDocuments, &metrics.TotalChunks,
		&metrics.TotalSessions, &metrics.TotalMessages, &metrics.AverageResponseTime,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get metrics: %w", err)
	}
	
	return &metrics, nil
}
