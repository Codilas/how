package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Codilas/how/internal/cli/testfixtures"
	"github.com/spf13/cobra"
)

// TestHistoryCommandFullFlow_WithExistingHistory tests the full Cobra command flow for 'how history'
// with existing conversation history file.
func TestHistoryCommandFullFlow_WithExistingHistory(t *testing.T) {
	// Setup: Create a test environment with conversation history
	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(
			[]string{"How do I list files?", "What about hidden files?"},
			[]string{"Use the ls command", "Use the -a flag with ls"},
		).
		Build()

	// Create a custom command that uses our test setup
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	// Execute the command
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()

	// Verify
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "📜 Recent conversations:") {
		t.Errorf("expected header in output, got: %q", output)
	}
	if !strings.Contains(output, "How do I list files?") {
		t.Errorf("expected conversation content in output, got: %q", output)
	}
	if !strings.Contains(output, "Use the ls command") {
		t.Errorf("expected response content in output, got: %q", output)
	}

	if len(stderr.String()) > 0 {
		t.Errorf("expected no stderr output, got: %q", stderr.String())
	}
}

// TestHistoryCommandFullFlow_NoHistory tests the full Cobra command flow when no history exists.
func TestHistoryCommandFullFlow_NoHistory(t *testing.T) {
	// Setup: Create an empty test environment
	setup := testfixtures.EmptyHistoryTestSetup(t)

	// Create a custom command that uses our test setup
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	// Execute the command
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()

	// Verify
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "📜 No conversation history found.") {
		t.Errorf("expected no history message in output, got: %q", output)
	}

	if len(stderr.String()) > 0 {
		t.Errorf("expected no stderr output, got: %q", stderr.String())
	}
}

// TestHistoryCommandFullFlow_LargeConversationHistory tests with a larger history file.
func TestHistoryCommandFullFlow_LargeConversationHistory(t *testing.T) {
	// Create many conversation entries
	prompts := make([]string, 10)
	responses := make([]string, 10)
	for i := 0; i < 10; i++ {
		prompts[i] = "Prompt " + string(rune(48+i))
		responses[i] = "Response " + string(rune(48+i))
	}

	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(prompts, responses).
		Build()

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	output := stdout.String()
	// Verify all prompts and responses are present
	for i := 0; i < 10; i++ {
		if !strings.Contains(output, prompts[i]) {
			t.Errorf("expected prompt %q in output", prompts[i])
		}
		if !strings.Contains(output, responses[i]) {
			t.Errorf("expected response %q in output", responses[i])
		}
	}
}

// TestHistoryCommandFullFlow_MultilineContent tests history with multiline content.
func TestHistoryCommandFullFlow_MultilineContent(t *testing.T) {
	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(
			[]string{"How do I write a Go function?\nWith multiple lines?"},
			[]string{"Use this syntax:\nfunc name() {\n  // code\n}"},
		).
		Build()

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "How do I write a Go function?") {
		t.Errorf("expected multiline content preserved in output")
	}
	if !strings.Contains(output, "func name()") {
		t.Errorf("expected code block in output")
	}
}

