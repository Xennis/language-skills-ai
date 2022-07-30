package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Error struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%d:type=%s: %s", e.StatusCode, e.Type, e.Message)
}

type ErrorResponse struct {
	Error Error `json:"error"`
}

// Completions https://beta.openai.com/docs/api-reference/completions
type Completions struct {
	Model            string   `json:"model"`
	Prompt           []string `json:"prompt"`
	MaxTokens        int      `json:"max_tokens,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	TopP             float32  `json:"top_p,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	User             string   `json:"user,omitempty"`
}

type CompletionResponseChoice struct {
	Text         string `json:"text"`
	Index        int    `json:"index"`
	FinishReason string `json:"finish_reason"`
}

type CompletionResponseUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type CompletionResponse struct {
	ID      string                     `json:"id"`
	Object  string                     `json:"object"`
	Created int                        `json:"created"`
	Model   string                     `json:"model"`
	Choices []CompletionResponseChoice `json:"choices"`
	Usage   CompletionResponseUsage    `json:"usage"`
}

type Client interface {
	Completion(ctx context.Context, request Completions) (*CompletionResponse, error)
}

type client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) Client {
	httpClient := &http.Client{
		Timeout: time.Duration(10 * time.Second),
	}

	return &client{
		apiKey:     apiKey,
		baseURL:    "https://api.openai.com/v1",
		httpClient: httpClient,
	}
}

func (c *client) Completion(ctx context.Context, request Completions) (*CompletionResponse, error) {
	// create request
	raw, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}
	url := c.baseURL + "/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(raw))
	if err != nil {
		return nil, fmt.Errorf("http new request: %w", err)
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// do request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http client do: %w", err)
	}

	// read response
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http read body: %w", err)
		}
		var res ErrorResponse
		if err := json.Unmarshal(data, &res); err != nil {
			return nil, fmt.Errorf("unmarshal body: %w", err)
		}
		return nil, res.Error
	}
	defer resp.Body.Close()
	out := new(CompletionResponse)
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return nil, fmt.Errorf("decode body: %w", err)
	}
	return out, nil
}
