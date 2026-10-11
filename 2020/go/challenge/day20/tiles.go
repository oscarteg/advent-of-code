package day20

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

// seaMonster is the shape to look for once the image is assembled.
var seaMonster = []string{
	"..................#.",
	"#....##....##....###",
	".#..#..#..#..#..#...",
}

// edge is one side of an image, read left to right or top to bottom. A parsed edge
// is never empty, so the empty edge stands for a neighbour that is not there yet.
type edge string

// edgeOwners maps the canonical form of an edge to the tiles owning it. An edge
// owned by a single tile cannot be joined to anything and so lies on the border.
type edgeOwners map[edge][]int

// image is a square bitmap of pixels, set where the input has a hash.
type image struct {
	rows [][]bool
}

// tile is one camera image together with its id.
type tile struct {
	id    int
	image image
}

func parseTiles(input *challenge.Input) ([]tile, error) {
	var tiles []tile
	var block []string

	// The appended empty line closes the last block.
	for _, line := range append(input.LinesSlice(), "") {
		if line = strings.TrimSpace(line); line != "" {
			block = append(block, line)
			continue
		}
		if len(block) == 0 {
			continue
		}
		parsed, err := parseTile(block)
		if err != nil {
			return nil, err
		}
		tiles = append(tiles, parsed)
		block = nil
	}

	if len(tiles) == 0 {
		return nil, errors.New("no tiles in input")
	}
	for _, t := range tiles {
		if t.image.size() != tiles[0].image.size() {
			return nil, errors.New("tiles differ in size")
		}
	}
	if side := isqrt(len(tiles)); side*side != len(tiles) {
		return nil, fmt.Errorf("%d tiles do not form a square image", len(tiles))
	}
	return tiles, nil
}

func parseTile(block []string) (tile, error) {
	id, err := parseID(block[0])
	if err != nil {
		return tile{}, err
	}
	parsed, err := parseImage(block[1:])
	if err != nil {
		return tile{}, fmt.Errorf("tile %d: %w", id, err)
	}
	return tile{id: id, image: parsed}, nil
}

func parseID(header string) (int, error) {
	if !strings.HasPrefix(header, "Tile ") || !strings.HasSuffix(header, ":") {
		return 0, fmt.Errorf("expected a \"Tile <id>:\" header, got %q", header)
	}
	id, err := strconv.Atoi(strings.TrimSuffix(strings.TrimPrefix(header, "Tile "), ":"))
	if err != nil {
		return 0, fmt.Errorf("invalid tile id in %q: %w", header, err)
	}
	return id, nil
}

func parseImage(lines []string) (image, error) {
	if len(lines) < 3 {
		return image{}, fmt.Errorf("a tile needs at least 3 rows, got %d", len(lines))
	}
	rows := make([][]bool, len(lines))
	for y, line := range lines {
		if len(line) != len(lines) {
			return image{}, fmt.Errorf("tile is not square: %d rows but a row of %d", len(lines), len(line))
		}
		row := make([]bool, len(line))
		for x := 0; x < len(line); x++ {
			switch line[x] {
			case '#':
				row[x] = true
			case '.':
				row[x] = false
			default:
				return image{}, fmt.Errorf("invalid pixel %q", line[x])
			}
		}
		rows[y] = row
	}
	return image{rows: rows}, nil
}

func (img image) size() int {
	return len(img.rows)
}

func (img image) top() edge {
	return toEdge(img.rows[0])
}

func (img image) bottom() edge {
	return toEdge(img.rows[len(img.rows)-1])
}

func (img image) left() edge {
	pixels := make([]bool, img.size())
	for y, row := range img.rows {
		pixels[y] = row[0]
	}
	return toEdge(pixels)
}

func (img image) right() edge {
	pixels := make([]bool, img.size())
	for y, row := range img.rows {
		pixels[y] = row[len(row)-1]
	}
	return toEdge(pixels)
}

