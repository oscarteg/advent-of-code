// Package main solves Advent of Code 2025 day 9, Movie Theater: red tiles
// listed as x,y in the order they are walked around the room. Both parts
// span a rectangle between two red tiles and count its tiles inclusively;
// part 2 also fills the loop green and rejects rectangles that leave it.
package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

var (
	errInvalidTile = errors.New("invalid tile")
	errNotALoop    = errors.New("tiles do not form a closed loop")
)

type tile struct {
	x, y int
}

// loop is a closed loop of red tiles: consecutive tiles, wrapping around,
// are joined by a horizontal or vertical line.
type loop []tile

func parseLoop(input string) (loop, error) {
	var tiles loop
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		fields := strings.Split(trimmed, ",")
		if len(fields) != 2 {
			return nil, fmt.Errorf("%w: %q is not x,y", errInvalidTile, trimmed)
		}
		x, err := strconv.Atoi(strings.TrimSpace(fields[0]))
		if err != nil {
			return nil, fmt.Errorf("%w %q: %w", errInvalidTile, trimmed, err)
		}
		y, err := strconv.Atoi(strings.TrimSpace(fields[1]))
		if err != nil {
			return nil, fmt.Errorf("%w %q: %w", errInvalidTile, trimmed, err)
		}
		tiles = append(tiles, tile{x: x, y: y})
	}
	if len(tiles) < 4 {
		return nil, fmt.Errorf("%w: got %d tiles", errNotALoop, len(tiles))
	}
	previous := tiles[len(tiles)-1]
	for _, t := range tiles {
		straightLine := (t.x == previous.x) != (t.y == previous.y)
		if !straightLine {
			return nil, fmt.Errorf("%w: %d,%d and %d,%d are not on one line",
				errNotALoop, previous.x, previous.y, t.x, t.y)
		}
		previous = t
	}
	return tiles, nil
}

// stripLines returns the sorted grid lines of one axis: every distinct
// coordinate, plus one representative for each gap between neighbouring
// coordinates. All red and green lines lie on coordinates of the loop, so
// the tiles a representative stands in for share its colour. Sorts values.
func stripLines(values []int) []int {
	slices.Sort(values)
	lines := make([]int, 0, 2*len(values))
	for i, value := range values {
		if i > 0 {
			previous := values[i-1]
			if value == previous {
				continue
			}
			if value > previous+1 {
				lines = append(lines, previous+1)
			}
		}
		lines = append(lines, value)
	}
	return lines
}

// coverage records where the loop covers the floor, on a grid compressed to
// the loop's own lines: tileColumn and tileRow hold the strip indices of
// every tile, and uncovered is a running count of uncovered strips, which
// turns the rectangle between two red tiles into four lookups.
type coverage struct {
	tileColumn, tileRow []int
	stride              int
	uncovered           []int
}

func newCoverage(l loop) coverage {
	xs := make([]int, len(l))
	ys := make([]int, len(l))
	for i, t := range l {
		xs[i], ys[i] = t.x, t.y
	}
	columns, rows := stripLines(xs), stripLines(ys)

	tileColumn := make([]int, len(l))
	tileRow := make([]int, len(l))
	for i, t := range l {
		tileColumn[i], _ = slices.BinarySearch(columns, t.x)
		tileRow[i], _ = slices.BinarySearch(rows, t.y)
	}

	stride := len(columns) + 1
	uncovered := make([]int, (len(rows)+1)*stride)
	crossings := make([]int, len(columns))
	onLine := make([]int, stride)
	for row, y := range rows {
		clear(crossings)
		clear(onLine)
		totalCrossings := 0

		previous := len(l) - 1
		for current, t := range l {
			other := l[previous]
			switch {
			case t.x == other.x:
				top, bottom := min(other.y, t.y), max(other.y, t.y)
				if y >= top && y <= bottom {
					onLine[tileColumn[current]]++
					onLine[tileColumn[current]+1]--
					// Half-open in y, so a corner counts as one crossing.
					if y < bottom {
						crossings[tileColumn[current]]++
						totalCrossings++
					}
				}
			case t.y == y:
				onLine[min(tileColumn[previous], tileColumn[current])]++
				onLine[max(tileColumn[previous], tileColumn[current])+1]--
			}
			previous = current
		}

		// A strip is covered when a line runs through it, or when an odd
		// number of vertical lines separates it from the right edge.
		crossingsLeft, linesHere := 0, 0
		for column := range columns {
			crossingsLeft += crossings[column]
			linesHere += onLine[column]
			strip := 1
			if linesHere > 0 || (totalCrossings-crossingsLeft)%2 == 1 {
				strip = 0
			}
			uncovered[(row+1)*stride+column+1] = uncovered[(row+1)*stride+column] +
				uncovered[row*stride+column+1] - uncovered[row*stride+column] + strip
		}
	}
	return coverage{tileColumn: tileColumn, tileRow: tileRow, stride: stride, uncovered: uncovered}
}

// fullyCovered reports whether every tile between the two red tiles is red
// or green.
func (c coverage) fullyCovered(a, b int) bool {
	left, right := min(c.tileColumn[a], c.tileColumn[b]), max(c.tileColumn[a], c.tileColumn[b])+1
	top, bottom := min(c.tileRow[a], c.tileRow[b]), max(c.tileRow[a], c.tileRow[b])+1
	at := func(row, column int) int { return c.uncovered[row*c.stride+column] }
	return at(bottom, right)-at(bottom, left)-at(top, right)+at(top, left) == 0
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

// areaBetween counts the tiles of the rectangle spanned by a and b,
// including both corner tiles.
func areaBetween(a, b tile) int {
	return (abs(a.x-b.x) + 1) * (abs(a.y-b.y) + 1)
}

func solve(input string) (part1, part2 int, err error) {
	tiles, err := parseLoop(input)
	if err != nil {
		return 0, 0, err
	}
	covered := newCoverage(tiles)
	for i, a := range tiles {
		for j := i + 1; j < len(tiles); j++ {
			area := areaBetween(a, tiles[j])
			part1 = max(part1, area)
			if area > part2 && covered.fullyCovered(i, j) {
				part2 = area
			}
		}
	}
	return part1, part2, nil
}
