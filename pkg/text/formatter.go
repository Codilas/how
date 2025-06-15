package text

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

// FormatterConfig controls how text is formatted for terminal output
type FormatterConfig struct {
	// Basic formatting options
	UseColors     bool   `json:"use_colors"`
	CommentPrefix string `json:"comment_prefix"`
	LineWidth     int    `json:"line_width"`
	IndentSize    int    `json:"indent_size"`

	// Advanced formatting options
	UseBoxes        bool `json:"use_boxes"`
	UseBullets      bool `json:"use_bullets"`
	HighlightCode   bool `json:"highlight_code"`
	WrapLongLines   bool `json:"wrap_long_lines"`
	RenderTables    bool `json:"render_tables"`
	ShowLineNumbers bool `json:"show_line_numbers"`
	CompactMode     bool `json:"compact_mode"`
	HighlightQuotes bool `json:"highlight_quotes"`
	ParseMarkdown   bool `json:"parse_markdown"`
}

// TerminalFormatter formats LLM response text for terminal display
type TerminalFormatter struct {
	config FormatterConfig

	// Compiled regexes for performance
	structuredRegex *regexp.Regexp
	codeBlockRegex  *regexp.Regexp
	inlineCodeRegex *regexp.Regexp
	markdownH1Regex *regexp.Regexp
	markdownH2Regex *regexp.Regexp
	markdownH3Regex *regexp.Regexp
	boldTextRegex   *regexp.Regexp
	italicTextRegex *regexp.Regexp
	linkRegex       *regexp.Regexp
	blockQuoteRegex *regexp.Regexp
	tableRowRegex   *regexp.Regexp
	listItemRegex   *regexp.Regexp

	// Color definitions
	colors struct {
		heading1   *color.Color
		heading2   *color.Color
		heading3   *color.Color
		code       *color.Color
		inlineCode *color.Color
		quote      *color.Color
		bullet     *color.Color
		text       *color.Color
		comment    *color.Color
		border     *color.Color
		link       *color.Color
		bold       *color.Color
		italic     *color.Color
	}
}

// NewTerminalFormatter creates a new terminal formatter with compiled regexes
func NewTerminalFormatter(config FormatterConfig) *TerminalFormatter {
	// Set defaults if not provided
	if config.CommentPrefix == "" {
		config.CommentPrefix = "# "
	}
	if config.LineWidth == 0 {
		config.LineWidth = 80
	}
	if config.IndentSize == 0 {
		config.IndentSize = 2
	}

	f := &TerminalFormatter{
		config: config,
	}

	// Compile regexes once for performance
	f.compileRegexes()

	// Initialize colors
	f.initColors()

	return f
}

// compileRegexes compiles all regular expressions used by the formatter
func (f *TerminalFormatter) compileRegexes() {
	f.structuredRegex = regexp.MustCompile(`(?s)<structured_commands>.*?</structured_commands>`)
	f.codeBlockRegex = regexp.MustCompile("(?s)```([a-zA-Z0-9+-]*)\n?(.*?)\n?```")
	f.inlineCodeRegex = regexp.MustCompile("`([^`\n]+)`")
	f.markdownH1Regex = regexp.MustCompile(`^# (.+)$`)
	f.markdownH2Regex = regexp.MustCompile(`^## (.+)$`)
	f.markdownH3Regex = regexp.MustCompile(`^### (.+)$`)
	f.boldTextRegex = regexp.MustCompile(`\*\*([^\*\n]+)\*\*`)
	f.italicTextRegex = regexp.MustCompile(`\*([^\*\n]+)\*`)
	f.linkRegex = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	f.blockQuoteRegex = regexp.MustCompile(`^> (.+)$`)
	f.tableRowRegex = regexp.MustCompile(`^\|(.+)\|$`)
	f.listItemRegex = regexp.MustCompile(`^(\s*)([-*+]|\d+\.)\s+(.+)$`)
}

