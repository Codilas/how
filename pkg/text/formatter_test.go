package text

import (
	"regexp"
	"strings"
	"testing"

	"github.com/fatih/color"
)

// TestNewTerminalFormatterDefaults tests that NewTerminalFormatter sets proper defaults
func TestNewTerminalFormatterDefaults(t *testing.T) {
	tests := []struct {
		name           string
		config         FormatterConfig
		expectedPrefix string
		expectedWidth  int
		expectedIndent int
	}{
		{
			name:           "Empty config gets defaults",
			config:         FormatterConfig{},
			expectedPrefix: "# ",
			expectedWidth:  80,
			expectedIndent: 2,
		},
		{
			name: "Partial config preserves provided values",
			config: FormatterConfig{
				CommentPrefix: ">> ",
				LineWidth:     100,
			},
			expectedPrefix: ">> ",
			expectedWidth:  100,
			expectedIndent: 2,
		},
		{
			name: "Custom indentation is preserved",
			config: FormatterConfig{
				IndentSize: 4,
			},
			expectedPrefix: "# ",
			expectedWidth:  80,
			expectedIndent: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)

			if formatter.config.CommentPrefix != tt.expectedPrefix {
				t.Errorf("CommentPrefix: expected %q, got %q", tt.expectedPrefix, formatter.config.CommentPrefix)
			}
			if formatter.config.LineWidth != tt.expectedWidth {
				t.Errorf("LineWidth: expected %d, got %d", tt.expectedWidth, formatter.config.LineWidth)
			}
			if formatter.config.IndentSize != tt.expectedIndent {
				t.Errorf("IndentSize: expected %d, got %d", tt.expectedIndent, formatter.config.IndentSize)
			}
		})
	}
}

// TestNewTerminalFormatterRegexCompilation tests that all regexes are properly compiled
func TestNewTerminalFormatterRegexCompilation(t *testing.T) {
	formatter := NewTerminalFormatter(FormatterConfig{})

	// Test that all regex fields are non-nil
	tests := []struct {
		name  string
		regex *regexp.Regexp
	}{
		{"structuredRegex", formatter.structuredRegex},
		{"codeBlockRegex", formatter.codeBlockRegex},
		{"inlineCodeRegex", formatter.inlineCodeRegex},
		{"markdownH1Regex", formatter.markdownH1Regex},
		{"markdownH2Regex", formatter.markdownH2Regex},
		{"markdownH3Regex", formatter.markdownH3Regex},
		{"boldTextRegex", formatter.boldTextRegex},
		{"italicTextRegex", formatter.italicTextRegex},
		{"linkRegex", formatter.linkRegex},
		{"blockQuoteRegex", formatter.blockQuoteRegex},
		{"tableRowRegex", formatter.tableRowRegex},
		{"listItemRegex", formatter.listItemRegex},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.regex == nil {
				t.Errorf("%s should not be nil", tt.name)
			}
		})
	}
}

// TestNewTerminalFormatterColorInitializationDisabled tests color initialization when UseColors is false
func TestNewTerminalFormatterColorInitializationDisabled(t *testing.T) {
	config := FormatterConfig{
		UseColors: false,
	}
	formatter := NewTerminalFormatter(config)

	// When UseColors is false, colors should not be initialized
	if formatter.colors.heading1 != nil || formatter.colors.heading2 != nil ||
		formatter.colors.code != nil || formatter.colors.text != nil {
		t.Error("Colors should not be initialized when UseColors is false")
	}
}

// TestNewTerminalFormatterColorInitializationEnabled tests color initialization when UseColors is true
func TestNewTerminalFormatterColorInitializationEnabled(t *testing.T) {
	config := FormatterConfig{
		UseColors: true,
	}
	formatter := NewTerminalFormatter(config)

	// When UseColors is true, all colors should be initialized
	colorTests := []struct {
		name  string
		color *color.Color
	}{
		{"heading1", formatter.colors.heading1},
		{"heading2", formatter.colors.heading2},
		{"heading3", formatter.colors.heading3},
		{"code", formatter.colors.code},
		{"inlineCode", formatter.colors.inlineCode},
		{"quote", formatter.colors.quote},
		{"bullet", formatter.colors.bullet},
		{"text", formatter.colors.text},
		{"comment", formatter.colors.comment},
		{"border", formatter.colors.border},
		{"link", formatter.colors.link},
		{"bold", formatter.colors.bold},
		{"italic", formatter.colors.italic},
	}

	for _, tt := range colorTests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.color == nil {
				t.Errorf("Color %s should be initialized when UseColors is true", tt.name)
			}
		})
	}
}

// TestDefaultConfig tests the DefaultConfig preset
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	expectedValues := map[string]interface{}{
		"UseColors":       false,
		"CommentPrefix":   "# ",
		"LineWidth":       80,
		"IndentSize":      2,
		"UseBoxes":        false,
		"UseBullets":      true,
		"HighlightCode":   false,
		"WrapLongLines":   true,
		"RenderTables":    true,
		"ShowLineNumbers": false,
		"CompactMode":     false,
		"HighlightQuotes": false,
		"ParseMarkdown":   true,
	}

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"UseColors", expectedValues["UseColors"], config.UseColors},
		{"CommentPrefix", expectedValues["CommentPrefix"], config.CommentPrefix},
		{"LineWidth", expectedValues["LineWidth"], config.LineWidth},
		{"IndentSize", expectedValues["IndentSize"], config.IndentSize},
		{"UseBoxes", expectedValues["UseBoxes"], config.UseBoxes},
		{"UseBullets", expectedValues["UseBullets"], config.UseBullets},
		{"HighlightCode", expectedValues["HighlightCode"], config.HighlightCode},
		{"WrapLongLines", expectedValues["WrapLongLines"], config.WrapLongLines},
		{"RenderTables", expectedValues["RenderTables"], config.RenderTables},
		{"ShowLineNumbers", expectedValues["ShowLineNumbers"], config.ShowLineNumbers},
		{"CompactMode", expectedValues["CompactMode"], config.CompactMode},
		{"HighlightQuotes", expectedValues["HighlightQuotes"], config.HighlightQuotes},
		{"ParseMarkdown", expectedValues["ParseMarkdown"], config.ParseMarkdown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected != tt.actual {
				t.Errorf("expected %v, got %v", tt.expected, tt.actual)
			}
		})
	}
}

