package main

import (
	"fmt"
	"strings"
)

type Pos struct {
	x, y int
}

type State struct {
	elves        map[Pos]bool
	directionIdx int
}

func newState() *State {
	return &State{
		elves:        make(map[Pos]bool),
		directionIdx: 0,
	}
}

func (s *State) hasNeighbors(pos Pos) bool {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if s.elves[Pos{pos.x + dx, pos.y + dy}] {
				return true
			}
		}
	}
	return false
}

func (s *State) canMoveNorth(pos Pos) bool {
	return !s.elves[Pos{pos.x - 1, pos.y - 1}] &&
		!s.elves[Pos{pos.x, pos.y - 1}] &&
		!s.elves[Pos{pos.x + 1, pos.y - 1}]
}

func (s *State) canMoveSouth(pos Pos) bool {
	return !s.elves[Pos{pos.x - 1, pos.y + 1}] &&
		!s.elves[Pos{pos.x, pos.y + 1}] &&
		!s.elves[Pos{pos.x + 1, pos.y + 1}]
}

func (s *State) canMoveWest(pos Pos) bool {
	return !s.elves[Pos{pos.x - 1, pos.y - 1}] &&
		!s.elves[Pos{pos.x - 1, pos.y}] &&
		!s.elves[Pos{pos.x - 1, pos.y + 1}]
}

func (s *State) canMoveEast(pos Pos) bool {
	return !s.elves[Pos{pos.x + 1, pos.y - 1}] &&
		!s.elves[Pos{pos.x + 1, pos.y}] &&
		!s.elves[Pos{pos.x + 1, pos.y + 1}]
}

func (s *State) proposeMove(pos Pos) *Pos {
	if !s.hasNeighbors(pos) {
		return nil
	}

	for i := 0; i < 4; i++ {
		dir := (s.directionIdx + i) % 4
		switch dir {
		case 0:
			if s.canMoveNorth(pos) {
				newPos := Pos{pos.x, pos.y - 1}
				return &newPos
			}
		case 1:
			if s.canMoveSouth(pos) {
				newPos := Pos{pos.x, pos.y + 1}
				return &newPos
			}
		case 2:
			if s.canMoveWest(pos) {
				newPos := Pos{pos.x - 1, pos.y}
				return &newPos
			}
		case 3:
			if s.canMoveEast(pos) {
				newPos := Pos{pos.x + 1, pos.y}
				return &newPos
			}
		}
	}
	return nil
}

func (s *State) round() bool {
	proposals := make(map[Pos][]Pos)

	for elf := range s.elves {
		if newPos := s.proposeMove(elf); newPos != nil {
			proposals[*newPos] = append(proposals[*newPos], elf)
		}
	}

	moved := false
	for newPos, elves := range proposals {
		if len(elves) == 1 {
			oldPos := elves[0]
			delete(s.elves, oldPos)
			s.elves[newPos] = true
			moved = true
		}
	}

	s.directionIdx = (s.directionIdx + 1) % 4
	return moved
}

func (s *State) countEmptyTiles() int {
	if len(s.elves) == 0 {
		return 0
	}

	minX, maxX := 1000000, -1000000
	minY, maxY := 1000000, -1000000

	for elf := range s.elves {
		if elf.x < minX {
			minX = elf.x
		}
		if elf.x > maxX {
			maxX = elf.x
		}
		if elf.y < minY {
			minY = elf.y
		}
		if elf.y > maxY {
			maxY = elf.y
		}
	}

	width := maxX - minX + 1
	height := maxY - minY + 1
	total := width * height
	return total - len(s.elves)
}

func parseInput(input string) (*State, error) {
	state := newState()
	lines := strings.Split(strings.TrimSpace(input), "\n")

	for y, line := range lines {
		if line == "" {
			continue
		}
		for x, ch := range line {
			switch ch {
			case '#':
				state.elves[Pos{x, y}] = true
			case '.':
			default:
				return nil, fmt.Errorf("invalid character: %c", ch)
			}
		}
	}

	if len(state.elves) == 0 {
		return nil, fmt.Errorf("no elves found")
	}

	return state, nil
}

func solve(input string) (int, int, error) {
	state, err := parseInput(input)
	if err != nil {
		return 0, 0, fmt.Errorf("parse error: %w", err)
	}

	part2 := 0
	for roundNum := 1; roundNum <= 10; roundNum++ {
		if !state.round() {
			part2 = roundNum
			break
		}
	}

	part1 := state.countEmptyTiles()

	if part2 == 0 {
		roundNum := 10
		for {
			roundNum++
			if !state.round() {
				part2 = roundNum
				break
			}
		}
	}

	return part1, part2, nil
}
