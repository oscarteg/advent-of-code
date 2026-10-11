package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	startValve          = "AA"
	minutesAlone        = 30
	minutesWithElephant = 26
	// maxOpenableValves bounds the bitmask table below; real scans stay far under it.
	maxOpenableValves = 20
	// unreachableMinutes is the travel time standing in for "no tunnel route at all".
	unreachableMinutes = 1 << 20
)

var (
	errNoValves     = errors.New("input contains no valves")
	errMissingStart = errors.New("scan has no " + startValve + " valve to start from")
)

// scannedValve is one line of the scan: a valve, its flow rate and the valves
// its tunnels reach.
type scannedValve struct {
	name    string
	flow    int
	tunnels []string
}

// cutTunnelPrefix splits the valve list off its prefix, which is singular for a
// lone tunnel.
func cutTunnelPrefix(text string) ([]string, bool) {
	list, ok := strings.CutPrefix(text, "tunnels lead to valves ")
	if !ok {
		list, ok = strings.CutPrefix(text, "tunnel leads to valve ")
	}
	if !ok {
		return nil, false
	}
	tunnels := strings.Split(list, ", ")
	for _, tunnel := range tunnels {
		if tunnel == "" {
			return nil, false
		}
	}
	return tunnels, true
}

func parseValve(line string) (scannedValve, error) {
	rest, ok := strings.CutPrefix(line, "Valve ")
	if !ok {
		return scannedValve{}, fmt.Errorf("malformed valve %q", line)
	}
	name, rest, ok := strings.Cut(rest, " has flow rate=")
	if !ok || name == "" {
		return scannedValve{}, fmt.Errorf("missing flow rate in %q", line)
	}
	flowText, tunnelText, ok := strings.Cut(rest, "; ")
	if !ok {
		return scannedValve{}, fmt.Errorf("missing tunnel list in %q", line)
	}
	flow, err := strconv.Atoi(flowText)
	if err != nil {
		return scannedValve{}, fmt.Errorf("invalid flow rate in %q: %w", line, err)
	}
	tunnels, ok := cutTunnelPrefix(tunnelText)
	if !ok {
		return scannedValve{}, fmt.Errorf("malformed tunnel list in %q", line)
	}
	return scannedValve{name: name, flow: flow, tunnels: tunnels}, nil
}

func parse(input string) ([]scannedValve, error) {
	var valves []scannedValve
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		valve, err := parseValve(line)
		if err != nil {
			return nil, err
		}
		valves = append(valves, valve)
	}
	if len(valves) == 0 {
		return nil, errNoValves
	}
	return valves, nil
}

// shortestWalkingTimes returns the walking time between every pair of valves, by
// Floyd-Warshall.
func shortestWalkingTimes(valves []scannedValve, index map[string]int) ([][]int, error) {
	minutes := make([][]int, len(valves))
	for from := range minutes {
		minutes[from] = make([]int, len(valves))
		for to := range minutes[from] {
			minutes[from][to] = unreachableMinutes
		}
	}
	for from, valve := range valves {
		for _, tunnel := range valve.tunnels {
			to, ok := index[tunnel]
			if !ok {
				return nil, fmt.Errorf("tunnel from %s leads to unknown valve %s", valve.name, tunnel)
			}
			minutes[from][to] = 1
		}
		minutes[from][from] = 0
	}
	for via := range valves {
		for from := range valves {
			hop := minutes[from][via]
			if hop >= unreachableMinutes {
				continue
			}
			for to, viaTo := range minutes[via] {
				if detour := hop + viaTo; detour < minutes[from][to] {
					minutes[from][to] = detour
				}
			}
		}
	}
	return minutes, nil
}

// network holds the valves worth opening and the walking times the search needs;
// every index is a bit position in the sets of opened valves.
type network struct {
	flows     []int
	travel    [][]int
	fromStart []int
}

func buildNetwork(valves []scannedValve) (network, error) {
	index := make(map[string]int, len(valves))
	for position, valve := range valves {
		if _, duplicate := index[valve.name]; duplicate {
			return network{}, fmt.Errorf("valve %s is listed twice", valve.name)
		}
		index[valve.name] = position
	}
	start, ok := index[startValve]
	if !ok {
		return network{}, errMissingStart
	}
	minutes, err := shortestWalkingTimes(valves, index)
	if err != nil {
		return network{}, err
	}

	var openable []int
	for position, valve := range valves {
		if valve.flow > 0 {
			openable = append(openable, position)
		}
	}
	if len(openable) > maxOpenableValves {
		return network{}, fmt.Errorf("%d valves with flow is more than this solver handles (%d)",
			len(openable), maxOpenableValves)
	}

	net := network{
		flows:     make([]int, len(openable)),
		travel:    make([][]int, len(openable)),
		fromStart: make([]int, len(openable)),
	}
	for from, position := range openable {
		net.flows[from] = valves[position].flow
		net.fromStart[from] = minutes[start][position]
		net.travel[from] = make([]int, len(openable))
		for to, other := range openable {
			net.travel[from][to] = minutes[position][other]
		}
	}
	return net, nil
}

// explore walks depth-first over the valves still closed, recording for every set
// of opened valves the most pressure that set can release.
func explore(net network, best []int, valve, minutesLeft, opened, released int) {
	best[opened] = max(best[opened], released)
	for next, flow := range net.flows {
		bit := 1 << next
		if opened&bit != 0 {
			continue
		}
		cost := net.travel[valve][next] + 1
		if cost >= minutesLeft {
			continue
		}
		left := minutesLeft - cost
		explore(net, best, next, left, opened|bit, released+left*flow)
	}
}

// bestPerOpenedSet returns the most pressure released per set of opened valves
// within the minutes given.
func bestPerOpenedSet(net network, minutes int) []int {
	best := make([]int, 1<<len(net.flows))
	for valve, flow := range net.flows {
		cost := net.fromStart[valve] + 1
		if cost >= minutes {
			continue
		}
		left := minutes - cost
		explore(net, best, valve, left, 1<<valve, left*flow)
	}
	return best
}

// maximizeOverSubsets turns best[set] into the best over every subset of set, so
// two disjoint routes can be combined by looking up a set and its complement.
func maximizeOverSubsets(best []int, valveCount int) {
	for position := 0; position < valveCount; position++ {
		bit := 1 << position
		for set := range best {
			if set&bit != 0 {
				best[set] = max(best[set], best[set^bit])
			}
		}
	}
}

func solve(input string) (part1, part2 int, err error) {
	valves, err := parse(input)
	if err != nil {
		return 0, 0, err
	}
	net, err := buildNetwork(valves)
	if err != nil {
		return 0, 0, err
	}

	for _, released := range bestPerOpenedSet(net, minutesAlone) {
		part1 = max(part1, released)
	}

	shared := bestPerOpenedSet(net, minutesWithElephant)
	maximizeOverSubsets(shared, len(net.flows))
	everyValve := len(shared) - 1
	for mine, released := range shared {
		part2 = max(part2, released+shared[everyValve^mine])
	}
	return part1, part2, nil
}
