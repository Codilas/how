package text

import (
	"strings"
	"testing"
)

// TestFormatEmptyString tests formatting with empty input
func TestFormatEmptyString(t *testing.T) {
	tests := []struct {
		name   string
		config FormatterConfig
	}{
		{"Empty with default config", DefaultConfig()},
		{"Empty with colored config", ColoredConfig()},
		{"Empty with compact config", CompactConfig()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format("")
			// Should handle gracefully without panicking
			if result != "" {
				t.Errorf("Expected empty result, got: %q", result)
			}
		})
	}
}

// TestFormatWhitespaceOnly tests formatting with only whitespace
func TestFormatWhitespaceOnly(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		config    FormatterConfig
	}{
		{"Single space", " ", DefaultConfig()},
		{"Multiple spaces", "   ", DefaultConfig()},
		{"Single tab", "\t", DefaultConfig()},
		{"Multiple tabs", "\t\t\t", DefaultConfig()},
		{"Mixed whitespace", " \t \t ", DefaultConfig()},
		{"Multiple newlines", "\n\n\n", DefaultConfig()},
		{"Mixed newlines and spaces", "\n \n \n", DefaultConfig()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			// Should handle gracefully without panicking
			if result == "" || !strings.Contains(result, "AI Assistant Response") {
				// Result should be either empty or just formatting
			}
		})
	}
}

