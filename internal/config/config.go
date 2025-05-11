package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Provider    string
	APIKey      string
	Model       string
	DefaultStyle string
}

func Load() (*Config, error) {
	config := &Config{
		Provider:    "ollama", // Default provider
		Model:       "deepseek-coder", // Default model
		DefaultStyle: "conventional",
	}

	viper.SetConfigName(".zeusrc")
	viper.SetConfigType("yaml")
	
	// Look for config in home directory
	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
	}
	
	// Look for config in current directory and all parent directories
	currentDir, err := os.Getwd()
	if err == nil {
		for {
			viper.AddConfigPath(currentDir)
			parent := filepath.Dir(currentDir)
			if parent == currentDir {
				break
			}
			currentDir = parent
		}
	}

	// Environment variables
	viper.SetEnvPrefix("ZEUS")
	viper.AutomaticEnv()

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there is no config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Override with environment variables or config file values
	if viper.IsSet("provider") {
		config.Provider = viper.GetString("provider")
	}
	if viper.IsSet("api_key") {
		config.APIKey = viper.GetString("api_key")
	}
	if viper.IsSet("model") {
		config.Model = viper.GetString("model")
	}
	if viper.IsSet("default_style") {
		config.DefaultStyle = viper.GetString("default_style")
	}

	// Check for environment variables
	if os.Getenv("ZEUS_PROVIDER") != "" {
		config.Provider = os.Getenv("ZEUS_PROVIDER")
	}
	if os.Getenv("ZEUS_API_KEY") != "" {
		config.APIKey = os.Getenv("ZEUS_API_KEY")
	}
	if os.Getenv("ZEUS_MODEL") != "" {
		config.Model = os.Getenv("ZEUS_MODEL")
	}
	if os.Getenv("ZEUS_DEFAULT_STYLE") != "" {
		config.DefaultStyle = os.Getenv("ZEUS_DEFAULT_STYLE")
	}

	return config, nil
}