package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/open-pm/open-pm/server/internal/config"
)

// LLMMessage represents a chat message for LLM providers.
type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMProvider defines the interface for LLM chat providers.
type LLMProvider interface {
	Chat(ctx context.Context, systemPrompt string, messages []LLMMessage) (string, error)
}

// NewLLMProvider creates an LLM provider based on configuration.
func NewLLMProvider(cfg *config.LLMConfig) LLMProvider {
	client := &http.Client{Timeout: 60 * time.Second}
	switch cfg.Provider {
	case "anthropic":
		model := cfg.Model
		if model == "" || model == "gpt-4o-mini" {
			model = "claude-sonnet-4-20250514"
		}
		return &AnthropicProvider{
			apiKey:    cfg.APIKey,
			model:     model,
			maxTokens: cfg.MaxTokens,
			client:    client,
		}
	default:
		model := cfg.Model
		if model == "" {
			model = "gpt-4o-mini"
		}
		return &OpenAIProvider{
			apiKey:    cfg.APIKey,
			model:     model,
			maxTokens: cfg.MaxTokens,
			client:    client,
		}
	}
}

// --- OpenAI Provider ---

type OpenAIProvider struct {
	apiKey    string
	model     string
	maxTokens int
	client    *http.Client
}

type openAIRequest struct {
	Model     string           `json:"model"`
	Messages  []openAIMessage  `json:"messages"`
	MaxTokens int              `json:"max_tokens,omitempty"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *OpenAIProvider) Chat(ctx context.Context, systemPrompt string, messages []LLMMessage) (string, error) {
	oaiMessages := []openAIMessage{
		{Role: "system", Content: systemPrompt},
	}
	for _, m := range messages {
		oaiMessages = append(oaiMessages, openAIMessage{Role: m.Role, Content: m.Content})
	}

	reqBody := openAIRequest{
		Model:     p.model,
		Messages:  oaiMessages,
		MaxTokens: p.maxTokens,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var oaiResp openAIResponse
	if err := json.Unmarshal(respBody, &oaiResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if oaiResp.Error != nil {
		return "", fmt.Errorf("OpenAI error: %s", oaiResp.Error.Message)
	}

	if len(oaiResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return oaiResp.Choices[0].Message.Content, nil
}

// --- Anthropic Provider ---

type AnthropicProvider struct {
	apiKey    string
	model     string
	maxTokens int
	client    *http.Client
}

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *AnthropicProvider) Chat(ctx context.Context, systemPrompt string, messages []LLMMessage) (string, error) {
	var antMessages []anthropicMessage
	for _, m := range messages {
		antMessages = append(antMessages, anthropicMessage{Role: m.Role, Content: m.Content})
	}

	maxTokens := p.maxTokens
	if maxTokens <= 0 {
		maxTokens = 1024
	}

	reqBody := anthropicRequest{
		Model:     p.model,
		MaxTokens: maxTokens,
		System:    systemPrompt,
		Messages:  antMessages,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var antResp anthropicResponse
	if err := json.Unmarshal(respBody, &antResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if antResp.Error != nil {
		return "", fmt.Errorf("Anthropic error: %s", antResp.Error.Message)
	}

	if len(antResp.Content) == 0 {
		return "", fmt.Errorf("no response from Anthropic")
	}

	return antResp.Content[0].Text, nil
}
