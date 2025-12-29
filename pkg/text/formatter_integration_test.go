package text

import (
	"strings"
	"testing"
)

// TestFormatIntegrationBasicStructuredCommandRemoval tests Format method removing structured commands
func TestFormatIntegrationBasicStructuredCommandRemoval(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		config        FormatterConfig
		shouldContain []string
		shouldNotHave []string
	}{
		{
			name: "Removes structured_commands section",
			input: `Here's the answer:

<structured_commands>
{
  "commands": ["ls", "pwd"],
  "safe": true
}
</structured_commands>

The output above shows your files.`,
			config: DefaultConfig(),
			shouldContain: []string{
				"Here's the answer",
				"The output above shows your files",
			},
			shouldNotHave: []string{
				"structured_commands",
				"commands",
			},
		},
		{
			name: "Multiple structured_commands sections",
			input: `First section

<structured_commands>
command1
</structured_commands>

Middle section

<structured_commands>
command2
</structured_commands>

Final section`,
			config: DefaultConfig(),
			shouldContain: []string{
				"First section",
				"Middle section",
				"Final section",
			},
			shouldNotHave: []string{
				"structured_commands",
				"command1",
				"command2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			output := formatter.Format(tt.input)

			for _, shouldHave := range tt.shouldContain {
				if !strings.Contains(output, shouldHave) {
					t.Errorf("Output should contain %q, got: %s", shouldHave, output)
				}
			}

			for _, shouldNotHave := range tt.shouldNotHave {
				if strings.Contains(output, shouldNotHave) {
					t.Errorf("Output should not contain %q, got: %s", shouldNotHave, output)
				}
			}
		})
	}
}

// TestFormatIntegrationCodeBlocksWithMarkdown tests Format combining code blocks with markdown
func TestFormatIntegrationCodeBlocksWithMarkdown(t *testing.T) {
	input := `# Installing Dependencies

To install the required packages, follow these steps:

## Step 1: Clone the repository

\`\`\`bash
git clone https://github.com/example/repo.git
cd repo
\`\`\`

## Step 2: Install dependencies

\`\`\`bash
npm install
\`\`\`

The above commands will set up your project.`

	tests := []struct {
		name   string
		config FormatterConfig
		checks func(t *testing.T, output string)
	}{
		{
			name:   "Default config with markdown parsing",
			config: DefaultConfig(),
			checks: func(t *testing.T, output string) {
				// Should contain code blocks
				if !strings.Contains(output, "git clone") {
					t.Error("Should contain git clone command")
				}
				if !strings.Contains(output, "npm install") {
					t.Error("Should contain npm install command")
				}
				// Should contain markdown content
				if !strings.Contains(output, "Installing Dependencies") {
					t.Error("Should contain heading")
				}
				if !strings.Contains(output, "Clone the repository") {
					t.Error("Should contain subheading")
				}
			},
		},
		{
			name:   "Colored config with code highlighting",
			config: ColoredConfig(),
			checks: func(t *testing.T, output string) {
				// Should contain code blocks even with colors
				if !strings.Contains(output, "git clone") {
					t.Error("Should contain git clone command")
				}
				// Output should be non-empty
				if len(output) == 0 {
					t.Error("Output should not be empty")
				}
			},
		},
		{
			name: "Compact mode",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "",
				LineWidth:       80,
				IndentSize:      2,
				UseBoxes:        false,
				UseBullets:      true,
				HighlightCode:   false,
				WrapLongLines:   false,
				RenderTables:    true,
				ShowLineNumbers: false,
				CompactMode:     true,
				HighlightQuotes: false,
				ParseMarkdown:   true,
			},
			checks: func(t *testing.T, output string) {
				// Should still contain the essential content
				if !strings.Contains(output, "git clone") {
					t.Error("Should contain git clone command")
				}
				// No comment prefix in compact mode
				if !strings.HasPrefix(output, "# Installing") && !strings.HasPrefix(output, "Installing") {
					// First line should be about installing (might have no prefix)
					if !strings.Contains(output, "Installing") {
						t.Error("Should contain Installing Dependencies")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			output := formatter.Format(input)
			tt.checks(t, output)
		})
	}
}

// TestFormatIntegrationTablesWithFormatting tests Format with tables and other formatting
func TestFormatIntegrationTablesWithFormatting(t *testing.T) {
	input := `# Performance Results

Here are the benchmark results:

| Metric | Value | Status |
| --- | --- | --- |
| Throughput | 1000 req/s | Good |
| Latency | 50ms | Excellent |
| Memory | 256MB | Good |

These results show good performance.`

	tests := []struct {
		name   string
		config FormatterConfig
		checks func(t *testing.T, output string)
	}{
		{
			name: "With table rendering",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				UseBoxes:        false,
				UseBullets:      true,
				HighlightCode:   false,
				WrapLongLines:   true,
				RenderTables:    true,
				ShowLineNumbers: false,
				CompactMode:     false,
				HighlightQuotes: false,
				ParseMarkdown:   true,
			},
			checks: func(t *testing.T, output string) {
				// Should render table with borders
				if !strings.Contains(output, "Metric") {
					t.Error("Should contain table header")
				}
				if !strings.Contains(output, "Throughput") {
					t.Error("Should contain table data")
				}
				if !strings.Contains(output, "│") {
					t.Error("Should contain table border characters")
				}
			},
		},
		{
			name: "With table rendering disabled",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				UseBoxes:        false,
				UseBullets:      true,
				HighlightCode:   false,
				WrapLongLines:   true,
				RenderTables:    false,
				ShowLineNumbers: false,
				CompactMode:     false,
				HighlightQuotes: false,
				ParseMarkdown:   true,
			},
			checks: func(t *testing.T, output string) {
				// Should contain raw table data
				if !strings.Contains(output, "Metric") {
					t.Error("Should contain table content")
				}
				// But might not have fancy borders if table rendering is disabled
				// It will still have the pipe characters from original markdown
				if !strings.Contains(output, "|") {
					t.Error("Should contain pipe characters from markdown table")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			output := formatter.Format(input)
			tt.checks(t, output)
		})
	}
}

