package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/Codilas/how/internal/config"
	"github.com/Codilas/how/internal/context"
	"github.com/Codilas/how/internal/manager"
	"github.com/Codilas/how/pkg/extractor"
	"github.com/Codilas/how/pkg/providers"
	"github.com/Codilas/how/pkg/providers/anthropic"
	"github.com/Codilas/how/pkg/text"
	"github.com/Codilas/how/pkg/version"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	verbose   bool
	useStream bool
	provider  string
	cfg       *config.Config
	mng       *manager.Manager
)

var rootCmd = &cobra.Command{
	Use:   "how [prompt...]",
	Short: "AI-powered shell assistant",
	Long:  `HOW is an AI-powered shell assistant that helps you with commands, explanations, and code generation.`,
	Args:  cobra.ArbitraryArgs,
	Run:   handlePrompt,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/how/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&useStream, "stream", "s", false, "stream response")
	rootCmd.PersistentFlags().StringVarP(&provider, "provider", "p", "", "AI provider to use")

	// Add version flag
	rootCmd.Flags().BoolP("version", "V", false, "show version")

	// Add subcommands
	rootCmd.AddCommand(setupCmd)
	rootCmd.AddCommand(historyCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(providersCmd)
}

func initConfig() {

	manager.RegisterProvider(anthropic.ProviderName, anthropic.NewProvider)

	var err error
	// Load configuration
	cfg, err = config.Load(cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize provider manager
	mng = manager.NewManager()

	// Ensure we have at least the mock provider for testing
	// ensureDefaultProviders(cfg)

	// Load providers from config
	if err := mng.LoadProviders(cfg.Providers); err != nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "Warning: failed to load some providers: %v\n", err)
		}
		// Fallback: ensure mock provider is available
		// if err := ensureMockProvider(); err != nil {
		// 	fmt.Fprintf(os.Stderr, "Failed to load fallback mock provider: %v\n", err)
		// }
	}

	// TODO: Figure out this mock stuff
	// Verify we have at least one working provider
	// if len(manager.ListProviders()) == 0 {
	// 	if err := ensureMockProvider(); err != nil {
	// 		fmt.Fprintf(os.Stderr, "No providers available and failed to create mock provider: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// }
}

// ensureDefaultProviders makes sure we have basic providers configured
func ensureDefaultProviders(cfg *config.Config) {
	if cfg.Providers == nil {
		cfg.Providers = make(map[string]config.ProviderConfig)
	}

	// Always ensure mock provider exists for testing
	// if _, exists := cfg.Providers["mock"]; !exists {
	// 	cfg.Providers["mock"] = config.ProviderConfig{
	// 		Type:      "mock",
	// 		Model:     "mock-model",
	// 		MaxTokens: 1000,
	// 	}
	// }

	// If no current provider is set and we have no real providers, default to mock
	// if cfg.CurrentProvider == "" || (cfg.CurrentProvider != "mock" && !hasRealProvider(cfg.Providers)) {
	// 	cfg.CurrentProvider = "mock"
	// }
}

// hasRealProvider checks if we have any non-mock providers configured
func hasRealProvider(providers map[string]config.ProviderConfig) bool {
	for name, provider := range providers {
		if name != "mock" && provider.APIKey != "" {
			return true
		}
	}
	return false
}

func handlePrompt(cmd *cobra.Command, args []string) {
	// Handle version flag - fixed
	if versionFlag, _ := cmd.Flags().GetBool("version"); versionFlag {
		fmt.Printf("HOW version %s\n", version.Version)
		return
	}

	if len(args) == 0 {
		cmd.Help()
		return
	}

	prompt := strings.Join(args, " ")

	// Determine which provider to use
	providerName := provider
	if providerName == "" {
		providerName = cfg.CurrentProvider
	}

	// Get the provider
	aiProvider, err := mng.GetProvider(providerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Available providers: %v\n", getProviderNames())
		os.Exit(1)
	}

	// Show provider info if verbose
	if verbose {
		info := aiProvider.GetInfo()
		fmt.Printf("Using %s (%s)\n", info.Name, info.Model)
	}

	// Gather context
	ctx, err := context.Gather(cfg.Context)
	if err != nil && verbose {
		fmt.Fprintf(os.Stderr, "Warning: failed to gather context: %v\n", err)
	}

	// Show context if verbose
	if verbose && ctx != nil {
		showContext(ctx)
	}

	// Send prompt
	if useStream {
		handleStreamingPrompt(aiProvider, prompt, ctx)
		return
	}

	handleRegularPrompt(aiProvider, prompt, ctx)
}

