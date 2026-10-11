package main

import (
	"fmt"
	"strings"
)

func parseGrid(input string) ([][]rune, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, fmt.Errorf("empty grid")
	}

	lines := strings.Split(trimmed, "\n")
	grid := make([][]rune, len(lines))
	for i, line := range lines {
		if len(line) != len(lines[0]) {
			return nil, fmt.Errorf("row %d has width %d, want %d", i+1, len(line), len(lines[0]))
		}
		if strings.ContainsFunc(line, func(r rune) bool { return !strings.ContainsRune("O#.", r) }) {
			return nil, fmt.Errorf("row %d: only 'O', '#' and '.' are allowed", i+1)
		}
		grid[i] = []rune(line)
	}
	return grid, nil
}

func tiltNorth(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for col := 0; col < cols; col++ {
		writePos := 0
		for row := 0; row < rows; row++ {
			switch grid[row][col] {
			case 'O':
				grid[row][col] = '.'
				grid[writePos][col] = 'O'
				writePos++
			case '#':
				writePos = row + 1
			}
		}
	}
}

func tiltWest(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for row := 0; row < rows; row++ {
		writePos := 0
		for col := 0; col < cols; col++ {
			switch grid[row][col] {
			case 'O':
				grid[row][col] = '.'
				grid[row][writePos] = 'O'
				writePos++
			case '#':
				writePos = col + 1
			}
		}
	}
}

func tiltSouth(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for col := 0; col < cols; col++ {
		writePos := rows
		for row := rows - 1; row >= 0; row-- {
			switch grid[row][col] {
			case 'O':
				grid[row][col] = '.'
				writePos--
				grid[writePos][col] = 'O'
			case '#':
				writePos = row
			}
		}
	}
}

func tiltEast(grid [][]rune) {
	rows := len(grid)
	cols := len(grid[0])

	for row := 0; row < rows; row++ {
		writePos := cols
		for col := cols - 1; col >= 0; col-- {
			switch grid[row][col] {
			case 'O':
				grid[row][col] = '.'
				writePos--
				grid[row][writePos] = 'O'
			case '#':
				writePos = col
			}
		}
	}
}

func spinCycle(grid [][]rune) {
	tiltNorth(grid)
	tiltWest(grid)
	tiltSouth(grid)
	tiltEast(grid)
}

func gridToString(grid [][]rune) string {
	lines := make([]string, len(grid))
	for i, row := range grid {
		lines[i] = string(row)
	}
	return strings.Join(lines, "\n")
}

func calculateLoad(grid [][]rune) int {
	rows := len(grid)
	load := 0

	for row := 0; row < rows; row++ {
		rowValue := rows - row
		for _, ch := range grid[row] {
			if ch == 'O' {
				load += rowValue
			}
		}
	}
	return load
}

func solve(input string) (int, int, error) {
	grid, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}

	gridCopy := make([][]rune, len(grid))
	for i := range grid {
		gridCopy[i] = make([]rune, len(grid[i]))
		copy(gridCopy[i], grid[i])
	}

	tiltNorth(gridCopy)
	part1 := calculateLoad(gridCopy)

	targetCycles := 1_000_000_000
	seen := make(map[string]int)

	for cycleNum := 0; cycleNum < targetCycles; cycleNum++ {
		gridStr := gridToString(grid)
		if prevCycle, exists := seen[gridStr]; exists {
			cycleLength := cycleNum - prevCycle
			remaining := (targetCycles - cycleNum) % cycleLength
			for i := 0; i < remaining; i++ {
				spinCycle(grid)
			}
			part2 := calculateLoad(grid)
			return part1, part2, nil
		}
		seen[gridStr] = cycleNum
		spinCycle(grid)
	}

	part2 := calculateLoad(grid)
	return part1, part2, nil
}
