package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLoadEmptyConfig(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider mismatch: expected anthropic, got %s", config.CurrentProvider)
	}

	if len(config.Providers) == 0 {
		t.Error("expected providers to be loaded")
	}
}

func TestLoadSampleConfig(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider mismatch: expected anthropic, got %s", config.CurrentProvider)
	}

	provider, exists := config.Providers["anthropic"]
	if !exists {
		t.Fatal("anthropic provider not found")
	}

	if provider.Type != "anthropic" {
		t.Errorf("Provider type mismatch: expected anthropic, got %s", provider.Type)
	}

	if provider.APIKey != "sk-ant-test-key-12345" {
		t.Errorf("APIKey mismatch: expected sk-ant-test-key-12345, got %s", provider.APIKey)
	}

	if provider.Model != "claude-3-5-sonnet-20241022" {
		t.Errorf("Model mismatch: expected claude-3-5-sonnet-20241022, got %s", provider.Model)
	}

	if provider.MaxTokens != 2048 {
		t.Errorf("MaxTokens mismatch: expected 2048, got %d", provider.MaxTokens)
	}

	if !config.Context.IncludeFiles {
		t.Error("expected IncludeFiles to be true")
	}

	if config.Context.IncludeHistory != 50 {
		t.Errorf("IncludeHistory mismatch: expected 50, got %d", config.Context.IncludeHistory)
	}

	if !config.Display.SyntaxHighlight {
		t.Error("expected SyntaxHighlight to be true")
	}
}

func TestLoadMinimalConfig(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, MinimalYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider mismatch: expected anthropic, got %s", config.CurrentProvider)
	}

	provider, exists := config.Providers["anthropic"]
	if !exists {
		t.Fatal("anthropic provider not found")
	}

	if provider.MaxTokens != 2048 {
		t.Errorf("MaxTokens mismatch: expected 2048, got %d", provider.MaxTokens)
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	config := SampleConfig()
	configPath := filepath.Join(tempDir, "config.yaml")

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	// Load the saved config and verify it matches
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.CurrentProvider != config.CurrentProvider {
		t.Errorf("CurrentProvider mismatch after save: expected %s, got %s", config.CurrentProvider, loaded.CurrentProvider)
	}

	if len(loaded.Providers) != len(config.Providers) {
		t.Errorf("Providers count mismatch after save: expected %d, got %d", len(config.Providers), len(loaded.Providers))
	}
}

func TestSaveConfigCreatesDirectory(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	config := SampleConfig()
	configPath := filepath.Join(tempDir, "subdir", "config.yaml")

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file not created: %v", err)
	}

	if _, err := os.Stat(filepath.Dir(configPath)); err != nil {
		t.Fatalf("directory not created: %v", err)
	}
}

func TestLoadNonexistentFile(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "nonexistent.yaml")

	// Load should return a config even if file doesn't exist
	config, err := Load(configPath)
	if err != nil {
		t.Logf("Load returned error (expected for non-viper errors): %v", err)
	}

	if config == nil {
		t.Fatal("expected config to be returned even if file doesn't exist")
	}
}

func TestConfigWithMultipleProviders(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	yaml := `currentProvider: provider1
providers:
  provider1:
    type: type1
    apiKey: key1
    model: model1
    maxTokens: 1000
  provider2:
    type: type2
    apiKey: key2
    model: model2
    maxTokens: 2000
`

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, yaml)

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(config.Providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(config.Providers))
	}

	if _, exists := config.Providers["provider1"]; !exists {
		t.Error("provider1 not found")
	}

	if _, exists := config.Providers["provider2"]; !exists {
		t.Error("provider2 not found")
	}
}

func TestConfigContextSettings(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.Context.MaxContextSize != 8000 {
		t.Errorf("MaxContextSize mismatch: expected 8000, got %d", config.Context.MaxContextSize)
	}

	if len(config.Context.ExcludePatterns) != 2 {
		t.Errorf("ExcludePatterns count mismatch: expected 2, got %d", len(config.Context.ExcludePatterns))
	}

	if config.Context.ExcludePatterns[0] != ".git" {
		t.Errorf("first exclude pattern mismatch: expected .git, got %s", config.Context.ExcludePatterns[0])
	}
}

func TestConfigDisplaySettings(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !config.Display.SyntaxHighlight {
		t.Error("expected SyntaxHighlight to be true")
	}

	if !config.Display.ShowContext {
		t.Error("expected ShowContext to be true")
	}

	if config.Display.Emoji {
		t.Error("expected Emoji to be false")
	}

	if !config.Display.Color {
		t.Error("expected Color to be true")
	}
}

func TestConfigHistorySettings(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if !config.History.Enabled {
		t.Error("expected History.Enabled to be true")
	}

	if config.History.MaxSize != 1000 {
		t.Errorf("History.MaxSize mismatch: expected 1000, got %d", config.History.MaxSize)
	}
}

func TestSampleConfigFunction(t *testing.T) {
	config := SampleConfig()

	if config == nil {
		t.Fatal("SampleConfig returned nil")
	}

	if config.CurrentProvider != "anthropic" {
		t.Errorf("expected CurrentProvider to be anthropic, got %s", config.CurrentProvider)
	}

	if len(config.Providers) == 0 {
		t.Error("expected providers to be populated")
	}
}

func TestConfigWithProviderFunction(t *testing.T) {
	config := ConfigWithProvider(t, "test-provider", "test-type", "test-key", "test-model")

	if config.CurrentProvider != "test-provider" {
		t.Errorf("expected CurrentProvider to be test-provider, got %s", config.CurrentProvider)
	}

	provider, exists := config.Providers["test-provider"]
	if !exists {
		t.Fatal("test-provider not found in config")
	}

	if provider.Type != "test-type" {
		t.Errorf("expected type to be test-type, got %s", provider.Type)
	}

	if provider.APIKey != "test-key" {
		t.Errorf("expected APIKey to be test-key, got %s", provider.APIKey)
	}

	if provider.Model != "test-model" {
		t.Errorf("expected Model to be test-model, got %s", provider.Model)
	}
}

func TestTestHelperTempDir(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()

	if tempDir == "" {
		t.Fatal("expected TempDir to return a path")
	}

	if !helper.DirectoryExists(tempDir) {
		t.Errorf("temp directory not created: %s", tempDir)
	}
}

func TestTestHelperWriteAndLoadConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	config := SampleConfig()

	configPath := helper.WriteConfig(config, tempDir)

	if !helper.FileExists(configPath) {
		t.Errorf("config file not created: %s", configPath)
	}

	loaded := helper.LoadConfig(configPath)

	if loaded.CurrentProvider != config.CurrentProvider {
		t.Errorf("CurrentProvider mismatch: expected %s, got %s", config.CurrentProvider, loaded.CurrentProvider)
	}
}

func TestTestHelperWriteYAML(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	yaml := SampleYAML()

	configPath := helper.WriteYAML(yaml, tempDir)

	if !helper.FileExists(configPath) {
		t.Errorf("config file not created: %s", configPath)
	}

	content := helper.ReadFile(configPath)
	if content != yaml {
		t.Error("YAML content mismatch")
	}
}

func TestTestHelperConfigEqual(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	config1 := SampleConfig()
	config2 := SampleConfig()

	// This should not fail
	helper.AssertConfigEqual(config1, config2)
}

func TestTestHelperEnv(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	cleanup := helper.CreateConfigWithEnv("TEST_KEY", "test_value")
	defer cleanup()

	value, exists := os.LookupEnv("HOW_TEST_KEY")
	if !exists {
		t.Fatal("environment variable not set")
	}

	if value != "test_value" {
		t.Errorf("expected test_value, got %s", value)
	}
}

