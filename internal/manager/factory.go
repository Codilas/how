package manager

import (
	"fmt"

	"github.com/Codilas/how/pkg/providers"
)

// ProviderFactory creates providers based on configuration
type ProviderFactory struct {
	providers map[string]func(providers.Config) (providers.Provider, error)
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]func(providers.Config) (providers.Provider, error)),
	}
}

// RegisterProvider registers a provider constructor
func (f *ProviderFactory) RegisterProvider(providerType string, constructor func(providers.Config) (providers.Provider, error)) {
	f.providers[providerType] = constructor
}

// CreateProvider creates a provider instance from configuration
func (f *ProviderFactory) CreateProvider(cfg providers.Config) (providers.Provider, error) {
	constructor, exists := f.providers[cfg.Type]
	if !exists {
		return nil, fmt.Errorf("unknown provider type: %s", cfg.Type)
	}

	provider, err := constructor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider %s: %w", cfg.Type, err)
	}

	// Validate the provider configuration
	if err := provider.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("invalid configuration for provider %s: %w", cfg.Type, err)
	}

	return provider, nil
}

// GetSupportedProviders returns a list of supported provider types
func (f *ProviderFactory) GetSupportedProviders() []string {
	var types []string
	for providerType := range f.providers {
		types = append(types, providerType)
	}
	return types
}

// Global factory instance
var defaultFactory = NewProviderFactory()

// RegisterProvider registers a provider with the default factory
func RegisterProvider(providerType string, constructor func(providers.Config) (providers.Provider, error)) {
	defaultFactory.RegisterProvider(providerType, constructor)
}

// GetProvider creates a provider instance using the default factory
func GetProvider(providerName string, providers map[string]providers.Config) (providers.Provider, error) {
	providerConfig, exists := providers[providerName]
	if !exists {
		return nil, fmt.Errorf("provider %s not found in configuration", providerName)
	}

	return defaultFactory.CreateProvider(providerConfig)
}

// GetSupportedProviders returns supported providers from the default factory
func GetSupportedProviders() []string {
	return defaultFactory.GetSupportedProviders()
}
