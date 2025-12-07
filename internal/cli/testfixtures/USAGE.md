# Test Fixtures Usage Guide

This guide provides practical examples for using the test fixtures in the `how` CLI project.

## Table of Contents

1. [Testing History Commands](#testing-history-commands)
2. [Testing Shell History Parsing](#testing-shell-history-parsing)
3. [Testing Conversation History](#testing-conversation-history)
4. [Advanced Scenarios](#advanced-scenarios)

## Testing History Commands

### Example 1: Testing the History Command Display

```go
package cli

import (
	"testing"
	"bytes"
	"io"
	"os"

	"github.com/Codilas/how/internal/cli/testfixtures"
)

func TestHistoryCommandDisplay(t *testing.T) {
	// Set up test environment
	setup := testfixtures.NewHistoryTestSetup(t)

	// Create test conversation history
	prompts := []string{
		"How do I list files in Go?",
		"How do I handle errors?",
	}
	responses := []string{
		"Use filepath.Walk() or ioutil.ReadDir()",
		"Use if err != nil pattern",
	}

	setup.CreateConversationHistory(prompts, responses)

	// Verify the file was created
	if !setup.ConversationHistoryExists() {
		t.Fatal("conversation history file not created")
	}

	// Read and verify content
	content := setup.ReadConversationHistory()
	if len(content) == 0 {
		t.Fatal("conversation history is empty")
	}

	// Verify all prompts are in the history
	for _, prompt := range prompts {
		if !strings.Contains(content, prompt) {
			t.Errorf("prompt not found in history: %s", prompt)
		}
	}
}
```

### Example 2: Testing Clear History Command

```go
func TestClearHistoryCommand(t *testing.T) {
	// Set up with existing history
	setup := testfixtures.NewHistoryTestSetup(t)
	setup.CreateConversationHistory(
		[]string{"test prompt"},
		[]string{"test response"},
	)

	// Verify history exists
	if !setup.ConversationHistoryExists() {
		t.Fatal("expected history to exist")
	}

	// Simulate clearing history
	err := setup.DeleteConversationHistory()
	if err != nil {
		t.Fatalf("failed to delete history: %v", err)
	}

	// Verify it's deleted
	if setup.ConversationHistoryExists() {
		t.Fatal("expected history to be deleted")
	}
}
```

## Testing Shell History Parsing

### Example 3: Testing Bash History Parsing

```go
func TestBashHistoryParsing(t *testing.T) {
	// Create bash history with test data
	setup := testfixtures.NewHistoryTestSetup(t)
	commands := []string{
		"git clone https://github.com/example/repo.git",
		"cd repo",
		"make build",
		"./bin/app --help",
	}

	bashPath := setup.CreateBashHistory(commands)

	// Verify file location
	if bashPath != setup.GetBashHistoryPath() {
		t.Errorf("unexpected path: got %s, want %s", bashPath, setup.GetBashHistoryPath())
	}

	// Test parsing (your actual parser function)
	// result, err := parseHistory.BashHistory(bashPath, 10)
	// if err != nil {
	//     t.Fatalf("parse error: %v", err)
	// }
	// assert command count matches
}
```

### Example 4: Testing Zsh History Parsing

```go
func TestZshHistoryParsing(t *testing.T) {
	// Create zsh history with specific timestamps
	setup := testfixtures.NewHistoryTestSetup(t)
	commands := []string{
		"ls -la /home",
		"cd /home/user/projects",
		"git pull origin main",
		"go test ./...",
	}

	zshPath := setup.CreateZshHistory(commands)

	// Read and verify format
	content := setup.ReadZshHistory()

	// Verify zsh format (: timestamp:duration;command)
	if !strings.Contains(content, ": ") {
		t.Fatal("expected zsh history format")
	}

	// Verify all commands are present
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Errorf("command not found: %s", cmd)
		}
	}
}
```

### Example 5: Testing Fish History Parsing

```go
func TestFishHistoryParsing(t *testing.T) {
	setup := testfixtures.NewHistoryTestSetup(t)
	commands := []string{
		"docker ps",
		"docker run -it ubuntu bash",
		"docker logs container-id",
	}

	fishPath := setup.CreateFishHistory(commands)
	content := setup.ReadFishHistory()

	// Verify fish format
	if !strings.Contains(content, "- cmd: ") {
		t.Fatal("expected fish history format '- cmd:'")
	}
	if !strings.Contains(content, "  when: ") {
		t.Fatal("expected fish history format '  when:'")
	}

	// Verify all commands present
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Errorf("command not found in fish history: %s", cmd)
		}
	}
}
```

## Testing Conversation History

### Example 6: Multiple Conversation Entries

```go
func TestMultipleConversations(t *testing.T) {
	setup := testfixtures.NewHistoryTestSetup(t)

	// Create multiple conversation entries
	prompts := []string{
		"How do I write a test in Go?",
		"What's the difference between * and & in Go?",
		"How do I handle JSON in Go?",
		"What are goroutines?",
	}

	responses := []string{
		"Use the testing package and create *_test.go files",
		"* dereferences pointers, & takes addresses",
		"Use encoding/json with json struct tags",
		"Goroutines are lightweight threads managed by Go runtime",
	}

	setup.CreateConversationHistory(prompts, responses)

	content := setup.ReadConversationHistory()
	lines := strings.Split(strings.TrimSpace(content), "\n")

	if len(lines) != len(prompts) {
		t.Errorf("expected %d lines, got %d", len(prompts), len(lines))
	}
}
```

### Example 7: Empty Conversation History

```go
func TestEmptyConversationHistory(t *testing.T) {
	setup := EmptyHistoryTestSetup(t)

	// Verify no history files exist initially
	if setup.ConversationHistoryExists() {
		t.Error("expected no conversation history")
	}

	// Create empty history
	setup.CreateConversationHistory([]string{}, []string{})

	// File should exist but might be empty
	if !setup.ConversationHistoryExists() {
		t.Error("expected conversation history file to be created")
	}
}
```

## Advanced Scenarios

### Example 8: Using Builder Pattern for Complex Setup

```go
func TestComplexHistoryScenario(t *testing.T) {
	setup := testfixtures.NewHistoryTestBuilder(t).
		WithBashHistory([]string{
			"ls -la",
			"cd /home/user/projects",
			"git init",
			"git add .",
			"git commit -m 'initial commit'",
		}).
		WithZshHistory([]string{
			"pwd",
			"echo $HOME",
			"export PATH=/usr/local/go/bin:$PATH",
		}).
		WithFishHistory([]string{
			"fish",
			"echo hello",
		}).
		WithConversationHistory(
			[]string{"How do I initialize a git repo?", "How do I commit?"},
			[]string{"Use git init", "Use git commit -m"},
		).
		Build()

	// Now all history types are available
	bashContent := setup.ReadBashHistory()
	zshContent := setup.ReadZshHistory()
	fishContent := setup.ReadFishHistory()
	conversationContent := setup.ReadConversationHistory()

	// All should have content
	if len(bashContent) == 0 || len(zshContent) == 0 {
		t.Fatal("expected all history content to be populated")
	}
}
```

### Example 9: Using Sample Data

```go
func TestWithStandardSampleData(t *testing.T) {
	// Create setup with all sample data pre-populated
	setup := testfixtures.StandardHistoryTestSetup(t)

	// All history files exist with sample content
	tests := []struct {
		name  string
		check func() bool
		path  string
	}{
		{"bash history", setup.BashHistoryExists, setup.GetBashHistoryPath()},
		{"zsh history", setup.ZshHistoryExists, setup.GetZshHistoryPath()},
		{"fish history", setup.FishHistoryExists, setup.GetFishHistoryPath()},
		{"conversation history", setup.ConversationHistoryExists, setup.GetConversationHistoryPath()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.check() {
				t.Errorf("expected %s to exist at %s", tt.name, tt.path)
			}
		})
	}
}
```

### Example 10: Testing File Operations

```go
func TestHistoryFileOperations(t *testing.T) {
	setup := testfixtures.NewHistoryTestSetup(t)

	// Test creating and reading
	setup.CreateBashHistory([]string{"echo test"})
	if !setup.BashHistoryExists() {
		t.Fatal("bash history should exist after creation")
	}

	// Test reading
	content := setup.ReadBashHistory()
	if !strings.Contains(content, "echo test") {
		t.Fatal("bash history should contain the command")
	}

	// Test deleting
	err := setup.DeleteBashHistory()
	if err != nil {
		t.Fatalf("failed to delete: %v", err)
	}

	if setup.BashHistoryExists() {
		t.Fatal("bash history should not exist after deletion")
	}
}
```

### Example 11: Testing with Custom Time

```go
func TestHistoryWithSpecificTimestamps(t *testing.T) {
	// Create generator with specific start time
	startTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	gen := testfixtures.NewSimpleHistoryGeneratorWithTime(startTime)

	commands := []string{
		"ls",
		"cd home",
		"pwd",
	}

	// Generate zsh history with known timestamps
	zshContent := gen.ZshHistory(commands)

	// Verify timestamps are in order
	if !strings.Contains(zshContent, "1704110400") { // 2024-01-01 12:00:00 UTC
		t.Error("expected start timestamp in zsh history")
	}
}
```

### Example 12: Testing Path Handling

```go
func TestHistoryPathHandling(t *testing.T) {
	setup := testfixtures.NewHistoryTestSetup(t)

	// Verify all paths are correct
	basePath := setup.BaseDir

	paths := map[string]string{
		"bash":         setup.GetBashHistoryPath(),
		"zsh":          setup.GetZshHistoryPath(),
		"fish":         setup.GetFishHistoryPath(),
		"conversation": setup.GetConversationHistoryPath(),
	}

	for name, path := range paths {
		if !strings.HasPrefix(path, basePath) {
			t.Errorf("%s path not under base directory: %s", name, path)
		}
	}

	// Create files and verify paths are valid
	setup.CreateBashHistory([]string{"test"})
	if _, err := os.Stat(setup.GetBashHistoryPath()); err != nil {
		t.Errorf("bash history path invalid: %v", err)
	}
}
```

## Common Patterns

### Pattern 1: Testing with Fresh Setup Per Test

```go
func TestFeature1(t *testing.T) {
	setup := testfixtures.NewHistoryTestSetup(t)
	// Each test gets a fresh environment
}

func TestFeature2(t *testing.T) {
	setup := testfixtures.NewHistoryTestSetup(t)
	// Clean slate for this test
}
```

### Pattern 2: Subtests with Shared Setup

```go
func TestHistoryFeatures(t *testing.T) {
	setup := testfixtures.StandardHistoryTestSetup(t)

	t.Run("read history", func(t *testing.T) {
		content := setup.ReadConversationHistory()
		if len(content) == 0 {
			t.Error("expected history content")
		}
	})

	t.Run("clear history", func(t *testing.T) {
		_ = setup.DeleteConversationHistory()
		if setup.ConversationHistoryExists() {
			t.Error("expected history to be deleted")
		}
	})
}
```

### Pattern 3: Error Handling in Tests

```go
func TestHistoryErrorHandling(t *testing.T) {
	setup := testfixtures.NewHistoryTestSetup(t)

	// Try to read non-existent file
	historyPath := setup.GetBashHistoryPath()
	_, err := os.ReadFile(historyPath)
	if !os.IsNotExist(err) {
		t.Error("expected file not found error")
	}
}
```

## Tips and Best Practices

1. **Use StandardHistoryTestSetup** when you need all formats with sample data
2. **Use Builder pattern** for fine-grained control over what gets created
3. **Use SampleCommands(), SamplePrompts(), SampleResponses()** for consistent test data
4. **Let t.Cleanup() handle cleanup** - don't manually call Cleanup()
5. **Use subtests** for testing multiple scenarios with the same setup
6. **Name your tests clearly** indicating what scenario they test
7. **Use table-driven tests** for testing multiple commands or formats
8. **Create helper functions** in your test files for repeated patterns

## Running Tests

```bash
# Run all tests in the testfixtures package
go test ./internal/cli/testfixtures -v

# Run a specific test
go test ./internal/cli/testfixtures -run TestBashHistoryParsing -v

# Run with coverage
go test ./internal/cli/testfixtures -cover

# Run your tests that use testfixtures
go test ./internal/cli -v
```
