package ai_model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"backend/domain/gm_session"
)

type OpenRouterAgent struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

func NewOpenRouterAgent() (*OpenRouterAgent, error) {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENROUTER_API_KEY environment variable is not set")
	}

	return &OpenRouterAgent{
		apiKey:     apiKey,
		baseURL:    "https://openrouter.ai/api/v1",
		httpClient: &http.Client{},
	}, nil
}

type chatCompletionRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (a *OpenRouterAgent) GetGMResponse(
	ctx context.Context,
	week int,
	categories, tickers []string,
	lastRatings map[string]string,
) (*gm_session.GMWeekData, error) {
	// Prepare template data
	templateData := map[string]interface{}{
		"Week":        week,
		"Categories":  categories,
		"Tickers":     tickers,
		"LastRatings": lastRatings,
	}

	prompt, err := LoadPrompt("infrastructure/ai_model/gm_prompt.txt", templateData)
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt: %w", err)
	}

	reqBody := chatCompletionRequest{
		Model: os.Getenv("OPENROUTER_MODEL_NAME"),
		Messages: []chatMessage{
			{
				Role:    "system",
				Content: "You are a stock market game master that provides realistic market simulation data.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	content, err := a.makeOpenRouterRequest(ctx, jsonBody)
	if err != nil {
		return nil, fmt.Errorf("failed to make OpenRouter request: %w", err)
	}

	// Parse the response
	var aiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(content, &aiResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(aiResp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned from API")
	}

	parsedData := &gm_session.GMWeekData{
		// Parse the aiResp.Choices[0].Message.Content into your structure
		// TODO: Implement parsing logic for Headlines and Stocks
	}

	return parsedData, nil
}

func (a *OpenRouterAgent) makeOpenRouterRequest(ctx context.Context, jsonBody []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/chat/completions", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	req.Header.Set("HTTP-Referer", os.Getenv("OPENROUTER_REFERER"))
	req.Header.Set("X-Title", "Stock Market Game Simulation")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return content, nil
}