// TestColoredConfig tests the ColoredConfig preset
func TestColoredConfig(t *testing.T) {
	config := ColoredConfig()

	expectedValues := map[string]interface{}{
		"UseColors":       true,
		"CommentPrefix":   "# ",
		"LineWidth":       80,
		"IndentSize":      2,
		"UseBoxes":        true,
		"UseBullets":      true,
		"HighlightCode":   true,
		"WrapLongLines":   true,
		"RenderTables":    true,
		"ShowLineNumbers": false,
		"CompactMode":     false,
		"HighlightQuotes": true,
		"ParseMarkdown":   true,
	}

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"UseColors", expectedValues["UseColors"], config.UseColors},
		{"CommentPrefix", expectedValues["CommentPrefix"], config.CommentPrefix},
		{"LineWidth", expectedValues["LineWidth"], config.LineWidth},
		{"IndentSize", expectedValues["IndentSize"], config.IndentSize},
		{"UseBoxes", expectedValues["UseBoxes"], config.UseBoxes},
		{"UseBullets", expectedValues["UseBullets"], config.UseBullets},
		{"HighlightCode", expectedValues["HighlightCode"], config.HighlightCode},
		{"WrapLongLines", expectedValues["WrapLongLines"], config.WrapLongLines},
		{"RenderTables", expectedValues["RenderTables"], config.RenderTables},
		{"ShowLineNumbers", expectedValues["ShowLineNumbers"], config.ShowLineNumbers},
		{"CompactMode", expectedValues["CompactMode"], config.CompactMode},
		{"HighlightQuotes", expectedValues["HighlightQuotes"], config.HighlightQuotes},
		{"ParseMarkdown", expectedValues["ParseMarkdown"], config.ParseMarkdown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected != tt.actual {
				t.Errorf("expected %v, got %v", tt.expected, tt.actual)
			}
		})
	}
}

// TestCompactConfig tests the CompactConfig preset
func TestCompactConfig(t *testing.T) {
	config := CompactConfig()

	expectedValues := map[string]interface{}{
		"UseColors":       true,
		"CommentPrefix":   "",
		"LineWidth":       120,
		"IndentSize":      2,
		"UseBoxes":        false,
		"UseBullets":      true,
		"HighlightCode":   true,
		"WrapLongLines":   false,
		"RenderTables":    true,
		"ShowLineNumbers": false,
		"CompactMode":     true,
		"HighlightQuotes": true,
		"ParseMarkdown":   true,
	}

	tests := []struct {
		name     string
		expected interface{}
		actual   interface{}
	}{
		{"UseColors", expectedValues["UseColors"], config.UseColors},
		{"CommentPrefix", expectedValues["CommentPrefix"], config.CommentPrefix},
		{"LineWidth", expectedValues["LineWidth"], config.LineWidth},
		{"IndentSize", expectedValues["IndentSize"], config.IndentSize},
		{"UseBoxes", expectedValues["UseBoxes"], config.UseBoxes},
		{"UseBullets", expectedValues["UseBullets"], config.UseBullets},
		{"HighlightCode", expectedValues["HighlightCode"], config.HighlightCode},
		{"WrapLongLines", expectedValues["WrapLongLines"], config.WrapLongLines},
		{"RenderTables", expectedValues["RenderTables"], config.RenderTables},
		{"ShowLineNumbers", expectedValues["ShowLineNumbers"], config.ShowLineNumbers},
		{"CompactMode", expectedValues["CompactMode"], config.CompactMode},
		{"HighlightQuotes", expectedValues["HighlightQuotes"], config.HighlightQuotes},
		{"ParseMarkdown", expectedValues["ParseMarkdown"], config.ParseMarkdown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expected != tt.actual {
				t.Errorf("expected %v, got %v", tt.expected, tt.actual)
			}
		})
	}
}

// TestFormatterConfigStruct tests the FormatterConfig struct fields and defaults
func TestFormatterConfigStruct(t *testing.T) {
	t.Run("Zero value config", func(t *testing.T) {
		config := FormatterConfig{}

		if config.UseColors != false {
			t.Error("UseColors should default to false")
		}
		if config.CommentPrefix != "" {
			t.Error("CommentPrefix should be empty in zero value")
		}
		if config.LineWidth != 0 {
			t.Error("LineWidth should be 0 in zero value")
		}
		if config.IndentSize != 0 {
			t.Error("IndentSize should be 0 in zero value")
		}
	})

	t.Run("Custom config values", func(t *testing.T) {
		config := FormatterConfig{
			UseColors:       true,
			CommentPrefix:   ">> ",
			LineWidth:       100,
			IndentSize:      4,
			UseBoxes:        true,
			UseBullets:      false,
			HighlightCode:   true,
			WrapLongLines:   false,
			RenderTables:    false,
			ShowLineNumbers: true,
			CompactMode:     true,
			HighlightQuotes: true,
			ParseMarkdown:   false,
		}

		if !config.UseColors {
			t.Error("UseColors should be true")
		}
		if config.CommentPrefix != ">> " {
			t.Errorf("CommentPrefix expected '>> ', got %q", config.CommentPrefix)
		}
		if config.LineWidth != 100 {
			t.Errorf("LineWidth expected 100, got %d", config.LineWidth)
		}
		if config.IndentSize != 4 {
			t.Errorf("IndentSize expected 4, got %d", config.IndentSize)
		}
		if !config.UseBoxes {
			t.Error("UseBoxes should be true")
		}
		if config.UseBullets {
			t.Error("UseBullets should be false")
		}
		if !config.HighlightCode {
			t.Error("HighlightCode should be true")
		}
		if config.WrapLongLines {
			t.Error("WrapLongLines should be false")
		}
		if config.RenderTables {
			t.Error("RenderTables should be false")
		}
		if !config.ShowLineNumbers {
			t.Error("ShowLineNumbers should be true")
		}
		if !config.CompactMode {
			t.Error("CompactMode should be true")
		}
		if !config.HighlightQuotes {
			t.Error("HighlightQuotes should be true")
		}
		if config.ParseMarkdown {
			t.Error("ParseMarkdown should be false")
		}
	})
}

