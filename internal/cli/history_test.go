package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestRunHistoryFileNotFound tests the runHistoryWithWriter function when the history file does not exist.
func TestRunHistoryFileNotFound(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock getHistoryFile that points to a non-existent file
	nonExistentFile := filepath.Join(tempDir, ".config", "how", "history.txt")

	// Capture stdout and stderr with buffers
	var stdout, stderr bytes.Buffer

	// Call runHistoryWithWriter
	runHistoryWithWriter(&stdout, &stderr, func() string {
		return nonExistentFile
	})

	// Verify output
	expectedMsg := "📜 No conversation history found."
	if !strings.Contains(stdout.String(), expectedMsg) {
		t.Fatalf("expected output to contain %q, got: %q", expectedMsg, stdout.String())
	}
}

// TestRunHistoryEmptyFile tests the runHistoryWithWriter function when the history file is empty.
func TestRunHistoryEmptyFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the config directory structure
	configDir := filepath.Join(tempDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Create an empty history file
	historyFile := filepath.Join(configDir, "history.txt")
	if err := os.WriteFile(historyFile, []byte(""), 0644); err != nil {
		t.Fatalf("failed to create history file: %v", err)
	}

	// Capture stdout and stderr with buffers
	var stdout, stderr bytes.Buffer

	// Call runHistoryWithWriter
	runHistoryWithWriter(&stdout, &stderr, func() string {
		return historyFile
	})

	// Verify output
	expectedMsg := "📜 No conversation history found."
	if !strings.Contains(stdout.String(), expectedMsg) {
		t.Fatalf("expected output to contain %q, got: %q", expectedMsg, stdout.String())
	}
}

