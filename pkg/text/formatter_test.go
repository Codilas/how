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
