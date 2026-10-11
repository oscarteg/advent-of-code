package main

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidMap = errors.New("invalid guard map")

type direction int

const (
	up direction = iota
	right
	down
	left
)

func (d direction) turnRight() direction {
	return (d + 1) % 4
}

func (d direction) delta() (dr, dc int) {
	switch d {
	case up:
		return -1, 0
	case right:
		return 0, 1
	case down:
		return 1, 0
	default:
		return 0, -1
	}
}

func directionFromByte(b byte) (direction, bool) {
	switch b {
	case '^':
		return up, true
	case '>':
		return right, true
	case 'v':
		return down, true
	case '<':
		return left, true
	default:
		return 0, false
	}
}

type position struct {
	row, col int
}

type labMap struct {
	rows      []string
	rowCount  int
	colCount  int
	start     position
	startFace direction
}

func parseMap(input string) (labMap, error) {
	var rows []string
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		if line == "" {
			continue
		}
		rows = append(rows, line)
	}
	if len(rows) == 0 {
		return labMap{}, fmt.Errorf("%w: empty input", errInvalidMap)
	}

	cols := len(rows[0])
	for _, row := range rows {
		if len(row) != cols {
			return labMap{}, fmt.Errorf("%w: irregular grid", errInvalidMap)
		}
	}

	found := false
	var start position
	var startFace direction
	normalized := make([]string, len(rows))

	for r, row := range rows {
		cells := []byte(row)
		for c, ch := range cells {
			if ch == '.' || ch == '#' {
				continue
			}
			direction, ok := directionFromByte(ch)
			if !ok {
				return labMap{}, fmt.Errorf("%w: invalid character %q", errInvalidMap, ch)
			}
			if found {
				return labMap{}, fmt.Errorf("%w: multiple guards", errInvalidMap)
			}
			found = true
			start = position{row: r, col: c}
			startFace = direction
			cells[c] = '.'
		}
		normalized[r] = string(cells)
	}

	if !found {
		return labMap{}, fmt.Errorf("%w: missing guard", errInvalidMap)
	}

	return labMap{
		rows:      normalized,
		rowCount:  len(normalized),
		colCount:  cols,
		start:     start,
		startFace: startFace,
	}, nil
}

func (m labMap) inBounds(r, c int) bool {
	return r >= 0 && r < m.rowCount && c >= 0 && c < m.colCount
}

func (m labMap) isWall(r, c int) bool {
	return m.rows[r][c] == '#'
}

// walk simulates the patrol, recording every visited position in visited
// (when non-nil) and treating extraWall as an additional obstruction. It
// reports whether the guard loops forever instead of leaving the map.
func (m labMap) walk(extraWall *position, visited map[position]bool) bool {
	type state struct {
		pos  position
		face direction
	}
	seen := make(map[state]bool)

	pos := m.start
	face := m.startFace

	for {
		s := state{pos, face}
		if seen[s] {
			return true
		}
		seen[s] = true
		if visited != nil {
			visited[pos] = true
		}

		dr, dc := face.delta()
		next := position{row: pos.row + dr, col: pos.col + dc}
		if !m.inBounds(next.row, next.col) {
			return false
		}

		blocked := m.isWall(next.row, next.col) || (extraWall != nil && *extraWall == next)
		if blocked {
			face = face.turnRight()
		} else {
			pos = next
		}
	}
}

func solve(input string) (part1, part2 int, err error) {
	m, err := parseMap(input)
	if err != nil {
		return 0, 0, err
	}

	visited := make(map[position]bool)
	m.walk(nil, visited)
	part1 = len(visited)

	for pos := range visited {
		if pos == m.start {
			continue
		}
		candidate := pos
		if m.walk(&candidate, nil) {
			part2++
		}
	}

	return part1, part2, nil
}
