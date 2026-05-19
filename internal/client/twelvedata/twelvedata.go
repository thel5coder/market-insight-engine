package twelvedata

import (
	"market-insight-engine/internal/config"
	"market-insight-engine/internal/domain/interfaces"
	"net/http"
)

type twelvedataClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// Compile-time check
var _ interfaces.ForexClient = (*twelvedataClient)(nil)

func NewTwelveDataClient(cfg config.ExternalMarketAPI, httpClient *http.Client) interfaces.ForexClient {
	return &twelvedataClient{
		apiKey:     cfg.TwelveDataApi.APIKey,
		baseURL:    cfg.TwelveDataApi.BaseUrl,
		httpClient: httpClient,
	}
}
