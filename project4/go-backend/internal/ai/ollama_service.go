package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// OllamaService handles communication with Ollama API
type OllamaService struct {
	baseURL    string
	genModel   string
	httpClient *http.Client
	gpuEnabled bool
	gpuType    string
}

// NewOllamaService creates a new Ollama service
func NewOllamaService(baseURL string) *OllamaService {
	service := &OllamaService{
		baseURL:  baseURL,
		genModel: "llama3.1:8b",
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // Reduced from 120s for faster failure
		},
	}
	
	// Detect GPU availability
	service.detectGPU()
	
	// Preload models for better performance
	service.preloadModels()
	
	return service
}

// detectGPU checks for GPU availability and sets GPU configuration
func (service *OllamaService) detectGPU() {
	service.gpuEnabled = false
	service.gpuType = "cpu"
	
	// Check if GPU usage is explicitly disabled
	if strings.ToLower(os.Getenv("DISABLE_GPU")) == "true" {
		logrus.Info("GPU usage explicitly disabled via DISABLE_GPU environment variable")
		return
	}
	
	// Check for NVIDIA GPU
	if service.checkNvidiaGPU() {
		service.gpuEnabled = true
		service.gpuType = "nvidia"
		logrus.Info("NVIDIA GPU detected and enabled for AI acceleration")
		return
	}
	
	// Check for AMD GPU on Linux
	if runtime.GOOS == "linux" && service.checkAMDGPU() {
		service.gpuEnabled = true
		service.gpuType = "amd"
		logrus.Info("AMD GPU detected and enabled for AI acceleration")
		return
	}
	
	// Check for Apple Silicon GPU
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		service.gpuEnabled = true
		service.gpuType = "apple_silicon"
		logrus.Info("Apple Silicon GPU detected and enabled for AI acceleration")
		return
	}
	
	logrus.Info("No compatible GPU detected, using CPU for AI processing")
}

// preloadModels preloads commonly used models to improve cold start performance
func (service *OllamaService) preloadModels() {
	logrus.Info("Preloading models for optimal performance...")
	
	// First check if Ollama is running
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Reduced from 5s
	defer cancel()
	
	// Test connection to Ollama
	url := fmt.Sprintf("%s/api/tags", service.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logrus.Warnf("Failed to create Ollama health check request: %v", err)
		return
	}
	
	resp, err := service.httpClient.Do(req)
	if err != nil {
		logrus.Warnf("Ollama is not running, skipping model preloading: %v", err)
		return
	}
	resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		logrus.Warnf("Ollama health check failed with status %d, skipping model preloading", resp.StatusCode)
		return
	}
	
	models := []string{"llama3.1:8b", "nomic-embed-text"}
	
	for _, model := range models {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Reduced from 30s
		
		// Create a minimal request to trigger model loading
		request := GenerateRequest{
			Model:  model,
			Prompt: "hi", // Shorter prompt for faster preloading
			Stream: false,
			Options: map[string]interface{}{
				"temperature": 0.1,
				"max_tokens":  5, // Reduced from 10 for faster preloading
			},
		}
		
		reqBody, err := json.Marshal(request)
		if err != nil {
			logrus.Warnf("Failed to marshal preload request for %s: %v", model, err)
			cancel()
			continue
		}
		
		url := fmt.Sprintf("%s/api/generate", service.baseURL)
		req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
		if err != nil {
			logrus.Warnf("Failed to create preload request for %s: %v", model, err)
			cancel()
			continue
		}
		
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := service.httpClient.Do(req)
		if err != nil {
			logrus.Warnf("Failed to preload model %s: %v", model, err)
			cancel()
			continue
		}
		resp.Body.Close()
		
		if resp.StatusCode == http.StatusOK {
			logrus.Infof("Successfully preloaded model: %s", model)
		} else {
			logrus.Warnf("Failed to preload model %s, status: %d", model, resp.StatusCode)
		}
		
		cancel()
	}
	
	logrus.Info("Model preloading completed")
}

