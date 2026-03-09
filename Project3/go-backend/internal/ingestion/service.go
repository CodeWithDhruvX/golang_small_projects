package ingestion

import (
	"context"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"private-knowledge-base-go/internal/rag"
	"private-knowledge-base-go/internal/storage"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Service handles document ingestion and processing
type Service struct {
	db           *storage.PostgresDB
	logger       *logrus.Logger
	pdfParser    *PDFParser
	mdParser     *MarkdownParser
	textParser   *TextParser
	goParser     *GoFileParser
	ragService   *rag.Service
}

// NewService creates a new ingestion service
func NewService(db *storage.PostgresDB, logger *logrus.Logger, ragService *rag.Service) *Service {
	return &Service{
		db:         db,
		logger:     logger,
		pdfParser:  NewPDFParser(logger),
		mdParser:   NewMarkdownParser(logger),
		textParser: NewTextParser(logger),
		goParser:   NewGoFileParser(logger),
		ragService: ragService,
	}
}

// ProcessDocument processes an uploaded document and stores it in the database
func (s *Service) ProcessDocument(ctx context.Context, filename string, contentType string, content []byte) (*storage.UploadResponse, error) {
	// Create document record
	docID := uuid.New()
	doc := &storage.Document{
		ID:          docID,
		Filename:    filename,
		ContentType: contentType,
		FileSize:    len(content),
		Processed:   false,
		Metadata:    make(map[string]interface{}),
	}

	// Store document metadata
	if err := s.db.CreateDocument(ctx, doc); err != nil {
		return nil, fmt.Errorf("failed to create document record: %w", err)
	}

	// Process document based on type
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch {
	case ext == ".pdf":
		return s.processPDF(ctx, doc, content)
	case ext == ".md" || ext == ".markdown":
		return s.processMarkdown(ctx, doc, content)
	case ext == ".txt":
		return s.processText(ctx, doc, content)
	case ext == ".go":
		return s.processGoFile(ctx, doc, content)
	default:
		return s.processText(ctx, doc, content) // Default to text processing
	}
}

// processPDF handles PDF document processing
func (s *Service) processPDF(ctx context.Context, doc *storage.Document, content []byte) (*storage.UploadResponse, error) {
	// Validate PDF
	if err := s.pdfParser.ValidatePDF(content); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Invalid PDF format: %v", err),
		}, nil
	}

	// Extract metadata
	var err error
	_, err = s.pdfParser.GetPDFMetadata(ctx, content)
	if err != nil {
		s.logger.Warnf("Failed to extract PDF metadata: %v", err)
	}

	// Parse PDF content
	pages, err := s.pdfParser.ParsePDF(ctx, content)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to parse PDF: %v", err),
		}, nil
	}

	// Create chunks from pages
	chunks, err := s.createChunksFromPages(ctx, doc.ID, pages)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to create chunks: %v", err),
		}, nil
	}

	// Store chunks
	if err := s.storeChunks(ctx, chunks); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to store chunks: %v", err),
		}, nil
	}

	// Update document as processed
	if err := s.db.UpdateDocumentProcessed(ctx, doc.ID, true); err != nil {
		s.logger.Warnf("Failed to update document processed status: %v", err)
	}

	return &storage.UploadResponse{
		DocumentID: doc.ID,
		Filename:   doc.Filename,
		Status:     "success",
		Message:    fmt.Sprintf("Successfully processed PDF with %d pages and %d chunks", len(pages), len(chunks)),
	}, nil
}

// processMarkdown handles Markdown document processing
func (s *Service) processMarkdown(ctx context.Context, doc *storage.Document, content []byte) (*storage.UploadResponse, error) {
	// Validate Markdown
	if err := s.mdParser.ValidateMarkdown(content); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Invalid Markdown format: %v", err),
		}, nil
	}

	// Extract metadata
	var err error
	_, err = s.mdParser.GetMarkdownMetadata(ctx, content)
	if err != nil {
		s.logger.Warnf("Failed to extract Markdown metadata: %v", err)
	}

	// Parse Markdown content
	sections, err := s.mdParser.ParseMarkdown(ctx, content)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to parse Markdown: %v", err),
		}, nil
	}

	// Create chunks from sections
	chunks, err := s.createChunksFromSections(ctx, doc.ID, sections)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to create chunks: %v", err),
		}, nil
	}

	// Store chunks
	if err := s.storeChunks(ctx, chunks); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to store chunks: %v", err),
		}, nil
	}

	// Update document as processed
	if err := s.db.UpdateDocumentProcessed(ctx, doc.ID, true); err != nil {
		s.logger.Warnf("Failed to update document processed status: %v", err)
	}

	return &storage.UploadResponse{
		DocumentID: doc.ID,
		Filename:   doc.Filename,
		Status:     "success",
		Message:    fmt.Sprintf("Successfully processed Markdown with %d sections and %d chunks", len(sections), len(chunks)),
	}, nil
}

