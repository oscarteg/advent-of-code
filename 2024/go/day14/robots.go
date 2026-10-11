package main

import (
	"fmt"
	"regexp"
	"strconv"
)

type robot struct {
	px, py, vx, vy int64
}

type position struct {
	x, y int64
}

var robotPattern = regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

func parseRobots(input string) ([]robot, error) {
	matches := robotPattern.FindAllStringSubmatch(input, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("parse: no valid robots found")
	}

	robots := make([]robot, 0, len(matches))
	for _, m := range matches {
		fields := make([]int64, 4)
		for i, text := range m[1:] {
			value, err := strconv.ParseInt(text, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("parse: invalid field %q: %w", text, err)
			}
			fields[i] = value
		}
		robots = append(robots, robot{px: fields[0], py: fields[1], vx: fields[2], vy: fields[3]})
	}

	return robots, nil
}

func wrap(value, modulus int64) int64 {
	m := value % modulus
	if m < 0 {
		m += modulus
	}
	return m
}

func positionAt(r robot, seconds, width, height int64) position {
	return position{
		x: wrap(r.px+r.vx*seconds, width),
		y: wrap(r.py+r.vy*seconds, height),
	}
}

func safetyFactor(robots []robot, width, height, seconds int64) int64 {
	midX, midY := width/2, height/2
	var quadrants [4]int64

	for _, r := range robots {
		pos := positionAt(r, seconds, width, height)
		if pos.x == midX || pos.y == midY {
			continue
		}
		index := 0
		if pos.x > midX {
			index += 1
		}
		if pos.y > midY {
			index += 2
		}
		quadrants[index]++
	}

	return quadrants[0] * quadrants[1] * quadrants[2] * quadrants[3]
}

func positionsUnique(positions []position, width int64) bool {
	seen := make(map[int64]struct{}, len(positions))
	for _, pos := range positions {
		key := pos.y*width + pos.x
		if _, ok := seen[key]; ok {
			return false
		}
		seen[key] = struct{}{}
	}
	return true
}

// The robots form a recognizable picture when no two of them overlap.
// Positions repeat with period width*height, which bounds the search.
func findPictureSecond(robots []robot, width, height int64) (int64, error) {
	period := width * height
	positions := make([]position, len(robots))

	for t := int64(0); t < period; t++ {
		for i, r := range robots {
			positions[i] = positionAt(r, t, width, height)
		}
		if positionsUnique(positions, width) {
			return t, nil
		}
	}

	return 0, fmt.Errorf("no non-overlapping arrangement found within one period")
}

// The input never states the grid size: the puzzle's example uses 11x7 and
// every real input uses 101x103, so pick between those two.
func inferDimensions(robots []robot) (width, height int64) {
	for _, r := range robots {
		if r.px >= 11 || r.py >= 7 {
			return 101, 103
		}
	}
	return 11, 7
}

func solveWith(input string, width, height int64) (int64, int64, error) {
	robots, err := parseRobots(input)
	if err != nil {
		return 0, 0, err
	}

	part1 := safetyFactor(robots, width, height, 100)
	part2, err := findPictureSecond(robots, width, height)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}

func solve(input string) (int64, int64, error) {
	robots, err := parseRobots(input)
	if err != nil {
		return 0, 0, err
	}

	width, height := inferDimensions(robots)
	part1 := safetyFactor(robots, width, height, 100)
	part2, err := findPictureSecond(robots, width, height)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}