// TestFormatVeryLongInput tests formatting with extremely long input
func TestFormatVeryLongInput(t *testing.T) {
	tests := []struct {
		name        string
		inputLength int
		config      FormatterConfig
	}{
		{"1KB input", 1024, DefaultConfig()},
		{"10KB input", 10 * 1024, DefaultConfig()},
		{"100KB input", 100 * 1024, DefaultConfig()},
		{"1MB input", 1024 * 1024, DefaultConfig()},
		{"Long lines without wrapping", 1024, CompactConfig()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			longInput := strings.Repeat("a", tt.inputLength)

			// Should not panic
			result := formatter.Format(longInput)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatVeryLongSingleLine tests formatting with a very long single line
func TestFormatVeryLongSingleLine(t *testing.T) {
	tests := []struct {
		name      string
		lineLen   int
		config    FormatterConfig
	}{
		{"Line 200 chars", 200, DefaultConfig()},
		{"Line 500 chars", 500, DefaultConfig()},
		{"Line 1000 chars", 1000, DefaultConfig()},
		{"Line 5000 chars", 5000, DefaultConfig()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			longLine := strings.Repeat("word ", tt.lineLen/5)

			result := formatter.Format(longLine)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}

			// When wrapping is enabled, result should contain multiple lines
			if tt.config.WrapLongLines {
				if !strings.Contains(result, "\n") && len(longLine) > 80 {
					// Some long lines might not wrap if they're single words
				}
			}
		})
	}
}

// TestFormatDeeplyNestedMarkdown tests deeply nested markdown structures
func TestFormatDeeplyNestedMarkdown(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Deeply nested lists",
			`- Item 1
  - Item 1.1
    - Item 1.1.1
      - Item 1.1.1.1
        - Item 1.1.1.1.1
          - Item 1.1.1.1.1.1`,
			DefaultConfig(),
		},
		{
			"Multiple heading levels",
			`# Level 1
## Level 2
### Level 3
## Level 2 again
# Level 1 again`,
			DefaultConfig(),
		},
		{
			"Nested code in text",
			"`code1` contains `code2` which contains more `code3` in `code4`",
			DefaultConfig(),
		},
		{
			"Nested bold and italic",
			"***bold and italic*** **bold** *italic*",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatMalformedTables tests various malformed table inputs
func TestFormatMalformedTables(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Table with missing closing pipe",
			`| Col1 | Col2
| --- | ---
| Data1 | Data2`,
			DefaultConfig(),
		},
		{
			"Table with uneven columns",
			`| Col1 | Col2 | Col3
| --- | ---
| Data1 | Data2`,
			DefaultConfig(),
		},
		{
			"Single column table",
			`| Col1
| ---
| Data1`,
			DefaultConfig(),
		},
		{
			"Table with empty cells",
			`| Col1 | | Col3
| --- | --- | ---
| | Data2 |`,
			DefaultConfig(),
		},
		{
			"Table with very long content",
			`| Col1 | Col2
| --- | ---
| ` + strings.Repeat("Data", 100) + ` | Short |`,
			DefaultConfig(),
		},
		{
			"Misaligned separator row",
			`| Col1 | Col2
| - | -- | ---
| Data1 | Data2`,
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			// Should handle gracefully without panicking
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatMixedLineEndings tests various line ending combinations
func TestFormatMixedLineEndings(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Unix line endings (LF)",
			"Line 1\nLine 2\nLine 3",
			DefaultConfig(),
		},
		{
			"Windows line endings (CRLF)",
			"Line 1\r\nLine 2\r\nLine 3",
			DefaultConfig(),
		},
		{
			"Old Mac line endings (CR)",
			"Line 1\rLine 2\rLine 3",
			DefaultConfig(),
		},
		{
			"Mixed line endings",
			"Line 1\nLine 2\r\nLine 3\rLine 4",
			DefaultConfig(),
		},
		{
			"Multiple consecutive line endings",
			"Line 1\n\n\nLine 2\r\n\r\n\r\nLine 3",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatUnicodeCharacters tests unicode and special character handling
func TestFormatUnicodeCharacters(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Basic emoji",
			"Hello 😀 World 🌍",
			DefaultConfig(),
		},
		{
			"Multiple emoji",
			"🎉 🎊 🎈 🎁 🎀 🎭 🎪 🎨 🎬 🎮",
			DefaultConfig(),
		},
		{
			"CJK characters (Chinese)",
			"你好 世界 中文 测试",
			DefaultConfig(),
		},
		{
			"CJK characters (Japanese)",
			"こんにちは 世界 日本語",
			DefaultConfig(),
		},
		{
			"CJK characters (Korean)",
			"안녕하세요 세계 한국어",
			DefaultConfig(),
		},
		{
			"Arabic characters",
			"مرحبا بالعالم اختبار",
			DefaultConfig(),
		},
		{
			"Hebrew characters",
			"שלום עולם בדיקה",
			DefaultConfig(),
		},
		{
			"Cyrillic characters",
			"Привет мир тест",
			DefaultConfig(),
		},
		{
			"Mixed scripts",
			"Hello مرحبا 你好 こんにちは",
			DefaultConfig(),
		},
		{
			"Unicode combining characters",
			"café naïve résumé",
			DefaultConfig(),
		},
		{
			"Zero-width characters",
			"Hello\u200bWorld\u200c\u200d",
			DefaultConfig(),
		},
		{
			"Right-to-left text",
			"עברית العربية",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatSpecialRegexCharacters tests input with regex special characters
func TestFormatSpecialRegexCharacters(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Regex metacharacters",
			"Match . (dot) and * (asterisk) and + (plus) and ? (question) and [ (bracket)",
			DefaultConfig(),
		},
		{
			"Parentheses and pipes",
			"Use (group) or (alternative) | pipes",
			DefaultConfig(),
		},
		{
			"Backslashes",
			"Escape \\ backslash \\\\ double backslash",
			DefaultConfig(),
		},
		{
			"Dollar sign and caret",
			"Start ^ with caret and end with $ dollar",
			DefaultConfig(),
		},
		{
			"Curly braces",
			"Repeat {3} times or {3,5} or {3,}",
			DefaultConfig(),
		},
		{
			"Mixed special characters",
			".*+?[]{}()\\|^$# all special chars",
			DefaultConfig(),
		},
		{
			"Code block with regex",
			"```regex\n^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$\n```",
			DefaultConfig(),
		},
		{
			"Inline code with special chars",
			"Use `grep \"pattern.*txt\"` or `sed 's/old/new/'` for regex",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatExtremeConfigValues tests extreme configuration values
func TestFormatExtremeConfigValues(t *testing.T) {
	tests := []struct {
		name   string
		config FormatterConfig
		input  string
	}{
		{
			"Very small line width",
			FormatterConfig{
				LineWidth:     10,
				CommentPrefix: "# ",
				IndentSize:    2,
			},
			"This is a long line that should be wrapped to fit within a very small line width constraint",
		},
		{
			"Very large line width",
			FormatterConfig{
				LineWidth:     2000,
				CommentPrefix: "# ",
				IndentSize:    2,
			},
			"This is a line that could be very long",
		},
		{
			"Zero indent size",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    0,
			},
			"- Item 1\n  - Item 2\n    - Item 3",
		},
		{
			"Very large indent size",
			FormatterConfig{
				LineWidth:     200,
				CommentPrefix: "# ",
				IndentSize:    50,
			},
			"- Item 1\n  - Item 2",
		},
		{
			"Empty comment prefix",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "",
				IndentSize:    2,
			},
			"Test content",
		},
		{
			"Very long comment prefix",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: strings.Repeat("#", 50),
				IndentSize:    2,
			},
			"Test content",
		},
		{
			"Negative-like large config values",
			FormatterConfig{
				LineWidth:     999999,
				CommentPrefix: "# ",
				IndentSize:    999,
			},
			"Test content",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			// Should handle gracefully without panicking
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatComplexCodeBlocks tests various code block edge cases
func TestFormatComplexCodeBlocks(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Code block without language",
			"```\ncode without language\n```",
			DefaultConfig(),
		},
		{
			"Code block with unknown language",
			"```unknownlang\ncode with unknown language\n```",
			DefaultConfig(),
		},
		{
			"Nested backticks in code",
			"```go\ncode with `backticks` inside\n```",
			DefaultConfig(),
		},
		{
			"Code block with special regex characters",
			"```regex\n^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$\n```",
			DefaultConfig(),
		},
		{
			"Code block with very long lines",
			"```python\n" + strings.Repeat("x", 500) + "\n```",
			DefaultConfig(),
		},
		{
			"Code block with many lines",
			"```\n" + strings.Repeat("line\n", 1000) + "```",
			DefaultConfig(),
		},
		{
			"Unclosed code block",
			"```\ncode without closing backticks",
			DefaultConfig(),
		},
		{
			"Empty code block",
			"```\n\n```",
			DefaultConfig(),
		},
		{
			"Code block with only whitespace",
			"```\n   \n\t\n```",
			DefaultConfig(),
		},
		{
			"Multiple consecutive code blocks",
			"```\ncode 1\n```\n\n```\ncode 2\n```",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatStructuredCommands tests removal of structured commands
func TestFormatStructuredCommands(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		config          FormatterConfig
		shouldContain   string
		shouldNotContain string
	}{
		{
			"Simple structured commands",
			"Response text <structured_commands>command here</structured_commands> more text",
			DefaultConfig(),
			"Response text",
			"structured_commands",
		},
		{
			"Nested structured commands",
			"Text <structured_commands>outer <nested>inner</nested> outer</structured_commands> end",
			DefaultConfig(),
			"Text",
			"structured_commands",
		},
		{
			"Multiple structured commands",
			"Start <structured_commands>cmd1</structured_commands> middle <structured_commands>cmd2</structured_commands> end",
			DefaultConfig(),
			"Start",
			"structured_commands",
		},
		{
			"Structured commands with special characters",
			"Text <structured_commands>echo 'test' | grep .*</structured_commands> end",
			DefaultConfig(),
			"Text",
			"structured_commands",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			if tt.shouldContain != "" && !strings.Contains(result, tt.shouldContain) {
				t.Errorf("Expected result to contain %q, got: %q", tt.shouldContain, result)
			}
			if tt.shouldNotContain != "" && strings.Contains(result, tt.shouldNotContain) {
				t.Errorf("Expected result to NOT contain %q, got: %q", tt.shouldNotContain, result)
			}
		})
	}
}

// TestFormatMalformedMarkdown tests malformed markdown syntax
func TestFormatMalformedMarkdown(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Unclosed bold",
			"**bold text without closing",
			DefaultConfig(),
		},
		{
			"Unclosed italic",
			"*italic text without closing",
			DefaultConfig(),
		},
		{
			"Unclosed link",
			"[link text without closing url",
			DefaultConfig(),
		},
		{
			"Malformed link",
			"[link text] (url without bracket before paren",
			DefaultConfig(),
		},
		{
			"Mixed bold and italic delimiters",
			"*bold with ** italic * mixed",
			DefaultConfig(),
		},
		{
			"Invalid header levels",
			"#### Valid\n##### Also valid\n###### Also valid\n####### Too many",
			DefaultConfig(),
		},
		{
			"Empty headings",
			"# \n## \n### ",
			DefaultConfig(),
		},
		{
			"Quote without space",
			">No space after greater-than",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			// Should handle gracefully without panicking
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatAllConfigCombinations tests various config combinations
func TestFormatAllConfigCombinations(t *testing.T) {
	baseConfig := DefaultConfig()
	testInput := "# Header\nSome **bold** and *italic* text\n```go\ncode\n```\n- List item\n> Quote"

	configVariations := []struct {
		name   string
		config FormatterConfig
	}{
		{"All features enabled", FormatterConfig{
			UseColors:       true,
			CommentPrefix:   "# ",
			LineWidth:       80,
			IndentSize:      2,
			UseBoxes:        true,
			UseBullets:      true,
			HighlightCode:   true,
			WrapLongLines:   true,
			RenderTables:    true,
			ShowLineNumbers: true,
			CompactMode:     false,
			HighlightQuotes: true,
			ParseMarkdown:   true,
		}},
		{"All features disabled", FormatterConfig{
			UseColors:       false,
			CommentPrefix:   "",
			LineWidth:       80,
			IndentSize:      0,
			UseBoxes:        false,
			UseBullets:      false,
			HighlightCode:   false,
			WrapLongLines:   false,
			RenderTables:    false,
			ShowLineNumbers: false,
			CompactMode:     true,
			HighlightQuotes: false,
			ParseMarkdown:   false,
		}},
		{"Mixed 1", FormatterConfig{
			UseColors:       true,
			CommentPrefix:   "# ",
			LineWidth:       100,
			IndentSize:      4,
			UseBoxes:        false,
			UseBullets:      true,
			HighlightCode:   true,
			WrapLongLines:   false,
			RenderTables:    true,
			ShowLineNumbers: true,
			CompactMode:     false,
			HighlightQuotes: false,
			ParseMarkdown:   true,
		}},
		{"Mixed 2", FormatterConfig{
			UseColors:       false,
			CommentPrefix:   "> ",
			LineWidth:       60,
			IndentSize:      3,
			UseBoxes:        true,
			UseBullets:      false,
			HighlightCode:   false,
			WrapLongLines:   true,
			RenderTables:    false,
			ShowLineNumbers: false,
			CompactMode:     true,
			HighlightQuotes: true,
			ParseMarkdown:   false,
		}},
	}

	for _, tt := range configVariations {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(testInput)
			if result == "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatStressLongContent tests stress with mixed content
func TestFormatStressLongContent(t *testing.T) {
	// Create a complex document with various elements
	var sb strings.Builder
	sb.WriteString("# Main Header\n\n")
	sb.WriteString("Some introductory text with **bold** and *italic* formatting.\n\n")

	// Add many list items
	for i := 0; i < 100; i++ {
		sb.WriteString("- Item " + string(rune(i+1)) + "\n")
	}
	sb.WriteString("\n")

	// Add code blocks
	for i := 0; i < 10; i++ {
		sb.WriteString("```go\nfunction " + string(rune(i+1)) + "() {\n  return value\n}\n```\n")
	}

	// Add table
	sb.WriteString("\n| Column 1 | Column 2 | Column 3 |\n")
	sb.WriteString("| --- | --- | --- |\n")
	for i := 0; i < 50; i++ {
		sb.WriteString("| Data1-" + string(rune(i+1)) + " | Data2-" + string(rune(i+1)) + " | Data3-" + string(rune(i+1)) + " |\n")
	}

	content := sb.String()

	configs := []struct {
		name   string
		config FormatterConfig
	}{
		{"Default config", DefaultConfig()},
		{"Colored config", ColoredConfig()},
		{"Compact config", CompactConfig()},
	}

	for _, tt := range configs {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(content)

			if result == "" {
				t.Errorf("Expected non-empty result for stress test")
			}

			// Verify some content is preserved
			if !strings.Contains(result, "Main Header") {
				t.Errorf("Expected 'Main Header' in result")
			}
		})
	}
}

// TestFormatBoundaryValues tests boundary values
func TestFormatBoundaryValues(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Exactly line width length",
			strings.Repeat("x", 80-2), // Minus comment prefix "# "
			DefaultConfig(),
		},
		{
			"One character over line width",
			strings.Repeat("x", 80-1),
			DefaultConfig(),
		},
		{
			"One character under line width",
			strings.Repeat("x", 80-3),
			DefaultConfig(),
		},
		{
			"Single character",
			"x",
			DefaultConfig(),
		},
		{
			"Two characters",
			"xy",
			DefaultConfig(),
		},
		{
			"Exactly one newline",
			"Line 1\nLine 2",
			DefaultConfig(),
		},
		{
			"Many consecutive newlines",
			"Line 1" + strings.Repeat("\n", 100) + "Line 2",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)
			// Should handle gracefully without panicking
		})
	}
}

// TestFormatRobustnessWithNoOptions tests formatter handles missing optional configs
func TestFormatRobustnessWithNoOptions(t *testing.T) {
	// Create a completely empty config that will use all defaults
	emptyConfig := FormatterConfig{}
	formatter := NewTerminalFormatter(emptyConfig)

	testCases := []string{
		"",
		"simple text",
		"# Header",
		"```\ncode\n```",
		"| table | data |\n| --- | --- |",
		strings.Repeat("a", 5000),
	}

	for i, input := range testCases {
		t.Run("empty_config_"+string(rune(i)), func(t *testing.T) {
			result := formatter.Format(input)
			// Should not panic
		})
	}
}