// initColors initializes color schemes
func (f *TerminalFormatter) initColors() {
	if !f.config.UseColors {
		return
	}

	f.colors.heading1 = color.New(color.FgYellow, color.Bold, color.Underline)
	f.colors.heading2 = color.New(color.FgYellow, color.Bold)
	f.colors.heading3 = color.New(color.FgYellow)
	f.colors.code = color.New(color.FgGreen)
	f.colors.inlineCode = color.New(color.FgGreen, color.Bold)
	f.colors.quote = color.New(color.FgBlue, color.Italic)
	f.colors.bullet = color.New(color.FgCyan, color.Bold)
	f.colors.text = color.New(color.FgWhite)
	f.colors.comment = color.New(color.FgHiBlack)
	f.colors.border = color.New(color.FgCyan, color.Bold)
	f.colors.link = color.New(color.FgBlue, color.Underline)
	f.colors.bold = color.New(color.Bold)
	f.colors.italic = color.New(color.Italic)
}

// Format formats the raw LLM response text for terminal display
func (f *TerminalFormatter) Format(text string) string {
	// Step 1: Remove structured commands section
	cleanText := f.structuredRegex.ReplaceAllString(text, "")

	// Step 2: Clean up extra whitespace
	cleanText = f.cleanWhitespace(cleanText)

	// Step 3: Handle code blocks first (they need special processing)
	blocks := f.parseCodeBlocks(cleanText)

	// Step 4: Apply formatting
	var result strings.Builder

	if f.config.UseBoxes {
		f.writeHeader(&result)
	}

	for _, block := range blocks {
		if block.IsCode {
			f.writeCodeBlock(&result, block.Content, block.Language)
		} else {
			f.writeTextBlock(&result, block.Content)
		}
	}

	return result.String()
}

// ContentBlock represents a block of content (either code or text)
type ContentBlock struct {
	Content  string
	IsCode   bool
	Language string
}

// parseCodeBlocks separates text into code and non-code blocks
func (f *TerminalFormatter) parseCodeBlocks(text string) []ContentBlock {
	var blocks []ContentBlock

	parts := f.codeBlockRegex.Split(text, -1)
	matches := f.codeBlockRegex.FindAllStringSubmatch(text, -1)

	for i, part := range parts {
		if strings.TrimSpace(part) != "" {
			blocks = append(blocks, ContentBlock{
				Content: part,
				IsCode:  false,
			})
		}

		if i < len(matches) {
			blocks = append(blocks, ContentBlock{
				Content:  matches[i][2],
				IsCode:   true,
				Language: matches[i][1],
			})
		}
	}

	return blocks
}

// writeHeader writes a formatted header box
func (f *TerminalFormatter) writeHeader(result *strings.Builder) {
	title := "AI Assistant Response"
	width := f.config.LineWidth
	if width > 80 {
		width = 80
	}

	if f.config.UseColors {
		border := strings.Repeat("═", width-len(f.config.CommentPrefix))
		result.WriteString(f.config.CommentPrefix + f.colors.border.Sprint(border) + "\n")

		padding := (width - len(f.config.CommentPrefix) - len(title)) / 2
		result.WriteString(f.config.CommentPrefix + strings.Repeat(" ", padding) +
			f.colors.heading1.Sprint(title) + "\n")

		result.WriteString(f.config.CommentPrefix + f.colors.border.Sprint(border) + "\n\n")
	} else {
		border := strings.Repeat("=", width-len(f.config.CommentPrefix))
		result.WriteString(f.config.CommentPrefix + border + "\n")

		padding := (width - len(f.config.CommentPrefix) - len(title)) / 2
		result.WriteString(f.config.CommentPrefix + strings.Repeat(" ", padding) +
			strings.ToUpper(title) + "\n")

		result.WriteString(f.config.CommentPrefix + border + "\n\n")
	}
}

