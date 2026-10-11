package main

import (
	"fmt"
	"strconv"
	"strings"
)

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return 0, 0, fmt.Errorf("no depth measurements found")
	}

	depths := make([]int, 0, len(lines))
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		depth, err := strconv.Atoi(line)
		if err != nil {
			return 0, 0, fmt.Errorf("line %d: invalid depth %q: %w", i+1, line, err)
		}
		depths = append(depths, depth)
	}

	if len(depths) == 0 {
		return 0, 0, fmt.Errorf("no valid depth measurements found")
	}

	increases := 0
	for i := 1; i < len(depths); i++ {
		if depths[i] > depths[i-1] {
			increases++
		}
	}
	part1 := increases

	windowIncreases := 0
	for i := 3; i < len(depths); i++ {
		if depths[i] > depths[i-3] {
			windowIncreases++
		}
	}
	part2 := windowIncreases

	return part1, part2, nil
}
