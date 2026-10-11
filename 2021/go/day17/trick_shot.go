package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type rect struct {
	xMin, xMax, yMin, yMax int
}

func parseTarget(input string) (rect, error) {
	re := regexp.MustCompile(`target area: x=(-?\d+)\.\.(-?\d+), y=(-?\d+)\.\.(-?\d+)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(input))
	if len(matches) != 5 {
		return rect{}, fmt.Errorf("invalid target format")
	}

	xMin, _ := strconv.Atoi(matches[1])
	xMax, _ := strconv.Atoi(matches[2])
	yMin, _ := strconv.Atoi(matches[3])
	yMax, _ := strconv.Atoi(matches[4])

	return rect{xMin, xMax, yMin, yMax}, nil
}

func (r rect) contains(x, y int) bool {
	return x >= r.xMin && x <= r.xMax && y >= r.yMin && y <= r.yMax
}

func simulate(vx, vy int, target rect) (int, bool) {
	x, y := 0, 0
	maxY := 0

	for i := 0; i < 1000; i++ {
		x += vx
		y += vy
		if y > maxY {
			maxY = y
		}

		if target.contains(x, y) {
			return maxY, true
		}

		if x > target.xMax && (y < target.yMin || vy < 0) {
			return maxY, false
		}

		vx--
		if vx < 0 {
			vx = 0
		}
		vy--

		if y < target.yMin-1000 {
			return maxY, false
		}
	}

	return maxY, false
}

func solve(input string) (int, int, error) {
	target, err := parseTarget(input)
	if err != nil {
		return 0, 0, err
	}

	maxHeight := 0
	validCount := 0

	vxMin := 1
	vxMax := target.xMax + 1
	vyMin := target.yMin - 1
	vyMax := -target.yMin * 2

	for vx := vxMin; vx <= vxMax; vx++ {
		for vy := vyMin; vy <= vyMax; vy++ {
			h, hit := simulate(vx, vy, target)
			if hit {
				validCount++
				if h > maxHeight {
					maxHeight = h
				}
			}
		}
	}

	return maxHeight, validCount, nil
}
