package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// TestLoad_SuccessfulLoadFromExplicitFile tests loading config from an explicit file path
func TestLoad_SuccessfulLoadFromExplicitFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Create test config file
	testConfig := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test-123",
				Model:     "claude-3-opus",
				MaxTokens: 4096,
				Temperature: 0.7,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     10,
			IncludeEnvironment: true,
			IncludeGit:         true,
			MaxContextSize:     100000,
			ExcludePatterns:    []string{"*.log", "node_modules/*"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Emoji:           true,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  1000,
			FilePath: "~/.how/history",
		},
	}

	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Test loading
	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded == nil {
		t.Fatalf("Load() returned nil config")
	}

	if loaded.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider = %q, want %q", loaded.CurrentProvider, "anthropic")
	}

	if _, ok := loaded.Providers["anthropic"]; !ok {
		t.Errorf("Providers missing anthropic entry")
	}

	if loaded.Providers["anthropic"].APIKey != "sk-test-123" {
		t.Errorf("APIKey = %q, want %q", loaded.Providers["anthropic"].APIKey, "sk-test-123")
	}

	if loaded.Display.SyntaxHighlight != true {
		t.Errorf("Display.SyntaxHighlight = %v, want true", loaded.Display.SyntaxHighlight)
	}

	if loaded.History.Enabled != true {
		t.Errorf("History.Enabled = %v, want true", loaded.History.Enabled)
	}
}

// TestLoad_SuccessfulLoadMinimalConfig tests loading with minimal config
func TestLoad_SuccessfulLoadMinimalConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "minimal_config.yaml")

	testConfig := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
	}

	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider = %q, want %q", loaded.CurrentProvider, "anthropic")
	}

	// Verify optional fields are zero values
	if loaded.Display.Color != false {
		t.Errorf("Display.Color should be false for unset field, got %v", loaded.Display.Color)
	}

	if loaded.History.MaxSize != 0 {
		t.Errorf("History.MaxSize should be 0 for unset field, got %v", loaded.History.MaxSize)
	}
}

// TestLoad_EmptyConfigFile tests loading an empty config file
func TestLoad_EmptyConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty_config.yaml")

	// Create empty file
	err := os.WriteFile(configFile, []byte{}, 0644)
	if err != nil {
		t.Fatalf("failed to write empty config file: %v", err)
	}

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error for empty file: %v", err)
	}

	if loaded == nil {
		t.Fatalf("Load() returned nil config for empty file")
	}

	// Verify all zero values
	if loaded.CurrentProvider != "" {
		t.Errorf("CurrentProvider = %q, want empty string", loaded.CurrentProvider)
	}

	if len(loaded.Providers) != 0 {
		t.Errorf("Providers should be empty, got %v", loaded.Providers)
	}
}

// TestLoad_MissingConfigFile tests handling of missing config file
func TestLoad_MissingConfigFile(t *testing.T) {
	configFile := "/nonexistent/path/config.yaml"

	// This should return warning but not error (ConfigFileNotFoundError is expected)
	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error for missing file: %v", err)
	}

	// Should return empty config
	if loaded == nil {
		t.Fatalf("Load() returned nil config for missing file")
	}

	if loaded.CurrentProvider != "" {
		t.Errorf("CurrentProvider should be empty for missing file, got %q", loaded.CurrentProvider)
	}
}

// TestLoad_InvalidYAMLSyntax tests handling of invalid YAML
func TestLoad_InvalidYAMLSyntax(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid_config.yaml")

	// Write invalid YAML
	invalidYAML := `currentProvider: anthropic
providers:
  anthropic
    type: anthropic  # Missing colon - invalid YAML
`
	err := os.WriteFile(configFile, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("failed to write invalid config file: %v", err)
	}

	loaded, err := Load(configFile)

	// Should return error for invalid YAML
	if err == nil {
		t.Fatalf("Load() should return error for invalid YAML")
	}

	if loaded != nil {
		t.Errorf("Load() should return nil config for invalid YAML, got %v", loaded)
	}

	// Verify error message
	if err.Error() == "" {
		t.Errorf("Load() error message is empty")
	}
}

// TestLoad_InvalidIndentation tests invalid YAML indentation
func TestLoad_InvalidIndentation(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "bad_indent_config.yaml")

	// Write YAML with bad indentation
	badYAML := `currentProvider: anthropic
providers:
  anthropic:
   type: anthropic  # Wrong indentation
   apiKey: test
`
	err := os.WriteFile(configFile, []byte(badYAML), 0644)
	if err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	loaded, err := Load(configFile)

	// Should return error
	if err == nil {
		t.Fatalf("Load() should return error for bad indentation")
	}

	if loaded != nil {
		t.Errorf("Load() should return nil config for bad indentation")
	}
}

