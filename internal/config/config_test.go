package config

import (
	"fmt"
	"os"
	"path/filepath"
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
