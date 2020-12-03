package day2

import (
	"fmt"
	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func aCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "a",
		Short: "Day 2, Problem A",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", a(challenge.FromFile()))
		},
	}
}
func a(challenge *challenge.Input) int {

	var valid int

	for line := range challenge.Lines() {
		split := strings.Fields(line)

		boundaries := strings.Split(split[0], "-")
		minChar,_ := strconv.Atoi(boundaries[0])
		maxChar,_ := strconv.Atoi(boundaries[1])

		find := strings.Replace(split[1], ":", "", -1)


		if numberOfOccurrence := strings.Count(split[2], find); numberOfOccurrence >= minChar && numberOfOccurrence <= maxChar{
			valid += 1
		}
	}

	return valid

}
