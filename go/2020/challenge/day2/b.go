package day2

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/spf13/cobra"
)

func bCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "b",
		Short: "Day 2, Problem B",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", b(challenge.FromFile()))
		},
	}
}

func b(challenge *challenge.Input) int {
	var valid int

	for line := range challenge.Lines() {
		split := strings.Fields(line)

		boundaries := strings.Split(split[0], "-")
		minChar, _ := strconv.Atoi(boundaries[0])
		maxChar, _ := strconv.Atoi(boundaries[1])

		find := strings.Replace(split[1], ":", "", -1)

		password := split[2]

		firstChar := string(password[minChar-1])
		secondChar := string(password[maxChar-1])

		if firstChar == find && secondChar == find {
			continue
		}

		if firstChar == find || secondChar == find {
			valid += 1
		}

	}

	return valid
}
