package main

import (
	"fmt"
	"strings"
)

var errInvalidGrid = fmt.Errorf("invalid grid")

type grid struct {
	heights []int
	width   int
	height  int
	start   int
	end     int
}

func parseGrid(input string) (grid, error) {
	rows := strings.Split(strings.TrimRight(input, "\n"), "\n")
	if len(rows) == 0 || rows[0] == "" {
		return grid{}, fmt.Errorf("%w: grid must have at least one row", errInvalidGrid)
	}
	width := len(rows[0])

	heights := make([]int, 0, width*len(rows))
	start, end := -1, -1
	for _, row := range rows {
		if len(row) != width {
			return grid{}, fmt.Errorf("%w: all rows must have the same length", errInvalidGrid)
		}
		for _, ch := range row {
			index := len(heights)
			var elevation int
			switch {
			case ch == 'S':
				if start != -1 {
					return grid{}, fmt.Errorf("%w: duplicate start", errInvalidGrid)
				}
				start = index
				elevation = int('a' - 'a')
			case ch == 'E':
				if end != -1 {
					return grid{}, fmt.Errorf("%w: duplicate end", errInvalidGrid)
				}
				end = index
				elevation = int('z' - 'a')
			case ch >= 'a' && ch <= 'z':
				elevation = int(ch - 'a')
			default:
				return grid{}, fmt.Errorf("%w: invalid elevation character %q", errInvalidGrid, ch)
			}
			heights = append(heights, elevation)
		}
	}
	if start == -1 {
		return grid{}, fmt.Errorf("%w: missing start", errInvalidGrid)
	}
	if end == -1 {
		return grid{}, fmt.Errorf("%w: missing end", errInvalidGrid)
	}

	return grid{heights: heights, width: width, height: len(rows), start: start, end: end}, nil
}

func (g grid) neighbors(index int) []int {
	row, col := index/g.width, index%g.width
	neighbors := make([]int, 0, 4)
	if row > 0 {
		neighbors = append(neighbors, index-g.width)
	}
	if row+1 < g.height {
		neighbors = append(neighbors, index+g.width)
	}
	if col > 0 {
		neighbors = append(neighbors, index-1)
	}
	if col+1 < g.width {
		neighbors = append(neighbors, index+1)
	}
	return neighbors
}

// shortestPathFromEnd runs a breadth-first search from the end over the
// reversed graph: edge current -> neighbor here is the forward edge
// neighbor -> current, valid when it climbs by at most one elevation.
func (g grid) shortestPathFromEnd(starts []int) (int, error) {
	const unvisited = -1
	distances := make([]int, len(g.heights))
	for i := range distances {
		distances[i] = unvisited
	}
	distances[g.end] = 0
	queue := []int{g.end}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, neighbor := range g.neighbors(current) {
			if distances[neighbor] != unvisited {
				continue
			}
			if g.heights[current]-g.heights[neighbor] > 1 {
				continue
			}
			distances[neighbor] = distances[current] + 1
			queue = append(queue, neighbor)
		}
	}

	best := -1
	for _, start := range starts {
		if distances[start] == unvisited {
			continue
		}
		if best == -1 || distances[start] < best {
			best = distances[start]
		}
	}
	if best == -1 {
		return 0, fmt.Errorf("%w: no path to end", errInvalidGrid)
	}
	return best, nil
}

func solve(input string) (int, int, error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}

	part1, err := g.shortestPathFromEnd([]int{g.start})
	if err != nil {
		return 0, 0, err
	}

	var lowest []int
	for index, elevation := range g.heights {
		if elevation == 0 {
			lowest = append(lowest, index)
		}
	}
	part2, err := g.shortestPathFromEnd(lowest)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}
