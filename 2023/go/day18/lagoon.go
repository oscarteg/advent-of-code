package main

import (
	"fmt"
	"strconv"
	"strings"
)

type direction int

const (
	right direction = iota
	down
	left
	up
)

func directionFromLetter(letter byte) (direction, error) {
	switch letter {
	case 'R':
		return right, nil
	case 'D':
		return down, nil
	case 'L':
		return left, nil
	case 'U':
		return up, nil
	default:
		return 0, fmt.Errorf("unknown direction letter %q", letter)
	}
}

func directionFromDigit(digit byte) (direction, error) {
	switch digit {
	case '0':
		return right, nil
	case '1':
		return down, nil
	case '2':
		return left, nil
	case '3':
		return up, nil
	default:
		return 0, fmt.Errorf("unknown direction digit %q", digit)
	}
}

func (d direction) offset() (int64, int64) {
	switch d {
	case up:
		return 0, -1
	case down:
		return 0, 1
	case left:
		return -1, 0
	default:
		return 1, 0
	}
}

type step struct {
	dir      direction
	distance int64
}

func parseStep(line string, useColour bool) (step, error) {
	fields := strings.Fields(line)
	if len(fields) != 3 {
		return step{}, fmt.Errorf("dig instruction needs three fields: %q", line)
	}

	if useColour {
		return parseColourStep(fields[2])
	}

	if len(fields[0]) != 1 {
		return step{}, fmt.Errorf("direction must be one letter: %q", fields[0])
	}
	dir, err := directionFromLetter(fields[0][0])
	if err != nil {
		return step{}, err
	}
	distance, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return step{}, fmt.Errorf("distance: %w", err)
	}
	return step{dir: dir, distance: distance}, nil
}

func parseColourStep(colour string) (step, error) {
	hex, ok := strings.CutPrefix(colour, "(#")
	hex, ok2 := strings.CutSuffix(hex, ")")
	if !ok || !ok2 || len(hex) != 6 {
		return step{}, fmt.Errorf("malformed colour: %q", colour)
	}
	distance, err := strconv.ParseInt(hex[:5], 16, 64)
	if err != nil {
		return step{}, fmt.Errorf("colour distance: %w", err)
	}
	dir, err := directionFromDigit(hex[5])
	if err != nil {
		return step{}, err
	}
	return step{dir: dir, distance: distance}, nil
}

func parseDigPlan(input string, useColour bool) ([]step, error) {
	var steps []step
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		s, err := parseStep(line, useColour)
		if err != nil {
			return nil, err
		}
		steps = append(steps, s)
	}
	if len(steps) == 0 {
		return nil, fmt.Errorf("empty dig plan")
	}
	return steps, nil
}

// lagoonArea computes the trench polygon's signed area with the shoelace
// formula, then applies Pick's theorem (area = interior + boundary/2 - 1)
// to recover the total number of dug-out tiles, interior + boundary.
func lagoonArea(steps []step) int64 {
	var x, y, shoelaceSum, boundary int64
	for _, s := range steps {
		dx, dy := s.dir.offset()
		nextX, nextY := x+dx*s.distance, y+dy*s.distance
		shoelaceSum += x*nextY - nextX*y
		boundary += s.distance
		x, y = nextX, nextY
	}
	area := shoelaceSum
	if area < 0 {
		area = -area
	}
	area /= 2
	return area + boundary/2 + 1
}

func solve(input string) (int64, int64, error) {
	plan1, err := parseDigPlan(input, false)
	if err != nil {
		return 0, 0, err
	}
	plan2, err := parseDigPlan(input, true)
	if err != nil {
		return 0, 0, err
	}
	return lagoonArea(plan1), lagoonArea(plan2), nil
}
