package day12

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	instructions, err := parseInstructions(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	// Ship starts at the origin, waypoint starts 10 east and 1 north of it.
	x, y := 0, 0
	wx, wy := 10, 1

	for _, instruction := range instructions {
		switch instruction.action {
		case north:
			wy += instruction.value
		case south:
			wy -= instruction.value
		case east:
			wx += instruction.value
		case west:
			wx -= instruction.value
		case left:
			wx, wy = rotate(wx, wy, instruction.value)
		case right:
			wx, wy = rotate(wx, wy, -instruction.value)
		case forward:
			x += wx * instruction.value
			y += wy * instruction.value
		}
	}

	return abs(x) + abs(y), nil
}
