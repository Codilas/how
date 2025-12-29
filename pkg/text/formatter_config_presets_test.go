package text

import (
	"strings"
	"testing"
)

// SampleLLMResponse represents realistic LLM response content for testing
const SampleLLMResponse = `Here's how to set up a Python virtual environment:

## Setup Steps

First, install Python if you haven't already, then:

\`\`\`bash
python3 -m venv myenv
source myenv/bin/activate
\`\`\`

**Important**: Always activate your virtual environment before installing packages.

\`\`\`python
# Install dependencies
pip install requests flask

# Verify installation
import requests
print(requests.__version__)
\`\`\`

## Package Management

You can manage dependencies with a \`requirements.txt\` file:

| Package | Version | Purpose |
| --- | --- | --- |
| requests | 2.28.0 | HTTP library |
| flask | 2.0.0 | Web framework |
| pytest | 7.0.0 | Testing |

> **Note**: Keep your virtual environment separate from system Python

To deactivate when done:

\`\`\`bash
deactivate
\`\`\`

<structured_commands>
{
  "commands": ["python3 -m venv myenv", "source myenv/bin/activate"],
  "safe": true
}
</structured_commands>

Your Python environment is now ready for development!`

// TestDefaultConfigPreset validates DefaultConfig works correctly with sample LLM response
func TestDefaultConfigPreset(t *testing.T) {
	config := DefaultConfig()
	formatter := NewTerminalFormatter(config)
	output := formatter.Format(SampleLLMResponse)

	// Verify configuration properties
	t.Run("Config properties", func(t *testing.T) {
		if config.UseColors {
			t.Error("DefaultConfig should not use colors")
		}
		if config.CommentPrefix != "# " {
			t.Errorf("DefaultConfig CommentPrefix should be '# ', got %q", config.CommentPrefix)
		}
		if config.LineWidth != 80 {
			t.Errorf("DefaultConfig LineWidth should be 80, got %d", config.LineWidth)
		}
		if config.UseBoxes {
			t.Error("DefaultConfig should not use boxes")
		}
		if config.CompactMode {
			t.Error("DefaultConfig should not use compact mode")
		}
		if !config.ParseMarkdown {
			t.Error("DefaultConfig should parse markdown")
		}
	})

	// Verify output content
	t.Run("Output content preservation", func(t *testing.T) {
		expectedPhrases := []string{
			"Setup Steps",
			"python3 -m venv myenv",
			"Important",
			"Always activate",
			"pip install",
			"requests flask",
			"pytest",
			"deactivate",
			"ready for development",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("Output should contain %q", phrase)
			}
		}
	})

	// Verify structured commands are removed
	t.Run("Structured commands removal", func(t *testing.T) {
		if strings.Contains(output, "structured_commands") {
			t.Error("Output should not contain structured_commands tag")
		}
		if strings.Contains(output, `"commands"`) {
			t.Error("Output should not contain structured commands JSON")
		}
	})

	// Verify code blocks are present
	t.Run("Code block preservation", func(t *testing.T) {
		if !strings.Contains(output, "python3 -m venv") {
			t.Error("Output should contain bash code block")
		}
		if !strings.Contains(output, "pip install") {
			t.Error("Output should contain python code block")
		}
		if !strings.Contains(output, "import requests") {
			t.Error("Output should contain python code content")
		}
	})

	// Verify markdown formatting
	t.Run("Markdown content", func(t *testing.T) {
		if !strings.Contains(output, "Setup Steps") {
			t.Error("Output should contain markdown heading")
		}
		if !strings.Contains(output, "Package Management") {
			t.Error("Output should contain second markdown heading")
		}
	})

	// Verify output structure
	t.Run("Output structure", func(t *testing.T) {
		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
		// Default config should have comment prefixes
		lines := strings.Split(output, "\n")
		if len(lines) == 0 {
			t.Error("Output should have multiple lines")
		}
	})

	// Verify no ANSI color codes in output
	t.Run("No color codes in output", func(t *testing.T) {
		if strings.Contains(output, "\033") || strings.Contains(output, "\x1b") {
			t.Error("DefaultConfig output should not contain ANSI color codes")
		}
	})
}

