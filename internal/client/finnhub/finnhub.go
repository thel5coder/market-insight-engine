package finnhub

import (
	"market-insight-engine/internal/config"
	"market-insight-engine/internal/domain/interfaces"
	"net/http"
)

type finnhubClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

var _ interfaces.NewsClient = (*finnhubClient)(nil)

func NewFinnhubClient(cfg config.ExternalMarketAPI, httpClient *http.Client) interfaces.NewsClient {
	return &finnhubClient{
		apiKey:     cfg.FinhubApi.APIKey,
		baseURL:    cfg.FinhubApi.BaseUrl,
		httpClient: httpClient,
	}
}
