package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClaudeProvider struct {
	APIKey string
	Model  string
}

func NewClaudeProvider(apiKey string, model string) *ClaudeProvider {
	if model == "" {
		model = "claude-3-sonnet-20240229"
	}

	return &ClaudeProvider{
		APIKey: apiKey,
		Model:  model,
	}
}

type ClaudeRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens"`
}

type ClaudeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

func (p *ClaudeProvider) GenerateSuggestions(diff string, includeBody bool, style string) ([]string, error) {
	// Build the prompt
	prompt := buildPrompt(diff, includeBody, style)

	// Create the request
	reqBody := ClaudeRequest{
		Model: p.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 2000,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("User-Agent", "zeus-ai/1.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned error: %s", string(respBody))
	}

	var respObj ClaudeResponse
	err = json.Unmarshal(respBody, &respObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(respObj.Content) == 0 {
		return nil, fmt.Errorf("API returned no content")
	}

	// Parse the response into individual suggestions
	var content string
	for _, chunk := range respObj.Content {
		if chunk.Type == "text" {
			content += chunk.Text
		}
	}
	return parseSuggestions(content), nil
}
