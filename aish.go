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
	resp, err := a.client.Generate(a.model, text)
	if err != nil {
		return "", fmt.Errorf("generate response: %w", err)
	}

	resp = replacer.Replace(resp)

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

var replacer = strings.NewReplacer(
	"```bash", "",
	"```shell", "",
	"```sh", "",
	"```", "",
)