// TestRunHistoryWithContent tests the runHistoryWithWriter function when the history file contains content.
func TestRunHistoryWithContent(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the config directory structure
	configDir := filepath.Join(tempDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Create a history file with content
	historyFile := filepath.Join(configDir, "history.txt")
	testContent := "User: How do I list files?\nAssistant: Use ls command\n"
	if err := os.WriteFile(historyFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create history file: %v", err)
	}

	// Capture stdout and stderr with buffers
	var stdout, stderr bytes.Buffer

	// Call runHistoryWithWriter
	runHistoryWithWriter(&stdout, &stderr, func() string {
		return historyFile
	})

	// Verify output contains the expected header and content
	if !strings.Contains(stdout.String(), "📜 Recent conversations:") {
		t.Fatalf("expected output to contain %q, got: %q", "📜 Recent conversations:", stdout.String())
	}
	if !strings.Contains(stdout.String(), testContent) {
		t.Fatalf("expected output to contain file content %q, got: %q", testContent, stdout.String())
	}
}

// TestRunHistoryOutputFormatting tests the output formatting with io.Writer.
func TestRunHistoryOutputFormatting(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the config directory structure
	configDir := filepath.Join(tempDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Create a history file with multiple lines
	historyFile := filepath.Join(configDir, "history.txt")
	testContent := "Conversation 1\nConversation 2\nConversation 3\n"
	if err := os.WriteFile(historyFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create history file: %v", err)
	}

	// Capture stdout and stderr using buffers
	var stdout, stderr bytes.Buffer

	// Call runHistoryWithWriter
	runHistoryWithWriter(&stdout, &stderr, func() string {
		return historyFile
	})

	// Verify that all content is present
	if !strings.Contains(stdout.String(), "Conversation 1") {
		t.Errorf("expected output to contain 'Conversation 1', got: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Conversation 2") {
		t.Errorf("expected output to contain 'Conversation 2', got: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Conversation 3") {
		t.Errorf("expected output to contain 'Conversation 3', got: %q", stdout.String())
	}

	// Verify no errors were written to stderr
	if len(stderr.String()) > 0 {
		t.Errorf("expected no error output, got: %q", stderr.String())
	}
}

// TestRunHistoryMultilineContent tests that multiline content is properly displayed.
func TestRunHistoryMultilineContent(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the config directory structure
	configDir := filepath.Join(tempDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Create a history file with multiline content
	historyFile := filepath.Join(configDir, "history.txt")
	testContent := `User: How do I list files?
Assistant: Use the ls command:
  ls [options] [path]

User: What about hidden files?
Assistant: Use the -a flag:
  ls -a
`
	if err := os.WriteFile(historyFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create history file: %v", err)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Call runHistoryWithWriter
	runHistoryWithWriter(&stdout, &stderr, func() string {
		return historyFile
	})

	// Verify multiline content is preserved
	if !strings.Contains(stdout.String(), "User: How do I list files?") {
		t.Fatalf("expected multiline content to be preserved")
	}
	if !strings.Contains(stdout.String(), "  ls [options] [path]") {
		t.Fatalf("expected indented content to be preserved")
	}
}

// TestRunClearHistoryFileExists tests the runClearHistoryWithWriter function when the history file exists.
func TestRunClearHistoryFileExists(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "clear-history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the config directory structure
	configDir := filepath.Join(tempDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Create a history file
	historyFile := filepath.Join(configDir, "history.txt")
	testContent := "User: How do I list files?\nAssistant: Use ls command\n"
	if err := os.WriteFile(historyFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to create history file: %v", err)
	}

	// Verify the file exists before deletion
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		t.Fatalf("history file should exist before deletion")
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Call runClearHistoryWithWriter
	runClearHistoryWithWriter(&stdout, &stderr, func() string {
		return historyFile
	})

	// Verify the file was deleted
	if _, err := os.Stat(historyFile); !os.IsNotExist(err) {
		t.Fatalf("history file should be deleted after runClearHistoryWithWriter")
	}

	// Verify output message
	expectedMsg := "🗑️ Conversation history cleared."
	if !strings.Contains(stdout.String(), expectedMsg) {
		t.Fatalf("expected output to contain %q, got: %q", expectedMsg, stdout.String())
	}

	// Verify no error output
	if len(stderr.String()) > 0 {
		t.Fatalf("expected no error output, got: %q", stderr.String())
	}
}

// TestRunClearHistoryFileNotFound tests the runClearHistoryWithWriter function when the history file does not exist.
func TestRunClearHistoryFileNotFound(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "clear-history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a mock getHistoryFile that points to a non-existent file
	nonExistentFile := filepath.Join(tempDir, ".config", "how", "history.txt")

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Call runClearHistoryWithWriter
	runClearHistoryWithWriter(&stdout, &stderr, func() string {
		return nonExistentFile
	})

	// Verify success message is printed (no error for non-existent file)
	expectedMsg := "🗑️ Conversation history cleared."
	if !strings.Contains(stdout.String(), expectedMsg) {
		t.Fatalf("expected output to contain %q, got: %q", expectedMsg, stdout.String())
	}

	// Verify no error output
	if len(stderr.String()) > 0 {
		t.Fatalf("expected no error output when file doesn't exist, got: %q", stderr.String())
	}
}

// TestRunClearHistoryDeletionError tests the runClearHistoryWithWriter function when file deletion fails.
func TestRunClearHistoryDeletionError(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "clear-history-test-*")
	if err != nil {
		t.Fatalf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the config directory structure
	configDir := filepath.Join(tempDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	// Create a history directory (not a file) to cause deletion error
	historyFile := filepath.Join(configDir, "history.txt")
	if err := os.Mkdir(historyFile, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// Capture stdout and stderr
	var stdout, stderr bytes.Buffer

	// Call runClearHistoryWithWriter
	runClearHistoryWithWriter(&stdout, &stderr, func() string {
		return historyFile
	})

	// Verify error message is written to stderr
	if !strings.Contains(stderr.String(), "Error clearing history:") {
		t.Fatalf("expected error output to contain 'Error clearing history:', got: %q", stderr.String())
	}

	// Verify success message is NOT printed
	if strings.Contains(stdout.String(), "🗑️ Conversation history cleared.") {
		t.Fatalf("expected no success message when deletion fails, got: %q", stdout.String())
	}
}
