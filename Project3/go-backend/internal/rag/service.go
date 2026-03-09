package rag

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"private-knowledge-base-go/internal/storage"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Service handles Retrieval-Augmented Generation operations
type Service struct {
	db         *storage.PostgresDB
	ollamaURL  string
	logger     *logrus.Logger
	aiClient   *OllamaClient
}

// NewService creates a new RAG service
func NewService(db *storage.PostgresDB, ollamaURL string, logger *logrus.Logger) *Service {
	return &Service{
		db:        db,
		ollamaURL: ollamaURL,
		logger:    logger,
		aiClient:  NewOllamaClient(ollamaURL, logger),
	}
}

// Chat handles a chat request with RAG
func (s *Service) Chat(ctx context.Context, req *storage.ChatRequest) (*storage.ChatResponse, error) {
	startTime := time.Now()

	// Ensure session exists
	session, err := s.getOrCreateSession(ctx, req.SessionID, req.Model)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create session: %w", err)
	}

	// Store user message
	userMessageID := uuid.New()
	userMessage := &storage.ChatMessage{
		ID:          userMessageID,
		SessionID:    req.SessionID,
		MessageType:  "user",
		Content:      req.Message,
		CreatedAt:    time.Now(),
	}

	if err := s.db.CreateChatMessage(ctx, userMessage); err != nil {
		return nil, fmt.Errorf("failed to store user message: %w", err)
	}

	// Generate embedding for user query
	queryEmbedding, err := s.aiClient.GenerateEmbedding(ctx, req.Message)
	if err != nil {
		s.logger.Warnf("Failed to generate embedding for query: %v", err)
		// Continue without embedding for now
		queryEmbedding = make([]float32, 768) // Placeholder
	}

	// Retrieve relevant document chunks
	searchResults, err := s.db.SearchSimilarChunks(ctx, queryEmbedding, 5) // Retrieve top 5 chunks
	if err != nil {
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	// Build context from retrieved chunks
	context := s.buildContext(searchResults)

	// Generate AI response
	response, err := s.generateResponse(ctx, req.Message, context, session.ModelName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AI response: %w", err)
	}

	// Store assistant message
	assistantMessageID := uuid.New()
	citations := s.buildCitations(searchResults)
	responseTime := int(time.Since(startTime).Milliseconds())

	assistantMessage := &storage.ChatMessage{
		ID:            assistantMessageID,
		SessionID:     req.SessionID,
		MessageType:   "assistant",
		Content:       response.Content,
		Citations:     s.citationsToJSON(citations),
		ModelUsed:     &response.ModelUsed,
		TokensUsed:    &response.TokensUsed,
		ResponseTime:  &responseTime,
		CreatedAt:     time.Now(),
	}

	if err := s.db.CreateChatMessage(ctx, assistantMessage); err != nil {
		return nil, fmt.Errorf("failed to store assistant message: %w", err)
	}

	return &storage.ChatResponse{
		SessionID:    req.SessionID,
		MessageID:    assistantMessageID,
		Message:      response.Content,
		Citations:    citations,
		ModelUsed:    response.ModelUsed,
		TokensUsed:   response.TokensUsed,
		ResponseTime: responseTime,
		Timestamp:    time.Now(),
	}, nil
}

// StreamChat handles streaming chat responses
func (s *Service) StreamChat(ctx context.Context, req *storage.ChatRequest) (<-chan StreamChunk, error) {
	chunkChan := make(chan StreamChunk, 10)

	go func() {
		defer close(chunkChan)

		startTime := time.Now()

		// Send initial chunk
		chunkChan <- StreamChunk{
			Type:      "start",
			SessionID: req.SessionID,
			Timestamp: time.Now(),
		}

		// Ensure session exists
		session, err := s.getOrCreateSession(ctx, req.SessionID, req.Model)
		if err != nil {
			chunkChan <- StreamChunk{
				Type:    "error",
				Message: fmt.Sprintf("Failed to get/create session: %v", err),
			}
			return
		}

		// Store user message
		userMessageID := uuid.New()
		userMessage := &storage.ChatMessage{
			ID:          userMessageID,
			SessionID:    req.SessionID,
			MessageType:  "user",
			Content:      req.Message,
			CreatedAt:    time.Now(),
		}

		if err := s.db.CreateChatMessage(ctx, userMessage); err != nil {
			chunkChan <- StreamChunk{
				Type:    "error",
				Message: fmt.Sprintf("Failed to store user message: %v", err),
			}
			return
		}

		// Generate embedding for user query
		queryEmbedding, err := s.aiClient.GenerateEmbedding(ctx, req.Message)
		if err != nil {
			s.logger.Warnf("Failed to generate embedding for query: %v", err)
			queryEmbedding = make([]float32, 768) // Placeholder
		}

		// Retrieve relevant document chunks
		searchResults, err := s.db.SearchSimilarChunks(ctx, queryEmbedding, 5)
		if err != nil {
			chunkChan <- StreamChunk{
				Type:    "error",
				Message: fmt.Sprintf("Failed to search documents: %v", err),
			}
			return
		}

		// Build context from retrieved chunks
		context := s.buildContext(searchResults)

		// Stream AI response
		responseChan, err := s.aiClient.StreamResponse(ctx, req.Message, context, session.ModelName)
		if err != nil {
			chunkChan <- StreamChunk{
				Type:    "error",
				Message: fmt.Sprintf("Failed to generate AI response: %v", err),
			}
			return
		}

		var fullResponse strings.Builder
		var tokensUsed int
		var modelUsed string

		// Stream response chunks
		for chunk := range responseChan {
			if chunk.Error != "" {
				chunkChan <- StreamChunk{
					Type:    "error",
					Message: chunk.Error,
				}
				return
			}

			fullResponse.WriteString(chunk.Content)
			tokensUsed += chunk.TokensUsed
			modelUsed = chunk.Model

			chunkChan <- StreamChunk{
				Type:      "content",
				Content:   chunk.Content,
				SessionID: req.SessionID,
				Timestamp: time.Now(),
			}

			// Check context for cancellation
			select {
			case <-ctx.Done():
				return
			default:
			}
		}

		// Store final assistant message
		assistantMessageID := uuid.New()
		citations := s.buildCitations(searchResults)
		responseTime := int(time.Since(startTime).Milliseconds())

		assistantMessage := &storage.ChatMessage{
			ID:            assistantMessageID,
			SessionID:     req.SessionID,
			MessageType:   "assistant",
			Content:       fullResponse.String(),
			Citations:     s.citationsToJSON(citations),
			ModelUsed:     &modelUsed,
			TokensUsed:    &tokensUsed,
			ResponseTime:  &responseTime,
			CreatedAt:     time.Now(),
		}

		if err := s.db.CreateChatMessage(ctx, assistantMessage); err != nil {
			s.logger.Warnf("Failed to store assistant message: %v", err)
		}

		// Send final chunk with metadata
		chunkChan <- StreamChunk{
			Type:       "end",
			SessionID:  req.SessionID,
			MessageID:  assistantMessageID,
			Citations:  citations,
			ModelUsed:  modelUsed,
			TokensUsed: tokensUsed,
			Timestamp:  time.Now(),
		}
	}()

	return chunkChan, nil
}

