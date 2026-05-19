package binance

import (
	"market-insight-engine/internal/config"
	"market-insight-engine/internal/domain/interfaces"
	"net/http"
)

type binanceClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

var _ interfaces.CryptoClient = (*binanceClient)(nil)

func NewBinanceClient(cfg config.ExternalMarketAPI, httpClient *http.Client) interfaces.CryptoClient {
	return &binanceClient{
		apiKey:     cfg.BinanceApi.APIKey,
		baseURL:    cfg.BinanceApi.BaseUrl,
		httpClient: httpClient,
	}
}
