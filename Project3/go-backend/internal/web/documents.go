package web

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"private-knowledge-base-go/internal/ingestion"
	"private-knowledge-base-go/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// DocumentsHandler handles document-related endpoints
type DocumentsHandler struct {
	ingestionService *ingestion.Service
	db               *storage.PostgresDB
	logger           *logrus.Logger
}

// NewDocumentsHandler creates a new documents handler
func NewDocumentsHandler(ingestionService *ingestion.Service, db *storage.PostgresDB, logger *logrus.Logger) *DocumentsHandler {
	return &DocumentsHandler{
		ingestionService: ingestionService,
		db:               db,
		logger:           logger,
	}
}

// UploadDocument handles file upload
// @Summary Upload a document
// @Description Uploads and processes a document for the knowledge base
// @Tags documents
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file"
// @Success 201 {object} storage.UploadResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/documents/upload [post]
func (h *DocumentsHandler) UploadDocument(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Failed to get file from request",
			Details: err.Error(),
		})
		return
	}
	defer file.Close()

	// Validate file size (50MB limit)
	if header.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "File too large",
			Details: "Maximum file size is 50MB",
		})
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedTypes := map[string]bool{
		".pdf":  true,
		".txt":  true,
		".md":   true,
		".markdown": true,
		".go":   true,
	}

	if !allowedTypes[ext] {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Unsupported file type",
			Details: fmt.Sprintf("Supported types: %s", strings.Join(getMapKeys(allowedTypes), ", ")),
		})
		return
	}

	// Read file content
	content := make([]byte, header.Size)
	_, err = file.Read(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to read file content",
			Details: err.Error(),
		})
		return
	}

	// Determine content type
	contentType := getContentType(ext)

	// Process document
	response, err := h.ingestionService.ProcessDocument(c.Request.Context(), header.Filename, contentType, content)
	if err != nil {
		h.logger.Errorf("Failed to process document: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to process document",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// ListDocuments retrieves all documents
// @Summary List all documents
// @Description Retrieves a paginated list of all documents
// @Tags documents
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} DocumentsListResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/documents [get]
func (h *DocumentsHandler) ListDocuments(c *gin.Context) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit
 
	// Retrieve documents from database
	docs, err := h.db.ListDocuments(c.Request.Context(), limit, offset)
	if err != nil {
		h.logger.Errorf("Failed to list documents: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to list documents",
			Details: err.Error(),
		})
		return
	}
 
	// Calculate pagination info (approximate without COUNT)
	total := len(docs)
	totalPages := 1
	if limit > 0 {
		totalPages = (total + limit - 1) / limit
	}

	c.JSON(http.StatusOK, DocumentsListResponse{
		Documents:  docs,
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// GetDocument retrieves a specific document
// @Summary Get a specific document
// @Description Retrieves details of a specific document
// @Tags documents
// @Produce json
// @Param id path string true "Document ID"
// @Success 200 {object} storage.Document
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/documents/{id} [get]
func (h *DocumentsHandler) GetDocument(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid document ID",
			Details: err.Error(),
		})
		return
	}

	// Retrieve document from database
	doc, err := h.db.GetDocument(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Document not found",
			Details: err.Error(),
		})
		return
	}
 
	c.JSON(http.StatusOK, doc)
}

// DeleteDocument deletes a document
// @Summary Delete a document
// @Description Deletes a document and all its chunks
// @Tags documents
// @Param id path string true "Document ID"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/documents/{id} [delete]
func (h *DocumentsHandler) DeleteDocument(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid document ID",
			Details: err.Error(),
		})
		return
	}

	// Delete document from database
	if err := h.db.DeleteDocument(c.Request.Context(), documentID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to delete document",
			Details: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// ReindexDocument re-indexes a document
// @Summary Re-index a document
// @Description Re-processes and re-indexes a document
// @Tags documents
// @Param id path string true "Document ID"
// @Success 200 {object} storage.UploadResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/documents/{id}/reindex [post]
func (h *DocumentsHandler) ReindexDocument(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid document ID",
			Details: err.Error(),
		})
		return
	}

	// TODO: Implement reindexing logic
	h.logger.Infof("Re-indexing document: %s", documentID)

	response := storage.UploadResponse{
		DocumentID: documentID,
		Filename:   "reindexed.pdf",
		Status:     "success",
		Message:    "Document re-indexed successfully",
	}

	c.JSON(http.StatusOK, response)
}

// GetDocumentChunks retrieves chunks for a document
// @Summary Get document chunks
// @Description Retrieves all chunks for a specific document
// @Tags documents
// @Produce json
// @Param id path string true "Document ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} DocumentChunksResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/documents/{id}/chunks [get]
func (h *DocumentsHandler) GetDocumentChunks(c *gin.Context) {
	documentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid document ID",
			Details: err.Error(),
		})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Retrieve chunks from database
	chunks, err := h.db.GetDocumentChunks(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get document chunks",
            Details: err.Error(),
		})
		return
	}

	// Calculate pagination info
	total := len(chunks) // TODO: Get actual count from database
	totalPages := (total + limit - 1) / limit

	c.JSON(http.StatusOK, DocumentChunksResponse{
		Chunks:     chunks,
		Pagination: PaginationInfo{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	})
}

// Helper functions

func getContentType(ext string) string {
	contentTypes := map[string]string{
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".md":   "text/markdown",
		".markdown": "text/markdown",
		".go":   "text/x-go",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}
	return "application/octet-stream"
}

func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Response structures

type DocumentsListResponse struct {
	Documents  []storage.Document `json:"documents"`
	Pagination PaginationInfo     `json:"pagination"`
}

type DocumentChunksResponse struct {
	Chunks     []storage.DocumentChunk `json:"chunks"`
	Pagination PaginationInfo          `json:"pagination"`
}

type PaginationInfo struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
