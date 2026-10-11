package main

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"strings"
)

var (
	errInvalidMaze = errors.New("invalid maze map")
	errNoPath      = errors.New("no path from the start to the end tile")
)

type step struct {
	row, col int
}

// facings are in clockwise order, so a 90 degree turn is +-1 modulo four.
var facings = [4]step{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

const (
	facingCount = len(facings)
	east        = 1
	stepScore   = 1
	turnScore   = 1000
	unreachable = math.MaxInt
)

// maze is a row-major grid of tiles; start and end are tile indices.
type maze struct {
	tiles  []byte
	width  int
	height int
	start  int
	end    int
}

// move is one state of the search: a (tile, facing) pair and what reaching it
// costs.
type move struct {
	score int
	state int
}

// frontier is a min-heap of moves ordered by score.
type frontier []move

func (f frontier) Len() int           { return len(f) }
func (f frontier) Less(i, j int) bool { return f[i].score < f[j].score }
func (f frontier) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func (f *frontier) Push(x any) { *f = append(*f, x.(move)) }

func (f *frontier) Pop() any {
	cheapest := (*f)[len(*f)-1]
	*f = (*f)[:len(*f)-1]
	return cheapest
}

func parseMaze(input string) (maze, error) {
	var parsed maze
	start, end := -1, -1
	for _, raw := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		if parsed.height == 0 {
			parsed.width = len(line)
		} else if len(line) != parsed.width {
			return maze{}, fmt.Errorf("%w: row %d is %d tiles wide, want %d", errInvalidMaze, parsed.height, len(line), parsed.width)
		}
		for col := range parsed.width {
			tile := parsed.height*parsed.width + col
			switch line[col] {
			case '#', '.':
			case 'S':
				if start >= 0 {
					return maze{}, fmt.Errorf("%w: two start tiles", errInvalidMaze)
				}
				start = tile
			case 'E':
				if end >= 0 {
					return maze{}, fmt.Errorf("%w: two end tiles", errInvalidMaze)
				}
				end = tile
			default:
				return maze{}, fmt.Errorf("%w: tile %q is not #, ., S or E", errInvalidMaze, line[col])
			}
		}
		parsed.tiles = append(parsed.tiles, line...)
		parsed.height++
	}
	switch {
	case parsed.height == 0:
		return maze{}, fmt.Errorf("%w: no maze rows", errInvalidMaze)
	case start < 0:
		return maze{}, fmt.Errorf("%w: no start tile", errInvalidMaze)
	case end < 0:
		return maze{}, fmt.Errorf("%w: no end tile", errInvalidMaze)
	}
	parsed.start, parsed.end = start, end
	return parsed, nil
}

func (m maze) isWall(tile int) bool { return m.tiles[tile] == '#' }

// stepFrom returns the tile one step along the facing, or against it when
// walking backwards, and whether that tile exists and is not a wall.
func (m maze) stepFrom(tile, facing int, backwards bool) (int, bool) {
	forward := facings[facing]
	if backwards {
		forward = step{-forward.row, -forward.col}
	}
	row, col := tile/m.width+forward.row, tile%m.width+forward.col
	if row < 0 || col < 0 || row >= m.height || col >= m.width {
		return 0, false
	}
	next := row*m.width + col
	return next, !m.isWall(next)
}

// scoresFrom runs Dijkstra over (tile, facing) states, where a turn costs 1000
// and a step forward 1. Walking backwards explores the reverse graph, which
// scores every state by its distance to the nearest source instead of from it.
func (m maze) scoresFrom(sources []int, backwards bool) []int {
	best := make([]int, len(m.tiles)*facingCount)
	for state := range best {
		best[state] = unreachable
	}
	queue := make(frontier, 0, len(sources))
	for _, source := range sources {
		best[source] = 0
		queue = append(queue, move{score: 0, state: source})
	}
	heap.Init(&queue)

	for queue.Len() > 0 {
		current := heap.Pop(&queue).(move)
		if current.score > best[current.state] {
			continue
		}
		tile, facing := current.state/facingCount, current.state%facingCount

		moves := [3]move{
			{score: current.score + turnScore, state: tile*facingCount + (facing+1)%facingCount},
			{score: current.score + turnScore, state: tile*facingCount + (facing+3)%facingCount},
		}
		moveCount := 2
		if next, open := m.stepFrom(tile, facing, backwards); open {
			moves[moveCount] = move{score: current.score + stepScore, state: next*facingCount + facing}
			moveCount++
		}

		for _, candidate := range moves[:moveCount] {
			if candidate.score >= best[candidate.state] {
				continue
			}
			best[candidate.state] = candidate.score
			heap.Push(&queue, candidate)
		}
	}
	return best
}

func solve(input string) (part1, part2 int, err error) {
	parsed, err := parseMaze(input)
	if err != nil {
		return 0, 0, err
	}

	fromStart := parsed.scoresFrom([]int{parsed.start*facingCount + east}, false)
	endStates := make([]int, 0, facingCount)
	for facing := range facingCount {
		endStates = append(endStates, parsed.end*facingCount+facing)
	}
	toEnd := parsed.scoresFrom(endStates, true)

	best := unreachable
	for _, state := range endStates {
		best = min(best, fromStart[state])
	}
	if best == unreachable {
		return 0, 0, errNoPath
	}

	// A state lies on a best path when reaching it and leaving it towards the
	// end together cost exactly the best score.
	tilesOnBestPaths := 0
	for tile := range len(parsed.tiles) {
		for facing := range facingCount {
			state := tile*facingCount + facing
			if fromStart[state] == unreachable || toEnd[state] == unreachable {
				continue
			}
			if fromStart[state]+toEnd[state] != best {
				continue
			}
			tilesOnBestPaths++
			break
		}
	}
	return best, tilesOnBestPaths, nil
}
