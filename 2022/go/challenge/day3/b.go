package day3

import (
	"fmt"
	"sync"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/spf13/cobra"
)

func bCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "b",
		Short: "Day 3, Problem B",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", b(challenge.FromFile()))
		},
	}
}
func b(challenge *challenge.Input) int {

	grid := challenge.ToGrid()

	slopes := [][]int{{1, 1}, {3, 1}, {5, 1}, {7, 1}, {1, 2}}
	treesMultiplied := 1
	var wg sync.WaitGroup

	for _, slope := range slopes {
		wg.Add(1)
		go func(slope []int) {
			defer wg.Done()
			treesMultiplied *= grid.CountTreesInPath(slope[0], slope[1])
		}(slope)
	}

	wg.Wait()
	return treesMultiplied

}