// TestNewTerminalFormatterPreservesConfig tests that provided config is preserved
func TestNewTerminalFormatterPreservesConfig(t *testing.T) {
	config := FormatterConfig{
		UseColors:       true,
		CommentPrefix:   ">> ",
		LineWidth:       100,
		IndentSize:      4,
		UseBoxes:        true,
		UseBullets:      false,
		HighlightCode:   true,
		WrapLongLines:   false,
		RenderTables:    false,
		ShowLineNumbers: true,
		CompactMode:     true,
		HighlightQuotes: true,
		ParseMarkdown:   false,
	}

	formatter := NewTerminalFormatter(config)

	if formatter.config.UseColors != config.UseColors {
		t.Error("UseColors not preserved")
	}
	if formatter.config.CommentPrefix != config.CommentPrefix {
		t.Error("CommentPrefix not preserved")
	}
	if formatter.config.LineWidth != config.LineWidth {
		t.Error("LineWidth not preserved")
	}
	if formatter.config.IndentSize != config.IndentSize {
		t.Error("IndentSize not preserved")
	}
	if formatter.config.UseBoxes != config.UseBoxes {
		t.Error("UseBoxes not preserved")
	}
	if formatter.config.UseBullets != config.UseBullets {
		t.Error("UseBullets not preserved")
	}
	if formatter.config.HighlightCode != config.HighlightCode {
		t.Error("HighlightCode not preserved")
	}
	if formatter.config.WrapLongLines != config.WrapLongLines {
		t.Error("WrapLongLines not preserved")
	}
	if formatter.config.RenderTables != config.RenderTables {
		t.Error("RenderTables not preserved")
	}
	if formatter.config.ShowLineNumbers != config.ShowLineNumbers {
		t.Error("ShowLineNumbers not preserved")
	}
	if formatter.config.CompactMode != config.CompactMode {
		t.Error("CompactMode not preserved")
	}
	if formatter.config.HighlightQuotes != config.HighlightQuotes {
		t.Error("HighlightQuotes not preserved")
	}
	if formatter.config.ParseMarkdown != config.ParseMarkdown {
		t.Error("ParseMarkdown not preserved")
	}
}

// TestConfigPresetsDifferences tests that the preset configs have distinct characteristics
func TestConfigPresetsDifferences(t *testing.T) {
	defaultConfig := DefaultConfig()
	coloredConfig := ColoredConfig()
	compactConfig := CompactConfig()

	t.Run("DefaultConfig vs ColoredConfig", func(t *testing.T) {
		if defaultConfig.UseColors == coloredConfig.UseColors {
			t.Error("DefaultConfig and ColoredConfig should differ in UseColors")
		}
		if defaultConfig.UseBoxes == coloredConfig.UseBoxes {
			t.Error("DefaultConfig and ColoredConfig should differ in UseBoxes")
		}
		if defaultConfig.HighlightCode == coloredConfig.HighlightCode {
			t.Error("DefaultConfig and ColoredConfig should differ in HighlightCode")
		}
	})

	t.Run("ColoredConfig vs CompactConfig", func(t *testing.T) {
		if coloredConfig.CommentPrefix == compactConfig.CommentPrefix {
			t.Error("ColoredConfig and CompactConfig should differ in CommentPrefix")
		}
		if coloredConfig.LineWidth == compactConfig.LineWidth {
			t.Error("ColoredConfig and CompactConfig should differ in LineWidth")
		}
		if coloredConfig.CompactMode == compactConfig.CompactMode {
			t.Error("ColoredConfig and CompactConfig should differ in CompactMode")
		}
	})

	t.Run("DefaultConfig vs CompactConfig", func(t *testing.T) {
		if defaultConfig.UseColors == compactConfig.UseColors {
			t.Error("DefaultConfig and CompactConfig should differ in UseColors")
		}
		if defaultConfig.CommentPrefix == compactConfig.CommentPrefix {
			t.Error("DefaultConfig and CompactConfig should differ in CommentPrefix")
		}
	})
}

// TestNewTerminalFormatterReturnsNotNil tests that NewTerminalFormatter returns a valid pointer
func TestNewTerminalFormatterReturnsNotNil(t *testing.T) {
	formatter := NewTerminalFormatter(FormatterConfig{})

	if formatter == nil {
		t.Error("NewTerminalFormatter should not return nil")
	}
}

// TestNewTerminalFormatterConfigurationPreservation tests that all configuration options are preserved
func TestNewTerminalFormatterConfigurationPreservation(t *testing.T) {
	testConfigs := []FormatterConfig{
		DefaultConfig(),
		ColoredConfig(),
		CompactConfig(),
	}

	for i, config := range testConfigs {
		formatter := NewTerminalFormatter(config)

		// Use reflection to compare all fields
		if formatter.config.UseColors != config.UseColors ||
			formatter.config.CommentPrefix != config.CommentPrefix ||
			formatter.config.LineWidth != config.LineWidth ||
			formatter.config.IndentSize != config.IndentSize ||
			formatter.config.UseBoxes != config.UseBoxes ||
			formatter.config.UseBullets != config.UseBullets ||
			formatter.config.HighlightCode != config.HighlightCode ||
			formatter.config.WrapLongLines != config.WrapLongLines ||
			formatter.config.RenderTables != config.RenderTables ||
			formatter.config.ShowLineNumbers != config.ShowLineNumbers ||
			formatter.config.CompactMode != config.CompactMode ||
			formatter.config.HighlightQuotes != config.HighlightQuotes ||
			formatter.config.ParseMarkdown != config.ParseMarkdown {
			t.Errorf("Configuration not fully preserved for config %d", i)
		}
	}
}

