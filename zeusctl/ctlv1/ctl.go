package ctlv1

import (
	"os"

	"github.com/amosehiguese/zeus-ai/pkg/cobrautil"
	"github.com/amosehiguese/zeus-ai/pkg/terminal"
	"github.com/amosehiguese/zeus-ai/zeusctl/ctlv1/command"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

const (
	cliName = "zeusctl"
	cliDescription = "A Git-aware CLI tool that helps generate smart commit messages"
)

var (
	rootCmd = &cobra.Command{
		Use:   cliName,
		Short: cliDescription,
		Long: terminal.TitleColor.Sprint(`zeus-ai is a Git-aware CLI tool that helps developers generate smart commit messages using LLM API.
		The tool interacts with the Git diff, asks the LLM for commit message suggestions, and provides
		a terminal interface for confirming, editing, and optionally signing the commit.`),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Disable colors if output isn't a terminal
			color.NoColor = !isTerminal()
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	rootCmd.AddCommand(
		command.NewInitCommand(),
		command.NewVersionCommand(),
		command.NewSuggestCommand(),
	)
}

func Start() error {
	return rootCmd.Execute()
}

func MustStart() {
	if err := Start(); err != nil {
		if rootCmd.SilenceErrors {
			cobrautil.ExitWithError(cobrautil.ExitError, err)
		}
		os.Exit(cobrautil.ExitError)
	}
}

func isTerminal() bool {
	fileInfo, _ := os.Stdout.Stat()
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}