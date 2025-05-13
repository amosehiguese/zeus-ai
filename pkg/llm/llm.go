package llm

import (
	"fmt"
	"regexp"
	"strings"
)

// Provider is an interface for different LLM providers
type Provider interface {
	GenerateSuggestions(diff string, includeBody bool, style string) ([]string, error)
}

func NewProvider(providerType string, apiKey string, model string) (Provider, error) {
	switch strings.ToLower(providerType) {
	case "deepseek":
		return NewDeepSeekProvider(apiKey, model), nil
	case "claude":
		return NewClaudeProvider(apiKey, model), nil
	case "ollama":
		return NewOllamaProvider(model), nil
	case "openrouter":
		return NewOpenRouterProvider(apiKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerType)
	}
}

// Helper functions for all providers
func buildPrompt(diff string, includeBody bool, style string) string {
	prompt := fmt.Sprintf(`You are an expert developer analyzing a git diff to create commit messages. 
Below is the git diff:

%s

Please generate 5 concise and descriptive git commit messages based on the changes in the diff.`, diff)

	if style == "conventional" {
		prompt += " Follow the Conventional Commits format (type(scope): description)."
	}

	if includeBody {
		prompt += " For each message, include a title and a detailed body explaining the changes."
	} else {
		prompt += " Each message should be a single line (title only)."
	}

	prompt += " Number each suggestion from 1 to 5."

	return prompt
}

func parseSuggestions(content string) []string {
	lines := strings.Split(content, "\n")
	suggestions := make([]string, 0)

	var currentSuggestion strings.Builder
	inSuggestion := false

	r := regexp.MustCompile(`^[1-5][.:]`)
	for _, line := range lines {
		// Check if line starts a new suggestion (1. or 1:, etc.)
		if matched := r.MatchString(strings.TrimSpace(line)); matched {
			// If we were already in a suggestion, add it to the list
			if inSuggestion && currentSuggestion.Len() > 0 {
				suggestions = append(suggestions, strings.TrimSpace(currentSuggestion.String()))
				currentSuggestion.Reset()
			}

			// Start a new suggestion, removing the number prefix
			parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
			if len(parts) > 1 {
				currentSuggestion.WriteString(strings.TrimSpace(parts[1]))
				inSuggestion = true
			}
		} else if inSuggestion {
			// Continue the current suggestion
			currentSuggestion.WriteString("\n")
			currentSuggestion.WriteString(line)
		}
	}

	// Add the last suggestion if there is one
	if inSuggestion && currentSuggestion.Len() > 0 {
		suggestions = append(suggestions, strings.TrimSpace(currentSuggestion.String()))
	}

	// If we couldn't parse any suggestions, just split by numbered lines
	if len(suggestions) == 0 {
		// Try alternative parsing strategy
		regex := regexp.MustCompile(`(?m)^[1-5][.:](.+)(?:\n(?:[^\n1-5].*(?:\n|$))*)`)
		matches := regex.FindAllStringSubmatch(content, -1)

		for _, match := range matches {
			if len(match) > 1 {
				suggestions = append(suggestions, strings.TrimSpace(match[0]))
			}
		}
	}

	// Clean up suggestions by removing the number prefixes
	for i, suggestion := range suggestions {
		suggestion = strings.TrimSpace(suggestion)
		if matched := r.MatchString(suggestion); matched {
			parts := strings.SplitN(suggestion, " ", 2)
			if len(parts) > 1 {
				suggestions[i] = strings.TrimSpace(parts[1])
			}
		}
	}

	return suggestions
}
