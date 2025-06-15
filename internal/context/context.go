package context

import (
	"os"
	"path/filepath"

	"github.com/tzvonimir/how/internal/config"
	"github.com/tzvonimir/how/pkg/providers"
)

// Gatherer coordinates collecting context information
type Gatherer struct {
	config config.ContextConfig
}

// NewGatherer creates a new context gatherer
func NewGatherer(cfg config.ContextConfig) *Gatherer {
	return &Gatherer{config: cfg}
}

// Gather collects all context information based on configuration
func Gather(cfg config.ContextConfig) (*providers.Context, error) {
	gatherer := NewGatherer(cfg)
	return gatherer.GatherAll()
}

// GatherAll collects all available context information
func (g *Gatherer) GatherAll() (*providers.Context, error) {
	ctx := &providers.Context{}

	// Get current working directory
	if wd, err := os.Getwd(); err == nil {
		ctx.WorkingDirectory = wd
	}

	// Detect shell
	if shell := detectShell(); shell != "" {
		ctx.Shell = shell
	}

	// Gather file system context
	if g.config.IncludeFiles {
		if files, err := g.gatherFileContext(); err == nil {
			ctx.Files = files
		}
	}

	// Gather command history
	if g.config.IncludeHistory > 0 {
		if commands, err := g.gatherCommandHistory(g.config.IncludeHistory); err == nil {
			ctx.RecentCommands = commands
		}
	}

	// Gather environment information
	if g.config.IncludeEnvironment {
		ctx.Environment = g.gatherEnvironment()
	}

	// Gather git context
	if g.config.IncludeGit {
		if git, err := g.gatherGitContext(); err == nil {
			ctx.Git = git
		}
	}

	// Detect project context
	if project, err := g.gatherProjectContext(); err == nil {
		ctx.Project = project
	}

	return ctx, nil
}

// gatherFileContext analyzes files in current directory
func (g *Gatherer) gatherFileContext() ([]providers.FileContext, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var files []providers.FileContext

	// Read directory contents
	entries, err := os.ReadDir(wd)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		// Skip hidden files and excluded patterns
		if g.shouldExclude(entry.Name()) {
			continue
		}

		fileCtx := providers.FileContext{
			Path: entry.Name(),
			Type: getFileType(entry),
		}

		// Check if file is important (README, config files, etc.)
		fileCtx.IsImportant = isImportantFile(entry.Name())

		if entry.IsDir() {
			fileCtx.Type = "directory"
		} else {
			info, err := entry.Info()
			if err == nil {
				fileCtx.Size = info.Size()
				fileCtx.Language = detectLanguage(entry.Name())

				// Read content for small, important files
				if fileCtx.IsImportant && fileCtx.Size < 2048 {
					if content, err := os.ReadFile(filepath.Join(wd, entry.Name())); err == nil {
						fileCtx.Content = string(content)
					}
				}
			}
		}

		files = append(files, fileCtx)

		// Limit number of files to prevent context overflow
		if len(files) >= 50 {
			break
		}
	}

	return files, nil
}

// gatherCommandHistory reads recent shell commands
func (g *Gatherer) gatherCommandHistory(count int) ([]providers.CommandHistory, error) {
	shell := detectShell()
	return getRecentCommands(shell, count)
}

// gatherEnvironment collects relevant environment variables
func (g *Gatherer) gatherEnvironment() map[string]string {
	env := make(map[string]string)

	// Include relevant environment variables
	relevantVars := []string{
		"PATH", "HOME", "USER", "SHELL", "PWD",
		"NODE_ENV", "PYTHON_VERSION", "GOPATH", "GOROOT",
		"DOCKER_HOST", "KUBERNETES_NAMESPACE",
		"AWS_REGION", "AWS_PROFILE",
		"GIT_AUTHOR_NAME", "GIT_AUTHOR_EMAIL",
	}

	for _, varName := range relevantVars {
		if value := os.Getenv(varName); value != "" {
			env[varName] = value
		}
	}

	return env
}

// gatherGitContext collects git repository information
func (g *Gatherer) gatherGitContext() (*providers.GitContext, error) {
	return getGitContext()
}

// gatherProjectContext detects project type and configuration
func (g *Gatherer) gatherProjectContext() (*providers.ProjectContext, error) {
	return detectProjectType()
}

// shouldExclude checks if a file should be excluded from context
func (g *Gatherer) shouldExclude(name string) bool {
	// Default exclusions
	defaultExclusions := []string{
		".", "..", ".git", ".svn", ".hg",
		"node_modules", "__pycache__", ".pytest_cache",
		"target", "build", "dist", ".next",
		".DS_Store", "Thumbs.db",
	}

	for _, exclusion := range defaultExclusions {
		if name == exclusion {
			return true
		}
	}

	// User-defined exclusions
	for _, pattern := range g.config.ExcludePatterns {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}

	return false
}

// Helper functions
func getFileType(entry os.DirEntry) string {
	if entry.IsDir() {
		return "directory"
	}
	return "file"
}

func isImportantFile(name string) bool {
	importantFiles := map[string]bool{
		"README.md": true, "README.txt": true, "README.rst": true, "README": true,
		"package.json": true, "package-lock.json": true,
		"go.mod": true, "go.sum": true,
		"requirements.txt": true, "setup.py": true, "pyproject.toml": true,
		"Cargo.toml": true, "Cargo.lock": true,
		"Dockerfile": true, "docker-compose.yml": true, "docker-compose.yaml": true,
		"Makefile": true, "makefile": true,
		".gitignore": true, ".env": true, ".env.example": true,
		"tsconfig.json": true, "babel.config.js": true, "webpack.config.js": true,
		"pom.xml": true, "build.gradle": true,
	}

	return importantFiles[name]
}

func detectLanguage(filename string) string {
	ext := filepath.Ext(filename)

	languages := map[string]string{
		".go":   "go",
		".js":   "javascript",
		".ts":   "typescript",
		".py":   "python",
		".rs":   "rust",
		".java": "java",
		".cpp":  "cpp",
		".c":    "c",
		".h":    "c",
		".sh":   "bash",
		".bash": "bash",
		".zsh":  "zsh",
		".fish": "fish",
		".yaml": "yaml",
		".yml":  "yaml",
		".json": "json",
		".xml":  "xml",
		".html": "html",
		".css":  "css",
		".sql":  "sql",
		".md":   "markdown",
		".txt":  "text",
	}

	if lang, exists := languages[ext]; exists {
		return lang
	}

	return "unknown"
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}
	return filepath.Base(shell)
}