// getOrCreateSession gets an existing session or creates a new one
func (s *Service) getOrCreateSession(ctx context.Context, sessionID uuid.UUID, modelName string) (*storage.ChatSession, error) {
	session, err := s.db.GetChatSession(ctx, sessionID)
	if err == nil {
		// Session exists, update model if different
		if modelName != "" && modelName != session.ModelName {
			session.ModelName = modelName
			// TODO: Update session in database
		}
		return session, nil
	}

	// Create new session
	newSession := &storage.ChatSession{
		ID:           sessionID,
		ModelName:    modelName,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	if err := s.db.CreateChatSession(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return newSession, nil
}

// buildContext creates context string from search results
func (s *Service) buildContext(results []storage.SearchResult) string {
	if len(results) == 0 {
		return "No relevant documents found."
	}

	var context strings.Builder
	context.WriteString("Context from relevant documents:\n\n")

	for i, result := range results {
		context.WriteString(fmt.Sprintf("Document %d (%s):\n", i+1, result.Document.Filename))
		if result.PageNumber != nil {
			context.WriteString(fmt.Sprintf("Page %d:\n", *result.PageNumber))
		}
		context.WriteString(result.Content)
		context.WriteString("\n\n")

		// Limit context length to avoid token limits
		if context.Len() > 4000 {
			break
		}
	}

	return context.String()
}

// generateResponse generates AI response using the retrieved context
func (s *Service) generateResponse(ctx context.Context, query, context, model string) (*AIResponse, error) {
	// Build system prompt
	systemPrompt := fmt.Sprintf(`You are a helpful AI assistant. Answer the user's question based ONLY on the provided context. 
If the context doesn't contain enough information to answer the question, say so clearly.
Do not use any external knowledge or make up information.
Cite the documents you use in your answer.

Context:
%s

Question: %s`, context, query)

	return s.aiClient.GenerateResponse(ctx, systemPrompt, query, model)
}

// buildCitations creates citation objects from search results
func (s *Service) buildCitations(results []storage.SearchResult) []storage.DocumentCitation {
	var citations []storage.DocumentCitation

	for _, result := range results {
		citation := storage.DocumentCitation{
			DocumentID: result.Document.ID,
			Filename:   result.Document.Filename,
			PageNumber: result.PageNumber,
			ChunkIndex: result.ChunkIndex,
			Content:    result.Content,
			Similarity: result.Similarity,
		}
		citations = append(citations, citation)
	}

	return citations
}

// citationsToJSON converts citations to JSON format for storage
func (s *Service) citationsToJSON(citations []storage.DocumentCitation) interface{} {
	// Convert citations to JSON and return as interface{}
	data, err := json.Marshal(citations)
	if err != nil {
		s.logger.Errorf("Failed to marshal citations: %v", err)
		return []interface{}{}
	}
	
	// Parse back to interface{} for proper JSONB handling
	var result interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		s.logger.Errorf("Failed to unmarshal citations: %v", err)
		return []interface{}{}
	}
	
	return result
}

// StreamChunk represents a chunk of streaming response
type StreamChunk struct {
	Type       string                    `json:"type"`        // "start", "content", "end", "error"
	Content    string                    `json:"content"`     // for "content" type
	SessionID  uuid.UUID                 `json:"session_id"`
	MessageID  uuid.UUID                 `json:"message_id,omitempty"`
	Citations  []storage.DocumentCitation `json:"citations,omitempty"`
	ModelUsed  string                    `json:"model_used,omitempty"`
	TokensUsed int                       `json:"tokens_used,omitempty"`
	Timestamp  time.Time                 `json:"timestamp"`
	Message    string                    `json:"message,omitempty"` // for "error" type
	Error      string                    `json:"error,omitempty"`
}

// AIResponse represents a response from the AI
type AIResponse struct {
	Content   string `json:"content"`
	ModelUsed string `json:"model_used"`
	TokensUsed int   `json:"tokens_used"`
}

// GenerateEmbedding generates embedding for text content
func (s *Service) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	return s.aiClient.GenerateEmbedding(ctx, text)
}