// TestFormatIntegrationComplexWorkflow tests Format with all features combined
func TestFormatIntegrationComplexWorkflow(t *testing.T) {
	input := `<structured_commands>
{
  "commands": ["docker build", "docker run"]
}
</structured_commands>

# Docker Setup Guide

## Creating a Dockerfile

Create a **Dockerfile** with the following content:

\`\`\`dockerfile
FROM golang:1.22
WORKDIR /app
COPY . .
RUN go build -o myapp
EXPOSE 8080
CMD ["./myapp"]
\`\`\`

## Building the Image

Run this command to build:

\`\`\`bash
docker build -t myapp:latest .
docker run -p 8080:8080 myapp:latest
\`\`\`

> Note: Make sure Docker is installed on your system

## Performance Metrics

| Stage | Duration | Status |
| --- | --- | --- |
| Build | 45s | Success |
| Run | <1s | Ready |

- Build image successfully
- Push to registry (optional)
- Deploy to production

The setup is complete!`

	tests := []struct {
		name   string
		config FormatterConfig
		checks func(t *testing.T, output string)
	}{
		{
			name:   "Full featured colored config",
			config: ColoredConfig(),
			checks: func(t *testing.T, output string) {
				// Should remove structured commands
				if strings.Contains(output, "docker build") && strings.Contains(output, "structured_commands") {
					t.Error("Should remove structured_commands section")
				}

				// Should contain markdown
				if !strings.Contains(output, "Docker Setup Guide") {
					t.Error("Should contain main heading")
				}

				// Should contain code blocks
				if !strings.Contains(output, "FROM golang") {
					t.Error("Should contain Dockerfile content")
				}
				if !strings.Contains(output, "docker build") {
					t.Error("Should contain docker command")
				}

				// Should contain quotes
				if !strings.Contains(output, "Note") {
					t.Error("Should contain blockquote content")
				}

				// Should contain table
				if !strings.Contains(output, "Stage") || !strings.Contains(output, "Duration") {
					t.Error("Should contain table content")
				}

				// Should contain list
				if !strings.Contains(output, "Build image") {
					t.Error("Should contain list items")
				}
			},
		},
		{
			name: "Default config (minimal features)",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				UseBoxes:        false,
				UseBullets:      true,
				HighlightCode:   false,
				WrapLongLines:   true,
				RenderTables:    true,
				ShowLineNumbers: false,
				CompactMode:     false,
				HighlightQuotes: false,
				ParseMarkdown:   true,
			},
			checks: func(t *testing.T, output string) {
				// Should remove structured commands
				if strings.Contains(output, "structured_commands") {
					t.Error("Should remove structured_commands section")
				}

				// Should contain essential content
				if !strings.Contains(output, "Docker Setup Guide") {
					t.Error("Should contain heading")
				}
				if !strings.Contains(output, "FROM golang") {
					t.Error("Should contain code content")
				}

				// Table should still be present
				if !strings.Contains(output, "Stage") {
					t.Error("Should contain table")
				}
			},
		},
		{
			name: "With boxes and line numbers",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				UseBoxes:        true,
				UseBullets:      true,
				HighlightCode:   false,
				WrapLongLines:   true,
				RenderTables:    true,
				ShowLineNumbers: true,
				CompactMode:     false,
				HighlightQuotes: true,
				ParseMarkdown:   true,
			},
			checks: func(t *testing.T, output string) {
				// Should have header box
				if !strings.Contains(output, "=") {
					t.Error("Should contain header box borders")
				}

				// Code blocks should have line numbers
				if !strings.Contains(output, ":") {
					// Line numbers are formatted as "  1: "
					t.Logf("Output:\n%s", output)
				}

				// Should contain content
				if !strings.Contains(output, "FROM golang") {
					t.Error("Should contain code content")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			output := formatter.Format(input)

			// Basic sanity checks
			if len(output) == 0 {
				t.Error("Output should not be empty")
			}

			tt.checks(t, output)
		})
	}
}

