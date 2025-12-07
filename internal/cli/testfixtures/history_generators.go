package testfixtures

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

// BashHistoryData represents bash history test data.
type BashHistoryData struct {
	Commands []string
}

// GenerateBashHistory generates bash history content.
// Each command is on a separate line, matching bash history format.
func GenerateBashHistory(commands []string) string {
	return strings.Join(commands, "\n") + "\n"
}

// CreateBashHistoryFile creates a bash history file with the given commands.
func CreateBashHistoryFile(t *testing.T, dir string, commands []string) string {
	t.Helper()
	content := GenerateBashHistory(commands)
	return WriteFile(t, dir, ".bash_history", content)
}

// ZshHistoryData represents zsh history test data.
type ZshHistoryData struct {
	Commands  []string
	Durations []int // in seconds, optional
}

// GenerateZshHistory generates zsh extended history content.
// Format: : <timestamp>:<duration>;<command>
func GenerateZshHistory(commands []string, startTime time.Time) string {
	var lines []string
	for i, cmd := range commands {
		// Generate timestamp for each command (1 minute apart)
		ts := startTime.Add(time.Duration(i) * time.Minute).Unix()
		duration := 0 // default duration in seconds

		line := fmt.Sprintf(": %d:%d;%s", ts, duration, cmd)
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n") + "\n"
}

// CreateZshHistoryFile creates a zsh history file with the given commands.
func CreateZshHistoryFile(t *testing.T, dir string, commands []string, startTime time.Time) string {
	t.Helper()
	content := GenerateZshHistory(commands, startTime)
	return WriteFile(t, dir, ".zsh_history", content)
}

// FishHistoryData represents fish shell history test data.
type FishHistoryData struct {
	Commands   []string
	Timestamps []int64 // unix timestamps
}

// GenerateFishHistory generates fish shell history content in YAML format.
// Format:
// - cmd: <command>
//   when: <timestamp>
func GenerateFishHistory(commands []string, startTime time.Time) string {
	var lines []string
	for i, cmd := range commands {
		// Generate timestamp for each command (1 minute apart)
		ts := startTime.Add(time.Duration(i) * time.Minute).Unix()

		lines = append(lines, fmt.Sprintf("- cmd: %s", cmd))
		lines = append(lines, fmt.Sprintf("  when: %d", ts))
	}
	return strings.Join(lines, "\n") + "\n"
}

// CreateFishHistoryFile creates a fish shell history file with the given commands.
func CreateFishHistoryFile(t *testing.T, dir string, commands []string, startTime time.Time) string {
	t.Helper()
	content := GenerateFishHistory(commands, startTime)
	return WriteFile(t, dir, "fish_history", content)
}

// ConversationHistoryData represents conversation history test data.
type ConversationHistoryData struct {
	Prompts   []string
	Responses []string
}

// GenerateConversationHistory generates conversation history content.
// Format: one entry per line with format "PROMPT: <prompt> | RESPONSE: <response>"
func GenerateConversationHistory(prompts []string, responses []string) string {
	var lines []string
	maxLen := len(prompts)
	if len(responses) > maxLen {
		maxLen = len(responses)
	}

	for i := 0; i < maxLen; i++ {
		prompt := ""
		response := ""

		if i < len(prompts) {
			prompt = prompts[i]
		}
		if i < len(responses) {
			response = responses[i]
		}

		entry := fmt.Sprintf("PROMPT: %s | RESPONSE: %s", prompt, response)
		lines = append(lines, entry)
	}

	return strings.Join(lines, "\n") + "\n"
}

// CreateConversationHistoryFile creates a conversation history file.
func CreateConversationHistoryFile(t *testing.T, dir, filename string, prompts []string, responses []string) string {
	t.Helper()
	content := GenerateConversationHistory(prompts, responses)
	return WriteFile(t, dir, filename, content)
}

// SimpleHistoryGenerator provides convenient methods for generating test history data.
type SimpleHistoryGenerator struct {
	startTime time.Time
}

// NewSimpleHistoryGenerator creates a new history generator with the current time.
func NewSimpleHistoryGenerator() *SimpleHistoryGenerator {
	return &SimpleHistoryGenerator{
		startTime: time.Now(),
	}
}

// NewSimpleHistoryGeneratorWithTime creates a new history generator with a specific start time.
func NewSimpleHistoryGeneratorWithTime(startTime time.Time) *SimpleHistoryGenerator {
	return &SimpleHistoryGenerator{
		startTime: startTime,
	}
}

// BashHistory generates bash history content.
func (g *SimpleHistoryGenerator) BashHistory(commands []string) string {
	return GenerateBashHistory(commands)
}

// ZshHistory generates zsh history content.
func (g *SimpleHistoryGenerator) ZshHistory(commands []string) string {
	return GenerateZshHistory(commands, g.startTime)
}

// FishHistory generates fish history content.
func (g *SimpleHistoryGenerator) FishHistory(commands []string) string {
	return GenerateFishHistory(commands, g.startTime)
}

// ConversationHistory generates conversation history content.
func (g *SimpleHistoryGenerator) ConversationHistory(prompts []string, responses []string) string {
	return GenerateConversationHistory(prompts, responses)
}

// SampleCommands returns a slice of sample shell commands for testing.
func SampleCommands() []string {
	return []string{
		"ls -la",
		"cd /home/user/projects",
		"git status",
		"git commit -m 'feat: add new feature'",
		"go test ./...",
		"make build",
		"docker ps",
		"curl https://api.example.com",
	}
}

// SamplePrompts returns a slice of sample AI prompts for testing.
func SamplePrompts() []string {
	return []string{
		"How do I list files in Go?",
		"How do I run tests?",
		"What does this error mean?",
		"How do I commit changes to git?",
	}
}

// SampleResponses returns a slice of sample AI responses for testing.
func SampleResponses() []string {
	return []string{
		"You can use ioutil.ReadDir() or filepath.Walk() to list files.",
		"Use 'go test ./...' to run all tests in the current module.",
		"This error typically means a file was not found. Check your path.",
		"Use 'git add <files>' followed by 'git commit -m \"message\"'.",
	}
}
