package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"market-insight-engine/internal/domain/model"
	"net/http"
	"net/url"
	"unsafe"
)

func (g *geminiClient) AnalyzeMarket(ctx context.Context, prompt string) (*model.InsightResult, error) {
	endpoint, err := url.Parse(fmt.Sprintf("%s/%s:generateContent", g.baseURL, g.model))
	if err != nil {
		return nil, fmt.Errorf("failed to build gemini endpoint URL: %w", err)
	}
	reqURL := endpoint.String()

	payload := model.GeminiRequest{
		Contents: []model.GeminiContent{
			{Parts: []model.GeminiPart{{Text: prompt}}},
		},

		SystemInstruction: model.GeminiContent{
			Parts: []model.GeminiPart{{Text: `You are a professional trading analyst. You must respond ONLY with a valid JSON object containing the keys: "signal" (BULLISH, BEARISH, NEUTRAL), "confidence" (an integer from 0 to 100), and "reasoning" (a brief explanation up to 2 sentences).`}},
		},
		GenerationConfig: model.GenerationConfig{
			ResponseMIMEType: "application/json",
			Temperature:      0.2,
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, fmt.Errorf("failed to marshal gemini payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", g.apiKey)

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var apiErr struct {
			Error struct {
				Code    int    `json:"code"`
				Message string `json:"message"`
				Status  string `json:"status"`
			} `json:"error"`
		}
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		if jsonErr := json.Unmarshal(body, &apiErr); jsonErr == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("gemini api error %d (%s): %s",
				resp.StatusCode, apiErr.Error.Status, apiErr.Error.Message)
		}
		return nil, fmt.Errorf("gemini api error: status %d, body: %s", resp.StatusCode, body)
	}

	var response struct {
		Candidates []struct {
			FinishReason string `json:"finishReason"`
			Content      struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	const maxResponseBytes = 1 << 20
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseBytes)).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode gemini response: %w", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini returned empty response")
	}

	candidate := response.Candidates[0]
	if candidate.FinishReason != "" && candidate.FinishReason != "STOP" {
		return nil, fmt.Errorf("gemini stopped with reason: %s", candidate.FinishReason)
	}

	rawText := candidate.Content.Parts[0].Text
	var insight model.InsightResult
	if err := json.Unmarshal(unsafe.Slice(unsafe.StringData(rawText), len(rawText)), &insight); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini JSON into struct: %w", err)
	}

	return &insight, nil
}
