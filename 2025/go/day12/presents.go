package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// presentSize is the width and height every present pattern is given in, so
// each orientation of a present fits a presentSize x presentSize box.
const presentSize = 3

var errInvalidFarm = errors.New("invalid christmas tree farm")

type pattern [presentSize][presentSize]bool

// offset is one cell of an oriented present, relative to the cell the present
// is anchored at: its own first cell in row-major order, so rows never go up.
type offset struct {
	row    int
	column int
}

// variant is a present in one of its distinct orientations.
type variant struct {
	cells  []offset
	height int
}

type present struct {
	area     int
	variants []variant
}

// region is the space under one tree: its size plus how many of each present
// must fit into it.
type region struct {
	width  int
	height int
	counts []int
}

type farm struct {
	presents []present
	regions  []region
}

func solve(input string) (int, error) {
	parsed, err := parseFarm(input)
	if err != nil {
		return 0, err
	}
	fitting := 0
	for _, tree := range parsed.regions {
		if regionFits(parsed.presents, tree) {
			fitting++
		}
	}
	return fitting, nil
}

func parseFarm(input string) (farm, error) {
	var parsed farm
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	for index := 0; index < len(lines); index++ {
		line := strings.TrimSpace(lines[index])
		switch {
		case line == "":
			continue
		case strings.HasSuffix(line, ":"):
			if len(parsed.regions) > 0 {
				return farm{}, fmt.Errorf("%w: present after a region", errInvalidFarm)
			}
			shape, err := strconv.Atoi(strings.TrimSuffix(line, ":"))
			if err != nil {
				return farm{}, fmt.Errorf("%w: present header %q: %w", errInvalidFarm, line, err)
			}
			if shape != len(parsed.presents) {
				return farm{}, fmt.Errorf("%w: present %d is out of order", errInvalidFarm, shape)
			}
			if index+presentSize >= len(lines) {
				return farm{}, fmt.Errorf("%w: present %d is truncated", errInvalidFarm, shape)
			}
			parsedPresent, err := parsePresent(lines[index+1 : index+1+presentSize])
			if err != nil {
				return farm{}, err
			}
			parsed.presents = append(parsed.presents, parsedPresent)
			index += presentSize
		default:
			if len(parsed.presents) == 0 {
				return farm{}, fmt.Errorf("%w: region before any present", errInvalidFarm)
			}
			parsedRegion, err := parseRegion(line, len(parsed.presents))
			if err != nil {
				return farm{}, err
			}
			parsed.regions = append(parsed.regions, parsedRegion)
		}
	}
	if len(parsed.presents) == 0 {
		return farm{}, fmt.Errorf("%w: no presents", errInvalidFarm)
	}
	return parsed, nil
}

func parsePresent(rows []string) (present, error) {
	var cells pattern
	area := 0
	for row, raw := range rows {
		line := strings.TrimSpace(raw)
		if len(line) != presentSize {
			return present{}, fmt.Errorf("%w: present row %q needs %d cells", errInvalidFarm, line, presentSize)
		}
		for column := 0; column < len(line); column++ {
			switch line[column] {
			case '#':
				cells[row][column] = true
				area++
			case '.':
			default:
				return present{}, fmt.Errorf("%w: unknown present cell %q", errInvalidFarm, line[column])
			}
		}
	}
	if area == 0 {
		return present{}, fmt.Errorf("%w: present has no cells", errInvalidFarm)
	}
	return present{area: area, variants: variantsOf(cells)}, nil
}

func parseRegion(line string, presents int) (region, error) {
	size, counts, found := strings.Cut(line, ":")
	if !found {
		return region{}, fmt.Errorf("%w: region %q needs a colon", errInvalidFarm, line)
	}
	width, height, found := strings.Cut(strings.TrimSpace(size), "x")
	if !found {
		return region{}, fmt.Errorf("%w: region size %q needs an x", errInvalidFarm, size)
	}
	parsed := region{}
	var err error
	if parsed.width, err = parseSize(width); err != nil {
		return region{}, err
	}
	if parsed.height, err = parseSize(height); err != nil {
		return region{}, err
	}
	for _, field := range strings.Fields(counts) {
		count, err := strconv.Atoi(field)
		if err != nil || count < 0 {
			return region{}, fmt.Errorf("%w: present count %q", errInvalidFarm, field)
		}
		parsed.counts = append(parsed.counts, count)
	}
	if len(parsed.counts) != presents {
		return region{}, fmt.Errorf("%w: region %q lists %d of %d presents", errInvalidFarm, line, len(parsed.counts), presents)
	}
	return parsed, nil
}

func parseSize(field string) (int, error) {
	size, err := strconv.Atoi(field)
	if err != nil {
		return 0, fmt.Errorf("%w: region size %q: %w", errInvalidFarm, field, err)
	}
	if size <= 0 {
		return 0, fmt.Errorf("%w: region size %q must be positive", errInvalidFarm, field)
	}
	return size, nil
}

