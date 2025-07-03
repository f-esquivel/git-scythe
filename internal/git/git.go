package git

import (
	"os/exec"
	"strings"
)

func getDefaultBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "origin/HEAD")

	out, err := cmd.Output()

	if err != nil {
		return "", err
	}

	replacedName := strings.ReplaceAll(string(out), "origin/", "")
	return strings.TrimSpace(replacedName), nil
}

func GetMergedBranches() ([]string, error) {
	defaultBranch, err := getDefaultBranch()

	if err != nil {
		return nil, err
	}

	gitMergedBranchesCmd := exec.Command("git", "branch", "--merged")
	grepCmd := exec.Command("grep", "-v", "-e", "*", "-e", defaultBranch)

	pipe, err := gitMergedBranchesCmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

	defer pipe.Close()

	grepCmd.Stdin = pipe

	if err := gitMergedBranchesCmd.Start(); err != nil {
		return nil, err
	}

	out, err := grepCmd.Output()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	var trimmedBranches []string
	for _, line := range lines {
		if line != "" {
			trimmedBranches = append(trimmedBranches, strings.TrimSpace(line))
		}
	}

	return trimmedBranches, nil
}
