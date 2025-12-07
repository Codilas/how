package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestHelper provides utility functions for config testing.
type TestHelper struct {
	t        *testing.T
	tempDirs []string
}

// NewTestHelper creates a new TestHelper instance for managing test cleanup.
func NewTestHelper(t *testing.T) *TestHelper {
	t.Helper()
	return &TestHelper{
		t:        t,
		tempDirs: []string{},
	}
}

// Cleanup removes all temporary directories and files created during testing.
func (h *TestHelper) Cleanup() {
	h.t.Helper()
	for _, dir := range h.tempDirs {
		if err := os.RemoveAll(dir); err != nil {
			h.t.Errorf("failed to remove temp directory %s: %v", dir, err)
		}
	}
}

// TempDir creates a temporary directory and registers it for cleanup.
// Returns the path to the temporary directory.
func (h *TestHelper) TempDir() string {
	h.t.Helper()

	tempDir, err := os.MkdirTemp("", "how-config-test-*")
	if err != nil {
		h.t.Fatalf("failed to create temp directory: %v", err)
	}

	h.tempDirs = append(h.tempDirs, tempDir)
	return tempDir
}

// WriteConfig writes a Config struct to a YAML file in the specified directory.
// Returns the full path to the created config file.
func (h *TestHelper) WriteConfig(config *Config, dir string) string {
	h.t.Helper()

	configPath := filepath.Join(dir, "config.yaml")
	if err := config.Save(configPath); err != nil {
		h.t.Fatalf("failed to save config: %v", err)
	}

	return configPath
}

// WriteYAML writes raw YAML content to a config file in the specified directory.
// Returns the full path to the created config file.
func (h *TestHelper) WriteYAML(yamlContent string, dir string) string {
	h.t.Helper()

	configPath := filepath.Join(dir, "config.yaml")
	WriteConfigFile(h.t, configPath, yamlContent)
	return configPath
}

// LoadConfig loads a Config from the specified file path.
// Returns the loaded Config or fails the test on error.
func (h *TestHelper) LoadConfig(configPath string) *Config {
	h.t.Helper()

	config, err := Load(configPath)
	if err != nil {
		h.t.Fatalf("failed to load config: %v", err)
	}

	return config
}

// AssertConfigEqual compares two Config structs for equality.
// Fails the test if they don't match.
func (h *TestHelper) AssertConfigEqual(expected, actual *Config) {
	h.t.Helper()

	if expected.CurrentProvider != actual.CurrentProvider {
		h.t.Errorf("CurrentProvider mismatch: expected %q, got %q", expected.CurrentProvider, actual.CurrentProvider)
	}

	if len(expected.Providers) != len(actual.Providers) {
		h.t.Errorf("number of providers mismatch: expected %d, got %d", len(expected.Providers), len(actual.Providers))
	}

	for name, expectedProvider := range expected.Providers {
		actualProvider, exists := actual.Providers[name]
		if !exists {
			h.t.Errorf("provider %q not found in actual config", name)
			continue
		}

		if expectedProvider.Type != actualProvider.Type {
			h.t.Errorf("provider %q type mismatch: expected %q, got %q", name, expectedProvider.Type, actualProvider.Type)
		}
		if expectedProvider.APIKey != actualProvider.APIKey {
			h.t.Errorf("provider %q APIKey mismatch: expected %q, got %q", name, expectedProvider.APIKey, actualProvider.APIKey)
		}
		if expectedProvider.Model != actualProvider.Model {
			h.t.Errorf("provider %q Model mismatch: expected %q, got %q", name, expectedProvider.Model, actualProvider.Model)
		}
		if expectedProvider.MaxTokens != actualProvider.MaxTokens {
			h.t.Errorf("provider %q MaxTokens mismatch: expected %d, got %d", name, expectedProvider.MaxTokens, actualProvider.MaxTokens)
		}
	}

	if expected.Context.IncludeFiles != actual.Context.IncludeFiles {
		h.t.Errorf("Context.IncludeFiles mismatch: expected %v, got %v", expected.Context.IncludeFiles, actual.Context.IncludeFiles)
	}
}

// CreateConfigWithEnv creates a config with environment variable overrides.
// Sets up environment variables prefixed with HOW_ and returns a cleanup function.
func (h *TestHelper) CreateConfigWithEnv(key, value string) func() {
	h.t.Helper()

	envKey := "HOW_" + key
	oldValue, wasSet := os.LookupEnv(envKey)

	if err := os.Setenv(envKey, value); err != nil {
		h.t.Fatalf("failed to set environment variable: %v", err)
	}

	return func() {
		if wasSet {
			if err := os.Setenv(envKey, oldValue); err != nil {
				h.t.Errorf("failed to restore environment variable: %v", err)
			}
		} else {
			if err := os.Unsetenv(envKey); err != nil {
				h.t.Errorf("failed to unset environment variable: %v", err)
			}
		}
	}
}

// FileExists checks if a file exists at the given path.
func (h *TestHelper) FileExists(path string) bool {
	h.t.Helper()
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads the contents of a file and returns it as a string.
func (h *TestHelper) ReadFile(path string) string {
	h.t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		h.t.Fatalf("failed to read file %s: %v", path, err)
	}

	return string(content)
}

// DirectoryExists checks if a directory exists at the given path.
func (h *TestHelper) DirectoryExists(path string) bool {
	h.t.Helper()

	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
