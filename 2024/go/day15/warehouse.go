package main

import (
	"errors"
	"fmt"
	"strings"
)

var (
	errInvalidWarehouse = errors.New("invalid warehouse map")
	errInvalidMoves     = errors.New("invalid robot moves")
)

type vector struct {
	row, col int
}

func (v vector) plus(other vector) vector {
	return vector{v.row + other.row, v.col + other.col}
}

func (v vector) minus(other vector) vector {
	return vector{v.row - other.row, v.col - other.col}
}

// The halves of a widened box sit one column apart.
var (
	rightHalf = vector{0, 1}
	leftHalf  = vector{0, -1}
)

// warehouse is a rectangular map enclosed by walls, stored row-major, plus the
// robot's position. That wall border is what makes every push stop instead of
// leaving the map.
type warehouse struct {
	tiles  []byte
	width  int
	height int
	robot  vector
}

type puzzle struct {
	warehouse warehouse
	moves     []vector
}

func (w warehouse) indexOf(at vector) int {
	return at.row*w.width + at.col
}

func (w warehouse) tile(at vector) byte {
	return w.tiles[w.indexOf(at)]
}

func (w *warehouse) setTile(at vector, value byte) {
	w.tiles[w.indexOf(at)] = value
}

func (w warehouse) enclosed() bool {
	for col := range w.width {
		if w.tile(vector{0, col}) != '#' || w.tile(vector{w.height - 1, col}) != '#' {
			return false
		}
	}
	for row := range w.height {
		if w.tile(vector{row, 0}) != '#' || w.tile(vector{row, w.width - 1}) != '#' {
			return false
		}
	}
	return true
}

func (w warehouse) gpsSum(box byte) int {
	sum := 0
	for row := range w.height {
		for col := range w.width {
			if w.tile(vector{row, col}) == box {
				sum += 100*row + col
			}
		}
	}
	return sum
}

func directionOf(move byte) (vector, error) {
	switch move {
	case '^':
		return vector{-1, 0}, nil
	case 'v':
		return vector{1, 0}, nil
	case '<':
		return vector{0, -1}, nil
	case '>':
		return vector{0, 1}, nil
	default:
		return vector{}, fmt.Errorf("%w: move %q", errInvalidMoves, move)
	}
}

// splitSections separates the map from the moves: the first blank line after
// the map ends it, and line breaks inside the moves carry no meaning.
func splitSections(input string) (mapLines, moveLines []string) {
	readingMap := true
	for _, raw := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case line == "":
			readingMap = readingMap && len(mapLines) == 0
		case readingMap:
			mapLines = append(mapLines, line)
		default:
			moveLines = append(moveLines, line)
		}
	}
	return mapLines, moveLines
}

func parseWarehouse(lines []string) (warehouse, error) {
	if len(lines) == 0 {
		return warehouse{}, fmt.Errorf("%w: no map rows", errInvalidWarehouse)
	}

	parsed := warehouse{width: len(lines[0]), height: len(lines)}
	robots := 0
	for row, line := range lines {
		if len(line) != parsed.width {
			return warehouse{}, fmt.Errorf("%w: row %d is %d tiles wide, want %d",
				errInvalidWarehouse, row, len(line), parsed.width)
		}
		for col := range len(line) {
			switch line[col] {
			case '#', '.', 'O':
			case '@':
				robots++
				parsed.robot = vector{row, col}
			default:
				return warehouse{}, fmt.Errorf("%w: tile %q", errInvalidWarehouse, line[col])
			}
		}
		parsed.tiles = append(parsed.tiles, line...)
	}

	if robots != 1 {
		return warehouse{}, fmt.Errorf("%w: found %d robots, want exactly one",
			errInvalidWarehouse, robots)
	}
	if !parsed.enclosed() {
		return warehouse{}, fmt.Errorf("%w: map is not enclosed by walls", errInvalidWarehouse)
	}
	return parsed, nil
}

func parseMoves(lines []string) ([]vector, error) {
	var moves []vector
	for _, line := range lines {
		for i := range len(line) {
			direction, err := directionOf(line[i])
			if err != nil {
				return nil, err
			}
			moves = append(moves, direction)
		}
	}
	if len(moves) == 0 {
		return nil, fmt.Errorf("%w: no moves", errInvalidMoves)
	}
	return moves, nil
}

