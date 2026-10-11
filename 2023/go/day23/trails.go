package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errInvalidMap = errors.New("invalid trail map")
	errNoHike     = errors.New("no hike reaches the exit")
)

type tile byte

const (
	path   tile = '.'
	forest tile = '#'
)

// slopes tells a hike how to treat slope tiles.
type slopes int

const (
	// icy forces the step off a slope to go downhill (part 1).
	icy slopes = iota
	// dry treats slopes as ordinary path tiles (part 2).
	dry
)

type direction struct{ row, column int }

var directions = [...]direction{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

// downhill holds the only step allowed off each slope tile.
var downhill = map[tile]direction{'^': {-1, 0}, '>': {0, 1}, 'v': {1, 0}, '<': {0, -1}}

type position struct{ row, column int }

// trailMap is a rectangular tile grid entered through the single path tile of
// the top row and left through the single path tile of the bottom row.
type trailMap struct {
	tiles      [][]tile
	start, end position
}

func parseMap(input string) (trailMap, error) {
	var tiles [][]tile
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		row := make([]tile, 0, len(line))
		for _, character := range []byte(line) {
			if !known(tile(character)) {
				return trailMap{}, fmt.Errorf("%w: unknown tile %q", errInvalidMap, character)
			}
			row = append(row, tile(character))
		}
		if len(tiles) > 0 && len(row) != len(tiles[0]) {
			return trailMap{}, fmt.Errorf("%w: rows differ in length", errInvalidMap)
		}
		tiles = append(tiles, row)
	}
	if len(tiles) == 0 {
		return trailMap{}, fmt.Errorf("%w: empty map", errInvalidMap)
	}
	start, ok := pathColumn(tiles[0])
	if !ok {
		return trailMap{}, fmt.Errorf("%w: no path tile in the top row", errInvalidMap)
	}
	end, ok := pathColumn(tiles[len(tiles)-1])
	if !ok {
		return trailMap{}, fmt.Errorf("%w: no path tile in the bottom row", errInvalidMap)
	}
	return trailMap{
		tiles: tiles,
		start: position{row: 0, column: start},
		end:   position{row: len(tiles) - 1, column: end},
	}, nil
}

func known(t tile) bool {
	if t == path || t == forest {
		return true
	}
	_, ok := downhill[t]
	return ok
}

func pathColumn(row []tile) (int, bool) {
	for column, t := range row {
		if t == path {
			return column, true
		}
	}
	return 0, false
}

// at reports the tile at p; tiles outside the grid count as forest, which keeps
// every walk inside it.
func (m trailMap) at(p position) tile {
	if p.row < 0 || p.row >= len(m.tiles) || p.column < 0 || p.column >= len(m.tiles[p.row]) {
		return forest
	}
	return m.tiles[p.row][p.column]
}

func (m trailMap) walkable(p position) bool { return m.at(p) != forest }

// steps lists the tiles reachable from p in a single step.
func (m trailMap) steps(p position, s slopes) []position {
	allowed := directions[:]
	if forced, ok := downhill[m.at(p)]; ok && s == icy {
		allowed = []direction{forced}
	}
	next := make([]position, 0, len(allowed))
	for _, d := range allowed {
		to := position{row: p.row + d.row, column: p.column + d.column}
		if m.walkable(to) {
			next = append(next, to)
		}
	}
	return next
}

// longestHike counts the steps of the longest hike from the entrance to the exit.
func (m trailMap) longestHike(s slopes) (int, error) {
	steps, ok := m.compress(s).longestHike()
	if !ok {
		return 0, errNoHike
	}
	return steps, nil
}

type edge struct{ to, steps int }

// graph is the trail map with its corridors collapsed: the nodes are the
// entrance, the exit and the forks, and every edge carries the steps along one
// corridor.
type graph struct {
	edges      [][]edge
	start, end int
}

// nodes collects the entrance, the exit and every tile where the trails fork.
// Slopes never change which tiles fork, so both parts share one set of nodes.
func (m trailMap) nodes() []position {
	nodes := []position{m.start}
	if m.end != m.start {
		nodes = append(nodes, m.end)
	}
	for row, tiles := range m.tiles {
		for column, t := range tiles {
			p := position{row: row, column: column}
			if t != forest && p != m.start && p != m.end && len(m.steps(p, dry)) > 2 {
				nodes = append(nodes, p)
			}
		}
	}
	return nodes
}

// compress replaces the corridors between the nodes by weighted edges.
func (m trailMap) compress(s slopes) graph {
	nodes := m.nodes()
	indices := make(map[position]int, len(nodes))
	for index, p := range nodes {
		indices[p] = index
	}
	edges := make([][]edge, len(nodes))
	for from, p := range nodes {
		for _, first := range m.steps(p, s) {
			if corridor, ok := m.follow(p, first, s, indices); ok {
				edges[from] = append(edges[from], corridor)
			}
		}
	}
	// nodes lists the entrance first and the exit second; a single-row map is
	// entered and left through the same tile and so has one node.
	end := 0
	if m.end != m.start {
		end = 1
	}
	return graph{edges: edges, start: 0, end: end}
}

// follow walks the corridor leaving a node through first up to the next node. A
// corridor that dead-ends, or that leads into a slope pointing back, yields no
// edge.
func (m trailMap) follow(node, first position, s slopes, indices map[position]int) (edge, bool) {
	previous, current, steps := node, first, 1
	for {
		if to, ok := indices[current]; ok {
			return edge{to: to, steps: steps}, true
		}
		next, ok := onward(m.steps(current, s), previous)
		if !ok {
			return edge{}, false
		}
		previous, current, steps = current, next, steps+1
	}
}

func onward(options []position, previous position) (position, bool) {
	for _, option := range options {
		if option != previous {
			return option, true
		}
	}
	return position{}, false
}

func (g graph) longestHike() (int, bool) {
	return g.search(g.start, make([]bool, len(g.edges)))
}

// search enumerates every hike that visits each node at most once: finding the
// longest path is NP-hard, so there is no shortcut.
func (g graph) search(node int, visited []bool) (int, bool) {
	if node == g.end {
		return 0, true
	}
	visited[node] = true
	longest, found := 0, false
	for _, e := range g.edges[node] {
		if visited[e.to] {
			continue
		}
		rest, ok := g.search(e.to, visited)
		if !ok {
			continue
		}
		if total := e.steps + rest; !found || total > longest {
			longest, found = total, true
		}
	}
	visited[node] = false
	return longest, found
}

func solve(input string) (int, int, error) {
	trails, err := parseMap(input)
	if err != nil {
		return 0, 0, err
	}
	part1, err := trails.longestHike(icy)
	if err != nil {
		return 0, 0, err
	}
	part2, err := trails.longestHike(dry)
	if err != nil {
		return 0, 0, err
	}
	return part1, part2, nil
}
