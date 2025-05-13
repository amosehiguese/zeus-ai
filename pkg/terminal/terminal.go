package terminal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Confirm asks the user for confirmation
func Confirm(prompt string) (bool, error) {
	fmt.Printf("%s (y/n): ", prompt)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}

// DisplayAndSelectSuggestion displays the suggestions and allows the user to select one
func DisplayAndSelectSuggestion(suggestions []string) (int, error) {
	fmt.Println("\n🔮 Commit message suggestions:")
	fmt.Println("-----------------------------------------")

	for i, suggestion := range suggestions {
		// For multi-line suggestions, only show the first line in the list
		firstLine := strings.SplitN(suggestion, "\n", 2)[0]
		fmt.Printf("%d. %s\n", i+1, firstLine)
	}

	fmt.Println("-----------------------------------------")
	fmt.Println("e. Edit manually")
	fmt.Println("q. Quit without committing")
	fmt.Print("\nSelect an option (1-5, e, q): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)

	if input == "e" {
		return -1, nil
	}

	if input == "q" {
		fmt.Println("Exiting without committing.")
		os.Exit(0)
	}

	// Try to parse the input as a number
	var selectedIdx int
	_, err = fmt.Sscanf(input, "%d", &selectedIdx)
	if err != nil || selectedIdx < 1 || selectedIdx > len(suggestions) {
		return -1, fmt.Errorf("invalid selection")
	}

	// Show full suggestion if it has multiple lines
	fmt.Println("\nSelected commit message:")
	fmt.Println("-----------------------------------------")
	fmt.Println(suggestions[selectedIdx-1])
	fmt.Println("-----------------------------------------")

	return selectedIdx - 1, nil
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
