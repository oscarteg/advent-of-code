package day10

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	chain, err := adapterChain(input)
	if err != nil {
		return 0, err
	}

	ones, threes, err := joltageGaps(chain)
	if err != nil {
		return 0, err
	}

	return ones * threes, nil
}

// parseJoltages reads one integer joltage per line.
func parseJoltages(input *challenge.Input) ([]int, error) {
	var joltages []int
	for _, line := range input.LinesSlice() {
		joltage, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid joltage %q: %w", line, err)
		}
		joltages = append(joltages, joltage)
	}
	return joltages, nil
}

// adapterChain sorts the adapter joltages and appends the device joltage
// (highest + 3), with the built-in outlet (0) implied as the starting point.
func adapterChain(input *challenge.Input) ([]int, error) {
	joltages, err := parseJoltages(input)
	if err != nil {
		return nil, err
	}
	if len(joltages) == 0 {
		return nil, fmt.Errorf("no adapters given")
	}

	sort.Ints(joltages)
	chain := append(joltages, joltages[len(joltages)-1]+3)
	return chain, nil
}

// joltageGaps counts the 1-jolt and 3-jolt gaps between consecutive
// adapters in the chain, starting from the outlet at 0.
func joltageGaps(chain []int) (ones, threes int, err error) {
	previous := 0
	for _, joltage := range chain {
		switch joltage - previous {
		case 1:
			ones++
		case 2:
		case 3:
			threes++
		default:
			return 0, 0, fmt.Errorf("unreachable gap of %d jolts after %d", joltage-previous, previous)
		}
		previous = joltage
	}
	return ones, threes, nil
}
