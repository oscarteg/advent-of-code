package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var errInvalidNotes = errors.New("invalid notes")

// Facings are the puzzle's own codes: 0 right, 1 down, 2 left, 3 up.
var deltas = [4]point{{1, 0}, {0, 1}, {-1, 0}, {0, -1}}

type tile byte

const (
	emptyTile tile = iota
	openTile
	wallTile
)

type point struct{ x, y int }

// move is one path instruction: walk forward tiles, or turn in place by
// quarters clockwise when quarters is non-zero.
type move struct {
	forward  int
	quarters int
}

type notes struct {
	tiles  []tile
	width  int
	height int
	path   []move
}

func splitNotes(input string) (rows []string, path string, err error) {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	for i, line := range lines {
		if line != "" {
			continue
		}
		for _, rest := range lines[i+1:] {
			if rest != "" {
				return lines[:i], rest, nil
			}
		}
		break
	}
	return nil, "", fmt.Errorf("%w: missing path", errInvalidNotes)
}

func parseNotes(input string) (notes, error) {
	rows, path, err := splitNotes(input)
	if err != nil {
		return notes{}, err
	}
	width := 0
	for _, row := range rows {
		width = max(width, len(row))
	}
	if width == 0 {
		return notes{}, fmt.Errorf("%w: missing map", errInvalidNotes)
	}

	parsed := notes{tiles: make([]tile, width*len(rows)), width: width, height: len(rows)}
	for y, row := range rows {
		for x := 0; x < len(row); x++ {
			switch row[x] {
			case ' ':
			case '.':
				parsed.tiles[y*width+x] = openTile
			case '#':
				parsed.tiles[y*width+x] = wallTile
			default:
				return notes{}, fmt.Errorf("%w: map character %q", errInvalidNotes, row[x])
			}
		}
	}
	if parsed.path, err = parsePath(path); err != nil {
		return notes{}, err
	}
	return parsed, nil
}

func parsePath(text string) ([]move, error) {
	var path []move
	forward, digits := 0, false
	for i := 0; i < len(text); i++ {
		ch := text[i]
		switch {
		case ch >= '0' && ch <= '9':
			forward = forward*10 + int(ch-'0')
			digits = true
		case ch == 'L' || ch == 'R':
			if digits {
				path = append(path, move{forward: forward})
				forward, digits = 0, false
			}
			quarters := 1
			if ch == 'L' {
				quarters = 3
			}
			path = append(path, move{quarters: quarters})
		default:
			return nil, fmt.Errorf("%w: path character %q", errInvalidNotes, ch)
		}
	}
	if digits {
		path = append(path, move{forward: forward})
	}
	if len(path) == 0 {
		return nil, fmt.Errorf("%w: path is empty", errInvalidNotes)
	}
	return path, nil
}

func (n notes) tile(p point) tile {
	if p.x < 0 || p.y < 0 || p.x >= n.width || p.y >= n.height {
		return emptyTile
	}
	return n.tiles[p.y*n.width+p.x]
}

func (n notes) start() (point, error) {
	for x := 0; x < n.width; x++ {
		if n.tile(point{x, 0}) == openTile {
			return point{x, 0}, nil
		}
	}
	return point{}, fmt.Errorf("%w: top row of the map has no open tile", errInvalidNotes)
}

// wrapFlat is part 1: leave the board and reappear on the opposite edge of the
// same row or column, keeping the facing.
func (n notes) wrapFlat(from point, facing int) (point, int, error) {
	delta := deltas[facing]
	pos := from
	for {
		pos = point{floorMod(pos.x+delta.x, n.width), floorMod(pos.y+delta.y, n.height)}
		if n.tile(pos) != emptyTile {
			return pos, facing, nil
		}
	}
}

func floorMod(value, modulus int) int {
	return ((value % modulus) + modulus) % modulus
}

type vec3 [3]int

func (a vec3) add(b vec3) vec3 { return vec3{a[0] + b[0], a[1] + b[1], a[2] + b[2]} }
func (a vec3) sub(b vec3) vec3 { return vec3{a[0] - b[0], a[1] - b[1], a[2] - b[2]} }
func (a vec3) negate() vec3    { return vec3{-a[0], -a[1], -a[2]} }
func (a vec3) dot(b vec3) int  { return a[0]*b[0] + a[1]*b[1] + a[2]*b[2] }

func (a vec3) scale(factor int) vec3 {
	return vec3{a[0] * factor, a[1] * factor, a[2] * factor}
}

// axis is the slot of one of the six axis-aligned unit normals.
func (a vec3) axis() int {
	for axis, part := range a {
		if part < 0 {
			return axis*2 + 1
		}
		if part > 0 {
			return axis * 2
		}
	}
	return 0
}

// face is one side of the folded cube: where it sits in the net plus the 3D
// directions its grid axes point in. Every face keeps the invariant
// normal = down x right, which is what makes the edge crossings exact.
type face struct {
	origin point
	right  vec3
	down   vec3
	normal vec3
}

// foldAcross rotates this basis a quarter turn around the edge shared with the
// net neighbour in the given facing.
func (f face) foldAcross(facing int, origin point) face {
	switch facing {
	case 0:
		return face{origin, f.normal.negate(), f.down, f.right}
	case 1:
		return face{origin, f.right, f.normal.negate(), f.down}
	case 2:
		return face{origin, f.normal, f.down, f.right.negate()}
	default:
		return face{origin, f.right, f.normal, f.down.negate()}
	}
}

type cube struct {
	size     int
	faces    []face
	byNormal [6]int
}

