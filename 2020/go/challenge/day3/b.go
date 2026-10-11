package day3

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	lines := input.LinesSlice()
	if len(lines) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	slopes := [][2]int{{1, 1}, {1, 3}, {1, 5}, {1, 7}, {2, 1}}
	result := 1

	for _, slope := range slopes {
		down := slope[0]
		right := slope[1]
		trees := slopeTreeCount(lines, right, down)
		result *= trees
	}

	return result, nil
}
