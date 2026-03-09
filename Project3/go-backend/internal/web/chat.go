package web

import (
	"net/http"
	"strconv"
	"time"

	"private-knowledge-base-go/internal/rag"
	"private-knowledge-base-go/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var allowedModels = map[string]bool{
	"llama3.1:8b":      true,
	"qwen2.5-coder:3b": true,
	"phi3:latest":      true,
}

// ChatHandler handles chat-related endpoints
type ChatHandler struct {
	ragService *rag.Service
	db         *storage.PostgresDB
	logger     *logrus.Logger
}

// NewChatHandler creates a new chat handler
func NewChatHandler(ragService *rag.Service, db *storage.PostgresDB, logger *logrus.Logger) *ChatHandler {
	return &ChatHandler{
		ragService: ragService,
		db:         db,
		logger:     logger,
	}
}

// CreateSession creates a new chat session
// @Summary Create a new chat session
// @Description Creates a new chat session with the specified model
// @Tags chat
// @Accept json
// @Produce json
// @Param request body CreateSessionRequest true "Session creation request"
// @Success 201 {object} ChatSessionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/sessions [post]
func (h *ChatHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	sessionID := uuid.New()
	session := &storage.ChatSession{
		ID:           sessionID,
		SessionName:  &req.SessionName,
		ModelName:    req.Model,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
	}

	if !allowedModels[session.ModelName] {
		session.ModelName = "llama3.1:8b"
	}

	if err := h.db.CreateChatSession(c.Request.Context(), session); err != nil {
		h.logger.Errorf("Failed to create session: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create chat session",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, ChatSessionResponse{
		SessionID:   sessionID,
		SessionName: req.SessionName,
		Model:       session.ModelName,
		CreatedAt:   session.CreatedAt,
		LastActivity: session.LastActivity,
	})
}

// GetSessions retrieves all chat sessions
// @Summary Get all chat sessions
// @Description Retrieves a list of all chat sessions
// @Tags chat
// @Produce json
// @Success 200 {array} ChatSessionResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/sessions [get]
func (h *ChatHandler) GetSessions(c *gin.Context) {
	list, err := h.db.ListChatSessions(c.Request.Context())
	if err != nil {
		h.logger.Errorf("Failed to list sessions: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve chat sessions",
			Details: err.Error(),
		})
		return
	}

	resp := make([]ChatSessionResponse, 0, len(list))
	for _, s := range list {
		name := ""
		if s.SessionName != nil {
			name = *s.SessionName
		}
		resp = append(resp, ChatSessionResponse{
			SessionID:   s.ID,
			SessionName: name,
			Model:       s.ModelName,
			CreatedAt:   s.CreatedAt,
			LastActivity: s.LastActivity,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": resp,
	})
}

// GetSession retrieves a specific chat session
// @Summary Get a specific chat session
// @Description Retrieves details of a specific chat session
// @Tags chat
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} ChatSessionResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/sessions/{id} [get]
func (h *ChatHandler) GetSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid session ID",
			Details: err.Error(),
		})
		return
	}

	s, err := h.db.GetChatSession(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error: "Chat session not found",
			Details: err.Error(),
		})
		return
	}
	name := ""
	if s.SessionName != nil {
		name = *s.SessionName
	}
	session := ChatSessionResponse{
		SessionID:   s.ID,
		SessionName: name,
		Model:       s.ModelName,
		CreatedAt:   s.CreatedAt,
		LastActivity: s.LastActivity,
	}

	c.JSON(http.StatusOK, session)
}

// GetMessages retrieves messages for a chat session
// @Summary Get messages for a session
// @Description Retrieves all messages for a specific chat session
// @Tags chat
// @Produce json
// @Param id path string true "Session ID"
// @Param limit query int false "Limit number of messages" default(50)
// @Success 200 {array} ChatMessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/sessions/{id}/messages [get]
func (h *ChatHandler) GetMessages(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid session ID",
			Details: err.Error(),
		})
		return
	}

	limit := 50
	page := 1
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	items, err := h.db.GetChatMessages(c.Request.Context(), sessionID, limit)
	if err != nil {
		h.logger.Errorf("Failed to get messages: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve messages",
			Details: err.Error(),
		})
		return
	}
	resp := make([]ChatMessageResponse, 0, len(items))
	for _, m := range items {
		resp = append(resp, ChatMessageResponse{
			ID:          m.ID,
			SessionID:   m.SessionID,
			MessageType: m.MessageType,
			Content:     m.Content,
			Citations:   nil,
			ModelUsed:   m.ModelUsed,
			TokensUsed:  m.TokensUsed,
			ResponseTime: m.ResponseTime,
			CreatedAt:    m.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"messages": resp,
		"pagination": gin.H{
			"page":        page,
			"limit":       limit,
			"total":       len(resp),
			"total_pages": 1,
		},
	})
}

