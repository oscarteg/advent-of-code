package day16

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type valueRange struct {
	low, high int
}

func (r valueRange) contains(value int) bool {
	return value >= r.low && value <= r.high
}

type rule struct {
	name   string
	ranges []valueRange
}

func (r rule) accepts(value int) bool {
	return slices.ContainsFunc(r.ranges, func(within valueRange) bool { return within.contains(value) })
}

func (r rule) acceptsColumn(tickets [][]int, column int) bool {
	for _, ticket := range tickets {
		if !r.accepts(ticket[column]) {
			return false
		}
	}
	return true
}

// notes holds the three sections of the input: the field rules, your own ticket and the
// nearby tickets.
type notes struct {
	rules  []rule
	ticket []int
	nearby [][]int
}

func (n *notes) acceptedByAnyRule(value int) bool {
	return slices.ContainsFunc(n.rules, func(r rule) bool { return r.accepts(value) })
}

func (n *notes) isValid(ticket []int) bool {
	for _, value := range ticket {
		if !n.acceptedByAnyRule(value) {
			return false
		}
	}
	return true
}

// errorRate sums every nearby value that no rule accepts.
func (n *notes) errorRate() int {
	rate := 0
	for _, ticket := range n.nearby {
		for _, value := range ticket {
			if !n.acceptedByAnyRule(value) {
				rate += value
			}
		}
	}
	return rate
}

// fieldOrder matches every ticket column to a field name. A column starts out with every field
// that accepts its value on all valid tickets; the column left with a single candidate is
// assigned and that field is struck from the other columns, until every column has its field.
func (n *notes) fieldOrder() ([]string, error) {
	valid := make([][]int, 0, len(n.nearby))
	for _, ticket := range n.nearby {
		if n.isValid(ticket) {
			valid = append(valid, ticket)
		}
	}

	candidates := make([][]int, len(n.rules))
	for column := range candidates {
		for field, r := range n.rules {
			if r.acceptsColumn(valid, column) {
				candidates[column] = append(candidates[column], field)
			}
		}
	}

	order := make([]string, len(n.rules))
	for range n.rules {
		column := slices.IndexFunc(candidates, func(fields []int) bool { return len(fields) == 1 })
		if column < 0 {
			return nil, errors.New("the notes do not pin down a single field order")
		}
		field := candidates[column][0]
		order[column] = n.rules[field].name
		for other := range candidates {
			candidates[other] = slices.DeleteFunc(candidates[other], func(candidate int) bool { return candidate == field })
		}
	}
	return order, nil
}

func parseRange(text string) (valueRange, error) {
	low, high, found := strings.Cut(strings.TrimSpace(text), "-")
	if !found {
		return valueRange{}, fmt.Errorf("invalid range %q", text)
	}
	from, err := strconv.Atoi(low)
	if err != nil {
		return valueRange{}, fmt.Errorf("invalid range start %q: %w", low, err)
	}
	to, err := strconv.Atoi(high)
	if err != nil {
		return valueRange{}, fmt.Errorf("invalid range end %q: %w", high, err)
	}
	if from > to {
		return valueRange{}, fmt.Errorf("range %d-%d is empty", from, to)
	}
	return valueRange{low: from, high: to}, nil
}

func parseRule(line string) (rule, error) {
	name, ranges, found := strings.Cut(line, ": ")
	if !found {
		return rule{}, fmt.Errorf("invalid field rule %q", line)
	}
	parts := strings.Split(ranges, " or ")
	parsed := make([]valueRange, 0, len(parts))
	for _, part := range parts {
		within, err := parseRange(part)
		if err != nil {
			return rule{}, err
		}
		parsed = append(parsed, within)
	}
	return rule{name: name, ranges: parsed}, nil
}

func parseTicket(line string) ([]int, error) {
	fields := strings.Split(line, ",")
	ticket := make([]int, 0, len(fields))
	for _, field := range fields {
		value, err := strconv.Atoi(strings.TrimSpace(field))
		if err != nil {
			return nil, fmt.Errorf("invalid ticket value %q: %w", field, err)
		}
		ticket = append(ticket, value)
	}
	return ticket, nil
}

func parseTickets(section []string, header string) ([][]int, error) {
	if section[0] != header {
		return nil, fmt.Errorf("expected the %q header, found %q", header, section[0])
	}
	tickets := make([][]int, 0, len(section)-1)
	for _, line := range section[1:] {
		ticket, err := parseTicket(line)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

// splitSections groups the non-empty lines into the blocks separated by blank lines.
func splitSections(lines []string) [][]string {
	var sections [][]string
	var current []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if len(current) > 0 {
				sections = append(sections, current)
				current = nil
			}
			continue
		}
		current = append(current, trimmed)
	}
	if len(current) > 0 {
		sections = append(sections, current)
	}
	return sections
}

func checkValueCount(tickets [][]int, fields int) error {
	for _, ticket := range tickets {
		if len(ticket) != fields {
			return fmt.Errorf("ticket has %d values but there are %d fields", len(ticket), fields)
		}
	}
	return nil
}

func parseNotes(lines []string) (*notes, error) {
	sections := splitSections(lines)
	if len(sections) != 3 {
		return nil, fmt.Errorf("expected 3 sections, found %d", len(sections))
	}

	rules := make([]rule, 0, len(sections[0]))
	for _, line := range sections[0] {
		parsed, err := parseRule(line)
		if err != nil {
			return nil, err
		}
		rules = append(rules, parsed)
	}

	own, err := parseTickets(sections[1], "your ticket:")
	if err != nil {
		return nil, err
	}
	if len(own) != 1 {
		return nil, fmt.Errorf("expected exactly one ticket of your own, found %d", len(own))
	}

	nearby, err := parseTickets(sections[2], "nearby tickets:")
	if err != nil {
		return nil, err
	}
	if len(nearby) == 0 {
		return nil, errors.New("missing the nearby tickets")
	}

	if err := checkValueCount(own, len(rules)); err != nil {
		return nil, err
	}
	if err := checkValueCount(nearby, len(rules)); err != nil {
		return nil, err
	}

	return &notes{rules: rules, ticket: own[0], nearby: nearby}, nil
}
