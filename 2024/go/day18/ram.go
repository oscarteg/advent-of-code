package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidByte = errors.New("invalid ram run input")

const (
	exampleGridSize    = 7
	exampleFallenBytes = 12
)

type point struct {
	x, y int
}

func parsePoints(input string) ([]point, error) {
	var points []point
	for _, line := range strings.Split(strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n")), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		x, y, found := strings.Cut(line, ",")
		if !found {
			return nil, fmt.Errorf("%w: missing ',' in line %q", errInvalidByte, line)
		}
		px, err := strconv.Atoi(strings.TrimSpace(x))
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errInvalidByte, err)
		}
		py, err := strconv.Atoi(strings.TrimSpace(y))
		if err != nil {
			return nil, fmt.Errorf("%w: %w", errInvalidByte, err)
		}
		points = append(points, point{x: px, y: py})
	}
	if len(points) == 0 {
		return nil, fmt.Errorf("%w: no byte coordinates", errInvalidByte)
	}
	return points, nil
}

// shortestPath runs a BFS from (0,0) to (size-1,size-1) over a grid where every
// point in corrupted is a wall, returning -1 if the exit is unreachable.
func shortestPath(corrupted []point, size int) int {
	blocked := make(map[point]bool, len(corrupted))
	for _, p := range corrupted {
		if p.x >= 0 && p.x < size && p.y >= 0 && p.y < size {
			blocked[p] = true
		}
	}

	start := point{0, 0}
	target := point{size - 1, size - 1}
	distance := map[point]int{start: 0}
	queue := []point{start}
	deltas := []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == target {
			return distance[current]
		}
		for _, delta := range deltas {
			next := point{current.x + delta.x, current.y + delta.y}
			if next.x < 0 || next.x >= size || next.y < 0 || next.y >= size {
				continue
			}
			if blocked[next] {
				continue
			}
			if _, seen := distance[next]; seen {
				continue
			}
			distance[next] = distance[current] + 1
			queue = append(queue, next)
		}
	}
	return -1
}

// firstBlockingByte binary-searches the byte index: the exit stays reachable for
// every prefix up to the first blocker and unreachable for every prefix after it.
func firstBlockingByte(points []point, size int) (point, error) {
	low, high := 0, len(points)
	for low < high {
		mid := low + (high-low)/2
		if shortestPath(points[:mid+1], size) == -1 {
			high = mid
		} else {
			low = mid + 1
		}
	}
	if low >= len(points) {
		return point{}, fmt.Errorf("%w: no byte blocks the exit", errInvalidByte)
	}
	return points[low], nil
}

func solveWith(input string, size, fallenBytes int) (int, string, error) {
	points, err := parsePoints(input)
	if err != nil {
		return 0, "", err
	}
	if fallenBytes > len(points) {
		return 0, "", fmt.Errorf("%w: fewer byte coordinates than requested fallen bytes", errInvalidByte)
	}

	part1 := shortestPath(points[:fallenBytes], size)
	if part1 == -1 {
		return 0, "", fmt.Errorf("%w: exit unreachable after the given number of fallen bytes", errInvalidByte)
	}

	blocker, err := firstBlockingByte(points, size)
	if err != nil {
		return 0, "", err
	}
	part2 := fmt.Sprintf("%d,%d", blocker.x, blocker.y)
	return part1, part2, nil
}

func solve(input string) (int, string, error) {
	return solveWith(input, exampleGridSize, exampleFallenBytes)
}
