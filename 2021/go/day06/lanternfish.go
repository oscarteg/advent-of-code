package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidTimer = errors.New("invalid fish timer")

func parseTimers(input string) ([]uint8, error) {
	fields := strings.Split(strings.TrimSpace(input), ",")
	timers := make([]uint8, 0, len(fields))
	for _, field := range fields {
		value, err := strconv.ParseUint(strings.TrimSpace(field), 10, 8)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errInvalidTimer, err)
		}
		if value > 8 {
			return nil, fmt.Errorf("%w: must be between 0 and 8", errInvalidTimer)
		}
		timers = append(timers, uint8(value))
	}
	if len(timers) == 0 {
		return nil, fmt.Errorf("%w: missing fish timers", errInvalidTimer)
	}
	return timers, nil
}

// simulate tracks how many fish share each timer value instead of simulating
// every fish individually, since the population grows exponentially.
func simulate(timers []uint8, days int) uint64 {
	var counts [9]uint64
	for _, timer := range timers {
		counts[timer]++
	}
	for day := 0; day < days; day++ {
		spawning := counts[0]
		for i := 0; i < 8; i++ {
			counts[i] = counts[i+1]
		}
		counts[6] += spawning
		counts[8] = spawning
	}
	var total uint64
	for _, count := range counts {
		total += count
	}
	return total
}

func solve(input string) (uint64, uint64, error) {
	timers, err := parseTimers(input)
	if err != nil {
		return 0, 0, err
	}
	return simulate(timers, 80), simulate(timers, 256), nil
}
