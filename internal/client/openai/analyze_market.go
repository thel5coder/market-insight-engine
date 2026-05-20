package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"market-insight-engine/internal/domain/model"
	"net/http"
)

func (o *openAIClient) AnalyzeMarket(ctx context.Context, prompt string) (*model.InsightResult, error) {
	url := fmt.Sprintf("%s/chat/completions", o.baseURL)

	type message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}

	payload := map[string]interface{}{
		"model": "gpt-4o-mini",
		"messages": []message{
			{
				Role:    "system",
				Content: `You are a professional trading analyst. You must respond ONLY with a valid JSON object containing the keys: "signal" (BULLISH, BEARISH, NEUTRAL), "confidence" (an integer from 0 to 100), and "reasoning" (a brief explanation up to 2 sentences).`,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		"response_format": map[string]string{"type": "json_object"},
		"temperature":     0.2,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal openai payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+o.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openai error: status %d", resp.StatusCode)
	}

	var response struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode openai response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("openai returned no choices")
	}

	var insight model.InsightResult
	contentStr := response.Choices[0].Message.Content

	if err := json.Unmarshal([]byte(contentStr), &insight); err != nil {
		return nil, fmt.Errorf("failed to parse AI content into struct: %w", err)
	}

	return &insight, nil
}
