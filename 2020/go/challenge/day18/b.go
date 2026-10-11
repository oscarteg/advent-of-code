package day18

import (
	"fmt"

	"github.com/oscarteg/advent-of-code-2020/challenge"
)

func B(input *challenge.Input) (int, error) {
	total := 0
	found := false
	for line := range input.Lines() {
		found = true
		value, err := evalExpression(line, true)
		if err != nil {
			return 0, err
		}
		total += value
	}
	if !found {
		return 0, fmt.Errorf("no expressions found in input")
	}
	return total, nil
}