// checkNvidiaGPU checks for NVIDIA GPU availability
func (service *OllamaService) checkNvidiaGPU() bool {
	// Check if nvidia-smi is available
	cmd := exec.Command("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	if strings.TrimSpace(string(output)) != "" {
		// Test if Ollama can access GPU
		return service.testOllamaGPU()
	}
	
	return false
}

// checkAMDGPU checks for AMD GPU availability on Linux
func (service *OllamaService) checkAMDGPU() bool {
	// Check for ROCm tools
	cmd := exec.Command("rocm-smi", "--showproductname")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	
	return strings.Contains(string(output), "GPU")
}

// testOllamaGPU tests if Ollama can actually use the GPU
func (service *OllamaService) testOllamaGPU() bool {
	// Create a simple test request to check if GPU is working
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	testRequest := GenerateRequest{
		Model:  service.genModel,
		Prompt: "Test",
		Stream: false,
	}
	
	reqBody, err := json.Marshal(testRequest)
	if err != nil {
		return false
	}
	
	url := fmt.Sprintf("%s/api/generate", service.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return false
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := service.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == http.StatusOK
}

// GetGPUStatus returns the current GPU status
func (service *OllamaService) GetGPUStatus() map[string]interface{} {
	return map[string]interface{}{
		"enabled":   service.gpuEnabled,
		"gpu_type":  service.gpuType,
		"detected":  service.gpuEnabled,
	}
}

// EmbeddingRequest represents the request for embedding generation
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse represents the response from embedding generation
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// GenerateRequest represents the request for text generation
type GenerateRequest struct {
	Model     string            `json:"model"`
	Prompt    string            `json:"prompt"`
	Stream    bool              `json:"stream"`
	Options   map[string]interface{} `json:"options,omitempty"`
}

// GenerateResponse represents the response from text generation
type GenerateResponse struct {
	Response   string `json:"response"`
	Done       bool   `json:"done"`
	Model      string `json:"model"`
	CreatedAt  string `json:"created_at"`
	TotalDur   int64  `json:"total_duration,omitempty"`
	LoadDur    int64  `json:"load_duration,omitempty"`
	PromptEval int    `json:"prompt_eval_count,omitempty"`
	EvalCount  int    `json:"eval_count,omitempty"`
	EvalDur    int64  `json:"eval_duration,omitempty"`
}

// GenerateEmbedding generates embeddings for text using Nomic Embed Text
func (service *OllamaService) GenerateEmbedding(ctx context.Context, text string) ([]float64, error) {
	startTime := time.Now()
	
	request := EmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/embeddings", service.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := service.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	logrus.Infof("Generated embedding in %v, dimensions: %d", time.Since(startTime), len(response.Embedding))
	return response.Embedding, nil
}

// GenerateText generates text using specified model
func (service *OllamaService) GenerateText(ctx context.Context, model, prompt string) (*GenerateResponse, error) {
	return service.GenerateTextWithOptions(ctx, model, prompt, map[string]interface{}{})
}

// GenerateTextStream generates text using streaming for better UX
func (service *OllamaService) GenerateTextStream(ctx context.Context, model, prompt string, options map[string]interface{}) (<-chan string, error) {
	startTime := time.Now()
	
	// Log GPU usage
	gpuStatus := "CPU"
	if service.gpuEnabled {
		gpuStatus = fmt.Sprintf("GPU (%s)", service.gpuType)
	}
	
	logrus.Infof("Starting streaming text generation using %s with model: %s", gpuStatus, model)
	
	// Default optimization options
	defaultOptions := map[string]interface{}{
		"temperature": 0.7,
		"top_p":       0.9,
		"max_tokens":  120, // Optimized for streaming
		"stream":      true, // Enable streaming
	}
	
	// Merge custom options with defaults
	for k, v := range options {
		defaultOptions[k] = v
	}
	
	request := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
		Options: defaultOptions,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", service.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	// Add GPU usage header
	if service.gpuEnabled {
		req.Header.Set("X-GPU-Accelerated", "true")
		req.Header.Set("X-GPU-Type", service.gpuType)
	}

	resp, err := service.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Create channel for streaming responses
	ch := make(chan string, 10)
	
	go func() {
		defer resp.Body.Close()
		defer close(ch)
		
		decoder := json.NewDecoder(resp.Body)
		for {
			var streamResp GenerateResponse
			if err := decoder.Decode(&streamResp); err != nil {
				if err == io.EOF {
					break
				}
				logrus.Warnf("Stream decode error: %v", err)
				break
			}
			
			if streamResp.Response != "" {
				ch <- streamResp.Response
			}
			
			if streamResp.Done {
				duration := time.Since(startTime)
				logrus.Infof("Streaming completed in %v using %s", duration, gpuStatus)
				break
			}
		}
	}()

	return ch, nil
}
func (service *OllamaService) GenerateTextWithOptions(ctx context.Context, model, prompt string, options map[string]interface{}) (*GenerateResponse, error) {
	startTime := time.Now()
	
	// Log GPU usage
	gpuStatus := "CPU"
	if service.gpuEnabled {
		gpuStatus = fmt.Sprintf("GPU (%s)", service.gpuType)
	}
	
	logrus.Infof("Generating text using %s with model: %s", gpuStatus, model)
	
	// Default optimization options based on Project 3 results
	defaultOptions := map[string]interface{}{
		"temperature": 0.7,
		"top_p":       0.9,
		"max_tokens":  150, // Reduced from 300 for faster responses
		"num_predict": 150, // Explicit token limit for some models
	}
	
	// Merge custom options with defaults
	for k, v := range options {
		defaultOptions[k] = v
	}
	
	request := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Options: defaultOptions,
	}

	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", service.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	
	// Add GPU usage header for monitoring
	if service.gpuEnabled {
		req.Header.Set("X-GPU-Accelerated", "true")
		req.Header.Set("X-GPU-Type", service.gpuType)
	}

	resp, err := service.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		logrus.Errorf("Ollama API error: status %d, url: %s, request: %s, body: %s", resp.StatusCode, url, string(reqBody), string(body))
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	duration := time.Since(startTime)
	
	// Enhanced performance logging
	loadTime := time.Duration(response.LoadDur) * time.Nanosecond
	evalTime := time.Duration(response.EvalDur) * time.Nanosecond
	
	logrus.Infof("Generated text in %v (load: %v, eval: %v) using %s with model: %s, tokens: %d", 
		duration, loadTime, evalTime, gpuStatus, model, response.EvalCount)
	
	return &response, nil
}

