package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

func solve(input string) (int64, int64, error) {
	fields := strings.Split(strings.TrimSpace(input), ",")
	if len(fields) == 0 || fields[0] == "" {
		return 0, 0, fmt.Errorf("no positions provided")
	}

	positions := make([]int64, 0, len(fields))
	for _, field := range fields {
		position, err := strconv.ParseInt(strings.TrimSpace(field), 10, 64)
		if err != nil || position < 0 {
			return 0, 0, fmt.Errorf("invalid position %q", field)
		}
		positions = append(positions, position)
	}

	sort.Slice(positions, func(i, j int) bool {
		return positions[i] < positions[j]
	})

	median := positions[len(positions)/2]
	var linearFuel int64
	for _, position := range positions {
		linearFuel += abs(position - median)
	}

	var sum int64
	for _, position := range positions {
		sum += position
	}
	mean := float64(sum) / float64(len(positions))
	lower := int64(math.Floor(mean))
	upper := int64(math.Ceil(mean))
	increasingFuel := min(triangularFuel(positions, lower), triangularFuel(positions, upper))
	return linearFuel, increasingFuel, nil
}

func triangularFuel(positions []int64, target int64) int64 {
	var fuel int64
	for _, position := range positions {
		distance := abs(position - target)
		fuel += distance * (distance + 1) / 2
	}
	return fuel
}

func abs(value int64) int64 {
	if value < 0 {
		return -value
	}
	return value
}