// writeCodeBlock formats and writes a code block
func (f *TerminalFormatter) writeCodeBlock(result *strings.Builder, code, language string) {
	lines := strings.Split(strings.TrimSpace(code), "\n")

	// Write code block header
	if language != "" {
		header := fmt.Sprintf("```%s", language)
		if f.config.UseColors {
			result.WriteString(f.config.CommentPrefix + f.colors.comment.Sprint(header) + "\n")
		} else {
			result.WriteString(f.config.CommentPrefix + header + "\n")
		}
	}

	// Write code lines with optional line numbers
	for i, line := range lines {
		indent := strings.Repeat(" ", f.config.IndentSize)

		var formattedLine string
		if f.config.ShowLineNumbers {
			lineNum := fmt.Sprintf("%3d: ", i+1)
			if f.config.UseColors {
				formattedLine = f.config.CommentPrefix + indent +
					f.colors.comment.Sprint(lineNum) + f.colors.code.Sprint(line)
			} else {
				formattedLine = f.config.CommentPrefix + indent + lineNum + line
			}
		} else {
			if f.config.UseColors {
				formattedLine = f.config.CommentPrefix + indent + f.colors.code.Sprint(line)
			} else {
				formattedLine = f.config.CommentPrefix + indent + line
			}
		}

		result.WriteString(formattedLine + "\n")
	}

	// Write code block footer
	if language != "" {
		footer := "```"
		if f.config.UseColors {
			result.WriteString(f.config.CommentPrefix + f.colors.comment.Sprint(footer) + "\n")
		} else {
			result.WriteString(f.config.CommentPrefix + footer + "\n")
		}
	}

	if !f.config.CompactMode {
		result.WriteString("\n")
	}
}

// writeTextBlock formats and writes a text block
func (f *TerminalFormatter) writeTextBlock(result *strings.Builder, text string) {
	lines := strings.Split(text, "\n")

	for i, line := range lines {
		line = strings.TrimRightFunc(line, unicode.IsSpace)

		if line == "" {
			if !f.config.CompactMode {
				result.WriteString("\n")
			}
			continue
		}

		// Check for tables
		if f.config.RenderTables && f.isTableRow(line) {
			tableLines := f.extractTable(lines, i)
			f.writeTable(result, tableLines)
			// Skip the processed table lines
			i += len(tableLines) - 1
			continue
		}

		formattedLine := f.formatTextLine(line)
		result.WriteString(formattedLine)

		if i < len(lines)-1 || !f.config.CompactMode {
			result.WriteString("\n")
		}
	}
}

// formatTextLine formats a single line of text with markdown support
func (f *TerminalFormatter) formatTextLine(line string) string {
	// Handle markdown headers
	if f.config.ParseMarkdown {
		if matches := f.markdownH1Regex.FindStringSubmatch(line); matches != nil {
			return f.formatHeading(matches[1], 1)
		}
		if matches := f.markdownH2Regex.FindStringSubmatch(line); matches != nil {
			return f.formatHeading(matches[1], 2)
		}
		if matches := f.markdownH3Regex.FindStringSubmatch(line); matches != nil {
			return f.formatHeading(matches[1], 3)
		}
	}

	// Handle block quotes
	if f.config.HighlightQuotes {
		if matches := f.blockQuoteRegex.FindStringSubmatch(line); matches != nil {
			return f.formatBlockQuote(matches[1])
		}
	}

	// Handle list items
	if f.config.UseBullets {
		if matches := f.listItemRegex.FindStringSubmatch(line); matches != nil {
			return f.formatListItem(matches[1], matches[2], matches[3])
		}
	}

	// Handle inline formatting
	if f.config.ParseMarkdown && f.config.UseColors {
		line = f.applyInlineFormatting(line)
	}

	// Handle regular text with potential inline code
	if f.config.HighlightCode {
		line = f.highlightInlineCode(line)
	}

	// Wrap long lines if needed
	formattedLine := f.config.CommentPrefix + line
	if f.config.WrapLongLines && len(formattedLine) > f.config.LineWidth {
		return f.wrapLine(formattedLine)
	}

	return formattedLine
}