// variantsOf returns every rotation and flip of a pattern, without duplicates.
func variantsOf(cells pattern) []variant {
	var variants []variant
	oriented := cells
	for flip := range 2 {
		if flip == 1 {
			oriented = mirrored(oriented)
		}
		for range 4 {
			oriented = rotated(oriented)
			candidate := variantFrom(oriented)
			if !containsShape(variants, candidate.cells) {
				variants = append(variants, candidate)
			}
		}
	}
	return variants
}

func rotated(cells pattern) pattern {
	var result pattern
	for row := range presentSize {
		for column := range presentSize {
			result[row][column] = cells[presentSize-1-column][row]
		}
	}
	return result
}

func mirrored(cells pattern) pattern {
	var result pattern
	for row := range presentSize {
		for column := range presentSize {
			result[row][column] = cells[row][presentSize-1-column]
		}
	}
	return result
}

func variantFrom(cells pattern) variant {
	oriented := variant{}
	anchorRow, anchorColumn, found := 0, 0, false
	for row := range presentSize {
		for column := range presentSize {
			if !cells[row][column] {
				continue
			}
			if !found {
				anchorRow, anchorColumn, found = row, column, true
			}
			oriented.cells = append(oriented.cells, offset{row: row - anchorRow, column: column - anchorColumn})
			oriented.height = row - anchorRow + 1
		}
	}
	return oriented
}

func containsShape(variants []variant, cells []offset) bool {
	for _, existing := range variants {
		if len(existing.cells) != len(cells) {
			continue
		}
		same := true
		for index, cell := range existing.cells {
			same = same && cell == cells[index]
		}
		if same {
			return true
		}
	}
	return false
}

// regionFits decides whether one region can hold every present listed for it.
func regionFits(presents []present, tree region) bool {
	wanted, area, minHeight := 0, 0, presentSize
	for shape, count := range tree.counts {
		if count == 0 {
			continue
		}
		wanted += count
		area += count * presents[shape].area
		for _, oriented := range presents[shape].variants {
			minHeight = min(minHeight, oriented.height)
		}
	}
	if wanted == 0 {
		return true
	}
	cells := tree.width * tree.height
	if area > cells {
		return false
	}
	// Every present fits a 3x3 box, so this many need no puzzling at all.
	if wanted <= (tree.width/presentSize)*(tree.height/presentSize) {
		return true
	}
	grid := &packer{
		width:     tree.width,
		height:    tree.height,
		minHeight: minHeight,
		presents:  presents,
		counts:    append([]int(nil), tree.counts...),
		remaining: wanted,
		occupied:  make([]bool, cells),
	}
	return grid.fits(0, cells-area)
}

type packer struct {
	width     int
	height    int
	minHeight int
	presents  []present
	counts    []int
	remaining int
	occupied  []bool
}

// fits walks the grid in row-major order: the first undecided cell is either
// left empty, spending one cell of the waste budget, or anchors a present. A
// present covering that cell must start there, so every packing is tried
// exactly once and the budget keeps the search tight.
func (p *packer) fits(start, waste int) bool {
	if p.remaining == 0 {
		return true
	}
	for index := start; index < len(p.occupied); index++ {
		if p.occupied[index] {
			continue
		}
		row, column := index/p.width, index%p.width
		if row+p.minHeight > p.height {
			return false
		}
		for shape, count := range p.counts {
			if count == 0 {
				continue
			}
			for _, oriented := range p.presents[shape].variants {
				if !p.canPlace(oriented, row, column) {
					continue
				}
				p.place(oriented, row, column, true)
				p.counts[shape], p.remaining = count-1, p.remaining-1
				if p.fits(index+1, waste) {
					return true
				}
				p.counts[shape], p.remaining = count, p.remaining+1
				p.place(oriented, row, column, false)
			}
		}
		if waste == 0 {
			return false
		}
		waste--
	}
	return false
}

func (p *packer) canPlace(oriented variant, row, column int) bool {
	for _, cell := range oriented.cells {
		index, inside := p.cellIndex(cell, row, column)
		if !inside || p.occupied[index] {
			return false
		}
	}
	return true
}

func (p *packer) place(oriented variant, row, column int, occupied bool) {
	for _, cell := range oriented.cells {
		index, _ := p.cellIndex(cell, row, column)
		p.occupied[index] = occupied
	}
}

func (p *packer) cellIndex(cell offset, row, column int) (int, bool) {
	targetRow, targetColumn := row+cell.row, column+cell.column
	if targetRow >= p.height || targetColumn < 0 || targetColumn >= p.width {
		return 0, false
	}
	return targetRow*p.width + targetColumn, true
}
