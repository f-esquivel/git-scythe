package cmd_utils

import (
	"errors"
	"io"
	"testing"

	"github.com/franklinesquivel/git-scythe/internal/types"
)

// mockCmd is a mock implementation of the types.Cmd interface for testing.
type mockCmd struct {
	stdoutPipeFunc func() (io.ReadCloser, error)
	startFunc      func() error
	waitFunc       func() error
	outputFunc     func() ([]byte, error)
	setStdinFunc   func(io.Reader)
}

func (m *mockCmd) StdoutPipe() (io.ReadCloser, error) {
	if m.stdoutPipeFunc != nil {
		return m.stdoutPipeFunc()
	}
	return nil, nil
}

func (m *mockCmd) Start() error {
	if m.startFunc != nil {
		return m.startFunc()
	}
	return nil
}

func (m *mockCmd) Wait() error {
	if m.waitFunc != nil {
		return m.waitFunc()
	}
	return nil
}

func (m *mockCmd) Output() ([]byte, error) {
	if m.outputFunc != nil {
		return m.outputFunc()
	}
	return nil, nil
}

func (m *mockCmd) SetStdin(stdin io.Reader) {
	if m.setStdinFunc != nil {
		m.setStdinFunc(stdin)
	}
}

func compareErrors(a, b error) bool {
	if a == nil && b == nil {
		return true
	}
	if a != nil && b != nil {
		return a.Error() == b.Error()
	}
	return false
}

func TestPipeCmds(t *testing.T) {
	tests := []struct {
		name        string
		cmds        []types.Cmd
		expectedOut string
		expectedErr error
	}{
		{
			name:        "no commands",
			cmds:        []types.Cmd{},
			expectedOut: "",
			expectedErr: nil,
		},
		{
			name: "single command success",
			cmds: []types.Cmd{
				&mockCmd{
					outputFunc: func() ([]byte, error) {
						return []byte("hello"), nil
					},
				},
			},
			expectedOut: "hello",
			expectedErr: nil,
		},
		{
			name: "single command error",
			cmds: []types.Cmd{
				&mockCmd{
					outputFunc: func() ([]byte, error) {
						return nil, errors.New("cmd error")
					},
				},
			},
			expectedOut: "",
			expectedErr: errors.New("cmd error"),
		},
		{
			name: "pipe two commands success",
			cmds: []types.Cmd{
				&mockCmd{
					stdoutPipeFunc: func() (io.ReadCloser, error) {
						return io.NopCloser(nil), nil
					},
				},
				&mockCmd{
					outputFunc: func() ([]byte, error) {
						return []byte("world"), nil
					},
				},
			},
			expectedOut: "world",
			expectedErr: nil,
		},
		{
			name: "error in stdoutpipe",
			cmds: []types.Cmd{
				&mockCmd{
					stdoutPipeFunc: func() (io.ReadCloser, error) {
						return nil, errors.New("pipe error")
					},
				},
				&mockCmd{},
			},
			expectedOut: "",
			expectedErr: errors.New("pipe error"),
		},
		{
			name: "error in start",
			cmds: []types.Cmd{
				&mockCmd{
					stdoutPipeFunc: func() (io.ReadCloser, error) {
						return io.NopCloser(nil), nil
					},
					startFunc: func() error {
						return errors.New("start error")
					},
				},
				&mockCmd{},
			},
			expectedOut: "",
			expectedErr: errors.New("start error"),
		},
		{
			name: "error in wait",
			cmds: []types.Cmd{
				&mockCmd{
					stdoutPipeFunc: func() (io.ReadCloser, error) {
						return io.NopCloser(nil), nil
					},
					waitFunc: func() error {
						return errors.New("wait error")
					},
				},
				&mockCmd{
					outputFunc: func() ([]byte, error) {
						return []byte("output"), nil
					},
				},
			},
			expectedOut: "",
			expectedErr: errors.New("wait error"),
		},
		{
			name: "empty stdout in pipe",
			cmds: []types.Cmd{
				&mockCmd{
					stdoutPipeFunc: func() (io.ReadCloser, error) {
						return io.NopCloser(nil), nil
					},
					outputFunc: func() ([]byte, error) {
						return []byte(""), nil
					},
				},
				&mockCmd{
					outputFunc: func() ([]byte, error) {
						return []byte("second command output"), nil
					},
				},
			},
			expectedOut: "second command output",
			expectedErr: nil,
		},
		{
			name: "stderr with success",
			cmds: []types.Cmd{
				&mockCmd{
					stdoutPipeFunc: func() (io.ReadCloser, error) {
						return io.NopCloser(nil), nil
					},
					waitFunc: func() error {
						// Simulate a command writing to stderr but exiting successfully
						return nil
					},
				},
				&mockCmd{
					outputFunc: func() ([]byte, error) {
						return []byte("final output"), nil
					},
				},
			},
			expectedOut: "final output",
			expectedErr: nil,
		},
		{
			name: "grep not found",
			cmds: []types.Cmd{
				&mockCmd{
					stdoutPipeFunc: func() (io.ReadCloser, error) {
						return io.NopCloser(nil), nil
					},
				},
				&mockCmd{
					outputFunc: func() ([]byte, error) {
						return nil, errors.New(`exec: "grep": executable file not found in $PATH`)
					},
				},
			},
			expectedOut: "",
			expectedErr: errors.New(`exec: "grep": executable file not found in $PATH`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := PipeCmds(tt.cmds)

			if out != tt.expectedOut {
				t.Errorf("expected output %q, got %q", tt.expectedOut, out)
			}

			if !compareErrors(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
