package rag

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// OllamaClient handles communication with Ollama API
type OllamaClient struct {
	baseURL string
	logger  *logrus.Logger
	client  *http.Client
}

// NewOllamaClient creates a new Ollama client
func NewOllamaClient(baseURL string, logger *logrus.Logger) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		logger:  logger,
		client: &http.Client{
			Timeout: 120 * time.Second, // 2 minute timeout for AI responses
		},
	}
}

// GenerateResponse generates a response from the AI model
func (c *OllamaClient) GenerateResponse(ctx context.Context, systemPrompt, userPrompt, model string) (*AIResponse, error) {
	if model == "" {
		model = "llama3.1:8b"
	}

	// Build the prompt
	fullPrompt := fmt.Sprintf("%s\n\nUser: %s\nAssistant: ", systemPrompt, userPrompt)

	// Prepare request
	requestBody := GenerateRequest{
		Model:  model,
		Prompt: fullPrompt,
		Stream: false,
		Options: Options{
			Temperature: 0.7,
			TopP:        0.9,
			MaxTokens:   300,
		},
	}

	// Send request
	resp, err := c.sendRequest(ctx, "/api/generate", requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var generateResp GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&generateResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &AIResponse{
		Content:   generateResp.Response,
		ModelUsed: model,
		TokensUsed: generateResp.PromptEvalCount + generateResp.EvalCount,
	}, nil
}

// StreamResponse streams a response from the AI model
func (c *OllamaClient) StreamResponse(ctx context.Context, systemPrompt, userPrompt, model string) (<-chan AIStreamChunk, error) {
	chunkChan := make(chan AIStreamChunk, 10)

	if model == "" {
		model = "llama3.1:8b"
	}

	// Build the prompt
	fullPrompt := fmt.Sprintf("%s\n\nUser: %s\nAssistant: ", systemPrompt, userPrompt)

	// Prepare request
	requestBody := GenerateRequest{
		Model:  model,
		Prompt: fullPrompt,
		Stream: true,
		Options: Options{
			Temperature: 0.7,
			TopP:        0.9,
			MaxTokens:   300,
		},
	}

	go func() {
		defer close(chunkChan)

		resp, err := c.sendRequest(ctx, "/api/generate", requestBody)
		if err != nil {
			chunkChan <- AIStreamChunk{
				Error: fmt.Sprintf("Failed to send request: %v", err),
			}
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			chunkChan <- AIStreamChunk{
				Error: fmt.Sprintf("API request failed with status %d: %s", resp.StatusCode, string(body)),
			}
			return
		}

		// Stream response
		decoder := json.NewDecoder(resp.Body)
		for {
			var chunk GenerateResponse
			if err := decoder.Decode(&chunk); err != nil {
				if err == io.EOF {
					break
				}
				chunkChan <- AIStreamChunk{
					Error: fmt.Sprintf("Failed to decode chunk: %v", err),
				}
				return
			}

			// Check if this is the final chunk
			if chunk.Done {
				break
			}

			chunkChan <- AIStreamChunk{
				Content:   chunk.Response,
				Model:     model,
				TokensUsed: chunk.EvalCount,
			}

			// Check context for cancellation
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return chunkChan, nil
}

// GenerateEmbedding generates embeddings for text using Nomic-Embed-Text
func (c *OllamaClient) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Prepare request
	requestBody := EmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	// Send request
	resp, err := c.sendRequest(ctx, "/api/embeddings", requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to send embedding request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("embedding request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %w", err)
	}

	return embeddingResp.Embedding, nil
}

// ListModels lists available models
func (c *OllamaClient) ListModels(ctx context.Context) ([]Model, error) {
	resp, err := c.sendRequest(ctx, "/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list models request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var listResp ListModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResp); err != nil {
		return nil, fmt.Errorf("failed to decode models response: %w", err)
	}

	return listResp.Models, nil
}

// PullModel pulls a model from Ollama
func (c *OllamaClient) PullModel(ctx context.Context, modelName string) error {
	requestBody := PullRequest{
		Name: modelName,
	}

	resp, err := c.sendRequest(ctx, "/api/pull", requestBody)
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("pull model request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Stream the pull progress
	decoder := json.NewDecoder(resp.Body)
	for {
		var progress PullProgress
		if err := decoder.Decode(&progress); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode pull progress: %w", err)
		}

		if progress.Status == "success" {
			c.logger.Infof("Successfully pulled model: %s", modelName)
			break
		}

		c.logger.Infof("Pulling model %s: %s", modelName, progress.Status)

		// Check context for cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}

	return nil
}

// sendRequest sends a request to the Ollama API
func (c *OllamaClient) sendRequest(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.client.Do(req)
}

// API request/response structures

type GenerateRequest struct {
	Model    string  `json:"model"`
	Prompt   string  `json:"prompt"`
	Stream   bool    `json:"stream"`
	Options  Options `json:"options"`
}

type Options struct {
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	MaxTokens   int     `json:"num_predict"`
}

type GenerateResponse struct {
	Model              string  `json:"model"`
	CreatedAt          string  `json:"created_at"`
	Response           string  `json:"response"`
	Done               bool    `json:"done"`
	TotalDuration      int64   `json:"total_duration"`
	PromptEvalCount    int     `json:"prompt_eval_count"`
	PromptEvalDuration int64   `json:"prompt_eval_duration"`
	EvalCount          int     `json:"eval_count"`
	EvalDuration       int64   `json:"eval_duration"`
}

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

type ListModelsResponse struct {
	Models []Model `json:"models"`
}

type Model struct {
	Name       string    `json:"name"`
	ModifiedAt time.Time `json:"modified_at"`
	Size       int64     `json:"size"`
}

type PullRequest struct {
	Name string `json:"name"`
}

type PullProgress struct {
	Status   string `json:"status"`
	Digest   string `json:"digest,omitempty"`
	Total    int64  `json:"total,omitempty"`
	Download int64  `json:"download,omitempty"`
}

type AIStreamChunk struct {
	Content   string `json:"content"`
	Model     string `json:"model"`
	TokensUsed int   `json:"tokens_used"`
	Error     string `json:"error,omitempty"`
}
