package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"market-insight-engine/internal/domain/model"
	"net/http"
)

func (f *finnhubClient) FetchLatestNews(ctx context.Context, category string) ([]*model.News, error) {
	// Category could be 'general', 'forex', 'crypto', etc.
	url := fmt.Sprintf("%s/api/v1/news?category=%s&token=%s", f.baseURL, category, f.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("finnhub api error: status %d", resp.StatusCode)
	}

	type finnhubNews struct {
		Headline string `json:"headline"`
		Summary  string `json:"summary"`
		Source   string `json:"source"`
		URL      string `json:"url"`
		Time     int64  `json:"datetime"`
	}

	var rawNews []finnhubNews
	if err := json.NewDecoder(resp.Body).Decode(&rawNews); err != nil {
		return nil, fmt.Errorf("failed to decode finnhub news: %w", err)
	}

	newsList := make([]*model.News, 0, len(rawNews))
	for _, n := range rawNews {
		newsList = append(newsList, &model.News{
			Title:       n.Headline,
			Content:     n.Summary,
			Source:      n.Source,
			URL:         n.URL,
			Timestamp:   n.Time,
			ImpactLevel: model.IMPACT_LEVEL_UNKNOWN,
		})
	}

	return newsList, nil
}
