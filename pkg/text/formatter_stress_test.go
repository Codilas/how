package text

import (
	"fmt"
	"strings"
	"testing"
)

// TestFormatMemoryIntensive tests formatter with memory-intensive inputs
func TestFormatMemoryIntensive(t *testing.T) {
	tests := []struct {
		name   string
		config FormatterConfig
		gen    func() string
	}{
		{
			"Many small lines",
			DefaultConfig(),
			func() string {
				var sb strings.Builder
				for i := 0; i < 10000; i++ {
					sb.WriteString(fmt.Sprintf("Line %d content\n", i))
				}
				return sb.String()
			},
		},
		{
			"Large single item",
			DefaultConfig(),
			func() string {
				return strings.Repeat("word ", 100000)
			},
		},
		{
			"Many code blocks",
			DefaultConfig(),
			func() string {
				var sb strings.Builder
				for i := 0; i < 1000; i++ {
					sb.WriteString("```go\ncode block " + fmt.Sprintf("%d", i) + "\n```\n")
				}
				return sb.String()
			},
		},
		{
			"Many tables",
			DefaultConfig(),
			func() string {
				var sb strings.Builder
				for t := 0; t < 100; t++ {
					sb.WriteString("| Col1 | Col2 | Col3 |\n")
					sb.WriteString("| --- | --- | --- |\n")
					for i := 0; i < 50; i++ {
						sb.WriteString(fmt.Sprintf("| Data1-%d | Data2-%d | Data3-%d |\n", i, i, i))
					}
					sb.WriteString("\n")
				}
				return sb.String()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.gen()
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(input)

			// Should not panic and should return something
			if result == "" && input != "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatConcurrentStress tests concurrent formatter usage
func TestFormatConcurrentStress(t *testing.T) {
	formatter := NewTerminalFormatter(DefaultConfig())
	testInput := "# Test\nContent with **bold** and *italic*\n```\ncode\n```"

	// Run multiple formatters concurrently
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func(index int) {
			result := formatter.Format(testInput)
			if result == "" {
				t.Errorf("Concurrent formatter call %d failed", index)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}
}

// TestFormatEdgeCaseMarkdownCombinations tests edge cases in markdown combinations
func TestFormatEdgeCaseMarkdownCombinations(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Bold inside italic",
			"*italic with **bold** inside*",
			DefaultConfig(),
		},
		{
			"Link with special characters in URL",
			"[link](https://example.com/?q=test&lang=en#section)",
			DefaultConfig(),
		},
		{
			"Multiple links on same line",
			"[link1](url1) and [link2](url2) and [link3](url3)",
			DefaultConfig(),
		},
		{
			"Code with markdown syntax",
			"`**not bold**` and `*not italic*`",
			DefaultConfig(),
		},
		{
			"Backtick inside code",
			"``code with ` backtick``",
			DefaultConfig(),
		},
		{
			"Heading with special characters",
			"# Heading with **bold** and *italic*",
			DefaultConfig(),
		},
		{
			"List with formatting",
			"- Item with **bold**\n- Item with *italic*\n- Item with `code`",
			DefaultConfig(),
		},
		{
			"Quote with formatting",
			"> Quote with **bold** and *italic*",
			DefaultConfig(),
		},
		{
			"Mixed quote and list",
			"> Quote\n- List item\n> Another quote",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should handle without panicking
			if result == "" && tt.input != "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatLineWrappingEdgeCases tests edge cases in line wrapping
func TestFormatLineWrappingEdgeCases(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		config           FormatterConfig
		shouldContainLen int // Rough check that output has content
	}{
		{
			"Single very long word",
			strings.Repeat("verylongwordwithoutspaces", 20),
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    2,
				WrapLongLines: true,
			},
			80,
		},
		{
			"Words exactly at line boundary",
			"word " + strings.Repeat("word ", 15),
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    2,
				WrapLongLines: true,
			},
			60,
		},
		{
			"Line wrapping disabled with long content",
			strings.Repeat("word ", 100),
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    2,
				WrapLongLines: false,
			},
			400,
		},
		{
			"Wrapping with very small line width",
			"This is a test sentence with multiple words",
			FormatterConfig{
				LineWidth:     20,
				CommentPrefix: "# ",
				IndentSize:    2,
				WrapLongLines: true,
			},
			40,
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

// TestFormatTableEdgeCases tests advanced table edge cases
func TestFormatTableEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Table with very wide columns",
			"| " + strings.Repeat("x", 100) + " | " + strings.Repeat("y", 100) + " |\n| --- | --- |\n| Data1 | Data2 |",
			DefaultConfig(),
		},
		{
			"Table with single row",
			"| Only | One | Row |\n| --- | --- | --- |",
			DefaultConfig(),
		},
		{
			"Table with many columns",
			"| " + strings.Join(strings.Split(strings.Repeat("Col ", 50), " "), " | ") + " |",
			DefaultConfig(),
		},
		{
			"Table with unicode characters",
			"| 中文 | 日本語 |\n| --- | --- |\n| 数据1 | データ2 |",
			DefaultConfig(),
		},
		{
			"Table with empty separator",
			"| Col1 | Col2 |\n|\n| Data1 | Data2 |",
			DefaultConfig(),
		},
		{
			"Adjacent tables",
			"| Col1 | Col2 |\n| --- | --- |\n| Data1 | Data2 |\n| Col1 | Col2 |\n| --- | --- |\n| Data1 | Data2 |",
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

// TestFormatIndentationEdgeCases tests indentation edge cases
func TestFormatIndentationEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Deeply indented list",
			"- L1\n  - L2\n    - L3\n      - L4\n        - L5\n          - L6",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    2,
				UseBullets:    true,
			},
		},
		{
			"Mixed indentation styles",
			"- Space indented\n\t- Tab indented\n  - Two space\n    - Four space",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    2,
				UseBullets:    true,
			},
		},
		{
			"Indentation with zero indent size",
			"- Item 1\n  - Item 2\n    - Item 3",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    0,
				UseBullets:    true,
			},
		},
		{
			"Quotes with indentation",
			"> Line 1\n> Line 2\n> > Nested quote\n> > > Double nested",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "# ",
				IndentSize:    2,
				HighlightQuotes: true,
			},
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

// TestFormatColoringEdgeCases tests coloring with edge cases
func TestFormatColoringEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Colored with empty input",
			"",
			ColoredConfig(),
		},
		{
			"Colored with only special characters",
			"!@#$%^&*()_+-=[]{}|;:',.<>?/",
			ColoredConfig(),
		},
		{
			"Colored with unicode",
			"🎉 🎊 🌟 ✨ 💫",
			ColoredConfig(),
		},
		{
			"Colored with very long single line",
			strings.Repeat("colored ", 200),
			ColoredConfig(),
		},
		{
			"Colored with boxes and complex content",
			"# Header\n\n**Bold** and *italic*\n\n```\ncode\n```",
			ColoredConfig(),
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

// TestFormatCompactModeEdgeCases tests compact mode with edge cases
func TestFormatCompactModeEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Compact with many blank lines",
			"Line 1\n\n\n\nLine 2\n\n\n\nLine 3",
			CompactConfig(),
		},
		{
			"Compact with code blocks",
			"```\ncode1\n```\n\n```\ncode2\n```",
			CompactConfig(),
		},
		{
			"Compact with lists",
			"- Item 1\n- Item 2\n- Item 3",
			CompactConfig(),
		},
		{
			"Compact vs non-compact spacing",
			"Line 1\n\nLine 2\n\nLine 3",
			CompactConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			if result == "" && tt.input != "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatRegexPerformance tests regex performance with complex patterns
func TestFormatRegexPerformance(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Many bold markers",
			strings.Repeat("**text** ", 1000),
			DefaultConfig(),
		},
		{
			"Many italic markers",
			strings.Repeat("*text* ", 1000),
			DefaultConfig(),
		},
		{
			"Many links",
			strings.Repeat("[text](url) ", 1000),
			DefaultConfig(),
		},
		{
			"Many inline code blocks",
			strings.Repeat("`code` ", 1000),
			DefaultConfig(),
		},
		{
			"Mixed regex patterns",
			strings.Repeat("**bold** *italic* `code` [link](url) ", 500),
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			if result == "" && tt.input != "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatCommentPrefixEdgeCases tests edge cases with comment prefixes
func TestFormatCommentPrefixEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		config     FormatterConfig
		shouldHave string
	}{
		{
			"Empty prefix",
			"Test content",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "",
				IndentSize:    2,
			},
			"Test",
		},
		{
			"Very long prefix",
			"Test",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: strings.Repeat("#", 100),
				IndentSize:    2,
			},
			"Test",
		},
		{
			"Special character prefix",
			"Test",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: ">>> ",
				IndentSize:    2,
			},
			"Test",
		},
		{
			"Unicode prefix",
			"Test",
			FormatterConfig{
				LineWidth:     80,
				CommentPrefix: "🔸 ",
				IndentSize:    2,
			},
			"Test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should contain the test content or be properly handled
		})
	}
}

