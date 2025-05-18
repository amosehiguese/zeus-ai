package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func setupGitRepo(t *testing.T) string {
	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "zeus-git-test-*")
	require.NoError(t, err, "Failed to create temp directory")

	// Initialize git repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		os.RemoveAll(tmpDir)
		require.NoError(t, err, "Failed to initialize git repository")
	}

	// Setup git config
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		os.RemoveAll(tmpDir)
		require.NoError(t, err, "Failed to set git user.name")
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	err = cmd.Run()
	if err != nil {
		os.RemoveAll(tmpDir)
		require.NoError(t, err, "Failed to set git user.email")
	}

	return tmpDir
}

// createAndAddFile creates a file with content and stages it
func createAndAddFile(t *testing.T, repoDir, filename, content string) {
	filePath := filepath.Join(repoDir, filename)
	err := os.WriteFile(filePath, []byte(content), 0o644)
	require.NoError(t, err, "Failed to write file")

	cmd := exec.Command("git", "add", filename)
	cmd.Dir = repoDir
	err = cmd.Run()
	require.NoError(t, err, "Failed to stage file")
}

func TestIsGitRepository(t *testing.T) {
	// Setup
	tmpDir := setupGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to git repository directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Test
	require.True(t, IsGitRepository(), "Expected directory to be a git repository")

	// Change to a non-git directory
	nonGitDir, err := os.MkdirTemp("", "zeus-non-git-*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(nonGitDir)

	err = os.Chdir(nonGitDir)
	require.NoError(t, err, "Failed to change directory")

	require.False(t, IsGitRepository(), "Expected directory to not be a git repository")
}

func TestGetDiff(t *testing.T) {
	// Setup
	tmpDir := setupGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to git repository directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Create and stage a file
	createAndAddFile(t, tmpDir, "test.txt", "Hello, Zeus!")

	// Test staged diff
	diff, err := GetDiff(true)
	require.NoError(t, err, "Failed to get staged diff")
	require.Contains(t, diff, "Hello, Zeus!", "Expected diff to contain file content")

	// Modify existing file
	err = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("Hello, Zeus!\nAnother content"), 0o644)
	require.NoError(t, err, "Failed to write file")

	// Test unstaged diff
	diff, err = GetDiff(false)
	require.NoError(t, err, "Failed to get unstaged diff")
	require.Contains(t, diff, "Another content", "Expected diff to contain unstaged file content")
}

func TestHasUnstagedChanges(t *testing.T) {
	// Setup
	tmpDir := setupGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to git repository directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Initially there should be no unstaged changes
	hasChanges, err := HasUnstagedChanges()
	require.NoError(t, err, "Failed to check for unstaged changes")
	require.False(t, hasChanges, "Expected no unstaged changes initially")

	// Create and stage a file to be tracked
	createAndAddFile(t, tmpDir, "test.txt", "Hello, Zeus!")

	// Modify file to be unstaged
	err = os.WriteFile(filepath.Join(tmpDir, "test.txt"), []byte("Hello, Zeus\nUnstaged content"), 0o644)
	require.NoError(t, err, "Failed to write file")

	// Now there should be unstaged changes
	hasChanges, err = HasUnstagedChanges()
	require.NoError(t, err, "Failed to check for unstaged changes")
	require.True(t, hasChanges, "Expected to find unstaged changes")
}

func TestStageAllChanges(t *testing.T) {
	// Setup
	tmpDir := setupGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to git repository directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Create an unstaged file
	err = os.WriteFile(filepath.Join(tmpDir, "unstaged.txt"), []byte("To be staged"), 0o644)
	require.NoError(t, err, "Failed to write file")

	// Stage all changes
	err = StageAllChanges()
	require.NoError(t, err, "Failed to stage all changes")

	// Verify no unstaged changes remain
	hasChanges, err := HasUnstagedChanges()
	require.NoError(t, err, "Failed to check for unstaged changes")
	require.False(t, hasChanges, "Expected no unstaged changes after staging all")

	// Check if file is staged
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	require.NoError(t, err, "Failed to check staged files")
	require.Contains(t, string(output), "unstaged.txt", "Expected file to be staged")
}

func TestGetDiffStats(t *testing.T) {
	// Setup
	tmpDir := setupGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to git repository directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Create and stage a file
	createAndAddFile(t, tmpDir, "stats-test.txt", "Line 1\nLine 2\nLine 3\n")

	// Test staged diff stats
	stats, err := GetDiffStats(true)
	require.NoError(t, err, "Failed to get diff stats")
	require.Contains(t, stats, "stats-test.txt", "Expected diff stats to mention the filename")
}

func TestCommit(t *testing.T) {
	// Setup
	tmpDir := setupGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to git repository directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Create and stage a file
	createAndAddFile(t, tmpDir, "commit-test.txt", "Test content for commit")

	// Make a commit
	commitMsg := "Test commit message"
	err = Commit(commitMsg, false)
	require.NoError(t, err, "Failed to commit")

	// Verify the commit was made
	cmd := exec.Command("git", "log", "--oneline", "-1")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	require.NoError(t, err, "Failed to get git log")
	require.Contains(t, string(output), commitMsg, "Expected git log to contain commit message")
}

// TestSignedCommit tests commit signing, but skips if GPG is not configured
func TestSignedCommit(t *testing.T) {
	// Check if git can sign commits
	canSign := exec.Command("git", "config", "--get", "user.signingkey").Run() == nil

	if !canSign {
		t.Skip("Skipping signed commit test as GPG signing is not configured")
	}

	// Setup
	tmpDir := setupGitRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to git repository directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Configure signing
	cmd := exec.Command("git", "config", "commit.gpgsign", "true")
	cmd.Dir = tmpDir
	err = cmd.Run()
	require.NoError(t, err, "Failed to configure git signing")

	// Create and stage a file
	createAndAddFile(t, tmpDir, "signed-commit-test.txt", "Test content for signed commit")

	// Make a signed commit
	commitMsg := "Signed commit message"
	err = Commit(commitMsg, true)
	require.NoError(t, err, "Failed to create signed commit")

	// Verify the commit was made and signed
	cmd = exec.Command("git", "log", "--show-signature", "-1")
	cmd.Dir = tmpDir
	output, err := cmd.Output()
	require.NoError(t, err, "Failed to get git log with signature")
	require.Contains(t, string(output), commitMsg, "Expected git log to contain commit message")
}
