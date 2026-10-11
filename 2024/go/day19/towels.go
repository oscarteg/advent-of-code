package main

import (
	"fmt"
	"strings"
)

func solve(input string) (int, int, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return 0, 0, fmt.Errorf("empty input")
	}
	lines := strings.Split(trimmed, "\n")
	if len(lines) == 0 {
		return 0, 0, fmt.Errorf("empty input")
	}

	patternStrs := strings.Split(lines[0], ",")
	patterns := make([]string, 0, len(patternStrs))
	for _, p := range patternStrs {
		p = strings.TrimSpace(p)
		if len(p) > 0 {
			patterns = append(patterns, p)
		}
	}

	var part1 int
	var part2 int

	for i := 1; i < len(lines); i++ {
		design := strings.TrimSpace(lines[i])
		if len(design) == 0 {
			continue
		}

		memo := make(map[string]int64)
		ways := countWays(design, patterns, memo)
		if ways > 0 {
			part1++
			part2 += int(ways)
		}
	}

	return part1, part2, nil
}

func countWays(remaining string, patterns []string, memo map[string]int64) int64 {
	if remaining == "" {
		return 1
	}

	if cached, ok := memo[remaining]; ok {
		return cached
	}

	var ways int64
	for _, pattern := range patterns {
		if strings.HasPrefix(remaining, pattern) {
			ways += countWays(remaining[len(pattern):], patterns, memo)
		}
	}

	memo[remaining] = ways
	return ways
}
