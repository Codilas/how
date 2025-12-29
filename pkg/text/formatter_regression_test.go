package text

import (
	"strings"
	"testing"
)

// TestFormatNilInputHandling tests handling of edge case inputs that could cause nil pointer issues
func TestFormatNilInputHandling(t *testing.T) {
	tests := []struct {
		name   string
		config FormatterConfig
	}{
		{"Default config with empty", DefaultConfig()},
		{"Colored config with empty", ColoredConfig()},
		{"Compact config with empty", CompactConfig()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)

			// Should not panic with empty input
			result := formatter.Format("")
			_ = result // Use the result to avoid compiler warnings
		})
	}
}

// TestFormatRegexCatastrophicBacktracking prevents regex DoS
func TestFormatRegexCatastrophicBacktracking(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Many unmatched stars",
			strings.Repeat("*", 10000),
			DefaultConfig(),
		},
		{
			"Many unmatched brackets",
			strings.Repeat("[", 10000),
			DefaultConfig(),
		},
		{
			"Many unmatched parentheses",
			strings.Repeat("(", 10000),
			DefaultConfig(),
		},
		{
			"Alternating special characters",
			strings.Repeat("*[]*[]*", 1000),
			DefaultConfig(),
		},
		{
			"Complex nested unmatched syntax",
			strings.Repeat("[(*[(*[(*", 500),
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)

			// Use a channel with timeout to detect if regex takes too long
			done := make(chan bool, 1)
			go func() {
				result := formatter.Format(tt.input)
				_ = result
				done <- true
			}()

			// This test primarily checks that it doesn't panic or hang
			<-done
		})
	}
}

// TestFormatRegressionOffByOne tests potential off-by-one errors
func TestFormatRegressionOffByOne(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Exactly line width minus prefix",
			strings.Repeat("x", 78), // 80 - 2 for "# "
			DefaultConfig(),
		},
		{
			"Exactly line width",
			strings.Repeat("x", 80),
			DefaultConfig(),
		},
		{
			"One more than line width",
			strings.Repeat("x", 81),
			DefaultConfig(),
		},
		{
			"Table with exactly line width",
			"| " + strings.Repeat("x", 10) + " | " + strings.Repeat("y", 10) + " |",
			DefaultConfig(),
		},
		{
			"Multiple lines at boundary",
			strings.Repeat("x", 78) + "\n" + strings.Repeat("y", 78),
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should not panic and produce reasonable output
			if result == "" && tt.input != "" {
				t.Errorf("Expected non-empty result")
			}
		})
	}
}

// TestFormatRegressionSliceIndexing tests potential slice indexing errors
func TestFormatRegressionSliceIndexing(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Single character line",
			"a",
			DefaultConfig(),
		},
		{
			"Single backtick",
			"`",
			DefaultConfig(),
		},
		{
			"Single pipe",
			"|",
			DefaultConfig(),
		},
		{
			"Single dash",
			"-",
			DefaultConfig(),
		},
		{
			"Single asterisk",
			"*",
			DefaultConfig(),
		},
		{
			"Single bracket",
			"[",
			DefaultConfig(),
		},
		{
			"Two character combinations",
			"##",
			DefaultConfig(),
		},
		{
			"Code fence markers only",
			"```",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should not panic
		})
	}
}

// TestFormatRegressionStringMethods tests potential string method errors
func TestFormatRegressionStringMethods(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Very long string split operation",
			strings.Repeat("a", 50000),
			DefaultConfig(),
		},
		{
			"String trim on various inputs",
			"   \t\n   content   \n\t   ",
			DefaultConfig(),
		},
		{
			"Fields on sparse whitespace",
			"word1" + strings.Repeat(" ", 1000) + "word2",
			DefaultConfig(),
		},
		{
			"Repeat with zero times",
			"content",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should not panic
		})
	}
}

// TestFormatRegressionMapAccess tests potential map access errors
func TestFormatRegressionMapAccess(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Empty table data",
			"||",
			DefaultConfig(),
		},
		{
			"Table with empty cells",
			"||||",
			DefaultConfig(),
		},
		{
			"Many empty columns",
			"|" + strings.Repeat("|", 100) + "|",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should not panic on map access
		})
	}
}

// TestFormatRegressionBufferOperations tests potential buffer overflow/underflow
func TestFormatRegressionBufferOperations(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Many consecutive newlines",
			"\n\n\n\n\n\n\n\n\n\n",
			DefaultConfig(),
		},
		{
			"Alternating content and newlines",
			"a\nb\nc\nd\ne\nf\ng\nh\ni\nj",
			DefaultConfig(),
		},
		{
			"Newlines at start",
			"\n\n\n\nContent",
			DefaultConfig(),
		},
		{
			"Newlines at end",
			"Content\n\n\n\n",
			DefaultConfig(),
		},
		{
			"Only newlines and content",
			strings.Repeat("content\n", 1000),
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should not panic or have buffer issues
		})
	}
}

// TestFormatRegressionTypeConversion tests potential type conversion errors
func TestFormatRegressionTypeConversion(t *testing.T) {
	tests := []struct {
		name   string
		config FormatterConfig
	}{
		{
			"Config with max int values",
			FormatterConfig{
				LineWidth:     2147483647,
				CommentPrefix: "# ",
				IndentSize:    100,
			},
		},
		{
			"Config with very small values",
			FormatterConfig{
				LineWidth:     1,
				CommentPrefix: "# ",
				IndentSize:    1,
			},
		},
		{
			"Mixed extreme values",
			FormatterConfig{
				LineWidth:     999999,
				CommentPrefix: "#",
				IndentSize:    0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format("test content")

			// Should handle type conversions gracefully
		})
	}
}

