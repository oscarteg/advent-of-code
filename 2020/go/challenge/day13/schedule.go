package day13

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

type bus struct {
	offset int
	id     int
}

type schedule struct {
	earliest int
	buses    []bus
}

func parseSchedule(input *challenge.Input) (schedule, error) {
	lines := input.LinesSlice()
	if len(lines) < 2 {
		return schedule{}, fmt.Errorf("expected an earliest timestamp and a bus id list, got %d lines", len(lines))
	}

	earliest, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return schedule{}, fmt.Errorf("invalid earliest departure timestamp: %w", err)
	}

	var buses []bus
	for offset, field := range strings.Split(strings.TrimSpace(lines[1]), ",") {
		if field == "x" {
			continue
		}
		id, err := strconv.Atoi(field)
		if err != nil {
			return schedule{}, fmt.Errorf("invalid bus id %q: %w", field, err)
		}
		if id <= 0 {
			return schedule{}, fmt.Errorf("bus id must be positive, got %d", id)
		}
		buses = append(buses, bus{offset: offset, id: id})
	}
	if len(buses) == 0 {
		return schedule{}, fmt.Errorf("bus id list must contain at least one bus")
	}

	return schedule{earliest: earliest, buses: buses}, nil
}
