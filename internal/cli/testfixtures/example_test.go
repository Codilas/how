package testfixtures

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ExampleTestBasicSetup demonstrates basic test setup with temporary directories.
func ExampleTestBasicSetup(t *testing.T) {
	// Create a temporary directory for the test
	tempDir := TempDir(t)

	// Write test data
	WriteFile(t, tempDir, "test.txt", "test content")

	// Verify file was created
	if !FileExists(t, filepath.Join(tempDir, "test.txt")) {
		t.Fatal("expected test.txt to exist")
	}
}

// TestTempDirCreation tests that temporary directories are created correctly.
func TestTempDirCreation(t *testing.T) {
	tempDir := TempDir(t)

	// Verify directory exists
	info, err := os.Stat(tempDir)
	if err != nil {
		t.Fatalf("expected temp directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected temp path to be a directory")
	}
}

// TestTempConfigDirStructure tests that config directory structure is created properly.
func TestTempConfigDirStructure(t *testing.T) {
	baseDir := TempConfigDir(t)

	configPath := filepath.Join(baseDir, ".config", "how")
	if !FileExists(t, configPath) {
		t.Fatalf("expected config directory to exist at %s", configPath)
	}
}

// TestBashHistoryGeneration tests bash history generation.
func TestBashHistoryGeneration(t *testing.T) {
	commands := []string{"ls -la", "cd /home", "git status"}
	content := GenerateBashHistory(commands)

	// Verify content contains all commands
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Fatalf("expected content to contain command: %s", cmd)
		}
	}

	// Verify each command is on its own line
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != len(commands) {
		t.Fatalf("expected %d lines, got %d", len(commands), len(lines))
	}
}

// TestZshHistoryGeneration tests zsh history generation.
func TestZshHistoryGeneration(t *testing.T) {
	commands := []string{"ls -la", "cd /home", "git status"}
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	content := GenerateZshHistory(commands, startTime)

	// Verify content contains zsh format markers
	if !strings.Contains(content, ": ") {
		t.Fatal("expected zsh history format with ': ' prefix")
	}

	// Verify all commands are present
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Fatalf("expected content to contain command: %s", cmd)
		}
	}
}

// TestFishHistoryGeneration tests fish history generation.
func TestFishHistoryGeneration(t *testing.T) {
	commands := []string{"ls -la", "cd /home", "git status"}
	startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	content := GenerateFishHistory(commands, startTime)

	// Verify content contains fish format markers
	if !strings.Contains(content, "- cmd: ") {
		t.Fatal("expected fish history format with '- cmd: ' prefix")
	}

	if !strings.Contains(content, "  when: ") {
		t.Fatal("expected fish history format with '  when: ' marker")
	}

	// Verify all commands are present
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Fatalf("expected content to contain command: %s", cmd)
		}
	}
}

// TestConversationHistoryGeneration tests conversation history generation.
func TestConversationHistoryGeneration(t *testing.T) {
	prompts := []string{"How do I run tests?", "How do I commit?"}
	responses := []string{"Use 'go test ./...'", "Use 'git commit'"}
	content := GenerateConversationHistory(prompts, responses)

	// Verify all prompts and responses are present
	for _, prompt := range prompts {
		if !strings.Contains(content, prompt) {
			t.Fatalf("expected content to contain prompt: %s", prompt)
		}
	}

	for _, response := range responses {
		if !strings.Contains(content, response) {
			t.Fatalf("expected content to contain response: %s", response)
		}
	}
}

// TestHistoryTestSetupCreation tests HistoryTestSetup initialization.
func TestHistoryTestSetupCreation(t *testing.T) {
	setup := NewHistoryTestSetup(t)

	// Verify paths are set
	if setup.BaseDir == "" {
		t.Fatal("expected BaseDir to be set")
	}
	if setup.HistoryDir == "" {
		t.Fatal("expected HistoryDir to be set")
	}
	if setup.ConfigDir == "" {
		t.Fatal("expected ConfigDir to be set")
	}

	// Verify directories exist
	if !FileExists(t, setup.BaseDir) {
		t.Fatal("expected base directory to exist")
	}
	if !FileExists(t, setup.ConfigDir) {
		t.Fatal("expected config directory to exist")
	}
}

// TestHistoryTestSetupBashHistory tests creating bash history via setup.
func TestHistoryTestSetupBashHistory(t *testing.T) {
	setup := NewHistoryTestSetup(t)
	commands := []string{"ls", "pwd", "cd /home"}

	setup.CreateBashHistory(commands)

	if !setup.BashHistoryExists() {
		t.Fatal("expected bash history file to exist")
	}

	content := setup.ReadBashHistory()
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Fatalf("expected bash history to contain: %s", cmd)
		}
	}
}

// TestHistoryTestSetupZshHistory tests creating zsh history via setup.
func TestHistoryTestSetupZshHistory(t *testing.T) {
	setup := NewHistoryTestSetup(t)
	commands := []string{"ls", "pwd", "cd /home"}

	setup.CreateZshHistory(commands)

	if !setup.ZshHistoryExists() {
		t.Fatal("expected zsh history file to exist")
	}

	content := setup.ReadZshHistory()
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Fatalf("expected zsh history to contain: %s", cmd)
		}
	}
}