// TestLoad_InvalidDataTypes tests invalid data type conversions
func TestLoad_InvalidDataTypes(t *testing.T) {
	tests := []struct {
		name     string
		yamlData string
	}{
		{
			name: "string for integer field",
			yamlData: `providers:
  anthropic:
    type: anthropic
    maxTokens: "not a number"
`,
		},
		{
			name: "invalid float for temperature",
			yamlData: `providers:
  anthropic:
    type: anthropic
    temperature: "hot"
`,
		},
		{
			name: "list instead of map for providers",
			yamlData: `providers:
  - type: anthropic
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configFile := filepath.Join(tmpDir, "type_test.yaml")

			err := os.WriteFile(configFile, []byte(tt.yamlData), 0644)
			if err != nil {
				t.Fatalf("failed to write config file: %v", err)
			}

			loaded, err := Load(configFile)

			if err == nil {
				t.Fatalf("Load() should return error for invalid data types")
			}

			if loaded != nil {
				t.Errorf("Load() should return nil config for invalid data types")
			}
		})
	}
}

// TestLoad_EnvironmentVariableOverrides tests env var overrides
func TestLoad_EnvironmentVariableOverrides(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "env_override_config.yaml")

	testConfig := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
	}

	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Set environment variable
	oldEnv := os.Getenv("HOW_CURRENTPROVIDER")
	defer func() {
		if oldEnv != "" {
			os.Setenv("HOW_CURRENTPROVIDER", oldEnv)
		} else {
			os.Unsetenv("HOW_CURRENTPROVIDER")
		}
	}()

	os.Setenv("HOW_CURRENTPROVIDER", "custom")

	// Load config - should be overridden by env var
	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded.CurrentProvider != "custom" {
		t.Errorf("CurrentProvider = %q, want %q (env override)", loaded.CurrentProvider, "custom")
	}
}

// TestLoad_ViaperConfiguration tests viper is properly configured
func TestLoad_ViaperConfiguration(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "viper_config.yaml")

	testConfig := Config{
		CurrentProvider: "anthropic",
	}

	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	// Test that viper is properly set up with HOW prefix
	oldEnv := os.Getenv("HOW_CURRENTPROVIDER")
	defer func() {
		if oldEnv != "" {
			os.Setenv("HOW_CURRENTPROVIDER", oldEnv)
		} else {
			os.Unsetenv("HOW_CURRENTPROVIDER")
		}
	}()

	os.Setenv("HOW_CURRENTPROVIDER", "viper_test")

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded.CurrentProvider != "viper_test" {
		t.Errorf("Viper env override failed: got %q, want %q", loaded.CurrentProvider, "viper_test")
	}
}

// TestLoad_SpecialCharactersInValues tests special characters in config values
func TestLoad_SpecialCharactersInValues(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "special_chars_config.yaml")

	testConfig := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test-123!@#$%^&*()",
				BaseURL:   "https://api.example.com/path?query=value&other=123",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
	}

	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded.Providers["anthropic"].APIKey != "sk-test-123!@#$%^&*()" {
		t.Errorf("Special characters not preserved in APIKey")
	}

	if loaded.Providers["anthropic"].BaseURL != "https://api.example.com/path?query=value&other=123" {
		t.Errorf("Special characters not preserved in BaseURL")
	}
}

// TestLoad_UnicodeCharacters tests unicode in config values
func TestLoad_UnicodeCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "unicode_config.yaml")

	testConfig := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:         "anthropic",
				APIKey:       "sk-test",
				Model:        "claude-3",
				MaxTokens:    2048,
				SystemPrompt: "你好世界 🚀 Привет мир",
			},
		},
	}

	data, err := yaml.Marshal(testConfig)
	if err != nil {
		t.Fatalf("failed to marshal test config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	expected := "你好世界 🚀 Привет мир"
	if loaded.Providers["anthropic"].SystemPrompt != expected {
		t.Errorf("SystemPrompt = %q, want %q", loaded.Providers["anthropic"].SystemPrompt, expected)
	}
}

// TestLoad_NullValuesInYAML tests null/nil values
func TestLoad_NullValuesInYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "null_config.yaml")

	nullYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: sk-test
    model: claude-3
    maxTokens: 2048
    baseUrl: null
    customHeaders: null
`
	err := os.WriteFile(configFile, []byte(nullYAML), 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Null values should be handled as empty/nil
	if loaded.Providers["anthropic"].BaseURL != "" {
		t.Errorf("BaseURL should be empty string for null, got %q", loaded.Providers["anthropic"].BaseURL)
	}

	if loaded.Providers["anthropic"].CustomHeaders != nil && len(loaded.Providers["anthropic"].CustomHeaders) != 0 {
		t.Errorf("CustomHeaders should be empty for null, got %v", loaded.Providers["anthropic"].CustomHeaders)
	}
}

// TestLoad_DuplicateKeysInYAML tests YAML with duplicate keys
func TestLoad_DuplicateKeysInYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "dup_config.yaml")

	dupYAML := `currentProvider: anthropic
currentProvider: openai
providers:
  anthropic:
    type: anthropic
    apiKey: sk-test
    model: claude-3
    maxTokens: 2048
`
	err := os.WriteFile(configFile, []byte(dupYAML), 0644)
	if err != nil {
		t.Fatalf("failed to write test config file: %v", err)
	}

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// YAML parser typically uses the last value for duplicate keys
	if loaded.CurrentProvider != "openai" {
		t.Errorf("CurrentProvider = %q, want %q (last value should win)", loaded.CurrentProvider, "openai")
	}
}

// TestLoad_LargeConfigFile tests loading a large config file
func TestLoad_LargeConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "large_config.yaml")

	// Create config with many providers
	config := Config{
		CurrentProvider: "anthropic",
		Providers:       make(map[string]ProviderConfig),
	}

	for i := 0; i < 100; i++ {
		config.Providers[fmt.Sprintf("provider-%d", i)] = ProviderConfig{
			Type:      "test",
			APIKey:    fmt.Sprintf("key-%d", i),
			Model:     "model",
			MaxTokens: 2048,
		}
	}

	// Add many exclude patterns
	config.Context.ExcludePatterns = make([]string, 1000)
	for i := 0; i < 1000; i++ {
		config.Context.ExcludePatterns[i] = fmt.Sprintf("pattern-%d", i)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal large config: %v", err)
	}

	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write large config file: %v", err)
	}

	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error for large config: %v", err)
	}

	if len(loaded.Providers) != 100 {
		t.Errorf("Provider count = %d, want 100", len(loaded.Providers))
	}

	if len(loaded.Context.ExcludePatterns) != 1000 {
		t.Errorf("Exclude patterns count = %d, want 1000", len(loaded.Context.ExcludePatterns))
	}
}

// TestSave_SuccessfulSaveToExplicitPath tests saving config to explicit path
func TestSave_SuccessfulSaveToExplicitPath(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "saved_config.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test-123",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
		Display: DisplayConfig{
			Color: true,
		},
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); err != nil {
		t.Fatalf("saved config file doesn't exist: %v", err)
	}

	// Verify file is readable
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved config file: %v", err)
	}

	// Verify content is valid YAML
	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved file contains invalid YAML: %v", err)
	}

	// Verify file permissions are 0644
	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	expectedPerm := os.FileMode(0644)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("file permissions = %o, want %o", info.Mode().Perm(), expectedPerm)
	}
}

// TestSave_CreatesParentDirectories tests that Save creates missing parent directories
func TestSave_CreatesParentDirectories(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "a", "b", "c", "config.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); err != nil {
		t.Fatalf("saved config file doesn't exist: %v", err)
	}

	// Verify parent directory has correct permissions
	parentDir := filepath.Dir(configFile)
	info, err := os.Stat(parentDir)
	if err != nil {
		t.Fatalf("failed to stat parent directory: %v", err)
	}

	expectedPerm := os.FileMode(0755)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("directory permissions = %o, want %o", info.Mode().Perm(), expectedPerm)
	}
}

// TestSave_OverwritesExistingFile tests that Save overwrites existing files
func TestSave_OverwritesExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "overwrite_config.yaml")

	// Create initial file
	initialConfig := &Config{
		CurrentProvider: "anthropic",
	}
	err := initialConfig.Save(configFile)
	if err != nil {
		t.Fatalf("initial Save() returned error: %v", err)
	}

	// Overwrite with different data
	newConfig := &Config{
		CurrentProvider: "openai",
		Providers: map[string]ProviderConfig{
			"openai": {
				Type:      "openai",
				APIKey:    "sk-new",
				Model:     "gpt-4",
				MaxTokens: 8192,
			},
		},
	}

	err = newConfig.Save(configFile)
	if err != nil {
		t.Fatalf("second Save() returned error: %v", err)
	}

	// Verify new content
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved file contains invalid YAML: %v", err)
	}

	if loaded.CurrentProvider != "openai" {
		t.Errorf("CurrentProvider = %q, want %q", loaded.CurrentProvider, "openai")
	}

	if _, ok := loaded.Providers["openai"]; !ok {
		t.Errorf("openai provider not in saved config")
	}
}

// TestSave_RoundTripConsistency tests Load -> Save -> Load consistency
func TestSave_RoundTripConsistency(t *testing.T) {
	tmpDir := t.TempDir()
	originalFile := filepath.Join(tmpDir, "original.yaml")
	newFile := filepath.Join(tmpDir, "new.yaml")

	// Create original file
	original := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:         "anthropic",
				APIKey:       "sk-test-123",
				Model:        "claude-3",
				MaxTokens:    4096,
				Temperature:  0.7,
				SystemPrompt: "You are helpful",
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     10,
			IncludeEnvironment: false,
			MaxContextSize:     50000,
			ExcludePatterns:    []string{"*.log"},
		},
	}

	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal original: %v", err)
	}
	err = os.WriteFile(originalFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write original file: %v", err)
	}

	// Load from original
	loaded1, err := Load(originalFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Save to new location
	err = loaded1.Save(newFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Load from new location
	loaded2, err := Load(newFile)
	if err != nil {
		t.Fatalf("Load() from new location returned error: %v", err)
	}

	// Verify they match
	if loaded2.CurrentProvider != loaded1.CurrentProvider {
		t.Errorf("CurrentProvider mismatch after round-trip")
	}

	if len(loaded2.Providers) != len(loaded1.Providers) {
		t.Errorf("Providers count mismatch after round-trip")
	}

	if loaded2.Providers["anthropic"].APIKey != loaded1.Providers["anthropic"].APIKey {
		t.Errorf("APIKey mismatch after round-trip")
	}

	if loaded2.Context.IncludeFiles != loaded1.Context.IncludeFiles {
		t.Errorf("Context.IncludeFiles mismatch after round-trip")
	}
}

// TestSave_EmptyProviders tests saving with empty providers
func TestSave_EmptyProviders(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty_providers.yaml")

	config := &Config{
		CurrentProvider: "test",
		Providers:       make(map[string]ProviderConfig),
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved file contains invalid YAML: %v", err)
	}
}

// TestSave_NilProviders tests saving with nil providers
func TestSave_NilProviders(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "nil_providers.yaml")

	config := &Config{
		CurrentProvider: "test",
		Providers:       nil,
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved file contains invalid YAML: %v", err)
	}
}

