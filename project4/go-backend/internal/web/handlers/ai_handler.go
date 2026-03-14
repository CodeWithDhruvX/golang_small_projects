package handlers

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"ai-recruiter-assistant/internal/ai"
	"ai-recruiter-assistant/internal/auth"
	"ai-recruiter-assistant/internal/rag"
	"ai-recruiter-assistant/internal/storage"
)

// CacheEntry represents a cached reply
type CacheEntry struct {
	Reply       string    `json:"reply"`
	ModelUsed   string    `json:"model_used"`
	ResponseTime int      `json:"response_time"`
	CreatedAt   time.Time `json:"created_at"`
	WordCount   int       `json:"word_count"`
}

// AIHandler handles AI-related HTTP requests
type AIHandler struct {
	storage     storage.StorageInterface
	authService *auth.AuthService
	ollama      *ai.OllamaService
	rag         *rag.RAGService
	cache       map[string]*CacheEntry
	cacheMutex  sync.RWMutex
}

// NewAIHandler creates a new AI handler
func NewAIHandler(storage storage.StorageInterface, authService *auth.AuthService) *AIHandler {
	// Initialize Ollama service
	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}
	
	ollama := ai.NewOllamaService(ollamaURL)
	ragService := rag.NewRAGService(storage, ollama)
	
	handler := &AIHandler{
		storage:     storage,
		authService: authService,
		ollama:      ollama,
		rag:         ragService,
		cache:       make(map[string]*CacheEntry),
	}
	
	// Start cache cleanup goroutine
	go handler.cleanupCache()
	
	return handler
}

// generateCacheKey creates a cache key for email content
func (h *AIHandler) generateCacheKey(body, model string) string {
	// Normalize email body for better cache hits
	normalizedBody := strings.ToLower(strings.TrimSpace(body))
	// Remove extra whitespace and normalize common variations
	normalizedBody = strings.Join(strings.Fields(normalizedBody), " ")
	
	// Remove common recruiter email variations that don't affect reply content
	replacements := map[string]string{
		"dear candidate": "dear candidate",
		"hi candidate": "hi candidate", 
		"hello candidate": "hello candidate",
		"greetings": "greetings",
		"regarding": "about",
		"concerning": "about",
		"with reference to": "about",
	}
	
	for old, new := range replacements {
		normalizedBody = strings.ReplaceAll(normalizedBody, old, new)
	}
	
	// Create hash
	hasher := md5.New()
	hasher.Write([]byte(normalizedBody + model))
	return hex.EncodeToString(hasher.Sum(nil))
}

// getSimilarFromCache retrieves cached reply for similar emails
func (h *AIHandler) getSimilarFromCache(body, model string) (*CacheEntry, bool) {
	h.cacheMutex.RLock()
	defer h.cacheMutex.RUnlock()
	
	// Extract key phrases from email for similarity matching
	phrases := h.extractKeyPhrases(body)
	
	for _, entry := range h.cache {
		// Check if cache entry is still valid (30 minutes)
		if time.Since(entry.CreatedAt) > 30*time.Minute {
			continue
		}
		
		// Check semantic similarity
		if h.isSimilarEmail(phrases, entry.Reply) {
			logrus.Infof("Found similar cached reply for similar email pattern")
			return entry, true
		}
	}
	
	return nil, false
}

// extractKeyPhrases extracts important phrases from email
func (h *AIHandler) extractKeyPhrases(body string) []string {
	lowerBody := strings.ToLower(body)
	phrases := []string{}
	
	// Common recruiter email patterns
	patterns := []string{
		"job opportunity", "position", "role", "hiring", "recruit",
		"resume", "cv", "interview", "salary", "compensation",
		"experience", "skills", "qualification", "notice period",
		"available", "join", "start date", "location",
	}
	
	for _, pattern := range patterns {
		if strings.Contains(lowerBody, pattern) {
			phrases = append(phrases, pattern)
		}
	}
	
	return phrases
}

// isSimilarEmail checks if cached reply matches email pattern
func (h *AIHandler) isSimilarEmail(phrases []string, cachedReply string) bool {
	// Simple heuristic: if email has similar key phrases, reply pattern should be similar
	if len(phrases) == 0 {
		return false
	}
	
	// Check if cached reply contains relevant response patterns
	replyLower := strings.ToLower(cachedReply)
	responsePatterns := []string{
		"interested", "opportunity", "discuss", "interview",
		"resume", "experience", "skills", "available",
	}
	
	matchedPatterns := 0
	for _, phrase := range phrases {
		for _, pattern := range responsePatterns {
			if strings.Contains(replyLower, pattern) && strings.Contains(replyLower, phrase) {
				matchedPatterns++
				break
			}
		}
	}
	
	// Consider similar if at least 30% of patterns match
	return matchedPatterns >= (len(phrases)+2)/3
}

