package main

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidNotes = errors.New("invalid pattern notes")

type pattern []string

func parsePatterns(input string) ([]pattern, error) {
	normalized := strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n")
	if normalized == "" {
		return nil, fmt.Errorf("%w: notes are empty", errInvalidNotes)
	}
	blocks := strings.Split(normalized, "\n\n")
	patterns := make([]pattern, 0, len(blocks))
	for _, block := range blocks {
		rows := strings.Split(block, "\n")
		width := len(rows[0])
		for _, row := range rows {
			if len(row) != width {
				return nil, fmt.Errorf("%w: pattern rows have inconsistent width", errInvalidNotes)
			}
			for _, c := range row {
				if c != '#' && c != '.' {
					return nil, fmt.Errorf("%w: pattern contains unknown symbol %q", errInvalidNotes, c)
				}
			}
		}
		patterns = append(patterns, pattern(rows))
	}
	return patterns, nil
}

func (p pattern) column(index int) string {
	column := make([]byte, len(p))
	for i, row := range p {
		column[i] = row[index]
	}
	return string(column)
}

func diffCount(a, b string) int {
	diff := 0
	for i := range a {
		if a[i] != b[i] {
			diff++
		}
	}
	return diff
}

func horizontalReflection(p pattern, smudges int) (int, bool) {
	for line := 1; line < len(p); line++ {
		span := line
		if len(p)-line < span {
			span = len(p) - line
		}
		diff := 0
		for offset := 0; offset < span; offset++ {
			diff += diffCount(p[line-1-offset], p[line+offset])
		}
		if diff == smudges {
			return line, true
		}
	}
	return 0, false
}

func verticalReflection(p pattern, smudges int) (int, bool) {
	width := len(p[0])
	for line := 1; line < width; line++ {
		span := line
		if width-line < span {
			span = width - line
		}
		diff := 0
		for offset := 0; offset < span; offset++ {
			diff += diffCount(p.column(line-1-offset), p.column(line+offset))
		}
		if diff == smudges {
			return line, true
		}
	}
	return 0, false
}

// summarize returns 100 times the rows above a horizontal mirror, or the
// columns left of a vertical mirror whose reflected cells differ by exactly
// smudges in total.
func summarize(p pattern, smudges int) (int, error) {
	if line, ok := horizontalReflection(p, smudges); ok {
		return 100 * line, nil
	}
	if line, ok := verticalReflection(p, smudges); ok {
		return line, nil
	}
	return 0, fmt.Errorf("%w: pattern has no mirror line", errInvalidNotes)
}

func solve(input string) (part1, part2 int, err error) {
	patterns, err := parsePatterns(input)
	if err != nil {
		return 0, 0, err
	}
	for _, p := range patterns {
		exact, err := summarize(p, 0)
		if err != nil {
			return 0, 0, err
		}
		smudged, err := summarize(p, 1)
		if err != nil {
			return 0, 0, err
		}
		part1 += exact
		part2 += smudged
	}
	return part1, part2, nil
}
