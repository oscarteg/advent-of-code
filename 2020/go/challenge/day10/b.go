package day10

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	chain, err := adapterChain(input)
	if err != nil {
		return 0, err
	}

	return countArrangements(chain), nil
}

// countArrangements counts the distinct valid adapter chains reaching the
// device with dynamic programming over joltage values: ways[j] is the
// number of ways to reach joltage j, the sum of ways[j-1], ways[j-2] and
// ways[j-3] for whichever of those joltages exist in the chain.
func countArrangements(chain []int) int {
	device := chain[len(chain)-1]
	ways := make([]int, device+1)
	ways[0] = 1

	for _, joltage := range chain {
		for delta := 1; delta <= 3; delta++ {
			if lower := joltage - delta; lower >= 0 {
				ways[joltage] += ways[lower]
			}
		}
	}

	return ways[device]
}
