package main

import (
	"fmt"
	"strings"
)

type grid [][]byte

type direction struct {
	dx, dy int
}

var allDirections = []direction{
	{1, 0}, {-1, 0}, {0, 1}, {0, -1},
	{1, 1}, {1, -1}, {-1, 1}, {-1, -1},
}

func parseGrid(input string) (grid, error) {
	var rows grid
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if line == "" {
			continue
		}
		rows = append(rows, []byte(line))
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("word search: input has no grid rows")
	}
	width := len(rows[0])
	for _, row := range rows {
		if len(row) != width {
			return nil, fmt.Errorf("word search: grid rows have inconsistent width")
		}
	}
	return rows, nil
}

func (g grid) at(row, col int) (byte, bool) {
	if row < 0 || col < 0 || row >= len(g) || col >= len(g[0]) {
		return 0, false
	}
	return g[row][col], true
}

func (g grid) matchesWord(row, col int, dir direction, word string) bool {
	for i := 0; i < len(word); i++ {
		actual, ok := g.at(row+dir.dy*i, col+dir.dx*i)
		if !ok || actual != word[i] {
			return false
		}
	}
	return true
}

func (g grid) countXmas() int {
	count := 0
	for row := range g {
		for col := range g[row] {
			if g[row][col] != 'X' {
				continue
			}
			for _, dir := range allDirections {
				if g.matchesWord(row, col, dir, "XMAS") {
					count++
				}
			}
		}
	}
	return count
}

func (g grid) diagonalIsMas(row, col int, dir direction) bool {
	a, ok1 := g.at(row-dir.dy, col-dir.dx)
	b, ok2 := g.at(row+dir.dy, col+dir.dx)
	if !ok1 || !ok2 {
		return false
	}
	return (a == 'M' && b == 'S') || (a == 'S' && b == 'M')
}

func (g grid) countXMas() int {
	count := 0
	for row := range g {
		for col := range g[row] {
			if g[row][col] != 'A' {
				continue
			}
			if g.diagonalIsMas(row, col, direction{1, 1}) && g.diagonalIsMas(row, col, direction{1, -1}) {
				count++
			}
		}
	}
	return count
}

func solve(input string) (part1, part2 int, err error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}
	return g.countXmas(), g.countXMas(), nil
}
