# Test Fixtures Quick Reference

## Imports

```go
import "github.com/Codilas/how/internal/cli/testfixtures"
```

## Basic Setup (Most Common)

```go
// Standard setup with all sample data
setup := testfixtures.StandardHistoryTestSetup(t)

// Empty setup
setup := testfixtures.EmptyHistoryTestSetup(t)

// Custom setup
setup := testfixtures.NewHistoryTestSetup(t)
```

## Using the Builder

```go
setup := testfixtures.NewHistoryTestBuilder(t).
    WithBashHistory([]string{"ls", "pwd"}).
    WithZshHistory([]string{"cd /home"}).
    WithConversationHistory(prompts, responses).
    Build()
```

## Creating History Files

```go
setup := testfixtures.NewHistoryTestSetup(t)

// Create individual history types
setup.CreateBashHistory(commands)      // []string
setup.CreateZshHistory(commands)       // []string
setup.CreateFishHistory(commands)      // []string
setup.CreateConversationHistory(prompts, responses)  // []string, []string
```

## Getting File Paths

```go
setup := testfixtures.NewHistoryTestSetup(t)

bashPath := setup.GetBashHistoryPath()
zshPath := setup.GetZshHistoryPath()
fishPath := setup.GetFishHistoryPath()
conversationPath := setup.GetConversationHistoryPath()
```

## Checking File Existence

```go
if setup.BashHistoryExists() { ... }
if setup.ZshHistoryExists() { ... }
if setup.FishHistoryExists() { ... }
if setup.ConversationHistoryExists() { ... }
```

## Reading Files

```go
content := setup.ReadBashHistory()
content := setup.ReadZshHistory()
content := setup.ReadFishHistory()
content := setup.ReadConversationHistory()
```

## Deleting Files

```go
setup.DeleteBashHistory()
setup.DeleteZshHistory()
setup.DeleteFishHistory()
setup.DeleteConversationHistory()
setup.ClearAllHistory()  // Delete all
```

## Direct File Operations

```go
// Create temp directory
tempDir := testfixtures.TempDir(t)

// Create config directory
baseDir := testfixtures.TempConfigDir(t)

// Write file
filePath := testfixtures.WriteFile(t, dir, "filename", "content")

// Read file
content := testfixtures.ReadFile(t, filePath)

// Check existence
if testfixtures.FileExists(t, filePath) { ... }
if testfixtures.FileDoesNotExist(t, filePath) { ... }
```

## Generating Test Data

```go
// Generate content without file creation
bashContent := testfixtures.GenerateBashHistory(commands)
zshContent := testfixtures.GenerateZshHistory(commands, startTime)
fishContent := testfixtures.GenerateFishHistory(commands, startTime)
conversationContent := testfixtures.GenerateConversationHistory(prompts, responses)

// Create files directly
path := testfixtures.CreateBashHistoryFile(t, dir, commands)
path := testfixtures.CreateZshHistoryFile(t, dir, commands, startTime)
path := testfixtures.CreateFishHistoryFile(t, dir, commands, startTime)
path := testfixtures.CreateConversationHistoryFile(t, dir, "filename", prompts, responses)
```

## Using Generator Utility

```go
// Current time
gen := testfixtures.NewSimpleHistoryGenerator()

// Specific time
gen := testfixtures.NewSimpleHistoryGeneratorWithTime(startTime)

// Generate content
bash := gen.BashHistory(commands)
zsh := gen.ZshHistory(commands)
fish := gen.FishHistory(commands)
conv := gen.ConversationHistory(prompts, responses)
```

## Sample Data

```go
commands := testfixtures.SampleCommands()     // []string
prompts := testfixtures.SamplePrompts()       // []string
responses := testfixtures.SampleResponses()   // []string
```

## Common Patterns

### Pattern 1: Single History Type
```go
func TestBashHistoryReading(t *testing.T) {
    setup := testfixtures.NewHistoryTestSetup(t)
    setup.CreateBashHistory([]string{"ls", "pwd"})

    content := setup.ReadBashHistory()
    // Test content
}
```

### Pattern 2: Multiple History Types
```go
func TestMultipleFormats(t *testing.T) {
    setup := testfixtures.NewHistoryTestBuilder(t).
        WithBashHistory([]string{"ls"}).
        WithZshHistory([]string{"pwd"}).
        Build()

    bash := setup.ReadBashHistory()
    zsh := setup.ReadZshHistory()
}
```

### Pattern 3: Using Paths
```go
func TestWithPaths(t *testing.T) {
    setup := testfixtures.NewHistoryTestSetup(t)
    setup.CreateBashHistory([]string{"test"})

    path := setup.GetBashHistoryPath()
    // Use path in code being tested
}
```

### Pattern 4: Testing File Operations
```go
func TestFileOperations(t *testing.T) {
    setup := testfixtures.NewHistoryTestSetup(t)
    setup.CreateBashHistory([]string{"cmd"})

    if !setup.BashHistoryExists() {
        t.Fatal("file not created")
    }

    err := setup.DeleteBashHistory()
    if err != nil {
        t.Fatal(err)
    }

    if setup.BashHistoryExists() {
        t.Fatal("file not deleted")
    }
}
```

## Struct Properties

```go
setup.BaseDir      // Root temporary directory
setup.HistoryDir   // History subdirectory
setup.ConfigDir    // Config directory
```

## Methods by Category

### File Creation
- `CreateBashHistory([]string) string`
- `CreateZshHistory([]string) string`
- `CreateFishHistory([]string) string`
- `CreateConversationHistory([]string, []string) string`

### File Checking
- `BashHistoryExists() bool`
- `ZshHistoryExists() bool`
- `FishHistoryExists() bool`
- `ConversationHistoryExists() bool`

### File Reading
- `ReadBashHistory() string`
- `ReadZshHistory() string`
- `ReadFishHistory() string`
- `ReadConversationHistory() string`

### File Deletion
- `DeleteBashHistory() error`
- `DeleteZshHistory() error`
- `DeleteFishHistory() error`
- `DeleteConversationHistory() error`
- `ClearAllHistory()`

### Path Access
- `GetBashHistoryPath() string`
- `GetZshHistoryPath() string`
- `GetFishHistoryPath() string`
- `GetConversationHistoryPath() string`

## Tips

1. **Cleanup is automatic** - Use `t.Cleanup()`, don't call `Cleanup()` manually
2. **Use builder for complex scenarios** - More readable than multiple Create calls
3. **Use StandardHistoryTestSetup for quick tests** - Pre-populated with sample data
4. **SampleCommands/Prompts/Responses** - Great for consistent test data
5. **FileExists checks** - Useful for assertions

## Running Tests

```bash
# All tests in package
go test ./internal/cli/testfixtures -v

# Specific test
go test ./internal/cli/testfixtures -run TestBashHistoryGeneration -v

# With coverage
go test ./internal/cli/testfixtures -cover

# Verbose output
go test ./internal/cli/testfixtures -v -race
```

## Integration Example

```go
package cli

import (
	"testing"
	"github.com/Codilas/how/internal/cli/testfixtures"
)

func TestHistoryCommand(t *testing.T) {
	// Set up test environment
	setup := testfixtures.StandardHistoryTestSetup(t)

	// Get path to history file
	historyPath := setup.GetConversationHistoryPath()

	// Use in code being tested
	// (You would need to mock getHistoryFile() to return historyPath)

	// Verify results
	content := setup.ReadConversationHistory()
	if len(content) == 0 {
		t.Fatal("expected conversation history")
	}
}
```

## See Also

- `README.md` - Full documentation
- `USAGE.md` - 12 practical examples
- `example_test.go` - All features demonstrated