// TestClearHistoryCommandFullFlow_WithExistingHistory tests 'how history clear' with existing history.
func TestClearHistoryCommandFullFlow_WithExistingHistory(t *testing.T) {
	// Setup: Create a test environment with conversation history
	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(
			[]string{"Test prompt"},
			[]string{"Test response"},
		).
		Build()

	// Verify file exists before clearing
	if !setup.ConversationHistoryExists() {
		t.Fatalf("history file should exist before clearing")
	}

	// Create a custom command that uses our test setup
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Test clear history command",
		Run: func(cmd *cobra.Command, args []string) {
			runClearHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	// Execute the command
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()

	// Verify
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	// Verify the file was deleted
	if setup.ConversationHistoryExists() {
		t.Errorf("history file should be deleted after clear command")
	}

	output := stdout.String()
	if !strings.Contains(output, "🗑️ Conversation history cleared.") {
		t.Errorf("expected success message in output, got: %q", output)
	}

	if len(stderr.String()) > 0 {
		t.Errorf("expected no stderr output, got: %q", stderr.String())
	}
}

// TestClearHistoryCommandFullFlow_NoHistory tests 'how history clear' when no history exists.
func TestClearHistoryCommandFullFlow_NoHistory(t *testing.T) {
	// Setup: Create an empty test environment
	setup := testfixtures.EmptyHistoryTestSetup(t)

	// Verify file doesn't exist
	if setup.ConversationHistoryExists() {
		t.Fatalf("history file should not exist in empty setup")
	}

	// Create a custom command that uses our test setup
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Test clear history command",
		Run: func(cmd *cobra.Command, args []string) {
			runClearHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	// Execute the command
	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()

	// Verify
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "🗑️ Conversation history cleared.") {
		t.Errorf("expected success message even when file doesn't exist, got: %q", output)
	}

	if len(stderr.String()) > 0 {
		t.Errorf("expected no stderr output, got: %q", stderr.String())
	}
}

// TestHistoryAndClearCommandSequence tests the sequence: show history, then clear it, then verify empty.
func TestHistoryAndClearCommandSequence(t *testing.T) {
	// Setup: Create a test environment with conversation history
	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(
			[]string{"Command 1", "Command 2"},
			[]string{"Response 1", "Response 2"},
		).
		Build()

	// Step 1: Show history
	historyCmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	var historyOut bytes.Buffer
	historyCmd.SetOut(&historyOut)
	historyCmd.SetErr(&bytes.Buffer{})

	err := historyCmd.Execute()
	if err != nil {
		t.Fatalf("history command execution failed: %v", err)
	}

	historyOutput := historyOut.String()
	if !strings.Contains(historyOutput, "Command 1") {
		t.Errorf("expected history to show content before clearing")
	}

	// Step 2: Clear history
	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Test clear history command",
		Run: func(cmd *cobra.Command, args []string) {
			runClearHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	var clearOut bytes.Buffer
	clearCmd.SetOut(&clearOut)
	clearCmd.SetErr(&bytes.Buffer{})

	err = clearCmd.Execute()
	if err != nil {
		t.Fatalf("clear history command execution failed: %v", err)
	}

	// Step 3: Show history again (should be empty)
	historyCmd2 := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	var historyOut2 bytes.Buffer
	historyCmd2.SetOut(&historyOut2)
	historyCmd2.SetErr(&bytes.Buffer{})

	err = historyCmd2.Execute()
	if err != nil {
		t.Fatalf("second history command execution failed: %v", err)
	}

	historyOutput2 := historyOut2.String()
	if !strings.Contains(historyOutput2, "📜 No conversation history found.") {
		t.Errorf("expected empty history message after clearing, got: %q", historyOutput2)
	}
}

// TestHistoryCommandStderr_FileReadError tests error handling when history file cannot be read.
func TestHistoryCommandStderr_FileReadError(t *testing.T) {
	// Create a directory where we expect a file (to cause read error)
	setup := testfixtures.EmptyHistoryTestSetup(t)
	historyPath := setup.GetConversationHistoryPath()

	// Create a directory at the history file path
	if err := os.MkdirAll(historyPath, 0755); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return historyPath
			})
		},
	}

	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	// Execute the command - should not panic
	cmd.Execute()

	// Verify error was reported to stderr
	stderrOutput := stderr.String()
	if !strings.Contains(stderrOutput, "Error reading history:") {
		t.Errorf("expected error message in stderr, got: %q", stderrOutput)
	}
}

// TestClearHistoryCommandStderr_PermissionError tests error handling for clear history permission issues.
func TestClearHistoryCommandStderr_PermissionError(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	// Setup: Create a test environment with conversation history
	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(
			[]string{"Test prompt"},
			[]string{"Test response"},
		).
		Build()

	// Make the config directory read-only to prevent deletion
	if err := os.Chmod(setup.ConfigDir, 0555); err != nil {
		t.Fatalf("failed to change directory permissions: %v", err)
	}
	defer os.Chmod(setup.ConfigDir, 0755)

	cmd := &cobra.Command{
		Use:   "clear",
		Short: "Test clear history command",
		Run: func(cmd *cobra.Command, args []string) {
			runClearHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	cmd.Execute()

	// Verify error was reported to stderr
	stderrOutput := stderr.String()
	if !strings.Contains(stderrOutput, "Error clearing history:") {
		t.Errorf("expected error message in stderr, got: %q", stderrOutput)
	}

	// Verify no success message
	if strings.Contains(stdout.String(), "🗑️ Conversation history cleared.") {
		t.Errorf("expected no success message when deletion fails")
	}
}

// TestHistoryCommandPathGetterFunction tests that custom path getter is used correctly.
func TestHistoryCommandPathGetterFunction(t *testing.T) {
	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(
			[]string{"Custom path test"},
			[]string{"It works!"},
		).
		Build()

	pathGetterCalled := false

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				pathGetterCalled = true
				return setup.GetConversationHistoryPath()
			})
		},
	}

	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	if !pathGetterCalled {
		t.Errorf("expected path getter function to be called")
	}

	if !strings.Contains(stdout.String(), "Custom path test") {
		t.Errorf("expected history content to be displayed")
	}
}

