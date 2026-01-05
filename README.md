# How - AI Shell Assistant

An intelligent shell assistant that helps you with commands, explanations, and code generation directly from your terminal.

## Quick Start

```bash
# Clone and build
git clone https://github.com/Codilas/how.git
cd how
make install

# Setup (configure AI provider)
how setup

# Start using
how "how to find large files?"
how write a Python function to validate emails
```

## Features

- **Context-aware AI assistance** - Understands your current directory, recent commands, and environment
- **Multiple input methods** - Natural prompts, command explanations, code generation
- **Multi-provider support** - Works with Claude, GPT-4, and local models
- **Conversation memory** - Remembers context across related questions
- **Command integration** - Suggests and optionally runs commands
- **Rich formatting** - Syntax highlighting, structured output

## Installation

### Quick Install
```bash
curl -sSL https://raw.githubusercontent.com/Codilas/how/main/scripts/install.sh | sh
```

### From Source
```bash
git clone https://github.com/Codilas/how.git
cd how
make install
```

## Configuration

Run the setup wizard:
```bash
how setup
```

## Usage Examples

```bash
# Ask questions
how explain kubernetes pods
how what is the difference between TCP and UDP

# Generate code
how write a bash script to backup MySQL
how create a Python function to parse JSON

# Get command help
how find files larger than 1GB
how what is best way to check memory usage

# Context-aware assistance
cd my-project/
how "how do I deploy this?"
```

## Development

### Prerequisites

- Go 1.22.4 or later
- `golangci-lint` for code linting (optional, auto-installed by `make setup-dev`)
- Git (for version tagging and commit information)

### Development Workflow

1. **Setup development environment:**
   ```bash
   make setup-dev
   ```
   This installs required development tools including golangci-lint.

2. **Make changes:**
   Edit code in the `cmd/`, `internal/`, and `pkg/` directories.

3. **Build and test:**
   ```bash
   make dev          # Build with race detection for development
   make test         # Run all tests
   make test-race    # Run tests with race detection
   ```

4. **Format and lint code:**
   ```bash
   make fmt          # Format all code
   make lint         # Run linter checks
   ```

5. **Run the application:**
   ```bash
   ./bin/how "test prompt"
   ```

## Make Targets

The project uses a Makefile with the following targets:

### Building

#### `make build`
**Description:** Build the application for the current platform.

Compiles the `how` binary with version information embedded via ldflags (version tag, git commit, and build date).

**Usage:**
```bash
make build
```

**Output:** `bin/how` binary

**Environment variables:**
- `VERSION` - Override the version string (default: git tag or "dev")

**Example:**
```bash
VERSION=1.2.3 make build
```

#### `make dev`
**Description:** Build the application with race detection enabled for development.

Useful for detecting race conditions during development. Same as `build` but with `-race` flag.

**Usage:**
```bash
make dev
```

**Output:** `bin/how` binary with race detection

#### `make release`
**Description:** Create cross-platform release builds.

Builds optimized binaries for multiple platforms:
- Linux (amd64)
- macOS (amd64, arm64)
- Windows (amd64)

**Usage:**
```bash
make release
```

**Output:** Release binaries in `releases/` directory
- `releases/how-linux-amd64`
- `releases/how-darwin-amd64`
- `releases/how-darwin-arm64`
- `releases/how-windows-amd64.exe`

### Installation

#### `make install`
**Description:** Build and install the application to `~/.local/bin/`.

This is the recommended installation method for development and testing. Requires no elevated permissions.

**Usage:**
```bash
make install
```

**Requirements:**
- `~/.local/bin` directory (created automatically if missing)
- `~/.local/bin` should be in your `PATH`

**Next steps after installation:**
```bash
# Add to PATH if not already present
export PATH="$HOME/.local/bin:$PATH"

# Configure the application
how setup
```

#### `make install-system`
**Description:** Build and install the application to `/usr/local/bin/`.

Requires `sudo` privileges. Use this for system-wide installation.

**Usage:**
```bash
make install-system
```

**Requirements:**
- sudo access
- `/usr/local/bin` must be in your PATH (typically is by default)

### Testing

#### `make test`
**Description:** Run all unit tests with verbose output.

Executes all tests in the project using `go test -v`.

**Usage:**
```bash
make test
```

**Output:** Test results with pass/fail status for each test

#### `make test-race`
**Description:** Run all tests with race detection enabled.

Detects race conditions in concurrent code. Useful before committing or before releases.

**Usage:**
```bash
make test-race
```

**Note:** Takes longer to run than regular tests due to instrumentation.

### Code Quality

#### `make fmt`
**Description:** Format all Go code according to standard conventions.

Runs `go fmt` on the entire codebase.

**Usage:**
```bash
make fmt
```

**Pre-commit tip:** Run this before committing to maintain consistent code style.

#### `make lint`
**Description:** Run static code analysis with `golangci-lint`.

Checks code for style, potential bugs, performance issues, and other problems.

**Usage:**
```bash
make lint
```

**Requirements:**
- `golangci-lint` must be installed (install with `make setup-dev`)

**Common lint errors:** Check the output and fix issues before committing.

### Dependency Management

#### `make deps`
**Description:** Download and tidy Go module dependencies.

Runs `go mod tidy` to clean up unused dependencies and `go mod download` to cache them.

**Usage:**
```bash
make deps
```

**When to use:**
- After pulling code changes
- After adding new imports
- When dependencies seem out of sync

### Maintenance

#### `make clean`
**Description:** Remove all build artifacts and clean build cache.

Deletes the `bin/` and `releases/` directories and runs `go clean`.

**Usage:**
```bash
make clean
```

**When to use:**
- Before fresh builds
- To free up disk space
- If builds seem stale or cached

#### `make setup-dev`
**Description:** Configure development environment.

Checks for and installs required development tools (specifically `golangci-lint`).

**Usage:**
```bash
make setup-dev
```

**What it does:**
1. Checks if `golangci-lint` is installed
2. If not found, installs it via `go install`
3. Prints confirmation message

**Run once when setting up development environment.**

## Complete Development Workflow Example

Here's a typical development session:

```bash
# Clone the repository
git clone https://github.com/Codilas/how.git
cd how

# Setup development environment (one-time)
make setup-dev

# Make changes to the code
# ...edit files...

# Format code
make fmt

# Run tests locally
make test

# Build for testing
make dev

# Test the binary
./bin/how "test query"

# Check code quality
make lint

# Clean up before committing
make clean

# Build for distribution
make release
```

## Build Information

The application embeds the following build metadata:

- **Version:** Git tag or "dev" (overridable via `VERSION` variable)
- **Git Commit:** Short SHA of the current commit
- **Build Date:** UTC timestamp of when binary was built

Access this information at runtime via the `pkg/version` package.

## License

MIT License - see [LICENSE](LICENSE) file.