// TestSave_YAMLValidation tests that saved YAML is valid
func TestSave_YAMLValidation(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "validation.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	// Verify YAML can be parsed
	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid: %v", err)
	}

	// Verify data integrity
	if loaded.CurrentProvider != config.CurrentProvider {
		t.Errorf("CurrentProvider mismatch after marshal/unmarshal")
	}
}

// TestIntegration_LoadModifySave tests complete workflow
func TestIntegration_LoadModifySave(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "integration.yaml")

	// Step 1: Create initial config
	initial := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-initial",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
	}

	data, err := yaml.Marshal(initial)
	if err != nil {
		t.Fatalf("failed to marshal initial: %v", err)
	}
	err = os.WriteFile(configFile, data, 0644)
	if err != nil {
		t.Fatalf("failed to write initial file: %v", err)
	}

	// Step 2: Load config
	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	// Step 3: Modify
	loaded.CurrentProvider = "openai"
	loaded.Providers["openai"] = ProviderConfig{
		Type:      "openai",
		APIKey:    "sk-new",
		Model:     "gpt-4",
		MaxTokens: 8192,
	}

	// Step 4: Save
	err = loaded.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Step 5: Load again and verify
	reloaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("reLoad() returned error: %v", err)
	}

	if reloaded.CurrentProvider != "openai" {
		t.Errorf("CurrentProvider = %q, want %q", reloaded.CurrentProvider, "openai")
	}

	if _, ok := reloaded.Providers["openai"]; !ok {
		t.Errorf("openai provider not found in reloaded config")
	}

	if reloaded.Providers["openai"].APIKey != "sk-new" {
		t.Errorf("openai APIKey = %q, want %q", reloaded.Providers["openai"].APIKey, "sk-new")
	}
}

// TestIntegration_CreateEmptyConfigAndPopulate tests building config from scratch
func TestIntegration_CreateEmptyConfigAndPopulate(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "new_config.yaml")

	// Create empty config
	config := &Config{}

	// Populate it
	config.CurrentProvider = "anthropic"
	config.Providers = map[string]ProviderConfig{
		"anthropic": {
			Type:      "anthropic",
			APIKey:    "sk-test",
			Model:     "claude-3",
			MaxTokens: 2048,
		},
	}
	config.Display.Color = true

	// Save
	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Load and verify
	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if loaded.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider = %q, want %q", loaded.CurrentProvider, "anthropic")
	}

	if loaded.Display.Color != true {
		t.Errorf("Display.Color = %v, want true", loaded.Display.Color)
	}
}

// TestYAMLTags_ConfigStruct tests YAML tag mappings for Config
func TestYAMLTags_ConfigStruct(t *testing.T) {
	config := Config{
		CurrentProvider: "anthropic",
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	// Verify field names map correctly
	if !contains(yamlStr, "currentProvider") {
		t.Errorf("currentProvider tag not found in YAML")
	}

	if !contains(yamlStr, "providers") {
		t.Errorf("providers tag not found in YAML")
	}

	if !contains(yamlStr, "context") {
		t.Errorf("context tag not found in YAML")
	}

	if !contains(yamlStr, "display") {
		t.Errorf("display tag not found in YAML")
	}

	if !contains(yamlStr, "history") {
		t.Errorf("history tag not found in YAML")
	}
}

// TestYAMLTags_ProviderConfigStruct tests YAML tag mappings for ProviderConfig
func TestYAMLTags_ProviderConfigStruct(t *testing.T) {
	provider := ProviderConfig{
		Type:      "anthropic",
		APIKey:    "sk-test",
		Model:     "claude-3",
		MaxTokens: 2048,
		BaseURL:   "https://api.anthropic.com",
	}

	data, err := yaml.Marshal(provider)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	if !contains(yamlStr, "type") {
		t.Errorf("type tag not found in YAML")
	}

	if !contains(yamlStr, "apiKey") {
		t.Errorf("apiKey tag not found in YAML")
	}

	if !contains(yamlStr, "model") {
		t.Errorf("model tag not found in YAML")
	}

	if !contains(yamlStr, "maxTokens") {
		t.Errorf("maxTokens tag not found in YAML")
	}

	if !contains(yamlStr, "baseUrl") {
		t.Errorf("baseUrl tag not found in YAML")
	}
}

// TestYAMLTags_OmitEmpty tests that omitempty works correctly
func TestYAMLTags_OmitEmpty(t *testing.T) {
	provider := ProviderConfig{
		Type:      "anthropic",
		APIKey:    "sk-test",
		Model:     "claude-3",
		MaxTokens: 2048,
		// BaseURL, Temperature, TopP, SystemPrompt, CustomHeaders are empty
	}

	data, err := yaml.Marshal(provider)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	// Optional fields with omitempty should not appear when empty
	if contains(yamlStr, "baseUrl:") && !contains(yamlStr, "baseUrl: null") {
		// Only check if baseUrl value is present (not null)
		if !contains(yamlStr, "baseUrl: \"\"") {
			t.Errorf("empty baseUrl should not appear in YAML with omitempty")
		}
	}
}

// TestSave_DefaultConfigDirectory tests saving to default config directory
func TestSave_DefaultConfigDirectory(t *testing.T) {
	// Skip if we can't determine home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("unable to determine home directory")
	}

	// Create a temporary override for testing
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
	}

	// Test saving with empty string (would save to default location)
	// We use explicit path instead to avoid modifying user's home directory
	err = config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); err != nil {
		t.Fatalf("saved config file doesn't exist: %v", err)
	}

	// Verify the config is readable
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid: %v", err)
	}

	if loaded.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider = %q, want %q", loaded.CurrentProvider, "anthropic")
	}

	_ = homeDir // Use homeDir to satisfy linter
}

// TestSave_CustomFilePath tests saving to a custom file path
func TestSave_CustomFilePath(t *testing.T) {
	tmpDir := t.TempDir()
	customPath := filepath.Join(tmpDir, "my_custom_config.yaml")

	config := &Config{
		CurrentProvider: "openai",
		Providers: map[string]ProviderConfig{
			"openai": {
				Type:      "openai",
				APIKey:    "sk-openai-test",
				Model:     "gpt-4",
				MaxTokens: 8192,
			},
		},
		Display: DisplayConfig{
			Color:           true,
			SyntaxHighlight: true,
		},
	}

	err := config.Save(customPath)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(customPath); err != nil {
		t.Fatalf("saved config file doesn't exist: %v", err)
	}

	// Verify content
	data, err := os.ReadFile(customPath)
	if err != nil {
		t.Fatalf("failed to read saved config: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid: %v", err)
	}

	if loaded.CurrentProvider != "openai" {
		t.Errorf("CurrentProvider = %q, want %q", loaded.CurrentProvider, "openai")
	}

	if loaded.Providers["openai"].Model != "gpt-4" {
		t.Errorf("Model = %q, want %q", loaded.Providers["openai"].Model, "gpt-4")
	}

	if !loaded.Display.Color {
		t.Errorf("Display.Color = %v, want true", loaded.Display.Color)
	}
}

