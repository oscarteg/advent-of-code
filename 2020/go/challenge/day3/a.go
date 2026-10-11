package day3

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	lines := input.LinesSlice()
	if len(lines) == 0 {
		return 0, fmt.Errorf("empty input")
	}

	trees := slopeTreeCount(lines, 3, 1)
	return trees, nil
}

func slopeTreeCount(grid []string, right, down int) int {
	trees := 0
	row := 0
	col := 0

	for row < len(grid) {
		line := grid[row]
		width := len(line)
		if line[col%width] == '#' {
			trees++
		}
		row += down
		col += right
	}

	return trees
}
