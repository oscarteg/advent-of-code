package main

import (
	"fmt"
	"strings"
)

func solve(input string) (int, int, error) {
	galaxies, err := parseGalaxies(input)
	if err != nil {
		return 0, 0, err
	}

	part1 := solveWithExpansion(galaxies, 2)
	part2 := solveWithExpansion(galaxies, 1_000_000)

	return part1, part2, nil
}

func parseGalaxies(input string) ([]Pos, error) {
	var galaxies []Pos
	for row, line := range strings.Split(strings.TrimSpace(input), "\n") {
		for col, ch := range line {
			if ch == '#' {
				galaxies = append(galaxies, Pos{row, col})
			}
		}
	}
	if len(galaxies) == 0 {
		return nil, fmt.Errorf("no galaxies found")
	}
	return galaxies, nil
}

type Pos struct {
	row, col int
}

func solveWithExpansion(galaxies []Pos, factor int) int {
	// Find grid dimensions
	maxRow := 0
	maxCol := 0
	for _, g := range galaxies {
		if g.row > maxRow {
			maxRow = g.row
		}
		if g.col > maxCol {
			maxCol = g.col
		}
	}

	// Mark empty rows and columns
	emptyRows := make([]bool, maxRow+1)
	emptyCols := make([]bool, maxCol+1)
	for i := range emptyRows {
		emptyRows[i] = true
	}
	for i := range emptyCols {
		emptyCols[i] = true
	}

	for _, g := range galaxies {
		emptyRows[g.row] = false
		emptyCols[g.col] = false
	}

	// Calculate sum of distances with expansion
	totalDistance := 0
	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			r1, c1 := galaxies[i].row, galaxies[i].col
			r2, c2 := galaxies[j].row, galaxies[j].col

			rowStart := r1
			if r2 < r1 {
				rowStart = r2
			}
			rowEnd := r1
			if r1 < r2 {
				rowEnd = r2
			}

			colStart := c1
			if c2 < c1 {
				colStart = c2
			}
			colEnd := c1
			if c1 < c2 {
				colEnd = c2
			}

			// Count empty rows and columns crossed
			emptyRowsCrossed := 0
			for r := rowStart; r < rowEnd; r++ {
				if emptyRows[r] {
					emptyRowsCrossed++
				}
			}
			emptyColsCrossed := 0
			for c := colStart; c < colEnd; c++ {
				if emptyCols[c] {
					emptyColsCrossed++
				}
			}

			// Calculate distance with expansion
			baseDistance := (rowEnd - rowStart) + (colEnd - colStart)
			extraDistance := (emptyRowsCrossed + emptyColsCrossed) * (factor - 1)
			distance := baseDistance + extraDistance

			totalDistance += distance
		}
	}

	return totalDistance
}
