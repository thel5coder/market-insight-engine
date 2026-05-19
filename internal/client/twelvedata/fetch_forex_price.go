package twelvedata

import (
	"context"
	"encoding/json"
	"fmt"
	"market-insight-engine/internal/domain/model"
	"net/http"
	"strconv"
)

func (t *twelvedataClient) FetchForexPrice(ctx context.Context, symbol string) (*model.MarketData, error) {
	url := fmt.Sprintf("%s/price?symbol=%s&apikey=%s", t.baseURL, symbol, t.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("twelvedata error: status %d", resp.StatusCode)
	}

	var data struct {
		Price string `json:"price"`
		Code  int    `json:"code"`
		Msg   string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if data.Code > 0 && data.Code != 200 {
		return nil, fmt.Errorf("twelvedata api error: %s", data.Msg)
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid price format from twelvedata: %w", err)
	}

	return &model.MarketData{
		Asset:        symbol,
		CurrentPrice: price,
	}, nil
}
