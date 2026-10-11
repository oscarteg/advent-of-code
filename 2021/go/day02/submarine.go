package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidCommand = errors.New("invalid command")

type direction string

const (
	forward direction = "forward"
	down    direction = "down"
	up      direction = "up"
)

type command struct {
	dir      direction
	distance uint32
}

func parseCommand(line string) (command, error) {
	action, value, ok := strings.Cut(line, " ")
	if !ok {
		return command{}, fmt.Errorf("%w: %q", errInvalidCommand, line)
	}
	distance, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		return command{}, fmt.Errorf("%w: %q: %w", errInvalidCommand, line, err)
	}
	switch direction(action) {
	case forward, down, up:
		return command{dir: direction(action), distance: uint32(distance)}, nil
	default:
		return command{}, fmt.Errorf("%w: unknown direction %q", errInvalidCommand, action)
	}
}

func parseCommands(input string) ([]command, error) {
	var commands []command
	for _, line := range strings.Split(strings.TrimSpace(input), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		cmd, err := parseCommand(line)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

func planCourse(commands []command) uint64 {
	var horizontal, depth uint32
	for _, cmd := range commands {
		switch cmd.dir {
		case forward:
			horizontal += cmd.distance
		case down:
			depth += cmd.distance
		case up:
			if cmd.distance > depth {
				depth = 0
			} else {
				depth -= cmd.distance
			}
		}
	}
	return uint64(horizontal) * uint64(depth)
}

func planCourseWithAim(commands []command) uint64 {
	var horizontal, depth, aim uint32
	for _, cmd := range commands {
		switch cmd.dir {
		case forward:
			horizontal += cmd.distance
			depth += aim * cmd.distance
		case down:
			aim += cmd.distance
		case up:
			if cmd.distance > aim {
				aim = 0
			} else {
				aim -= cmd.distance
			}
		}
	}
	return uint64(horizontal) * uint64(depth)
}

func solve(input string) (uint64, uint64, error) {
	commands, err := parseCommands(input)
	if err != nil {
		return 0, 0, err
	}
	return planCourse(commands), planCourseWithAim(commands), nil
}
