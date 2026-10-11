package main

import (
	"fmt"
	"strings"
)

func solve(input string) (part1, part2 int64, err error) {
	template, rules, err := parseInput(input)
	if err != nil {
		return 0, 0, err
	}

	part1, err = simulateSteps(template, rules, 10)
	if err != nil {
		return 0, 0, err
	}

	part2, err = simulateSteps(template, rules, 40)
	if err != nil {
		return 0, 0, err
	}

	return part1, part2, nil
}

func parseInput(input string) (string, map[string]rune, error) {
	sections := strings.Split(strings.TrimSpace(input), "\n\n")
	if len(sections) != 2 {
		return "", nil, fmt.Errorf("expected template and rules sections")
	}

	template := sections[0]
	if template == "" {
		return "", nil, fmt.Errorf("template cannot be empty")
	}

	rules := make(map[string]rune)
	for _, line := range strings.Split(sections[1], "\n") {
		if line == "" {
			continue
		}
		parts := strings.Split(line, " -> ")
		if len(parts) != 2 {
			return "", nil, fmt.Errorf("invalid rule format: %s", line)
		}
		pair := parts[0]
		if len(pair) != 2 {
			return "", nil, fmt.Errorf("pair must be exactly 2 characters: %s", pair)
		}
		insertion := []rune(parts[1])
		if len(insertion) != 1 {
			return "", nil, fmt.Errorf("insertion must be exactly 1 character: %s", parts[1])
		}
		rules[pair] = insertion[0]
	}

	if len(rules) == 0 {
		return "", nil, fmt.Errorf("no insertion rules found")
	}

	return template, rules, nil
}

func simulateSteps(template string, rules map[string]rune, steps int) (int64, error) {
	pairCounts := make(map[string]int64)
	for i := 0; i < len(template)-1; i++ {
		pair := template[i : i+2]
		pairCounts[pair]++
	}

	firstChar := rune(template[0])

	for step := 0; step < steps; step++ {
		newPairCounts := make(map[string]int64)
		for pair, count := range pairCounts {
			if insertion, ok := rules[pair]; ok {
				left := string([]rune{rune(pair[0]), insertion})
				right := string([]rune{insertion, rune(pair[1])})
				newPairCounts[left] += count
				newPairCounts[right] += count
			} else {
				newPairCounts[pair] += count
			}
		}
		pairCounts = newPairCounts
	}

	charCounts := make(map[rune]int64)
	charCounts[firstChar]++

	for pair, count := range pairCounts {
		secondChar := rune(pair[1])
		charCounts[secondChar] += count
	}

	var maxCount, minCount int64
	first := true
	for _, count := range charCounts {
		if first {
			maxCount = count
			minCount = count
			first = false
		} else {
			if count > maxCount {
				maxCount = count
			}
			if count < minCount {
				minCount = count
			}
		}
	}

	return maxCount - minCount, nil
}
