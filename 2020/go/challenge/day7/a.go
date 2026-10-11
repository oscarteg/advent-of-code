package day7

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

type Rule struct {
	BagType  string
	Contents []struct {
		Count int
		Type  string
	}
}

func parseRules(lines []string) ([]Rule, error) {
	var rules []Rule
	bagRegex := regexp.MustCompile(`^(\w+ \w+) bags contain (.+)$`)
	contentRegex := regexp.MustCompile(`(\d+) (\w+ \w+) bags?`)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := bagRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			return nil, fmt.Errorf("failed to parse bag rule: %s", line)
		}

		bagType := matches[1]
		contentsStr := matches[2]

		var contents []struct {
			Count int
			Type  string
		}

		if contentsStr != "no other bags." {
			for _, match := range contentRegex.FindAllStringSubmatch(contentsStr, -1) {
				count, err := strconv.Atoi(match[1])
				if err != nil {
					return nil, fmt.Errorf("failed to parse count: %s", match[1])
				}
				contents = append(contents, struct {
					Count int
					Type  string
				}{count, match[2]})
			}
		}

		rules = append(rules, Rule{bagType, contents})
	}

	if len(rules) == 0 {
		return nil, fmt.Errorf("no rules parsed")
	}

	return rules, nil
}

func buildReverseIndex(rules []Rule) map[string]map[string]bool {
	canContain := make(map[string]map[string]bool)

	for _, rule := range rules {
		for _, content := range rule.Contents {
			if canContain[content.Type] == nil {
				canContain[content.Type] = make(map[string]bool)
			}
			canContain[content.Type][rule.BagType] = true
		}
	}

	return canContain
}

func countContainers(bagType string, canContain map[string]map[string]bool) int {
	visited := make(map[string]bool)
	queue := []string{bagType}
	visited[bagType] = true

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		if containers, ok := canContain[current]; ok {
			for container := range containers {
				if !visited[container] {
					visited[container] = true
					queue = append(queue, container)
				}
			}
		}
	}

	return len(visited) - 1
}

func A(input *challenge.Input) (int, error) {
	rules, err := parseRules(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	canContain := buildReverseIndex(rules)
	return countContainers("shiny gold", canContain), nil
}