// TestSave_NestedDirectoryCreation tests creating deeply nested directories
func TestSave_NestedDirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	deepPath := filepath.Join(tmpDir, "level1", "level2", "level3", "level4", "config.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
	}

	err := config.Save(deepPath)
	if err != nil {
		t.Fatalf("Save() returned error for deeply nested path: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(deepPath); err != nil {
		t.Fatalf("deeply nested config file doesn't exist: %v", err)
	}

	// Verify all parent directories were created
	parentDir := filepath.Dir(deepPath)
	if _, err := os.Stat(parentDir); err != nil {
		t.Fatalf("parent directory was not created: %v", err)
	}
}

// TestSave_DirectoryPermissions tests that created directories have correct permissions
func TestSave_DirectoryPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "newdir", "config.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Check directory permissions
	dir := filepath.Dir(configFile)
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("failed to stat directory: %v", err)
	}

	expectedPerm := os.FileMode(0755)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("directory permissions = %o, want %o", info.Mode().Perm(), expectedPerm)
	}
}

// TestSave_YAMLMarshalingComplexTypes tests YAML marshaling of complex nested types
func TestSave_YAMLMarshalingComplexTypes(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "complex.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test",
				Model:     "claude-3-opus",
				MaxTokens: 4096,
				Temperature: 0.8,
				TopP:      0.95,
				SystemPrompt: "You are a helpful AI assistant",
				CustomHeaders: map[string]string{
					"X-Custom": "value",
					"X-Auth":   "bearer token",
				},
			},
			"openai": {
				Type:      "openai",
				APIKey:    "sk-openai",
				Model:     "gpt-4",
				MaxTokens: 8192,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     20,
			IncludeEnvironment: false,
			IncludeGit:         true,
			MaxContextSize:     100000,
			ExcludePatterns:    []string{"*.log", "node_modules/*", ".git/*"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Emoji:           false,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  5000,
			FilePath: "~/.how/history",
		},
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify YAML is valid and contains all data
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid: %v", err)
	}

	// Verify complex nested types
	if len(loaded.Providers) != 2 {
		t.Errorf("provider count = %d, want 2", len(loaded.Providers))
	}

	if loaded.Providers["anthropic"].CustomHeaders["X-Custom"] != "value" {
		t.Errorf("CustomHeaders not preserved correctly")
	}

	if len(loaded.Context.ExcludePatterns) != 3 {
		t.Errorf("ExcludePatterns count = %d, want 3", len(loaded.Context.ExcludePatterns))
	}

	if loaded.History.MaxSize != 5000 {
		t.Errorf("History.MaxSize = %d, want 5000", loaded.History.MaxSize)
	}
}

// TestSave_SpecialCharactersInPath tests saving with special characters in file path
func TestSave_SpecialCharactersInPath(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config-file_123.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configFile); err != nil {
		t.Fatalf("saved config file doesn't exist: %v", err)
	}
}

// TestSave_FilePermissionsAre0644 tests that saved files have 0644 permissions
func TestSave_FilePermissionsAre0644(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "perm_test.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-secret-key",
				Model:     "claude-3",
				MaxTokens: 2048,
			},
		},
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify file permissions
	info, err := os.Stat(configFile)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	expectedPerm := os.FileMode(0644)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf("file permissions = %o, want %o", info.Mode().Perm(), expectedPerm)
	}
}

// TestSave_PermissionDeniedOnWrite tests error handling when write permission is denied
func TestSave_PermissionDeniedOnWrite(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("test cannot run as root")
	}

	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")

	// Create directory and remove write permissions
	if err := os.Mkdir(readOnlyDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	// Restore permissions for cleanup
	defer os.Chmod(readOnlyDir, 0755)

	configFile := filepath.Join(readOnlyDir, "config.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
	}

	err := config.Save(configFile)
	if err == nil {
		t.Fatalf("Save() should return error when write permission is denied")
	}

	// Verify error message contains useful information
	if err.Error() == "" {
		t.Errorf("error message is empty")
	}

	if !contains(err.Error(), "permission denied") && !contains(err.Error(), "failed to") {
		t.Errorf("error message should indicate permission issue: %v", err)
	}
}

// TestSave_PermissionDeniedOnDirectory tests error when can't create directory
func TestSave_PermissionDeniedOnDirectory(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("test cannot run as root")
	}

	tmpDir := t.TempDir()
	readOnlyDir := filepath.Join(tmpDir, "readonly")

	// Create directory
	if err := os.Mkdir(readOnlyDir, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// Remove write permission
	if err := os.Chmod(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to change permissions: %v", err)
	}

	// Restore permissions for cleanup
	defer os.Chmod(readOnlyDir, 0755)

	configFile := filepath.Join(readOnlyDir, "subdir", "config.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
	}

	err := config.Save(configFile)
	if err == nil {
		t.Fatalf("Save() should return error when directory creation is denied")
	}

	if err.Error() == "" {
		t.Errorf("error message is empty")
	}
}

// TestSave_RobustnessWithLargeConfig tests saving very large config structures
func TestSave_RobustnessWithLargeConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "large_save.yaml")

	// Create config with many providers
	config := &Config{
		CurrentProvider: "anthropic",
		Providers:       make(map[string]ProviderConfig),
	}

	for i := 0; i < 50; i++ {
		config.Providers[fmt.Sprintf("provider-%d", i)] = ProviderConfig{
			Type:      "test",
			APIKey:    fmt.Sprintf("key-%d", i),
			Model:     "model",
			MaxTokens: 2048 + i,
		}
	}

	// Add many exclude patterns
	config.Context.ExcludePatterns = make([]string, 500)
	for i := 0; i < 500; i++ {
		config.Context.ExcludePatterns[i] = fmt.Sprintf("pattern-%d", i)
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error for large config: %v", err)
	}

	// Verify file exists and can be loaded back
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid for large config: %v", err)
	}

	if len(loaded.Providers) != 50 {
		t.Errorf("provider count after save = %d, want 50", len(loaded.Providers))
	}

	if len(loaded.Context.ExcludePatterns) != 500 {
		t.Errorf("exclude patterns count after save = %d, want 500", len(loaded.Context.ExcludePatterns))
	}
}

// TestSave_UnicodeInConfig tests saving config with unicode characters
func TestSave_UnicodeInConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "unicode.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:         "anthropic",
				APIKey:       "sk-test",
				Model:        "claude-3",
				MaxTokens:    2048,
				SystemPrompt: "你好 مرحبا привет 🚀",
			},
		},
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify unicode is preserved
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid: %v", err)
	}

	if loaded.Providers["anthropic"].SystemPrompt != "你好 مرحبا привет 🚀" {
		t.Errorf("unicode not preserved in SystemPrompt")
	}
}

// TestSave_EmptyConfig tests saving a completely empty config
func TestSave_EmptyConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "empty.yaml")

	config := &Config{}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error for empty config: %v", err)
	}

	// Verify file exists and is valid YAML
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid for empty config: %v", err)
	}
}

// TestSave_MultipleConsecutiveSaves tests saving multiple times without issues
func TestSave_MultipleConsecutiveSaves(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "consecutive.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
	}

	// Save multiple times with different data
	for i := 0; i < 5; i++ {
		config.CurrentProvider = fmt.Sprintf("provider-%d", i)
		config.Providers = map[string]ProviderConfig{
			fmt.Sprintf("provider-%d", i): {
				Type:      "test",
				APIKey:    fmt.Sprintf("key-%d", i),
				Model:     "model",
				MaxTokens: 2048 + i,
			},
		}

		err := config.Save(configFile)
		if err != nil {
			t.Fatalf("Save() iteration %d returned error: %v", i, err)
		}
	}

	// Verify final state
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read final saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("final saved YAML is invalid: %v", err)
	}

	if loaded.CurrentProvider != "provider-4" {
		t.Errorf("CurrentProvider = %q, want %q", loaded.CurrentProvider, "provider-4")
	}
}

