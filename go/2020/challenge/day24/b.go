package day24

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/spf13/cobra"
)

func bCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "b",
		Short: "Day 24, Problem B",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", b(challenge.FromFile()))
		},
	}
}
func b(challenge *challenge.Input) int {
	return 0
}
