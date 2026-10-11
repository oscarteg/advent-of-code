package main

import (
	"errors"
	"fmt"
	"strings"
)

type graph map[string][]string

var errEmptyInput = errors.New("empty input")

func parseGraph(input string) (graph, error) {
	g := graph{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		device, outputs, found := strings.Cut(line, ":")
		device = strings.TrimSpace(device)
		if !found || device == "" {
			return nil, fmt.Errorf("invalid line %q: want \"device: outputs\"", line)
		}
		g[device] = strings.Fields(outputs)
	}
	if len(g) == 0 {
		return nil, errEmptyInput
	}
	return g, nil
}

func (g graph) countPaths(from, to string, memo map[string]int) int {
	if from == to {
		return 1
	}
	if known, ok := memo[from]; ok {
		return known
	}
	total := 0
	for _, next := range g[from] {
		total += g.countPaths(next, to, memo)
	}
	memo[from] = total
	return total
}

func (g graph) pathsVia(waypoints ...string) int {
	total := 1
	for i := 1; i < len(waypoints); i++ {
		total *= g.countPaths(waypoints[i-1], waypoints[i], map[string]int{})
	}
	return total
}

func solve(input string) (part1, part2 int, err error) {
	g, err := parseGraph(input)
	if err != nil {
		return 0, 0, err
	}
	// The graph is acyclic, so at most one of the two waypoint orders has any paths.
	part2 = g.pathsVia("svr", "dac", "fft", "out") + g.pathsVia("svr", "fft", "dac", "out")
	return g.pathsVia("you", "out"), part2, nil
}
