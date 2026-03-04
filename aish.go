package aish

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type client interface {
	Generate(model string, prompt string) (string, error)
}

const basePrompt = `
You are a strict Linux bash command generator.

Rules:
- You MUST return exactly ONE valid bash command.
- Output ONLY the command.
- No explanations, no markdown, no comments, no backticks, no extra text.
- Never include sudo
- If the request is not related to system or terminal operations,
  return exactly this command: echo "Unsupported request"

User request:
%s
`

const unsupportedRequest = `echo "Unsupported request"`

type Aish struct {
	model       string
	client      client
	execCommand func(name string, arg ...string) *exec.Cmd
}

func New(model string, client client) *Aish {
	return &Aish{
		model:       model,
		client:      client,
		execCommand: exec.Command,
	}
}

func (a *Aish) Query(text string) (string, error) {
	resp, err := a.client.Generate(a.model, fmt.Sprintf(basePrompt, text))
	if err != nil {
		return "", fmt.Errorf("generate response: %w", err)
	}

	if strings.Contains(resp, unsupportedRequest) {
		return "", fmt.Errorf("unsupported request")
	}

	return resp, nil
}

func (a *Aish) Execute(cmd string) error {
	command := a.execCommand("/bin/bash", "-c", cmd) //nolint:gosec
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if err := command.Run(); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}
