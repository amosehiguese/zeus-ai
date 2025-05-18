package command

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/amosehiguese/zeus-ai/internal/config"
	"github.com/amosehiguese/zeus-ai/internal/git"
	"github.com/amosehiguese/zeus-ai/internal/llm"
	"github.com/amosehiguese/zeus-ai/internal/terminal"
)

var (
	bodyFlag      bool
	editFlag      bool
	signFlag      bool
	dryRunFlag    bool
	autoStageFlag bool
)

func NewSuggestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suggest",
		Short: "Suggest commit messages based on git diff",
		Long:  `Suggest commit messages based on staged changes (or unstaged if nothing is staged)`,
		RunE:  suggestCommandFunc,
	}

	cmd.Flags().BoolVar(&bodyFlag, "body", false, "Include detailed body text in suggestions")
	cmd.Flags().BoolVar(&editFlag, "edit", false, "Open the selected message in default editor")
	cmd.Flags().BoolVar(&signFlag, "sign", false, "Sign the commit message")
	cmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Display message suggestions but don't run commit")
	cmd.Flags().BoolVar(&autoStageFlag, "auto-stage", false, "Automatically stage all changes")
	cmd.Flags().StringVar(&styleFlag, "style", "conventional", "Specify commit style (e.g., conventional, simple)")

	return cmd
}

func suggestCommandFunc(cmd *cobra.Command, args []string) error {
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
		if err = git.StageAllChanges(); err != nil {
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
		terminal.ShowWarning("No staged changes found.")
		var unstaged bool
		unstaged, err = git.HasUnstagedChanges()
		if err != nil {
			return fmt.Errorf("failed to check for unstaged changes: %w", err)
		}

		if unstaged {
			var shouldUseUnstaged bool
			shouldUseUnstaged, err = terminal.Confirm("Would you like to use unstaged changes instead?")
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
	if stats, statErr := git.GetDiffStats(true); statErr == nil {
		terminal.ShowDiffStats(stats)
	}

	stopSpinner := terminal.ShowSpinner("Generating comit message suggestions...")
	suggestions, err := provider.GenerateSuggestions(diff, bodyFlag, styleFlag)
	stopSpinner()
	if err != nil {
		log.Printf("Got an error while generating suggestions: %v", err)
		return err
	}
	terminal.ShowSuccess("Generated 3 suggesions")

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

	if dryRunFlag {
		terminal.ShowSuccess("Dry run - would commit:")
		fmt.Println(commitMsg)
		return nil
	}

	// Perform the commit
	if err := git.Commit(commitMsg, signFlag); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	terminal.ShowSuccess("Commit created successfully")
	return nil
}