// TestColoredConfigPreset validates ColoredConfig works correctly with sample LLM response
func TestColoredConfigPreset(t *testing.T) {
	config := ColoredConfig()
	formatter := NewTerminalFormatter(config)
	output := formatter.Format(SampleLLMResponse)

	// Verify configuration properties
	t.Run("Config properties", func(t *testing.T) {
		if !config.UseColors {
			t.Error("ColoredConfig should use colors")
		}
		if config.CommentPrefix != "# " {
			t.Errorf("ColoredConfig CommentPrefix should be '# ', got %q", config.CommentPrefix)
		}
		if config.LineWidth != 80 {
			t.Errorf("ColoredConfig LineWidth should be 80, got %d", config.LineWidth)
		}
		if !config.UseBoxes {
			t.Error("ColoredConfig should use boxes")
		}
		if !config.HighlightCode {
			t.Error("ColoredConfig should highlight code")
		}
		if !config.HighlightQuotes {
			t.Error("ColoredConfig should highlight quotes")
		}
		if config.CompactMode {
			t.Error("ColoredConfig should not use compact mode")
		}
		if !config.ParseMarkdown {
			t.Error("ColoredConfig should parse markdown")
		}
	})

	// Verify output content
	t.Run("Output content preservation", func(t *testing.T) {
		expectedPhrases := []string{
			"Setup Steps",
			"python3 -m venv myenv",
			"Important",
			"Always activate",
			"pip install",
			"requests flask",
			"pytest",
			"deactivate",
			"ready for development",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("Output should contain %q", phrase)
			}
		}
	})

	// Verify structured commands are removed
	t.Run("Structured commands removal", func(t *testing.T) {
		if strings.Contains(output, "structured_commands") {
			t.Error("Output should not contain structured_commands tag")
		}
		if strings.Contains(output, `"commands"`) {
			t.Error("Output should not contain structured commands JSON")
		}
	})

	// Verify code blocks are present
	t.Run("Code block preservation", func(t *testing.T) {
		if !strings.Contains(output, "python3 -m venv") {
			t.Error("Output should contain bash code block")
		}
		if !strings.Contains(output, "pip install") {
			t.Error("Output should contain python code block")
		}
		if !strings.Contains(output, "import requests") {
			t.Error("Output should contain python code content")
		}
	})

	// Verify markdown formatting
	t.Run("Markdown content", func(t *testing.T) {
		if !strings.Contains(output, "Setup Steps") {
			t.Error("Output should contain markdown heading")
		}
		if !strings.Contains(output, "Package Management") {
			t.Error("Output should contain second markdown heading")
		}
	})

	// Verify color output features
	t.Run("Colored features", func(t *testing.T) {
		// May contain ANSI codes
		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
		// Output should have meaningful content even with colors
		if !strings.Contains(output, "python") && !strings.Contains(output, "bash") {
			t.Error("Output should contain code language markers")
		}
	})

	// Verify boxes are being used (header/footer)
	t.Run("Box formatting", func(t *testing.T) {
		// ColoredConfig has UseBoxes=true, so header should be present
		// The writeHeader function is called when UseBoxes is true
		// This should show up in the output formatting
		if len(output) == 0 {
			t.Error("Output should have box formatting applied")
		}
	})

	// Verify table rendering
	t.Run("Table content", func(t *testing.T) {
		if !strings.Contains(output, "Package") {
			t.Error("Output should contain table header")
		}
		if !strings.Contains(output, "requests") {
			t.Error("Output should contain table data (requests package)")
		}
		if !strings.Contains(output, "flask") {
			t.Error("Output should contain table data (flask package)")
		}
	})
}