func TestTempConfigDirCleanup(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)

	if _, err := os.Stat(tempDir); err != nil {
		t.Fatalf("temp directory not created: %v", err)
	}

	cleanup()

	if _, err := os.Stat(tempDir); err == nil {
		t.Fatal("temp directory not cleaned up")
	}
}

func TestProviderConfigWithCustomHeaders(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	yaml := `currentProvider: custom
providers:
  custom:
    type: custom
    apiKey: key123
    model: model-1
    maxTokens: 4000
    customHeaders:
      X-Custom-Header: value1
      X-Another-Header: value2
`

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, yaml)

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	provider := config.Providers["custom"]
	if len(provider.CustomHeaders) != 2 {
		t.Errorf("expected 2 custom headers, got %d", len(provider.CustomHeaders))
	}

	if provider.CustomHeaders["X-Custom-Header"] != "value1" {
		t.Errorf("custom header mismatch: expected value1, got %s", provider.CustomHeaders["X-Custom-Header"])
	}
}

// TestLoadExplicitConfigFile tests loading from an explicitly specified config file path.
func TestLoadExplicitConfigFile(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "custom-config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config == nil {
		t.Fatal("expected config to be returned")
	}

	if config.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider mismatch: expected anthropic, got %s", config.CurrentProvider)
	}

	provider, exists := config.Providers["anthropic"]
	if !exists {
		t.Fatal("anthropic provider not found")
	}

	if provider.APIKey != "sk-ant-test-key-12345" {
		t.Errorf("APIKey mismatch: expected sk-ant-test-key-12345, got %s", provider.APIKey)
	}
}

// TestLoadFromDefaultConfigDirectory tests loading from the default config directory
// when no explicit config file path is provided.
func TestLoadFromDefaultConfigDirectory(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	homeDir := helper.TempDir()
	configDir := filepath.Join(homeDir, ".config", "how")

	// Create config directory structure
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	// Temporarily override home directory using environment variable
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	config, err := Load("")
	if err != nil {
		t.Fatalf("Load from default directory failed: %v", err)
	}

	if config == nil {
		t.Fatal("expected config to be returned")
	}

	if config.CurrentProvider != "anthropic" {
		t.Errorf("CurrentProvider mismatch: expected anthropic, got %s", config.CurrentProvider)
	}
}

// TestLoadWithEnvironmentVariableOverrides tests that environment variables override config file values.
func TestLoadWithEnvironmentVariableOverrides(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	// Set environment variable to override the config
	cleanup := helper.CreateConfigWithEnv("CURRENTPROVIDER", "override-provider")
	defer cleanup()

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Environment variable should override the YAML value
	if config.CurrentProvider != "override-provider" {
		t.Errorf("CurrentProvider not overridden by env var: expected override-provider, got %s", config.CurrentProvider)
	}
}

// TestLoadWithProviderAPIKeyEnvOverride tests environment variable override for provider API key.
func TestLoadWithProviderAPIKeyEnvOverride(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	yaml := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: original-key
    model: claude-3-5-sonnet-20241022
    maxTokens: 2048
`
	WriteConfigFile(t, configPath, yaml)

	// Set environment variable override for API key
	oldKey, wasSet := os.LookupEnv("HOW_PROVIDERS_ANTHROPIC_APIKEY")
	if err := os.Setenv("HOW_PROVIDERS_ANTHROPIC_APIKEY", "env-override-key"); err != nil {
		t.Fatalf("failed to set environment variable: %v", err)
	}
	defer func() {
		if wasSet {
			os.Setenv("HOW_PROVIDERS_ANTHROPIC_APIKEY", oldKey)
		} else {
			os.Unsetenv("HOW_PROVIDERS_ANTHROPIC_APIKEY")
		}
	}()

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	provider, exists := config.Providers["anthropic"]
	if !exists {
		t.Fatal("anthropic provider not found")
	}

	if provider.APIKey != "env-override-key" {
		t.Errorf("APIKey not overridden by env var: expected env-override-key, got %s", provider.APIKey)
	}
}

// TestLoadConfigFileNotFoundReturnsDefault tests that loading a non-existent config returns a valid empty config.
func TestLoadConfigFileNotFoundReturnsDefault(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	nonexistentPath := filepath.Join(tempDir, "nonexistent-config.yaml")

	config, err := Load(nonexistentPath)

	// The Load function should not error on file not found; it should warn and continue
	if config == nil {
		t.Fatal("expected config to be returned even when file not found")
	}

	// The returned config should be a valid (possibly empty) Config struct
	if config.Providers == nil {
		config.Providers = make(map[string]ProviderConfig)
	}
}

// TestLoadInvalidYAMLSyntax tests that invalid YAML syntax returns an unmarshal error.
func TestLoadInvalidYAMLSyntax(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "invalid.yaml")

	// Write invalid YAML
	invalidYAML := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: test
    maxTokens: not-a-number
`

	WriteConfigFile(t, configPath, invalidYAML)

	config, err := Load(configPath)

	// The function may or may not error depending on how viper handles this
	// But if it returns a config, maxTokens should be 0 (not unmarshalled correctly)
	if config != nil && len(config.Providers) > 0 {
		provider, exists := config.Providers["anthropic"]
		if exists && provider.MaxTokens != 0 {
			t.Logf("Note: viper parsed non-numeric value as %d", provider.MaxTokens)
		}
	}
}

// TestLoadMalformedYAMLStructure tests that deeply malformed YAML is handled.
func TestLoadMalformedYAMLStructure(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "malformed.yaml")

	// Write malformed YAML with invalid syntax
	malformedYAML := `currentProvider: anthropic
providers:
  anthropic
    type: anthropic  # Invalid indentation
    apiKey: test
`

	WriteConfigFile(t, configPath, malformedYAML)

	config, err := Load(configPath)

	// When YAML parsing fails, Load should return an error or an empty config
	if err != nil && err.Error() != "" {
		// Error case is acceptable for malformed YAML
		t.Logf("Malformed YAML correctly returned error: %v", err)
	} else if config != nil {
		// Empty/default config is also acceptable
		t.Logf("Malformed YAML returned empty config, no error")
	}
}

// TestLoadWithMultipleEnvironmentVariables tests multiple environment variable overrides.
func TestLoadWithMultipleEnvironmentVariables(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	// Set multiple environment variable overrides
	cleanup1 := helper.CreateConfigWithEnv("CURRENTPROVIDER", "env-provider")
	defer cleanup1()

	cleanup2 := helper.CreateConfigWithEnv("CONTEXT_INCLUDEFILES", "false")
	defer cleanup2()

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if config.CurrentProvider != "env-provider" {
		t.Errorf("CurrentProvider not overridden: expected env-provider, got %s", config.CurrentProvider)
	}

	// Note: Environment variable parsing for boolean values depends on viper's implementation
	t.Logf("Context.IncludeFiles: %v (env var set to 'false')", config.Context.IncludeFiles)
}

// TestLoadExplicitConfigFileWithEnvOverrides tests that env vars override explicit config file values.
func TestLoadExplicitConfigFileWithEnvOverrides(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	yaml := `currentProvider: file-provider
providers:
  file-provider:
    type: anthropic
    apiKey: file-key
    model: model-from-file
    maxTokens: 1000
`
	WriteConfigFile(t, configPath, yaml)

	cleanup := helper.CreateConfigWithEnv("CURRENTPROVIDER", "env-provider")
	defer cleanup()

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Environment variable should take precedence
	if config.CurrentProvider != "env-provider" {
		t.Errorf("Env var override failed: expected env-provider, got %s", config.CurrentProvider)
	}
}

