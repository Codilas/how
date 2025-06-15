package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install shell integration",
	Long:  `Set up shell integration for seamless AI assistance.`,
	Run:   runInstall,
}

func runInstall(cmd *cobra.Command, args []string) {
	fmt.Println("Shell Integration Setup")
	fmt.Println()

	// Detect current shell
	shell := detectShell()
	fmt.Printf("Detected shell: %s\n", shell)
	fmt.Println()

	// Get binary path
	binaryPath, err := os.Executable()
	if err != nil {
		binaryPath = "how"
	}

	switch shell {
	case "bash":
		showBashIntegration(binaryPath)
	case "zsh":
		showZshIntegration(binaryPath)
	case "fish":
		showFishIntegration(binaryPath)
	default:
		showGenericIntegration(binaryPath)
	}

	fmt.Println()
	fmt.Println("After adding the integration, reload your shell:")
	fmt.Printf("  source ~/%s\n", getShellConfigFile(shell))
	fmt.Println()
}

func detectShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "unknown"
	}
	return filepath.Base(shell)
}

func getShellConfigFile(shell string) string {
	switch shell {
	case "bash":
		return ".bashrc"
	case "zsh":
		return ".zshrc"
	case "fish":
		return ".config/fish/config.fish"
	default:
		return ".profile"
	}
}

func showBashIntegration(binaryPath string) {
}

func showZshIntegration(binaryPath string) {
}

func showFishIntegration(binaryPath string) {
}

func showGenericIntegration(binaryPath string) {
}
