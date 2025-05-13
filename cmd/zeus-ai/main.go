package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	bodyFlag      bool
	editFlag      bool
	signFlag      bool
	dryRunFlag    bool
	autoStageFlag bool
	styleFlag     string
)


func main() {
	var rootCmd = &cobra.Command{
		Use:   "zeus-ai",
		Short: "A Git-aware CLI tool that helps generate smart commit messages",
		Long: `zeus-ai is a Git-aware CLI tool that helps developers generate smart commit messages using LLM API.
The tool interacts with the Git diff, asks the LLM for commit message suggestions, and provides
a terminal interface for confirming, editing, and optionally signing the commit.`,
	}

	var suggestCmd = &cobra.Command{
		Use:   "suggest",
		Short: "Suggest commit messages based on git diff",
		Long:  `Suggest commit messages based on staged changes (or unstaged if nothing is staged)`,
		RunE:  handleSuggestCommand,
	}

	suggestCmd.Flags().BoolVar(&bodyFlag, "body", false, "Include detailed body text in suggestions")
	suggestCmd.Flags().BoolVar(&editFlag, "edit", false, "Open the selected message in default editor")
	suggestCmd.Flags().BoolVar(&signFlag, "sign", false, "GPG-sign the commit message")
	suggestCmd.Flags().BoolVar(&dryRunFlag, "dry-run", false, "Display message suggestions but don't run commit")
	suggestCmd.Flags().BoolVar(&autoStageFlag, "auto-stage", false, "Automatically stage all changes")
	suggestCmd.Flags().StringVar(&styleFlag, "style", "conventional", "Specify commit style (e.g., conventional, simple)")

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize zeus-ai configuration",
		Long:  `Initialize zeus-ai configuration by creating a .zeusrc file in the current directory`,
		RunE: handleInitCommand,
	}

	initCmd.Flags().String("provider", "ollama", "LLM provider (deepseek, claude, ollama, openrouter)")
	initCmd.Flags().String("api-key", "", "API key for the provider")
	initCmd.Flags().String("model", "deepseek-coder", "Model to use")
	initCmd.Flags().String("style", "conventional", "Default commit style")

	rootCmd.AddCommand(suggestCmd)
	rootCmd.AddCommand(initCmd)

	// Add version command
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of zeus-ai",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("zeus-ai v1.0.0")
		},
	}
	rootCmd.AddCommand(versionCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
