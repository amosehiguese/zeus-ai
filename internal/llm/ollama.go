package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OllamaProvider struct {
	Model string
}

func NewOllamaProvider(model string) *OllamaProvider {
	if model == "" {
		model = "deepseek-coder"
	}

	return &OllamaProvider{
		Model: model,
	}
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Format string `json:"format"`
	Stream bool   `json:"stream"`
	Options struct {
		Temperature float64 `json:"temperature"`
	} `json:"options"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func (p *OllamaProvider) GenerateSuggestions(diff string, includeBody bool, style string) ([]string, error) {
	// Check if Ollama is running
	_, err := http.Get("http://localhost:11434/api/version")
	if err != nil {
		return nil, fmt.Errorf("ollama server not running at http://localhost:11434. Start Ollama or use a different provider")
	}

	// Build the prompt
	prompt := buildPrompt(diff, includeBody, style)

	// Create the request
	reqBody := OllamaRequest{
		Model:  p.Model,
		Prompt: prompt,
		Format: "json",
		Stream: false,
		Options: struct{Temperature float64 "json:\"temperature\""}{
			Temperature: 0.7,
		},
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make the API request
	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Ollama: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned error: %s", string(respBody))
	}

	var respObj OllamaResponse
	err = json.Unmarshal(respBody, &respObj)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return parseJSONResponse(respObj.Response, includeBody)
}