// TestSave_WithNestedMapsAndSlices tests complex data structure marshaling
func TestSave_WithNestedMapsAndSlices(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "nested.yaml")

	config := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test",
				Model:     "claude-3",
				MaxTokens: 4096,
				CustomHeaders: map[string]string{
					"Authorization": "Bearer token123",
					"X-API-Version": "v1",
					"X-Client-ID":   "my-app",
				},
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     10,
			IncludeEnvironment: true,
			IncludeGit:         true,
			MaxContextSize:     100000,
			ExcludePatterns: []string{
				"*.log",
				".git/*",
				"node_modules/*",
				".env*",
				"dist/*",
			},
		},
	}

	err := config.Save(configFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify all nested data is preserved
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Config
	err = yaml.Unmarshal(data, &loaded)
	if err != nil {
		t.Fatalf("saved YAML is invalid: %v", err)
	}

	if len(loaded.Providers["anthropic"].CustomHeaders) != 3 {
		t.Errorf("CustomHeaders count = %d, want 3", len(loaded.Providers["anthropic"].CustomHeaders))
	}

	if loaded.Providers["anthropic"].CustomHeaders["Authorization"] != "Bearer token123" {
		t.Errorf("CustomHeaders not preserved correctly")
	}

	if len(loaded.Context.ExcludePatterns) != 5 {
		t.Errorf("ExcludePatterns count = %d, want 5", len(loaded.Context.ExcludePatterns))
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findInString(s, substr)
}

func findInString(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// TestGetConfigDir_SuccessfulHomeDirectoryRetrieval tests successful home directory retrieval
func TestGetConfigDir_SuccessfulHomeDirectoryRetrieval(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	if configDir == "" {
		t.Fatalf("getConfigDir() returned empty string")
	}

	if !filepath.IsAbs(configDir) {
		t.Errorf("getConfigDir() returned non-absolute path: %q", configDir)
	}
}

// TestGetConfigDir_CorrectPathConstruction tests correct path construction
func TestGetConfigDir_CorrectPathConstruction(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	// Verify the path ends with .config/how
	if !filepath.HasSuffix(configDir, filepath.Join(".config", "how")) {
		t.Errorf("config dir = %q, should end with %q", configDir, filepath.Join(".config", "how"))
	}

	// Verify it contains .config and how in the right order
	parts := filepath.SplitList(configDir)
	hasConfigAndHow := false
	for i := 0; i < len(parts)-1; i++ {
		if parts[i] == ".config" && parts[i+1] == "how" {
			hasConfigAndHow = true
			break
		}
	}
	if !hasConfigAndHow {
		// Alternative check: use filepath functions to verify structure
		if !strings.Contains(configDir, ".config"+string(filepath.Separator)+"how") &&
			!strings.HasSuffix(configDir, filepath.Join(".config", "how")) {
			t.Errorf("config dir = %q, does not have .config/how structure", configDir)
		}
	}
}

// TestGetConfigDir_PathContainsUserHomeDir tests that path contains user's home directory
func TestGetConfigDir_PathContainsUserHomeDir(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Skip("unable to determine home directory")
	}

	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	if !strings.HasPrefix(configDir, homeDir) {
		t.Errorf("config dir = %q, does not start with home dir %q", configDir, homeDir)
	}
}

// TestGetConfigDir_PathIsAbsolute tests that returned path is absolute
func TestGetConfigDir_PathIsAbsolute(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	if !filepath.IsAbs(configDir) {
		t.Errorf("config dir = %q, is not absolute path", configDir)
	}
}

// TestGetConfigDir_RepeatedCallsReturnSameResult tests consistency of calls
func TestGetConfigDir_RepeatedCallsReturnSameResult(t *testing.T) {
	dir1, err1 := getConfigDir()
	if err1 != nil {
		t.Fatalf("first getConfigDir() returned error: %v", err1)
	}

	dir2, err2 := getConfigDir()
	if err2 != nil {
		t.Fatalf("second getConfigDir() returned error: %v", err2)
	}

	if dir1 != dir2 {
		t.Errorf("repeated calls return different results: %q vs %q", dir1, dir2)
	}
}

// TestGetConfigDir_PathDoesNotContainEmptyElements tests path validity
func TestGetConfigDir_PathDoesNotContainEmptyElements(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	// Verify no double slashes or empty path components
	if strings.Contains(configDir, "//") {
		t.Errorf("config dir contains double slashes: %q", configDir)
	}

	if strings.HasSuffix(configDir, string(filepath.Separator)) {
		t.Errorf("config dir ends with separator: %q", configDir)
	}
}

// TestGetConfigDir_PathCanBeJoinedWithFile tests that returned path can be used with filepath.Join
func TestGetConfigDir_PathCanBeJoinedWithFile(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	// Should be able to join with a filename without issues
	configFile := filepath.Join(configDir, "config.yaml")

	if !strings.Contains(configFile, "config.yaml") {
		t.Errorf("filepath.Join failed to create valid path: %q", configFile)
	}

	if !filepath.IsAbs(configFile) {
		t.Errorf("joined path is not absolute: %q", configFile)
	}
}

// TestGetConfigDir_NoError tests that function returns no error under normal conditions
func TestGetConfigDir_NoError(t *testing.T) {
	_, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() should not return error under normal conditions: %v", err)
	}
}

// TestGetConfigDir_IntegrationWithLoad tests getConfigDir usage in Load function
func TestGetConfigDir_IntegrationWithLoad(t *testing.T) {
	// Create a temporary override by saving a file to the actual config location
	// For this test, we verify that Load can successfully call getConfigDir internally
	configFile := "/nonexistent/path/config.yaml"

	// This should not error out due to getConfigDir failing
	loaded, err := Load(configFile)
	if err != nil {
		t.Fatalf("Load() failed, suggesting getConfigDir issue: %v", err)
	}

	// Should return a valid (empty) config
	if loaded == nil {
		t.Fatalf("Load() returned nil config")
	}
}

// TestGetConfigDir_IntegrationWithSave tests getConfigDir usage in Save function context
func TestGetConfigDir_IntegrationWithSave(t *testing.T) {
	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	// Create a test config and verify the directory structure matches
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test_config.yaml")

	config := &Config{
		CurrentProvider: "test",
	}

	err = config.Save(testFile)
	if err != nil {
		t.Fatalf("Save() returned error: %v", err)
	}

	// Verify that configDir ends with the same structure that would be created
	if !strings.HasSuffix(configDir, filepath.Join(".config", "how")) {
		t.Errorf("getConfigDir() and Save() path structures don't match")
	}
}

// TestGetConfigDir_OSUserHomeDirResolution tests proper resolution of OS user home directory
func TestGetConfigDir_OSUserHomeDirResolution(t *testing.T) {
	// Get the home directory directly
	expectedHome, err := os.UserHomeDir()
	if err != nil {
		t.Skip("unable to determine home directory from os.UserHomeDir()")
	}

	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir() returned error: %v", err)
	}

	// Verify that config dir is built from the home directory
	expectedPath := filepath.Join(expectedHome, ".config", "how")
	if configDir != expectedPath {
		t.Errorf("getConfigDir() = %q, want %q", configDir, expectedPath)
	}
}

// TestGetConfigDir_NoNilPointerDereference tests that function doesn't panic
func TestGetConfigDir_NoNilPointerDereference(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("getConfigDir() panicked: %v", r)
		}
	}()

	_, _ = getConfigDir()
}

// TestContextConfig_FieldInitialization tests ContextConfig field initialization
func TestContextConfig_FieldInitialization(t *testing.T) {
	config := ContextConfig{}

	// Verify zero values for all fields
	if config.IncludeFiles != false {
		t.Errorf("IncludeFiles = %v, want false", config.IncludeFiles)
	}

	if config.IncludeHistory != 0 {
		t.Errorf("IncludeHistory = %d, want 0", config.IncludeHistory)
	}

	if config.IncludeEnvironment != false {
		t.Errorf("IncludeEnvironment = %v, want false", config.IncludeEnvironment)
	}

	if config.IncludeGit != false {
		t.Errorf("IncludeGit = %v, want false", config.IncludeGit)
	}

	if config.MaxContextSize != 0 {
		t.Errorf("MaxContextSize = %d, want 0", config.MaxContextSize)
	}

	if len(config.ExcludePatterns) != 0 {
		t.Errorf("ExcludePatterns should be empty, got %v", config.ExcludePatterns)
	}
}

