package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	// Parts enter the system through this workflow and leave it accepted or
	// rejected; every rating lies between minRating and maxRating inclusive.
	startWorkflow = "in"
	acceptTarget  = "A"
	rejectTarget  = "R"
	minRating     = 1
	maxRating     = 4000
)

type category int

const (
	categoryX category = iota
	categoryM
	categoryA
	categoryS
	categoryCount
)

func parseCategory(name string) (category, error) {
	switch name {
	case "x":
		return categoryX, nil
	case "m":
		return categoryM, nil
	case "a":
		return categoryA, nil
	case "s":
		return categoryS, nil
	default:
		return 0, fmt.Errorf("unknown rating category %q", name)
	}
}

// part holds one machine part's x, m, a and s ratings.
type part [categoryCount]int64

func (p part) totalRating() int64 {
	var total int64
	for _, rating := range p {
		total += rating
	}
	return total
}

func parsePart(line string) (part, error) {
	body, openend := strings.CutPrefix(line, "{")
	body, closed := strings.CutSuffix(body, "}")
	if !openend || !closed {
		return part{}, fmt.Errorf("part must be wrapped in braces: %q", line)
	}

	var (
		p    part
		seen [categoryCount]bool
	)
	for _, assignment := range strings.Split(body, ",") {
		name, value, found := strings.Cut(assignment, "=")
		if !found {
			return part{}, fmt.Errorf("rating needs category=value: %q", assignment)
		}
		cat, err := parseCategory(strings.TrimSpace(name))
		if err != nil {
			return part{}, err
		}
		if seen[cat] {
			return part{}, fmt.Errorf("duplicate rating category %q", name)
		}
		rating, err := parseRating(value)
		if err != nil {
			return part{}, fmt.Errorf("rating %q: %w", assignment, err)
		}
		p[cat], seen[cat] = rating, true
	}
	for _, found := range seen {
		if !found {
			return part{}, fmt.Errorf("part misses a rating category: %q", line)
		}
	}
	return p, nil
}

// parseRating limits ratings and rule thresholds to 32 bits so that the range
// arithmetic in rule.split cannot overflow an int64.
func parseRating(text string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(text), 10, 32)
}

// rule routes a part to target when its condition holds. Only the last rule of
// a workflow is unconditional (operator 0), so it always applies once reached.
type rule struct {
	cat      category
	operator byte
	value    int64
	target   string
}

func parseRule(field string) (rule, error) {
	condition, target, found := strings.Cut(field, ":")
	if !found {
		if field == "" {
			return rule{}, fmt.Errorf("rule has no target")
		}
		return rule{target: field}, nil
	}
	if target == "" {
		return rule{}, fmt.Errorf("rule has no target: %q", field)
	}
	operator := strings.IndexAny(condition, "<>")
	if operator < 0 {
		return rule{}, fmt.Errorf("condition needs < or >: %q", condition)
	}
	cat, err := parseCategory(strings.TrimSpace(condition[:operator]))
	if err != nil {
		return rule{}, err
	}
	value, err := parseRating(condition[operator+1:])
	if err != nil {
		return rule{}, fmt.Errorf("condition %q: %w", condition, err)
	}
	return rule{cat: cat, operator: condition[operator], value: value, target: target}, nil
}

func (r rule) unconditional() bool {
	return r.operator == 0
}

func (r rule) matches(p part) bool {
	switch r.operator {
	case '<':
		return p[r.cat] < r.value
	case '>':
		return p[r.cat] > r.value
	default:
		return true
	}
}

// split divides a search space into the ratings that satisfy the rule's
// condition and the ratings that fall through to the next rule.
func (r rule) split(space ratings) (matched, rest ratings) {
	if r.operator == '<' {
		return space.clamp(r.cat, minRating, r.value-1), space.clamp(r.cat, r.value, maxRating)
	}
	return space.clamp(r.cat, r.value+1, maxRating), space.clamp(r.cat, minRating, r.value)
}

type workflow []rule

func parseWorkflow(line string) (string, workflow, error) {
	name, body, found := strings.Cut(line, "{")
	if !found {
		return "", nil, fmt.Errorf("workflow needs a rule list: %q", line)
	}
	if name == "" {
		return "", nil, fmt.Errorf("workflow has no name: %q", line)
	}
	body, closed := strings.CutSuffix(body, "}")
	if !closed {
		return "", nil, fmt.Errorf("workflow rule list is unclosed: %q", line)
	}

	fields := strings.Split(body, ",")
	w := make(workflow, 0, len(fields))
	for _, field := range fields {
		r, err := parseRule(field)
		if err != nil {
			return "", nil, fmt.Errorf("workflow %q: %w", name, err)
		}
		w = append(w, r)
	}
	for i, r := range w {
		if r.unconditional() != (i == len(w)-1) {
			return "", nil, fmt.Errorf("workflow %q: only the last rule may be unconditional", name)
		}
	}
	return name, w, nil
}

