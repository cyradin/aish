package aish

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAish_Query(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		client    client
		text      string
		want      string
		wantError bool
	}{
		{
			name: "normal response",
			client: &mockClient{
				resp: `ls -l`,
			},
			text:      "list files",
			want:      "ls -l",
			wantError: false,
		},
		{
			name: "client error",
			client: &mockClient{
				err: errors.New("network error"),
			},
			text:      "list files",
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := New("test-model", tt.client)

			out, err := a.Query(tt.text)

			if tt.wantError {
				require.Error(t, err)
				require.Empty(t, out)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, out)
		})
	}
}

func TestAish_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		cmd       string
		wantError bool
	}{
		{
			name:      "successful command",
			cmd:       `echo hello`,
			wantError: false,
		},
		{
			name:      "failing command",
			cmd:       `false`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			a := New("", nil)
			a.execCommand = func(name string, args ...string) *exec.Cmd {
				return exec.Command("/bin/bash", "-c", tt.cmd)
			}

			err := a.Execute(tt.cmd)
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

type mockClient struct {
	resp string
	err  error
}

func (m *mockClient) Generate(model, prompt string) (string, error) {
	return m.resp, m.err
}
