# Test Fixtures for History Package

The `testfixtures` package provides comprehensive utilities for testing the history commands and functionality in the `how` CLI application. It includes helpers for setting up temporary directories, generating test data in various shell history formats, and managing test environments.

## Overview

This package provides:

- **Temporary Directory Management**: Create and manage temporary directories for testing
- **History Data Generators**: Generate test data in bash, zsh, and fish shell history formats
- **Test Setup Helpers**: `HistoryTestSetup` struct for managing test environments
- **Builder Pattern**: `HistoryTestBuilder` for fluent test scenario creation
- **Sample Data**: Pre-defined sample commands, prompts, and responses for testing

## Quick Start

### Basic Test Setup

```go
func TestMyHistoryFeature(t *testing.T) {
    // Create a basic test setup
    setup := NewHistoryTestSetup(t)

    // Create bash history with test data
    setup.CreateBashHistory([]string{
        "ls -la",
        "cd /home",
        "git status",
    })

    // Verify history was created
    if !setup.BashHistoryExists() {
        t.Fatal("bash history not created")
    }

    // Read and use the history
    content := setup.ReadBashHistory()
    t.Log("History content:", content)
}
```

### Using the Builder Pattern

```go
func TestComplexScenario(t *testing.T) {
    // Build a complete test environment
    setup := NewHistoryTestBuilder(t).
        WithBashHistory([]string{"ls", "pwd", "git status"}).
        WithZshHistory([]string{"cd /home", "make build"}).
        WithConversationHistory(
            []string{"How do I test?"},
            []string{"Use 'go test ./...'"},
        ).
        Build()

    // Use setup for testing
    bashPath := setup.GetBashHistoryPath()
    zshPath := setup.GetZshHistoryPath()
}
```

### Using Sample Data

```go
func TestWithSampleData(t *testing.T) {
    setup := StandardHistoryTestSetup(t)

    // All history files are created with sample data
    bashContent := setup.ReadBashHistory()
    zshContent := setup.ReadZshHistory()
    fishContent := setup.ReadFishHistory()
    conversationContent := setup.ReadConversationHistory()
}
```

## Available Functions

### Temporary Directory Management

| Function | Description |
|----------|-------------|
| `TempDir(t *testing.T) string` | Creates a temporary directory |
| `TempConfigDir(t *testing.T) string` | Creates a temp directory with `.config/how` structure |
| `TempHistoryDir(t *testing.T) (baseDir, historyDir string)` | Creates directories for shell history files |
| `WriteFile(t, dir, filename, content string) string` | Writes a test file |
| `ReadFile(t, filePath string) string` | Reads a test file |
| `FileExists(t, filePath string) bool` | Checks if a file exists |
| `FileDoesNotExist(t, filePath string) bool` | Checks if a file doesn't exist |

### History Generators

#### Bash History
- `GenerateBashHistory(commands []string) string` - Generate bash history format
- `CreateBashHistoryFile(t, dir string, commands []string) string` - Create a bash history file

#### Zsh History
- `GenerateZshHistory(commands []string, startTime time.Time) string` - Generate zsh history format
- `CreateZshHistoryFile(t, dir string, commands []string, startTime time.Time) string` - Create a zsh history file

#### Fish History
- `GenerateFishHistory(commands []string, startTime time.Time) string` - Generate fish history format
- `CreateFishHistoryFile(t, dir string, commands []string, startTime time.Time) string` - Create a fish history file

#### Conversation History
- `GenerateConversationHistory(prompts, responses []string) string` - Generate conversation history format
- `CreateConversationHistoryFile(t, dir, filename string, prompts, responses []string) string` - Create conversation history

### SimpleHistoryGenerator Utility

```go
// Create a generator with current time
gen := NewSimpleHistoryGenerator()

// Or with a specific time
gen := NewSimpleHistoryGeneratorWithTime(startTime)

// Generate different formats
bashContent := gen.BashHistory(commands)
zshContent := gen.ZshHistory(commands)
fishContent := gen.FishHistory(commands)
convContent := gen.ConversationHistory(prompts, responses)
```

### HistoryTestSetup

The `HistoryTestSetup` struct manages test environments:

```go
type HistoryTestSetup struct {
    BaseDir     string          // Temporary base directory
    HistoryDir  string          // History subdirectory
    ConfigDir   string          // Config directory
}
```

#### Methods

