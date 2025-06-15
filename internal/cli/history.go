package cli

import (
	"fmt"
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
	historyFile := getHistoryFile()

	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		fmt.Println("📜 No conversation history found.")
		return
	}

	content, err := os.ReadFile(historyFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading history: %v\n", err)
		return
	}

	if len(content) == 0 {
		fmt.Println("📜 No conversation history found.")
		return
	}

	fmt.Println("📜 Recent conversations:")
	fmt.Println()
	fmt.Print(string(content))
}

func runClearHistory(cmd *cobra.Command, args []string) {
	historyFile := getHistoryFile()

	if err := os.Remove(historyFile); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error clearing history: %v\n", err)
		return
	}

	fmt.Println("🗑️ Conversation history cleared.")
}

func getHistoryFile() string {
	configDir, _ := os.UserHomeDir()
	return filepath.Join(configDir, ".config", "how", "history.txt")
}
