package terminal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func DisplayAndSelectSuggestion(suggestions []string) (int, error) {
	if len(suggestions) != 3 {
		return -1, fmt.Errorf("expected 3 suggestions, got %d", len(suggestions))
	}

	printHeader("COMMIT MESSAGE SUGGESTIONS")

	for i, suggestion := range suggestions {
		parts := strings.SplitN(suggestion, "\n\n", 2)
		TitleColor.Printf("  %d. %s\n", i+1, parts[0])

		if len(parts) > 1 {
			BodyColor.Printf("     %s\n", indentBody(parts[1]))
		}
	}

	printOptions()
	return getSelection(len(suggestions))
}

func ShowDiff(diff string) {
	lines := strings.Split(diff, "\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "+"):
			DiffAddColor.Println(line)
		case strings.HasPrefix(line, "-"):
			DiffRemoveColor.Println(line)
		default:
			fmt.Println(line)
		}
	}
}

func ShowDiffStats(stats string) {
	DividerColor.Println("\nGIT DIFF STATS:")
	BodyColor.Println(stats)
	DividerColor.Println(strings.Repeat("─", 40))
}

func ShowSpinner(message string) func() {
	stop := make(chan bool)
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

	go func() {
		i := 0
		for {
			select {
			case <-stop:
				return
			default:
				SpinnerColor.Printf("\r%s %s", frames[i], message)
				i = (i + 1) % len(frames)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return func() {
		stop <- true
		fmt.Printf("\r%s\r", strings.Repeat(" ", 60))
	}
}

func Confirm(prompt string) (bool, error) {
	PromptColor.Printf("%s (y/N): ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("input error: %w", err)
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}

func ShowSuccess(message string) {
	SuccessColor.Println("✓", message)
}

func ShowWarning(message string) {
	WarningColor.Println("!", message)
}

func ShowError(message string) {
	ErrorColor.Println("✖", message)
}

// Helper functions
func printHeader(title string) {
	DividerColor.Println("\n┌───────────────────────────────────────────────────────┐")
	TitleColor.Printf("  %s\n", title)
	DividerColor.Println("├───────────────────────────────────────────────────────┤")
}

func printOptions() {
	DividerColor.Println("├───────────────────────────────────────────────────────┤")
	OptionColor.Println("  e - Edit manually")
	OptionColor.Println("  q - Quit without committing")
	DividerColor.Println("└───────────────────────────────────────────────────────┘")
}

func getSelection(max int) (int, error) {
	PromptColor.Print("\n  Select an option (1-3/e/q): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, fmt.Errorf("input error: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))
	switch input {
	case "e":
		return -1, nil // Edit mode
	case "q":
		os.Exit(0)
	}

	idx, err := strconv.Atoi(input)
	if err != nil || idx < 1 || idx > max {
		ShowError("Invalid selection. Please choose 1-3, e, or q")
		return getSelection(max)
	}
	return idx - 1, nil
}

func indentBody(body string) string {
	return strings.ReplaceAll(body, "\n", "\n     ")
}

// EditMessage opens the default editor to edit the message
func EditMessage(initialContent string, includeBody bool) (string, error) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "zeus-commit-msg-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the initial content to the file
	_, err = tmpFile.WriteString(initialContent)
	if err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	// Add a template/guidance for the commit message if it's empty
	if initialContent == "" {
		template := "# Enter your commit message here\n"
		if includeBody {
			template += "# First line: a brief summary (50 chars or less)\n\n# Body: more detailed explanatory text if needed\n"
		}
		_, err = tmpFile.WriteString(template)
		if err != nil {
			return "", fmt.Errorf("failed to write template to temp file: %w", err)
		}
	}

	// Close the file to ensure all data is written
	if err = tmpFile.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %w", err)
	}

	// Get editor from environment or use a default
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try to find a default editor based on the OS
		if _, err := exec.LookPath("nano"); err == nil {
			editor = "nano"
		} else if _, err := exec.LookPath("vim"); err == nil {
			editor = "vim"
		} else if _, err := exec.LookPath("vi"); err == nil {
			editor = "vi"
		} else if _, err := exec.LookPath("notepad"); err == nil {
			editor = "notepad"
		} else {
			return "", fmt.Errorf("no suitable editor found, please set EDITOR environment variable")
		}
	}

	// Open the editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("editor command failed: %w", err)
	}

	// Read the updated content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}

	// Remove comments and trailing empty lines
	lines := strings.Split(string(content), "\n")
	var result []string
	for _, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			result = append(result, line)
		}
	}

	// Trim empty lines from the end
	for len(result) > 0 && strings.TrimSpace(result[len(result)-1]) == "" {
		result = result[:len(result)-1]
	}

	return strings.Join(result, "\n"), nil
}
