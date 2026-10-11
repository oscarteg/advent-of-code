package main

import (
	"fmt"
	"strings"
)

// priority returns an item's priority: a-z are 1-26, A-Z are 27-52.
func priority(c byte) (int, error) {
	switch {
	case c >= 'a' && c <= 'z':
		return int(c-'a') + 1, nil
	case c >= 'A' && c <= 'Z':
		return int(c-'A') + 27, nil
	default:
		return 0, fmt.Errorf("invalid item type: %q", c)
	}
}

// itemSet returns the set of distinct item types present in s.
func itemSet(s string) map[byte]struct{} {
	set := make(map[byte]struct{}, len(s))
	for i := 0; i < len(s); i++ {
		set[s[i]] = struct{}{}
	}
	return set
}

// commonItem finds the single item type shared by every one of the given strings.
func commonItem(groups []string) (byte, error) {
	if len(groups) == 0 {
		return 0, fmt.Errorf("no compartments/rucksacks to compare")
	}

	common := itemSet(groups[0])
	for _, group := range groups[1:] {
		next := itemSet(group)
		for item := range common {
			if _, ok := next[item]; !ok {
				delete(common, item)
			}
		}
	}

	if len(common) == 0 {
		return 0, fmt.Errorf("no item common to all compartments/rucksacks")
	}
	if len(common) > 1 {
		return 0, fmt.Errorf("more than one common item found")
	}
	for item := range common {
		return item, nil
	}
	panic("unreachable")
}

func parseLines(input string) []string {
	var lines []string
	for _, line := range strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n") {
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func part1(input string) (int, error) {
	sum := 0
	for _, rucksack := range parseLines(input) {
		if len(rucksack)%2 != 0 {
			return 0, fmt.Errorf("rucksack must have an even number of items: %q", rucksack)
		}
		half := len(rucksack) / 2
		item, err := commonItem([]string{rucksack[:half], rucksack[half:]})
		if err != nil {
			return 0, fmt.Errorf("rucksack %q has no shared item: %w", rucksack, err)
		}
		p, err := priority(item)
		if err != nil {
			return 0, err
		}
		sum += p
	}
	return sum, nil
}

func part2(input string) (int, error) {
	lines := parseLines(input)
	if len(lines) == 0 || len(lines)%3 != 0 {
		return 0, fmt.Errorf("input must contain a multiple of 3 rucksacks, got %d", len(lines))
	}

	sum := 0
	for i := 0; i < len(lines); i += 3 {
		group := lines[i : i+3]
		badge, err := commonItem(group)
		if err != nil {
			return 0, fmt.Errorf("group %v has no shared badge item: %w", group, err)
		}
		p, err := priority(badge)
		if err != nil {
			return 0, err
		}
		sum += p
	}
	return sum, nil
}

func solve(input string) (part1Answer, part2Answer int, err error) {
	part1Answer, err = part1(input)
	if err != nil {
		return 0, 0, err
	}
	part2Answer, err = part2(input)
	if err != nil {
		return 0, 0, err
	}
	return part1Answer, part2Answer, nil
}