// formatHeading formats markdown-style headings
func (f *TerminalFormatter) formatHeading(text string, level int) string {
	if !f.config.UseColors {
		return f.config.CommentPrefix + strings.ToUpper(text)
	}

	switch level {
	case 1:
		return f.config.CommentPrefix + f.colors.heading1.Sprint(text)
	case 2:
		return f.config.CommentPrefix + f.colors.heading2.Sprint(text)
	case 3:
		return f.config.CommentPrefix + f.colors.heading3.Sprint(text)
	default:
		return f.config.CommentPrefix + f.colors.heading3.Sprint(text)
	}
}

// formatBlockQuote formats block quote lines
func (f *TerminalFormatter) formatBlockQuote(text string) string {
	indent := strings.Repeat(" ", f.config.IndentSize)

	if f.config.UseColors {
		return f.config.CommentPrefix + indent + f.colors.quote.Sprint("│ "+text)
	}

	return f.config.CommentPrefix + indent + "> " + text
}

// formatListItem formats list items with proper indentation
func (f *TerminalFormatter) formatListItem(indentStr, marker, text string) string {
	baseIndent := strings.Repeat(" ", f.config.IndentSize)
	listIndent := indentStr

	var bullet string
	if strings.Contains(marker, ".") {
		// Numbered list
		bullet = marker
	} else {
		// Bullet list
		bullet = "•"
	}

	if f.config.UseColors {
		coloredBullet := f.colors.bullet.Sprint(bullet)
		coloredText := f.colors.text.Sprint(text)
		return f.config.CommentPrefix + baseIndent + listIndent + coloredBullet + " " + coloredText
	}

	return f.config.CommentPrefix + baseIndent + listIndent + bullet + " " + text
}

// applyInlineFormatting applies bold, italic, and link formatting
func (f *TerminalFormatter) applyInlineFormatting(line string) string {
	// Handle links
	line = f.linkRegex.ReplaceAllStringFunc(line, func(match string) string {
		parts := f.linkRegex.FindStringSubmatch(match)
		if len(parts) == 3 {
			if f.config.UseColors {
				return f.colors.link.Sprint(parts[1]) + f.colors.comment.Sprint(" ("+parts[2]+")")
			}
			return parts[1] + " (" + parts[2] + ")"
		}
		return match
	})

	// Handle bold text
	line = f.boldTextRegex.ReplaceAllStringFunc(line, func(match string) string {
		text := strings.Trim(match, "*")
		if f.config.UseColors {
			return f.colors.bold.Sprint(text)
		}
		return strings.ToUpper(text)
	})

	// Handle italic text
	line = f.italicTextRegex.ReplaceAllStringFunc(line, func(match string) string {
		text := strings.Trim(match, "*")
		if f.config.UseColors {
			return f.colors.italic.Sprint(text)
		}
		return "_" + text + "_"
	})

	return line
}

// highlightInlineCode highlights inline code blocks
func (f *TerminalFormatter) highlightInlineCode(line string) string {
	if !f.config.UseColors {
		return line
	}

	return f.inlineCodeRegex.ReplaceAllStringFunc(line, func(match string) string {
		code := strings.Trim(match, "`")
		return f.colors.inlineCode.Sprint("`" + code + "`")
	})
}

// Table handling methods
func (f *TerminalFormatter) isTableRow(line string) bool {
	return f.tableRowRegex.MatchString(strings.TrimSpace(line))
}

func (f *TerminalFormatter) extractTable(lines []string, startIndex int) []string {
	var tableLines []string
	for i := startIndex; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			break
		}
		if f.isTableRow(line) {
			tableLines = append(tableLines, line)
		} else {
			break
		}
	}
	return tableLines
}