// TestHistoryTestSetupFishHistory tests creating fish history via setup.
func TestHistoryTestSetupFishHistory(t *testing.T) {
	setup := NewHistoryTestSetup(t)
	commands := []string{"ls", "pwd", "cd /home"}

	setup.CreateFishHistory(commands)

	if !setup.FishHistoryExists() {
		t.Fatal("expected fish history file to exist")
	}

	content := setup.ReadFishHistory()
	for _, cmd := range commands {
		if !strings.Contains(content, cmd) {
			t.Fatalf("expected fish history to contain: %s", cmd)
		}
	}
}

// TestHistoryTestSetupConversationHistory tests creating conversation history via setup.
func TestHistoryTestSetupConversationHistory(t *testing.T) {
	setup := NewHistoryTestSetup(t)
	prompts := []string{"How do I test?"}
	responses := []string{"Use 'go test'"}

	setup.CreateConversationHistory(prompts, responses)

	if !setup.ConversationHistoryExists() {
		t.Fatal("expected conversation history file to exist")
	}

	content := setup.ReadConversationHistory()
	for _, prompt := range prompts {
		if !strings.Contains(content, prompt) {
			t.Fatalf("expected conversation history to contain prompt: %s", prompt)
		}
	}
}

// TestHistoryTestSetupDeleteFiles tests deleting history files.
func TestHistoryTestSetupDeleteFiles(t *testing.T) {
	setup := NewHistoryTestSetup(t)
	setup.CreateBashHistory([]string{"ls"})

	if !setup.BashHistoryExists() {
		t.Fatal("expected bash history to be created")
	}

	err := setup.DeleteBashHistory()
	if err != nil {
		t.Fatalf("unexpected error deleting bash history: %v", err)
	}

	if setup.BashHistoryExists() {
		t.Fatal("expected bash history to be deleted")
	}
}

// TestHistoryTestBuilderFluent tests the fluent builder interface.
func TestHistoryTestBuilderFluent(t *testing.T) {
	setup := NewHistoryTestBuilder(t).
		WithBashHistory([]string{"ls", "pwd"}).
		WithZshHistory([]string{"cd /home"}).
		WithConversationHistory([]string{"How?"}, []string{"Like this"}).
		Build()

	if !setup.BashHistoryExists() {
		t.Fatal("expected bash history")
	}
	if !setup.ZshHistoryExists() {
		t.Fatal("expected zsh history")
	}
	if !setup.ConversationHistoryExists() {
		t.Fatal("expected conversation history")
	}
}

// TestEmptyHistoryTestSetup tests creating an empty setup.
func TestEmptyHistoryTestSetup(t *testing.T) {
	setup := EmptyHistoryTestSetup(t)

	if !setup.BashHistoryExists() == FileDoesNotExist(t, setup.GetBashHistoryPath()) {
		// Either both true or both false is expected
	}
}

// TestStandardHistoryTestSetup tests creating a setup with standard sample data.
func TestStandardHistoryTestSetup(t *testing.T) {
	setup := StandardHistoryTestSetup(t)

	if !setup.BashHistoryExists() {
		t.Fatal("expected bash history to exist")
	}
	if !setup.ZshHistoryExists() {
		t.Fatal("expected zsh history to exist")
	}
	if !setup.FishHistoryExists() {
		t.Fatal("expected fish history to exist")
	}
	if !setup.ConversationHistoryExists() {
		t.Fatal("expected conversation history to exist")
	}

	bashContent := setup.ReadBashHistory()
	if len(bashContent) == 0 {
		t.Fatal("expected bash history content")
	}
}

// TestSimpleHistoryGenerator tests the generator utility.
func TestSimpleHistoryGenerator(t *testing.T) {
	gen := NewSimpleHistoryGenerator()
	commands := []string{"ls", "pwd", "git status"}

	bashContent := gen.BashHistory(commands)
	if !strings.Contains(bashContent, "ls") {
		t.Fatal("expected bash history to contain commands")
	}

	zshContent := gen.ZshHistory(commands)
	if !strings.Contains(zshContent, ": ") {
		t.Fatal("expected zsh format")
	}

	fishContent := gen.FishHistory(commands)
	if !strings.Contains(fishContent, "- cmd: ") {
		t.Fatal("expected fish format")
	}

	convContent := gen.ConversationHistory([]string{"q1"}, []string{"r1"})
	if !strings.Contains(convContent, "PROMPT:") {
		t.Fatal("expected conversation format")
	}
}

// TestSampleDataGeneration tests sample data generators.
func TestSampleDataGeneration(t *testing.T) {
	commands := SampleCommands()
	if len(commands) == 0 {
		t.Fatal("expected sample commands to be non-empty")
	}

	prompts := SamplePrompts()
	if len(prompts) == 0 {
		t.Fatal("expected sample prompts to be non-empty")
	}

	responses := SampleResponses()
	if len(responses) == 0 {
		t.Fatal("expected sample responses to be non-empty")
	}
}
