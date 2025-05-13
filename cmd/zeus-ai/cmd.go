package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amosehiguese/zeus-ai/internal/config"
	"github.com/amosehiguese/zeus-ai/pkg/git"
	"github.com/amosehiguese/zeus-ai/pkg/llm"
	"github.com/amosehiguese/zeus-ai/pkg/terminal"
	"github.com/spf13/cobra"
)

func handleSuggestCommand(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Check if git repo
	if !git.IsGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// Auto-stage if flag is set
	if autoStageFlag {
		if err := git.StageAllChanges(); err != nil {
			return fmt.Errorf("failed to stage changes: %w", err)
		}
	}

	// Get diff
	diff, err := git.GetDiff(true) // Get staged diff first
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	// If no staged changes, check if there are unstaged changes
	if diff == "" {
		fmt.Println("No staged changes found.")

		unstaged, err := git.HasUnstagedChanges()
		if err != nil {
			return fmt.Errorf("failed to check for unstaged changes: %w", err)
		}

		if unstaged {
			shouldUseUnstaged, err := terminal.Confirm("Would you like to use unstaged changes instead?")
			if err != nil {
				return fmt.Errorf("failed to get confirmation: %w", err)
			}

			if shouldUseUnstaged {
				diff, err = git.GetDiff(false)
				if err != nil {
					return fmt.Errorf("failed to get unstaged diff: %w", err)
				}
			} else {
				return fmt.Errorf("no changes to commit")
			}
		} else {
			return fmt.Errorf("no changes to commit")
		}
	}

	// Create LLM provider
	provider, err := llm.NewProvider(cfg.Provider, cfg.APIKey, cfg.Model)
	if err != nil {
		return fmt.Errorf("failed to create LLM provider: %w", err)
	}

	// Show diff stats
	stats, err := git.GetDiffStats(true)
	if err == nil && stats != "" {
		fmt.Println("\n📊 Git diff stats:")
		fmt.Println("-----------------------------------------")
		fmt.Println(stats)
		fmt.Println("-----------------------------------------")
	}

	// Generate suggestions with spinner
	fmt.Print("🧠 Generating commit message suggestions... ")

	// Create a spinner in a goroutine
	stopSpinner := make(chan bool)
	go func() {
		spinChars := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-stopSpinner:
				fmt.Print("\r")
				return
			default:
				fmt.Printf("\r🧠 Generating commit message suggestions... %s", spinChars[i])
				i = (i + 1) % len(spinChars)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	suggestions, err := provider.GenerateSuggestions(diff, bodyFlag, styleFlag)
	if err != nil {
		log.Printf("Got an error while generating suggestions: %v", err)
		return err
	}
	// Stop the spinner
	stopSpinner <- true
	fmt.Println("\r✅ Generated suggestions!                  ")

	// Display suggestions
	selectedIdx, err := terminal.DisplayAndSelectSuggestion(suggestions)
	if err != nil {
		return fmt.Errorf("failed to select suggestion: %w", err)
	}

	var commitMsg string
	if selectedIdx == -1 {
		// User wants to edit manually
		commitMsg, err = terminal.EditMessage("", bodyFlag)
		if err != nil {
			return fmt.Errorf("failed to edit message: %w", err)
		}
	} else {
		commitMsg = suggestions[selectedIdx]

		// If edit flag is set, open the selected message in editor
		if editFlag {
			commitMsg, err = terminal.EditMessage(commitMsg, bodyFlag)
			if err != nil {
				return fmt.Errorf("failed to edit message: %w", err)
			}
		}
	}

	// Perform the commit
	if !dryRunFlag {
		err = git.Commit(commitMsg, signFlag)
		if err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}
		fmt.Println("✅ Commit successful!")
	} else {
		fmt.Println("📝 Dry run - message that would be committed:")
		fmt.Println("-----")
		fmt.Println(commitMsg)
		fmt.Println("-----")
	}

	return nil
}

func handleInitCommand(cmd *cobra.Command, args []string) error {
	provider, _ := cmd.Flags().GetString("provider")
	apiKey, _ := cmd.Flags().GetString("api-key")
	model, _ := cmd.Flags().GetString("model")
	style, _ := cmd.Flags().GetString("style")

	config := fmt.Sprintf(`# zeus-ai configuration
provider: %s
api_key: %s
model: %s
default_style: %s
`, provider, apiKey, model, style)

	err := os.WriteFile(".zeusrc", []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Println("✅ Configuration file .zeusrc created successfully!")
	return nil
}
