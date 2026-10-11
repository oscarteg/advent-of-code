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

type blizzard struct {
	row, col int
	dir      direction
}

type position struct{ row, col int }

type grid struct {
	width, height    int
	startCol, endCol int
	blizzards        []blizzard
}

func parseGrid(input string) (grid, error) {
	lines := []string{}
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) < 3 {
		return grid{}, fmt.Errorf("parse grid: need a top wall, interior and bottom wall")
	}

	height := len(lines) - 2
	if len(lines[0]) < 2 {
		return grid{}, fmt.Errorf("parse grid: top wall too short")
	}
	width := len(lines[0]) - 2
	if width == 0 {
		return grid{}, fmt.Errorf("parse grid: interior must be non-empty")
	}

	top, bottom := lines[0], lines[len(lines)-1]
	if len(top) != width+2 || len(bottom) != width+2 {
		return grid{}, fmt.Errorf("parse grid: wall width mismatch")
	}

	startIdx := strings.IndexByte(top, '.')
	if startIdx <= 0 || startIdx-1 >= width {
		return grid{}, fmt.Errorf("parse grid: missing or invalid start opening")
	}
	endIdx := strings.IndexByte(bottom, '.')
	if endIdx <= 0 || endIdx-1 >= width {
		return grid{}, fmt.Errorf("parse grid: missing or invalid end opening")
	}

	var blizzards []blizzard
	for row, line := range lines[1 : 1+height] {
		if len(line) != width+2 {
			return grid{}, fmt.Errorf("parse grid: row %d has inconsistent width", row)
		}
		if line[0] != '#' || line[width+1] != '#' {
			return grid{}, fmt.Errorf("parse grid: row %d missing side walls", row)
		}
		for col, ch := range line[1 : width+1] {
			var dir direction
			switch ch {
			case '^':
				dir = up
			case 'v':
				dir = down
			case '<':
				dir = left
			case '>':
				dir = right
			case '.':
				continue
			default:
				return grid{}, fmt.Errorf("parse grid: unexpected character %q at row %d", ch, row)
			}
			blizzards = append(blizzards, blizzard{row: row, col: col, dir: dir})
		}
	}

	return grid{
		width:     width,
		height:    height,
		startCol:  startIdx - 1,
		endCol:    endIdx - 1,
		blizzards: blizzards,
	}, nil
}

func mod(a, b int) int {
	return ((a % b) + b) % b
}

func (g grid) occupiedAt(t int) map[position]bool {
	occupied := make(map[position]bool, len(g.blizzards))
	for _, b := range g.blizzards {
		var dr, dc int
		switch b.dir {
		case up:
			dr, dc = -1, 0
		case down:
			dr, dc = 1, 0
		case left:
			dr, dc = 0, -1
		case right:
			dr, dc = 0, 1
		}
		row := mod(b.row+dr*t, g.height)
		col := mod(b.col+dc*t, g.width)
		occupied[position{row, col}] = true
	}
	return occupied
}

func (g grid) isValid(row, col int) bool {
	if row == -1 {
		return col == g.startCol
	}
	if row == g.height {
		return col == g.endCol
	}
	return row >= 0 && col >= 0 && col < g.width
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

type state struct {
	row, col, timeMod int
}

func shortestTime(g grid, occupied []map[position]bool, period, startTime int, from, to position) (int, error) {
	type queued struct {
		pos  position
		time int
	}
	queue := []queued{{pos: from, time: startTime}}
	visited := map[state]bool{{from.row, from.col, startTime % period}: true}

	for head := 0; head < len(queue); head++ {
		current := queue[head]
		if current.pos == to {
			return current.time, nil
		}
		nextTime := current.time + 1
		occupiedNext := occupied[nextTime%period]
		deltas := [5][2]int{{0, 0}, {-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, d := range deltas {
			next := position{row: current.pos.row + d[0], col: current.pos.col + d[1]}
			if !g.isValid(next.row, next.col) {
				continue
			}
			if next.row >= 0 && next.row < g.height && occupiedNext[next] {
				continue
			}
			st := state{next.row, next.col, nextTime % period}
			if !visited[st] {
				visited[st] = true
				queue = append(queue, queued{pos: next, time: nextTime})
			}
		}
	}

	return 0, fmt.Errorf("shortest time: no path found")
}

// solve parses the grid and returns the elapsed time for the first trip
// (part1) and the total elapsed time for there-back-there (part2).
func solve(input string) (part1, part2 int, err error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}

	period := lcm(g.width, g.height)
	occupied := make([]map[position]bool, period)
	for t := 0; t < period; t++ {
		occupied[t] = g.occupiedAt(t)
	}

	start := position{row: -1, col: g.startCol}
	end := position{row: g.height, col: g.endCol}

	there, err := shortestTime(g, occupied, period, 0, start, end)
	if err != nil {
		return 0, 0, err
	}
	back, err := shortestTime(g, occupied, period, there, end, start)
	if err != nil {
		return 0, 0, err
	}
	thereAgain, err := shortestTime(g, occupied, period, back, start, end)
	if err != nil {
		return 0, 0, err
	}

	return there, thereAgain, nil
}
