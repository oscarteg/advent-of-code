package main

import (
	"fmt"
	"strconv"
	"strings"
)

type RotationDir int

const (
	Left RotationDir = iota
	Right
)

type Rotation struct {
	dir      RotationDir
	distance int
}

func parseRotation(line string) (Rotation, error) {
	if len(line) < 2 {
		return Rotation{}, fmt.Errorf("invalid format: line too short")
	}

	var dir RotationDir
	switch line[0] {
	case 'L':
		dir = Left
	case 'R':
		dir = Right
	default:
		return Rotation{}, fmt.Errorf("invalid direction: %c", line[0])
	}

	distance, err := strconv.Atoi(line[1:])
	if err != nil {
		return Rotation{}, fmt.Errorf("invalid distance: %w", err)
	}

	return Rotation{dir: dir, distance: distance}, nil
}

func (r Rotation) applyAndCountCrossings(start int) (int, int) {
	switch r.dir {
	case Right:
		newPos := (start + r.distance) % 100
		crossings := (start + r.distance) / 100
		return newPos, crossings

	case Left:
		// Modulo with negative numbers requires special handling in Go
		newPos := ((start-r.distance)%100 + 100) % 100
		var crossings int

		if start == 0 {
			crossings = r.distance / 100
		} else if r.distance < start {
			crossings = 0
		} else {
			crossings = 1 + (r.distance-start)/100
		}

		return newPos, crossings
	}

	return start, 0
}

func solve(input string) (uint32, uint32, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	position := 50
	var part1, part2 uint32

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		rotation, err := parseRotation(line)
		if err != nil {
			return 0, 0, err
		}

		newPos, crossings := rotation.applyAndCountCrossings(position)
		position = newPos
		part2 += uint32(crossings)

		if position == 0 {
			part1++
		}
	}

	return part1, part2, nil
}
