package testfixtures

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// HistoryTestSetup manages setup and teardown of history test environments.
type HistoryTestSetup struct {
	BaseDir     string
	HistoryDir  string
	ConfigDir   string
	TempDirs    []string
	t            *testing.T
}

// NewHistoryTestSetup creates a new test setup for history tests.
func NewHistoryTestSetup(t *testing.T) *HistoryTestSetup {
	t.Helper()
	baseDir, historyDir := TempHistoryDir(t)
	configDir := filepath.Join(baseDir, ".config", "how")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("failed to create config directory: %v", err)
	}

	return &HistoryTestSetup{
		BaseDir:    baseDir,
		HistoryDir: historyDir,
		ConfigDir:  configDir,
		t:           t,
	}
}

// GetBashHistoryPath returns the path to the bash history file.
func (h *HistoryTestSetup) GetBashHistoryPath() string {
	return filepath.Join(h.BaseDir, ".bash_history")
}

// GetZshHistoryPath returns the path to the zsh history file.
func (h *HistoryTestSetup) GetZshHistoryPath() string {
	return filepath.Join(h.BaseDir, ".zsh_history")
}

// GetFishHistoryPath returns the path to the fish history file.
func (h *HistoryTestSetup) GetFishHistoryPath() string {
	return filepath.Join(h.HistoryDir, "fish_history")
}

// GetConversationHistoryPath returns the path to the conversation history file.
func (h *HistoryTestSetup) GetConversationHistoryPath() string {
	return filepath.Join(h.ConfigDir, "history.txt")
}

// CreateBashHistory creates a bash history file with the given commands.
func (h *HistoryTestSetup) CreateBashHistory(commands []string) string {
	h.t.Helper()
	return CreateBashHistoryFile(h.t, h.BaseDir, commands)
}

// CreateZshHistory creates a zsh history file with the given commands.
func (h *HistoryTestSetup) CreateZshHistory(commands []string) string {
	h.t.Helper()
	startTime := time.Now().Add(-time.Duration(len(commands)) * time.Minute)
	return CreateZshHistoryFile(h.t, h.BaseDir, commands, startTime)
}

// CreateFishHistory creates a fish history file with the given commands.
func (h *HistoryTestSetup) CreateFishHistory(commands []string) string {
	h.t.Helper()
	startTime := time.Now().Add(-time.Duration(len(commands)) * time.Minute)
	return CreateFishHistoryFile(h.t, h.HistoryDir, commands, startTime)
}

// CreateConversationHistory creates a conversation history file.
func (h *HistoryTestSetup) CreateConversationHistory(prompts []string, responses []string) string {
	h.t.Helper()
	return CreateConversationHistoryFile(h.t, h.ConfigDir, "history.txt", prompts, responses)
}

// BashHistoryExists checks if bash history file exists.
func (h *HistoryTestSetup) BashHistoryExists() bool {
	h.t.Helper()
	return FileExists(h.t, h.GetBashHistoryPath())
}

// ZshHistoryExists checks if zsh history file exists.
func (h *HistoryTestSetup) ZshHistoryExists() bool {
	h.t.Helper()
	return FileExists(h.t, h.GetZshHistoryPath())
}

// FishHistoryExists checks if fish history file exists.
func (h *HistoryTestSetup) FishHistoryExists() bool {
	h.t.Helper()
	return FileExists(h.t, h.GetFishHistoryPath())
}

// ConversationHistoryExists checks if conversation history file exists.
func (h *HistoryTestSetup) ConversationHistoryExists() bool {
	h.t.Helper()
	return FileExists(h.t, h.GetConversationHistoryPath())
}

// ReadBashHistory reads the bash history file.
func (h *HistoryTestSetup) ReadBashHistory() string {
	h.t.Helper()
	return ReadFile(h.t, h.GetBashHistoryPath())
}

// ReadZshHistory reads the zsh history file.
func (h *HistoryTestSetup) ReadZshHistory() string {
	h.t.Helper()
	return ReadFile(h.t, h.GetZshHistoryPath())
}

// ReadFishHistory reads the fish history file.
func (h *HistoryTestSetup) ReadFishHistory() string {
	h.t.Helper()
	return ReadFile(h.t, h.GetFishHistoryPath())
}

