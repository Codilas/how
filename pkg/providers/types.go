package providers

import "time"

// Context contains contextual information to include with prompts
type Context struct {
	// Current working directory and files
	WorkingDirectory string        `json:"working_directory,omitempty"`
	Files            []FileContext `json:"files,omitempty"`

	// Shell and command history
	Shell          string           `json:"shell,omitempty"`
	RecentCommands []CommandHistory `json:"recent_commands,omitempty"`

	// Environment information
	Environment map[string]string `json:"environment,omitempty"`

	// Git repository information
	Git *GitContext `json:"git,omitempty"`

	// Project information
	Project *ProjectContext `json:"project,omitempty"`

	// Conversation history for multi-turn conversations
	ConversationID  string         `json:"conversation_id,omitempty"`
	PreviousPrompts []HistoryEntry `json:"previous_prompts,omitempty"`
}

// FileContext represents a file in the current context
type FileContext struct {
	Path        string `json:"path"`
	Type        string `json:"type"` // "file", "directory"
	Size        int64  `json:"size,omitempty"`
	Content     string `json:"content,omitempty"` // For small, relevant files
	Summary     string `json:"summary,omitempty"` // For large files
	Language    string `json:"language,omitempty"`
	IsImportant bool   `json:"is_important"` // README, package.json, etc.
}

// CommandHistory represents a recent shell command
type CommandHistory struct {
	Command   string    `json:"command"`
	ExitCode  int       `json:"exit_code"`
	Timestamp time.Time `json:"timestamp"`
	Output    string    `json:"output,omitempty"` // Last few lines of output
}

// GitContext contains git repository information
type GitContext struct {
	Repository    string   `json:"repository,omitempty"`
	Branch        string   `json:"branch,omitempty"`
	CommitHash    string   `json:"commit_hash,omitempty"`
	Status        string   `json:"status,omitempty"` // git status output
	RecentCommits []string `json:"recent_commits,omitempty"`
	RemoteURL     string   `json:"remote_url,omitempty"`
}

// ProjectContext contains detected project information
type ProjectContext struct {
	Type         string            `json:"type"` // "nodejs", "python", "go", "rust", etc.
	Name         string            `json:"name,omitempty"`
	Version      string            `json:"version,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
	Scripts      map[string]string `json:"scripts,omitempty"`
	Framework    string            `json:"framework,omitempty"` // "react", "vue", "django", etc.
}

// HistoryEntry represents a previous conversation entry
type HistoryEntry struct {
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
}

// Response represents the response from an AI provider
type Response struct {
	// The main response text
	Text string `json:"text"`

	// Metadata
	Model          string        `json:"model"`
	Provider       string        `json:"provider"`
	TokensUsed     int           `json:"tokens_used,omitempty"`
	ResponseTime   time.Duration `json:"response_time"`
	ConversationID string        `json:"conversation_id,omitempty"`

	// Additional metadata
	Confidence float32  `json:"confidence,omitempty"`
	Tags       []string `json:"tags,omitempty"`
	Language   string   `json:"language,omitempty"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	Text     string                 `json:"text"`
	Done     bool                   `json:"done"`
	Error    error                  `json:"error,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SuggestedCommand represents a command the AI suggests
type SuggestedCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
	Safe        bool   `json:"safe"`     // Whether it's safe to auto-execute
	Category    string `json:"category"` // "file", "network", "system", etc.
}

// CodeBlock represents a code snippet in the response
type CodeBlock struct {
	Language    string `json:"language"`
	Code        string `json:"code"`
	Description string `json:"description,omitempty"`
	Filename    string `json:"filename,omitempty"`
}

// Reference represents a reference or link in the response
type Reference struct {
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
	Type  string `json:"type"` // "documentation", "tutorial", "example", etc.
}

// ProviderInfo contains information about a provider
type ProviderInfo struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // "anthropic", "openai", "local", etc.
	Model       string `json:"model"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
}

// Capabilities defines what features a provider supports
type Capabilities struct {
	Streaming          bool `json:"streaming"`
	FunctionCalling    bool `json:"function_calling"`
	CodeExecution      bool `json:"code_execution"`
	ImageAnalysis      bool `json:"image_analysis"`
	ConversationMemory bool `json:"conversation_memory"`
	MaxContextSize     int  `json:"max_context_size"`
	MaxTokens          int  `json:"max_tokens"`
}

// Config represents configuration for any provider
type Config struct {
	Type      string `json:"type"`
	APIKey    string `json:"api_key"`
	Model     string `json:"model"`
	BaseURL   string `json:"base_url"`
	MaxTokens int    `json:"max_tokens"`
}