// Chat handles a chat request
// @Summary Send a chat message
// @Description Sends a message and gets AI response
// @Tags chat
// @Accept json
// @Produce json
// @Param request body ChatRequest true "Chat request"
// @Success 200 {object} ChatResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat [post]
func (h *ChatHandler) Chat(c *gin.Context) {
	var req storage.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	if req.Model != "" && !allowedModels[req.Model] {
		req.Model = "llama3.1:8b"
	}

	response, err := h.ragService.Chat(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorf("Chat request failed: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to process chat request",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":       response.SessionID,
		"message_id":       response.MessageID,
		"content":          response.Message,
		"citations":        response.Citations,
		"model_used":       response.ModelUsed,
		"tokens_used":      response.TokensUsed,
		"response_time_ms": response.ResponseTime,
		"created_at":       response.Timestamp,
	})
}

// StreamChat handles streaming chat responses
// @Summary Stream chat response
// @Description Sends a message and streams AI response using Server-Sent Events
// @Tags chat
// @Accept json
// @Produce text/event-stream
// @Param request body ChatRequest true "Chat request"
// @Success 200 {string} string "Event stream"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/stream [post]
func (h *ChatHandler) StreamChat(c *gin.Context) {
	var req storage.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	if req.Model != "" && !allowedModels[req.Model] {
		req.Model = "llama3.1:8b"
	}

	// Set SSE headers
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, Cache-Control, Accept")

	// Get stream from RAG service
	streamChan, err := h.ragService.StreamChat(c.Request.Context(), &req)
	if err != nil {
		h.logger.Errorf("Stream chat request failed: %v", err)
		c.SSEvent("error", gin.H{
			"type":    "error",
			"error":   "Failed to process chat request",
			"details": err.Error(),
		})
		return
	}

	// Stream response
	for chunk := range streamChan {
		switch chunk.Type {
		case "start":
			c.SSEvent("start", gin.H{
				"type":       "start",
				"session_id": chunk.SessionID,
			})
		case "content":
			c.SSEvent("chunk", gin.H{
				"type":       "chunk",
				"content":    chunk.Content,
				"session_id": chunk.SessionID,
			})
		case "end":
			c.SSEvent("end", gin.H{
				"type":        "end",
				"session_id":  chunk.SessionID,
				"message_id":  chunk.MessageID,
				"citations":   chunk.Citations,
				"tokens_used": chunk.TokensUsed,
			})
			return
		case "error":
			c.SSEvent("error", gin.H{
				"type":  "error",
				"error": chunk.Message,
			})
			return
		}

		// Flush response to client
		c.Writer.Flush()
	}
}

// DeleteSession deletes a chat session
// @Summary Delete a chat session
// @Description Deletes a chat session and all its messages
// @Tags chat
// @Param id path string true "Session ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/chat/sessions/{id} [delete]
func (h *ChatHandler) DeleteSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid session ID",
			Details: err.Error(),
		})
		return
	}

	if err := h.db.DeleteChatSession(c.Request.Context(), sessionID); err != nil {
		h.logger.Errorf("Failed to delete session: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to delete chat session",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Chat session deleted successfully",
	})
}

// Request/Response structures

type CreateSessionRequest struct {
	SessionName string `json:"session_name" binding:"required"`
	Model string `json:"model" binding:"required"`
}

type ChatSessionResponse struct {
	SessionID   uuid.UUID `json:"session_id"`
	SessionName string    `json:"session_name"`
	Model       string    `json:"model"`
	CreatedAt   time.Time `json:"created_at"`
	LastActivity time.Time `json:"last_activity"`
}

type ChatMessageResponse struct {
	ID           uuid.UUID                 `json:"id"`
	SessionID    uuid.UUID                 `json:"session_id"`
	MessageType  string                    `json:"message_type"`
	Content      string                    `json:"content"`
	Citations    []storage.DocumentCitation `json:"citations,omitempty"`
	ModelUsed    *string                   `json:"model_used,omitempty"`
	TokensUsed   *int                      `json:"tokens_used,omitempty"`
	ResponseTime *int                      `json:"response_time_ms,omitempty"`
	CreatedAt    time.Time                 `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