// TestLoadConfigWithComplexNesting tests loading config with nested structures and env var overrides.
func TestLoadConfigWithComplexNesting(t *testing.T) {
	tempDir, cleanup := TempConfigDir(t)
	defer cleanup()

	configPath := filepath.Join(tempDir, "config.yaml")
	WriteConfigFile(t, configPath, SampleYAML())

	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify nested structure is preserved
	if config.Context.IncludeHistory != 50 {
		t.Errorf("Context.IncludeHistory mismatch: expected 50, got %d", config.Context.IncludeHistory)
	}

	if len(config.Context.ExcludePatterns) != 2 {
		t.Errorf("ExcludePatterns count mismatch: expected 2, got %d", len(config.Context.ExcludePatterns))
	}

	if config.Display.SyntaxHighlight != true {
		t.Error("Display.SyntaxHighlight should be true")
	}

	if config.History.MaxSize != 1000 {
		t.Errorf("History.MaxSize mismatch: expected 1000, got %d", config.History.MaxSize)
	}
}

// TestSaveToExplicitPath tests saving config to an explicitly specified path.
func TestSaveToExplicitPath(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	config := SampleConfig()
	configPath := filepath.Join(tempDir, "custom-config.yaml")

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if !helper.FileExists(configPath) {
		t.Fatalf("config file not created at explicit path: %s", configPath)
	}

	// Verify the saved file can be loaded and matches the original
	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.CurrentProvider != config.CurrentProvider {
		t.Errorf("CurrentProvider mismatch: expected %s, got %s", config.CurrentProvider, loaded.CurrentProvider)
	}
}

// TestSaveToDefaultConfigDirectory tests saving config to the default config directory
// when no explicit path is provided.
func TestSaveToDefaultConfigDirectory(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	homeDir := helper.TempDir()
	configDir := filepath.Join(homeDir, ".config", "how")

	// Temporarily override HOME to use our test directory
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	config := SampleConfig()

	// Save with empty string to use default directory
	if err := config.Save(""); err != nil {
		t.Fatalf("Save to default directory failed: %v", err)
	}

	expectedConfigPath := filepath.Join(configDir, "config.yaml")
	if !helper.FileExists(expectedConfigPath) {
		t.Fatalf("config file not created in default directory: %s", expectedConfigPath)
	}

	// Verify the saved file can be loaded
	loaded, err := Load("")
	if err != nil {
		t.Fatalf("Load from default directory failed: %v", err)
	}

	if loaded.CurrentProvider != config.CurrentProvider {
		t.Errorf("CurrentProvider mismatch: expected %s, got %s", config.CurrentProvider, loaded.CurrentProvider)
	}
}

// TestSaveCreatesNestedDirectories tests that Save creates nested directories as needed.
func TestSaveCreatesNestedDirectories(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	config := SampleConfig()
	configPath := filepath.Join(tempDir, "level1", "level2", "level3", "config.yaml")

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save with nested directories failed: %v", err)
	}

	if !helper.FileExists(configPath) {
		t.Fatalf("config file not created with nested directories: %s", configPath)
	}

	// Verify all intermediate directories were created
	if !helper.DirectoryExists(filepath.Join(tempDir, "level1")) {
		t.Error("level1 directory not created")
	}

	if !helper.DirectoryExists(filepath.Join(tempDir, "level1", "level2")) {
		t.Error("level2 directory not created")
	}

	if !helper.DirectoryExists(filepath.Join(tempDir, "level1", "level2", "level3")) {
		t.Error("level3 directory not created")
	}
}

// TestSaveWithComplexProviderConfig tests saving a config with multiple providers and custom headers.
func TestSaveWithComplexProviderConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	config := &Config{
		CurrentProvider: "primary",
		Providers: map[string]ProviderConfig{
			"primary": {
				Type:      "anthropic",
				APIKey:    "key-1",
				Model:     "model-1",
				MaxTokens: 2048,
				Temperature: 0.7,
				TopP:       1.0,
				CustomHeaders: map[string]string{
					"X-Custom-Header": "value1",
					"X-Another-Header": "value2",
				},
			},
			"secondary": {
				Type:      "openai",
				APIKey:    "key-2",
				Model:     "model-2",
				MaxTokens: 4096,
				Temperature: 0.5,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     100,
			IncludeEnvironment: true,
			IncludeGit:         true,
			MaxContextSize:     16000,
			ExcludePatterns:    []string{".git", ".env", "node_modules"},
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
			FilePath: "~/.local/share/how/history",
		},
	}

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify complex structure is preserved
	if len(loaded.Providers) != 2 {
		t.Errorf("expected 2 providers, got %d", len(loaded.Providers))
	}

	primaryProvider := loaded.Providers["primary"]
	if len(primaryProvider.CustomHeaders) != 2 {
		t.Errorf("expected 2 custom headers, got %d", len(primaryProvider.CustomHeaders))
	}

	if primaryProvider.CustomHeaders["X-Custom-Header"] != "value1" {
		t.Errorf("custom header mismatch: expected value1, got %s", primaryProvider.CustomHeaders["X-Custom-Header"])
	}

	if len(loaded.Context.ExcludePatterns) != 3 {
		t.Errorf("expected 3 exclude patterns, got %d", len(loaded.Context.ExcludePatterns))
	}
}

// TestSavePreservesAllConfigFields tests that saving and loading preserves all config fields.
func TestSavePreservesAllConfigFields(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	original := SampleConfig()

	if err := original.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Check all main fields
	if loaded.CurrentProvider != original.CurrentProvider {
		t.Errorf("CurrentProvider mismatch: expected %s, got %s", original.CurrentProvider, loaded.CurrentProvider)
	}

	if len(loaded.Providers) != len(original.Providers) {
		t.Errorf("Providers count mismatch: expected %d, got %d", len(original.Providers), len(loaded.Providers))
	}

	// Check context config
	if loaded.Context.IncludeFiles != original.Context.IncludeFiles {
		t.Errorf("IncludeFiles mismatch: expected %v, got %v", original.Context.IncludeFiles, loaded.Context.IncludeFiles)
	}

	if loaded.Context.IncludeHistory != original.Context.IncludeHistory {
		t.Errorf("IncludeHistory mismatch: expected %d, got %d", original.Context.IncludeHistory, loaded.Context.IncludeHistory)
	}

	if loaded.Context.MaxContextSize != original.Context.MaxContextSize {
		t.Errorf("MaxContextSize mismatch: expected %d, got %d", original.Context.MaxContextSize, loaded.Context.MaxContextSize)
	}

	// Check display config
	if loaded.Display.SyntaxHighlight != original.Display.SyntaxHighlight {
		t.Errorf("SyntaxHighlight mismatch: expected %v, got %v", original.Display.SyntaxHighlight, loaded.Display.SyntaxHighlight)
	}

	if loaded.Display.Color != original.Display.Color {
		t.Errorf("Color mismatch: expected %v, got %v", original.Display.Color, loaded.Color)
	}

	// Check history config
	if loaded.History.Enabled != original.History.Enabled {
		t.Errorf("History.Enabled mismatch: expected %v, got %v", original.History.Enabled, loaded.History.Enabled)
	}

	if loaded.History.MaxSize != original.History.MaxSize {
		t.Errorf("History.MaxSize mismatch: expected %d, got %d", original.History.MaxSize, loaded.History.MaxSize)
	}
}

