package day11

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	grid, err := parseGrid(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	for {
		next := stepGrid(grid, 4, countAdjacentOccupied)
		if equal(grid, next) {
			break
		}
		grid = next
	}

	return countOccupied(grid), nil
}
