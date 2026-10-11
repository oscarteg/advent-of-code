package main

import (
	"fmt"
	"strings"
)

type Grid struct {
	trees [][]int
	rows  int
	cols  int
}

func parseGrid(input string) (*Grid, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty grid")
	}

	var trees [][]int
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var row []int
		for _, ch := range line {
			if ch >= '0' && ch <= '9' {
				row = append(row, int(ch-'0'))
			} else {
				return nil, fmt.Errorf("invalid character: %c", ch)
			}
		}

		if len(row) == 0 {
			continue
		}

		if len(trees) > 0 && len(row) != len(trees[0]) {
			return nil, fmt.Errorf("inconsistent row length")
		}

		trees = append(trees, row)
	}

	if len(trees) == 0 {
		return nil, fmt.Errorf("empty grid")
	}

	return &Grid{
		trees: trees,
		rows:  len(trees),
		cols:  len(trees[0]),
	}, nil
}

func (g *Grid) isVisible(r, c int) bool {
	height := g.trees[r][c]

	visibleLeft := true
	for col := 0; col < c; col++ {
		if g.trees[r][col] >= height {
			visibleLeft = false
			break
		}
	}
	if visibleLeft {
		return true
	}

	visibleRight := true
	for col := c + 1; col < g.cols; col++ {
		if g.trees[r][col] >= height {
			visibleRight = false
			break
		}
	}
	if visibleRight {
		return true
	}

	visibleUp := true
	for row := 0; row < r; row++ {
		if g.trees[row][c] >= height {
			visibleUp = false
			break
		}
	}
	if visibleUp {
		return true
	}

	visibleDown := true
	for row := r + 1; row < g.rows; row++ {
		if g.trees[row][c] >= height {
			visibleDown = false
			break
		}
	}
	if visibleDown {
		return true
	}

	return false
}

func (g *Grid) scenicScore(r, c int) int {
	height := g.trees[r][c]

	left := 0
	for col := c - 1; col >= 0; col-- {
		left++
		if g.trees[r][col] >= height {
			break
		}
	}

	right := 0
	for col := c + 1; col < g.cols; col++ {
		right++
		if g.trees[r][col] >= height {
			break
		}
	}

	up := 0
	for row := r - 1; row >= 0; row-- {
		up++
		if g.trees[row][c] >= height {
			break
		}
	}

	down := 0
	for row := r + 1; row < g.rows; row++ {
		down++
		if g.trees[row][c] >= height {
			break
		}
	}

	return left * right * up * down
}

func solve(input string) (int, int, error) {
	grid, err := parseGrid(input)
	if err != nil {
		return 0, 0, fmt.Errorf("parse error: %w", err)
	}

	visibleCount := 0
	for r := 0; r < grid.rows; r++ {
		for c := 0; c < grid.cols; c++ {
			if grid.isVisible(r, c) {
				visibleCount++
			}
		}
	}

	maxScore := 0
	for r := 0; r < grid.rows; r++ {
		for c := 0; c < grid.cols; c++ {
			score := grid.scenicScore(r, c)
			if score > maxScore {
				maxScore = score
			}
		}
	}

	return visibleCount, maxScore, nil
}
