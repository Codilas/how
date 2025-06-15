package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	CurrentProvider string                    `yaml:"currentProvider"`
	Providers       map[string]ProviderConfig `yaml:"providers"`
	Context         ContextConfig             `yaml:"context"`
	Display         DisplayConfig             `yaml:"display"`
	History         HistoryConfig             `yaml:"history"`
}

type ProviderConfig struct {
	Type      string `yaml:"type"`
	APIKey    string `yaml:"apiKey"`
	Model     string `yaml:"model"`
	BaseURL   string `yaml:"baseUrl,omitempty"`
	MaxTokens int    `yaml:"maxTokens"`

	// Additional provider-specific settings
	Temperature   float32           `yaml:"temperature,omitempty"`
	TopP          float32           `yaml:"topP,omitempty"`
	SystemPrompt  string            `yaml:"systemPrompt,omitempty"`
	CustomHeaders map[string]string `yaml:"customHeaders,omitempty"`
}

type ContextConfig struct {
	IncludeFiles       bool     `yaml:"includeFiles"`
	IncludeHistory     int      `yaml:"includeHistory"`
	IncludeEnvironment bool     `yaml:"includeEnvironment"`
	IncludeGit         bool     `yaml:"includeGit"`
	MaxContextSize     int      `yaml:"maxContextSize"`
	ExcludePatterns    []string `yaml:"excludePatterns"`
}

type DisplayConfig struct {
	SyntaxHighlight bool `yaml:"syntaxHighlight"`
	ShowContext     bool `yaml:"showContext"`
	Emoji           bool `yaml:"emoji"`
	Color           bool `yaml:"color"`
}

type HistoryConfig struct {
	Enabled  bool   `yaml:"enabled"`
	MaxSize  int    `yaml:"maxSize"`
	FilePath string `yaml:"filePath"`
}

func Load(configFile string) (*Config, error) {
	// Set up viper
	v := viper.New()

	// Determine config file path
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		configDir, err := getConfigDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get config directory: %w", err)
		}

		v.AddConfigPath(configDir)
		v.SetConfigName("config")
		v.SetConfigType("yaml")
	}

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Warning: %v\n", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Environment variable overrides
	v.SetEnvPrefix("HOW")
	v.AutomaticEnv()

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func (c *Config) Save(configFile string) error {
	if configFile == "" {
		configDir, err := getConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %w", err)
		}
		configFile = filepath.Join(configDir, "config.yaml")
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configFile), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "how"), nil
}