// processText handles plain text document processing
func (s *Service) processText(ctx context.Context, doc *storage.Document, content []byte) (*storage.UploadResponse, error) {
	// Validate text
	if err := s.textParser.ValidateText(content); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Invalid text format: %v", err),
		}, nil
	}

	// Extract metadata
	var err error
	_, err = s.textParser.GetTextMetadata(ctx, content)
	if err != nil {
		s.logger.Warnf("Failed to extract text metadata: %v", err)
	}

	// Parse text content
	paragraphs, err := s.textParser.ParseText(ctx, content)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to parse text: %v", err),
		}, nil
	}

	// Create chunks from paragraphs
	chunks, err := s.createChunksFromParagraphs(ctx, doc.ID, paragraphs)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to create chunks: %v", err),
		}, nil
	}

	// Store chunks
	if err := s.storeChunks(ctx, chunks); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to store chunks: %v", err),
		}, nil
	}

	// Update document as processed
	if err := s.db.UpdateDocumentProcessed(ctx, doc.ID, true); err != nil {
		s.logger.Warnf("Failed to update document processed status: %v", err)
	}

	return &storage.UploadResponse{
		DocumentID: doc.ID,
		Filename:   doc.Filename,
		Status:     "success",
		Message:    fmt.Sprintf("Successfully processed text with %d paragraphs and %d chunks", len(paragraphs), len(chunks)),
	}, nil
}

// processGoFile handles Go source file processing
func (s *Service) processGoFile(ctx context.Context, doc *storage.Document, content []byte) (*storage.UploadResponse, error) {
	// Validate Go code
	if err := s.goParser.ValidateGoFile(content); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Invalid Go code: %v", err),
		}, nil
	}

	// Extract metadata
	var err error
	_, err = s.goParser.GetGoFileMetadata(ctx, content)
	if err != nil {
		s.logger.Warnf("Failed to extract Go file metadata: %v", err)
	}

	// Parse Go content
	elements, err := s.goParser.ParseGoFile(ctx, content)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to parse Go file: %v", err),
		}, nil
	}

	// Create chunks from elements
	chunks, err := s.createChunksFromGoElements(ctx, doc.ID, elements)
	if err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to create chunks: %v", err),
		}, nil
	}

	// Store chunks
	if err := s.storeChunks(ctx, chunks); err != nil {
		return &storage.UploadResponse{
			DocumentID: doc.ID,
			Filename:   doc.Filename,
			Status:     "error",
			Message:    fmt.Sprintf("Failed to store chunks: %v", err),
		}, nil
	}

	// Update document as processed
	if err := s.db.UpdateDocumentProcessed(ctx, doc.ID, true); err != nil {
		s.logger.Warnf("Failed to update document processed status: %v", err)
	}

	return &storage.UploadResponse{
		DocumentID: doc.ID,
		Filename:   doc.Filename,
		Status:     "success",
		Message:    fmt.Sprintf("Successfully processed Go file with %d elements and %d chunks", len(elements), len(chunks)),
	}, nil
}

