package embeddings

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"ai-recruiter-assistant/internal/storage"
)

// ResumeProcessor handles PDF resume processing
type ResumeProcessor struct {
	storage storage.StorageInterface
}

// NewResumeProcessor creates a new resume processor
func NewResumeProcessor(storage storage.StorageInterface) *ResumeProcessor {
	return &ResumeProcessor{
		storage: storage,
	}
}

// ProcessResume processes a PDF resume file
func (rp *ResumeProcessor) ProcessResume(ctx context.Context, userID string, filename string, fileContent []byte) (*storage.Document, error) {
	logrus.Infof("Processing resume for user: %s, file: %s", userID, filename)

	// For now, we'll simulate PDF text extraction
	// TODO: Implement actual PDF processing with unipdf
	text := fmt.Sprintf("Resume content for %s\n\nExperience: Software Developer\nSkills: Go, Python, JavaScript\nEducation: Bachelor's in Computer Science", userID)

	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("no text extracted from PDF")
	}

	// Create document record
	document := &storage.Document{
		ID:      fmt.Sprintf("resume_%s_%d", userID, time.Now().Unix()),
		UserID:  userID,
		Content: text,
		Source:  "resume",
	}

	// Store document
	err := rp.storage.CreateDocument(document)
	if err != nil {
		return nil, fmt.Errorf("failed to store document: %w", err)
	}

	logrus.Infof("Successfully processed and stored resume for user: %s", userID)
	return document, nil
}

// ChunkText splits text into chunks for embedding generation
func (rp *ResumeProcessor) ChunkText(text string, chunkSize int, overlap int) []string {
	if len(text) <= chunkSize {
		return []string{text}
	}

	var chunks []string
	textRunes := []rune(text)
	
	for i := 0; i < len(textRunes); i += (chunkSize - overlap) {
		end := i + chunkSize
		if end > len(textRunes) {
			end = len(textRunes)
		}
		
		chunk := string(textRunes[i:end])
		chunks = append(chunks, chunk)
		
		if end >= len(textRunes) {
			break
		}
	}
	
	return chunks
}