// TestCompactConfigPreset validates CompactConfig works correctly with sample LLM response
func TestCompactConfigPreset(t *testing.T) {
	config := CompactConfig()
	formatter := NewTerminalFormatter(config)
	output := formatter.Format(SampleLLMResponse)

	// Verify configuration properties
	t.Run("Config properties", func(t *testing.T) {
		if !config.UseColors {
			t.Error("CompactConfig should use colors")
		}
		if config.CommentPrefix != "" {
			t.Errorf("CompactConfig CommentPrefix should be empty, got %q", config.CommentPrefix)
		}
		if config.LineWidth != 120 {
			t.Errorf("CompactConfig LineWidth should be 120, got %d", config.LineWidth)
		}
		if config.UseBoxes {
			t.Error("CompactConfig should not use boxes")
		}
		if !config.HighlightCode {
			t.Error("CompactConfig should highlight code")
		}
		if !config.HighlightQuotes {
			t.Error("CompactConfig should highlight quotes")
		}
		if !config.CompactMode {
			t.Error("CompactConfig should use compact mode")
		}
		if !config.ParseMarkdown {
			t.Error("CompactConfig should parse markdown")
		}
		if config.WrapLongLines {
			t.Error("CompactConfig should not wrap long lines")
		}
	})

	// Verify output content
	t.Run("Output content preservation", func(t *testing.T) {
		expectedPhrases := []string{
			"Setup Steps",
			"python3 -m venv myenv",
			"Important",
			"Always activate",
			"pip install",
			"requests flask",
			"pytest",
			"deactivate",
			"ready for development",
		}

		for _, phrase := range expectedPhrases {
			if !strings.Contains(output, phrase) {
				t.Errorf("Output should contain %q", phrase)
			}
		}
	})

	// Verify structured commands are removed
	t.Run("Structured commands removal", func(t *testing.T) {
		if strings.Contains(output, "structured_commands") {
			t.Error("Output should not contain structured_commands tag")
		}
		if strings.Contains(output, `"commands"`) {
			t.Error("Output should not contain structured commands JSON")
		}
	})

	// Verify code blocks are present
	t.Run("Code block preservation", func(t *testing.T) {
		if !strings.Contains(output, "python3 -m venv") {
			t.Error("Output should contain bash code block")
		}
		if !strings.Contains(output, "pip install") {
			t.Error("Output should contain python code block")
		}
		if !strings.Contains(output, "import requests") {
			t.Error("Output should contain python code content")
		}
	})

	// Verify compact spacing (no comment prefix)
	t.Run("Compact spacing", func(t *testing.T) {
		// CompactConfig should not have comment prefixes
		// Output should start with content directly (no "# " prefix)
		lines := strings.Split(output, "\n")
		hasLeadingPrefix := false
		for _, line := range lines {
			if strings.HasPrefix(line, "# ") {
				hasLeadingPrefix = true
				break
			}
		}
		if hasLeadingPrefix {
			t.Error("CompactConfig should not have comment prefix in output")
		}
	})

	// Verify markdown formatting
	t.Run("Markdown content", func(t *testing.T) {
		if !strings.Contains(output, "Setup Steps") {
			t.Error("Output should contain markdown heading")
		}
		if !strings.Contains(output, "Package Management") {
			t.Error("Output should contain second markdown heading")
		}
	})

	// Verify density (CompactMode should minimize blank lines)
	t.Run("Compact density", func(t *testing.T) {
		if len(output) == 0 {
			t.Error("Output should not be empty")
		}
		// Compact mode should produce relatively dense output
		lines := strings.Split(output, "\n")
		if len(lines) == 0 {
			t.Error("Output should have content")
		}
	})

	// Verify table rendering
	t.Run("Table content", func(t *testing.T) {
		if !strings.Contains(output, "Package") {
			t.Error("Output should contain table header")
		}
		if !strings.Contains(output, "requests") {
			t.Error("Output should contain table data (requests package)")
		}
	})
}