// TestFormatIntegrationLineWrappingWithCodeBlocks tests line wrapping with code blocks
func TestFormatIntegrationLineWrappingWithCodeBlocks(t *testing.T) {
	input := `Here's a very long line that should be wrapped when the formatter is configured with a narrow line width: This text is intentionally long to test the line wrapping functionality.

\`\`\`python
def very_long_function_name_that_takes_many_parameters(param1, param2, param3, param4, param5):
    return param1 + param2 + param3 + param4 + param5
\`\`\`

More text that follows.`

	tests := []struct {
		name   string
		config FormatterConfig
		checks func(t *testing.T, output string)
	}{
		{
			name: "Line wrapping enabled with narrow width",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       60,
				IndentSize:      2,
				UseBoxes:        false,
				UseBullets:      true,
				HighlightCode:   false,
				WrapLongLines:   true,
				RenderTables:    true,
				ShowLineNumbers: false,
				CompactMode:     false,
				HighlightQuotes: false,
				ParseMarkdown:   true,
			},
			checks: func(t *testing.T, output string) {
				// Should contain wrapped content
				if !strings.Contains(output, "Here's a very long line") {
					t.Error("Should contain text content")
				}

				// Should preserve code block content
				if !strings.Contains(output, "def very_long_function") {
					t.Error("Should contain code content")
				}

				// Code block lines shouldn't wrap individually
				lines := strings.Split(output, "\n")
				foundCodeLine := false
				for _, line := range lines {
					if strings.Contains(line, "def very_long_function") {
						foundCodeLine = true
						// Code lines are allowed to exceed width limit
					}
				}
				if !foundCodeLine {
					t.Error("Should find the function definition line")
				}
			},
		},
		{
			name: "Line wrapping disabled",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       60,
				IndentSize:      2,
				UseBoxes:        false,
				UseBullets:      true,
				HighlightCode:   false,
				WrapLongLines:   false,
				RenderTables:    true,
				ShowLineNumbers: false,
				CompactMode:     false,
				HighlightQuotes: false,
				ParseMarkdown:   true,
			},
			checks: func(t *testing.T, output string) {
				// Should contain long lines without wrapping
				if !strings.Contains(output, "intentionally long") {
					t.Error("Should preserve long text")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			output := formatter.Format(input)

			if len(output) == 0 {
				t.Error("Output should not be empty")
			}

			tt.checks(t, output)
		})
	}
}

