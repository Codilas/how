package extractor

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Command represents a single command
type Command struct {
	Command     string `json:"command"`
	Description string `json:"description"`
	Order       int    `json:"order"`
	Safe        bool   `json:"safe"`
}

// Workflow represents a multi-step workflow
type Workflow struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Steps       []Command `json:"steps"`
}

// ExtractedCommands represents the final output
type ExtractedCommands struct {
	Commands  []Command  `json:"commands"`
	Workflows []Workflow `json:"workflows"`
}

// CommandExtractor extracts commands and workflows from LLM text
type CommandExtractor struct {
	structuredRegex *regexp.Regexp
	codeBlockRegex  *regexp.Regexp
}

// NewCommandExtractor creates a new command extractor
func NewCommandExtractor() *CommandExtractor {
	return &CommandExtractor{
		structuredRegex: regexp.MustCompile(`(?s)<structured_commands>\s*(.*?)\s*</structured_commands>`),
		codeBlockRegex:  regexp.MustCompile("(?s)```(?:bash|shell|sh)?\n?(.*?)\n?```"),
	}
}

// Extract extracts commands and workflows from LLM response text
func (e *CommandExtractor) Extract(text string) (*ExtractedCommands, error) {
	// Try structured extraction first
	if structured := e.extractStructured(text); structured != nil {
		return structured, nil
	}

	// Fall back to extracting from code blocks
	return e.extractFromCodeBlocks(text), nil
}

// extractStructured extracts from <structured_commands> section
func (e *CommandExtractor) extractStructured(text string) *ExtractedCommands {
	matches := e.structuredRegex.FindStringSubmatch(text)
	if len(matches) < 2 {
		return nil
	}

	jsonStr := strings.TrimSpace(matches[1])

	var result ExtractedCommands
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil
	}

	return &result
}

// extractFromCodeBlocks extracts commands from bash code blocks
func (e *CommandExtractor) extractFromCodeBlocks(text string) *ExtractedCommands {
	result := &ExtractedCommands{
		Commands:  []Command{},
		Workflows: []Workflow{},
	}

	// Find all bash code blocks
	matches := e.codeBlockRegex.FindAllStringSubmatch(text, -1)

	order := 1
	for _, match := range matches {
		if len(match) >= 2 {
			codeContent := strings.TrimSpace(match[1])
			lines := strings.Split(codeContent, "\n")

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}

				// Remove common prompt prefixes
				for _, prefix := range []string{"$ ", "# ", "> "} {
					if strings.HasPrefix(line, prefix) {
						line = strings.TrimSpace(strings.TrimPrefix(line, prefix))
						break
					}
				}

				if len(line) > 0 {
					result.Commands = append(result.Commands, Command{
						Command:     line,
						Description: "", // Let LLM provide descriptions
						Order:       order,
					})
					order++
				}
			}
		}
	}

	return result
}

// ToJSON converts the extracted commands to JSON string
func (e *ExtractedCommands) ToJSON() (string, error) {
	jsonBytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// ToJSONCompact converts the extracted commands to compact JSON string
func (e *ExtractedCommands) ToJSONCompact() (string, error) {
	jsonBytes, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// GetAllCommands returns all commands (individual + workflow steps) as a flat list
func (e *ExtractedCommands) GetAllCommands() []Command {
	var allCommands []Command

	// Add individual commands
	allCommands = append(allCommands, e.Commands...)

	// Add workflow steps
	for _, workflow := range e.Workflows {
		allCommands = append(allCommands, workflow.Steps...)
	}

	return allCommands
}

// HasCommands returns true if any commands were extracted
func (e *ExtractedCommands) HasCommands() bool {
	return len(e.Commands) > 0 || len(e.Workflows) > 0
}

// Count returns the total number of commands (individual + workflow steps)
func (e *ExtractedCommands) Count() int {
	count := len(e.Commands)
	for _, workflow := range e.Workflows {
		count += len(workflow.Steps)
	}
	return count
}
