package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TempConfigDir creates a temporary directory for config file testing and returns the path.
// The caller is responsible for cleaning up the directory using cleanup().
func TempConfigDir(t *testing.T) (string, func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "how-config-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	cleanup := func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Errorf("failed to remove temp directory: %v", err)
		}
	}

	return tempDir, cleanup
}

// SampleConfig returns a valid Config struct with default test values.
func SampleConfig() *Config {
	return &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-ant-test-key-12345",
				Model:     "claude-3-5-sonnet-20241022",
				MaxTokens: 2048,
				Temperature: 0.7,
				TopP:       1.0,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     50,
			IncludeEnvironment: false,
			IncludeGit:         true,
			MaxContextSize:     8000,
			ExcludePatterns:    []string{".git", "node_modules"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Emoji:           false,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  1000,
			FilePath: filepath.Join("~", ".local", "share", "how", "history"),
		},
	}
}

// SampleYAML returns a sample YAML configuration as a string.
func SampleYAML() string {
	return `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: sk-ant-test-key-12345
    model: claude-3-5-sonnet-20241022
    maxTokens: 2048
    temperature: 0.7
    topP: 1
context:
  includeFiles: true
  includeHistory: 50
  includeEnvironment: false
  includeGit: true
  maxContextSize: 8000
  excludePatterns:
    - .git
    - node_modules
display:
  syntaxHighlight: true
  showContext: true
  emoji: false
  color: true
history:
  enabled: true
  maxSize: 1000
  filePath: ~/.local/share/how/history
`
}

// MinimalYAML returns a minimal valid YAML configuration with only required fields.
func MinimalYAML() string {
	return `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: sk-ant-test-key
    model: claude-3-5-sonnet-20241022
    maxTokens: 2048
`
}

// WriteConfigFile writes the provided YAML content to a config file at the specified path.
func WriteConfigFile(t *testing.T, configPath string, yamlContent string) {
	t.Helper()

	// Ensure parent directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Write the config file
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
}

// ConfigWithProvider returns a Config struct with a custom provider configuration.
func ConfigWithProvider(t *testing.T, providerName, providerType, apiKey, model string) *Config {
	t.Helper()

	return &Config{
		CurrentProvider: providerName,
		Providers: map[string]ProviderConfig{
			providerName: {
				Type:      providerType,
				APIKey:    apiKey,
				Model:     model,
				MaxTokens: 2048,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     50,
			IncludeEnvironment: false,
			IncludeGit:         true,
			MaxContextSize:     8000,
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled: true,
			MaxSize: 1000,
		},
	}
}

// EmptyYAML returns an empty but valid YAML string.
func EmptyYAML() string {
	return ""
}