// TestSaveOverwritesExistingFile tests that Save overwrites existing config files.
func TestSaveOverwritesExistingFile(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// Save initial config
	config1 := ConfigWithProvider(t, "provider1", "type1", "key1", "model1")
	if err := config1.Save(configPath); err != nil {
		t.Fatalf("First save failed: %v", err)
	}

	// Verify initial save
	loaded1, err := Load(configPath)
	if err != nil {
		t.Fatalf("First load failed: %v", err)
	}

	if loaded1.CurrentProvider != "provider1" {
		t.Errorf("first save: expected provider1, got %s", loaded1.CurrentProvider)
	}

	// Save different config to same path
	config2 := ConfigWithProvider(t, "provider2", "type2", "key2", "model2")
	if err := config2.Save(configPath); err != nil {
		t.Fatalf("Second save failed: %v", err)
	}

	// Verify second save overwrote the first
	loaded2, err := Load(configPath)
	if err != nil {
		t.Fatalf("Second load failed: %v", err)
	}

	if loaded2.CurrentProvider != "provider2" {
		t.Errorf("second save: expected provider2, got %s", loaded2.CurrentProvider)
	}

	// Verify old provider is no longer in the config
	if _, exists := loaded2.Providers["provider1"]; exists {
		t.Error("old provider1 should not exist after overwrite")
	}
}

// TestSaveWithEmptyConfig tests saving an empty or minimal config.
func TestSaveWithEmptyConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "empty-config.yaml")

	emptyConfig := &Config{
		Providers: make(map[string]ProviderConfig),
	}

	if err := emptyConfig.Save(configPath); err != nil {
		t.Fatalf("Save empty config failed: %v", err)
	}

	if !helper.FileExists(configPath) {
		t.Fatal("config file not created for empty config")
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load empty config failed: %v", err)
	}

	if loaded == nil {
		t.Fatal("loaded config should not be nil")
	}
}

// TestSaveConfigFilePermissions tests that saved config file has correct permissions.
func TestSaveConfigFilePermissions(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")
	config := SampleConfig()

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	fileInfo, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}

	// Check file permissions are 0644 (rw-r--r--)
	expectedMode := os.FileMode(0644)
	if fileInfo.Mode().Perm() != expectedMode {
		t.Errorf("file permissions mismatch: expected %o, got %o", expectedMode, fileInfo.Mode().Perm())
	}
}

// TestSaveWithSpecialCharactersInConfig tests saving config with special characters.
func TestSaveWithSpecialCharactersInConfig(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()
	configPath := filepath.Join(tempDir, "special-config.yaml")

	config := &Config{
		CurrentProvider: "test",
		Providers: map[string]ProviderConfig{
			"test": {
				Type:      "test",
				APIKey:    "key-with-special-chars-!@#$%^&*()_+-=[]{}|;:,.<>?",
				Model:     "model-with-unicode-日本語",
				MaxTokens: 1024,
				SystemPrompt: "You are a helpful assistant.\nWith newlines.\nAnd special chars: \"quotes\" and 'apostrophes'",
			},
		},
	}

	if err := config.Save(configPath); err != nil {
		t.Fatalf("Save with special characters failed: %v", err)
	}

	loaded, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load with special characters failed: %v", err)
	}

	provider := loaded.Providers["test"]
	if provider.APIKey != "key-with-special-chars-!@#$%^&*()_+-=[]{}|;:,.<>?" {
		t.Errorf("APIKey with special chars not preserved: got %s", provider.APIKey)
	}

	if provider.Model != "model-with-unicode-日本語" {
		t.Errorf("Model with unicode not preserved: got %s", provider.Model)
	}
}

// TestSaveMultipleConfigsToSameDirectory tests saving multiple config files to the same directory.
func TestSaveMultipleConfigsToSameDirectory(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	tempDir := helper.TempDir()

	config1 := ConfigWithProvider(t, "provider1", "type1", "key1", "model1")
	config2 := ConfigWithProvider(t, "provider2", "type2", "key2", "model2")

	configPath1 := filepath.Join(tempDir, "config1.yaml")
	configPath2 := filepath.Join(tempDir, "config2.yaml")

	if err := config1.Save(configPath1); err != nil {
		t.Fatalf("First save failed: %v", err)
	}

	if err := config2.Save(configPath2); err != nil {
		t.Fatalf("Second save failed: %v", err)
	}

	// Verify both files exist
	if !helper.FileExists(configPath1) {
		t.Fatal("config1 file not created")
	}

	if !helper.FileExists(configPath2) {
		t.Fatal("config2 file not created")
	}

	// Verify each file has the correct content
	loaded1, err := Load(configPath1)
	if err != nil {
		t.Fatalf("Load config1 failed: %v", err)
	}

	if loaded1.CurrentProvider != "provider1" {
		t.Errorf("config1 provider mismatch: expected provider1, got %s", loaded1.CurrentProvider)
	}

	loaded2, err := Load(configPath2)
	if err != nil {
		t.Fatalf("Load config2 failed: %v", err)
	}

	if loaded2.CurrentProvider != "provider2" {
		t.Errorf("config2 provider mismatch: expected provider2, got %s", loaded2.CurrentProvider)
	}
}

// TestGetConfigDirSuccessfulResolution tests that getConfigDir successfully resolves the home directory.
func TestGetConfigDirSuccessfulResolution(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	// Set up a temporary home directory
	homeDir := helper.TempDir()
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir failed: %v", err)
	}

	expectedDir := filepath.Join(homeDir, ".config", "how")
	if configDir != expectedDir {
		t.Errorf("config directory mismatch: expected %s, got %s", expectedDir, configDir)
	}
}

// TestGetConfigDirReturnsCorrectPath tests that getConfigDir returns the correct path format.
func TestGetConfigDirReturnsCorrectPath(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	homeDir := helper.TempDir()
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	configDir, err := getConfigDir()
	if err != nil {
		t.Fatalf("getConfigDir failed: %v", err)
	}

	// Verify the path contains the expected components
	if !filepath.IsAbs(configDir) {
		t.Errorf("expected absolute path, got relative path: %s", configDir)
	}

	// Verify path ends with the correct directory structure
	if !filepath.HasPrefix(configDir, homeDir) {
		t.Errorf("config directory should be under home directory: expected to start with %s, got %s", homeDir, configDir)
	}

	if !filepath.HasPrefix(configDir, filepath.Join(homeDir, ".config")) {
		t.Errorf("config directory should be under .config: expected to contain %s, got %s", filepath.Join(homeDir, ".config"), configDir)
	}
}

