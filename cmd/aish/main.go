package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/cyradin/aish"
	"github.com/fatih/color"
)

func main() {
	if err := run(); err != nil {
		_, _ = color.New(color.FgRed).Fprintf(os.Stderr, "error: %s\n", err)

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

	_, _ = color.New(color.FgWhite).Fprintf(os.Stdout, "suggested command:\n")
	_, _ = color.New(color.FgYellow).Fprintln(os.Stdout, suggestedCommand)
	_, _ = color.New(color.FgWhite).Fprintf(os.Stdout, "execute? (")
	_, _ = color.New(color.FgGreen).Fprintf(os.Stdout, "Y")
	_, _ = color.New(color.FgWhite).Fprintf(os.Stdout, "/n)\n")

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')

	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "" && answer != "y" && answer != "yes" {
		_, _ = color.New(color.FgRed).Fprintf(os.Stdout, "Aborted\n")

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