func foldCube(n notes) (cube, error) {
	count := 0
	for _, t := range n.tiles {
		if t != emptyTile {
			count++
		}
	}
	if count == 0 || count%6 != 0 {
		return cube{}, fmt.Errorf("%w: %d tiles are not six square faces", errInvalidNotes, count)
	}
	area := count / 6
	size := 1
	for size*size < area {
		size++
	}
	if size*size != area {
		return cube{}, fmt.Errorf("%w: %d tiles per face is not a square", errInvalidNotes, area)
	}

	var cells []point
	for y := 0; y < n.height; y++ {
		for x := 0; x < n.width; x++ {
			if n.tile(point{x, y}) == emptyTile {
				continue
			}
			if cell := (point{x / size, y / size}); !slices.Contains(cells, cell) {
				cells = append(cells, cell)
			}
		}
	}
	if len(cells) != 6 {
		return cube{}, fmt.Errorf("%w: net has %d faces, expected six", errInvalidNotes, len(cells))
	}
	for _, cell := range cells {
		for j := 0; j < size; j++ {
			for i := 0; i < size; i++ {
				if n.tile(point{cell.x*size + i, cell.y*size + j}) == emptyTile {
					return cube{}, fmt.Errorf("%w: face at net cell (%d, %d) has holes", errInvalidNotes, cell.x, cell.y)
				}
			}
		}
	}

	originOf := func(cell point) point { return point{cell.x * size, cell.y * size} }
	faces := make([]face, len(cells))
	assigned := make([]bool, len(cells))
	faces[0] = face{originOf(cells[0]), vec3{1, 0, 0}, vec3{0, 1, 0}, vec3{0, 0, -1}}
	assigned[0] = true
	queue := []int{0}
	for head := 0; head < len(queue); head++ {
		index := queue[head]
		for facing, delta := range deltas {
			neighbour := point{cells[index].x + delta.x, cells[index].y + delta.y}
			next := slices.Index(cells, neighbour)
			if next < 0 || assigned[next] {
				continue
			}
			faces[next] = faces[index].foldAcross(facing, originOf(cells[next]))
			assigned[next] = true
			queue = append(queue, next)
		}
	}
	if len(queue) != len(cells) {
		return cube{}, fmt.Errorf("%w: net is not a connected cube net", errInvalidNotes)
	}

	folded := cube{size: size, faces: faces}
	for slot := range folded.byNormal {
		folded.byNormal[slot] = -1
	}
	for index, f := range faces {
		slot := f.normal.axis()
		if folded.byNormal[slot] >= 0 {
			return cube{}, fmt.Errorf("%w: net folds two faces onto the same side", errInvalidNotes)
		}
		folded.byNormal[slot] = index
	}
	return folded, nil
}

func (c cube) faceAt(p point) (face, bool) {
	for _, f := range c.faces {
		inside := p.x >= f.origin.x && p.x < f.origin.x+c.size &&
			p.y >= f.origin.y && p.y < f.origin.y+c.size
		if inside {
			return f, true
		}
	}
	return face{}, false
}

// wrap is part 2: step over a cube edge. Cell (i, j) of a face sits at the
// doubled lattice point size*normal + (2i-size+1)*right + (2j-size+1)*down.
// Crossing an edge swaps the components along the old normal and the travel
// direction, which reduces to adding travel-normal, and the walk continues
// along -normal on the face whose normal is the travel direction.
func (c cube) wrap(from point, facing int) (point, int, error) {
	f, ok := c.faceAt(from)
	if !ok {
		return point{}, 0, fmt.Errorf("%w: walked off the cube at (%d, %d)", errInvalidNotes, from.x, from.y)
	}
	travel := [4]vec3{f.right, f.down, f.right.negate(), f.down.negate()}[facing]
	i, j := from.x-f.origin.x, from.y-f.origin.y
	here := f.normal.scale(c.size).
		add(f.right.scale(2*i - c.size + 1)).
		add(f.down.scale(2*j - c.size + 1))
	there := here.add(travel.sub(f.normal))

	dest := c.faces[c.byNormal[travel.axis()]]
	local := there.sub(dest.normal.scale(c.size))
	entryFacing := 3
	switch entry := f.normal.negate(); entry {
	case dest.right:
		entryFacing = 0
	case dest.down:
		entryFacing = 1
	case dest.right.negate():
		entryFacing = 2
	}
	pos := point{
		dest.origin.x + (local.dot(dest.right)+c.size-1)/2,
		dest.origin.y + (local.dot(dest.down)+c.size-1)/2,
	}
	return pos, entryFacing, nil
}

// wrapper reappears on the board after a step that would leave it.
type wrapper func(from point, facing int) (point, int, error)

func (n notes) walk(wrap wrapper) (int, error) {
	pos, err := n.start()
	if err != nil {
		return 0, err
	}
	facing := 0
	for _, m := range n.path {
		if m.quarters != 0 {
			facing = (facing + m.quarters) % 4
			continue
		}
		for step := 0; step < m.forward; step++ {
			delta := deltas[facing]
			ahead := point{pos.x + delta.x, pos.y + delta.y}
			next, nextFacing := ahead, facing
			if n.tile(ahead) == emptyTile {
				if next, nextFacing, err = wrap(pos, facing); err != nil {
					return 0, err
				}
			}
			if n.tile(next) == wallTile {
				break
			}
			pos, facing = next, nextFacing
		}
	}
	return 1000*(pos.y+1) + 4*(pos.x+1) + facing, nil
}

func solve(input string) (part1, part2 int, err error) {
	parsed, err := parseNotes(input)
	if err != nil {
		return 0, 0, err
	}
	if part1, err = parsed.walk(parsed.wrapFlat); err != nil {
		return 0, 0, err
	}
	folded, err := foldCube(parsed)
	if err != nil {
		return 0, 0, err
	}
	if part2, err = parsed.walk(folded.wrap); err != nil {
		return 0, 0, err
	}
	return part1, part2, nil
}