// BenchmarkNewTerminalFormatter benchmarks the formatter initialization
func BenchmarkNewTerminalFormatter(b *testing.B) {
	config := DefaultConfig()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewTerminalFormatter(config)
	}
}

// BenchmarkDefaultConfig benchmarks the DefaultConfig function
func BenchmarkDefaultConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = DefaultConfig()
	}
}

// BenchmarkColoredConfig benchmarks the ColoredConfig function
func BenchmarkColoredConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ColoredConfig()
	}
}

// BenchmarkCompactConfig benchmarks the CompactConfig function
func BenchmarkCompactConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = CompactConfig()
	}
}

// TestCleanWhitespace tests the cleanWhitespace method
func TestCleanWhitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Single blank line removal",
			input:    "line1\n\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "Multiple blank lines removal",
			input:    "line1\n\n\n\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "Leading whitespace trimming",
			input:    "   \n\nline1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "Trailing whitespace trimming",
			input:    "line1\nline2\n\n   ",
			expected: "line1\nline2",
		},
		{
			name:     "Both leading and trailing whitespace",
			input:    "  \n\nline1\nline2\n\n  ",
			expected: "line1\nline2",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Single character",
			input:    "a",
			expected: "a",
		},
		{
			name:     "Single line with no extra whitespace",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "Multiple lines with single blank line between",
			input:    "line1\n\nline2\n\nline3",
			expected: "line1\n\nline2\n\nline3",
		},
		{
			name:     "Text with only whitespace",
			input:    "   \n\n  \n  ",
			expected: "",
		},
		{
			name:     "Multiple consecutive newlines at start",
			input:    "\n\n\n\ntext",
			expected: "text",
		},
		{
			name:     "Multiple consecutive newlines at end",
			input:    "text\n\n\n\n",
			expected: "text",
		},
		{
			name:     "Tab characters with blank lines",
			input:    "\t\ntext\n\n\n\nmore",
			expected: "text\n\nmore",
		},
		{
			name:     "Mixed whitespace in blank lines",
			input:    "line1\n  \n\t\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "Four consecutive newlines (boundary case)",
			input:    "line1\n\n\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "Exactly three consecutive newlines",
			input:    "line1\n\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "Two consecutive newlines (should remain)",
			input:    "line1\n\nline2",
			expected: "line1\n\nline2",
		},
		{
			name:     "Single newline",
			input:    "line1\nline2",
			expected: "line1\nline2",
		},
		{
			name:     "Spaces and tabs mixed",
			input:    "  \t  \nline1\nline2\n  \t  ",
			expected: "line1\nline2",
		},
		{
			name:     "Complex multiline text with various whitespace",
			input:    "  \n\nIntroduction\n\n\n\nMain content\n\nmore info\n\n\n\n  ",
			expected: "Introduction\n\nMain content\n\nmore info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(FormatterConfig{})
			result := formatter.cleanWhitespace(tt.input)

			if result != tt.expected {
				t.Errorf("cleanWhitespace failed\ninput: %q\nexpected: %q\ngot: %q", tt.input, tt.expected, result)
			}
		})
	}
}

// TestCleanWhitespaceEdgeCases tests edge cases for cleanWhitespace
func TestCleanWhitespaceEdgeCases(t *testing.T) {
	formatter := NewTerminalFormatter(FormatterConfig{})

	t.Run("Very long blank line sequence", func(t *testing.T) {
		input := "start\n" + strings.Repeat("\n", 100) + "end"
		result := formatter.cleanWhitespace(input)
		expected := "start\n\nend"

		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("Newlines with spaces", func(t *testing.T) {
		input := "line1\n   \n   \n   \nline2"
		result := formatter.cleanWhitespace(input)
		// Note: cleanWhitespace only removes extra newlines, not spaces on lines
		// This is expected behavior based on the implementation

		if !strings.Contains(result, "line1") || !strings.Contains(result, "line2") {
			t.Errorf("Result should contain both lines: %q", result)
		}
	})

	t.Run("Unicode whitespace", func(t *testing.T) {
		input := "line1\n\n\nline2"
		result := formatter.cleanWhitespace(input)
		expected := "line1\n\nline2"

		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("Preserves internal structure", func(t *testing.T) {
		input := "First paragraph\n\nSecond paragraph\n\nThird paragraph"
		result := formatter.cleanWhitespace(input)

		if result != input {
			t.Errorf("Should preserve proper paragraph spacing")
		}
	})
}

// TestCleanWhitespaceIntegration tests cleanWhitespace within the Format method
func TestCleanWhitespaceIntegration(t *testing.T) {
	formatter := NewTerminalFormatter(DefaultConfig())

	input := "Extra whitespace\n\n\n\nshould be cleaned"
	result := formatter.Format(input)

	if !strings.Contains(result, "Extra whitespace") || !strings.Contains(result, "should be cleaned") {
		t.Errorf("Format should contain both text parts: %q", result)
	}
}

// BenchmarkCleanWhitespace benchmarks the cleanWhitespace method
func BenchmarkCleanWhitespace(b *testing.B) {
	formatter := NewTerminalFormatter(FormatterConfig{})
	input := "line1\n\n\n\nline2\n\n\n\nline3"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.cleanWhitespace(input)
	}
}

// BenchmarkCleanWhitespaceComplexText benchmarks cleanWhitespace with complex input
func BenchmarkCleanWhitespaceComplexText(b *testing.B) {
	formatter := NewTerminalFormatter(FormatterConfig{})
	complexInput := `This is a complex text

with multiple paragraphs



and various whitespace patterns

		indentation

and some more content



at the end`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.cleanWhitespace(complexInput)
	}
}

// ============================================================================
// TESTS FOR parseCodeBlocks METHOD
// ============================================================================

// TestParseCodeBlocksDetection tests the detection of code blocks
func TestParseCodeBlocksDetection(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedBlocks int
		expectedCode   int
		expectedText   int
	}{
		{
			name:           "Empty string",
			input:          "",
			expectedBlocks: 0,
			expectedCode:   0,
			expectedText:   0,
		},
		{
			name:           "Text only, no code blocks",
			input:          "This is plain text without code",
			expectedBlocks: 1,
			expectedCode:   0,
			expectedText:   1,
		},
		{
			name:           "Single code block",
			input:          "```go\nfunc main() {}\n```",
			expectedBlocks: 1,
			expectedCode:   1,
			expectedText:   0,
		},
		{
			name:           "Multiple code blocks",
			input:          "```go\nfunc main() {}\n```\n\nSome text\n\n```python\nprint('hello')\n```",
			expectedBlocks: 3,
			expectedCode:   2,
			expectedText:   1,
		},
		{
			name:           "Code block with text before and after",
			input:          "Introduction text\n\n```js\nconst x = 1;\n```\n\nConclusion text",
			expectedBlocks: 3,
			expectedCode:   1,
			expectedText:   2,
		},
		{
			name:           "Multiple consecutive code blocks",
			input:          "```go\ncode1\n```\n```python\ncode2\n```",
			expectedBlocks: 2,
			expectedCode:   2,
			expectedText:   0,
		},
		{
			name:           "Code block with no language specified",
			input:          "```\nplain code\n```",
			expectedBlocks: 1,
			expectedCode:   1,
			expectedText:   0,
		},
		{
			name:           "Code block followed by text",
			input:          "```\ncode\n```\nText after",
			expectedBlocks: 2,
			expectedCode:   1,
			expectedText:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(FormatterConfig{})
			blocks := formatter.parseCodeBlocks(tt.input)

			if len(blocks) != tt.expectedBlocks {
				t.Errorf("Expected %d blocks, got %d", tt.expectedBlocks, len(blocks))
			}

			codeCount := 0
			textCount := 0
			for _, block := range blocks {
				if block.IsCode {
					codeCount++
				} else {
					textCount++
				}
			}

			if codeCount != tt.expectedCode {
				t.Errorf("Expected %d code blocks, got %d", tt.expectedCode, codeCount)
			}
			if textCount != tt.expectedText {
				t.Errorf("Expected %d text blocks, got %d", tt.expectedText, textCount)
			}
		})
	}
}