// ClassifyEmail classifies whether an email is from a recruiter
func (service *OllamaService) ClassifyEmail(ctx context.Context, emailText string) (bool, float64, error) {
	prompt := fmt.Sprintf(`Analyze the following email and determine if it is from a recruiter or related to a job opportunity. 
Respond ONLY in JSON format: {"is_recruiter": true/false, "confidence": 0.0-1.0}

Email:
%s`, emailText)

	// Use optimized options for classification
	classificationOptions := map[string]interface{}{
		"temperature": 0.1, // Low temperature for consistent classification
		"max_tokens":  50,  // Short responses for classification
	}

	resp, err := service.GenerateTextWithOptions(ctx, service.genModel, prompt, classificationOptions)
	if err != nil {
		logrus.Warnf("AI classification failed: %v, using fallback", err)
		return service.fallbackClassification(emailText)
	}

	var result struct {
		IsRecruiter bool    `json:"is_recruiter"`
		Confidence  float64 `json:"confidence"`
	}
	
	// Try to find JSON in the response if it's wrapped in text
	cleanedResponse := resp.Response
	if strings.Contains(cleanedResponse, "{") {
		start := strings.Index(cleanedResponse, "{")
		end := strings.LastIndex(cleanedResponse, "}") + 1
		cleanedResponse = cleanedResponse[start:end]
	}

	if err := json.Unmarshal([]byte(cleanedResponse), &result); err != nil {
		logrus.Warnf("Failed to parse AI classification JSON: %v, using fallback", err)
		return service.fallbackClassification(emailText)
	}

	return result.IsRecruiter, result.Confidence, nil
}

