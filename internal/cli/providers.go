package cli

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var providersCmd = &cobra.Command{
	Use:   "providers",
	Short: "Manage AI providers",
	Long:  `List, test, and manage AI providers.`,
}

var listProvidersCmd = &cobra.Command{
	Use:   "list",
	Short: "List available providers",
	Run:   runListProviders,
}

var testProvidersCmd = &cobra.Command{
	Use:   "test",
	Short: "Test provider connectivity",
	Run:   runTestProviders,
}

var capabilitiesCmd = &cobra.Command{
	Use:   "capabilities [provider]",
	Short: "Show provider capabilities",
	Args:  cobra.MaximumNArgs(1),
	Run:   runCapabilities,
}

func init() {
	providersCmd.AddCommand(listProvidersCmd)
	providersCmd.AddCommand(testProvidersCmd)
	providersCmd.AddCommand(capabilitiesCmd)
}

func runListProviders(cmd *cobra.Command, args []string) {
	infos := mng.ListProviders()

	if len(infos) == 0 {
		fmt.Println("No providers configured.")
		fmt.Println("Run 'how setup' to configure a provider.")
		return
	}

	fmt.Println("Available providers:")
	fmt.Println()

	for _, info := range infos {
		status := color.GreenString("✓")

		// Check if provider is current
		current := ""
		if info.Type == cfg.CurrentProvider {
			current = color.BlueString(" (current)")
		}

		// Fixed: Use color.New instead of color.BoldString
		bold := color.New(color.Bold)
		fmt.Printf("  %s %s%s\n", status, bold.Sprint(info.Name), current)
		fmt.Printf("    Type: %s\n", info.Type)
		fmt.Printf("    Model: %s\n", info.Model)
		if info.Description != "" {
			fmt.Printf("    Description: %s\n", info.Description)
		}
		fmt.Println()
	}
}

func runTestProviders(cmd *cobra.Command, args []string) {
	fmt.Println("Testing provider connectivity...")
	fmt.Println()

	results := mng.HealthCheck()

	for name, err := range results {
		if err == nil {
			fmt.Printf("  %s %s\n", color.GreenString("✓"), name)
		} else {
			fmt.Printf("  %s %s: %v\n", color.RedString("✗"), name, err)
		}
	}
}

func runCapabilities(cmd *cobra.Command, args []string) {
	var providerName string
	if len(args) > 0 {
		providerName = args[0]
	} else {
		providerName = cfg.CurrentProvider
	}

	caps, err := mng.GetProviderCapabilities(providerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Capabilities for %s:\n\n", providerName)

	showCapability("Streaming", caps.Streaming)
	showCapability("Function Calling", caps.FunctionCalling)
	showCapability("Code Execution", caps.CodeExecution)
	showCapability("Image Analysis", caps.ImageAnalysis)
	showCapability("Conversation Memory", caps.ConversationMemory)

	fmt.Printf("  Max Context Size: %s\n", formatNumber(caps.MaxContextSize))
	fmt.Printf("  Max Tokens: %s\n", formatNumber(caps.MaxTokens))
}

func showCapability(name string, supported bool) {
	status := color.RedString("✗")
	if supported {
		status = color.GreenString("✓")
	}
	fmt.Printf("  %s %s\n", status, name)
}

func formatNumber(n int) string {
	if n >= 1000000 {
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	} else if n >= 1000 {
		return fmt.Sprintf("%.1fK", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}