func (img image) edges() []edge {
	return []edge{img.top(), img.bottom(), img.left(), img.right()}
}

func toEdge(pixels []bool) edge {
	var text strings.Builder
	for _, lit := range pixels {
		if lit {
			text.WriteByte('#')
		} else {
			text.WriteByte('.')
		}
	}
	return edge(text.String())
}

// canonical is the same value for both reading directions of an edge.
func canonical(e edge) edge {
	reversed := make([]byte, len(e))
	for i := 0; i < len(e); i++ {
		reversed[i] = e[len(e)-1-i]
	}
	if backward := edge(reversed); backward < e {
		return backward
	}
	return e
}

// rotated is a quarter turn clockwise.
func (img image) rotated() image {
	size := img.size()
	rows := make([][]bool, size)
	for y := 0; y < size; y++ {
		row := make([]bool, size)
		for x := 0; x < size; x++ {
			row[x] = img.rows[size-1-x][y]
		}
		rows[y] = row
	}
	return image{rows: rows}
}

// flipped is mirrored left to right.
func (img image) flipped() image {
	rows := make([][]bool, img.size())
	for y, row := range img.rows {
		mirrored := make([]bool, len(row))
		for x, pixel := range row {
			mirrored[len(row)-1-x] = pixel
		}
		rows[y] = mirrored
	}
	return image{rows: rows}
}

// orientations are the four rotations of the image followed by those of its mirror.
func (img image) orientations() []image {
	orientations := make([]image, 0, 8)
	for _, start := range []image{img, img.flipped()} {
		current := start
		for i := 0; i < 4; i++ {
			orientations = append(orientations, current)
			current = current.rotated()
		}
	}
	return orientations
}

// interior is the image without its one pixel wide border.
func (img image) interior() image {
	rows := make([][]bool, 0, img.size()-2)
	for _, row := range img.rows[1 : img.size()-1] {
		rows = append(rows, row[1:len(row)-1])
	}
	return image{rows: rows}
}

func (img image) litPixels() int {
	lit := 0
	for _, row := range img.rows {
		for _, pixel := range row {
			if pixel {
				lit++
			}
		}
	}
	return lit
}

// pattern is a shape that need not be square, stored as the offsets of its lit pixels.
type pattern struct {
	lit           [][2]int
	height, width int
}

func newPattern(lines []string) pattern {
	shape := pattern{height: len(lines)}
	for y, line := range lines {
		if len(line) > shape.width {
			shape.width = len(line)
		}
		for x := 0; x < len(line); x++ {
			if line[x] == '#' {
				shape.lit = append(shape.lit, [2]int{y, x})
			}
		}
	}
	return shape
}

// occurrences counts how often the shape sits entirely on lit pixels of the image.
func (img image) occurrences(shape pattern) int {
	size := img.size()
	if size < shape.height || size < shape.width {
		return 0
	}
	found := 0
	for row := 0; row <= size-shape.height; row++ {
		for column := 0; column <= size-shape.width; column++ {
			if img.contains(shape, row, column) {
				found++
			}
		}
	}
	return found
}

func (img image) contains(shape pattern, row, column int) bool {
	for _, offset := range shape.lit {
		if !img.rows[row+offset[0]][column+offset[1]] {
			return false
		}
	}
	return true
}

func ownersByEdge(tiles []tile) edgeOwners {
	owners := edgeOwners{}
	for index, t := range tiles {
		for _, e := range t.image.edges() {
			key := canonical(e)
			owners[key] = append(owners[key], index)
		}
	}
	return owners
}

// isBorder reports whether no other tile can be joined to this edge.
func (owners edgeOwners) isBorder(e edge) bool {
	return len(owners[canonical(e)]) == 1
}

