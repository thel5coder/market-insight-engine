package openai

import (
	"market-insight-engine/internal/config"
	"market-insight-engine/internal/domain/interfaces"
	"net/http"
)

type openAIClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

var _ interfaces.AIClient = (*openAIClient)(nil)

func NewOpenAIClient(cfg config.ExternalMarketAPI, httpClient *http.Client) interfaces.AIClient {
	return &openAIClient{
		apiKey:     cfg.OpenAI.APIKey,
		baseURL:    cfg.OpenAI.BaseUrl, // Jangan lupa bisa taruh di .env juga
		httpClient: httpClient,
	}
}
