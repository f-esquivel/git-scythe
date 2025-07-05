package main

import (
	"fmt"
	"os"

	"github.com/franklinesquivel/git-scythe/internal/git_utils"
)

func main() {
	// Get the merged branches.
	branches, err := git_utils.GetMergedBranches()

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
