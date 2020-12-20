package day3

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/spf13/cobra"
)

func aCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "a",
		Short: "Day 3, Problem A",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", a(challenge.FromFile()))
		},
	}
}
func a(challenge *challenge.Input) int {
	return 0
}
