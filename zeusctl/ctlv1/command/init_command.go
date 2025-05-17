package command

import (
	"fmt"
	"os"

	"github.com/amosehiguese/zeus-ai/internal/terminal"
	"github.com/spf13/cobra"
)

var (
	providerFlag string
	apiKeyFlag string
	modelFlag string
	styleFlag string
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize zeus-ai configuration",
		Long:  `Initialize zeus-ai configuration by creating a .zeusrc file in the current directory`,
		RunE:  initCommandFunc,
	}

	cmd.Flags().StringVar(&providerFlag, "provider", "ollama", "LLM provider (ollama, openrouter)")
	cmd.Flags().StringVar(&apiKeyFlag, "api-key", "", "API key for the provider")
	cmd.Flags().StringVar(&modelFlag, "model", "deepseek-coder", "Model to use")
	cmd.Flags().StringVar(&styleFlag, "style", "conventional", "Default commit style")

	return cmd
}

func initCommandFunc(cmd *cobra.Command, args []string) error {
	config := fmt.Sprintf(`# zeus-ai configuration
provider: %s
api_key: %s
model: %s
default_style: %s
`, providerFlag, apiKeyFlag, modelFlag, styleFlag)

	err := os.WriteFile(".zeusrc", []byte(config), 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	terminal.ShowSuccess("✅ Configuration file .zeusrc created successfully!")
	return nil
}