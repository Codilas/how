package manager

import (
	"fmt"
	"sort"

	"github.com/Codilas/how/internal/config"
	"github.com/Codilas/how/pkg/providers"
)

// Manager handles provider lifecycle and selection
type Manager struct {
	providers map[string]providers.Provider
	factory   *ProviderFactory
}

// NewManager creates a new provider manager
func NewManager() *Manager {
	return &Manager{
		providers: make(map[string]providers.Provider),
		factory:   defaultFactory,
	}
}

// LoadProviders initializes all configured providers
func (m *Manager) LoadProviders(cfg map[string]config.ProviderConfig) error {

	for name, providerCfg := range cfg {
		provider, err := m.factory.CreateProvider(convertCfg(providerCfg))
		if err != nil {
			return fmt.Errorf("failed to load provider %s: %w", name, err)
		}

		m.providers[name] = provider
	}

	return nil
}

// GetProvider retrieves a provider by name
func (m *Manager) GetProvider(name string) (providers.Provider, error) {
	provider, exists := m.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider %s not found", name)
	}

	return provider, nil
}

// ListProviders returns information about all loaded providers
func (m *Manager) ListProviders() []providers.ProviderInfo {
	var infos []providers.ProviderInfo

	for _, provider := range m.providers {
		infos = append(infos, provider.GetInfo())
	}

	// Sort by name for consistent output
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name < infos[j].Name
	})

	return infos
}

// ValidateProviders checks all providers are properly configured
func (m *Manager) ValidateProviders() map[string]error {
	errors := make(map[string]error)

	for name, provider := range m.providers {
		if err := provider.ValidateConfig(); err != nil {
			errors[name] = err
		}
	}

	return errors
}

// GetProviderCapabilities returns capabilities for a specific provider
func (m *Manager) GetProviderCapabilities(name string) (providers.Capabilities, error) {
	provider, err := m.GetProvider(name)
	if err != nil {
		return providers.Capabilities{}, err
	}

	return provider.GetCapabilities(), nil
}

// SelectBestProvider chooses the best provider for a given task
func (m *Manager) SelectBestProvider(requirements ProviderRequirements) (string, providers.Provider, error) {
	var bestProvider string
	var bestScore int

	for name, provider := range m.providers {
		score := m.scoreProvider(provider, requirements)
		if score > bestScore {
			bestScore = score
			bestProvider = name
		}
	}

	if bestProvider == "" {
		return "", nil, fmt.Errorf("no suitable provider found for requirements")
	}

	return bestProvider, m.providers[bestProvider], nil
}

// ProviderRequirements defines what features are needed
type ProviderRequirements struct {
	Streaming       bool
	FunctionCalling bool
	ImageAnalysis   bool
	MinContextSize  int
	PreferredTypes  []string // e.g., ["anthropic", "openai"]
}

// scoreProvider rates how well a provider matches requirements
func (m *Manager) scoreProvider(provider providers.Provider, req ProviderRequirements) int {
	caps := provider.GetCapabilities()
	info := provider.GetInfo()
	score := 0

	// Check required capabilities
	if req.Streaming && caps.Streaming {
		score += 10
	}
	if req.FunctionCalling && caps.FunctionCalling {
		score += 10
	}
	if req.ImageAnalysis && caps.ImageAnalysis {
		score += 10
	}
	if caps.MaxContextSize >= req.MinContextSize {
		score += 5
	}

	// Prefer certain provider types
	for _, preferred := range req.PreferredTypes {
		if info.Type == preferred {
			score += 20
			break
		}
	}

	// Validate configuration
	if provider.ValidateConfig() != nil {
		score = 0 // Invalid providers get zero score
	}

	return score
}

// HealthCheck tests all providers
func (m *Manager) HealthCheck() map[string]error {
	results := make(map[string]error)

	for name, provider := range m.providers {
		// Simple validation check
		err := provider.ValidateConfig()
		results[name] = err
	}

	return results
}

// ReloadProvider reloads a specific provider with new configuration
func (m *Manager) ReloadProvider(name string, cfg config.ProviderConfig) error {
	provider, err := m.factory.CreateProvider(convertCfg(cfg))
	if err != nil {
		return fmt.Errorf("failed to reload provider %s: %w", name, err)
	}

	m.providers[name] = provider
	return nil
}

// RemoveProvider removes a provider from the manager
func (m *Manager) RemoveProvider(name string) {
	delete(m.providers, name)
}

// GetAvailableTypes returns all available provider types
func (m *Manager) GetAvailableTypes() []string {
	return m.factory.GetSupportedProviders()
}

// convertCfg converts a config.ProviderConfig to providers.Config
func convertCfg(cfg config.ProviderConfig) providers.Config {
	return providers.Config{
		Type:      cfg.Type,
		APIKey:    cfg.APIKey,
		Model:     cfg.Model,
		BaseURL:   cfg.BaseURL,
		MaxTokens: cfg.MaxTokens,
	}
}
