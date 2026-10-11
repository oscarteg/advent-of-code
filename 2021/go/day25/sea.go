package main

import (
	"fmt"
	"strings"
)

type cell byte

const (
	empty cell = '.'
	east  cell = '>'
	south cell = 'v'
)

func solve(input string) (int, error) {
	grid, err := parseGrid(input)
	if err != nil {
		return 0, err
	}

	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0, fmt.Errorf("grid must not be empty")
	}

	height := len(grid)
	width := len(grid[0])

	step := 0
	for {
		step++
		moved := simulateStep(grid, height, width)
		if !moved {
			return step, nil
		}
	}
}

func parseGrid(input string) ([][]cell, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	grid := make([][]cell, len(lines))
	for i, line := range lines {
		row := make([]cell, len(line))
		for j, ch := range line {
			c := cell(ch)
			if c != empty && c != east && c != south {
				return nil, fmt.Errorf("invalid character: %c", ch)
			}
			row[j] = c
		}
		grid[i] = row
	}
	return grid, nil
}

func simulateStep(grid [][]cell, height, width int) bool {
	moved := false

	moved = moveEast(grid, height, width) || moved
	moved = moveSouth(grid, height, width) || moved

	return moved
}

func moveEast(grid [][]cell, height, width int) bool {
	moved := false
	newGrid := make([][]cell, height)
	for i := range grid {
		newGrid[i] = make([]cell, width)
		copy(newGrid[i], grid[i])
	}

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			if grid[row][col] == east {
				nextCol := (col + 1) % width
				if grid[row][nextCol] == empty {
					newGrid[row][col] = empty
					newGrid[row][nextCol] = east
					moved = true
				}
			}
		}
	}

	copy2D(grid, newGrid)
	return moved
}

func moveSouth(grid [][]cell, height, width int) bool {
	moved := false
	newGrid := make([][]cell, height)
	for i := range grid {
		newGrid[i] = make([]cell, width)
		copy(newGrid[i], grid[i])
	}

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			if grid[row][col] == south {
				nextRow := (row + 1) % height
				if grid[nextRow][col] == empty {
					newGrid[row][col] = empty
					newGrid[nextRow][col] = south
					moved = true
				}
			}
		}
	}

	copy2D(grid, newGrid)
	return moved
}

func copy2D(dest, src [][]cell) {
	for i := range src {
		copy(dest[i], src[i])
	}
}
