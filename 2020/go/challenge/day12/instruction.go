package day12

import (
	"fmt"
	"strconv"
)

type action int

const (
	north action = iota
	south
	east
	west
	left
	right
	forward
)

type instruction struct {
	action action
	value  int
}

func parseInstructions(lines []string) ([]instruction, error) {
	var instructions []instruction

	for _, line := range lines {
		if line == "" {
			continue
		}

		instr, err := parseInstruction(line)
		if err != nil {
			return nil, err
		}

		instructions = append(instructions, instr)
	}

	if len(instructions) == 0 {
		return nil, fmt.Errorf("input contains no instructions")
	}

	return instructions, nil
}

func parseInstruction(line string) (instruction, error) {
	if len(line) < 2 {
		return instruction{}, fmt.Errorf("instruction line too short: %q", line)
	}

	value, err := strconv.Atoi(line[1:])
	if err != nil {
		return instruction{}, fmt.Errorf("invalid instruction value in line %q: %w", line, err)
	}

	var act action
	switch line[0] {
	case 'N':
		act = north
	case 'S':
		act = south
	case 'E':
		act = east
	case 'W':
		act = west
	case 'L':
		if value%90 != 0 {
			return instruction{}, fmt.Errorf("left turn must be a multiple of 90 degrees: %q", line)
		}
		act = left
	case 'R':
		if value%90 != 0 {
			return instruction{}, fmt.Errorf("right turn must be a multiple of 90 degrees: %q", line)
		}
		act = right
	case 'F':
		act = forward
	default:
		return instruction{}, fmt.Errorf("unknown instruction action %q in line %q", line[0:1], line)
	}

	return instruction{action: act, value: value}, nil
}

// rotate rotates a vector (x, y) counter-clockwise by degrees (a multiple of 90).
func rotate(x, y, degrees int) (int, int) {
	steps := ((degrees/90)%4 + 4) % 4
	for i := 0; i < steps; i++ {
		x, y = -y, x
	}
	return x, y
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
