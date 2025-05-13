package git

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

// IsGitRepository checks if the current directory is a git repository
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}

// GetDiff returns the git diff (staged or unstaged)
func GetDiff(staged bool) (string, error) {
	args := []string{"diff"}
	if staged {
		args = append(args, "--cached")
	}

	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w", err)
	}

	return out.String(), nil
}

// HasUnstagedChanges checks if there are any unstaged changes
func HasUnstagedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, fmt.Errorf("git diff failed: %w", err)
	}

	return out.String() != "", nil
}

// GetDiffStats returns statistics about the git diff
func GetDiffStats(staged bool) (string, error) {
	args := []string{"diff", "--stat"}
	if staged {
		args = append(args, "--cached")
	}

	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("git diff --stat failed: %w", err)
	}

	return out.String(), nil
}

// StageAllChanges stages all changes
func StageAllChanges() error {
	cmd := exec.Command("git", "add", "-A")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	return nil
}

// Commit performs a git commit with the given message
func Commit(message string, sign bool) error {
	args := []string{"commit", "-m", message}
	if sign {
		args = append(args, "-S")
	}

	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	return nil
}