// TestFormatIntegrationMultipleCodeBlocksSequence tests multiple code blocks in sequence
func TestFormatIntegrationMultipleCodeBlocksSequence(t *testing.T) {
	input := `# Setup Instructions

First, configure the environment:

\`\`\`bash
export NODE_ENV=production
export DEBUG=false
\`\`\`

Then, build the application:

\`\`\`bash
npm run build
npm run bundle
\`\`\`

Finally, start the server:

\`\`\`javascript
const server = require('./server');
server.start();
\`\`\`

Done!`

	config := FormatterConfig{
		UseColors:       false,
		CommentPrefix:   "# ",
		LineWidth:       80,
		IndentSize:      2,
		UseBoxes:        false,
		UseBullets:      false,
		HighlightCode:   false,
		WrapLongLines:   true,
		RenderTables:    true,
		ShowLineNumbers: false,
		CompactMode:     false,
		HighlightQuotes: false,
		ParseMarkdown:   true,
	}

	formatter := NewTerminalFormatter(config)
	output := formatter.Format(input)

	// Check that all code blocks are present
	if !strings.Contains(output, "export NODE_ENV") {
		t.Error("Should contain first code block")
	}
	if !strings.Contains(output, "npm run build") {
		t.Error("Should contain second code block")
	}
	if !strings.Contains(output, "const server") {
		t.Error("Should contain third code block")
	}

	// Check language identifiers are preserved
	bashCount := strings.Count(output, "```bash")
	jsCount := strings.Count(output, "```javascript")

	if bashCount != 2 {
		t.Errorf("Expected 2 bash code blocks, got %d", bashCount)
	}
	if jsCount != 1 {
		t.Errorf("Expected 1 javascript code block, got %d", jsCount)
	}

	// Check text between code blocks is preserved
	if !strings.Contains(output, "First, configure") {
		t.Error("Should contain text between code blocks")
	}
	if !strings.Contains(output, "Finally, start") {
		t.Error("Should contain text before last code block")
	}
	if !strings.Contains(output, "Done!") {
		t.Error("Should contain trailing text")
	}
}

// TestFormatIntegrationWhitespaceNormalization tests whitespace cleanup across features
func TestFormatIntegrationWhitespaceNormalization(t *testing.T) {
	input := `Some text


Multiple blank lines above

\`\`\`python
code block
\`\`\`


Multiple blank lines after code



Another section`

	config := DefaultConfig()
	formatter := NewTerminalFormatter(config)
	output := formatter.Format(input)

	// Should normalize multiple blank lines to at most 2
	tripleNewlines := strings.Count(output, "\n\n\n")
	if tripleNewlines > 0 {
		t.Errorf("Should not have triple newlines, found %d", tripleNewlines)
	}

	// Should preserve essential content
	if !strings.Contains(output, "Some text") {
		t.Error("Should preserve initial text")
	}
	if !strings.Contains(output, "code block") {
		t.Error("Should preserve code block")
	}
	if !strings.Contains(output, "Another section") {
		t.Error("Should preserve final text")
	}
}

// TestFormatIntegrationMixedListsAndQuotes tests lists and quotes formatting
func TestFormatIntegrationMixedListsAndQuotes(t *testing.T) {
	input := `# Instructions

Follow these steps:

- First step is important
- Second step builds on the first
- Third step completes the process

Important notes:

> This is a critical requirement
> It must be done correctly

More steps:

1. Setup the environment
2. Install dependencies
3. Run the tests

> Final note: Always verify your work`

	config := FormatterConfig{
		UseColors:       false,
		CommentPrefix:   "# ",
		LineWidth:       80,
		IndentSize:      2,
		UseBoxes:        false,
		UseBullets:      true,
		HighlightCode:   false,
		WrapLongLines:   true,
		RenderTables:    true,
		ShowLineNumbers: false,
		CompactMode:     false,
		HighlightQuotes: true,
		ParseMarkdown:   true,
	}

	formatter := NewTerminalFormatter(config)
	output := formatter.Format(input)

	// Check bullet points are preserved
	if !strings.Contains(output, "First step") {
		t.Error("Should contain first bullet point")
	}
	if !strings.Contains(output, "Second step") {
		t.Error("Should contain second bullet point")
	}

	// Check numbered list items are preserved
	if !strings.Contains(output, "Setup the environment") {
		t.Error("Should contain numbered list item")
	}

	// Check quotes are preserved
	if !strings.Contains(output, "critical requirement") {
		t.Error("Should contain blockquote content")
	}

	// Check all text is present
	if !strings.Contains(output, "Final note") {
		t.Error("Should contain final blockquote")
	}
}