// createChunksFromPages creates chunks from PDF pages with overlap
func (s *Service) createChunksFromPages(ctx context.Context, docID uuid.UUID, pages []PageContent) ([]*storage.DocumentChunk, error) {
	var chunks []*storage.DocumentChunk
	chunkIndex := 0
	
	// Combine pages into larger chunks for better context
	const chunkSize = 1000 // Target chunk size in characters
	const overlap = 200     // Overlap between chunks
	
	for _, page := range pages {
		pageContent := page.Content
		
		// Split page content into smaller chunks if it's too large
		for start := 0; start < len(pageContent); start += chunkSize - overlap {
			end := start + chunkSize
			if end > len(pageContent) {
				end = len(pageContent)
			}

			chunk := &storage.DocumentChunk{
				ID:         uuid.New(),
				DocumentID: docID,
				ChunkIndex: chunkIndex,
				Content:    pageContent[start:end],
				PageNumber: &page.PageNumber,
				Metadata:   make(map[string]interface{}),
			}

			chunks = append(chunks, chunk)
			chunkIndex++

			// Check context for cancellation
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
	}

	return chunks, nil
}

// createChunksFromSections creates chunks from Markdown sections
func (s *Service) createChunksFromSections(ctx context.Context, docID uuid.UUID, sections []SectionContent) ([]*storage.DocumentChunk, error) {
	var chunks []*storage.DocumentChunk
	chunkIndex := 0
	
	const chunkSize = 1000
	const overlap = 200

	for _, section := range sections {
		sectionContent := section.Content
		
		// Split section content into chunks
		for start := 0; start < len(sectionContent); start += chunkSize - overlap {
			end := start + chunkSize
			if end > len(sectionContent) {
				end = len(sectionContent)
			}

			chunk := &storage.DocumentChunk{
				ID:         uuid.New(),
				DocumentID: docID,
				ChunkIndex: chunkIndex,
				Content:    sectionContent[start:end],
				Metadata:   make(map[string]interface{}),
			}

			chunks = append(chunks, chunk)
			chunkIndex++

			// Check context for cancellation
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
	}

	return chunks, nil
}

// createChunksFromParagraphs creates chunks from text paragraphs
func (s *Service) createChunksFromParagraphs(ctx context.Context, docID uuid.UUID, paragraphs []ParagraphContent) ([]*storage.DocumentChunk, error) {
	var chunks []*storage.DocumentChunk
	chunkIndex := 0
	
	const chunkSize = 1000
	const overlap = 200

	for _, paragraph := range paragraphs {
		paraContent := paragraph.Content
		
		// Split paragraph content into chunks
		for start := 0; start < len(paraContent); start += chunkSize - overlap {
			end := start + chunkSize
			if end > len(paraContent) {
				end = len(paraContent)
			}

			chunk := &storage.DocumentChunk{
				ID:         uuid.New(),
				DocumentID: docID,
				ChunkIndex: chunkIndex,
				Content:    paraContent[start:end],
				Metadata:   make(map[string]interface{}),
			}

			chunks = append(chunks, chunk)
			chunkIndex++

			// Check context for cancellation
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
	}

	return chunks, nil
}

// createChunksFromGoElements creates chunks from Go code elements
func (s *Service) createChunksFromGoElements(ctx context.Context, docID uuid.UUID, elements []GoCodeElement) ([]*storage.DocumentChunk, error) {
	var chunks []*storage.DocumentChunk
	chunkIndex := 0

	for _, element := range elements {
		// Each Go element becomes its own chunk for better code analysis
		chunk := &storage.DocumentChunk{
			ID:         uuid.New(),
			DocumentID: docID,
			ChunkIndex: chunkIndex,
			Content:    element.Content,
			Metadata:   make(map[string]interface{}),
		}

		chunks = append(chunks, chunk)
		chunkIndex++

		// Check context for cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}

	return chunks, nil
}

// storeChunks stores document chunks in the database using concurrent processing
func (s *Service) storeChunks(ctx context.Context, chunks []*storage.DocumentChunk) error {
	// Use worker pool for concurrent processing
	chunkChan := make(chan *storage.DocumentChunk, 100)
	errorChan := make(chan error, runtime.NumCPU())

	// Start workers
	for i := 0; i < runtime.NumCPU(); i++ {
			go func() {
				for chunk := range chunkChan {
					// Generate embedding for chunk content
					embedding, err := s.generateEmbedding(ctx, chunk.Content)
					if err != nil {
						s.logger.Warnf("Failed to generate embedding for chunk %d: %v", chunk.ChunkIndex, err)
						// Continue with nil embedding for now
					} else {
						chunk.Embedding = embedding
					}

					if err := s.db.CreateDocumentChunk(ctx, chunk); err != nil {
						errorChan <- fmt.Errorf("failed to store chunk %d: %w", chunk.ChunkIndex, err)
						return
					}
				}
				errorChan <- nil
			}()
		}

	// Send chunks to workers
	go func() {
		defer close(chunkChan)
		for _, chunk := range chunks {
			select {
			case chunkChan <- chunk:
			case <-ctx.Done():
				return
			}
		}
	}()

	// Wait for workers to complete
	for i := 0; i < runtime.NumCPU(); i++ {
		if err := <-errorChan; err != nil {
			return err
		}
	}

	s.logger.Infof("Successfully stored %d chunks", len(chunks))
	return nil
}

// generateEmbedding generates embedding for text content using Ollama
func (s *Service) generateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if s.ragService == nil {
		return nil, fmt.Errorf("rag service not available")
	}
	
	// Use the RAG service's AI client to generate embedding
	return s.ragService.GenerateEmbedding(ctx, text)
}

// float32ArrayToVectorString converts []float32 to PostgreSQL vector format string
func (s *Service) float32ArrayToVectorString(embedding []float32) string {
	if embedding == nil {
		return "[]"
	}
	
	strValues := make([]string, len(embedding))
	for i, val := range embedding {
		strValues[i] = strconv.FormatFloat(float64(val), 'f', -1, 64)
	}
	
	return "[" + strings.Join(strValues, ",") + "]"
}
