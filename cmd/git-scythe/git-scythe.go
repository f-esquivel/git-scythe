package main

import (
	"fmt"

	"github.com/franklinesquivel/git-scythe/internal/git"
)

func main() {
	branches, err := git.GetMergedBranches()

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, branch := range branches {
		fmt.Println(branch)
	}
}
