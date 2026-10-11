package main

import (
	"fmt"
	"strings"
	"unicode"
)

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	part1 := 0
	part2 := 0

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		v1, err := processLinePart1(line)
		if err != nil {
			return 0, 0, err
		}
		part1 += v1

		v2, err := processLinePart2(line)
		if err != nil {
			return 0, 0, err
		}
		part2 += v2
	}

	return part1, part2, nil
}

func processLinePart1(line string) (int, error) {
	var first, last rune

	for _, ch := range line {
		if unicode.IsDigit(ch) {
			if first == 0 {
				first = ch
			}
			last = ch
		}
	}

	if first == 0 {
		return 0, nil
	}

	return int((first-'0')*10 + (last - '0')), nil
}

func processLinePart2(line string) (int, error) {
	digits := []int{}

	for i := 0; i < len(line); i++ {
		if ch := line[i]; ch >= '0' && ch <= '9' {
			digits = append(digits, int(ch-'0'))
		} else if strings.HasPrefix(line[i:], "one") {
			digits = append(digits, 1)
		} else if strings.HasPrefix(line[i:], "two") {
			digits = append(digits, 2)
		} else if strings.HasPrefix(line[i:], "three") {
			digits = append(digits, 3)
		} else if strings.HasPrefix(line[i:], "four") {
			digits = append(digits, 4)
		} else if strings.HasPrefix(line[i:], "five") {
			digits = append(digits, 5)
		} else if strings.HasPrefix(line[i:], "six") {
			digits = append(digits, 6)
		} else if strings.HasPrefix(line[i:], "seven") {
			digits = append(digits, 7)
		} else if strings.HasPrefix(line[i:], "eight") {
			digits = append(digits, 8)
		} else if strings.HasPrefix(line[i:], "nine") {
			digits = append(digits, 9)
		}
	}

	if len(digits) == 0 {
		return 0, fmt.Errorf("no digits found in line: %s", line)
	}

	return digits[0]*10 + digits[len(digits)-1], nil
}