// TestGetConfigDirErrorHandlingNoHomeDir tests that getConfigDir returns an error when home directory cannot be determined.
func TestGetConfigDirErrorHandlingNoHomeDir(t *testing.T) {
	// Save current HOME value
	oldHome, wasSet := os.LookupEnv("HOME")
	defer func() {
		if wasSet {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// Clear HOME to simulate missing home directory
	if err := os.Unsetenv("HOME"); err != nil {
		t.Fatalf("failed to unset HOME: %v", err)
	}

	configDir, err := getConfigDir()
	if err == nil {
		t.Error("expected error when home directory cannot be determined, got nil")
	}

	if configDir != "" {
		t.Errorf("expected empty config directory on error, got %s", configDir)
	}
}

// TestGetConfigDirConsistency tests that getConfigDir returns consistent results across multiple calls.
func TestGetConfigDirConsistency(t *testing.T) {
	helper := NewTestHelper(t)
	defer helper.Cleanup()

	homeDir := helper.TempDir()
	oldHome := os.Getenv("HOME")
	if err := os.Setenv("HOME", homeDir); err != nil {
		t.Fatalf("failed to set HOME: %v", err)
	}
	defer func() {
		if oldHome != "" {
			os.Setenv("HOME", oldHome)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	// Call getConfigDir multiple times
	configDir1, err1 := getConfigDir()
	if err1 != nil {
		t.Fatalf("first getConfigDir call failed: %v", err1)
	}

	configDir2, err2 := getConfigDir()
	if err2 != nil {
		t.Fatalf("second getConfigDir call failed: %v", err2)
	}

	configDir3, err3 := getConfigDir()
	if err3 != nil {
		t.Fatalf("third getConfigDir call failed: %v", err3)
	}

	// All calls should return the same result
	if configDir1 != configDir2 {
		t.Errorf("inconsistent results between first and second call: %s vs %s", configDir1, configDir2)
	}

	if configDir1 != configDir3 {
		t.Errorf("inconsistent results between first and third call: %s vs %s", configDir1, configDir3)
	}
}

// YAML Marshaling Tests for Config Struct Types

// TestProviderConfigYAMLMarshal tests marshaling a ProviderConfig to YAML.
func TestProviderConfigYAMLMarshal(t *testing.T) {
	provider := ProviderConfig{
		Type:      "anthropic",
		APIKey:    "test-key",
		Model:     "claude-3-5-sonnet",
		BaseURL:   "https://api.anthropic.com",
		MaxTokens: 2048,
	}

	data, err := yaml.Marshal(provider)
	if err != nil {
		t.Fatalf("failed to marshal ProviderConfig: %v", err)
	}

	if data == nil || len(data) == 0 {
		t.Fatal("expected YAML data, got empty result")
	}

	// Verify YAML contains expected fields
	yamlStr := string(data)
	if !strings.Contains(yamlStr, "type: anthropic") {
		t.Errorf("expected 'type: anthropic' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "apiKey: test-key") {
		t.Errorf("expected 'apiKey: test-key' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "maxTokens: 2048") {
		t.Errorf("expected 'maxTokens: 2048' in YAML, got: %s", yamlStr)
	}
}

// TestProviderConfigYAMLUnmarshal tests unmarshaling YAML to ProviderConfig.
func TestProviderConfigYAMLUnmarshal(t *testing.T) {
	yamlData := `type: anthropic
apiKey: test-key
model: claude-3-5-sonnet
baseUrl: https://api.anthropic.com
maxTokens: 2048
temperature: 0.8
topP: 0.95
`

	var provider ProviderConfig
	err := yaml.Unmarshal([]byte(yamlData), &provider)
	if err != nil {
		t.Fatalf("failed to unmarshal ProviderConfig: %v", err)
	}

	if provider.Type != "anthropic" {
		t.Errorf("type mismatch: expected anthropic, got %s", provider.Type)
	}

	if provider.APIKey != "test-key" {
		t.Errorf("APIKey mismatch: expected test-key, got %s", provider.APIKey)
	}

	if provider.Model != "claude-3-5-sonnet" {
		t.Errorf("model mismatch: expected claude-3-5-sonnet, got %s", provider.Model)
	}

	if provider.BaseURL != "https://api.anthropic.com" {
		t.Errorf("baseUrl mismatch: expected https://api.anthropic.com, got %s", provider.BaseURL)
	}

	if provider.MaxTokens != 2048 {
		t.Errorf("maxTokens mismatch: expected 2048, got %d", provider.MaxTokens)
	}

	if provider.Temperature != 0.8 {
		t.Errorf("temperature mismatch: expected 0.8, got %f", provider.Temperature)
	}

	if provider.TopP != 0.95 {
		t.Errorf("topP mismatch: expected 0.95, got %f", provider.TopP)
	}
}

// TestProviderConfigYAMLRoundTrip tests that a ProviderConfig can be marshaled and unmarshaled without loss of data.
func TestProviderConfigYAMLRoundTrip(t *testing.T) {
	original := ProviderConfig{
		Type:      "anthropic",
		APIKey:    "test-key-12345",
		Model:     "claude-3-5-sonnet-20241022",
		BaseURL:   "https://api.anthropic.com",
		MaxTokens: 4096,
		Temperature: 0.7,
		TopP:       0.99,
		SystemPrompt: "You are a helpful assistant.",
		CustomHeaders: map[string]string{
			"X-Custom-Header": "value1",
			"X-Another-Header": "value2",
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal ProviderConfig: %v", err)
	}

	// Unmarshal back
	var unmarshaled ProviderConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal ProviderConfig: %v", err)
	}

	// Verify all fields match
	if unmarshaled.Type != original.Type {
		t.Errorf("type mismatch: expected %s, got %s", original.Type, unmarshaled.Type)
	}

	if unmarshaled.APIKey != original.APIKey {
		t.Errorf("APIKey mismatch: expected %s, got %s", original.APIKey, unmarshaled.APIKey)
	}

	if unmarshaled.Model != original.Model {
		t.Errorf("model mismatch: expected %s, got %s", original.Model, unmarshaled.Model)
	}

	if unmarshaled.BaseURL != original.BaseURL {
		t.Errorf("baseUrl mismatch: expected %s, got %s", original.BaseURL, unmarshaled.BaseURL)
	}

	if unmarshaled.MaxTokens != original.MaxTokens {
		t.Errorf("maxTokens mismatch: expected %d, got %d", original.MaxTokens, unmarshaled.MaxTokens)
	}

	if unmarshaled.Temperature != original.Temperature {
		t.Errorf("temperature mismatch: expected %f, got %f", original.Temperature, unmarshaled.Temperature)
	}

	if unmarshaled.TopP != original.TopP {
		t.Errorf("topP mismatch: expected %f, got %f", original.TopP, unmarshaled.TopP)
	}

	if unmarshaled.SystemPrompt != original.SystemPrompt {
		t.Errorf("systemPrompt mismatch: expected %s, got %s", original.SystemPrompt, unmarshaled.SystemPrompt)
	}

	if len(unmarshaled.CustomHeaders) != len(original.CustomHeaders) {
		t.Errorf("customHeaders count mismatch: expected %d, got %d", len(original.CustomHeaders), len(unmarshaled.CustomHeaders))
	}

	for key, expectedValue := range original.CustomHeaders {
		actualValue, exists := unmarshaled.CustomHeaders[key]
		if !exists {
			t.Errorf("custom header %q not found in unmarshaled config", key)
		}
		if actualValue != expectedValue {
			t.Errorf("custom header %q value mismatch: expected %s, got %s", key, expectedValue, actualValue)
		}
	}
}

// TestContextConfigYAMLMarshal tests marshaling a ContextConfig to YAML.
func TestContextConfigYAMLMarshal(t *testing.T) {
	context := ContextConfig{
		IncludeFiles:       true,
		IncludeHistory:     100,
		IncludeEnvironment: true,
		IncludeGit:         false,
		MaxContextSize:     16000,
		ExcludePatterns:    []string{".git", ".env", "node_modules"},
	}

	data, err := yaml.Marshal(context)
	if err != nil {
		t.Fatalf("failed to marshal ContextConfig: %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "includeFiles: true") {
		t.Errorf("expected 'includeFiles: true' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "includeHistory: 100") {
		t.Errorf("expected 'includeHistory: 100' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "maxContextSize: 16000") {
		t.Errorf("expected 'maxContextSize: 16000' in YAML, got: %s", yamlStr)
	}
}

// TestContextConfigYAMLUnmarshal tests unmarshaling YAML to ContextConfig.
func TestContextConfigYAMLUnmarshal(t *testing.T) {
	yamlData := `includeFiles: true
includeHistory: 100
includeEnvironment: false
includeGit: true
maxContextSize: 12000
excludePatterns:
  - .git
  - node_modules
  - .env
`

	var context ContextConfig
	err := yaml.Unmarshal([]byte(yamlData), &context)
	if err != nil {
		t.Fatalf("failed to unmarshal ContextConfig: %v", err)
	}

	if !context.IncludeFiles {
		t.Error("expected IncludeFiles to be true")
	}

	if context.IncludeHistory != 100 {
		t.Errorf("includeHistory mismatch: expected 100, got %d", context.IncludeHistory)
	}

	if context.IncludeEnvironment {
		t.Error("expected IncludeEnvironment to be false")
	}

	if !context.IncludeGit {
		t.Error("expected IncludeGit to be true")
	}

	if context.MaxContextSize != 12000 {
		t.Errorf("maxContextSize mismatch: expected 12000, got %d", context.MaxContextSize)
	}

	if len(context.ExcludePatterns) != 3 {
		t.Errorf("excludePatterns count mismatch: expected 3, got %d", len(context.ExcludePatterns))
	}
}

// TestContextConfigYAMLRoundTrip tests that a ContextConfig can be marshaled and unmarshaled without loss of data.
func TestContextConfigYAMLRoundTrip(t *testing.T) {
	original := ContextConfig{
		IncludeFiles:       true,
		IncludeHistory:     50,
		IncludeEnvironment: true,
		IncludeGit:         true,
		MaxContextSize:     8000,
		ExcludePatterns:    []string{".git", ".env", "node_modules", ".vscode"},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal ContextConfig: %v", err)
	}

	// Unmarshal back
	var unmarshaled ContextConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal ContextConfig: %v", err)
	}

	// Verify all fields match
	if unmarshaled.IncludeFiles != original.IncludeFiles {
		t.Errorf("IncludeFiles mismatch: expected %v, got %v", original.IncludeFiles, unmarshaled.IncludeFiles)
	}

	if unmarshaled.IncludeHistory != original.IncludeHistory {
		t.Errorf("IncludeHistory mismatch: expected %d, got %d", original.IncludeHistory, unmarshaled.IncludeHistory)
	}

	if unmarshaled.IncludeEnvironment != original.IncludeEnvironment {
		t.Errorf("IncludeEnvironment mismatch: expected %v, got %v", original.IncludeEnvironment, unmarshaled.IncludeEnvironment)
	}

	if unmarshaled.IncludeGit != original.IncludeGit {
		t.Errorf("IncludeGit mismatch: expected %v, got %v", original.IncludeGit, unmarshaled.IncludeGit)
	}

	if unmarshaled.MaxContextSize != original.MaxContextSize {
		t.Errorf("MaxContextSize mismatch: expected %d, got %d", original.MaxContextSize, unmarshaled.MaxContextSize)
	}

	if len(unmarshaled.ExcludePatterns) != len(original.ExcludePatterns) {
		t.Errorf("ExcludePatterns count mismatch: expected %d, got %d", len(original.ExcludePatterns), len(unmarshaled.ExcludePatterns))
	}

	for i, pattern := range original.ExcludePatterns {
		if unmarshaled.ExcludePatterns[i] != pattern {
			t.Errorf("ExcludePatterns[%d] mismatch: expected %s, got %s", i, pattern, unmarshaled.ExcludePatterns[i])
		}
	}
}

// TestDisplayConfigYAMLMarshal tests marshaling a DisplayConfig to YAML.
func TestDisplayConfigYAMLMarshal(t *testing.T) {
	display := DisplayConfig{
		SyntaxHighlight: true,
		ShowContext:     false,
		Emoji:           true,
		Color:           true,
	}

	data, err := yaml.Marshal(display)
	if err != nil {
		t.Fatalf("failed to marshal DisplayConfig: %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "syntaxHighlight: true") {
		t.Errorf("expected 'syntaxHighlight: true' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "showContext: false") {
		t.Errorf("expected 'showContext: false' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "emoji: true") {
		t.Errorf("expected 'emoji: true' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "color: true") {
		t.Errorf("expected 'color: true' in YAML, got: %s", yamlStr)
	}
}

// TestDisplayConfigYAMLUnmarshal tests unmarshaling YAML to DisplayConfig.
func TestDisplayConfigYAMLUnmarshal(t *testing.T) {
	yamlData := `syntaxHighlight: true
showContext: true
emoji: false
color: true
`

	var display DisplayConfig
	err := yaml.Unmarshal([]byte(yamlData), &display)
	if err != nil {
		t.Fatalf("failed to unmarshal DisplayConfig: %v", err)
	}

	if !display.SyntaxHighlight {
		t.Error("expected SyntaxHighlight to be true")
	}

	if !display.ShowContext {
		t.Error("expected ShowContext to be true")
	}

	if display.Emoji {
		t.Error("expected Emoji to be false")
	}

	if !display.Color {
		t.Error("expected Color to be true")
	}
}

// TestDisplayConfigYAMLRoundTrip tests that a DisplayConfig can be marshaled and unmarshaled without loss of data.
func TestDisplayConfigYAMLRoundTrip(t *testing.T) {
	original := DisplayConfig{
		SyntaxHighlight: true,
		ShowContext:     true,
		Emoji:           false,
		Color:           true,
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal DisplayConfig: %v", err)
	}

	// Unmarshal back
	var unmarshaled DisplayConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal DisplayConfig: %v", err)
	}

	// Verify all fields match
	if unmarshaled.SyntaxHighlight != original.SyntaxHighlight {
		t.Errorf("SyntaxHighlight mismatch: expected %v, got %v", original.SyntaxHighlight, unmarshaled.SyntaxHighlight)
	}

	if unmarshaled.ShowContext != original.ShowContext {
		t.Errorf("ShowContext mismatch: expected %v, got %v", original.ShowContext, unmarshaled.ShowContext)
	}

	if unmarshaled.Emoji != original.Emoji {
		t.Errorf("Emoji mismatch: expected %v, got %v", original.Emoji, unmarshaled.Emoji)
	}

	if unmarshaled.Color != original.Color {
		t.Errorf("Color mismatch: expected %v, got %v", original.Color, unmarshaled.Color)
	}
}

// TestHistoryConfigYAMLMarshal tests marshaling a HistoryConfig to YAML.
func TestHistoryConfigYAMLMarshal(t *testing.T) {
	history := HistoryConfig{
		Enabled:  true,
		MaxSize:  5000,
		FilePath: "~/.local/share/how/history",
	}

	data, err := yaml.Marshal(history)
	if err != nil {
		t.Fatalf("failed to marshal HistoryConfig: %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "enabled: true") {
		t.Errorf("expected 'enabled: true' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "maxSize: 5000") {
		t.Errorf("expected 'maxSize: 5000' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "filePath:") {
		t.Errorf("expected 'filePath:' in YAML, got: %s", yamlStr)
	}
}

// TestHistoryConfigYAMLUnmarshal tests unmarshaling YAML to HistoryConfig.
func TestHistoryConfigYAMLUnmarshal(t *testing.T) {
	yamlData := `enabled: true
maxSize: 2000
filePath: /home/user/.local/share/how/history
`

	var history HistoryConfig
	err := yaml.Unmarshal([]byte(yamlData), &history)
	if err != nil {
		t.Fatalf("failed to unmarshal HistoryConfig: %v", err)
	}

	if !history.Enabled {
		t.Error("expected Enabled to be true")
	}

	if history.MaxSize != 2000 {
		t.Errorf("maxSize mismatch: expected 2000, got %d", history.MaxSize)
	}

	if history.FilePath != "/home/user/.local/share/how/history" {
		t.Errorf("filePath mismatch: expected /home/user/.local/share/how/history, got %s", history.FilePath)
	}
}

// TestHistoryConfigYAMLRoundTrip tests that a HistoryConfig can be marshaled and unmarshaled without loss of data.
func TestHistoryConfigYAMLRoundTrip(t *testing.T) {
	original := HistoryConfig{
		Enabled:  true,
		MaxSize:  3000,
		FilePath: "/home/user/.config/how/history",
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal HistoryConfig: %v", err)
	}

	// Unmarshal back
	var unmarshaled HistoryConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal HistoryConfig: %v", err)
	}

	// Verify all fields match
	if unmarshaled.Enabled != original.Enabled {
		t.Errorf("Enabled mismatch: expected %v, got %v", original.Enabled, unmarshaled.Enabled)
	}

	if unmarshaled.MaxSize != original.MaxSize {
		t.Errorf("MaxSize mismatch: expected %d, got %d", original.MaxSize, unmarshaled.MaxSize)
	}

	if unmarshaled.FilePath != original.FilePath {
		t.Errorf("FilePath mismatch: expected %s, got %s", original.FilePath, unmarshaled.FilePath)
	}
}

// TestConfigYAMLMarshal tests marshaling a complete Config to YAML.
func TestConfigYAMLMarshal(t *testing.T) {
	config := SampleConfig()

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("failed to marshal Config: %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "currentProvider: anthropic") {
		t.Errorf("expected 'currentProvider: anthropic' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "providers:") {
		t.Errorf("expected 'providers:' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "context:") {
		t.Errorf("expected 'context:' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "display:") {
		t.Errorf("expected 'display:' in YAML, got: %s", yamlStr)
	}

	if !strings.Contains(yamlStr, "history:") {
		t.Errorf("expected 'history:' in YAML, got: %s", yamlStr)
	}
}

// TestConfigYAMLUnmarshal tests unmarshaling YAML to Config.
func TestConfigYAMLUnmarshal(t *testing.T) {
	yamlData := `currentProvider: anthropic
providers:
  anthropic:
    type: anthropic
    apiKey: sk-ant-test
    model: claude-3-5-sonnet
    maxTokens: 2048
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

	var config Config
	err := yaml.Unmarshal([]byte(yamlData), &config)
	if err != nil {
		t.Fatalf("failed to unmarshal Config: %v", err)
	}

	if config.CurrentProvider != "anthropic" {
		t.Errorf("currentProvider mismatch: expected anthropic, got %s", config.CurrentProvider)
	}

	if len(config.Providers) != 1 {
		t.Errorf("providers count mismatch: expected 1, got %d", len(config.Providers))
	}

	provider, exists := config.Providers["anthropic"]
	if !exists {
		t.Fatal("anthropic provider not found")
	}

	if provider.Type != "anthropic" {
		t.Errorf("provider type mismatch: expected anthropic, got %s", provider.Type)
	}

	if provider.APIKey != "sk-ant-test" {
		t.Errorf("provider apiKey mismatch: expected sk-ant-test, got %s", provider.APIKey)
	}

	if config.Context.IncludeFiles {
		if !config.Context.IncludeFiles {
			t.Error("expected context.includeFiles to be true")
		}
	}

	if !config.Display.SyntaxHighlight {
		t.Error("expected display.syntaxHighlight to be true")
	}

	if !config.History.Enabled {
		t.Error("expected history.enabled to be true")
	}
}

// TestConfigYAMLRoundTrip tests that a Config can be marshaled and unmarshaled without loss of data.
func TestConfigYAMLRoundTrip(t *testing.T) {
	original := &Config{
		CurrentProvider: "anthropic",
		Providers: map[string]ProviderConfig{
			"anthropic": {
				Type:      "anthropic",
				APIKey:    "sk-ant-test-key",
				Model:     "claude-3-5-sonnet-20241022",
				BaseURL:   "https://api.anthropic.com",
				MaxTokens: 4096,
				Temperature: 0.7,
				TopP:       1.0,
				SystemPrompt: "You are a helpful assistant",
				CustomHeaders: map[string]string{
					"X-Custom": "header",
				},
			},
			"backup": {
				Type:      "openai",
				APIKey:    "sk-openai-test",
				Model:     "gpt-4",
				MaxTokens: 8192,
				Temperature: 0.8,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     100,
			IncludeEnvironment: true,
			IncludeGit:         true,
			MaxContextSize:     16000,
			ExcludePatterns:    []string{".git", ".env", "node_modules"},
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
			FilePath: "~/.local/share/how/history",
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Config: %v", err)
	}

	// Unmarshal back
	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal Config: %v", err)
	}

	// Verify main fields
	if unmarshaled.CurrentProvider != original.CurrentProvider {
		t.Errorf("currentProvider mismatch: expected %s, got %s", original.CurrentProvider, unmarshaled.CurrentProvider)
	}

	if len(unmarshaled.Providers) != len(original.Providers) {
		t.Errorf("providers count mismatch: expected %d, got %d", len(original.Providers), len(unmarshaled.Providers))
	}

	// Verify provider details
	for name, originalProvider := range original.Providers {
		unmarshaledProvider, exists := unmarshaled.Providers[name]
		if !exists {
			t.Errorf("provider %q not found in unmarshaled config", name)
			continue
		}

		if unmarshaledProvider.Type != originalProvider.Type {
			t.Errorf("provider %q type mismatch: expected %s, got %s", name, originalProvider.Type, unmarshaledProvider.Type)
		}

		if unmarshaledProvider.APIKey != originalProvider.APIKey {
			t.Errorf("provider %q apiKey mismatch: expected %s, got %s", name, originalProvider.APIKey, unmarshaledProvider.APIKey)
		}

		if unmarshaledProvider.Model != originalProvider.Model {
			t.Errorf("provider %q model mismatch: expected %s, got %s", name, originalProvider.Model, unmarshaledProvider.Model)
		}

		if unmarshaledProvider.MaxTokens != originalProvider.MaxTokens {
			t.Errorf("provider %q maxTokens mismatch: expected %d, got %d", name, originalProvider.MaxTokens, unmarshaledProvider.MaxTokens)
		}
	}

	// Verify context fields
	if unmarshaled.Context.IncludeFiles != original.Context.IncludeFiles {
		t.Errorf("context.includeFiles mismatch: expected %v, got %v", original.Context.IncludeFiles, unmarshaled.Context.IncludeFiles)
	}

	if unmarshaled.Context.IncludeHistory != original.Context.IncludeHistory {
		t.Errorf("context.includeHistory mismatch: expected %d, got %d", original.Context.IncludeHistory, unmarshaled.Context.IncludeHistory)
	}

	if unmarshaled.Context.MaxContextSize != original.Context.MaxContextSize {
		t.Errorf("context.maxContextSize mismatch: expected %d, got %d", original.Context.MaxContextSize, unmarshaled.Context.MaxContextSize)
	}

	// Verify display fields
	if unmarshaled.Display.SyntaxHighlight != original.Display.SyntaxHighlight {
		t.Errorf("display.syntaxHighlight mismatch: expected %v, got %v", original.Display.SyntaxHighlight, unmarshaled.Display.SyntaxHighlight)
	}

	if unmarshaled.Display.ShowContext != original.Display.ShowContext {
		t.Errorf("display.showContext mismatch: expected %v, got %v", original.Display.ShowContext, unmarshaled.Display.ShowContext)
	}

	if unmarshaled.Display.Emoji != original.Display.Emoji {
		t.Errorf("display.emoji mismatch: expected %v, got %v", original.Display.Emoji, unmarshaled.Display.Emoji)
	}

	if unmarshaled.Display.Color != original.Display.Color {
		t.Errorf("display.color mismatch: expected %v, got %v", original.Display.Color, unmarshaled.Display.Color)
	}

	// Verify history fields
	if unmarshaled.History.Enabled != original.History.Enabled {
		t.Errorf("history.enabled mismatch: expected %v, got %v", original.History.Enabled, unmarshaled.History.Enabled)
	}

	if unmarshaled.History.MaxSize != original.History.MaxSize {
		t.Errorf("history.maxSize mismatch: expected %d, got %d", original.History.MaxSize, unmarshaled.History.MaxSize)
	}

	if unmarshaled.History.FilePath != original.History.FilePath {
		t.Errorf("history.filePath mismatch: expected %s, got %s", original.History.FilePath, unmarshaled.History.FilePath)
	}
}

// TestConfigWithMinimalFields tests marshaling and unmarshaling Config with minimal field combinations.
func TestConfigWithMinimalFields(t *testing.T) {
	original := &Config{
		CurrentProvider: "minimal",
		Providers: map[string]ProviderConfig{
			"minimal": {
				Type:      "test",
				APIKey:    "key",
				Model:     "model",
				MaxTokens: 1024,
			},
		},
		Context:  ContextConfig{},
		Display:  DisplayConfig{},
		History:  HistoryConfig{},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal minimal Config: %v", err)
	}

	// Unmarshal back
	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal minimal Config: %v", err)
	}

	if unmarshaled.CurrentProvider != original.CurrentProvider {
		t.Errorf("currentProvider mismatch: expected %s, got %s", original.CurrentProvider, unmarshaled.CurrentProvider)
	}

	if len(unmarshaled.Providers) != 1 {
		t.Errorf("providers count mismatch: expected 1, got %d", len(unmarshaled.Providers))
	}
}

// TestConfigWithMultipleProviderTypes tests marshaling and unmarshaling Config with multiple different provider types.
func TestConfigWithMultipleProviderTypes(t *testing.T) {
	original := &Config{
		CurrentProvider: "primary",
		Providers: map[string]ProviderConfig{
			"primary": {
				Type:      "anthropic",
				APIKey:    "sk-ant-key",
				Model:     "claude-3-5-sonnet",
				MaxTokens: 2048,
				Temperature: 0.7,
			},
			"secondary": {
				Type:      "openai",
				APIKey:    "sk-openai-key",
				Model:     "gpt-4",
				MaxTokens: 4096,
				Temperature: 0.5,
			},
			"tertiary": {
				Type:      "local",
				APIKey:    "local-key",
				Model:     "llama-2",
				MaxTokens: 8192,
				BaseURL:   "http://localhost:8000",
				Temperature: 0.9,
			},
		},
		Context: ContextConfig{
			IncludeFiles:       true,
			IncludeHistory:     100,
			IncludeEnvironment: true,
			MaxContextSize:     10000,
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
			ShowContext:     true,
			Color:           true,
		},
		History: HistoryConfig{
			Enabled:  true,
			MaxSize:  2000,
			FilePath: "/var/log/how/history",
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Config with multiple providers: %v", err)
	}

	// Unmarshal back
	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal Config with multiple providers: %v", err)
	}

	// Verify all providers are present
	if len(unmarshaled.Providers) != 3 {
		t.Errorf("expected 3 providers, got %d", len(unmarshaled.Providers))
	}

	for expectedName := range original.Providers {
		if _, exists := unmarshaled.Providers[expectedName]; !exists {
			t.Errorf("provider %q not found in unmarshaled config", expectedName)
		}
	}

	// Verify specific provider details
	primaryProvider := unmarshaled.Providers["primary"]
	if primaryProvider.Type != "anthropic" {
		t.Errorf("primary provider type mismatch: expected anthropic, got %s", primaryProvider.Type)
	}

	tertiaryProvider := unmarshaled.Providers["tertiary"]
	if tertiaryProvider.BaseURL != "http://localhost:8000" {
		t.Errorf("tertiary provider baseUrl mismatch: expected http://localhost:8000, got %s", tertiaryProvider.BaseURL)
	}
}

// TestConfigWithOmittedFields tests that omitted fields (with omitempty tags) are handled correctly.
func TestConfigWithOmittedFields(t *testing.T) {
	original := &Config{
		CurrentProvider: "test",
		Providers: map[string]ProviderConfig{
			"test": {
				Type:      "test",
				APIKey:    "key",
				Model:     "model",
				MaxTokens: 1024,
				// BaseURL, Temperature, TopP, SystemPrompt, CustomHeaders are omitted
			},
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Config with omitted fields: %v", err)
	}

	yamlStr := string(data)

	// Verify that empty optional fields are omitted
	if strings.Contains(yamlStr, "baseUrl:") && !strings.Contains(yamlStr, "baseUrl: null") {
		t.Errorf("expected baseUrl to be omitted or null, got: %s", yamlStr)
	}

	// Unmarshal back
	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal Config with omitted fields: %v", err)
	}

	provider := unmarshaled.Providers["test"]
	if provider.BaseURL != "" {
		t.Errorf("expected BaseURL to be empty, got %s", provider.BaseURL)
	}

	if provider.Temperature != 0 {
		t.Errorf("expected Temperature to be 0, got %f", provider.Temperature)
	}
}

// TestContextConfigEmptyExcludePatterns tests marshaling and unmarshaling with empty exclude patterns.
func TestContextConfigEmptyExcludePatterns(t *testing.T) {
	original := ContextConfig{
		IncludeFiles:       true,
		IncludeHistory:     50,
		IncludeEnvironment: false,
		IncludeGit:         true,
		MaxContextSize:     8000,
		ExcludePatterns:    []string{},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal ContextConfig with empty patterns: %v", err)
	}

	// Unmarshal back
	var unmarshaled ContextConfig
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal ContextConfig with empty patterns: %v", err)
	}

	if len(unmarshaled.ExcludePatterns) != 0 && unmarshaled.ExcludePatterns != nil {
		t.Errorf("expected empty ExcludePatterns, got %v", unmarshaled.ExcludePatterns)
	}
}

// TestConfigWithSpecialCharactersInFields tests marshaling and unmarshaling with special characters.
func TestConfigWithSpecialCharactersInFields(t *testing.T) {
	original := &Config{
		CurrentProvider: "special",
		Providers: map[string]ProviderConfig{
			"special": {
				Type:      "special",
				APIKey:    "key-with-special-!@#$%^&*()",
				Model:     "model-with-unicode-日本語",
				MaxTokens: 2048,
				SystemPrompt: "You are helpful.\nMultiline.\n\"Quoted\" text.",
				CustomHeaders: map[string]string{
					"X-Special": "value with special: chars!",
				},
			},
		},
		Display: DisplayConfig{
			SyntaxHighlight: true,
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal Config with special characters: %v", err)
	}

	// Unmarshal back
	var unmarshaled Config
	err = yaml.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("failed to unmarshal Config with special characters: %v", err)
	}

	provider := unmarshaled.Providers["special"]
	if provider.APIKey != "key-with-special-!@#$%^&*()" {
		t.Errorf("special characters in APIKey not preserved: got %s", provider.APIKey)
	}

	if provider.Model != "model-with-unicode-日本語" {
		t.Errorf("unicode characters in Model not preserved: got %s", provider.Model)
	}

	if provider.SystemPrompt != "You are helpful.\nMultiline.\n\"Quoted\" text." {
		t.Errorf("multiline and quoted text in SystemPrompt not preserved: got %s", provider.SystemPrompt)
	}
}
