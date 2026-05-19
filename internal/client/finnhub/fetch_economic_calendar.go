package finnhub

import (
	"context"
	"encoding/json"
	"fmt"
	"market-insight-engine/internal/domain/model"
	"net/http"
)

func (f *finnhubClient) FetchEconomicCalendar(ctx context.Context, from, to string) ([]*model.EconomicEvent, error) {
	// Endpoint: /economic-calendar?from=YYYY-MM-DD&to=YYYY-MM-DD
	url := fmt.Sprintf("%s/api/v1/calendar/economic?from=%s&to=%s&token=%s", f.baseURL, from, to, f.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data struct {
		EconomicCalendar []struct {
			Event    string  `json:"event"`
			Country  string  `json:"country"`
			Actual   float64 `json:"actual"`
			Estimate float64 `json:"estimate"`
			Prev     float64 `json:"prev"`
			Impact   int     `json:"impact"`
			Unit     string  `json:"unit"`
			Time     string  `json:"time"`
		} `json:"economicCalendar"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	events := make([]*model.EconomicEvent, 0, len(data.EconomicCalendar))
	for _, e := range data.EconomicCalendar {
		events = append(events, &model.EconomicEvent{
			Event:    e.Event,
			Country:  e.Country,
			Actual:   e.Actual,
			Estimate: e.Estimate,
			Previous: e.Prev,
			Impact:   e.Impact,
			Unit:     e.Unit,
			// Conversion logic for time string to Unix would go here
		})
	}

	return events, nil
}
