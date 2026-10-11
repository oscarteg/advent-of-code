package main

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidGrid = errors.New("invalid paper roll grid")

// grid is a mutable rectangular grid of paper rolls ('@') and empty cells ('.').
type grid struct {
	cells  [][]byte
	height int
	width  int
}

func parseGrid(input string) (grid, error) {
	var result grid
	for _, line := range strings.Split(strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n")), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if result.width != 0 && len(line) != result.width {
			return result, fmt.Errorf("%w: inconsistent row width", errInvalidGrid)
		}
		row := make([]byte, len(line))
		for i := 0; i < len(line); i++ {
			if line[i] != '@' && line[i] != '.' {
				return result, fmt.Errorf("%w: invalid character %q", errInvalidGrid, line[i])
			}
			row[i] = line[i]
		}
		result.width = len(row)
		result.cells = append(result.cells, row)
	}
	result.height = len(result.cells)
	if result.height == 0 {
		return result, fmt.Errorf("%w: empty grid", errInvalidGrid)
	}
	return result, nil
}

func (g grid) clone() grid {
	cells := make([][]byte, g.height)
	for i, row := range g.cells {
		cells[i] = append([]byte(nil), row...)
	}
	return grid{cells: cells, height: g.height, width: g.width}
}

func (g grid) countNeighbours(row, col int) int {
	count := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			r, c := row+dr, col+dc
			if r < 0 || c < 0 || r >= g.height || c >= g.width {
				continue
			}
			if g.cells[r][c] == '@' {
				count++
			}
		}
	}
	return count
}

func (g grid) isAccessible(row, col int) bool {
	return g.cells[row][col] == '@' && g.countNeighbours(row, col) < 4
}

func solve(input string) (int, int, error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}

	part1 := 0
	for row := 0; row < g.height; row++ {
		for col := 0; col < g.width; col++ {
			if g.isAccessible(row, col) {
				part1++
			}
		}
	}

	working := g.clone()
	part2 := 0
	for {
		type cell struct{ row, col int }
		var removable []cell
		for row := 0; row < working.height; row++ {
			for col := 0; col < working.width; col++ {
				if working.isAccessible(row, col) {
					removable = append(removable, cell{row, col})
				}
			}
		}
		if len(removable) == 0 {
			break
		}
		for _, c := range removable {
			working.cells[c.row][c.col] = '.'
		}
		part2 += len(removable)
	}

	return part1, part2, nil
}
