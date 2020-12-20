package day1

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/util"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/spf13/cobra"
)

func aCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "a",
		Short: "Day 1, Problem A",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", a(challenge.FromFile()))
		},
	}
}
func a(challenge *challenge.Input) int {

	var entries []int

	for v := range challenge.Lines() {
		entries = append(entries, util.MustAtoI(v))
	}

	for i := 0; i <= len(entries); i++ {
		for j := i; j <= len(entries)-1; j++ {
			sum := entries[i] + entries[j]
			if sum == 2020 {
				return entries[i] * entries[j]
			}
		}
	}

	panic("There is no solution")
}