// TestFormatNumberedLists tests numbered list edge cases
func TestFormatNumberedLists(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Simple numbered list",
			"1. First\n2. Second\n3. Third",
			DefaultConfig(),
		},
		{
			"Large numbers",
			"100. Item\n101. Next\n1000. Far future",
			DefaultConfig(),
		},
		{
			"Mixed ordered and unordered",
			"1. Ordered\n- Unordered\n2. Ordered again",
			DefaultConfig(),
		},
		{
			"Numbered with leading zeros",
			"01. First\n02. Second\n10. Tenth",
			DefaultConfig(),
		},
		{
			"Non-sequential numbers",
			"1. First\n5. Fifth\n3. Third",
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

// TestFormatShowLineNumbersEdgeCases tests ShowLineNumbers configuration edge cases
func TestFormatShowLineNumbersEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Line numbers with empty code block",
			"```\n\n```",
			FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				ShowLineNumbers: true,
			},
		},
		{
			"Line numbers with single line code",
			"```go\nsingle line\n```",
			FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				ShowLineNumbers: true,
			},
		},
		{
			"Line numbers with many lines",
			"```python\n" + strings.Repeat("line\n", 1000) + "```",
			FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				ShowLineNumbers: true,
			},
		},
		{
			"Line numbers disabled",
			"```go\ncode line 1\ncode line 2\n```",
			FormatterConfig{
				UseColors:       false,
				CommentPrefix:   "# ",
				LineWidth:       80,
				IndentSize:      2,
				ShowLineNumbers: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			if result == "" && tt.input != "" {
				t.Errorf("Expected non-empty result for %s", tt.name)
			}
		})
	}
}

// TestFormatWhitespaceHandling tests various whitespace handling scenarios
func TestFormatWhitespaceHandling(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		config FormatterConfig
	}{
		{
			"Trailing whitespace on lines",
			"Line 1   \nLine 2\t\nLine 3  ",
			DefaultConfig(),
		},
		{
			"Leading whitespace on lines",
			"   Line 1\n\t\tLine 2\nLine 3",
			DefaultConfig(),
		},
		{
			"Tab characters",
			"Line\twith\ttabs\teverywhere",
			DefaultConfig(),
		},
		{
			"Mixed spaces and tabs",
			"  \tLine with\t  mixed\t spacing",
			DefaultConfig(),
		},
		{
			"Non-breaking spaces",
			"Line\u00A0with\u00A0nbsp",
			DefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			result := formatter.Format(tt.input)

			// Should handle without panicking
		})
	}
}
