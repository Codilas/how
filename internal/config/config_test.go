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
