package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OpenRouterProvider struct {
	APIKey string
	Model  string
}

func NewOpenRouterProvider(apiKey string, model string) *OpenRouterProvider {
	if model == "" {
		model = "deepseek/deepseek-coder"
	}

	return &OpenRouterProvider{
		APIKey: apiKey,
		Model:  model,
	}
}

type OpenRouterRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	ResponseFormat struct {
		Type string `json:"type"`
	} `json:"response_format"`
	Stream bool `json:"stream"`
}

type OpenRouterResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (p *OpenRouterProvider) GenerateSuggestions(diff string, includeBody bool, style string) ([]string, error) {
	// Build the prompt
	prompt := buildPrompt(diff, includeBody, style)

	// Create the request
	reqBody := OpenRouterRequest{
		Model: p.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		ResponseFormat: struct{Type string "json:\"type\""}{
			Type: "json_object",
		},
		Stream: false,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	req, err := http.NewRequest("POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)
	req.Header.Set("HTTP-Referer", "https://github.com/amosehiguese/zeus-ai")

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

	var respObj OpenRouterResponse
	err = json.Unmarshal(respBody, &respObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(respObj.Choices) == 0 {
		return nil, fmt.Errorf("API returned no suggestions")
	}

	// Parse the response into individual suggestions
	content := respObj.Choices[0].Message.Content
	return parseJSONResponse(content, includeBody)
}
