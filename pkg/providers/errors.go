package providers

import "fmt"

// Common errors
var (
	ErrInvalidAPIKey      = fmt.Errorf("invalid or missing API key")
	ErrInvalidModel       = fmt.Errorf("invalid or unsupported model")
	ErrInvalidBaseURL     = fmt.Errorf("invalid base URL")
	ErrProviderNotFound   = fmt.Errorf("provider not found")
	ErrContextTooLarge    = fmt.Errorf("context exceeds maximum size")
	ErrRateLimitExceeded  = fmt.Errorf("rate limit exceeded")
	ErrQuotaExceeded      = fmt.Errorf("quota exceeded")
	ErrServiceUnavailable = fmt.Errorf("service unavailable")
)