func parse(input string) (puzzle, error) {
	mapLines, moveLines := splitSections(input)
	parsedMap, err := parseWarehouse(mapLines)
	if err != nil {
		return puzzle{}, err
	}
	moves, err := parseMoves(moveLines)
	if err != nil {
		return puzzle{}, err
	}
	return puzzle{warehouse: parsedMap, moves: moves}, nil
}

// push is part 1: the robot shoves the whole run of boxes ahead of it one tile
// along the move, and only moves at all when the tile past that run is free.
// Hopping the near box to the far end leaves the same map as shifting them all.
func (w *warehouse) push(direction vector) {
	target := w.robot.plus(direction)
	free := target
	for w.tile(free) == 'O' {
		free = free.plus(direction)
	}
	if w.tile(free) == '#' {
		return
	}

	w.setTile(free, 'O')
	w.setTile(target, '@')
	w.setTile(w.robot, '.')
	w.robot = target
}

// widened doubles every tile's width for part 2: walls and floor repeat, a box
// becomes the pair "[]", and the robot keeps a single tile of its own.
func (w warehouse) widened() warehouse {
	tiles := make([]byte, 0, len(w.tiles)*2)
	for _, current := range w.tiles {
		switch current {
		case 'O':
			tiles = append(tiles, '[', ']')
		case '@':
			tiles = append(tiles, '@', '.')
		default:
			tiles = append(tiles, current, current)
		}
	}
	return warehouse{
		tiles:  tiles,
		width:  w.width * 2,
		height: w.height,
		robot:  vector{w.robot.row, w.robot.col * 2},
	}
}

func (w *warehouse) pushWide(direction vector) {
	if direction.row == 0 {
		w.pushWideRow(direction)
		return
	}
	w.pushWideColumn(direction)
}

// pushWideRow shifts the run of box halves beside the robot one tile sideways;
// the halves stay paired, so no box can come apart.
func (w *warehouse) pushWideRow(direction vector) {
	target := w.robot.plus(direction)
	free := target
	for w.tile(free) == '[' || w.tile(free) == ']' {
		free = free.plus(direction)
	}
	if w.tile(free) == '#' {
		return
	}

	for at := free; at != target; at = at.minus(direction) {
		w.setTile(at, w.tile(at.minus(direction)))
	}
	w.setTile(target, '@')
	w.setTile(w.robot, '.')
	w.robot = target
}

// pushedBoxes collects the left halves of every box a vertical push would move,
// nearest first, and reports whether the push is possible at all. A box can be
// reached through either half and the list stays small, so a linear duplicate
// check beats tracking a grid of visited tiles.
func (w warehouse) pushedBoxes(direction vector) (boxes []vector, movable bool) {
	appendBox := func(left vector) {
		for _, existing := range boxes {
			if existing == left {
				return
			}
		}
		boxes = append(boxes, left)
	}

	target := w.robot.plus(direction)
	switch w.tile(target) {
	case '#':
		return nil, false
	case '[':
		appendBox(target)
	case ']':
		appendBox(target.plus(leftHalf))
	default:
		return nil, true
	}

	for index := 0; index < len(boxes); index++ {
		left := boxes[index]
		for _, half := range [2]vector{left, left.plus(rightHalf)} {
			next := half.plus(direction)
			switch w.tile(next) {
			case '#':
				return nil, false
			case '[':
				appendBox(next)
			case ']':
				appendBox(next.plus(leftHalf))
			}
		}
	}
	return boxes, true
}

func (w *warehouse) pushWideColumn(direction vector) {
	boxes, movable := w.pushedBoxes(direction)
	if !movable {
		return
	}

	// Farthest box first, so each one lands on tiles the previous box left.
	for index := len(boxes) - 1; index >= 0; index-- {
		left := boxes[index]
		right := left.plus(rightHalf)
		w.setTile(left, '.')
		w.setTile(right, '.')
		w.setTile(left.plus(direction), '[')
		w.setTile(right.plus(direction), ']')
	}
	w.setTile(w.robot, '.')
	w.robot = w.robot.plus(direction)
	w.setTile(w.robot, '@')
}

func solve(input string) (part1, part2 int, err error) {
	parsed, err := parse(input)
	if err != nil {
		return 0, 0, fmt.Errorf("parse input: %w", err)
	}

	narrow := parsed.warehouse
	wide := narrow.widened()
	for _, direction := range parsed.moves {
		narrow.push(direction)
		wide.pushWide(direction)
	}
	return narrow.gpsSum('O'), wide.gpsSum('['), nil
}