// getFromCache retrieves a cached reply if available
func (h *AIHandler) getFromCache(cacheKey string) (*CacheEntry, bool) {
	h.cacheMutex.RLock()
	defer h.cacheMutex.RUnlock()
	
	entry, exists := h.cache[cacheKey]
	if !exists {
		return nil, false
	}
	
	// Check if cache entry is still valid (30 minutes)
	if time.Since(entry.CreatedAt) > 30*time.Minute {
		return nil, false
	}
	
	return entry, true
}

// setCache stores a reply in cache
func (h *AIHandler) setCache(cacheKey string, reply *CacheEntry) {
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock()
	
	h.cache[cacheKey] = reply
}

// cleanupCache removes old cache entries
func (h *AIHandler) cleanupCache() {
	ticker := time.NewTicker(10 * time.Minute) // Clean every 10 minutes
	defer ticker.Stop()
	
	for range ticker.C {
		h.cacheMutex.Lock()
		for key, entry := range h.cache {
			if time.Since(entry.CreatedAt) > 30*time.Minute {
				delete(h.cache, key)
			}
		}
		h.cacheMutex.Unlock()
	}
}

// ClassifyEmailRequest represents the request to classify an email
type ClassifyEmailRequest struct {
	EmailID string `json:"email_id" binding:"required"`
	Body    string `json:"body" binding:"required"`
}

// ClassifyEmailResponse represents the response from email classification
type ClassifyEmailResponse struct {
	IsRecruiter bool    `json:"is_recruiter"`
	Confidence  float64 `json:"confidence"`
	Reason      string  `json:"reason"`
}

// ClassifyEmail classifies whether an email is from a recruiter
// @Summary Classify email
// @Description Classify whether an email is from a recruiter using AI
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ClassifyEmailRequest true "Email classification request"
// @Success 200 {object} ClassifyEmailResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/classify-email [post]
func (h *AIHandler) ClassifyEmail(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ClassifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use Ollama to classify the email
	isRecruiter, confidence, err := h.ollama.ClassifyEmail(c.Request.Context(), req.Body)
	if err != nil {
		logrus.Errorf("Failed to classify email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to classify email"})
		return
	}

	// Update email record with classification result
	if req.EmailID != "" {
		email, err := h.storage.GetEmailByID(req.EmailID, userID)
		if err == nil {
			email.IsRecruiter = isRecruiter
			h.storage.UpdateEmail(email)
		}
	}

	response := ClassifyEmailResponse{
		IsRecruiter: isRecruiter,
		Confidence:  confidence,
		Reason:      "Classified using AI model",
	}

	c.JSON(http.StatusOK, response)
}

// GenerateReplyRequest represents the request for reply generation
type GenerateReplyRequest struct {
	EmailID string `json:"email_id"`
	Body    string `json:"body"`
	Model   string `json:"model"`
}

// GenerateReplyResponse represents the response from reply generation
type GenerateReplyResponse struct {
	Reply        string  `json:"reply"`
	ModelUsed    string  `json:"model_used"`
	ResponseTime int     `json:"response_time"`
	Confidence   float64 `json:"confidence"`
	WordCount    int     `json:"word_count"`
	IsTemplate   bool    `json:"is_template"`
	GPUEnabled   bool    `json:"gpu_enabled"`
	GPUType      string  `json:"gpu_type,omitempty"`
}

