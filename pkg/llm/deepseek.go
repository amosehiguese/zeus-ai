package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DeepSeekProvider struct {
	APIKey string
	Model  string
}

func NewDeepSeekProvider(apiKey string, model string) *DeepSeekProvider {
	if model == "" {
		model = "deepseek-coder"
	}

	return &DeepSeekProvider{
		APIKey: apiKey,
		Model:  model,
	}
}

type DeepSeekRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type DeepSeekResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (p *DeepSeekProvider) GenerateSuggestions(diff string, includeBody bool, style string) ([]string, error) {
	// Build the prompt
	prompt := buildPrompt(diff, includeBody, style)

	// Create the request
	reqBody := DeepSeekRequest{
		Model: p.Model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.APIKey)

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

	var respObj DeepSeekResponse
	err = json.Unmarshal(respBody, &respObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(respObj.Choices) == 0 {
		return nil, fmt.Errorf("API returned no suggestions")
	}

	// Parse the response into individual suggestions
	content := respObj.Choices[0].Message.Content
	return parseSuggestions(content), nil
}
