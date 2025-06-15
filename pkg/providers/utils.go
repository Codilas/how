package providers

import (
	"fmt"
	"strings"
)

// CommandCategory represents a category of shell commands
type CommandCategory string

const (
	CategoryFile    CommandCategory = "file"
	CategoryNetwork CommandCategory = "network"
	CategorySystem  CommandCategory = "system"
	CategoryGit     CommandCategory = "git"
	CategoryPackage CommandCategory = "package"
	CategoryBuild   CommandCategory = "build"
	CategoryGeneral CommandCategory = "general"
)

// IsCommandSafe determines if a command is safe to auto-execute
func IsCommandSafe(command string) bool {
	dangerousCommands := []string{
		"rm", "sudo", "chmod", "mv", "dd", "mkfs", "fdisk",
	}
	cmd := strings.Fields(command)[0]

	for _, dangerous := range dangerousCommands {
		if cmd == dangerous {
			return false
		}
	}

	return true
}

// CategorizeCommand categorizes a command into a specific category
func CategorizeCommand(command string) CommandCategory {
	cmd := strings.Fields(command)[0]

	categories := map[string]CommandCategory{
		"ls":   CategoryFile,
		"cat":  CategoryFile,
		"grep": CategoryFile,
		"find": CategoryFile,
		"curl": CategoryNetwork,
		"wget": CategoryNetwork,
		"ping": CategoryNetwork,
		"ps":   CategorySystem,
		"top":  CategorySystem,
		"df":   CategorySystem,
		"free": CategorySystem,
		"git":  CategoryGit,
		"npm":  CategoryPackage,
		"pip":  CategoryPackage,
		"go":   CategoryBuild,
	}

	if category, exists := categories[cmd]; exists {
		return category
	}

	return CategoryGeneral
}

// ExtractCodeBlocks finds code blocks in the response text
func ExtractCodeBlocks(text string) []CodeBlock {
	var blocks []CodeBlock

	// Simple regex-based extraction for markdown code blocks
	// TODO: use a proper markdown parser
	lines := strings.Split(text, "\n")
	var currentBlock *CodeBlock

	for _, line := range lines {
		if strings.HasPrefix(line, "```") {
			if currentBlock == nil {
				// Start of code block
				language := strings.TrimSpace(strings.TrimPrefix(line, "```"))
				currentBlock = &CodeBlock{
					Language: language,
					Code:     "",
				}
			} else {
				// End of code block
				blocks = append(blocks, *currentBlock)
				currentBlock = nil
			}
		} else if currentBlock != nil {
			currentBlock.Code += line + "\n"
		}
	}

	fmt.Println("\n---------------------------\n")
	fmt.Printf("Extracted code blocks: %+v", blocks)
	fmt.Println("\n---------------------------\n")

	return blocks
}

// ExtractCommands finds command suggestions in the response text
func ExtractCommands(text string) []SuggestedCommand {
	var commands []SuggestedCommand

	// TODO: Simple pattern matching for common command patterns we should extract
	lines := strings.Split(text, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for lines that start with $ or contain common command patterns
		if strings.HasPrefix(line, "$ ") {
			cmd := strings.TrimPrefix(line, "$ ")
			commands = append(commands, SuggestedCommand{
				Command:     cmd,
				Description: "Suggested command",
				Safe:        IsCommandSafe(cmd),
				Category:    string(CategorizeCommand(cmd)),
			})
		}
	}

	fmt.Println("\n---------------------------\n")
	fmt.Printf("Extracted commands: %+v", commands)
	fmt.Println("\n---------------------------\n")

	return commands
}
