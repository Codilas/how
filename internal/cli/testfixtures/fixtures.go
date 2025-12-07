// Package testfixtures provides test utilities and helpers for testing the CLI package.
package testfixtures

import (
	"os"
	"path/filepath"
	"testing"
)

// TempDir creates a temporary directory for testing and returns its path.
// The directory is automatically cleaned up when the test completes.
func TempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "how-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}

	// Register cleanup function
	t.Cleanup(func() {
		_ = os.RemoveAll(dir)
	})

	return dir
}

// TempConfigDir creates a temporary config directory structure for testing.
// It simulates the ~/.config/how directory structure used by the how CLI.
func TempConfigDir(t *testing.T) string {
	t.Helper()
	baseDir := TempDir(t)
	configDir := filepath.Join(baseDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create temp config directory: %v", err)
	}
	return baseDir
}

// TempHistoryDir creates a temporary directory for shell history files.
// Returns the base directory and the specific history directory path.
func TempHistoryDir(t *testing.T) (baseDir, historyDir string) {
	t.Helper()
	baseDir = TempDir(t)
	historyDir = filepath.Join(baseDir, ".local", "share", "fish")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		t.Fatalf("failed to create temp history directory: %v", err)
	}
	bashHistoryDir := filepath.Join(baseDir)
	zshHistoryDir := filepath.Join(baseDir)
	if err := os.MkdirAll(bashHistoryDir, 0755); err != nil {
		t.Fatalf("failed to create temp bash history directory: %v", err)
	}
	if err := os.MkdirAll(zshHistoryDir, 0755); err != nil {
		t.Fatalf("failed to create temp zsh history directory: %v", err)
	}
	return baseDir, historyDir
}

// WriteFile writes test data to a file in the specified directory.
// It ensures the parent directory exists.
func WriteFile(t *testing.T, dir, filename string, content string) string {
	t.Helper()
	filePath := filepath.Join(dir, filename)

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	return filePath
}

// ReadFile reads a test file and returns its content.
func ReadFile(t *testing.T, filePath string) string {
	t.Helper()
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}
	return string(content)
}

// FileExists checks if a file exists at the given path.
func FileExists(t *testing.T, filePath string) bool {
	t.Helper()
	_, err := os.Stat(filePath)
	return err == nil
}

// FileDoesNotExist checks if a file does not exist at the given path.
func FileDoesNotExist(t *testing.T, filePath string) bool {
	t.Helper()
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}
