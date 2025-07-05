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
	remoteOut, err := g.executor.Command("git", "remote").Output()
	if err != nil {
		return "", err
	}

	// The purpose of this block is to identify the default branch of the
	// repository. To do this, we first retrieve the list of remotes, then
	// pick the first one as the default. If the list is not empty and the
	// first remote is not origin or upstream, we use origin or upstream
	// as the default remote to get the default branch from. If none of
	// these are available, we just take the first remote.
	remotes := strings.Split(strings.TrimSpace(string(remoteOut)), "\n")
	remote := "origin"

	if len(remotes) > 0 && remotes[0] != "" {
		// Prefer origin, then upstream, then the first one in the list
		hasOrigin := false
		hasUpstream := false
		for _, r := range remotes {
			if r == "origin" {
				hasOrigin = true
				break
			}
			if r == "upstream" {
				hasUpstream = true
			}
		}

		if hasOrigin {
			remote = "origin"
		} else if hasUpstream {
			remote = "upstream"
		} else {
			remote = remotes[0]
		}
	}

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
		g.executor.Command("git", "branch", "--merged"),
		g.executor.Command("grep", "-v", "-e", "*", "-e", defaultBranch),
	}

	out, err := cmd_utils.PipeCmds(cmds)
	if err != nil {
		// If `grep` returns a non-zero exit code, it means no matches were
		// found. In this case, we can assume that no branches need to be
		// deleted, so we return an empty list of branches.
		return []string{}, nil
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
