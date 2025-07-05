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

func selectPreferredRemote(remotes []string) string {
	if len(remotes) == 0 || remotes[0] == "" {
		return "origin"
	}

	// Check for preferred remotes in order
	preferredRemotes := []string{"origin", "upstream"}
	remoteSet := make(map[string]bool)
	for _, r := range remotes {
		remoteSet[r] = true
	}

	for _, preferred := range preferredRemotes {
		if remoteSet[preferred] {
			return preferred
		}
	}

	return remotes[0]
}

func (g *GitUtils) GetDefaultBranch() (string, error) {
	remoteOut, err := g.executor.Command("git", "remote").Output()
	if err != nil {
		return "", err
	}

	remotes := strings.Split(strings.TrimSpace(string(remoteOut)), "\n")
	remote := selectPreferredRemote(remotes)

	cmd := g.executor.Command("git", "rev-parse", "--abbrev-ref", remote+"/HEAD")

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	sanitizedName := strings.TrimSpace(strings.ReplaceAll(string(out), remote+"/", ""))
	return sanitizedName, nil
}

func (g *GitUtils) GetMergedBranches() ([]string, error) {
	defaultBranch, err := g.GetDefaultBranch()
	if err != nil {
		return nil, err
	}

	cmds := []types.Cmd{
		g.executor.Command("git", "branch", "--merged", defaultBranch),
		g.executor.Command("grep", "-v", "-e", "*", "-e", defaultBranch),
	}

	out, err := cmd_utils.PipeCmds(cmds)
	if err != nil {
		// grep returns exit code 1 when no lines match, which is expected
		// when there are no merged branches. Other errors should be propagated.
		if strings.Contains(err.Error(), "exit status 1") {
			return []string{}, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	trimmedBranches := make([]string, 0)
	for _, line := range lines {
		if line != "" {
			trimmedBranches = append(trimmedBranches, strings.TrimSpace(line))
		}
	}
	return trimmedBranches, nil
}

func (g *GitUtils) DeleteBranch(branchName string, force bool) error {
	args := []string{"branch"}
	if force {
		args = append(args, "-D")
	} else {
		args = append(args, "-d")
	}
	args = append(args, branchName)

	cmd := g.executor.Command("git", args...)
	_, err := cmd.Output()
	return err
}
