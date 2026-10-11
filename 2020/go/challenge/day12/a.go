package day12

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	instructions, err := parseInstructions(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	// Ship starts at the origin facing east.
	x, y := 0, 0
	dx, dy := 1, 0

	for _, instruction := range instructions {
		switch instruction.action {
		case north:
			y += instruction.value
		case south:
			y -= instruction.value
		case east:
			x += instruction.value
		case west:
			x -= instruction.value
		case left:
			dx, dy = rotate(dx, dy, instruction.value)
		case right:
			dx, dy = rotate(dx, dy, -instruction.value)
		case forward:
			x += dx * instruction.value
			y += dy * instruction.value
		}
	}

	return abs(x) + abs(y), nil
}
