package main

import (
	"errors"
	"fmt"
	"strings"
)

const gridSize = 10

var errInvalidGrid = errors.New("invalid octopus grid")

type grid struct {
	energy [gridSize * gridSize]uint8
}

func parseGrid(input string) (grid, error) {
	var g grid
	rows := strings.Split(strings.TrimSpace(input), "\n")
	if len(rows) != gridSize {
		return g, fmt.Errorf("%w: grid needs %d rows", errInvalidGrid, gridSize)
	}
	for row, line := range rows {
		if len(line) != gridSize {
			return g, fmt.Errorf("%w: row needs %d columns", errInvalidGrid, gridSize)
		}
		for column, ch := range line {
			if ch < '0' || ch > '9' {
				return g, fmt.Errorf("%w: invalid energy level %q", errInvalidGrid, ch)
			}
			g.energy[row*gridSize+column] = uint8(ch - '0')
		}
	}
	return g, nil
}

func neighbors(index int) []int {
	row, column := index/gridSize, index%gridSize
	result := make([]int, 0, 8)
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			newRow, newColumn := row+dr, column+dc
			if newRow >= 0 && newRow < gridSize && newColumn >= 0 && newColumn < gridSize {
				result = append(result, newRow*gridSize+newColumn)
			}
		}
	}
	return result
}

// step advances the grid by one step and returns the number of octopuses
// that flashed.
func (g *grid) step() int {
	var flashed [gridSize * gridSize]bool
	stack := make([]int, 0)

	for i := range g.energy {
		g.energy[i]++
		if g.energy[i] > 9 {
			stack = append(stack, i)
		}
	}

	for len(stack) > 0 {
		index := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if flashed[index] {
			continue
		}
		flashed[index] = true
		for _, neighbor := range neighbors(index) {
			g.energy[neighbor]++
			if g.energy[neighbor] > 9 && !flashed[neighbor] {
				stack = append(stack, neighbor)
			}
		}
	}

	count := 0
	for i, f := range flashed {
		if f {
			g.energy[i] = 0
			count++
		}
	}
	return count
}

func solveWith(input string, steps int) (int, int, error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}
	totalCells := len(g.energy)
	totalFlashes := 0
	synchronizedStep := 0
	for step := 1; synchronizedStep == 0 || step <= steps; step++ {
		flashed := g.step()
		if step <= steps {
			totalFlashes += flashed
		}
		if flashed == totalCells && synchronizedStep == 0 {
			synchronizedStep = step
		}
	}
	return totalFlashes, synchronizedStep, nil
}

func solve(input string) (int, int, error) {
	return solveWith(input, 100)
}
