package twelvedata

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func (t *twelvedataClient) FetchRSI(ctx context.Context, symbol string, interval string) (float64, error) {
	// TwelveData RSI endpoint: /rsi?symbol=...&interval=...
	url := fmt.Sprintf("%s/rsi?symbol=%s&interval=%s&outputsize=1&apikey=%s",
		t.baseURL, symbol, interval, t.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Values []struct {
			RSI string `json:"rsi"`
		} `json:"values"`
		Status string `json:"status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	if data.Status != "ok" || len(data.Values) == 0 {
		return 0, fmt.Errorf("twelvedata rsi error or no data")
	}

	rsi, err := strconv.ParseFloat(data.Values[0].RSI, 64)
	if err != nil {
		return 0, err
	}

	return rsi, nil
}
