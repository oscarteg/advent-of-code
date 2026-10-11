package day11

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	grid, err := parseGrid(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	for {
		next := stepGrid(grid, 5, countVisibleOccupied)
		if equal(grid, next) {
			break
		}
		grid = next
	}

	return countOccupied(grid), nil
}
