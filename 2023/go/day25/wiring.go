package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

var errInvalidDiagram = errors.New("invalid wiring diagram")

// diagram holds components connected by undirected wires as paired
// unit-capacity arcs, so that a maximum flow run yields a minimum cut. Arc i
// and arc i^1 are the two directions of the same wire; outgoing[c] holds the
// arcs leaving component c and heads[a] the component arc a enters.
type diagram struct {
	outgoing [][]int
	heads    []int
}

func parseDiagram(input string) (*diagram, error) {
	d := &diagram{}
	ids := make(map[string]int)
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		name, connections, found := strings.Cut(line, ":")
		if !found {
			return nil, fmt.Errorf("%w: line %q has no ':'", errInvalidDiagram, line)
		}
		component, err := d.componentID(ids, strings.TrimSpace(name))
		if err != nil {
			return nil, err
		}
		neighbours := strings.Fields(connections)
		if len(neighbours) == 0 {
			return nil, fmt.Errorf("%w: component %q has no connections", errInvalidDiagram, name)
		}
		for _, neighbour := range neighbours {
			id, err := d.componentID(ids, neighbour)
			if err != nil {
				return nil, err
			}
			if id == component {
				return nil, fmt.Errorf("%w: component %q is wired to itself", errInvalidDiagram, neighbour)
			}
			d.connect(component, id)
		}
	}
	if d.components() < 2 {
		return nil, fmt.Errorf("%w: fewer than two components", errInvalidDiagram)
	}
	return d, nil
}

func (d *diagram) componentID(ids map[string]int, name string) (int, error) {
	if name == "" {
		return 0, fmt.Errorf("%w: component without a name", errInvalidDiagram)
	}
	id, known := ids[name]
	if !known {
		id = d.components()
		ids[name] = id
		d.outgoing = append(d.outgoing, nil)
	}
	return id, nil
}

func (d *diagram) connect(left, right int) {
	d.outgoing[left] = append(d.outgoing[left], len(d.heads))
	d.heads = append(d.heads, right)
	d.outgoing[right] = append(d.outgoing[right], len(d.heads))
	d.heads = append(d.heads, left)
}

func (d *diagram) components() int {
	return len(d.outgoing)
}

// minimumCut returns the sizes of the two groups left by the cheapest set of
// wires that splits the diagram. Keeping the cheapest maximum flow from one
// fixed component to every other one yields the global minimum cut, because the
// fixed component always lies on one of its two sides.
func (d *diagram) minimumCut() (int, int) {
	cheapest, sourceSide := math.MaxInt, 1
	for sink := 1; sink < d.components(); sink++ {
		if wires, reachable, under := d.flowUnder(sink, cheapest); under {
			cheapest, sourceSide = wires, reachable
		}
	}
	return sourceSide, d.components() - sourceSide
}

// flowUnder pushes unit flow from component 0 to sink along augmenting paths and
// reports the flow plus the number of components on the source side of the
// resulting cut. It gives up, reporting false, once the flow can no longer stay
// below limit.
func (d *diagram) flowUnder(sink, limit int) (wires, sourceSide int, under bool) {
	capacity := make([]int8, len(d.heads))
	for arc := range capacity {
		capacity[arc] = 1
	}
	incoming := make([]int, d.components())
	for {
		reached := make([]bool, d.components())
		reached[0] = true
		queue := append(make([]int, 0, d.components()), 0)
		for next := 0; next < len(queue); next++ {
			for _, arc := range d.outgoing[queue[next]] {
				head := d.heads[arc]
				if capacity[arc] == 0 || reached[head] {
					continue
				}
				reached[head] = true
				incoming[head] = arc
				queue = append(queue, head)
			}
		}
		if !reached[sink] {
			return wires, len(queue), true
		}
		if wires+1 >= limit {
			return 0, 0, false
		}
		for component := sink; component != 0; {
			arc := incoming[component]
			capacity[arc]--
			capacity[arc^1]++
			component = d.heads[arc^1]
		}
		wires++
	}
}

func solve(input string) (int, error) {
	d, err := parseDiagram(input)
	if err != nil {
		return 0, err
	}
	left, right := d.minimumCut()
	return left * right, nil
}
