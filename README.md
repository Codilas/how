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

```bash
# Setup development environment
make setup-dev

# Build and test
make build
make test

# Run in development mode
make dev
./bin/how "test prompt"
```

### Testing

The project includes comprehensive test coverage for configuration management:

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run tests with coverage report
make test-coverage

# Generate HTML coverage report
make test-coverage-html
```

For detailed testing documentation, see [TESTING.md](./TESTING.md).

## License

MIT License - see [LICENSE](LICENSE) file.
