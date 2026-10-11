package day5

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func A(input *challenge.Input) (int, error) {
	lines := input.LinesSlice()
	var seats []int

	for _, line := range lines {
		if line == "" {
			continue
		}

		seat, err := parseSeat(line)
		if err != nil {
			return 0, err
		}
		seats = append(seats, seat)
	}

	maxID := 0
	for _, id := range seats {
		if id > maxID {
			maxID = id
		}
	}

	return maxID, nil
}

func parseSeat(line string) (int, error) {
	if len(line) != 10 {
		return 0, fmt.Errorf("invalid seat format: expected 10 chars, got %d", len(line))
	}

	rowPart := line[:7]
	colPart := line[7:]

	for _, c := range rowPart {
		if c != 'F' && c != 'B' {
			return 0, fmt.Errorf("invalid row char: %c", c)
		}
	}

	for _, c := range colPart {
		if c != 'L' && c != 'R' {
			return 0, fmt.Errorf("invalid col char: %c", c)
		}
	}

	row := 0
	for _, c := range rowPart {
		row = row*2 + boolToInt(c == 'B')
	}

	col := 0
	for _, c := range colPart {
		col = col*2 + boolToInt(c == 'R')
	}

	if row >= 128 {
		return 0, fmt.Errorf("row out of range: %d", row)
	}
	if col >= 8 {
		return 0, fmt.Errorf("col out of range: %d", col)
	}

	return row*8 + col, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
