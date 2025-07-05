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

func run(git types.Gitter, reader io.Reader, writer io.Writer) error {
	branches, err := git.GetMergedBranches()
	if err != nil {
		return err
	}

	if len(branches) == 0 {
		fmt.Fprintln(writer, "No merged branches to display.")
		return nil
	}

	fmt.Fprintln(writer, "The following merged branches will be deleted:")
	for _, branch := range branches {
		fmt.Fprintln(writer, "  >", branch)
	}

	fmt.Fprint(writer, "\nDo you want to continue? (Y/n) ")
	bufReader := bufio.NewReader(reader)
	input, _ := bufReader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "Y" || input == "y" || input == "" {
		fmt.Fprintf(writer, "\nDeleting %d branches...\n", len(branches))
		successfulDeletions := 0
		for _, branch := range branches {
			err := git.DeleteBranch(branch, false)
			if err != nil {
				fmt.Fprintf(writer, "  > Error deleting branch %s: %v\n", branch, err)
			} else {
				successfulDeletions++
			}
		}
		fmt.Fprintf(writer, "Deleted %d branches.\n", successfulDeletions)
	} else {
		fmt.Fprintln(writer, "Operation cancelled.")
	}

	return nil
}

func main() {
	git := git_utils.New(OsCommandExecutor{})
	if err := run(git, os.Stdin, os.Stdout); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