func handleRegularPrompt(aiProvider providers.Provider, prompt string, ctx *providers.Context) {
	// Show spinner
	s := spinner.New(spinner.CharSets[14], 100)
	s.Suffix = " Thinking..."
	if cfg.Display.Emoji {
		s.Suffix = " 🤔" + s.Suffix
	}
	s.Start()

	// Send prompt
	response, err := aiProvider.SendPrompt(prompt, ctx)
	s.Stop()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Display response
	displayResponse(response)
}

func handleStreamingPrompt(aiProvider providers.Provider, prompt string, ctx *providers.Context) {
	fmt.Print("🤖 ")

	// Send streaming prompt
	responseChan, err := aiProvider.SendPromptStream(prompt, ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Handle streaming response
	var fullText strings.Builder
	for chunk := range responseChan {
		if chunk.Error != nil {
			fmt.Fprintf(os.Stderr, "\nError: %v\n", chunk.Error)
			os.Exit(1)
		}

		if chunk.Done {
			break
		}

		fmt.Print(chunk.Text)
		fullText.WriteString(chunk.Text)
	}

	fmt.Println() // New line after streaming

	// Show metadata if verbose
	// TODO: ass more metadata
	if verbose {
		fmt.Printf("\nResponse completed\n")
	}
}

func displayResponse(resp *providers.Response) {

	// TODO: get this and put it in a better place
	commandExtractor := extractor.NewCommandExtractor()

	commands, err := commandExtractor.Extract(resp.Text)
	if err != nil {
		fmt.Print("Failed to parse commands")
		return
	}

	terminalFormatter := text.NewTerminalFormatter(text.CompactConfig())

	output := terminalFormatter.Format(resp.Text)

	fmt.Println()
	fmt.Print(output)
	fmt.Println()

	// Show suggested commands - TODO: Improve command handler and figure out formatter here
	if commands.Count() > 0 {
		fmt.Println()

		fmt.Println("Suggested commands:")

		for _, cmd := range commands.Commands {
			safety := ""
			if !cmd.Safe {
				safety = color.YellowString(" ⚠ ")
			}

			fmt.Printf("  %s%s %s\n", safety, color.CyanString("$"), cmd.Command)
			if cmd.Description != "" {
				fmt.Printf("    %s\n", color.HiBlackString(cmd.Description))
			}
		}
	}

	// Show metadata if verbose
	if verbose {
		fmt.Printf("\n%s\n", color.HiBlackString(formatMetadata(resp)))
	}
}

func showContext(ctx *providers.Context) {
	fmt.Printf("Context:\n")
	if ctx.WorkingDirectory != "" {
		fmt.Printf("  Directory: %s\n", ctx.WorkingDirectory)
	}
	if len(ctx.RecentCommands) > 0 {
		fmt.Printf("  Recent commands: %d\n", len(ctx.RecentCommands))
	}
	if len(ctx.Files) > 0 {
		fmt.Printf("  Files: %d\n", len(ctx.Files))
	}
	if ctx.Git != nil {
		fmt.Printf("  Git: %s (%s)\n", ctx.Git.Repository, ctx.Git.Branch)
	}
	fmt.Println()
}

func formatMetadata(resp *providers.Response) string {
	return fmt.Sprintf("Provider: %s | Model: %s | Tokens: %d | Time: %v",
		resp.Provider,
		resp.Model,
		resp.TokensUsed,
		resp.ResponseTime,
	)
}

func getProviderNames() []string {
	var names []string
	for _, info := range mng.ListProviders() {
		names = append(names, info.Type)
	}
	return names
}