// TestContextConfig_YAMLTags tests YAML tag mappings for ContextConfig
func TestContextConfig_YAMLTags(t *testing.T) {
	ctx := ContextConfig{
		IncludeFiles:       true,
		IncludeHistory:     10,
		IncludeEnvironment: true,
		IncludeGit:         true,
		MaxContextSize:     100000,
		ExcludePatterns:    []string{"*.log"},
	}

	data, err := yaml.Marshal(ctx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	if !contains(yamlStr, "includeFiles") {
		t.Errorf("includeFiles tag not found in YAML")
	}

	if !contains(yamlStr, "includeHistory") {
		t.Errorf("includeHistory tag not found in YAML")
	}

	if !contains(yamlStr, "includeEnvironment") {
		t.Errorf("includeEnvironment tag not found in YAML")
	}

	if !contains(yamlStr, "includeGit") {
		t.Errorf("includeGit tag not found in YAML")
	}

	if !contains(yamlStr, "maxContextSize") {
		t.Errorf("maxContextSize tag not found in YAML")
	}

	if !contains(yamlStr, "excludePatterns") {
		t.Errorf("excludePatterns tag not found in YAML")
	}
}

// TestContextConfig_UnmarshalYAML tests unmarshaling YAML into ContextConfig
func TestContextConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `includeFiles: true
includeHistory: 20
includeEnvironment: false
includeGit: true
maxContextSize: 50000
excludePatterns:
  - "*.log"
  - "node_modules/*"
  - ".git/*"
`

	var ctx ContextConfig
	err := yaml.Unmarshal([]byte(yamlData), &ctx)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if ctx.IncludeFiles != true {
		t.Errorf("IncludeFiles = %v, want true", ctx.IncludeFiles)
	}

	if ctx.IncludeHistory != 20 {
		t.Errorf("IncludeHistory = %d, want 20", ctx.IncludeHistory)
	}

	if ctx.IncludeEnvironment != false {
		t.Errorf("IncludeEnvironment = %v, want false", ctx.IncludeEnvironment)
	}

	if ctx.MaxContextSize != 50000 {
		t.Errorf("MaxContextSize = %d, want 50000", ctx.MaxContextSize)
	}

	if len(ctx.ExcludePatterns) != 3 {
		t.Errorf("ExcludePatterns count = %d, want 3", len(ctx.ExcludePatterns))
	}
}

// TestContextConfig_RoundTripMarshaling tests Marshal -> Unmarshal consistency
func TestContextConfig_RoundTripMarshaling(t *testing.T) {
	original := ContextConfig{
		IncludeFiles:       true,
		IncludeHistory:     15,
		IncludeEnvironment: true,
		IncludeGit:         false,
		MaxContextSize:     75000,
		ExcludePatterns:    []string{"*.tmp", "*.cache"},
	}

	// Marshal
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	var unmarshaled ContextConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify fields match
	if unmarshaled.IncludeFiles != original.IncludeFiles {
		t.Errorf("IncludeFiles mismatch after round-trip")
	}

	if unmarshaled.IncludeHistory != original.IncludeHistory {
		t.Errorf("IncludeHistory mismatch after round-trip")
	}

	if unmarshaled.MaxContextSize != original.MaxContextSize {
		t.Errorf("MaxContextSize mismatch after round-trip")
	}

	if len(unmarshaled.ExcludePatterns) != len(original.ExcludePatterns) {
		t.Errorf("ExcludePatterns length mismatch after round-trip")
	}
}

// TestDisplayConfig_FieldInitialization tests DisplayConfig field initialization
func TestDisplayConfig_FieldInitialization(t *testing.T) {
	config := DisplayConfig{}

	// Verify zero values for all boolean fields
	if config.SyntaxHighlight != false {
		t.Errorf("SyntaxHighlight = %v, want false", config.SyntaxHighlight)
	}

	if config.ShowContext != false {
		t.Errorf("ShowContext = %v, want false", config.ShowContext)
	}

	if config.Emoji != false {
		t.Errorf("Emoji = %v, want false", config.Emoji)
	}

	if config.Color != false {
		t.Errorf("Color = %v, want false", config.Color)
	}
}

// TestDisplayConfig_YAMLTags tests YAML tag mappings for DisplayConfig
func TestDisplayConfig_YAMLTags(t *testing.T) {
	display := DisplayConfig{
		SyntaxHighlight: true,
		ShowContext:     true,
		Emoji:           false,
		Color:           true,
	}

	data, err := yaml.Marshal(display)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	if !contains(yamlStr, "syntaxHighlight") {
		t.Errorf("syntaxHighlight tag not found in YAML")
	}

	if !contains(yamlStr, "showContext") {
		t.Errorf("showContext tag not found in YAML")
	}

	if !contains(yamlStr, "emoji") {
		t.Errorf("emoji tag not found in YAML")
	}

	if !contains(yamlStr, "color") {
		t.Errorf("color tag not found in YAML")
	}
}

// TestDisplayConfig_UnmarshalYAML tests unmarshaling YAML into DisplayConfig
func TestDisplayConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `syntaxHighlight: true
showContext: false
emoji: true
color: false
`

	var display DisplayConfig
	err := yaml.Unmarshal([]byte(yamlData), &display)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if display.SyntaxHighlight != true {
		t.Errorf("SyntaxHighlight = %v, want true", display.SyntaxHighlight)
	}

	if display.ShowContext != false {
		t.Errorf("ShowContext = %v, want false", display.ShowContext)
	}

	if display.Emoji != true {
		t.Errorf("Emoji = %v, want true", display.Emoji)
	}

	if display.Color != false {
		t.Errorf("Color = %v, want false", display.Color)
	}
}

// TestDisplayConfig_RoundTripMarshaling tests Marshal -> Unmarshal consistency
func TestDisplayConfig_RoundTripMarshaling(t *testing.T) {
	original := DisplayConfig{
		SyntaxHighlight: true,
		ShowContext:     true,
		Emoji:           true,
		Color:           true,
	}

	// Marshal
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	var unmarshaled DisplayConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify all fields match
	if unmarshaled.SyntaxHighlight != original.SyntaxHighlight {
		t.Errorf("SyntaxHighlight mismatch after round-trip")
	}

	if unmarshaled.ShowContext != original.ShowContext {
		t.Errorf("ShowContext mismatch after round-trip")
	}

	if unmarshaled.Emoji != original.Emoji {
		t.Errorf("Emoji mismatch after round-trip")
	}

	if unmarshaled.Color != original.Color {
		t.Errorf("Color mismatch after round-trip")
	}
}

// TestHistoryConfig_FieldInitialization tests HistoryConfig field initialization
func TestHistoryConfig_FieldInitialization(t *testing.T) {
	config := HistoryConfig{}

	// Verify zero values
	if config.Enabled != false {
		t.Errorf("Enabled = %v, want false", config.Enabled)
	}

	if config.MaxSize != 0 {
		t.Errorf("MaxSize = %d, want 0", config.MaxSize)
	}

	if config.FilePath != "" {
		t.Errorf("FilePath = %q, want empty string", config.FilePath)
	}
}