// fallbackClassification provides basic keyword-based classification as fallback
func (service *OllamaService) fallbackClassification(emailText string) (bool, float64, error) {
	text := strings.ToLower(emailText)
	
	recruiterKeywords := []string{
		"job opportunity", "position", "role", "hiring", "recruit", "resume", "cv",
		"interview", "salary", "compensation", "candidate", "talent", "vacancy",
		"opening", "employment", "career", "apply", "application",
	}
	
	keywordCount := 0
	for _, keyword := range recruiterKeywords {
		if strings.Contains(text, keyword) {
			keywordCount++
		}
	}
	
	confidence := float64(keywordCount) / float64(len(recruiterKeywords))
	isRecruiter := keywordCount >= 2
	
	return isRecruiter, confidence, nil
}

// ExtractRequirements extracts requested information from recruiter email
func (service *OllamaService) ExtractRequirements(ctx context.Context, emailText string) (map[string]bool, error) {
	prompt := fmt.Sprintf(`Extract the information requested by the recruiter in this email. 
Check for: resume, experience details, expected salary (CTC), notice period, specific skills, portfolio/GitHub, cover letter, and availability for a call.
Respond ONLY in JSON format: {"resume": true/false, "experience": true/false, ...}

Email:
%s`, emailText)

	// Use optimized options for extraction
	extractionOptions := map[string]interface{}{
		"temperature": 0.2, // Low temperature for consistent extraction
		"max_tokens":  100, // Moderate length for JSON response
	}

	resp, err := service.GenerateTextWithOptions(ctx, service.genModel, prompt, extractionOptions)
	if err != nil {
		logrus.Warnf("AI extraction failed: %v, using fallback", err)
		return service.fallbackRequirementExtraction(emailText)
	}

	var requirements map[string]bool
	
	cleanedResponse := resp.Response
	if strings.Contains(cleanedResponse, "{") {
		start := strings.Index(cleanedResponse, "{")
		end := strings.LastIndex(cleanedResponse, "}") + 1
		cleanedResponse = cleanedResponse[start:end]
	}

	if err := json.Unmarshal([]byte(cleanedResponse), &requirements); err != nil {
		logrus.Warnf("Failed to parse AI extraction JSON: %v, using fallback", err)
		return service.fallbackRequirementExtraction(emailText)
	}

	return requirements, nil
}

// fallbackRequirementExtraction provides basic keyword-based extraction as fallback
func (service *OllamaService) fallbackRequirementExtraction(emailText string) (map[string]bool, error) {
	text := strings.ToLower(emailText)
	
	requirements := map[string]bool{
		"resume":        strings.Contains(text, "resume") || strings.Contains(text, "cv"),
		"experience":    strings.Contains(text, "experience") || strings.Contains(text, "years"),
		"expected_ctc":  strings.Contains(text, "salary") || strings.Contains(text, "ctc") || strings.Contains(text, "compensation"),
		"notice_period": strings.Contains(text, "notice") || strings.Contains(text, "joining"),
		"skills":        strings.Contains(text, "skills") || strings.Contains(text, "technologies"),
		"portfolio":     strings.Contains(text, "portfolio") || strings.Contains(text, "github"),
		"cover_letter":  strings.Contains(text, "cover letter"),
		"availability":  strings.Contains(text, "available") || strings.Contains(text, "join"),
	}
	
	return requirements, nil
}
	

