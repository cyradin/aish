package main

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Model       string        `envconfig:"AISH_MODEL" required:"true"`
	HTTPTimeout time.Duration `envconfig:"AISH_HTTP_REQUEST_TIMEOUT" default:"120s"`
	Ollama      OllamaConfig
}

type OllamaConfig struct {
	URL string `envconfig:"AISH_OLLAMA_URL" required:"true"`
}

func parseConfig() (*Config, error) {
	var cfg Config

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("configuration error: %w", err)
	}

	return &cfg, nil
}
