package main

import (
	"fmt"
	"strings"
)

// Entry represents one line of input with signal patterns and output digits.
type Entry struct {
	Patterns []string
	Output   []string
}

func solve(input string) (int, int, error) {
	lines := strings.TrimSpace(input)
	if lines == "" {
		return 0, 0, fmt.Errorf("empty input")
	}

	entries, err := parseInput(lines)
	if err != nil {
		return 0, 0, err
	}

	part1 := countEasyDigits(entries)
	part2, err := sumDecodedOutputs(entries)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}

func parseInput(input string) ([]Entry, error) {
	var entries []Entry
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid entry format: %s", line)
		}

		patterns := strings.Fields(parts[0])
		output := strings.Fields(parts[1])

		if len(patterns) != 10 {
			return nil, fmt.Errorf("expected 10 patterns, got %d", len(patterns))
		}
		if len(output) != 4 {
			return nil, fmt.Errorf("expected 4 output digits, got %d", len(output))
		}

		for i := range patterns {
			patterns[i] = sortString(patterns[i])
		}
		for i := range output {
			output[i] = sortString(output[i])
		}

		entries = append(entries, Entry{Patterns: patterns, Output: output})
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("no valid entries found")
	}

	return entries, nil
}

func sortString(s string) string {
	chars := []rune(s)
	for i := 0; i < len(chars); i++ {
		for j := i + 1; j < len(chars); j++ {
			if chars[i] > chars[j] {
				chars[i], chars[j] = chars[j], chars[i]
			}
		}
	}
	return string(chars)
}

func countEasyDigits(entries []Entry) int {
	count := 0
	easyLengths := map[int]bool{2: true, 3: true, 4: true, 7: true}
	for _, entry := range entries {
		for _, output := range entry.Output {
			if easyLengths[len(output)] {
				count++
			}
		}
	}
	return count
}

func sumDecodedOutputs(entries []Entry) (int, error) {
	sum := 0
	for _, entry := range entries {
		mapping, err := deduceMapping(entry.Patterns)
		if err != nil {
			return 0, err
		}

		value := 0
		for _, output := range entry.Output {
			digit, ok := mapping[output]
			if !ok {
				return 0, fmt.Errorf("unknown pattern: %s", output)
			}
			value = value*10 + digit
		}
		sum += value
	}
	return sum, nil
}

func deduceMapping(patterns []string) (map[string]int, error) {
	mapping := make(map[string]int)
	byLength := make(map[int][]string)

	for _, p := range patterns {
		byLength[len(p)] = append(byLength[len(p)], p)
	}

	one := byLength[2][0]
	seven := byLength[3][0]
	four := byLength[4][0]
	eight := byLength[7][0]

	mapping[one] = 1
	mapping[seven] = 7
	mapping[four] = 4
	mapping[eight] = 8

	// 6-segment deduction: 9 has all of 4, 0 has 1 but not 4, 6 lacks 1
	var zero, six, nine string
	for _, p := range byLength[6] {
		if containsAll(p, four) {
			nine = p
		}
	}

	for _, p := range byLength[6] {
		if p == nine {
			continue
		}
		if containsAll(p, seven) && !containsAll(p, four) {
			zero = p
		} else if !containsAll(p, one) {
			six = p
		}
	}

	mapping[zero] = 0
	mapping[six] = 6
	mapping[nine] = 9

	// 5-segment deduction: 3 has 1, 5 fits in 6, 2 is the rest
	var two, three, five string
	for _, p := range byLength[5] {
		if containsAll(p, one) {
			three = p
		}
	}

	for _, p := range byLength[5] {
		if p == three {
			continue
		}
		if containsAll(six, p) {
			five = p
		} else {
			two = p
		}
	}

	mapping[two] = 2
	mapping[three] = 3
	mapping[five] = 5

	return mapping, nil
}

func containsAll(a, b string) bool {
	for _, c := range b {
		found := false
		for _, ac := range a {
			if ac == c {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
