package command

import (
	"github.com/spf13/cobra"

	"github.com/amosehiguese/zeus-ai/api/version"
	"github.com/amosehiguese/zeus-ai/internal/terminal"
)

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of zeus-ai",
		Run:   versionCommandFunc,
	}
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
	terminal.TitleColor.Println("zeusctl version:", version.CtlVersion)
	terminal.TitleColor.Println("zeus-ai version:", version.Version)
}
