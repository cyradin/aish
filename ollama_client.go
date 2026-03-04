package aish

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type OllamaClient struct {
	baseURL string
	inner   *http.Client
}

func NewOllamaClient(baseURL string, inner *http.Client) *OllamaClient {
	return &OllamaClient{
		baseURL: baseURL,
		inner:   inner,
	}
}

func (c *OllamaClient) Generate(model string, prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.inner.Post(
		strings.TrimRight(c.baseURL, "/")+"/api/generate",
		"application/json",
		bytes.NewBuffer(data),
	)
	if err != nil {
		return "", err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Response, nil
}