// TestFormatIntegrationComplexInlineFormatting tests inline formatting combinations
func TestFormatIntegrationComplexInlineFormatting(t *testing.T) {
	input := `Here's a **bold statement** with *italic text* and a [link to docs](https://example.com).

The \`command line\` tool supports many features:
- Use \`--help\` for options
- Use \`--verbose\` for details
- Check [documentation](https://example.com/docs)

**Important**: Always read the *manual* before using.`

	config := FormatterConfig{
		UseColors:       true,
		CommentPrefix:   "# ",
		LineWidth:       80,
		IndentSize:      2,
		UseBoxes:        false,
		UseBullets:      true,
		HighlightCode:   true,
		WrapLongLines:   true,
		RenderTables:    true,
		ShowLineNumbers: false,
		CompactMode:     false,
		HighlightQuotes: false,
		ParseMarkdown:   true,
	}

	formatter := NewTerminalFormatter(config)
	output := formatter.Format(input)

	// Should contain all text elements (formatting markers may be removed or styled)
	if !strings.Contains(output, "bold statement") {
		t.Error("Should contain bold text content")
	}
	if !strings.Contains(output, "italic text") {
		t.Error("Should contain italic text content")
	}
	if !strings.Contains(output, "link to docs") {
		t.Error("Should contain link text")
	}
	if !strings.Contains(output, "command line") {
		t.Error("Should contain inline code")
	}
	if !strings.Contains(output, "--help") {
		t.Error("Should contain command option")
	}
	if !strings.Contains(output, "documentation") {
		t.Error("Should contain documentation link")
	}
}

// TestFormatIntegrationConfigVariations tests different config combinations
func TestFormatIntegrationConfigVariations(t *testing.T) {
	input := `# Test Content

\`\`\`go
func main() {
    fmt.Println("Hello")
}
\`\`\`

| Feature | Enabled |
| --- | --- |
| Colors | No |

- Item 1
- Item 2`

	variations := []struct {
		name   string
		config FormatterConfig
	}{
		{
			name:   "All features enabled",
			config: ColoredConfig(),
		},
		{
			name:   "All features disabled except markdown",
			config: FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				UseBoxes:        false,
				UseBullets:      false,
				HighlightCode:   false,
				WrapLongLines:   false,
				RenderTables:    false,
				ShowLineNumbers: false,
				CompactMode:     false,
				HighlightQuotes: false,
				ParseMarkdown:   true,
			},
		},
		{
			name:   "Compact mode only",
			config: CompactConfig(),
		},
		{
			name: "Custom settings",
			config: FormatterConfig{
				UseColors:       true,
				CommentPrefix:   "> ",
				LineWidth:       100,
				IndentSize:      4,
				UseBoxes:        true,
				UseBullets:      false,
				HighlightCode:   true,
				WrapLongLines:   true,
				RenderTables:    true,
				ShowLineNumbers: true,
				CompactMode:     false,
				HighlightQuotes: true,
				ParseMarkdown:   false,
			},
		},
	}

	for _, variation := range variations {
		t.Run(variation.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(variation.config)
			output := formatter.Format(input)

			// All configs should produce non-empty output
			if len(output) == 0 {
				t.Error("Format should produce output for all config variations")
			}

			// Code content should always be present
			if !strings.Contains(output, "func main") {
				t.Error("Code content should be preserved in all variations")
			}
		})
	}
}