// ReadConversationHistory reads the conversation history file.
func (h *HistoryTestSetup) ReadConversationHistory() string {
	h.t.Helper()
	return ReadFile(h.t, h.GetConversationHistoryPath())
}

// DeleteBashHistory deletes the bash history file.
func (h *HistoryTestSetup) DeleteBashHistory() error {
	return os.Remove(h.GetBashHistoryPath())
}

// DeleteZshHistory deletes the zsh history file.
func (h *HistoryTestSetup) DeleteZshHistory() error {
	return os.Remove(h.GetZshHistoryPath())
}

// DeleteFishHistory deletes the fish history file.
func (h *HistoryTestSetup) DeleteFishHistory() error {
	return os.Remove(h.GetFishHistoryPath())
}

// DeleteConversationHistory deletes the conversation history file.
func (h *HistoryTestSetup) DeleteConversationHistory() error {
	return os.Remove(h.GetConversationHistoryPath())
}

// ClearAllHistory deletes all history files.
func (h *HistoryTestSetup) ClearAllHistory() {
	h.t.Helper()
	_ = h.DeleteBashHistory()
	_ = h.DeleteZshHistory()
	_ = h.DeleteFishHistory()
	_ = h.DeleteConversationHistory()
}

// Cleanup removes all temporary directories created during setup.
// This is automatically called by the test framework via t.Cleanup().
func (h *HistoryTestSetup) Cleanup() {
	_ = os.RemoveAll(h.BaseDir)
}

// HistoryTestBuilder provides a fluent interface for building history test scenarios.
type HistoryTestBuilder struct {
	setup     *HistoryTestSetup
	t          *testing.T
	bashCmds  []string
	zshCmds   []string
	fishCmds  []string
	prompts   []string
	responses []string
}

// NewHistoryTestBuilder creates a new builder for history test scenarios.
func NewHistoryTestBuilder(t *testing.T) *HistoryTestBuilder {
	t.Helper()
	return &HistoryTestBuilder{
		setup: NewHistoryTestSetup(t),
		t:      t,
	}
}

// WithBashHistory adds bash history commands to the test setup.
func (b *HistoryTestBuilder) WithBashHistory(commands []string) *HistoryTestBuilder {
	b.bashCmds = commands
	return b
}

// WithZshHistory adds zsh history commands to the test setup.
func (b *HistoryTestBuilder) WithZshHistory(commands []string) *HistoryTestBuilder {
	b.zshCmds = commands
	return b
}

// WithFishHistory adds fish history commands to the test setup.
func (b *HistoryTestBuilder) WithFishHistory(commands []string) *HistoryTestBuilder {
	b.fishCmds = commands
	return b
}

// WithConversationHistory adds conversation history to the test setup.
func (b *HistoryTestBuilder) WithConversationHistory(prompts, responses []string) *HistoryTestBuilder {
	b.prompts = prompts
	b.responses = responses
	return b
}

// Build creates the actual history files based on the builder configuration.
func (b *HistoryTestBuilder) Build() *HistoryTestSetup {
	b.t.Helper()

	if len(b.bashCmds) > 0 {
		b.setup.CreateBashHistory(b.bashCmds)
	}
	if len(b.zshCmds) > 0 {
		b.setup.CreateZshHistory(b.zshCmds)
	}
	if len(b.fishCmds) > 0 {
		b.setup.CreateFishHistory(b.fishCmds)
	}
	if len(b.prompts) > 0 {
		b.setup.CreateConversationHistory(b.prompts, b.responses)
	}

	return b.setup
}

// EmptyHistoryTestSetup creates a test setup with no history files.
func EmptyHistoryTestSetup(t *testing.T) *HistoryTestSetup {
	t.Helper()
	return NewHistoryTestSetup(t)
}

// StandardHistoryTestSetup creates a test setup with standard sample history.
func StandardHistoryTestSetup(t *testing.T) *HistoryTestSetup {
	t.Helper()
	return NewHistoryTestBuilder(t).
		WithBashHistory(SampleCommands()).
		WithZshHistory(SampleCommands()).
		WithFishHistory(SampleCommands()).
		WithConversationHistory(SamplePrompts(), SampleResponses()).
		Build()
}
