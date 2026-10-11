package main

import (
	"errors"
	"fmt"
	"strings"
)

type point struct{ row, col int }

type topoMap struct {
	heights    []string
	rows, cols int
}

var directions = [4]point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

func parseMap(input string) (topoMap, error) {
	lines := strings.Fields(input)
	if len(lines) == 0 {
		return topoMap{}, errors.New("empty input")
	}
	cols := len(lines[0])
	for _, line := range lines {
		if len(line) != cols {
			return topoMap{}, errors.New("irregular grid")
		}
		for _, c := range line {
			if c < '0' || c > '9' {
				return topoMap{}, fmt.Errorf("invalid character in grid: %c", c)
			}
		}
	}
	return topoMap{heights: lines, rows: len(lines), cols: cols}, nil
}

func (m topoMap) inside(p point) bool {
	return p.row >= 0 && p.row < m.rows && p.col >= 0 && p.col < m.cols
}

// climb follows every strictly uphill path from p, records each summit (height 9)
// it reaches and returns how many distinct trails end there.
func (m topoMap) climb(p point, summits map[point]bool) int {
	height := m.heights[p.row][p.col]
	if height == '9' {
		summits[p] = true
		return 1
	}
	trails := 0
	for _, d := range directions {
		next := point{p.row + d.row, p.col + d.col}
		if m.inside(next) && m.heights[next.row][next.col] == height+1 {
			trails += m.climb(next, summits)
		}
	}
	return trails
}

func solve(input string) (int, int, error) {
	m, err := parseMap(input)
	if err != nil {
		return 0, 0, err
	}
	score, rating := 0, 0
	for row := range m.rows {
		for col := range m.cols {
			if m.heights[row][col] != '0' {
				continue
			}
			summits := map[point]bool{}
			rating += m.climb(point{row, col}, summits)
			score += len(summits)
		}
	}
	return score, rating, nil
}
