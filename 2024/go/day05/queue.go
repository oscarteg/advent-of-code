package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Rule struct {
	from int
	to   int
}

type Update struct {
	pages []int
}

func solve(input string) (int, int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")

	var rules []Rule
	var updates []Update
	parsingRules := true

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			parsingRules = false
			continue
		}

		if parsingRules {
			parts := strings.Split(line, "|")
			if len(parts) != 2 {
				return 0, 0, fmt.Errorf("invalid rule format: %s", line)
			}
			from, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return 0, 0, fmt.Errorf("invalid page number in rule: %w", err)
			}
			to, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return 0, 0, fmt.Errorf("invalid page number in rule: %w", err)
			}
			rules = append(rules, Rule{from, to})
		} else {
			var pages []int
			for _, pageStr := range strings.Split(line, ",") {
				page, err := strconv.Atoi(strings.TrimSpace(pageStr))
				if err != nil {
					return 0, 0, fmt.Errorf("invalid page number in update: %w", err)
				}
				pages = append(pages, page)
			}
			if len(pages) > 0 {
				updates = append(updates, Update{pages})
			}
		}
	}

	part1 := 0
	part2 := 0

	for _, update := range updates {
		if isValidOrder(update, rules) {
			part1 += getMiddle(update.pages)
		} else {
			sorted := sortUpdate(update, rules)
			part2 += getMiddle(sorted)
		}
	}

	return part1, part2, nil
}

func isValidOrder(update Update, rules []Rule) bool {
	positions := make(map[int]int)
	for i, page := range update.pages {
		positions[page] = i
	}

	for _, rule := range rules {
		fromPos, fromOk := positions[rule.from]
		toPos, toOk := positions[rule.to]

		if !fromOk || !toOk {
			continue
		}
		if fromPos > toPos {
			return false
		}
	}
	return true
}

func sortUpdate(update Update, rules []Rule) []int {
	pageSet := make(map[int]bool)
	for _, page := range update.pages {
		pageSet[page] = true
	}

	adj := make(map[int][]int)
	inDegree := make(map[int]int)

	for _, page := range update.pages {
		adj[page] = []int{}
		inDegree[page] = 0
	}

	for _, rule := range rules {
		if pageSet[rule.from] && pageSet[rule.to] {
			adj[rule.from] = append(adj[rule.from], rule.to)
		}
	}

	for _, neighbors := range adj {
		for _, neighbor := range neighbors {
			inDegree[neighbor]++
		}
	}

	var queue []int
	for page, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, page)
		}
	}

	var sorted []int
	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]
		sorted = append(sorted, current)

		for _, neighbor := range adj[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	return sorted
}

func getMiddle(pages []int) int {
	if len(pages) == 0 {
		return 0
	}
	return pages[len(pages)/2]
}
