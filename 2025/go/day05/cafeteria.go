package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type Range struct {
	start int
	end   int
}

func solve(input string) (int, int, error) {
	parts := strings.Split(input, "\n\n")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("invalid input format: expected ranges and ids separated by blank line")
	}

	rangeLines := strings.Split(strings.TrimSpace(parts[0]), "\n")
	ranges, err := parseRanges(rangeLines)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing ranges: %w", err)
	}

	idLines := strings.Split(strings.TrimSpace(parts[1]), "\n")
	ids, err := parseIDs(idLines)
	if err != nil {
		return 0, 0, fmt.Errorf("parsing ids: %w", err)
	}

	part1 := solvePart1(ranges, ids)
	part2 := solvePart2(ranges)

	return part1, part2, nil
}

func parseRanges(lines []string) ([]Range, error) {
	var ranges []Range
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid range format: %s", line)
		}

		start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return nil, fmt.Errorf("invalid start: %w", err)
		}

		end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("invalid end: %w", err)
		}

		ranges = append(ranges, Range{start: start, end: end})
	}
	return ranges, nil
}

func parseIDs(lines []string) ([]int, error) {
	var ids []int
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		id, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("invalid id: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func isInRange(id int, ranges []Range) bool {
	for _, r := range ranges {
		if id >= r.start && id <= r.end {
			return true
		}
	}
	return false
}

func solvePart1(ranges []Range, ids []int) int {
	count := 0
	for _, id := range ids {
		if isInRange(id, ranges) {
			count++
		}
	}
	return count
}

func solvePart2(ranges []Range) int {
	if len(ranges) == 0 {
		return 0
	}

	sortedRanges := make([]Range, len(ranges))
	copy(sortedRanges, ranges)
	slices.SortFunc(sortedRanges, func(a, b Range) int {
		return a.start - b.start
	})

	merged := []Range{sortedRanges[0]}
	for _, r := range sortedRanges[1:] {
		last := &merged[len(merged)-1]

		if r.start <= last.end+1 {
			if r.end > last.end {
				last.end = r.end
			}
		} else {
			merged = append(merged, r)
		}
	}

	total := 0
	for _, r := range merged {
		total += r.end - r.start + 1
	}
	return total
}
