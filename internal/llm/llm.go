package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Provider is an interface for different LLM providers
type Provider interface {
	GenerateSuggestions(diff string, includeBody bool, style string) ([]string, error)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Suggestion struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
}

type LLMResponse struct {
	Suggestions []Suggestion `json:"suggestions"`
}

func NewProvider(providerType string, apiKey string, model string) (Provider, error) {
	switch strings.ToLower(providerType) {
	case "ollama":
		return NewOllamaProvider(model), nil
	case "openrouter":
		return NewOpenRouterProvider(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}
}

func buildPrompt(diff string, includeBody bool, style string) string {
	var prompt strings.Builder

	prompt.WriteString(`You are a commit message generator. Analyze this git diff and respond with JSON containing exactly 3 commit message suggestions in the following format:
    
{
  "suggestions": [
    {
      "title": "commit title",
      "body": "commit body (optional)"
    }
  ]
}

STRICT REQUIREMENTS:
1. Response must be valid JSON
2. Include exactly 3 suggestions
3. Title must follow Conventional Commits format when requested
4. Omit "body" field when not requested
5. Escape all special JSON characters
6. Do NOT include the git diff in your response
7. Do NOT include any commentary or markdown

Git Diff:
`)
	prompt.WriteString("```diff\n")
	prompt.WriteString(diff)
	prompt.WriteString("\n```\n\n")

	if style == "conventional" {
		prompt.WriteString(`CONVENTIONAL COMMITS RULES:
- Title format: "type(scope): description"
- Types: feat, fix, docs, style, refactor, test, chore
- Scope: optional component name
- Description: imperative mood, lowercase, no period
`)
	}

	if includeBody {
		prompt.WriteString(`BODY REQUIREMENTS:
- Separate from title by blank line
- Explain "what" and "why" not "how"
- Wrap lines at 72 characters
`)
	}

	prompt.WriteString("\nRespond ONLY with valid JSON in this exact format. Do not include any commentary or markdown.")

	return prompt.String()
}


func parseJSONResponse(content string, includeBody bool) ([]string, error) {
	jsonStart := strings.Index(content, "```json")
	if jsonStart >= 0 {
		content = content[jsonStart+7:]
	}

	jsonEnd := strings.Index(content, "```")
	if jsonEnd >= 0 {
		content = content[:jsonEnd]
	}

	content = strings.TrimSpace(content)

	var response LLMResponse
	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("invalid JSON response: %w", err)
	}

	if len(response.Suggestions) != 3 {
		return nil, fmt.Errorf("expected 3 suggestions, got %d", len(response.Suggestions))
	}

	var suggestions []string
	for _, s := range response.Suggestions {
		if includeBody && s.Body != "" {
			suggestions = append(suggestions, fmt.Sprintf("%s\n\n%s", s.Title, s.Body))
		} else {
			suggestions = append(suggestions, s.Title)
		}
	}

	return suggestions, nil
}