// GenerateReply generates an AI reply to a recruiter email
// @Summary Generate reply
// @Description Generate a professional reply to a recruiter email using AI
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body GenerateReplyRequest true "Reply generation request"
// @Success 200 {object} GenerateReplyResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/generate-reply [post]
func (h *AIHandler) GenerateReply(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req GenerateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logrus.Errorf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that body is not empty
	if strings.TrimSpace(req.Body) == "" {
		logrus.Warnf("Empty body received for EmailID: %s", req.EmailID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email body is required for reply generation"})
		return
	}

	// Validate email body length to prevent very long inputs
	if len(req.Body) > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email body is too long (max 10,000 characters)"})
		return
	}

	startTime := time.Now()

	// Determine model to use for record-keeping
	targetModel := req.Model
	if targetModel == "" {
		targetModel = "llama3.1:8b"
	}

	logrus.Infof("Generating reply for user: %s, email_id: %s, model: %s", userID, req.EmailID, targetModel)

	// Get GPU status for response
	gpuStatus := h.ollama.GetGPUStatus()
	gpuEnabled := gpuStatus["enabled"].(bool)
	gpuType := ""
	if gpuEnabled {
		gpuType = gpuStatus["gpu_type"].(string)
	}

	// Check cache first for similar emails
	cacheKey := h.generateCacheKey(req.Body, targetModel)
	if cachedEntry, found := h.getFromCache(cacheKey); found {
		logrus.Infof("Returning exact cached reply for user: %s, cache age: %v", userID, time.Since(cachedEntry.CreatedAt))
		
		response := GenerateReplyResponse{
			Reply:        cachedEntry.Reply,
			ModelUsed:    cachedEntry.ModelUsed,
			ResponseTime: cachedEntry.ResponseTime,
			Confidence:   0.85, // Slightly lower confidence for cached responses
			WordCount:    cachedEntry.WordCount,
			IsTemplate:   false,
			GPUEnabled:   gpuEnabled,
			GPUType:      gpuType,
		}
		
		c.JSON(http.StatusOK, response)
		return
	}

	// Check for similar emails if exact cache miss
	if similarEntry, found := h.getSimilarFromCache(req.Body, targetModel); found {
		logrus.Infof("Returning similar cached reply for user: %s", userID)
		
		response := GenerateReplyResponse{
			Reply:        similarEntry.Reply,
			ModelUsed:    similarEntry.ModelUsed,
			ResponseTime: similarEntry.ResponseTime,
			Confidence:   0.75, // Lower confidence for similar responses
			WordCount:    similarEntry.WordCount,
			IsTemplate:   false,
			GPUEnabled:   gpuEnabled,
			GPUType:      gpuType,
		}
		
		c.JSON(http.StatusOK, response)
		return
	}

	// Use RAG service to generate contextual reply
	reply, err := h.rag.GenerateContextualReply(c.Request.Context(), userID, req.Body, req.Model)
	if err != nil {
		logrus.Errorf("Failed to generate reply: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reply"})
		return
	}

	responseTime := int(time.Since(startTime).Milliseconds())
	
	// Calculate word count and detect if it's likely a template
	wordCount := len(strings.Fields(reply))
	isTemplate := wordCount < 50 || strings.Contains(strings.ToLower(reply), "dear recruiter,") && strings.Contains(strings.ToLower(reply), "best regards")

	// Store in cache for future requests
	cacheEntry := &CacheEntry{
		Reply:        reply,
		ModelUsed:    targetModel,
		ResponseTime: responseTime,
		CreatedAt:    time.Now(),
		WordCount:    wordCount,
	}
	h.setCache(cacheKey, cacheEntry)

	// Calculate confidence based on various factors
	confidence := 0.92 // Base confidence
	if isTemplate {
		confidence -= 0.2 // Reduce confidence for template responses
	}
	if wordCount < 30 {
		confidence -= 0.1 // Reduce confidence for very short responses
	}
	if wordCount > 200 {
		confidence -= 0.1 // Reduce confidence for overly long responses
	}

	// Store the generated reply
	aiReply := &storage.AIReply{
		UserID:       userID,
		EmailID:      req.EmailID,
		ReplyContent: reply,
		ModelUsed:    targetModel,
		ResponseTime: responseTime,
		IsSent:       false,
	}

	err = h.storage.CreateAIReply(aiReply)
	if err != nil {
		logrus.Warnf("Failed to store AI reply: %v", err)
	}

	logrus.Infof("Successfully generated reply for user: %s, words: %d, confidence: %.2f", userID, wordCount, confidence)

	response := GenerateReplyResponse{
		Reply:        reply,
		ModelUsed:    targetModel,
		ResponseTime: responseTime,
		Confidence:   confidence,
		WordCount:    wordCount,
		IsTemplate:   isTemplate,
		GPUEnabled:   gpuEnabled,
		GPUType:      gpuType,
	}

	c.JSON(http.StatusOK, response)
}

// GPUStatusResponse represents the response from GPU status check
type GPUStatusResponse struct {
	Enabled  bool   `json:"enabled"`
	GPUType  string `json:"gpu_type"`
	Detected bool   `json:"detected"`
	Message  string `json:"message"`
}

// GetGPUStatus returns the current GPU status
// @Summary Get GPU status
// @Description Get the current GPU acceleration status and type
// @Tags ai
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GPUStatusResponse
// @Failure 401 {object} map[string]interface{}
// @Router /ai/gpu-status [get]
func (h *AIHandler) GetGPUStatus(c *gin.Context) {
	_, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	gpuStatus := h.ollama.GetGPUStatus()
	
	response := GPUStatusResponse{
		Enabled:  gpuStatus["enabled"].(bool),
		GPUType:  gpuStatus["gpu_type"].(string),
		Detected: gpuStatus["detected"].(bool),
		Message:  getGPUMessage(gpuStatus),
	}

	c.JSON(http.StatusOK, response)
}

