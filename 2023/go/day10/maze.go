package main

import (
	"fmt"
	"strings"
)

var errInvalidMaze = fmt.Errorf("invalid pipe maze")

type direction struct {
	dr, dc int
}

var (
	north = direction{-1, 0}
	south = direction{1, 0}
	east  = direction{0, 1}
	west  = direction{0, -1}
)

func (d direction) opposite() direction {
	return direction{-d.dr, -d.dc}
}

type point struct {
	row, col int
}

func pipeDirections(tile byte) ([2]direction, bool) {
	switch tile {
	case '|':
		return [2]direction{north, south}, true
	case '-':
		return [2]direction{east, west}, true
	case 'L':
		return [2]direction{north, east}, true
	case 'J':
		return [2]direction{north, west}, true
	case '7':
		return [2]direction{south, west}, true
	case 'F':
		return [2]direction{south, east}, true
	default:
		return [2]direction{}, false
	}
}

func connects(tile byte, dir direction) bool {
	dirs, ok := pipeDirections(tile)
	return ok && (dirs[0] == dir || dirs[1] == dir)
}

type grid struct {
	tiles []string
	rows  int
	cols  int
}

func parseGrid(input string) (grid, error) {
	lines := strings.Split(strings.TrimRight(strings.ReplaceAll(input, "\r\n", "\n"), "\n"), "\n")
	if len(lines) == 0 || lines[0] == "" {
		return grid{}, fmt.Errorf("%w: grid must not be empty", errInvalidMaze)
	}
	width := len(lines[0])
	for _, line := range lines {
		if len(line) != width {
			return grid{}, fmt.Errorf("%w: every row must have the same width", errInvalidMaze)
		}
		for _, c := range line {
			if _, ok := pipeDirections(byte(c)); !ok && c != '.' && c != 'S' {
				return grid{}, fmt.Errorf("%w: invalid tile character %q", errInvalidMaze, c)
			}
		}
	}
	return grid{tiles: lines, rows: len(lines), cols: width}, nil
}

func (g grid) at(p point) byte {
	return g.tiles[p.row][p.col]
}

func (g grid) inBounds(p point) bool {
	return p.row >= 0 && p.row < g.rows && p.col >= 0 && p.col < g.cols
}

func (g grid) step(p point, dir direction) (point, bool) {
	next := point{p.row + dir.dr, p.col + dir.dc}
	return next, g.inBounds(next)
}

func (g grid) start() (point, error) {
	for r, line := range g.tiles {
		if c := strings.IndexByte(line, 'S'); c != -1 {
			return point{r, c}, nil
		}
	}
	return point{}, fmt.Errorf("%w: grid has no start tile", errInvalidMaze)
}

func (g grid) startDirections(start point) []direction {
	var dirs []direction
	for _, dir := range []direction{north, south, east, west} {
		if next, ok := g.step(start, dir); ok && connects(g.at(next), dir.opposite()) {
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

// findLoop walks the main loop starting at S, inferring its real pipe shape
// from which neighbors connect back to it, and returns every tile visited in
// order.
func (g grid) findLoop() ([]point, error) {
	start, err := g.start()
	if err != nil {
		return nil, err
	}
	startDirs := g.startDirections(start)
	if len(startDirs) != 2 {
		return nil, fmt.Errorf("%w: start tile must connect to exactly two pipes", errInvalidMaze)
	}

	pos := start
	dir := startDirs[0]
	var loop []point
	for {
		loop = append(loop, pos)
		next, ok := g.step(pos, dir)
		if !ok {
			return nil, fmt.Errorf("%w: loop runs off the edge of the grid", errInvalidMaze)
		}
		pos = next
		if pos == start {
			return loop, nil
		}
		dirs, ok := pipeDirections(g.at(pos))
		if !ok {
			return nil, fmt.Errorf("%w: loop passes through a non-pipe tile", errInvalidMaze)
		}
		incoming := dir.opposite()
		switch {
		case dirs[0] == incoming:
			dir = dirs[1]
		case dirs[1] == incoming:
			dir = dirs[0]
		default:
			return nil, fmt.Errorf("%w: pipe does not connect back to the previous tile", errInvalidMaze)
		}
	}
}

// enclosedArea counts tiles enclosed by a simple loop using the shoelace
// formula for its area and Pick's theorem (area = interior + boundary/2 - 1)
// to recover the interior point count.
func enclosedArea(loop []point) int {
	n := len(loop)
	area2 := 0
	for i := 0; i < n; i++ {
		a, b := loop[i], loop[(i+1)%n]
		area2 += a.col*b.row - b.col*a.row
	}
	if area2 < 0 {
		area2 = -area2
	}
	return area2/2 - n/2 + 1
}

func solve(input string) (int, int, error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}
	loop, err := g.findLoop()
	if err != nil {
		return 0, 0, err
	}
	return len(loop) / 2, enclosedArea(loop), nil
}
