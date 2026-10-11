package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidRecords = errors.New("invalid condition records")

// unfoldCopies is how often part two repeats every row.
const unfoldCopies = 5

// record is one row of the condition records: the spring states ('.', '#' and
// '?') and the sizes of the contiguous groups of damaged springs.
type record struct {
	springs string
	groups  []int
}

func parseRecord(line string) (record, error) {
	fields := strings.Fields(line)
	if len(fields) != 2 {
		return record{}, fmt.Errorf("%w: row %q needs springs and group sizes", errInvalidRecords, line)
	}
	for _, symbol := range fields[0] {
		if symbol != '.' && symbol != '#' && symbol != '?' {
			return record{}, fmt.Errorf("%w: unknown spring symbol %q", errInvalidRecords, symbol)
		}
	}
	sizes := strings.Split(fields[1], ",")
	groups := make([]int, 0, len(sizes))
	for _, size := range sizes {
		group, err := strconv.Atoi(size)
		if err != nil {
			return record{}, fmt.Errorf("group size: %w", err)
		}
		if group < 1 {
			return record{}, fmt.Errorf("%w: group sizes must be positive, got %d", errInvalidRecords, group)
		}
		groups = append(groups, group)
	}
	return record{springs: fields[0], groups: groups}, nil
}

func parseRecords(input string) ([]record, error) {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	records := make([]record, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parsed, err := parseRecord(line)
		if err != nil {
			return nil, err
		}
		records = append(records, parsed)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("%w: input holds no rows", errInvalidRecords)
	}
	return records, nil
}

// unfold repeats the row copies times, joining the spring states with '?'.
func (r record) unfold(copies int) record {
	springs := make([]string, copies)
	groups := make([]int, 0, len(r.groups)*copies)
	for i := range springs {
		springs[i] = r.springs
		groups = append(groups, r.groups...)
	}
	return record{springs: strings.Join(springs, "?"), groups: groups}
}

// arrangements counts the ways the unknown springs can be resolved so that the
// damaged groups match the listed sizes. Every state is a spring position
// paired with the group that still has to be placed, so memoising those pairs
// keeps the walk polynomial instead of exponential.
func (r record) arrangements() int {
	width := len(r.groups) + 1
	cache := make([]int, (len(r.springs)+1)*width)
	for i := range cache {
		cache[i] = -1
	}
	var count func(spring, group int) int
	count = func(spring, group int) int {
		if spring == len(r.springs) {
			if group == len(r.groups) {
				return 1
			}
			return 0
		}
		key := spring*width + group
		if cache[key] >= 0 {
			return cache[key]
		}
		arrangements := 0
		if r.springs[spring] != '#' {
			arrangements += count(spring+1, group)
		}
		// A group can start here when it fits in non-operational springs and is
		// followed by an operational spring or the end of the row.
		if r.springs[spring] != '.' && group < len(r.groups) {
			end := spring + r.groups[group]
			if end <= len(r.springs) &&
				!strings.Contains(r.springs[spring:end], ".") &&
				(end == len(r.springs) || r.springs[end] != '#') {
				arrangements += count(min(end+1, len(r.springs)), group+1)
			}
		}
		cache[key] = arrangements
		return arrangements
	}
	return count(0, 0)
}

func solve(input string) (part1, part2 int, err error) {
	return solveWith(input, unfoldCopies)
}

// solveWith sums the arrangements of every row as written and after unfolding
// it copies times.
func solveWith(input string, copies int) (part1, part2 int, err error) {
	records, err := parseRecords(input)
	if err != nil {
		return 0, 0, err
	}
	for _, r := range records {
		part1 += r.arrangements()
		part2 += r.unfold(copies).arrangements()
	}
	return part1, part2, nil
}
