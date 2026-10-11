package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func elfTotals(input string) ([]int, error) {
	trimmed := strings.Trim(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	if trimmed == "" {
		return nil, fmt.Errorf("input is empty")
	}

	blocks := strings.Split(trimmed, "\n\n")
	totals := make([]int, 0, len(blocks))
	for _, block := range blocks {
		total := 0
		for _, line := range strings.Split(block, "\n") {
			calories, err := strconv.Atoi(strings.TrimSpace(line))
			if err != nil {
				return nil, fmt.Errorf("invalid calorie value %q: %w", line, err)
			}
			total += calories
		}
		totals = append(totals, total)
	}
	return totals, nil
}

func solve(input string) (int, int, error) {
	totals, err := elfTotals(input)
	if err != nil {
		return 0, 0, err
	}

	sort.Sort(sort.Reverse(sort.IntSlice(totals)))

	part1 := totals[0]
	part2 := 0
	for _, total := range totals[:min(3, len(totals))] {
		part2 += total
	}
	return part1, part2, nil
}
