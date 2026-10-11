// Package main solves Advent of Code 2025 day 7, "Laboratories".
package main

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidGrid = errors.New("invalid grid")

// grid is the parsed laboratory layout: rows of equal width containing '.'
// and '^', plus the column the beam starts from in the first row it occupies.
type grid struct {
	rows     []string
	width    int
	startRow int
	startCol int
}

func parseGrid(input string) (grid, error) {
	var rows []string
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		rows = append(rows, line)
	}
	if len(rows) == 0 {
		return grid{}, fmt.Errorf("%w: empty input", errInvalidGrid)
	}

	width := len(rows[0])
	startRow, startCol := -1, -1
	for r, row := range rows {
		if len(row) != width {
			return grid{}, fmt.Errorf("%w: row %d has inconsistent width", errInvalidGrid, r)
		}
		for c, ch := range row {
			switch ch {
			case '.', '^':
			case 'S':
				if startRow != -1 {
					return grid{}, fmt.Errorf("%w: more than one start", errInvalidGrid)
				}
				startRow, startCol = r, c
			default:
				return grid{}, fmt.Errorf("%w: unexpected character %q", errInvalidGrid, ch)
			}
		}
	}
	if startRow == -1 {
		return grid{}, fmt.Errorf("%w: missing start", errInvalidGrid)
	}

	return grid{rows: rows, width: width, startRow: startRow, startCol: startCol}, nil
}

// solve simulates the beam: part1 counts how many times a splitter fires,
// part2 counts how many timelines (beams, with multiplicity) reach the
// bottom row.
func solve(input string) (part1, part2 uint64, err error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}

	counts := map[int]uint64{g.startCol: 1}
	for row := g.startRow + 1; row < len(g.rows); row++ {
		next := make(map[int]uint64, len(counts))
		for col, count := range counts {
			if g.rows[row][col] == '^' {
				part1++
				if col > 0 {
					next[col-1] += count
				}
				if col+1 < g.width {
					next[col+1] += count
				}
			} else {
				next[col] += count
			}
		}
		counts = next
	}

	for _, count := range counts {
		part2 += count
	}
	return part1, part2, nil
}
