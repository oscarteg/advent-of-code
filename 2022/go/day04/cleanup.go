package main

import (
	"fmt"
	"strings"
)

type Range struct {
	start int
	end   int
}

func parseRange(s string) (Range, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return Range{}, fmt.Errorf("invalid range format: %s", s)
	}

	var start, end int
	_, err := fmt.Sscanf(parts[0], "%d", &start)
	if err != nil {
		return Range{}, fmt.Errorf("invalid start number: %w", err)
	}

	_, err = fmt.Sscanf(parts[1], "%d", &end)
	if err != nil {
		return Range{}, fmt.Errorf("invalid end number: %w", err)
	}

	return Range{start: start, end: end}, nil
}

func contains(a, b Range) bool {
	return a.start <= b.start && a.end >= b.end
}

func overlaps(a, b Range) bool {
	max := a.start
	if b.start > max {
		max = b.start
	}
	min := a.end
	if b.end < min {
		min = b.end
	}
	return max <= min
}

func solve(input string) (int, int, error) {
	var part1, part2 int

	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid line format: %s", line)
		}

		a, err := parseRange(parts[0])
		if err != nil {
			return 0, 0, err
		}

		b, err := parseRange(parts[1])
		if err != nil {
			return 0, 0, err
		}

		if contains(a, b) || contains(b, a) {
			part1++
		}

		if overlaps(a, b) {
			part2++
		}
	}

	return part1, part2, nil
}