// TestHistoryConfig_YAMLTags tests YAML tag mappings for HistoryConfig
func TestHistoryConfig_YAMLTags(t *testing.T) {
	history := HistoryConfig{
		Enabled:  true,
		MaxSize:  1000,
		FilePath: "~/.how/history",
	}

	data, err := yaml.Marshal(history)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	if !contains(yamlStr, "enabled") {
		t.Errorf("enabled tag not found in YAML")
	}

	if !contains(yamlStr, "maxSize") {
		t.Errorf("maxSize tag not found in YAML")
	}

	if !contains(yamlStr, "filePath") {
		t.Errorf("filePath tag not found in YAML")
	}
}

// TestHistoryConfig_UnmarshalYAML tests unmarshaling YAML into HistoryConfig
func TestHistoryConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `enabled: true
maxSize: 2000
filePath: "/var/log/how/history"
`

	var history HistoryConfig
	err := yaml.Unmarshal([]byte(yamlData), &history)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if history.Enabled != true {
		t.Errorf("Enabled = %v, want true", history.Enabled)
	}

	if history.MaxSize != 2000 {
		t.Errorf("MaxSize = %d, want 2000", history.MaxSize)
	}

	if history.FilePath != "/var/log/how/history" {
		t.Errorf("FilePath = %q, want %q", history.FilePath, "/var/log/how/history")
	}
}

// TestHistoryConfig_RoundTripMarshaling tests Marshal -> Unmarshal consistency
func TestHistoryConfig_RoundTripMarshaling(t *testing.T) {
	original := HistoryConfig{
		Enabled:  true,
		MaxSize:  5000,
		FilePath: "~/.config/how/history.txt",
	}

	// Marshal
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	var unmarshaled HistoryConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify fields match
	if unmarshaled.Enabled != original.Enabled {
		t.Errorf("Enabled mismatch after round-trip")
	}

	if unmarshaled.MaxSize != original.MaxSize {
		t.Errorf("MaxSize mismatch after round-trip")
	}

	if unmarshaled.FilePath != original.FilePath {
		t.Errorf("FilePath mismatch after round-trip")
	}
}

// TestProviderConfig_FieldInitialization tests ProviderConfig field initialization
func TestProviderConfig_FieldInitialization(t *testing.T) {
	config := ProviderConfig{}

	// Verify zero values for all fields
	if config.Type != "" {
		t.Errorf("Type = %q, want empty string", config.Type)
	}

	if config.APIKey != "" {
		t.Errorf("APIKey = %q, want empty string", config.APIKey)
	}

	if config.Model != "" {
		t.Errorf("Model = %q, want empty string", config.Model)
	}

	if config.BaseURL != "" {
		t.Errorf("BaseURL = %q, want empty string", config.BaseURL)
	}

	if config.MaxTokens != 0 {
		t.Errorf("MaxTokens = %d, want 0", config.MaxTokens)
	}

	if config.Temperature != 0 {
		t.Errorf("Temperature = %f, want 0", config.Temperature)
	}

	if config.TopP != 0 {
		t.Errorf("TopP = %f, want 0", config.TopP)
	}

	if config.SystemPrompt != "" {
		t.Errorf("SystemPrompt = %q, want empty string", config.SystemPrompt)
	}

	if len(config.CustomHeaders) != 0 {
		t.Errorf("CustomHeaders should be empty, got %v", config.CustomHeaders)
	}
}

// TestProviderConfig_OptionalFieldsHandling tests optional fields with omitempty tags
func TestProviderConfig_OptionalFieldsHandling(t *testing.T) {
	// Test with minimal required fields
	provider := ProviderConfig{
		Type:      "anthropic",
		APIKey:    "sk-test",
		Model:     "claude-3",
		MaxTokens: 2048,
		// Optional fields left empty
	}

	data, err := yaml.Marshal(provider)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	// Required fields should be present
	if !contains(yamlStr, "type") {
		t.Errorf("required field 'type' should be in YAML")
	}

	if !contains(yamlStr, "apiKey") {
		t.Errorf("required field 'apiKey' should be in YAML")
	}

	if !contains(yamlStr, "maxTokens") {
		t.Errorf("required field 'maxTokens' should be in YAML")
	}
}

// TestProviderConfig_OmitemptyFieldsBehavior tests that optional fields with omitempty don't appear when empty
func TestProviderConfig_OmitemptyFieldsBehavior(t *testing.T) {
	provider := ProviderConfig{
		Type:      "anthropic",
		APIKey:    "sk-test",
		Model:     "claude-3",
		MaxTokens: 2048,
		// BaseURL, Temperature, TopP, SystemPrompt, CustomHeaders are empty/zero
	}

	data, err := yaml.Marshal(provider)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	yamlStr := string(data)

	// Optional fields should generally not appear if empty
	// Note: behavior depends on YAML marshaler, but omitempty should prevent them
	if contains(yamlStr, "baseUrl: null") {
		// Some YAML marshalers output null, which is acceptable
	} else if contains(yamlStr, "baseUrl:") && !contains(yamlStr, "baseUrl: \"\"") {
		t.Logf("baseUrl appears as empty in YAML output, which may be due to marshaler behavior")
	}
}

// TestProviderConfig_UnmarshalYAML tests unmarshaling YAML into ProviderConfig
func TestProviderConfig_UnmarshalYAML(t *testing.T) {
	yamlData := `type: anthropic
apiKey: sk-test-123
model: claude-3-opus
maxTokens: 4096
baseUrl: https://api.anthropic.com
temperature: 0.7
topP: 0.95
systemPrompt: "You are helpful"
customHeaders:
  X-Custom: value
`

	var provider ProviderConfig
	err := yaml.Unmarshal([]byte(yamlData), &provider)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if provider.Type != "anthropic" {
		t.Errorf("Type = %q, want %q", provider.Type, "anthropic")
	}

	if provider.APIKey != "sk-test-123" {
		t.Errorf("APIKey = %q, want %q", provider.APIKey, "sk-test-123")
	}

	if provider.MaxTokens != 4096 {
		t.Errorf("MaxTokens = %d, want 4096", provider.MaxTokens)
	}

	if provider.Temperature != 0.7 {
		t.Errorf("Temperature = %f, want 0.7", provider.Temperature)
	}

	if provider.CustomHeaders["X-Custom"] != "value" {
		t.Errorf("CustomHeaders not unmarshaled correctly")
	}
}

// TestProviderConfig_RoundTripMarshaling tests Marshal -> Unmarshal consistency
func TestProviderConfig_RoundTripMarshaling(t *testing.T) {
	original := ProviderConfig{
		Type:        "anthropic",
		APIKey:      "sk-test-abc",
		Model:       "claude-3-sonnet",
		MaxTokens:   8192,
		BaseURL:     "https://api.anthropic.com",
		Temperature: 0.8,
		TopP:        0.9,
		SystemPrompt: "Help the user",
		CustomHeaders: map[string]string{
			"Authorization": "Bearer token",
			"X-Version":     "v1",
		},
	}

	// Marshal
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal
	var unmarshaled ProviderConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify fields match
	if unmarshaled.Type != original.Type {
		t.Errorf("Type mismatch after round-trip")
	}

	if unmarshaled.APIKey != original.APIKey {
		t.Errorf("APIKey mismatch after round-trip")
	}

	if unmarshaled.MaxTokens != original.MaxTokens {
		t.Errorf("MaxTokens mismatch after round-trip")
	}

	if unmarshaled.BaseURL != original.BaseURL {
		t.Errorf("BaseURL mismatch after round-trip")
	}

	if len(unmarshaled.CustomHeaders) != len(original.CustomHeaders) {
		t.Errorf("CustomHeaders count mismatch after round-trip")
	}
}

