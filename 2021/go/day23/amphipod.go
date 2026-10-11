package main

import (
	"container/heap"
	"errors"
	"fmt"
	"slices"
	"strings"
)

const (
	hallLen   = 11
	roomCount = 4
	maxDepth  = 4
)

// doors are the hallway cells in front of the rooms; amphipods pass through them
// but never stop there.
var doors = [roomCount]int{2, 4, 6, 8}

// foldedRows are the room rows part two reveals between the two rows of the diagram.
var foldedRows = [2][roomCount]amphipod{
	{desert, copper, bronze, amber},
	{desert, bronze, amber, copper},
}

var errInvalidBurrow = errors.New("invalid burrow")

type amphipod uint8

const (
	empty amphipod = iota
	amber
	bronze
	copper
	desert
)

var stepCosts = [...]int{amber: 1, bronze: 10, copper: 100, desert: 1000}

func parseAmphipod(cell byte) (amphipod, error) {
	if cell < 'A' || cell > 'D' {
		return empty, fmt.Errorf("%w: %q is not an amphipod", errInvalidBurrow, cell)
	}
	return amber + amphipod(cell-'A'), nil
}

func (a amphipod) String() string {
	if a < amber || a > desert {
		return "."
	}
	return string(rune('A' + a - amber))
}

func (a amphipod) home() int { return int(a - amber) }

func (a amphipod) stepCost() int { return stepCosts[a] }

// burrow is a comparable snapshot of the hallway and the first depth slots of each
// room. rooms[room][0] is the slot next to the hallway; rooms fill from the back.
type burrow struct {
	hall  [hallLen]amphipod
	rooms [roomCount][maxDepth]amphipod
	depth int
}

func parseBurrow(input string) (burrow, error) {
	var rows []string
	for _, row := range strings.Split(input, "\n") {
		if trimmed := strings.TrimSpace(row); trimmed != "" {
			rows = append(rows, trimmed)
		}
	}
	if len(rows) < 2 {
		return burrow{}, fmt.Errorf("%w: missing diagram", errInvalidBurrow)
	}
	if strings.Trim(rows[0], "#") != "" {
		return burrow{}, fmt.Errorf("%w: the diagram starts with a wall row", errInvalidBurrow)
	}
	if len(rows[1]) != hallLen+2 || strings.Count(rows[1], ".") != hallLen {
		return burrow{}, fmt.Errorf("%w: the hallway must be %d empty cells wide", errInvalidBurrow, hallLen)
	}
	var parsed burrow
	for _, row := range rows[2:] {
		cells := strings.ReplaceAll(row, "#", "")
		if cells == "" {
			break
		}
		if parsed.depth == maxDepth {
			return burrow{}, fmt.Errorf("%w: a burrow has at most %d room rows", errInvalidBurrow, maxDepth)
		}
		if len(cells) != roomCount {
			return burrow{}, fmt.Errorf("%w: a room row holds %d amphipods, found %d", errInvalidBurrow, roomCount, len(cells))
		}
		for room := range roomCount {
			resident, err := parseAmphipod(cells[room])
			if err != nil {
				return burrow{}, err
			}
			parsed.rooms[room][parsed.depth] = resident
		}
		parsed.depth++
	}
	if parsed.depth == 0 {
		return burrow{}, fmt.Errorf("%w: the burrow has no room rows", errInvalidBurrow)
	}
	for kind := amber; kind <= desert; kind++ {
		found := 0
		for _, slots := range parsed.rooms {
			for _, slot := range slots {
				if slot == kind {
					found++
				}
			}
		}
		if found != parsed.depth {
			return burrow{}, fmt.Errorf("%w: the burrow needs %d %s amphipods, found %d", errInvalidBurrow, parsed.depth, kind, found)
		}
	}
	return parsed, nil
}

func (b *burrow) unfold() (burrow, error) {
	if b.depth != 2 {
		return burrow{}, fmt.Errorf("%w: only a two-row burrow can be unfolded", errInvalidBurrow)
	}
	unfolded := *b
	unfolded.depth = maxDepth
	for room := range unfolded.rooms {
		unfolded.rooms[room] = [maxDepth]amphipod{
			b.rooms[room][0],
			foldedRows[0][room],
			foldedRows[1][room],
			b.rooms[room][1],
		}
	}
	return unfolded, nil
}

