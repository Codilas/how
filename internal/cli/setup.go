package cli

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/Codilas/how/internal/config"
	"github.com/Codilas/how/pkg/providers"
	"github.com/Codilas/how/pkg/providers/anthropic"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure the AI assistant",
	Long:  `Interactive setup wizard to configure AI providers and preferences.`,
	Run:   runSetup,
}

func runSetup(cmd *cobra.Command, args []string) {
	fmt.Println("Setting up 'how' AI Shell Assistant")
	fmt.Println()

	// Provider selection
	var provider string
	providerPrompt := &survey.Select{
		Message: "Which AI provider would you like to use?",
		Options: []string{"Anthropic (Claude)", "OpenAI (GPT)", "Local Model"},
		Default: "Anthropic (Claude)",
	}
	survey.AskOne(providerPrompt, &provider)

	var apiKey string
	var model string

	switch provider {
	case "Anthropic (Claude)":
		cfg.CurrentProvider = "anthropic"

		apiKeyPrompt := &survey.Password{
			Message: "Enter your Anthropic API key:",
		}
		survey.AskOne(apiKeyPrompt, &apiKey)

		tmpProvider, err := anthropic.NewProvider(providers.Config{
			APIKey: apiKey,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing Anthropic provider: %v\n", err)
			os.Exit(1)
		}

		models, err := tmpProvider.GetModels()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching models: %v\n", err)
			os.Exit(1)
		}

		modelPrompt := &survey.Select{
			Message: "Choose Claude model:",
			Options: models,
			Default: models[0],
		}
		survey.AskOne(modelPrompt, &model)

		if cfg.Providers == nil {
			cfg.Providers = make(map[string]config.ProviderConfig)
		}
		cfg.Providers[anthropic.ProviderName] = config.ProviderConfig{
			Type:      anthropic.ProviderName,
			APIKey:    apiKey,
			Model:     model,
			MaxTokens: 1000,
		}

	case "OpenAI (GPT)":
		fmt.Println("OpenAI provider is not yet implemented.")
		os.Exit(1)
	case "Local Model":
		fmt.Println("Local model provider is not yet implemented.")
		os.Exit(1)
	}

	// Context preferences
	var includeFiles bool
	filesPrompt := &survey.Confirm{
		Message: "Include current directory files in context?",
		Default: true,
	}
	survey.AskOne(filesPrompt, &includeFiles)
	cfg.Context.IncludeFiles = includeFiles

	var includeHistory int
	historyPrompt := &survey.Select{
		Message: "How many recent commands to include?",
		Options: []string{"0", "3", "5", "10"},
		Default: "5",
	}
	var historyStr string
	survey.AskOne(historyPrompt, &historyStr)
	switch historyStr {
	case "0":
		includeHistory = 0
	case "3":
		includeHistory = 3
	case "5":
		includeHistory = 5
	case "10":
		includeHistory = 10
	}
	cfg.Context.IncludeHistory = includeHistory

	// Display preferences
	var syntaxHighlight bool
	highlightPrompt := &survey.Confirm{
		Message: "Enable syntax highlighting?",
		Default: true,
	}
	survey.AskOne(highlightPrompt, &syntaxHighlight)
	cfg.Display.SyntaxHighlight = syntaxHighlight

	var emoji bool
	emojiPrompt := &survey.Confirm{
		Message: "Use emoji in output?",
		Default: true,
	}
	survey.AskOne(emojiPrompt, &emoji)
	cfg.Display.Emoji = emoji

	// Save configuration
	if err := cfg.Save(""); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving configuration: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Configuration saved!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("• Run 'how install' to set up shell integration")
	fmt.Println("• Try: how \"how to use grep?\"")
	fmt.Println("• Try: how \"write a Python function to reverse a string\"")
}
