package main

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidSchematic = errors.New("invalid schematic")

const (
	rows     = 7
	columns  = 5
	pinSpace = rows - 2
)

type heights [columns]int

func parseSchematic(block string) (heights heights, isLock bool, err error) {
	lines := strings.Split(block, "\n")
	if len(lines) != rows {
		return heights, false, fmt.Errorf("%w: schematic must have %d rows", errInvalidSchematic, rows)
	}
	for _, line := range lines {
		if len(line) != columns {
			return heights, false, fmt.Errorf("%w: row must have %d columns", errInvalidSchematic, columns)
		}
	}

	topFull := !strings.Contains(lines[0], ".")
	bottomFull := !strings.Contains(lines[rows-1], ".")
	if topFull == bottomFull {
		return heights, false, fmt.Errorf("%w: schematic is neither lock nor key", errInvalidSchematic)
	}

	for _, line := range lines {
		for column, cell := range line {
			if cell == '#' {
				heights[column]++
			}
		}
	}
	for column := range heights {
		heights[column]--
	}
	return heights, topFull, nil
}

func fits(lock, key heights) bool {
	for column := 0; column < columns; column++ {
		if lock[column]+key[column] > pinSpace {
			return false
		}
	}
	return true
}

func solve(input string) (int, error) {
	blocks := strings.Split(strings.Trim(strings.ReplaceAll(input, "\r\n", "\n"), "\n"), "\n\n")

	var locks, keys []heights
	for _, block := range blocks {
		h, isLock, err := parseSchematic(block)
		if err != nil {
			return 0, err
		}
		if isLock {
			locks = append(locks, h)
		} else {
			keys = append(keys, h)
		}
	}
	if len(locks) == 0 || len(keys) == 0 {
		return 0, fmt.Errorf("%w: input needs locks and keys", errInvalidSchematic)
	}

	fittingPairs := 0
	for _, lock := range locks {
		for _, key := range keys {
			if fits(lock, key) {
				fittingPairs++
			}
		}
	}
	return fittingPairs, nil
}