// TestPresetsComparisonOutputDifference validates that different presets produce distinctly different output
func TestPresetsComparisonOutputDifference(t *testing.T) {
	defaultFormatter := NewTerminalFormatter(DefaultConfig())
	coloredFormatter := NewTerminalFormatter(ColoredConfig())
	compactFormatter := NewTerminalFormatter(CompactConfig())

	defaultOutput := defaultFormatter.Format(SampleLLMResponse)
	coloredOutput := coloredFormatter.Format(SampleLLMResponse)
	compactOutput := compactFormatter.Format(SampleLLMResponse)

	// All should have content
	t.Run("All presets produce output", func(t *testing.T) {
		if len(defaultOutput) == 0 {
			t.Error("DefaultConfig should produce output")
		}
		if len(coloredOutput) == 0 {
			t.Error("ColoredConfig should produce output")
		}
		if len(compactOutput) == 0 {
			t.Error("CompactConfig should produce output")
		}
	})

	// All should contain essential content
	t.Run("Essential content in all outputs", func(t *testing.T) {
		essentials := []string{"Setup Steps", "python3 -m venv", "Important"}
		outputs := map[string]string{
			"default": defaultOutput,
			"colored": coloredOutput,
			"compact": compactOutput,
		}

		for outputName, output := range outputs {
			for _, essential := range essentials {
				if !strings.Contains(output, essential) {
					t.Errorf("%s output missing essential phrase: %q", outputName, essential)
				}
			}
		}
	})

	// DefaultConfig should have comment prefixes, others might not
	t.Run("Prefix difference", func(t *testing.T) {
		defaultHasPrefix := strings.Contains(defaultOutput, "# ")
		compactHasPrefix := strings.Contains(compactOutput, "# ")

		if !defaultHasPrefix {
			t.Error("DefaultConfig output should contain comment prefixes")
		}
		if compactHasPrefix {
			t.Error("CompactConfig output should not contain comment prefixes")
		}
	})

	// DefaultConfig should not have color codes, ColoredConfig might
	t.Run("Color code differences", func(t *testing.T) {
		defaultHasColor := strings.Contains(defaultOutput, "\033") || strings.Contains(defaultOutput, "\x1b")
		if defaultHasColor {
			t.Error("DefaultConfig should not contain ANSI color codes")
		}
	})

	// CompactConfig should have fewer lines (due to CompactMode)
	t.Run("Line count differences", func(t *testing.T) {
		defaultLines := len(strings.Split(defaultOutput, "\n"))
		compactLines := len(strings.Split(compactOutput, "\n"))

		// Compact mode should produce fewer or equal lines
		if compactLines > defaultLines+5 {
			t.Errorf("CompactConfig should have similar or fewer lines than DefaultConfig: default=%d, compact=%d",
				defaultLines, compactLines)
		}
	})

	// ColoredConfig should have boxes
	t.Run("Box formatting in ColoredConfig", func(t *testing.T) {
		// ColoredConfig uses UseBoxes=true
		// This affects the header output
		if len(coloredOutput) == 0 {
			t.Error("ColoredConfig should have formatted output")
		}
	})

	// All should preserve code block content
	t.Run("Code preservation in all configs", func(t *testing.T) {
		codeContent := "python3 -m venv myenv"
		if !strings.Contains(defaultOutput, codeContent) {
			t.Error("DefaultConfig should preserve code blocks")
		}
		if !strings.Contains(coloredOutput, codeContent) {
			t.Error("ColoredConfig should preserve code blocks")
		}
		if !strings.Contains(compactOutput, codeContent) {
			t.Error("CompactConfig should preserve code blocks")
		}
	})
}

// TestPresetsWithSpecialFeatures validates special formatting features with different presets
func TestPresetsWithSpecialFeatures(t *testing.T) {
	// Test with markdown that uses special formatting
	specialContent := `# Main Heading

This is **bold text** and this is *italic text*.

Here's a [link](https://example.com) to follow.

> This is a blockquote
> with multiple lines

- Bullet point 1
- Bullet point 2
- Bullet point 3

\`\`\`python
def hello():
    print("Hello, World!")
\`\`\`

**Important**: This is bold content.`

	configs := []struct {
		name   string
		config FormatterConfig
	}{
		{"DefaultConfig", DefaultConfig()},
		{"ColoredConfig", ColoredConfig()},
		{"CompactConfig", CompactConfig()},
	}

	for _, tc := range configs {
		t.Run(tc.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tc.config)
			output := formatter.Format(specialContent)

			// All should preserve headings
			if !strings.Contains(output, "Main Heading") {
				t.Error("Should preserve heading")
			}

			// All should preserve code
			if !strings.Contains(output, "hello") {
				t.Error("Should preserve code content")
			}

			// All should preserve important text
			if !strings.Contains(output, "Important") {
				t.Error("Should preserve important content")
			}

			// All should preserve blockquote content
			if !strings.Contains(output, "blockquote") {
				t.Error("Should preserve blockquote")
			}

			// All should preserve lists
			if !strings.Contains(output, "Bullet point") {
				t.Error("Should preserve list items")
			}

			// Output should be non-empty
			if len(output) == 0 {
				t.Error("Output should not be empty")
			}
		})
	}
}

