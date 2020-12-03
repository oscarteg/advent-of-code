package day1

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/util"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/spf13/cobra"
)

func bCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "b",
		Short: "Day 1, Problem B",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", b(challenge.FromFile()))
		},
	}
}

func b(challenge *challenge.Input) int {
	var entries []int

	for v := range challenge.Lines() {
		entries = append(entries, util.MustAtoI(v))
	}

	elementMap := util.SliceToMap(entries)

	for index, element := range entries {
		for j := index; j <= len(entries)-1; j++ {
			sum := element + entries[j]
			if missingElement, exists := elementMap[2020-sum]; exists {
				return element * missingElement * entries[j]
			}
		}
	}

	panic("There is no answer")
}