// topmost returns depth and empty when the room is empty.
func (b *burrow) topmost(room int) (int, amphipod) {
	for slot := range b.depth {
		if resident := b.rooms[room][slot]; resident != empty {
			return slot, resident
		}
	}
	return b.depth, empty
}

// settled reports whether room holds only its own kind, so nobody has to leave it
// and that kind is allowed to move in.
func (b *burrow) settled(room int) bool {
	for slot := range b.depth {
		if resident := b.rooms[room][slot]; resident != empty && resident.home() != room {
			return false
		}
	}
	return true
}

func (b *burrow) organized() bool {
	for room := range roomCount {
		if b.rooms[room][0] == empty || !b.settled(room) {
			return false
		}
	}
	return true
}

// hallFree excludes from, where the walking amphipod itself stands.
func (b *burrow) hallFree(from, to int) bool {
	for cell := min(from, to); cell <= max(from, to); cell++ {
		if cell != from && b.hall[cell] != empty {
			return false
		}
	}
	return true
}

// state pairs a burrow with an energy total: the cost of a move while moves are
// generated, the energy spent so far once the state is queued.
type state struct {
	energy int
	burrow burrow
}

// Walking an amphipod home from the hallway is always part of an optimal solution:
// the move costs the same whenever it happens, it frees a hallway cell and it fills
// a slot that amphipod has to occupy anyway. Taking that move on its own keeps the
// search small.
func (b *burrow) moves() []state {
	if home, ok := b.moveHome(); ok {
		return []state{home}
	}
	moves := make([]state, 0, roomCount*(hallLen-roomCount))
	for room, door := range doors {
		if b.settled(room) {
			continue
		}
		slot, resident := b.topmost(room)
		for cell := range hallLen {
			if slices.Contains(doors[:], cell) || !b.hallFree(door, cell) {
				continue
			}
			next := *b
			next.rooms[room][slot] = empty
			next.hall[cell] = resident
			steps := slot + 1 + abs(door-cell)
			moves = append(moves, state{resident.stepCost() * steps, next})
		}
	}
	return moves
}

// moveHome returns the first move that walks an amphipod from the hallway straight
// into its own room.
func (b *burrow) moveHome() (state, bool) {
	for cell, resident := range b.hall {
		if resident == empty {
			continue
		}
		room := resident.home()
		top, _ := b.topmost(room)
		if top == 0 || !b.settled(room) || !b.hallFree(cell, doors[room]) {
			continue
		}
		next := *b
		next.hall[cell] = empty
		next.rooms[room][top-1] = resident
		steps := abs(cell-doors[room]) + top
		return state{resident.stepCost() * steps, next}, true
	}
	return state{}, false
}

func leastEnergy(start burrow) (int, error) {
	best := map[burrow]int{start: 0}
	pending := &queue{{0, start}}
	for pending.Len() > 0 {
		current := heap.Pop(pending).(state)
		if current.burrow.organized() {
			return current.energy, nil
		}
		if shortest, seen := best[current.burrow]; seen && shortest < current.energy {
			continue
		}
		for _, move := range current.burrow.moves() {
			energy := current.energy + move.energy
			if shortest, seen := best[move.burrow]; !seen || energy < shortest {
				best[move.burrow] = energy
				heap.Push(pending, state{energy, move.burrow})
			}
		}
	}
	return 0, fmt.Errorf("%w: the amphipods cannot be organized", errInvalidBurrow)
}

func solve(input string) (int, int, error) {
	start, err := parseBurrow(input)
	if err != nil {
		return 0, 0, err
	}
	part1, err := leastEnergy(start)
	if err != nil {
		return 0, 0, err
	}
	unfolded, err := start.unfold()
	if err != nil {
		return 0, 0, err
	}
	part2, err := leastEnergy(unfolded)
	if err != nil {
		return 0, 0, err
	}
	return part1, part2, nil
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

type queue []state

func (q queue) Len() int           { return len(q) }
func (q queue) Less(i, j int) bool { return q[i].energy < q[j].energy }
func (q queue) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }

func (q *queue) Push(item any) { *q = append(*q, item.(state)) }

func (q *queue) Pop() any {
	last := len(*q) - 1
	popped := (*q)[last]
	*q = (*q)[:last]
	return popped
}