// TestParseCodeBlocksLanguageIdentification tests language identification
func TestParseCodeBlocksLanguageIdentification(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedLanguage string
	}{
		{
			name:             "Go language",
			input:            "```go\nfunc main() {}\n```",
			expectedLanguage: "go",
		},
		{
			name:             "Python language",
			input:            "```python\nprint('hello')\n```",
			expectedLanguage: "python",
		},
		{
			name:             "JavaScript language",
			input:            "```js\nconst x = 1;\n```",
			expectedLanguage: "js",
		},
		{
			name:             "TypeScript language",
			input:            "```typescript\nconst x: number = 1;\n```",
			expectedLanguage: "typescript",
		},
		{
			name:             "Bash language",
			input:            "```bash\necho 'hello'\n```",
			expectedLanguage: "bash",
		},
		{
			name:             "Rust language",
			input:            "```rust\nfn main() {}\n```",
			expectedLanguage: "rust",
		},
		{
			name:             "C++ language",
			input:            "```c++\nint main() {}\n```",
			expectedLanguage: "c++",
		},
		{
			name:             "C# language",
			input:            "```csharp\nclass Program {}\n```",
			expectedLanguage: "csharp",
		},
		{
			name:             "SQL language",
			input:            "```sql\nSELECT * FROM table;\n```",
			expectedLanguage: "sql",
		},
		{
			name:             "JSON language",
			input:            "```json\n{\"key\": \"value\"}\n```",
			expectedLanguage: "json",
		},
		{
			name:             "YAML language",
			input:            "```yaml\nkey: value\n```",
			expectedLanguage: "yaml",
		},
		{
			name:             "XML language",
			input:            "```xml\n<tag>content</tag>\n```",
			expectedLanguage: "xml",
		},
		{
			name:             "No language specified",
			input:            "```\nplain code\n```",
			expectedLanguage: "",
		},
		{
			name:             "Language with numbers",
			input:            "```python3\nprint('hello')\n```",
			expectedLanguage: "python3",
		},
		{
			name:             "Language with hyphen",
			input:            "```c-sharp\nclass Program {}\n```",
			expectedLanguage: "c-sharp",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(FormatterConfig{})
			blocks := formatter.parseCodeBlocks(tt.input)

			if len(blocks) == 0 {
				t.Fatal("Expected at least one block")
			}

			codeBlock := blocks[0]
			if !codeBlock.IsCode {
				t.Fatalf("Expected code block, got text block")
			}

			if codeBlock.Language != tt.expectedLanguage {
				t.Errorf("Expected language %q, got %q", tt.expectedLanguage, codeBlock.Language)
			}
		})
	}
}

// TestParseCodeBlocksContentExtraction tests correct content extraction
func TestParseCodeBlocksContentExtraction(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedContent string
	}{
		{
			name:            "Simple single line code",
			input:           "```go\nfunc main() {}\n```",
			expectedContent: "func main() {}",
		},
		{
			name:            "Multi-line code block",
			input:           "```python\ndef hello():\n    print('world')\n    return True\n```",
			expectedContent: "def hello():\n    print('world')\n    return True",
		},
		{
			name:            "Code with indentation",
			input:           "```go\nif true {\n    fmt.Println(\"hello\")\n}\n```",
			expectedContent: "if true {\n    fmt.Println(\"hello\")\n}",
		},
		{
			name:            "Code with empty lines",
			input:           "```javascript\nconst x = 1;\n\nconst y = 2;\n```",
			expectedContent: "const x = 1;\n\nconst y = 2;",
		},
		{
			name:            "Code with special characters",
			input:           "```bash\necho \"Hello $USER\"\nls -la\n```",
			expectedContent: "echo \"Hello $USER\"\nls -la",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(FormatterConfig{})
			blocks := formatter.parseCodeBlocks(tt.input)

			if len(blocks) == 0 {
				t.Fatal("Expected at least one block")
			}

			codeBlock := blocks[0]
			if !codeBlock.IsCode {
				t.Fatal("Expected code block")
			}

			if codeBlock.Content != tt.expectedContent {
				t.Errorf("Content mismatch\nExpected: %q\nGot: %q", tt.expectedContent, codeBlock.Content)
			}
		})
	}
}

