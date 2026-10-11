package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidLine = errors.New("invalid line format")

type point struct {
	x, y int
}

type line struct {
	start, end point
}

func parseLine(s string) (line, error) {
	parts := strings.Split(s, "->")
	if len(parts) != 2 {
		return line{}, fmt.Errorf("%w: expected 'x1,y1 -> x2,y2'", errInvalidLine)
	}

	start, err := parsePoint(strings.TrimSpace(parts[0]))
	if err != nil {
		return line{}, fmt.Errorf("start point: %w", err)
	}

	end, err := parsePoint(strings.TrimSpace(parts[1]))
	if err != nil {
		return line{}, fmt.Errorf("end point: %w", err)
	}

	return line{start, end}, nil
}

func parsePoint(s string) (point, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return point{}, fmt.Errorf("%w: expected 'x,y'", errInvalidLine)
	}

	x, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return point{}, fmt.Errorf("x coordinate: %w", err)
	}

	y, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return point{}, fmt.Errorf("y coordinate: %w", err)
	}

	return point{x, y}, nil
}

func (l line) isHorizontalOrVertical() bool {
	return l.start.x == l.end.x || l.start.y == l.end.y
}

func (l line) points() []point {
	var result []point
	dx := 0
	dy := 0

	if l.end.x > l.start.x {
		dx = 1
	} else if l.end.x < l.start.x {
		dx = -1
	}

	if l.end.y > l.start.y {
		dy = 1
	} else if l.end.y < l.start.y {
		dy = -1
	}

	current := l.start
	for {
		result = append(result, current)
		if current == l.end {
			break
		}
		current.x += dx
		current.y += dy
	}

	return result
}

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 {
		return 0, 0, errors.New("empty input")
	}

	var vents []line
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		v, err := parseLine(l)
		if err != nil {
			return 0, 0, err
		}
		vents = append(vents, v)
	}

	if len(vents) == 0 {
		return 0, 0, errors.New("no valid lines found")
	}

	part1 := countOverlaps(filterHorizontalVertical(vents))
	part2 := countOverlaps(vents)

	return part1, part2, nil
}

func filterHorizontalVertical(vents []line) []line {
	var result []line
	for _, v := range vents {
		if v.isHorizontalOrVertical() {
			result = append(result, v)
		}
	}
	return result
}

func countOverlaps(vents []line) int {
	counts := make(map[point]int)
	for _, v := range vents {
		for _, p := range v.points() {
			counts[p]++
		}
	}

	overlap := 0
	for _, count := range counts {
		if count >= 2 {
			overlap++
		}
	}
	return overlap
}
