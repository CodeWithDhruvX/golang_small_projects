package rag

import (
	"context"
	"fmt"
	"strings"
	"time"

	"ai-recruiter-assistant/internal/auth"
	"github.com/sirupsen/logrus"
	"ai-recruiter-assistant/internal/ai"
	"ai-recruiter-assistant/internal/storage"
)

// RAGService handles Retrieval-Augmented Generation operations
type RAGService struct {
	storage     storage.StorageInterface
	ollama      *ai.OllamaService
}

// NewRAGService creates a new RAG service
func NewRAGService(storage storage.StorageInterface, ollama *ai.OllamaService) *RAGService {
	return &RAGService{
		storage: storage,
		ollama:  ollama,
	}
}

// SearchContext performs semantic search to retrieve relevant context
func (rs *RAGService) SearchContext(ctx context.Context, userID, query string, topK int) ([]storage.Document, error) {
	logrus.Infof("Performing semantic search for user: %s, query: %s", userID, query)

	// Generate embedding for the query
	queryEmbedding, err := rs.ollama.GenerateEmbedding(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search for similar documents
	documents, err := rs.storage.SearchSimilar(queryEmbedding, userID, topK, "documents")
	if err != nil {
		return nil, fmt.Errorf("failed to search similar documents: %w", err)
	}

	logrus.Infof("Found %d relevant documents for query", len(documents))
	return documents, nil
}

// GenerateContextualReply generates a reply using retrieved context
func (rs *RAGService) GenerateContextualReply(ctx context.Context, userID, emailText string, model string) (string, error) {
	logrus.Infof("Generating contextual reply for user: %s", userID)

	// Use a timeout context for the entire operation
	ctx, cancel := context.WithTimeout(ctx, 45*time.Second) // 45s total timeout
	defer cancel()

	// For short emails, skip RAG to improve speed
	if len(strings.Fields(emailText)) < 50 {
		logrus.Infof("Short email detected, skipping RAG for speed")
		return rs.generateFastReply(ctx, userID, emailText, model)
	}

	// Retrieve relevant context with reduced top_k for speed
	contextDocs, err := rs.SearchContext(ctx, userID, emailText, 3) // Reduced from 5
	if err != nil {
		logrus.Warnf("Failed to retrieve context: %v", err)
		// Continue without context
		contextDocs = []storage.Document{}
	}

	// Get user profile information
	user, err := rs.storage.GetUserByID(userID)
	if err != nil {
		logrus.Warnf("Failed to get user profile: %v", err)
		// Create a default user profile
		user = &auth.User{
			Name: "Candidate",
			Experience: "Software Developer",
			Skills: []string{"Programming", "Problem Solving"},
		}
	}

	// Build context string
	contextStr := rs.buildContextString(contextDocs, user)

	// Generate reply using context
	reply, err := rs.ollama.GenerateReply(ctx, emailText, contextStr, model)
	if err != nil {
		logrus.Warnf("Failed to generate reply with AI, using template: %v", err)
		return rs.generateTemplateReply(user), nil
	}

	logrus.Infof("Generated contextual reply with %d context documents", len(contextDocs))
	return reply, nil
}

// generateFastReply generates a quick reply without RAG for short emails
func (rs *RAGService) generateFastReply(ctx context.Context, userID, emailText, model string) (string, error) {
	// Get user profile information
	user, err := rs.storage.GetUserByID(userID)
	if err != nil {
		logrus.Warnf("Failed to get user profile: %v", err)
		// Create a default user profile
		user = &auth.User{
			Name: "Candidate",
			Experience: "Software Developer",
			Skills: []string{"Programming", "Problem Solving"},
		}
	}

	// Build minimal context string
	contextStr := fmt.Sprintf("Name: %s\nExperience: %s", user.Name, user.Experience)

	// Generate reply with minimal context
	reply, err := rs.ollama.GenerateReply(ctx, emailText, contextStr, model)
	if err != nil {
		logrus.Warnf("Failed to generate fast reply: %v", err)
		return rs.generateTemplateReply(user), nil
	}

	logrus.Infof("Generated fast reply without RAG")
	return reply, nil
}

// generateTemplateReply creates a basic template reply as fallback
func (rs *RAGService) generateTemplateReply(user *auth.User) string {
	// Use user's actual name if available, otherwise use "Candidate"
	candidateName := user.Name
	if candidateName == "" {
		candidateName = "Candidate"
	}

	// Build experience/skills context
	experienceContext := ""
	if user.Experience != "" {
		experienceContext = fmt.Sprintf("With my background as a %s, ", user.Experience)
	}

	skillsContext := ""
	if len(user.Skills) > 0 {
		skillsContext = fmt.Sprintf(" I bring expertise in areas such as %s,", strings.Join(user.Skills, ", "))
	}

	return fmt.Sprintf(`Dear Recruiter,

Thank you for reaching out to me regarding this exciting opportunity. I am genuinely interested in learning more about this position and how my background aligns with your requirements.

%s%s I believe I could be a strong fit for this role and would welcome the chance to discuss my qualifications in more detail.

I am available for a conversation at your convenience and would be happy to provide any additional information you may need, including my resume, portfolio, or references.

Could you please suggest a suitable time for us to connect? I look forward to hearing from you.

Best regards,
%s`, experienceContext, skillsContext, candidateName)
}

// buildContextString builds a context string from documents and user info
func (rs *RAGService) buildContextString(documents []storage.Document, user *auth.User) string {
	var context strings.Builder

	// Add user profile information
	context.WriteString("Candidate Profile:\n")
	context.WriteString(fmt.Sprintf("Name: %s\n", user.Name))
	if user.Experience != "" {
		context.WriteString(fmt.Sprintf("Experience: %s\n", user.Experience))
	}
	if len(user.Skills) > 0 {
		context.WriteString(fmt.Sprintf("Skills: %s\n", strings.Join(user.Skills, ", ")))
	}
	if user.ExpectedSalary > 0 {
		context.WriteString(fmt.Sprintf("Expected Salary: %.2f\n", user.ExpectedSalary))
	}
	if user.NoticePeriod > 0 {
		context.WriteString(fmt.Sprintf("Notice Period: %d days\n", user.NoticePeriod))
	}
	if user.Location != "" {
		context.WriteString(fmt.Sprintf("Location: %s\n", user.Location))
	}
	context.WriteString("\n")

	// Add relevant document excerpts
	if len(documents) > 0 {
		context.WriteString("Relevant Information:\n")
		for i, doc := range documents {
			context.WriteString(fmt.Sprintf("%d. %s\n", i+1, doc.Content))
			if i < len(documents)-1 {
				context.WriteString("\n")
			}
		}
	}

	return context.String()
}

// ProcessAndStoreDocument processes a document and stores it with embeddings
func (rs *RAGService) ProcessAndStoreDocument(ctx context.Context, document *storage.Document) error {
	logrus.Infof("Processing document for user: %s, source: %s", document.UserID, document.Source)

	// Generate embedding for the document content
	embedding, err := rs.ollama.GenerateEmbedding(ctx, document.Content)
	if err != nil {
		return fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Store embedding with the document
	document.Embedding = embedding
	err = rs.storage.StoreEmbedding(document.ID, embedding, "documents")
	if err != nil {
		return fmt.Errorf("failed to store embedding: %w", err)
	}

	logrus.Infof("Successfully processed and stored document: %s", document.ID)
	return nil
}

// ProcessResume processes resume text and stores it with embeddings
func (rs *RAGService) ProcessResume(ctx context.Context, userID, resumeText string) error {
	logrus.Infof("Processing resume for user: %s", userID)

	// Create document chunks for better retrieval
	chunks := rs.chunkText(resumeText, 1000, 200) // 1000 chars with 200 overlap

	for i, chunk := range chunks {
		document := &storage.Document{
			ID:      fmt.Sprintf("resume_chunk_%s_%d", userID, i),
			UserID:  userID,
			Content: chunk,
			Source:  "resume",
		}

		// Process and store each chunk
		err := rs.ProcessAndStoreDocument(ctx, document)
		if err != nil {
			logrus.Errorf("Failed to process resume chunk %d: %v", i, err)
			continue
		}
	}

	logrus.Infof("Successfully processed resume into %d chunks for user: %s", len(chunks), userID)
	return nil
}

// chunkText splits text into chunks with overlap
func (rs *RAGService) chunkText(text string, chunkSize, overlap int) []string {
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

// GetCandidateInfo retrieves formatted candidate information for AI generation
func (rs *RAGService) GetCandidateInfo(ctx context.Context, userID string) (string, error) {
	user, err := rs.storage.GetUserByID(userID)
	if err != nil {
		return "", fmt.Errorf("failed to get user profile: %w", err)
	}

	var info strings.Builder
	info.WriteString(fmt.Sprintf("Name: %s\n", user.Name))
	
	if user.Experience != "" {
		info.WriteString(fmt.Sprintf("Experience: %s\n", user.Experience))
	}
	
	if len(user.Skills) > 0 {
		info.WriteString(fmt.Sprintf("Skills: %s\n", strings.Join(user.Skills, ", ")))
	}
	
	if user.CurrentSalary > 0 {
		info.WriteString(fmt.Sprintf("Current Salary: %.2f\n", user.CurrentSalary))
	}
	
	if user.ExpectedSalary > 0 {
		info.WriteString(fmt.Sprintf("Expected Salary: %.2f\n", user.ExpectedSalary))
	}
	
	if user.NoticePeriod > 0 {
		info.WriteString(fmt.Sprintf("Notice Period: %d days\n", user.NoticePeriod))
	}
	
	if user.Location != "" {
		info.WriteString(fmt.Sprintf("Location: %s\n", user.Location))
	}
	
	if user.LinkedInURL != "" {
		info.WriteString(fmt.Sprintf("LinkedIn: %s\n", user.LinkedInURL))
	}
	
	if user.GitHubURL != "" {
		info.WriteString(fmt.Sprintf("GitHub: %s\n", user.GitHubURL))
	}

	return info.String(), nil
}