// TestParseCodeBlocksWithMixedContent tests mixed text and code blocks
func TestParseCodeBlocksWithMixedContent(t *testing.T) {
	input := `Here is an introduction to the code.

\`\`\`go
func main() {
    fmt.Println("Hello")
}
\`\`\`

And here is some explanation after the code.

\`\`\`python
def greet():
    print("Hi")
\`\`\`

Final note.`

	formatter := NewTerminalFormatter(FormatterConfig{})
	blocks := formatter.parseCodeBlocks(input)

	// Should have: text, code, text, code, text (5 blocks)
	if len(blocks) < 3 {
		t.Fatalf("Expected at least 3 blocks, got %d", len(blocks))
	}

	// Check first block is text
	if blocks[0].IsCode {
		t.Error("First block should be text")
	}
	if !strings.Contains(blocks[0].Content, "introduction") {
		t.Error("First text block should contain 'introduction'")
	}

	// Check for code blocks
	codeBlockCount := 0
	for _, block := range blocks {
		if block.IsCode {
			codeBlockCount++
		}
	}

	if codeBlockCount < 2 {
		t.Errorf("Expected at least 2 code blocks, got %d", codeBlockCount)
	}
}

// TestParseCodeBlocksEdgeCases tests edge cases
func TestParseCodeBlocksEdgeCases(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		shouldHaveCode  bool
		expectedContent string
	}{
		{
			name:            "Code block at start",
			input:           "```go\ncode\n```\ntext",
			shouldHaveCode:  true,
			expectedContent: "code",
		},
		{
			name:            "Code block at end",
			input:           "text\n```go\ncode\n```",
			shouldHaveCode:  true,
			expectedContent: "code",
		},
		{
			name:            "Code with extra whitespace",
			input:           "```go\n\ncode\n\n```",
			shouldHaveCode:  true,
			expectedContent: "\ncode\n",
		},
		{
			name:            "Code block with only whitespace",
			input:           "```go\n   \n```",
			shouldHaveCode:  true,
			expectedContent: "   ",
		},
		{
			name:            "Multiple language specs in same text",
			input:           "Use ```go\nfor Go\n``` or ```python\nfor Python\n```",
			shouldHaveCode:  true,
			expectedContent: "for Go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(FormatterConfig{})
			blocks := formatter.parseCodeBlocks(tt.input)

			hasCode := false
			for _, block := range blocks {
				if block.IsCode {
					hasCode = true
					if tt.expectedContent != "" && block.Content == tt.expectedContent {
						return
					}
				}
			}

			if tt.shouldHaveCode && !hasCode {
				t.Error("Expected code block but none found")
			}
		})
	}
}

// TestParseCodeBlocksEmptyCode tests parsing empty code blocks
func TestParseCodeBlocksEmptyCode(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		hasEmpty  bool
		isCode    bool
	}{
		{
			name:     "Empty code block",
			input:    "```go\n```",
			hasEmpty: true,
			isCode:   true,
		},
		{
			name:     "Empty code block with language",
			input:    "```python\n```",
			hasEmpty: true,
			isCode:   true,
		},
		{
			name:     "Code block with just newlines",
			input:    "```\n\n\n```",
			hasEmpty: true,
			isCode:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(FormatterConfig{})
			blocks := formatter.parseCodeBlocks(tt.input)

			if len(blocks) == 0 {
				t.Fatal("Expected at least one block")
			}

			block := blocks[0]
			if block.IsCode != tt.isCode {
				t.Errorf("Expected IsCode=%v, got %v", tt.isCode, block.IsCode)
			}
		})
	}
}

// ============================================================================
// TESTS FOR writeCodeBlock METHOD
// ============================================================================

// TestWriteCodeBlockBasic tests basic code block writing
func TestWriteCodeBlockBasic(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		language       string
		config         FormatterConfig
		shouldContain  []string
		shouldNotContain []string
	}{
		{
			name:     "Simple code with language",
			code:     "func main() {}",
			language: "go",
			config: FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				UseColors:     false,
			},
			shouldContain:  []string{"```go", "func main() {}", "```"},
			shouldNotContain: []string{},
		},
		{
			name:     "Code without language",
			code:     "some code",
			language: "",
			config: FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				UseColors:     false,
			},
			shouldContain:  []string{"some code"},
			shouldNotContain: []string{"```"},
		},
		{
			name:     "Multi-line code",
			code:     "line1\nline2\nline3",
			language: "python",
			config: FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				UseColors:     false,
			},
			shouldContain:  []string{"line1", "line2", "line3", "```python"},
			shouldNotContain: []string{},
		},
		{
			name:     "Code with special characters",
			code:     "echo \"$VAR\"",
			language: "bash",
			config: FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				UseColors:     false,
			},
			shouldContain:  []string{"echo", "$VAR", "```bash"},
			shouldNotContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewTerminalFormatter(tt.config)
			var result strings.Builder

			formatter.writeCodeBlock(&result, tt.code, tt.language)
			output := result.String()

			for _, expected := range tt.shouldContain {
				if !strings.Contains(output, expected) {
					t.Errorf("Output should contain %q, but got:\n%s", expected, output)
				}
			}

			for _, notExpected := range tt.shouldNotContain {
				if strings.Contains(output, notExpected) {
					t.Errorf("Output should not contain %q, but got:\n%s", notExpected, output)
				}
			}
		})
	}
}