| Method | Description |
|--------|-------------|
| `GetBashHistoryPath() string` | Returns path to bash history file |
| `GetZshHistoryPath() string` | Returns path to zsh history file |
| `GetFishHistoryPath() string` | Returns path to fish history file |
| `GetConversationHistoryPath() string` | Returns path to conversation history file |
| `CreateBashHistory(commands []string) string` | Creates bash history file |
| `CreateZshHistory(commands []string) string` | Creates zsh history file |
| `CreateFishHistory(commands []string) string` | Creates fish history file |
| `CreateConversationHistory(prompts, responses []string) string` | Creates conversation history file |
| `BashHistoryExists() bool` | Checks if bash history exists |
| `ZshHistoryExists() bool` | Checks if zsh history exists |
| `FishHistoryExists() bool` | Checks if fish history exists |
| `ConversationHistoryExists() bool` | Checks if conversation history exists |
| `ReadBashHistory() string` | Reads bash history content |
| `ReadZshHistory() string` | Reads zsh history content |
| `ReadFishHistory() string` | Reads fish history content |
| `ReadConversationHistory() string` | Reads conversation history content |
| `DeleteBashHistory() error` | Deletes bash history file |
| `DeleteZshHistory() error` | Deletes zsh history file |
| `DeleteFishHistory() error` | Deletes fish history file |
| `DeleteConversationHistory() error` | Deletes conversation history file |
| `ClearAllHistory()` | Deletes all history files |
| `Cleanup()` | Removes all temporary directories |

### HistoryTestBuilder

Fluent builder for creating test scenarios:

```go
builder := NewHistoryTestBuilder(t).
    WithBashHistory([]string{"cmd1", "cmd2"}).
    WithZshHistory([]string{"cmd3"}).
    WithFishHistory([]string{"cmd4"}).
    WithConversationHistory(prompts, responses).
    Build()  // Returns *HistoryTestSetup
```

### Sample Data Functions

| Function | Description |
|----------|-------------|
| `SampleCommands() []string` | Returns sample shell commands |
| `SamplePrompts() []string` | Returns sample AI prompts |
| `SampleResponses() []string` | Returns sample AI responses |
| `EmptyHistoryTestSetup(t) *HistoryTestSetup` | Creates setup with no history files |
| `StandardHistoryTestSetup(t) *HistoryTestSetup` | Creates setup with all sample data |

## History File Formats

### Bash History Format
Simple newline-separated commands:
```
ls -la
cd /home
git status
```

### Zsh History Format
Extended history with timestamp and duration:
```
: 1640000000:0;ls -la
: 1640000060:0;cd /home
: 1640000120:0;git status
```

### Fish History Format
YAML-like format with command and timestamp:
```
- cmd: ls -la
  when: 1640000000
- cmd: cd /home
  when: 1640000060
- cmd: git status
  when: 1640000120
```

### Conversation History Format
Pipe-separated prompt and response:
```
PROMPT: How do I run tests? | RESPONSE: Use 'go test ./...'
PROMPT: How do I commit? | RESPONSE: Use 'git commit'
```

## Automatic Cleanup

All temporary directories and files are automatically cleaned up via `t.Cleanup()` when the test completes. No manual cleanup is required.

## Example Tests

See `example_test.go` for comprehensive examples of using the test fixtures:

- `TestTempDirCreation` - Basic directory creation
- `TestBashHistoryGeneration` - Bash history format
- `TestZshHistoryGeneration` - Zsh history format
- `TestFishHistoryGeneration` - Fish history format
- `TestHistoryTestSetupCreation` - Using HistoryTestSetup
- `TestHistoryTestBuilderFluent` - Using the builder pattern
- `TestStandardHistoryTestSetup` - Using sample data

## Integration with CLI Commands

When testing CLI history commands, you'll need to mock or override `getHistoryFile()` to use the test setup paths:

```go
func TestHistoryCommand(t *testing.T) {
    setup := NewHistoryTestSetup(t)
    setup.CreateConversationHistory(
        []string{"prompt1"},
        []string{"response1"},
    )

    // Configure your test to use setup.GetConversationHistoryPath()
    historyPath := setup.GetConversationHistoryPath()
    // ... test code using historyPath
}
```

## Best Practices

1. **Use StandardHistoryTestSetup for simple tests**: If you need all history formats with sample data
2. **Use NewHistoryTestBuilder for complex scenarios**: For fine-grained control
3. **Let cleanup happen automatically**: Don't call `Cleanup()` manually
4. **Use constants for repeated test data**: Create helper functions for common scenarios
5. **Document your setup**: Add comments explaining what each test is testing

## Running Tests

```bash
# Run all tests in the package
go test ./internal/cli/testfixtures -v

# Run specific test
go test ./internal/cli/testfixtures -v -run TestBashHistoryGeneration

# Run with coverage
go test ./internal/cli/testfixtures -cover
```
