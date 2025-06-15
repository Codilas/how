package anthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Codilas/how/pkg/providers"
)

// Provider implements the providers.Provider interface for Anthropic's Claude API
type Provider struct {
	httpClient *http.Client
	baseURL    string
	cfg        providers.Config
}

// API request/response structures
type request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []message `json:"messages"`
	System    string    `json:"system,omitempty"`
	Stream    bool      `json:"stream,omitempty"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Role         string         `json:"role"`
	Content      []contentBlock `json:"content"`
	Model        string         `json:"model"`
	StopReason   string         `json:"stop_reason"`
	StopSequence string         `json:"stop_sequence"`
	Usage        usage          `json:"usage"`
}

type contentBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type apiError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

// modelsResponse represents the response from the models API endpoint
type modelsResponse struct {
	Data    []model `json:"data"`
	FirstID string  `json:"first_id"`
	HasMore bool    `json:"has_more"`
	LastID  string  `json:"last_id"`
}

type model struct {
	CreatedAt   string `json:"created_at"`
	DisplayName string `json:"display_name"`
	ID          string `json:"id"`
	Type        string `json:"type"`
}

const (
	// ProviderName is anthropic claude provider name
	ProviderName = "anthropic"
	baseURL      = "https://api.anthropic.com/v1"
	version      = "2023-06-01"
	displayName  = "Anthropic Claude AI"
	description  = "Anthropic's Claude AI assistant, designed for safe and helpful interactions."
)

// NewProvider creates a new Anthropic provider instance
func NewProvider(cfg providers.Config) (providers.Provider, error) {
	url := baseURL
	if cfg.BaseURL != "" {
		url = cfg.BaseURL
	}

	return &Provider{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		cfg:     cfg,
		baseURL: url,
	}, nil
}

// SendPrompt implements the providers.Provider interface
func (p *Provider) SendPrompt(prompt string, context *providers.Context) (*providers.Response, error) {
	startTime := time.Now()

	// Build the request
	req, err := p.buildRequest(prompt, context, false)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Create HTTP request
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", fmt.Sprint(p.baseURL, "/messages"), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Send the request
	var apiResp response
	if err := p.doRequest(httpReq, &apiResp); err != nil {
		return nil, err
	}

	var text string
	if len(apiResp.Content) > 0 {
		text = apiResp.Content[0].Text
	}

	response := &providers.Response{
		Text:         text,
		Model:        apiResp.Model,
		Provider:     ProviderName,
		TokensUsed:   apiResp.Usage.InputTokens + apiResp.Usage.OutputTokens,
		ResponseTime: time.Since(startTime),
	}

	return response, nil
}

// SendPromptStream implements streaming for the providers.Provider interface
func (p *Provider) SendPromptStream(prompt string, context *providers.Context) (<-chan providers.StreamResponse, error) {
	return nil, fmt.Errorf("streaming not yet supported by %s provider", ProviderName)
}

// ValidateConfig implements the providers.Provider interface
func (p *Provider) ValidateConfig() error {
	if p.cfg.APIKey == "" {
		return providers.ErrInvalidAPIKey
	}

	if p.cfg.Model == "" {
		return providers.ErrInvalidModel
	}

	if p.cfg.MaxTokens <= 0 || p.cfg.MaxTokens > 4096 {
		return fmt.Errorf("max_tokens must be between 1 and 4096, got %d", p.cfg.MaxTokens)
	}

	return nil
}

// GetInfo implements the providers.Provider interface
func (p *Provider) GetInfo() providers.ProviderInfo {
	return providers.ProviderInfo{
		Name:        displayName,
		Type:        ProviderName,
		Model:       p.cfg.Model,
		Description: description,
	}
}

// GetCapabilities implements the providers.Provider interface
func (p *Provider) GetCapabilities() providers.Capabilities {
	return providers.Capabilities{
		Streaming:          false,
		FunctionCalling:    false,
		CodeExecution:      false,
		ImageAnalysis:      false,
		ConversationMemory: true,
		MaxContextSize:     200000,
		MaxTokens:          4096,
	}
}

// GetModels implements the providers.Provider interface
func (p *Provider) GetModels() ([]string, error) {
	httpReq, err := http.NewRequest("GET", fmt.Sprint(p.baseURL, "/models"), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Send request and parse response
	var modelsResp modelsResponse
	if err := p.doRequest(httpReq, &modelsResp); err != nil {
		return nil, err
	}

	models := make([]string, len(modelsResp.Data))
	for i, model := range modelsResp.Data {
		models[i] = model.ID
	}

	return models, nil
}

// buildRequest creates an API request
func (p *Provider) buildRequest(prompt string, context *providers.Context, stream bool) (*request, error) {
	// Build system prompt with context
	systemPrompt, err := p.buildSystemPrompt(context)
	if err != nil {
		return nil, err
	}

	// Create the user message
	userMessage := message{
		Role:    "user",
		Content: prompt,
	}

	// Add conversation history if available
	messages := []message{}
	if context != nil && len(context.PreviousPrompts) > 0 {
		for _, entry := range context.PreviousPrompts {
			messages = append(messages,
				message{Role: "user", Content: entry.Prompt},
				message{Role: "assistant", Content: entry.Response},
			)
		}
	}

	messages = append(messages, userMessage)

	return &request{
		Model:     p.cfg.Model,
		MaxTokens: p.cfg.MaxTokens,
		Messages:  messages,
		System:    systemPrompt,
		Stream:    stream,
	}, nil
}

// buildSystemPrompt creates a system prompt with context information
func (p *Provider) buildSystemPrompt(context *providers.Context) (string, error) {
	// Build the system context string
	systemContext := buildSystemContext(context)

	// Create template data
	data := TemplateData{
		SystemContext: systemContext,
	}

	// Process the template
	prompt, err := processTemplate(systemPromptTemplate, data)
	if err != nil {
		return "", fmt.Errorf("failed to process system prompt template: %w", err)
	}

	return prompt, nil
}

// doRequest sends an HTTP request and parses the response into the provided struct
func (p *Provider) doRequest(httpReq *http.Request, response interface{}) error {
	// Add common headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.cfg.APIKey)
	httpReq.Header.Set("anthropic-version", version)

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		var errorResp errorResponse
		if json.Unmarshal(body, &errorResp) == nil {
			return fmt.Errorf("API error (%d): %s", resp.StatusCode, errorResp.Error.Message)
		}
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
