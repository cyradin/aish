package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cyradin/aish"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	//nolint:mnd
	if len(os.Args) < 2 {
		return fmt.Errorf(`prompt is required (i.e. %q)`, "aish list top 5 processes by cpu usage")
	}

	cfg, err := parseConfig()
	if err != nil {
		return err
	}

	aish := initAish(cfg)

	suggestedCommand, err := aish.Query(strings.Join(os.Args[1:], " "))
	if err != nil {
		return err
	}

	suggestedCommand = strings.TrimSpace(suggestedCommand)

	_, _ = fmt.Fprintln(os.Stdout, "Suggested command:")
	_, _ = fmt.Fprintln(os.Stdout, suggestedCommand) //nolint:gosec
	_, _ = fmt.Fprint(os.Stdout, "\nExecute? Y/N: ")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		_, _ = fmt.Fprintln(os.Stdout, "Aborted")
		return nil
	}

	if err := aish.Execute(suggestedCommand); err != nil {
		return err
	}

	return nil
}

func initAish(cfg *Config) *aish.Aish {
	client := aish.NewOllamaClient(cfg.Ollama.URL, &http.Client{
		Timeout: cfg.HTTPTimeout,
	})

	return aish.New(cfg.Model, client)
}