// cornerIndices are the tiles with exactly two edges no other tile can be joined to.
func cornerIndices(tiles []tile, owners edgeOwners) ([]int, error) {
	var corners []int
	for index, t := range tiles {
		borders := 0
		for _, e := range t.image.edges() {
			if owners.isBorder(e) {
				borders++
			}
		}
		if borders == 2 {
			corners = append(corners, index)
		}
	}
	if len(corners) != 4 {
		return nil, fmt.Errorf("expected 4 corner tiles, found %d", len(corners))
	}
	return corners, nil
}

// slot is one position of the layout and the edges a tile must present there.
type slot struct {
	row, column, side int
	above, left       edge
}

// joinedEdge is an edge shared with a placed neighbour, owned by just one other tile.
func (s slot) joinedEdge() edge {
	if s.left != "" {
		return s.left
	}
	return s.above
}

func (s slot) accepts(img image, owners edgeOwners) bool {
	joins := func(side, neighbour edge) bool {
		if neighbour == "" {
			return owners.isBorder(side)
		}
		return side == neighbour
	}
	return joins(img.top(), s.above) &&
		joins(img.left(), s.left) &&
		(s.row+1 < s.side || owners.isBorder(img.bottom())) &&
		(s.column+1 < s.side || owners.isBorder(img.right()))
}

// placement is a tile in the orientation it takes in the layout.
type placement struct {
	tile  int
	image image
}

// assemble lays the tiles out row by row and returns the stitched image.
func assemble(tiles []tile, owners edgeOwners, corners []int) (image, error) {
	side := isqrt(len(tiles))
	layout := make([]image, 0, len(tiles))
	used := make([]bool, len(tiles))

	for position := range tiles {
		current := slot{row: position / side, column: position % side, side: side}
		if current.row > 0 {
			current.above = layout[position-side].bottom()
		}
		if current.column > 0 {
			current.left = layout[position-1].right()
		}

		candidates := corners
		if joined := current.joinedEdge(); joined != "" {
			candidates = owners[canonical(joined)]
		}

		var fitting []placement
		for _, index := range candidates {
			if used[index] {
				continue
			}
			for _, oriented := range tiles[index].image.orientations() {
				if current.accepts(oriented, owners) {
					fitting = append(fitting, placement{tile: index, image: oriented})
				}
			}
		}

		if len(fitting) == 0 {
			return image{}, fmt.Errorf("no tile fits at row %d, column %d", current.row, current.column)
		}
		// Every slot but the first is pinned by a neighbour edge, so a second fit there
		// means the tiles can be laid out in more than one way. The first slot accepts any
		// corner in either of its two border-aligned orientations; all eight choices
		// assemble the same image up to rotation and reflection, which both parts ignore.
		if position > 0 && len(fitting) > 1 {
			return image{}, fmt.Errorf("%d tiles fit at row %d, column %d", len(fitting), current.row, current.column)
		}
		used[fitting[0].tile] = true
		layout = append(layout, fitting[0].image)
	}

	return stitch(layout, side), nil
}

// stitch joins the laid out tiles into one image, dropping the borders that lined them up.
func stitch(layout []image, side int) image {
	var rows [][]bool
	for band := 0; band < len(layout); band += side {
		interiors := make([]image, 0, side)
		for _, img := range layout[band : band+side] {
			interiors = append(interiors, img.interior())
		}
		for row := 0; row < interiors[0].size(); row++ {
			var merged []bool
			for _, interior := range interiors {
				merged = append(merged, interior.rows[row]...)
			}
			rows = append(rows, merged)
		}
	}
	return image{rows: rows}
}

// waterRoughness is the lit pixels left once every sea monster is marked. Monsters
// appear in exactly one orientation, so the orientation with the most is the right one.
func waterRoughness(img image) int {
	monster := newPattern(seaMonster)
	most := 0
	for _, view := range img.orientations() {
		if found := view.occurrences(monster); found > most {
			most = found
		}
	}
	return img.litPixels() - most*len(monster.lit)
}

func isqrt(n int) int {
	root := 0
	for (root+1)*(root+1) <= n {
		root++
	}
	return root
}
