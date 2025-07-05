package cmd_utils

import (
	"github.com/franklinesquivel/git-scythe/internal/types"
)

// PipeCmds executes a series of commands, piping the stdout of each command
// to the stdin of the next. It is a powerful utility for chaining command-line
// tools, similar to how pipes work in a shell (e.g., `cmd1 | cmd2 | cmd3`).
//
// The function takes a slice of Cmd objects, where each command is
// expected to be fully configured with its arguments. It then orchestrates
// the execution of these commands in a pipeline.
//
// Key behaviors:
//   - It connects the stdout of cmds[i] to the stdin of cmds[i+1].
//   - It starts all commands but waits to capture the output of only the final
//     command in the chain.
//   - It explicitly calls `Wait()` on all preceding commands to ensure that
//     all resources are properly released and to check for any errors that
//     may have occurred in the middle of the pipeline.
//
// This approach ensures that the entire pipeline runs efficiently and that
// the program cleans up all child processes, preventing resource leaks.
//
// It returns the final output as a string and an error if any command
// in the pipeline fails.
func PipeCmds(cmds []types.Cmd) (string, error) {
	if len(cmds) == 0 {
		return "", nil
	}

	// Connect the stdout of each command to the stdin of the next.
	for i := 0; i < len(cmds)-1; i++ {
		stdout, err := cmds[i].StdoutPipe()
		if err != nil {
			return "", err
		}
		cmds[i+1].SetStdin(stdout)
	}

	// Start all commands except for the last one.
	// The last command will be started by the Output() method.
	for i := 0; i < len(cmds)-1; i++ {
		if err := cmds[i].Start(); err != nil {
			return "", err
		}
	}

	// Run the last command and capture its output.
	// The Output() method waits for the command to complete.
	out, err := cmds[len(cmds)-1].Output()
	if err != nil {
		return "", err
	}

	// Wait for all other commands to complete.
	// This is important to release resources and check for errors.
	for i := 0; i < len(cmds)-1; i++ {
		if err := cmds[i].Wait(); err != nil {
			return "", err
		}
	}

	return string(out), nil
}