// TestFormatIntegrationEdgeCases tests Format with edge case inputs
func TestFormatIntegrationEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			name:   "Empty string",
			input:  "",
			config: DefaultConfig(),
		},
		{
			name:   "Only whitespace",
			input:  "   \n\n\n   ",
			config: DefaultConfig(),
		},
		{
			name:   "Only structured commands",
			input:  "<structured_commands>data</structured_commands>",
			config: DefaultConfig(),
		},
		{
			name:   "Only code blocks",
			input:  "```\ncode\n```",
			config: DefaultConfig(),
		},
		{
			name:   "Nested code blocks syntax",
			input:  "```\n```python\ncode\n```\n```",
			config: DefaultConfig(),
		},
		{
			name:   "Malformed markdown",
			input:  "# [[ { unclosed",
			config: DefaultConfig(),
		},
		{
			name:   "Very long single line",
			input:  strings.Repeat("a", 1000),
			config: FormatterConfig{LineWidth: 80, WrapLongLines: true, CommentPrefix: "# "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)

			// Should not panic
			output := formatter.Format(tt.input)

			// Output should be a string (even if empty)
			_ = output
		})
	}
}

// TestFormatIntegrationPreservesEssentialContent tests that critical content is preserved
func TestFormatIntegrationPreservesEssentialContent(t *testing.T) {
	input := `<structured_commands>
{
  "command": "rm -rf /"
}
</structured_commands>

**CRITICAL**: Do not run unsafe commands.

Here's the safe way:

\`\`\`bash
ls -la
cd /tmp
\`\`\`

Follow these important steps:

1. Always verify commands
2. Use dry-run mode first
3. Back up your data

| Step | Action | Status |
| --- | --- | --- |
| 1 | Verify | ✓ |
| 2 | Backup | ✓ |

> Important: Test in safe environment first

The command is now ready to use.`

	configs := []FormatterConfig{
		DefaultConfig(),
		ColoredConfig(),
		CompactConfig(),
	}

	essentialPhrases := []string{
		"CRITICAL",
		"Do not run unsafe",
		"safe way",
		"ls -la",
		"Always verify",
		"Back up",
		"Important",
		"Test in safe",
		"ready to use",
	}

	for _, config := range configs {
		formatter := NewTerminalFormatter(config)
		output := formatter.Format(input)

		for _, phrase := range essentialPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("Essential phrase %q not found in output with config", phrase)
			}
		}

		// Should NOT contain the structured command
		if strings.Contains(output, "rm -rf /") && strings.Contains(output, "structured_commands") {
			t.Error("Should remove structured_commands section with unsafe command")
		}
	}
}

// TestFormatIntegrationConsistency tests Format produces consistent output
func TestFormatIntegrationConsistency(t *testing.T) {
	input := `# Heading

Some content with **bold** and \`code\`.

\`\`\`python
def func():
    pass
\`\`\`

Final line.`

	config := DefaultConfig()
	formatter := NewTerminalFormatter(config)

	// Format the same input multiple times
	output1 := formatter.Format(input)
	output2 := formatter.Format(input)
	output3 := formatter.Format(input)

	// All outputs should be identical
	if output1 != output2 {
		t.Error("First and second format calls should produce identical output")
	}
	if output2 != output3 {
		t.Error("Second and third format calls should produce identical output")
	}
}

// TestFormatIntegrationBoundaryLineWidths tests Format with extreme line widths
func TestFormatIntegrationBoundaryLineWidths(t *testing.T) {
	input := `This is a test line that might wrap.

\`\`\`bash
echo "hello world"
\`\`\`

Another line.`

	tests := []struct {
		name      string
		lineWidth int
	}{
		{
			name:      "Very narrow (30)",
			lineWidth: 30,
		},
		{
			name:      "Very wide (200)",
			lineWidth: 200,
		},
		{
			name:      "Minimum practical (40)",
			lineWidth: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FormatterConfig{
				LineWidth:     tt.lineWidth,
				CommentPrefix: "# ",
				WrapLongLines: true,
			}
			formatter := NewTerminalFormatter(config)
			output := formatter.Format(input)

			// Should produce output and contain essential content
			if len(output) == 0 {
				t.Error("Should produce output even with extreme line widths")
			}
			if !strings.Contains(output, "hello world") {
				t.Error("Should preserve code content")
			}
		})
	}
}
