package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	part1Minutes    = 24
	part2Minutes    = 32
	part2Blueprints = 3
)

const (
	oreRobot = iota
	clayRobot
	obsidianRobot
	geodeRobot
	robotTypes
)

const (
	ore = iota
	clay
	obsidian
	resourceCount
)

var buildOrder = [robotTypes]int{geodeRobot, obsidianRobot, clayRobot, oreRobot}

var errInvalidBlueprint = errors.New("invalid blueprint")

// scanner reads a blueprint line literal by literal, keeping the first error
// so malformed input is rejected instead of silently misparsed.
type scanner struct {
	rest string
	err  error
}

func (s *scanner) literal(expected string) {
	if s.err != nil {
		return
	}
	rest := strings.TrimLeft(s.rest, " \t\r")
	if !strings.HasPrefix(rest, expected) {
		s.err = fmt.Errorf("%w: expected %q, got %q", errInvalidBlueprint, expected, rest)
		return
	}
	s.rest = rest[len(expected):]
}

func (s *scanner) number() int {
	if s.err != nil {
		return 0
	}
	rest := strings.TrimLeft(s.rest, " \t\r")
	digits := 0
	for digits < len(rest) && rest[digits] >= '0' && rest[digits] <= '9' {
		digits++
	}
	if digits == 0 {
		s.err = fmt.Errorf("%w: expected a number, got %q", errInvalidBlueprint, rest)
		return 0
	}
	value, err := strconv.Atoi(rest[:digits])
	if err != nil {
		s.err = fmt.Errorf("%w: %w", errInvalidBlueprint, err)
		return 0
	}
	s.rest = rest[digits:]
	return value
}

func (s *scanner) finish() {
	if s.err != nil {
		return
	}
	if trailing := strings.TrimSpace(s.rest); trailing != "" {
		s.err = fmt.Errorf("%w: unexpected trailing text %q", errInvalidBlueprint, trailing)
	}
}

// blueprint holds the robot costs in ore, clay and obsidian, indexed by robot
// type, plus the highest useful number of ore, clay and obsidian robots: no
// recipe can spend more of a resource per minute than its largest cost.
type blueprint struct {
	id    int
	costs [robotTypes][resourceCount]int
	caps  [resourceCount]int
}

func parseBlueprint(line string) (blueprint, error) {
	s := &scanner{rest: line}
	s.literal("Blueprint")
	id := s.number()
	s.literal(":")
	s.literal("Each ore robot costs")
	oreRobotOre := s.number()
	s.literal("ore.")
	s.literal("Each clay robot costs")
	clayRobotOre := s.number()
	s.literal("ore.")
	s.literal("Each obsidian robot costs")
	obsidianRobotOre := s.number()
	s.literal("ore and")
	obsidianRobotClay := s.number()
	s.literal("clay.")
	s.literal("Each geode robot costs")
	geodeRobotOre := s.number()
	s.literal("ore and")
	geodeRobotObsidian := s.number()
	s.literal("obsidian.")
	s.finish()
	if s.err != nil {
		return blueprint{}, s.err
	}

	return blueprint{
		id: id,
		costs: [robotTypes][resourceCount]int{
			{oreRobotOre, 0, 0},
			{clayRobotOre, 0, 0},
			{obsidianRobotOre, obsidianRobotClay, 0},
			{geodeRobotOre, 0, geodeRobotObsidian},
		},
		caps: [resourceCount]int{
			max(oreRobotOre, clayRobotOre, obsidianRobotOre, geodeRobotOre),
			obsidianRobotClay,
			geodeRobotObsidian,
		},
	}, nil
}

// state is a moment in a blueprint's schedule. Geode robots are not tracked:
// the geodes they will have cracked once the clock runs out are added to
// geodes the minute they are finished.
type state struct {
	minutesLeft int
	resources   [resourceCount]int
	robots      [resourceCount]int
	geodes      int
}

// build waits until robot is affordable and builds it, reporting false when it
// can never pay off before the clock runs out.
func (s state) build(robot int, cost [resourceCount]int) (state, bool) {
	wait := 0
	for resource, price := range cost {
		missing := price - s.resources[resource]
		if missing <= 0 {
			continue
		}
		rate := s.robots[resource]
		if rate == 0 {
			return state{}, false
		}
		wait = max(wait, (missing+rate-1)/rate)
	}

	elapsed := wait + 1
	if elapsed >= s.minutesLeft {
		return state{}, false
	}

	next := s
	next.minutesLeft -= elapsed
	for resource := range next.resources {
		next.resources[resource] += s.robots[resource]*elapsed - cost[resource]
	}
	if robot == geodeRobot {
		next.geodes += next.minutesLeft
	} else {
		next.robots[robot]++
	}
	return next, true
}

func (b blueprint) maxGeodes(minutes int) int {
	best := 0
	b.search(state{minutesLeft: minutes, robots: [resourceCount]int{1, 0, 0}}, &best)
	return best
}

// search explores the next robot to build, geode robots first, dropping every
// branch whose optimistic bound cannot beat the best schedule found so far.
func (b blueprint) search(current state, best *int) {
	if current.geodes > *best {
		*best = current.geodes
	}
	if b.optimisticGeodes(current) <= *best {
		return
	}
	for _, robot := range buildOrder {
		if robot != geodeRobot && current.robots[robot] >= b.caps[robot] {
			continue
		}
		if next, ok := current.build(robot, b.costs[robot]); ok {
			b.search(next, best)
		}
	}
}

// optimisticGeodes bounds the geodes reachable from current with a relaxed
// game that ignores ore costs and lets one robot of every type be built per
// minute. Building everything affordable as early as possible is optimal there
// (clay is only spent on obsidian robots, obsidian only on geode robots), so
// no real schedule can beat the result.
func (b blueprint) optimisticGeodes(current state) int {
	obsidianPrice := b.costs[geodeRobot][obsidian]
	clayPrice := b.costs[obsidianRobot][clay]
	clayAmount := current.resources[clay]
	obsidianAmount := current.resources[obsidian]
	clayRobots := current.robots[clay]
	obsidianRobots := current.robots[obsidian]
	geodes := current.geodes

	for remaining := current.minutesLeft; remaining > 0; remaining-- {
		builtObsidianRobot := 0
		if obsidianAmount >= obsidianPrice {
			obsidianAmount -= obsidianPrice
			geodes += remaining - 1
		}
		if clayAmount >= clayPrice {
			clayAmount -= clayPrice
			builtObsidianRobot = 1
		}
		clayAmount += clayRobots
		obsidianAmount += obsidianRobots
		clayRobots++
		obsidianRobots += builtObsidianRobot
	}
	return geodes
}

func parseBlueprints(input string) ([]blueprint, error) {
	var blueprints []blueprint
	for number, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parsed, err := parseBlueprint(line)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", number+1, err)
		}
		blueprints = append(blueprints, parsed)
	}
	if len(blueprints) == 0 {
		return nil, fmt.Errorf("%w: input has no blueprints", errInvalidBlueprint)
	}
	return blueprints, nil
}

func solve(input string) (int, int, error) {
	blueprints, err := parseBlueprints(input)
	if err != nil {
		return 0, 0, err
	}

	part1, part2 := 0, 1
	for index, b := range blueprints {
		part1 += b.id * b.maxGeodes(part1Minutes)
		if index < part2Blueprints {
			part2 *= b.maxGeodes(part2Minutes)
		}
	}
	return part1, part2, nil
}
