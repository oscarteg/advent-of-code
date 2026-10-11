package main

import (
	"fmt"
	"strconv"
	"strings"
)

type move struct {
	count, from, to int
}

type drawing struct {
	stacks [][]byte
	moves  []move
}

// parseDrawing reads crate letters by fixed character column (index 1, 5,
// 9, ... every 4th column) rather than by splitting on whitespace, since
// empty stack slots leave gaps that whitespace splitting would misalign.
func parseDrawing(input string) (drawing, error) {
	stacksInput, movesInput, ok := strings.Cut(input, "\n\n")
	if !ok {
		return drawing{}, fmt.Errorf("parse drawing: missing blank line separating stacks from moves")
	}

	lines := strings.Split(stacksInput, "\n")
	if len(lines) == 0 {
		return drawing{}, fmt.Errorf("parse drawing: missing stack label row")
	}
	labels := lines[len(lines)-1]
	lines = lines[:len(lines)-1]
	numStacks := len(strings.Fields(labels))
	if numStacks == 0 {
		return drawing{}, fmt.Errorf("parse drawing: no stack labels found")
	}

	stacks := make([][]byte, numStacks)
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		for stackIndex := range stacks {
			column := 1 + stackIndex*4
			if column < len(line) && line[column] != ' ' {
				stacks[stackIndex] = append(stacks[stackIndex], line[column])
			}
		}
	}

	var moves []move
	for _, line := range strings.Split(movesInput, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		mv, err := parseMove(line)
		if err != nil {
			return drawing{}, err
		}
		moves = append(moves, mv)
	}

	return drawing{stacks: stacks, moves: moves}, nil
}

func parseMove(line string) (move, error) {
	parts := strings.Fields(line)
	if len(parts) != 6 || parts[0] != "move" || parts[2] != "from" || parts[4] != "to" {
		return move{}, fmt.Errorf("parse move: malformed line %q", line)
	}
	count, err := strconv.Atoi(parts[1])
	if err != nil {
		return move{}, fmt.Errorf("parse move: invalid count in %q: %w", line, err)
	}
	from, err := strconv.Atoi(parts[3])
	if err != nil {
		return move{}, fmt.Errorf("parse move: invalid source stack in %q: %w", line, err)
	}
	to, err := strconv.Atoi(parts[5])
	if err != nil {
		return move{}, fmt.Errorf("parse move: invalid destination stack in %q: %w", line, err)
	}
	if from < 1 || to < 1 {
		return move{}, fmt.Errorf("parse move: stack numbers are 1-indexed in %q", line)
	}
	return move{count: count, from: from - 1, to: to - 1}, nil
}

func cloneStacks(stacks [][]byte) [][]byte {
	clone := make([][]byte, len(stacks))
	for i, stack := range stacks {
		clone[i] = append([]byte(nil), stack...)
	}
	return clone
}

func applyMoves(stacks [][]byte, moves []move, preserveOrder bool) ([][]byte, error) {
	for _, mv := range moves {
		if mv.from < 0 || mv.from >= len(stacks) || mv.to < 0 || mv.to >= len(stacks) {
			return nil, fmt.Errorf("apply move: unknown stack in %+v", mv)
		}
		from := stacks[mv.from]
		if mv.count > len(from) {
			return nil, fmt.Errorf("apply move: not enough crates to move in %+v", mv)
		}
		splitAt := len(from) - mv.count
		moved := append([]byte(nil), from[splitAt:]...)
		stacks[mv.from] = from[:splitAt]
		if !preserveOrder {
			for i, j := 0, len(moved)-1; i < j; i, j = i+1, j-1 {
				moved[i], moved[j] = moved[j], moved[i]
			}
		}
		stacks[mv.to] = append(stacks[mv.to], moved...)
	}
	return stacks, nil
}

func topCrates(stacks [][]byte) string {
	top := make([]byte, 0, len(stacks))
	for _, stack := range stacks {
		if len(stack) > 0 {
			top = append(top, stack[len(stack)-1])
		}
	}
	return string(top)
}

func solve(input string) (part1, part2 string, err error) {
	d, err := parseDrawing(input)
	if err != nil {
		return "", "", err
	}

	part1Stacks, err := applyMoves(cloneStacks(d.stacks), d.moves, false)
	if err != nil {
		return "", "", err
	}
	part2Stacks, err := applyMoves(d.stacks, d.moves, true)
	if err != nil {
		return "", "", err
	}

	return topCrates(part1Stacks), topCrates(part2Stacks), nil
}
