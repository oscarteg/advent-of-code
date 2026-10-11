package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// Directional-keypad robots between the door's numeric keypad and the human.
const (
	part1Robots = 2
	part2Robots = 25
)

var errInvalidCode = errors.New("invalid door code")

type pad int

const (
	numericPad pad = iota
	directionalPad
)

// Keypad rows, with ' ' marking the gap no robot arm may cross.
var padRows = [...][]string{
	numericPad:     {"789", "456", "123", " 0A"},
	directionalPad: {" ^A", "<v>"},
}

type position struct {
	row, column int
}

func (p pad) locate(key byte) (position, error) {
	for row, buttons := range padRows[p] {
		if column := strings.IndexByte(buttons, key); column >= 0 {
			return position{row: row, column: column}, nil
		}
	}
	return position{}, fmt.Errorf("%w: key %q is not on the keypad", errInvalidCode, key)
}

func (p pad) hasKey(at position) bool {
	rows := padRows[p]
	if at.row < 0 || at.row >= len(rows) {
		return false
	}
	buttons := rows[at.row]
	return at.column >= 0 && at.column < len(buttons) && buttons[at.column] != ' '
}

// shortestRoutes returns the shortest routes for one move, each ending with the A that presses
// the key. Only two orders can be shortest - every horizontal step first, or every vertical
// step first - and an order is impossible when its turning point lands on the keypad's gap.
func (p pad) shortestRoutes(from, to byte) ([]string, error) {
	start, err := p.locate(from)
	if err != nil {
		return nil, err
	}
	target, err := p.locate(to)
	if err != nil {
		return nil, err
	}

	horizontalButton, verticalButton := byte('<'), byte('^')
	if target.column > start.column {
		horizontalButton = '>'
	}
	if target.row > start.row {
		verticalButton = 'v'
	}
	horizontal := strings.Repeat(string(horizontalButton), abs(target.column-start.column))
	vertical := strings.Repeat(string(verticalButton), abs(target.row-start.row))

	var routes []string
	if p.hasKey(position{row: start.row, column: target.column}) {
		routes = append(routes, horizontal+vertical+"A")
	}
	if horizontal != "" && vertical != "" && p.hasKey(position{row: target.row, column: start.column}) {
		routes = append(routes, vertical+horizontal+"A")
	}
	if len(routes) == 0 {
		return nil, fmt.Errorf("%w: the gap blocks every route from %q to %q", errInvalidCode, from, to)
	}
	return routes, nil
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}

type moveKey struct {
	pad      pad
	from, to byte
	robots   int
}

// solver counts human button presses, memoised per keypad, key pair and robot count.
type solver struct {
	cache map[moveKey]int64
}

func newSolver() *solver {
	return &solver{cache: make(map[moveKey]int64)}
}

// typeOn returns the presses the human needs to have sequence typed on p.
func (s *solver) typeOn(p pad, sequence string, robots int) (int64, error) {
	var presses int64
	from := byte('A') // every arm starts on its own A button
	for i := range len(sequence) {
		move, err := s.moveAndPress(p, from, sequence[i], robots)
		if err != nil {
			return 0, err
		}
		presses += move
		from = sequence[i]
	}
	return presses, nil
}

func (s *solver) moveAndPress(p pad, from, to byte, robots int) (int64, error) {
	key := moveKey{pad: p, from: from, to: to, robots: robots}
	if cached, ok := s.cache[key]; ok {
		return cached, nil
	}

	routes, err := p.shortestRoutes(from, to)
	if err != nil {
		return 0, err
	}

	cheapest := int64(math.MaxInt64)
	for _, route := range routes {
		// The human types the last keypad itself, one press per button of the route.
		presses := int64(len(route))
		if robots > 0 {
			if presses, err = s.typeOn(directionalPad, route, robots-1); err != nil {
				return 0, err
			}
		}
		cheapest = min(cheapest, presses)
	}

	s.cache[key] = cheapest
	return cheapest, nil
}

type code struct {
	buttons     string
	numericPart int64
}

func parseCode(line string) (code, error) {
	digits, ok := strings.CutSuffix(line, "A")
	if !ok {
		return code{}, fmt.Errorf("%w: %q must be digits followed by A", errInvalidCode, line)
	}
	numericPart, err := strconv.ParseUint(digits, 10, 63)
	if err != nil {
		return code{}, fmt.Errorf("%w: %q has no numeric part: %w", errInvalidCode, line, err)
	}
	return code{buttons: line, numericPart: int64(numericPart)}, nil
}

func parseCodes(input string) ([]code, error) {
	var codes []code
	for _, raw := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		parsed, err := parseCode(line)
		if err != nil {
			return nil, err
		}
		codes = append(codes, parsed)
	}
	if len(codes) == 0 {
		return nil, fmt.Errorf("%w: input has no door codes", errInvalidCode)
	}
	return codes, nil
}

// solve sums the complexity of every code: the presses the human needs for the shortest
// sequence that types it, times the code's numeric part.
func solve(input string) (int64, int64, error) {
	codes, err := parseCodes(input)
	if err != nil {
		return 0, 0, err
	}

	keypads := newSolver()
	var part1, part2 int64
	for _, c := range codes {
		withTwoRobots, err := keypads.typeOn(numericPad, c.buttons, part1Robots)
		if err != nil {
			return 0, 0, err
		}
		withAllRobots, err := keypads.typeOn(numericPad, c.buttons, part2Robots)
		if err != nil {
			return 0, 0, err
		}
		part1 += withTwoRobots * c.numericPart
		part2 += withAllRobots * c.numericPart
	}
	return part1, part2, nil
}
