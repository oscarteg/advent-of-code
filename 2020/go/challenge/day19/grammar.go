package day19

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

const rootRule = 0

// A rule matches a single literal byte when alternatives is nil, otherwise any one
// of the alternative sequences of rule IDs.
type rule struct {
	literal      byte
	alternatives [][]int
}

type grammar map[int]rule

type cacheKey struct {
	id, start int
}

func parse(input *challenge.Input) (grammar, []string, error) {
	rules := grammar{}
	var messages []string
	readingRules := true

	for _, line := range input.LinesSlice() {
		line = strings.TrimSpace(line)
		switch {
		case line == "":
			readingRules = false
		case !readingRules:
			messages = append(messages, line)
		default:
			id, parsed, err := parseRule(line)
			if err != nil {
				return nil, nil, err
			}
			if _, exists := rules[id]; exists {
				return nil, nil, fmt.Errorf("duplicate rule %d", id)
			}
			rules[id] = parsed
		}
	}

	if len(messages) == 0 {
		return nil, nil, errors.New("no messages after the rules")
	}
	if err := rules.validate(); err != nil {
		return nil, nil, err
	}

	return rules, messages, nil
}

func parseRule(line string) (int, rule, error) {
	head, body, found := strings.Cut(line, ":")
	if !found {
		return 0, rule{}, fmt.Errorf("rule without ':': %q", line)
	}
	id, err := strconv.Atoi(strings.TrimSpace(head))
	if err != nil {
		return 0, rule{}, fmt.Errorf("invalid rule id in %q: %w", line, err)
	}
	parsed, err := parseBody(strings.TrimSpace(body))
	if err != nil {
		return 0, rule{}, fmt.Errorf("rule %d: %w", id, err)
	}

	return id, parsed, nil
}

func parseBody(body string) (rule, error) {
	if strings.HasPrefix(body, `"`) {
		literal := strings.TrimSuffix(strings.TrimPrefix(body, `"`), `"`)
		if len(literal) != 1 {
			return rule{}, fmt.Errorf("literal must be one character: %q", body)
		}
		return rule{literal: literal[0]}, nil
	}

	var alternatives [][]int
	for _, alternative := range strings.Split(body, "|") {
		fields := strings.Fields(alternative)
		if len(fields) == 0 {
			return rule{}, fmt.Errorf("empty alternative in %q", body)
		}
		ids := make([]int, 0, len(fields))
		for _, field := range fields {
			id, err := strconv.Atoi(field)
			if err != nil {
				return rule{}, fmt.Errorf("invalid rule id %q: %w", field, err)
			}
			ids = append(ids, id)
		}
		alternatives = append(alternatives, ids)
	}

	return rule{alternatives: alternatives}, nil
}

// Matching resolves rule IDs without failing, so every reference must exist.
func (g grammar) validate() error {
	if _, exists := g[rootRule]; !exists {
		return fmt.Errorf("rule %d is missing", rootRule)
	}
	for id, r := range g {
		for _, alternative := range r.alternatives {
			for _, referenced := range alternative {
				if _, exists := g[referenced]; !exists {
					return fmt.Errorf("rule %d references unknown rule %d", id, referenced)
				}
			}
		}
	}

	return nil
}

// Part 2 makes rules 8 and 11 recursive. A rule set without them (the part 1
// example) is left alone, so its answer stays the same.
func (g grammar) withLoops() (grammar, error) {
	patched := maps.Clone(g)
	loops := []struct {
		id   int
		body string
	}{
		{id: 8, body: "42 | 42 8"},
		{id: 11, body: "42 31 | 42 11 31"},
	}

	for _, loop := range loops {
		if _, exists := patched[loop.id]; !exists {
			continue
		}
		parsed, err := parseBody(loop.body)
		if err != nil {
			return nil, err
		}
		patched[loop.id] = parsed
	}
	if err := patched.validate(); err != nil {
		return nil, err
	}

	return patched, nil
}

func (g grammar) countMatching(messages []string) int {
	count := 0
	for _, message := range messages {
		if g.matches(message) {
			count++
		}
	}

	return count
}

func (g grammar) matches(message string) bool {
	cache := map[cacheKey][]int{}

	return slices.Contains(g.ruleEnds(rootRule, 0, message, cache), len(message))
}

// ruleEnds returns every position the rule can end at when it starts at start.
// Alternatives of different lengths make this a set instead of a single position,
// which is what makes the recursive part 2 rules work without special cases.
func (g grammar) ruleEnds(id, start int, message string, cache map[cacheKey][]int) []int {
	key := cacheKey{id: id, start: start}
	if ends, cached := cache[key]; cached {
		return ends
	}

	r := g[id]
	var ends []int
	if r.alternatives == nil {
		if start < len(message) && message[start] == r.literal {
			ends = []int{start + 1}
		}
	} else {
		for _, alternative := range r.alternatives {
			ends = append(ends, g.sequenceEnds(alternative, start, message, cache)...)
		}
		ends = sortedUnique(ends)
	}
	cache[key] = ends

	return ends
}

func (g grammar) sequenceEnds(ids []int, start int, message string, cache map[cacheKey][]int) []int {
	positions := []int{start}
	for _, id := range ids {
		var next []int
		for _, position := range positions {
			next = append(next, g.ruleEnds(id, position, message, cache)...)
		}
		if len(next) == 0 {
			return nil
		}
		positions = sortedUnique(next)
	}

	return positions
}

func sortedUnique(values []int) []int {
	slices.Sort(values)

	return slices.Compact(values)
}
