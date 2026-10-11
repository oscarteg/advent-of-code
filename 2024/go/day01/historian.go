package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return 0, 0, fmt.Errorf("empty input")
	}

	var left, right []int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return 0, 0, fmt.Errorf("invalid line format: %s", line)
		}

		l, err := strconv.Atoi(fields[0])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid left number: %w", err)
		}
		r, err := strconv.Atoi(fields[1])
		if err != nil {
			return 0, 0, fmt.Errorf("invalid right number: %w", err)
		}

		left = append(left, l)
		right = append(right, r)
	}

	slices.Sort(left)
	slices.Sort(right)

	distance := 0
	for i := 0; i < len(left); i++ {
		diff := left[i] - right[i]
		if diff < 0 {
			diff = -diff
		}
		distance += diff
	}

	rightCounts := make(map[int]int)
	for _, r := range right {
		rightCounts[r]++
	}

	score := 0
	for _, l := range left {
		score += l * rightCounts[l]
	}

	return distance, score, nil
}
