package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/franklinesquivel/git-scythe/internal/git_utils"
	"github.com/franklinesquivel/git-scythe/internal/types"
)

// OsCommandExecutor implements CommandExecutor using the real os/exec package.
type OsCommandExecutor struct{}

// Command runs a command and returns a os/exec.Cmd object.
func (e OsCommandExecutor) Command(name string, arg ...string) types.Cmd {
	return &execCmd{exec.Command(name, arg...)}
}

type execCmd struct {
	*exec.Cmd
}

func (c *execCmd) SetStdin(stdin io.Reader) {
	c.Stdin = stdin
}

func main() {
	git := git_utils.New(OsCommandExecutor{})

	// Get the merged branches.
	branches, err := git.GetMergedBranches()

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(branches) == 0 {
		fmt.Println("No merged branches to display.")
		os.Exit(0)
	}

	fmt.Println("Merged branches:")
	for _, branch := range branches {
		fmt.Println("  > ", branch)
	}
}
