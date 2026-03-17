package ai

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// AIService defines the interface for AI operations
type AIService interface {
	GenerateEmailResponse(ctx context.Context, request EmailRequest) (*EmailResponse, error)
	ValidateRequest(request EmailRequest) error
}

// EmailRequest represents a request to generate an email response
type EmailRequest struct {
	EmailBody string `json:"email_body"`
	Subject   string `json:"subject"`
	Sender    string `json:"sender"`
	Tone      string `json:"tone"`      // professional, casual, friendly, formal
	Length    string `json:"length"`    // short, medium, detailed
	UserID    string `json:"user_id"`
}

// EmailResponse represents the response from email generation
type EmailResponse struct {
	Response   string `json:"response"`
	Tone       string `json:"tone"`
	Length     string `json:"length"`
	TokensUsed int    `json:"tokens_used"`
	Model      string `json:"model"`
}

// DefaultAIService implements AIService using Ollama
type DefaultAIService struct {
	ollamaService *OllamaService
}

// NewAIService creates a new AI service
func NewAIService(ollamaService *OllamaService) AIService {
	return &DefaultAIService{
		ollamaService: ollamaService,
	}
}

// GenerateEmailResponse generates an AI response for an email
func (s *DefaultAIService) GenerateEmailResponse(ctx context.Context, request EmailRequest) (*EmailResponse, error) {
	// Validate request
	if err := s.ValidateRequest(request); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Generate prompt
	prompt := s.generatePrompt(request)

	logrus.Infof("Generating email response for user %s using model %s", request.UserID, s.ollamaService.genModel)

	// Generate response using Ollama
	startTime := time.Now()
	response, err := s.ollamaService.GenerateText(ctx, s.ollamaService.genModel, prompt)
	if err != nil {
		logrus.Errorf("Failed to generate response using Ollama: %v", err)
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}
	
	duration := time.Since(startTime)
	logrus.Infof("Email response generated in %v for user %s", duration, request.UserID)

	// Clean up the response
	cleanedResponse := s.cleanResponse(response.Response)

	return &EmailResponse{
		Response:   cleanedResponse,
		Tone:       request.Tone,
		Length:     request.Length,
		TokensUsed: response.PromptEval + response.EvalCount,
		Model:      response.Model,
	}, nil
}

// ValidateRequest validates the email request
func (s *DefaultAIService) ValidateRequest(request EmailRequest) error {
	if strings.TrimSpace(request.EmailBody) == "" {
		return fmt.Errorf("email body cannot be empty")
	}

	if request.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	// Validate tone
	validTones := map[string]bool{
		"professional": true,
		"casual":       true,
		"friendly":     true,
		"formal":       true,
	}
	if !validTones[request.Tone] {
		return fmt.Errorf("invalid tone: %s. Valid options: professional, casual, friendly, formal", request.Tone)
	}

	// Validate length
	validLengths := map[string]bool{
		"short":    true,
		"medium":   true,
		"detailed": true,
	}
	if !validLengths[request.Length] {
		return fmt.Errorf("invalid length: %s. Valid options: short, medium, detailed", request.Length)
	}

	return nil
}

// generatePrompt creates a prompt for the AI model based on the request
func (s *DefaultAIService) generatePrompt(request EmailRequest) string {
	toneInstructions := map[string]string{
		"professional": "Write in a professional and formal tone suitable for business communication. Use proper business etiquette.",
		"casual":       "Write in a casual and relaxed tone as if talking to a colleague. Be friendly but maintain professionalism.",
		"friendly":     "Write in a warm and friendly tone. Be approachable and personable while remaining appropriate.",
		"formal":       "Write in a very formal tone using proper business language and etiquette. Be respectful and diplomatic.",
	}

	lengthInstructions := map[string]string{
		"short":    "Keep the response concise and to the point. Aim for 2-3 sentences maximum.",
		"medium":   "Write a moderately detailed response with 3-5 sentences. Include key points but avoid unnecessary details.",
		"detailed": "Write a comprehensive response with multiple paragraphs. Include relevant details and explanations.",
	}

	toneInstr := toneInstructions[request.Tone]
	lengthInstr := lengthInstructions[request.Length]

	prompt := fmt.Sprintf(`You are an AI email assistant helping to draft a professional email response.

CONTEXT:
- Original Email Subject: %s
- Sender: %s
- Email Content: %s

INSTRUCTIONS:
1. %s
2. %s
3. Respond as if you are the recipient drafting a reply
4. Do not include any greetings like "Dear [Name]" or sign-offs like "Best regards"
5. Focus only on the email body content
6. Be helpful, professional, and appropriate
7. Address the key points from the original email
8. Ask for clarification if needed

Please generate the email response:`,
		request.Subject,
		request.Sender,
		request.EmailBody,
		toneInstr,
		lengthInstr,
	)

	return prompt
}

// cleanResponse cleans up the AI-generated response
func (s *DefaultAIService) cleanResponse(response string) string {
	// Remove common AI response artifacts
	cleaned := response

	// Remove leading/trailing whitespace
	cleaned = strings.TrimSpace(cleaned)

	// Remove common prefixes that AI models might add
	prefixesToRemove := []string{
		"Response:",
		"Email Response:",
		"Here is the response:",
		"Here's the response:",
		"The response is:",
	}

	for _, prefix := range prefixesToRemove {
		if strings.HasPrefix(strings.ToLower(cleaned), strings.ToLower(prefix)) {
			cleaned = strings.TrimSpace(cleaned[len(prefix):])
		}
	}

	// Remove quotation marks if the entire response is quoted
	if strings.HasPrefix(cleaned, `"`) && strings.HasSuffix(cleaned, `"`) {
		cleaned = strings.TrimSpace(cleaned[1 : len(cleaned)-1])
	}

	// Remove any remaining AI-like explanations
	lines := strings.Split(cleaned, "\n")
	var cleanedLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && 
		   !strings.HasPrefix(strings.ToLower(line), "here is") &&
		   !strings.HasPrefix(strings.ToLower(line), "the response") &&
		   !strings.HasPrefix(strings.ToLower(line), "this response") {
			cleanedLines = append(cleanedLines, line)
		}
	}

	cleaned = strings.Join(cleanedLines, "\n")
	return cleaned
}

// OllamaResponse represents the response from Ollama API
type OllamaResponse struct {
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	Model      string `json:"model"`
	TokensUsed int    `json:"total_tokens,omitempty"`
}

// Generate is a helper method to use OllamaService for generation
func (s *DefaultAIService) Generate(ctx context.Context, prompt string, model string) (*OllamaResponse, error) {
	// Use GenerateText method from OllamaService
	response, err := s.ollamaService.GenerateText(ctx, model, prompt)
	if err != nil {
		return nil, err
	}
	
	// Convert to OllamaResponse format
	return &OllamaResponse{
		Response:   response.Response,
		Model:      response.Model,
		Done:       response.Done,
		TokensUsed: response.PromptEval + response.EvalCount,
	}, nil
}
