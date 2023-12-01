package example

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
	"github.com/oscarteg/advent-of-code-2020/util"
	"github.com/spf13/cobra"
)

func aCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "a",
		Short: "Example Day, Problem A",
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Answer: %d\n", a(challenge.FromFile()))
		},
	}
}

func a(input *challenge.Input) (result int) {
	return util.MustAtoI(<-input.Lines())
}
