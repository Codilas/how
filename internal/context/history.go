package context

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tzvonimir/how/pkg/providers"
)

// getRecentCommands reads recent commands from shell history
func getRecentCommands(shell string, count int) ([]providers.CommandHistory, error) {
	switch shell {
	case "bash":
		return getBashHistory(count)
	case "zsh":
		return getZshHistory(count)
	case "fish":
		return getFishHistory(count)
	default:
		return getBashHistory(count) // fallback to bash format
	}
}

// getBashHistory reads from ~/.bash_history
func getBashHistory(count int) ([]providers.CommandHistory, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyFile := filepath.Join(homeDir, ".bash_history")
	return readSimpleHistory(historyFile, count)
}

// getZshHistory reads from ~/.zsh_history (extended format)
func getZshHistory(count int) ([]providers.CommandHistory, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyFile := filepath.Join(homeDir, ".zsh_history")
	return readZshHistoryFormat(historyFile, count)
}

// getFishHistory reads from ~/.local/share/fish/fish_history
func getFishHistory(count int) ([]providers.CommandHistory, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	historyFile := filepath.Join(homeDir, ".local", "share", "fish", "fish_history")
	return readFishHistoryFormat(historyFile, count)
}

// readSimpleHistory reads bash-style history (one command per line)
func readSimpleHistory(filename string, count int) ([]providers.CommandHistory, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			lines = append(lines, line)
		}
	}

	// Get last 'count' commands
	start := len(lines) - count
	if start < 0 {
		start = 0
	}

	var commands []providers.CommandHistory
	for i := start; i < len(lines); i++ {
		commands = append(commands, providers.CommandHistory{
			Command:   lines[i],
			ExitCode:  0, // Unknown for simple format
			Timestamp: time.Now().Add(-time.Duration(len(lines)-i) * time.Minute),
		})
	}

	return commands, nil
}

// readZshHistoryFormat reads zsh extended history format
func readZshHistoryFormat(filename string, count int) ([]providers.CommandHistory, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []providers.CommandHistory
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Zsh format: : <timestamp>:<duration>;<command>
		if strings.HasPrefix(line, ": ") {
			parts := strings.SplitN(line[2:], ";", 2)
			if len(parts) == 2 {
				// Parse timestamp and duration
				timestampPart := strings.Split(parts[0], ":")[0]
				if timestamp, err := strconv.ParseInt(timestampPart, 10, 64); err == nil {
					commands = append(commands, providers.CommandHistory{
						Command:   parts[1],
						ExitCode:  0, // Not stored in zsh history
						Timestamp: time.Unix(timestamp, 0),
					})
				}
			}
		}
	}

	// Return last 'count' commands
	if len(commands) > count {
		commands = commands[len(commands)-count:]
	}

	return commands, nil
}

// readFishHistoryFormat reads fish shell history format
func readFishHistoryFormat(filename string, count int) ([]providers.CommandHistory, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var commands []providers.CommandHistory
	scanner := bufio.NewScanner(file)

	var currentCommand providers.CommandHistory
	var inCommand bool

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "- cmd: ") {
			// Save previous command if exists
			if inCommand && currentCommand.Command != "" {
				commands = append(commands, currentCommand)
			}

			// Start new command
			currentCommand = providers.CommandHistory{
				Command:  strings.TrimPrefix(line, "- cmd: "),
				ExitCode: 0,
			}
			inCommand = true

		} else if strings.HasPrefix(line, "  when: ") && inCommand {
			timestampStr := strings.TrimPrefix(line, "  when: ")
			if timestamp, err := strconv.ParseInt(timestampStr, 10, 64); err == nil {
				currentCommand.Timestamp = time.Unix(timestamp, 0)
			}
		}
	}

	// Add last command
	if inCommand && currentCommand.Command != "" {
		commands = append(commands, currentCommand)
	}

	// Return last 'count' commands
	if len(commands) > count {
		commands = commands[len(commands)-count:]
	}

	return commands, nil
}