// TestHistoryCommandSubcommandIntegration tests 'history' as a Cobra subcommand.
func TestHistoryCommandSubcommandIntegration(t *testing.T) {
	// Create parent command with history subcommand
	rootCmd := &cobra.Command{
		Use: "how",
		Run: func(cmd *cobra.Command, args []string) {
			// no-op
		},
	}

	setup := testfixtures.NewHistoryTestBuilder(t).
		WithConversationHistory(
			[]string{"How do I use Go?"},
			[]string{"Go is a programming language"},
		).
		Build()

	// Create history subcommand
	historySubCmd := &cobra.Command{
		Use:   "history",
		Short: "Show conversation history",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	// Create clear subcommand
	clearSubCmd := &cobra.Command{
		Use:   "clear",
		Short: "Clear conversation history",
		Run: func(cmd *cobra.Command, args []string) {
			runClearHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return setup.GetConversationHistoryPath()
			})
		},
	}

	historySubCmd.AddCommand(clearSubCmd)
	rootCmd.AddCommand(historySubCmd)

	// Test: how history
	var out1 bytes.Buffer
	rootCmd.SetOut(&out1)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"history"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("'how history' command failed: %v", err)
	}

	if !strings.Contains(out1.String(), "How do I use Go?") {
		t.Errorf("expected history content in output")
	}

	// Test: how history clear
	var out2 bytes.Buffer
	rootCmd.SetOut(&out2)
	rootCmd.SetErr(&bytes.Buffer{})
	rootCmd.SetArgs([]string{"history", "clear"})

	err = rootCmd.Execute()
	if err != nil {
		t.Fatalf("'how history clear' command failed: %v", err)
	}

	if !strings.Contains(out2.String(), "🗑️ Conversation history cleared.") {
		t.Errorf("expected clear success message in output")
	}

	// Verify file was actually deleted
	if setup.ConversationHistoryExists() {
		t.Errorf("history file should be deleted after clear")
	}
}

// TestHistoryCommand_ActualFileIO tests with real file I/O using actual history file paths.
func TestHistoryCommand_ActualFileIO(t *testing.T) {
	// Create a temporary directory to simulate real user config structure
	tempDir, err := os.MkdirTemp("", "history-integration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the .config/how directory structure
	configDir := filepath.Join(tempDir, ".config", "how")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config dir: %v", err)
	}

	historyFile := filepath.Join(configDir, "history.txt")
	testContent := `PROMPT: How do I list files?
RESPONSE: Use the ls command

PROMPT: What about hidden files?
RESPONSE: Use the -a flag
`

	// Write the history file
	if err := os.WriteFile(historyFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("failed to write history file: %v", err)
	}

	// Create command that uses the actual file
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Test history command",
		Run: func(cmd *cobra.Command, args []string) {
			runHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return historyFile
			})
		},
	}

	var stdout, stderr bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	err = cmd.Execute()
	if err != nil {
		t.Fatalf("command execution failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "How do I list files?") {
		t.Errorf("expected prompt in output, got: %q", output)
	}
	if !strings.Contains(output, "Use the ls command") {
		t.Errorf("expected response in output, got: %q", output)
	}

	// Now test clearing
	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Test clear history command",
		Run: func(cmd *cobra.Command, args []string) {
			runClearHistoryWithWriter(cmd.OutOrStdout(), cmd.OutOrStderr(), func() string {
				return historyFile
			})
		},
	}

	var clearOut bytes.Buffer
	clearCmd.SetOut(&clearOut)
	clearCmd.SetErr(&bytes.Buffer{})

	err = clearCmd.Execute()
	if err != nil {
		t.Fatalf("clear command execution failed: %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(historyFile); !os.IsNotExist(err) {
		t.Errorf("history file should be deleted after clear")
	}
}
