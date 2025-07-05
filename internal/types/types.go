package types

import "io"

// CommandExecutor defines an interface for creating and running commands.
// This allows for mocking the command execution in tests.
type CommandExecutor interface {
	Command(name string, arg ...string) Cmd
}

// Cmd defines an interface for a command that can be executed.
// This abstracts the *exec.Cmd type for testing purposes.
type Cmd interface {
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	Output() ([]byte, error)
	SetStdin(stdin io.Reader)
}