package day7

import (
	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func countContainedBags(
	bagType string,
	rulesMap map[string][]struct {
		Count int
		Type  string
	},
	memo map[string]int,
) int {
	if cached, ok := memo[bagType]; ok {
		return cached
	}

	memo[bagType] = 0

	var result int
	if contents, ok := rulesMap[bagType]; ok {
		for _, content := range contents {
			result += content.Count * (1 + countContainedBags(content.Type, rulesMap, memo))
		}
	}

	memo[bagType] = result
	return result
}

func B(input *challenge.Input) (int, error) {
	rules, err := parseRules(input.LinesSlice())
	if err != nil {
		return 0, err
	}

	rulesMap := make(map[string][]struct {
		Count int
		Type  string
	})
	for _, rule := range rules {
		rulesMap[rule.BagType] = rule.Contents
	}

	memo := make(map[string]int)
	return countContainedBags("shiny gold", rulesMap, memo), nil
}
