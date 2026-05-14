package client

import (
	"context"
	"encoding/json"
	"fmt"
	"market-insight-engine/internal/config"
	"market-insight-engine/internal/domain/interfaces"
	"market-insight-engine/internal/domain/model"
	"net/http"
	"strconv"
)

type BinanceClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

var _ interfaces.CryptoClient = (*BinanceClient)(nil)

func NewBinanceClient(cfg config.ExternalMarketAPI, httpClient *http.Client) interfaces.CryptoClient {
	return &BinanceClient{
		apiKey:     cfg.BinanceApi.APIKey,
		baseURL:    cfg.BinanceApi.BaseUrl,
		httpClient: httpClient,
	}
}

func (b *BinanceClient) FetchPrice(ctx context.Context, symbol string) (*model.MarketData, error) {
	url := fmt.Sprintf("%s/api/v3/ticker/price?symbol=%s", b.baseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-MBX-APIKEY", b.apiKey)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance error: status %d", resp.StatusCode)
	}

	var data struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price format: %w", err)
	}

	return &model.MarketData{
		Asset:        data.Symbol,
		CurrentPrice: price,
	}, nil
}
