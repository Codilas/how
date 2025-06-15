package providers

// Provider defines the interface that all AI providers must implement
type Provider interface {
	// SendPrompt sends a prompt with context to the AI provider
	SendPrompt(prompt string, context *Context) (*Response, error)

	// SendPromptStream sends a prompt and returns a streaming response
	SendPromptStream(prompt string, context *Context) (<-chan StreamResponse, error)

	// ValidateConfig checks if the provider configuration is valid
	ValidateConfig() error

	// GetInfo returns provider information
	GetInfo() ProviderInfo

	// GetCapabilities returns what this provider supports
	GetCapabilities() Capabilities

	// GetModels returns a list of available models
	GetModels() ([]string, error)
}