// GenerateReply generates a professional reply to recruiter email
func (service *OllamaService) GenerateReply(ctx context.Context, emailText, candidateInfo string, model string) (string, error) {
	// Smart prompt optimization based on email complexity
	emailComplexity := service.analyzeEmailComplexity(emailText)
	
	var prompt string
	switch emailComplexity {
	case "simple":
		prompt = fmt.Sprintf(`Write a concise professional reply to this recruiter email.
CANDIDATE: %s
EMAIL: %s
Requirements: Professional, 2-3 sentences max, express interest, suggest next step.`, candidateInfo, emailText)
	case "medium":
		prompt = fmt.Sprintf(`Write a professional reply to this recruiter email.
CANDIDATE PROFILE: %s
RECRUITER EMAIL: %s
Structure: Greeting + Interest + Specific Response + Next Steps + Closing
Keep concise but comprehensive.`, candidateInfo, emailText)
	default: // complex
		prompt = fmt.Sprintf(`You are an AI assistant helping a job candidate craft a professional reply to a recruiter's email.

CONTEXT AND GUIDELINES:
- You are writing on behalf of a job candidate responding to a recruiter
- The tone should be professional, confident, and enthusiastic but not desperate
- Keep reply concise (1-2 paragraphs maximum) and easy to read
- Always express genuine interest in opportunity
- Include a clear call-to-action (next steps)
- Address any specific questions or requirements mentioned in the recruiter's email
- Highlight relevant skills/experience that match the opportunity if mentioned

CANDIDATE PROFILE AND CONTEXT:
%s

RECRUITER'S EMAIL:
%s

REPLY STRUCTURE:
1. Professional greeting and thank you
2. Brief expression of interest and skills alignment
3. Response to specific questions (if any)
4. Clear next steps
5. Professional closing

Write a concise reply (max 100 words):`, candidateInfo, emailText)
	}

	targetModel := model
	if targetModel == "" {
		// Smart model selection based on email complexity
		switch emailComplexity {
		case "simple":
			targetModel = "phi3:latest" // Faster model for simple emails
		case "medium":
			targetModel = "qwen2.5-coder:3b" // Balanced model
		default:
			targetModel = "llama3.1:8b" // Best quality for complex emails
		}
	}

	logrus.Infof("Generating reply using model: %s, complexity: %s", targetModel, emailComplexity)

	// Use optimized options for faster, concise replies
	replyOptions := map[string]interface{}{
		"temperature": 0.6,
		"top_p":       0.85,
		"max_tokens":  120, // Reduced from 200 for much faster replies
		"num_predict": 120, // Explicit token limit
		"repeat_penalty": 1.1, // Reduce repetition
	}

	resp, err := service.GenerateTextWithOptions(ctx, targetModel, prompt, replyOptions)
	if err != nil {
		logrus.Errorf("AI reply generation failed (context: %v): %v", ctx.Err(), err)
		return service.generateTemplateReply(emailText, candidateInfo), nil
	}

	return resp.Response, nil
}

// analyzeEmailComplexity determines email complexity for prompt optimization
func (service *OllamaService) analyzeEmailComplexity(emailText string) string {
	wordCount := len(strings.Fields(emailText))
	
	// Check for complexity indicators
	complexityIndicators := []string{
		"salary", "compensation", "ctc", "requirements", "qualifications",
		"experience", "skills", "responsibilities", "deadline", "interview process",
	}
	
	indicatorCount := 0
	lowerText := strings.ToLower(emailText)
	for _, indicator := range complexityIndicators {
		if strings.Contains(lowerText, indicator) {
			indicatorCount++
		}
	}
	
	// Determine complexity
	if wordCount < 50 && indicatorCount <= 1 {
		return "simple"
	} else if wordCount < 150 && indicatorCount <= 3 {
		return "medium"
	}
	return "complex"
}

// generateTemplateReply creates a basic template reply as fallback
func (service *OllamaService) generateTemplateReply(emailText, candidateInfo string) string {
	// Extract candidate name from info if available
	candidateName := "Candidate"
	if strings.Contains(candidateInfo, "Name:") {
		nameStart := strings.Index(candidateInfo, "Name:") + 5
		nameEnd := strings.Index(candidateInfo[nameStart:], "\n")
		if nameEnd > 0 {
			candidateName = strings.TrimSpace(candidateInfo[nameStart : nameStart+nameEnd])
		}
	}

	return fmt.Sprintf(`Dear Recruiter,

Thank you for reaching out to me regarding this exciting opportunity. I am genuinely interested in learning more about this position and how my background aligns with your requirements.

Based on my experience and skills, I believe I could be a strong fit for this role. I would welcome the chance to discuss this opportunity further and share more details about my qualifications.

I am available for a conversation at your convenience and would be happy to provide any additional information you may need, including my resume, portfolio, or references.

Could you please suggest a suitable time for us to connect? I look forward to hearing from you.

Best regards,
%s`, candidateName)
}
