package main

import (
	"container/heap"
	"fmt"
	"strings"
)

var errInvalidGrid = fmt.Errorf("invalid city grid")

type grid struct {
	costs  []int
	width  int
	height int
}

func parseGrid(input string) (grid, error) {
	rows := strings.Split(strings.TrimRight(strings.ReplaceAll(input, "\r\n", "\n"), "\n"), "\n")
	if len(rows) == 0 || rows[0] == "" {
		return grid{}, fmt.Errorf("%w: grid is empty", errInvalidGrid)
	}
	width := len(rows[0])
	costs := make([]int, 0, width*len(rows))
	for _, row := range rows {
		if len(row) != width {
			return grid{}, fmt.Errorf("%w: rows have mismatched widths", errInvalidGrid)
		}
		for _, ch := range row {
			if ch < '0' || ch > '9' {
				return grid{}, fmt.Errorf("%w: contains a non-digit %q", errInvalidGrid, ch)
			}
			costs = append(costs, int(ch-'0'))
		}
	}
	return grid{costs: costs, width: width, height: len(rows)}, nil
}

func (g grid) cost(x, y int) int {
	return g.costs[y*g.width+x]
}

func (g grid) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < g.width && y < g.height
}

type direction struct{ dx, dy int }

var directions = []direction{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}

func (d direction) opposite() direction {
	return direction{-d.dx, -d.dy}
}

type crucibleState struct {
	x, y      int
	direction direction
	run       int
}

type heapEntry struct {
	heatLoss int
	state    crucibleState
}

type priorityQueue []heapEntry

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].heatLoss < pq[j].heatLoss }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(heapEntry)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	entry := old[n-1]
	*pq = old[:n-1]
	return entry
}

// leastHeatLoss runs Dijkstra over (position, last direction, run length), where the
// crucible may continue straight while run < maxRun, or turn once run >= minRun, and
// the goal only counts once it has travelled at least minRun blocks in its final direction.
func leastHeatLoss(g grid, minRun, maxRun int) (int, error) {
	goalX, goalY := g.width-1, g.height-1
	best := make(map[crucibleState]int)

	pq := &priorityQueue{{heatLoss: 0, state: crucibleState{x: 0, y: 0}}}
	heap.Init(pq)

	for pq.Len() > 0 {
		entry := heap.Pop(pq).(heapEntry)
		state := entry.state

		if state.x == goalX && state.y == goalY && state.run >= minRun {
			return entry.heatLoss, nil
		}

		for _, d := range directions {
			if state.direction == d.opposite() {
				continue
			}
			continuing := state.direction == d
			if continuing && state.run >= maxRun {
				continue
			}
			zeroDirection := direction{}
			if !continuing && state.direction != zeroDirection && state.run < minRun {
				continue
			}
			nx, ny := state.x+d.dx, state.y+d.dy
			if !g.inBounds(nx, ny) {
				continue
			}

			run := 1
			if continuing {
				run = state.run + 1
			}
			next := crucibleState{x: nx, y: ny, direction: d, run: run}
			heatLoss := entry.heatLoss + g.cost(nx, ny)
			if seen, ok := best[next]; ok && seen <= heatLoss {
				continue
			}
			best[next] = heatLoss
			heap.Push(pq, heapEntry{heatLoss: heatLoss, state: next})
		}
	}

	return 0, fmt.Errorf("%w: no path to the factory", errInvalidGrid)
}

func solve(input string) (int, int, error) {
	g, err := parseGrid(input)
	if err != nil {
		return 0, 0, err
	}
	part1, err := leastHeatLoss(g, 0, 3)
	if err != nil {
		return 0, 0, err
	}
	part2, err := leastHeatLoss(g, 4, 10)
	if err != nil {
		return 0, 0, err
	}
	return part1, part2, nil
}
