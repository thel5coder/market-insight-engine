package gemini

import (
	"market-insight-engine/internal/config"
	"market-insight-engine/internal/domain/interfaces"
	"net/http"
)

type geminiClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
	model      string
}

var _ interfaces.AIClient = (*geminiClient)(nil)

func NewGeminiClient(cfg config.ExternalMarketAPI, httpClient *http.Client) interfaces.AIClient {
	return &geminiClient{
		apiKey:     cfg.GeminiAI.APIKey,
		baseURL:    cfg.GeminiAI.BaseUrl,
		httpClient: httpClient,
	}
}
