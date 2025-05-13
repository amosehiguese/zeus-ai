package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestLoadWithEnvironmentVariables(t *testing.T) {
	viper.Reset()

	// Set environment variables
	os.Setenv("ZEUS_PROVIDER", "claude")
	os.Setenv("ZEUS_API_KEY", "test-api-key")
	os.Setenv("ZEUS_MODEL", "test-model")
	os.Setenv("ZEUS_DEFAULT_STYLE", "simple")
	defer func() {
		os.Unsetenv("ZEUS_PROVIDER")
		os.Unsetenv("ZEUS_API_KEY")
		os.Unsetenv("ZEUS_MODEL")
		os.Unsetenv("ZEUS_DEFAULT_STYLE")
	}()

	// Load configuration
	cfg, err := Load()
	require.NoError(t, err, "Failed to load config")

	// Verify configuration values
	require.Equal(t, "claude", cfg.Provider, "Wrong provider value")
	require.Equal(t, "test-api-key", cfg.APIKey, "Wrong API key value")
	require.Equal(t, "test-model", cfg.Model, "Wrong model value")
	require.Equal(t, "simple", cfg.DefaultStyle, "Wrong style value")
}

func TestLoadWithConfigFile(t *testing.T) {
	viper.Reset()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "zeus-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(currentDir)

	// Change to temporary directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create a config file
	configContent := `
provider: deepseek
api_key: file-api-key
model: file-model
default_style: conventional
`
	err = os.WriteFile(".zeusrc", []byte(configContent), 0644)
	require.NoError(t, err, "Failed to write config file")

	// Load configuration
	cfg, err := Load()
	require.NoError(t, err, "Failed to load config")

	// Verify configuration values
	require.Equal(t, "deepseek", cfg.Provider, "Wrong provider value")
	require.Equal(t, "file-api-key", cfg.APIKey, "Wrong API key value")
	require.Equal(t, "file-model", cfg.Model, "Wrong model value")
	require.Equal(t, "conventional", cfg.DefaultStyle, "Wrong style value")
}

func TestLoadWithParentConfigFile(t *testing.T) {
	viper.Reset()

	// Create a temporary directory structure
	tmpDir, err := os.MkdirTemp("", "zeus-config-test-*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tmpDir)

	// Create parent and child directories
	parentDir := filepath.Join(tmpDir, "parent")
	childDir := filepath.Join(parentDir, "child")
	err = os.MkdirAll(childDir, 0755)
	require.NoError(t, err, "Failed to create directory structure")

	// Create a config file in the parent directory
	configContent := `
provider: openrouter
api_key: parent-api-key
model: parent-model
default_style: simple
`
	err = os.WriteFile(filepath.Join(parentDir, ".zeusrc"), []byte(configContent), 0644)
	require.NoError(t, err, "Failed to write config file")

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to child directory
	err = os.Chdir(childDir)
	require.NoError(t, err, "Failed to change directory")

	// Load configuration
	cfg, err := Load()
	require.NoError(t, err, "Failed to load config")

	// Verify configuration values from parent directory
	require.Equal(t, "openrouter", cfg.Provider, "Wrong provider value")
	require.Equal(t, "parent-api-key", cfg.APIKey, "Wrong API key value")
	require.Equal(t, "parent-model", cfg.Model, "Wrong model value")
	require.Equal(t, "simple", cfg.DefaultStyle, "Wrong style value")
}

func TestEnvironmentOverridesFile(t *testing.T) {
	viper.Reset()

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "zeus-config-test-*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to temporary directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Create a config file
	configContent := `
provider: deepseek
api_key: file-api-key
model: file-model
default_style: conventional
`
	err = os.WriteFile(".zeusrc", []byte(configContent), 0644)
	require.NoError(t, err, "Failed to write config file")

	// Set environment variables that should override the file
	os.Setenv("ZEUS_PROVIDER", "claude")
	os.Setenv("ZEUS_MODEL", "env-model")
	defer func() {
		os.Unsetenv("ZEUS_PROVIDER")
		os.Unsetenv("ZEUS_MODEL")
	}()

	// Load configuration
	cfg, err := Load()
	require.NoError(t, err, "Failed to load config")

	// Verify environment variables override file
	require.Equal(t, "claude", cfg.Provider, "Environment variable should override file for Provider")
	require.Equal(t, "file-api-key", cfg.APIKey, "APIKey should remain from file")
	require.Equal(t, "env-model", cfg.Model, "Environment variable should override file for Model")
	require.Equal(t, "conventional", cfg.DefaultStyle, "DefaultStyle should remain from file")
}

func TestLoadDefaultValues(t *testing.T) {
	viper.Reset()

	// Create a temporary directory with no config file
	tmpDir, err := os.MkdirTemp("", "zeus-config-test-*")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tmpDir)

	// Save current directory
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(currentDir)

	// Change to temporary directory
	err = os.Chdir(tmpDir)
	require.NoError(t, err, "Failed to change directory")

	// Clear any existing environment variables
	os.Unsetenv("ZEUS_PROVIDER")
	os.Unsetenv("ZEUS_API_KEY")
	os.Unsetenv("ZEUS_MODEL")
	os.Unsetenv("ZEUS_DEFAULT_STYLE")

	// Load configuration
	cfg, err := Load()
	require.NoError(t, err, "Failed to load config")

	// Verify default values
	require.Equal(t, "ollama", cfg.Provider, "Wrong default Provider")
	require.Equal(t, "deepseek-coder", cfg.Model, "Wrong default Model")
	require.Equal(t, "conventional", cfg.DefaultStyle, "Wrong default Style")
}
