package main

import (
	"fmt"
	"strconv"
	"strings"
)

type report []int

func parseReport(line string) (report, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, fmt.Errorf("report: empty line")
	}
	levels := make(report, len(fields))
	for i, field := range fields {
		level, err := strconv.Atoi(field)
		if err != nil {
			return nil, fmt.Errorf("report: level %q: %w", field, err)
		}
		levels[i] = level
	}
	return levels, nil
}

func (r report) isSafe() bool {
	if len(r) < 2 {
		return true
	}
	increasing := r[1] > r[0]
	for i := 1; i < len(r); i++ {
		diff := r[i] - r[i-1]
		magnitude := diff
		if magnitude < 0 {
			magnitude = -magnitude
		}
		if magnitude < 1 || magnitude > 3 {
			return false
		}
		if increasing && diff < 0 {
			return false
		}
		if !increasing && diff > 0 {
			return false
		}
	}
	return true
}

func (r report) isSafeWithDampener() bool {
	if r.isSafe() {
		return true
	}
	without := make(report, len(r)-1)
	for skip := range r {
		without = without[:0]
		for i, level := range r {
			if i != skip {
				without = append(without, level)
			}
		}
		if without.isSafe() {
			return true
		}
	}
	return false
}

func solve(input string) (safeCount, safeWithDampenerCount int, err error) {
	for _, line := range strings.Split(strings.ReplaceAll(strings.TrimRight(input, "\n"), "\r\n", "\n"), "\n") {
		if line == "" {
			continue
		}
		levels, err := parseReport(line)
		if err != nil {
			return 0, 0, err
		}
		if levels.isSafe() {
			safeCount++
		}
		if levels.isSafeWithDampener() {
			safeWithDampenerCount++
		}
	}
	return safeCount, safeWithDampenerCount, nil
}