// getGPUMessage generates a user-friendly GPU status message
func getGPUMessage(gpuStatus map[string]interface{}) string {
	enabled := gpuStatus["enabled"].(bool)
	gpuType := gpuStatus["gpu_type"].(string)
	
	if !enabled {
		return "GPU acceleration is not available. Using CPU for AI processing."
	}
	
	switch gpuType {
	case "nvidia":
		return "NVIDIA GPU acceleration is enabled for faster AI processing."
	case "amd":
		return "AMD GPU acceleration is enabled for faster AI processing."
	case "apple_silicon":
		return "Apple Silicon GPU acceleration is enabled for faster AI processing."
	default:
		return "GPU acceleration is enabled for faster AI processing."
	}
}

// ExtractRequirementsRequest represents the request for requirement extraction
type ExtractRequirementsRequest struct {
	EmailID string `json:"email_id"`
	Body    string `json:"body" binding:"required"`
}

// ExtractRequirementsResponse represents the response from requirement extraction
type ExtractRequirementsResponse struct {
	Resume        bool `json:"resume"`
	Experience    bool `json:"experience"`
	ExpectedCTC   bool `json:"expected_ctc"`
	NoticePeriod  bool `json:"notice_period"`
	Skills        bool `json:"skills"`
	Portfolio     bool `json:"portfolio"`
	CoverLetter   bool `json:"cover_letter"`
	Availability  bool `json:"availability"`
}

// ExtractRequirements extracts requested information from recruiter email
// @Summary Extract requirements
// @Description Extract what information the recruiter is requesting from the email
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ExtractRequirementsRequest true "Requirement extraction request"
// @Success 200 {object} ExtractRequirementsResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/extract-requirements [post]
func (h *AIHandler) ExtractRequirements(c *gin.Context) {
	_, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req ExtractRequirementsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use Ollama to extract requirements
	requirements, err := h.ollama.ExtractRequirements(c.Request.Context(), req.Body)
	if err != nil {
		logrus.Errorf("Failed to extract requirements: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract requirements"})
		return
	}

	response := ExtractRequirementsResponse{
		Resume:        requirements["resume"],
		Experience:    requirements["experience"],
		ExpectedCTC:   requirements["expected_ctc"],
		NoticePeriod:  requirements["notice_period"],
		Skills:        requirements["skills"],
		Portfolio:     requirements["portfolio"],
		CoverLetter:   requirements["cover_letter"],
		Availability:  requirements["availability"],
	}

	c.JSON(http.StatusOK, response)
}

// SearchContextRequest represents the request for context search
type SearchContextRequest struct {
	Query string `json:"query" binding:"required"`
	TopK  int    `json:"top_k"`
}

// ContextResult represents a single context result
type ContextResult struct {
	ID       string                 `json:"id"`
	Content  string                 `json:"content"`
	Source   string                 `json:"source"`
	Score    float64                `json:"score"`
	Metadata map[string]interface{} `json:"metadata"`
}

// SearchContextResponse represents the response from context search
type SearchContextResponse struct {
	Results []ContextResult `json:"results"`
	Query   string          `json:"query"`
	Total   int             `json:"total"`
}

// SearchContext performs semantic search for relevant context
// @Summary Search context
// @Description Perform semantic search to retrieve relevant context using RAG
// @Tags ai
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SearchContextRequest true "Context search request"
// @Success 200 {object} SearchContextResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /ai/search-context [post]
func (h *AIHandler) SearchContext(c *gin.Context) {
	userID, exists := auth.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req SearchContextRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default top_k if not provided
	if req.TopK <= 0 {
		req.TopK = 5
	}

	// Use RAG service to perform semantic search
	documents, err := h.rag.SearchContext(c.Request.Context(), userID, req.Query, req.TopK)
	if err != nil {
		logrus.Errorf("Failed to search context: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search context"})
		return
	}

	// Convert to response format
	results := make([]ContextResult, len(documents))
	for i, doc := range documents {
		results[i] = ContextResult{
			ID:       doc.ID,
			Content:  doc.Content,
			Source:   doc.Source,
			Score:    0.95, // TODO: Get actual similarity score from RAG service
			Metadata: map[string]interface{}{}, // TODO: Get actual metadata
		}
	}

	response := SearchContextResponse{
		Results: results,
		Query:   req.Query,
		Total:   len(results),
	}

	c.JSON(http.StatusOK, response)
}
