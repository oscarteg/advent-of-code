package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var errInvalidHeightMap = errors.New("invalid height map")

type heightMap struct {
	heights []uint8
	width   int
	height  int
}

func parseHeightMap(input string) (heightMap, error) {
	var result heightMap
	rows := strings.Split(strings.TrimSpace(input), "\n")
	if len(rows) == 0 || rows[0] == "" {
		return result, fmt.Errorf("%w: must have at least one row", errInvalidHeightMap)
	}
	width := len(rows[0])
	heights := make([]uint8, 0, width*len(rows))
	for _, row := range rows {
		if len(row) != width {
			return result, fmt.Errorf("%w: all rows must have the same width", errInvalidHeightMap)
		}
		for _, ch := range row {
			if ch < '0' || ch > '9' {
				return result, fmt.Errorf("%w: invalid height digit %q", errInvalidHeightMap, ch)
			}
			heights = append(heights, uint8(ch-'0'))
		}
	}
	return heightMap{heights: heights, width: width, height: len(rows)}, nil
}

type point struct{ x, y int }

func (m heightMap) at(p point) uint8 {
	return m.heights[p.y*m.width+p.x]
}

func (m heightMap) neighbors(p point) []point {
	candidates := []point{{p.x, p.y - 1}, {p.x, p.y + 1}, {p.x - 1, p.y}, {p.x + 1, p.y}}
	neighbors := make([]point, 0, 4)
	for _, c := range candidates {
		if c.x >= 0 && c.x < m.width && c.y >= 0 && c.y < m.height {
			neighbors = append(neighbors, c)
		}
	}
	return neighbors
}

func (m heightMap) isLowPoint(p point) bool {
	here := m.at(p)
	for _, n := range m.neighbors(p) {
		if m.at(n) <= here {
			return false
		}
	}
	return true
}

func (m heightMap) lowPoints() []point {
	var points []point
	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			p := point{x, y}
			if m.isLowPoint(p) {
				points = append(points, p)
			}
		}
	}
	return points
}

func (m heightMap) basinSize(start point) int {
	visited := make([]bool, len(m.heights))
	stack := []point{start}
	size := 0
	for len(stack) > 0 {
		p := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		index := p.y*m.width + p.x
		if visited[index] || m.at(p) == 9 {
			continue
		}
		visited[index] = true
		size++
		stack = append(stack, m.neighbors(p)...)
	}
	return size
}

func solve(input string) (uint32, int, error) {
	m, err := parseHeightMap(input)
	if err != nil {
		return 0, 0, err
	}

	lowPoints := m.lowPoints()

	var riskLevel uint32
	for _, p := range lowPoints {
		riskLevel += uint32(m.at(p)) + 1
	}

	basinSizes := make([]int, len(lowPoints))
	for i, p := range lowPoints {
		basinSizes[i] = m.basinSize(p)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(basinSizes)))
	if len(basinSizes) < 3 {
		return 0, 0, fmt.Errorf("%w: fewer than three basins found", errInvalidHeightMap)
	}
	largestBasins := basinSizes[0] * basinSizes[1] * basinSizes[2]

	return riskLevel, largestBasins, nil
}