// route returns the target of the first rule that matches; the unconditional
// last rule guarantees that one always does.
func (w workflow) route(p part) string {
	for _, r := range w {
		if r.matches(p) {
			return r.target
		}
	}
	return ""
}

type system struct {
	workflows map[string]workflow
	parts     []part
}

func parseSystem(input string) (system, error) {
	s := system{workflows: make(map[string]workflow)}
	readingParts := false
	for _, raw := range strings.Split(input, "\n") {
		line := strings.TrimSpace(raw)
		switch {
		case line == "":
			readingParts = len(s.workflows) > 0
		case readingParts:
			p, err := parsePart(line)
			if err != nil {
				return system{}, err
			}
			s.parts = append(s.parts, p)
		default:
			name, w, err := parseWorkflow(line)
			if err != nil {
				return system{}, err
			}
			if _, duplicate := s.workflows[name]; duplicate {
				return system{}, fmt.Errorf("duplicate workflow %q", name)
			}
			s.workflows[name] = w
		}
	}

	if _, ok := s.workflows[startWorkflow]; !ok {
		return system{}, fmt.Errorf("missing workflow %q", startWorkflow)
	}
	if len(s.parts) == 0 {
		return system{}, fmt.Errorf("missing part ratings")
	}
	if err := s.validateTargets(); err != nil {
		return system{}, err
	}
	return s, nil
}

func (s system) validateTargets() error {
	for name, w := range s.workflows {
		for _, r := range w {
			if r.target == acceptTarget || r.target == rejectTarget {
				continue
			}
			if _, ok := s.workflows[r.target]; !ok {
				return fmt.Errorf("workflow %q routes to unknown workflow %q", name, r.target)
			}
		}
	}
	return nil
}

// accepts walks a part through the workflows until it is accepted or rejected.
func (s system) accepts(p part) (bool, error) {
	target := startWorkflow
	// Visiting more workflows than exist means the rules form a cycle.
	for hop := 0; hop <= len(s.workflows); hop++ {
		switch target {
		case acceptTarget:
			return true, nil
		case rejectTarget:
			return false, nil
		}
		w, ok := s.workflows[target]
		if !ok {
			return false, fmt.Errorf("unknown workflow %q", target)
		}
		target = w.route(p)
	}
	return false, fmt.Errorf("workflows form a cycle")
}

// span is an inclusive range of still-possible ratings; it is empty when
// high < low.
type span struct {
	low, high int64
}

func (s span) size() int64 {
	if s.high < s.low {
		return 0
	}
	return s.high - s.low + 1
}

// ratings is a search space of part ratings: one span per category.
type ratings [categoryCount]span

func fullRatings() ratings {
	var space ratings
	for cat := range space {
		space[cat] = span{low: minRating, high: maxRating}
	}
	return space
}

func (r ratings) combinations() int64 {
	total := int64(1)
	for _, s := range r {
		total *= s.size()
	}
	return total
}

func (r ratings) clamp(cat category, low, high int64) ratings {
	clamped := r
	clamped[cat] = span{low: max(r[cat].low, low), high: min(r[cat].high, high)}
	return clamped
}

// countAccepted pushes a whole search space through target instead of a single
// part. Every condition splits the space into the ratings it routes away and
// the ratings the next rule sees, so the accepted fragments are disjoint and
// their combination counts can simply be summed.
func (s system) countAccepted(target string, space ratings, depth int) (int64, error) {
	combinations := space.combinations()
	if combinations == 0 {
		return 0, nil
	}
	switch target {
	case acceptTarget:
		return combinations, nil
	case rejectTarget:
		return 0, nil
	}
	if depth > len(s.workflows) {
		return 0, fmt.Errorf("workflows form a cycle")
	}
	w, ok := s.workflows[target]
	if !ok {
		return 0, fmt.Errorf("unknown workflow %q", target)
	}

	var accepted int64
	remaining := space
	for _, r := range w {
		if r.unconditional() {
			count, err := s.countAccepted(r.target, remaining, depth+1)
			if err != nil {
				return 0, err
			}
			return accepted + count, nil
		}
		matched, rest := r.split(remaining)
		count, err := s.countAccepted(r.target, matched, depth+1)
		if err != nil {
			return 0, err
		}
		accepted, remaining = accepted+count, rest
	}
	return accepted, nil
}

func solve(input string) (int64, int64, error) {
	s, err := parseSystem(input)
	if err != nil {
		return 0, 0, err
	}

	var totalRating int64
	for _, p := range s.parts {
		accepted, err := s.accepts(p)
		if err != nil {
			return 0, 0, err
		}
		if accepted {
			totalRating += p.totalRating()
		}
	}

	combinations, err := s.countAccepted(startWorkflow, fullRatings(), 0)
	if err != nil {
		return 0, 0, err
	}
	return totalRating, combinations, nil
}
