package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Cube struct {
	x, y, z int
}

func (c Cube) neighbors() [6]Cube {
	return [6]Cube{
		{c.x + 1, c.y, c.z},
		{c.x - 1, c.y, c.z},
		{c.x, c.y + 1, c.z},
		{c.x, c.y - 1, c.z},
		{c.x, c.y, c.z + 1},
		{c.x, c.y, c.z - 1},
	}
}

func parseInput(input string) ([]Cube, error) {
	var cubes []Cube
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid cube format: %s", line)
		}

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid x coordinate: %w", err)
		}
		y, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid y coordinate: %w", err)
		}
		z, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid z coordinate: %w", err)
		}

		cubes = append(cubes, Cube{x, y, z})
	}

	if len(cubes) == 0 {
		return nil, fmt.Errorf("no cubes in input")
	}

	return cubes, nil
}

func surfaceArea(cubes []Cube, set map[Cube]bool) int {
	area := 0
	for _, cube := range cubes {
		for _, neighbor := range cube.neighbors() {
			if !set[neighbor] {
				area++
			}
		}
	}
	return area
}

func externalSurfaceArea(cubes []Cube, set map[Cube]bool) int {
	minX, maxX := 1<<31-1, -(1 << 31)
	minY, maxY := 1<<31-1, -(1 << 31)
	minZ, maxZ := 1<<31-1, -(1 << 31)

	for _, c := range cubes {
		if c.x < minX {
			minX = c.x
		}
		if c.x > maxX {
			maxX = c.x
		}
		if c.y < minY {
			minY = c.y
		}
		if c.y > maxY {
			maxY = c.y
		}
		if c.z < minZ {
			minZ = c.z
		}
		if c.z > maxZ {
			maxZ = c.z
		}
	}

	minX--
	maxX++
	minY--
	maxY++
	minZ--
	maxZ++

	external := make(map[Cube]bool)
	queue := []Cube{{minX, minY, minZ}}
	external[Cube{minX, minY, minZ}] = true

	for len(queue) > 0 {
		cube := queue[0]
		queue = queue[1:]

		for _, neighbor := range cube.neighbors() {
			if neighbor.x >= minX && neighbor.x <= maxX &&
				neighbor.y >= minY && neighbor.y <= maxY &&
				neighbor.z >= minZ && neighbor.z <= maxZ &&
				!set[neighbor] && !external[neighbor] {
				external[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}

	area := 0
	for _, cube := range cubes {
		for _, neighbor := range cube.neighbors() {
			if external[neighbor] {
				area++
			}
		}
	}

	return area
}

func solve(input string) (int, int, error) {
	cubes, err := parseInput(input)
	if err != nil {
		return 0, 0, fmt.Errorf("parse: %w", err)
	}

	set := make(map[Cube]bool)
	for _, c := range cubes {
		set[c] = true
	}

	part1 := surfaceArea(cubes, set)
	part2 := externalSurfaceArea(cubes, set)

	return part1, part2, nil
}