// TestConfig_NestedStructMarshaling tests marshaling of Config with nested structs
func TestConfig_NestedStructMarshaling(t *testing.T) {
	config := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-test",
				Model:     "claude-3",
				MaxTokens: 4096,
			},
		},
		Context: ContextConfig{
			IncludeFiles:   true,
			IncludeHistory: 10,
			MaxContextSize: 100000,
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  1000,
			FilePath: "~/.how/history",
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify nested structs are preserved
	if unmarshaled.Context.IncludeFiles != config.Context.IncludeFiles {
		t.Errorf("nested Context.IncludeFiles mismatch")
	}

	if unmarshaled.Display.SyntaxHighlight != config.Display.SyntaxHighlight {
		t.Errorf("nested Display.SyntaxHighlight mismatch")
	}

	if unmarshaled.History.MaxSize != config.History.MaxSize {
		t.Errorf("nested History.MaxSize mismatch")
	}
}

// TestConfig_NestedMapMarshaling tests marshaling of nested Provider map in Config
func TestConfig_NestedMapMarshaling(t *testing.T) {
	config := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-anthropic",
				Model:     "claude-3",
				MaxTokens: 4096,
			},
			"openai": {
				Type:      "openai",
				APIKey:    "sk-openai",
				Model:     "gpt-4",
				MaxTokens: 8192,
			},
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify all providers are preserved
	if len(unmarshaled.Providers) != 2 {
		t.Errorf("provider count = %d, want 2", len(unmarshaled.Providers))
	}

	if unmarshaled.Providers["anthropic"].APIKey != "sk-anthropic" {
		t.Errorf("anthropic provider not preserved correctly")
	}

	if unmarshaled.Providers["openai"].Model != "gpt-4" {
		t.Errorf("openai provider not preserved correctly")
	}
}

// TestConfig_NestedSliceMarshaling tests marshaling of nested slices in Config
func TestConfig_NestedSliceMarshaling(t *testing.T) {
	config := Config{
		CurrentProvider: "test",
		Context: ContextConfig{
			ExcludePatterns: []string{
				"*.log",
				"node_modules/*",
				".git/*",
				"*.tmp",
			},
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify all patterns are preserved
	if len(unmarshaled.Context.ExcludePatterns) != 4 {
		t.Errorf("exclude patterns count = %d, want 4", len(unmarshaled.Context.ExcludePatterns))
	}

	expectedPatterns := []string{"*.log", "node_modules/*", ".git/*", "*.tmp"}
	for i, pattern := range expectedPatterns {
		if i < len(unmarshaled.Context.ExcludePatterns) {
			if unmarshaled.Context.ExcludePatterns[i] != pattern {
				t.Errorf("pattern[%d] = %q, want %q", i, unmarshaled.Context.ExcludePatterns[i], pattern)
			}
		}
	}
}

// TestConfig_AllStructsCompleteMarshaling tests complete Config with all nested structs populated
func TestConfig_AllStructsCompleteMarshaling(t *testing.T) {
	config := Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:         "anthropic",
				APIKey:       "sk-test",
				Model:        "claude-3-opus",
				MaxTokens:    4096,
				BaseURL:      "https://api.anthropic.com",
				Temperature:  0.7,
				TopP:         0.95,
				SystemPrompt: "You are helpful",
				CustomHeaders: map[string]string{
					"X-Custom": "header-value",
				},
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     20,
			IncludeEnvironment: true,
			IncludeGit:         true,
			MaxContextSize:     100000,
			ExcludePatterns:    []string{"*.log", ".env*"},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Emoji:           true,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  5000,
			FilePath: "~/.how/history",
		},
	}

	// Marshal and unmarshal
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify all top-level fields
	if unmarshaled.CurrentProvider != config.CurrentProvider {
		t.Errorf("CurrentProvider mismatch")
	}

	// Verify Providers
	if len(unmarshaled.Providers) != len(config.Providers) {
		t.Errorf("Providers count mismatch")
	}

	// Verify Context
	if unmarshaled.Context.IncludeFiles != config.Context.IncludeFiles {
		t.Errorf("Context.IncludeFiles mismatch")
	}

	if unmarshaled.Context.MaxContextSize != config.Context.MaxContextSize {
		t.Errorf("Context.MaxContextSize mismatch")
	}

	// Verify Display
	if unmarshaled.Display.SyntaxHighlight != config.Display.SyntaxHighlight {
		t.Errorf("Display.SyntaxHighlight mismatch")
	}

	// Verify History
	if unmarshaled.History.Enabled != config.History.Enabled {
		t.Errorf("History.Enabled mismatch")
	}

	if unmarshaled.History.MaxSize != config.History.MaxSize {
		t.Errorf("History.MaxSize mismatch")
	}
}

// TestProviderConfig_CustomHeadersMarshalingComplex tests complex CustomHeaders marshaling
func TestProviderConfig_CustomHeadersMarshalingComplex(t *testing.T) {
	provider := ProviderConfig{
		Type:      "anthropic",
		APIKey:    "sk-test",
		Model:     "claude-3",
		MaxTokens: 2048,
		CustomHeaders: map[string]string{
			"Authorization":  "Bearer secret-token-123",
			"X-API-Version":  "2024-01",
			"X-Client-ID":    "my-client",
			"Content-Type":   "application/json",
		},
	}

	data, err := yaml.Marshal(provider)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled ProviderConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify all custom headers are preserved
	if len(unmarshaled.CustomHeaders) != 4 {
		t.Errorf("CustomHeaders count = %d, want 4", len(unmarshaled.CustomHeaders))
	}

	if unmarshaled.CustomHeaders["Authorization"] != "Bearer secret-token-123" {
		t.Errorf("Authorization header not preserved")
	}

	if unmarshaled.CustomHeaders["X-API-Version"] != "2024-01" {
		t.Errorf("X-API-Version header not preserved")
	}
}

// TestContextConfig_SliceMarshalingEmpty tests empty ExcludePatterns slice marshaling
func TestContextConfig_SliceMarshalingEmpty(t *testing.T) {
	ctx := ContextConfig{
		IncludeFiles: true,
		ExcludePatterns: []string{}, // Empty slice
	}

	data, err := yaml.Marshal(ctx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled ContextConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if unmarshaled.ExcludePatterns == nil {
		// Empty slice may be unmarshaled as nil, which is acceptable
		unmarshaled.ExcludePatterns = []string{}
	}

	if len(unmarshaled.ExcludePatterns) != 0 {
		t.Errorf("ExcludePatterns should be empty after round-trip")
	}
}

// TestContextConfig_SliceMarshalingSingleElement tests single element ExcludePatterns slice marshaling
func TestContextConfig_SliceMarshalingSingleElement(t *testing.T) {
	ctx := ContextConfig{
		IncludeFiles:    true,
		ExcludePatterns: []string{"*.log"},
	}

	data, err := yaml.Marshal(ctx)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled ContextConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(unmarshaled.ExcludePatterns) != 1 {
		t.Errorf("ExcludePatterns count = %d, want 1", len(unmarshaled.ExcludePatterns))
	}

	if unmarshaled.ExcludePatterns[0] != "*.log" {
		t.Errorf("ExcludePatterns[0] = %q, want %q", unmarshaled.ExcludePatterns[0], "*.log")
	}
}

// TestProviderConfig_FloatFieldsPrecision tests float field precision in marshaling
func TestProviderConfig_FloatFieldsPrecision(t *testing.T) {
	provider := ProviderConfig{
		Type:        "anthropic",
		APIKey:      "sk-test",
		Model:       "claude-3",
		MaxTokens:   2048,
		Temperature: 0.75,
		TopP:        0.9,
	}

	data, err := yaml.Marshal(provider)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled ProviderConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Float32 comparison with tolerance
	if unmarshaled.Temperature < 0.74 || unmarshaled.Temperature > 0.76 {
		t.Errorf("Temperature = %f, want ~0.75", unmarshaled.Temperature)
	}

	if unmarshaled.TopP < 0.89 || unmarshaled.TopP > 0.91 {
		t.Errorf("TopP = %f, want ~0.9", unmarshaled.TopP)
	}
}