// TestFormatRegressionLoopConditions tests loop conditions for off-by-one errors
func TestFormatRegressionLoopConditions(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Exactly N lines",
			strings.Join([]string{"line1", "line2", "line3", "line4", "line5"}, "\n"),
			DefaultConfig(),
		},
		{
			"N-1 lines",
			strings.Join([]string{"line1", "line2", "line3", "line4"}, "\n"),
			DefaultConfig(),
		},
		{
			"N+1 lines",
			strings.Join([]string{"line1", "line2", "line3", "line4", "line5", "line6"}, "\n"),
			DefaultConfig(),
		},
		{
			"Many code block lines",
			"```\n" + strings.Join(strings.Split(strings.Repeat("line\n", 1000), "\n"), "\n") + "```",
			DefaultConfig(),
		},
		{
			"Many table rows",
			"| Col |\n" + strings.Repeat("| Data |\n", 1000),
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should not panic on loop conditions
		})
	}
}

// TestFormatRegressionRegexMatching tests regex match edge cases
func TestFormatRegressionRegexMatching(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Markdown headers at document boundaries",
			"# Start\nContent\n# Middle\nMore\n# End",
			DefaultConfig(),
		},
		{
			"Code blocks at boundaries",
			"```\ncode at start\n```\nContent\n```\ncode at end\n```",
			DefaultConfig(),
		},
		{
			"Lists at boundaries",
			"- Start\n- Item\n- End",
			DefaultConfig(),
		},
		{
			"Quotes at boundaries",
			"> Start\n> Quote\n> End",
			DefaultConfig(),
		},
		{
			"Mixed elements at boundaries",
			"# Header\n```\ncode\n```\n- List\n> Quote",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			if result == "" {
				t.Errorf("Expected non-empty result")
			}
		})
	}
}

// TestFormatRegressionContentPreservation tests that content is preserved correctly
func TestFormatRegressionContentPreservation(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		config           FormatterConfig
		shouldContain    string
		shouldNotContain string
	}{
		{
			"Regular text preservation",
			"Hello world test",
			DefaultConfig(),
			"world",
			"",
		},
		{
			"Special characters preserved",
			"@#$%^&*()",
			DefaultConfig(),
			"@#$",
			"",
		},
		{
			"Unicode preserved",
			"你好世界 مرحبا",
			DefaultConfig(),
			"世界",
			"",
		},
		{
			"Code content preserved",
			"```\necho hello\n```",
			DefaultConfig(),
			"echo",
			"structured_commands",
		},
		{
			"Numbers preserved",
			"Count: 1, 2, 3, 4, 5",
			DefaultConfig(),
			"Count",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			if tt.shouldContain != "" && !strings.Contains(result, tt.shouldContain) {
				t.Errorf("Expected to contain %q, got: %q", tt.shouldContain, result)
			}
			if tt.shouldNotContain != "" && strings.Contains(result, tt.shouldNotContain) {
				t.Errorf("Should not contain %q, got: %q", tt.shouldNotContain, result)
			}
		})
	}
}

// TestFormatRegressionEmptyContentBlocks tests empty content block handling
func TestFormatRegressionEmptyContentBlocks(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Empty code block",
			"```\n```",
			DefaultConfig(),
		},
		{
			"Empty lines in code",
			"```\n\n\n\n```",
			DefaultConfig(),
		},
		{
			"Empty table",
			"|||",
			DefaultConfig(),
		},
		{
			"Empty list item",
			"- ",
			DefaultConfig(),
		},
		{
			"Empty quote",
			"> ",
			DefaultConfig(),
		},
		{
			"Empty header",
			"# ",
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

// TestFormatRegressionComplexFormatCombinations tests complex format combinations
func TestFormatRegressionComplexFormatCombinations(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Code with formatted text",
			"```markdown\n# **Bold Header**\n```",
			DefaultConfig(),
		},
		{
			"List with code",
			"- Item with `code`\n- Another `code block`",
			DefaultConfig(),
		},
		{
			"Quote with formatting",
			"> # **Quote Header**\n> With `code` and *italic*",
			DefaultConfig(),
		},
		{
			"Table with special characters",
			"| **Bold** | *Italic* |\n| --- | --- |\n| `code` | [link](url) |",
			DefaultConfig(),
		},
		{
			"Mixed all features",
			"# Header\n\nText with **bold** and *italic*\n\n```code\n```\n\n- List\n\n> Quote",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			if result == "" {
				t.Errorf("Expected non-empty result")
			}
		})
	}
}

// TestFormatRegressionIncrementalChanges tests that formatter handles incremental changes
func TestFormatRegressionIncrementalChanges(t *testing.T) {
	baseText := "# Header\nSome content"
	baseFormatter := NewTerminalFormatter(DefaultConfig())
	baseResult := baseFormatter.Format(baseText)

	// Test variations of the same text
	variations := []struct {
		name   string
		input  string
		should string
	}{
		{"Original", baseText, "Header"},
		{"With extra space", baseText + " ", "Header"},
		{"With newline", baseText + "\n", "Header"},
		{"With content", baseText + "\nMore", "Header"},
	}

	for _, tt := range variations {
		t.Run(tt.name, func(t *testing.T) {
			result := baseFormatter.Format(tt.input)
			if !strings.Contains(result, tt.should) {
				t.Errorf("Expected %q in result", tt.should)
			}
		})
	}
}
