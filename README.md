# aish

[![CI](https://github.com/cyradin/aish/actions/workflows/default.yaml/badge.svg)](https://github.com/cyradin/aish/actions/workflows/default.yaml)
[![Go Reference](https://pkg.go.dev/badge/github.com/cyradin/aish.svg)](https://pkg.go.dev/github.com/cyradin/aish)
[![Go Version](https://img.shields.io/badge/go-1.26-blue.svg)](https://golang.org/)
[![Release](https://img.shields.io/github/v/release/cyradin/aish?sort=semver)](https://github.com/cyradin/aish/releases)

## Description

**aish** is a command-line tool that generates bash commands based on a text description.
It uses an Ollama API model and always returns exactly one valid command.
Aish can also safely execute the generated command after asking for user confirmation.

## How it works

- The user provides a textual description of the task.
- The tool sends a request to the Ollama API with a prompt that enforces returning only a command.
- The model responds with a bash command.
- The tool displays the command and asks if it should be executed. The default option is (Y). Pressing Enter executes the command.
- The command is executed in a shell.

__Aish does not add sudo automatically. If a command requires elevated privileges, the user is expected to rerun it using sudo.__

## Installation

Install the latest version using:

```bash
go install github.com/cyradin/aish/cmd/aish@latest
```

Make sure `$GOPATH/bin` or `$HOME/go/bin` is in your `PATH`.

## Environment Variables

Aish uses the following environment variables:

| Variable | Description | Default / Required |
|----------|-------------|--------------------|
| AISH_MODEL | The Ollama model name to use | Required |
| AISH_HTTP_REQUEST_TIMEOUT | HTTP request timeout for the Ollama API | 120s |
| AISH_OLLAMA_URL | Base URL of the Ollama API | Required |

You can set them in your shell configuration (e.g., `.bashrc` or `.zshrc`):

```bash
export AISH_MODEL="qwen2.5-coder:7b"
export AISH_OLLAMA_URL="http://localhost:11434"
export AISH_HTTP_REQUEST_TIMEOUT="120s"
```

## Usage

Run Aish with a text description:

```bash
aish "list top 5 processes by cpu usage"
```

Example output:
```bash
Suggested command:
ps aux --sort=-%cpu | head -n 6

Execute? (Y)/n:
```
Press Enter or Y to execute the command, or N to abort.
