package main

import (
	"fmt"
	"strings"
)

func solve(input string) (int64, int64, error) {
	if strings.TrimSpace(input) == "" {
		return 0, 0, fmt.Errorf("empty input")
	}

	part1Sum := int64(0)
	part2Sum := int64(0)
	enabled := true
	pos := 0

	for pos < len(input) {
		nextMul := strings.Index(input[pos:], "mul(")
		nextDo := strings.Index(input[pos:], "do()")
		nextDont := strings.Index(input[pos:], "don't()")

		nextPos := len(input) - pos
		action := -1

		if nextMul >= 0 {
			nextPos = nextMul
			action = 0
		}
		if nextDo >= 0 && (action < 0 || nextDo < nextPos) {
			nextPos = nextDo
			action = 1
		}
		if nextDont >= 0 && (action < 0 || nextDont < nextPos) {
			nextPos = nextDont
			action = 2
		}

		if action < 0 {
			break
		}

		actualPos := pos + nextPos

		switch action {
		case 0:
			substr := input[actualPos:]
			if pair, ok := tryParseMul(substr); ok {
				part1Sum += int64(pair.x * pair.y)
				if enabled {
					part2Sum += int64(pair.x * pair.y)
				}
				pos = actualPos + pair.len
			} else {
				pos = actualPos + 4
			}
		case 1:
			enabled = true
			pos = actualPos + 4
		case 2:
			enabled = false
			pos = actualPos + 7
		}
	}

	return part1Sum, part2Sum, nil
}

type mulResult struct {
	x   int
	y   int
	len int
}

func tryParseMul(substr string) (mulResult, bool) {
	if len(substr) < 8 {
		return mulResult{}, false
	}
	if !strings.HasPrefix(substr, "mul(") {
		return mulResult{}, false
	}

	i := 4
	var x int
	digitCount := 0

	for i < len(substr) && digitCount < 3 {
		c := substr[i]
		if c >= '0' && c <= '9' {
			x = x*10 + int(c-'0')
			digitCount++
			i++
		} else {
			break
		}
	}

	if digitCount == 0 || digitCount > 3 || i >= len(substr) {
		return mulResult{}, false
	}
	if substr[i] != ',' {
		return mulResult{}, false
	}

	i++
	var y int
	digitCount = 0

	for i < len(substr) && digitCount < 3 {
		c := substr[i]
		if c >= '0' && c <= '9' {
			y = y*10 + int(c-'0')
			digitCount++
			i++
		} else {
			break
		}
	}

	if digitCount == 0 || digitCount > 3 || i >= len(substr) {
		return mulResult{}, false
	}
	if substr[i] != ')' {
		return mulResult{}, false
	}

	return mulResult{x: x, y: y, len: i + 1}, true
}