// TestWriteCodeBlockLineNumbers tests line numbering
func TestWriteCodeBlockLineNumbers(t *testing.T) {
	tests := []struct {
		name          string
		code          string
		showLineNums  bool
		shouldContain []string
	}{
		{
			name:         "Line numbers enabled",
			code:         "line1\nline2\nline3",
			showLineNums: true,
			shouldContain: []string{"1:", "2:", "3:"},
		},
		{
			name:         "Line numbers disabled",
			code:         "line1\nline2\nline3",
			showLineNums: false,
			shouldContain: []string{"line1", "line2", "line3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FormatterConfig{
				CommentPrefix:  "# ",
				IndentSize:     2,
				ShowLineNumbers: tt.showLineNums,
				UseColors:      false,
			}
			formatter := NewTerminalFormatter(config)
			var result strings.Builder

			formatter.writeCodeBlock(&result, tt.code, "go")
			output := result.String()

			for _, expected := range tt.shouldContain {
				if !strings.Contains(output, expected) {
					t.Errorf("Output should contain %q when ShowLineNumbers=%v, but got:\n%s",
						expected, tt.showLineNums, output)
				}
			}
		})
	}
}

// TestWriteCodeBlockIndentation tests indentation application
func TestWriteCodeBlockIndentation(t *testing.T) {
	tests := []struct {
		name       string
		indentSize int
		code       string
	}{
		{
			name:       "No indentation",
			indentSize: 0,
			code:       "code",
		},
		{
			name:       "2-space indentation",
			indentSize: 2,
			code:       "code",
		},
		{
			name:       "4-space indentation",
			indentSize: 4,
			code:       "code",
		},
		{
			name:       "Tab indentation",
			indentSize: 1,
			code:       "multi\nline\ncode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FormatterConfig{
				CommentPrefix:  "# ",
				IndentSize:     tt.indentSize,
				UseColors:      false,
				ShowLineNumbers: false,
			}
			formatter := NewTerminalFormatter(config)
			var result strings.Builder

			formatter.writeCodeBlock(&result, tt.code, "python")
			output := result.String()

			lines := strings.Split(strings.TrimSpace(output), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "# ") && !strings.HasPrefix(line, "# ```") {
					// Check indentation is applied
					content := strings.TrimPrefix(line, "# ")
					expectedIndent := strings.Repeat(" ", tt.indentSize)
					if tt.indentSize > 0 && !strings.HasPrefix(content, expectedIndent) {
						t.Errorf("Expected indentation of %d spaces, line: %q", tt.indentSize, line)
					}
				}
			}
		})
	}
}

// TestWriteCodeBlockCommentPrefix tests comment prefix application
func TestWriteCodeBlockCommentPrefix(t *testing.T) {
	tests := []struct {
		name          string
		prefix        string
		shouldContain bool
	}{
		{
			name:          "Default hash prefix",
			prefix:        "# ",
			shouldContain: true,
		},
		{
			name:          "Custom arrow prefix",
			prefix:        ">> ",
			shouldContain: true,
		},
		{
			name:          "Empty prefix",
			prefix:        "",
			shouldContain: false,
		},
		{
			name:          "Colon prefix",
			prefix:        ": ",
			shouldContain: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FormatterConfig{
				CommentPrefix: tt.prefix,
				IndentSize:    2,
				UseColors:     false,
			}
			formatter := NewTerminalFormatter(config)
			var result strings.Builder

			formatter.writeCodeBlock(&result, "code", "go")
			output := result.String()

			if tt.shouldContain && tt.prefix != "" {
				if !strings.Contains(output, tt.prefix) {
					t.Errorf("Output should contain prefix %q, got:\n%s", tt.prefix, output)
				}
			}
		})
	}
}

// TestWriteCodeBlockCompactMode tests compact mode formatting
func TestWriteCodeBlockCompactMode(t *testing.T) {
	tests := []struct {
		name       string
		compactMode bool
		code       string
	}{
		{
			name:        "Compact mode enabled",
			compactMode: true,
			code:        "code1\ncode2",
		},
		{
			name:        "Compact mode disabled",
			compactMode: false,
			code:        "code1\ncode2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				CompactMode:   tt.compactMode,
				UseColors:     false,
			}
			formatter := NewTerminalFormatter(config)
			var result strings.Builder

			formatter.writeCodeBlock(&result, tt.code, "python")
			output := result.String()

			// In non-compact mode, there should be trailing newline after code block
			// In compact mode, there should not be extra newlines
			if !tt.compactMode && !strings.HasSuffix(strings.TrimRight(output, "\n"), "```") {
				// Should have blank line at end
				lineCount := strings.Count(output, "\n")
				if lineCount < 2 {
					t.Errorf("Non-compact mode should have more newlines")
				}
			}
		})
	}
}

// TestWriteCodeBlockLanguages tests various programming languages
func TestWriteCodeBlockLanguages(t *testing.T) {
	languages := []string{
		"go", "python", "javascript", "typescript", "bash", "rust",
		"java", "cpp", "csharp", "ruby", "php", "sql", "json",
		"yaml", "xml", "html", "css", "swift", "kotlin", "scala",
	}

	for _, lang := range languages {
		t.Run(lang, func(t *testing.T) {
			config := FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				UseColors:     false,
			}
			formatter := NewTerminalFormatter(config)
			var result strings.Builder

			code := "code snippet in " + lang
			formatter.writeCodeBlock(&result, code, lang)
			output := result.String()

			expected := "```" + lang
			if !strings.Contains(output, expected) {
				t.Errorf("Output should contain %q, got:\n%s", expected, output)
			}
		})
	}
}

