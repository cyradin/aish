package aish

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const ollamaSystemPrompt = `
You generate Linux bash commands.

Output Rules (STRICT):
1. Output EXACTLY ONE line of text.
2. The output must be a valid bash command that can be executed directly.
3. If the task requires multiple operations, chain them using && or ; or || — this still counts as ONE command.
4. NEVER use newlines in the output.

Forbidden Content:
- NO backticks
- NO markdown
- NO code fences or triple quotes
- NO explanations, comments, or extra text
- NO "sudo" keyword (even if root is needed)
- NO echo-wrapped instructions (e.g., echo "do this" is only valid if the task is literally to print text)

Formatting Examples:

User: mount a disk
Output:
mount /dev/sdb1 /mnt

User: update and upgrade system
Output:
apt update && apt upgrade -y

User: create dir and enter it
Output:
mkdir -p mydir && cd mydir

Now generate the command for the user's request.
`

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	System string `json:"system"`
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
		System: ollamaSystemPrompt,
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
