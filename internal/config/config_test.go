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
