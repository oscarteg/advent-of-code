package main

import (
	"errors"
	"fmt"
	"strings"
)

var errInvalidNetwork = errors.New("invalid module configuration")

// errNoImplicitSink means the graph has no module referenced only as a
// destination, so part 2's rx-style target does not apply to this input.
var errNoImplicitSink = errors.New("no undefined sink module found")

type moduleKind int

const (
	broadcasterKind moduleKind = iota
	flipFlopKind
	conjunctionKind
)

type module struct {
	kind         moduleKind
	destinations []string
	flipFlopOn   bool
	memory       map[string]bool
}

// network is a graph of pulse-propagation modules, keyed by name.
type network map[string]*module

func parseNetwork(input string) (network, error) {
	net := network{}
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(input), "\r\n", "\n"), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		namePart, destPart, ok := strings.Cut(line, "->")
		if !ok {
			return nil, fmt.Errorf("%w: module line is missing '->'", errInvalidNetwork)
		}
		namePart = strings.TrimSpace(namePart)
		kind := broadcasterKind
		name := namePart
		switch {
		case strings.HasPrefix(namePart, "%"):
			kind, name = flipFlopKind, namePart[1:]
		case strings.HasPrefix(namePart, "&"):
			kind, name = conjunctionKind, namePart[1:]
		}
		if name == "" {
			return nil, fmt.Errorf("%w: module is missing a name", errInvalidNetwork)
		}
		var destinations []string
		for _, dest := range strings.Split(destPart, ",") {
			dest = strings.TrimSpace(dest)
			if dest != "" {
				destinations = append(destinations, dest)
			}
		}
		if len(destinations) == 0 {
			return nil, fmt.Errorf("%w: module has no destinations", errInvalidNetwork)
		}
		net[name] = &module{kind: kind, destinations: destinations, memory: map[string]bool{}}
	}
	if len(net) == 0 {
		return nil, fmt.Errorf("%w: module configuration is empty", errInvalidNetwork)
	}

	for name, m := range net {
		for _, dest := range m.destinations {
			if target, ok := net[dest]; ok && target.kind == conjunctionKind {
				target.memory[name] = false
			}
		}
	}
	return net, nil
}

// pulse is one delivered signal: source sent a high (true) or low pulse to target.
type pulse struct {
	source, target string
	high           bool
}

// pressButton pushes the button once and returns every pulse sent, in
// delivery order, including the initial low pulse the button sends to
// broadcaster.
func (net network) pressButton() []pulse {
	var log []pulse
	queue := []pulse{{source: "button", target: "broadcaster", high: false}}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		log = append(log, p)

		m, ok := net[p.target]
		if !ok {
			continue
		}
		var send bool
		var fire bool
		switch m.kind {
		case broadcasterKind:
			send, fire = p.high, true
		case flipFlopKind:
			if !p.high {
				m.flipFlopOn = !m.flipFlopOn
				send, fire = m.flipFlopOn, true
			}
		case conjunctionKind:
			m.memory[p.source] = p.high
			allHigh := true
			for _, remembered := range m.memory {
				allHigh = allHigh && remembered
			}
			send, fire = !allHigh, true
		}
		if fire {
			for _, dest := range m.destinations {
				queue = append(queue, pulse{source: p.target, target: dest, high: send})
			}
		}
	}
	return log
}

// inputsOf returns the names of every module that lists target as a destination.
func (net network) inputsOf(target string) []string {
	var inputs []string
	for name, m := range net {
		for _, dest := range m.destinations {
			if dest == target {
				inputs = append(inputs, name)
				break
			}
		}
	}
	return inputs
}

// implicitSink returns the single module name that is referenced as a
// destination but never defined with its own line (e.g. rx in the real
// puzzle input, or output in the official examples).
func (net network) implicitSink() (string, error) {
	seen := map[string]bool{}
	for _, m := range net {
		for _, dest := range m.destinations {
			if _, defined := net[dest]; !defined {
				seen[dest] = true
			}
		}
	}
	switch len(seen) {
	case 1:
		for name := range seen {
			return name, nil
		}
	case 0:
		return "", errNoImplicitSink
	}
	return "", fmt.Errorf("%w: more than one undefined sink module", errInvalidNetwork)
}
