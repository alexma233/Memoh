package chat

import (
	"fmt"
	"strings"
	"time"

	dbsqlc "github.com/memohai/memoh/internal/db/sqlc"
)

// CreateProvider creates a chat provider instance.
func CreateProvider(provider dbsqlc.LlmProvider, timeout time.Duration) (Provider, error) {
	clientType := strings.ToLower(strings.TrimSpace(provider.ClientType))
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	switch clientType {
	case ProviderOpenAI, ProviderOpenAICompat:
		if strings.TrimSpace(provider.ApiKey) == "" {
			return nil, fmt.Errorf("openai api key is required")
		}
		return NewOpenAIProvider(provider.ApiKey, provider.BaseUrl, timeout)
	case ProviderAnthropic:
		if strings.TrimSpace(provider.ApiKey) == "" {
			return nil, fmt.Errorf("anthropic api key is required")
		}
		return NewAnthropicProvider(provider.ApiKey, timeout)
	case ProviderGoogle:
		if strings.TrimSpace(provider.ApiKey) == "" {
			return nil, fmt.Errorf("google api key is required")
		}
		return NewGoogleProvider(provider.ApiKey, timeout)
	case ProviderOllama:
		return NewOllamaProvider(provider.BaseUrl, timeout)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", clientType)
	}
}
