package git_utils

import (
	"strings"

	"github.com/franklinesquivel/git-scythe/internal/cmd_utils"
	"github.com/franklinesquivel/git-scythe/internal/types"
)

type GitUtils struct {
	executor types.CommandExecutor
}

func New(executor types.CommandExecutor) *GitUtils {
	return &GitUtils{executor: executor}
}

func (g *GitUtils) GetDefaultBranch() (string, error) {
	cmd := g.executor.Command("git", "rev-parse", "--abbrev-ref", "origin/HEAD")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	sanitizedName := strings.TrimSpace(strings.ReplaceAll(string(out), "origin/", ""))
	return sanitizedName, nil
}

func (g *GitUtils) GetMergedBranches() ([]string, error) {
	defaultBranch, err := g.GetDefaultBranch()
	if err != nil {
		return nil, err
	}

	cmds := []types.Cmd{
		g.executor.Command("git", "branch", "--merged"),
		g.executor.Command("grep", "-v", "-e", "*", "-e", defaultBranch),
	}

	out, err := cmd_utils.PipeCmds(cmds)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	var trimmedBranches []string
	for _, line := range lines {
		if line != "" {
			trimmedBranches = append(trimmedBranches, strings.TrimSpace(line))
		}
	}
	return trimmedBranches, nil
}