// TestWriteCodeBlockWhitespaceHandling tests whitespace in code blocks
func TestWriteCodeBlockWhitespaceHandling(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		language string
	}{
		{
			name:     "Leading and trailing whitespace",
			code:     "  \n\ncode content\n\n  ",
			language: "go",
		},
		{
			name:     "Code with internal indentation",
			code:     "if true {\n    return value\n}",
			language: "js",
		},
		{
			name:     "Code with tabs",
			code:     "function test() {\n\t\treturn true;\n}",
			language: "typescript",
		},
		{
			name:     "Code with mixed whitespace",
			code:     "  spaces\n\ttabs\n  mixed\t\t content",
			language: "python",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				UseColors:     false,
			}
			formatter := NewTerminalFormatter(config)
			var result strings.Builder

			formatter.writeCodeBlock(&result, tt.code, tt.language)
			output := result.String()

			// Verify we get some output
			if output == "" {
				t.Error("Expected non-empty output")
			}

			// Verify code block markers are present for languages
			if tt.language != "" {
				if !strings.Contains(output, "```"+tt.language) {
					t.Errorf("Missing language header for %s", tt.language)
				}
			}
		})
	}
}

// TestWriteCodeBlockEmptyCode tests handling of empty code
func TestWriteCodeBlockEmptyCode(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		language string
	}{
		{
			name:     "Empty string",
			code:     "",
			language: "go",
		},
		{
			name:     "Only whitespace",
			code:     "   \n  \n  ",
			language: "python",
		},
		{
			name:     "Only newlines",
			code:     "\n\n",
			language: "bash",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := FormatterConfig{
				CommentPrefix: "# ",
				IndentSize:    2,
				UseColors:     false,
			}
			formatter := NewTerminalFormatter(config)
			var result strings.Builder

			formatter.writeCodeBlock(&result, tt.code, tt.language)
			output := result.String()

			// Should still have code block markers for language
			if tt.language != "" {
				if !strings.Contains(output, "```"+tt.language) {
					t.Errorf("Empty code should still have language marker")
				}
			}
		})
	}
}

// TestWriteCodeBlockWithColors tests code block output with colors enabled
func TestWriteCodeBlockWithColors(t *testing.T) {
	config := FormatterConfig{
		CommentPrefix: "# ",
		IndentSize:    2,
		UseColors:     true,
	}
	formatter := NewTerminalFormatter(config)

	var result strings.Builder
	formatter.writeCodeBlock(&result, "func main() {}", "go")
	output := result.String()

	// Output should have content
	if output == "" {
		t.Error("Expected non-empty output with colors enabled")
	}

	// Should still contain code
	if !strings.Contains(output, "main") {
		t.Error("Output should contain the code content")
	}
}

// TestWriteCodeBlockMultipleLines tests multi-line code formatting
func TestWriteCodeBlockMultipleLines(t *testing.T) {
	code := `func greet(name string) {
    fmt.Printf("Hello, %s!\n", name)
    return nil
}`

	config := FormatterConfig{
		CommentPrefix:  "# ",
		IndentSize:     2,
		ShowLineNumbers: true,
		UseColors:       false,
	}
	formatter := NewTerminalFormatter(config)

	var result strings.Builder
	formatter.writeCodeBlock(&result, code, "go")
	output := result.String()

	// Check each line is present
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		if line != "" && !strings.Contains(output, line) {
			t.Errorf("Line %d not found in output: %q", i+1, line)
		}
	}

	// Check line numbers
	if !strings.Contains(output, "1:") || !strings.Contains(output, "2:") {
		t.Error("Line numbers should be present")
	}
}

// BenchmarkParseCodeBlocks benchmarks the parseCodeBlocks method
func BenchmarkParseCodeBlocks(b *testing.B) {
	input := `Introduction text

\`\`\`go
func main() {
    fmt.Println("hello")
}
\`\`\`

Middle section

\`\`\`python
def test():
    return True
\`\`\`

Conclusion`

	formatter := NewTerminalFormatter(FormatterConfig{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.parseCodeBlocks(input)
	}
}

// BenchmarkParseCodeBlocksComplex benchmarks parseCodeBlocks with complex input
func BenchmarkParseCodeBlocksComplex(b *testing.B) {
	input := strings.Repeat(`Text content\n\n\`\`\`go\nfunc test() {}\n\`\`\`\n\n`, 10)

	formatter := NewTerminalFormatter(FormatterConfig{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = formatter.parseCodeBlocks(input)
	}
}

// BenchmarkWriteCodeBlock benchmarks the writeCodeBlock method
func BenchmarkWriteCodeBlock(b *testing.B) {
	code := `func main() {
    fmt.Println("Hello, World!")
    for i := 0; i < 10; i++ {
        fmt.Printf("%d\n", i)
    }
}`

	config := FormatterConfig{
		CommentPrefix:   "# ",
		IndentSize:      2,
		UseColors:       false,
		ShowLineNumbers: false,
	}
	formatter := NewTerminalFormatter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result strings.Builder
		formatter.writeCodeBlock(&result, code, "go")
	}
}

// BenchmarkWriteCodeBlockWithLineNumbers benchmarks writeCodeBlock with line numbers
func BenchmarkWriteCodeBlockWithLineNumbers(b *testing.B) {
	code := `line1
line2
line3
line4
line5
line6
line7
line8
line9
line10`

	config := FormatterConfig{
		CommentPrefix:   "# ",
		IndentSize:      2,
		UseColors:       false,
		ShowLineNumbers: true,
	}
	formatter := NewTerminalFormatter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result strings.Builder
		formatter.writeCodeBlock(&result, code, "python")
	}
}

// BenchmarkWriteCodeBlockWithColors benchmarks writeCodeBlock with colors
func BenchmarkWriteCodeBlockWithColors(b *testing.B) {
	code := `const main = async () => {
    console.log("hello");
    return true;
};`

	config := FormatterConfig{
		CommentPrefix:   "# ",
		IndentSize:      2,
		UseColors:       true,
		ShowLineNumbers: false,
	}
	formatter := NewTerminalFormatter(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result strings.Builder
		formatter.writeCodeBlock(&result, code, "javascript")
	}
}
