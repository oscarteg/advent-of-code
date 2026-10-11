package main

import (
	"fmt"
	"strings"
)

type direction int

const (
	up direction = iota
	down
	left
	right
)

func (d direction) dx() int {
	switch d {
	case up:
		return -1
	case down:
		return 1
	default:
		return 0
	}
}

func (d direction) dy() int {
	switch d {
	case left:
		return -1
	case right:
		return 1
	default:
		return 0
	}
}

func reflectForward(d direction) direction {
	switch d {
	case up:
		return right
	case down:
		return left
	case left:
		return down
	case right:
		return up
	}
	return d
}

func reflectBackslash(d direction) direction {
	switch d {
	case up:
		return left
	case down:
		return right
	case left:
		return up
	case right:
		return down
	}
	return d
}

func solve(input string) (int, int, error) {
	lines := strings.TrimSpace(input)
	if lines == "" {
		return 0, 0, nil
	}

	grid := make([][]rune, 0)
	for i, line := range strings.Split(lines, "\n") {
		if strings.ContainsFunc(line, func(r rune) bool { return !strings.ContainsRune(`./\|-`, r) }) {
			return 0, 0, fmt.Errorf(`row %d: only '.', '/', '\', '|' and '-' are allowed`, i+1)
		}
		if len(grid) > 0 && len(line) != len(grid[0]) {
			return 0, 0, fmt.Errorf("row %d has width %d, want %d", i+1, len(line), len(grid[0]))
		}
		grid = append(grid, []rune(line))
	}

	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0, 0, nil
	}

	part1 := len(simulate(grid, 0, 0, right))

	rows := len(grid)
	cols := len(grid[0])
	maxEnergized := 0

	for col := 0; col < cols; col++ {
		energized := simulate(grid, 0, col, down)
		if len(energized) > maxEnergized {
			maxEnergized = len(energized)
		}
	}

	for col := 0; col < cols; col++ {
		energized := simulate(grid, rows-1, col, up)
		if len(energized) > maxEnergized {
			maxEnergized = len(energized)
		}
	}

	for row := 0; row < rows; row++ {
		energized := simulate(grid, row, 0, right)
		if len(energized) > maxEnergized {
			maxEnergized = len(energized)
		}
	}

	for row := 0; row < rows; row++ {
		energized := simulate(grid, row, cols-1, left)
		if len(energized) > maxEnergized {
			maxEnergized = len(energized)
		}
	}

	return part1, maxEnergized, nil
}

func simulate(grid [][]rune, startRow, startCol int, startDir direction) map[[2]int]bool {
	visited := make(map[[2]int]bool)
	seen := make(map[[3]int]bool)
	queue := [][3]int{{startRow, startCol, int(startDir)}}

	for len(queue) > 0 {
		state := queue[0]
		queue = queue[1:]

		row, col, dirVal := state[0], state[1], state[2]
		d := direction(dirVal)

		if row < 0 || col < 0 || row >= len(grid) || col >= len(grid[0]) {
			continue
		}

		if seen[[3]int{row, col, dirVal}] {
			continue
		}
		seen[[3]int{row, col, dirVal}] = true

		visited[[2]int{row, col}] = true

		cell := grid[row][col]

		switch cell {
		case '.':
			nextRow := row + d.dx()
			nextCol := col + d.dy()
			queue = append(queue, [3]int{nextRow, nextCol, int(d)})

		case '/':
			newDir := reflectForward(d)
			nextRow := row + newDir.dx()
			nextCol := col + newDir.dy()
			queue = append(queue, [3]int{nextRow, nextCol, int(newDir)})

		case '\\':
			newDir := reflectBackslash(d)
			nextRow := row + newDir.dx()
			nextCol := col + newDir.dy()
			queue = append(queue, [3]int{nextRow, nextCol, int(newDir)})

		case '-':
			if d == left || d == right {
				nextRow := row + d.dx()
				nextCol := col + d.dy()
				queue = append(queue, [3]int{nextRow, nextCol, int(d)})
			} else {
				queue = append(queue, [3]int{row, col - 1, int(left)})
				queue = append(queue, [3]int{row, col + 1, int(right)})
			}

		case '|':
			if d == up || d == down {
				nextRow := row + d.dx()
				nextCol := col + d.dy()
				queue = append(queue, [3]int{nextRow, nextCol, int(d)})
			} else {
				queue = append(queue, [3]int{row - 1, col, int(up)})
				queue = append(queue, [3]int{row + 1, col, int(down)})
			}
		}
	}

	return visited
}
