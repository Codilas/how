package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Show conversation history",
	Long:  `Display recent AI conversations and prompts.`,
	Run:   runHistory,
}

var clearHistoryCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear conversation history",
	Long:  `Clear all stored conversation history.`,
	Run:   runClearHistory,
}

func init() {
	historyCmd.AddCommand(clearHistoryCmd)
}

func runHistory(cmd *cobra.Command, args []string) {
	runHistoryWithWriter(os.Stdout, os.Stderr, getHistoryFile)
}

// runHistoryWithWriter is the testable implementation of runHistory.
// It accepts output writers and a function to get the history file path.
func runHistoryWithWriter(stdout, stderr io.Writer, getHistoryFileFn func() string) {
	historyFile := getHistoryFileFn()

	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		fmt.Fprintln(stdout, "📜 No conversation history found.")
		return
	}

	content, err := os.ReadFile(historyFile)
	if err != nil {
		fmt.Fprintf(stderr, "Error reading history: %v\n", err)
		return
	}

	if len(content) == 0 {
		fmt.Fprintln(stdout, "📜 No conversation history found.")
		return
	}

	fmt.Fprintln(stdout, "📜 Recent conversations:")
	fmt.Fprintln(stdout)
	fmt.Fprint(stdout, string(content))
}

func runClearHistory(cmd *cobra.Command, args []string) {
	runClearHistoryWithWriter(os.Stdout, os.Stderr, getHistoryFile)
}

// runClearHistoryWithWriter is the testable implementation of runClearHistory.
// It accepts output writers and a function to get the history file path.
func runClearHistoryWithWriter(stdout, stderr io.Writer, getHistoryFileFn func() string) {
	historyFile := getHistoryFileFn()

	if err := os.Remove(historyFile); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(stderr, "Error clearing history: %v\n", err)
		return
	}

	fmt.Fprintln(stdout, "🗑️ Conversation history cleared.")
}

func getHistoryFile() string {
	configDir, _ := os.UserHomeDir()
	return filepath.Join(configDir, ".config", "how", "history.txt")
}
