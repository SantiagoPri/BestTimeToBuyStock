package ai_model

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"backend/domain/gm_session"
	"backend/domain/stock"
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

func sanitizeJSONStrict(raw string) string {
	var out strings.Builder
	inString := false
	escapeNext := false
	inNumber := false

	for i := 0; i < len(raw); i++ {
		c := raw[i]

		if escapeNext {
			out.WriteByte('\\')
			out.WriteByte(c)
			escapeNext = false
			continue
		}

		if inString {
			out.WriteByte(c)
			if c == '\\' {
				escapeNext = true
			} else if c == '"' {
				inString = false
			}
		} else if inNumber {
			if (c >= '0' && c <= '9') || c == '.' || c == '-' {
				out.WriteByte(c)
			} else {
				inNumber = false
				if c == ',' || c == '}' || c == ']' {
					out.WriteByte(c)
				}
			}
		} else {
			switch c {
			case '"':
				inString = true
				out.WriteByte(c)
			case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				inNumber = true
				out.WriteByte(c)
			case '{', '}', '[', ']', ':', ',':
				out.WriteByte(c)
			case ' ', '\n', '\t', '\r':
				out.WriteByte(c)
			}
		}
	}

	return out.String()
}

func extractFirstJSONObject(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty response")
	}

	// Find the first opening brace
	start := strings.Index(raw, "{")
	if start == -1 {
		return "", fmt.Errorf("no JSON object found in response")
	}

	// Track nested braces
	braces := 0
	inString := false
	escapeNext := false

	for i := start; i < len(raw); i++ {
		c := raw[i]

		if escapeNext {
			escapeNext = false
			continue
		}

		if inString {
			if c == '\\' {
				escapeNext = true
			} else if c == '"' {
				inString = false
			}
			continue
		}

		if c == '"' {
			inString = true
		} else if c == '{' {
			braces++
		} else if c == '}' {
			braces--
			if braces == 0 {
				// Found matching closing brace
				extracted := raw[start : i+1]
				sanitized := sanitizeJSONStrict(extracted)

				// Validate JSON structure
				var js json.RawMessage
				if err := json.Unmarshal([]byte(sanitized), &js); err != nil {
					return "", fmt.Errorf("invalid JSON structure: %w\nProblematic content: %s", err, sanitized)
				}

				return sanitized, nil
			}
		}
	}

	return "", fmt.Errorf("unclosed JSON object")
}

func (a *OpenRouterAgent) GetGMResponse(
	ctx context.Context,
	categories []string,
	stocks []stock.Stock,
) (map[string]*gm_session.GMWeekData, error) {
	// Prepare template data
	stocksData := make([]map[string]interface{}, len(stocks))
	for i, s := range stocks {
		stocksData[i] = map[string]interface{}{
			"ticker":     s.Ticker,
			"company":    s.Company,
			"category":   s.Category,
			"ratingFrom": s.RatingFrom,
			"ratingTo":   s.RatingTo,
		}
	}

	templateData := map[string]interface{}{
		"Categories": categories,
		"stocks":     stocksData,
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

	// Parse the AI response content into weeks data
	var response struct {
		Weeks map[string]*gm_session.GMWeekData `json:"weeks"`
	}

	cleanedContent, err := extractFirstJSONObject(aiResp.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to extract first JSON object: %w", err)
	}

	if err := json.Unmarshal([]byte(cleanedContent), &response); err != nil {
		return nil, fmt.Errorf("failed to parse AI response content: %w", err)
	}

	return response.Weeks, nil
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