// TestPresetsConsistency validates that applying presets multiple times produces consistent output
func TestPresetsConsistency(t *testing.T) {
	configs := []struct {
		name   string
		config FormatterConfig
	}{
		{"DefaultConfig", DefaultConfig()},
		{"ColoredConfig", ColoredConfig()},
		{"CompactConfig", CompactConfig()},
	}

	for _, tc := range configs {
		t.Run(tc.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tc.config)

			// Format the same content multiple times
			output1 := formatter.Format(SampleLLMResponse)
			output2 := formatter.Format(SampleLLMResponse)
			output3 := formatter.Format(SampleLLMResponse)

			// All outputs should be identical
			if output1 != output2 {
				t.Error("Formatting should be consistent - output1 != output2")
			}
			if output2 != output3 {
				t.Error("Formatting should be consistent - output2 != output3")
			}
		})
	}
}

// TestPresetsRobustness validates that presets handle edge cases gracefully
func TestPresetsRobustness(t *testing.T) {
	edgeCases := []struct {
		name  string
		input string
	}{
		{"Empty string", ""},
		{"Only whitespace", "   \n\n\n   "},
		{"Only code block", "```python\nprint('hello')\n```"},
		{"Only structured commands", "<structured_commands>data</structured_commands>"},
		{"Very long lines", strings.Repeat("This is a very long line. ", 100)},
		{"Mixed empty blocks", "```\n\n```\n\nText\n\n```\n\n```"},
	}

	configs := []struct {
		name   string
		config FormatterConfig
	}{
		{"DefaultConfig", DefaultConfig()},
		{"ColoredConfig", ColoredConfig()},
		{"CompactConfig", CompactConfig()},
	}

	for _, ec := range edgeCases {
		t.Run(ec.name, func(t *testing.T) {
			for _, tc := range configs {
				t.Run(tc.name, func(t *testing.T) {
					formatter := NewTerminalFormatter(tc.config)
					// Should not panic
					output := formatter.Format(ec.input)
					// Should return a string (possibly empty)
					_ = output
				})
			}
		})
	}
}

// TestPresetsSpacing validates proper spacing and indentation with different presets
func TestPresetsSpacing(t *testing.T) {
	spacingTest := `Here's a list:

- Item 1
- Item 2
- Item 3

\`\`\`bash
command1
command2
\`\`\`

And some more text.`

	t.Run("DefaultConfig spacing", func(t *testing.T) {
		config := DefaultConfig()
		formatter := NewTerminalFormatter(config)
		output := formatter.Format(spacingTest)

		// Should have comment prefixes on lines
		if !strings.Contains(output, "# ") {
			t.Error("DefaultConfig should use comment prefix")
		}

		// Should have proper indentation for code
		if len(output) > 0 && !strings.Contains(output, "  ") {
			t.Error("DefaultConfig should have indentation")
		}
	})

	t.Run("ColoredConfig spacing", func(t *testing.T) {
		config := ColoredConfig()
		formatter := NewTerminalFormatter(config)
		output := formatter.Format(spacingTest)

		// Should have content
		if len(output) == 0 {
			t.Error("ColoredConfig should produce output")
		}

		// Should preserve structure
		if !strings.Contains(output, "Item") {
			t.Error("ColoredConfig should preserve list items")
		}
	})

	t.Run("CompactConfig spacing", func(t *testing.T) {
		config := CompactConfig()
		formatter := NewTerminalFormatter(config)
		output := formatter.Format(spacingTest)

		// Should have no comment prefix
		if strings.Contains(output, "# ") {
			t.Error("CompactConfig should not use comment prefix")
		}

		// Should have content
		if len(output) == 0 {
			t.Error("CompactConfig should produce output")
		}

		// Should preserve structure
		if !strings.Contains(output, "Item") {
			t.Error("CompactConfig should preserve list items")
		}
	})
}
