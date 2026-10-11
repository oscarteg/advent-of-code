package day11

import "fmt"

type Cell int

const (
	Floor Cell = iota
	Empty
	Occupied
)

type Grid [][]Cell

func parseGrid(lines []string) (Grid, error) {
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty input")
	}
	grid := make(Grid, len(lines))
	for i, line := range lines {
		row := make([]Cell, len(line))
		for j, ch := range line {
			switch ch {
			case '.':
				row[j] = Floor
			case 'L':
				row[j] = Empty
			case '#':
				row[j] = Occupied
			default:
				return nil, fmt.Errorf("invalid cell: %c", ch)
			}
		}
		grid[i] = row
	}
	return grid, nil
}

func (g Grid) rows() int {
	return len(g)
}

func (g Grid) cols() int {
	if g.rows() == 0 {
		return 0
	}
	return len(g[0])
}

func (g Grid) get(r, c int) Cell {
	if r < 0 || r >= g.rows() || c < 0 || c >= g.cols() {
		return Floor
	}
	return g[r][c]
}

func countAdjacentOccupied(g Grid, r, c int) int {
	count := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			if g.get(r+dr, c+dc) == Occupied {
				count++
			}
		}
	}
	return count
}

func countVisibleOccupied(g Grid, r, c int) int {
	count := 0
	for dr := -1; dr <= 1; dr++ {
		for dc := -1; dc <= 1; dc++ {
			if dr == 0 && dc == 0 {
				continue
			}
			nr, nc := r+dr, c+dc
			for nr >= 0 && nr < g.rows() && nc >= 0 && nc < g.cols() {
				cell := g[nr][nc]
				if cell == Floor {
					nr += dr
					nc += dc
				} else if cell == Empty {
					break
				} else {
					count++
					break
				}
			}
		}
	}
	return count
}

func stepGrid(grid Grid, threshold int, countFn func(Grid, int, int) int) Grid {
	rows := grid.rows()
	cols := grid.cols()
	newGrid := make(Grid, rows)
	for r := 0; r < rows; r++ {
		newGrid[r] = make([]Cell, cols)
		for c := 0; c < cols; c++ {
			switch grid[r][c] {
			case Floor:
				newGrid[r][c] = Floor
			case Empty:
				if countFn(grid, r, c) == 0 {
					newGrid[r][c] = Occupied
				} else {
					newGrid[r][c] = Empty
				}
			case Occupied:
				if countFn(grid, r, c) >= threshold {
					newGrid[r][c] = Empty
				} else {
					newGrid[r][c] = Occupied
				}
			}
		}
	}
	return newGrid
}

func equal(g1, g2 Grid) bool {
	if len(g1) != len(g2) {
		return false
	}
	for i := range g1 {
		if len(g1[i]) != len(g2[i]) {
			return false
		}
		for j := range g1[i] {
			if g1[i][j] != g2[i][j] {
				return false
			}
		}
	}
	return true
}

func countOccupied(grid Grid) int {
	count := 0
	for _, row := range grid {
		for _, cell := range row {
			if cell == Occupied {
				count++
			}
		}
	}
	return count
}