func (f *TerminalFormatter) writeTable(result *strings.Builder, tableLines []string) {
	if len(tableLines) == 0 {
		return
	}

	// Parse table data
	var rows [][]string
	for _, line := range tableLines {
		cells := strings.Split(strings.Trim(line, "|"), "|")
		var cleanCells []string
		for _, cell := range cells {
			cleanCells = append(cleanCells, strings.TrimSpace(cell))
		}
		rows = append(rows, cleanCells)
	}

	// Calculate column widths
	if len(rows) == 0 {
		return
	}

	colWidths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			if i < len(colWidths) && len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Write table
	for i, row := range rows {
		var line strings.Builder
		line.WriteString(f.config.CommentPrefix)
		line.WriteString("│")

		for j, cell := range row {
			if j < len(colWidths) {
				padding := colWidths[j] - len(cell)
				line.WriteString(" " + cell + strings.Repeat(" ", padding) + " │")
			}
		}

		if f.config.UseColors {
			result.WriteString(f.colors.text.Sprint(line.String()) + "\n")
		} else {
			result.WriteString(line.String() + "\n")
		}

		// Add separator after header (first row)
		if i == 0 {
			var separator strings.Builder
			separator.WriteString(f.config.CommentPrefix)
			separator.WriteString("├")
			for j, width := range colWidths {
				separator.WriteString(strings.Repeat("─", width+2))
				if j < len(colWidths)-1 {
					separator.WriteString("┼")
				}
			}
			separator.WriteString("┤")

			if f.config.UseColors {
				result.WriteString(f.colors.border.Sprint(separator.String()) + "\n")
			} else {
				result.WriteString(separator.String() + "\n")
			}
		}
	}

	result.WriteString("\n")
}

// wrapLine wraps long lines intelligently
func (f *TerminalFormatter) wrapLine(line string) string {
	if len(line) <= f.config.LineWidth {
		return line
	}

	var result strings.Builder
	words := strings.Fields(line)
	if len(words) == 0 {
		return line
	}

	currentLine := f.config.CommentPrefix
	continuationIndent := f.config.CommentPrefix + strings.Repeat(" ", f.config.IndentSize)

	for i, word := range words {
		if i == 0 {
			currentLine += word
		} else if len(currentLine)+len(word)+1 <= f.config.LineWidth {
			currentLine += " " + word
		} else {
			result.WriteString(currentLine + "\n")
			currentLine = continuationIndent + word
		}
	}

	if len(currentLine) > len(continuationIndent) {
		result.WriteString(currentLine)
	}

	return result.String()
}

// cleanWhitespace removes extra whitespace and normalizes line breaks
func (f *TerminalFormatter) cleanWhitespace(text string) string {
	// Remove extra blank lines (more than 2 consecutive newlines)
	cleanText := regexp.MustCompile(`\n{3,}`).ReplaceAllString(text, "\n\n")

	// Trim leading and trailing whitespace
	cleanText = strings.TrimSpace(cleanText)

	return cleanText
}

// Configuration presets
func DefaultConfig() FormatterConfig {
	return FormatterConfig{
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
	}
}

func ColoredConfig() FormatterConfig {
	return FormatterConfig{
		UseColors:       true,
		CommentPrefix:   "# ",
		LineWidth:       80,
		IndentSize:      2,
		UseBoxes:        true,
		UseBullets:      true,
		HighlightCode:   true,
		WrapLongLines:   true,
		RenderTables:    true,
		ShowLineNumbers: false,
		CompactMode:     false,
		HighlightQuotes: true,
		ParseMarkdown:   true,
	}
}

func CompactConfig() FormatterConfig {
	return FormatterConfig{
		UseColors:       true,
		CommentPrefix:   "",
		LineWidth:       120,
		IndentSize:      2,
		UseBoxes:        false,
		UseBullets:      true,
		HighlightCode:   true,
		WrapLongLines:   false,
		RenderTables:    true,
		ShowLineNumbers: false,
		CompactMode:     true,
		HighlightQuotes: true,
		ParseMarkdown:   true,
	}
}
