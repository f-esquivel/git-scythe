package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

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

	branches, err := git.GetMergedBranches()

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	if len(branches) == 0 {
		fmt.Println("No merged branches to display.")
		os.Exit(0)
	}

	fmt.Println("The following merged branches will be deleted:")
	for _, branch := range branches {
		fmt.Println("  >", branch)
	}

	fmt.Print("\nDo you want to continue? (Y/n) ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "Y" || input == "y" || input == "" {
		sucessfulDeletions := 0
		fmt.Printf("\nDeleting %d branches...\n", len(branches))

		for _, branch := range branches {
			err := git.DeleteBranch(branch, false)
			if err != nil {
				fmt.Printf("  > Error deleting branch %s: %v\n", branch, err)
			} else {
				sucessfulDeletions++
			}
		}

		fmt.Printf("Deleted %d branches.\n", sucessfulDeletions)
	} else {
		fmt.Println("Operation cancelled.")
	}
}
