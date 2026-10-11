package main

import (
	"errors"
	"fmt"
)

const buttonPresses = 1000

// solve pushes the button 1000 times for part 1, then finds how many
// presses it takes the puzzle's implicit sink module to go low for part 2.
func solve(input string) (uint64, uint64, error) {
	net, err := parseNetwork(input)
	if err != nil {
		return 0, 0, err
	}

	var low, high uint64
	for i := 0; i < buttonPresses; i++ {
		for _, p := range net.pressButton() {
			if p.high {
				high++
			} else {
				low++
			}
		}
	}
	part1 := low * high

	net2, err := parseNetwork(input)
	if err != nil {
		return 0, 0, err
	}
	target, err := net2.implicitSink()
	if errors.Is(err, errNoImplicitSink) {
		// This graph has no rx-style target, as in the puzzle's own part 1
		// examples; part 2 does not apply.
		return part1, 0, nil
	}
	if err != nil {
		return 0, 0, err
	}
	part2, err := pressesUntilTargetGoesLow(net2, target)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}

const maxButtonPresses = 1_000_000

// pressesUntilTargetGoesLow finds target's single feeder module, then
// returns the least common multiple of the first button press at which
// each of that feeder's inputs sends it a high pulse.
func pressesUntilTargetGoesLow(net network, target string) (uint64, error) {
	feeders := net.inputsOf(target)
	switch len(feeders) {
	case 0:
		return 0, fmt.Errorf("%w: target module has no input", errInvalidNetwork)
	case 1:
	default:
		return 0, fmt.Errorf("%w: target module needs exactly one input", errInvalidNetwork)
	}
	feeder := feeders[0]

	watched := net.inputsOf(feeder)
	if len(watched) == 0 {
		return 0, fmt.Errorf("%w: feeder module has no inputs", errInvalidNetwork)
	}

	firstHigh := map[string]uint64{}
	for press := uint64(1); press <= maxButtonPresses; press++ {
		if len(firstHigh) == len(watched) {
			break
		}
		for _, p := range net.pressButton() {
			if p.high && p.target == feeder {
				if _, seen := firstHigh[p.source]; !seen {
					firstHigh[p.source] = press
				}
			}
		}
	}
	if len(firstHigh) != len(watched) {
		return 0, fmt.Errorf("%w: no high-pulse cycle found", errInvalidNetwork)
	}

	answer := uint64(1)
	for _, n := range firstHigh {
		answer = lcm(answer, n)
	}
	return answer, nil
}

func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b uint64) uint64 {
	return a / gcd(a, b) * b
}
