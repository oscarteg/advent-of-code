package main

import (
	"errors"
	"fmt"
	"strings"
)

const (
	chamberWidth = 7
	part1Rocks   = 2022
	part2Rocks   = 1_000_000_000_000
	profileDepth = 50
)

var errInvalidJetPattern = errors.New("invalid jet pattern")

// rock describes a shape as bottom-to-top row bitmasks (bit 0 is the leftmost column) and its width.
type rock struct {
	rows  []uint8
	width int
}

var rocks = [5]rock{
	{rows: []uint8{0b1111}, width: 4},
	{rows: []uint8{0b010, 0b111, 0b010}, width: 3},
	{rows: []uint8{0b111, 0b100, 0b100}, width: 3},
	{rows: []uint8{0b1, 0b1, 0b1, 0b1}, width: 1},
	{rows: []uint8{0b11, 0b11}, width: 2},
}

func parseJets(input string) ([]int, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, fmt.Errorf("%w: pattern is empty", errInvalidJetPattern)
	}
	jets := make([]int, len(trimmed))
	for i, c := range trimmed {
		switch c {
		case '<':
			jets[i] = -1
		case '>':
			jets[i] = 1
		default:
			return nil, fmt.Errorf("%w: invalid jet character %q", errInvalidJetPattern, c)
		}
	}
	return jets, nil
}

func collides(shape rock, x, y int, rows []uint8) bool {
	if x < 0 || x+shape.width > chamberWidth || y < 0 {
		return true
	}
	for i, rowMask := range shape.rows {
		rowY := y + i
		if rowY < len(rows) && uint16(rows[rowY])&(uint16(rowMask)<<x) != 0 {
			return true
		}
	}
	return false
}

func settle(shape rock, x, y int, rows []uint8) []uint8 {
	for i, rowMask := range shape.rows {
		rowY := y + i
		for len(rows) <= rowY {
			rows = append(rows, 0)
		}
		rows[rowY] |= rowMask << uint(x)
	}
	return rows
}

// surfaceProfile returns the top profileDepth rows of the tower, right-aligned and padded with
// the sentinel 0xFF (never a valid 7-bit row mask) so short towers still produce a stable key.
func surfaceProfile(rows []uint8) [profileDepth]uint8 {
	var profile [profileDepth]uint8
	for i := range profile {
		profile[i] = 0xFF
	}
	depth := len(rows)
	if depth > profileDepth {
		depth = profileDepth
	}
	copy(profile[profileDepth-depth:], rows[len(rows)-depth:])
	return profile
}

type profileKey struct {
	rockIdx int
	jetIdx  int
	profile [profileDepth]uint8
}

type seenEntry struct {
	rockCount uint64
	height    uint64
}

// towerHeight returns the tower height after dropping target rocks. It detects a cycle on
// (rock shape index, jet index, surface profile) so that huge target values stay fast: once a
// key repeats, the gap between occurrences gives a repeating (rock count, height) cycle, and the
// remaining rocks are resolved as full cycles plus a small simulated remainder.
func towerHeight(jets []int, target uint64) uint64 {
	var rows []uint8
	heights := []uint64{0}
	seen := make(map[profileKey]seenEntry)

	jetIndex := 0
	var n uint64
	for n < target {
		rockIdx := int(n % 5)
		shape := rocks[rockIdx]
		x, y := 2, len(rows)+3
		for {
			dir := jets[jetIndex]
			jetIndex = (jetIndex + 1) % len(jets)
			if nx := x + dir; !collides(shape, nx, y, rows) {
				x = nx
			}
			if !collides(shape, x, y-1, rows) {
				y--
			} else {
				rows = settle(shape, x, y, rows)
				break
			}
		}
		n++
		heights = append(heights, uint64(len(rows)))

		key := profileKey{rockIdx: rockIdx, jetIdx: jetIndex, profile: surfaceProfile(rows)}
		if prev, ok := seen[key]; ok {
			cycleLen := n - prev.rockCount
			cycleGain := uint64(len(rows)) - prev.height
			remaining := target - n
			cycles := remaining / cycleLen
			remainder := remaining % cycleLen
			remainderGain := heights[prev.rockCount+remainder] - prev.height
			return uint64(len(rows)) + cycles*cycleGain + remainderGain
		}
		seen[key] = seenEntry{rockCount: n, height: uint64(len(rows))}
	}
	return uint64(len(rows))
}

func solve(input string) (uint64, uint64, error) {
	jets, err := parseJets(input)
	if err != nil {
		return 0, 0, err
	}
	part1 := towerHeight(jets, part1Rocks)
	part2 := towerHeight(jets, part2Rocks)
	return part1, part2, nil
}
