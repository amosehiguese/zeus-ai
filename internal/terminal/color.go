package terminal

import "github.com/fatih/color"

var (
	// Primary colors
	TitleColor   = color.New(color.FgHiCyan, color.Bold)
	BodyColor    = color.New(color.FgHiWhite)
	OptionColor  = color.New(color.FgHiYellow)
	PromptColor  = color.New(color.FgHiBlue, color.Bold)
	DividerColor = color.New(color.FgHiMagenta)
	SpinnerColor = color.New(color.FgHiMagenta)

	// Status colors
	SuccessColor = color.New(color.FgHiGreen)
	WarningColor = color.New(color.FgHiYellow)
	ErrorColor   = color.New(color.FgHiRed)

	// Git diff colors
	DiffAddColor    = color.New(color.FgGreen)
	DiffRemoveColor = color.New(color.FgRed)
)
