package main

import (
	"container/heap"
	"errors"
	"fmt"
	"strings"
)

var errInvalidGrid = errors.New("invalid chiton grid")

type grid struct {
	risks  []uint8
	width  int
	height int
}

func parseGrid(input string) (grid, error) {
	rows := strings.Split(strings.TrimSpace(strings.ReplaceAll(input, "\r\n", "\n")), "\n")
	if len(rows) == 0 || rows[0] == "" {
		return grid{}, fmt.Errorf("%w: missing rows", errInvalidGrid)
	}
	width := len(rows[0])
	risks := make([]uint8, 0, width*len(rows))
	for _, row := range rows {
		if len(row) != width {
			return grid{}, fmt.Errorf("%w: rows must have equal length", errInvalidGrid)
		}
		for _, ch := range row {
			if ch < '0' || ch > '9' {
				return grid{}, fmt.Errorf("%w: invalid risk digit %q", errInvalidGrid, ch)
			}
			risks = append(risks, uint8(ch-'0'))
		}
	}
	return grid{risks: risks, width: width, height: len(rows)}, nil
}

func (g grid) riskAt(x, y int) uint32 {
	return uint32(g.risks[y*g.width+x])
}

// tiledRiskAt computes the risk of cell (x, y) in a grid tiled `tiles` times
// in each direction, increasing risk by the tile distance and wrapping back
// to 1 after 9, without materializing the tiled grid.
func (g grid) tiledRiskAt(x, y int) uint32 {
	baseRisk := g.riskAt(x%g.width, y%g.height)
	tileDistance := uint32(x/g.width + y/g.height)
	return (baseRisk-1+tileDistance)%9 + 1
}

type heapItem struct {
	risk uint32
	x, y int
}

type priorityQueue []heapItem

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].risk < pq[j].risk }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(heapItem)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func lowestTotalRisk(width, height int, riskAt func(x, y int) uint32) (uint32, error) {
	const unvisited = ^uint32(0)
	best := make([]uint32, width*height)
	for i := range best {
		best[i] = unvisited
	}
	best[0] = 0

	pq := &priorityQueue{{risk: 0, x: 0, y: 0}}
	targetX, targetY := width-1, height-1

	for pq.Len() > 0 {
		item := heap.Pop(pq).(heapItem)
		if item.x == targetX && item.y == targetY {
			return item.risk, nil
		}
		if item.risk > best[item.y*width+item.x] {
			continue
		}
		type point struct{ x, y int }
		neighbors := make([]point, 0, 4)
		if item.x > 0 {
			neighbors = append(neighbors, point{item.x - 1, item.y})
		}
		if item.x+1 < width {
			neighbors = append(neighbors, point{item.x + 1, item.y})
		}
		if item.y > 0 {
			neighbors = append(neighbors, point{item.x, item.y - 1})
		}
		if item.y+1 < height {
			neighbors = append(neighbors, point{item.x, item.y + 1})
		}
		for _, n := range neighbors {
			nextRisk := item.risk + riskAt(n.x, n.y)
			index := n.y*width + n.x
			if nextRisk < best[index] {
				best[index] = nextRisk
				heap.Push(pq, heapItem{risk: nextRisk, x: n.x, y: n.y})
			}
		}
	}
	return 0, fmt.Errorf("%w: no path to the bottom-right corner", errInvalidGrid)
}

func solveWith(input string, tiles int) (uint32, uint32, error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}
	part1, err := lowestTotalRisk(g.width, g.height, g.riskAt)
	if err != nil {
		return 0, 0, err
	}
	if tiles < 1 {
		return 0, 0, fmt.Errorf("%w: tile count must be positive", errInvalidGrid)
	}
	part2, err := lowestTotalRisk(g.width*tiles, g.height*tiles, g.tiledRiskAt)
	if err != nil {
		return 0, 0, err
	}
	return part1, part2, nil
}

func solve(input string) (uint32, uint32, error) {
	return solveWith(input, 5)
}
