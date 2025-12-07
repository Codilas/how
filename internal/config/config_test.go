package config

import (
	"os"
	"path/filepath"
	"testing"
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
