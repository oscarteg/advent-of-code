package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	part1Steps = 64
	part2Steps = 26_501_365
)

type point struct {
	x, y int
}

type garden struct {
	rocks         []bool
	width, height int
	start         point
}

func solve(input string) (int, int, error) {
	return solveWith(input, part1Steps, part2Steps)
}

func solveWith(input string, boundedSteps, tiledSteps int) (int, int, error) {
	g, err := parseGarden(input)
	if err != nil {
		return 0, 0, err
	}
	tiled, err := g.reachableExtrapolated(tiledSteps)
	if err != nil {
		return 0, 0, err
	}
	return g.reachable(boundedSteps), tiled, nil
}

func parseGarden(input string) (*garden, error) {
	var lines []string
	for _, line := range strings.Split(input, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	if len(lines) == 0 {
		return nil, errors.New("parse garden: no map")
	}

	g := &garden{width: len(lines[0]), height: len(lines), start: point{x: -1, y: -1}}
	g.rocks = make([]bool, 0, g.width*g.height)
	for y, line := range lines {
		if len(line) != g.width {
			return nil, fmt.Errorf("parse garden: row %d holds %d tiles, want %d", y, len(line), g.width)
		}
		for x := 0; x < len(line); x++ {
			switch line[x] {
			case '.':
				g.rocks = append(g.rocks, false)
			case '#':
				g.rocks = append(g.rocks, true)
			case 'S':
				if g.start.x >= 0 {
					return nil, errors.New("parse garden: more than one starting plot")
				}
				g.start = point{x: x, y: y}
				g.rocks = append(g.rocks, false)
			default:
				return nil, fmt.Errorf("parse garden: unknown tile %q at row %d", rune(line[x]), y)
			}
		}
	}
	if g.start.x < 0 {
		return nil, errors.New("parse garden: no starting plot")
	}
	return g, nil
}

// reachable counts the plots reachable in exactly steps steps, the map edges
// acting as rocks.
func (g *garden) reachable(steps int) int {
	return g.flood(steps, false)
}

// reachableTiled counts the plots reachable in exactly steps steps on the map
// repeated infinitely in every direction.
func (g *garden) reachableTiled(steps int) int {
	return g.flood(steps, true)
}

// reachableExtrapolated answers step counts too large to flood. Every width
// steps the flood swallows one more copy of the map in each direction, so on a
// square map with an open border and an open row and column through the start
// the counts after offset+tiles*width steps are a quadratic in tiles. Fit that
// quadratic through the first three tile-aligned floods and evaluate it.
func (g *garden) reachableExtrapolated(steps int) (int, error) {
	if steps < 2*g.width {
		return g.reachableTiled(steps), nil
	}
	if g.width != g.height {
		return 0, fmt.Errorf("extrapolation needs a square map, got %dx%d", g.width, g.height)
	}
	offset := steps % g.width
	first := g.reachableTiled(offset)
	second := g.reachableTiled(offset + g.width)
	third := g.reachableTiled(offset + 2*g.width)
	tiles := steps / g.width
	return first + tiles*(second-first) + tiles*(tiles-1)/2*(third-2*second+first), nil
}

func (g *garden) flood(steps int, tiled bool) int {
	// The elf can step off a plot and back onto it, so every plot whose
	// distance has the same parity as the step count is reachable - unless he
	// is walled in and cannot take a single step.
	if steps > 0 && g.isWalledIn(tiled) {
		return 0
	}

	seen := map[point]bool{g.start: true}
	frontier := []point{g.start}
	reachable := 0
	for remaining := steps; remaining >= 0; remaining-- {
		if remaining%2 == 0 {
			reachable += len(frontier)
		}
		var next []point
		for _, plot := range frontier {
			for _, neighbour := range neighbours(plot) {
				if !g.isPlot(neighbour, tiled) || seen[neighbour] {
					continue
				}
				seen[neighbour] = true
				next = append(next, neighbour)
			}
		}
		frontier = next
	}
	return reachable
}

func (g *garden) isWalledIn(tiled bool) bool {
	for _, neighbour := range neighbours(g.start) {
		if g.isPlot(neighbour, tiled) {
			return false
		}
	}
	return true
}

func (g *garden) isPlot(p point, tiled bool) bool {
	x, y := p.x, p.y
	if tiled {
		x, y = floorMod(x, g.width), floorMod(y, g.height)
	} else if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return false
	}
	return !g.rocks[y*g.width+x]
}

func neighbours(p point) [4]point {
	return [4]point{{p.x + 1, p.y}, {p.x - 1, p.y}, {p.x, p.y + 1}, {p.x, p.y - 1}}
}

// floorMod wraps negative coordinates onto the map; Go's % keeps the sign of
// the dividend.
func floorMod(value, divisor int) int {
	return ((value % divisor) + divisor) % divisor
}
