package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidReservoir = errors.New("invalid reservoir input")

type point struct {
	x, y int
}

var source = point{x: 500, y: 0}

func parsePoint(text string) (point, error) {
	text = strings.TrimSpace(text)
	fields := strings.Split(text, ",")
	if len(fields) != 2 {
		return point{}, fmt.Errorf("%w: point %q is missing a comma", errInvalidReservoir, text)
	}
	x, err := strconv.Atoi(strings.TrimSpace(fields[0]))
	if err != nil {
		return point{}, fmt.Errorf("%w: invalid x in %q: %w", errInvalidReservoir, text, err)
	}
	y, err := strconv.Atoi(strings.TrimSpace(fields[1]))
	if err != nil {
		return point{}, fmt.Errorf("%w: invalid y in %q: %w", errInvalidReservoir, text, err)
	}
	return point{x: x, y: y}, nil
}

func addSegment(rocks map[point]bool, a, b point) error {
	switch {
	case a.x == b.x:
		lo, hi := min(a.y, b.y), max(a.y, b.y)
		for y := lo; y <= hi; y++ {
			rocks[point{x: a.x, y: y}] = true
		}
	case a.y == b.y:
		lo, hi := min(a.x, b.x), max(a.x, b.x)
		for x := lo; x <= hi; x++ {
			rocks[point{x: x, y: a.y}] = true
		}
	default:
		return fmt.Errorf("%w: segment from %v to %v is not axis-aligned", errInvalidReservoir, a, b)
	}
	return nil
}

func parseRocks(input string) (map[point]bool, error) {
	rocks := make(map[point]bool)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, "->")
		if len(fields) < 2 {
			return nil, fmt.Errorf("%w: a rock path needs at least two points", errInvalidReservoir)
		}
		previous, err := parsePoint(fields[0])
		if err != nil {
			return nil, err
		}
		for _, raw := range fields[1:] {
			current, err := parsePoint(raw)
			if err != nil {
				return nil, err
			}
			if err := addSegment(rocks, previous, current); err != nil {
				return nil, err
			}
			previous = current
		}
	}
	if len(rocks) == 0 {
		return nil, fmt.Errorf("%w: no rock paths found", errInvalidReservoir)
	}
	return rocks, nil
}

func isBlocked(settled map[point]bool, floorY *int, p point) bool {
	if floorY != nil && p.y == *floorY {
		return true
	}
	return settled[p]
}

// simulate drops sand from the source until a unit either falls past maxY
// (no floor) or settles on the source itself (with a floor), returning the
// count of settled units.
func simulate(rocks map[point]bool, maxY int, floorY *int) uint64 {
	settled := make(map[point]bool, len(rocks))
	for p := range rocks {
		settled[p] = true
	}

	var count uint64
	for !settled[source] {
		sand := source
		for {
			if floorY == nil && sand.y > maxY {
				return count
			}
			below := point{x: sand.x, y: sand.y + 1}
			belowLeft := point{x: sand.x - 1, y: sand.y + 1}
			belowRight := point{x: sand.x + 1, y: sand.y + 1}
			switch {
			case !isBlocked(settled, floorY, below):
				sand = below
			case !isBlocked(settled, floorY, belowLeft):
				sand = belowLeft
			case !isBlocked(settled, floorY, belowRight):
				sand = belowRight
			default:
				settled[sand] = true
				count++
				goto next
			}
		}
	next:
	}
	return count
}

func solve(input string) (uint64, uint64, error) {
	rocks, err := parseRocks(input)
	if err != nil {
		return 0, 0, err
	}
	maxY := 0
	for p := range rocks {
		if p.y > maxY {
			maxY = p.y
		}
	}
	floorY := maxY + 2
	part1 := simulate(rocks, maxY, nil)
	part2 := simulate(rocks, maxY, &floorY)
	return part1, part2, nil
}
